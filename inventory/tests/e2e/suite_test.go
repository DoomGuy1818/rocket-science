package e2e

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/DoomGuy1818/rocket-science/platform/pkg/logger"
)

const testsTimeout = 5 * time.Minute

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(
	func() {
		err := logger.Init(loggerLevelValue, true)
		if err != nil {
			panic(fmt.Sprintf("Не удалось инициализировать логгер: %v", err))
		}

		suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

		envVars, err := godotenv.Read(filepath.Join("..", "..", "deploy", "compose", "payment"))
		if err != nil {
			logger.Fatal(suiteCtx, "Не удалось загрузить .env файл", zap.Error(err))
		}

		for key, value := range envVars {
			_ = os.Setenv(key, value)
		}

		logger.Info(suiteCtx, "Запуск тестового окружения...")
		env = setupTestEnvironment(suiteCtx)
	},
)

var _ = AfterSuite(
	func() {
		logger.Info(context.Background(), "Завершение набора тестов")
		if env != nil {
			teardownTestEnvironment(suiteCtx, env)
		}
		suiteCancel()
	},
)
