package extractor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractSettings_WithHooks(t *testing.T) {
	input := []byte(`{
  "permissions": {
    "allow": ["Read", "Write"]
  },
  "outputStyle": "pair",
  "hooks": {
    "PreToolUse": [
      {
        "command": "godo hook PreToolUse",
        "timeout": 5000
      }
    ],
    "PostToolUse": [
      {
        "command": "godo hook PostToolUse"
      }
    ]
  }
}`)

	core, persona, err := ExtractSettings(input)
	if err != nil {
		t.Fatalf("ExtractSettings returned error: %v", err)
	}

	// Core should have permissions and outputStyle, but not hooks
	if _, ok := core["permissions"]; !ok {
		t.Error("core missing 'permissions' key")
	}
	if _, ok := core["outputStyle"]; !ok {
		t.Error("core missing 'outputStyle' key")
	}
	if _, ok := core["hooks"]; ok {
		t.Error("core should not contain 'hooks' key")
	}

	// Persona should have hooks only
	if _, ok := persona["hooks"]; !ok {
		t.Error("persona missing 'hooks' key")
	}
	if len(persona) != 1 {
		t.Errorf("persona should only have 'hooks', got %d keys: %v", len(persona), keys(persona))
	}

	// Verify hooks structure is preserved
	hooks, ok := persona["hooks"].(map[string]interface{})
	if !ok {
		t.Fatal("persona hooks is not a map")
	}
	if _, ok := hooks["PreToolUse"]; !ok {
		t.Error("persona hooks missing 'PreToolUse'")
	}
	if _, ok := hooks["PostToolUse"]; !ok {
		t.Error("persona hooks missing 'PostToolUse'")
	}
}

func TestExtractSettings_WithoutHooks(t *testing.T) {
	input := []byte(`{
  "permissions": {
    "allow": ["Read"]
  },
  "outputStyle": "direct"
}`)

	core, persona, err := ExtractSettings(input)
	if err != nil {
		t.Fatalf("ExtractSettings returned error: %v", err)
	}

	if len(core) != 2 {
		t.Errorf("core should have 2 keys, got %d: %v", len(core), keys(core))
	}
	if _, ok := core["permissions"]; !ok {
		t.Error("core missing 'permissions'")
	}
	if _, ok := core["outputStyle"]; !ok {
		t.Error("core missing 'outputStyle'")
	}

	if len(persona) != 0 {
		t.Errorf("persona should be empty, got %d keys: %v", len(persona), keys(persona))
	}
}

func TestExtractSettings_Empty(t *testing.T) {
	input := []byte(`{}`)

	core, persona, err := ExtractSettings(input)
	if err != nil {
		t.Fatalf("ExtractSettings returned error: %v", err)
	}

	if len(core) != 0 {
		t.Errorf("core should be empty, got %d keys: %v", len(core), keys(core))
	}
	if len(persona) != 0 {
		t.Errorf("persona should be empty, got %d keys: %v", len(persona), keys(persona))
	}
}

func TestExtractSettings_RoundTrip(t *testing.T) {
	input := []byte(`{
  "permissions": {
    "allow": ["Read", "Write", "Bash"]
  },
  "outputStyle": "sprint",
  "hooks": {
    "PreToolUse": [
      {
        "command": "godo hook PreToolUse",
        "timeout": 5000
      }
    ]
  },
  "experimental": true
}`)

	core, persona, err := ExtractSettings(input)
	if err != nil {
		t.Fatalf("ExtractSettings returned error: %v", err)
	}

	// Merge core + persona and verify all original keys are present
	merged := make(map[string]interface{})
	for k, v := range core {
		merged[k] = v
	}
	for k, v := range persona {
		merged[k] = v
	}

	// Parse original for comparison
	var original map[string]interface{}
	if err := json.Unmarshal(input, &original); err != nil {
		t.Fatalf("failed to parse original: %v", err)
	}

	if len(merged) != len(original) {
		t.Errorf("round-trip key count mismatch: merged=%d, original=%d", len(merged), len(original))
	}
	for k := range original {
		if _, ok := merged[k]; !ok {
			t.Errorf("round-trip missing key: %q", k)
		}
	}
}

