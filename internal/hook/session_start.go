package hook

import (
	"fmt"
	"os"
	"strings"

	"github.com/yejune/godo/internal/mode"
	"github.com/yejune/godo/internal/persona"
)

// HandleSessionStart handles the SessionStart hook event.
// It loads the current mode and persona, returning a system message.
func HandleSessionStart(input *Input) *Output {
	currentMode := mode.ReadState()
	userName := os.Getenv("DO_USER_NAME")
	personaType := os.Getenv("DO_PERSONA")
	if personaType == "" {
		personaType = "young-f"
	}

	var parts []string

	// Load persona
	personaDir := persona.ResolveDir()
	if personaDir != "" {
		if pd, err := persona.LoadCharacter(personaDir, personaType); err == nil {
			honorific := pd.BuildHonorific(userName)
			if honorific != "" {
				parts = append(parts, fmt.Sprintf("Persona: %s (호칭: %s)", pd.Name, honorific))
			}
			if pd.FullContent != "" {
				parts = append(parts, pd.FullContent)
			}
		}
	}

	// Mode info
	modePrefix := strings.ToUpper(currentMode[:1]) + currentMode[1:]
	parts = append(parts, fmt.Sprintf("현재 실행 모드: %s (응답 접두사: [%s])", currentMode, modePrefix))

	message := strings.Join(parts, "\n\n")
	return NewSessionOutput(true, message)
}
