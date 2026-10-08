package conventer

import (
	"github.com/google/uuid"

	domainModel "github.com/DoomGuy1818/rocket-science/order/internal/model"
	repoModel "github.com/DoomGuy1818/rocket-science/order/internal/repository/model"
)

// OrderToRepoModel конвертирует доменную модель в модель репозитория.
func OrderToRepoModel(order domainModel.Order) repoModel.Order {
	return repoModel.Order{
		OrderID:       order.OrderID,
		UserID:        order.UserID,
		Parts:         copyUUIDs(order.PartIDs),
		TotalPrice:    order.TotalPrice,
		TransactionID: order.TransactionID,
		PaymentMethod: repoModel.PaymentMethod(order.PaymentMethod),
		Status:        repoModel.Status(order.Status),
	}
}

// OrderToDomainModel конвертирует модель репозитория в доменную модель.
func OrderToDomainModel(order repoModel.Order) domainModel.Order {
	return domainModel.Order{
		OrderID:       order.OrderID,
		UserID:        order.UserID,
		PartIDs:       copyUUIDs(order.Parts),
		TotalPrice:    order.TotalPrice,
		TransactionID: order.TransactionID,
		PaymentMethod: domainModel.PaymentMethod(order.PaymentMethod),
		Status:        domainModel.Status(order.Status),
	}
}

// copyUUIDs создаёт копию слайса, чтобы модели не делили один и тот же
// underlying array. Сохраняет nil для nil-слайса.
func copyUUIDs(src []uuid.UUID) []uuid.UUID {
	if src == nil {
		return nil
	}
	dst := make([]uuid.UUID, len(src))
	copy(dst, src)
	return dst
}
