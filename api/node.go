package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewNode(rg *gin.RouterGroup) {
	rg.GET("/nodes", listNodes)
	rg.GET("/node/:name", getNode)
}

func listNodes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func getNode(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
