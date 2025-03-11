package message

type Message struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	CreatedAt int64       `json:"created_at"`
	Data      interface{} `json:"data"`
	Metadata  interface{} `json:"metadata"`
}
