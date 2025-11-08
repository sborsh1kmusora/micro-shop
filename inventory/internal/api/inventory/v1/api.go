package v1

import (
	"github.com/sborsh1kmusora/micro-shop/inventory/internal/service"
	inventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	inventoryService service.InventoryService
}

func NewApi(inventoryService service.InventoryService) *api {
	return &api{
		inventoryService: inventoryService,
	}
}
