package message

import "encoding/json"

type Message struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Source    string      `json:"source"`
	Target    string      `json:"target"`
	CreatedAt int64       `json:"created_at"`
	Data      interface{} `json:"data"`
	Metadata  interface{} `json:"metadata"`
}

// MarshalJSON serializes the Message struct into JSON.
func (m *Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// UnmarshalJSON deserializes JSON data into a Message struct.
func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	return json.Unmarshal(data, aux)
}
