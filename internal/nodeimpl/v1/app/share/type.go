package share

type Symbol struct {
	Base  string
	Quote string
}

type Exchange string

const (
	ExchangeBinance     Exchange = "binance"
	ExchangeBinancePerp Exchange = "binance_perp"
	ExchangeBybit       Exchange = "bybit"
)

type Instrument string

const (
	InstrumentSpot Instrument = "spot"
	InstrumentPerp Instrument = "perp"
)

type Side string

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"
)
