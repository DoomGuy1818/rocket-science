package app

import (
	"context"

	v1 "github.com/DoomGuy1818/rocket-science/payment/internal/api/v1"
	"github.com/DoomGuy1818/rocket-science/payment/internal/service"
	"github.com/DoomGuy1818/rocket-science/payment/internal/service/payment"
	paymentV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API   paymentV1.PaymentServiceServer
	paymentService service.Service
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) NewPaymentAPIV1(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = v1.NewAPI(d.PaymentService(ctx))
	}
	return d.paymentV1API
}

func (d *diContainer) PaymentService(_ context.Context) service.Service {
	if d.paymentService == nil {
		d.paymentService = payment.NewService()
	}

	return d.paymentService
}
