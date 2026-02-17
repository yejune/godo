package hook

import (
	"encoding/json"
	"testing"
)

func Test_IsValidEventType_known(t *testing.T) {
	for _, et := range ValidEventTypes() {
		if !IsValidEventType(et) {
			t.Errorf("expected %q to be valid", et)
		}
	}
}

func Test_IsValidEventType_unknown(t *testing.T) {
	if IsValidEventType("Bogus") {
		t.Error("expected 'Bogus' to be invalid")
	}
}

func Test_ValidEventTypes_completeness(t *testing.T) {
	expected := []EventType{
		EventSessionStart, EventPreToolUse, EventPostToolUse,
		EventSessionEnd, EventStop, EventSubagentStop, EventPreCompact,
		EventUserPromptSubmit,
	}
	got := ValidEventTypes()
	if len(got) != len(expected) {
		t.Fatalf("expected %d event types, got %d", len(expected), len(got))
	}
	for i, e := range expected {
		if got[i] != e {
			t.Errorf("index %d: expected %q, got %q", i, e, got[i])
		}
	}
}

func Test_Input_JSON_roundtrip(t *testing.T) {
	input := Input{
		SessionID:  "sess-123",
		ToolName:   "Bash",
		CWD:        "/tmp",
		ToolInput:  json.RawMessage(`{"command":"ls"}`),
	}

	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Input
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.SessionID != "sess-123" {
		t.Errorf("SessionID: got %q, want %q", decoded.SessionID, "sess-123")
	}
	if decoded.ToolName != "Bash" {
		t.Errorf("ToolName: got %q, want %q", decoded.ToolName, "Bash")
	}
	if decoded.CWD != "/tmp" {
		t.Errorf("CWD: got %q, want %q", decoded.CWD, "/tmp")
	}
}

func Test_Output_JSON_roundtrip(t *testing.T) {
	output := NewDenyOutput("not allowed")
	data, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Output
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.HookSpecificOutput == nil {
		t.Fatal("HookSpecificOutput is nil")
	}
	if decoded.HookSpecificOutput.PermissionDecision != DecisionDeny {
		t.Errorf("decision: got %q, want %q", decoded.HookSpecificOutput.PermissionDecision, DecisionDeny)
	}
}

func Test_NewAllowOutput(t *testing.T) {
	out := NewAllowOutput()
	if out.HookSpecificOutput == nil {
		t.Fatal("HookSpecificOutput is nil")
	}
	if out.HookSpecificOutput.PermissionDecision != DecisionAllow {
		t.Errorf("got %q, want %q", out.HookSpecificOutput.PermissionDecision, DecisionAllow)
	}
	if out.HookSpecificOutput.HookEventName != "PreToolUse" {
		t.Errorf("got %q, want %q", out.HookSpecificOutput.HookEventName, "PreToolUse")
	}
}

func Test_NewDenyOutput(t *testing.T) {
	out := NewDenyOutput("dangerous")
	if out.HookSpecificOutput.PermissionDecision != DecisionDeny {
		t.Errorf("decision: got %q, want %q", out.HookSpecificOutput.PermissionDecision, DecisionDeny)
	}
	if out.HookSpecificOutput.PermissionDecisionReason != "dangerous" {
		t.Errorf("reason: got %q, want %q", out.HookSpecificOutput.PermissionDecisionReason, "dangerous")
	}
}

func Test_NewAskOutput(t *testing.T) {
	out := NewAskOutput("confirm?")
	if out.HookSpecificOutput.PermissionDecision != DecisionAsk {
		t.Errorf("decision: got %q, want %q", out.HookSpecificOutput.PermissionDecision, DecisionAsk)
	}
}

func Test_NewAllowOutputWithWarning(t *testing.T) {
	out := NewAllowOutputWithWarning("be careful")
	if out.HookSpecificOutput.PermissionDecision != DecisionAllow {
		t.Errorf("decision: got %q, want %q", out.HookSpecificOutput.PermissionDecision, DecisionAllow)
	}
	if out.HookSpecificOutput.AdditionalContext != "be careful" {
		t.Errorf("context: got %q, want %q", out.HookSpecificOutput.AdditionalContext, "be careful")
	}
}

func Test_NewSessionOutput(t *testing.T) {
	out := NewSessionOutput(true, "hello")
	if !out.Continue {
		t.Error("expected Continue=true")
	}
	if out.SystemMessage != "hello" {
		t.Errorf("SystemMessage: got %q, want %q", out.SystemMessage, "hello")
	}
}

func Test_NewSuppressOutput(t *testing.T) {
	out := NewSuppressOutput()
	if !out.SuppressOutput {
		t.Error("expected SuppressOutput=true")
	}
}

func Test_NewPostToolOutput(t *testing.T) {
	out := NewPostToolOutput("extra context")
	if out.HookSpecificOutput.HookEventName != "PostToolUse" {
		t.Errorf("event: got %q, want %q", out.HookSpecificOutput.HookEventName, "PostToolUse")
	}
	if out.HookSpecificOutput.AdditionalContext != "extra context" {
		t.Errorf("context: got %q, want %q", out.HookSpecificOutput.AdditionalContext, "extra context")
	}
}

func Test_NewStopBlockOutput(t *testing.T) {
	out := NewStopBlockOutput("not yet")
	if out.Decision != DecisionBlock {
		t.Errorf("decision: got %q, want %q", out.Decision, DecisionBlock)
	}
	if out.Reason != "not yet" {
		t.Errorf("reason: got %q, want %q", out.Reason, "not yet")
	}
}

func Test_NewPostToolBlockOutput_with_context(t *testing.T) {
	out := NewPostToolBlockOutput("blocked", "detail")
	if out.Decision != DecisionBlock {
		t.Errorf("decision: got %q, want %q", out.Decision, DecisionBlock)
	}
	if out.HookSpecificOutput == nil {
		t.Fatal("HookSpecificOutput is nil")
	}
	if out.HookSpecificOutput.AdditionalContext != "detail" {
		t.Errorf("context: got %q, want %q", out.HookSpecificOutput.AdditionalContext, "detail")
	}
}

func Test_NewPostToolBlockOutput_empty_context(t *testing.T) {
	out := NewPostToolBlockOutput("blocked", "")
	if out.HookSpecificOutput != nil {
		t.Error("expected HookSpecificOutput to be nil when context is empty")
	}
}
