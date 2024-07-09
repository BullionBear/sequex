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
	"strings"

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

func (b *Backtest) OpenPosition(symbol string, quoteQty decimal.Decimal, leverage int) (int, error) {
	margin := quoteQty.Div(decimal.NewFromInt(int64(leverage)))
	parts := strings.Split(symbol, "-")
	_, quote := parts[0], parts[1]
	if err := b.account.Lock(quote, margin); err != nil {
		return 0, err
	}
	price := decimal.NewFromFloat(b.currentKline.Close)
	openTime := b.currentKline.OpenTime
	baseQty := quoteQty.Div(price)
	orderId := b.position.OpenPosition(symbol, openTime, baseQty, price, leverage)
	return orderId, nil
}

func (b *Backtest) ClosePosition(orderId int) error {
	perp, err := b.position.ClosePosition(orderId)
	if err != nil {
		return err
	}
	parts := strings.Split(perp.Symbol, "-")
	_, quote := parts[0], parts[1]
	openPrice := perp.OpenPrice
	closePrice := decimal.NewFromFloat(b.currentKline.Close)
	pnl := perp.BaseQty.Mul(closePrice.Sub(openPrice))
	b.account.Unlock(quote, perp.Margin)
	asset := b.account.GetAsset(quote)
	asset.Free = asset.Free.Add(pnl)
	return nil
}
