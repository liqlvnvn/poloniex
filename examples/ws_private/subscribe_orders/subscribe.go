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
	ws := poloniex.NewPrivateWSClient(wsObserver, apiKey, apiSecret)
	err := ws.Run()
	if err != nil {
		return
	}

	ch, err := ws.ListeningReports()

	go func() {
		polo := poloniex.NewPrivateClient(wsObserver, apiKey, apiSecret)

		resp, _ := polo.GetBalances()
		fmt.Println("BTT", resp["BTT"], "\nUSDT", resp["USDT"])

		buy, err := polo.Sell("USDT_BTT", 0.00243776, 450.0)
		if err != nil {
			fmt.Println("error while tried to buy:", err)
			return
		}
		fmt.Println("buy order sent")

		fmt.Println("order number:", buy.OrderNumber)
		fmt.Println("resulting trades:", buy.ResultingTrades)
	}()

	for {
		fmt.Printf("[EXECUTED] %#v\n", <-ch)
	}
}
