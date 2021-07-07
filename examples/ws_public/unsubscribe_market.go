package main

import (
	"fmt"
	"time"

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

	go func() {
		time.Sleep(time.Second * 10)
		ws.UnsubscribeMarket("USDT_BTC")
	}()

	for {
		fmt.Println(<-ws.Subs["USDT_BTC"])
	}
}
