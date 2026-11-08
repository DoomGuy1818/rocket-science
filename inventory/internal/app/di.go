package app

import (
	"context"
	"fmt"
	"platform/closer"

	inventoryV1API "github.com/DoomGuy1818/rocket-science/inventory/internal/api/inventory/v1"
	"github.com/DoomGuy1818/rocket-science/inventory/internal/config"
	"github.com/DoomGuy1818/rocket-science/inventory/internal/repository"
	inventoryRepository "github.com/DoomGuy1818/rocket-science/inventory/internal/repository/inventory"
	"github.com/DoomGuy1818/rocket-science/inventory/internal/service"
	inventoryService "github.com/DoomGuy1818/rocket-science/inventory/internal/service/inventory"
	inventoryV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type diContainer struct {
	inventoryV1API      inventoryV1.InventoryServiceServer
	inventoryService    service.Service
	inventoryRepository repository.InventoryRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) NewInventoryAPIV1(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.InventoryService(ctx))
	}
	return d.inventoryV1API
}

func (d *diContainer) InventoryService(ctx context.Context) service.Service {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.InventoryRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) InventoryRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = inventoryRepository.NewRepository(ctx, d.MongoDBHandle(ctx))
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		if err = client.Ping(ctx, nil); nil != err {
			panic(fmt.Sprintf("failed to ping MongoDB: %s\n", err.Error()))
		}

		closer.AddNamed(
			"MongoDB client", func(ctx context.Context) error {
				return client.Disconnect(ctx)
			},
		)

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}
