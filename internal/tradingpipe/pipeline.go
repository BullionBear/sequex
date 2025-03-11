package tradingpipe

import (
	"github.com/BullionBear/sequex/internal/metadata"
	"github.com/BullionBear/sequex/internal/strategy"
)

type TradingPipeline struct {
	name string
	st   strategy.Strategy
}

func NewTradingPipeline(name string, strategy strategy.Strategy) *TradingPipeline {
	return &TradingPipeline{
		name: name,
		st:   strategy,
	}
}

func (t *TradingPipeline) OnKLineUpdate(klineUpdate metadata.KLineUpdate) {
	t.st.OnKLineUpdate(klineUpdate)
}

func (t *TradingPipeline) Status() string {
	return ""
}

func (t *TradingPipeline) Name() string {
	return t.name
}
