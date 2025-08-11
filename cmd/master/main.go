package main

import (
	"github.com/BullionBear/sequex/pkg/log"
)

func main() {
	// Initialize structured logger
	logger := log.New(
		log.WithLevel(log.LevelInfo),
		log.WithEncoder(log.NewTextEncoder()),
	)

	logger.Info("Hello, World!",
		log.String("component", "master"),
		log.String("version", "1.0.0"),
	)
}
