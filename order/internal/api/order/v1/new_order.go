package v1

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/sborsh1kmusora/micro-shop/order/internal/converter"
	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(
	ctx context.Context,
	req *orderV1.CreateOrderRequest,
) (orderV1.CreateOrderRes, error) {
	res, err := a.orderService.Create(ctx, converter.OrderToModel(req))
	if err != nil {
		if errors.Is(err, model.ErrItemsNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Один или несколько товаров не найдены",
			}, nil
		}
		return nil, err
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  res.OrderUUID,
		TotalPrice: res.TotalPrice,
	}, nil
}
