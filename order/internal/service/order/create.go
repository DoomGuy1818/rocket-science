package order

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
)

func (s *service) Create(ctx context.Context, userID uuid.UUID, partIDs []uuid.UUID) (uuid.UUID, float64, error) {
	var totalPrice float64

	parts, err := s.inventoryClient.ListParts(ctx, partIDs)
	if err != nil {
		return uuid.Nil, 0, err
	}

	log.Println(parts)
	log.Println(len(partIDs))

	if len(parts) != len(partIDs) {
		return uuid.Nil, 0, model.ErrCannotProcessOrder
	}

	for _, part := range parts {
		totalPrice += part.Price
	}

	orderID, price, err := s.orderStorage.Create(ctx, userID, partIDs, totalPrice)
	if err != nil {
		return uuid.Nil, 0, err
	}

	return orderID, price, nil
}
