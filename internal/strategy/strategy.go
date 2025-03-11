package strategy

type IStrategy interface {
	OnMarketUpdate()
	OnSignal()
	OnExecute()
}
