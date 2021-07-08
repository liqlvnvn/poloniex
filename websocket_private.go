package poloniex

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// subscription and unsubscription on account notification
type authenticatedSubscription struct {
	Command string `json:"command"`
	Channel string `json:"channel"`
	Sign    string `json:"sign"`
	Key     string `json:"key"`
	Payload string `json:"payload"`
}

func (s *authenticatedSubscription) toJSON() ([]byte, bool) {
	j, err := json.Marshal(s)
	if err != nil {
		return j, false
	}

	return j, true
}

// SubscribeAccount make subscription to account notification.
func (ws *WSClient) SubscribeAccount() error {
	return ws.subscribeToAccountNotification(ACCOUNT, "ACCOUNT")
}

// UnsubscribeAccount make unsubscription from account notification.
func (ws *WSClient) UnsubscribeAccount() error {
	return ws.unsubscribe("ACCOUNT")
}

// sub-function for subscription.
func (ws *WSClient) subscribeToAccountNotification(chID int, chName string) (err error) {
	ws.Lock()
	defer ws.Unlock()

	if ws.Subs[chName] == nil {
		ws.Subs[chName] = make(chan interface{}, SUBSBUFFER)
	}

	now := time.Now().UnixNano()

	parameters := make(map[string]string)
	parameters["nonce"] = strconv.FormatInt(now, 10)

	formValues := url.Values{}

	for k, v := range parameters {
		formValues.Set(k, v)
	}

	formData := formValues.Encode()

	sign, err := ws.sign(formData)
	if err != nil {
		return
	}

	authSub := &authenticatedSubscription{
		Command: "subscribe",
		Channel: strconv.Itoa(chID),
		Key:     ws.key,
		Sign:    sign,
		Payload: fmt.Sprintf("nonce=%v", now),
	}
	subsMsg, _ := authSub.toJSON()

	err = ws.writeMessage(subsMsg)
	if err != nil {
		return
	}

	return nil
}

// AccountUpdate represent a single message on an account.
type AccountUpdate struct {
	Data       interface{}
	TypeUpdate string `json:"type"`
}

// ListeningReports make subscription to account executed orders notification.
func (ws *WSClient) ListeningReports() (ch chan Fill, err error) {
	if _, isClientSubscribedToAccountNotification := ws.Subs["ACCOUNT"]; !isClientSubscribedToAccountNotification {
		if err := ws.subscribeToAccountNotification(ACCOUNT, "ACCOUNT"); err != nil {
			return nil, err
		}
	}

	ch = make(chan Fill, SUBSBUFFER)

	go func(ch chan Fill) {
		for updates := range ws.Subs["ACCOUNT"] {
			for _, msg := range updates.([]AccountUpdate) {
				switch msg.TypeUpdate {
				case "Trade":
					trade := msg.Data.(Trade)
					if ws.observer.IsObservable(trade.OrderNumber) {
						servObj, err := ws.observer.Items(trade.OrderNumber)
						if err != nil {
							return
						}

						f := Fill{
							trade.OrderNumber,
							trade.TradeID,
							servObj.symbol,
							trade.Rate,
							trade.Amount,
							servObj.side,
							trade.Date,
						}

						ch <- f
					}
				default:
					continue
				}
			}
		}
	}(ch)

	return ch, nil
}

