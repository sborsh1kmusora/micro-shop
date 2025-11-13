package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
)

func (r *repo) Create(ctx context.Context, order *model.Order) (string, error) {
	orderUUID := uuid.NewString()
	order.UUID = orderUUID

	query, args, err := r.sb.
		Insert(tableName).
		Columns(uuidColumn, userUUIDColumn, itemsUUIDsColumn, statusColumn, paymentMethodColumn, transactionUUIDColumn, totalPriceColumn).
		Values(order.UUID, order.UserUUID, order.ItemUuids, order.Status, order.PaymentMethod, order.TransactionUUID, order.TotalPrice).
		ToSql()
	if err != nil {
		return "", err
	}

	_, execErr := r.pool.Exec(ctx, query, args...)
	if execErr != nil {
		return "", execErr
	}

	return orderUUID, nil
}
