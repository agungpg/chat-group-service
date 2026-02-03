package main

import (
	"context"
	"fmt"
	"os"

	"github.com/agungpg/group-chat-service/config"
	"github.com/agungpg/group-chat-service/internal/app"
	"github.com/agungpg/group-chat-service/pkg/database"
	"github.com/agungpg/group-chat-service/pkg/storage"
)

func main() {
	config.LoadEnv()

	if err := database.Connect(); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	cfg := config.Load()

	if err := storage.Connect(context.Background()); err != nil {
		fmt.Printf("Failed to connect to storage: %v\n", err)
		os.Exit(1)
	}

	container := app.NewContainer(database.DB, storage.S3Client, &cfg.Storage)

	server := app.NewServer(container)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := server.Listen(":" + port); err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
}
