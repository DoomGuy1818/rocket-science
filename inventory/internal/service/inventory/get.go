package inventory

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
)

func (i *service) GetByID(ctx context.Context, uuid uuid.UUID) (model.Part, error) {
	part, err := i.inventoryRepository.Get(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
