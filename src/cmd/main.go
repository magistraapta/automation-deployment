package main

import (
	"log/slog"
	"src/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	userHandler := handler.NewUserHandler()
	router.GET("/users", userHandler.GetUsers)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	slog.Info("Starting server on port 8080")
	if err := router.Run(":8080"); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
