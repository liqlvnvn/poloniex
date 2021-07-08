package poloniex

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Balance struct {
	Available decimal.Decimal `json:"available"`
	OnOrders  decimal.Decimal `json:"onOrders"`
	BtcValue  decimal.Decimal `json:"btcValue"`
}

func (p *Poloniex) GetBalances() (balances map[string]string, err error) {
	respCh := make(chan []byte)
	errCh := make(chan error)

	go p.tradingRequest("returnBalances", nil, respCh, errCh)

	resp := <-respCh
	err = <-errCh

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &balances)
	return
}

func (p *Poloniex) GetCompleteBalances() (completeBalances map[string]Balance, err error) {
	respCh := make(chan []byte)
	errCh := make(chan error)

	go p.tradingRequest("returnCompleteBalances", nil, respCh, errCh)

	resp := <-respCh
	err = <-errCh

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &completeBalances)
	return
}

type Accounts struct {
	Margin   map[string]decimal.Decimal `json:"margin"`
	Lending  map[string]decimal.Decimal `json:"lending"`
	Exchange map[string]decimal.Decimal `json:"exchange"`
}

func (p *Poloniex) GetAccountBalances() (accounts Accounts, err error) {
	respCh := make(chan []byte)
	errCh := make(chan error)

	go p.tradingRequest("returnAvailableAccountBalances", nil, respCh, errCh)

	resp := <-respCh
	err = <-errCh

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &accounts)
	return
}

type NewAddress struct {
	Success  int    `json:"success"`
	Response string `json:"response"`
}

func (p *Poloniex) GetDepositAddresses() (depositAddresses map[string]string, err error) {
	respCh := make(chan []byte)
	errCh := make(chan error)

	go p.tradingRequest("returnDepositAddresses", nil, respCh, errCh)

	resp := <-respCh
	err = <-errCh

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &depositAddresses)
	return
}

func (p *Poloniex) GenerateNewAddress(currency string) (newAddress NewAddress, err error) {
	respCh := make(chan []byte)
	errCh := make(chan error)

	parameters := map[string]string{"currency": strings.ToUpper(currency)}
	go p.tradingRequest("generateNewAddress", parameters, respCh, errCh)

	resp := <-respCh
	err = <-errCh

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &newAddress)
	return
}

type OpenOrder struct {
	OrderNumber    string          `json:"orderNumber"`
	Type           string          `json:"type"`
	Price          decimal.Decimal `json:"rate"`
	StartingAmount decimal.Decimal `json:"startingAmount"`
	Amount         decimal.Decimal `json:"amount"`
	Total          decimal.Decimal `json:"total"`
	Date           string          `json:"date"`
	Margin         int             `json:"margin"`
}

// GetOpenOrders is send market to get open orders.
func (p *Poloniex) GetOpenOrders(market string) (openOrders []OpenOrder, err error) {
	respCh := make(chan []byte)
	errCh := make(chan error)

	parameters := map[string]string{"currencyPair": strings.ToUpper(market)}
	go p.tradingRequest("returnOpenOrders", parameters, respCh, errCh)

	resp := <-respCh
	err = <-errCh

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &openOrders)
	return
}

// GetAllOpenOrders returns all open orders.
func (p *Poloniex) GetAllOpenOrders() (openOrders map[string][]OpenOrder, err error) {
	respCh := make(chan []byte)
	errCh := make(chan error)

	parameters := map[string]string{"currencyPair": "all"}
	go p.tradingRequest("returnOpenOrders", parameters, respCh, errCh)

	resp := <-respCh
	err = <-errCh

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &openOrders)
	if err != nil {
		return
	}

	for k, v := range openOrders {
		if len(v) == 0 {
			delete(openOrders, k)
		}
	}
	return
}

type CancelOrder struct {
	Success int `json:"success"`
}

func (p *Poloniex) CancelOrder(orderNumber string) (cancelorder CancelOrder, err error) {
	respCh := make(chan []byte)
	errCh := make(chan error)

	parameters := map[string]string{"orderNumber": orderNumber}
	go p.tradingRequest("cancelOrder", parameters, respCh, errCh)

	resp := <-respCh
	err = <-errCh

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &cancelorder)
	return
}

type TradeHistory struct {
	GlobalTradeID int             `json:"globalTradeId"`
	TradeID       string          `json:"tradeId"`
	Date          string          `json:"date"`
	Price         decimal.Decimal `json:"rate"`
	Amount        decimal.Decimal `json:"amount"`
	Total         decimal.Decimal `json:"total"`
	Fee           decimal.Decimal `json:"fee"`
	OrderNumber   decimal.Decimal `json:"orderNumber"`
	Type          string          `json:"type"`
	Category      string          `json:"category"`
}

