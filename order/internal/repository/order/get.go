package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	repoModel "github.com/DoomGuy1818/rocket-science/order/internal/repository/model"
)

func (s *repository) Get(_ context.Context, orderID uuid.UUID) (repoModel.Order, error) {
	order, ok := s.orders[orderID]
	if !ok {
		return repoModel.Order{}, model.ErrOrderNotFound
	}

	return *order, nil
}
