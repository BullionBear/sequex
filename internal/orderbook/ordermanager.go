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

type BinanceOrderManager struct {
	SpotOrderBook map[string]*BinanceOrderBook
	PerpOrderBook map[string]*BinancePerpOrderBook
}

func NewBinanceOrderManager() *BinanceOrderManager {
	return &BinanceOrderManager{
		SpotOrderBook: make(map[string]*BinanceOrderBook),
		PerpOrderBook: make(map[string]*BinancePerpOrderBook),
	}
}

/*
	func getTickSize(symbol string, instrumentType InstrumentType) (decimal.Decimal, error) {
		switch instrumentType {
		case Spot:
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
			return decimal.Decimal{}, errors.New("unable to find tick size")
		case Perpetual:
			client := futures.NewClient("", "")
			info, err := client.NewExchangeInfoService().Do(context.Background())
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
			return decimal.Decimal{}, errors.New("unable to find tick size")
		default:
			log.Printf("Unsupported instrument type: %d", instrumentType)
			return decimal.Decimal{}, errors.New("unsupported instrument type")
		}
	}
*/
func (bom *BinanceOrderManager) CreateSpotBook(symbol string, updateSpeed UpdateSpeed) error {

	if _, exists := bom.SpotOrderBook[symbol]; !exists {
		bom.SpotOrderBook[symbol] = NewBinanceOrderBook(symbol, 500)
		go bom.SpotOrderBook[symbol].Run(updateSpeed)
	}
	return nil
}

func (bom *BinanceOrderManager) CloseSpotBook(symbol string) {
	if ob, exists := bom.SpotOrderBook[symbol]; exists {
		ob.Close()
		delete(bom.SpotOrderBook, symbol)
	}
}

func (bom *BinanceOrderManager) GetSpotOrderBook(symbol string, depth int) (*OrderBook, error) {
	if ob, exists := bom.SpotOrderBook[symbol]; exists {
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

func (bom *BinanceOrderManager) CreatePerpBook(symbol string, updateSpeed UpdateSpeed) error {
	if _, exists := bom.PerpOrderBook[symbol]; !exists {
		bom.PerpOrderBook[symbol] = NewBinancePerpOrderBook(symbol, 500)
		go bom.PerpOrderBook[symbol].Run(updateSpeed)
	}
	return nil
}

func (bom *BinanceOrderManager) ClosePerpBook(symbol string) {
	if ob, exists := bom.PerpOrderBook[symbol]; exists {
		ob.Close()
		delete(bom.PerpOrderBook, symbol)
	}
}

func (bom *BinanceOrderManager) GetPerpOrderBook(symbol string, depth int) (*OrderBook, error) {
	if ob, exists := bom.PerpOrderBook[symbol]; exists {
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
