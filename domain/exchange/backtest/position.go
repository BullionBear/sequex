package backtest

import (
	"errors"

	"github.com/shopspring/decimal"
)

/*
	Position is regarding as a unified position for manage perpetuals.
*/

var (
	ErrPositionNotFound = errors.New("position not found")
)

type Perpetual struct {
	Id        int
	Symbol    string
	OpenTime  int64
	BaseQty   decimal.Decimal
	OpenPrice decimal.Decimal
	Leverage  int
}

func NewPerpetual(id int, symbol string, openTime int64, baseQty, openPrice decimal.Decimal, leverage int) *Perpetual {
	return &Perpetual{
		Id:        id,
		Symbol:    symbol,
		OpenTime:  openTime,
		BaseQty:   baseQty,
		OpenPrice: openPrice,
		Leverage:  leverage,
	}
}

type Position struct {
	Id              int
	PerpetualMap    map[int]*Perpetual
	PerpetualGroups map[string][]*Perpetual
}

func NewPosition() *Position {
	return &Position{
		Id:              0,
		PerpetualMap:    make(map[int]*Perpetual),
		PerpetualGroups: make(map[string][]*Perpetual),
	}
}

func (pos *Position) OpenPosition(symbol string, openTime int64, baseQty, openPrice decimal.Decimal, leverage int) int {
	if _, ok := pos.PerpetualGroups[symbol]; !ok {
		pos.PerpetualGroups[symbol] = make([]*Perpetual, 0)
	}
	pos.PerpetualMap[pos.Id] = NewPerpetual(pos.Id, symbol, openTime, baseQty, openPrice, leverage)
	pos.PerpetualGroups[symbol] = append(pos.PerpetualGroups[symbol], pos.PerpetualMap[pos.Id])
	Id := pos.Id
	pos.Id++
	return Id
}

func (pos *Position) ClosePosition(Id int) (*Perpetual, error) {
	if _, ok := pos.PerpetualMap[Id]; !ok {
		return nil, ErrPositionNotFound
	}
	symbol := pos.PerpetualMap[Id].Symbol
	delete(pos.PerpetualMap, Id)
	for i, p := range pos.PerpetualGroups[symbol] {
		if p.Id == Id {
			pos.PerpetualGroups[symbol] = append(pos.PerpetualGroups[symbol][:i], pos.PerpetualGroups[symbol][i+1:]...)
			return p, nil
		}
	}
	return nil, ErrPositionNotFound
}
