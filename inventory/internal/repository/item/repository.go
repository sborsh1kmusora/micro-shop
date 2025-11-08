package item

import (
	"sync"

	def "github.com/sborsh1kmusora/micro-shop/inventory/internal/repository"
	"github.com/sborsh1kmusora/micro-shop/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repo)(nil)

type repo struct {
	mu   sync.RWMutex
	data map[string]*model.Item
}

func NewRepository() *repo {
	return &repo{
		data: make(map[string]*model.Item),
	}
}
