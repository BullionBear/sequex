package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

const (
	subject = "sequex.trades"
	url     = nats.DefaultURL
)

func main() {
	// Connect to NATS
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	log.Printf("Connected to NATS at %s", url)

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping subscriber...")
		cancel()
	}()

	// Subscribe to messages
	subscription, err := nc.Subscribe(subject, func(msg *nats.Msg) {
		log.Printf("Received message: %s", string(msg.Data))

		// Acknowledge the message
		msg.Ack()
	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	defer subscription.Unsubscribe()

	log.Printf("Subscribed to %s", subject)

	// Keep the subscriber running
	<-ctx.Done()
	log.Println("Subscriber stopped")
}
