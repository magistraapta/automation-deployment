package main

import (
	"log/slog"
	"os"
	"src/internal/config"
	"src/internal/handler"
	apirouter "src/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	db, err := config.ConnectDatabase()
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}

	router := gin.Default()
	userHandler := handler.NewUserHandler(db)
	apirouter.UserRouter(router, userHandler)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	slog.Info("Starting server on port 8080")
	if err := router.Run(":8080"); err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
