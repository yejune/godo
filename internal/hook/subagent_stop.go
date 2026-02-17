package hook

import "fmt"

// HandleSubagentStop handles the SubagentStop hook event.
// Checks for uncommitted changes when a subagent stops.
func HandleSubagentStop(input *Input) *Output {
	hasChanges, summary := GitStatus()
	if hasChanges {
		warning := fmt.Sprintf("[HARD] 서브에이전트 종료 시 미커밋 변경사항 감지. 오케스트레이터는 커밋을 확인하라.\n미커밋 파일:\n%s", summary)
		return &Output{
			Continue: true,
			HookSpecificOutput: &SpecificOutput{
				HookEventName:    "SubagentStop",
				AdditionalContext: warning,
			},
		}
	}
	return &Output{Continue: true}
}
