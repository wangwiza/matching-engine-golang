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
	"time"
)

type Engine struct {
	wg *wg.WaitGroup
	ob *OrderBook
}

func (e *Engine) Init(ctx context.Context, wg *wg.WaitGroup) {
	e.wg = wg
	e.ob = NewOrderBook()
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
		case utils.InputBuy:
		case utils.InputSell:
			var orderType OrderType
			if in.OrderType == utils.InputBuy {
				orderType = BUY
			} else {
				orderType = SELL
			}
			newOrder := &Order{
				ID:          in.OrderId,
				Instrument:  in.Instrument,
				Price:       in.Price,
				Count:       in.Count,
				Type:        orderType,
				ExecutionID: 1,
			}
			ob.ordersChan <- newOrder
		default:
			fmt.Fprintf(os.Stderr, "Got order: %c %v x %v @ %v ID: %v\n",
				in.OrderType, in.Instrument, in.Count, in.Price, in.OrderId)
			utils.OutputOrderAdded(in, GetCurrentTimestamp())
		}
		utils.OutputOrderExecuted(123, 124, 1, 2000, 10, GetCurrentTimestamp())
	}
}

func GetCurrentTimestamp() uint64 {
	return uint64(time.Now().UnixNano())
}
