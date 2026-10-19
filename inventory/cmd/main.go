package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	v1 "github.com/DoomGuy1818/rocket-science/inventory/internal/api/inventory/v1"
	"github.com/DoomGuy1818/rocket-science/inventory/internal/repository/inventory"
	inventoryService "github.com/DoomGuy1818/rocket-science/inventory/internal/service/inventory"
	invV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

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

	ctx := context.Background()

	dbURI := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}

	defer func() {
		cerr := client.Disconnect(ctx)
		if cerr != nil {
			log.Printf("failed to disconnect: %v\n", cerr)
		}
	}()

	if err = client.Ping(ctx, nil); err != nil {
		log.Printf("failed to ping database: %v\n", err)
	}

	db := client.Database("parts")

	repo := inventory.NewRepository(db)
	serv := inventoryService.NewService(repo)
	api := v1.NewAPI(serv)

	invV1.RegisterInventoryServiceServer(s, api)

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
