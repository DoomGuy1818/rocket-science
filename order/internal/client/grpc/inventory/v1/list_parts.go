package v1

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/DoomGuy1818/rocket-science/order/internal/client/converter"
	"github.com/DoomGuy1818/rocket-science/order/internal/model"
	inventoryV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, ids []uuid.UUID) ([]model.Part, error) {
	uuids := make([]string, 0, len(ids))
	for _, id := range ids {
		uuids = append(uuids, id.String())
	}

	log.Println(uuids)

	protoResp, err := c.invClient.ListParts(
		ctx, &inventoryV1.ListPartsRequest{
			Filter: &inventoryV1.PartFilters{
				Uuids: uuids,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	parts, err := converter.ProtoModelsToDomain(protoResp.GetParts())
	if err != nil {
		return nil, err
	}

	return parts, nil
}
