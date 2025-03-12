package metadata

type KLineUpdate struct {
	Symbol    string `json:"symbol"`
	Interval  string `json:"interval"`
	Timestamp int64  `json:"timestamp"`
}
