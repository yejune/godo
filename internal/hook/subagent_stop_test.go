package hook

import "testing"

func Test_HandleSubagentStop_clean_repo_returns_continue(t *testing.T) {
	original := GitStatus
	defer func() { GitStatus = original }()
	GitStatus = func() (bool, string) { return false, "" }

	input := &Input{SessionID: "test-session"}
	output := HandleSubagentStop(input)

	if !output.Continue {
		t.Error("expected Continue=true")
	}
	if output.HookSpecificOutput != nil {
		t.Error("expected nil HookSpecificOutput when repo is clean")
	}
}

func Test_HandleSubagentStop_uncommitted_changes_warns(t *testing.T) {
	original := GitStatus
	defer func() { GitStatus = original }()
	GitStatus = func() (bool, string) { return true, "M dirty.go" }

	input := &Input{}
	output := HandleSubagentStop(input)

	if !output.Continue {
		t.Error("expected Continue=true even with warning")
	}
	if output.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput with warning")
	}
	if output.HookSpecificOutput.AdditionalContext == "" {
		t.Error("expected non-empty AdditionalContext")
	}
}

func Test_HandleSubagentStop_returns_non_nil(t *testing.T) {
	original := GitStatus
	defer func() { GitStatus = original }()
	GitStatus = func() (bool, string) { return false, "" }

	output := HandleSubagentStop(&Input{})
	if output == nil {
		t.Fatal("expected non-nil output")
	}
}
