package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const configFile = ".do/.claude-launch.json"

// LaunchConfig stores default flags for godo claude.
type LaunchConfig struct {
	Bypass   bool `json:"bypass"`
	Chrome   bool `json:"chrome"`
	Continue bool `json:"continue"`
	AutoSync bool `json:"auto_sync"`
}

// LoadConfig reads the launch config from .do/.claude-launch.json.
func LoadConfig() LaunchConfig {
	cfg := LaunchConfig{Chrome: false} // default: no-chrome
	data, err := os.ReadFile(configFile)
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, &cfg)
	return cfg
}

// SaveConfig writes the launch config to .do/.claude-launch.json.
func SaveConfig(cfg LaunchConfig) error {
	_ = os.MkdirAll(".do", 0755)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, append(data, '\n'), 0644)
}

// BuildClaudeFlags generates claude CLI flags from config.
func BuildClaudeFlags(cfg LaunchConfig) []string {
	var flags []string
	if cfg.Bypass {
		flags = append(flags, "--dangerously-skip-permissions")
	}
	if !cfg.Chrome {
		flags = append(flags, "--no-chrome")
	}
	if cfg.Continue {
		flags = append(flags, "--continue")
	}
	return flags
}

// ParseOverrides parses godo-specific flags from args and overrides config.
// Returns the updated config and filtered args (without godo-specific flags).
func ParseOverrides(cfg LaunchConfig, args []string) (LaunchConfig, []string) {
	var filtered []string
	for _, arg := range args {
		switch arg {
		case "--chrome":
			cfg.Chrome = true
		case "--no-chrome":
			cfg.Chrome = false
		case "-b", "--bypass":
			cfg.Bypass = true
		case "-c", "--continue":
			cfg.Continue = true
		default:
			filtered = append(filtered, arg)
		}
	}
	return cfg, filtered
}

// GetInstalledVersion runs "godo version" and parses the output.
func GetInstalledVersion(exePath string) string {
	out, err := exec.Command(exePath, "version").Output()
	if err != nil {
		return ""
	}
	s := strings.TrimSpace(string(out))
	parts := strings.Fields(s)
	if len(parts) >= 3 {
		return parts[2]
	}
	if len(parts) >= 1 {
		return parts[len(parts)-1]
	}
	return s
}

// CheckBrewOutdated checks if godo has an update available via brew.
func CheckBrewOutdated() bool {
	outdatedOut, err := exec.Command("brew", "outdated", "godo").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(outdatedOut), "godo")
}

// RunSelfUpdate runs "godo selfupdate" using the given executable path.
func RunSelfUpdate(exePath string) error {
	cmd := exec.Command(exePath, "selfupdate")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunSyncCmd runs "godo sync" using the given executable path.
func RunSyncCmd(exePath string) error {
	cmd := exec.Command(exePath, "sync")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PrintAutoUpdateStatus prints auto-update progress messages.
func PrintAutoUpdateStatus(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
}
