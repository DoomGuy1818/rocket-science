package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, userID uuid.UUID, parts []uuid.UUID) (uuid.UUID, float64, error)
	PayOrder(
		ctx context.Context,
		orderID uuid.UUID,
		userID uuid.UUID,
		paymentMethod string,
	) (uuid.UUID, error)
	Cancel(ctx context.Context, orderID uuid.UUID) error
	Get(ctx context.Context, orderID uuid.UUID) (model.Order, error)
}
