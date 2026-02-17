package hook

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestHandlePreTool_AllowsNormalFileWrite(t *testing.T) {
	input := &Input{
		ToolName:  "Write",
		ToolInput: json.RawMessage(`{"file_path": "/home/user/project/main.go"}`),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	if output.HookSpecificOutput.PermissionDecision != DecisionAllow {
		t.Errorf("expected decision %q, got %q", DecisionAllow, output.HookSpecificOutput.PermissionDecision)
	}
}

func TestHandlePreTool_DeniesSSHKeyWrite(t *testing.T) {
	input := &Input{
		ToolName:  "Write",
		ToolInput: json.RawMessage(`{"file_path": "/home/user/.ssh/id_rsa"}`),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	if output.HookSpecificOutput.PermissionDecision != DecisionDeny {
		t.Errorf("expected decision %q, got %q", DecisionDeny, output.HookSpecificOutput.PermissionDecision)
	}
	if output.HookSpecificOutput.PermissionDecisionReason == "" {
		t.Error("expected non-empty deny reason")
	}
}

func TestHandlePreTool_DeniesDangerousBashCommand(t *testing.T) {
	input := &Input{
		ToolName:  "Bash",
		ToolInput: json.RawMessage(`{"command": "rm -rf /"}`),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	if output.HookSpecificOutput.PermissionDecision != DecisionDeny {
		t.Errorf("expected decision %q, got %q", DecisionDeny, output.HookSpecificOutput.PermissionDecision)
	}
}

func TestHandlePreTool_AsksForLockFileWrite(t *testing.T) {
	input := &Input{
		ToolName:  "Write",
		ToolInput: json.RawMessage(`{"file_path": "/project/package-lock.json"}`),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	if output.HookSpecificOutput.PermissionDecision != DecisionAsk {
		t.Errorf("expected decision %q, got %q", DecisionAsk, output.HookSpecificOutput.PermissionDecision)
	}
}

func TestHandlePreTool_AllowsUnknownTool(t *testing.T) {
	input := &Input{
		ToolName:  "WebSearch",
		ToolInput: json.RawMessage(`{"query": "test"}`),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	if output.HookSpecificOutput.PermissionDecision != DecisionAllow {
		t.Errorf("expected decision %q, got %q", DecisionAllow, output.HookSpecificOutput.PermissionDecision)
	}
}

func TestHandlePreTool_ReadDoesNotTriggerAskPatterns(t *testing.T) {
	input := &Input{
		ToolName:  "Read",
		ToolInput: json.RawMessage(`{"file_path": "/project/package-lock.json"}`),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	// Read should allow lock files without asking
	if output.HookSpecificOutput.PermissionDecision != DecisionAllow {
		t.Errorf("expected decision %q for Read on lock file, got %q", DecisionAllow, output.HookSpecificOutput.PermissionDecision)
	}
}

func TestHandlePreTool_DeniesSensitiveContentInBash(t *testing.T) {
	// Build a fake sensitive token dynamically to avoid triggering
	// the project's own PreToolUse hook on this source file.
	fakeToken := "sk-" + strings.Repeat("a", 40)
	cmdJSON := `{"command": "echo ` + fakeToken + `"}`

	input := &Input{
		ToolName:  "Bash",
		ToolInput: json.RawMessage(cmdJSON),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	if output.HookSpecificOutput.PermissionDecision != DecisionDeny {
		t.Errorf("expected decision %q, got %q", DecisionDeny, output.HookSpecificOutput.PermissionDecision)
	}
}

func TestExtractFilePath_ValidInput(t *testing.T) {
	raw := json.RawMessage(`{"file_path": "/home/user/test.go"}`)
	result := extractFilePath(raw)
	if result != "/home/user/test.go" {
		t.Errorf("expected %q, got %q", "/home/user/test.go", result)
	}
}

func TestExtractFilePath_FallsBackToPath(t *testing.T) {
	raw := json.RawMessage(`{"path": "/home/user/dir"}`)
	result := extractFilePath(raw)
	if result != "/home/user/dir" {
		t.Errorf("expected %q, got %q", "/home/user/dir", result)
	}
}

func TestExtractFilePath_EmptyInput(t *testing.T) {
	result := extractFilePath(nil)
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestExtractFilePath_InvalidJSON(t *testing.T) {
	raw := json.RawMessage(`{invalid}`)
	result := extractFilePath(raw)
	if result != "" {
		t.Errorf("expected empty string for invalid JSON, got %q", result)
	}
}

func TestExtractCommand_ValidInput(t *testing.T) {
	raw := json.RawMessage(`{"command": "  ls -la  "}`)
	result := extractCommand(raw)
	if result != "ls -la" {
		t.Errorf("expected %q, got %q", "ls -la", result)
	}
}

func TestExtractCommand_EmptyInput(t *testing.T) {
	result := extractCommand(nil)
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestExtractCommand_InvalidJSON(t *testing.T) {
	raw := json.RawMessage(`not json`)
	result := extractCommand(raw)
	if result != "" {
		t.Errorf("expected empty string for invalid JSON, got %q", result)
	}
}

func TestExtractCommand_NoCommandField(t *testing.T) {
	raw := json.RawMessage(`{"other": "value"}`)
	result := extractCommand(raw)
	if result != "" {
		t.Errorf("expected empty string when no command field, got %q", result)
	}
}

func TestHandlePreTool_AsksForDangerousBashCommand(t *testing.T) {
	input := &Input{
		ToolName:  "Bash",
		ToolInput: json.RawMessage(`{"command": "git push --force"}`),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	if output.HookSpecificOutput.PermissionDecision != DecisionAsk {
		t.Errorf("expected decision %q for git push --force, got %q", DecisionAsk, output.HookSpecificOutput.PermissionDecision)
	}
}

func TestHandlePreTool_AllowsSafeBashCommand(t *testing.T) {
	input := &Input{
		ToolName:  "Bash",
		ToolInput: json.RawMessage(`{"command": "go test ./..."}`),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	if output.HookSpecificOutput.PermissionDecision != DecisionAllow {
		t.Errorf("expected decision %q for safe command, got %q", DecisionAllow, output.HookSpecificOutput.PermissionDecision)
	}
}

func TestHandlePreTool_EmptyBashCommand(t *testing.T) {
	input := &Input{
		ToolName:  "Bash",
		ToolInput: json.RawMessage(`{"command": ""}`),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	if output.HookSpecificOutput.PermissionDecision != DecisionAllow {
		t.Errorf("expected decision %q for empty command, got %q", DecisionAllow, output.HookSpecificOutput.PermissionDecision)
	}
}

func TestHandlePreTool_EditDeniesSSHKey(t *testing.T) {
	input := &Input{
		ToolName:  "Edit",
		ToolInput: json.RawMessage(`{"file_path": "/home/user/.ssh/id_rsa"}`),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	if output.HookSpecificOutput.PermissionDecision != DecisionDeny {
		t.Errorf("expected decision %q, got %q", DecisionDeny, output.HookSpecificOutput.PermissionDecision)
	}
}

func TestHandlePreTool_EditAsksForLockFile(t *testing.T) {
	input := &Input{
		ToolName:  "Edit",
		ToolInput: json.RawMessage(`{"file_path": "/project/package-lock.json"}`),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	if output.HookSpecificOutput.PermissionDecision != DecisionAsk {
		t.Errorf("expected decision %q for Edit on lock file, got %q", DecisionAsk, output.HookSpecificOutput.PermissionDecision)
	}
}

func TestHandlePreTool_GlobAllowsNormalPath(t *testing.T) {
	input := &Input{
		ToolName:  "Glob",
		ToolInput: json.RawMessage(`{"file_path": "/home/user/project/src"}`),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	if output.HookSpecificOutput.PermissionDecision != DecisionAllow {
		t.Errorf("expected decision %q, got %q", DecisionAllow, output.HookSpecificOutput.PermissionDecision)
	}
}

func TestHandlePreTool_WriteEmptyInput(t *testing.T) {
	input := &Input{
		ToolName:  "Write",
		ToolInput: json.RawMessage(`{}`),
	}
	output := HandlePreTool(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput to be non-nil")
	}
	// Empty file path should allow
	if output.HookSpecificOutput.PermissionDecision != DecisionAllow {
		t.Errorf("expected decision %q for empty input, got %q", DecisionAllow, output.HookSpecificOutput.PermissionDecision)
	}
}
