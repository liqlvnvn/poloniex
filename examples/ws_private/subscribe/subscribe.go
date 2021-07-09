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
	wsObserver := poloniex.NewWebsocketObserver()

	ws := poloniex.NewPrivateWSClient(wsObserver, apiKey, apiSecret)
	err := ws.Run()
	if err != nil {
		return
	}
	err = ws.SubscribeAccount()
	if err != nil {
		return
	}
	go func() {
		time.Sleep(time.Second * 10)
		ws.UnsubscribeAccount()
	}()
	for {
		fmt.Println(<-ws.Subs["ACCOUNT"])
	}
}
