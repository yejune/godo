package hook

import (
	"fmt"
	"os"
	"strings"

	"github.com/yejune/godo/internal/mode"
	"github.com/yejune/godo/internal/persona"
)

// HandleSessionStart handles the SessionStart hook event.
// Injects mode info and persona content into SystemMessage.
func HandleSessionStart(input *Input) *Output {
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

	// Persona content
	personaDir := persona.ResolveDir()
	if personaDir != "" {
		if pd, err := persona.LoadCharacter(personaDir, personaType); err == nil {
			// Add persona name
			if pd.Name != "" {
				parts = append(parts, pd.Name)
			}
			reminder := pd.BuildReminder(userName)
			if reminder != "" {
				parts = append(parts, reminder)
			}
			if pd.FullContent != "" {
				parts = append(parts, pd.FullContent)
			}
		}
	}

	message := strings.Join(parts, "\n")
	return NewSessionOutput(true, message)
}
