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
	StringValue *string  `bson:"string_value,omitempty"`
	IntValue    *int64   `bson:"int_value,omitempty"`
	DoubleValue *float64 `bson:"double_value,omitempty"`
	BoolValue   *bool    `bson:"bool_value,omitempty"`
}

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type Item struct {
	UUID          string            `bson:"uuid"`
	Name          string            `bson:"name"`
	Description   string            `bson:"description"`
	Price         float32           `bson:"price"`
	StockQuantity int64             `bson:"stock_quantity"`
	Category      Category          `bson:"category"`
	Dimensions    *Dimensions       `bson:"dimensions"`
	Manufacturer  *Manufacturer     `bson:"manufacturer"`
	Tags          []string          `bson:"tags"`
	Metadata      map[string]*Value `bson:"metadata"`
	CreatedAt     time.Time         `bson:"created_at"`
	UpdatedAt     *time.Time        `bson:"updated_at"`
}
