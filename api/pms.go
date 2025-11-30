package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewPMS(rg *gin.RouterGroup) {
	rg.POST("/symbol", createSymbol)
}

type CreateSymbolRequest struct {
	Exchange string `json:"exchange" binding:"required"`
	Symbol   string `json:"symbol" binding:"required"`
}

type CreateSymbolResponse struct {
	SymbolID string `json:"symbol_id"`
}

// @Summary Create a symbol
// @Description Create a new symbol
// @Accept json
// @Produce json
// @Success 200 {object} string "Symbol"
// @Router /symbol [post]
func createSymbol(c *gin.Context) {
	var req CreateSymbolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Symbol created successfully"})
}
