// Package server provides the HTTP server for the Do Worker Service.
package server

import (
	"os"
	"time"

	"github.com/do-focus/worker/internal/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server.
type Server struct {
	router *gin.Engine
	db     db.Adapter
}

// New creates a new server instance.
func New() (*Server, error) {
	// Configure database
	dbCfg := db.Config{
		Type:     os.Getenv("DO_DB_TYPE"),
		Path:     os.Getenv("DO_DB_PATH"),
		Host:     os.Getenv("DO_DB_HOST"),
		Port:     os.Getenv("DO_DB_PORT"),
		User:     os.Getenv("DO_DB_USER"),
		Password: os.Getenv("DO_DB_PASSWORD"),
		Database: os.Getenv("DO_DB_DATABASE"),
	}

	// Set defaults
	if dbCfg.Type == "" {
		dbCfg.Type = "sqlite"
	}
	if dbCfg.Path == "" {
		dbCfg.Path = ".do/memory.db"
	}
	if dbCfg.Port == "" {
		dbCfg.Port = "3306"
	}

	// Initialize database
	dbAdapter, err := db.New(dbCfg)
	if err != nil {
		return nil, err
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// CORS configuration for web viewer
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3777", "http://127.0.0.1:3777"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s := &Server{
		router: router,
		db:     dbAdapter,
	}

	// Setup routes
	s.setupRoutes()

	return s, nil
}

// Run starts the HTTP server.
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// Close closes the server and its resources.
func (s *Server) Close() error {
	return s.db.Close()
}
