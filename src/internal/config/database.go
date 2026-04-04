package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"src/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	dbURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is empty: set it in the environment or in a .env file (loaded from .env or ../.env when using go run from cmd/)")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		slog.Error("Failed to connect to database: ", "error", err)
		return nil, err
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		slog.Error("Failed to migrate database: ", "error", err)
		return nil, err
	}

	slog.Info("Database migrated successfully")

	return db, nil
}
