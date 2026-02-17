package hook

import "fmt"

// HandleStop handles the Stop hook event.
// It blocks stopping if there are in-progress checklist items or uncommitted changes.
func HandleStop(input *Input) *Output {
	var reasons []string

	// Check 1: Checklist incomplete items
	checklistPath := FindLatestChecklist()
	if checklistPath != "" {
		stats, err := ParseChecklistFile(checklistPath)
		if err == nil && stats.Total > 0 && stats.HasIncomplete() {
			reasons = append(reasons, fmt.Sprintf("체크리스트에 미완료 항목이 있습니다: %s\n파일: %s", stats.Summary(), checklistPath))
		}
	}

	// Check 2: Uncommitted git changes
	hasChanges, summary := GitStatus()
	if hasChanges {
		reasons = append(reasons, fmt.Sprintf("미커밋 변경사항이 있습니다:\n%s", summary))
	}

	if len(reasons) > 0 {
		combined := ""
		for i, r := range reasons {
			if i > 0 {
				combined += "\n\n"
			}
			combined += r
		}
		return NewStopBlockOutput(combined)
	}

	return &Output{}
}
