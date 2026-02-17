package hook

import (
	"encoding/json"
	"io"
	"os"
)

// ReadStdin reads JSON from stdin into a raw message.
func ReadStdin() map[string]interface{} {
	data, err := io.ReadAll(os.Stdin)
	if err != nil || len(data) == 0 {
		return nil
	}
	var result map[string]interface{}
	if json.Unmarshal(data, &result) != nil {
		return nil
	}
	return result
}

// ReadInput reads JSON from stdin into a structured Input.
func ReadInput() *Input {
	data, err := io.ReadAll(os.Stdin)
	if err != nil || len(data) == 0 {
		return &Input{}
	}
	var input Input
	if json.Unmarshal(data, &input) != nil {
		return &Input{}
	}
	return &input
}

// WriteResult writes a JSON hook result to stdout.
func WriteResult(result map[string]interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(result)
}

// WriteOutput marshals and writes an Output to stdout.
func WriteOutput(output *Output) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(output)
}

// GetStringField extracts a string field from a map with a fallback default.
func GetStringField(data map[string]interface{}, key, fallback string) string {
	if data == nil {
		return fallback
	}
	if v, ok := data[key].(string); ok && v != "" {
		return v
	}
	return fallback
}
