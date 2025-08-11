package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BullionBear/sequex/pkg/log"
	"github.com/nats-io/nats.go"
)

const (
	subject = "sequex.trades"
	url     = nats.DefaultURL
)

func main() {
	// Initialize structured logger
	logger := log.New(
		log.WithLevel(log.LevelInfo),
		log.WithEncoder(log.NewTextEncoder()),
		log.WithTimeRotation("./logs", "pubsub_publisher.log", 24*time.Hour, 7),
	)

	// Connect to NATS
	nc, err := nats.Connect(url)
	if err != nil {
		logger.Fatal("Failed to connect to NATS",
			log.String("nats_url", url),
			log.Error(err),
		)
	}
	defer nc.Close()

	logger.Info("Connected to NATS",
		log.String("nats_url", url),
		log.String("component", "pubsub_publisher"),
	)

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Received shutdown signal, stopping publisher")
		cancel()
	}()

	// Start publishing messages
	publishMessages(ctx, nc, logger)
}

func publishMessages(ctx context.Context, nc *nats.Conn, logger log.Logger) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	messageCount := 0
	publisherLogger := logger.With(log.String("component", "message_publisher"))

	for {
		select {
		case <-ctx.Done():
			publisherLogger.Info("Publisher stopped",
				log.Int("total_messages", messageCount),
			)
			return
		case <-ticker.C:
			messageCount++
			message := fmt.Sprintf("Trade message #%d at %s", messageCount, time.Now().Format(time.RFC3339))

			// Publish message
			err := nc.Publish(subject, []byte(message))
			if err != nil {
				publisherLogger.Error("Failed to publish message",
					log.Int("message_number", messageCount),
					log.String("subject", subject),
					log.String("message", message),
					log.Error(err),
				)
				continue
			}

			// Flush to ensure message is sent
			err = nc.Flush()
			if err != nil {
				publisherLogger.Error("Failed to flush connection",
					log.Int("message_number", messageCount),
					log.Error(err),
				)
				continue
			}

			publisherLogger.Info("Message published successfully",
				log.Int("message_number", messageCount),
				log.String("subject", subject),
				log.String("message", message),
			)
		}
	}
}
