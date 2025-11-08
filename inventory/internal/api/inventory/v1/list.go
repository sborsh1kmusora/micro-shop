package v1

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/converter"
	inventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
)

func (a *api) ListItems(
	ctx context.Context,
	req *inventoryV1.ListItemsRequest,
) (*inventoryV1.ListItemsResponse, error) {
	items, err := a.inventoryService.List(ctx, converter.FilterProtoToModel(req.Filter))
	if err != nil {
		return nil, err
	}

	return &inventoryV1.ListItemsResponse{
		Items: converter.ListItemsToProto(items),
	}, nil
}
