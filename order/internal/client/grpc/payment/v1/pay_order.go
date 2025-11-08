package v1

import (
	"context"

	genPaymentV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {
	res, err := c.generatedClient.PayOrder(ctx, &genPaymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: genPaymentV1.PaymentMethod(genPaymentV1.PaymentMethod_value[paymentMethod]),
	})
	if err != nil {
		return "", err
	}

	return res.TransactionUuid, nil
}
