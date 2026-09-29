package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	orderV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/payment/v1"
)

const (
	InventoryServiceAddress = "localhost:50051"
	PaymentServiceAddress   = "localhost:50052"
	HttpPort                = "8080"
	readHeaderTimeout       = 5 * time.Second
	shutdownTimeout         = 10 * time.Second
)

type OrderStorage struct {
	mu      sync.RWMutex
	storage map[string]*orderV1.Order
}

type PartsFilter struct {
	partUuids    []string
	names        []string
	categories   []inventory_v1.Category
	manufacturer []string
	Tags         []string
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		storage: make(map[string]*orderV1.Order),
	}
}

func (o *OrderStorage) GetOrder(uuid uuid.UUID) (*orderV1.Order, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	order, ok := o.storage[uuid.String()]
	if !ok {
		return nil, fmt.Errorf("order not found")
	}

	return order, nil
}

func (o *OrderStorage) CreateOrder(order *orderV1.Order) (uuid.UUID, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.storage[order.OrderUUID.String()] = order

	return order.OrderUUID, nil
}

func (o *OrderStorage) UpdateOrderTransaction(order *orderV1.Order) string {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.storage[order.OrderUUID.String()] = order

	return order.TransactionUUID.String()
}

func (o *OrderStorage) CancelOrder(order *orderV1.Order) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.storage[order.OrderUUID.String()] = order
}

type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventory_v1.InventoryServiceClient
	paymentClient   payment_v1.PaymentServiceClient
}

func NewOrderHandler(
	storage *OrderStorage,
	invClient inventory_v1.InventoryServiceClient,
	payClient payment_v1.PaymentServiceClient,
) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: invClient,
		paymentClient:   payClient,
	}
}

func (h *OrderHandler) CancelOrderByID(
	_ context.Context,
	params orderV1.CancelOrderByIDParams,
) (orderV1.CancelOrderByIDRes, error) {
	order, err := h.storage.GetOrder(params.OrderID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Cannot get order with id" + params.OrderID.String(),
		}, nil
	}

	h.storage.CancelOrder(
		&orderV1.Order{
			OrderUUID:       order.GetOrderUUID(),
			UserUUID:        order.GetUserUUID(),
			PartUuids:       order.GetPartUuids(),
			TotalPrice:      order.GetTotalPrice(),
			TransactionUUID: order.GetTransactionUUID(),
			PaymentMethod:   order.GetPaymentMethod(),
			Status:          orderV1.PaymentStatusCANCELLED,
		},
	)

	return &orderV1.CancelOrderByIDNoContent{}, nil
}

func (h *OrderHandler) CreateOrderByID(
	ctx context.Context,
	req *orderV1.OrderCreateRequest,
) (orderV1.CreateOrderByIDRes, error) {
	partsUUID := req.GetPartUuids()
	if len(partsUUID) == 0 {
		return &orderV1.UnprocessableEntityError{
			Code:    http.StatusBadRequest,
			Message: "part_uuids must not be empty",
		}, nil
	}

	uuidStrings := make([]string, len(partsUUID))
	for i, u := range partsUUID {
		uuidStrings[i] = u.String()
	}

	invResp, err := GetPartsByFilters(
		ctx, PartsFilter{partUuids: uuidStrings}, h.inventoryClient,
	)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "failed to check parts availability",
		}, nil
	}

	parts := invResp

	if len(parts) < len(partsUUID) {
		return &orderV1.UnprocessableEntityError{
			Code:    http.StatusUnprocessableEntity,
			Message: "one or more parts do not exist",
		}, nil
	}

	totalPrice := float32(0)

	for _, part := range parts {
		totalPrice += float32(part.GetPrice())
	}

	order := orderV1.Order{
		OrderUUID:     uuid.New(),
		UserUUID:      req.GetUserUUID(),
		PartUuids:     req.GetPartUuids(),
		TotalPrice:    totalPrice,
		PaymentMethod: orderV1.PaymentMethodPAYMENTMETHODUNKNOWN,
		Status:        orderV1.PaymentStatusPENDINGPAYMENT,
	}

	orderID, err := h.storage.CreateOrder(&order)
	if err != nil {
		return &orderV1.UnprocessableEntityError{
			Code:    http.StatusUnprocessableEntity,
			Message: "one or more parts do not exist",
		}, nil
	}

	return &orderV1.OrderCreateResponse{
		OrderUUID: orderID, TotalPrice: totalPrice,
	}, nil
}

