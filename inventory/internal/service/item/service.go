package item

import (
	"github.com/sborsh1kmusora/micro-shop/inventory/internal/repository"
	def "github.com/sborsh1kmusora/micro-shop/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	orderRepo repository.InventoryRepository
}

func NewService(orderRepo repository.InventoryRepository) *service {
	return &service{orderRepo: orderRepo}
}
