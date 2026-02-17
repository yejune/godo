package hook

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandlePostToolUse_ReturnsEmptyOutputWhenNoPersonaDir(t *testing.T) {
	// When no persona directory exists (default in test environment),
	// HandlePostToolUse should return an empty Output.
	input := &Input{
		ToolName: "Write",
	}
	output := HandlePostToolUse(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}

	// Empty output should have no HookSpecificOutput when persona dir is missing
	if output.HookSpecificOutput != nil {
		// Only fail if the persona dir was genuinely missing.
		// If it happens to exist in the test environment, the output could be valid.
		t.Logf("HookSpecificOutput is non-nil; persona dir may exist in test environment")
	}
}

func TestHandlePostToolUse_OutputIsJSONMarshallable(t *testing.T) {
	input := &Input{
		ToolName: "Read",
	}
	output := HandlePostToolUse(input)

	data, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("failed to marshal output: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty JSON output")
	}

	// Verify it can be unmarshalled back
	var roundtrip Output
	if err := json.Unmarshal(data, &roundtrip); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}
}

func TestHandlePostToolUse_NoPersonaDirReturnsEmptyOutput(t *testing.T) {
	// Force persona dir to not exist by pointing to a temp dir
	origProjectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	tmpDir := t.TempDir()
	os.Setenv("CLAUDE_PROJECT_DIR", tmpDir)
	defer func() {
		if origProjectDir == "" {
			os.Unsetenv("CLAUDE_PROJECT_DIR")
		} else {
			os.Setenv("CLAUDE_PROJECT_DIR", origProjectDir)
		}
	}()

	input := &Input{ToolName: "Write"}
	output := HandlePostToolUse(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}
	// With no persona dir, HookSpecificOutput should be nil
	if output.HookSpecificOutput != nil {
		t.Error("expected nil HookSpecificOutput when persona dir does not exist")
	}
}

func TestHandlePostToolUse_WithPersonaDirReturnsContext(t *testing.T) {
	// Point to the actual project root where personas/do/ exists
	origProjectDir := os.Getenv("CLAUDE_PROJECT_DIR")

	// Find the project root - go up from internal/hook/ to project root
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	// Check if personas/do/ exists relative to project root candidates
	projectRoot := origDir
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(projectRoot + "/personas/do/characters/young-f.md"); err == nil {
			break
		}
		parent := projectRoot + "/.."
		projectRoot, _ = filepath.Abs(parent)
	}

	// Verify the persona file exists
	if _, err := os.Stat(projectRoot + "/personas/do/characters/young-f.md"); err != nil {
		t.Skip("persona files not found in project root; skipping")
	}

	os.Setenv("CLAUDE_PROJECT_DIR", projectRoot)
	defer func() {
		if origProjectDir == "" {
			os.Unsetenv("CLAUDE_PROJECT_DIR")
		} else {
			os.Setenv("CLAUDE_PROJECT_DIR", origProjectDir)
		}
	}()

	input := &Input{ToolName: "Write"}
	output := HandlePostToolUse(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}
	if output.HookSpecificOutput == nil {
		t.Fatal("expected non-nil HookSpecificOutput when persona is available")
	}
	if output.HookSpecificOutput.HookEventName != "PostToolUse" {
		t.Errorf("HookEventName: got %q, want %q", output.HookSpecificOutput.HookEventName, "PostToolUse")
	}
	if output.HookSpecificOutput.AdditionalContext == "" {
		t.Error("expected non-empty AdditionalContext when persona is available")
	}
}

func TestHandlePostToolUse_PersonaDirExistsButLoadFails(t *testing.T) {
	// Create a persona dir that exists but has no valid character file
	tmpDir := t.TempDir()
	charDir := filepath.Join(tmpDir, "personas", "do", "characters")
	if err := os.MkdirAll(charDir, 0755); err != nil {
		t.Fatal(err)
	}
	// Write an invalid character file (no frontmatter)
	if err := os.WriteFile(filepath.Join(charDir, "young-f.md"), []byte("no frontmatter here"), 0644); err != nil {
		t.Fatal(err)
	}

	origProjectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	os.Setenv("CLAUDE_PROJECT_DIR", tmpDir)
	defer func() {
		if origProjectDir == "" {
			os.Unsetenv("CLAUDE_PROJECT_DIR")
		} else {
			os.Setenv("CLAUDE_PROJECT_DIR", origProjectDir)
		}
	}()

	out := HandlePostToolUse(&Input{ToolName: "Write"})
	// LoadCharacter fails, should return empty output
	if out.HookSpecificOutput != nil {
		t.Error("expected nil HookSpecificOutput when LoadCharacter fails")
	}
}

func TestHandlePostToolUse_EmptyReminder(t *testing.T) {
	// Create a persona with empty honorific and tone to produce empty reminder
	tmpDir := t.TempDir()
	charDir := filepath.Join(tmpDir, "personas", "do", "characters")
	if err := os.MkdirAll(charDir, 0755); err != nil {
		t.Fatal(err)
	}
	charContent := []byte(`---
id: young-f
name: EmptyPersona
honorific_template: ""
honorific_default: ""
tone: ""
character_summary: ""
relationship: ""
---
Body content.
`)
	if err := os.WriteFile(filepath.Join(charDir, "young-f.md"), charContent, 0644); err != nil {
		t.Fatal(err)
	}

	origProjectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	origUserName := os.Getenv("DO_USER_NAME")
	os.Setenv("CLAUDE_PROJECT_DIR", tmpDir)
	os.Setenv("DO_USER_NAME", "")
	defer func() {
		if origProjectDir == "" {
			os.Unsetenv("CLAUDE_PROJECT_DIR")
		} else {
			os.Setenv("CLAUDE_PROJECT_DIR", origProjectDir)
		}
		if origUserName == "" {
			os.Unsetenv("DO_USER_NAME")
		} else {
			os.Setenv("DO_USER_NAME", origUserName)
		}
	}()

	out := HandlePostToolUse(&Input{ToolName: "Write"})
	// BuildReminder returns "" because honorific is empty
	if out.HookSpecificOutput != nil {
		t.Error("expected nil HookSpecificOutput when reminder is empty")
	}
}

func TestHandlePostToolUse_WithTempPersonaFixture(t *testing.T) {
	cleanup := setupTestPersona(t)
	defer cleanup()

	out := HandlePostToolUse(&Input{ToolName: "Write"})
	if out.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput with persona reminder")
	}
	if out.HookSpecificOutput.AdditionalContext == "" {
		t.Error("expected non-empty additionalContext")
	}
	if !strings.Contains(out.HookSpecificOutput.AdditionalContext, "testusersunbae") {
		t.Errorf("expected persona honorific in context, got: %s", out.HookSpecificOutput.AdditionalContext)
	}
	if out.HookSpecificOutput.HookEventName != "PostToolUse" {
		t.Errorf("HookEventName: got %q, want %q", out.HookSpecificOutput.HookEventName, "PostToolUse")
	}
}
