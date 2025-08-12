package market

import "github.com/BullionBear/sequex/pkg/node"

type PublicTradeConfig struct {
	Symbol string `json:"symbol"`
	Market string `json:"market"`
}

type PublicTradeNode struct {
	*node.BaseNode
	// Configurable parameters
	cfg PublicTradeConfig

	SubscribeTrade(market string, symbol string, options TradeSubscriptionOptions)
}