func TestExtractSettings_InvalidJSON(t *testing.T) {
	input := []byte(`{invalid json`)

	_, _, err := ExtractSettings(input)
	if err == nil {
		t.Error("ExtractSettings should return error for invalid JSON")
	}
}

func TestWriteSettingsFiles(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()

	core := map[string]interface{}{
		"permissions": map[string]interface{}{
			"allow": []interface{}{"Read", "Write"},
		},
		"outputStyle": "pair",
	}
	persona := map[string]interface{}{
		"hooks": map[string]interface{}{
			"PreToolUse": []interface{}{
				map[string]interface{}{
					"command": "godo hook PreToolUse",
				},
			},
		},
	}

	err := WriteSettingsFiles(coreDir, personaDir, core, persona)
	if err != nil {
		t.Fatalf("WriteSettingsFiles returned error: %v", err)
	}

	// Verify core file
	coreData, err := os.ReadFile(filepath.Join(coreDir, "settings-core.json"))
	if err != nil {
		t.Fatalf("failed to read settings-core.json: %v", err)
	}
	var coreRead map[string]interface{}
	if err := json.Unmarshal(coreData, &coreRead); err != nil {
		t.Fatalf("failed to parse settings-core.json: %v", err)
	}
	if _, ok := coreRead["permissions"]; !ok {
		t.Error("settings-core.json missing 'permissions'")
	}
	if _, ok := coreRead["outputStyle"]; !ok {
		t.Error("settings-core.json missing 'outputStyle'")
	}

	// Verify persona file
	personaData, err := os.ReadFile(filepath.Join(personaDir, "settings-hooks.json"))
	if err != nil {
		t.Fatalf("failed to read settings-hooks.json: %v", err)
	}
	var personaRead map[string]interface{}
	if err := json.Unmarshal(personaData, &personaRead); err != nil {
		t.Fatalf("failed to parse settings-hooks.json: %v", err)
	}
	if _, ok := personaRead["hooks"]; !ok {
		t.Error("settings-hooks.json missing 'hooks'")
	}
}

func TestWriteSettingsFiles_EmptyPersona(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()

	core := map[string]interface{}{
		"outputStyle": "direct",
	}
	persona := map[string]interface{}{}

	err := WriteSettingsFiles(coreDir, personaDir, core, persona)
	if err != nil {
		t.Fatalf("WriteSettingsFiles returned error: %v", err)
	}

	// Core file should exist
	if _, err := os.Stat(filepath.Join(coreDir, "settings-core.json")); err != nil {
		t.Errorf("settings-core.json should exist: %v", err)
	}

	// Persona file should still be written (empty object)
	personaData, err := os.ReadFile(filepath.Join(personaDir, "settings-hooks.json"))
	if err != nil {
		t.Fatalf("failed to read settings-hooks.json: %v", err)
	}
	var personaRead map[string]interface{}
	if err := json.Unmarshal(personaData, &personaRead); err != nil {
		t.Fatalf("failed to parse settings-hooks.json: %v", err)
	}
	if len(personaRead) != 0 {
		t.Errorf("settings-hooks.json should be empty object, got %d keys", len(personaRead))
	}
}

func TestWriteSettingsFiles_IndentFormat(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()

	core := map[string]interface{}{
		"outputStyle": "pair",
	}
	persona := map[string]interface{}{}

	err := WriteSettingsFiles(coreDir, personaDir, core, persona)
	if err != nil {
		t.Fatalf("WriteSettingsFiles returned error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(coreDir, "settings-core.json"))
	if err != nil {
		t.Fatalf("failed to read settings-core.json: %v", err)
	}

	// Verify 2-space indentation (json.MarshalIndent with "", "  ")
	expected := "{\n  \"outputStyle\": \"pair\"\n}\n"
	if string(data) != expected {
		t.Errorf("settings-core.json format mismatch\ngot:  %q\nwant: %q", string(data), expected)
	}
}

// keys returns the keys of a map for error messages.
func keys(m map[string]interface{}) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}
