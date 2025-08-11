package main

import (
	"github.com/BullionBear/sequex/internal/config"
	"github.com/BullionBear/sequex/pkg/log"
)

func main() {
	// Initialize with default config
	defaultConfig := config.LoggerConfig{
		Format: "text",
		Level:  "info",
		Path:   "master.log",
	}

	if err := config.InitializeLogger(defaultConfig); err != nil {
		panic(err)
	}

	config.Info("Hello, World!",
		log.String("component", "master"),
		log.String("version", "1.0.0"),
	)
}
