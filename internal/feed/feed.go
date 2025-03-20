package feed

import "github.com/BullionBear/sequex/internal/payload"

type Feed interface {
	SubscribeKlineUpdate(symbol string, handler func(*payload.KLineUpdate)) (unsubscribe func(), err error)
}
