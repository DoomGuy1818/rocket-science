package v1

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	orderV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/openapi/order/v1"
)

func (h *orderHandler) CancelOrderByID(
	ctx context.Context,
	params orderV1.CancelOrderByIDParams,
) (orderV1.CancelOrderByIDRes, error) {
	if err := h.orderService.Cancel(ctx, params.OrderID); err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "cannot get order with id" + params.OrderID.String(),
			}, nil
		}

		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "cannot cancel order with id" + params.OrderID.String(),
		}, nil
	}

	return &orderV1.CancelOrderByIDNoContent{}, nil
}
