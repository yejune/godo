package mode

import (
	"encoding/json"
	"os"
	"strings"
)

// StateFile is the path to the mode state file relative to project root.
const StateFile = ".do/.current-mode"

// PermissionModes maps short names to Claude Code's defaultMode setting values.
var PermissionModes = map[string]string{
	"bypass":  "bypassPermissions",
	"accept":  "acceptEdits",
	"default": "default",
	"plan":    "plan",
}

// ExecutionModes defines valid do execution modes.
var ExecutionModes = map[string]bool{
	"do": true, "focus": true, "team": true,
}

// ReadState reads the current execution mode from the state file.
// Falls back to DO_MODE env var, then defaults to "do".
func ReadState() string {
	data, err := os.ReadFile(StateFile)
	if err != nil {
		m := os.Getenv("DO_MODE")
		if m == "" {
			return "do"
		}
		return m
	}
	m := strings.TrimSpace(string(data))
	if m == "" {
		m = os.Getenv("DO_MODE")
		if m == "" {
			return "do"
		}
		return m
	}
	return m
}

// WriteState persists the execution mode to the state file.
func WriteState(mode string) {
	os.MkdirAll(".do", 0755)
	os.WriteFile(StateFile, []byte(mode+"\n"), 0644)
}

// SetDefaultMode updates defaultMode in .claude/settings.local.json.
func SetDefaultMode(defaultMode string) error {
	settingsPath := ".claude/settings.local.json"

	// Read existing settings
	var settings map[string]interface{}
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		settings = make(map[string]interface{})
	} else {
		if err := json.Unmarshal(data, &settings); err != nil {
			settings = make(map[string]interface{})
		}
	}

	// Ensure permissions object exists
	permissions, ok := settings["permissions"].(map[string]interface{})
	if !ok {
		permissions = make(map[string]interface{})
	}

	// Set or remove defaultMode inside permissions
	if defaultMode == "default" {
		delete(permissions, "defaultMode")
	} else {
		permissions["defaultMode"] = defaultMode
	}
	settings["permissions"] = permissions

	// Clean up old top-level defaultMode if exists
	delete(settings, "defaultMode")

	// Write back
	out, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, append(out, '\n'), 0644)
}
