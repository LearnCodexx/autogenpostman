package main

import "github.com/gin-gonic/gin"

// @title Flexible API Demo
// @version 1.0
// @description Clean API definition separate from business logic
// @host localhost:8080
// @BasePath /api/v1

// This file is dedicated ONLY for API documentation
// Keep business logic in main.go

// @Summary Health check
// @Description Check if API service is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/health [get]
func healthCheck(c *gin.Context) {
    c.JSON(200, gin.H{"status": "healthy", "service": "flexible-api"})
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
    c.JSON(200, users)
}

// @Summary Create new user
// @Description Create a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param user body User true "User data"
// @Success 201 {object} User
// @Router /api/v1/users [post]
func createUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(201, user)
}

// @Summary Get user by ID
// @Description Get specific user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Router /api/v1/users/{id} [get]
func getUserByID(c *gin.Context) {
    id := c.Param("id")
    user := User{ID: 1, Name: "John Doe", Email: "john@example.com"}
    c.JSON(200, user)
}

type User struct {
    ID    int    `json:"id" example:"1"`
    Name  string `json:"name" example:"John Doe"`
    Email string `json:"email" example:"john@example.com"`
}