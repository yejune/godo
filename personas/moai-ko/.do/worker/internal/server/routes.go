package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/do-focus/worker/pkg/models"
	"github.com/gin-gonic/gin"
)

// Version is set by main package
var Version = "dev"

// setupRoutes configures all API routes.
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.handleHealth)

	// API routes
	api := s.router.Group("/api")
	{
		// Context injection for SessionStart hook
		api.GET("/context/inject", s.handleContextInject)

		// Session management
		api.GET("/sessions", s.handleGetSessions)
		api.GET("/sessions/:id", s.handleGetSession)
		api.POST("/sessions", s.handleCreateSession)
		api.PUT("/sessions/:id/end", s.handleEndSession)

		// Observations
		api.GET("/observations", s.handleGetObservations)
		api.GET("/observations/search", s.handleSearchObservations)
		api.POST("/observations", s.handleCreateObservation)

		// Summaries
		api.GET("/summaries", s.handleGetSummaries)
		api.POST("/summaries", s.handleCreateSummary)
		api.POST("/summaries/generate", s.handleGenerateSummary)

		// User Prompts
		api.GET("/prompts", s.handleGetUserPrompts)
		api.POST("/prompts", s.handleCreateUserPrompt)
		api.PUT("/prompts/latest/response", s.handleUpdateLatestPromptResponse)

		// FTS5 Search
		api.GET("/search", s.handleSearch)

		// Plans
		api.GET("/plans", s.handleGetPlans)
		api.POST("/plans", s.handleCreatePlan)

		// Team context
		api.GET("/team/context", s.handleTeamContext)

		// Projects
		api.GET("/projects", s.getProjects)
	}
}

// handleHealth handles the health check endpoint.
func (s *Server) handleHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	dbStatus := "ok"
	if err := s.db.Health(ctx); err != nil {
		dbStatus = "error: " + err.Error()
	}

	dbType := os.Getenv("DO_DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}

	c.JSON(http.StatusOK, models.HealthResponse{
		Status:   "ok",
		DBType:   dbType,
		DBStatus: dbStatus,
		Version:  Version,
	})
}

// handleContextInject handles context injection for SessionStart hook.
// level parameter controls the amount of data returned:
// - level 1: minimal (session only)
// - level 2: standard (session + observations) [default]
// - level 3: full (session + observations + plan + team)
func (s *Server) handleContextInject(c *gin.Context) {
	ctx := c.Request.Context()

	userName := c.Query("user")
	if userName == "" {
		userName = os.Getenv("DO_USER_NAME")
	}
	if userName == "" {
		userName = "default"
	}

	// Parse level parameter (1-3, default 2)
	levelStr := c.DefaultQuery("level", "2")
	level, _ := strconv.Atoi(levelStr)
	if level < 1 {
		level = 1
	}
	if level > 3 {
		level = 3
	}

	// Get latest session (always included)
	session, err := s.db.GetLatestSession(ctx, userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	var observations []models.Observation
	var plan *models.Plan
	var teamContext []models.TeamContext

	// Level 2+: Include observations
	if level >= 2 {
		limitStr := c.DefaultQuery("obs_limit", "20")
		limit, _ := strconv.Atoi(limitStr)
		if limit <= 0 {
			limit = 20
		}
		// Adjust limit based on level
		if level == 2 {
			if limit > 20 {
				limit = 20
			}
		}

		observations, err = s.db.GetRecentObservations(ctx, userName, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "database_error",
				Message: err.Error(),
			})
			return
		}
	}

	// Level 3: Include plan and team context
	if level >= 3 {
		plan, err = s.db.GetActivePlan(ctx, userName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "database_error",
				Message: err.Error(),
			})
			return
		}

		teamContext, err = s.db.GetTeamContext(ctx, userName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "database_error",
				Message: err.Error(),
			})
			return
		}
	}

	// Build markdown response
	markdown := buildContextMarkdown(session, observations, plan, teamContext)

	c.JSON(http.StatusOK, models.ContextInjectResponse{
		Session:      session,
		Observations: observations,
		ActivePlan:   plan,
		TeamContext:  teamContext,
		Markdown:     markdown,
	})
}

