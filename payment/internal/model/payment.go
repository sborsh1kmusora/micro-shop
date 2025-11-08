package model

type PaymentMethod int

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

func (m PaymentMethod) String() string {
	switch m {
	case PaymentMethodCard:
		return "CARD"
	case PaymentMethodSBP:
		return "SBP"
	case PaymentMethodCreditCard:
		return "CREDIT_CARD"
	case PaymentMethodInvestorMoney:
		return "INVESTOR_MONEY"
	default:
		return "UNKNOWN"
	}
}
