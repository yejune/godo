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

	// Core should have permissions, but not hooks
	if _, ok := core["permissions"]; !ok {
		t.Error("core missing 'permissions' key")
	}
	if _, ok := core["hooks"]; ok {
		t.Error("core should not contain 'hooks' key")
	}

	// Persona should have hooks only
	if _, ok := persona["hooks"]; !ok {
		t.Error("persona missing 'hooks' key")
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
  }
}`)

	core, persona, err := ExtractSettings(input)
	if err != nil {
		t.Fatalf("ExtractSettings returned error: %v", err)
	}

	if len(core) != 1 {
		t.Errorf("core should have 1 key (permissions), got %d: %v", len(core), keys(core))
	}
	if _, ok := core["permissions"]; !ok {
		t.Error("core missing 'permissions'")
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

func TestExtractSettings_PersonaKeyClassification(t *testing.T) {
	input := []byte(`{
  "permissions": {
    "allow": ["Read", "Write", "Bash"]
  },
  "outputStyle": "pair",
  "plansDirectory": ".do/plans",
  "attribution": {
    "commit": "Do-Focus AI <do@example.com>"
  },
  "statusLine": {
    "command": ".do/status_line.sh"
  },
  "hooks": {
    "PreToolUse": [{"command": "godo hook pre-tool"}]
  },
  "experimental": true,
  "cleanupPeriodDays": 30
}`)

	core, persona, err := ExtractSettings(input)
	if err != nil {
		t.Fatalf("ExtractSettings returned error: %v", err)
	}

	// Core should have: permissions, experimental, cleanupPeriodDays
	expectedCore := []string{"permissions", "experimental", "cleanupPeriodDays"}
	for _, k := range expectedCore {
		if _, ok := core[k]; !ok {
			t.Errorf("core missing expected key %q", k)
		}
	}

	// Core should NOT have persona keys
	personaKeys := []string{"outputStyle", "plansDirectory", "attribution", "statusLine", "hooks"}
	for _, k := range personaKeys {
		if _, ok := core[k]; ok {
			t.Errorf("core should not contain persona key %q", k)
		}
	}

	// Persona should have all persona keys
	for _, k := range personaKeys {
		if _, ok := persona[k]; !ok {
			t.Errorf("persona missing expected key %q", k)
		}
	}

	// Persona should NOT have core keys
	for _, k := range expectedCore {
		if _, ok := persona[k]; ok {
			t.Errorf("persona should not contain core key %q", k)
		}
	}

	// Verify persona values preserved
	if persona["outputStyle"] != "pair" {
		t.Errorf("expected persona outputStyle 'pair', got %v", persona["outputStyle"])
	}
	if persona["plansDirectory"] != ".do/plans" {
		t.Errorf("expected persona plansDirectory '.do/plans', got %v", persona["plansDirectory"])
	}
}

func TestExtractSettings_EnvFieldSplitting(t *testing.T) {
	input := []byte(`{
  "permissions": {"allow": ["Read"]},
  "env": {
    "DO_MODE": "focus",
    "DO_USER_NAME": "max",
    "MOAI_CONFIG_SOURCE": "yaml"
  }
}`)

	core, persona, err := ExtractSettings(input)
	if err != nil {
		t.Fatalf("ExtractSettings returned error: %v", err)
	}

	// Core env should have DO_MODE and DO_USER_NAME
	coreEnv, ok := core["env"].(map[string]interface{})
	if !ok {
		t.Fatal("core should have 'env' as map")
	}
	if _, ok := coreEnv["DO_MODE"]; !ok {
		t.Error("core env missing 'DO_MODE'")
	}
	if _, ok := coreEnv["DO_USER_NAME"]; !ok {
		t.Error("core env missing 'DO_USER_NAME'")
	}
	if _, ok := coreEnv["MOAI_CONFIG_SOURCE"]; ok {
		t.Error("core env should not contain 'MOAI_CONFIG_SOURCE'")
	}

	// Persona env should have MOAI_CONFIG_SOURCE
	personaEnv, ok := persona["env"].(map[string]interface{})
	if !ok {
		t.Fatal("persona should have 'env' as map")
	}
	if _, ok := personaEnv["MOAI_CONFIG_SOURCE"]; !ok {
		t.Error("persona env missing 'MOAI_CONFIG_SOURCE'")
	}
	if len(personaEnv) != 1 {
		t.Errorf("persona env should have 1 key, got %d: %v", len(personaEnv), personaEnv)
	}
}

func TestExtractSettings_EnvAllCore(t *testing.T) {
	input := []byte(`{
  "env": {
    "DO_MODE": "focus",
    "DO_USER_NAME": "max"
  }
}`)

	core, persona, err := ExtractSettings(input)
	if err != nil {
		t.Fatalf("ExtractSettings returned error: %v", err)
	}

	// All env keys are core, so persona should have no env
	if _, ok := core["env"]; !ok {
		t.Error("core should have 'env'")
	}
	if _, ok := persona["env"]; ok {
		t.Error("persona should not have 'env' when all env keys are core")
	}
}

func TestExtractSettings_EnvAllPersona(t *testing.T) {
	input := []byte(`{
  "env": {
    "MOAI_CONFIG_SOURCE": "yaml"
  }
}`)

	core, persona, err := ExtractSettings(input)
	if err != nil {
		t.Fatalf("ExtractSettings returned error: %v", err)
	}

	// All env keys are persona, so core should have no env
	if _, ok := core["env"]; ok {
		t.Error("core should not have 'env' when all env keys are persona")
	}
	if _, ok := persona["env"]; !ok {
		t.Error("persona should have 'env'")
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
  "experimental": true,
  "plansDirectory": ".do/plans",
  "env": {
    "DO_MODE": "focus",
    "MOAI_CONFIG_SOURCE": "yaml"
  }
}`)

	core, persona, err := ExtractSettings(input)
	if err != nil {
		t.Fatalf("ExtractSettings returned error: %v", err)
	}

	// Merge core + persona env fields and top-level fields
	// All original keys should be accounted for
	var original map[string]interface{}
	if err := json.Unmarshal(input, &original); err != nil {
		t.Fatalf("failed to parse original: %v", err)
	}

	// Check all non-env top-level keys are in either core or persona
	for k := range original {
		if k == "env" {
			continue // env is split, checked separately
		}
		_, inCore := core[k]
		_, inPersona := persona[k]
		if !inCore && !inPersona {
			t.Errorf("round-trip missing key: %q", k)
		}
	}

	// Check env sub-keys are all accounted for
	if origEnv, ok := original["env"].(map[string]interface{}); ok {
		coreEnv, _ := core["env"].(map[string]interface{})
		personaEnv, _ := persona["env"].(map[string]interface{})
		for k := range origEnv {
			_, inCore := coreEnv[k]
			_, inPersona := personaEnv[k]
			if !inCore && !inPersona {
				t.Errorf("round-trip missing env key: %q", k)
			}
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

func TestExtractSettings_BackwardCompat_HooksOnly(t *testing.T) {
	// Verify backward compatibility: settings with only hooks
	// still works the same way as the original implementation.
	input := []byte(`{
  "permissions": {"allow": ["Read"]},
  "outputStyle": "pair",
  "hooks": {"PreToolUse": [{"command": "test"}]}
}`)

	core, persona, err := ExtractSettings(input)
	if err != nil {
		t.Fatalf("ExtractSettings returned error: %v", err)
	}

	// With the new classification, outputStyle is now persona too
	if _, ok := core["permissions"]; !ok {
		t.Error("core missing 'permissions'")
	}
	if _, ok := persona["hooks"]; !ok {
		t.Error("persona missing 'hooks'")
	}
	if _, ok := persona["outputStyle"]; !ok {
		t.Error("persona missing 'outputStyle'")
	}
}

func TestWriteSettingsFiles(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()

	core := map[string]interface{}{
		"permissions": map[string]interface{}{
			"allow": []interface{}{"Read", "Write"},
		},
	}
	persona := map[string]interface{}{
		"hooks": map[string]interface{}{
			"PreToolUse": []interface{}{
				map[string]interface{}{
					"command": "godo hook PreToolUse",
				},
			},
		},
		"outputStyle": "pair",
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
		"permissions": map[string]interface{}{
			"allow": []interface{}{"Read"},
		},
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

func TestSplitEnvField_NonMapValue(t *testing.T) {
	// env value that isn't a map should return nil, nil
	core, persona := splitEnvField("not a map")
	if core != nil {
		t.Errorf("expected nil core for non-map env, got %v", core)
	}
	if persona != nil {
		t.Errorf("expected nil persona for non-map env, got %v", persona)
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
