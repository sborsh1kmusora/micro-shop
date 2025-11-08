package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/sborsh1kmusora/micro-shop/payment/internal/model"
)

func (s *service) Pay(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error) {
	transactionUUID := uuid.NewString()

	log.Printf(
		"Заказ %s был успешно оплачен пользователем с uuid %s c помощью %s, transaction uuid %s ",
		orderUUID, userUUID, paymentMethod, transactionUUID,
	)

	return transactionUUID, nil
}
