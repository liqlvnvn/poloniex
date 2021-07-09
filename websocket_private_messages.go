package poloniex

import "time"

// List of constants for parsing poloniex messages.
const (
	OrderTypeSell = 0
	OrderTypeBuy  = 1

	OrderTypeSellValue = "sell"
	OrderTypeBuyValue  = "buy"

	OrderTypeFill      = "f"
	OrderTypeSelfTrade = "s"
	OrderTypeCanceled  = "c"

	MessageTypePending     = "Pending"
	MessageTypeOrderUpdate = "OrderUpdate"
	MessageTypeTrade       = "Trade"
	MessageTypeBalance     = "BalanceUpdate"
	MessageTypeNewOrder    = "NewOrder"
	MessageTypeMargin      = "MarginPositionUpdate"
	MessageTypeKill        = "Kill"

	WalletTypeExchange = "e"
	WalletTypeMargin   = "m"
	WalletTypeLending  = "l"
)

// Pending represent "p" messages which an acknowledgement that an order is pending.
// ["p", <order number>, <currency pair id>, "<rate>", "<amount>", "<order type>", "<clientOrderId>", "<epoch_ms>"]
type Pending struct {
	OrderNumber    string
	CurrencyPairID string
	Rate           float64
	Amount         float64
	OrderType      string
	ClientOrderID  string
	EpochMS        string
}

// BalanceUpdate represent "b" messages - an available balance update.
// ["b", 28, "e", "-0.06000000", "0.94000000"]
// The wallet can be e (exchange), m (margin), or l (lending).
type BalanceUpdate struct {
	CurrencyID string
	Wallet     string
	Amount     float64
	Balance    float64
}

// NewOrder represent "n" updates - a newly created limit order.
// ["n", 148, 6083059, 1, "0.03000000", "2.00000000", "2018-09-08 04:54:09", "2.00000000", "12345"]
// OrderType type can either be 0 (sell) or 1 (buy)
type NewOrder struct {
	CurrencyPairID        string
	OrderNumber           string
	OrderType             string
	Rate                  float64
	Amount                float64
	Date                  string
	OriginalAmountOrdered float64
	ClientOrderID         string
}

// OrderUpdate represent "o" messages.
// OrderType is one of: f, s, or c, corresponding to a fill, self-trade, or canceled order,
type OrderUpdate struct {
	OrderNumber    string
	NewAmount      float64
	OrderType      string
	ClientOrderID  string
	CanceledAmount float64
}

// MarginPositionUpdate represent "m" messages.
type MarginPositionUpdate struct {
	OrderNumber   string
	Currency      string
	Amount        float64
	ClientOrderID string
}

// Trade represent "t" messages - a trade notification.
// The funding type represents the funding used for the trade,
// which may be 0 (exchange wallet), 1 (borrowed funds), 2 (margin funds), or 3 (lending funds).
type Trade struct {
	TradeID       string    `json:"tradeID"`
	Rate          float64   `json:"rate"`
	Amount        float64   `json:"amount"`
	FeeMultiplier float64   `json:"feeMultiplier"`
	FundingType   string    `json:"fundingType"`
	OrderNumber   string    `json:"orderNumber"`
	TotalFee      float64   `json:"totalFee"`
	Date          time.Time `json:"date"`
	ClientOrderID string    `json:"clientOrderID"`
	TradeTotal    float64   `json:"tradeTotal"`
	EpochMS       string    `json:"epochMS"`
}

// Kill represent "k" messages which indicating that an API order has been killed,
// due to specified constraints not being matched.
// A postOnly or fillOrKill order that doesn't successfully execute will generate a k message.
type Kill struct {
	OrderNumber   string
	ClientOrderID string
}

type Fill struct {
	OrderID  string
	TradeID  string
	Symbol   string
	Price    float64
	Size     float64
	Side     string
	FilledAt time.Time
}
