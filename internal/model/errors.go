package model

import "fmt"

// ErrParse indicates a markdown/frontmatter parsing failure.
type ErrParse struct {
	File    string
	Line    int
	Message string
}

func (e *ErrParse) Error() string {
	return fmt.Sprintf("parse error in %s at line %d: %s", e.File, e.Line, e.Message)
}

// ErrDetection indicates a pattern detection issue.
type ErrDetection struct {
	File       string
	Pattern    string
	Message    string
	Confidence float64
}

func (e *ErrDetection) Error() string {
	return fmt.Sprintf("detection error in %s (pattern: %s, confidence: %.2f): %s",
		e.File, e.Pattern, e.Confidence, e.Message)
}

// ErrSlot indicates a template slot resolution failure.
type ErrSlot struct {
	SlotID  string
	File    string
	Message string
}

func (e *ErrSlot) Error() string {
	return fmt.Sprintf("slot error for %s in %s: %s", e.SlotID, e.File, e.Message)
}

// ErrAssembly indicates a merge/assembly failure.
type ErrAssembly struct {
	Phase   string // "copy", "fill_slots", "patch_agent", "merge"
	File    string
	Message string
}

func (e *ErrAssembly) Error() string {
	return fmt.Sprintf("assembly error during %s for %s: %s", e.Phase, e.File, e.Message)
}
