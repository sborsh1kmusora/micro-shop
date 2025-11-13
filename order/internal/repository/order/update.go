package order

import (
	"context"

	"github.com/Masterminds/squirrel"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
)

func (r *repo) Update(ctx context.Context, uuid string, orderUpdate *model.OrderUpdate) error {
	builder := r.sb.Update(tableName)

	if orderUpdate.TransactionUUID != nil {
		builder = builder.Set(transactionUUIDColumn, *orderUpdate.TransactionUUID)
	}
	if orderUpdate.PaymentMethod != nil {
		builder = builder.Set(paymentMethodColumn, *orderUpdate.PaymentMethod)
	}
	if orderUpdate.Status != nil {
		builder = builder.Set(statusColumn, *orderUpdate.Status)
	}

	builder = builder.Where(squirrel.Eq{uuidColumn: uuid})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	cmdTag, execErr := r.pool.Exec(ctx, query, args...)
	if execErr != nil {
		return execErr
	}

	if cmdTag.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
