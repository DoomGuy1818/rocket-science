package order

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	repoModel "github.com/DoomGuy1818/rocket-science/order/internal/repository/model"
)

func (s *ApiSuite) TestPaymentService() {
	ctx := context.Background()

	t := s.T()

	t.Run(
		"Successful payed order", func(t *testing.T) {
			orderID := uuid.New()
			userID := uuid.New()
			transactionID := uuid.New()
			paymentMethod := "card"

			s.repo.
				On("Get", ctx, orderID).
				Return(
					repoModel.Order{
						OrderID: orderID,
					}, nil,
				).
				Once()

			s.paymentClient.
				On(
					"Pay",
					ctx,
					orderID,
					userID,
					paymentMethod,
				).
				Return(transactionID, nil).
				Once()

			s.repo.
				On(
					"Pay",
					ctx,
					orderID,
					transactionID,
					paymentMethod,
				).
				Return(nil).
				Once()

			actualTransactionID, err := s.service.PayOrder(
				ctx,
				orderID,
				userID,
				paymentMethod,
			)

			require.NoError(t, err)
			require.Equal(t, transactionID, actualTransactionID)
		},
	)

	t.Run(
		"Failed payed order by order not found", func(t *testing.T) {
			orderID := uuid.New()
			userID := uuid.New()
			paymentMethod := "card"

			s.repo.
				On("Get", ctx, orderID).
				Return(repoModel.Order{}, model.ErrOrderNotFound).
				Once()

			actualTransactionID, err := s.service.PayOrder(
				ctx,
				orderID,
				userID,
				paymentMethod,
			)

			require.ErrorIs(t, err, model.ErrOrderNotFound)
			require.Equal(t, uuid.Nil, actualTransactionID)
		},
	)

	t.Run(
		"Failed payed order by payment service", func(t *testing.T) {
			orderID := uuid.New()
			userID := uuid.New()
			paymentMethod := "card"

			s.repo.
				On("Get", ctx, orderID).
				Return(
					repoModel.Order{
						OrderID: orderID,
					}, nil,
				).
				Once()

			s.paymentClient.
				On(
					"Pay",
					ctx,
					orderID,
					userID,
					paymentMethod,
				).
				Return(uuid.Nil, model.ErrInternalServerError).
				Once()

			actualTransactionID, err := s.service.PayOrder(
				ctx,
				orderID,
				userID,
				paymentMethod,
			)

			require.ErrorIs(t, err, model.ErrInternalServerError)
			require.Equal(t, uuid.Nil, actualTransactionID)
		},
	)

	t.Run(
		"Failed payed order by storage", func(t *testing.T) {
			orderID := uuid.New()
			userID := uuid.New()
			transactionID := uuid.New()
			paymentMethod := "card"

			s.repo.
				On("Get", ctx, orderID).
				Return(
					repoModel.Order{
						OrderID: orderID,
					}, nil,
				).
				Once()

			s.paymentClient.
				On(
					"Pay",
					ctx,
					orderID,
					userID,
					paymentMethod,
				).
				Return(transactionID, nil).
				Once()

			s.repo.
				On(
					"Pay",
					ctx,
					orderID,
					transactionID,
					paymentMethod,
				).
				Return(model.ErrInternalServerError).
				Once()

			actualTransactionID, err := s.service.PayOrder(
				ctx,
				orderID,
				userID,
				paymentMethod,
			)

			require.ErrorIs(t, err, model.ErrInternalServerError)
			require.Equal(t, uuid.Nil, actualTransactionID)
		},
	)
}