// handleCreateSession handles session creation (idempotent).
func (s *Server) handleCreateSession(c *gin.Context) {
	var req models.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Check if session already exists (idempotent)
	existing, _ := s.db.GetSession(c.Request.Context(), req.ID)
	if existing != nil {
		c.JSON(http.StatusOK, existing)
		return
	}

	session := &models.Session{
		ID:        req.ID,
		UserName:  req.UserName,
		ProjectID: req.ProjectID,
		StartedAt: time.Now(),
	}

	if err := s.db.CreateSession(c.Request.Context(), session); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// handleEndSession handles session ending.
func (s *Server) handleEndSession(c *gin.Context) {
	id := c.Param("id")

	var req models.EndSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Allow empty body
		req = models.EndSessionRequest{}
	}

	if err := s.db.EndSession(c.Request.Context(), id, req.Summary); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ended"})
}

// handleGetSessions handles session list retrieval.
func (s *Server) handleGetSessions(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 20
	}

	sessions, err := s.db.GetRecentSessions(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// handleGetSession handles single session retrieval.
func (s *Server) handleGetSession(c *gin.Context) {
	id := c.Param("id")

	session, err := s.db.GetSession(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Session not found",
		})
		return
	}

	c.JSON(http.StatusOK, session)
}

// handleGetObservations handles observation list retrieval.
func (s *Server) handleGetObservations(c *gin.Context) {
	sessionID := c.Query("session_id")
	obsType := c.Query("type")
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)

	observations, err := s.db.GetObservationsFiltered(c.Request.Context(), sessionID, obsType, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, observations)
}

// handleSearchObservations handles observation search.
func (s *Server) handleSearchObservations(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Query parameter 'q' is required",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	results, err := s.db.SearchObservations(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}

// handleGetSummaries handles summary list retrieval.
func (s *Server) handleGetSummaries(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "7")
	days, _ := strconv.Atoi(daysStr)
	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	summaries, err := s.db.GetAllSummaries(c.Request.Context(), days, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, summaries)
}

// handleGetPlans handles plan list retrieval.
func (s *Server) handleGetPlans(c *gin.Context) {
	sessionID := c.Query("session_id")
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	plans, err := s.db.GetAllPlans(c.Request.Context(), sessionID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// handleCreateObservation handles observation creation.
func (s *Server) handleCreateObservation(c *gin.Context) {
	var req models.CreateObservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Default importance to 3
	if req.Importance <= 0 || req.Importance > 5 {
		req.Importance = 3
	}

	// Convert tags to JSON
	var tagsJSON string
	if len(req.Tags) > 0 {
		tagsBytes, _ := json.Marshal(req.Tags)
		tagsJSON = string(tagsBytes)
	}

	obs := &models.Observation{
		SessionID:  req.SessionID,
		AgentName:  req.AgentName,
		Type:       req.Type,
		Content:    req.Content,
		Importance: req.Importance,
		Tags:       tagsJSON,
	}

	if err := s.db.CreateObservation(c.Request.Context(), obs); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, obs)
}

// handleCreateSummary handles summary creation.
func (s *Server) handleCreateSummary(c *gin.Context) {
	var req models.CreateSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	summary := &models.Summary{
		SessionID: req.SessionID,
		Type:      req.Type,
		Content:   req.Content,
	}

	if err := s.db.CreateSummary(c.Request.Context(), summary); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, summary)
}

// handleCreatePlan handles plan creation.
func (s *Server) handleCreatePlan(c *gin.Context) {
	var req models.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	plan := &models.Plan{
		SessionID:     req.SessionID,
		Title:         req.Title,
		Content:       req.Content,
		Status:        "draft",
		FilePath:      req.FilePath,
		RequestPrompt: req.RequestPrompt,
	}

	if err := s.db.CreatePlan(c.Request.Context(), plan); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

// handleTeamContext handles team context retrieval.
func (s *Server) handleTeamContext(c *gin.Context) {
	userName := c.Query("exclude_user")
	if userName == "" {
		userName = os.Getenv("DO_USER_NAME")
	}

	contexts, err := s.db.GetTeamContext(c.Request.Context(), userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"team": contexts})
}

