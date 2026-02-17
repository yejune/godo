// Package db provides database abstraction for the Do Worker Service.
package db

import (
	"context"

	"github.com/do-focus/worker/pkg/models"
)

// Adapter defines the database interface.
type Adapter interface {
	// Health checks database connectivity.
	Health(ctx context.Context) error

	// Close closes the database connection.
	Close() error

	// Session operations
	CreateSession(ctx context.Context, session *models.Session) error
	GetSession(ctx context.Context, id string) (*models.Session, error)
	GetLatestSession(ctx context.Context, userName string) (*models.Session, error)
	EndSession(ctx context.Context, id string, summary string) error

	// Observation operations
	CreateObservation(ctx context.Context, obs *models.Observation) error
	GetObservations(ctx context.Context, sessionID string) ([]models.Observation, error)
	GetRecentObservations(ctx context.Context, userName string, limit int) ([]models.Observation, error)
	GetObservationsFiltered(ctx context.Context, sessionID string, obsType string, limit int, offset int) ([]models.Observation, error)
	SearchObservations(ctx context.Context, query string, limit int) ([]models.Observation, error)

	// Summary operations
	CreateSummary(ctx context.Context, summary *models.Summary) error
	GetSummaries(ctx context.Context, summaryType string, limit int) ([]models.Summary, error)
	GetAllSummaries(ctx context.Context, days int, limit int) ([]models.Summary, error)
	GetLatestSummary(ctx context.Context, userName string) (*models.Summary, error)

	// Plan operations
	CreatePlan(ctx context.Context, plan *models.Plan) error
	GetActivePlan(ctx context.Context, userName string) (*models.Plan, error)
	GetAllPlans(ctx context.Context, sessionID string, limit int) ([]models.Plan, error)
	UpdatePlanStatus(ctx context.Context, id int64, status string) error

	// Session list operations
	GetRecentSessions(ctx context.Context, limit int) ([]models.Session, error)

	// Team operations
	GetTeamContext(ctx context.Context, excludeUser string) ([]models.TeamContext, error)

	// Project operations
	GetProjects(ctx context.Context) ([]models.Project, error)

	// UserPrompt operations
	CreateUserPrompt(ctx context.Context, prompt *models.UserPrompt) error
	GetUserPrompts(ctx context.Context, sessionID string, limit int) ([]models.UserPrompt, error)
	UpdateLatestPromptResponse(ctx context.Context, sessionID string, response string) error

	// FTS5 Search operations
	SearchFTS(ctx context.Context, query string, types []string, limit int) ([]models.SearchResult, error)
}

// Config holds database configuration.
type Config struct {
	Type     string // sqlite or mysql
	Path     string // for sqlite
	Host     string // for mysql
	Port     string // for mysql
	User     string // for mysql
	Password string // for mysql
	Database string // for mysql
}

// New creates a new database adapter based on configuration.
func New(cfg Config) (Adapter, error) {
	switch cfg.Type {
	case "mysql":
		return NewMySQL(cfg)
	default:
		return NewSQLite(cfg)
	}
}
