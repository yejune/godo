// Package main is the entry point for the Do Worker Service.
// Worker Service handles background DB operations for memory management
// without consuming Claude tokens.
package main

import (
	"log"
	"os"

	"github.com/do-focus/worker/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file if exists
	_ = godotenv.Load()

	// Get port from environment or use default
	port := os.Getenv("DO_WORKER_PORT")
	if port == "" {
		port = "3778"
	}

	// Initialize and start server
	srv, err := server.New()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	log.Printf("Do Worker Service starting on port %s", port)
	if err := srv.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
