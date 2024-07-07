package backtest

/*
	Backtest is a struct for simulating exchange behavior for backtesting.
	The basic components include:
	- Data: Connecting crypto-feed gRPC
	- Account: Managing the balance of assets
	- Order: Managing the order status (ignore)
	- Position: Managing the position status
*/

import (
	"github.com/BullionBear/crypto-trade/domain/feedclient"

	"github.com/BullionBear/crypto-trade/domain/models"
	"github.com/shopspring/decimal"
)

type Backtest struct {
	account      *Account
	feed         *feedclient.FeedClient
	position     *Position
	currentKline models.Kline
}

func NewBacktest(acc *Account, feed *feedclient.FeedClient) *Backtest {
	position := NewPosition()
	return &Backtest{
		account:      acc,
		feed:         feed,
		position:     position,
		currentKline: models.Kline{},
	}
}

func (b *Backtest) OpenPosition(symbol string, side bool, quoteQty decimal.Decimal, openTime int64) (int, error) {
	kline, err := b.db.QueryKline(openTime)
	if err != nil {
		return err
	}
	price := decimal.NewFromFloat(kline.Close)
	baseQty := quoteQty.Div(price)
	return b.wallet.Trade(symbol, side, price, baseQty)
}
