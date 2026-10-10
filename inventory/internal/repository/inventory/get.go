package inventory

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
	modelConverter "github.com/DoomGuy1818/rocket-science/inventory/internal/repository/converter"
)

func (i *repository) Get(_ context.Context, uuid uuid.UUID) (model.Part, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	resp, ok := i.parts[uuid]
	if !ok {
		return model.Part{}, model.ErrNotFound
	}

	return modelConverter.PartToModel(*resp), nil
}
