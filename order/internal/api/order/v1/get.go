package v1

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/sborsh1kmusora/micro-shop/order/internal/converter"
	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(
	ctx context.Context,
	params orderV1.GetOrderParams,
) (orderV1.GetOrderRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Товар не найден",
			}, nil
		}
		return nil, err
	}

	return &orderV1.GetOrderResponse{
		Order: converter.OrderToOpenApi(order),
	}, nil
}
