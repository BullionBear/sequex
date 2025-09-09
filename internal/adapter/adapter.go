package adapter

import (
	"fmt"

	"github.com/BullionBear/sequex/internal/model/sqx"
)

var (
	AdapterMap = make(map[sqx.Exchange]map[sqx.DataType]Adapter)
)

type Callback func(data []byte) error

type Adapter interface {
	Subscribe(symbol sqx.Symbol, instrumentType sqx.InstrumentType, callback Callback) (func(), error)
}

func CreateAdapter(exchange sqx.Exchange, dataType sqx.DataType) (Adapter, error) {
	if _, ok := AdapterMap[exchange]; !ok {
		return nil, fmt.Errorf("adapter not found for exchange: %s", exchange)
	}
	if _, ok := AdapterMap[exchange][dataType]; !ok {
		return nil, fmt.Errorf("adapter not found for exchange: %s and data type: %s", exchange, dataType)
	}
	return AdapterMap[exchange][dataType], nil
}

func RegisterAdapter(exchange sqx.Exchange, dataType sqx.DataType, adapter Adapter) {
	if _, ok := AdapterMap[exchange]; !ok {
		AdapterMap[exchange] = make(map[sqx.DataType]Adapter)
	}
	AdapterMap[exchange][dataType] = adapter
}
