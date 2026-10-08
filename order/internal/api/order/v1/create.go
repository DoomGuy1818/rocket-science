package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	orderV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/openapi/order/v1"
)

func (h *orderHandler) CreateOrderByID(
	ctx context.Context,
	req *orderV1.OrderCreateRequest,
) (orderV1.CreateOrderByIDRes, error) {
	id, price, err := h.orderService.Create(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		if errors.Is(err, model.ErrCannotProcessOrder) {
			return &orderV1.UnprocessableEntityError{
				Code:    http.StatusUnprocessableEntity,
				Message: "one or more parts do not exist",
			}, nil
		}

		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error:" + err.Error(),
		}, nil
	}

	return &orderV1.OrderCreateResponse{
		OrderUUID:  id,
		TotalPrice: float32(price),
	}, nil
}
