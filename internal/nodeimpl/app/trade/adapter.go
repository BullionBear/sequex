package trade

import (
	"fmt"

	"github.com/BullionBear/sequex/internal/nodeimpl/app/share"
)

var adapterMap = map[share.Exchange]map[share.Instrument]SubscribeTradeAdapter{}

func CreateAdapter(exchange share.Exchange, instrument share.Instrument) (SubscribeTradeAdapter, error) {
	if adapter, ok := adapterMap[exchange][instrument]; ok {
		return adapter, nil
	}
	return nil, fmt.Errorf("adapter not found")
}

func RegisterAdapter(exchange share.Exchange, instrument share.Instrument, adapter SubscribeTradeAdapter) {
	if _, ok := adapterMap[exchange]; !ok {
		adapterMap[exchange] = make(map[share.Instrument]SubscribeTradeAdapter)
	}
	adapterMap[exchange][instrument] = adapter
}

type Trade struct {
	Symbol    share.Symbol
	ID        int64
	Price     float64
	Qty       float64
	Time      int64
	TakerSide share.Side
}

type TradeSubscriptionOptions struct {
	OnConnect    func()            // Called when connection is established
	OnReconnect  func()            // Called when connection is reestablished
	OnError      func(err error)   // Called when an error occurs
	OnTrade      func(trade Trade) // Called when trade data is received
	OnDisconnect func()            // Called when connection is disconnected
}

type SubscribeTradeAdapter interface {
	Subscribe(symbol share.Symbol, options TradeSubscriptionOptions) (func(), error)
}
