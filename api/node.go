package api

import (
	"net/http"

	"github.com/BullionBear/sequex/internal/nodeimpl/master"
	"github.com/gin-gonic/gin"
)

func NewNode(rg *gin.RouterGroup, masterRPCClient *master.MasterRPCClient) {
	rg.GET("/nodes", listNodes)
	rg.GET("/node/:name", getNode)
	rg.POST("/node/register", registerNode)
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

type RegisterNodeRequest struct {
	Name string `json:"name"`
}

// @Summary Register a node
// @Description Register a node
// @Accept json
// @Produce json
// @Success 200 {object} string "Node"
// @Router /node/register [post]
func registerNode(c *gin.Context) {
	var req RegisterNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Node registered successfully"})
}
