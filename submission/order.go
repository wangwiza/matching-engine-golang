package main

import (
  "fmt"
)

type OrderType int

const (
  BUY OrderType = iota
  SELL
)

// Order represents a trading order with its properties
type Order struct {
  ID          uint32      // Changed from uint64 to uint32
  Instrument  string
  Price       uint32      // Changed from uint64 to uint32
  Count       uint32      // Changed from uint64 to uint32
  Type        OrderType   // Changed OrderType to OrderType
  Timestamp   uint64
  ExecutionID uint32
  Cancelled   bool
}

func (o *Order) Init(id uint32, instrument string, price uint32, count uint32, orderType OrderType, timestamp uint64) {
    o.ID = id
    o.Instrument = instrument
    o.Price = price
    o.Count = count
    o.Type = orderType
    o.Timestamp = timestamp
    o.ExecutionID = 1
    o.Cancelled = false
}

// Available checks if the order is still available for execution
func (o *Order) Available() bool {
  return !o.Cancelled && o.Count > 0
}

// String implements the Stringer interface for pretty-printing
func (o Order) String() string {
  return fmt.Sprintf("%d %s %d %d %d %d %d %t",
    o.ID, o.Instrument, o.Price, o.Count, o.Type,
    o.Timestamp, o.ExecutionID, o.Cancelled)
}

// Equal checks if two orders are equal based on their ID
func (o Order) Equal(other Order) bool {
  return o.ID == other.ID
}

