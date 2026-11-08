package e2e

import (
	"context"

	"platform/pkg/logger"

	"go.uber.org/zap"
)

func teardownTestEnvironment(ctx context.Context, environment *TestEnvironment) {
	log := logger.Logger()
	log.Info(ctx, "🧹 Очистка тестового окружения...")

	cleanupTestEnvironment(ctx, environment)
}

func cleanupTestEnvironment(ctx context.Context, environment *TestEnvironment) {
	if environment.App != nil {
		if err := environment.App.Terminate(ctx); err != nil {
			logger.Error(ctx, "не удалось остановить контейнер приложения", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Контейнер приложения остановлен")
		}
	}

	if environment.Network != nil {
		if err := environment.Network.Remove(ctx); err != nil {
			logger.Error(ctx, "не удалось удалить сеть", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Сеть удалена")
		}
	}
}
