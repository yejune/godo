package hook

import (
	"os"

	"github.com/yejune/godo/internal/persona"
)

// HandlePostToolUse handles the PostToolUse hook event.
// It injects persona reminders as additionalContext.
func HandlePostToolUse(input *Input) *Output {
	userName := os.Getenv("DO_USER_NAME")
	personaType := os.Getenv("DO_PERSONA")
	if personaType == "" {
		personaType = "young-f"
	}

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

	return NewPostToolOutput(reminder)
}
