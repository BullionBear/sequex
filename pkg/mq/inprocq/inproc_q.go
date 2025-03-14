package inprocq

import (
	"errors"
	"sync"
	"time"

	"github.com/BullionBear/sequex/pkg/message"
	"github.com/BullionBear/sequex/pkg/mq"
)

// InprocQueue is an in-memory implementation of the MessageQueue interface.
type InprocQueue struct {
	mu    sync.Mutex
	cond  *sync.Cond
	queue []*message.Message
}

// NewInprocQueue creates a new instance of InprocQueue.
func NewInprocQueue() *InprocQueue {
	q := &InprocQueue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Send adds a message to the queue.
func (q *InprocQueue) Send(msg *message.Message) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.queue = append(q.queue, msg)
	q.cond.Signal()
	return nil
}

// Recv blocks until a message is received from the queue.
func (q *InprocQueue) Recv() (*message.Message, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.queue) == 0 {
		q.cond.Wait()
	}

	msg := q.queue[0]
	q.queue = q.queue[1:]
	return msg, nil
}

// RecvTimeout waits for a message with a specified timeout, returning an error if the timeout is exceeded.
func (q *InprocQueue) RecvTimeout(timeout time.Duration) (*message.Message, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	deadline := time.Now().Add(timeout)
	for len(q.queue) == 0 {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return nil, errors.New("timeout")
		}

		// Create a timer to trigger broadcast on timeout
		timer := time.AfterFunc(remaining, q.cond.Broadcast)
		defer timer.Stop()

		q.cond.Wait()
		if time.Now().After(deadline) {
			return nil, errors.New("timeout")
		}
	}

	msg := q.queue[0]
	q.queue = q.queue[1:]
	return msg, nil
}

// Size returns the number of messages in the queue.
func (q *InprocQueue) Size() uint64 {
	q.mu.Lock()
	defer q.mu.Unlock()
	return uint64(len(q.queue))
}

// Clear removes all messages from the queue.
func (q *InprocQueue) Clear() error {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = nil
	return nil
}

// IsEmpty checks if the queue is empty.
func (q *InprocQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue) == 0
}

var _ mq.MessageQueue = (*InprocQueue)(nil)
