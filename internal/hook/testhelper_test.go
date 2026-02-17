package hook

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestPersona creates a temporary persona directory structure
// and sets CLAUDE_PROJECT_DIR to point to it.
// Returns a cleanup function that must be deferred.
func setupTestPersona(t *testing.T) func() {
	t.Helper()

	tmpDir := t.TempDir()
	charDir := filepath.Join(tmpDir, "personas", "do", "characters")
	if err := os.MkdirAll(charDir, 0755); err != nil {
		t.Fatal(err)
	}

	charContent := []byte(`---
id: young-f
name: TestPersona
honorific_template: "{{name}}sunbae"
honorific_default: "sunbae"
tone: "bright and energetic"
character_summary: "A bright test persona"
relationship: "colleague"
---
This is the full persona content for testing.
`)
	if err := os.WriteFile(filepath.Join(charDir, "young-f.md"), charContent, 0644); err != nil {
		t.Fatal(err)
	}

	oldProjectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	oldPersona := os.Getenv("DO_PERSONA")
	oldUserName := os.Getenv("DO_USER_NAME")
	os.Setenv("CLAUDE_PROJECT_DIR", tmpDir)
	os.Setenv("DO_PERSONA", "young-f")
	os.Setenv("DO_USER_NAME", "testuser")

	return func() {
		if oldProjectDir == "" {
			os.Unsetenv("CLAUDE_PROJECT_DIR")
		} else {
			os.Setenv("CLAUDE_PROJECT_DIR", oldProjectDir)
		}
		if oldPersona == "" {
			os.Unsetenv("DO_PERSONA")
		} else {
			os.Setenv("DO_PERSONA", oldPersona)
		}
		if oldUserName == "" {
			os.Unsetenv("DO_USER_NAME")
		} else {
			os.Setenv("DO_USER_NAME", oldUserName)
		}
	}
}