// getProjects handles project list retrieval.
func (s *Server) getProjects(c *gin.Context) {
	projects, err := s.db.GetProjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// handleGenerateSummary generates a rule-based summary from session observations.
func (s *Server) handleGenerateSummary(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.GenerateSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// 1. Verify session exists
	session, err := s.db.GetSession(ctx, req.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}
	if session == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "session_not_found",
			Message: "Session not found: " + req.SessionID,
		})
		return
	}

	// 2. Get observations for the session
	observations, err := s.db.GetObservationsFiltered(ctx, req.SessionID, "", 100, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	// 3. Get user prompts for request extraction
	userPrompts, _ := s.db.GetUserPrompts(ctx, req.SessionID, 10)
	var promptTexts []string
	for _, p := range userPrompts {
		promptTexts = append(promptTexts, p.PromptText)
	}

	// 4. Generate structured summary (LLM if available, else rule-based)
	var structured StructuredSummary
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey != "" && req.LastAssistantMessage != "" {
		// Try LLM-based summary
		llmSummary, err := generateLLMSummary(ctx, req.LastAssistantMessage, promptTexts, apiKey)
		if err == nil && llmSummary != nil {
			structured = *llmSummary
		} else {
			// Fallback to rule-based
			structured = generateStructuredSummary(observations, req.LastAssistantMessage, promptTexts)
		}
	} else {
		structured = generateStructuredSummary(observations, req.LastAssistantMessage, promptTexts)
	}

	// 5. Convert to JSON for file lists
	filesReadJSON, _ := json.Marshal(structured.FilesRead)
	filesEditedJSON, _ := json.Marshal(structured.FilesEdited)

	// 6. Save summary to DB with structured fields
	// Truncate source message to 50KB for storage (includes tool_use)
	sourceMsg := req.LastAssistantMessage
	if len(sourceMsg) > 50000 {
		sourceMsg = sourceMsg[:50000] + "\n...(truncated)"
	}

	// Truncate full transcript to 500KB for storage
	fullTranscript := req.FullTranscript
	if len(fullTranscript) > 500000 {
		fullTranscript = fullTranscript[:500000] + "\n...(truncated)"
	}

	summary := &models.Summary{
		SessionID:      req.SessionID,
		Type:           "session",
		Content:        formatSummaryAsMarkdown(structured),
		Request:        strPtr(structured.Request),
		Investigated:   strPtr(structured.Investigated),
		Learned:        strPtr(structured.Learned),
		Completed:      strPtr(structured.Completed),
		NextSteps:      strPtr(structured.NextSteps),
		FilesRead:      string(filesReadJSON),
		FilesEdited:    string(filesEditedJSON),
		SourceMessage:  sourceMsg,
		FullTranscript: fullTranscript,
	}

	if err := s.db.CreateSummary(ctx, summary); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, summary)
}

// StructuredSummary holds parsed summary fields.
type StructuredSummary struct {
	Request      string
	Investigated string
	Learned      string
	Completed    string
	NextSteps    string
	FilesRead    []string
	FilesEdited  []string
}

