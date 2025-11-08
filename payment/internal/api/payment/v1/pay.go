package v1

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/payment/internal/model"
	paymentV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(
	ctx context.Context,
	req *paymentV1.PayOrderRequest,
) (*paymentV1.PayOrderResponse, error) {
	transactionUUID, err := a.paymentService.Pay(ctx, req.GetOrderUuid(), req.GetUserUuid(), model.PaymentMethod(req.GetPaymentMethod()))
	if err != nil {
		return nil, err
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
