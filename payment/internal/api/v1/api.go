package v1

import (
	"github.com/DoomGuy1818/rocket-science/payment/internal/service"
	payV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/payment/v1"
)

type api struct {
	payV1.UnimplementedPaymentServiceServer

	inventoryService service.Service
}

func NewAPI(inventoryService service.Service) *api {
	return &api{
		inventoryService: inventoryService,
	}
}
