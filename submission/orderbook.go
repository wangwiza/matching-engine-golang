package submission

import (
	"assign2/wg"
	"context"
)

type OrderBook struct {
	ctx                 context.Context
	parWg               *wg.WaitGroup
	childWg             *wg.WaitGroup
	ordersChan          chan *Order
	orderIDToInstrument map[uint32]string
	instrumentToWorker  map[string]Worker
}

func NewOrderBook(ctx context.Context, parWg *wg.WaitGroup) *OrderBook {
	ob := &OrderBook{
		ctx:                 ctx,
		parWg:               parWg,
		childWg:             &wg.WaitGroup{},
		ordersChan:          make(chan *Order),
		orderIDToInstrument: make(map[uint32]string),
		instrumentToWorker:  make(map[string]Worker),
	}
	parWg.Add(1)

	go ob.dispatchOrders()

	return ob
}

func (ob *OrderBook) dispatchOrders() {
	for {
		select {
		case order := <-ob.ordersChan: // receive order from input chan
			var instrument string
			if order.Type == CANCEL {
				// fetch instrument string from cache
				instrument = ob.orderIDToInstrument[order.ID]
			} else {
				// cache instrument string to order ID
				ob.orderIDToInstrument[order.ID] = order.Instrument
				instrument = order.Instrument
			}

			// get worker
			worker, ok := ob.instrumentToWorker[instrument]

			// if worker for instrument does not exist, add new worker
			if !ok {
				worker = ob.AddNewWorker(instrument)
			}
			// send order to worker
			worker.queue <- order // send order to worker
		case <-ob.ctx.Done():
			ob.childWg.Wait()
			ob.parWg.Done()
			return
		}
	}
}

func (ob *OrderBook) AddNewWorker(instrument string) Worker {
	var newWorker Worker
	newWorker.Init(ob.ctx, ob.childWg)
	ob.instrumentToWorker[instrument] = newWorker
	return newWorker
}
