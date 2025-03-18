package strategy

import "github.com/BullionBear/sequex/internal/payload"

type Strategy interface {
	OnKLineUpdate(payload.KLineUpdate)
}
