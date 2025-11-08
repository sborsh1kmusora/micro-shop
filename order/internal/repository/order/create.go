package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
	"github.com/sborsh1kmusora/micro-shop/order/internal/repository/converter"
)

func (r *repo) Create(ctx context.Context, order *model.Order) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderUUID := uuid.NewString()

	order.UUID = orderUUID

	r.data[orderUUID] = converter.OrderToRepoModel(order)

	return orderUUID, nil
}
