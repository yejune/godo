package rank

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_SaveCredentials_and_LoadCredentials_roundtrip(t *testing.T) {
	// Override home directory to use temp dir
	tmp := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", tmp)
	defer os.Setenv("HOME", origHome)

	creds := &Credentials{
		APIKey:    "test-api-key-123",
		Username:  "testuser",
		UserID:    "user-456",
		CreatedAt: "2026-02-17T00:00:00Z",
	}

	if err := SaveCredentials(creds); err != nil {
		t.Fatalf("SaveCredentials: %v", err)
	}

	loaded, err := LoadCredentials()
	if err != nil {
		t.Fatalf("LoadCredentials: %v", err)
	}
	if loaded == nil {
		t.Fatal("LoadCredentials returned nil")
	}

	if loaded.APIKey != creds.APIKey {
		t.Errorf("APIKey: got %q, want %q", loaded.APIKey, creds.APIKey)
	}
	if loaded.Username != creds.Username {
		t.Errorf("Username: got %q, want %q", loaded.Username, creds.Username)
	}
	if loaded.UserID != creds.UserID {
		t.Errorf("UserID: got %q, want %q", loaded.UserID, creds.UserID)
	}
}

func Test_LoadCredentials_missing_file(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	loaded, err := LoadCredentials()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded != nil {
		t.Error("expected nil for missing credentials")
	}
}

func Test_DeleteCredentials(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	// Save then delete
	creds := &Credentials{APIKey: "key", Username: "user"}
	SaveCredentials(creds)

	if err := DeleteCredentials(); err != nil {
		t.Fatalf("DeleteCredentials: %v", err)
	}

	loaded, _ := LoadCredentials()
	if loaded != nil {
		t.Error("expected nil after deletion")
	}
}

func Test_DeleteCredentials_no_file(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	// Should not error when file doesn't exist
	if err := DeleteCredentials(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func Test_HasCredentials(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	if HasCredentials() {
		t.Error("expected false before saving")
	}

	SaveCredentials(&Credentials{APIKey: "key"})

	if !HasCredentials() {
		t.Error("expected true after saving")
	}
}

func Test_GetCredentialsDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := GetCredentialsDir()
	expected := filepath.Join(tmp, ".do", "rank")
	if dir != expected {
		t.Errorf("got %q, want %q", dir, expected)
	}
}

func Test_SaveCredentials_file_permissions(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	SaveCredentials(&Credentials{APIKey: "secret"})

	info, err := os.Stat(GetCredentialsPath())
	if err != nil {
		t.Fatalf("stat: %v", err)
	}

	// File should have 0600 permissions
	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("permissions: got %o, want 0600", perm)
	}
}
