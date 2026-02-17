package hook

import "testing"

func Test_HandleSubagentStop_returns_continue(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleSubagentStop(input)

	if !output.Continue {
		t.Error("expected Continue=true")
	}
}

func Test_HandleSubagentStop_returns_non_nil(t *testing.T) {
	input := &Input{}
	output := HandleSubagentStop(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}
}

func Test_HandleSubagentStop_no_decision(t *testing.T) {
	input := &Input{}
	output := HandleSubagentStop(input)

	if output.Decision != "" {
		t.Errorf("expected empty Decision, got %q", output.Decision)
	}
}

func Test_HandleSubagentStop_no_hook_specific_output(t *testing.T) {
	input := &Input{}
	output := HandleSubagentStop(input)

	if output.HookSpecificOutput != nil {
		t.Error("expected nil HookSpecificOutput")
	}
}
