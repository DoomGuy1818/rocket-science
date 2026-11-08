package app

import (
	"context"
	"fmt"
	"os"
	"platform/closer"
	"platform/migrator"
	"platform/pkg/logger"

	orderV1API "github.com/DoomGuy1818/rocket-science/order/internal/api/order/v1"
	inventoryV1Client "github.com/DoomGuy1818/rocket-science/order/internal/client/grpc/inventory/v1"
	paymentV1Client "github.com/DoomGuy1818/rocket-science/order/internal/client/grpc/payment/v1"
	"github.com/DoomGuy1818/rocket-science/order/internal/config"
	"github.com/DoomGuy1818/rocket-science/order/internal/repository"
	orderRepo "github.com/DoomGuy1818/rocket-science/order/internal/repository/order"
	"github.com/DoomGuy1818/rocket-science/order/internal/service"
	"github.com/DoomGuy1818/rocket-science/order/internal/service/order"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc/credentials/insecure"

	orderV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/payment/v1"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

type diContainer struct {
	orderV1API      orderV1.Handler
	orderService    service.Service
	orderRepository repository.Repository
	migrator        *migrator.Migrator

	paymentClient   *paymentV1Client.Client
	inventoryClient *inventoryV1Client.Client

	pgConn              *pgx.Conn
	paymentClientConn   *grpc.ClientConn
	inventoryClientConn *grpc.ClientConn
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1API.NewOrderHandler(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.Service {
	if d.orderService == nil {
		d.orderService = order.NewService(d.OrderRepository(ctx), d.InventoryClient(ctx), d.PaymentClient(ctx))
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.Repository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepo.NewRepository(d.NewPgConn(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) NewPgConn(ctx context.Context) *pgx.Conn {
	if d.pgConn == nil {
		conn, err := pgx.Connect(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("failed to connect to database: %v\n", err))
			os.Exit(1)
		}

		err = conn.Ping(ctx)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("failed to ping database: %v\n", err))
			os.Exit(1)
		}

		closer.AddNamed(
			"PG client", func(ctx context.Context) error {
				return conn.Close(ctx)
			},
		)

		d.pgConn = conn
	}

	return d.pgConn
}

func (d *diContainer) PaymentClient(ctx context.Context) *paymentV1Client.Client {
	if d.paymentClient == nil {
		conn := d.PaymentClientConn(ctx)

		d.paymentClient = paymentV1Client.NewClient(
			paymentV1.NewPaymentServiceClient(conn),
		)
	}

	return d.paymentClient
}

func (d *diContainer) InventoryClient(ctx context.Context) *inventoryV1Client.Client {
	if d.inventoryClient == nil {
		conn := d.InventoryClientConn(ctx)

		d.inventoryClient = inventoryV1Client.NewClient(
			inventoryV1.NewInventoryServiceClient(conn),
		)
	}

	return d.inventoryClient
}

func (d *diContainer) PaymentClientConn(ctx context.Context) *grpc.ClientConn {
	if d.paymentClientConn == nil {
		paymentConn, err := grpc.NewClient(
			config.AppConfig().PaymentGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("cannot connect to payment service: %v", err))
		}

		closer.AddNamed(
			"Payment client", func(ctx context.Context) error {
				return paymentConn.Close()
			},
		)

		d.paymentClientConn = paymentConn
	}

	return d.paymentClientConn
}

func (d *diContainer) InventoryClientConn(ctx context.Context) *grpc.ClientConn {
	if d.inventoryClientConn == nil {
		inventoryConn, err := grpc.NewClient(
			config.AppConfig().InventoryGRPC.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("cannot connect to payment service: %v", err))
		}

		closer.AddNamed(
			"inventory client", func(ctx context.Context) error {
				return inventoryConn.Close()
			},
		)

		d.inventoryClientConn = inventoryConn
	}

	return d.inventoryClientConn
}

func (d *diContainer) Migrator(ctx context.Context) *migrator.Migrator {
	if d.migrator == nil {
		d.migrator = migrator.NewMigrator(
			stdlib.OpenDB(*d.NewPgConn(ctx).Config().Copy()),
			config.AppConfig().Migrator.Path(),
		)
	}

	return d.migrator
}
