package hook

import "fmt"

// HandleStop handles the Stop hook event.
// It blocks stopping if there are in-progress or blocked checklist items.
func HandleStop(input *Input) *Output {
	checklistPath := FindLatestChecklist()
	if checklistPath == "" {
		// No checklist, allow stopping
		return &Output{}
	}

	stats, err := ParseChecklistFile(checklistPath)
	if err != nil || stats.Total == 0 {
		return &Output{}
	}

	if stats.HasIncomplete() {
		reason := fmt.Sprintf("체크리스트에 미완료 항목이 있습니다: %s\n파일: %s", stats.Summary(), checklistPath)
		return NewStopBlockOutput(reason)
	}

	return &Output{}
}
