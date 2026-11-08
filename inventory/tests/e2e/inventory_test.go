package e2e

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/DoomGuy1818/rocket-science/shared/pkg/proto/inventory/v1"
)

var _ = Describe(
	"InventoryService", func() {
		var (
			ctx             context.Context
			cancel          context.CancelFunc
			inventoryClient inventoryV1.InventoryServiceClient
		)

		BeforeEach(
			func() {
				ctx, cancel = context.WithCancel(suiteCtx)

				// Создаём gRPC клиент
				conn, err := grpc.NewClient(
					env.App.Address(),
					grpc.WithTransportCredentials(insecure.NewCredentials()),
				)
				Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

				inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
			},
		)

		AfterEach(
			func() {
				// Чистим коллекцию после теста
				err := env.ClearInventoryCollection(ctx)
				Expect(err).ToNot(HaveOccurred(), "ожидали успешную очистку коллекции sightings")

				cancel()
			},
		)

		Describe(
			"List", func() {
				var partUUID string
				BeforeEach(
					func() {
						// Вставляем тестовое наблюдение
						var err error
						partUUID, err = env.InsertTestPart(ctx)
						Expect(err).ToNot(HaveOccurred(), "ожидали успешную вставку тестовых деталей в MongoDB")
					},
				)

				It(
					"должен успешно возвращасть список запчастей", func() {
						resp, err := inventoryClient.ListParts(
							ctx, &inventoryV1.ListPartsRequest{
								Filter: &inventoryV1.PartFilters{
									Uuids: []string{partUUID},
								},
							},
						)
						Expect(err).ToNot(HaveOccurred(), "ожидали успешный поиск тестовых деталей в MongoDB")

						Expect(resp).ToNot(BeNil())
						Expect(resp.GetParts()).To(HaveLen(1))

						part := resp.GetParts()[0]

						Expect(part.GetUuid()).To(Equal(partUUID))
						Expect(part.GetName()).To(Equal("Ионный двигатель X-200"))
						Expect(part.GetCategory()).To(Equal(inventoryV1.Category_ENGINE))
						Expect(part.GetPrice()).To(Equal(1250000.00))
						Expect(part.GetStockQuantity()).To(Equal(int64(5)))
					},
				)
			},
		)
	},
)
