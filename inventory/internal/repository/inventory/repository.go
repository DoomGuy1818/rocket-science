package inventory

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type repository struct {
	mu    sync.RWMutex
	parts *mongo.Collection
}

func NewRepository(ctx context.Context, db *mongo.Database) *repository {
	collection := db.Collection("inventory")

	indexModel := []mongo.IndexModel{
		{
			Keys: bson.D{
				{
					Key: "name", Value: 1,
				},
			},
			Options: options.Index().SetUnique(false),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModel)
	if err != nil {
		panic(err)
	}

	return &repository{
		parts: collection,
	}
}
