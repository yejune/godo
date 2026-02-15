package extractor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// personaSettingsKeys lists settings.json top-level keys that are persona-specific.
// During extraction, these keys are moved from core to persona.
// During assembly, they are injected from manifest.Settings into output settings.
var personaSettingsKeys = map[string]bool{
	"hooks":          true, // persona-specific hook commands
	"outputStyle":    true, // persona's output style selection
	"plansDirectory": true, // persona's plans directory path
	"attribution":    true, // persona's commit attribution
	"statusLine":     true, // persona's status line script
}

// personaEnvKeys lists env sub-keys within the "env" settings field
// that are persona-specific. The "env" field itself is mixed:
// some keys are persona, others are core.
var personaEnvKeys = map[string]bool{
	"MOAI_CONFIG_SOURCE": true, // moai config system reference
}

// ExtractSettings splits a settings.json into core and persona parts.
// Persona keys: hooks, outputStyle, plansDirectory, attribution, statusLine.
// The env field is mixed: keys in personaEnvKeys go to persona, rest stays core.
// All other keys are core.
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
		if personaSettingsKeys[k] {
			persona[k] = v
		} else if k == "env" {
			// Split env field into core and persona sub-keys
			coreEnv, personaEnv := splitEnvField(v)
			if len(coreEnv) > 0 {
				core[k] = coreEnv
			}
			if len(personaEnv) > 0 {
				persona[k] = personaEnv
			}
		} else {
			core[k] = v
		}
	}

	return core, persona, nil
}

// splitEnvField splits the env map into core and persona sub-keys.
func splitEnvField(v interface{}) (core, persona map[string]interface{}) {
	envMap, ok := v.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	core = make(map[string]interface{})
	persona = make(map[string]interface{})

	for k, val := range envMap {
		if personaEnvKeys[k] {
			persona[k] = val
		} else {
			core[k] = val
		}
	}
	return core, persona
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
