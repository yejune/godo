package memory

import (
	"context"
	"time"

	"github.com/do-focus/worker/internal/db"
	"github.com/do-focus/worker/pkg/models"
)

// SessionManager handles session lifecycle.
type SessionManager struct {
	db    db.Adapter
	store *Store
}

// NewSessionManager creates a new session manager.
func NewSessionManager(adapter db.Adapter, store *Store) *SessionManager {
	return &SessionManager{
		db:    adapter,
		store: store,
	}
}

// StartSession creates a new session.
func (m *SessionManager) StartSession(ctx context.Context, id, userName string) (*models.Session, error) {
	session := &models.Session{
		ID:        id,
		UserName:  userName,
		StartedAt: time.Now(),
	}

	if err := m.db.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	// Cache the session
	m.store.SetCache(CacheKeySession, id, session)

	return session, nil
}

// GetSession retrieves a session by ID.
func (m *SessionManager) GetSession(ctx context.Context, id string) (*models.Session, error) {
	// Check cache first
	if cached, ok := m.store.GetCache(CacheKeySession, id); ok {
		return cached.(*models.Session), nil
	}

	session, err := m.db.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}

	if session != nil {
		m.store.SetCache(CacheKeySession, id, session)
	}

	return session, nil
}

// GetLatestSession retrieves the latest session for a user.
func (m *SessionManager) GetLatestSession(ctx context.Context, userName string) (*models.Session, error) {
	return m.db.GetLatestSession(ctx, userName)
}

// EndSession ends a session with an optional summary.
func (m *SessionManager) EndSession(ctx context.Context, id, summary string) error {
	if err := m.db.EndSession(ctx, id, summary); err != nil {
		return err
	}

	// Invalidate cache
	m.store.DeleteCache(CacheKeySession, id)

	return nil
}

// GetActiveSession returns the current active session for a user.
func (m *SessionManager) GetActiveSession(ctx context.Context, userName string) (*models.Session, error) {
	session, err := m.db.GetLatestSession(ctx, userName)
	if err != nil {
		return nil, err
	}

	// Only return if session is still active
	if session != nil && session.EndedAt == nil {
		return session, nil
	}

	return nil, nil
}
