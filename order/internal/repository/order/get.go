package order

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	repoModel "github.com/DoomGuy1818/rocket-science/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, orderID uuid.UUID) (repoModel.Order, error) {
	row := r.orders.QueryRow(
		ctx,
		`SELECT o.id, o.user_id, ARRAY_AGG(oi.part_id) AS parts, o.total_price, o.transaction_id, o.payment_method, o.status FROM orders o LEFT JOIN order_items oi ON oi.order_id = o.id WHERE o.id = $1 GROUP BY o.id, o.user_id, o.total_price, o.transaction_id, o.payment_method, o.status`,
		orderID,
	)

	var order repoModel.Order

	err := row.Scan(
		&order.OrderID,
		&order.UserID,
		&order.Parts,
		&order.TotalPrice,
		&order.TransactionID,
		&order.PaymentMethod,
		&order.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("failed to scan order: %v\n", err)
			return repoModel.Order{}, model.ErrOrderNotFound
		}
		log.Printf("failed to scan order: %v\n", err)
		return repoModel.Order{}, err
	}

	return order, nil
}
