package order

import (
	"context"

	"github.com/samber/lo"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
)

func (s *service) Pay(ctx context.Context, uuid, paymentMethod string) (string, error) {
	order, err := s.orderRepo.Get(ctx, uuid)
	if err != nil {
		return "", err
	}

	if order.Status != model.OrderStatusPendingPayment {
		return "", model.ErrInvalidPaymentStatus
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, order.UUID, order.UserUUID, order.PaymentMethod)
	if err != nil {
		return "", err
	}

	err = s.orderRepo.Update(ctx, order.UUID, &model.OrderUpdate{
		TransactionUUID: &transactionUUID,
		Status:          lo.ToPtr(model.OrderStatusPaid),
		PaymentMethod:   &paymentMethod,
	})
	if err != nil {
		return "", err
	}

	return transactionUUID, nil
}
