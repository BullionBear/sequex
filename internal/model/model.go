package model

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

type Tick struct {
	TradeID int64
	Time    int64
	Price   float64
	IsValid bool
}
