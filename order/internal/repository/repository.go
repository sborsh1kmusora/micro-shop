package repository

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (string, error)
	Get(ctx context.Context, uuid string) (*model.Order, error)
	Update(ctx context.Context, uuid string, orderUpdate *model.OrderUpdate) error
}
