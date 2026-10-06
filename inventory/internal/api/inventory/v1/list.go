package v1

import (
	"context"

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

	parts := a.inventoryService.ListByFilters(ctx, filter)

	return &invV1.ListPartsResponse{
		Parts: conventer.ModelPartsToProto(parts),
	}, nil
}
