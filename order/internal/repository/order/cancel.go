package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	repoModel "github.com/DoomGuy1818/rocket-science/order/internal/repository/model"
)

func (s *repository) Cancel(_ context.Context, orderID uuid.UUID) error {
	order, ok := s.orders[orderID]
	if !ok {
		return model.ErrOrderNotFound
	}

	updatedOrder := &repoModel.Order{
		OrderID:       order.OrderID,
		UserID:        order.UserID,
		Parts:         copyPartUUIDs(order.Parts),
		TotalPrice:    order.TotalPrice,
		TransactionID: order.TransactionID,
		PaymentMethod: order.PaymentMethod,
		Status:        repoModel.StatusCanceled,
	}

	s.orders[orderID] = updatedOrder

	return nil
}
