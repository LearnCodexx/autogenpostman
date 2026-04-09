package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @title API Documentation
// @version 1.0
// @description This is a sample API documentation
// @host localhost:8080
// @BasePath /api/v1
func main() {
	r := gin.Default()

	// @Summary Health check
	// @Description Get health status of the API
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router /health [get]
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "Service is healthy",
		})
	})

	// @Summary Get users
	// @Description Get all users
	// @Tags users  
	// @Accept json
	// @Produce json
	// @Success 200 {array} map[string]interface{}
	// @Router /api/v1/users [get]
	r.GET("/api/v1/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, []map[string]interface{}{
			{"id": 1, "name": "John Doe", "email": "john@example.com"},
			{"id": 2, "name": "Jane Smith", "email": "jane@example.com"},
		})
	})

	log.Println("Starting server on :8080")
	r.Run(":8080")
}