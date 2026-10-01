package inventory

import (
	"sync"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/repository/model"
)

type Repository struct {
	mu    sync.RWMutex
	parts map[string]*model.Part
}

func NewRepository() *Repository {
	return &Repository{
		parts: make(map[string]*model.Part),
	}
}
