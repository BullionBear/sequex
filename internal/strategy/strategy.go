package strategy

import "github.com/BullionBear/sequex/internal/metadata"

type Strategy interface {
	OnKLineUpdate(metadata.KLineUpdate)
}
