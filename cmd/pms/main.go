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

// registerPMSRoutes registers PMS-specific routes
func registerPMSRoutes(rg *gin.RouterGroup) {
	// Portfolio routes
	rg.GET("/portfolios", listPortfolios)
	rg.GET("/portfolio/:id", getPortfolio)
	rg.POST("/portfolio", createPortfolio)
	rg.PUT("/portfolio/:id", updatePortfolio)
	rg.DELETE("/portfolio/:id", deletePortfolio)

	// Position routes
	rg.GET("/portfolio/:id/positions", listPositions)
	rg.GET("/position/:id", getPosition)
	rg.POST("/position", createPosition)
	rg.PUT("/position/:id", updatePosition)
	rg.DELETE("/position/:id", deletePosition)

	// Health check
	rg.GET("/health", healthCheck)
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status" example:"ok"`
	Timestamp int64  `json:"timestamp" example:"1699999999"`
}

// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().Unix(),
	})
}

// Portfolio represents a trading portfolio
type Portfolio struct {
	ID          string    `json:"id" example:"portfolio-001"`
	Name        string    `json:"name" example:"Main Portfolio"`
	Description string    `json:"description" example:"Primary trading portfolio"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreatePortfolioRequest represents a request to create a portfolio
type CreatePortfolioRequest struct {
	Name        string `json:"name" binding:"required" example:"Main Portfolio"`
	Description string `json:"description" example:"Primary trading portfolio"`
}

// Position represents a trading position
type Position struct {
	ID          string  `json:"id" example:"position-001"`
	PortfolioID string  `json:"portfolio_id" example:"portfolio-001"`
	Symbol      string  `json:"symbol" example:"BTCUSDT"`
	Side        string  `json:"side" example:"LONG"`
	EntryPrice  float64 `json:"entry_price" example:"45000.50"`
	Quantity    float64 `json:"quantity" example:"0.5"`
	CurrentPnL  float64 `json:"current_pnl" example:"1250.00"`
}

// CreatePositionRequest represents a request to create a position
type CreatePositionRequest struct {
	PortfolioID string  `json:"portfolio_id" binding:"required" example:"portfolio-001"`
	Symbol      string  `json:"symbol" binding:"required" example:"BTCUSDT"`
	Side        string  `json:"side" binding:"required" example:"LONG"`
	EntryPrice  float64 `json:"entry_price" binding:"required" example:"45000.50"`
	Quantity    float64 `json:"quantity" binding:"required" example:"0.5"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"resource not found"`
}

// @Summary List all portfolios
// @Description Get a list of all portfolios
// @Tags portfolios
// @Accept json
// @Produce json
// @Success 200 {array} Portfolio
// @Router /portfolios [get]
func listPortfolios(c *gin.Context) {
	// TODO: Implement actual portfolio listing from database
	portfolios := []Portfolio{
		{
			ID:          "portfolio-001",
			Name:        "Main Portfolio",
			Description: "Primary trading portfolio",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	c.JSON(http.StatusOK, portfolios)
}

// @Summary Get a portfolio
// @Description Get a portfolio by ID
// @Tags portfolios
// @Accept json
// @Produce json
// @Param id path string true "Portfolio ID"
// @Success 200 {object} Portfolio
// @Failure 404 {object} ErrorResponse
// @Router /portfolio/{id} [get]
func getPortfolio(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement actual portfolio retrieval from database
	portfolio := Portfolio{
		ID:          id,
		Name:        "Main Portfolio",
		Description: "Primary trading portfolio",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	c.JSON(http.StatusOK, portfolio)
}

// @Summary Create a portfolio
// @Description Create a new portfolio
// @Tags portfolios
// @Accept json
// @Produce json
// @Param request body CreatePortfolioRequest true "Portfolio creation request"
// @Success 201 {object} Portfolio
// @Failure 400 {object} ErrorResponse
// @Router /portfolio [post]
func createPortfolio(c *gin.Context) {
	var req CreatePortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// TODO: Implement actual portfolio creation in database
	portfolio := Portfolio{
		ID:          "portfolio-new",
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	c.JSON(http.StatusCreated, portfolio)
}

// @Summary Update a portfolio
// @Description Update an existing portfolio
// @Tags portfolios
// @Accept json
// @Produce json
// @Param id path string true "Portfolio ID"
// @Param request body CreatePortfolioRequest true "Portfolio update request"
// @Success 200 {object} Portfolio
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /portfolio/{id} [put]
func updatePortfolio(c *gin.Context) {
	id := c.Param("id")
	var req CreatePortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// TODO: Implement actual portfolio update in database
	portfolio := Portfolio{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	c.JSON(http.StatusOK, portfolio)
}

// @Summary Delete a portfolio
// @Description Delete a portfolio by ID
// @Tags portfolios
// @Accept json
// @Produce json
// @Param id path string true "Portfolio ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /portfolio/{id} [delete]
func deletePortfolio(c *gin.Context) {
	// id := c.Param("id")
	// TODO: Implement actual portfolio deletion from database
	c.Status(http.StatusNoContent)
}

// @Summary List positions in a portfolio
// @Description Get all positions in a specific portfolio
// @Tags positions
// @Accept json
// @Produce json
// @Param id path string true "Portfolio ID"
// @Success 200 {array} Position
// @Failure 404 {object} ErrorResponse
// @Router /portfolio/{id}/positions [get]
func listPositions(c *gin.Context) {
	portfolioID := c.Param("id")
	// TODO: Implement actual position listing from database
	positions := []Position{
		{
			ID:          "position-001",
			PortfolioID: portfolioID,
			Symbol:      "BTCUSDT",
			Side:        "LONG",
			EntryPrice:  45000.50,
			Quantity:    0.5,
			CurrentPnL:  1250.00,
		},
	}
	c.JSON(http.StatusOK, positions)
}

// @Summary Get a position
// @Description Get a position by ID
// @Tags positions
// @Accept json
// @Produce json
// @Param id path string true "Position ID"
// @Success 200 {object} Position
// @Failure 404 {object} ErrorResponse
// @Router /position/{id} [get]
func getPosition(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement actual position retrieval from database
	position := Position{
		ID:          id,
		PortfolioID: "portfolio-001",
		Symbol:      "BTCUSDT",
		Side:        "LONG",
		EntryPrice:  45000.50,
		Quantity:    0.5,
		CurrentPnL:  1250.00,
	}
	c.JSON(http.StatusOK, position)
}

// @Summary Create a position
// @Description Create a new position
// @Tags positions
// @Accept json
// @Produce json
// @Param request body CreatePositionRequest true "Position creation request"
// @Success 201 {object} Position
// @Failure 400 {object} ErrorResponse
// @Router /position [post]
func createPosition(c *gin.Context) {
	var req CreatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// TODO: Implement actual position creation in database
	position := Position{
		ID:          "position-new",
		PortfolioID: req.PortfolioID,
		Symbol:      req.Symbol,
		Side:        req.Side,
		EntryPrice:  req.EntryPrice,
		Quantity:    req.Quantity,
		CurrentPnL:  0,
	}
	c.JSON(http.StatusCreated, position)
}

// @Summary Update a position
// @Description Update an existing position
// @Tags positions
// @Accept json
// @Produce json
// @Param id path string true "Position ID"
// @Param request body CreatePositionRequest true "Position update request"
// @Success 200 {object} Position
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /position/{id} [put]
func updatePosition(c *gin.Context) {
	id := c.Param("id")
	var req CreatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// TODO: Implement actual position update in database
	position := Position{
		ID:          id,
		PortfolioID: req.PortfolioID,
		Symbol:      req.Symbol,
		Side:        req.Side,
		EntryPrice:  req.EntryPrice,
		Quantity:    req.Quantity,
		CurrentPnL:  0,
	}
	c.JSON(http.StatusOK, position)
}

// @Summary Delete a position
// @Description Delete a position by ID
// @Tags positions
// @Accept json
// @Produce json
// @Param id path string true "Position ID"
// @Success 204 "No Content"
// @Failure 404 {object} ErrorResponse
// @Router /position/{id} [delete]
func deletePosition(c *gin.Context) {
	// id := c.Param("id")
	// TODO: Implement actual position deletion from database
	c.Status(http.StatusNoContent)
}
