package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	invV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type inventoryService struct {
	invV1.UnimplementedInventoryServiceServer

	mu    sync.RWMutex
	parts map[string]*invV1.Part
}

func (i *inventoryService) GetPart(_ context.Context, req *invV1.GetPartRequest) (*invV1.GetPartResponse, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	part, ok := i.parts[req.Uuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part with uuid %s not found", req.Uuid)
	}

	return &invV1.GetPartResponse{
		Part: part,
	}, nil
}

func (i *inventoryService) ListParts(_ context.Context, req *invV1.ListPartsRequest) (*invV1.ListPartsResponse, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	response := &invV1.ListPartsResponse{
		Parts: make([]*invV1.Part, 0),
	}

	filter := req.GetFilter()

	for _, part := range i.parts {
		if filter != nil && !i.matchesFilter(part, filter) {
			continue
		}

		response.Parts = append(response.Parts, part)
	}

	return response, nil
}

func GetPartsById(ctx context.Context, partID uuid.UUID, client invV1.InventoryServiceClient) (*invV1.Part, error) {
	partUuid := invV1.GetPartRequest{
		Uuid: partID.String(),
	}

	resp, err := client.GetPart(ctx, &partUuid)
	if err != nil {
		return nil, err
	}

	return resp.GetPart(), nil
}

func (i *inventoryService) matchesFilter(part *invV1.Part, filter *invV1.PartFilters) bool {
	if len(filter.GetUuids()) > 0 && !contains(filter.GetUuids(), part.GetUuid()) {
		return false
	}

	if len(filter.GetNames()) > 0 && !contains(filter.GetNames(), part.GetName()) {
		return false
	}

	if len(filter.GetCategories()) > 0 && !contains(filter.GetCategories(), part.GetCategory()) {
		return false
	}

	if len(filter.GetManufacturerCountries()) > 0 {
		if part.GetManufacturer() == nil ||
			!contains(filter.GetManufacturerCountries(), part.GetManufacturer().GetCountry()) {
			return false
		}
	}

	if len(filter.GetTags()) > 0 && !hasAnyTag(part.GetTags(), filter.GetTags()) {
		return false
	}

	return true
}

func contains[T comparable](values []T, target T) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}

func hasAnyTag(partTags, filterTags []string) bool {
	for _, partTag := range partTags {
		for _, filterTag := range filterTags {
			if partTag == filterTag {
				return true
			}
		}
	}

	return false
}

func initTestParts() map[string]*invV1.Part {
	testParts := []*invV1.Part{
		{
			Uuid:     "124f3ba2-b266-4814-985c-f8bf5cf16793",
			Name:     "Ракетный двигатель RD-180",
			Category: invV1.Category_ENGINE,
			Price:    12500000.00,
			Manufacturer: &invV1.Manufacturer{
				Name:    "Энергомаш",
				Country: "Россия",
			},
			Tags: []string{"engine", "liquid-fuel", "heavy-lift"},
		},
		{
			Uuid:     "523c2f78-20d0-404f-96ce-fe43d7675283",
			Name:     "Топливный бак Х-1",
			Category: invV1.Category_FUEL,
			Price:    850000.00,
			Manufacturer: &invV1.Manufacturer{
				Name:    "SpaceTank Inc.",
				Country: "США",
			},
			Tags: []string{"fuel", "tank", "aluminum"},
		},
		{
			Uuid:     "e900893c-f486-464d-b888-f671a2822962",
			Name:     "Иллюминатор кабины",
			Category: invV1.Category_PORTHOLE,
			Price:    45000.00,
			Manufacturer: &invV1.Manufacturer{
				Name:    "GlassCraft",
				Country: "Германия",
			},
			Tags: []string{"porthole", "cabin", "glass"},
		},
		{
			Uuid:     "a489bcb0-245c-4f2e-9559-f02dc38ebc43",
			Name:     "Крыло стабилизации",
			Category: invV1.Category_WING,
			Price:    2300000.00,
			Manufacturer: &invV1.Manufacturer{
				Name:    "AeroWings",
				Country: "Франция",
			},
			Tags: []string{"wing", "stabilizer", "aluminum"},
		},
		{
			Uuid:     "71b062a5-4dd2-48d3-b830-bf689b9b4c5d",
			Name:     "Ракетный двигатель Merlin",
			Category: invV1.Category_ENGINE,
			Price:    9800000.00,
			Manufacturer: &invV1.Manufacturer{
				Name:    "SpaceX",
				Country: "США",
			},
			Tags: []string{"engine", "liquid-fuel", "reusable"},
		},
		{
			Uuid:     "02f6e5fe-015d-4765-a464-8abdd0791ca9",
			Name:     "Топливный бак Y-2",
			Category: invV1.Category_FUEL,
			Price:    920000.00,
			Manufacturer: &invV1.Manufacturer{
				Name:    "Энергомаш",
				Country: "Россия",
			},
			Tags: []string{"fuel", "tank", "titanium"},
		},
	}

	parts := make(map[string]*invV1.Part, len(testParts))
	for _, part := range testParts {
		parts[part.GetUuid()] = part
	}

	return parts
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

	service := &inventoryService{
		parts: initTestParts(),
	}

	invV1.RegisterInventoryServiceServer(s, service)

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
