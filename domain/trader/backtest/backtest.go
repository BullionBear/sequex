package trader

import (
	"github.com/BullionBear/crypto-trade/domain/pgdb"
	"github.com/BullionBear/crypto-trade/domain/wallet"
	"github.com/shopspring/decimal"
)

type Backtest struct {
	db     *pgdb.PgDatabase
	wallet *wallet.Wallet
}

func NewBacktest(db *pgdb.PgDatabase, wallet *wallet.Wallet) *Backtest {
	return &Backtest{
		db:     db,
		wallet: wallet,
	}

}

func (b *Backtest) CreateMarketOrder(symbol string, side bool, quoteQty decimal.Decimal, openTime int64) error {
	kline, err := b.db.GetKline(symbol, openTime)
	if err != nil {
		return err
	}
	return b.wallet.Trade(symbol, side, kline.Close, quoteQty)
}
