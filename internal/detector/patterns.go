package detector

import (
	"regexp"
	"strings"
)

// ContentPattern matches inline text strings in document bodies.
// Unlike HeaderPattern (which matches section headers for section-level slotting),
// ContentPattern matches arbitrary text within body content for inline replacement
// with {{slot:SLOT_ID}} markers.
//
// Use case: "Follow TRUST 5 quality gates" in language rule files should become
// "Follow {{slot:QUALITY_GATE_TEXT}}" so a different persona can inject
// its own quality framework name.
type ContentPattern struct {
	Pattern     string // Regex pattern to match in body text
	SlotID      string // Inline slot ID to replace with: {{slot:SLOT_ID}}
	Category    string // "quality_framework", "methodology", etc.
	Description string // Human-readable description
}

// PatternRegistry holds all detection patterns organized by category.
// Patterns are matched against markdown headers, content text, and
// frontmatter fields to identify persona-specific content.
type PatternRegistry struct {
	HeaderPatterns       []HeaderPattern
	PathPatterns         []PathPattern
	SkillPatterns        []SkillPattern
	ContentPatterns      []ContentPattern      // Inline text patterns for body replacement
	PartialSkillPatterns []PartialSkillPattern  // Skills with module-level persona classification
	WholeFileAgents      []string               // agent names that are 100% persona
	WholeFileSkills      []string               // skill names that are 100% persona
	WholeFileSkillDirs   []string               // skill directory names where ALL contents are persona
	WholeFileRules       []string               // rule files that are 100% persona
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
	Replacement string // Template replacement: "{{slot:SLOT_ID}}"
}

// SkillPattern identifies persona-specific skills in frontmatter.
type SkillPattern struct {
	SkillName string // Exact skill name to match
	Category  string // Why it is persona-specific
}

// PartialSkillPattern identifies specific modules within a skill as persona.
// Unlike WholeFileSkills (entire skill is persona), this allows module-level
// granularity: some modules are persona, the rest remain core.
type PartialSkillPattern struct {
	SkillName      string   // Skill directory name (e.g., "moai-workflow-testing")
	PersonaModules []string // Path prefixes relative to skill dir (e.g., "modules/ddd")
	Category       string   // Why these modules are persona-specific
}

// NewDefaultRegistry creates a PatternRegistry pre-loaded with
// known moai-adk persona patterns.
func NewDefaultRegistry() *PatternRegistry {
	return &PatternRegistry{
		HeaderPatterns: []HeaderPattern{
			{
				Pattern:     `(?i)^TRUST\s*5\s+(Compliance|Validation|Framework)`,
				SlotID:      "QUALITY_FRAMEWORK",
				Category:    "quality_framework",
				Description: "TRUST 5 quality compliance section in agent body",
			},
			{
				Pattern:     `(?i)^Security\s+&\s+TRUST\s+5`,
				SlotID:      "QUALITY_SECURITY_FRAMEWORK",
				Category:    "quality_framework",
				Description: "Combined security and TRUST 5 section",
			},
			{
				Pattern:     `(?i)^TAG\s+Chain`,
				SlotID:      "TRACEABILITY_SYSTEM",
				Category:    "methodology",
				Description: "TAG chain traceability system section",
			},
			{
				Pattern:     `(?i)^Research\s+TAG\s+System`,
				SlotID:      "TRACEABILITY_SYSTEM",
				Category:    "methodology",
				Description: "Research TAG system integration",
			},
		},
		PathPatterns: []PathPattern{
			{
				Pattern:     `\.moai/specs/SPEC-\{?[A-Z0-9_]*\}?/?[a-z._]*`,
				SlotID:      "SPEC_PATH_PATTERN",
				Replacement: "{{slot:SPEC_PATH_PATTERN}}",
			},
			{
				Pattern:     `\.moai/config/sections/quality\.yaml`,
				SlotID:      "QUALITY_CONFIG_PATH",
				Replacement: "{{slot:QUALITY_CONFIG_PATH}}",
			},
			{
				Pattern:     `\.moai/docs/`,
				SlotID:      "DOCS_PATH_PATTERN",
				Replacement: "{{slot:DOCS_PATH_PATTERN}}",
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
		// Inline content patterns for text replacement in rule bodies
		ContentPatterns: []ContentPattern{
			{
				Pattern:     `TRUST\s*5\s+quality\s+gates`,
				SlotID:      "QUALITY_GATE_TEXT",
				Category:    "quality_framework",
				Description: "TRUST 5 quality gates reference in rule/workflow text",
			},
			{
				Pattern:     `TRUST\s*5\s+principles`,
				SlotID:      "QUALITY_PRINCIPLES_TEXT",
				Category:    "quality_framework",
				Description: "TRUST 5 principles reference in output style text",
			},
		},
		PartialSkillPatterns: []PartialSkillPattern{
			{
				SkillName:      "moai-workflow-testing",
				PersonaModules: []string{"modules/ddd"},
				Category:       "DDD methodology modules within testing skill",
			},
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
		WholeFileSkillDirs: []string{
			"moai", // skills/moai/ and all subdirectories (workflows/, references/, etc.)
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

// IsWholeFilePersonaSkillDir returns true if the given skill directory name
// is a whole-directory persona skill dir. All files under skills/<dirName>/
// (including subdirectories) should be moved entirely to persona output.
//
// Matching rules:
//   - Exact match: dirName == entry (e.g., "moai" matches "moai")
//   - Prefix match: dirName starts with entry + "-" (e.g., "moai-workflow-project" matches "moai")
//
// This ensures that skills like "moai-workflow-project", "moai-docs-generation",
// etc. are all classified as persona when "moai" is in WholeFileSkillDirs.
func (r *PatternRegistry) IsWholeFilePersonaSkillDir(dirName string) bool {
	for _, d := range r.WholeFileSkillDirs {
		if d == dirName {
			return true
		}
		// Prefix match: "moai" matches "moai-workflow-project", "moai-docs-generation", etc.
		if strings.HasPrefix(dirName, d+"-") {
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

// IsPartialSkill returns true if the given skill name has partial persona
// module patterns defined (i.e., some modules are persona, others are core).
func (r *PatternRegistry) IsPartialSkill(skillName string) bool {
	for _, p := range r.PartialSkillPatterns {
		if p.SkillName == skillName {
			return true
		}
	}
	return false
}

// IsPartialPersonaModule returns true if the given file path (relative to the
// skill directory) matches a persona module pattern for the given skill.
// The moduleRelPath should be relative to the skill directory, e.g.,
// "modules/ddd-context7/advanced-features.md" or "modules/ddd/core-classes.md".
func (r *PatternRegistry) IsPartialPersonaModule(skillName, moduleRelPath string) bool {
	for _, p := range r.PartialSkillPatterns {
		if p.SkillName != skillName {
			continue
		}
		for _, prefix := range p.PersonaModules {
			if strings.HasPrefix(moduleRelPath, prefix) {
				return true
			}
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

// CompileContentPatterns compiles all content pattern regexes and returns them.
// Returns an error if any pattern fails to compile.
func (r *PatternRegistry) CompileContentPatterns() ([]*regexp.Regexp, error) {
	compiled := make([]*regexp.Regexp, len(r.ContentPatterns))
	for i, cp := range r.ContentPatterns {
		re, err := regexp.Compile(cp.Pattern)
		if err != nil {
			return nil, err
		}
		compiled[i] = re
	}
	return compiled, nil
}
