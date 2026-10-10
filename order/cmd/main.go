package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/DoomGuy1818/rocket-science/order/internal/api/order/v1"
	invClient "github.com/DoomGuy1818/rocket-science/order/internal/client/grpc/inventory/v1"
	payClient "github.com/DoomGuy1818/rocket-science/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/DoomGuy1818/rocket-science/order/internal/repository/order"
	orderService "github.com/DoomGuy1818/rocket-science/order/internal/service/order"
	orderV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/payment/v1"
)

const (
	InventoryServiceAddress = "localhost:50051"
	PaymentServiceAddress   = "localhost:50052"
	HttpPort                = "8080"
	readHeaderTimeout       = 5 * time.Second
	shutdownTimeout         = 10 * time.Second
)

func main() {
	storage := orderRepo.NewStorage()

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

	inventoryClient := invClient.NewClient(inventoryV1.NewInventoryServiceClient(inventoryConn))

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

	paymentClient := payClient.NewClient(paymentV1.NewPaymentServiceClient(paymentConn))

	service := orderService.NewService(storage, inventoryClient, paymentClient)

	orderHandler := v1.NewOrderHandler(service)

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
