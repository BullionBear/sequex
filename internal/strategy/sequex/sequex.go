package sequex

import (
	"fmt"

	"github.com/BullionBear/sequex/internal/metadata"
	"github.com/BullionBear/sequex/internal/strategy"
)

var _ strategy.Strategy = (*Sequex)(nil)

type Sequex struct {
}

func NewSequex() *Sequex {
	return &Sequex{}
}

func (s *Sequex) OnKLineUpdate(meta metadata.KLineUpdate) {
	fmt.Printf("Sequex strategy received kline update %+v\n", meta)
}
