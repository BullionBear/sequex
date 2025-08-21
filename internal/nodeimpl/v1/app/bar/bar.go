package bar

import "github.com/BullionBear/sequex/internal/nodeimpl/v1/app/share"

type Bar struct {
	Symbol        share.Symbol
	Instrument    share.Instrument
	Exchange      share.Exchange
	StartSeq      int64
	EndSeq        int64
	NextSeq       int64
	StartTime     int64
	EndTime       int64
	Open          float64
	High          float64
	Low           float64
	Close         float64
	Mean          float64
	Std           float64
	Median        float64
	FirstQuartile float64
	ThirdQuartile float64
	VolumeBase    float64
	VolumeQuote   float64
	Count         int64
}
