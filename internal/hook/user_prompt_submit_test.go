package hook

import "testing"

func Test_HandleUserPromptSubmit_returns_output(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleUserPromptSubmit(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}
}

func Test_HandleUserPromptSubmit_has_hook_specific_output(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleUserPromptSubmit(input)

	// The function always appends at least the mode reminder,
	// so HookSpecificOutput should be set.
	if output.HookSpecificOutput == nil {
		t.Fatal("expected non-nil HookSpecificOutput")
	}
}

func Test_HandleUserPromptSubmit_event_name(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleUserPromptSubmit(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("HookSpecificOutput is nil")
	}
	if output.HookSpecificOutput.HookEventName != "UserPromptSubmit" {
		t.Errorf("HookEventName: got %q, want %q", output.HookSpecificOutput.HookEventName, "UserPromptSubmit")
	}
}

func Test_HandleUserPromptSubmit_contains_mode_info(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleUserPromptSubmit(input)

	if output.HookSpecificOutput == nil {
		t.Fatal("HookSpecificOutput is nil")
	}
	if output.HookSpecificOutput.AdditionalContext == "" {
		t.Fatal("AdditionalContext is empty")
	}
	if !containsSubstring(output.HookSpecificOutput.AdditionalContext, "실행 모드") {
		t.Errorf("AdditionalContext should contain mode info, got: %q", output.HookSpecificOutput.AdditionalContext)
	}
}
