package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	"github.com/DoomGuy1818/rocket-science/order/internal/repository/conventer"
)

func (s *service) Get(ctx context.Context, orderID uuid.UUID) (model.Order, error) {
	order, err := s.orderStorage.Get(ctx, orderID)
	if err != nil {
		return model.Order{}, err
	}

	return conventer.OrderToDomainModel(order), nil
}
