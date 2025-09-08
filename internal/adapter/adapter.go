package adapter

import (
	"fmt"
)

var (
	AdapterMap = make(map[string]map[string]Adapter)
)

type Callback func(data []byte) error

type Adapter interface {
	Subscribe(callback Callback) (func(), error)
}

func CreateAdapter(exchange string, dataType string) (Adapter, error) {
	if _, ok := AdapterMap[exchange]; !ok {
		return nil, fmt.Errorf("adapter not found for exchange: %s", exchange)
	}
	if _, ok := AdapterMap[exchange][dataType]; !ok {
		return nil, fmt.Errorf("adapter not found for exchange: %s and data type: %s", exchange, dataType)
	}
	return AdapterMap[exchange][dataType], nil
}

func RegisterAdapter(exchange string, dataType string, adapter Adapter) {
	AdapterMap[exchange][dataType] = adapter
}
