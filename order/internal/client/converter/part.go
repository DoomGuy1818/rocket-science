package converter

import (
	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	inventoryV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
)

func ProtoModelsToDomain(protoParts []*inventoryV1.Part) ([]model.Part, error) {
	domainParts := make([]model.Part, 0, len(protoParts))

	for _, p := range protoParts {
		id, err := uuid.Parse(p.GetUuid())
		if err != nil {
			return nil, err
		}

		domainParts = append(
			domainParts, model.Part{
				UUID:  id,
				Price: p.GetPrice(),
			},
		)
	}

	return domainParts, nil
}