// generateStructuredSummary extracts key information from lastMessage and observations.
// Enhanced version with better pattern matching and observation utilization.
func generateStructuredSummary(observations []models.Observation, lastMessage string, userPrompts []string) StructuredSummary {
	summary := StructuredSummary{}

	// 1. Request: 첫 번째 user prompt 사용 (가장 신뢰할 수 있는 소스)
	if len(userPrompts) > 0 {
		req := userPrompts[0]
		if len(req) > 500 {
			req = req[:500] + "..."
		}
		summary.Request = req
	}

	// 2. observations에서 풍부한 정보 추출 (lastMessage보다 우선)
	filesReadMap := make(map[string]bool)
	filesEditedMap := make(map[string]bool)
	var completedItems []string
	var investigatedItems []string
	var learnedItems []string
	var nextStepItems []string

	for _, obs := range observations {
		// FilesRead, FilesModified 필드 활용 (JSON 배열)
		if obs.FilesRead != "" && obs.FilesRead != "[]" {
			var files []string
			if err := json.Unmarshal([]byte(obs.FilesRead), &files); err == nil {
				for _, f := range files {
					if f != "" && !filesReadMap[f] {
						filesReadMap[f] = true
						summary.FilesRead = append(summary.FilesRead, f)
					}
				}
			}
		}
		if obs.FilesModified != "" && obs.FilesModified != "[]" {
			var files []string
			if err := json.Unmarshal([]byte(obs.FilesModified), &files); err == nil {
				for _, f := range files {
					if f != "" && !filesEditedMap[f] {
						filesEditedMap[f] = true
						summary.FilesEdited = append(summary.FilesEdited, f)
					}
				}
			}
		}

		// Title/Subtitle 활용 (더 구조화된 정보)
		title := ""
		if obs.Title != nil && *obs.Title != "" {
			title = *obs.Title
		}
		subtitle := ""
		if obs.Subtitle != nil && *obs.Subtitle != "" {
			subtitle = *obs.Subtitle
		}

		// 타입별 정보 추출
		switch obs.Type {
		case "read", "exploration":
			path := extractFilePath(obs.Content)
			if path != "" && !filesReadMap[path] {
				filesReadMap[path] = true
				summary.FilesRead = append(summary.FilesRead, path)
			}
			// 탐색 내용을 investigated에 추가
			if title != "" {
				investigatedItems = append(investigatedItems, title)
			}
		case "feature", "edit", "write", "implementation":
			path := extractFilePath(obs.Content)
			if path != "" && !filesEditedMap[path] {
				filesEditedMap[path] = true
				summary.FilesEdited = append(summary.FilesEdited, path)
			}
			// 구현된 기능을 completed에 추가
			if title != "" {
				completedItems = append(completedItems, title)
			} else if obs.Content != "" {
				completedItems = append(completedItems, obs.Content)
			}
		case "decision":
			if title != "" {
				completedItems = append(completedItems, title)
			} else {
				completedItems = append(completedItems, obs.Content)
			}
		case "bugfix", "fix":
			item := "Fixed: "
			if title != "" {
				item += title
			} else {
				item += obs.Content
			}
			completedItems = append(completedItems, item)
		case "learning", "insight", "discovery":
			item := ""
			if title != "" {
				item = title
				if subtitle != "" {
					item += " - " + subtitle
				}
			} else {
				item = obs.Content
			}
			if item != "" {
				learnedItems = append(learnedItems, item)
			}
		case "commit":
			// 커밋 메시지는 작업 완료의 좋은 요약
			if obs.Content != "" {
				completedItems = append(completedItems, obs.Content)
			}
		case "test", "testing":
			if title != "" {
				completedItems = append(completedItems, "Tested: "+title)
			}
		case "analysis", "investigation":
			if title != "" {
				investigatedItems = append(investigatedItems, title)
			} else if obs.Content != "" {
				investigatedItems = append(investigatedItems, obs.Content)
			}
		case "todo", "next", "plan":
			if title != "" {
				nextStepItems = append(nextStepItems, title)
			} else if obs.Content != "" {
				nextStepItems = append(nextStepItems, obs.Content)
			}
		}

		// ResultPreview 활용 (있으면)
		if obs.ResultPreview != nil && *obs.ResultPreview != "" {
			preview := *obs.ResultPreview
			// 결과 미리보기에서 추가 정보 추출
			if strings.Contains(strings.ToLower(preview), "error") || strings.Contains(strings.ToLower(preview), "failed") {
				// 에러가 있으면 investigated에 추가
				investigatedItems = append(investigatedItems, "Investigated issue: "+truncateString(preview, 100))
			}
		}

		// Narrative 활용 (상세 설명)
		if obs.Narrative != nil && *obs.Narrative != "" {
			narrative := *obs.Narrative
			// 긴 narrative에서 핵심 추출
			if len(narrative) > 200 {
				narrative = narrative[:200] + "..."
			}
			// investigation 타입이면 investigated에 추가
			if obs.Type == "analysis" || obs.Type == "investigation" || obs.Type == "exploration" {
				if !containsString(investigatedItems, narrative) {
					investigatedItems = append(investigatedItems, narrative)
				}
			}
		}
	}

	// 3. lastMessage에서 추가 정보 추출 (observations에서 못 찾은 것만)
	if lastMessage != "" {
		parsed := parseMarkdownStructure(lastMessage)

		// lastMessage에서 추출한 내용 보완
		if summary.Request == "" && parsed.Request != "" {
			summary.Request = parsed.Request
		}

		// observations에서 이미 수집한 것과 병합
		if parsed.Investigated != "" && len(investigatedItems) == 0 {
			investigatedItems = append(investigatedItems, parsed.Investigated)
		}
		if parsed.Learned != "" && len(learnedItems) == 0 {
			learnedItems = append(learnedItems, parsed.Learned)
		}
		if parsed.Completed != "" && len(completedItems) == 0 {
			completedItems = append(completedItems, parsed.Completed)
		}
		if parsed.NextSteps != "" && len(nextStepItems) == 0 {
			nextStepItems = append(nextStepItems, parsed.NextSteps)
		}
	}

	// 4. 수집된 항목들을 요약 필드로 변환
	if len(investigatedItems) > 0 {
		summary.Investigated = formatItems(investigatedItems, 5)
	}
	if len(learnedItems) > 0 {
		summary.Learned = formatItems(learnedItems, 5)
	}
	if len(completedItems) > 0 {
		summary.Completed = formatItems(completedItems, 10)
	}
	if len(nextStepItems) > 0 {
		summary.NextSteps = formatItems(nextStepItems, 5)
	}

	// 5. 아무 정보도 없으면 lastMessage에서 동사 기반 추출 시도
	if summary.Completed == "" && lastMessage != "" {
		summary.Completed = extractActionItems(lastMessage)
	}

	return summary
}

