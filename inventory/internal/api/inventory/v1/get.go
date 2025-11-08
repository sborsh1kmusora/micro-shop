package v1

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/converter"
	inventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
)

func (a *api) GetItem(
	ctx context.Context,
	req *inventoryV1.GetItemRequest,
) (*inventoryV1.GetItemResponse, error) {
	item, err := a.inventoryService.Get(ctx, req.GetUuid())
	if err != nil {
		return nil, err
	}

	return &inventoryV1.GetItemResponse{
		Item: converter.ItemToProto(item),
	}, nil
}
