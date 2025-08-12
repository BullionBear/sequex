package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/BullionBear/sequex/env"
	"github.com/BullionBear/sequex/internal/config"
	_ "github.com/BullionBear/sequex/internal/nodeimpl/init" // Import to register all nodes
	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/node"
	"github.com/BullionBear/sequex/pkg/shutdown"
)

var logger log.Logger

func main() {
	// Parse command line arguments
	var configFile string
	flag.StringVar(&configFile, "c", "config/node/example.yml", "Configuration file path")
	flag.Parse()

	fmt.Println("Starting services with BuildTime:", env.BuildTime)
	fmt.Println("Starting services with Version:", env.Version)
	fmt.Println("Starting services with CommitHash:", env.CommitHash)

	// Load configuration
	cfg, err := config.LoadConfig[config.Config](configFile)
	if err != nil {
		// Use fmt for error before logger is initialized
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize global logger from config
	logger, err = config.CreateLogger(cfg.Logger)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logger.Info("Starting services with",
		log.String("build_time", env.BuildTime),
		log.String("version", env.Version),
		log.String("commit_hash", env.CommitHash),
	)

	logger.Info("Starting NATS Microservices",
		log.String("config_file", configFile),
		log.String("component", "node_deployer"),
	)

	// Create shutdown
	shutdown := shutdown.NewShutdown(logger)

	// Create a single NATS connection for the entire process
	nc, err := config.CreateNATSConnection(cfg.NATS.URL)
	if err != nil {
		logger.Fatalf("Failed to connect to NATS %s: %v", cfg.NATS.URL, err)
	}
	defer nc.Close()

	logger.Infof("Successfully connected to NATS %s", cfg.NATS.URL)

	// Create deployer
	d := node.NewDeployer(&logger)

	// Create and register nodes based on configuration
	for _, nodeConfig := range cfg.Deployer.Nodes {
		nodeName := nodeConfig["name"].(string)
		logger.Infof("Creating node %s", nodeName)

		node, err := config.CreateNode(nodeConfig, nc, &logger)
		if err != nil {
			logger.Fatalf("Failed to create node %s: %v", nodeName, err)
		}

		if err := d.RegisterNode(node); err != nil {
			logger.Fatalf("Failed to register node %s: %v", nodeName, err)
		}

		logger.Infof("Successfully registered node %s", nodeName)
	}

	// Start all nodes
	for _, nodeConfig := range cfg.Deployer.Nodes {
		nodeName := nodeConfig["name"].(string)
		if err := d.Start(nodeName); err != nil {
			logger.Fatalf("Failed to start node %s: %v", nodeName, err)
		}
		logger.Infof("Successfully started node %s", nodeName)
	}

	logger.Infof("All nodes deployed successfully, %d nodes", len(cfg.Deployer.Nodes))

	// Wait for shutdown signal
	shutdown.WaitForShutdown(syscall.SIGINT, syscall.SIGTERM)

	// Stop all nodes
	for _, nodeConfig := range cfg.Deployer.Nodes {
		nodeName := nodeConfig["name"].(string)
		if err := d.Stop(nodeName); err != nil {
			logger.Errorf("Error stopping node %s: %v", nodeName, err)
		}
	}

	logger.Info("All nodes stopped")
}
