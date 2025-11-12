package item

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
	"github.com/sborsh1kmusora/micro-shop/inventory/internal/repository/converter"
	repoModel "github.com/sborsh1kmusora/micro-shop/inventory/internal/repository/model"
)

func (r *repo) List(ctx context.Context) ([]*model.Item, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Printf("failed to close cursor: %v\n", err)
		}
	}()

	var items []*repoModel.Item
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return converter.ItemListToModel(items), nil
}
