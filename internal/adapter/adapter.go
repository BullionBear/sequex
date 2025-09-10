package adapter

import (
	"fmt"

	"github.com/BullionBear/sequex/internal/model/sqx"
)

var (
	TradeAdapterMap = make(map[sqx.Exchange]TradeAdapter)
)

type TradeCallback func(trade sqx.Trade) error

// type DepthCallback func(depth sqx.Depth) error

type TradeAdapter interface {
	Subscribe(symbol sqx.Symbol, instrumentType sqx.InstrumentType, callback TradeCallback) (func(), error)
}

func CreateTradeAdapter(exchange sqx.Exchange) (TradeAdapter, error) {
	if _, ok := TradeAdapterMap[exchange]; !ok {
		return nil, fmt.Errorf("adapter not found for exchange: %s", exchange)
	}
	if _, ok := TradeAdapterMap[exchange]; !ok {
		return nil, fmt.Errorf("adapter not found for exchange: %s", exchange)
	}
	return TradeAdapterMap[exchange], nil
}

func RegisterTradeAdapter(exchange sqx.Exchange, adapter TradeAdapter) {
	if _, ok := TradeAdapterMap[exchange]; !ok {
		TradeAdapterMap[exchange] = adapter
	}
}
