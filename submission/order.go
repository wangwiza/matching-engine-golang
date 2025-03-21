package submission

import (
	"fmt"
)

type OrderType byte

const (
	BUY OrderType = iota
	SELL
	CANCEL
)

type Order struct {
	ID          uint32
	Instrument  string
	Price       uint32
	Count       uint32
	Type        OrderType
	ExecutionID uint32
	Processed   chan struct{}
}

func (o *Order) Init(id uint32, instrument string, price uint32, count uint32, orderType OrderType) {
	o.ID = id
	o.Instrument = instrument
	o.Price = price
	o.Count = count
	o.Type = orderType
	o.ExecutionID = 1
	o.Processed = make(chan struct{})
}

func (o *Order) Available() bool {
	return o.Count > 0
}

func (o Order) String() string {
	return fmt.Sprintf("%d %s %d %d %d %d",
		o.ID, o.Instrument, o.Price, o.Count, o.Type,
		o.ExecutionID)
}

func (o Order) Equal(other Order) bool {
	return o.ID == other.ID
}
