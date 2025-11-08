package converter

import (
	"github.com/sborsh1kmusora/micro-shop/order/internal/model"
	repoModel "github.com/sborsh1kmusora/micro-shop/order/internal/repository/model"
)

func OrderToRepoModel(order *model.Order) *repoModel.Order {
	return &repoModel.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		ItemUuids:       order.ItemUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}
}

func OrderRepoToModel(order *repoModel.Order) *model.Order {
	return &model.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		ItemUuids:       order.ItemUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}
}
