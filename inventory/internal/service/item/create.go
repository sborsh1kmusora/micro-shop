package item

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/inventory/internal/model"
)

func (s *service) Create(ctx context.Context, item *model.Item) (string, error) {
	uuid, err := s.orderRepo.Create(ctx, item)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
