package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/do-focus/worker/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

// SQLite implements the Adapter interface for SQLite.
type SQLite struct {
	db *sql.DB
}

// NewSQLite creates a new SQLite adapter.
func NewSQLite(cfg Config) (*SQLite, error) {
	path := cfg.Path
	if path == "" {
		path = ".do/memory.db"
	}

	db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1) // SQLite only supports one writer
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	s := &SQLite{db: db}

	// Initialize schema
	if err := s.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return s, nil
}

// initSchema creates the database tables if they don't exist.
func (s *SQLite) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		user_name TEXT NOT NULL,
		project_id TEXT,
		started_at DATETIME NOT NULL,
		ended_at DATETIME,
		summary TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_sessions_user_name ON sessions(user_name);
	CREATE INDEX IF NOT EXISTS idx_sessions_started_at ON sessions(started_at);
	CREATE INDEX IF NOT EXISTS idx_sessions_project_id ON sessions(project_id);

	CREATE TABLE IF NOT EXISTS observations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT NOT NULL,
		project_id TEXT,
		agent_name TEXT,
		type TEXT NOT NULL,
		content TEXT NOT NULL,
		importance INTEGER DEFAULT 3,
		tags TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (session_id) REFERENCES sessions(id)
	);

	CREATE INDEX IF NOT EXISTS idx_observations_session_id ON observations(session_id);
	CREATE INDEX IF NOT EXISTS idx_observations_type ON observations(type);
	CREATE INDEX IF NOT EXISTS idx_observations_importance ON observations(importance);
	CREATE INDEX IF NOT EXISTS idx_observations_project_id ON observations(project_id);

	CREATE TABLE IF NOT EXISTS summaries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT,
		project_id TEXT,
		type TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (session_id) REFERENCES sessions(id)
	);

	CREATE INDEX IF NOT EXISTS idx_summaries_type ON summaries(type);
	CREATE INDEX IF NOT EXISTS idx_summaries_project_id ON summaries(project_id);

	CREATE TABLE IF NOT EXISTS plans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		session_id TEXT,
		project_id TEXT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'draft',
		file_path TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (session_id) REFERENCES sessions(id)
	);

	CREATE INDEX IF NOT EXISTS idx_plans_status ON plans(status);
	CREATE INDEX IF NOT EXISTS idx_plans_project_id ON plans(project_id);
	`

	_, err := s.db.Exec(schema)
	if err != nil {
		return err
	}

	// Run migrations for existing databases
	return s.runMigrations()
}

// runMigrations applies schema migrations for existing databases.
func (s *SQLite) runMigrations() error {
	// Migration 004: Add project_id column to existing tables
	migrations := []string{
		// Check and add project_id to sessions
		`ALTER TABLE sessions ADD COLUMN project_id TEXT`,
		// Check and add project_id to observations
		`ALTER TABLE observations ADD COLUMN project_id TEXT`,
		// Check and add project_id to summaries
		`ALTER TABLE summaries ADD COLUMN project_id TEXT`,
		// Check and add project_id to plans
		`ALTER TABLE plans ADD COLUMN project_id TEXT`,
	}

	for _, migration := range migrations {
		// SQLite will error if column already exists, which is fine
		_, _ = s.db.Exec(migration)
	}

	// Migrate existing data: copy user_name to project_id where project_id is NULL
	_, _ = s.db.Exec(`UPDATE sessions SET project_id = user_name WHERE project_id IS NULL OR project_id = ''`)

	// Create indexes if they don't exist (already in schema, but ensure for migrated DBs)
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_sessions_project_id ON sessions(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_observations_project_id ON observations(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_summaries_project_id ON summaries(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_plans_project_id ON plans(project_id)`,
	}

	for _, idx := range indexes {
		if _, err := s.db.Exec(idx); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	// Migration 005: Extend observations table with rich metadata columns
	observationColumns := []string{
		`ALTER TABLE observations ADD COLUMN title TEXT`,
		`ALTER TABLE observations ADD COLUMN subtitle TEXT`,
		`ALTER TABLE observations ADD COLUMN narrative TEXT`,
		`ALTER TABLE observations ADD COLUMN facts TEXT`,
		`ALTER TABLE observations ADD COLUMN concepts TEXT`,
		`ALTER TABLE observations ADD COLUMN files_read TEXT`,
		`ALTER TABLE observations ADD COLUMN files_modified TEXT`,
		`ALTER TABLE observations ADD COLUMN result_preview TEXT`,
		`ALTER TABLE observations ADD COLUMN prompt_number INTEGER`,
		`ALTER TABLE observations ADD COLUMN discovery_tokens INTEGER DEFAULT 0`,
	}
	for _, col := range observationColumns {
		_, _ = s.db.Exec(col)
	}

	// Migration 006: Create user_prompts table
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS user_prompts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id TEXT NOT NULL,
			prompt_number INTEGER NOT NULL,
			prompt_text TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_at_epoch INTEGER NOT NULL,
			FOREIGN KEY (session_id) REFERENCES sessions(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create user_prompts table: %w", err)
	}
	_, _ = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_user_prompts_session ON user_prompts(session_id)`)
	_, _ = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_user_prompts_prompt_number ON user_prompts(session_id, prompt_number)`)

	// Migration 007: Extend summaries table with structured fields
	summaryColumns := []string{
		`ALTER TABLE summaries ADD COLUMN request TEXT`,
		`ALTER TABLE summaries ADD COLUMN investigated TEXT`,
		`ALTER TABLE summaries ADD COLUMN learned TEXT`,
		`ALTER TABLE summaries ADD COLUMN completed TEXT`,
		`ALTER TABLE summaries ADD COLUMN next_steps TEXT`,
		`ALTER TABLE summaries ADD COLUMN files_read TEXT`,
		`ALTER TABLE summaries ADD COLUMN files_edited TEXT`,
		`ALTER TABLE summaries ADD COLUMN discovery_tokens INTEGER DEFAULT 0`,
	}
	for _, col := range summaryColumns {
		_, _ = s.db.Exec(col)
	}

	// Migration 008: Add request_prompt column to plans table
	_, _ = s.db.Exec(`ALTER TABLE plans ADD COLUMN request_prompt TEXT`)

	// Migration 009: Add source_message column to summaries table
	_, _ = s.db.Exec(`ALTER TABLE summaries ADD COLUMN source_message TEXT`)

	// Migration 010: Add full_transcript column to summaries table (stores complete session transcript)
	_, _ = s.db.Exec(`ALTER TABLE summaries ADD COLUMN full_transcript TEXT`)

	// Migration 011: Add response column to user_prompts table (stores assistant response with tool_use)
	_, _ = s.db.Exec(`ALTER TABLE user_prompts ADD COLUMN response TEXT`)

	return nil
}

