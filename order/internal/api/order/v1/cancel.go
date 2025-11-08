package v1

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(
	ctx context.Context,
	params orderV1.CancelOrderParams,
) (orderV1.CancelOrderRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderCannotBeCanceled) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: "Заказ оплачен или уже отменен",
			}, nil
		}
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
