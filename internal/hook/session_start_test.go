package hook

import (
	"strings"
	"testing"
)

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

func Test_HandleSessionStart_contains_mode(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleSessionStart(input)

	if !strings.Contains(output.SystemMessage, "current_mode:") {
		t.Errorf("SystemMessage should contain current_mode, got: %q", output.SystemMessage)
	}
}

func Test_HandleSessionStart_output_is_not_nil(t *testing.T) {
	input := &Input{}
	output := HandleSessionStart(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}
}
