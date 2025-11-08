package v1

import (
	"github.com/sborsh1kmusora/micro-shop/payment/internal/service"
	paymentV1 "github.com/sborsh1kmusora/micro-shop/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewApi(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
