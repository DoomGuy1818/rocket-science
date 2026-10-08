package order

import (
	"github.com/google/uuid"

	repoModel "github.com/DoomGuy1818/rocket-science/order/internal/repository/model"
)

type repository struct {
	orders map[uuid.UUID]*repoModel.Order
}

func NewStorage() *repository {
	return &repository{
		orders: make(map[uuid.UUID]*repoModel.Order),
	}
}

func copyPartUUIDs(src []uuid.UUID) []uuid.UUID {
	if src == nil {
		return nil
	}

	dst := make([]uuid.UUID, len(src))

	copy(dst, src)

	return dst
}
