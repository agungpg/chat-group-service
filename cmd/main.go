package main

import (
	"fmt"
	"os"

	"github.com/agungpg/group-chat-service/config"
	"github.com/agungpg/group-chat-service/internal/app"
	"github.com/agungpg/group-chat-service/pkg/database"
)

func main() {
	config.LoadEnv()

	if err := database.Connect(); err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	container := app.NewContainer(database.DB)

	server := app.NewServer(container)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := server.Listen(":" + port); err != nil {
		fmt.Printf("Failed to start server: %v", err)
	}
}
