package payment

import (
	"context"

	"github.com/google/uuid"
)

func (s *ApiSuite) TestOrderPayment() {
	ctx := context.Background()
	userID := uuid.New()
	orderID := uuid.New()
	paymentMethod := "CARD"

	ID, err := s.service.PayOrder(ctx, orderID, userID, paymentMethod)

	s.Require().NoError(err)
	s.Require().NotNil(ID)
}
