package hook

import "testing"

func Test_HandleCompact_returns_continue(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleCompact(input)

	if !output.Continue {
		t.Error("expected Continue=true")
	}
}

func Test_HandleCompact_has_system_message(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleCompact(input)

	if output.SystemMessage == "" {
		t.Error("expected non-empty SystemMessage")
	}
}

func Test_HandleCompact_contains_mode_info(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleCompact(input)

	if output.SystemMessage == "" {
		t.Fatal("SystemMessage is empty")
	}
	if !containsSubstring(output.SystemMessage, "실행 모드") {
		t.Errorf("SystemMessage should contain mode info, got: %q", output.SystemMessage)
	}
}

func Test_HandleCompact_output_is_not_nil(t *testing.T) {
	input := &Input{}
	output := HandleCompact(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}
}
