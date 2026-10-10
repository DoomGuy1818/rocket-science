package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
)

type InventoryRepository interface {
	Get(ctx context.Context, uuid uuid.UUID) (model.Part, error)
	List(ctx context.Context, filterParts *model.PartsFilter) []model.Part
}
