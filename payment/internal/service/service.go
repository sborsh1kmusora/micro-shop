package service

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/payment/internal/model"
)

type PaymentService interface {
	Pay(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error)
}
