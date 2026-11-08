package e2e

import (
	"context"
	"os"
	"platform/pkg/testcontainers"
	"platform/pkg/testcontainers/mongo"
	"time"

	"platform/pkg/logger"
	"platform/pkg/testcontainers/app"
	"platform/pkg/testcontainers/network"
	"platform/pkg/testcontainers/path"

	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

const (
	orderAppName    = "order-app"
	orderDockerfile = "deploy/docker/order/Dockerfile"

	// Переменные окружения приложения
	grpcPortKey = "GRPC_PORT"

	loggerLevelValue = "debug"

	startupTimeout = 3 * time.Minute
)

type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "Подготовка тестого окружения")

	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "сеть успешно создана")

	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogging(ctx, testcontainers.MongoDatabaseKey)

	grpcPort := getEnvWithLogging(ctx, grpcPortKey)

	generatedMongo, err := mongo.NewContainer(
		ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер MongoDB", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер MongoDB успешно запущен")

	projectRoot := path.GetPathRoot()

	appEnv := map[string]string{
		// Переопределяем хост MongoDB для подключения к контейнеру из testcontainers
		testcontainers.MongoHostKey: generatedMongo.Config().ContainerName,
	}

	waitStrategy := wait.ForListeningPort(grpcPort + "/tcp").
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(
		ctx,
		app.WithName(orderAppName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, orderDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		logger.Fatal(ctx, "не удалось инициализировать контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "Контейнер приложения успешно запущен!")

	return &TestEnvironment{
		Network: generatedNetwork,
		App:     appContainer,
	}
}

func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
	}

	return value
}
