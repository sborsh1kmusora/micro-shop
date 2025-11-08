package v1

import (
	"context"

	clientConverter "github.com/sborsh1kmusora/micro-shop/order/internal/client/converter"
	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
	genInventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
)

func (c *client) ListItems(ctx context.Context, filter *model.Filter) ([]*model.Item, error) {
	items, err := c.generatedClient.ListItems(ctx, &genInventoryV1.ListItemsRequest{
		Filter: clientConverter.ItemsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}

	return clientConverter.ItemListToModel(items.Items), nil
}
