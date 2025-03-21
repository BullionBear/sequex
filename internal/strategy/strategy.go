package strategy

type Strategy interface {
	OnKLineUpdate(symbol string, timestamp int64) error
}
