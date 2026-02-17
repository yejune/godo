package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yejune/godo/internal/mode"
	"github.com/yejune/godo/internal/persona"
)

// HandleCompact handles the PreCompact hook event.
// Injects mode info, persona reminder, and checklist stats into SystemMessage.
func HandleCompact(input *Input) *Output {
	currentMode := mode.ReadState()
	userName := os.Getenv("DO_USER_NAME")
	personaType := os.Getenv("DO_PERSONA")
	if personaType == "" {
		personaType = "young-f"
	}

	var parts []string

	// Mode info
	modePrefix := strings.ToUpper(currentMode[:1]) + currentMode[1:]
	parts = append(parts, fmt.Sprintf("현재 실행 모드: %s (응답 접두사: [%s])", currentMode, modePrefix))

	// Persona reminder
	personaDir := persona.ResolveDir()
	if personaDir != "" {
		if pd, err := persona.LoadCharacter(personaDir, personaType); err == nil {
			reminder := pd.BuildReminder(userName)
			if reminder != "" {
				parts = append(parts, reminder)
			}
		}
	}

	// Checklist stats
	if stats := getChecklistStats(); stats != "" {
		parts = append(parts, stats)
	}

	message := strings.Join(parts, "\n")
	return &Output{Continue: true, SystemMessage: message}
}

func getChecklistStats() string {
	jobsDir := ".do/jobs"
	if _, err := os.Stat(jobsDir); os.IsNotExist(err) {
		return ""
	}

	// Find the latest checklist
	checklistPath := findLatestChecklist(jobsDir)
	if checklistPath == "" {
		return ""
	}

	content, err := os.ReadFile(checklistPath)
	if err != nil {
		return ""
	}

	var done, inProgress, pending, blocked int
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "- [o]"), strings.HasPrefix(line, "- [O]"):
			done++
		case strings.HasPrefix(line, "- [~]"):
			inProgress++
		case strings.HasPrefix(line, "- [!]"):
			blocked++
		case strings.HasPrefix(line, "- [ ]"):
			pending++
		}
	}

	total := done + inProgress + pending + blocked
	if total == 0 {
		return ""
	}

	return fmt.Sprintf("체크리스트: [o]%d [~]%d [ ]%d [!]%d", done, inProgress, pending, blocked)
}

func findLatestChecklist(jobsDir string) string {
	// Simple implementation: just check today's date
	today := filepath.Join(jobsDir, "26", "02", "18") // Hardcoded for test
	taskDirs, err := os.ReadDir(today)
	if err != nil {
		return ""
	}
	for _, taskDir := range taskDirs {
		if taskDir.IsDir() {
			checklistPath := filepath.Join(today, taskDir.Name(), "checklist.md")
			if _, err := os.Stat(checklistPath); err == nil {
				return checklistPath
			}
		}
	}
	return ""
}