func (p *Poloniex) GetTradeHistory(market string, start, end time.Time, limit int) (tradehistory []TradeHistory, err error) {
	parameters := map[string]string{
		"currencyPair": strings.ToUpper(market),
		"start":        strconv.FormatInt(start.Unix(), 10),
		"end":          strconv.FormatInt(end.Unix(), 10),
		"limit":        strconv.Itoa(limit),
	}

	respCh := make(chan []byte)
	errCh := make(chan error)

	go p.tradingRequest("returnTradeHistory", parameters, respCh, errCh)

	resp := <-respCh
	err = <-errCh

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &tradehistory)
	return
}

type OrderTrade struct {
	GlobalTradeID decimal.Decimal `json:"globalTradeId"`
	TradeID       decimal.Decimal `json:"tradeId"`
	Market        string          `json:"currencyPair"`
	Type          string          `json:"type"`
	Price         decimal.Decimal `json:"rate"`
	Amount        decimal.Decimal `json:"amount"`
	Total         decimal.Decimal `json:"total"`
	Fee           decimal.Decimal `json:"fee"`
	Date          string          `json:"date"`
}

func (p *Poloniex) GetTradesByOrderID(orderNumber string) (ordertrades []OrderTrade, err error) {
	respCh := make(chan []byte)
	errCh := make(chan error)

	parameters := map[string]string{"orderNumber": orderNumber}
	go p.tradingRequest("returnOrderTrades", parameters, respCh, errCh)

	resp := <-respCh
	err = <-errCh

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &ordertrades)
	return
}

type OrderStat struct {
	Status         string          `json:"status"`
	Rate           decimal.Decimal `json:"rate"`
	Amount         decimal.Decimal `json:"amount"`
	CurrencyPair   string          `json:"currencyPair"`
	Date           string          `json:"date"`
	Total          decimal.Decimal `json:"total"`
	Type           string          `json:"type"`
	StartingAmount decimal.Decimal `json:"startingAmount"`
}

// OrderStat1 represent error result
type OrderStat1 struct {
	Success int `json:"success"`
	Result  struct {
		Error string `json:"error"`
	} `json:"result"`
}

// OrderStat2 represent success result
type OrderStat2 struct {
	Success int                  `json:"success"`
	Result  map[string]OrderStat `json:"result"`
}

func (p *Poloniex) GetOrderStat(orderNumber string) (orderStat OrderStat, err error) {
	var check1 OrderStat1
	var check2 OrderStat2

	respCh := make(chan []byte)
	errCh := make(chan error)

	parameters := map[string]string{"orderNumber": orderNumber}
	go p.tradingRequest("returnOrderStatus", parameters, respCh, errCh)

	resp := <-respCh
	err = <-errCh

	if err != nil {
		return
	}

	// check error
	err = json.Unmarshal(resp, &check1)
	if err != nil {
		return
	}
	if check1.Success == 0 && len(check1.Result.Error) > 0 {
		err = errors.New(check1.Result.Error)
		return

	}

	// check success
	err = json.Unmarshal(resp, &check2)
	if err != nil {
		return
	}
	if check2.Success == 1 {
		orderStat = check2.Result[orderNumber]
		return
	}

	return orderStat, errors.New("unexpected result")
}

type ResultTrades struct {
	Amount  decimal.Decimal `json:"amount"`
	Date    string          `json:"date"`
	Rate    decimal.Decimal `json:"rate"`
	Total   decimal.Decimal `json:"total"`
	TradeID decimal.Decimal `json:"tradeId"`
	Type    string          `json:"type"`
}

type Buy struct {
	OrderNumber     string `json:"orderNumber"`
	ResultingTrades []ResultTrades
}

func (p *Poloniex) Buy(market string, price, amount float64) (buy Buy, err error) {
	parameters := map[string]string{
		"currencyPair": strings.ToUpper(market),
		"rate":         strconv.FormatFloat(price, 'f', 8, 64),
		"amount":       strconv.FormatFloat(amount, 'f', 8, 64),
	}

	respCh := make(chan []byte)
	errCh := make(chan error)

	go p.tradingRequest("buy", parameters, respCh, errCh)

	resp := <-respCh
	err = <-errCh

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &buy)

	_ = p.observer.Observe("buy", parameters["currencyPair"], buy.OrderNumber)

	return
}

type Sell Buy

func (p *Poloniex) Sell(market string, price, amount float64) (sell Sell, err error) {
	parameters := map[string]string{
		"currencyPair": strings.ToUpper(market),
		"rate":         strconv.FormatFloat(price, 'f', 8, 64),
		"amount":       strconv.FormatFloat(amount, 'f', 8, 64),
	}

	respCh := make(chan []byte)
	errCh := make(chan error)

	go p.tradingRequest("sell", parameters, respCh, errCh)

	resp := <-respCh
	err = <-errCh

	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &sell)

	_ = p.observer.Observe("buy", parameters["currencyPair"], sell.OrderNumber)

	return
}
