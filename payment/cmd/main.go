package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	payV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

type paymentService struct {
	payV1.UnimplementedPaymentServiceServer
}

func (p *paymentService) PayOrder(_ context.Context, req *payV1.PayOrderRequest) (*payV1.PayOrderResponse, error) {
	ID := uuid.New()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", ID.String())

	return &payV1.PayOrderResponse{
		TransactionUuid: ID.String(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}
	defer func() {
		if cerr := lis.Close(); err != nil {
			log.Printf("failed to close listener: %v", cerr)
		}
	}()

	s := grpc.NewServer()

	payV1.RegisterPaymentServiceServer(s, &paymentService{})

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)

		if err = s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutdown Server ...")

	s.GracefulStop()

	log.Printf("Shutdown Server exiting")
}
