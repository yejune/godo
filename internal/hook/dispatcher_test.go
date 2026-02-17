package hook

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func Test_GetStringField_found(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
		"mode": "do",
	}
	got := GetStringField(data, "name", "default")
	if got != "test" {
		t.Errorf("got %q, want %q", got, "test")
	}
}

func Test_GetStringField_missing_key(t *testing.T) {
	data := map[string]interface{}{
		"name": "test",
	}
	got := GetStringField(data, "missing", "fallback")
	if got != "fallback" {
		t.Errorf("got %q, want %q", got, "fallback")
	}
}

func Test_GetStringField_empty_string(t *testing.T) {
	data := map[string]interface{}{
		"name": "",
	}
	got := GetStringField(data, "name", "default")
	if got != "default" {
		t.Errorf("got %q, want %q (empty string should return fallback)", got, "default")
	}
}

func Test_GetStringField_non_string_value(t *testing.T) {
	data := map[string]interface{}{
		"count": 42,
	}
	got := GetStringField(data, "count", "default")
	if got != "default" {
		t.Errorf("got %q, want %q (non-string should return fallback)", got, "default")
	}
}

func Test_GetStringField_nil_map(t *testing.T) {
	got := GetStringField(nil, "key", "fallback")
	if got != "fallback" {
		t.Errorf("got %q, want %q", got, "fallback")
	}
}

func Test_WriteOutput_writes_valid_json(t *testing.T) {
	// Capture stdout
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	output := NewAllowOutput()
	WriteOutput(output)

	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)

	// Verify it is valid JSON
	var decoded Output
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("WriteOutput did not produce valid JSON: %v, output: %q", err, buf.String())
	}
	if decoded.HookSpecificOutput == nil {
		t.Fatal("expected HookSpecificOutput in output")
	}
	if decoded.HookSpecificOutput.PermissionDecision != DecisionAllow {
		t.Errorf("decision: got %q, want %q", decoded.HookSpecificOutput.PermissionDecision, DecisionAllow)
	}
}

func Test_WriteResult_writes_valid_json(t *testing.T) {
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	result := map[string]interface{}{
		"continue": true,
		"message":  "test",
	}
	WriteResult(result)

	w.Close()
	os.Stdout = origStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)

	var decoded map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("WriteResult did not produce valid JSON: %v, output: %q", err, buf.String())
	}
	if decoded["message"] != "test" {
		t.Errorf("message: got %q, want %q", decoded["message"], "test")
	}
}

func Test_ReadInput_empty_stdin(t *testing.T) {
	// Redirect stdin to an empty reader
	origStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	w.Close() // Close immediately to simulate empty stdin
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	input := ReadInput()
	if input == nil {
		t.Fatal("expected non-nil input")
	}
	// Empty stdin returns empty Input
	if input.SessionID != "" {
		t.Errorf("expected empty SessionID, got %q", input.SessionID)
	}
}

func Test_ReadInput_valid_json(t *testing.T) {
	origStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	jsonData := `{"session_id": "sess-456", "tool_name": "Write"}`
	w.WriteString(jsonData)
	w.Close()

	input := ReadInput()
	if input.SessionID != "sess-456" {
		t.Errorf("SessionID: got %q, want %q", input.SessionID, "sess-456")
	}
	if input.ToolName != "Write" {
		t.Errorf("ToolName: got %q, want %q", input.ToolName, "Write")
	}
}

func Test_ReadInput_invalid_json(t *testing.T) {
	origStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	w.WriteString("not valid json")
	w.Close()

	input := ReadInput()
	if input == nil {
		t.Fatal("expected non-nil input even for invalid JSON")
	}
	// Invalid JSON returns empty Input
	if input.SessionID != "" {
		t.Errorf("expected empty SessionID for invalid JSON, got %q", input.SessionID)
	}
}

func Test_ReadStdin_valid_json(t *testing.T) {
	origStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	w.WriteString(`{"key": "value", "num": 42}`)
	w.Close()

	result := ReadStdin()
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result["key"] != "value" {
		t.Errorf("key: got %v, want %q", result["key"], "value")
	}
}

func Test_ReadStdin_empty(t *testing.T) {
	origStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	w.Close()
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	result := ReadStdin()
	if result != nil {
		t.Errorf("expected nil for empty stdin, got %v", result)
	}
}

func Test_ReadStdin_invalid_json(t *testing.T) {
	origStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	w.WriteString("not json")
	w.Close()

	result := ReadStdin()
	if result != nil {
		t.Errorf("expected nil for invalid JSON, got %v", result)
	}
}
