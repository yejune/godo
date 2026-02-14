package extractor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ExtractSettings splits a settings.json into core and persona parts.
// Core fields: permissions, outputStyle, and other non-hook settings.
// Persona fields: hooks section (references persona-specific binaries like godo/moai).
//
// Input: source settings.json content ([]byte)
// Output: coreSettings (map), personaSettings (map), error
func ExtractSettings(data []byte) (core map[string]interface{}, persona map[string]interface{}, err error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, fmt.Errorf("failed to parse settings JSON: %w", err)
	}

	core = make(map[string]interface{})
	persona = make(map[string]interface{})

	for k, v := range raw {
		if k == "hooks" {
			persona[k] = v
		} else {
			core[k] = v
		}
	}

	return core, persona, nil
}

// WriteSettingsFiles writes the split settings to their respective directories.
// Core settings are written to coreDir/settings-core.json.
// Persona settings are written to personaDir/settings-hooks.json.
func WriteSettingsFiles(coreDir, personaDir string, core, persona map[string]interface{}) error {
	if err := writeJSON(filepath.Join(coreDir, "settings-core.json"), core); err != nil {
		return fmt.Errorf("failed to write settings-core.json: %w", err)
	}

	if err := writeJSON(filepath.Join(personaDir, "settings-hooks.json"), persona); err != nil {
		return fmt.Errorf("failed to write settings-hooks.json: %w", err)
	}

	return nil
}

// writeJSON marshals data as indented JSON and writes to the given path.
func writeJSON(path string, data map[string]interface{}) error {
	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	buf = append(buf, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return os.WriteFile(path, buf, 0o644)
}
