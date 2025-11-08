package order

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/openapi/order/v1"
)

func (s *service) Create(ctx context.Context, orderInfo *model.OrderCreate) (*model.OrderCreateRes, error) {
	items, err := s.inventoryClient.ListItems(ctx, &model.Filter{UUIDs: orderInfo.ItemUUIDs})
	if err != nil {
		return nil, err
	}

	if len(items) != len(orderInfo.ItemUUIDs) {
		return nil, model.ErrItemsNotFound
	}

	var totalPrice float32
	for _, item := range items {
		totalPrice += item.Price
	}

	order := &model.Order{
		UserUUID:   orderInfo.UserUUID,
		ItemUuids:  orderInfo.ItemUUIDs,
		TotalPrice: totalPrice,
		Status:     string(orderV1.OrderStatusPendingPayment),
	}

	orderUUID, err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return &model.OrderCreateRes{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}