// Health checks database connectivity.
func (s *SQLite) Health(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// Close closes the database connection.
func (s *SQLite) Close() error {
	return s.db.Close()
}

// CreateSession creates a new session.
func (s *SQLite) CreateSession(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (id, user_name, project_id, started_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	_, err := s.db.ExecContext(ctx, query, session.ID, session.UserName, session.ProjectID, session.StartedAt, now, now)
	return err
}

// GetSession retrieves a session by ID.
func (s *SQLite) GetSession(ctx context.Context, id string) (*models.Session, error) {
	query := `SELECT id, user_name, COALESCE(project_id, ''), started_at, ended_at, COALESCE(summary, ''), created_at, updated_at FROM sessions WHERE id = ?`
	session := &models.Session{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID, &session.UserName, &session.ProjectID, &session.StartedAt, &session.EndedAt,
		&session.Summary, &session.CreatedAt, &session.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

// GetLatestSession retrieves the latest session for a user.
func (s *SQLite) GetLatestSession(ctx context.Context, userName string) (*models.Session, error) {
	query := `
		SELECT id, user_name, started_at, ended_at, COALESCE(summary, ''), created_at, updated_at
		FROM sessions
		WHERE user_name = ?
		ORDER BY started_at DESC
		LIMIT 1
	`
	session := &models.Session{}
	err := s.db.QueryRowContext(ctx, query, userName).Scan(
		&session.ID, &session.UserName, &session.StartedAt, &session.EndedAt,
		&session.Summary, &session.CreatedAt, &session.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

// EndSession ends a session with an optional summary.
func (s *SQLite) EndSession(ctx context.Context, id string, summary string) error {
	query := `UPDATE sessions SET ended_at = ?, summary = ?, updated_at = ? WHERE id = ?`
	now := time.Now()
	_, err := s.db.ExecContext(ctx, query, now, summary, now, id)
	return err
}

// CreateObservation creates a new observation.
func (s *SQLite) CreateObservation(ctx context.Context, obs *models.Observation) error {
	query := `
		INSERT INTO observations (session_id, agent_name, type, content, importance, tags, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := s.db.ExecContext(ctx, query, obs.SessionID, obs.AgentName, obs.Type, obs.Content, obs.Importance, obs.Tags, time.Now())
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		obs.ID = id
	}
	return nil
}

// GetObservations retrieves observations for a session.
func (s *SQLite) GetObservations(ctx context.Context, sessionID string) ([]models.Observation, error) {
	query := `
		SELECT id, session_id, COALESCE(agent_name, ''), type, content, importance, COALESCE(tags, ''), created_at
		FROM observations
		WHERE session_id = ?
		ORDER BY created_at DESC
	`
	rows, err := s.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var observations []models.Observation
	for rows.Next() {
		var obs models.Observation
		if err := rows.Scan(&obs.ID, &obs.SessionID, &obs.AgentName, &obs.Type, &obs.Content, &obs.Importance, &obs.Tags, &obs.CreatedAt); err != nil {
			return nil, err
		}
		observations = append(observations, obs)
	}
	return observations, rows.Err()
}

// GetRecentObservations retrieves recent observations across sessions for a user.
func (s *SQLite) GetRecentObservations(ctx context.Context, userName string, limit int) ([]models.Observation, error) {
	query := `
		SELECT o.id, o.session_id, COALESCE(o.agent_name, ''), o.type, o.content, o.importance, COALESCE(o.tags, ''), o.created_at
		FROM observations o
		JOIN sessions s ON o.session_id = s.id
		WHERE s.user_name = ?
		ORDER BY o.importance DESC, o.created_at DESC
		LIMIT ?
	`
	rows, err := s.db.QueryContext(ctx, query, userName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var observations []models.Observation
	for rows.Next() {
		var obs models.Observation
		if err := rows.Scan(&obs.ID, &obs.SessionID, &obs.AgentName, &obs.Type, &obs.Content, &obs.Importance, &obs.Tags, &obs.CreatedAt); err != nil {
			return nil, err
		}
		observations = append(observations, obs)
	}
	return observations, rows.Err()
}

// GetObservationsFiltered retrieves observations with optional filters and pagination.
func (s *SQLite) GetObservationsFiltered(ctx context.Context, sessionID string, obsType string, limit int, offset int) ([]models.Observation, error) {
	query := `
		SELECT id, session_id, COALESCE(agent_name, ''), type, content, importance, COALESCE(tags, ''), created_at
		FROM observations
		WHERE (? = '' OR session_id = ?) AND (? = '' OR type = ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := s.db.QueryContext(ctx, query, sessionID, sessionID, obsType, obsType, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var observations []models.Observation
	for rows.Next() {
		var obs models.Observation
		if err := rows.Scan(&obs.ID, &obs.SessionID, &obs.AgentName, &obs.Type, &obs.Content, &obs.Importance, &obs.Tags, &obs.CreatedAt); err != nil {
			return nil, err
		}
		observations = append(observations, obs)
	}
	return observations, rows.Err()
}

// SearchObservations searches observations by content.
func (s *SQLite) SearchObservations(ctx context.Context, query string, limit int) ([]models.Observation, error) {
	sqlQuery := `
		SELECT id, session_id, COALESCE(agent_name, ''), type, content, importance, COALESCE(tags, ''), created_at
		FROM observations
		WHERE content LIKE ?
		ORDER BY importance DESC, created_at DESC
		LIMIT ?
	`
	if limit <= 0 {
		limit = 50
	}
	searchPattern := "%" + query + "%"
	rows, err := s.db.QueryContext(ctx, sqlQuery, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var observations []models.Observation
	for rows.Next() {
		var obs models.Observation
		if err := rows.Scan(&obs.ID, &obs.SessionID, &obs.AgentName, &obs.Type, &obs.Content, &obs.Importance, &obs.Tags, &obs.CreatedAt); err != nil {
			return nil, err
		}
		observations = append(observations, obs)
	}
	return observations, rows.Err()
}

// CreateSummary creates a new summary with structured fields.
func (s *SQLite) CreateSummary(ctx context.Context, summary *models.Summary) error {
	query := `INSERT INTO summaries (
		session_id, type, content, created_at,
		request, investigated, learned, completed, next_steps,
		files_read, files_edited, discovery_tokens, source_message, full_transcript
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := s.db.ExecContext(ctx, query,
		summary.SessionID,
		summary.Type,
		summary.Content,
		time.Now(),
		summary.Request,
		summary.Investigated,
		summary.Learned,
		summary.Completed,
		summary.NextSteps,
		summary.FilesRead,
		summary.FilesEdited,
		summary.DiscoveryTokens,
		summary.SourceMessage,
		summary.FullTranscript,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		summary.ID = id
	}
	return nil
}

// GetSummaries retrieves summaries by type.
func (s *SQLite) GetSummaries(ctx context.Context, summaryType string, limit int) ([]models.Summary, error) {
	query := `
		SELECT id, session_id, type, content, created_at
		FROM summaries
		WHERE type = ?
		ORDER BY created_at DESC
		LIMIT ?
	`
	rows, err := s.db.QueryContext(ctx, query, summaryType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []models.Summary
	for rows.Next() {
		var sum models.Summary
		if err := rows.Scan(&sum.ID, &sum.SessionID, &sum.Type, &sum.Content, &sum.CreatedAt); err != nil {
			return nil, err
		}
		summaries = append(summaries, sum)
	}
	return summaries, rows.Err()
}

// GetAllSummaries retrieves all summaries within a date range.
func (s *SQLite) GetAllSummaries(ctx context.Context, days int, limit int) ([]models.Summary, error) {
	if days <= 0 {
		days = 7
	}
	if limit <= 0 {
		limit = 100
	}
	query := `
		SELECT id, COALESCE(session_id, ''), type, content, created_at,
			COALESCE(request, ''), COALESCE(investigated, ''), COALESCE(learned, ''),
			COALESCE(completed, ''), COALESCE(next_steps, ''), COALESCE(source_message, ''),
			COALESCE(full_transcript, '')
		FROM summaries
		WHERE created_at >= datetime('now', ? || ' days')
		ORDER BY created_at DESC
		LIMIT ?
	`
	daysArg := fmt.Sprintf("-%d", days)
	rows, err := s.db.QueryContext(ctx, query, daysArg, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []models.Summary
	for rows.Next() {
		var sum models.Summary
		var request, investigated, learned, completed, nextSteps, sourceMessage, fullTranscript string
		if err := rows.Scan(&sum.ID, &sum.SessionID, &sum.Type, &sum.Content, &sum.CreatedAt,
			&request, &investigated, &learned, &completed, &nextSteps, &sourceMessage, &fullTranscript); err != nil {
			return nil, err
		}
		if request != "" {
			sum.Request = &request
		}
		if investigated != "" {
			sum.Investigated = &investigated
		}
		if learned != "" {
			sum.Learned = &learned
		}
		if completed != "" {
			sum.Completed = &completed
		}
		if nextSteps != "" {
			sum.NextSteps = &nextSteps
		}
		if sourceMessage != "" {
			sum.SourceMessage = sourceMessage
		}
		if fullTranscript != "" {
			sum.FullTranscript = fullTranscript
		}
		summaries = append(summaries, sum)
	}
	return summaries, rows.Err()
}

// GetLatestSummary retrieves the latest session summary for a user.
func (s *SQLite) GetLatestSummary(ctx context.Context, userName string) (*models.Summary, error) {
	query := `
		SELECT su.id, COALESCE(su.session_id, ''), su.type, su.content, su.created_at
		FROM summaries su
		JOIN sessions s ON su.session_id = s.id
		WHERE s.user_name = ? AND su.type = 'session'
		ORDER BY su.created_at DESC
		LIMIT 1
	`
	summary := &models.Summary{}
	err := s.db.QueryRowContext(ctx, query, userName).Scan(
		&summary.ID, &summary.SessionID, &summary.Type, &summary.Content, &summary.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return summary, nil
}

// CreatePlan creates a new plan.
func (s *SQLite) CreatePlan(ctx context.Context, plan *models.Plan) error {
	query := `
		INSERT INTO plans (session_id, title, content, status, file_path, request_prompt, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := s.db.ExecContext(ctx, query, plan.SessionID, plan.Title, plan.Content, "draft", plan.FilePath, plan.RequestPrompt, now, now)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		plan.ID = id
	}
	return nil
}

// GetActivePlan retrieves the active plan for a user.
func (s *SQLite) GetActivePlan(ctx context.Context, userName string) (*models.Plan, error) {
	query := `
		SELECT p.id, p.session_id, p.title, p.content, p.status, COALESCE(p.file_path, ''), COALESCE(p.request_prompt, ''), p.created_at, p.updated_at
		FROM plans p
		JOIN sessions s ON p.session_id = s.id
		WHERE s.user_name = ? AND p.status = 'active'
		ORDER BY p.updated_at DESC
		LIMIT 1
	`
	plan := &models.Plan{}
	err := s.db.QueryRowContext(ctx, query, userName).Scan(
		&plan.ID, &plan.SessionID, &plan.Title, &plan.Content,
		&plan.Status, &plan.FilePath, &plan.RequestPrompt, &plan.CreatedAt, &plan.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return plan, err
}

// GetAllPlans retrieves all plans with optional session filter.
func (s *SQLite) GetAllPlans(ctx context.Context, sessionID string, limit int) ([]models.Plan, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
		SELECT id, COALESCE(session_id, ''), title, content, status, COALESCE(file_path, ''), COALESCE(request_prompt, ''), created_at, updated_at
		FROM plans
		WHERE ? = '' OR session_id = ?
		ORDER BY updated_at DESC
		LIMIT ?
	`
	rows, err := s.db.QueryContext(ctx, query, sessionID, sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.Plan
	for rows.Next() {
		var plan models.Plan
		if err := rows.Scan(&plan.ID, &plan.SessionID, &plan.Title, &plan.Content, &plan.Status, &plan.FilePath, &plan.RequestPrompt, &plan.CreatedAt, &plan.UpdatedAt); err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, rows.Err()
}

// UpdatePlanStatus updates a plan's status.
func (s *SQLite) UpdatePlanStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE plans SET status = ?, updated_at = ? WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

// GetRecentSessions retrieves recent sessions.
func (s *SQLite) GetRecentSessions(ctx context.Context, limit int) ([]models.Session, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `
		SELECT id, user_name, started_at, ended_at, COALESCE(summary, ''), created_at, updated_at
		FROM sessions
		ORDER BY started_at DESC
		LIMIT ?
	`
	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		if err := rows.Scan(&session.ID, &session.UserName, &session.StartedAt, &session.EndedAt, &session.Summary, &session.CreatedAt, &session.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, rows.Err()
}

// GetTeamContext retrieves context from other team members.
func (s *SQLite) GetTeamContext(ctx context.Context, excludeUser string) ([]models.TeamContext, error) {
	query := `
		SELECT
			s.user_name,
			MAX(s.started_at) as last_activity,
			COALESCE(s.summary, '') as summary,
			COALESCE(p.title, '') as active_plan
		FROM sessions s
		LEFT JOIN plans p ON p.session_id = s.id AND p.status = 'active'
		WHERE s.user_name != ? AND s.ended_at IS NOT NULL
		GROUP BY s.user_name
		ORDER BY last_activity DESC
		LIMIT 10
	`
	rows, err := s.db.QueryContext(ctx, query, excludeUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contexts []models.TeamContext
	for rows.Next() {
		var tc models.TeamContext
		if err := rows.Scan(&tc.UserName, &tc.LastActivity, &tc.Summary, &tc.ActivePlan); err != nil {
			return nil, err
		}
		contexts = append(contexts, tc)
	}
	return contexts, rows.Err()
}

// GetProjects retrieves all registered projects with session statistics.
func (s *SQLite) GetProjects(ctx context.Context) ([]models.Project, error) {
	query := `
		SELECT
			project_id,
			project_id as path,
			COUNT(*) as session_count,
			MAX(started_at) as last_activity
		FROM sessions
		WHERE project_id IS NOT NULL AND project_id != ''
		GROUP BY project_id
		ORDER BY last_activity DESC
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		var lastActivityStr string
		if err := rows.Scan(&p.ID, &p.Path, &p.SessionCount, &lastActivityStr); err != nil {
			return nil, err
		}
		// Parse SQLite datetime string (format: 2006-01-02 15:04:05.999999999+09:00)
		if t, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", lastActivityStr); err == nil {
			p.LastActivity = t
		} else if t, err := time.Parse("2006-01-02 15:04:05-07:00", lastActivityStr); err == nil {
			p.LastActivity = t
		} else if t, err := time.Parse("2006-01-02 15:04:05", lastActivityStr); err == nil {
			p.LastActivity = t
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// tagsToJSON converts a slice of strings to a JSON string.
func tagsToJSON(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	data, _ := json.Marshal(tags)
	return string(data)
}

// CreateUserPrompt creates a new user prompt record.
func (s *SQLite) CreateUserPrompt(ctx context.Context, prompt *models.UserPrompt) error {
	query := `
		INSERT INTO user_prompts (session_id, prompt_number, prompt_text, created_at, created_at_epoch)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := s.db.ExecContext(ctx, query, prompt.SessionID, prompt.PromptNumber, prompt.PromptText, now, now.Unix())
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		prompt.ID = id
		prompt.CreatedAt = now
		prompt.CreatedAtEpoch = now.Unix()
	}
	return nil
}

// GetUserPrompts retrieves user prompts for a session (or all if sessionID is empty).
func (s *SQLite) GetUserPrompts(ctx context.Context, sessionID string, limit int) ([]models.UserPrompt, error) {
	if limit <= 0 {
		limit = 100
	}

	var query string
	var rows *sql.Rows
	var err error

	if sessionID == "" {
		// Return all prompts (most recent first)
		query = `
			SELECT id, session_id, prompt_number, prompt_text, COALESCE(response, ''), created_at, created_at_epoch
			FROM user_prompts
			ORDER BY created_at DESC
			LIMIT ?
		`
		rows, err = s.db.QueryContext(ctx, query, limit)
	} else {
		// Return prompts for specific session
		query = `
			SELECT id, session_id, prompt_number, prompt_text, COALESCE(response, ''), created_at, created_at_epoch
			FROM user_prompts
			WHERE session_id = ?
			ORDER BY prompt_number ASC
			LIMIT ?
		`
		rows, err = s.db.QueryContext(ctx, query, sessionID, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prompts []models.UserPrompt
	for rows.Next() {
		var p models.UserPrompt
		if err := rows.Scan(&p.ID, &p.SessionID, &p.PromptNumber, &p.PromptText, &p.Response, &p.CreatedAt, &p.CreatedAtEpoch); err != nil {
			return nil, err
		}
		prompts = append(prompts, p)
	}
	return prompts, rows.Err()
}

// UpdateLatestPromptResponse updates the response for the latest prompt in a session.
func (s *SQLite) UpdateLatestPromptResponse(ctx context.Context, sessionID string, response string) error {
	query := `
		UPDATE user_prompts
		SET response = ?
		WHERE id = (
			SELECT id FROM user_prompts
			WHERE session_id = ?
			ORDER BY created_at DESC
			LIMIT 1
		)
	`
	_, err := s.db.ExecContext(ctx, query, response, sessionID)
	return err
}

// SearchFTS performs full-text search across observations and user_prompts using FTS5.
func (s *SQLite) SearchFTS(ctx context.Context, query string, types []string, limit int) ([]models.SearchResult, error) {
	// FTS5 disabled - Go sqlite3 driver doesn't support FTS5 by default
	// Use simple LIKE search instead
	if limit <= 0 {
		limit = 50
	}

	var results []models.SearchResult
	likeQuery := "%" + query + "%"

	// Search observations using LIKE
	obsQuery := `
		SELECT id, session_id, content, created_at
		FROM observations
		WHERE content LIKE ?
		ORDER BY created_at DESC
		LIMIT ?
	`
	rows, err := s.db.QueryContext(ctx, obsQuery, likeQuery, limit)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var r models.SearchResult
			r.Type = "observation"
			if err := rows.Scan(&r.ID, &r.SessionID, &r.Content, &r.CreatedAt); err != nil {
				continue
			}
			r.Snippet = r.Content
			if len(r.Snippet) > 100 {
				r.Snippet = r.Snippet[:100] + "..."
			}
			results = append(results, r)
		}
	}

	return results, nil
}
