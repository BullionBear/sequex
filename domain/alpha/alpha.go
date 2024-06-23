package alpha

import (
	"github.com/BullionBear/crypto-trade/api/gen/models"
)

type Alpha struct {
	// 5 min moving average
	ShortCloseMovingAvg *MovingAverage

	// 1 hour term moving average
	LongCloseMovingAvg *MovingAverage
}

func NewAlpha() *Alpha {
	return &Alpha{
		ShortCloseMovingAvg: NewMovingAverage(300),
		LongCloseMovingAvg:  NewMovingAverage(3600),
	}
}

func (a *Alpha) Append(kline models.Kline) {
	a.shortCloseMovingAvg.Append(kline.Close)
	a.longCloseMovingAvg.Append(kline.Close)
}
