package orderbook

import (
	"errors"
	"math"

	"github.com/shopspring/decimal"
)

const MaxPriceLevels = 5000

type PriceLevel struct {
	Price decimal.Decimal
	Size  decimal.Decimal
}

func NewPriceLevel(price, size decimal.Decimal) PriceLevel {
	return PriceLevel{
		Price: price,
		Size:  size,
	}
}

func (pl *PriceLevel) Empty() {
	pl.Price = decimal.Zero
	pl.Size = decimal.Zero
}

func (pl *PriceLevel) Set(price, size decimal.Decimal) {
	pl.Price = price
	pl.Size = size
}

type AskBookArray struct {
	PriceLevels  [MaxPriceLevels]PriceLevel // Static array with a fixed size of 100
	BestIndex    int
	PriceDecimal decimal.Decimal
}

func NewAskBookArray(priceDecimal int) *AskBookArray {
	return &AskBookArray{
		PriceLevels:  [MaxPriceLevels]PriceLevel{},
		BestIndex:    math.MaxInt,
		PriceDecimal: decimal.NewFromInt(1).Div(decimal.NewFromInt(int64(math.Pow10(priceDecimal)))),
	}
}

func (oa *AskBookArray) GetBestLayer() (PriceLevel, error) {
	if oa.BestIndex >= 0 && oa.BestIndex < MaxPriceLevels {
		return oa.PriceLevels[oa.BestIndex], nil
	}
	return PriceLevel{}, errors.New("best price not available")
}

func (oa *AskBookArray) GetBook(depth int) []PriceLevel {
	if depth <= 0 || depth > MaxPriceLevels {
		return nil
	}

	book := make([]PriceLevel, 0, depth)
	j := 0
	for i := 0; i < MaxPriceLevels; i++ {
		if !oa.PriceLevels[i].Size.IsZero() {
			book[j] = oa.PriceLevels[(oa.BestIndex+i)%MaxPriceLevels]
			j++
		}
		if j == depth {
			break
		}
	}
	return book
}

func (oa *AskBookArray) UpdateDiff(levels []PriceLevel) {
	for _, level := range levels {
		index := int(level.Price.Div(oa.PriceDecimal).IntPart())
		oa.PriceLevels[index%MaxPriceLevels].Set(level.Price, level.Size)
		oa.BestIndex = min(oa.BestIndex, index)
		if level.Size.IsZero() {
			oa.PriceLevels[index%MaxPriceLevels].Empty()
		} else {
			oa.PriceLevels[index%MaxPriceLevels].Set(level.Price, level.Size)
		}
	}
}

func (oa *AskBookArray) UpdateAll(levels []PriceLevel) {
	for i := 0; i < MaxPriceLevels; i++ {
		oa.PriceLevels[i].Empty()
	}
	oa.BestIndex = math.MaxInt
	oa.UpdateDiff(levels)
}

type BidBookArray struct {
	PriceLevels  [MaxPriceLevels]PriceLevel // Static array with a fixed size of 100
	BestIndex    int
	PriceDecimal decimal.Decimal
}

func NewBidBookArray(priceDecimal int) *BidBookArray {
	return &BidBookArray{
		PriceLevels:  [MaxPriceLevels]PriceLevel{},
		BestIndex:    math.MinInt,
		PriceDecimal: decimal.NewFromInt(1).Div(decimal.NewFromInt(int64(math.Pow10(priceDecimal)))),
	}
}

func (ob *BidBookArray) GetBestLayer() (PriceLevel, error) {
	if ob.BestIndex >= 0 && ob.BestIndex < MaxPriceLevels {
		return ob.PriceLevels[ob.BestIndex], nil
	}
	return PriceLevel{}, errors.New("best price not available")
}

func (ob *BidBookArray) GetBook(depth int) []PriceLevel {
	if depth <= 0 || depth > MaxPriceLevels {
		return nil
	}

	book := make([]PriceLevel, 0, depth)
	j := 0
	for i := 0; i < MaxPriceLevels; i++ {
		if !ob.PriceLevels[i].Size.IsZero() {
			book[j] = ob.PriceLevels[(ob.BestIndex-i)%MaxPriceLevels] // Note: Adjusted for bid book
			j++
		}
		if j == depth {
			break
		}
	}
	return book
}

func (ob *BidBookArray) UpdateDiff(levels []PriceLevel) {
	for _, level := range levels {
		index := int(level.Price.Div(ob.PriceDecimal).IntPart())
		ob.PriceLevels[index%MaxPriceLevels].Set(level.Price, level.Size)
		ob.BestIndex = max(ob.BestIndex, index)
		if level.Size.IsZero() {
			ob.PriceLevels[index%MaxPriceLevels].Empty()
		} else {
			ob.PriceLevels[index%MaxPriceLevels].Set(level.Price, level.Size)
		}
	}
}

func (ob *BidBookArray) UpdateAll(levels []PriceLevel) {
	for i := 0; i < MaxPriceLevels; i++ {
		ob.PriceLevels[i].Empty()
	}
	ob.BestIndex = math.MinInt
	ob.UpdateDiff(levels)
}

type OrderBook struct {
	Asks         AskBookArray
	Bids         BidBookArray
	timestamp    int64
	lastUpdateID int64
}

func NewOrderBook(priceDecimal int) *OrderBook {
	return &OrderBook{
		Asks:         *NewAskBookArray(priceDecimal),
		Bids:         *NewBidBookArray(priceDecimal),
		timestamp:    0,
		lastUpdateID: 0,
	}
}
