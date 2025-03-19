package solvexity

import (
	"fmt"

	"github.com/BullionBear/sequex/internal/payload"
	"github.com/BullionBear/sequex/internal/strategy"
)

var _ strategy.Strategy = (*Solvexity)(nil)

type Solvexity struct {
}

func NewSolvexity() *Solvexity {
	return &Solvexity{}
}

func (s *Solvexity) OnKLineUpdate(meta payload.KLineUpdate) {
	fmt.Printf("solvexity strategy received kline update %+v\n", meta)
}
