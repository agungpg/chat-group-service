package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads a local .env file if present; otherwise rely on environment variables.
func LoadEnv() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Failed to load .env file: %v", err)
		}
	}
}
