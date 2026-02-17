package hook

import (
	"encoding/json"
	"testing"
)

func TestHandlePostToolUse_ReturnsEmptyOutput(t *testing.T) {
	input := &Input{ToolName: "Write"}
	output := HandlePostToolUse(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}
}

func TestHandlePostToolUse_OutputIsJSONMarshallable(t *testing.T) {
	input := &Input{ToolName: "Read"}
	output := HandlePostToolUse(input)

	data, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("failed to marshal output: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty JSON output")
	}
}
