package rank

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// TranscriptUsage represents token usage extracted from a Claude Code transcript.
type TranscriptUsage struct {
	InputTokens         int64  `json:"input_tokens"`
	OutputTokens        int64  `json:"output_tokens"`
	CacheCreationTokens int64  `json:"cache_creation_tokens"`
	CacheReadTokens     int64  `json:"cache_read_tokens"`
	ModelName           string `json:"model_name"`
	StartedAt           string `json:"started_at,omitempty"`
	EndedAt             string `json:"ended_at,omitempty"`
	DurationSeconds     int    `json:"duration_seconds,omitempty"`
	TurnCount           int    `json:"turn_count,omitempty"`
}

// transcriptMessage represents a single line in the JSONL transcript file.
type transcriptMessage struct {
	Timestamp string        `json:"timestamp"`
	Type      string        `json:"type"`
	Message   transcriptMsg `json:"message"`
	Model     string        `json:"model"`
}

// transcriptMsg represents the message content with usage data.
type transcriptMsg struct {
	Usage *transcriptUsage `json:"usage"`
	Model string           `json:"model"`
}

// transcriptUsage represents token usage information from a transcript line.
type transcriptUsage struct {
	InputTokens              int64 `json:"input_tokens"`
	OutputTokens             int64 `json:"output_tokens"`
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
}

// ParseTranscript parses a Claude Code transcript JSONL file and extracts token usage.
func ParseTranscript(path string) (*TranscriptUsage, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open transcript: %w", err)
	}
	defer func() { _ = file.Close() }()

	usage := &TranscriptUsage{}
	var firstTimestamp, lastTimestamp string

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var msg transcriptMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}

		if msg.Timestamp != "" {
			if firstTimestamp == "" {
				firstTimestamp = msg.Timestamp
			}
			lastTimestamp = msg.Timestamp
		}

		if msg.Type == "user" {
			usage.TurnCount++
		}

		if usage.ModelName == "" {
			if msg.Model != "" {
				usage.ModelName = msg.Model
			} else if msg.Message.Model != "" {
				usage.ModelName = msg.Message.Model
			}
		}

		if msg.Message.Usage != nil {
			usage.InputTokens += msg.Message.Usage.InputTokens
			usage.OutputTokens += msg.Message.Usage.OutputTokens
			usage.CacheCreationTokens += msg.Message.Usage.CacheCreationInputTokens
			usage.CacheReadTokens += msg.Message.Usage.CacheReadInputTokens
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan transcript: %w", err)
	}

	usage.StartedAt = firstTimestamp
	usage.EndedAt = lastTimestamp

	if firstTimestamp != "" && lastTimestamp != "" {
		start, errS := time.Parse(time.RFC3339Nano, firstTimestamp)
		end, errE := time.Parse(time.RFC3339Nano, lastTimestamp)
		if errS == nil && errE == nil {
			usage.DurationSeconds = int(end.Sub(start).Seconds())
		}
	}

	return usage, nil
}

// ClaudeCodeDir returns the Claude Code CLI configuration directory (~/.claude/).
func ClaudeCodeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".claude")
}

// ClaudeDesktopConfigDir returns the Claude Desktop (Electron app) configuration directory.
func ClaudeDesktopConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "Claude")
	case "linux":
		return filepath.Join(homeDir, ".config", "Claude")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return ""
		}
		return filepath.Join(appData, "Claude")
	default:
		return ""
	}
}

// GlobJSONL collects .jsonl files matching the given pattern, ignoring glob errors.
func GlobJSONL(pattern string) []string {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil
	}
	return matches
}

// IsValidSessionID validates a session ID to prevent path traversal attacks.
func IsValidSessionID(sessionID string) bool {
	if sessionID == "" || len(sessionID) > 128 {
		return false
	}
	for _, c := range sessionID {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	return true
}

// FindTranscriptForSession finds the transcript file for a specific session ID.
func FindTranscriptForSession(sessionID string) string {
	if !IsValidSessionID(sessionID) {
		return ""
	}

	if codeDir := ClaudeCodeDir(); codeDir != "" {
		pattern := filepath.Join(codeDir, "projects", "*", sessionID+"*.jsonl")
		if matches := GlobJSONL(pattern); len(matches) > 0 {
			return matches[0]
		}
	}

	if codeDir := ClaudeCodeDir(); codeDir != "" {
		pattern := filepath.Join(codeDir, "transcripts", sessionID+"*.jsonl")
		if matches := GlobJSONL(pattern); len(matches) > 0 {
			return matches[0]
		}
	}

	if desktopDir := ClaudeDesktopConfigDir(); desktopDir != "" {
		pattern := filepath.Join(desktopDir, "*", "transcripts", sessionID+"*.jsonl")
		if matches := GlobJSONL(pattern); len(matches) > 0 {
			return matches[0]
		}
	}

	return ""
}

// FindAllTranscripts returns all Claude Code transcript JSONL files.
func FindAllTranscripts() []string {
	seen := make(map[string]struct{})
	var results []string

	addUnique := func(paths []string) {
		for _, p := range paths {
			if _, exists := seen[p]; !exists {
				seen[p] = struct{}{}
				results = append(results, p)
			}
		}
	}

	if codeDir := ClaudeCodeDir(); codeDir != "" {
		addUnique(GlobJSONL(filepath.Join(codeDir, "projects", "*", "*.jsonl")))
	}

	if codeDir := ClaudeCodeDir(); codeDir != "" {
		addUnique(GlobJSONL(filepath.Join(codeDir, "transcripts", "*.jsonl")))
	}

	if desktopDir := ClaudeDesktopConfigDir(); desktopDir != "" {
		addUnique(GlobJSONL(filepath.Join(desktopDir, "*", "transcripts", "*.jsonl")))
	}

	return results
}

// AnonymizeProjectPath converts a project path to a SHA-256 hash for privacy.
func AnonymizeProjectPath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	hash := sha256.Sum256([]byte(absPath))
	return hex.EncodeToString(hash[:])
}
