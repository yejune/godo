// Package context provides context building utilities for session injection.
package context

import (
	"context"

	"github.com/do-focus/worker/internal/db"
	"github.com/do-focus/worker/pkg/models"
)

// Builder constructs context for session injection.
type Builder struct {
	db       db.Adapter
	renderer *Renderer
}

// NewBuilder creates a new context builder.
func NewBuilder(adapter db.Adapter) *Builder {
	return &Builder{
		db:       adapter,
		renderer: NewRenderer(),
	}
}

// =============================================================================
// Level Configuration - Progressive Disclosure
// =============================================================================

// GetLevelConfig returns configuration for a given level.
// Implements progressive disclosure based on context level:
// - Minimal (Level 1): Essential context only, fastest injection
// - Standard (Level 2): Balanced context with plan and timeline
// - Full (Level 3): Complete context including team and all history
func GetLevelConfig(level models.ContextLevel) models.LevelConfig {
	switch level {
	case models.LevelMinimal:
		return models.LevelConfig{
			Level:             models.LevelMinimal,
			MaxTokens:         models.TokenBudgetMinimal,
			ObservationLimit:  5,
			IncludePlan:       false,
			IncludeTeam:       false,
			IncludeTimeline:   false,
			IncludePreviously: false,
		}
	case models.LevelFull:
		return models.LevelConfig{
			Level:             models.LevelFull,
			MaxTokens:         models.TokenBudgetFull,
			ObservationLimit:  50,
			IncludePlan:       true,
			IncludeTeam:       true,
			IncludeTimeline:   true,
			IncludePreviously: true,
		}
	default: // LevelStandard
		return models.LevelConfig{
			Level:             models.LevelStandard,
			MaxTokens:         models.TokenBudgetStandard,
			ObservationLimit:  20,
			IncludePlan:       true,
			IncludeTeam:       false,
			IncludeTimeline:   true,
			IncludePreviously: true,
		}
	}
}

// =============================================================================
// Build Request - Fluent API
// =============================================================================

// BuildRequest holds parameters for context building with fluent API.
type BuildRequest struct {
	UserName    string              `json:"user_name"`
	ProjectPath string              `json:"project_path"`
	SessionID   string              `json:"session_id"`
	Level       models.ContextLevel `json:"level"`
	Previously  string              `json:"previously"` // Last assistant response for continuity
}

// NewBuildRequest creates a new build request with defaults.
func NewBuildRequest(userName string) BuildRequest {
	return BuildRequest{
		UserName: userName,
		Level:    models.LevelStandard,
	}
}

// WithLevel sets the context level.
func (r BuildRequest) WithLevel(level models.ContextLevel) BuildRequest {
	r.Level = level
	return r
}

// WithProject sets the project path.
func (r BuildRequest) WithProject(path string) BuildRequest {
	r.ProjectPath = path
	return r
}

// WithSession sets the session ID.
func (r BuildRequest) WithSession(id string) BuildRequest {
	r.SessionID = id
	return r
}

// WithPreviously sets the previous assistant response.
func (r BuildRequest) WithPreviously(text string) BuildRequest {
	r.Previously = text
	return r
}

// =============================================================================
// Context Building - Level-based API
// =============================================================================

// BuildContextWithLevel assembles context based on level configuration.
// This is the primary method for level-based progressive disclosure.
func (b *Builder) BuildContextWithLevel(ctx context.Context, req BuildRequest) (*models.ContextInjectResponse, error) {
	cfg := GetLevelConfig(req.Level)

	resp := &models.ContextInjectResponse{
		Level: req.Level,
	}

	// 1. Session info (always included)
	session, err := b.db.GetLatestSession(ctx, req.UserName)
	if err != nil {
		return nil, err
	}
	resp.Session = session

	// 2. Observations based on level limit
	observations, err := b.db.GetRecentObservations(ctx, req.UserName, cfg.ObservationLimit)
	if err != nil {
		return nil, err
	}
	resp.Observations = observations

	// 3. Plan (if level allows)
	if cfg.IncludePlan {
		plan, _ := b.db.GetActivePlan(ctx, req.UserName)
		resp.ActivePlan = plan
	}

	// 4. Team context (if level allows - Full only)
	if cfg.IncludeTeam {
		team, _ := b.db.GetTeamContext(ctx, req.UserName)
		resp.TeamContext = team
	}

	// 5. Previously summary (if level allows and not provided in request)
	if cfg.IncludePreviously && req.Previously == "" {
		previously, err := b.db.GetLatestSummary(ctx, req.UserName)
		if err == nil && previously != nil {
			resp.Previously = previously.Content
		}
	} else if req.Previously != "" {
		resp.Previously = req.Previously
	}

	// 6. Calculate token economics
	resp.Economics = b.calculateEconomics(resp, req.Level)

	// 7. Render markdown with level-appropriate sections
	resp.Markdown = b.renderer.RenderContext(resp)

	return resp, nil
}

