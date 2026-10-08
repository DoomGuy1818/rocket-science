package order

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) PayOrder(
	ctx context.Context,
	orderID uuid.UUID,
	userID uuid.UUID,
	paymentMethod string,
) (
	uuid.UUID,
	error,
) {
	_, err := s.Get(ctx, orderID)
	if err != nil {
		return uuid.Nil, err
	}

	transID, err := s.paymentClient.Pay(ctx, orderID, userID, paymentMethod)
	if err != nil {
		return uuid.Nil, err
	}

	err = s.orderStorage.Pay(ctx, orderID, transID, paymentMethod)
	if err != nil {
		return uuid.Nil, err
	}

	return transID, nil
}
