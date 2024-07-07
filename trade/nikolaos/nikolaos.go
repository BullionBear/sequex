package nikolaos

import (
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

	beat int
}

func NewNikolaos(wallet *wallet.Wallet, alpha *alpha.Alpha, chronicler *chronicler.Chronicler) *Nikolaos {
	return &Nikolaos{
		Wallet:     wallet,
		Alpha:      alpha,
		Chronicler: chronicler,

		beat: 0,
	}
}

func (niko *Nikolaos) Prepare(kline *models.Kline) {

	niko.Alpha.Append(kline)
	niko.heartBeat("Load History")
}

func (niko *Nikolaos) MakeDecision(kline *models.Kline) {
	// Append the kline to the alpha
	niko.Alpha.Append(kline)
	// Retrieve the moving averages
	_ = niko.Alpha.CloseMovingAvg5Min.Mean()
	// closeMA30Min := niko.Alpha.CloseMovingAvg3Hr.Mean()
	// closeMA3Hr := niko.Alpha.CloseMovingAvg3Hr.Mean()
	// closeMA18Hr := niko.Alpha.CloseMovingAvg18Hr.Mean()

	// volumeMA5Min := niko.Alpha.VolumeMovingAvg5Min.Mean()
	// volumeMA30Min := niko.Alpha.VolumeMovingAvg30Min.Mean()
	// volumeMA3Hr := niko.Alpha.VolumeMovingAvg3Hr.Mean()
	// volumeMA18Hr := niko.Alpha.VolumeMovingAvg18Hr.Mean()

	// Record the data
	data := bson.M{
		"closeMA5Min":     niko.Alpha.CloseMovingAvg5Min.Mean(),
		"closeMa5MinStd":  niko.Alpha.CloseMovingAvg5Min.Std(),
		"closeMA30Min":    niko.Alpha.CloseMovingAvg30Min.Mean(),
		"closeMa30MinStd": niko.Alpha.CloseMovingAvg30Min.Std(),
		"closeMA3Hr":      niko.Alpha.CloseMovingAvg3Hr.Mean(),
		"closeMa3HrStd":   niko.Alpha.CloseMovingAvg3Hr.Std(),
		"closeMA18Hr":     niko.Alpha.CloseMovingAvg18Hr.Mean(),
		"closeMa18HrStd":  niko.Alpha.CloseMovingAvg18Hr.Std(),

		"volumeMA5Min":     niko.Alpha.VolumeMovingAvg5Min.Mean(),
		"volumeMa5MinStd":  niko.Alpha.VolumeMovingAvg5Min.Std(),
		"volumeMA30Min":    niko.Alpha.VolumeMovingAvg30Min.Mean(),
		"volumeMa30MinStd": niko.Alpha.VolumeMovingAvg30Min.Std(),
		"volumeMA3Hr":      niko.Alpha.VolumeMovingAvg3Hr.Mean(),
		"volumeMa3HrStd":   niko.Alpha.VolumeMovingAvg3Hr.Std(),
		"volumeMA18Hr":     niko.Alpha.VolumeMovingAvg18Hr.Mean(),
		"volumeMa18HrStd":  niko.Alpha.VolumeMovingAvg18Hr.Std(),
	}
	BTCAmount, _ := niko.Wallet.GetBalance("BTC").Float64()
	USDTAmount, _ := niko.Wallet.GetBalance("USDT").Float64()
	wallet := bson.M{
		"BTC":  BTCAmount,
		"USDT": USDTAmount,
	}
	history := chronicler.NewHistory(kline.OpenTime, kline.Close, data, wallet)
	niko.Chronicler.Record(history)

	niko.heartBeat("Make Desision")
}

func (niko *Nikolaos) heartBeat(name string) {
	if niko.beat%3600 == 0 {
		logrus.Infof("%s: Nikolaos 5Min price: %f", name, niko.Alpha.CloseMovingAvg5Min.Std())
	}
	niko.beat++
}
