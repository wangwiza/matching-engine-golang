# Go Matching Engine

[![Go Report Card](https://goreportcard.com/badge/github.com/wangwiza/matching-engine-golang)](https://goreportcard.com/report/github.com/wangwiza/matching-engine-golang)
[![Go Reference](https://pkg.go.dev/badge/github.com/wangwiza/matching-engine-golang.svg)](https://pkg.go.dev/github.com/wangwiza/matching-engine-golang)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A high-performance, low-latency financial order matching engine built in Go. This project implements a Central Limit Order Book (CLOB) and is designed for applications requiring efficient trade execution, such as cryptocurrency exchanges, stock trading platforms, or other financial marketplaces.

---

## Features

* **High Performance**: Built for low-latency order processing and high throughput.
* **Core Order Types**: Supports **Market** and **Limit** orders.
* **Concurrent & Safe**: Designed with concurrency in mind, ensuring thread-safe operations on the order book.
* **Standard Operations**: Full support for adding, canceling, and matching orders.
* **Clear Data Structures**: Clean and understandable implementation of key financial structures like the `Orderbook` and `Trade`.
* **Extensible**: Easily extendable to support more complex order types (e.g., Stop-Loss, IOC).

---

## Architecture

The matching engine's core is the `Orderbook`. It maintains two separate priority queues (implemented as heaps) for **buy (bid)** and **sell (ask)** orders.

1.  **Order Submission**: When a new order is submitted, it's sent to the matching engine.
2.  **Matching Logic**:
    * A new **buy** order is matched against the best-priced **sell** orders in the book.
    * A new **sell** order is matched against the best-priced **buy** orders.
    * Matching continues until the incoming order is fully filled or no more matches are possible at acceptable prices.
3.  **Order Book Update**: If the order is not fully filled, the remaining quantity is placed in the corresponding priority queue within the order book, waiting to be matched by future orders.
4.  **Trade Execution**: When a match occurs, a `Trade` is generated and can be broadcast to notify relevant parties.

This design ensures that orders are matched based on **price-time priority**:
* **Price**: Higher bids match lower asks.
* **Time**: For orders at the same price, the one that arrived earlier gets matched first.

---

## Getting Started

Follow these instructions to get a local copy up and running.

### Prerequisites

* **Go**: Version 1.18 or higher.
* **C++ Compiler**: A C++20 compatible compiler like clang++ or g++.

You can check your Go version with:
```sh
go version
```

### Build

1.  **Clone the repository:**

    ```sh
    git clone [https://github.com/wangwiza/matching-engine-golang.git](https://github.com/wangwiza/matching-engine-golang.git)
    cd matching-engine-golang
    ```

2.  **Build the matching engine:**

    ```sh
    make engine
    ```

3.  **Build the client:**

    ```sh
    make client
    ```

### Running

1.  **Run the matching engine:**

    ```sh
    ./engine
    ```

2.  **In a separate terminal, run the client with an input file:**

    ```sh
    ./client <path_to_socket> < <input_file>
    ```

-----

## Usage Example

Here is a simple example of how to create an order book, place a few orders, and see the matching logic in action.

```go
package main

import (
	"fmt"
	"[github.com/wangwiza/matching-engine-golang/matching](https://github.com/wangwiza/matching-engine-golang/matching)"
)

func main() {
	// 1. Create a new order book for the BTC/USD trading pair
	ob := matching.NewOrderbook("BTC-USD")

	// 2. Create some buy orders
	buyOrderA := matching.NewOrder(true, 5.0)  // Buy 5 BTC
	buyOrderB := matching.NewOrder(true, 2.0)  // Buy 2 BTC

	// 3. Place the buy orders as limit orders
	// The `PlaceLimitOrder` function takes price and the order object
	ob.PlaceLimitOrder(10000.0, buyOrderA)
	ob.PlaceLimitOrder(9900.0, buyOrderB)

	fmt.Println("--- Order Book after placing buy orders ---")
	fmt.Printf("Bids: %+v\n", ob.Bids())

	// 4. Create a sell order that will match the existing buy orders
	sellOrder := matching.NewOrder(false, 6.0) // Sell 6 BTC

	// 5. Place the sell order
	// This sell order (at price 9800) will match the highest bid (10000)
	trades := ob.PlaceLimitOrder(9800.0, sellOrder)

	// 6. Print the executed trades
	fmt.Println("\n--- Trades executed ---")
	for _, trade := range trades {
		fmt.Printf("Matched %f of %s at price %f\n", trade.Amount, trade.TakerOrderID, trade.Price)
	}

	// 7. Check the remaining size of the sell order
	fmt.Printf("\nRemaining size of sell order: %f\n", sellOrder.Size) // Will be 1.0

	// 8. Check the state of the order book
	// The first buy order (buyOrderA) should be gone
	fmt.Println("\n--- Final Order Book State ---")
	fmt.Printf("Bids: %+v\n", ob.Bids())
	fmt.Printf("Asks: %+v\n", ob.Asks())
}
```

### Key Data Structures

  * `Order`: Represents a single buy or sell order.
    ```go
    type Order struct {
        ID        string
        Size      float64
        Bid       bool // true for buy (bid), false for sell (ask)
        Timestamp int64
    }
    ```
  * `Orderbook`: Manages the collections of bids and asks.
    ```go
    type Orderbook struct {
        instrument string
        asks       *AskPriceLevel // Min-heap for asks
        bids       *BidPriceLevel // Max-heap for bids
    }
    ```
  * `Trade`: Represents a successful match between a bid and an ask.
    ```go
    type Trade struct {
        TakerOrderID string
        MakerOrderID string
        Amount       float64
        Price        float64
        Timestamp    int64
    }
    ```

-----

## Contributing

Contributions are welcome\! If you'd like to contribute, please follow these steps:

1.  **Fork** the repository.
2.  Create a new **feature branch** (`git checkout -b feature/AmazingFeature`).
3.  **Commit** your changes (`git commit -m 'Add some AmazingFeature'`).
4.  **Push** to the branch (`git push origin feature/AmazingFeature`).
5.  Open a **Pull Request**.

-----

## License

This project is distributed under the MIT License. See `LICENSE` for more information.
