package order

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) Cancel(ctx context.Context, orderID uuid.UUID) error {
	if err := s.orderStorage.Cancel(ctx, orderID); err != nil {
		return err
	}

	return nil
}
