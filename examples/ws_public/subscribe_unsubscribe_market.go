package main

import (
	"fmt"
	"log"
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
	log.Print("Subscribed to USDT_BTC channel.")
	go func() {
		for {
			fmt.Println(<-ws.Subs["USDT_BTC"], ws.Subs)
		}
	}()
	time.Sleep(time.Second * 10)

	err = ws.UnsubscribeMarket("USDT_BTC")
	if err != nil {
		return
	}
	log.Print("Unsubscribed from USDT_BTC channel.")
	time.Sleep(time.Second * 10)

	err = ws.SubscribeMarket("USDT_BTC")
	if err != nil {
		panic(err)
		return
	}
	log.Print("Subscribed to USDT_BTC channel.")
	time.Sleep(time.Second * 10)
}
