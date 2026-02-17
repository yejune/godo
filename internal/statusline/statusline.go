package statusline

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// ANSI color codes.
const (
	AnsiGreen  = "\033[32m"
	AnsiYellow = "\033[33m"
	AnsiRed    = "\033[31m"
	AnsiReset  = "\033[0m"
)

// Input represents the actual JSON that Claude Code sends via stdin.
type Input struct {
	Model struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"model"`
	ContextWindow struct {
		ContextWindowSize   int  `json:"context_window_size"`
		UsedPercentage      *int `json:"used_percentage"`
		RemainingPercentage *int `json:"remaining_percentage"`
		TotalInputTokens    int  `json:"total_input_tokens"`
		TotalOutputTokens   int  `json:"total_output_tokens"`
	} `json:"context_window"`
	Cost struct {
		TotalCostUSD      float64 `json:"total_cost_usd"`
		TotalDurationMS   int     `json:"total_duration_ms"`
		TotalLinesAdded   int     `json:"total_lines_added"`
		TotalLinesRemoved int     `json:"total_lines_removed"`
	} `json:"cost"`
	Agent *struct {
		Name string `json:"name"`
	} `json:"agent"`
}

// Config holds external dependencies for rendering the status line.
type Config struct {
	Version        string
	ReadModeState  func() string
	GetProfileName func() string
}

// Render reads JSON from stdin and prints the formatted status line.
func Render(cfg Config) {
	var input Input
	decoder := json.NewDecoder(os.Stdin)
	if err := decoder.Decode(&input); err != nil {
		fmt.Print("[Do]")
		return
	}

	mode := "do"
	if cfg.ReadModeState != nil {
		mode = cfg.ReadModeState()
	}
	var modePrefix string
	switch strings.ToLower(mode) {
	case "focus":
		modePrefix = "[Focus]"
	case "do":
		modePrefix = "[Do]"
	case "team":
		modePrefix = "[Team]"
	case "auto":
		modePrefix = "[Auto]"
	default:
		modePrefix = "[Do]"
	}

	modelShort := ShortenModel(input.Model.DisplayName)
	if modelShort == "" {
		modelShort = ShortenModel(input.Model.ID)
	}

	ctxPercent := 0
	if input.ContextWindow.UsedPercentage != nil {
		ctxPercent = *input.ContextWindow.UsedPercentage
	} else if input.ContextWindow.RemainingPercentage != nil {
		ctxPercent = 100 - *input.ContextWindow.RemainingPercentage
	}
	if ctxPercent > 100 {
		ctxPercent = 100
	}
	ctxStr := ColorizeContext(ctxPercent)

	branch, changes := GetGitInfo()
	costStr := FormatCost(input.Cost.TotalCostUSD)

	if icon := PersonaIcon(os.Getenv("DO_PERSONA")); icon != "" {
		modePrefix += icon
	}

	parts := []string{modePrefix}

	if cfg.GetProfileName != nil {
		if profile := cfg.GetProfileName(); profile != "default" {
			parts = append(parts, "🤖"+profile)
		}
	}

	if modelShort != "" {
		parts = append(parts, modelShort)
	}

	if input.Agent != nil && input.Agent.Name != "" {
		parts = append(parts, "🤖"+input.Agent.Name)
	}

	ctxPart := ctxStr
	if input.Cost.TotalDurationMS > 0 {
		mins := input.Cost.TotalDurationMS / 60000
		if mins >= 60 {
			ctxPart += fmt.Sprintf(" ⏰%dh%dm", mins/60, mins%60)
		} else if mins > 0 {
			ctxPart += fmt.Sprintf(" ⏰%dm", mins)
		}
	}
	parts = append(parts, ctxPart)

	parts = append(parts, TildeDir(getCwd()))

	if branch != "" {
		gitPart := branch
		if changes > 0 {
			gitPart += fmt.Sprintf(" +%d", changes)
		}
		parts = append(parts, gitPart)
	}

	if costStr != "" {
		parts = append(parts, costStr)
	}

	if os.Getenv("CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS") == "1" {
		parts = append(parts, "👥")
	}

	version := cfg.Version
	verStr := strings.TrimPrefix(version, "v")
	if latest := ReadLatestVersion(); latest != "" && version != "dev" && IsNewer(latest, version) {
		verStr += "🆙"
	}
	parts = append(parts, verStr)

	fmt.Print(strings.Join(parts, " | "))
}

// ShortenModel extracts a short model identifier.
func ShortenModel(model string) string {
	if model == "" {
		return ""
	}

	lower := strings.ToLower(model)

	switch {
	case strings.Contains(lower, "opus"):
		return "opus"
	case strings.Contains(lower, "sonnet"):
		return "sonnet"
	case strings.Contains(lower, "haiku"):
		return "haiku"
	case strings.Contains(lower, "gpt-4"):
		return "gpt4"
	case strings.Contains(lower, "gpt-3"):
		return "gpt3"
	default:
		if len(model) > 12 {
			return model[:12]
		}
		return model
	}
}

// ColorizeContext returns a colored context percentage string.
func ColorizeContext(percent int) string {
	label := fmt.Sprintf("used:%d%%", percent)

	switch {
	case percent >= 80:
		return AnsiRed + label + AnsiReset
	case percent >= 50:
		return AnsiYellow + label + AnsiReset
	default:
		return AnsiGreen + label + AnsiReset
	}
}

// GetGitInfo returns the current branch name and count of uncommitted changes.
func GetGitInfo() (string, int) {
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, err := branchCmd.Output()
	if err != nil {
		return "", 0
	}
	branch := strings.TrimSpace(string(branchOut))

	statusCmd := exec.Command("git", "status", "--porcelain")
	statusOut, err := statusCmd.Output()
	if err != nil {
		return branch, 0
	}

	changes := 0
	lines := strings.Split(strings.TrimSpace(string(statusOut)), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			changes++
		}
	}

	return branch, changes
}

// ReadLatestVersion reads the cached latest version from .do/.latest-version.
func ReadLatestVersion() string {
	projectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	if projectDir == "" {
		projectDir, _ = os.Getwd()
	}
	data, err := os.ReadFile(filepath.Join(projectDir, ".do", ".latest-version"))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// IsNewer returns true if latest version is newer than current.
func IsNewer(latest, current string) bool {
	parse := func(v string) []int {
		v = strings.TrimPrefix(v, "v")
		parts := strings.Split(v, ".")
		nums := make([]int, len(parts))
		for i, p := range parts {
			n, _ := strconv.Atoi(p)
			nums[i] = n
		}
		return nums
	}
	l, c := parse(latest), parse(current)
	for i := 0; i < len(l) && i < len(c); i++ {
		if l[i] > c[i] {
			return true
		}
		if l[i] < c[i] {
			return false
		}
	}
	return false
}

// TildeDir replaces home directory prefix with ~.
func TildeDir(dir string) string {
	if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(dir, home) {
		return "~" + dir[len(home):]
	}
	return dir
}

// PersonaIcon returns an emoji for the persona type.
func PersonaIcon(persona string) string {
	switch persona {
	case "young-f":
		return "🦋"
	case "young-m":
		return "🔥"
	case "senior-f":
		return "👩‍💻"
	case "senior-m":
		return "🧑‍💼"
	default:
		return ""
	}
}

// FormatCost formats a cost value as a dollar string.
func FormatCost(cost float64) string {
	if cost <= 0 {
		return ""
	}
	if cost < 0.01 {
		return "$" + strconv.FormatFloat(cost, 'f', 4, 64)
	}
	return "$" + strconv.FormatFloat(cost, 'f', 2, 64)
}

func getCwd() string {
	dir, _ := os.Getwd()
	return dir
}