// formatItems formats a list of items as bullet points, limiting to maxItems.
func formatItems(items []string, maxItems int) string {
	// 중복 제거
	seen := make(map[string]bool)
	unique := make([]string, 0)
	for _, item := range items {
		normalized := strings.TrimSpace(item)
		if normalized != "" && !seen[normalized] {
			seen[normalized] = true
			unique = append(unique, normalized)
		}
	}

	if len(unique) == 0 {
		return ""
	}

	if len(unique) > maxItems {
		unique = unique[:maxItems]
	}

	if len(unique) == 1 {
		return unique[0]
	}

	var result []string
	for _, item := range unique {
		result = append(result, "- "+item)
	}
	return strings.Join(result, "\n")
}

// truncateString truncates a string to maxLen characters.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// containsString checks if a slice contains a string.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// extractActionItems extracts action-oriented items from text using verb patterns.
func extractActionItems(text string) string {
	// 코드 블록 제거
	text = removeCodeBlocks(text)

	lines := strings.Split(text, "\n")
	var actions []string

	// 동사 패턴 (영어 + 한국어)
	actionPatterns := []string{
		// 영어 완료형
		"implemented", "added", "created", "fixed", "updated", "modified",
		"configured", "installed", "set up", "removed", "deleted",
		"refactored", "optimized", "enhanced", "improved", "resolved",
		"completed", "finished", "built", "deployed", "tested",
		// 한국어 완료형
		"구현", "추가", "생성", "수정", "삭제", "설정", "설치",
		"리팩토링", "최적화", "개선", "해결", "완료", "빌드", "배포", "테스트",
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 리스트 아이템 처리
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "• ") {
			line = strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* "), "• ")
		}

		lineLower := strings.ToLower(line)
		for _, pattern := range actionPatterns {
			if strings.Contains(lineLower, pattern) {
				// 너무 긴 라인은 자르기
				if len(line) > 150 {
					line = line[:150] + "..."
				}
				actions = append(actions, line)
				break
			}
		}
	}

	// 중복 제거 및 최대 10개로 제한
	return formatItems(actions, 10)
}

