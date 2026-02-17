package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/do-focus/worker/pkg/models"
	_ "github.com/go-sql-driver/mysql"
)

// MySQL implements the Adapter interface for MySQL.
type MySQL struct {
	db *sql.DB
}

// NewMySQL creates a new MySQL adapter.
func NewMySQL(cfg Config) (*MySQL, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	m := &MySQL{db: db}

	// Initialize schema
	if err := m.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return m, nil
}

// initSchema creates the database tables if they don't exist.
func (m *MySQL) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS sessions (
		id VARCHAR(255) PRIMARY KEY,
		user_name VARCHAR(255) NOT NULL,
		project_id VARCHAR(500),
		started_at DATETIME NOT NULL,
		ended_at DATETIME,
		summary TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_sessions_user_name (user_name),
		INDEX idx_sessions_started_at (started_at),
		INDEX idx_sessions_project_id (project_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	if _, err := m.db.Exec(schema); err != nil {
		return err
	}

	// Run migrations for existing tables
	if err := m.runMigrations(); err != nil {
		return err
	}

	schema = `
	CREATE TABLE IF NOT EXISTS observations (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		session_id VARCHAR(255) NOT NULL,
		agent_name VARCHAR(255),
		type VARCHAR(50) NOT NULL,
		content TEXT NOT NULL,
		importance INT DEFAULT 3,
		tags JSON,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_observations_session_id (session_id),
		INDEX idx_observations_type (type),
		INDEX idx_observations_importance (importance),
		FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	if _, err := m.db.Exec(schema); err != nil {
		return err
	}

	schema = `
	CREATE TABLE IF NOT EXISTS summaries (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		session_id VARCHAR(255),
		type VARCHAR(50) NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_summaries_type (type),
		FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE SET NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	if _, err := m.db.Exec(schema); err != nil {
		return err
	}

	schema = `
	CREATE TABLE IF NOT EXISTS plans (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		session_id VARCHAR(255),
		title VARCHAR(500) NOT NULL,
		content TEXT NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'draft',
		file_path VARCHAR(1000),
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_plans_status (status),
		FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE SET NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`
	_, err := m.db.Exec(schema)
	return err
}

// Health checks database connectivity.
func (m *MySQL) Health(ctx context.Context) error {
	return m.db.PingContext(ctx)
}

// Close closes the database connection.
func (m *MySQL) Close() error {
	return m.db.Close()
}

// CreateSession creates a new session.
func (m *MySQL) CreateSession(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (id, user_name, project_id, started_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`
	_, err := m.db.ExecContext(ctx, query, session.ID, session.UserName, session.ProjectID, session.StartedAt)
	return err
}

// GetSession retrieves a session by ID.
func (m *MySQL) GetSession(ctx context.Context, id string) (*models.Session, error) {
	query := `SELECT id, user_name, COALESCE(project_id, ''), started_at, ended_at, COALESCE(summary, ''), created_at, updated_at FROM sessions WHERE id = ?`
	session := &models.Session{}
	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID, &session.UserName, &session.ProjectID, &session.StartedAt, &session.EndedAt,
		&session.Summary, &session.CreatedAt, &session.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

// GetLatestSession retrieves the latest session for a user.
func (m *MySQL) GetLatestSession(ctx context.Context, userName string) (*models.Session, error) {
	query := `
		SELECT id, user_name, started_at, ended_at, summary, created_at, updated_at
		FROM sessions
		WHERE user_name = ?
		ORDER BY started_at DESC
		LIMIT 1
	`
	session := &models.Session{}
	err := m.db.QueryRowContext(ctx, query, userName).Scan(
		&session.ID, &session.UserName, &session.StartedAt, &session.EndedAt,
		&session.Summary, &session.CreatedAt, &session.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

// EndSession ends a session with an optional summary.
func (m *MySQL) EndSession(ctx context.Context, id string, summary string) error {
	query := `UPDATE sessions SET ended_at = NOW(), summary = ? WHERE id = ?`
	_, err := m.db.ExecContext(ctx, query, summary, id)
	return err
}

// CreateObservation creates a new observation.
func (m *MySQL) CreateObservation(ctx context.Context, obs *models.Observation) error {
	query := `
		INSERT INTO observations (session_id, agent_name, type, content, importance, tags, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
	`
	result, err := m.db.ExecContext(ctx, query, obs.SessionID, obs.AgentName, obs.Type, obs.Content, obs.Importance, obs.Tags)
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
func (m *MySQL) GetObservations(ctx context.Context, sessionID string) ([]models.Observation, error) {
	query := `
		SELECT id, session_id, agent_name, type, content, importance, tags, created_at
		FROM observations
		WHERE session_id = ?
		ORDER BY created_at DESC
	`
	rows, err := m.db.QueryContext(ctx, query, sessionID)
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
func (m *MySQL) GetRecentObservations(ctx context.Context, userName string, limit int) ([]models.Observation, error) {
	query := `
		SELECT o.id, o.session_id, o.agent_name, o.type, o.content, o.importance, o.tags, o.created_at
		FROM observations o
		JOIN sessions s ON o.session_id = s.id
		WHERE s.user_name = ?
		ORDER BY o.importance DESC, o.created_at DESC
		LIMIT ?
	`
	rows, err := m.db.QueryContext(ctx, query, userName, limit)
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
func (m *MySQL) GetObservationsFiltered(ctx context.Context, sessionID string, obsType string, limit int, offset int) ([]models.Observation, error) {
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
	rows, err := m.db.QueryContext(ctx, query, sessionID, sessionID, obsType, obsType, limit, offset)
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
func (m *MySQL) SearchObservations(ctx context.Context, query string, limit int) ([]models.Observation, error) {
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
	rows, err := m.db.QueryContext(ctx, sqlQuery, searchPattern, limit)
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

// CreateSummary creates a new summary.
func (m *MySQL) CreateSummary(ctx context.Context, summary *models.Summary) error {
	query := `INSERT INTO summaries (
		session_id, type, content, created_at,
		request, investigated, learned, completed, next_steps,
		files_read, files_edited, discovery_tokens, source_message, full_transcript
	) VALUES (?, ?, ?, NOW(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.db.ExecContext(ctx, query,
		summary.SessionID,
		summary.Type,
		summary.Content,
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
func (m *MySQL) GetSummaries(ctx context.Context, summaryType string, limit int) ([]models.Summary, error) {
	query := `
		SELECT id, session_id, type, content, created_at
		FROM summaries
		WHERE type = ?
		ORDER BY created_at DESC
		LIMIT ?
	`
	rows, err := m.db.QueryContext(ctx, query, summaryType, limit)
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
func (m *MySQL) GetAllSummaries(ctx context.Context, days int, limit int) ([]models.Summary, error) {
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
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		ORDER BY created_at DESC
		LIMIT ?
	`
	rows, err := m.db.QueryContext(ctx, query, days, limit)
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
func (m *MySQL) GetLatestSummary(ctx context.Context, userName string) (*models.Summary, error) {
	query := `
		SELECT su.id, COALESCE(su.session_id, ''), su.type, su.content, su.created_at
		FROM summaries su
		JOIN sessions s ON su.session_id = s.id
		WHERE s.user_name = ? AND su.type = 'session'
		ORDER BY su.created_at DESC
		LIMIT 1
	`
	summary := &models.Summary{}
	err := m.db.QueryRowContext(ctx, query, userName).Scan(
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
func (m *MySQL) CreatePlan(ctx context.Context, plan *models.Plan) error {
	query := `
		INSERT INTO plans (session_id, title, content, status, file_path, request_prompt, created_at, updated_at)
		VALUES (?, ?, ?, 'draft', ?, ?, NOW(), NOW())
	`
	result, err := m.db.ExecContext(ctx, query, plan.SessionID, plan.Title, plan.Content, plan.FilePath, plan.RequestPrompt)
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
func (m *MySQL) GetActivePlan(ctx context.Context, userName string) (*models.Plan, error) {
	query := `
		SELECT p.id, p.session_id, p.title, p.content, p.status, COALESCE(p.file_path, ''), COALESCE(p.request_prompt, ''), p.created_at, p.updated_at
		FROM plans p
		JOIN sessions s ON p.session_id = s.id
		WHERE s.user_name = ? AND p.status = 'active'
		ORDER BY p.updated_at DESC
		LIMIT 1
	`
	plan := &models.Plan{}
	err := m.db.QueryRowContext(ctx, query, userName).Scan(
		&plan.ID, &plan.SessionID, &plan.Title, &plan.Content,
		&plan.Status, &plan.FilePath, &plan.RequestPrompt, &plan.CreatedAt, &plan.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return plan, err
}

// GetAllPlans retrieves all plans with optional session filter.
func (m *MySQL) GetAllPlans(ctx context.Context, sessionID string, limit int) ([]models.Plan, error) {
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
	rows, err := m.db.QueryContext(ctx, query, sessionID, sessionID, limit)
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
func (m *MySQL) UpdatePlanStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE plans SET status = ? WHERE id = ?`
	_, err := m.db.ExecContext(ctx, query, status, id)
	return err
}

// GetRecentSessions retrieves recent sessions.
func (m *MySQL) GetRecentSessions(ctx context.Context, limit int) ([]models.Session, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `
		SELECT id, user_name, started_at, ended_at, COALESCE(summary, ''), created_at, updated_at
		FROM sessions
		ORDER BY started_at DESC
		LIMIT ?
	`
	rows, err := m.db.QueryContext(ctx, query, limit)
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
func (m *MySQL) GetTeamContext(ctx context.Context, excludeUser string) ([]models.TeamContext, error) {
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
	rows, err := m.db.QueryContext(ctx, query, excludeUser)
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

// runMigrations applies schema migrations for existing databases.
func (m *MySQL) runMigrations() error {
	// Check if project_id column exists and add if not
	var columnExists int
	err := m.db.QueryRow(`
		SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'sessions' AND COLUMN_NAME = 'project_id'
	`).Scan(&columnExists)
	if err != nil {
		return err
	}

	if columnExists == 0 {
		// Add project_id column to sessions
		if _, err := m.db.Exec(`ALTER TABLE sessions ADD COLUMN project_id VARCHAR(500) AFTER user_name`); err != nil {
			return fmt.Errorf("failed to add project_id to sessions: %w", err)
		}
		// Add index
		if _, err := m.db.Exec(`CREATE INDEX idx_sessions_project_id ON sessions(project_id)`); err != nil {
			// Index might already exist, ignore error
		}
		// Migrate existing data
		if _, err := m.db.Exec(`UPDATE sessions SET project_id = user_name WHERE project_id IS NULL OR project_id = ''`); err != nil {
			return fmt.Errorf("failed to migrate project_id data: %w", err)
		}
	}

	// Migration 005: Extend observations table with rich metadata columns
	observationColumns := []struct {
		name string
		def  string
	}{
		{"title", "TEXT"},
		{"subtitle", "TEXT"},
		{"narrative", "TEXT"},
		{"facts", "JSON"},
		{"concepts", "JSON"},
		{"files_read", "JSON"},
		{"files_modified", "JSON"},
		{"result_preview", "TEXT"},
		{"prompt_number", "INT"},
		{"discovery_tokens", "INT DEFAULT 0"},
	}
	for _, col := range observationColumns {
		var exists int
		m.db.QueryRow(`
			SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'observations' AND COLUMN_NAME = ?
		`, col.name).Scan(&exists)
		if exists == 0 {
			m.db.Exec(fmt.Sprintf(`ALTER TABLE observations ADD COLUMN %s %s`, col.name, col.def))
		}
	}

	// Migration 006: Create user_prompts table
	_, _ = m.db.Exec(`
		CREATE TABLE IF NOT EXISTS user_prompts (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			session_id VARCHAR(255) NOT NULL,
			prompt_number INT NOT NULL,
			prompt_text TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_at_epoch BIGINT NOT NULL,
			INDEX idx_user_prompts_session (session_id),
			INDEX idx_user_prompts_prompt_number (session_id, prompt_number),
			FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
	`)

	// Migration 007: Extend summaries table with structured fields
	summaryColumns := []struct {
		name string
		def  string
	}{
		{"request", "TEXT"},
		{"investigated", "TEXT"},
		{"learned", "TEXT"},
		{"completed", "TEXT"},
		{"next_steps", "TEXT"},
		{"files_read", "JSON"},
		{"files_edited", "JSON"},
		{"discovery_tokens", "INT DEFAULT 0"},
		{"source_message", "LONGTEXT"},
		{"full_transcript", "LONGTEXT"},
	}
	for _, col := range summaryColumns {
		var exists int
		m.db.QueryRow(`
			SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'summaries' AND COLUMN_NAME = ?
		`, col.name).Scan(&exists)
		if exists == 0 {
			m.db.Exec(fmt.Sprintf(`ALTER TABLE summaries ADD COLUMN %s %s`, col.name, col.def))
		}
	}

	// Migration 008: Add request_prompt column to plans table
	var requestPromptExists int
	m.db.QueryRow(`
		SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'plans' AND COLUMN_NAME = 'request_prompt'
	`).Scan(&requestPromptExists)
	if requestPromptExists == 0 {
		m.db.Exec(`ALTER TABLE plans ADD COLUMN request_prompt TEXT`)
	}

	// Migration 011: Add response column to user_prompts table
	var responseExists int
	m.db.QueryRow(`
		SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'user_prompts' AND COLUMN_NAME = 'response'
	`).Scan(&responseExists)
	if responseExists == 0 {
		m.db.Exec(`ALTER TABLE user_prompts ADD COLUMN response LONGTEXT`)
	}

	return nil
}

// GetProjects retrieves all registered projects with session statistics.
func (m *MySQL) GetProjects(ctx context.Context) ([]models.Project, error) {
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
	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Path, &p.SessionCount, &p.LastActivity); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// mysqlTagsToJSON converts a slice of strings to a JSON string.
func mysqlTagsToJSON(tags []string) string {
	if len(tags) == 0 {
		return "[]"
	}
	data, _ := json.Marshal(tags)
	return string(data)
}

// CreateUserPrompt creates a new user prompt record.
func (m *MySQL) CreateUserPrompt(ctx context.Context, prompt *models.UserPrompt) error {
	query := `
		INSERT INTO user_prompts (session_id, prompt_number, prompt_text, created_at, created_at_epoch)
		VALUES (?, ?, ?, NOW(), ?)
	`
	result, err := m.db.ExecContext(ctx, query, prompt.SessionID, prompt.PromptNumber, prompt.PromptText, time.Now().Unix())
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err == nil {
		prompt.ID = id
	}
	return nil
}

// GetUserPrompts retrieves user prompts for a session (or all if sessionID is empty).
func (m *MySQL) GetUserPrompts(ctx context.Context, sessionID string, limit int) ([]models.UserPrompt, error) {
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
		rows, err = m.db.QueryContext(ctx, query, limit)
	} else {
		// Return prompts for specific session
		query = `
			SELECT id, session_id, prompt_number, prompt_text, COALESCE(response, ''), created_at, created_at_epoch
			FROM user_prompts
			WHERE session_id = ?
			ORDER BY prompt_number ASC
			LIMIT ?
		`
		rows, err = m.db.QueryContext(ctx, query, sessionID, limit)
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
func (m *MySQL) UpdateLatestPromptResponse(ctx context.Context, sessionID string, response string) error {
	query := `
		UPDATE user_prompts
		SET response = ?
		WHERE id = (
			SELECT id FROM (
				SELECT id FROM user_prompts
				WHERE session_id = ?
				ORDER BY created_at DESC
				LIMIT 1
			) AS tmp
		)
	`
	_, err := m.db.ExecContext(ctx, query, response, sessionID)
	return err
}

// SearchFTS performs full-text search across observations and user_prompts.
// Note: MySQL uses FULLTEXT indexes instead of FTS5.
func (m *MySQL) SearchFTS(ctx context.Context, query string, types []string, limit int) ([]models.SearchResult, error) {
	if limit <= 0 {
		limit = 50
	}

	var results []models.SearchResult

	// Determine which types to search
	searchObservations := len(types) == 0
	searchPrompts := len(types) == 0
	for _, t := range types {
		if t == "observation" {
			searchObservations = true
		}
		if t == "prompt" {
			searchPrompts = true
		}
	}

	// Search observations using LIKE (MySQL FULLTEXT requires separate setup)
	if searchObservations {
		searchPattern := "%" + query + "%"
		obsQuery := `
			SELECT id, session_id, content, created_at
			FROM observations
			WHERE content LIKE ? OR COALESCE(title, '') LIKE ? OR COALESCE(narrative, '') LIKE ?
			ORDER BY created_at DESC
			LIMIT ?
		`
		rows, err := m.db.QueryContext(ctx, obsQuery, searchPattern, searchPattern, searchPattern, limit)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var r models.SearchResult
				r.Type = "observation"
				if err := rows.Scan(&r.ID, &r.SessionID, &r.Content, &r.CreatedAt); err != nil {
					continue
				}
				// Create snippet from content
				if len(r.Content) > 100 {
					r.Snippet = r.Content[:100] + "..."
				} else {
					r.Snippet = r.Content
				}
				results = append(results, r)
			}
		}
	}

	// Search user_prompts using LIKE
	if searchPrompts {
		searchPattern := "%" + query + "%"
		promptQuery := `
			SELECT id, session_id, prompt_text, created_at
			FROM user_prompts
			WHERE prompt_text LIKE ?
			ORDER BY created_at DESC
			LIMIT ?
		`
		rows, err := m.db.QueryContext(ctx, promptQuery, searchPattern, limit)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var r models.SearchResult
				r.Type = "prompt"
				if err := rows.Scan(&r.ID, &r.SessionID, &r.Content, &r.CreatedAt); err != nil {
					continue
				}
				// Create snippet from content
				if len(r.Content) > 100 {
					r.Snippet = r.Content[:100] + "..."
				} else {
					r.Snippet = r.Content
				}
				results = append(results, r)
			}
		}
	}

	return results, nil
}
