package rank

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	subdir       = "rank"
	credFilename = "credentials.json"
	credDirPerm  = 0700
	credFilePerm = 0600
)

// Credentials holds the user's Rank API authentication credentials.
type Credentials struct {
	APIKey    string `json:"api_key"`
	Username  string `json:"username"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

// GetCredentialsDir returns the directory path for rank credentials (~/.do/rank/).
func GetCredentialsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".do", subdir)
}

// GetCredentialsPath returns the full path to the credentials file (~/.do/rank/credentials.json).
func GetCredentialsPath() string {
	return filepath.Join(GetCredentialsDir(), credFilename)
}

// SaveCredentials persists credentials to disk atomically with secure file permissions.
func SaveCredentials(creds *Credentials) error {
	dir := GetCredentialsDir()
	if err := os.MkdirAll(dir, credDirPerm); err != nil {
		return fmt.Errorf("create rank credentials directory: %w", err)
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal rank credentials: %w", err)
	}

	credPath := GetCredentialsPath()
	tmpFile := credPath + ".tmp"

	if err := os.WriteFile(tmpFile, data, credFilePerm); err != nil {
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
	data, err := os.ReadFile(GetCredentialsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read rank credential file: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, nil
	}

	return &creds, nil
}

// DeleteCredentials removes the credentials file.
func DeleteCredentials() error {
	err := os.Remove(GetCredentialsPath())
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove rank credential file: %w", err)
	}
	return nil
}

// HasCredentials checks whether a credentials file exists.
func HasCredentials() bool {
	_, err := os.Stat(GetCredentialsPath())
	return err == nil
}
