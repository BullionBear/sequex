package sqx

type Trade struct {
	Id             int64
	Symbol         Symbol
	Exchange       Exchange
	InstrumentType InstrumentType
	TakerSide      Side
	Price          float64
	Quantity       float64
	Timestamp      int64
}
