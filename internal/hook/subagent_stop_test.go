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
	output := HandleSubagentStop(&Input{})
	if output == nil {
		t.Fatal("expected non-nil output")
	}
}
