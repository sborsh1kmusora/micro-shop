package v1

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	inventoryService *mocks.InventoryService

	api *api
}

func (s *APISuite) SetupTest() {
	s.inventoryService = mocks.NewInventoryService(s.T())

	s.api = NewApi(
		s.inventoryService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
