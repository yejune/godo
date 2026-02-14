package detector

import (
	"regexp"
)

// PatternRegistry holds all detection patterns organized by category.
// Patterns are matched against markdown headers, content text, and
// frontmatter fields to identify persona-specific content.
type PatternRegistry struct {
	HeaderPatterns  []HeaderPattern
	PathPatterns    []PathPattern
	SkillPatterns   []SkillPattern
	WholeFileAgents []string // agent names that are 100% persona
	WholeFileSkills []string // skill names that are 100% persona
	WholeFileRules  []string // rule files that are 100% persona
}

// HeaderPattern matches markdown section headers that indicate persona content.
// Detection is header-text based, not line-number based, for version resilience.
type HeaderPattern struct {
	Pattern     string // Regex pattern to match header text
	SlotID      string // Template slot to assign
	Category    string // "quality_framework", "spec_workflow", "methodology"
	Description string // What this pattern detects
}

// PathPattern matches hardcoded file paths in content text.
type PathPattern struct {
	Pattern     string // Regex for the path pattern
	SlotID      string // Template slot for replacement
	Replacement string // Template replacement: "{{SLOT_ID}}"
}

// SkillPattern identifies persona-specific skills in frontmatter.
type SkillPattern struct {
	SkillName string // Exact skill name to match
	Category  string // Why it is persona-specific
}

// NewDefaultRegistry creates a PatternRegistry pre-loaded with
// known moai-adk persona patterns.
func NewDefaultRegistry() *PatternRegistry {
	return &PatternRegistry{
		HeaderPatterns: []HeaderPattern{
			{
				Pattern:     `(?i)^###?\s+TRUST\s*5\s+(Compliance|Validation|Framework)`,
				SlotID:      "QUALITY_FRAMEWORK",
				Category:    "quality_framework",
				Description: "TRUST 5 quality compliance section in agent body",
			},
			{
				Pattern:     `(?i)^###?\s+Security\s+&\s+TRUST\s+5`,
				SlotID:      "QUALITY_SECURITY_FRAMEWORK",
				Category:    "quality_framework",
				Description: "Combined security and TRUST 5 section",
			},
			{
				Pattern:     `(?i)^###?\s+TAG\s+Chain`,
				SlotID:      "TRACEABILITY_SYSTEM",
				Category:    "methodology",
				Description: "TAG chain traceability system section",
			},
			{
				Pattern:     `(?i)^###?\s+Research\s+TAG\s+System`,
				SlotID:      "TRACEABILITY_SYSTEM",
				Category:    "methodology",
				Description: "Research TAG system integration",
			},
		},
		PathPatterns: []PathPattern{
			{
				Pattern:     `\.moai/specs/SPEC-\{?[A-Z0-9_]*\}?/?[a-z._]*`,
				SlotID:      "SPEC_PATH_PATTERN",
				Replacement: "{{SPEC_PATH_PATTERN}}",
			},
			{
				Pattern:     `\.moai/config/sections/quality\.yaml`,
				SlotID:      "QUALITY_CONFIG_PATH",
				Replacement: "{{QUALITY_CONFIG_PATH}}",
			},
			{
				Pattern:     `\.moai/docs/`,
				SlotID:      "DOCS_PATH_PATTERN",
				Replacement: "{{DOCS_PATH_PATTERN}}",
			},
		},
		SkillPatterns: []SkillPattern{
			{SkillName: "moai-foundation-core", Category: "TRUST5 + SPEC-First DDD"},
			{SkillName: "moai-foundation-quality", Category: "TRUST5 validation engine"},
			{SkillName: "moai-workflow-ddd", Category: "DDD methodology"},
			{SkillName: "moai-workflow-tdd", Category: "TDD methodology"},
			{SkillName: "moai-workflow-spec", Category: "SPEC workflow"},
			{SkillName: "moai-workflow-project", Category: "SPEC project init"},
		},
		WholeFileAgents: []string{
			"manager-spec",
			"manager-ddd",
			"manager-tdd",
			"manager-project",
			"manager-quality",
			"team-quality",
		},
		WholeFileSkills: []string{
			"moai-foundation-core",
			"moai-foundation-quality",
			"moai-workflow-ddd",
			"moai-workflow-tdd",
			"moai-workflow-spec",
			"moai-workflow-project",
		},
		WholeFileRules: []string{
			"spec-workflow.md",
			"workflow-modes.md",
		},
	}
}

// IsWholeFilePersonaAgent returns true if the given agent name is a
// whole-file persona agent that should be moved entirely to persona output.
func (r *PatternRegistry) IsWholeFilePersonaAgent(name string) bool {
	for _, a := range r.WholeFileAgents {
		if a == name {
			return true
		}
	}
	return false
}

// IsWholeFilePersonaSkill returns true if the given skill name is a
// whole-file persona skill that should be moved entirely to persona output.
func (r *PatternRegistry) IsWholeFilePersonaSkill(name string) bool {
	for _, s := range r.WholeFileSkills {
		if s == name {
			return true
		}
	}
	return false
}

// IsWholeFilePersonaRule returns true if the given rule filename is a
// whole-file persona rule that should be moved entirely to persona output.
func (r *PatternRegistry) IsWholeFilePersonaRule(filename string) bool {
	for _, ru := range r.WholeFileRules {
		if ru == filename {
			return true
		}
	}
	return false
}

// CompileHeaderPatterns compiles all header pattern regexes and returns them.
// Returns an error if any pattern fails to compile (RE2 incompatible).
func (r *PatternRegistry) CompileHeaderPatterns() ([]*regexp.Regexp, error) {
	compiled := make([]*regexp.Regexp, len(r.HeaderPatterns))
	for i, hp := range r.HeaderPatterns {
		re, err := regexp.Compile(hp.Pattern)
		if err != nil {
			return nil, err
		}
		compiled[i] = re
	}
	return compiled, nil
}
