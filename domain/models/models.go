package models

import "github.com/BullionBear/crypto-trade/api/gen/feed"

type Kline struct {
	OpenTime                 int64
	Open                     float64
	High                     float64
	Low                      float64
	Close                    float64
	Volume                   float64
	CloseTime                int64
	QuoteAssetVolume         float64
	TradeNum                 int64
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
}

func NewKlineFromPb(pbKline *feed.Kline) *Kline {
	return &Kline{
		OpenTime:                 pbKline.OpenTime,
		Open:                     pbKline.Open,
		High:                     pbKline.High,
		Low:                      pbKline.Low,
		Close:                    pbKline.Close,
		Volume:                   pbKline.Volume,
		CloseTime:                pbKline.CloseTime,
		QuoteAssetVolume:         pbKline.QuoteAssetVolume,
		TradeNum:                 pbKline.TradeNum,
		TakerBuyBaseAssetVolume:  pbKline.TakerBuyBaseAssetVolume,
		TakerBuyQuoteAssetVolume: pbKline.TakerBuyQuoteAssetVolume,
	}
}
