package hook

import (
	"fmt"
	"os"
	"strings"

	"github.com/yejune/godo/internal/mode"
	"github.com/yejune/godo/internal/persona"
)

// HandlePostToolUse handles the PostToolUse hook event.
// Injects mode and persona reminders as additionalContext when persona is available.
func HandlePostToolUse(input *Input) *Output {
	currentMode := mode.ReadState()
	userName := os.Getenv("DO_USER_NAME")
	personaType := os.Getenv("DO_PERSONA")
	if personaType == "" {
		personaType = "young-f"
	}

	// Check if persona is available
	personaDir := persona.ResolveDir()
	if personaDir == "" {
		return &Output{}
	}

	pd, err := persona.LoadCharacter(personaDir, personaType)
	if err != nil {
		return &Output{}
	}

	reminder := pd.BuildReminder(userName)
	if reminder == "" {
		return &Output{}
	}

	var parts []string

	// Mode reminder
	modePrefix := strings.ToUpper(currentMode[:1]) + currentMode[1:]
	parts = append(parts, fmt.Sprintf("현재 실행 모드: %s (응답 접두사: [%s])", currentMode, modePrefix))

	// Persona reminder
	parts = append(parts, reminder)

	return &Output{
		HookSpecificOutput: &SpecificOutput{
			HookEventName:     "PostToolUse",
			AdditionalContext: strings.Join(parts, "\n"),
		},
	}
}
