package order

import (
	"github.com/DoomGuy1818/rocket-science/order/internal/client/grpc"
	"github.com/DoomGuy1818/rocket-science/order/internal/repository"
)

type service struct {
	orderStorage    repository.OrderRepository
	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(
	orderStorage repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
) *service {
	return &service{
		orderStorage,
		inventoryClient,
		paymentClient,
	}
}
