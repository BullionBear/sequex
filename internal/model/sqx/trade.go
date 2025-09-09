package sqx

import "github.com/BullionBear/sequex/internal/model/protobuf"

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

func (t *Trade) ToProtobuf() *protobuf.Trade {
	symbol := t.Symbol.ToProtobuf()
	return &protobuf.Trade{
		Id:         t.Id,
		Symbol:     &symbol,
		Exchange:   t.Exchange.ToProtobuf(),
		Instrument: t.InstrumentType.ToProtobuf(),
		Side:       t.TakerSide.ToProtobuf(),
		Price:      t.Price,
		Quantity:   t.Quantity,
		Timestamp:  t.Timestamp,
	}
}
