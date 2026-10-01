package v1

import (
	inventoryV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
)

type InventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
}
