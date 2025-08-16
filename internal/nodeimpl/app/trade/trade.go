package market

type PublicTradeConfig struct {
	Symbol string `json:"symbol"`
	Market string `json:"market"`
}

/*
type PublicTradeNode struct {
	*base.BaseNode
	// Configurable parameters
	cfg PublicTradeConfig

	SubscribeTrade(market string, symbol string, options TradeSubscriptionOptions)
}
*/
