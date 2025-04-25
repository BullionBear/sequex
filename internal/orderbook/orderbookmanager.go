package orderbook

import (
	"errors"

	"github.com/BullionBear/sequex/pkg/log"
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
	logger     *log.Logger
}

func NewBinanceOrderBookManager(logger *log.Logger) *BinanceOrderBookManager {
	return &BinanceOrderBookManager{
		OrderBooks: make(map[string]*BinanceOrderBook),
		eventBus:   evbus.New(),
		logger:     logger,
	}
}

func (bom *BinanceOrderBookManager) CreateOrderBook(symbol string, spd UpdateSpeed) error {
	if _, exists := bom.OrderBooks[symbol]; !exists {
		bom.OrderBooks[symbol] = NewBinanceOrderBook(symbol, 500, bom.logger.WithKV(log.KV{Key: "orderbook", Value: symbol}))
		if err := bom.OrderBooks[symbol].Start(spd); err != nil {
			return err
		}
		bom.logger.Info("Registering OrderBook for %s", symbol)
		bom.OrderBooks[symbol].SubscribeBestDepth(func(ask, bid PriceLevel) {
			bom.logger.Info("Publishing Ask(%v)/Bid(%v)", ask, bid)
			bom.eventBus.Publish(bom.channelName(symbol), ask, bid)
		})
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

func (bom *BinanceOrderBookManager) SubscribeBestDepth(symbol string, callback func(ask, bid PriceLevel)) (func(), error) {
	chName := bom.channelName(symbol)
	bom.logger.Info("Subscribing to channel: %s\n", chName)
	if err := bom.eventBus.SubscribeAsync(chName, callback, false); err != nil {
		return nil, err
	}
	return func() {
		if err := bom.eventBus.Unsubscribe(chName, callback); err != nil {
			bom.logger.Info("Failed to unsubscribe from channel: %s\n", chName)
		}
	}, nil
}

func (bom *BinanceOrderBookManager) UnsubscribeBestDepth(symbol string, callback func(ask, bid PriceLevel)) error {
	if err := bom.eventBus.Unsubscribe(bom.channelName(symbol), callback); err != nil {
		return err
	}
	return nil
}

func (bom *BinanceOrderBookManager) channelName(symbol string) string {
	return symbol + ":depth1"
}

type BinancePerpOrderBookManager struct {
	OrderBooks map[string]*BinancePerpOrderBook
	logger     *log.Logger
}

func NewBinancePerpOrderBookManager(logger *log.Logger) *BinancePerpOrderBookManager {
	return &BinancePerpOrderBookManager{
		OrderBooks: make(map[string]*BinancePerpOrderBook),
		logger:     logger,
	}
}

func (bom *BinancePerpOrderBookManager) CreateOrderBook(symbol string, spd UpdateSpeed) error {
	if _, exists := bom.OrderBooks[symbol]; !exists {
		bom.OrderBooks[symbol] = NewBinancePerpOrderBook(symbol, 500, bom.logger.WithKV(log.KV{Key: "orderbookperp", Value: symbol}))
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
