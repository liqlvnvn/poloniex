package main

import (
	"fmt"

	"vcshl.b2broker.tech/common/golang-libs/poloniex"
)

const (
	apiKey    = ""
	apiSecret = ""
)

func main() {
	wsObserver := poloniex.NewWebsocketObserver()

	polo := poloniex.NewPrivateClient(wsObserver, apiKey, apiSecret)

	resp, _ := polo.GetBalances()
	fmt.Println("BTT", resp["BTT"], "\nUSDT", resp["USDT"])

	// Cancel open orders
	orders, _ := polo.GetAllOpenOrders()
	for value, ords := range orders {
		fmt.Println(value)
		fmt.Println(ords)
		for _, val := range ords {
			fmt.Println(val.OrderNumber)
			fmt.Println(polo.CancelOrder(val.OrderNumber))
		}
	}

	resp, _ = polo.GetBalances()
	fmt.Println("BTT", resp["BTT"], "\nUSDT", resp["USDT"])
}
