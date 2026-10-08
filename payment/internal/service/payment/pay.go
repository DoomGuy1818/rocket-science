package payment

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) PayOrder(_ context.Context, _, _ uuid.UUID, _ string) (
	uuid.UUID,
	error,
) {
	return uuid.New(), nil
}
