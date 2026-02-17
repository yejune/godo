package rank

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Rank API client constants.
const (
	DefaultBaseURL = "https://rank.mo.ai.kr"
	APIVersion     = "v1"
	UserAgent      = "godo/1.0"
	MaxTokens      = 100_000_000
)

// --- Error Types ---

// ClientError represents a general client-side error.
type ClientError struct {
	Message string
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("rank client error: %s", e.Message)
}

// AuthError represents an authentication failure.
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("rank auth error: %s", e.Message)
}

// APIError represents an API response error.
type APIError struct {
	Message    string
	StatusCode int
	Details    map[string]any
}

func (e *APIError) Error() string {
	return fmt.Sprintf("rank API error (status %d): %s", e.StatusCode, e.Message)
}

// --- Data Models ---

// SessionSubmission holds session data for submission to the Rank API.
type SessionSubmission struct {
	SessionHash         string `json:"session_hash"`
	EndedAt             string `json:"ended_at"`
	InputTokens         int64  `json:"input_tokens"`
	OutputTokens        int64  `json:"output_tokens"`
	CacheCreationTokens int64  `json:"cache_creation_tokens"`
	CacheReadTokens     int64  `json:"cache_read_tokens"`
	ModelName           string `json:"model_name,omitempty"`
	AnonymousProjectID  string `json:"anonymous_project_id,omitempty"`
	StartedAt           string `json:"started_at,omitempty"`
	DurationSeconds     int    `json:"duration_seconds,omitempty"`
	TurnCount           int    `json:"turn_count,omitempty"`
}

// UserInfo represents the full ranking information for a user.
type UserInfo struct {
	Username      string    `json:"username"`
	TotalTokens   int64     `json:"total_tokens"`
	TotalSessions int       `json:"total_sessions"`
	InputTokens   int64     `json:"input_tokens"`
	OutputTokens  int64     `json:"output_tokens"`
	Daily         *Position `json:"daily,omitempty"`
	Weekly        *Position `json:"weekly,omitempty"`
	Monthly       *Position `json:"monthly,omitempty"`
	AllTime       *Position `json:"all_time,omitempty"`
	LastUpdated   string    `json:"last_updated"`
}

// Position holds ranking position details for a time period.
type Position struct {
	Position          int     `json:"position"`
	CompositeScore    float64 `json:"composite_score"`
	TotalParticipants int     `json:"total_participants"`
}

// --- Client ---

// Client communicates with the MoAI Rank API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Client.
func NewClient(apiKey string) *Client {
	baseURL := DefaultBaseURL
	if envURL := os.Getenv("DO_RANK_API_URL"); envURL != "" {
		baseURL = envURL
	}

	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ComputeSignature calculates the HMAC-SHA256 signature for a request.
func ComputeSignature(apiKey, timestamp, body string) string {
	message := timestamp + ":" + body
	mac := hmac.New(sha256.New, []byte(apiKey))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// doAuthRequest performs an authenticated HTTP request to the Rank API.
func (c *Client) doAuthRequest(ctx context.Context, method, path string, body []byte) (*http.Response, error) {
	if c.apiKey == "" {
		return nil, &AuthError{Message: "API key not configured"}
	}

	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, &ClientError{Message: fmt.Sprintf("create request: %v", err)}
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/json")

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	bodyStr := ""
	if body != nil {
		bodyStr = string(body)
	}
	signature := ComputeSignature(c.apiKey, timestamp, bodyStr)

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Signature", signature)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &ClientError{Message: fmt.Sprintf("request failed: %v", err)}
	}

	return resp, nil
}

// SubmitSession submits a single session metric to the Rank API.
func (c *Client) SubmitSession(ctx context.Context, session *SessionSubmission) error {
	session.InputTokens = ClampTokens(session.InputTokens)
	session.OutputTokens = ClampTokens(session.OutputTokens)
	session.CacheCreationTokens = ClampTokens(session.CacheCreationTokens)
	session.CacheReadTokens = ClampTokens(session.CacheReadTokens)

	body, err := json.Marshal(session)
	if err != nil {
		return &ClientError{Message: fmt.Sprintf("marshal session: %v", err)}
	}

	path := "/api/" + APIVersion + "/sessions"
	resp, err := c.doAuthRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return &AuthError{Message: "authentication failed"}
	}

	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}

	return nil
}

// GetUserRank returns the current user's ranking information.
func (c *Client) GetUserRank(ctx context.Context) (*UserInfo, error) {
	path := "/api/" + APIVersion + "/rank"
	resp, err := c.doAuthRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, &AuthError{Message: "authentication failed"}
	}

	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ClientError{Message: fmt.Sprintf("read response: %v", err)}
	}

	var userRank UserInfo
	if err := json.Unmarshal(respBody, &userRank); err != nil {
		return nil, &ClientError{Message: fmt.Sprintf("parse rank response: %v", err)}
	}

	return &userRank, nil
}

// ComputeSessionHash generates a unique SHA-256 hash for a session.
func ComputeSessionHash(endedAt string, inputTokens, outputTokens int64) (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}

	data := fmt.Sprintf("%s:%d:%d:%x", endedAt, inputTokens, outputTokens, nonce)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:]), nil
}

// ClampTokens clamps a token value to MaxTokens (100M).
func ClampTokens(value int64) int64 {
	if value > MaxTokens {
		return MaxTokens
	}
	return value
}

// parseAPIError extracts error details from an HTTP response.
func parseAPIError(resp *http.Response) error {
	respBody, _ := io.ReadAll(resp.Body)
	apiErr := &APIError{
		Message:    fmt.Sprintf("API returned status %d", resp.StatusCode),
		StatusCode: resp.StatusCode,
	}

	var details map[string]any
	if json.Unmarshal(respBody, &details) == nil {
		apiErr.Details = details
		if msg, ok := details["message"].(string); ok {
			apiErr.Message = msg
		}
	}

	return apiErr
}
