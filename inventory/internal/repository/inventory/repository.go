package inventory

import (
	"sync"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/repository/model"
)

type Repository struct {
	mu    sync.RWMutex
	parts map[uuid.UUID]*model.Part
}

func NewRepository() *Repository {
	return &Repository{
		parts: InitRepository(),
	}
}
