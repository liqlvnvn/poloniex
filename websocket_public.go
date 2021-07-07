package poloniex

import (
	"encoding/json"
	"strconv"
	"strings"
)

// subscription and unsubscription
type subscription struct {
	Command string `json:"command"`
	Channel string `json:"channel"`
}

func (s subscription) toJSON() ([]byte, bool) {
	j, err := json.Marshal(s)
	if err != nil {
		return j, false
	}

	return j, true
}

// MarketUpdate is for market update.
type MarketUpdate struct {
	Data       interface{}
	TypeUpdate string `json:"type"`
}

// SubscribeTicker subscribes to ticker channel.
// It returns nil if successful.
func (ws *WSClient) SubscribeTicker() error {
	return ws.subscribe(TICKER, "TICKER")
}

// UnsubscribeTicker unsubscribes from ticker channel.
// It returns nil if successful.
func (ws *WSClient) UnsubscribeTicker() error {
	return ws.unsubscribe("TICKER")
}

// SubscribeMarket subscribes to market channel.
// It returns nil if successful.
func (ws *WSClient) SubscribeMarket(chName string) error {
	chName = strings.ToUpper(chName)
	chID, ok := channelsByName[chName]
	if !ok {
		return Error(ChannelError, chName)
	}

	return ws.subscribe(chID, chName)
}

// UnsubscribeMarket unsubscribes from market channel.
// It returns nil if successful.
func (ws *WSClient) UnsubscribeMarket(chName string) error {
	chName = strings.ToUpper(chName)
	_, ok := channelsByName[chName]
	if !ok {
		return Error(ChannelError, chName)
	}

	return ws.unsubscribe(chName)
}

// Convert ticker update arguments and fill wsTicker.
func convertArgsToTicker(args []interface{}) (wsTicker WSTicker, err error) {
	wsTicker.Symbol = channelsByID[int(args[0].(float64))]

	wsTicker.Last, err = strconv.ParseFloat(args[1].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "Last")
		return
	}

	wsTicker.LowestAsk, err = strconv.ParseFloat(args[2].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "LowestAsk")
		return
	}

	wsTicker.HighestBid, err = strconv.ParseFloat(args[3].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "HighestBid")
		return
	}

	wsTicker.PercentChange, err = strconv.ParseFloat(args[4].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "PercentChange")
		return
	}

	wsTicker.BaseVolume, err = strconv.ParseFloat(args[5].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "BaseVolume")
		return
	}

	wsTicker.QuoteVolume, err = strconv.ParseFloat(args[6].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "QuoteVolume")
		return
	}

	if v, ok := args[7].(float64); ok {
		if v == 0 {
			wsTicker.IsFrozen = false
		} else {
			wsTicker.IsFrozen = true
		}
	} else {
		err = Error(WSTickerError, "IsFrozen")
		return
	}

	wsTicker.High24hr, err = strconv.ParseFloat(args[8].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "High24hr")
		return
	}

	wsTicker.Low24hr, err = strconv.ParseFloat(args[9].(string), 64)
	if err != nil {
		err = Error(WSTickerError, "Low24hr")
		return
	}

	return wsTicker, nil
}

// Convert market update arguments and fill marketUpdate.
func convertArgsToMarketUpdate(args []interface{}) (res []MarketUpdate, err error) {
	res = make([]MarketUpdate, len(args))
	for i, val := range args {
		vals := val.([]interface{})
		var marketUpdate MarketUpdate

		switch vals[0].(string) {
		case "i":
			var orderDepth OrderDepth
			val := vals[1].(map[string]interface{})
			orderDepth.Symbol = val["currencyPair"].(string)

			asks := val["orderBook"].([]interface{})[0].(map[string]interface{})
			bids := val["orderBook"].([]interface{})[1].(map[string]interface{})

			for k, v := range bids {
				price, _ := strconv.ParseFloat(k, 64)
				quantity, _ := strconv.ParseFloat(v.(string), 64)
				book := Book{Price: price, Quantity: quantity}
				orderDepth.OrderBook.Bids = append(orderDepth.OrderBook.Bids, book)
			}

			for k, v := range asks {
				price, _ := strconv.ParseFloat(k, 64)
				quantity, _ := strconv.ParseFloat(v.(string), 64)
				book := Book{Price: price, Quantity: quantity}
				orderDepth.OrderBook.Asks = append(orderDepth.OrderBook.Asks, book)
			}

			marketUpdate.TypeUpdate = "OrderDepth"
			marketUpdate.Data = orderDepth

		case "o":
			var orderDataField WSOrderBook

			if vals[3].(string) == "0.00000000" {
				marketUpdate.TypeUpdate = "OrderBookRemove"
			} else {
				marketUpdate.TypeUpdate = "OrderBookModify"
			}

			if vals[1].(float64) == 1 {
				orderDataField.TypeOrder = "bid"
			} else {
				orderDataField.TypeOrder = "ask"
			}

			orderDataField.Rate, err = strconv.ParseFloat(vals[2].(string), 64)
			if err != nil {
				err = Error(WSOrderBookError, "Rate")
				return
			}

			orderDataField.Amount, err = strconv.ParseFloat(vals[3].(string), 64)
			if err != nil {
				err = Error(WSOrderBookError, "Amount")
				return
			}

			marketUpdate.Data = orderDataField

		case "t":
			var tradeDataField NewTrade

			tradeDataField.TradeID, err = strconv.ParseInt(vals[1].(string), 10, 64)
			if err != nil {
				err = Error(NewTradeError, "TradeID")
				return
			}

			if vals[2].(float64) == 1 {
				tradeDataField.TypeOrder = "buy"
			} else {
				tradeDataField.TypeOrder = "sell"
			}

			tradeDataField.Rate, err = strconv.ParseFloat(vals[3].(string), 64)
			if err != nil {
				err = Error(NewTradeError, "Rate")
				return
			}

			tradeDataField.Amount, err = strconv.ParseFloat(vals[4].(string), 64)
			if err != nil {
				err = Error(NewTradeError, "Amount")
				return
			}

			tradeDataField.Total = vals[5].(float64)

			marketUpdate.TypeUpdate = "NewTrade"
			marketUpdate.Data = tradeDataField
		}
		res[i] = marketUpdate
	}

	return res, nil
}
