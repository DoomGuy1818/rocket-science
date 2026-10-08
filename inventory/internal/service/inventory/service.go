package inventory

import "github.com/DoomGuy1818/rocket-science/inventory/internal/repository"

type service struct {
	inventoryRepository repository.InventoryRepository
}

func NewService(inventoryRepository repository.InventoryRepository) *service {
	return &service{
		inventoryRepository: inventoryRepository,
	}
}
