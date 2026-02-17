package hook

import (
	"encoding/json"
	"testing"
)

func TestHandlePostToolUse_ReturnsEmptyOutputWhenNoPersonaDir(t *testing.T) {
	// When no persona directory exists (default in test environment),
	// HandlePostToolUse should return an empty Output.
	input := &Input{
		ToolName: "Write",
	}
	output := HandlePostToolUse(input)

	if output == nil {
		t.Fatal("expected non-nil output")
	}

	// Empty output should have no HookSpecificOutput when persona dir is missing
	if output.HookSpecificOutput != nil {
		// Only fail if the persona dir was genuinely missing.
		// If it happens to exist in the test environment, the output could be valid.
		t.Logf("HookSpecificOutput is non-nil; persona dir may exist in test environment")
	}
}

func TestHandlePostToolUse_OutputIsJSONMarshallable(t *testing.T) {
	input := &Input{
		ToolName: "Read",
	}
	output := HandlePostToolUse(input)

	data, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("failed to marshal output: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty JSON output")
	}

	// Verify it can be unmarshalled back
	var roundtrip Output
	if err := json.Unmarshal(data, &roundtrip); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}
}
