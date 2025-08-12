package api

import (
	"net/http"

	"github.com/BullionBear/sequex/internal/nodeimpl/master"
	"github.com/gin-gonic/gin"
)

func NewNode(rg *gin.RouterGroup, masterRPCClient *master.MasterRPCClient) {
	rg.GET("/nodes", listNodes)
	rg.GET("/node/:name", getNode)
}

// @Summary List all nodes
// @Description List all nodes
// @Accept json
// @Produce json
// @Success 200 {array} string "List of nodes"
// @Router /nodes [get]
func listNodes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

// @Summary Get a node
// @Description Get a node
// @Accept json
// @Produce json
// @Success 200 {object} string "Node"
// @Router /node/{name} [get]
func getNode(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
