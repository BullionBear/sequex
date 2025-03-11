package sequex

import (
	"github.com/BullionBear/sequex/internal/metadata"
	"github.com/BullionBear/sequex/internal/strategy"
)

var _ strategy.Strategy = (*Sequex)(nil)

type Sequex struct {
}

func NewSequex() *Sequex {
	return &Sequex{}
}

func (s *Sequex) OnKLineUpdate(metadata.KLineUpdate) {
}
