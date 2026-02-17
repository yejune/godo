package hook

import (
	"os"
	"strings"
	"testing"
)

func Test_HandleSessionStart_returns_continue(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleSessionStart(input)

	if !output.Continue {
		t.Error("expected Continue=true")
	}
}

func Test_HandleSessionStart_has_system_message(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleSessionStart(input)

	if output.SystemMessage == "" {
		t.Error("expected non-empty SystemMessage")
	}
}

func Test_HandleSessionStart_contains_mode_info(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleSessionStart(input)

	// The message should contain mode information
	if output.SystemMessage == "" {
		t.Fatal("SystemMessage is empty")
	}
	// Mode info should reference the execution mode prefix pattern
	if !containsSubstring(output.SystemMessage, "실행 모드") {
		t.Errorf("SystemMessage should contain mode info, got: %q", output.SystemMessage)
	}
}

func Test_HandleSessionStart_no_persona_dir(t *testing.T) {
	// Ensure persona dir is not found by working from a temp directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	tmpDir := t.TempDir()
	origProjectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	os.Setenv("CLAUDE_PROJECT_DIR", tmpDir)
	defer func() {
		os.Chdir(origDir)
		if origProjectDir == "" {
			os.Unsetenv("CLAUDE_PROJECT_DIR")
		} else {
			os.Setenv("CLAUDE_PROJECT_DIR", origProjectDir)
		}
	}()

	input := &Input{SessionID: "test-session"}
	output := HandleSessionStart(input)

	if !output.Continue {
		t.Error("expected Continue=true even without persona dir")
	}
	if output.SystemMessage == "" {
		t.Error("expected non-empty SystemMessage even without persona dir")
	}
	// Should still contain mode info
	if !containsSubstring(output.SystemMessage, "실행 모드") {
		t.Errorf("SystemMessage should contain mode info, got: %q", output.SystemMessage)
	}
}

func Test_HandleSessionStart_output_is_not_nil(t *testing.T) {
	input := &Input{}
	output := HandleSessionStart(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}
}

func Test_HandleSessionStart_with_persona(t *testing.T) {
	cleanup := setupTestPersona(t)
	defer cleanup()

	out := HandleSessionStart(&Input{})
	if !out.Continue {
		t.Error("expected Continue=true")
	}
	if !strings.Contains(out.SystemMessage, "TestPersona") {
		t.Errorf("expected SystemMessage to contain persona name, got: %s", out.SystemMessage)
	}
	if !strings.Contains(out.SystemMessage, "testusersunbae") {
		t.Errorf("expected SystemMessage to contain honorific, got: %s", out.SystemMessage)
	}
	if !strings.Contains(out.SystemMessage, "full persona content") {
		t.Errorf("expected SystemMessage to contain full content, got: %s", out.SystemMessage)
	}
}
