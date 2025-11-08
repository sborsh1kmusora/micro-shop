package order

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
	"github.com/sborsh1kmusora/micro-shop/order/internal/repository/converter"
)

func (r *repo) Get(ctx context.Context, uuid string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.data[uuid]
	if !ok {
		return nil, model.ErrOrderNotFound
	}

	return converter.OrderRepoToModel(order), nil
}
