package v1

import (
	def "github.com/sborsh1kmusora/micro-shop/order/internal/client/grpc"
	genPaymentV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	generatedClient genPaymentV1.PaymentServiceClient
}

func NewClient(generatedClient genPaymentV1.PaymentServiceClient) *client {
	return &client{generatedClient}
}
