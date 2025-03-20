package submission

import (
	"context"
)

type OrderBook struct {
	ctx                     context.Context
	ordersChan              chan *Order
	orderIDToInstrument     map[uint32]string
	instrumentToWorkerQueue map[string]chan *Order
}

func NewOrderBook() *OrderBook {
	ob := &OrderBook{
		ctx:                 context.Background(),
		ordersChan:          make(chan *Order),
		orderIDToInstrument: make(map[uint32]string),
		instrumentToWorker:  make(map[string]Worker),
	}

	go ob.dispatchOrders()

	return ob
}

func (ob *OrderBook) dispatchOrders() {
	for {
		select {
		case <-ob.ctx.Done():
			return
		case order := <-ob.ordersChan: // receive order from input chan
			var instrument string
			if order.Type == CANCEL {
				// fetch instrument string from cache
				instrument = ob.orderIDToInstrument[order.ID]
			} else {
				// cache instrument string to order ID
				ob.orderIDToInstrument[order.ID] = order.Instrument
			}

			// get worker
			worker, ok := ob.instrumentToWorker[instrument]

			// if worker for instrument does not exist, add new worker
			if !ok {
				worker = ob.AddNewWorker(instrument)
			}
			// send order to worker
			worker.queue <- order // send order to worker
		}
	}
}

func (ob *OrderBook) AddNewWorker(instrument string) Worker {
	var newWorker Worker
	newWorker.Init(ctx)
	ob.instrumentToWorker[instrument] = newWorker
	return newWorker
}
