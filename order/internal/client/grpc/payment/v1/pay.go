package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	paymentV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/payment/v1"
)

func (c *client) Pay(ctx context.Context, orderID, userID uuid.UUID, paymentMethod string) (
	uuid.UUID,
	error,
) {
	payResp, err := c.payClient.PayOrder(
		ctx, &paymentV1.PayOrderRequest{
			OrderId:       orderID.String(),
			UserUuid:      userID.String(),
			PaymentMethod: paymentV1.PaymentMethod(paymentV1.PaymentMethod_value[paymentMethod]),
		},
	)
	if err != nil {
		return uuid.Nil, model.ErrInternalServerError
	}

	id, err := uuid.Parse(payResp.TransactionUuid)
	if err != nil {
		return uuid.Nil, model.ErrInternalServerError
	}

	return id, nil
}
