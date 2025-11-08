package service

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, orderInfo *model.OrderCreate) (*model.OrderCreateRes, error)
	Get(ctx context.Context, orderUUID string) (*model.Order, error)
	Pay(ctx context.Context, uuid, paymentMethod string) (string, error)
	Cancel(ctx context.Context, uuid string) error
}
