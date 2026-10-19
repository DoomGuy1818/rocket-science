package order

import (
	"context"
	"database/sql"
	"log"

	"github.com/go-faster/errors"
	"github.com/google/uuid"

	domainModel "github.com/DoomGuy1818/rocket-science/order/internal/model"
)

func (r *repository) Cancel(ctx context.Context, orderID uuid.UUID) error {
	res, err := r.orders.Exec(ctx, "UPDATE orders SET status = 'CANCELLED' WHERE id = $1", orderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domainModel.ErrOrderNotFound
		}
		log.Printf("failed to cancel order: %v", err)
		return err
	}

	log.Printf("cancelled %d rows\n", res.RowsAffected())
	return nil
}
