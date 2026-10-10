package v1

import paymentV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/payment/v1"

type client struct {
	payClient paymentV1.PaymentServiceClient
}

func NewClient(payClient paymentV1.PaymentServiceClient) *client {
	return &client{
		payClient: payClient,
	}
}
