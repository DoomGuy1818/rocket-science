package inventory

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
)

func (i *InventoryService) GetByID(ctx context.Context, uuid uuid.UUID) (model.Part, error) {
	part, err := i.inventoryRepository.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return model.Part{}, err
		}
		return model.Part{}, err
	}

	return part, nil
}
