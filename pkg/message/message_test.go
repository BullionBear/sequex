package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageMarshalUnmarshal(t *testing.T) {
	originalMessage := Message{
		ID:        "1234",
		Type:      "test_type",
		CreatedAt: 1710456789,
		Payload: map[string]interface{}{
			"key": "value",
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(originalMessage)
	assert.NoError(t, err, "MarshalJSON should not return an error")
	assert.NotEmpty(t, jsonData, "Marshalled JSON should not be empty")

	// Unmarshal back to struct
	var unmarshalledMessage Message
	err = json.Unmarshal(jsonData, &unmarshalledMessage)
	assert.NoError(t, err, "UnmarshalJSON should not return an error")

	// Verify fields match
	assert.Equal(t, originalMessage.ID, unmarshalledMessage.ID, "ID should match")
	assert.Equal(t, originalMessage.Type, unmarshalledMessage.Type, "Type should match")
	assert.Equal(t, originalMessage.CreatedAt, unmarshalledMessage.CreatedAt, "CreatedAt should match")
	assert.Equal(t, originalMessage.Payload, unmarshalledMessage.Payload, "Data should match")
}
