// Package models defines shared types for the Do Worker Service.
package models

import "time"

// =============================================================================
// Token Economics Constants
// =============================================================================

const (
	TokenBudgetMinimal  = 500
	TokenBudgetStandard = 2000
	TokenBudgetFull     = 5000
)

// ContextLevel defines the level of context to inject.
type ContextLevel int

const (
	LevelMinimal  ContextLevel = 1
	LevelStandard ContextLevel = 2
	LevelFull     ContextLevel = 3
)

// =============================================================================
// Core Entity Types
// =============================================================================

// Session represents a Claude session.
type Session struct {
	ID        string     `json:"id" db:"id"`
	UserName  string     `json:"user_name" db:"user_name"`
	ProjectID string     `json:"project_id,omitempty" db:"project_id"`
	StartedAt time.Time  `json:"started_at" db:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty" db:"ended_at"`
	Summary   string     `json:"summary,omitempty" db:"summary"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// Observation represents a recorded observation during a session.
type Observation struct {
	ID         int64     `json:"id" db:"id"`
	SessionID  string    `json:"session_id" db:"session_id"`
	AgentName  string    `json:"agent_name" db:"agent_name"`
	Type       string    `json:"type" db:"type"` // decision, pattern, learning, insight
	Content    string    `json:"content" db:"content"`
	Importance int       `json:"importance" db:"importance"` // 1-5
	Tags       string    `json:"tags,omitempty" db:"tags"`   // JSON array
	CreatedAt  time.Time `json:"created_at" db:"created_at"`

	// Extended fields for structured observations
	Title           *string `json:"title,omitempty" db:"title"`
	Subtitle        *string `json:"subtitle,omitempty" db:"subtitle"`
	Narrative       *string `json:"narrative,omitempty" db:"narrative"`
	Facts           string  `json:"facts,omitempty" db:"facts"`                   // JSON array
	Concepts        string  `json:"concepts,omitempty" db:"concepts"`             // JSON array
	FilesRead       string  `json:"files_read,omitempty" db:"files_read"`         // JSON array
	FilesModified   string  `json:"files_modified,omitempty" db:"files_modified"` // JSON array
	ResultPreview   *string `json:"result_preview,omitempty" db:"result_preview"`
	PromptNumber    *int    `json:"prompt_number,omitempty" db:"prompt_number"`
	DiscoveryTokens int     `json:"discovery_tokens" db:"discovery_tokens"`
}

// UserPrompt represents a user prompt within a session.
type UserPrompt struct {
	ID             int64     `json:"id" db:"id"`
	SessionID      string    `json:"session_id" db:"session_id"`
	PromptNumber   int       `json:"prompt_number" db:"prompt_number"`
	PromptText     string    `json:"prompt_text" db:"prompt_text"`
	Response       string    `json:"response,omitempty" db:"response"` // Assistant response with tool_use
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	CreatedAtEpoch int64     `json:"created_at_epoch" db:"created_at_epoch"`
}

// Summary represents a session or period summary.
type Summary struct {
	ID        int64     `json:"id" db:"id"`
	SessionID string    `json:"session_id,omitempty" db:"session_id"`
	Type      string    `json:"type" db:"type"` // session, daily, weekly
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Extended fields for structured summaries
	Request         *string `json:"request,omitempty" db:"request"`
	Investigated    *string `json:"investigated,omitempty" db:"investigated"`
	Learned         *string `json:"learned,omitempty" db:"learned"`
	Completed       *string `json:"completed,omitempty" db:"completed"`
	NextSteps       *string `json:"next_steps,omitempty" db:"next_steps"`
	FilesRead       string  `json:"files_read,omitempty" db:"files_read"`     // JSON array
	FilesEdited     string  `json:"files_edited,omitempty" db:"files_edited"` // JSON array
	DiscoveryTokens int     `json:"discovery_tokens" db:"discovery_tokens"`
	SourceMessage   string  `json:"source_message,omitempty" db:"source_message"`   // Original assistant message
	FullTranscript  string  `json:"full_transcript,omitempty" db:"full_transcript"` // Complete session transcript (JSONL)
}

// Plan represents a development plan.
type Plan struct {
	ID            int64     `json:"id" db:"id"`
	SessionID     string    `json:"session_id,omitempty" db:"session_id"`
	Title         string    `json:"title" db:"title"`
	Content       string    `json:"content" db:"content"`
	Status        string    `json:"status" db:"status"` // draft, active, completed
	FilePath      string    `json:"file_path,omitempty" db:"file_path"`
	RequestPrompt string    `json:"request_prompt,omitempty" db:"request_prompt"` // User request that triggered this plan
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// TeamContext represents context from team members.
type TeamContext struct {
	UserName     string    `json:"user_name" db:"user_name"`
	LastActivity time.Time `json:"last_activity" db:"last_activity"`
	Summary      string    `json:"summary" db:"summary"`
	ActivePlan   string    `json:"active_plan,omitempty" db:"active_plan"`
}

// Project represents a registered project with aggregated session data.
type Project struct {
	ID           string    `json:"id"`
	Path         string    `json:"path"`
	SessionCount int       `json:"session_count"`
	LastActivity time.Time `json:"last_activity"`
}

// =============================================================================
// Token Economics & Level Configuration
// =============================================================================

// LevelConfig defines configuration for a context level.
type LevelConfig struct {
	Level             ContextLevel
	MaxTokens         int
	ObservationLimit  int
	IncludePlan       bool
	IncludeTeam       bool
	IncludeTimeline   bool
	IncludePreviously bool
}

// TokenEconomics tracks token usage and efficiency.
type TokenEconomics struct {
	ReadTokens  int     `json:"read_tokens"`
	WorkTokens  int     `json:"work_tokens"`
	TotalBudget int     `json:"total_budget"`
	UsedTokens  int     `json:"used_tokens"`
	Savings     int     `json:"savings"`
	Efficiency  float64 `json:"efficiency"`
}

// =============================================================================
// API Response Types
// =============================================================================

// ContextInjectResponse is the response for context injection.
type ContextInjectResponse struct {
	Session      *Session       `json:"session,omitempty"`
	Observations []Observation  `json:"observations,omitempty"`
	ActivePlan   *Plan          `json:"active_plan,omitempty"`
	TeamContext  []TeamContext  `json:"team_context,omitempty"`
	Previously   string         `json:"previously,omitempty"`
	Economics    *TokenEconomics `json:"economics,omitempty"`
	Level        ContextLevel   `json:"level"`
	Markdown     string         `json:"markdown"`
}

// HealthResponse is the health check response.
type HealthResponse struct {
	Status   string `json:"status"`
	DBType   string `json:"db_type"`
	DBStatus string `json:"db_status"`
	Version  string `json:"version"`
}

// ErrorResponse is a standard error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SearchResult represents a single search result.
type SearchResult struct {
	Type      string    `json:"type"` // observation, summary, prompt
	ID        int64     `json:"id"`
	SessionID string    `json:"session_id,omitempty"`
	Content   string    `json:"content,omitempty"`
	Snippet   string    `json:"snippet"`
	Rank      float64   `json:"rank"`
	CreatedAt time.Time `json:"created_at"`
}

// SearchResponse is the response for FTS5 search.
type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Query   string         `json:"query"`
	Total   int            `json:"total"`
}

// =============================================================================
// API Request Types
// =============================================================================

// CreateSessionRequest is the request to create a new session.
type CreateSessionRequest struct {
	ID        string `json:"id" binding:"required"`
	UserName  string `json:"user_name" binding:"required"`
	ProjectID string `json:"project_id"`
}

// EndSessionRequest is the request to end a session.
type EndSessionRequest struct {
	Summary string `json:"summary,omitempty"`
}

// CreateObservationRequest is the request to create an observation.
type CreateObservationRequest struct {
	SessionID  string   `json:"session_id" binding:"required"`
	AgentName  string   `json:"agent_name"`
	Type       string   `json:"type" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Importance int      `json:"importance"`
	Tags       []string `json:"tags,omitempty"`
}

