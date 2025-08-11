package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BullionBear/sequex/internal/config"
	_ "github.com/BullionBear/sequex/internal/nodeimpl/init" // Import to register all nodes
	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/node"
)

func main() {
	// Parse command line arguments
	var configFile string
	flag.StringVar(&configFile, "c", "config/rng.yml", "Configuration file path")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		// Use fmt for error before logger is initialized
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize global logger from config
	if err := config.InitializeLogger(cfg.Logger); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	config.Info("Starting NATS Microservices",
		log.String("config_file", configFile),
		log.String("component", "node_deployer"),
	)

	// Create a single NATS connection for the entire process
	nc, err := config.CreateNATSConnection(cfg.NATS.URL)
	if err != nil {
		config.Fatal("Failed to connect to NATS",
			log.String("nats_url", cfg.NATS.URL),
			log.Error(err),
		)
	}
	defer nc.Close()

	config.Info("Successfully connected to NATS",
		log.String("nats_url", cfg.NATS.URL),
	)

	// Create deployer
	d := node.NewDeployer()

	// Create and register nodes based on configuration
	for _, nodeConfig := range cfg.Deployer.Nodes {
		nodeName := nodeConfig["name"].(string)
		config.Info("Creating node",
			log.String("node_name", nodeName),
			log.String("component", "node_deployer"),
		)

		node, err := config.CreateNode(nodeConfig, nc)
		if err != nil {
			config.Fatal("Failed to create node",
				log.String("node_name", nodeName),
				log.Error(err),
			)
		}

		if err := d.RegisterNode(node); err != nil {
			config.Fatal("Failed to register node",
				log.String("node_name", nodeName),
				log.Error(err),
			)
		}

		config.Info("Successfully registered node",
			log.String("node_name", nodeName),
		)
	}

	// Start all nodes
	for _, nodeConfig := range cfg.Deployer.Nodes {
		nodeName := nodeConfig["name"].(string)
		if err := d.Start(nodeName); err != nil {
			config.Fatal("Failed to start node",
				log.String("node_name", nodeName),
				log.Error(err),
			)
		}
		config.Info("Successfully started node",
			log.String("node_name", nodeName),
		)
	}

	config.Info("All nodes deployed successfully",
		log.Int("node_count", len(cfg.Deployer.Nodes)),
	)

	// Wait for shutdown signal
	waitForShutdown()

	// Stop all nodes
	for _, nodeConfig := range cfg.Deployer.Nodes {
		nodeName := nodeConfig["name"].(string)
		if err := d.Stop(nodeName); err != nil {
			config.Error("Error stopping node",
				log.String("node_name", nodeName),
				log.Error(err),
			)
		}
	}

	config.Info("All nodes stopped")
}

func waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	config.Info("Waiting for shutdown signal...")
	<-sigChan
	config.Info("Shutdown signal received")
}
