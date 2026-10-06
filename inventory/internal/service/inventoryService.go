package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
)

type InventoryService interface {
	GetByID(ctx context.Context, ID uuid.UUID) (model.Part, error)
	ListByFilters(ctx context.Context, filter *model.PartsFilter) []model.Part
}
