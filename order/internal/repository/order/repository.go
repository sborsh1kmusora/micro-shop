package order

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/sborsh1kmusora/micro-shop/order/internal/repository"
)

const (
	tableName = "orders"

	uuidColumn            = "uuid"
	userUUIDColumn        = "user_uuid"
	itemsUUIDsColumn      = "items_uuids"
	statusColumn          = "status"
	paymentMethodColumn   = "payment_method"
	transactionUUIDColumn = "transaction_uuid"
	totalPriceColumn      = "total_price"
)

var _ def.OrderRepository = (*repo)(nil)

type repo struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
}

func NewOrderRepository(pool *pgxpool.Pool) *repo {
	return &repo{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
