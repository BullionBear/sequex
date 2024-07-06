package alpha

import (
	"github.com/BullionBear/crypto-trade/domain/models"
)

const (
	Min = 60
	Hr  = 60 * Min
)

type Alpha struct {
	CloseMovingAvg5Min  *MovingAverage
	CloseMovingAvg30Min *MovingAverage
	CloseMovingAvg3Hr   *MovingAverage
	CloseMovingAvg18Hr  *MovingAverage

	VolumeMovingAvg5Min  *MovingAverage
	VolumeMovingAvg30Min *MovingAverage
	VolumeMovingAvg3Hr   *MovingAverage
	VolumeMovingAvg18Hr  *MovingAverage
}

func NewAlpha() *Alpha {
	return &Alpha{
		CloseMovingAvg5Min:  NewMovingAverage(5 * Min),
		CloseMovingAvg30Min: NewMovingAverage(30 * Min),
		CloseMovingAvg3Hr:   NewMovingAverage(3 * Hr),
		CloseMovingAvg18Hr:  NewMovingAverage(18 * Hr),

		VolumeMovingAvg5Min:  NewMovingAverage(5 * Min),
		VolumeMovingAvg30Min: NewMovingAverage(30 * Min),
		VolumeMovingAvg3Hr:   NewMovingAverage(3 * Hr),
		VolumeMovingAvg18Hr:  NewMovingAverage(18 * Hr),
	}
}

func (a *Alpha) Append(kline *models.Kline) {
	a.CloseMovingAvg5Min.Append(kline.Close)
	a.CloseMovingAvg30Min.Append(kline.Close)
	a.CloseMovingAvg3Hr.Append(kline.Close)
	a.CloseMovingAvg18Hr.Append(kline.Close)

	a.VolumeMovingAvg5Min.Append(kline.QuoteAssetVolume)
	a.VolumeMovingAvg30Min.Append(kline.QuoteAssetVolume)
	a.VolumeMovingAvg3Hr.Append(kline.QuoteAssetVolume)
	a.VolumeMovingAvg18Hr.Append(kline.QuoteAssetVolume)
}
