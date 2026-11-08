package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"platform/closer"
	"platform/migrator"
	"platform/pkg/logger"
	"time"

	inventoryV1Client "github.com/DoomGuy1818/rocket-science/order/internal/client/grpc/inventory/v1"
	paymentV1Client "github.com/DoomGuy1818/rocket-science/order/internal/client/grpc/payment/v1"
	"github.com/DoomGuy1818/rocket-science/order/internal/config"
	orderV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	di              *diContainer
	migrator        *migrator.Migrator
	serverApi       *orderV1.Server
	httpServer      *http.Server
	paymentClient   *paymentV1Client.Client
	inventoryClient *inventoryV1Client.Client
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runHttpServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initDI,
		a.initMigrator,
		a.initLogger,
		a.initCloser,
		a.initHttpServer,
		a.initPaymentGRPCClient,
		a.initInventoryGRPCClient,
	}

	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.di = NewDIContainer()

	return nil
}

func (a *App) initMigrator(ctx context.Context) error {
	a.migrator = a.di.Migrator(ctx)

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())

	return nil
}

func (a *App) initHttpServer(_ context.Context) error {
	httpServer, err := orderV1.NewServer(a.di.orderV1API)
	if err != nil {
		return err
	}

	a.serverApi = httpServer

	return nil
}

func (a *App) initPaymentGRPCClient(ctx context.Context) error {
	a.paymentClient = a.di.PaymentClient(ctx)

	return nil
}

func (a *App) initInventoryGRPCClient(ctx context.Context) error {
	a.inventoryClient = a.di.InventoryClient(ctx)

	closer.AddNamed(
		"http Server", func(ctx context.Context) error {

			if err := a.httpServer.Shutdown(ctx); err != nil {
				logger.Error(ctx, fmt.Sprintf("cannot shutdown server: %v", err))
				return err
			}

			return nil
		},
	)

	return nil
}

func (a *App) runHttpServer(ctx context.Context) error {
	logger.Info(
		ctx,
		fmt.Sprintf("🚀 HTTP OrderService server listening on %s", config.AppConfig().OrderHTTP.Address()),
	)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", a.serverApi)

	server := &http.Server{
		Addr:        config.AppConfig().OrderHTTP.Address(),
		Handler:     r,
		ReadTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	a.httpServer = server

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, fmt.Sprintf("cannot start server: %v", err))

		return err
	}

	return nil
}
