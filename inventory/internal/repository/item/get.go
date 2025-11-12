package item

import (
	"context"

	"github.com/go-faster/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
	"github.com/sborsh1kmusora/micro-shop/inventory/internal/repository/converter"
	repoModel "github.com/sborsh1kmusora/micro-shop/inventory/internal/repository/model"
)

func (r *repo) Get(ctx context.Context, uuid string) (*model.Item, error) {
	filter := bson.M{"uuid": uuid}

	var item repoModel.Item
	err := r.collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrItemNotFound
		}
		return nil, err
	}

	return converter.ItemRepoToModel(&item), nil
}
