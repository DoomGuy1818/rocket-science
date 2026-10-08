package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	repoModel "github.com/DoomGuy1818/rocket-science/order/internal/repository/model"
)

func (s *ApiSuite) TestGetOrder() {
	ctx := context.Background()
	id := uuid.New()

	userID := uuid.New()
	partIDs := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}
	transactionID := uuid.New()

	repoOrder := repoModel.Order{
		OrderID:       id,
		UserID:        userID,
		Parts:         partIDs,
		TransactionID: transactionID,
		PaymentMethod: repoModel.PaymentMethodCard,
		Status:        repoModel.StatusPendingPayment,
	}

	expectedOrder := model.Order{
		OrderID:       id,
		UserID:        userID,
		PartIDs:       partIDs,
		TransactionID: transactionID,
		PaymentMethod: model.PaymentMethodCard,
		Status:        model.StatusPendingPayment,
	}

	tests := []struct {
		name      string
		ID        uuid.UUID
		repoOrder repoModel.Order
		expected  model.Order
		expErr    error
	}{
		{
			name:      "Successful Get order",
			ID:        id,
			repoOrder: repoOrder,
			expected:  expectedOrder,
			expErr:    nil,
		},
		{
			name:      "Failed to get order by Not Found",
			ID:        id,
			repoOrder: repoModel.Order{},
			expected:  model.Order{},
			expErr:    model.ErrOrderNotFound,
		},
	}

	for _, tt := range tests {
		s.Run(
			tt.name, func() {
				s.repo.
					On("Get", ctx, tt.ID).
					Return(tt.repoOrder, tt.expErr).
					Once()

				actualOrder, err := s.service.Get(ctx, tt.ID)

				s.Require().ErrorIs(err, tt.expErr)

				if tt.expErr == nil {
					s.Require().Equal(tt.expected, actualOrder)
				} else {
					s.Require().Equal(model.Order{}, actualOrder)
				}

				s.repo.AssertExpectations(s.T())
			},
		)
	}
}
