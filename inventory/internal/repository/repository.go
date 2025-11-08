package repository

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
)

type InventoryRepository interface {
	Create(ctx context.Context, item *model.Item) (string, error)
	Get(ctx context.Context, uuid string) (*model.Item, error)
	List(ctx context.Context) ([]*model.Item, error)
}
