package hook

import "encoding/json"

// EventType represents a Claude Code hook event type.
type EventType string

const (
	EventSessionStart EventType = "SessionStart"
	EventPreToolUse   EventType = "PreToolUse"
	EventPostToolUse  EventType = "PostToolUse"
	EventSessionEnd   EventType = "SessionEnd"
	EventStop         EventType = "Stop"
	EventSubagentStop EventType = "SubagentStop"
	EventPreCompact   EventType = "PreCompact"
)

// ValidEventTypes returns all valid event types.
func ValidEventTypes() []EventType {
	return []EventType{
		EventSessionStart,
		EventPreToolUse,
		EventPostToolUse,
		EventSessionEnd,
		EventStop,
		EventSubagentStop,
		EventPreCompact,
	}
}

// IsValidEventType checks if the given event type is valid.
func IsValidEventType(et EventType) bool {
	for _, v := range ValidEventTypes() {
		if v == et {
			return true
		}
	}
	return false
}

// Permission decision constants for PreToolUse hooks (Claude Code protocol).
const (
	DecisionAllow = "allow"
	DecisionDeny  = "deny"
	DecisionAsk   = "ask"
)

// Top-level decision constant for Stop, PostToolUse, etc. (Claude Code protocol).
const (
	DecisionBlock = "block"
)

// Input represents the JSON payload received from Claude Code via stdin.
// Fields follow the official Claude Code hooks protocol.
type Input struct {
	// Common fields (all events)
	SessionID      string `json:"session_id,omitempty"`
	TranscriptPath string `json:"transcript_path,omitempty"`
	CWD            string `json:"cwd,omitempty"`
	PermissionMode string `json:"permission_mode,omitempty"`
	HookEventName  string `json:"hook_event_name,omitempty"`

	// Tool-related fields (PreToolUse, PostToolUse)
	ToolName     string          `json:"tool_name,omitempty"`
	ToolInput    json.RawMessage `json:"tool_input,omitempty"`
	ToolOutput   json.RawMessage `json:"tool_output,omitempty"`
	ToolResponse json.RawMessage `json:"tool_response,omitempty"`
	ToolUseID    string          `json:"tool_use_id,omitempty"`

	// SessionStart fields
	Source    string `json:"source,omitempty"`
	Model     string `json:"model,omitempty"`
	AgentType string `json:"agent_type,omitempty"`

	// SessionEnd fields
	Reason string `json:"reason,omitempty"`

	// Stop/SubagentStop fields
	StopHookActive bool `json:"stop_hook_active,omitempty"`

	// SubagentStart/SubagentStop fields
	AgentID             string `json:"agent_id,omitempty"`
	AgentTranscriptPath string `json:"agent_transcript_path,omitempty"`

	// PreCompact fields
	Trigger            string `json:"trigger,omitempty"`
	CustomInstructions string `json:"custom_instructions,omitempty"`

	// PostToolUseFailure fields
	Error       string `json:"error,omitempty"`
	IsInterrupt bool   `json:"is_interrupt,omitempty"`

	// UserPromptSubmit fields
	Prompt string `json:"prompt,omitempty"`

	// Notification fields
	Message          string `json:"message,omitempty"`
	Title            string `json:"title,omitempty"`
	NotificationType string `json:"notification_type,omitempty"`

	// Legacy field (deprecated, use CWD instead)
	ProjectDir string `json:"project_dir,omitempty"`
}

// SpecificOutput represents the hookSpecificOutput field for PreToolUse/PostToolUse.
type SpecificOutput struct {
	HookEventName            string `json:"hookEventName,omitempty"`
	PermissionDecision       string `json:"permissionDecision,omitempty"`
	PermissionDecisionReason string `json:"permissionDecisionReason,omitempty"`
	AdditionalContext        string `json:"additionalContext,omitempty"`
}

// Output represents the JSON payload written to stdout for Claude Code.
type Output struct {
	// Universal fields (all events)
	Continue       bool   `json:"continue,omitempty"`
	StopReason     string `json:"stopReason,omitempty"`
	SystemMessage  string `json:"systemMessage,omitempty"`
	SuppressOutput bool   `json:"suppressOutput,omitempty"`

	// Top-level decision fields (Stop, SubagentStop, PostToolUse)
	Decision string `json:"decision,omitempty"`
	Reason   string `json:"reason,omitempty"`

	// For PreToolUse/PostToolUse: hook-specific output
	HookSpecificOutput *SpecificOutput `json:"hookSpecificOutput,omitempty"`
}

// NewAllowOutput creates an Output with permissionDecision "allow" for PreToolUse.
func NewAllowOutput() *Output {
	return &Output{
		HookSpecificOutput: &SpecificOutput{
			HookEventName:      "PreToolUse",
			PermissionDecision: DecisionAllow,
		},
	}
}

// NewDenyOutput creates an Output with permissionDecision "deny" for PreToolUse.
func NewDenyOutput(reason string) *Output {
	return &Output{
		HookSpecificOutput: &SpecificOutput{
			HookEventName:            "PreToolUse",
			PermissionDecision:       DecisionDeny,
			PermissionDecisionReason: reason,
		},
	}
}

// NewAskOutput creates an Output with permissionDecision "ask" for PreToolUse.
func NewAskOutput(reason string) *Output {
	return &Output{
		HookSpecificOutput: &SpecificOutput{
			HookEventName:            "PreToolUse",
			PermissionDecision:       DecisionAsk,
			PermissionDecisionReason: reason,
		},
	}
}

// NewAllowOutputWithWarning creates an Output with permissionDecision "allow"
// but includes a warning message in additionalContext for the orchestrator to see.
func NewAllowOutputWithWarning(warning string) *Output {
	return &Output{
		HookSpecificOutput: &SpecificOutput{
			HookEventName:      "PreToolUse",
			PermissionDecision: DecisionAllow,
			AdditionalContext:  warning,
		},
	}
}

// NewBlockOutput creates an Output with permissionDecision "deny" for PreToolUse.
// For Stop/PostToolUse, use NewStopBlockOutput instead.
func NewBlockOutput(reason string) *Output {
	return NewDenyOutput(reason)
}

// NewSuppressOutput creates an Output that suppresses output.
func NewSuppressOutput() *Output {
	return &Output{SuppressOutput: true}
}

// NewSessionOutput creates an Output for SessionStart/SessionEnd events.
func NewSessionOutput(continueSession bool, message string) *Output {
	return &Output{
		Continue:      continueSession,
		SystemMessage: message,
	}
}

// NewPostToolOutput creates an Output with additionalContext for PostToolUse.
func NewPostToolOutput(additionalContext string) *Output {
	return &Output{
		HookSpecificOutput: &SpecificOutput{
			HookEventName:    "PostToolUse",
			AdditionalContext: additionalContext,
		},
	}
}

// NewStopBlockOutput creates an Output that prevents Claude from stopping.
// Per Claude Code protocol, Stop hooks use top-level decision/reason.
func NewStopBlockOutput(reason string) *Output {
	return &Output{
		Decision: DecisionBlock,
		Reason:   reason,
	}
}

// NewPostToolBlockOutput creates an Output that blocks after tool execution.
// Per Claude Code protocol, PostToolUse uses top-level decision/reason.
func NewPostToolBlockOutput(reason string, additionalContext string) *Output {
	output := &Output{
		Decision: DecisionBlock,
		Reason:   reason,
	}
	if additionalContext != "" {
		output.HookSpecificOutput = &SpecificOutput{
			HookEventName:    "PostToolUse",
			AdditionalContext: additionalContext,
		}
	}
	return output
}
