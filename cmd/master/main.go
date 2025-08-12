package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/BullionBear/sequex/api"
	_ "github.com/BullionBear/sequex/docs"
	"github.com/BullionBear/sequex/env"
	"github.com/BullionBear/sequex/internal/config"
	_ "github.com/BullionBear/sequex/internal/nodeimpl/init" // Import to register all nodes
	"github.com/BullionBear/sequex/internal/nodeimpl/master"
	"github.com/BullionBear/sequex/pkg/log"
	"github.com/BullionBear/sequex/pkg/shutdown"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var logger log.Logger

func main() {
	// Parse command line arguments
	var configFile string
	flag.StringVar(&configFile, "c", "config/master/app.yml", "Configuration file path")
	flag.Parse()

	fmt.Println("Starting services with BuildTime:", env.BuildTime)
	fmt.Println("Starting services with Version:", env.Version)
	fmt.Println("Starting services with CommitHash:", env.CommitHash)

	// Load configuration
	cfg, err := config.LoadConfig[config.MasterConfig](configFile)
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

	// Create shutdown
	shutdown := shutdown.NewShutdown(logger)

	nc, err := config.CreateNATSConnection(cfg.Nats.URL)
	if err != nil {
		fmt.Printf("Failed to create NATS connection: %v\n", err)
		os.Exit(1)
	}

	masterRPCClient := master.NewMasterRPCClient(nc)

	rg := gin.New()
	rg.Use(gin.Logger())
	rg.Use(api.AllowAllCors)
	v1rg := rg.Group("/v1", gin.Recovery())
	api.NewNode(v1rg, masterRPCClient)

	go func() {
		rg.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		logger.Infof("Server started on %s:%d", cfg.App.Host, cfg.App.Port)
		rg.Run(fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port))
	}()

	shutdown.WaitForShutdown(syscall.SIGINT, syscall.SIGTERM)
}
