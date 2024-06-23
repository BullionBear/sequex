package trader

import "github.com/shopspring/decimal"

type Backtest struct {
}

func NewBacktest() *Backtest {
}

func (b *Backtest) CreateMarketOrder(symbol string, quoteQty decimal.Decimal) error {

}
