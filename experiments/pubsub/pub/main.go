package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		log.Println("Received shutdown signal, stopping publisher...")
		cancel()
	}()

	// Start publishing messages
	publishMessages(ctx, nc)
}

func publishMessages(ctx context.Context, nc *nats.Conn) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	messageCount := 0

	for {
		select {
		case <-ctx.Done():
			log.Println("Publisher stopped")
			return
		case <-ticker.C:
			messageCount++
			message := fmt.Sprintf("Trade message #%d at %s", messageCount, time.Now().Format(time.RFC3339))

			// Publish message
			err := nc.Publish(subject, []byte(message))
			if err != nil {
				log.Printf("Failed to publish message: %v", err)
				continue
			}

			// Flush to ensure message is sent
			err = nc.Flush()
			if err != nil {
				log.Printf("Failed to flush: %v", err)
				continue
			}

			log.Printf("Published: %s", message)
		}
	}
}
