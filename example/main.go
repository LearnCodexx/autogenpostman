package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "example-api/docs"
)

type User struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
}

type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"Service is healthy"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// @title           Example API
// @version         1.0
// @description     Sample API for demonstrating autogenpostman package
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthCheck)
		v1.GET("/users", getUsers)
		v1.GET("/users/:id", getUserByID)
		v1.POST("/users", createUser)
	}

	r.Run(":8080")
}

// @Summary      Health check
// @Tags         health
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "ok",
		Message: "Service is healthy",
	})
}

// @Summary      List users
// @Tags         users
// @Success      200  {array}   User
// @Router       /users [get]
func getUsers(c *gin.Context) {
	users := []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}
	c.JSON(http.StatusOK, users)
}

// @Summary      Get a user
// @Tags         users
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  User
// @Router       /users/{id} [get]
func getUserByID(c *gin.Context) {
	id := c.Param("id")
	if id == "1" {
		c.JSON(http.StatusOK, User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		})
	} else {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "User not found",
		})
	}
}

// @Summary      Create a user
// @Tags         users
// @Param        user  body      User  true  "User data"
// @Success      201   {object}  User
// @Router       /users [post]
func createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid JSON",
		})
		return
	}

	user.ID = 99
	c.JSON(http.StatusCreated, user)
}
