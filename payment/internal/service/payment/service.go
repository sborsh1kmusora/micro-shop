package payment

import def "github.com/sborsh1kmusora/micro-shop/payment/internal/service"

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewPaymentService() *service {
	return &service{}
}