func convertArgsToAccountNotification(args []interface{}) (res []AccountUpdate, err error) {
	res = make([]AccountUpdate, len(args))
	for i, val := range args {
		vals := val.([]interface{})
		var accountUpdate AccountUpdate

		switch vals[0].(string) {
		case "p":
			var pending Pending

			orderNumber, ok := vals[1].(float64)
			if !ok {
				return nil, Error(WSAccountNotification, "pending.OrderNumber")
			}
			pending.OrderNumber = fmt.Sprintf("%.0f", orderNumber)

			currencyPairID, ok := vals[2].(float64)
			if !ok {
				return nil, Error(WSAccountNotification, "pending.CurrencyPairID")
			}
			pending.CurrencyPairID = fmt.Sprintf("%.0f", currencyPairID)

			rate, ok := vals[3].(string)
			if !ok {
				return nil, Error(WSAccountNotification, "pending.Rate")
			}

			pending.Rate, err = strconv.ParseFloat(rate, 64)
			if err != nil {
				return nil, Error(WSAccountNotification, "pending.Rate")
			}

			amount, ok := vals[4].(string)
			if !ok {
				return nil, Error(WSAccountNotification, "pending.Amount")
			}
			pending.Amount, err = strconv.ParseFloat(amount, 64)
			if err != nil {
				return nil, Error(WSAccountNotification, "pending.Amount")
			}

			orderType, ok := vals[5].(float64)
			if !ok {
				return nil, Error(WSAccountNotification, "pending.OrderType")
			}

			switch orderType {
			case OrderTypeBuy:
				pending.OrderType = OrderTypeBuyValue
			case OrderTypeSell:
				pending.OrderType = OrderTypeSellValue
			default:
				return nil, Error(WSWrongOrderType, "pending.OrderType")
			}

			pending.ClientOrderID = fmt.Sprintf("%v", vals[6])

			pending.EpochMS, ok = vals[7].(string)
			if !ok {
				return nil, Error(WSAccountNotification, "pending.EpochMS")
			}

			accountUpdate.TypeUpdate = MessageTypePending
			accountUpdate.Data = pending

		case "o":
			var orderUpdate OrderUpdate

			orderNumber, ok := vals[1].(float64)
			if !ok {
				return nil, Error(WSAccountNotification, "orderUpdate.OrderNumber")
			}
			orderUpdate.OrderNumber = fmt.Sprintf("%.0f", orderNumber)

			orderUpdate.NewAmount, err = strconv.ParseFloat(vals[2].(string), 64)
			if err != nil {
				return nil, Error(WSAccountNotification, "orderUpdate.NewAmount")
			}

			orderType, ok := vals[3].(string)
			if !ok {
				return nil, Error(WSAccountNotification, "orderUpdate.orderType")
			}

			switch orderType {
			case OrderTypeFill:
				orderUpdate.OrderType = "fill"
			case OrderTypeSelfTrade:
				orderUpdate.OrderType = "self-trade"
			case OrderTypeCanceled:
				orderUpdate.OrderType = "canceled"
			default:
				return nil, Error(WSWrongOrderType, "orderUpdate.OrderType")
			}

			orderUpdate.ClientOrderID = fmt.Sprintf("%v", vals[4])
			orderUpdate.CanceledAmount, err = strconv.ParseFloat(vals[5].(string), 64)
			if err != nil {
				return nil, Error(WSAccountNotification, "orderUpdate.CanceledAmount")
			}

			accountUpdate.TypeUpdate = MessageTypeOrderUpdate
			accountUpdate.Data = orderUpdate

		case "t":
			var trade Trade

			tradeID, ok := vals[1].(float64)
			if !ok {
				return nil, Error(WSAccountNotification, "trade.TradeId")
			}
			trade.TradeID = fmt.Sprintf("%0.f", tradeID)

			trade.Rate, err = strconv.ParseFloat(vals[2].(string), 64)
			if err != nil {
				err = Error(WSAccountNotification, "trade.Rate")
				return
			}

			trade.Amount, err = strconv.ParseFloat(vals[3].(string), 64)
			if err != nil {
				err = Error(WSAccountNotification, "trade.Amount")
				return
			}

			trade.FeeMultiplier, err = strconv.ParseFloat(vals[4].(string), 64)
			if err != nil {
				err = Error(WSAccountNotification, "trade.FeeMultiplier")
				return
			}

			trade.FundingType = fmt.Sprintf("%v", vals[5])

			orderNumber, ok := vals[6].(float64)
			if !ok {
				return nil, Error(WSAccountNotification, "trade.OrderNumber")
			}
			trade.OrderNumber = fmt.Sprintf("%.0f", orderNumber)

			trade.TotalFee, err = strconv.ParseFloat(vals[7].(string), 64)
			if err != nil {
				err = Error(WSAccountNotification, "trade.TotalFee")
				return
			}

			trade.Date = fmt.Sprintf("%v", vals[8])
			trade.ClientOrderID = fmt.Sprintf("%v", vals[9])

			trade.TradeTotal, err = strconv.ParseFloat(vals[10].(string), 64)
			if err != nil {
				err = Error(WSAccountNotification, "trade.TradeTotal")
				return
			}

			trade.EpochMS = fmt.Sprintf("%v", vals[11])

			accountUpdate.TypeUpdate = MessageTypeTrade
			accountUpdate.Data = trade

		case "b":
			var balance BalanceUpdate

			currencyID, ok := vals[1].(float64)
			if !ok {
				return nil, Error(WSAccountNotification, "balance.CurrencyID")
			}
			balance.CurrencyID = fmt.Sprintf("%0.f", currencyID)

			balance.Wallet = fmt.Sprintf("%v", vals[2])
			switch vals[2].(string) {
			case WalletTypeExchange:
				balance.Wallet = "exchange"
			case WalletTypeMargin:
				balance.Wallet = "margin"
			case WalletTypeLending:
				balance.Wallet = "lending"
			default:
				return nil, Error(WSAccountNotification, "unknown balance.Wallet type")
			}

			balance.Amount, err = strconv.ParseFloat(vals[3].(string), 64)
			if err != nil {
				return nil, Error(WSAccountNotification, "balance.Amount")
			}

			balance.Balance, err = strconv.ParseFloat(vals[4].(string), 64)
			if err != nil {
				return nil, Error(WSAccountNotification, "balance.Balance")
			}

			accountUpdate.TypeUpdate = MessageTypeBalance
			accountUpdate.Data = balance

		case "n":
			var order NewOrder

			order.CurrencyPairID = fmt.Sprintf("%v", vals[1])

			orderNumber, ok := vals[2].(float64)
			if !ok {
				return nil, Error(WSAccountNotification, "order.OrderNumber")
			}
			order.OrderNumber = fmt.Sprintf("%.0f", orderNumber)

			switch vals[3].(float64) {
			case OrderTypeBuy:
				order.OrderType = OrderTypeBuyValue
			case OrderTypeSell:
				order.OrderType = OrderTypeSellValue
			default:
				return nil, Error(WSWrongOrderType, "pending.OrderType")
			}

			order.Rate, err = strconv.ParseFloat(vals[4].(string), 64)
			if err != nil {
				return nil, Error(WSAccountNotification, "order.Rate")
			}

			order.Amount, err = strconv.ParseFloat(vals[5].(string), 64)
			if err != nil {
				return nil, Error(WSAccountNotification, "order.Amount")
			}

			order.Date = fmt.Sprintf("%v", vals[6])

			order.OriginalAmountOrdered, err = strconv.ParseFloat(vals[7].(string), 64)
			if err != nil {
				return nil, Error(WSAccountNotification, "order.OriginalAmountOrdered")
			}

			order.ClientOrderID = fmt.Sprintf("%v", vals[8])

			accountUpdate.TypeUpdate = MessageTypeNewOrder
			accountUpdate.Data = order

		case "m":
			var mpu MarginPositionUpdate

			orderNumber, ok := vals[1].(float64)
			if !ok {
				return nil, Error(WSAccountNotification, "mpu.OrderNumber")
			}
			mpu.OrderNumber = fmt.Sprintf("%.0f", orderNumber)

			mpu.Currency = fmt.Sprintf("%v", vals[2])

			mpu.Amount, err = strconv.ParseFloat(vals[3].(string), 64)
			if err != nil {
				return nil, Error(WSAccountNotification, "mpu.Amount")
			}

			mpu.ClientOrderID = fmt.Sprintf("%v", vals[4])

			accountUpdate.TypeUpdate = MessageTypeMargin
			accountUpdate.Data = mpu

		case "k":
			var kill Kill

			orderNumber, ok := vals[1].(float64)
			if !ok {
				return nil, Error("[ERROR] Account Notification Kill Parsing", "kill.OrderNumber")
			}
			kill.OrderNumber = fmt.Sprintf("%.0f", orderNumber)

			kill.ClientOrderID = fmt.Sprintf("%v", vals[2])

			accountUpdate.TypeUpdate = MessageTypeKill
			accountUpdate.Data = kill
		}

		res[i] = accountUpdate
	}

	return res, nil
}
