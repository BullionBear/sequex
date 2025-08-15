package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/BullionBear/sequex/env"
	"github.com/BullionBear/sequex/internal/config"
	"github.com/nats-io/nats.go"

	_ "github.com/BullionBear/sequex/internal/nodeimpl/init" // Import to register all nodes
	"github.com/BullionBear/sequex/pkg/eventbus"
	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/node"
	"github.com/BullionBear/sequex/pkg/shutdown"
	"github.com/spf13/cobra"
)

var (
	logger     log.Logger
	configFile string
)

func main() {
	// Initialize logger
	logger = log.New(
		log.WithLevel(log.LevelInfo),
		log.WithOutput(os.Stdout),
		log.WithEncoder(log.NewTextEncoder()),
	)

	// Create root command
	rootCmd := &cobra.Command{
		Use:   "sqx",
		Short: "Sequex Node - A distributed computing node",
		Long: `Sequex Node is a distributed computing node that can run as a server
or interact with other nodes as a client.

Examples:
  sqx serve -c config.yml     # Start a server with config
  sqx call rng --server localhost:8080 --input 10  # Call RNG service
  sqx list --server localhost:8080                  # List available services`,
		Version: fmt.Sprintf("Version: %s\nBuild Time: %s\nCommit Hash: %s",
			env.Version, env.BuildTime, env.CommitHash),
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Configuration file path")

	// Create serve command (server group)
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the node as a server",
		Long:  "Start the Sequex node as a server to handle incoming requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer()
		},
	}

	// Create client commands
	callCmd := &cobra.Command{
		Use:   "call [service]",
		Short: "Call a specific service",
		Long:  "Call a specific service on a remote node",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return callService(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available services",
		Long:  "List all available services on a remote node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listServices()
		},
	}

	// Add flags for client commands
	callCmd.Flags().String("server", "localhost:8080", "Server address to connect to")
	callCmd.Flags().String("input", "", "Input data for the service")
	listCmd.Flags().String("server", "localhost:8080", "Server address to connect to")

	// Add commands to root
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(callCmd)
	rootCmd.AddCommand(listCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runServer starts the node as a server
func runServer() error {
	logger.Info("Starting node server",
		log.String("build_time", env.BuildTime),
		log.String("version", env.Version),
		log.String("commit_hash", env.CommitHash),
		log.String("config_file", configFile),
	)

	// Load global config
	globalCfg, err := config.LoadConfig[config.GlobalConfig](config.PathGlobalConfig)
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	cfg, err := config.LoadConfig[node.NodeConfig](configFile)
	if err != nil {
		return err
	}

	conn, err := nats.Connect(globalCfg.EventBus.Url)
	if err != nil {
		return err
	}
	eventBus := eventbus.NewEventBus(conn)
	shutdown := shutdown.NewShutdown(logger)

	nodeInstance, err := node.CreateNode(cfg.Type, eventBus, cfg, logger)
	if err != nil {
		return err
	}

	if err := nodeInstance.Start(); err != nil {
		return err
	}

	logger.Info("Server started successfully",
		log.String("name", cfg.Name),
		log.String("type", cfg.Type),
	)
	shutdown.HookShutdownCallback(fmt.Sprintf("(%s).shutdown", cfg.Name), func() {
		nodeInstance.Shutdown()
	}, 10*time.Second)
	shutdown.WaitForShutdown(os.Interrupt, syscall.SIGTERM)

	return nil
}

// callService calls a specific service on a remote node
func callService(serviceName string) error {
	// For now, use default values since we can't access flags from here
	server := "localhost:8080"
	input := ""

	logger.Info("Calling service",
		log.String("service", serviceName),
		log.String("server", server),
		log.String("input", input),
	)

	// TODO: Implement service call logic
	logger.Info("TODO: Implement service call functionality")

	return nil
}

// listServices lists available services on a remote node
func listServices() error {
	// For now, use default values since we can't access flags from here
	server := "localhost:8080"

	logger.Info("Listing services",
		log.String("server", server),
	)

	// TODO: Implement service listing logic
	logger.Info("TODO: Implement service listing functionality")

	return nil
}