// BuildContext assembles the full context for a user.
// Deprecated: Use BuildContextWithLevel for level-based progressive disclosure.
func (b *Builder) BuildContext(ctx context.Context, userName string, opts BuildOptions) (*models.ContextInjectResponse, error) {
	resp := &models.ContextInjectResponse{
		Level: opts.Level,
	}

	// Get latest session
	session, err := b.db.GetLatestSession(ctx, userName)
	if err != nil {
		return nil, err
	}
	resp.Session = session

	// Get recent observations
	limit := opts.ObservationLimit
	if limit <= 0 {
		limit = 20
	}
	observations, err := b.db.GetRecentObservations(ctx, userName, limit)
	if err != nil {
		return nil, err
	}
	resp.Observations = observations

	// Get active plan
	if opts.IncludePlan {
		plan, err := b.db.GetActivePlan(ctx, userName)
		if err != nil {
			return nil, err
		}
		resp.ActivePlan = plan
	}

	// Get team context (Level 3 only)
	if opts.IncludeTeam && opts.Level >= models.LevelFull {
		team, err := b.db.GetTeamContext(ctx, userName)
		if err != nil {
			return nil, err
		}
		resp.TeamContext = team
	}

	// Get previously summary (Level 2+)
	if opts.IncludePreviously && opts.Level >= models.LevelStandard {
		previously, err := b.db.GetLatestSummary(ctx, userName)
		if err == nil && previously != nil {
			resp.Previously = previously.Content
		}
	}

	// Calculate token economics
	resp.Economics = b.calculateEconomics(resp, opts.Level)

	// Render markdown
	resp.Markdown = b.renderer.RenderContext(resp)

	return resp, nil
}

// calculateEconomics estimates token usage and efficiency.
func (b *Builder) calculateEconomics(resp *models.ContextInjectResponse, level models.ContextLevel) *models.TokenEconomics {
	// Token budget based on level
	var budget int
	switch level {
	case models.LevelMinimal:
		budget = models.TokenBudgetMinimal
	case models.LevelStandard:
		budget = models.TokenBudgetStandard
	case models.LevelFull:
		budget = models.TokenBudgetFull
	default:
		budget = models.TokenBudgetStandard
	}

	// Estimate used tokens (rough: ~4 chars per token)
	var usedTokens int

	if resp.Session != nil {
		usedTokens += 50 // Session header
	}

	for _, obs := range resp.Observations {
		usedTokens += len(obs.Content) / 4
		usedTokens += 20 // Observation metadata
	}

	if resp.ActivePlan != nil {
		usedTokens += len(resp.ActivePlan.Content) / 4
		usedTokens += 30 // Plan metadata
	}

	if resp.Previously != "" {
		usedTokens += len(resp.Previously) / 4
	}

	for range resp.TeamContext {
		usedTokens += 40 // Team member context
	}

	// Calculate efficiency
	var efficiency float64
	if budget > 0 {
		efficiency = float64(usedTokens) / float64(budget) * 100
		if efficiency > 100 {
			efficiency = 100
		}
	}

	return &models.TokenEconomics{
		TotalBudget: budget,
		UsedTokens:  usedTokens,
		Savings:     budget - usedTokens,
		Efficiency:  efficiency,
	}
}

// BuildOptions configures context building.
type BuildOptions struct {
	ObservationLimit  int
	IncludePlan       bool
	IncludeTeam       bool
	IncludeSession    bool
	IncludeTimeline   bool
	IncludePreviously bool
	Level             models.ContextLevel
}

// DefaultBuildOptions returns the default build options (Standard level).
func DefaultBuildOptions() BuildOptions {
	return BuildOptions{
		ObservationLimit:  20,
		IncludePlan:       true,
		IncludeTeam:       false,
		IncludeSession:    true,
		IncludeTimeline:   true,
		IncludePreviously: true,
		Level:             models.LevelStandard,
	}
}

// MinimalBuildOptions returns minimal build options.
func MinimalBuildOptions() BuildOptions {
	return BuildOptions{
		ObservationLimit:  10,
		IncludePlan:       false,
		IncludeTeam:       false,
		IncludeSession:    true,
		IncludeTimeline:   false,
		IncludePreviously: false,
		Level:             models.LevelMinimal,
	}
}

// FullBuildOptions returns full build options.
func FullBuildOptions() BuildOptions {
	return BuildOptions{
		ObservationLimit:  50,
		IncludePlan:       true,
		IncludeTeam:       true,
		IncludeSession:    true,
		IncludeTimeline:   true,
		IncludePreviously: true,
		Level:             models.LevelFull,
	}
}

// BuildOptionsFromLevel creates BuildOptions from a ContextLevel.
// Useful for migrating from Level-based API to options-based API.
func BuildOptionsFromLevel(level models.ContextLevel) BuildOptions {
	cfg := GetLevelConfig(level)
	return BuildOptions{
		ObservationLimit:  cfg.ObservationLimit,
		IncludePlan:       cfg.IncludePlan,
		IncludeTeam:       cfg.IncludeTeam,
		IncludeSession:    true,
		IncludeTimeline:   cfg.IncludeTimeline,
		IncludePreviously: cfg.IncludePreviously,
		Level:             level,
	}
}

// ToLevel infers the ContextLevel from BuildOptions.
// Used for migration from legacy API.
func (opts BuildOptions) ToLevel() models.ContextLevel {
	// If Level is already set, use it
	if opts.Level > 0 {
		return opts.Level
	}

	// Infer from options
	if opts.ObservationLimit <= 10 && !opts.IncludePlan && !opts.IncludeTeam {
		return models.LevelMinimal
	}
	if opts.ObservationLimit >= 50 && opts.IncludePlan && opts.IncludeTeam {
		return models.LevelFull
	}
	return models.LevelStandard
}
