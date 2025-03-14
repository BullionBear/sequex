package inprocq

import (
	"sync"
	"testing"
	"time"

	"github.com/BullionBear/sequex/pkg/message"
	"github.com/stretchr/testify/assert"
)

func TestSendAndRecv(t *testing.T) {
	q := NewInprocQueue()
	msg := &message.Message{ID: "1"}
	err := q.Send(msg)
	assert.NoError(t, err)

	received, err := q.Recv()
	assert.NoError(t, err)
	assert.Equal(t, msg, received)
}

func TestRecvBlocks(t *testing.T) {
	q := NewInprocQueue()

	go func() {
		time.Sleep(100 * time.Millisecond)
		q.Send(&message.Message{ID: "1"})
	}()

	start := time.Now()
	msg, err := q.Recv()
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, "1", msg.ID)
	assert.True(t, elapsed >= 100*time.Millisecond, "Recv returned too quickly")
}

func TestRecvTimeout(t *testing.T) {
	q := NewInprocQueue()

	// Test timeout
	start := time.Now()
	_, err := q.RecvTimeout(100 * time.Millisecond)
	elapsed := time.Since(start)
	assert.Error(t, err)
	assert.True(t, elapsed >= 100*time.Millisecond, "Timeout too short")

	// Test message received before timeout
	go func() {
		time.Sleep(50 * time.Millisecond)
		q.Send(&message.Message{ID: "2"})
	}()

	msg, err := q.RecvTimeout(100 * time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, "2", msg.ID)
}

func TestConcurrentSendAndRecv(t *testing.T) {
	q := NewInprocQueue()
	const numMessages = 100
	var wg sync.WaitGroup

	for i := 0; i < numMessages; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			q.Send(&message.Message{ID: string(rune(id))})
		}(i)
	}

	received := make(chan *message.Message, numMessages)
	for i := 0; i < numMessages; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			msg, err := q.Recv()
			assert.NoError(t, err)
			received <- msg
		}()
	}

	wg.Wait()
	close(received)

	count := 0
	for range received {
		count++
	}
	assert.Equal(t, numMessages, count)
}

func TestSize(t *testing.T) {
	q := NewInprocQueue()
	assert.Equal(t, uint64(0), q.Size())

	q.Send(&message.Message{})
	assert.Equal(t, uint64(1), q.Size())

	q.Recv()
	assert.Equal(t, uint64(0), q.Size())
}

func TestClear(t *testing.T) {
	q := NewInprocQueue()
	q.Send(&message.Message{})
	q.Send(&message.Message{})

	assert.Equal(t, uint64(2), q.Size())
	assert.NoError(t, q.Clear())
	assert.Equal(t, uint64(0), q.Size())
	assert.True(t, q.IsEmpty())
}

func TestIsEmpty(t *testing.T) {
	q := NewInprocQueue()
	assert.True(t, q.IsEmpty())

	q.Send(&message.Message{})
	assert.False(t, q.IsEmpty())
}
