package v1

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(
	ctx context.Context,
	req *orderV1.PayOrderRequest,
	params orderV1.PayOrderParams,
) (orderV1.PayOrderRes, error) {
	transactionUUID, err := a.orderService.Pay(ctx, params.OrderUUID, string(req.GetPaymentMethod()))
	if err != nil {
		if errors.Is(err, model.ErrInvalidPaymentStatus) {
			return &orderV1.BadRequestError{
				Code:    400,
				Message: "Заказ уже оплачен или был отменен",
			}, nil
		}
		return nil, err
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
