package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
)

func (s *ApiSuite) TestCreateOrder() {
	ctx := context.Background()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{
		uuid.New(),
		uuid.New(),
	}

	tests := []struct {
		name          string
		userID        uuid.UUID
		partIDs       []uuid.UUID
		parts         []model.Part
		err           error
		expectedPrice float64
	}{
		{
			name:    "Successful creation order",
			userID:  userUUID,
			partIDs: partUUIDs,
			parts: []model.Part{
				{
					UUID:  partUUIDs[0],
					Price: 100,
				},
				{
					UUID:  partUUIDs[1],
					Price: 200,
				},
			},
			err:           nil,
			expectedPrice: 300,
		},
		{
			name:    "Failed creation order by unprocessable entity",
			userID:  userUUID,
			partIDs: partUUIDs, // 2 ID
			parts: []model.Part{
				{UUID: partUUIDs[0], Price: 100}, // вернули только 1
			},
			err:           model.ErrCannotProcessOrder,
			expectedPrice: 0,
		},
	}

	for _, tt := range tests {
		s.Run(
			tt.name, func() {
				s.inventoryClient.On("ListParts", ctx, tt.partIDs).Return(tt.parts, nil).Once()

				if tt.err != nil {
					_, _, err := s.service.Create(ctx, tt.userID, tt.partIDs)

					s.Require().ErrorIs(err, tt.err)
					return
				}

				orderID := uuid.New()

				s.repo.
					On(
						"Create",
						ctx,
						tt.userID,
						tt.partIDs,
						tt.expectedPrice,
					).
					Return(orderID, tt.expectedPrice, nil)

				actualID, actualPrice, err := s.service.Create(ctx, tt.userID, tt.partIDs)

				s.Require().NoError(err)
				s.Require().Equal(orderID, actualID)
				s.Require().Equal(tt.expectedPrice, actualPrice)
			},
		)
	}
}
