package inventory

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
)

func (i *repository) Get(ctx context.Context, uuid uuid.UUID) (model.Part, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	var part model.Part

	err := i.parts.FindOne(ctx, bson.M{"_id": uuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrNotFound
		}

		return model.Part{}, err
	}

	return part, nil
}
