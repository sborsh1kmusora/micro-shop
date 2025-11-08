package grpc

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
)

type InventoryClient interface {
	ListItems(ctx context.Context, filter *model.Filter) ([]*model.Item, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error)
}
