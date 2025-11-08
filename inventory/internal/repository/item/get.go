package item

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
	"github.com/sborsh1kmusora/micro-shop/inventory/internal/repository/converter"
)

func (r *repo) Get(ctx context.Context, uuid string) (*model.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.data[uuid]
	if !ok {
		return nil, model.ErrItemNotFound
	}

	return converter.ItemRepoToModel(order), nil
}
