package orderbook

import (
	"context"
	"errors"

	"log"

	"github.com/adshao/go-binance/v2"
	"github.com/shopspring/decimal"
)

type OrderBook struct {
	Symbol       string       `json:"symbol"`
	Asks         []PriceLevel `json:"asks"`
	Bids         []PriceLevel `json:"bids"`
	LastUpdateID int64        `json:"lastUpdateId"`
	Timestamp    int64        `json:"timestamp"`
}

type BinanceOrderManager struct {
	OrderBook map[string]*BinanceOrderBook
}

func NewBinanceOrderManager() *BinanceOrderManager {
	return &BinanceOrderManager{
		OrderBook: make(map[string]*BinanceOrderBook),
	}
}

func getTickSize(symbol string) (decimal.Decimal, error) {
	client := binance.NewClient("", "")
	info, err := client.NewExchangeInfoService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return decimal.Decimal{}, err
	}
	// Extract the number of price decimals for the symbol
	for _, s := range info.Symbols {
		if s.Symbol == symbol {
			for _, filter := range s.Filters {
				if filter["filterType"] == "PRICE_FILTER" {
					priceTickSize := filter["tickSize"].(string)
					dTickSize, err := decimal.NewFromString(priceTickSize)
					if err != nil {
						log.Printf("Error parsing tick size: %+v", err)
						return decimal.Decimal{}, err
					}
					return dTickSize, nil
				}
			}
		}
	}
	return decimal.Decimal{}, nil
}

func (bom *BinanceOrderManager) CreateBook(symbol string) error {
	tickSize, err := getTickSize(symbol)
	if err != nil {
		log.Printf("Error getting tick size for symbol %s: %v", symbol, err)
	}
	if _, exists := bom.OrderBook[symbol]; !exists {
		bom.OrderBook[symbol] = NewBinanceOrderBook(symbol, tickSize, 1000)
		go bom.OrderBook[symbol].Run()
	}
	return nil
}

func (bom *BinanceOrderManager) CloseBook(symbol string) {
	if ob, exists := bom.OrderBook[symbol]; exists {
		ob.Close()
		delete(bom.OrderBook, symbol)
	}
}

func (bom *BinanceOrderManager) GetOrderBook(symbol string, depth int) (*OrderBook, error) {
	if ob, exists := bom.OrderBook[symbol]; exists {
		return &OrderBook{
			Symbol:       ob.Symbol,
			Asks:         ob.Asks.GetBook(depth),
			Bids:         ob.Bids.GetBook(depth),
			LastUpdateID: ob.lastUpdateID,
			Timestamp:    ob.timestamp,
		}, nil
	}
	return nil, errors.New("order book not found")
}
