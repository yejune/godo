package hook

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_HandleStop_no_checklist(t *testing.T) {
	origDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	input := &Input{}
	output := HandleStop(input)

	if output.Decision != "" {
		t.Errorf("expected empty Decision when no checklist, got %q", output.Decision)
	}
}

func Test_HandleStop_returns_output(t *testing.T) {
	input := &Input{}
	output := HandleStop(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}
}

func Test_HandleStop_blocks_with_in_progress_items(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

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
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

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

func Test_HandleStop_blocks_with_blocked_items(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

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

func Test_HandleSessionEnd_returns_continue(t *testing.T) {
	input := &Input{}
	output := HandleSessionEnd(input)

	if !output.Continue {
		t.Error("expected Continue=true")
	}
}
