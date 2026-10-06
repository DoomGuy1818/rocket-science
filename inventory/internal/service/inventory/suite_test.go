package inventory

import (
	"testing"

	"github.com/stretchr/testify/suite"

	repoMock "github.com/DoomGuy1818/rocket-science/inventory/internal/repository/mocks"
	"github.com/DoomGuy1818/rocket-science/inventory/internal/service"
)

type ApiSuite struct {
	suite.Suite

	repo *repoMock.MockInventoryRepository

	service service.InventoryService
}

func (s *ApiSuite) SetupSuite() {
	s.repo = repoMock.NewMockInventoryRepository(s.T())

	s.service = NewService(
		s.repo,
	)
}

func (s *ApiSuite) TearDownSuite() {}

func TestService(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
