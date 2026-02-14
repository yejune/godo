package model

// Slot represents a template placeholder that replaces persona-specific content.
// The slot ID is globally unique within a persona and maps to a content file.
type Slot struct {
	ID          string `yaml:"id"`          // e.g., "QUALITY_FRAMEWORK", "SPEC_PATH_PATTERN"
	Category    string `yaml:"category"`    // "section", "skill_ref", "path_pattern", "rule_block"
	Description string `yaml:"description"` // Human-readable description of what this slot provides
	Default     string `yaml:"default"`     // Optional default content when no persona fills it
}

// Slot marker syntax used in templatized files.
//
// Inline example:
//
//	Read specification files from {{SPEC_PATH_PATTERN}}
//
// Section example:
//
//	<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->
//	<!-- END_SLOT:QUALITY_FRAMEWORK -->
const (
	// InlineSlotPattern matches {{SLOT_ID}} placeholders in content.
	InlineSlotPattern = `\{\{([A-Z][A-Z0-9_]*)\}\}`

	// SectionSlotBegin is the format string for section slot begin markers.
	SectionSlotBegin = "<!-- BEGIN_SLOT:%s -->"

	// SectionSlotEnd is the format string for section slot end markers.
	SectionSlotEnd = "<!-- END_SLOT:%s -->"
)
