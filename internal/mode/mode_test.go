package mode

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func Test_ReadState_defaults_to_do(t *testing.T) {
	// Work in a temp dir so no .do/.current-mode exists
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(origDir)

	os.Unsetenv("DO_MODE")

	got := ReadState()
	if got != "do" {
		t.Errorf("got %q, want %q", got, "do")
	}
}

func Test_ReadState_from_env(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(origDir)

	t.Setenv("DO_MODE", "focus")

	got := ReadState()
	if got != "focus" {
		t.Errorf("got %q, want %q", got, "focus")
	}
}

func Test_WriteState_and_ReadState_roundtrip(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(origDir)

	os.Unsetenv("DO_MODE")

	WriteState("team")
	got := ReadState()
	if got != "team" {
		t.Errorf("got %q, want %q", got, "team")
	}
}

func Test_WriteState_creates_directory(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(origDir)

	WriteState("focus")

	// Verify the file exists
	data, err := os.ReadFile(filepath.Join(tmp, StateFile))
	if err != nil {
		t.Fatalf("read state file: %v", err)
	}
	if string(data) != "focus\n" {
		t.Errorf("file content: got %q, want %q", string(data), "focus\n")
	}
}

func Test_WriteState_overwrite(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(origDir)

	os.Unsetenv("DO_MODE")

	WriteState("do")
	WriteState("team")
	got := ReadState()
	if got != "team" {
		t.Errorf("got %q, want %q", got, "team")
	}
}

func Test_SetDefaultMode_creates_settings(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(origDir)

	os.MkdirAll(".claude", 0755)

	err := SetDefaultMode("bypassPermissions")
	if err != nil {
		t.Fatalf("SetDefaultMode: %v", err)
	}

	data, err := os.ReadFile(".claude/settings.local.json")
	if err != nil {
		t.Fatalf("read settings: %v", err)
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	perms, ok := settings["permissions"].(map[string]interface{})
	if !ok {
		t.Fatal("permissions key missing")
	}

	if perms["defaultMode"] != "bypassPermissions" {
		t.Errorf("defaultMode: got %v, want %q", perms["defaultMode"], "bypassPermissions")
	}
}

func Test_SetDefaultMode_default_removes_key(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmp)
	defer os.Chdir(origDir)

	os.MkdirAll(".claude", 0755)

	// First set a mode
	SetDefaultMode("bypassPermissions")
	// Then set to "default" which should remove the key
	err := SetDefaultMode("default")
	if err != nil {
		t.Fatalf("SetDefaultMode: %v", err)
	}

	data, err := os.ReadFile(".claude/settings.local.json")
	if err != nil {
		t.Fatalf("read settings: %v", err)
	}

	var settings map[string]interface{}
	json.Unmarshal(data, &settings)
	perms := settings["permissions"].(map[string]interface{})

	if _, exists := perms["defaultMode"]; exists {
		t.Error("defaultMode should be removed when mode is 'default'")
	}
}

func Test_ExecutionModes_valid(t *testing.T) {
	for _, m := range []string{"do", "focus", "team"} {
		if !ExecutionModes[m] {
			t.Errorf("%q should be a valid execution mode", m)
		}
	}
}

func Test_ExecutionModes_invalid(t *testing.T) {
	if ExecutionModes["bogus"] {
		t.Error("'bogus' should not be a valid execution mode")
	}
}
