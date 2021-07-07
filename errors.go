package poloniex

import (
	"errors"
	"fmt"
)

// List of errors
const (
	ConnectError          = "[ERROR] Connection could not be established!"
	RequestError          = "[ERROR] NewRequest Error!"
	SetAPIError           = "[ERROR] Set the API KEY and API SECRET!"
	PeriodError           = "[ERROR] Invalid Period!"
	TimePeriodError       = "[ERROR] Time Period incompatibility!"
	TimeError             = "[ERROR] Invalid Time!"
	StartTimeError        = "[ERROR] Start Time Format Error!"
	EndTimeError          = "[ERROR] End Time Format Error!"
	LimitError            = "[ERROR] Limit Format Error!"
	ChannelError          = "[ERROR] Unknown Channel Name: %s"
	SubscribeError        = "[ERROR] Already Subscribed!"
	WSTickerError         = "[ERROR] WSTicker Parsing %s"
	WSOrderBookError      = "[ERROR] MarketUpdate OrderBook Parsing %s"
	NewTradeError         = "[ERROR] MarketUpdate NewTrade Parsing %s"
	WSAccountNotification = "[ERROR] Account Notification Parsing %s"
	WSWrongOrderType      = "[ERROR] Account Notification Parsing: Wrong Order Type %s"
	ServerError           = "[SERVER ERROR] Response: %s"
)

func Error(msg string, args ...interface{}) error {
	if len(args) > 0 {
		return fmt.Errorf("%v %v", msg, args)
	}

	return errors.New(msg)
}
