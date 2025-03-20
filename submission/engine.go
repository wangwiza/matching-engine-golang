package submission

import "C"
import (
	"assign2/utils"
	"assign2/wg"
	"context"
	"fmt"
	"io"
	"net"
	"os"
)

type Engine struct {
	wg *wg.WaitGroup
	ob *OrderBook
}

func (e *Engine) Init(ctx context.Context, wg *wg.WaitGroup) {
	e.wg = wg
	e.ob = NewOrderBook(ctx, wg)
}

func (e *Engine) Shutdown(ctx context.Context) {
	e.wg.Wait()
}

func (e *Engine) Accept(ctx context.Context, conn net.Conn) {
	e.wg.Add(2)

	go func() {
		defer e.wg.Done()
		<-ctx.Done()
		conn.Close()
	}()

	// This goroutine handles the connection.
	go func() {
		defer e.wg.Done()
		handleConn(conn, e.ob)
	}()
}

func handleConn(conn net.Conn, ob *OrderBook) {
	defer conn.Close()
	for {
		in, err := utils.ReadInput(conn)
		if err != nil {
			if err != io.EOF {
				_, _ = fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			}
			return
		}
		switch in.OrderType {
		case utils.InputCancel:
			var newOrder Order
			newOrder.Init(in.OrderId, "", 0, 0, CANCEL)
			ob.ordersChan <- &newOrder
			<- newOrder.Processed
		case utils.InputBuy:
			var newOrder Order
			newOrder.Init(in.OrderId, in.Instrument, in.Price, in.Count, BUY)
			ob.ordersChan <- &newOrder
			<- newOrder.Processed
		case utils.InputSell:
			var newOrder Order
			newOrder.Init(in.OrderId, in.Instrument, in.Price, in.Count, SELL)
			ob.ordersChan <- &newOrder
			<- newOrder.Processed
		default:
			fmt.Fprintf(os.Stderr, "Got order: %c %v x %v @ %v ID: %v\n",
				in.OrderType, in.Instrument, in.Count, in.Price, in.OrderId)
		}
	}
}
