package v1

import (
	def "github.com/sborsh1kmusora/micro-shop/order/internal/client/grpc"
	genInventoryV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient genInventoryV1.InventoryServiceClient
}

func NewClient(generatedClient genInventoryV1.InventoryServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
