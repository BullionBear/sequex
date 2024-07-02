package strategy

import (
	"github.com/BullionBear/crypto-trade/domain/alpha"
	"github.com/BullionBear/crypto-trade/domain/models"
	"github.com/BullionBear/crypto-trade/domain/reporter"
	"github.com/BullionBear/crypto-trade/domain/wallet"
	"github.com/sirupsen/logrus"
)

type Strategy struct {
	Wallet   wallet.Wallet
	Alpha    alpha.Alpha
	Reporter reporter.Reporter
}

func NewStrategy(wallet *wallet.Wallet, alpha *alpha.Alpha, reporter *reporter.Reporter) *Strategy {
	return &Strategy{
		Wallet:   *wallet,
		Alpha:    *alpha,
		Reporter: *reporter,
	}
}

func (s *Strategy) Prepare(kline *models.Kline) {
	s.Alpha.Append(kline)
}

func (s *Strategy) MakeDecision(kline *models.Kline) {
	s.Alpha.Append(kline)
	lm := s.Alpha.LongCloseMovingAvg.Mean()
	sm := s.Alpha.ShortCloseMovingAvg.Mean()
	if lm > sm {
		logrus.Infof("Long moving average: %f, Short moving average: %f", lm, sm)
	}
}
