package hook

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_HandleStop_no_checklist(t *testing.T) {
	original := GitStatus
	defer func() { GitStatus = original }()
	GitStatus = func() (bool, string) { return false, "" }

	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	input := &Input{}
	output := HandleStop(input)

	// When no checklist exists and git is clean, should return empty output (allow stop)
	if output.Decision != "" {
		t.Errorf("expected empty Decision when no checklist, got %q", output.Decision)
	}
}

func Test_HandleStop_returns_output(t *testing.T) {
	original := GitStatus
	defer func() { GitStatus = original }()
	GitStatus = func() (bool, string) { return false, "" }

	input := &Input{}
	output := HandleStop(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}
}

func Test_ParseChecklistContent_counts_for_stop(t *testing.T) {
	// Verify the stats that HandleStop would use
	content := `- [~] work in progress
- [o] completed task
- [ ] pending task
`
	stats := ParseChecklistContent(content)

	if stats.Total != 3 {
		t.Errorf("Total: got %d, want 3", stats.Total)
	}
	if stats.InProgress != 1 {
		t.Errorf("InProgress: got %d, want 1", stats.InProgress)
	}
	if !stats.HasIncomplete() {
		t.Error("expected HasIncomplete=true with in-progress items")
	}
}

func Test_HandleSessionEnd_returns_continue(t *testing.T) {
	input := &Input{}
	output := HandleSessionEnd(input)

	if !output.Continue {
		t.Error("expected Continue=true")
	}
}

func Test_HandleStop_blocks_with_in_progress_items(t *testing.T) {
	original := GitStatus
	defer func() { GitStatus = original }()
	GitStatus = func() (bool, string) { return false, "" }

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create a checklist with in-progress items
	checklistDir := filepath.Join(tmpDir, ".do", "jobs", "26", "02", "18", "test-task")
	if err := os.MkdirAll(checklistDir, 0755); err != nil {
		t.Fatalf("failed to create dirs: %v", err)
	}
	content := "- [~] work in progress\n- [o] done task\n- [ ] pending\n"
	if err := os.WriteFile(filepath.Join(checklistDir, "checklist.md"), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write checklist: %v", err)
	}

	input := &Input{}
	output := HandleStop(input)

	if output.Decision != DecisionBlock {
		t.Errorf("expected Decision %q when in-progress items exist, got %q", DecisionBlock, output.Decision)
	}
	if output.Reason == "" {
		t.Error("expected non-empty Reason when blocking")
	}
}

func Test_HandleStop_allows_when_all_done(t *testing.T) {
	original := GitStatus
	defer func() { GitStatus = original }()
	GitStatus = func() (bool, string) { return false, "" }

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create a checklist with all items done
	checklistDir := filepath.Join(tmpDir, ".do", "jobs", "26", "02", "18", "test-task")
	if err := os.MkdirAll(checklistDir, 0755); err != nil {
		t.Fatalf("failed to create dirs: %v", err)
	}
	content := "- [o] done task 1\n- [o] done task 2\n- [o] done task 3\n"
	if err := os.WriteFile(filepath.Join(checklistDir, "checklist.md"), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write checklist: %v", err)
	}

	input := &Input{}
	output := HandleStop(input)

	if output.Decision != "" {
		t.Errorf("expected empty Decision when all done, got %q", output.Decision)
	}
}

func Test_HandleStop_empty_checklist_file(t *testing.T) {
	original := GitStatus
	defer func() { GitStatus = original }()
	GitStatus = func() (bool, string) { return false, "" }

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create a checklist file with no items (only header text, no checkboxes)
	checklistDir := filepath.Join(tmpDir, ".do", "jobs", "26", "02", "18", "test-empty")
	if err := os.MkdirAll(checklistDir, 0755); err != nil {
		t.Fatalf("failed to create dirs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(checklistDir, "checklist.md"), []byte("# Empty Checklist\nNo items here.\n"), 0644); err != nil {
		t.Fatalf("failed to write checklist: %v", err)
	}

	input := &Input{}
	output := HandleStop(input)

	// stats.Total == 0, so should allow stop (no block)
	if output.Decision == DecisionBlock {
		t.Error("should not block when checklist has no items")
	}
}

func Test_HandleStop_blocks_with_blocked_items(t *testing.T) {
	original := GitStatus
	defer func() { GitStatus = original }()
	GitStatus = func() (bool, string) { return false, "" }

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create a checklist with blocked items
	checklistDir := filepath.Join(tmpDir, ".do", "jobs", "26", "02", "18", "test-task")
	if err := os.MkdirAll(checklistDir, 0755); err != nil {
		t.Fatalf("failed to create dirs: %v", err)
	}
	content := "- [!] blocked task\n- [o] done task\n"
	if err := os.WriteFile(filepath.Join(checklistDir, "checklist.md"), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write checklist: %v", err)
	}

	input := &Input{}
	output := HandleStop(input)

	if output.Decision != DecisionBlock {
		t.Errorf("expected Decision %q when blocked items exist, got %q", DecisionBlock, output.Decision)
	}
}

func Test_HandleStop_blocks_with_uncommitted_changes(t *testing.T) {
	original := GitStatus
	defer func() { GitStatus = original }()
	GitStatus = func() (bool, string) { return true, "M uncommitted.go" }

	// No checklist (use a temp dir with no .do/jobs)
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	input := &Input{}
	output := HandleStop(input)

	if output.Decision != DecisionBlock {
		t.Errorf("expected Decision %q with uncommitted changes, got %q", DecisionBlock, output.Decision)
	}
	if output.Reason == "" {
		t.Error("expected non-empty Reason")
	}
}

func Test_HandleStop_allows_clean_state(t *testing.T) {
	original := GitStatus
	defer func() { GitStatus = original }()
	GitStatus = func() (bool, string) { return false, "" }

	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	input := &Input{}
	output := HandleStop(input)

	if output.Decision != "" {
		t.Errorf("expected empty Decision when clean, got %q", output.Decision)
	}
}
