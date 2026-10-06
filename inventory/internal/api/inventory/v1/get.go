package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DoomGuy1818/rocket-science/inventory/internal/conventer"
	"github.com/DoomGuy1818/rocket-science/inventory/internal/model"
	invV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *invV1.GetPartRequest) (*invV1.GetPartResponse, error) {
	ID, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid: %s", err)
	}

	part, err := a.inventoryService.GetByID(ctx, ID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", err)
		}
	}

	return &invV1.GetPartResponse{
		Part: conventer.ModelPartToProto(part),
	}, nil
}
