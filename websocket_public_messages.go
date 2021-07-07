package poloniex

// WSTicker is for ticker update.
type WSTicker struct {
	Symbol        string  `json:"symbol"`
	Last          float64 `json:"last"`
	LowestAsk     float64 `json:"lowestAsk"`
	HighestBid    float64 `json:"highestBid"`
	PercentChange float64 `json:"percentChange"`
	BaseVolume    float64 `json:"baseVolume"`
	QuoteVolume   float64 `json:"quoteVolume"`
	IsFrozen      bool    `json:"isFrozen"`
	High24hr      float64 `json:"high24hr"`
	Low24hr       float64 `json:"low24hr"`
}

// OrderDepth is for "i" messages.
type OrderDepth struct {
	Symbol    string `json:"symbol"`
	OrderBook struct {
		Asks []Book `json:"asks"`
		Bids []Book `json:"bids"`
	} `json:"orderBook"`
}

// WSOrderBook is for "o" messages
type WSOrderBook struct {
	Rate      float64 `json:"rate,string"`
	TypeOrder string  `json:"type"`
	Amount    float64 `json:"amount,string"`
}

// WSOrderBookModify is for "o" messages.
type WSOrderBookModify WSOrderBook

// WSOrderBookRemove is for "o" messages.
type WSOrderBookRemove struct {
	Rate      float64 `json:"rate,string"`
	TypeOrder string  `json:"type"`
}

// NewTrade - "t" messages.
type NewTrade struct {
	TradeID   int64   `json:"tradeID,string"`
	Rate      float64 `json:"rate,string"`
	Amount    float64 `json:"amount,string"`
	Total     float64 `json:"total,string"`
	TypeOrder string  `json:"type"`
}
