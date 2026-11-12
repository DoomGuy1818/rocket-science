package order

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
)

func (r *repository) Pay(ctx context.Context, orderID, transactionID uuid.UUID, payMethod string) error {
	res, err := r.orders.Exec(
		ctx,
		"UPDATE orders SET payment_method=$1, transaction_id=$2 WHERE id=$3",
		payMethod,
		transactionID,
		orderID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrOrderNotFound
		}

		return err
	}

	log.Printf("ipdated %d rows\n", res.RowsAffected())
	return nil
}
