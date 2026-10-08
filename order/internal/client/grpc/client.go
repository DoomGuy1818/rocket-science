package grpc

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, ids []uuid.UUID) ([]model.Part, error)
}

type PaymentClient interface {
	Pay(ctx context.Context, orderID, userID uuid.UUID, paymentMethod string) (
		uuid.UUID,
		error,
	)
}
