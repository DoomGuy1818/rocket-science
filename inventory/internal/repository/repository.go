package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
	repoModel "github.com/DoomGuy1818/rocket-science/inventory/internal/repository/model"
)

type InventoryRepository interface {
	Get(ctx context.Context, uuid uuid.UUID) (model.Part, error)
	List(ctx context.Context, filterParts *repoModel.PartsFilter) ([]model.Part, error)
}
