package v1

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/converter"
	inventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
)

func (a *api) AddItem(
	ctx context.Context,
	req *inventoryV1.AddItemRequest,
) (*inventoryV1.AddItemResponse, error) {
	orderUUID, err := a.inventoryService.Create(ctx, converter.ItemProtoToModel(req.GetItem()))
	if err != nil {
		return nil, err
	}

	return &inventoryV1.AddItemResponse{
		Uuid: orderUUID,
	}, nil
}
