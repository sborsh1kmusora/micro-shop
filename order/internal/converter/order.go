package converter

import (
	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
	orderV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/openapi/order/v1"
)

func OrderToModel(req *orderV1.CreateOrderRequest) *model.OrderCreate {
	return &model.OrderCreate{
		UserUUID:  req.UserUUID,
		ItemUUIDs: req.ItemUuids,
	}
}

func OrderToOpenApi(order *model.Order) orderV1.Order {
	return orderV1.Order{
		OrderUUID:       order.UUID,
		UserUUID:        order.UserUUID,
		ItemUuids:       order.ItemUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   orderV1.PaymentMethod(order.PaymentMethod),
		Status:          orderV1.OrderStatus(order.Status),
	}
}
