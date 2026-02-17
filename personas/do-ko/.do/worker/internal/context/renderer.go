package context

import (
	"fmt"
	"strings"
	"time"

	"github.com/do-focus/worker/pkg/models"
)

// Renderer converts context data to markdown format.
type Renderer struct {
	maxObservations int
	maxTeamMembers  int
	maxPreviously   int
}

// NewRenderer creates a new renderer with default settings.
func NewRenderer() *Renderer {
	return &Renderer{
		maxObservations: 20,
		maxTeamMembers:  5,
		maxPreviously:   500,
	}
}

// RenderContext generates markdown from context data.
func (r *Renderer) RenderContext(ctx *models.ContextInjectResponse) string {
	var sb strings.Builder

	// Header with Token Economics
	sb.WriteString("# Do Memory Context\n\n")

	// Render token economics (always show)
	if ctx.Economics != nil {
		r.renderEconomics(&sb, ctx.Economics, ctx.Level)
	}

	// Render session
	if ctx.Session != nil {
		r.renderSession(&sb, ctx.Session)
	}

	// Render previously section (Level 2+)
	if ctx.Previously != "" && ctx.Level >= models.LevelStandard {
		r.renderPreviously(&sb, ctx.Previously)
	}

	// Render timeline (observations grouped by date)
	if len(ctx.Observations) > 0 {
		r.renderTimeline(&sb, ctx.Observations)
	}

	// Render active plan
	if ctx.ActivePlan != nil {
		r.renderPlan(&sb, ctx.ActivePlan)
	}

	// Render team context (Level 3 only)
	if len(ctx.TeamContext) > 0 && ctx.Level >= models.LevelFull {
		r.renderTeamContext(&sb, ctx.TeamContext)
	}

	return sb.String()
}

// renderEconomics renders token economics header.
func (r *Renderer) renderEconomics(sb *strings.Builder, econ *models.TokenEconomics, level models.ContextLevel) {
	levelName := "Standard"
	switch level {
	case models.LevelMinimal:
		levelName = "Minimal"
	case models.LevelFull:
		levelName = "Full"
	}

	sb.WriteString(fmt.Sprintf("**Context Level**: %s | **Tokens**: %d/%d (%.0f%% used)\n\n",
		levelName, econ.UsedTokens, econ.TotalBudget, econ.Efficiency))
}

// renderSession renders session information.
func (r *Renderer) renderSession(sb *strings.Builder, session *models.Session) {
	sb.WriteString("## Current Session\n\n")
	sb.WriteString(fmt.Sprintf("- **ID**: `%s`\n", session.ID))
	sb.WriteString(fmt.Sprintf("- **User**: %s\n", session.UserName))
	sb.WriteString(fmt.Sprintf("- **Started**: %s\n", session.StartedAt.Format(time.RFC3339)))

	if session.EndedAt != nil {
		sb.WriteString(fmt.Sprintf("- **Ended**: %s\n", session.EndedAt.Format(time.RFC3339)))
		duration := session.EndedAt.Sub(session.StartedAt)
		sb.WriteString(fmt.Sprintf("- **Duration**: %s\n", formatDuration(duration)))
	} else {
		duration := time.Since(session.StartedAt)
		sb.WriteString(fmt.Sprintf("- **Active for**: %s\n", formatDuration(duration)))
	}

	if session.Summary != "" {
		sb.WriteString(fmt.Sprintf("\n**Summary**: %s\n", session.Summary))
	}
	sb.WriteString("\n")
}

// renderPreviously renders the previous session summary.
func (r *Renderer) renderPreviously(sb *strings.Builder, previously string) {
	sb.WriteString("## Previously\n\n")

	// Truncate if too long
	content := previously
	if len(content) > r.maxPreviously {
		content = content[:r.maxPreviously] + "..."
	}

	sb.WriteString(fmt.Sprintf("_%s_\n\n", content))
}

// renderTimeline renders observations grouped by date in chronological order.
func (r *Renderer) renderTimeline(sb *strings.Builder, observations []models.Observation) {
	sb.WriteString("## Timeline\n\n")

	// Group observations by date
	byDate := make(map[string][]models.Observation)
	var dateOrder []string
	dateSet := make(map[string]bool)

	for _, obs := range observations {
		obsDate := obs.CreatedAt.Format("2006-01-02")
		if !dateSet[obsDate] {
			dateSet[obsDate] = true
			dateOrder = append(dateOrder, obsDate)
		}
		byDate[obsDate] = append(byDate[obsDate], obs)
	}

	// Render by date (most recent first)
	for i := len(dateOrder) - 1; i >= 0; i-- {
		date := dateOrder[i]
		obs := byDate[date]

		// Relative date label
		dateLabel := r.getRelativeDateLabel(date)
		sb.WriteString(fmt.Sprintf("### %s\n\n", dateLabel))

		// Render observations for this date
		for _, o := range obs {
			r.renderTimelineObservation(sb, o)
		}
		sb.WriteString("\n")
	}
}

// renderTimelineObservation renders a single observation in timeline format.
func (r *Renderer) renderTimelineObservation(sb *strings.Builder, obs models.Observation) {
	timeStr := obs.CreatedAt.Format("15:04")
	typeLabel := fmt.Sprintf("[%s]", obs.Type)

	importance := ""
	if obs.Importance >= 4 {
		importance = " **[!]**"
	}

	agent := ""
	if obs.AgentName != "" {
		agent = fmt.Sprintf(" _%s_", obs.AgentName)
	}

	sb.WriteString(fmt.Sprintf("- `%s` %s %s%s%s\n", timeStr, typeLabel, obs.Content, importance, agent))
}

// renderPlan renders the active plan.
func (r *Renderer) renderPlan(sb *strings.Builder, plan *models.Plan) {
	sb.WriteString("## Active Plan\n\n")
	sb.WriteString(fmt.Sprintf("### %s\n\n", plan.Title))
	sb.WriteString(fmt.Sprintf("**Status**: %s\n", plan.Status))

	if plan.FilePath != "" {
		sb.WriteString(fmt.Sprintf("**File**: `%s`\n", plan.FilePath))
	}
	sb.WriteString("\n")
}

// renderTeamContext renders team member activity.
func (r *Renderer) renderTeamContext(sb *strings.Builder, team []models.TeamContext) {
	sb.WriteString("## Team Activity\n\n")

	for i, t := range team {
		if i >= r.maxTeamMembers {
			break
		}

		sb.WriteString(fmt.Sprintf("### %s\n\n", t.UserName))
		sb.WriteString(fmt.Sprintf("- **Last Active**: %s\n", t.LastActivity.Format("2006-01-02 15:04")))

		if t.Summary != "" {
			sb.WriteString(fmt.Sprintf("- **Last Work**: %s\n", t.Summary))
		}

		if t.ActivePlan != "" {
			sb.WriteString(fmt.Sprintf("- **Working On**: %s\n", t.ActivePlan))
		}

		sb.WriteString("\n")
	}
}

// Helper functions

// getRelativeDateLabel returns a human-readable date label.
func (r *Renderer) getRelativeDateLabel(dateStr string) string {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	daysAgo := int(today.Sub(date).Hours() / 24)

	switch daysAgo {
	case 0:
		return "Today"
	case 1:
		return "Yesterday"
	case 2, 3, 4, 5, 6:
		return fmt.Sprintf("%d days ago", daysAgo)
	default:
		return dateStr
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if minutes > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dh", hours)
}