// parseMarkdownStructure parses markdown text and extracts structured sections.
// Enhanced version with better pattern matching for unstructured text.
func parseMarkdownStructure(text string) StructuredSummary {
	summary := StructuredSummary{}

	// 코드 블록 제거
	text = removeCodeBlocks(text)

	// 헤딩별 섹션 추출
	sections := extractSections(text)

	for heading, content := range sections {
		headingLower := strings.ToLower(heading)
		contentClean := strings.TrimSpace(content)
		if contentClean == "" {
			continue
		}

		switch {
		case containsAny(headingLower, "request", "요청", "task"):
			summary.Request = contentClean
		case containsAny(headingLower, "investigated", "분석", "조사", "확인", "탐색", "analysis", "review"):
			summary.Investigated = contentClean
		case containsAny(headingLower, "learned", "학습", "발견", "insight", "배움", "discovery"):
			summary.Learned = contentClean
		case containsAny(headingLower, "completed", "완료", "구현", "수정", "done", "result", "summary", "changes"):
			summary.Completed = contentClean
		case containsAny(headingLower, "next", "다음", "todo", "후속", "계획", "follow"):
			summary.NextSteps = contentClean
		}
	}

	// 헤딩이 없거나 completed가 비어있으면 리스트 아이템에서 추출
	listItems := extractListItems(text)
	if len(listItems) > 0 {
		var completedItems []string
		var investigatedItems []string
		var learnedItems []string

		for _, item := range listItems {
			itemLower := strings.ToLower(item)
			// 완료 항목 패턴
			if containsAny(itemLower, "완료", "구현", "수정", "추가", "fix", "add", "implement", "update",
				"create", "modify", "change", "set", "configure", "install", "remove", "delete",
				"refactor", "optimize", "enhance", "improve", "resolve", "build", "deploy") {
				completedItems = append(completedItems, item)
			}
			// 분석/조사 패턴
			if containsAny(itemLower, "분석", "조사", "확인", "검토", "analyze", "review", "check", "examine", "investigate", "found", "discovered") {
				investigatedItems = append(investigatedItems, item)
			}
			// 학습/발견 패턴
			if containsAny(itemLower, "학습", "발견", "배움", "알게", "learn", "discover", "realize", "understand", "note") {
				learnedItems = append(learnedItems, item)
			}
		}

		if summary.Completed == "" && len(completedItems) > 0 {
			summary.Completed = formatItems(completedItems, 10)
		}
		if summary.Investigated == "" && len(investigatedItems) > 0 {
			summary.Investigated = formatItems(investigatedItems, 5)
		}
		if summary.Learned == "" && len(learnedItems) > 0 {
			summary.Learned = formatItems(learnedItems, 5)
		}
	}

	// 그래도 completed가 비어있으면 첫 몇 줄의 의미있는 문장 추출
	if summary.Completed == "" {
		summary.Completed = extractFirstMeaningfulSentences(text, 3)
	}

	return summary
}

// extractFirstMeaningfulSentences extracts the first N meaningful sentences from text.
func extractFirstMeaningfulSentences(text string, maxSentences int) string {
	lines := strings.Split(text, "\n")
	var sentences []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 헤딩 건너뛰기
		if strings.HasPrefix(line, "#") {
			continue
		}
		// 리스트 마커 제거
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "• ") {
			line = strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* "), "• ")
		}
		// 너무 짧은 라인 건너뛰기 (10자 미만)
		if len(line) < 10 {
			continue
		}
		// 의미없는 라인 건너뛰기
		lineLower := strings.ToLower(line)
		if containsAny(lineLower, "here is", "here's", "let me", "i will", "i'll", "i've", "이제", "그럼", "아래") {
			continue
		}

		sentences = append(sentences, line)
		if len(sentences) >= maxSentences {
			break
		}
	}

	if len(sentences) == 0 {
		return ""
	}

	return formatItems(sentences, maxSentences)
}

// removeCodeBlocks removes ```...``` code blocks from text.
func removeCodeBlocks(text string) string {
	result := text
	for {
		start := strings.Index(result, "```")
		if start == -1 {
			break
		}
		end := strings.Index(result[start+3:], "```")
		if end == -1 {
			// 닫히지 않은 코드 블록 - 끝까지 제거
			result = result[:start]
			break
		}
		result = result[:start] + result[start+3+end+3:]
	}
	return result
}

