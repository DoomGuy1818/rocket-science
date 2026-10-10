package payment

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ApiSuite struct {
	suite.Suite

	service *service
}

func (s *ApiSuite) SetupSuite() {
	s.service = NewService()
}

func (s *ApiSuite) TearDownSuite() {}

func TestService(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
