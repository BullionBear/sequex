package wallet

import (
	"strings"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	balance map[string]decimal.Decimal
}

func NewWallet() *Wallet {
	return &Wallet{
		balance: make(map[string]decimal.Decimal),
	}
}

func (w *Wallet) SetBalance(token string, balance decimal.Decimal) {
	w.balance[token] = balance
}

func (w *Wallet) GetBalance(symbol string) decimal.Decimal {
	if b, ok := w.balance[symbol]; ok {
		return b
	}
	return decimal.Zero
}

func (w *Wallet) Trade(symbol string, side bool, price, baseQty decimal.Decimal) error {
	parts := strings.Split(symbol, "-")
	base := parts[0]
	quote := parts[1]

	quoteQty := price.Mul(baseQty)
	// This ensure the data consistency
	baseBalance := w.GetBalance(base)
	quoteBalance := w.GetBalance(quote)
	if side {
		w.balance[base] = baseBalance.Sub(baseQty)
		w.balance[quote] = quoteBalance.Add(quoteQty)
	} else {
		w.balance[base] = baseBalance.Add(baseQty)
		w.balance[quote] = quoteBalance.Sub(quoteQty)
	}
	return nil
}
