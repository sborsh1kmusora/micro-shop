package item

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
	"github.com/sborsh1kmusora/micro-shop/inventory/internal/repository/converter"
	repoModel "github.com/sborsh1kmusora/micro-shop/inventory/internal/repository/model"
)

func (r *repo) List(ctx context.Context) ([]*model.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]*repoModel.Item, 0, len(r.data))
	for _, item := range r.data {
		items = append(items, item)
	}

	return converter.ItemListToModel(items), nil
}
