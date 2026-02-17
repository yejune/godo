package hook

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_HandleCompact_returns_continue(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleCompact(input)

	if !output.Continue {
		t.Error("expected Continue=true")
	}
}

func Test_HandleCompact_has_system_message(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleCompact(input)

	if output.SystemMessage == "" {
		t.Error("expected non-empty SystemMessage")
	}
}

func Test_HandleCompact_contains_mode_info(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleCompact(input)

	if output.SystemMessage == "" {
		t.Fatal("SystemMessage is empty")
	}
	if !containsSubstring(output.SystemMessage, "실행 모드") {
		t.Errorf("SystemMessage should contain mode info, got: %q", output.SystemMessage)
	}
}

func Test_HandleCompact_output_is_not_nil(t *testing.T) {
	input := &Input{}
	output := HandleCompact(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}
}

func Test_HandleCompact_with_persona(t *testing.T) {
	cleanup := setupTestPersona(t)
	defer cleanup()

	out := HandleCompact(&Input{})
	if !out.Continue {
		t.Error("expected Continue=true")
	}
	if !strings.Contains(out.SystemMessage, "testusersunbae") {
		t.Errorf("expected persona reminder in message, got: %s", out.SystemMessage)
	}
}

func Test_HandleCompact_with_checklist(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Ensure persona dir does not interfere
	origProjectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	os.Setenv("CLAUDE_PROJECT_DIR", tmpDir)
	defer func() {
		if origProjectDir == "" {
			os.Unsetenv("CLAUDE_PROJECT_DIR")
		} else {
			os.Setenv("CLAUDE_PROJECT_DIR", origProjectDir)
		}
	}()

	jobsDir := filepath.Join(tmpDir, ".do", "jobs", "26", "02", "18", "test-task")
	if err := os.MkdirAll(jobsDir, 0755); err != nil {
		t.Fatalf("failed to create dirs: %v", err)
	}
	checklist := "# Checklist\n- [o] Task 1\n- [~] Task 2\n- [ ] Task 3\n"
	if err := os.WriteFile(filepath.Join(jobsDir, "checklist.md"), []byte(checklist), 0644); err != nil {
		t.Fatalf("failed to write checklist: %v", err)
	}

	out := HandleCompact(&Input{})
	if !strings.Contains(out.SystemMessage, "[o]1") {
		t.Errorf("expected checklist stats in message, got: %s", out.SystemMessage)
	}
}
