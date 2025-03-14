package mq

import (
	"time"

	"github.com/BullionBear/sequex/pkg/message"
)

// MessageQueue defines the interface for a message queue system.
type MessageQueue interface {
	// Send a message to the queue.
	Send(msg *message.Message) error
	// Receive a message from the queue.
	Recv() (*message.Message, error)
	// Receive a message with a timeout.
	RecvTimeout(timeout time.Duration) (*message.Message, error)
	// Close the queue.
	Size() uint64
	// Clear the queue.
	Clear() error
	// is the queue empty.
	IsEmpty() bool
}
