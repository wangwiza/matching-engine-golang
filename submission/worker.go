package submission

import (
	"assign2/utils"
	"context"
	"time"
)

const QueueSize = 1000

type Worker struct {
	queue  chan *Order
	buyPQ  *SortedList
	sellPQ *SortedList
}

func (w *Worker) Init(ctx context.Context) {
	w.queue = make(chan *Order, QueueSize)
	w.buyPQ = NewSortedList(false)
	w.sellPQ = NewSortedList(true)

	go func() {
		for order := range w.queue {
			w.handleOrder(order)
		}
	}()
}

func (w *Worker) handleOrder(order *Order) {
	switch order.Type {
	case CANCEL:
		w.handleCancel(order)
	case BUY:
		handleBuySell(order, w.sellPQ)
		if order.Available() {
			w.buyPQ.Push(order)
			utils.OutputOrderAdded("B", order.ID, order.Instrument, order.Price, order.Count, GetCurrentTimestamp())
		}
	case SELL:
		handleBuySell(order, w.buyPQ)
		if order.Available() {
			w.sellPQ.Push(order)
			utils.OutputOrderAdded("S", order.ID, order.Instrument, order.Price, order.Count, GetCurrentTimestamp())
		}
	}
}

func priceMatched(activeOrder, bestOrder *Order) bool {
	if activeOrder.Type == BUY {
		return activeOrder.Price >= bestOrder.Price
	} else {
		return activeOrder.Price <= bestOrder.Price
	}
}

func handleBuySell(order *Order, pq *SortedList) {
	for order.Available() && !pq.IsEmpty() {
		bestOrder := pq.Peek()

		if !priceMatched(order, bestOrder) {
			break
		}

		m := min(order.Count, bestOrder.Count)
		order.Count -= m
		bestOrder.Count -= m
		utils.OutputOrderExecuted(bestOrder.ID, order.ID, bestOrder.ExecutionID,
			bestOrder.Price, m, GetCurrentTimestamp())
		bestOrder.ExecutionID += 1

		if bestOrder.Count == 0 {
			pq.Pop()
		}
	}
}

func (w *Worker) handleCancel(cOrder *Order) {
	targetID := cOrder.ID
	inBuy := w.buyPQ.HasID(targetID)
	inSell := w.sellPQ.HasID(targetID)

	if !(inBuy || inSell) {
		utils.OutputOrderDeleted(targetID, false, GetCurrentTimestamp())
		return
	}

	if inBuy {
		w.buyPQ.Remove(targetID)
	} else {
		w.sellPQ.Remove(targetID)
	}
	utils.OutputOrderDeleted(targetID, true, GetCurrentTimestamp())
}

func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano()
}
