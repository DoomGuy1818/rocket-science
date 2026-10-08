package repository

import (
	"context"

	"github.com/google/uuid"

	repoModel "github.com/DoomGuy1818/rocket-science/order/internal/repository/model"
)

type OrderRepository interface {
	Create(ctx context.Context, userID uuid.UUID, parts []uuid.UUID, totalPrice float64) (uuid.UUID, float64, error)
	Cancel(ctx context.Context, orderID uuid.UUID) error
	Get(ctx context.Context, orderID uuid.UUID) (repoModel.Order, error)
	Pay(ctx context.Context, orderID, transactionID uuid.UUID, paymentMethod string) error
}
