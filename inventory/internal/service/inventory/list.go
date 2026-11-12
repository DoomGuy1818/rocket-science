package inventory

import (
	"context"
	"log"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
)

func (i *service) ListByFilters(ctx context.Context, filters *model.PartsFilter) ([]model.Part, error) {
	parts, err := i.inventoryRepository.List(ctx, filters)
	if err != nil {
		return []model.Part{}, err
	}

	for _, part := range parts {
		log.Printf(
			"Part: %s, Metadata: %#v",
			part.Name,
			part.Metadata,
		)
	}

	return parts, nil
}
