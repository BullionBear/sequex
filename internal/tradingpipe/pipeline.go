package tradingpipe

import (
	"github.com/BullionBear/sequex/internal/metadata"
	"github.com/BullionBear/sequex/internal/strategy"
)

type TradingPipeline struct {
	name string
	st   strategy.IStrategy
}

func NewTradingPipeline(name string, strategy strategy.IStrategy) *TradingPipeline {
	return &TradingPipeline{
		name: name,
		st:   strategy,
	}
}

func (t *TradingPipeline) OnKLineUpdate(klineUpdate metadata.KLineUpdate) {
	st.OnMarketUpdate(event)
}

func (t *TradingPipeline) Run() {
}

func (t *TradingPipeline) Stop() {
}

func (t *TradingPipeline) Destroy() {
}

func (t *TradingPipeline) Status() string {
	return ""
}

func (t *TradingPipeline) Name() string {
	return t.name
}
