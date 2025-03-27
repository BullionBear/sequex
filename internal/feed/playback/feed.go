package playback

import (
	"github.com/BullionBear/sequex/internal/feed"
	"github.com/BullionBear/sequex/internal/payload"
)

var _ feed.Feed = (*PlaybackFeed)(nil)

type PlaybackFeed struct {
}

func NewPlaybackFeed() *PlaybackFeed {
	return &PlaybackFeed{}
}

func (p *PlaybackFeed) SubscribeKlineUpdate(symbol string, handler func(*payload.KLineUpdate)) (unsubscribe func(), err error) {
	// Implement me
	return nil, nil
}

func (p *PlaybackFeed) Next(symbol string) {
	// Implement me
}
