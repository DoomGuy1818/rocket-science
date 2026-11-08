package v1

import (
	"context"
	"net/http"

	"github.com/DoomGuy1818/rocket-science/order/internal/service"
	orderV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/openapi/order/v1"
)

type orderHandler struct {
	orderService service.Service
}

func NewOrderHandler(orderService service.Service) *orderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}

func (h *orderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}
