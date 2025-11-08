package order

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID string) (*model.Order, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		return nil, err
	}

	return order, nil
}
