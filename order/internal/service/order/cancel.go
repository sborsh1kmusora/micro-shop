package order

import (
	"context"

	"github.com/samber/lo"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, uuid string) error {
	order, err := s.orderRepo.Get(ctx, uuid)
	if err != nil {
		return err
	}

	if order.Status != model.OrderStatusPendingPayment {
		return model.ErrOrderCannotBeCanceled
	}

	err = s.orderRepo.Update(ctx, order.UUID, &model.OrderUpdate{
		Status: lo.ToPtr(model.OrderStatusCanceled),
	})
	if err != nil {
		return err
	}

	return nil
}
