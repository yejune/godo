package hook

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_ParseChecklistContent_all_statuses(t *testing.T) {
	content := `# Checklist
- [ ] pending item 1
- [ ] pending item 2
- [~] in progress item
- [*] testing item
- [!] blocked item
- [o] done item 1
- [o] done item 2
- [o] done item 3
- [x] failed item
`
	stats := ParseChecklistContent(content)

	if stats.Total != 9 {
		t.Errorf("Total: got %d, want 9", stats.Total)
	}
	if stats.Pending != 2 {
		t.Errorf("Pending: got %d, want 2", stats.Pending)
	}
	if stats.InProgress != 1 {
		t.Errorf("InProgress: got %d, want 1", stats.InProgress)
	}
	if stats.Testing != 1 {
		t.Errorf("Testing: got %d, want 1", stats.Testing)
	}
	if stats.Blocked != 1 {
		t.Errorf("Blocked: got %d, want 1", stats.Blocked)
	}
	if stats.Done != 3 {
		t.Errorf("Done: got %d, want 3", stats.Done)
	}
	if stats.Failed != 1 {
		t.Errorf("Failed: got %d, want 1", stats.Failed)
	}
}

func Test_ParseChecklistContent_empty(t *testing.T) {
	stats := ParseChecklistContent("")
	if stats.Total != 0 {
		t.Errorf("Total: got %d, want 0", stats.Total)
	}
}

func Test_ParseChecklistContent_no_checklist_items(t *testing.T) {
	content := `# Just a heading
Some text here
- regular list item without checkbox
`
	stats := ParseChecklistContent(content)
	if stats.Total != 0 {
		t.Errorf("Total: got %d, want 0", stats.Total)
	}
}

func Test_ParseChecklistContent_indented(t *testing.T) {
	content := `  - [ ] indented pending
    - [o] deeply indented done
`
	stats := ParseChecklistContent(content)
	if stats.Total != 2 {
		t.Errorf("Total: got %d, want 2", stats.Total)
	}
	if stats.Pending != 1 {
		t.Errorf("Pending: got %d, want 1", stats.Pending)
	}
	if stats.Done != 1 {
		t.Errorf("Done: got %d, want 1", stats.Done)
	}
}

func Test_ChecklistStats_HasIncomplete_with_in_progress(t *testing.T) {
	stats := &ChecklistStats{Total: 3, InProgress: 1, Done: 2}
	if !stats.HasIncomplete() {
		t.Error("expected HasIncomplete=true when InProgress > 0")
	}
}

func Test_ChecklistStats_HasIncomplete_with_blocked(t *testing.T) {
	stats := &ChecklistStats{Total: 3, Blocked: 1, Done: 2}
	if !stats.HasIncomplete() {
		t.Error("expected HasIncomplete=true when Blocked > 0")
	}
}

func Test_ChecklistStats_HasIncomplete_all_done(t *testing.T) {
	stats := &ChecklistStats{Total: 3, Done: 3}
	if stats.HasIncomplete() {
		t.Error("expected HasIncomplete=false when all items are done")
	}
}

func Test_ChecklistStats_HasIncomplete_pending_only(t *testing.T) {
	stats := &ChecklistStats{Total: 2, Pending: 2}
	if stats.HasIncomplete() {
		t.Error("expected HasIncomplete=false when only pending items (no in-progress or blocked)")
	}
}

func Test_ChecklistStats_Summary_mixed(t *testing.T) {
	stats := &ChecklistStats{
		Total:      5,
		Done:       2,
		InProgress: 1,
		Pending:    1,
		Failed:     1,
	}
	got := stats.Summary()
	// Should contain all non-zero counts
	if got == "" {
		t.Error("Summary should not be empty")
	}
	// Check specific substrings
	for _, sub := range []string{"[o]2", "[~]1", "[ ]1", "[x]1"} {
		if !contains(got, sub) {
			t.Errorf("Summary %q should contain %q", got, sub)
		}
	}
}

func Test_ChecklistStats_Summary_no_items(t *testing.T) {
	stats := &ChecklistStats{}
	got := stats.Summary()
	if got != "no items" {
		t.Errorf("Summary: got %q, want %q", got, "no items")
	}
}

func Test_ChecklistStats_Summary_all_done(t *testing.T) {
	stats := &ChecklistStats{Total: 3, Done: 3}
	got := stats.Summary()
	if got != "[o]3" {
		t.Errorf("Summary: got %q, want %q", got, "[o]3")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstring(s, sub))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func Test_ParseChecklistFile_valid_file(t *testing.T) {
	f, err := os.CreateTemp("", "checklist-test-*.md")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	content := `# Test Checklist
- [ ] pending task
- [~] in progress task
- [o] done task
`
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write: %v", err)
	}
	f.Close()

	stats, err := ParseChecklistFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Total != 3 {
		t.Errorf("Total: got %d, want 3", stats.Total)
	}
	if stats.Pending != 1 {
		t.Errorf("Pending: got %d, want 1", stats.Pending)
	}
	if stats.InProgress != 1 {
		t.Errorf("InProgress: got %d, want 1", stats.InProgress)
	}
	if stats.Done != 1 {
		t.Errorf("Done: got %d, want 1", stats.Done)
	}
}

func Test_ParseChecklistFile_nonexistent_file(t *testing.T) {
	_, err := ParseChecklistFile("/nonexistent/path/checklist.md")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func Test_ParseChecklistFile_empty_file(t *testing.T) {
	f, err := os.CreateTemp("", "checklist-test-*.md")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.Close()

	stats, err := ParseChecklistFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.Total != 0 {
		t.Errorf("Total: got %d, want 0", stats.Total)
	}
}

func Test_FindLatestChecklist_no_jobs_dir(t *testing.T) {
	// Save current directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	// Work in a temp dir with no .do/jobs/
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	result := FindLatestChecklist()
	if result != "" {
		t.Errorf("expected empty string when no .do/jobs/, got %q", result)
	}
}

func Test_FindLatestChecklist_with_checklist_file(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	defer os.Chdir(origDir)

	// Create .do/jobs/26/02/18/test-task/checklist.md
	checklistDir := filepath.Join(tmpDir, ".do", "jobs", "26", "02", "18", "test-task")
	if err := os.MkdirAll(checklistDir, 0755); err != nil {
		t.Fatalf("failed to create dirs: %v", err)
	}
	checklistPath := filepath.Join(checklistDir, "checklist.md")
	if err := os.WriteFile(checklistPath, []byte("- [ ] task 1\n- [o] task 2\n"), 0644); err != nil {
		t.Fatalf("failed to write checklist: %v", err)
	}

	result := FindLatestChecklist()
	if result == "" {
		t.Fatal("expected non-empty checklist path")
	}
}

func Test_ChecklistStats_HasIncomplete_with_testing(t *testing.T) {
	// Testing status alone does not count as incomplete
	stats := &ChecklistStats{Total: 2, Testing: 1, Done: 1}
	if stats.HasIncomplete() {
		t.Error("expected HasIncomplete=false when only testing items (no in-progress or blocked)")
	}
}

func Test_ChecklistStats_Summary_with_testing(t *testing.T) {
	stats := &ChecklistStats{Total: 1, Testing: 1}
	got := stats.Summary()
	if got != "[*]1" {
		t.Errorf("Summary: got %q, want %q", got, "[*]1")
	}
}

func Test_ChecklistStats_Summary_with_blocked(t *testing.T) {
	stats := &ChecklistStats{Total: 1, Blocked: 1}
	got := stats.Summary()
	if got != "[!]1" {
		t.Errorf("Summary: got %q, want %q", got, "[!]1")
	}
}
