package v1

import (
	"github.com/DoomGuy1818/rocket-science/inventory/internal/service"
	inventoryV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	inventoryService service.Service
}

func NewAPI(inventoryService service.Service) *api {
	return &api{
		inventoryService: inventoryService,
	}
}
