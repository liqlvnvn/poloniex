// the following code shows
// how to access NewTrade fields.
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

	var n polo.NewTrade

	for {
		receive := <-ws.Subs["USDT_BTC"]
		updates := receive.([]polo.MarketUpdate)
		for _, v := range updates {
			if v.TypeUpdate == "NewTrade" {
				n = v.Data.(polo.NewTrade)
				fmt.Printf("TradeId:%d, Rate:%f, Amount:%f, Total:%f, Type:%s\n",
					n.TradeID, n.Rate, n.Amount, n.Total, n.TypeOrder)
			}
		}
	}
}
