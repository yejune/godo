package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ChecklistStats holds parsed checklist status counts.
type ChecklistStats struct {
	Total      int
	Pending    int // [ ]
	InProgress int // [~]
	Testing    int // [*]
	Blocked    int // [!]
	Done       int // [o]
	Failed     int // [x]
}

// HasIncomplete returns true if there are in-progress or blocked items.
func (s *ChecklistStats) HasIncomplete() bool {
	return s.InProgress > 0 || s.Blocked > 0
}

// Summary returns a one-line status summary.
func (s *ChecklistStats) Summary() string {
	parts := []string{}
	if s.Done > 0 {
		parts = append(parts, fmt.Sprintf("[o]%d", s.Done))
	}
	if s.InProgress > 0 {
		parts = append(parts, fmt.Sprintf("[~]%d", s.InProgress))
	}
	if s.Blocked > 0 {
		parts = append(parts, fmt.Sprintf("[!]%d", s.Blocked))
	}
	if s.Pending > 0 {
		parts = append(parts, fmt.Sprintf("[ ]%d", s.Pending))
	}
	if s.Testing > 0 {
		parts = append(parts, fmt.Sprintf("[*]%d", s.Testing))
	}
	if s.Failed > 0 {
		parts = append(parts, fmt.Sprintf("[x]%d", s.Failed))
	}
	if len(parts) == 0 {
		return "no items"
	}
	return strings.Join(parts, " ")
}

var checklistItemRe = regexp.MustCompile(`^\s*-\s*\[(.)\]`)

// ParseChecklistFile reads a checklist file and returns stats.
func ParseChecklistFile(path string) (*ChecklistStats, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseChecklistContent(string(data)), nil
}

// ParseChecklistContent parses checklist content and returns stats.
func ParseChecklistContent(content string) *ChecklistStats {
	stats := &ChecklistStats{}
	for _, line := range strings.Split(content, "\n") {
		m := checklistItemRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		stats.Total++
		switch m[1] {
		case " ":
			stats.Pending++
		case "~":
			stats.InProgress++
		case "*":
			stats.Testing++
		case "!":
			stats.Blocked++
		case "o":
			stats.Done++
		case "x":
			stats.Failed++
		}
	}
	return stats
}

// FindLatestChecklist finds the most recent checklist.md in .do/jobs/.
func FindLatestChecklist() string {
	jobsDir := ".do/jobs"
	var checklists []string

	filepath.Walk(jobsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.Name() == "checklist.md" && !info.IsDir() {
			checklists = append(checklists, path)
		}
		return nil
	})

	if len(checklists) == 0 {
		return ""
	}

	// Sort by path (date-based paths sort chronologically)
	sort.Strings(checklists)
	return checklists[len(checklists)-1]
}
