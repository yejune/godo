package hook

import (
	"encoding/json"
	"os"
)

// JobState represents the workflow state tracked in state.json.
type JobState struct {
	JobID               string                `json:"job_id"`
	CreatedAt           string                `json:"created_at"`
	WorkflowType        string                `json:"workflow_type"`                   // "simple" or "complex"
	Phases              map[string]PhaseState  `json:"phases"`
	Agents              map[string]AgentState  `json:"agents"`
	AutoResolveAttempts map[string]bool        `json:"auto_resolve_attempts,omitempty"` // tracks auto-resolve attempts (max 1 per dep)
}

// PhaseState tracks the status of a workflow phase.
type PhaseState struct {
	Status      string `json:"status"`                 // pending, in_progress, complete
	StartedAt   string `json:"started_at,omitempty"`
	CompletedAt string `json:"completed_at,omitempty"`
}

// AgentState tracks the status of an agent's work.
type AgentState struct {
	Status      string   `json:"status"`                 // pending, in_progress, complete, failed, blocked
	Checklist   string   `json:"checklist,omitempty"`
	CompletedAt string   `json:"completed_at,omitempty"`
	BlockedBy   []string `json:"blocked_by,omitempty"`
}

// LoadJobState reads and parses a state.json file.
func LoadJobState(path string) (*JobState, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var state JobState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// SaveJobState writes a JobState to a JSON file with indentation.
func SaveJobState(path string, state *JobState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0644)
}
