package item

import (
	"context"

	"github.com/google/uuid"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
	"github.com/sborsh1kmusora/micro-shop/inventory/internal/repository/converter"
)

func (r *repo) Create(ctx context.Context, item *model.Item) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	itemUUID := uuid.NewString()

	item.UUID = itemUUID

	r.data[itemUUID] = converter.ItemToRepoModel(item)

	return itemUUID, nil
}
