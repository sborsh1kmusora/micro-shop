package item

import (
	"context"

	"github.com/google/uuid"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
	"github.com/sborsh1kmusora/micro-shop/inventory/internal/repository/converter"
)

func (r *repo) Create(ctx context.Context, item *model.Item) (string, error) {
	item.UUID = uuid.NewString()

	doc := converter.ItemToRepoModel(item)

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return "", err
	}

	return item.UUID, nil
}
