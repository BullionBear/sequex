package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BullionBear/sequex/internal/config"
	_ "github.com/BullionBear/sequex/internal/nodeimpl/init" // Import to register all nodes
	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/node"
)

func main() {
	// Initialize structured logger
	logger := log.New(
		log.WithLevel(log.LevelInfo),
		log.WithEncoder(log.NewTextEncoder()),
		log.WithTimeRotation("./logs", "node.log", 24*time.Hour, 7),
	)

	// Parse command line arguments
	var configFile string
	flag.StringVar(&configFile, "c", "config/rng.yml", "Configuration file path")
	flag.Parse()

	logger.Info("Starting NATS Microservices",
		log.String("config_file", configFile),
		log.String("component", "node_deployer"),
	)

	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		logger.Fatal("Failed to load configuration",
			log.String("config_file", configFile),
			log.Error(err),
		)
	}

	// Create a single NATS connection for the entire process
	nc, err := config.CreateNATSConnection(cfg.NATS.URL)
	if err != nil {
		logger.Fatal("Failed to connect to NATS",
			log.String("nats_url", cfg.NATS.URL),
			log.Error(err),
		)
	}
	defer nc.Close()

	logger.Info("Successfully connected to NATS",
		log.String("nats_url", cfg.NATS.URL),
	)

	// Create deployer
	d := node.NewDeployer()

	// Create and register nodes based on configuration
	for _, nodeConfig := range cfg.Deployer.Nodes {
		nodeName := nodeConfig["name"].(string)
		logger.Info("Creating node",
			log.String("node_name", nodeName),
			log.String("component", "node_deployer"),
		)

		node, err := config.CreateNode(nodeConfig, nc)
		if err != nil {
			logger.Fatal("Failed to create node",
				log.String("node_name", nodeName),
				log.Error(err),
			)
		}

		if err := d.RegisterNode(node); err != nil {
			logger.Fatal("Failed to register node",
				log.String("node_name", nodeName),
				log.Error(err),
			)
		}

		logger.Info("Successfully registered node",
			log.String("node_name", nodeName),
		)
	}

	// Start all nodes
	for _, nodeConfig := range cfg.Deployer.Nodes {
		nodeName := nodeConfig["name"].(string)
		if err := d.Start(nodeName); err != nil {
			logger.Fatal("Failed to start node",
				log.String("node_name", nodeName),
				log.Error(err),
			)
		}
		logger.Info("Successfully started node",
			log.String("node_name", nodeName),
		)
	}

	logger.Info("All nodes deployed successfully",
		log.Int("node_count", len(cfg.Deployer.Nodes)),
	)

	// Wait for shutdown signal
	waitForShutdown(logger)

	// Stop all nodes
	for _, nodeConfig := range cfg.Deployer.Nodes {
		nodeName := nodeConfig["name"].(string)
		if err := d.Stop(nodeName); err != nil {
			logger.Error("Error stopping node",
				log.String("node_name", nodeName),
				log.Error(err),
			)
		}
	}

	logger.Info("All nodes stopped")
}

func waitForShutdown(logger log.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Waiting for shutdown signal...")
	<-sigChan
	logger.Info("Shutdown signal received")
}
