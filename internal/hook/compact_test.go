package hook

import (
	"testing"
)

func Test_HandleCompact_returns_continue(t *testing.T) {
	input := &Input{SessionID: "test-session"}
	output := HandleCompact(input)

	if !output.Continue {
		t.Error("expected Continue=true")
	}
}

func Test_HandleCompact_output_is_not_nil(t *testing.T) {
	input := &Input{}
	output := HandleCompact(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}
}
