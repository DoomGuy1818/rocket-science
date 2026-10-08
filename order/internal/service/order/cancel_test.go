package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
)

func (s *ApiSuite) TestCancelingOrder() {
	testID := uuid.New()
	ctx := context.Background()

	tests := []struct {
		name string
		ID   uuid.UUID
		err  error
	}{
		{
			name: "Successful order cancellation",
			ID:   testID,
			err:  nil,
		},
		{
			name: "Failed order cancellation by not found error",
			ID:   testID,
			err:  model.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		s.repo.On("Cancel", ctx, tt.ID).Return(tt.err).Once()

		err := s.service.Cancel(ctx, tt.ID)

		s.Require().Equal(err, tt.err)
	}
}
