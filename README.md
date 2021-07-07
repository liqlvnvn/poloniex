# Poloniex Go

Poloniex Websocket, Public and Private APIs.

## Related URL's

- [Poloniex API docs](https://docs.poloniex.com/)
- [Confluence]()

## Websocket Private
Create websocket client.
#### NewAuthenticatedWSClient()
~~~go
ws := poloniex.NewPrivateWSClient(apiKey, apiSecret)
err := ws.Run()
if err != nil {
    return
}
~~~
* Methods
  * SubscribeAccount()
  * UnsubscribeAccount()
  
#### SubscribeAccount()
~~~go
err = ws.SubscribeAccount()
if err != nil {
    return
}
for {
    fmt.Println(<-ws.Subs["ACCOUNT"])
}
~~~
#### UnsubscribeAccount()
~~~go
err = ws.SubscribeAccount()
go func() {
    time.Sleep(time.Second * 10)
    ws.UnsubscribeAccount()
}()
~~~
#### ListeningReports()
~~~go
ch, _ := ws.ListeningReports()
for {
    fmt.Println(<-ch)
}
~~~

### Examples
* See `./example/ws_private`

## Websocket Public
Create websocket client.
#### NewWSClient()
~~~go
ws := poloniex.NewPublicWSClient()
err := ws.Run()
if err != nil {
    return
}
~~~
* Push Api Methods
    * SubscribeTicker()
    * SubscribeMarket()
    * UnsubscribeTicker()
    * UnsubscribeMarket()
  
### Ticker
#### SubscribeTicker()
~~~go
err = ws.SubscribeTicker()
if err != nil {
    return
}
for {
    fmt.Println(<-ws.Subs["TICKER"])
}
~~~
#### UnsubscribeTicker()
~~~go
err = ws.SubscribeTicker()
go func() {
    time.Sleep(time.Second * 10)
    ws.UnsubscribeTicker()
}()
for {
    fmt.Println(<-ws.Subs["TICKER"])
}
~~~

### OrderDepth, OrderBook and Trades
#### SubscribeMarket()
~~~go
err = ws.SubscribeMarket("USDT_BTC")
if err != nil {
    return
}
for {
    fmt.Println(<-ws.Subs["USDT_BTC"])
}
~~~
#### UnsubscribeMarket()
~~~go
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
~~~~

### Examples
* See `./example/ws_public`

## Public Api
~~~go
poloniex, err := poloniex.NewClient(api_key, api_secret)
~~~
* Public Api Methods
    * GetTickers()
    * Get24hVolumes()
    * GetOrderBook()
    * GetPublicTradeHistory()
    * GetChartData()
    * GetCurrencies()
    * GetLoanOrders()
    
#### Example
~~~go
resp, err := poloniex.GetTickers()
if err != nil{
    panic(err)
}
fmt.Println(resp)
~~~
* See `./example/public_api`

## Private Api
~~~go
const (
        APIKey    = ""
        APISecret = ""
)
~~~
~~~go
poloniex, err := poloniex.NewClient(api_key, api_secret)
~~~ 

* Trading Api Methods
    * GetBalances()
    * GetCompleteBalances()
    * GetAccountBalances()
    * GetDepositAddresses()
    * GenerateNewAddress()
    * GetOpenOrders()
    * GetAllOpenOrders()
    * CancelOrder()
    * GetTradeHistory()
    * GetTradesByOrderID()
    * GetOrderStat()
    * Buy()
    * Sell()


#### Example
~~~go
resp, err := poloniex.Buy("btc_dgb", 0.00000099, 10000)
if err != nil{
    panic(err)
}
fmt.Println(resp)
~~~
* See `./example/private_api`
