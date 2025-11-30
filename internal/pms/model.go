package pms

import "time"

// Portfolio represents a trading portfolio
type Portfolio struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Position represents a trading position within a portfolio
type Position struct {
	ID          string    `json:"id"`
	PortfolioID string    `json:"portfolio_id"`
	Symbol      string    `json:"symbol"`
	Side        string    `json:"side"` // LONG or SHORT
	EntryPrice  float64   `json:"entry_price"`
	Quantity    float64   `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreatePortfolioParams contains parameters for creating a portfolio
type CreatePortfolioParams struct {
	Name        string
	Description string
}

// UpdatePortfolioParams contains parameters for updating a portfolio
type UpdatePortfolioParams struct {
	Name        string
	Description string
}

// CreatePositionParams contains parameters for creating a position
type CreatePositionParams struct {
	PortfolioID string
	Symbol      string
	Side        string
	EntryPrice  float64
	Quantity    float64
}

// UpdatePositionParams contains parameters for updating a position
type UpdatePositionParams struct {
	Symbol     string
	Side       string
	EntryPrice float64
	Quantity   float64
}
