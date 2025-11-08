package order

import (
	"context"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
)

func (r *repo) Update(ctx context.Context, uuid string, orderUpdate *model.OrderUpdate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[uuid]
	if !ok {
		return model.ErrOrderNotFound
	}

	if orderUpdate.TransactionUUID != nil {
		order.TransactionUUID = *orderUpdate.TransactionUUID
	}

	if orderUpdate.PaymentMethod != nil {
		order.PaymentMethod = *orderUpdate.PaymentMethod
	}

	if orderUpdate.Status != nil {
		order.Status = *orderUpdate.Status
	}

	return nil
}
