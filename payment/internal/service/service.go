package service

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	PayOrder(ctx context.Context, orderID, userID uuid.UUID, payMethod string) (uuid.UUID, error)
}
