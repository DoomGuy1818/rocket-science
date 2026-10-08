package v1

import (
	"context"
	"net/http"

	converter "github.com/DoomGuy1818/rocket-science/order/internal/conventer"
	orderV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/openapi/order/v1"
)

func (h *orderHandler) GetOrderByID(ctx context.Context, params orderV1.GetOrderByIDParams) (
	orderV1.GetOrderByIDRes,
	error,
) {
	order, err := h.orderService.Get(ctx, params.OrderID)
	if err != nil {
		return &orderV1.NotFoundError{
			Code:    http.StatusNotFound,
			Message: "cannot get order with id" + params.OrderID.String(),
		}, nil
	}

	return converter.ModelOrderToApiResponse(order), nil
}
