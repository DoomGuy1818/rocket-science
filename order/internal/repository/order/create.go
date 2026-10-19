package order

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"

	domainModel "github.com/DoomGuy1818/rocket-science/order/internal/model"
)

func (r *repository) Create(ctx context.Context, userID uuid.UUID, parts []domainModel.Part, totalPrice float64) (
	uuid.UUID,
	float64,
	error,
) {
	var queryBuilder strings.Builder

	queryBuilder.WriteString(
		"INSERT INTO order_parts (id, order_id, part_id, price) VALUES ",
	)

	args := make([]any, 0, len(parts)*4)

	for i, part := range parts {
		if i > 0 {
			queryBuilder.WriteString(", ")
		}

		offset := i * 4

		fmt.Fprintf(
			&queryBuilder,
			"($%d, $%d, $%d, $%d)",
			offset+1, offset+2, offset+3, offset+4,
		)

		args = append(args, uuid.New(), uuid.New(), part.UUID, part.Price)
	}

	if len(parts) > 0 {
		_, err := r.orders.Exec(ctx, queryBuilder.String(), args...)
		if err != nil {
			log.Printf("failed to insert order parts: %v", err)
			return uuid.Nil, 0, domainModel.ErrInternalServerError
		}
	}

	return userID, totalPrice, nil
}
