package tradingpipe

type TradingPipeline struct {
	name string
}

func NewTradingPipeline(name string) *TradingPipeline {
	return &TradingPipeline{
		name: name,
	}
}

func (t *TradingPipeline) OnKLineUpdate(event interface{}) {
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
