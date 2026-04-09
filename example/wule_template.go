package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	// TODO: Ganti dengan actual module name dari go.mod
	// _ "your-wule-app/docs"
)

// @title WULE API untuk Limit Management System  
// @version 1.0
// @description WULE API untuk Limit Management System
// @host localhost:8080
// @BasePath /api/v1
func main() {
	r := gin.Default()
	
	api := r.Group("/api/v1")
	{
		api.GET("/health", healthCheck)
		api.GET("/users", getUsers)
		api.POST("/users", createUser)
		api.GET("/limits", getLimits)
	}
	
	r.Run(":8080")
}

// @Summary Health check
// @Description Check if API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "WULE API"})
}

// @Summary Get all users
// @Description Retrieve list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /api/v1/users [get]
func getUsers(c *gin.Context) {
	users := []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}
	c.JSON(http.StatusOK, users)
}

// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User data"
// @Success 201 {object} User
// @Router /api/v1/users [post]
func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// @Summary Get limits
// @Description Get limit management data
// @Tags limits
// @Accept json
// @Produce json
// @Success 200 {array} Limit
// @Router /api/v1/limits [get]
func getLimits(c *gin.Context) {
	limits := []Limit{
		{ID: 1, Type: "daily", Value: 1000},
		{ID: 2, Type: "monthly", Value: 30000},
	}
	c.JSON(http.StatusOK, limits)
}

type User struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

type Limit struct {
	ID    int    `json:"id" example:"1"`
	Type  string `json:"type" example:"daily"`
	Value int    `json:"value" example:"1000"`
}