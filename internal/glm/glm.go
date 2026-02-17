package glm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	subdir       = "glm"
	credFilename = "credentials.json"
	credDirPerm  = 0700
	credFilePerm = 0600
)

// Credentials holds the user's GLM API key.
type Credentials struct {
	APIKey    string `json:"api_key"`
	CreatedAt string `json:"created_at"`
}

// GetCredentialsDir returns the directory path for GLM credentials (~/.do/glm/).
func GetCredentialsDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".do", subdir)
}

// GetCredentialsPath returns the full path to the credentials file (~/.do/glm/credentials.json).
func GetCredentialsPath() string {
	return filepath.Join(GetCredentialsDir(), credFilename)
}

// SaveCredentials persists credentials to disk atomically with secure file permissions.
func SaveCredentials(creds *Credentials) error {
	dir := GetCredentialsDir()
	if err := os.MkdirAll(dir, credDirPerm); err != nil {
		return fmt.Errorf("create GLM credentials directory: %w", err)
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal GLM credentials: %w", err)
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
		return nil, fmt.Errorf("read GLM credential file: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, nil
	}

	return &creds, nil
}

// SetupCredentials saves a new API key.
func SetupCredentials(apiKey string) error {
	creds := &Credentials{
		APIKey:    apiKey,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	return SaveCredentials(creds)
}

// SetGLMEnv sets the environment variables for GLM backend in the current process.
func SetGLMEnv(apiKey string) {
	os.Setenv("ANTHROPIC_AUTH_TOKEN", apiKey)
	os.Setenv("ANTHROPIC_BASE_URL", "https://api.z.ai/api/anthropic")
	os.Setenv("ANTHROPIC_DEFAULT_HAIKU_MODEL", "glm-4.7-flash")
	os.Setenv("ANTHROPIC_DEFAULT_SONNET_MODEL", "glm-4.7")
	os.Setenv("ANTHROPIC_DEFAULT_OPUS_MODEL", "glm-5.0")
}

// MaskAPIKey masks an API key for display, showing only prefix and suffix.
func MaskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}