// extractSections extracts ## heading sections from markdown.
func extractSections(text string) map[string]string {
	sections := make(map[string]string)
	lines := strings.Split(text, "\n")

	var currentHeading string
	var currentContent []string

	for _, line := range lines {
		if strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "# ") {
			// 이전 섹션 저장
			if currentHeading != "" {
				sections[currentHeading] = strings.TrimSpace(strings.Join(currentContent, "\n"))
			}
			// 새 섹션 시작
			currentHeading = strings.TrimPrefix(strings.TrimPrefix(line, "## "), "# ")
			currentContent = nil
		} else if currentHeading != "" {
			currentContent = append(currentContent, line)
		}
	}

	// 마지막 섹션 저장
	if currentHeading != "" {
		sections[currentHeading] = strings.TrimSpace(strings.Join(currentContent, "\n"))
	}

	return sections
}

// extractListItems extracts - or * list items from text.
func extractListItems(text string) []string {
	var items []string
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			item := strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* ")
			if item != "" {
				items = append(items, item)
			}
		}
	}
	return items
}

// containsAny checks if text contains any of the keywords.
func containsAny(text string, keywords ...string) bool {
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

// extractFilePath extracts file path from content like "수정: /path/to/file" or "읽기: /path/to/file".
func extractFilePath(content string) string {
	// "수정: ", "읽기: ", "테스트: " 등의 접두사 제거
	prefixes := []string{"수정: ", "읽기: ", "테스트: ", "Implemented: ", "Fixed: "}
	for _, prefix := range prefixes {
		if strings.HasPrefix(content, prefix) {
			return strings.TrimPrefix(content, prefix)
		}
	}
	// 경로로 보이면 그대로 반환
	if strings.HasPrefix(content, "/") || strings.Contains(content, "/") {
		return content
	}
	return ""
}

// strPtr returns a pointer to the string, or nil if empty.
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// generateLLMSummary generates a summary using Anthropic Claude API.
func generateLLMSummary(ctx context.Context, lastMessage string, userPrompts []string, apiKey string) (*StructuredSummary, error) {
	// Build the prompt (claude-mem style)
	userRequest := ""
	if len(userPrompts) > 0 {
		userRequest = userPrompts[0]
		if len(userRequest) > 500 {
			userRequest = userRequest[:500] + "..."
		}
	}

	prompt := `PROGRESS SUMMARY CHECKPOINT
===========================
Write progress notes of what was done, what was learned, and what's next.
This is a checkpoint to capture progress so far.

User's Request:
` + userRequest + `

Claude's Response:
` + lastMessage + `

Respond in this XML format ONLY:
<summary>
  <request>What the user originally requested</request>
  <investigated>What was explored or analyzed</investigated>
  <learned>Key learnings or discoveries</learned>
  <completed>What was actually done/implemented</completed>
  <next_steps>Current trajectory of work</next_steps>
</summary>

IMPORTANT: Output ONLY the XML summary, nothing else.`

	// Call Anthropic API
	requestBody := map[string]interface{}{
		"model":      "claude-sonnet-4-20250514",
		"max_tokens": 1000,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Content) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	// Parse XML response
	return parseSummaryXML(result.Content[0].Text)
}

// parseSummaryXML parses the XML summary response from LLM.
func parseSummaryXML(text string) (*StructuredSummary, error) {
	summary := &StructuredSummary{}

	// Extract each field from XML
	extractXMLField := func(xml, tag string) string {
		start := strings.Index(xml, "<"+tag+">")
		end := strings.Index(xml, "</"+tag+">")
		if start == -1 || end == -1 || start >= end {
			return ""
		}
		return strings.TrimSpace(xml[start+len(tag)+2 : end])
	}

	summary.Request = extractXMLField(text, "request")
	summary.Investigated = extractXMLField(text, "investigated")
	summary.Learned = extractXMLField(text, "learned")
	summary.Completed = extractXMLField(text, "completed")
	summary.NextSteps = extractXMLField(text, "next_steps")

	// Check if we got anything useful
	if summary.Request == "" && summary.Completed == "" {
		return nil, fmt.Errorf("no useful content in response")
	}

	return summary, nil
}

// formatSummaryAsMarkdown converts StructuredSummary to markdown string for legacy compatibility.
func formatSummaryAsMarkdown(s StructuredSummary) string {
	var parts []string

	if s.Request != "" {
		parts = append(parts, "## Request\n"+s.Request)
	}
	if s.Investigated != "" {
		parts = append(parts, "## Investigated\n"+s.Investigated)
	}
	if s.Completed != "" {
		parts = append(parts, "## Completed\n"+s.Completed)
	}
	if s.Learned != "" {
		parts = append(parts, "## Learned\n"+s.Learned)
	}
	if s.NextSteps != "" {
		parts = append(parts, "## Next Steps\n"+s.NextSteps)
	}

	if len(s.FilesEdited) > 0 {
		filesStr := "## Files Edited\n"
		for _, f := range s.FilesEdited {
			filesStr += "- " + f + "\n"
		}
		parts = append(parts, filesStr)
	}

	if len(parts) == 0 {
		return "No significant work recorded in this session."
	}

	return strings.Join(parts, "\n\n")
}

// handleGetUserPrompts handles user prompt list retrieval.
func (s *Server) handleGetUserPrompts(c *gin.Context) {
	ctx := c.Request.Context()

	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 100
	}

	sessionID := c.Query("session_id")

	prompts, err := s.db.GetUserPrompts(ctx, sessionID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, prompts)
}

