package market

type SubscribeTradeAdapter interface {
	Subscribe(market string, symbol string)
}
