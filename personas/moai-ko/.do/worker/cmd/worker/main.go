// Package main provides the CLI entry point for the Do Worker Service.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/do-focus/worker/internal/server"
	"github.com/joho/godotenv"
)

var (
	version = "0.1.0"
	commit  = "dev"
)

// getDefaultDBPath returns the default database path.
// Priority: DO_DB_PATH env var > global path (~/.do/memory.db)
func getDefaultDBPath() string {
	// Environment variable takes priority
	if path := os.Getenv("DO_DB_PATH"); path != "" {
		return path
	}

	// Global path (~/.do/memory.db)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home not available
		return ".do/memory.db"
	}
	globalDir := filepath.Join(homeDir, ".do")
	if err := os.MkdirAll(globalDir, 0755); err != nil {
		log.Printf("Warning: failed to create global dir: %v", err)
	}
	return filepath.Join(globalDir, "memory.db")
}

func main() {
	// Parse flags
	showVersion := flag.Bool("version", false, "Show version")
	port := flag.String("port", "", "Port to listen on (default: 3778)")
	dbType := flag.String("db", "", "Database type: sqlite or mysql")
	dbPath := flag.String("db-path", "", "Database path for SQLite")
	flag.Parse()

	if *showVersion {
		fmt.Printf("do-worker %s (%s)\n", version, commit)
		os.Exit(0)
	}

	// Load environment variables
	_ = godotenv.Load()

	// Override with flags
	if *port != "" {
		os.Setenv("DO_WORKER_PORT", *port)
	}
	if *dbType != "" {
		os.Setenv("DO_DB_TYPE", *dbType)
	}
	if *dbPath != "" {
		os.Setenv("DO_DB_PATH", *dbPath)
	}

	// Set default DB path if not specified
	if os.Getenv("DO_DB_PATH") == "" {
		os.Setenv("DO_DB_PATH", getDefaultDBPath())
	}

	// Get port
	listenPort := os.Getenv("DO_WORKER_PORT")
	if listenPort == "" {
		listenPort = "3778"
	}

	// Set version for health endpoint
	server.Version = version

	// Initialize server
	srv, err := server.New()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Handle shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-done
		log.Println("Shutting down...")
		if err := srv.Close(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		os.Exit(0)
	}()

	// Start server
	log.Printf("Do Worker Service v%s starting on port %s", version, listenPort)
	if err := srv.Run(":" + listenPort); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
