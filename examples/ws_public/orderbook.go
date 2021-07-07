// the following code shows
// how to access OrderBook fields.
package main

import (
	"fmt"

	polo "vcshl.b2broker.tech/common/golang-libs/poloniex"
)

func main() {
	ws := polo.NewPublicWSClient()
	err := ws.Run()
	if err != nil {
		return
	}

	err = ws.SubscribeMarket("USDT_BTC")
	if err != nil {
		return
	}

	var m polo.WSOrderBook

	for {
		receive := <-ws.Subs["USDT_BTC"]
		updates := receive.([]polo.MarketUpdate)
		for _, v := range updates {
			if v.TypeUpdate == "OrderBookRemove" || v.TypeUpdate == "OrderBookModify" {
				m = v.Data.(polo.WSOrderBook)

				fmt.Printf("Rate:%f, Type:%s, Amount:%f\n",
					m.Rate, m.TypeOrder, m.Amount)
			}
		}
	}
}
