package inventory

import "github.com/DoomGuy1818/rocket-science/inventory/internal/repository"

type InventoryService struct {
	inventoryRepository repository.InventoryRepository
}

func NewService(inventoryRepository repository.InventoryRepository) *InventoryService {
	return &InventoryService{
		inventoryRepository: inventoryRepository,
	}
}
