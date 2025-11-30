package main

import (
	"flag"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/BullionBear/sequex/api"
	_ "github.com/BullionBear/sequex/docs"
	"github.com/BullionBear/sequex/env"
	"github.com/BullionBear/sequex/pkg/logger"
	"github.com/BullionBear/sequex/pkg/shutdown"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title PMS API
// @version 1.0
// @description Portfolio Management System API for managing trading positions and portfolios.
// @host localhost:8080
// @BasePath /api/v1

func main() {
	// Define flags
	var port string
	flag.StringVar(&port, "p", "8080", "Port to run the server on")

	// Custom usage function
	flag.Usage = func() {
		logger.Log.Info().Msg(`PMS is a Portfolio Management System API server.

Usage:
  pms [flags]

Flags:
  -p string   Port to run the server on (default "8080")

Examples:
  pms -p 8080
`)
		flag.PrintDefaults()
	}

	// Parse flags
	flag.Parse()

	// Output version information
	logger.Log.Info().
		Str("version", env.Version).
		Str("buildTime", env.BuildTime).
		Str("commitHash", env.CommitHash).
		Msg("PMS started")

	// Initialize shutdown handler
	sd := shutdown.NewShutdown(logger.Log)

	// Setup Gin router
	router := gin.Default()
	router.Use(api.AllowAllCors)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		api.NewPMS(v1)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		logger.Log.Info().Str("port", port).Msg("Starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error().Err(err).Msg("Failed to start server")
			os.Exit(1)
		}
	}()

	// Register shutdown callback for graceful server shutdown
	sd.HookShutdownCallback("http-server", func() {
		logger.Log.Info().Msg("Shutting down HTTP server...")
		if err := srv.Close(); err != nil {
			logger.Log.Error().Err(err).Msg("Error closing HTTP server")
		}
	}, 10*time.Second)

	// Wait for shutdown signal
	sd.WaitForShutdown(syscall.SIGINT, syscall.SIGTERM)
	logger.Log.Info().Msg("PMS server stopped gracefully")
}
