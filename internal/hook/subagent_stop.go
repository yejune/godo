package hook

import (
	"fmt"
)

// HandleSubagentStop handles the SubagentStop hook event.
// Warns about uncommitted changes but allows the subagent to stop.
func HandleSubagentStop(input *Input) *Output {
	// Check for uncommitted changes
	if hasChanges, summary := GitStatus(); hasChanges {
		return &Output{
			Continue: true,
			HookSpecificOutput: &SpecificOutput{
				HookEventName:     "SubagentStop",
				AdditionalContext: fmt.Sprintf("경고: 커밋되지 않은 변경사항이 있습니다:\n%s", summary),
			},
		}
	}
	return &Output{Continue: true}
}
