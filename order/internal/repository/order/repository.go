package order

import (
	"sync"

	def "github.com/sborsh1kmusora/micro-shop/order/internal/repository"
	repoModel "github.com/sborsh1kmusora/micro-shop/order/internal/repository/model"
)

var _ def.OrderRepository = (*repo)(nil)

type repo struct {
	mu   sync.RWMutex
	data map[string]*repoModel.Order
}

func NewOrderRepository() *repo {
	return &repo{
		data: make(map[string]*repoModel.Order),
	}
}
