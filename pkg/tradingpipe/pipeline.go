package tradingpipe

type TradingPipeline struct {
}

/*
func NewTradingPipeline() *TradingPipeline {
	bus := NewEventBus()

	pipeline := &TradingPipeline{
		bus:           bus,
		dataCollector: &DataCollector{bus: bus},
		strategy: &Strategy{
			bus:         bus,
			dataChannel: bus.Subscribe(DataEvent),
		},
		executor: &TradeExecutor{
			bus:           bus,
			signalChannel: bus.Subscribe(SignalEvent),
		},
		logger: &Logger{
			logChannel: bus.Subscribe(ExecutionEvent),
		},
	}

	return pipeline
}
*/
