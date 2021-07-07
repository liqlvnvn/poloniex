package main

import (
	"fmt"
	"time"

	"vcshl.b2broker.tech/common/golang-libs/poloniex"
)

const (
	apiKey    = ""
	apiSecret = ""
)

func main() {
	polo, err := poloniex.NewClient(apiKey, apiSecret)
	if err != nil {
		return
	}

	resp, _ := polo.GetBalances()
	fmt.Println("BTT", resp["BTT"], "\nUSDT", resp["USDT"])

	fmt.Println(time.Now(), "starting buy")
	buy, err := polo.Buy("USDT_BTT", 0.0270416, 450.0)
	if err != nil {
		fmt.Println("error while tried to buy:", err)
		return
	}
	fmt.Println("buy order sent")

	fmt.Println("order number:", buy.OrderNumber)
	fmt.Println("resulting trades:", buy.ResultingTrades)

	fmt.Println(time.Now())

	resp, _ = polo.GetBalances()
	fmt.Println("BTT", resp["BTT"], "\nUSDT", resp["USDT"])
}
