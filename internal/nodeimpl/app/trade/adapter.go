package trade

type SubscribeTradeAdapter interface {
	Subscribe(market string, symbol string)
}
