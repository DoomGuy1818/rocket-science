package v1

import (
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	payV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *payV1.PayOrderRequest) (*payV1.PayOrderResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	UserID, err := uuid.Parse(req.UserUuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	payMethod := req.PaymentMethod

	transactionID, err := a.inventoryService.PayOrder(ctx, orderID, UserID, payMethod.String())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Printf("%s", transactionID)

	return &payV1.PayOrderResponse{
		TransactionUuid: transactionID.String(),
	}, nil
}
