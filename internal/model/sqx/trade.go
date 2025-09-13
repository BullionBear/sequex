package sqx

import (
	"fmt"

	"github.com/BullionBear/sequex/internal/model/protobuf"
	"google.golang.org/protobuf/proto"
)

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

func (t *Trade) FromProtobuf(trade *protobuf.Trade) error {
	t.Id = trade.Id
	t.Symbol = NewSymbol(trade.Symbol.Base, trade.Symbol.Quote)
	t.Exchange = NewExchange(trade.Exchange.String())
	if t.Exchange == ExchangeUnknown {
		return fmt.Errorf("unknown exchange: %s", trade.Exchange.String())
	}
	t.InstrumentType = NewInstrumentType(trade.Instrument.String())
	if t.InstrumentType == InstrumentTypeUnknown {
		return fmt.Errorf("unknown instrument type: %s", trade.Instrument.String())
	}
	t.TakerSide = NewSide(trade.Side.String())
	if t.TakerSide == SideUnknown {
		return fmt.Errorf("unknown taker side: %s", trade.Side.String())
	}
	t.Price = trade.Price
	t.Quantity = trade.Quantity
	t.Timestamp = trade.Timestamp
	return nil
}

func (t *Trade) Marshal() ([]byte, error) {
	return proto.Marshal(t.ToProtobuf())
}

func Unmarshal(data []byte, trade *Trade) error {
	pbTrade := &protobuf.Trade{}
	err := proto.Unmarshal(data, pbTrade)
	if err != nil {
		return err
	}
	err = trade.FromProtobuf(pbTrade)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trade) IdStr() string {

	return fmt.Sprintf("%s-%s-%s-%d", t.Exchange.String(), t.InstrumentType.String(), t.Symbol.String(), t.Id)
}
