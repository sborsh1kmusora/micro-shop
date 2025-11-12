package order

import (
	"github.com/sborsh1kmusora/micro-shop/order/internal/client/grpc"
	"github.com/sborsh1kmusora/micro-shop/order/internal/repository"
	def "github.com/sborsh1kmusora/micro-shop/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepo repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(
	orderRepo repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
) *service {
	return &service{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
