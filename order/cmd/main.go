package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/DoomGuy1818/rocket-science/order/internal/app"
	"github.com/DoomGuy1818/rocket-science/order/internal/config"
	"github.com/DoomGuy1818/rocket-science/platform/closer"
	"github.com/DoomGuy1818/rocket-science/platform/pkg/logger"
)

const configPath = "./deploy/compose/order/.env"

const (
	InventoryServiceAddress = "localhost:50051"
	PaymentServiceAddress   = "localhost:50052"
	HttpPort                = "8080"
	readHeaderTimeout       = 5 * time.Second
	shutdownTimeout         = 10 * time.Second
)

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, cancelCtx := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer gracefulShutdown()
	defer cancelCtx()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось запустить приложение", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
