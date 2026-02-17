package memory

import (
	"context"
	"encoding/json"

	"github.com/do-focus/worker/internal/db"
	"github.com/do-focus/worker/pkg/models"
)

// ObservationType defines the type of observation.
type ObservationType string

const (
	ObservationDecision ObservationType = "decision"
	ObservationPattern  ObservationType = "pattern"
	ObservationLearning ObservationType = "learning"
	ObservationInsight  ObservationType = "insight"
)

// ObservationManager handles observation operations.
type ObservationManager struct {
	db    db.Adapter
	store *Store
}

// NewObservationManager creates a new observation manager.
func NewObservationManager(adapter db.Adapter, store *Store) *ObservationManager {
	return &ObservationManager{
		db:    adapter,
		store: store,
	}
}

// Record creates a new observation.
func (m *ObservationManager) Record(ctx context.Context, sessionID string, obsType ObservationType, content string, opts ...ObservationOption) (*models.Observation, error) {
	obs := &models.Observation{
		SessionID:  sessionID,
		Type:       string(obsType),
		Content:    content,
		Importance: 3, // Default importance
	}

	// Apply options
	for _, opt := range opts {
		opt(obs)
	}

	// Use batch queue for lower importance observations
	if obs.Importance < 4 {
		m.store.QueueObservation(*obs)
		return obs, nil
	}

	// Write high importance observations immediately
	if err := m.db.CreateObservation(ctx, obs); err != nil {
		return nil, err
	}

	return obs, nil
}

// ObservationOption configures an observation.
type ObservationOption func(*models.Observation)

// WithAgent sets the agent name for an observation.
func WithAgent(name string) ObservationOption {
	return func(o *models.Observation) {
		o.AgentName = name
	}
}

// WithImportance sets the importance level (1-5).
func WithImportance(level int) ObservationOption {
	return func(o *models.Observation) {
		if level < 1 {
			level = 1
		} else if level > 5 {
			level = 5
		}
		o.Importance = level
	}
}

// WithTags sets the tags for an observation.
func WithTags(tags ...string) ObservationOption {
	return func(o *models.Observation) {
		if len(tags) > 0 {
			data, _ := json.Marshal(tags)
			o.Tags = string(data)
		}
	}
}

// GetBySession retrieves all observations for a session.
func (m *ObservationManager) GetBySession(ctx context.Context, sessionID string) ([]models.Observation, error) {
	return m.db.GetObservations(ctx, sessionID)
}

// GetRecent retrieves recent observations for a user.
func (m *ObservationManager) GetRecent(ctx context.Context, userName string, limit int) ([]models.Observation, error) {
	if limit <= 0 {
		limit = 20
	}
	return m.db.GetRecentObservations(ctx, userName, limit)
}

// GetHighPriority retrieves high importance observations.
func (m *ObservationManager) GetHighPriority(ctx context.Context, userName string, limit int) ([]models.Observation, error) {
	// Get all recent and filter by importance >= 4
	all, err := m.db.GetRecentObservations(ctx, userName, limit*2)
	if err != nil {
		return nil, err
	}

	result := make([]models.Observation, 0, limit)
	for _, obs := range all {
		if obs.Importance >= 4 {
			result = append(result, obs)
			if len(result) >= limit {
				break
			}
		}
	}

	return result, nil
}
