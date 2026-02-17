package glm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	Subdir       = "glm"
	CredFilename = "credentials.json"
	CredDirPerm  = 0700
	CredFilePerm = 0600
)

// Credentials holds the user's GLM API key.
type Credentials struct {
	APIKey    string `json:"api_key"`
	CreatedAt string `json:"created_at"`
}

// CredentialsDir returns the directory path for GLM credentials (~/.do/glm/).
func CredentialsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".do", Subdir)
}

// CredentialsPath returns the full path to the credentials file (~/.do/glm/credentials.json).
func CredentialsPath() string {
	return filepath.Join(CredentialsDir(), CredFilename)
}

// SaveCredentials persists credentials to disk atomically with secure file permissions.
func SaveCredentials(creds *Credentials) error {
	dir := CredentialsDir()
	if err := os.MkdirAll(dir, CredDirPerm); err != nil {
		return fmt.Errorf("create GLM credentials directory: %w", err)
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal GLM credentials: %w", err)
	}

	credPath := CredentialsPath()
	tmpFile := credPath + ".tmp"

	if err := os.WriteFile(tmpFile, data, CredFilePerm); err != nil {
		return fmt.Errorf("write temp credential file: %w", err)
	}

	if err := os.Rename(tmpFile, credPath); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("rename credential file: %w", err)
	}

	return nil
}

// LoadCredentials reads credentials from disk.
// Returns (nil, nil) if the file does not exist or contains invalid JSON.
func LoadCredentials() (*Credentials, error) {
	data, err := os.ReadFile(CredentialsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read GLM credential file: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, nil
	}

	return &creds, nil
}

// MaskAPIKey masks an API key for display, showing only prefix and suffix.
func MaskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}
