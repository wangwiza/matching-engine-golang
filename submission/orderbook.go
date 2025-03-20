package submission

import (
	"context"
)

type OrderBook struct {
	ctx                 context.Context // added because may need for closing goroutines?
	ordersChan          chan *Order
	orderIDToInstrument map[uint32]string
	instrumentToWorker  map[string]Worker
}

func NewOrderBook() *OrderBook {
	ob := &OrderBook{
		ctx:                 context.Background(),
		ordersChan:          make(chan *Order),
		orderIDToInstrument: make(map[uint32]string),
		instrumentToWorker:  make(map[string]Worker), // currently screams because Worker not defined yet
	}

	go ob.dispatchOrders()

	return ob
}

func (ob *OrderBook) dispatchOrders() {
	for {
		select {
		case order := <-ob.ordersChan: // receive order from input chan
			// fmt.Fprintf(os.Stderr, "Received order: %v\n", order)
			var instrument string
			if order.Type == CANCEL {
				// fetch instrument string from cache
				instrument = ob.orderIDToInstrument[order.ID]
				// fmt.Fprintf(os.Stderr, "Cancel %v mapped to: %s\n", order.ID, instrument)
			} else {
				// cache instrument string to order ID
				ob.orderIDToInstrument[order.ID] = order.Instrument
				instrument = order.Instrument
			}

			// get worker
			worker, ok := ob.instrumentToWorker[instrument]

			// if worker for instrument does not exist, add new worker
			if !ok {
				// fmt.Fprintf(os.Stderr, "Adding new worker for instrument: %s\n", instrument)
				worker = ob.AddNewWorker(instrument)
			}
			// send order to worker
			worker.queue <- order // send order to worker
		}
	}
}

func (ob *OrderBook) AddNewWorker(instrument string) Worker {
	var newWorker Worker
	newWorker.Init(ob.ctx, instrument)
	ob.instrumentToWorker[instrument] = newWorker
	return newWorker
}
