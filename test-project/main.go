package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// @title		API
// @version		1.0
// @description	Auto-generated API
// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	r := gin.Default()
	
	// @Summary		Health check
	// @Tags		health
	// @Success		200	{object}	map[string]string
	// @Router		/health [get]
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	
	r.Run(":8080")
}