// handleCreateUserPrompt handles user prompt creation.
func (s *Server) handleCreateUserPrompt(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.CreateUserPromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	prompt := &models.UserPrompt{
		SessionID:      req.SessionID,
		PromptNumber:   req.PromptNumber,
		PromptText:     req.PromptText,
		CreatedAtEpoch: time.Now().Unix(),
	}

	if err := s.db.CreateUserPrompt(ctx, prompt); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, prompt)
}

// handleUpdateLatestPromptResponse updates the response for the latest prompt in a session.
func (s *Server) handleUpdateLatestPromptResponse(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		SessionID string `json:"session_id" binding:"required"`
		Response  string `json:"response" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Truncate response to 100KB
	response := req.Response
	if len(response) > 100000 {
		response = response[:100000] + "\n...(truncated)"
	}

	if err := s.db.UpdateLatestPromptResponse(ctx, req.SessionID, response); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// handleSearch handles FTS5 full-text search across multiple types.
func (s *Server) handleSearch(c *gin.Context) {
	ctx := c.Request.Context()

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Query parameter 'q' is required",
		})
		return
	}

	// Parse types parameter (comma-separated)
	typesStr := c.DefaultQuery("types", "observation,prompt")
	types := strings.Split(typesStr, ",")
	for i, t := range types {
		types[i] = strings.TrimSpace(t)
	}

	// Parse limit
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	results, err := s.db.SearchFTS(ctx, query, types, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SearchResponse{
		Results: results,
		Query:   query,
		Total:   len(results),
	})
}

// buildContextMarkdown builds a markdown representation of the context.
func buildContextMarkdown(session *models.Session, observations []models.Observation, plan *models.Plan, team []models.TeamContext) string {
	var md string

	md += "# Do Worker Context\n\n"

	// Session info
	if session != nil {
		md += "## Last Session\n"
		md += "- ID: " + session.ID + "\n"
		md += "- Started: " + session.StartedAt.Format(time.RFC3339) + "\n"
		if session.EndedAt != nil {
			md += "- Ended: " + session.EndedAt.Format(time.RFC3339) + "\n"
		}
		if session.Summary != "" {
			md += "- Summary: " + session.Summary + "\n"
		}
		md += "\n"
	}

	// Active plan
	if plan != nil {
		md += "## Active Plan\n"
		md += "**" + plan.Title + "**\n\n"
		md += plan.Content + "\n\n"
	}

	// Recent observations
	if len(observations) > 0 {
		md += "## Recent Observations\n"
		for _, obs := range observations {
			importance := ""
			if obs.Importance >= 4 {
				importance = " [HIGH]"
			}
			md += "- [" + obs.Type + "]" + importance + " " + obs.Content
			if obs.AgentName != "" {
				md += " (by " + obs.AgentName + ")"
			}
			md += "\n"
		}
		md += "\n"
	}

	// Team context
	if len(team) > 0 {
		md += "## Team Activity\n"
		for _, t := range team {
			md += "- **" + t.UserName + "**: " + t.Summary
			if t.ActivePlan != "" {
				md += " [Working on: " + t.ActivePlan + "]"
			}
			md += "\n"
		}
	}

	return md
}
