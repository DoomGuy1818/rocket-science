package order

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/DoomGuy1818/rocket-science/order/internal/client/grpc/mocks"
	repoMock "github.com/DoomGuy1818/rocket-science/order/internal/repository/mocks"
)

type ApiSuite struct {
	suite.Suite

	repo            *repoMock.MockOrderRepository
	inventoryClient *mocks.MockInventoryClient
	paymentClient   *mocks.MockPaymentClient

	service *service
}

func (s *ApiSuite) SetupTest() {
	s.repo = repoMock.NewMockOrderRepository(s.T())
	s.inventoryClient = mocks.NewMockInventoryClient(s.T())
	s.paymentClient = mocks.NewMockPaymentClient(s.T())

	s.service = NewService(
		s.repo,
		s.inventoryClient,
		s.paymentClient,
	)
}

func TestService(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
