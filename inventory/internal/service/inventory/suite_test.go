package inventory

import (
	"testing"

	"github.com/stretchr/testify/suite"

	repoMock "github.com/DoomGuy1818/rocket-science/inventory/internal/repository/mocks"
)

type ApiSuite struct {
	suite.Suite

	repo *repoMock.MockInventoryRepository

	service *service
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
