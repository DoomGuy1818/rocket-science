package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	v1 "github.com/DoomGuy1818/rocket-science/payment/internal/api/v1"
	"github.com/DoomGuy1818/rocket-science/payment/internal/service/payment"
	payV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

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
	serv := payment.NewService()
	api := v1.NewAPI(serv)

	payV1.RegisterPaymentServiceServer(s, api)

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
