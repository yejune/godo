package hook

import (
	"fmt"
	"os"
	"strings"

	"github.com/yejune/godo/internal/persona"
)

// HandlePostToolUse handles the PostToolUse hook event.
// It injects persona reminders and checks for uncommitted changes after Task tool.
func HandlePostToolUse(input *Input) *Output {
	userName := os.Getenv("DO_USER_NAME")
	personaType := os.Getenv("DO_PERSONA")
	if personaType == "" {
		personaType = "young-f"
	}

	var contextParts []string

	// Persona reminder
	personaDir := persona.ResolveDir()
	if personaDir != "" {
		if pd, err := persona.LoadCharacter(personaDir, personaType); err == nil {
			reminder := pd.BuildReminder(userName)
			if reminder != "" {
				contextParts = append(contextParts, reminder)
			}
		}
	}

	// Git status check after Task tool (agent completed)
	if input.ToolName == "Task" {
		hasChanges, summary := GitStatus()
		if hasChanges {
			warning := fmt.Sprintf("[HARD] 에이전트가 코드를 수정했으나 커밋하지 않았다. 에이전트는 자신의 변경사항을 반드시 커밋해야 한다. 미커밋 변경사항을 확인하고 커밋하라.\n미커밋 파일:\n%s", summary)
			contextParts = append(contextParts, warning)
		}
	}

	if len(contextParts) == 0 {
		return &Output{}
	}

	return NewPostToolOutput(strings.Join(contextParts, "\n"))
}
