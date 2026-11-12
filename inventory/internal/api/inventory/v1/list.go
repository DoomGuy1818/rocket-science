package v1

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/conventer"
	invV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *invV1.ListPartsRequest) (*invV1.ListPartsResponse, error) {
	filter, err := conventer.PartFilterToModel(req.GetFilter())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	parts, err := a.inventoryService.ListByFilters(ctx, filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, part := range parts {
		log.Printf(
			"Part: %s, Metadata: %#v",
			part.Name,
			part.Metadata,
		)
	}

	return &invV1.ListPartsResponse{
		Parts: conventer.ModelPartsToProto(parts),
	}, nil
}
