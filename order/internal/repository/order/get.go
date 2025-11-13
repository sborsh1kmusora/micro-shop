package order

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"

	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
)

func (r *repo) Get(ctx context.Context, uuid string) (*model.Order, error) {
	query, args, err := r.sb.
		Select(uuidColumn, userUUIDColumn, statusColumn, paymentMethodColumn, transactionUUIDColumn, itemsUUIDsColumn, totalPriceColumn).
		From(tableName).
		Where(squirrel.Eq{uuidColumn: uuid}).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.pool.QueryRow(ctx, query, args...)

	var order model.Order
	if err := row.Scan(
		&order.UUID,
		&order.UserUUID,
		&order.Status,
		&order.PaymentMethod,
		&order.TransactionUUID,
		&order.ItemUuids,
		&order.TotalPrice,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}
		return nil, err
	}

	return &order, nil
}