func (h *OrderHandler) CreateOrderPaymentByID(
	ctx context.Context,
	req *orderV1.OrderPaymentRequest,
	params orderV1.CreateOrderPaymentByIDParams,
) (orderV1.CreateOrderPaymentByIDRes, error) {
	order, err := h.storage.GetOrder(params.OrderID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "Cannot get order with id" + params.OrderID.String(),
		}, nil
	}

	payServiceResp, err := PayOrderByOrderID(
		ctx,
		h.paymentClient,
		order.OrderUUID.String(),
		order.UserUUID.String(),
		string(req.GetPaymentMethod()),
	)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Cannot get response from payService",
		}, nil
	}

	transactionUUID, err := uuid.Parse(payServiceResp)
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Cannot parse transaction uuid",
		}, nil
	}

	h.storage.UpdateOrderTransaction(
		&orderV1.Order{
			OrderUUID:       order.OrderUUID,
			UserUUID:        order.UserUUID,
			PartUuids:       order.PartUuids,
			TotalPrice:      order.TotalPrice,
			TransactionUUID: transactionUUID,
			PaymentMethod:   req.GetPaymentMethod(),
			Status:          orderV1.PaymentStatusPAID,
		},
	)

	return &orderV1.OrderPaymentResponse{TransactionUUID: transactionUUID}, nil
}

func (h *OrderHandler) GetOrderByID(_ context.Context, params orderV1.GetOrderByIDParams) (
	orderV1.GetOrderByIDRes,
	error,
) {
	order, err := h.storage.GetOrder(params.OrderID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "cannot get order with id" + params.OrderID.String(),
		}, nil
	}

	return &orderV1.Order{
		OrderUUID:       order.GetOrderUUID(),
		UserUUID:        order.GetUserUUID(),
		PartUuids:       order.GetPartUuids(),
		TotalPrice:      order.GetTotalPrice(),
		TransactionUUID: order.GetTransactionUUID(),
		PaymentMethod:   order.GetPaymentMethod(),
		Status:          order.GetStatus(),
	}, nil
}

func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func PayOrderByOrderID(
	ctx context.Context,
	client payment_v1.PaymentServiceClient,
	orderID string,
	userID string,
	paymentMethod string,
) (string, error) {
	val, ok := payment_v1.PaymentMethod_value[paymentMethod]
	if !ok {
		return "", status.Errorf(codes.InvalidArgument, "Sended not valid payment method: %s", paymentMethod)
	}

	paymentOrder := payment_v1.PayOrderRequest{
		OrderId:       orderID,
		UserUuid:      userID,
		PaymentMethod: payment_v1.PaymentMethod(val),
	}

	resp, err := client.PayOrder(ctx, &paymentOrder)
	if err != nil {
		return "", err
	}

	return resp.GetTransactionUuid(), nil
}

func GetPartsByFilters(
	ctx context.Context,
	filters PartsFilter,
	client inventory_v1.InventoryServiceClient,
) (
	[]*inventory_v1.Part,
	error,
) {
	partFilters := &inventory_v1.PartFilters{
		Uuids:                 filters.partUuids,
		Names:                 filters.names,
		Categories:            filters.categories,
		ManufacturerCountries: filters.manufacturer,
		Tags:                  filters.Tags,
	}

	resp, err := client.ListParts(ctx, &inventory_v1.ListPartsRequest{Filter: partFilters})
	if err != nil {
		return nil, err
	}

	return resp.GetParts(), nil
}

func main() {
	storage := NewOrderStorage()

	inventoryConn, err := grpc.NewClient(
		InventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("cannot connect to inventory service: %v", err)
	}
	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close inventory connection: %v", cerr)
		}
	}()

	inventoryClient := inventory_v1.NewInventoryServiceClient(inventoryConn)

	paymentConn, err := grpc.NewClient(
		PaymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("cannot connect to payment service: %v", err)
	}
	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close payment connection: %v", cerr)
		}
	}()

	paymentClient := payment_v1.NewPaymentServiceClient(paymentConn)

	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Printf("cannot start server: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:        net.JoinHostPort("localhost", HttpPort),
		Handler:     r,
		ReadTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("starting server on port %s", HttpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("cannot start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Printf("cannot shutdown server: %v", err)
	}

	log.Println("shut down successfully")
}
