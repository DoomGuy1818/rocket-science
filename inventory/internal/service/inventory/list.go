package inventory

import (
	"context"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
)

func (i *service) ListByFilters(ctx context.Context, filters *model.PartsFilter) []model.Part {
	parts := i.inventoryRepository.List(ctx, filters)

	return parts
}
