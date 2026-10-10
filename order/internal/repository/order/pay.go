package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	repoModel "github.com/DoomGuy1818/rocket-science/order/internal/repository/model"
)

func (s *repository) Pay(_ context.Context, orderID, transactionID uuid.UUID, payMethod string) error {
	order, ok := s.orders[orderID]
	if !ok {
		return model.ErrOrderNotFound
	}

	updatedOrder := &repoModel.Order{
		OrderID:       orderID,
		UserID:        order.UserID,
		Parts:         order.Parts,
		TotalPrice:    order.TotalPrice,
		TransactionID: transactionID,
		PaymentMethod: repoModel.PaymentMethod(payMethod),
		Status:        repoModel.StatusPaid,
	}

	s.orders[orderID] = updatedOrder

	return nil
}
