package persona

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GetClaudeSettingsPath returns the path to the Claude settings.json file.
func GetClaudeSettingsPath() string {
	if configDir := os.Getenv("CLAUDE_CONFIG_DIR"); configDir != "" {
		return filepath.Join(configDir, "settings.json")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", "settings.json")
	}
	return filepath.Join(homeDir, ".claude", "settings.json")
}

// ReadClaudeSettings reads and parses the Claude settings file at the given path.
// Returns an empty map if the file does not exist.
func ReadClaudeSettings(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]interface{}), nil
		}
		return nil, err
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("invalid JSON in %s: %w", path, err)
	}
	return settings, nil
}

// WriteClaudeSettings serializes settings to JSON and writes to path.
func WriteClaudeSettings(path string, settings map[string]interface{}) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}

// ApplySpinnerToSettings writes spinnerVerbs into the Claude settings file.
func ApplySpinnerToSettings(verbs []string) error {
	settingsPath := GetClaudeSettingsPath()
	settings, err := ReadClaudeSettings(settingsPath)
	if err != nil {
		return fmt.Errorf("reading settings: %w", err)
	}

	iverbs := make([]interface{}, len(verbs))
	for i, v := range verbs {
		iverbs[i] = v
	}

	settings["spinnerVerbs"] = map[string]interface{}{
		"mode":  "replace",
		"verbs": iverbs,
	}

	if err := WriteClaudeSettings(settingsPath, settings); err != nil {
		return fmt.Errorf("writing settings: %w", err)
	}
	return nil
}
