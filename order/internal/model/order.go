package model

const (
	OrderStatusPendingPayment string = "pending_payment"
	OrderStatusPaid           string = "paid"
	OrderStatusCanceled       string = "canceled"
)

type Order struct {
	UUID            string
	UserUUID        string
	ItemUuids       []string
	TotalPrice      float32
	TransactionUUID string
	PaymentMethod   string
	Status          string
}

type OrderUpdate struct {
	TransactionUUID *string
	Status          *string
	PaymentMethod   *string
}

type OrderCreate struct {
	UserUUID  string
	ItemUUIDs []string
}

type OrderCreateRes struct {
	OrderUUID  string
	TotalPrice float32
}
