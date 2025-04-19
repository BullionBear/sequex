package orderbook

import (
	"errors"
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
	OrderBooks map[InstrumentType]map[string]interface{}
}

func NewBinanceOrderBookManager() *BinanceOrderBookManager {
	return &BinanceOrderBookManager{
		OrderBooks: map[InstrumentType]map[string]interface{}{
			Spot:      make(map[string]interface{}),
			Perpetual: make(map[string]interface{}),
		},
	}
}

func (bom *BinanceOrderBookManager) CreateOrderBook(symbol string, instrumentType InstrumentType, updateSpeed UpdateSpeed) error {
	if _, exists := bom.OrderBooks[instrumentType][symbol]; !exists {
		if instrumentType == Spot {
			bom.OrderBooks[instrumentType][symbol] = NewBinanceOrderBook(symbol, 500)
			go bom.OrderBooks[instrumentType][symbol].(*BinanceOrderBook).Start(updateSpeed)
		} else if instrumentType == Perpetual {
			bom.OrderBooks[instrumentType][symbol] = NewBinancePerpOrderBook(symbol, 500)
			go bom.OrderBooks[instrumentType][symbol].(*BinancePerpOrderBook).Run(updateSpeed)
		}
	}
	return nil
}

func (bom *BinanceOrderBookManager) CloseOrderBook(symbol string, instrumentType InstrumentType) {
	if ob, exists := bom.OrderBooks[instrumentType][symbol]; exists {
		if instrumentType == Spot {
			ob.(*BinanceOrderBook).Close()
		} else if instrumentType == Perpetual {
			ob.(*BinancePerpOrderBook).Close()
		}
		delete(bom.OrderBooks[instrumentType], symbol)
	}
}

func (bom *BinanceOrderBookManager) GetOrderBook(symbol string, depth int, instrumentType InstrumentType) (*OrderBook, error) {
	if ob, exists := bom.OrderBooks[instrumentType][symbol]; exists {
		if instrumentType == Spot {
			spotBook := ob.(*BinanceOrderBook)
			return &OrderBook{
				Symbol:       spotBook.Symbol,
				Asks:         spotBook.Asks.GetBook(depth, true),
				Bids:         spotBook.Bids.GetBook(depth, false),
				LastUpdateID: spotBook.lastUpdateID,
				Timestamp:    spotBook.timestamp,
			}, nil
		} else if instrumentType == Perpetual {
			perpBook := ob.(*BinancePerpOrderBook)
			return &OrderBook{
				Symbol:       perpBook.Symbol,
				Asks:         perpBook.Asks.GetBook(depth, true),
				Bids:         perpBook.Bids.GetBook(depth, false),
				LastUpdateID: perpBook.lastUpdateID,
				Timestamp:    perpBook.timestamp,
			}, nil
		}
	}
	return nil, errors.New("order book not found")
}
