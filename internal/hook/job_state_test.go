package hook

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_LoadJobState_valid_file(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	content := `{
  "job_id": "test-123",
  "created_at": "2026-02-18T10:00:00Z",
  "workflow_type": "simple",
  "phases": {
    "plan": {"status": "complete"}
  },
  "agents": {
    "expert-backend": {"status": "in_progress", "checklist": "checklists/01_backend.md"}
  }
}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	state, err := LoadJobState(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.JobID != "test-123" {
		t.Errorf("JobID: got %q, want %q", state.JobID, "test-123")
	}
	if state.WorkflowType != "simple" {
		t.Errorf("WorkflowType: got %q, want %q", state.WorkflowType, "simple")
	}
	if state.Phases["plan"].Status != "complete" {
		t.Errorf("plan phase status: got %q, want %q", state.Phases["plan"].Status, "complete")
	}
	if state.Agents["expert-backend"].Status != "in_progress" {
		t.Errorf("agent status: got %q, want %q", state.Agents["expert-backend"].Status, "in_progress")
	}
}

func Test_LoadJobState_nonexistent_file(t *testing.T) {
	_, err := LoadJobState("/nonexistent/path/state.json")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func Test_LoadJobState_invalid_json(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	if err := os.WriteFile(path, []byte("not json"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	_, err := LoadJobState(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func Test_SaveJobState_and_reload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")

	state := &JobState{
		JobID:        "save-test",
		CreatedAt:    "2026-02-18T10:00:00Z",
		WorkflowType: "complex",
		Phases: map[string]PhaseState{
			"analysis": {Status: "complete", StartedAt: "2026-02-18T10:00:00Z", CompletedAt: "2026-02-18T10:05:00Z"},
		},
		Agents: map[string]AgentState{
			"expert-testing": {Status: "pending"},
		},
	}

	if err := SaveJobState(path, state); err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	// Reload and verify
	loaded, err := LoadJobState(path)
	if err != nil {
		t.Fatalf("failed to reload: %v", err)
	}
	if loaded.JobID != "save-test" {
		t.Errorf("JobID: got %q, want %q", loaded.JobID, "save-test")
	}
	if loaded.WorkflowType != "complex" {
		t.Errorf("WorkflowType: got %q, want %q", loaded.WorkflowType, "complex")
	}
	if loaded.Phases["analysis"].Status != "complete" {
		t.Errorf("phase status: got %q, want %q", loaded.Phases["analysis"].Status, "complete")
	}
}

func Test_SaveJobState_invalid_path(t *testing.T) {
	state := &JobState{JobID: "test"}
	err := SaveJobState("/nonexistent/dir/state.json", state)
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
