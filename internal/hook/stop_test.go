package hook

import "testing"

func Test_HandleStop_no_checklist(t *testing.T) {
	input := &Input{}
	output := HandleStop(input)

	// When no checklist exists, should return empty output (allow stop)
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
