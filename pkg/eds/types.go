package eds

type EventType string

const (
	KLineEvent EventType = "KLineEvent"
	TradeEvent EventType = "TradeEvent"
)
