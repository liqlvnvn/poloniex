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

	for {
		fmt.Println(<-ws.Subs["USDT_BTC"])
	}
}
