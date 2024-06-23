package models

type Kline struct {
	StartTime            int64
	EndTime              int64
	Symbol               string
	Interval             string
	FirstTradeID         int64
	LastTradeID          int64
	Open                 float64
	Close                float64
	High                 float64
	Low                  float64
	Volume               float64
	TradeNum             int64
	QuoteVolume          float64
	ActiveBuyVolume      float64
	ActiveBuyQuoteVolume float64
}


// Define a kline handler
KlineHandler func (kline Kline) 