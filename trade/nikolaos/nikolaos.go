package nikolaos

import (
	"math"

	"github.com/BullionBear/crypto-trade/domain/alpha"
	"github.com/BullionBear/crypto-trade/domain/chronicler"
	"github.com/BullionBear/crypto-trade/domain/models"
	"github.com/BullionBear/crypto-trade/domain/wallet"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

type Nikolaos struct {
	Wallet     *wallet.Wallet
	Alpha      *alpha.Alpha
	Chronicler *chronicler.Chronicler
}

func NewNikolaos(wallet *wallet.Wallet, alpha *alpha.Alpha, chronicler *chronicler.Chronicler) *Nikolaos {
	return &Nikolaos{
		Wallet:     wallet,
		Alpha:      alpha,
		Chronicler: chronicler,
	}
}

func (niko *Nikolaos) Prepare(kline *models.Kline) {
	niko.Alpha.Append(kline)
}

func (niko *Nikolaos) MakeDecision(kline *models.Kline) {
	niko.Alpha.Append(kline)
	lm := niko.Alpha.LongCloseMovingAvg.Mean()
	sm := niko.Alpha.ShortCloseMovingAvg.Mean()

	data := bson.M{
		"long_moving_avg":  lm,
		"short_moving_avg": sm,
	}
	BTCAmount, _ := niko.Wallet.GetBalance("BTC").Float64()
	USDTAmount, _ := niko.Wallet.GetBalance("USDT").Float64()
	wallet := bson.M{
		"BTC":  BTCAmount,
		"USDT": USDTAmount,
	}
	history := chronicler.NewHistory(kline.OpenTime, data, wallet)
	niko.Chronicler.Record(history)

	indicator := (lm - sm) / lm
	if math.Abs(indicator) > 0.05 {
		logrus.Infof("Long moving average: %f, Short moving average: %f, indicator: %f", lm, sm, indicator)
	}
}
