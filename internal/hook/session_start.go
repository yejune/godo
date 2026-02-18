package hook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yejune/godo/internal/mode"
	"github.com/yejune/godo/internal/persona"
)

// Version is the current godo version, set via SetVersion from main.
// When empty, version check is skipped.
var Version string

// SetVersion sets the version string used for update checks.
func SetVersion(v string) {
	Version = v
}

// HandleSessionStart handles the SessionStart hook event.
// Applies spinner verbs, checks for updates, and injects project info.
func HandleSessionStart(input *Input) *Output {
	personaType := os.Getenv("DO_PERSONA")
	if personaType == "" {
		personaType = "young-f"
	}

	currentMode := mode.ReadState()

	// Apply spinner verbs (file-based first, then hardcoded fallback)
	verbs := persona.GetSpinnerVerbs(personaType)
	if len(verbs) > 0 {
		_ = persona.ApplySpinnerToSettings(verbs)
	}

	var systemMsgParts []string

	// Version update check
	latestVer := checkLatestVersion()
	if latestVer != "" && latestVer != Version && Version != "dev" && Version != "" {
		systemMsgParts = append(systemMsgParts,
			fmt.Sprintf("godo 업데이트 있음: %s → %s (brew upgrade godo)", Version, latestVer))
	}

	// Project info from CWD
	cwd := ""
	if input != nil {
		cwd = input.CWD
	}

	systemMsgParts = append(systemMsgParts, "current_mode: "+currentMode)

	projectName, projectType, projectLang := detectProjectInfo(cwd)
	var projectInfoParts []string
	if projectName != "" {
		projectInfoParts = append(projectInfoParts, "project: "+projectName)
	}
	if projectType != "" {
		projectInfoParts = append(projectInfoParts, "type: "+projectType)
	}
	if projectLang != "" {
		projectInfoParts = append(projectInfoParts, "lang: "+projectLang)
	}
	if len(projectInfoParts) > 0 {
		systemMsgParts = append(systemMsgParts, strings.Join(projectInfoParts, ", "))
	}

	systemMsg := strings.Join(systemMsgParts, "\n")
	return NewSessionOutput(true, systemMsg)
}

// checkLatestVersion fetches the latest godo release tag from GitHub.
// Returns empty string on failure or timeout.
func checkLatestVersion() string {
	// Try to read cached version first
	cached := readCachedVersion()
	if cached != "" {
		return cached
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/yejune/godo/releases/latest")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return ""
	}

	tag := strings.TrimPrefix(release.TagName, "v")
	if tag != "" {
		writeCachedVersion(tag)
	}
	return tag
}

// readCachedVersion reads a previously cached latest version from .do/.latest-version.
// Returns empty string if not available or stale (older than 24h).
func readCachedVersion() string {
	cacheFile := latestVersionCachePath()
	if cacheFile == "" {
		return ""
	}

	info, err := os.Stat(cacheFile)
	if err != nil {
		return ""
	}
	// Invalidate cache after 24 hours
	if time.Since(info.ModTime()) > 24*time.Hour {
		return ""
	}

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// writeCachedVersion persists the latest version to .do/.latest-version.
func writeCachedVersion(ver string) {
	cacheFile := latestVersionCachePath()
	if cacheFile == "" {
		return
	}
	_ = os.MkdirAll(filepath.Dir(cacheFile), 0755)
	_ = os.WriteFile(cacheFile, []byte(ver+"\n"), 0644)
}

// latestVersionCachePath returns the path for the version cache file.
func latestVersionCachePath() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Join(cwd, ".do", ".latest-version")
}

// detectProjectInfo inspects the given directory to determine project name,
// type, and primary language. Falls back to environment variables first.
func detectProjectInfo(cwd string) (name, projectType, lang string) {
	// Environment variable overrides
	if v := os.Getenv("DO_PROJECT_NAME"); v != "" {
		name = v
	}
	if v := os.Getenv("DO_PROJECT_TYPE"); v != "" {
		projectType = v
	}
	if v := os.Getenv("DO_PROJECT_LANG"); v != "" {
		lang = v
	}
	if name != "" && projectType != "" && lang != "" {
		return
	}

	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return
		}
	}

	// Detect from go.mod
	if name == "" || lang == "" {
		if data, err := os.ReadFile(filepath.Join(cwd, "go.mod")); err == nil {
			lines := strings.SplitN(string(data), "\n", 5)
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "module ") {
					mod := strings.TrimPrefix(line, "module ")
					mod = strings.TrimSpace(mod)
					if name == "" {
						// Use the last path segment as project name
						parts := strings.Split(mod, "/")
						name = parts[len(parts)-1]
					}
					if lang == "" {
						lang = "go"
					}
					if projectType == "" {
						projectType = "go"
					}
					break
				}
			}
		}
	}

	// Detect from package.json
	if name == "" || lang == "" {
		if data, err := os.ReadFile(filepath.Join(cwd, "package.json")); err == nil {
			var pkg struct {
				Name string `json:"name"`
			}
			if err := json.Unmarshal(data, &pkg); err == nil && pkg.Name != "" {
				if name == "" {
					name = pkg.Name
				}
			}
			if lang == "" {
				// Check for TypeScript
				if _, err := os.Stat(filepath.Join(cwd, "tsconfig.json")); err == nil {
					lang = "typescript"
				} else {
					lang = "javascript"
				}
			}
			if projectType == "" {
				projectType = "node"
			}
		}
	}

	// Detect from Cargo.toml
	if name == "" || lang == "" {
		if data, err := os.ReadFile(filepath.Join(cwd, "Cargo.toml")); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "name") && strings.Contains(line, "=") {
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						n := strings.Trim(strings.TrimSpace(parts[1]), `"`)
						if name == "" {
							name = n
						}
					}
					break
				}
			}
			if lang == "" {
				lang = "rust"
			}
			if projectType == "" {
				projectType = "rust"
			}
		}
	}

	// Detect from pyproject.toml or setup.py
	if name == "" || lang == "" {
		if _, err := os.Stat(filepath.Join(cwd, "pyproject.toml")); err == nil {
			if lang == "" {
				lang = "python"
			}
			if projectType == "" {
				projectType = "python"
			}
		} else if _, err := os.Stat(filepath.Join(cwd, "setup.py")); err == nil {
			if lang == "" {
				lang = "python"
			}
			if projectType == "" {
				projectType = "python"
			}
		}
	}

	// Use directory name as fallback for project name
	if name == "" {
		name = filepath.Base(cwd)
	}

	return
}
