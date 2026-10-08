package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/repository/model"
)

func (s *repository) Create(_ context.Context, userID uuid.UUID, parts []uuid.UUID, totalPrice float64) (
	uuid.UUID,
	float64,
	error,
) {
	orderID := uuid.New()

	order := &model.Order{
		OrderID:    orderID,
		UserID:     userID,
		Parts:      parts,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
	}

	s.orders[orderID] = order

	return orderID, order.TotalPrice, nil
}
