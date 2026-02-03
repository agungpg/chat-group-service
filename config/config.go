package config

import (
	"log"
	"os"

	"github.com/agungpg/group-chat-service/pkg/utils"
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

type Config struct {
	Storage StorageConfig
}

type StorageConfig struct {
	PublicBucket  string
	PrivateBucket string
}

// Prefer required env for critical values
func Load() Config {
	publicBucket := utils.GetEnvOrDefault("STORAGE_BUCKET_PUBLIC", "")
	privateBucket := utils.GetEnvOrDefault("STORAGE_BUCKET_PRIVATE", "")

	if publicBucket == "" {
		log.Fatal("missing env: STORAGE_BUCKET_PUBLIC")
	}
	if privateBucket == "" {
		log.Fatal("missing env: STORAGE_BUCKET_PRIVATE")
	}

	return Config{
		Storage: StorageConfig{
			PublicBucket:  publicBucket,
			PrivateBucket: privateBucket,
		},
	}
}
