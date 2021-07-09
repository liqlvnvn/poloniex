package main

import (
	"fmt"
	"time"

	polo "vcshl.b2broker.tech/common/golang-libs/poloniex"
)

const (
	api_key    = ""
	api_secret = ""
)

func main() {
	poloniex := polo.NewPublicClient()

	// resp, err := poloniex.GetTickers()
	// resp, err := poloniex.Get24hVolumes()
	// resp, err := poloniex.GetOrderBook("btc_dgb", 1)
	// resp, err := poloniex.GetPublicTradeHistory("btc_dgb")
	// resp, err := poloniex.GetPublicTradeHistory("btc_sc", time.Now().AddDate(0, 0, -1), time.Now())
	resp, err := poloniex.GetChartData("USDT_BTC", time.Now().AddDate(0, 0, -1), time.Now(), "1d")
	// resp, err := poloniex.GetCurrencies()
	// resp, err := poloniex.GetLoanOrders("BTC")
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
