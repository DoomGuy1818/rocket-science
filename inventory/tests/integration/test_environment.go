package integration

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := "3f2b8c1e-7a4d-4e9b-9c3a-5d6e8f1a2b47"
	now := time.Now()

	partDoc := bson.M{
		"_id":            partUUID,
		"name":           "Ионный двигатель X-200",
		"description":    "Маломощный ионный двигатель для орбитальных манёвров",
		"price":          1250000.00,
		"stock_quantity": 5,
		"category":       1,
		"dimensions": bson.M{
			"length": 220.5,
			"width":  90.0,
			"height": 90.0,
			"weight": 340.0,
		},
		"manufacturer": bson.M{
			"name":    "AeroDyne",
			"country": "USA",
			"website": "",
		},
		"tags": []string{
			"ion",
			"orbital",
			"low-thrust",
		},
		"metadata": bson.M{
			"cert_number": bson.M{
				"string_value": "ENG-2024-0042",
			},
		},
		"created_at": bson.NewDateTimeFromTime(
			time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC),
		),
		"updated_at": bson.NewDateTimeFromTime(now),
	}

	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "parts-service"
	}

	_, err := env.Mongo.Client().
		Database(databaseName).
		Collection(orderCollectionName).
		InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

func (env *TestEnvironment) ClearInventoryCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(orderCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
