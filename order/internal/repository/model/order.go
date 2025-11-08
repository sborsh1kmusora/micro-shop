package model

type Order struct {
	UUID            string
	UserUUID        string
	ItemUuids       []string
	TotalPrice      float32
	TransactionUUID string
	PaymentMethod   string
	Status          string
}
