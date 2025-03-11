package inprocq

import (
	"testing"
	"time"

	"strconv"

	"github.com/BullionBear/sequex/pkg/message"
	"github.com/stretchr/testify/assert"
)

// TestPublishSubscribe checks if messages are correctly published and received.
func TestPublishSubscribe(t *testing.T) {
	queue := New(10)

	// Subscribe to a topic
	sub, err := queue.Subscribe("tech")
	assert.NoError(t, err)

	// Publish messages
	msg1 := message.Message{ID: "1", Data: "Go 1.20 released", CreatedAt: time.Now().Unix()}
	msg2 := message.Message{ID: "2", Data: "New AI model", CreatedAt: time.Now().Unix()}
	err = queue.Publish("tech", msg1)
	assert.NoError(t, err)
	err = queue.Publish("tech", msg2)
	assert.NoError(t, err)

	// Receive messages
	received1 := <-sub
	received2 := <-sub

	assert.Equal(t, msg1, received1)
	assert.Equal(t, msg2, received2)
}

func TestUnsubscribe(t *testing.T) {
	queue := New(10)

	// Subscribe to a topic
	sub, err := queue.Subscribe("sports")
	assert.NoError(t, err)

	// Unsubscribe from the topic
	err = queue.Unsubscribe("sports", sub)
	assert.NoError(t, err)

	// Try publishing after unsubscribe
	msg := message.Message{ID: "3", Data: "Football World Cup", CreatedAt: time.Now().Unix()}
	err = queue.Publish("sports", msg)
	assert.NoError(t, err)

	// Ensure no messages are received and the channel is closed
	received := false
	timeout := time.After(100 * time.Millisecond)
	for {
		select {
		case _, ok := <-sub:
			if !ok {
				// Channel is closed, pass the test
				return
			} else {
				received = true
			}
		case <-timeout:
			if received {
				t.Error("Received message after unsubscribe")
			}
			// Test passes if no message is received and timeout occurs
			return
		}
	}
}

// TestChannelFull checks behavior when subscriber channels are full.
func TestChannelFull(t *testing.T) {
	queue := New(10)

	// Subscribe with a small buffer
	_, err := queue.Subscribe("news") // Remove unused `sub` variable
	assert.NoError(t, err)

	// Fill the channel buffer
	for i := 0; i < 10; i++ {
		err := queue.Publish("news", message.Message{ID: strconv.Itoa(i), Data: "News", CreatedAt: time.Now().Unix()})
		assert.NoError(t, err)
	}

	// Try to publish one more message
	err = queue.Publish("news", message.Message{ID: "11", Data: "Overflow", CreatedAt: time.Now().Unix()})
	assert.Error(t, err, "Expected error when channel is full")
}

// TestConcurrentAccess checks thread safety of concurrent publish and subscribe.
func TestConcurrentAccess(t *testing.T) {
	queue := New(10)
	sub, err := queue.Subscribe("concurrent")
	assert.NoError(t, err)

	var receivedMessages []message.Message
	done := make(chan struct{})

	// Subscriber goroutine
	go func() {
		for i := 0; i < 100; i++ {
			msg := <-sub
			receivedMessages = append(receivedMessages, msg)
		}
		close(done)
	}()

	// Publisher goroutines
	for i := 0; i < 100; i++ {
		go queue.Publish("concurrent", message.Message{ID: strconv.Itoa(i), Data: "Message", CreatedAt: time.Now().Unix()})
	}

	// Wait for subscriber to finish
	<-done

	// Check all messages were received
	assert.Equal(t, 100, len(receivedMessages))
}
