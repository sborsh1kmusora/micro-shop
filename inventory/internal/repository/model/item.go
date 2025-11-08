package model

import (
	"time"
)

type Category int32

const (
	CategoryUnspecified Category = 0
	CategoryElectronics Category = 1
	CategoryClothing    Category = 2
	CategoryBooks       Category = 3
	CategoryBeauty      Category = 4
)

type Value struct {
	StringValue *string
	IntValue    *int64
	DoubleValue *float64
	BoolValue   *bool
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Item struct {
	UUID          string
	Name          string
	Description   string
	Price         float32
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}
