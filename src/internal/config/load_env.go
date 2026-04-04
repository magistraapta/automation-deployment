package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	wd, err := os.Getwd()
	if err != nil {
		slog.Warn("could not get working directory", "error", err)
		return
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		p := filepath.Join(dir, ".env")
		if _, err := os.Stat(p); err != nil {
			if parent := filepath.Dir(dir); parent == dir {
				slog.Warn("no .env found (searched from working directory upward)", "start", wd)
				return
			}
			continue
		}
		if err := godotenv.Load(p); err != nil {
			slog.Error("failed to load .env", "path", p, "error", err)
			return
		}
		slog.Info("loaded .env", "path", p)
		return
	}
}
