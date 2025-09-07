package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/BullionBear/sequex/env"
	"github.com/BullionBear/sequex/internal/config"
	"github.com/nats-io/nats.go"

	_ "github.com/BullionBear/sequex/internal/nodeimpl/v1/init" // Import to register all nodes
	"github.com/BullionBear/sequex/pkg/eventbus"
	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/node"
	"github.com/BullionBear/sequex/pkg/shutdown"
	"google.golang.org/protobuf/proto"

	pbCommon "github.com/BullionBear/sequex/internal/model/protobuf/common"
)

var (
	logger     log.Logger
	globalCfg  *config.GlobalConfig
	err        error
	configFile string
	version    bool
)

func main() {
	// Define flags
	flag.StringVar(&configFile, "config", "", "Configuration file path")
	flag.StringVar(&configFile, "c", "", "Configuration file path (shorthand)")
	flag.BoolVar(&version, "version", false, "Show version information")
	flag.BoolVar(&version, "v", false, "Show version information (shorthand)")

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Sequex Node is a distributed computing node that can run as a server
or interact with other nodes as a client.

Usage:
  sqx <command> [flags]

Commands:
  serve    Start the node as a server
  call     Call a specific method on a remote node

Examples:
  sqx serve -c config.yml     		# Start a server with config
  sqx call metadata -c config.yml   # Call RNG service

Flags:
`)
		flag.PrintDefaults()
	}

	// Parse flags
	flag.Parse()

	// Handle version flag
	if version {
		fmt.Printf("Version: %s\nBuild Time: %s\nCommit Hash: %s\n",
			env.Version, env.BuildTime, env.CommitHash)
		return
	}

	// Check if we have a command
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: command required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Initialize logger
	logger = log.DefaultLogger(
		log.WithLevel(log.LevelInfo),
		log.WithOutput(os.Stdout),
		log.WithEncoder(log.NewTextEncoder()),
	)

	// Load global config to get NATS URL
	globalCfg, err = config.LoadConfig[config.GlobalConfig](config.PathGlobalConfig)
	if err != nil {
		logger.Error("failed to load global config", log.Error(err))
		os.Exit(1)
	}

	// Handle commands
	command := args[0]
	switch command {
	case "serve":
		if err := runServer(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "call":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error: call command requires a service name\n")
			fmt.Fprintf(os.Stderr, "Usage: sqx call <service> [method]\n")
			fmt.Fprintf(os.Stderr, "Supported services: metadata, status, params\n")
			os.Exit(1)
		}
		serviceName := args[1]
		method := "Supported methods: metadata, status, params"
		if len(args) > 2 {
			method = args[2]
		}
		if err := callService(serviceName, method); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n\n", command)
		flag.Usage()
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
func callService(serviceName, method string) error {
	logger = log.DefaultLogger(
		log.WithLevel(log.LevelInfo),
		log.WithOutput(os.Stderr),
		log.WithEncoder(log.NewTextEncoder()),
	)
	// Log only in debug mode or remove for cleaner output
	// logger.Info("Calling service",
	// 	log.String("service", serviceName),
	// 	log.String("method", method),
	// 	log.String("config_file", configFile),
	// )

	// Load the node configuration to get RPC endpoints
	cfg, err := config.LoadConfig[node.NodeConfig](configFile)
	if err != nil {
		return fmt.Errorf("failed to load node config: %w", err)
	}

	// Connect to NATS
	conn, err := nats.Connect(globalCfg.EventBus.Url)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	defer conn.Close()

	eventBus := eventbus.NewEventBus(conn)

	// Map service names to RPC endpoints
	var rpcEndpoint string
	switch serviceName {
	case "metadata":
		rpcEndpoint = cfg.Rpc["req_metadata"]
	case "status":
		rpcEndpoint = cfg.Rpc["req_status"]
	case "params":
		rpcEndpoint = cfg.Rpc["req_parameters"]
	default:
		return fmt.Errorf("unknown service: %s. Supported services: metadata, status, params", serviceName)
	}
	// Create appropriate request based on service
	var request proto.Message
	var responseFactory func() proto.Message

	switch serviceName {
	case "metadata":
		request = &pbCommon.MetadataRequest{Id: time.Now().UnixNano()}
		responseFactory = func() proto.Message { return &pbCommon.MetadataResponse{} }
	case "status":
		request = &pbCommon.StatusRequest{Id: time.Now().UnixNano()}
		responseFactory = func() proto.Message { return &pbCommon.StatusResponse{} }
	case "params":
		request = &pbCommon.ParametersRequest{Id: time.Now().UnixNano()}
		responseFactory = func() proto.Message { return &pbCommon.ParametersResponse{} }
	}

	// Make RPC call
	// logger.Info("Making RPC call",
	// 	log.String("endpoint", rpcEndpoint),
	// 	log.String("service", serviceName),
	// )

	response, err := eventBus.CallRPC(rpcEndpoint, request, responseFactory, time.Second)
	if err != nil {
		return fmt.Errorf("RPC call failed: %w", err)
	}

	// Handle response based on service type
	switch serviceName {
	case "metadata":
		if resp, ok := response.(*pbCommon.MetadataResponse); ok {
			if resp.Code != pbCommon.ErrorCode_ERROR_CODE_OK {
				return fmt.Errorf("metadata request failed: %s", resp.Message)
			}
			metadataResult := map[string]interface{}{
				"name":       resp.Name,
				"created_at": resp.CreatedAt,
				"emit":       resp.Emit,
				"on":         resp.On,
				"rpc":        resp.Rpc,
			}
			jsonData, _ := json.MarshalIndent(metadataResult, "", "  ")
			fmt.Println(string(jsonData))
		}
	case "status":
		if resp, ok := response.(*pbCommon.StatusResponse); ok {
			if resp.Code != pbCommon.ErrorCode_ERROR_CODE_OK {
				return fmt.Errorf("status request failed: %s", resp.Message)
			}
			var statusData map[string]interface{}
			if err := json.Unmarshal(resp.Status, &statusData); err != nil {
				return fmt.Errorf("failed to unmarshal status data: %w", err)
			}
			jsonData, _ := json.MarshalIndent(statusData, "", "  ")
			fmt.Println(string(jsonData))
		}
	case "params":
		if resp, ok := response.(*pbCommon.ParametersResponse); ok {
			if resp.Code != pbCommon.ErrorCode_ERROR_CODE_OK {
				return fmt.Errorf("parameters request failed: %s", resp.Message)
			}
			var paramsData map[string]interface{}
			if err := json.Unmarshal(resp.Parameters, &paramsData); err != nil {
				return fmt.Errorf("failed to unmarshal parameters data: %w", err)
			}
			jsonData, _ := json.MarshalIndent(paramsData, "", "  ")
			fmt.Println(string(jsonData))
		}
	}

	return nil
}
