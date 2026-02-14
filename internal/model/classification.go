package model

// ClassificationResult holds the per-section core/persona verdict for a document.
type ClassificationResult struct {
	DocPath   string                  `yaml:"doc_path"`
	Sections  []SectionClassification `yaml:"sections"`
	SkillRefs []SkillClassification   `yaml:"skill_refs"`
	PathRefs  []PathMatch             `yaml:"path_refs"`
}

// SectionClassification is the verdict for a single section.
type SectionClassification struct {
	Section    *Section `yaml:"-"`
	IsPersona  bool     `yaml:"is_persona"`
	Reason     string   `yaml:"reason"`     // e.g., "header matches TRUST5 pattern"
	SlotID     string   `yaml:"slot_id"`    // Template slot to replace with (if persona)
	Confidence float64  `yaml:"confidence"` // 0.0-1.0 detection confidence
}

// SkillClassification records a persona-specific skill detected in frontmatter.
type SkillClassification struct {
	SkillName string `yaml:"skill_name"`
	Category  string `yaml:"category"` // Why it is persona-specific
}

// PathMatch describes a hardcoded persona path found in content.
type PathMatch struct {
	Original string `yaml:"original"` // e.g., ".moai/specs/SPEC-{ID}/spec.md"
	SlotID   string `yaml:"slot_id"`  // e.g., "SPEC_PATH_PATTERN"
	Line     int    `yaml:"line"`
	Column   int    `yaml:"column"`
}
