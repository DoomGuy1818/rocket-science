package model

import (
	"time"

	"github.com/google/uuid"
)

type Category string

const (
	CategoryUnknown  Category = "UNKNOWN"
	CategoryEngine   Category = "ENGINE"
	CategoryFuel     Category = "FUEL"
	CategoryPorthole Category = "PORTHOLE"
	CategoryWing     Category = "WING"
)

const (
	ValueKindUnknown ValueKind = iota
	ValueKindString
	ValueKindInt64
	ValueKindDouble
	ValueKindBool
)

type ValueKind int

type Value interface {
	isValue()
	Kind() ValueKind
}

type StringValue string

func (StringValue) isValue()        {}
func (StringValue) Kind() ValueKind { return ValueKindString }

type Int64Value int64

func (Int64Value) isValue()        {}
func (Int64Value) Kind() ValueKind { return ValueKindInt64 }

type DoubleValue float32

func (DoubleValue) isValue()        {}
func (DoubleValue) Kind() ValueKind { return ValueKindDouble }

type BoolValue bool

func (BoolValue) isValue()        {}
func (BoolValue) Kind() ValueKind { return ValueKindBool }

type Dimensions struct {
	Length float32
	Width  float32
	Height float32
	Weight float32
}

type Manufacturer struct {
	Name    string
	Country string
}

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         string
	StockQuantity int
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]Value
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
