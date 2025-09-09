package adapter

import (
	"fmt"

	"github.com/BullionBear/sequex/internal/model/sqx"
)

var (
	TradeAdapterMap = make(map[sqx.Exchange]map[sqx.DataType]TradeAdapter)
	DepthAdapterMap = make(map[sqx.Exchange]map[sqx.DataType]DepthAdapter)
)

type TradeCallback func(trade sqx.Trade) error

// type DepthCallback func(depth sqx.Depth) error

type TradeAdapter interface {
	Subscribe(symbol sqx.Symbol, instrumentType sqx.InstrumentType, callback TradeCallback) (func(), error)
}

func CreateTradeAdapter(exchange sqx.Exchange, dataType sqx.DataType) (TradeAdapter, error) {
	if _, ok := TradeAdapterMap[exchange]; !ok {
		return nil, fmt.Errorf("adapter not found for exchange: %s", exchange)
	}
	if _, ok := TradeAdapterMap[exchange][dataType]; !ok {
		return nil, fmt.Errorf("adapter not found for exchange: %s and data type: %s", exchange, dataType)
	}
	return TradeAdapterMap[exchange][dataType], nil
}

func RegisterTradeAdapter(exchange sqx.Exchange, dataType sqx.DataType, adapter TradeAdapter) {
	if _, ok := TradeAdapterMap[exchange]; !ok {
		TradeAdapterMap[exchange] = make(map[sqx.DataType]TradeAdapter)
	}
	TradeAdapterMap[exchange][dataType] = adapter
}
