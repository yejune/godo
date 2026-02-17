package hook

import (
	"fmt"
	"os"
	"strings"

	"github.com/yejune/godo/internal/mode"
	"github.com/yejune/godo/internal/persona"
)

// HandleUserPromptSubmit handles the UserPromptSubmit hook event.
// It injects mode and persona reminders as additionalContext.
func HandleUserPromptSubmit(input *Input) *Output {
	currentMode := mode.ReadState()
	userName := os.Getenv("DO_USER_NAME")
	personaType := os.Getenv("DO_PERSONA")
	if personaType == "" {
		personaType = "young-f"
	}

	var parts []string

	// Mode reminder
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

	return &Output{
		HookSpecificOutput: &SpecificOutput{
			HookEventName:    "UserPromptSubmit",
			AdditionalContext: strings.Join(parts, "\n"),
		},
	}
}
