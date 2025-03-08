package inprocq

import (
	"errors"
	"sync"

	"github.com/BullionBear/sequex/pkg/mq"
)

// Ensure InprocQueue implements the mq.MessageQueue interface.
var _ mq.MessageQueue = (*InprocQueue)(nil)

// InprocQueue is an in-process message queue with topic support.
type InprocQueue struct {
	topics sync.Map // Thread-safe map to store topic channels.
	size   uint     // Size determines if channels are buffered or unbuffered.
}

// New creates a new InprocQueue.
// If size > 0, creates buffered channels of given size.
// If size == 0, creates unbuffered channels.
func New(size uint) *InprocQueue {
	return &InprocQueue{size: size}
}

// Publish sends a message to all subscribers of a topic.
func (q *InprocQueue) Publish(topic string, msg mq.Message) error {
	if subs, ok := q.topics.Load(topic); ok {
		for _, ch := range subs.([]chan mq.Message) {
			select {
			case ch <- msg:
			default:
				// Drop message if subscriber's channel is full
				return errors.New("subscriber channel is full, message dropped")
			}
		}
	}
	return nil
}

// Subscribe allows a user to subscribe to a topic.
func (q *InprocQueue) Subscribe(topic string) (<-chan mq.Message, error) {
	var ch chan mq.Message

	// Create either buffered or unbuffered channel based on size
	if q.size > 0 {
		ch = make(chan mq.Message, q.size) // Buffered channel
	} else {
		ch = make(chan mq.Message) // Unbuffered channel
	}

	if subs, ok := q.topics.Load(topic); ok {
		q.topics.Store(topic, append(subs.([]chan mq.Message), ch))
	} else {
		q.topics.Store(topic, []chan mq.Message{ch})
	}

	return ch, nil
}

// Unsubscribe removes a subscriber's channel from a topic.
func (q *InprocQueue) Unsubscribe(topic string, ch <-chan mq.Message) error {
	if subs, ok := q.topics.Load(topic); ok {
		channels := subs.([]chan mq.Message)
		for i, sub := range channels {
			if sub == ch {
				// Remove the channel from the slice
				newChannels := append(channels[:i], channels[i+1:]...)
				if len(newChannels) == 0 {
					q.topics.Delete(topic)
				} else {
					q.topics.Store(topic, newChannels)
				}
				close(sub)
				break
			}
		}
	}
	return nil
}
