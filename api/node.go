package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewNode(rg *gin.RouterGroup) {
	prg := rg.Group("/node")
	prg.GET("/list", listNodes)
	prg.POST("/register", registerNode)
	prg.GET("/:name", getNode)
	prg.DELETE("/:name", unregisterNode)
}

func listNodes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func registerNode(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func getNode(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func unregisterNode(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
