package backtest

/*
	Account is regarding as a unified account for manage balance of assets.
*/

import (
	"errors"
	"strings"

	"github.com/shopspring/decimal"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type Asset struct {
	Coin   string
	Free   decimal.Decimal
	Locked decimal.Decimal
}

func NewAsset(coin string, free, locked decimal.Decimal) Asset {
	return Asset{
		Coin:   coin,
		Free:   free,
		Locked: locked,
	}
}

type Account struct {
	assets map[string]Asset
}

func NewAccount() *Account {
	return &Account{
		assets: make(map[string]Asset),
	}
}

func (acc *Account) SetBalance(coin string, balance decimal.Decimal) {
	acc.assets[coin] = NewAsset(coin, balance, decimal.Zero)
}

func (acc *Account) GetAsset(coin string) Asset {
	if b, ok := acc.assets[coin]; ok {
		return b
	}
	acc.assets[coin] = NewAsset(coin, decimal.Zero, decimal.Zero)
	return acc.assets[coin]
}

func (acc *Account) Swap(symbol string, side bool, price, baseQty decimal.Decimal) error {
	parts := strings.Split(symbol, "-")
	base := parts[0]
	quote := parts[1]

	quoteQty := price.Mul(baseQty)
	// This ensure the data consistency
	baseAsset := acc.GetAsset(base)
	quoteAsset := acc.GetAsset(quote)

	if side {
		if baseAsset.Free.LessThan(baseQty) {
			return ErrInsufficientBalance
		}
		baseAsset.Free = baseAsset.Free.Sub(baseQty)
		quoteAsset.Free = quoteAsset.Free.Add(quoteQty)
	} else {
		if quoteAsset.Free.LessThan(quoteQty) {
			return ErrInsufficientBalance
		}
		baseAsset.Free = baseAsset.Free.Add(baseQty)
		quoteAsset.Free = quoteAsset.Free.Sub(quoteQty)
	}
	return nil
}
