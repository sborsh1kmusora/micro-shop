package item

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (*model.Item, error) {
	item, err := s.orderRepo.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return item, nil
}
