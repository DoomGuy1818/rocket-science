package converter

import (
	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	orderV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/openapi/order/v1"
)

func domainPaymentMethodToResponse(method model.PaymentMethod) orderV1.PaymentMethod {
	return orderV1.PaymentMethod(method)
}

func domainStatusToResponse(status model.Status) orderV1.PaymentStatus {
	return orderV1.PaymentStatus(status)
}

func ModelOrderToApiResponse(order model.Order) *orderV1.Order {
	return &orderV1.Order{
		OrderUUID:       order.OrderID,
		UserUUID:        order.UserID,
		PartUuids:       order.PartIDs,
		TotalPrice:      float32(order.TotalPrice),
		TransactionUUID: order.TransactionID,
		PaymentMethod:   domainPaymentMethodToResponse(order.PaymentMethod),
		Status:          domainStatusToResponse(order.Status),
	}
}
