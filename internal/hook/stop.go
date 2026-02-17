package hook

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// HandleStop handles the Stop hook event.
// Match original behavior:
// 1) If stop_hook_active is true, allow stop to prevent loops.
// 2) Block if there are uncommitted changes.
// 3) Block when the most recent active checklist is in-progress/blocked.
func HandleStop(input *Input) *Output {
	if input != nil && input.StopHookActive {
		return &Output{}
	}

	// Check for uncommitted changes
	if hasChanges, summary := GitStatus(); hasChanges {
		return NewStopBlockOutput(fmt.Sprintf("커밋되지 않은 변경사항이 있습니다:\n%s\n변경사항을 커밋한 후 종료하세요.", summary))
	}

	if reason := checkActiveChecklist(); reason != "" {
		return NewStopBlockOutput(reason)
	}
	return &Output{}
}

func checkActiveChecklist() string {
	jobsDir := ".do/jobs"
	if _, err := os.Stat(jobsDir); os.IsNotExist(err) {
		return ""
	}

	yearDirs, err := os.ReadDir(jobsDir)
	if err != nil || len(yearDirs) == 0 {
		return ""
	}
	sort.Slice(yearDirs, func(i, j int) bool { return yearDirs[i].Name() > yearDirs[j].Name() })

	for _, yearDir := range yearDirs {
		if !yearDir.IsDir() || !stopIsDigits(yearDir.Name()) {
			continue
		}
		monthPath := filepath.Join(jobsDir, yearDir.Name())
		monthDirs, err := os.ReadDir(monthPath)
		if err != nil {
			continue
		}
		sort.Slice(monthDirs, func(i, j int) bool { return monthDirs[i].Name() > monthDirs[j].Name() })
		for _, monthDir := range monthDirs {
			if !monthDir.IsDir() || !stopIsDigits(monthDir.Name()) {
				continue
			}
			dayPath := filepath.Join(monthPath, monthDir.Name())
			dayDirs, err := os.ReadDir(dayPath)
			if err != nil {
				continue
			}
			sort.Slice(dayDirs, func(i, j int) bool { return dayDirs[i].Name() > dayDirs[j].Name() })
			for _, dayDir := range dayDirs {
				if !dayDir.IsDir() || !stopIsDigits(dayDir.Name()) {
					continue
				}
				taskRoot := filepath.Join(dayPath, dayDir.Name())
				taskDirs, err := os.ReadDir(taskRoot)
				if err != nil {
					continue
				}
				sort.Slice(taskDirs, func(i, j int) bool { return taskDirs[i].Name() > taskDirs[j].Name() })
				for _, taskDir := range taskDirs {
					if !taskDir.IsDir() {
						continue
					}
					checklistPath := filepath.Join(taskRoot, taskDir.Name(), "checklist.md")
					if summary := parseChecklistSummary(checklistPath); summary != "" {
						return summary
					}
				}
				break
			}
			break
		}
		break
	}
	return ""
}

func parseChecklistSummary(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	var total, done, inProgress, pending, blocked int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "- [o]"), strings.HasPrefix(line, "- [O]"):
			done++
			total++
		case strings.HasPrefix(line, "- [~]"):
			inProgress++
			total++
		case strings.HasPrefix(line, "- [!]"):
			blocked++
			total++
		case strings.HasPrefix(line, "- [*]"):
			total++
		case strings.HasPrefix(line, "- [ ]"):
			pending++
			total++
		}
	}
	if total == 0 || done == total {
		return ""
	}
	if inProgress == 0 && blocked == 0 {
		return ""
	}

	return fmt.Sprintf(
		"활성 체크리스트가 있습니다 (%d/%d 완료, %d 진행중, %d 대기, %d 블로커). 체크리스트 파일(%s)을 읽고 현재 상태를 사용자에게 표시한 뒤 종료하세요.",
		done, total, inProgress, pending, blocked, path,
	)
}

func stopIsDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