// CreateSummaryRequest is the request to create a summary.
type CreateSummaryRequest struct {
	SessionID string `json:"session_id,omitempty"`
	Type      string `json:"type" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// CreatePlanRequest is the request to create a plan.
type CreatePlanRequest struct {
	SessionID     string `json:"session_id,omitempty"`
	Title         string `json:"title" binding:"required"`
	Content       string `json:"content" binding:"required"`
	FilePath      string `json:"file_path,omitempty"`
	RequestPrompt string `json:"request_prompt,omitempty"` // User request that triggered this plan
}

// GenerateSummaryRequest is the request to generate a session summary.
type GenerateSummaryRequest struct {
	SessionID            string                 `json:"session_id" binding:"required"`
	LastAssistantMessage string                 `json:"last_assistant_message"`
	FullTranscript       string                 `json:"full_transcript,omitempty"` // Complete session transcript (JSONL)
	TranscriptStats      map[string]interface{} `json:"transcript_stats,omitempty"`
}

// CreateUserPromptRequest is the request to create a user prompt record.
type CreateUserPromptRequest struct {
	SessionID    string `json:"session_id" binding:"required"`
	PromptNumber int    `json:"prompt_number" binding:"required"`
	PromptText   string `json:"prompt_text" binding:"required"`
}

// SearchRequest is the request for searching observations and summaries.
type SearchRequest struct {
	Query string   `json:"query" form:"q" binding:"required"`
	Types []string `json:"types" form:"types"`
	Limit int      `json:"limit" form:"limit"`
}
