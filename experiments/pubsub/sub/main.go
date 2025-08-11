package main

import (
	"context"
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
		log.WithTimeRotation("./logs", "pubsub_subscriber.log", 24*time.Hour, 7),
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
		log.String("component", "pubsub_subscriber"),
	)

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("Received shutdown signal, stopping subscriber")
		cancel()
	}()

	// Subscribe to messages
	subscription, err := nc.Subscribe(subject, func(msg *nats.Msg) {
		logger.Info("Received message",
			log.String("subject", msg.Subject),
			log.String("message", string(msg.Data)),
			log.String("reply", msg.Reply),
		)

		// Acknowledge the message
		msg.Ack()
	})
	if err != nil {
		logger.Fatal("Failed to subscribe",
			log.String("subject", subject),
			log.Error(err),
		)
	}
	defer subscription.Unsubscribe()

	logger.Info("Subscribed to subject",
		log.String("subject", subject),
	)

	// Keep the subscriber running
	<-ctx.Done()
	logger.Info("Subscriber stopped")
}
