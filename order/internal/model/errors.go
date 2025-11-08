package model

import "github.com/go-faster/errors"

var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrInvalidPaymentStatus  = errors.New("invalid payment status")
	ErrItemsNotFound         = errors.New("items not found")
	ErrOrderCannotBeCanceled = errors.New("order already paid")
)
