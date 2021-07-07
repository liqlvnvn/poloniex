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
	ws := poloniex.NewPrivateWSClient(apiKey, apiSecret)
	err := ws.Run()
	if err != nil {
		return
	}

	ch, err := ws.ListeningReports()

	for {
		fmt.Printf("\n[EXECUTED] %#v\n", <-ch)
	}
}
