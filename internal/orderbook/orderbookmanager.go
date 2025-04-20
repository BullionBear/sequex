package orderbook

import (
	"errors"

	evbus "github.com/asaskevich/EventBus"
)

type InstrumentType int

const (
	Spot InstrumentType = iota
	Perpetual
)

type OrderBook struct {
	Symbol       string       `json:"symbol"`
	Asks         []PriceLevel `json:"asks"`
	Bids         []PriceLevel `json:"bids"`
	LastUpdateID int64        `json:"lastUpdateId"`
	Timestamp    int64        `json:"timestamp"`
}

type BinanceOrderBookManager struct {
	OrderBooks map[string]*BinanceOrderBook
	eventBus   evbus.Bus
}

func NewBinanceOrderBookManager() *BinanceOrderBookManager {
	return &BinanceOrderBookManager{
		OrderBooks: make(map[string]*BinanceOrderBook),
		eventBus:   evbus.New(),
	}
}

func (bom *BinanceOrderBookManager) CreateOrderBook(symbol string, spd UpdateSpeed) error {
	if _, exists := bom.OrderBooks[symbol]; !exists {
		bom.OrderBooks[symbol] = NewBinanceOrderBook(symbol, 500)
		if err := bom.OrderBooks[symbol].Start(spd); err != nil {
			return err
		}
	}
	return nil
}

func (bom *BinanceOrderBookManager) CloseOrderBook(symbol string) {
	if ob, exists := bom.OrderBooks[symbol]; exists {
		ob.Close()
		delete(bom.OrderBooks, symbol)
	}
}

func (bom *BinanceOrderBookManager) GetOrderBook(symbol string, depth int) (*OrderBook, error) {
	if ob, exists := bom.OrderBooks[symbol]; exists {
		ask, bid, err := ob.GetDepth(depth)
		if err != nil {
			return nil, err
		}
		return &OrderBook{
			Symbol:       symbol,
			Asks:         ask,
			Bids:         bid,
			LastUpdateID: ob.lastUpdateID,
			Timestamp:    ob.timestamp,
		}, nil
	}
	return nil, errors.New("order book not found")
}

type BinancePerpOrderBookManager struct {
	OrderBooks map[string]*BinancePerpOrderBook
}

func NewBinancePerpOrderBookManager() *BinancePerpOrderBookManager {
	return &BinancePerpOrderBookManager{
		OrderBooks: make(map[string]*BinancePerpOrderBook),
	}
}

func (bom *BinancePerpOrderBookManager) CreateOrderBook(symbol string, spd UpdateSpeed) error {
	if _, exists := bom.OrderBooks[symbol]; !exists {
		bom.OrderBooks[symbol] = NewBinancePerpOrderBook(symbol, 500)
		if err := bom.OrderBooks[symbol].Start(spd); err != nil {
			return err
		}
	}
	return nil
}

func (bom *BinancePerpOrderBookManager) CloseOrderBook(symbol string) {
	if ob, exists := bom.OrderBooks[symbol]; exists {
		ob.Close()
		delete(bom.OrderBooks, symbol)
	}
}

func (bom *BinancePerpOrderBookManager) GetOrderBook(symbol string, depth int) (*OrderBook, error) {
	if ob, exists := bom.OrderBooks[symbol]; exists {
		ask, bid, err := ob.GetDepth(depth)
		if err != nil {
			return nil, err
		}
		return &OrderBook{
			Symbol:       symbol,
			Asks:         ask,
			Bids:         bid,
			LastUpdateID: ob.lastUpdateID,
			Timestamp:    ob.timestamp,
		}, nil
	}
	return nil, errors.New("order book not found")
}
