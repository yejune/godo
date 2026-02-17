package hook

import "testing"

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
