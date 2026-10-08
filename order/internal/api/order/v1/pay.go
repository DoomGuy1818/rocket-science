package v1

import (
	"context"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	orderV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/openapi/order/v1"
)

func (h *orderHandler) CreateOrderPaymentByID(
	ctx context.Context,
	req *orderV1.OrderPaymentRequest,
	params orderV1.CreateOrderPaymentByIDParams,
) (orderV1.CreateOrderPaymentByIDRes, error) {
	id, err := h.orderService.PayOrder(ctx, params.OrderID, uuid.New(), string(req.GetPaymentMethod()))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "cannot get order with id" + params.OrderID.String(),
			}, nil
		}

		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "Cannot get response from payService",
		}, nil
	}

	return &orderV1.OrderPaymentResponse{TransactionUUID: id}, nil
}
