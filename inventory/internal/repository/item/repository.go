package item

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	def "github.com/sborsh1kmusora/micro-shop/inventory/internal/repository"
)

var _ def.InventoryRepository = (*repo)(nil)

type repo struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repo {
	collection := db.Collection("items")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := collection.Indexes().CreateMany(ctx, indexModels); err != nil {
		panic(err)
	}

	return &repo{
		collection: collection,
	}
}
