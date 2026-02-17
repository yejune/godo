package hook

import (
	"fmt"
	"os"
	"strings"

	"github.com/yejune/godo/internal/mode"
	"github.com/yejune/godo/internal/persona"
)

// HandleCompact handles the PreCompact hook event.
// It preserves mode and checklist context across compaction.
func HandleCompact(input *Input) *Output {
	currentMode := mode.ReadState()
	userName := os.Getenv("DO_USER_NAME")
	personaType := os.Getenv("DO_PERSONA")
	if personaType == "" {
		personaType = "young-f"
	}

	var parts []string

	// Mode
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

	// Checklist status
	checklistPath := FindLatestChecklist()
	if checklistPath != "" {
		if stats, err := ParseChecklistFile(checklistPath); err == nil && stats.Total > 0 {
			parts = append(parts, fmt.Sprintf("체크리스트 상태: %s (파일: %s)", stats.Summary(), checklistPath))
		}
	}

	message := strings.Join(parts, "\n")
	return NewSessionOutput(true, message)
}
