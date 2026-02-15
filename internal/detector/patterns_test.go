package detector

import (
	"regexp"
	"testing"
)

func TestNewDefaultRegistry_PatternCounts(t *testing.T) {
	reg := NewDefaultRegistry()

	if got := len(reg.HeaderPatterns); got != 4 {
		t.Errorf("HeaderPatterns count = %d, want 4", got)
	}
	if got := len(reg.PathPatterns); got != 3 {
		t.Errorf("PathPatterns count = %d, want 3", got)
	}
	if got := len(reg.SkillPatterns); got != 6 {
		t.Errorf("SkillPatterns count = %d, want 6", got)
	}
	if got := len(reg.WholeFileAgents); got != 6 {
		t.Errorf("WholeFileAgents count = %d, want 6", got)
	}
	if got := len(reg.WholeFileSkills); got != 6 {
		t.Errorf("WholeFileSkills count = %d, want 6", got)
	}
	if got := len(reg.WholeFileRules); got != 2 {
		t.Errorf("WholeFileRules count = %d, want 2", got)
	}
}

func TestNewDefaultRegistry_HeaderPatternSlotIDs(t *testing.T) {
	reg := NewDefaultRegistry()

	expectedSlots := map[string]bool{
		"QUALITY_FRAMEWORK":          false,
		"QUALITY_SECURITY_FRAMEWORK": false,
		"TRACEABILITY_SYSTEM":        false,
	}

	for _, hp := range reg.HeaderPatterns {
		if _, ok := expectedSlots[hp.SlotID]; ok {
			expectedSlots[hp.SlotID] = true
		}
	}

	for slotID, found := range expectedSlots {
		if !found {
			t.Errorf("expected slot ID %q not found in HeaderPatterns", slotID)
		}
	}
}

func TestNewDefaultRegistry_PathPatternSlotIDs(t *testing.T) {
	reg := NewDefaultRegistry()

	expectedSlots := []string{"SPEC_PATH_PATTERN", "QUALITY_CONFIG_PATH", "DOCS_PATH_PATTERN"}
	for i, pp := range reg.PathPatterns {
		if pp.SlotID != expectedSlots[i] {
			t.Errorf("PathPatterns[%d].SlotID = %q, want %q", i, pp.SlotID, expectedSlots[i])
		}
	}
}

func TestNewDefaultRegistry_SkillPatternNames(t *testing.T) {
	reg := NewDefaultRegistry()

	expectedSkills := []string{
		"moai-foundation-core",
		"moai-foundation-quality",
		"moai-workflow-ddd",
		"moai-workflow-tdd",
		"moai-workflow-spec",
		"moai-workflow-project",
	}

	for i, sp := range reg.SkillPatterns {
		if sp.SkillName != expectedSkills[i] {
			t.Errorf("SkillPatterns[%d].SkillName = %q, want %q", i, sp.SkillName, expectedSkills[i])
		}
	}
}

func TestCompileHeaderPatterns_AllCompile(t *testing.T) {
	reg := NewDefaultRegistry()

	compiled, err := reg.CompileHeaderPatterns()
	if err != nil {
		t.Fatalf("CompileHeaderPatterns() returned error: %v", err)
	}

	if len(compiled) != len(reg.HeaderPatterns) {
		t.Errorf("compiled count = %d, want %d", len(compiled), len(reg.HeaderPatterns))
	}

	for i, re := range compiled {
		if re == nil {
			t.Errorf("compiled[%d] is nil for pattern %q", i, reg.HeaderPatterns[i].Pattern)
		}
	}
}

func TestCompileHeaderPatterns_MatchExamples(t *testing.T) {
	reg := NewDefaultRegistry()
	compiled, err := reg.CompileHeaderPatterns()
	if err != nil {
		t.Fatalf("CompileHeaderPatterns() error: %v", err)
	}

	tests := []struct {
		name    string
		header  string
		wantIdx int // index in compiled that should match, -1 for no match
	}{
		{"TRUST 5 Compliance", "TRUST 5 Compliance", 0},
		{"TRUST5 Validation", "TRUST5 Validation", 0},
		{"trust 5 framework", "trust 5 framework", 0},
		{"TRUST 5 Framework with caps", "TRUST 5 Framework", 0},
		{"Security & TRUST 5", "Security & TRUST 5", 1},
		{"security & trust 5", "security & trust 5", 1},
		{"TAG Chain", "TAG Chain", 2},
		{"tag chain case insensitive", "TAG Chain Integrity", 2},
		{"Research TAG System", "Research TAG System", 3},
		{"No match", "Some Other Section", -1},
		{"Partial TRUST mention", "Using TRUST principles", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchIdx := -1
			for i, re := range compiled {
				if re.MatchString(tt.header) {
					matchIdx = i
					break
				}
			}
			if matchIdx != tt.wantIdx {
				t.Errorf("header %q: matched pattern index = %d, want %d", tt.header, matchIdx, tt.wantIdx)
			}
		})
	}
}

func TestCompilePathPatterns_AllCompile(t *testing.T) {
	reg := NewDefaultRegistry()

	for i, pp := range reg.PathPatterns {
		_, err := regexp.Compile(pp.Pattern)
		if err != nil {
			t.Errorf("PathPatterns[%d] pattern %q failed to compile: %v", i, pp.Pattern, err)
		}
	}
}

func TestCompilePathPatterns_MatchExamples(t *testing.T) {
	reg := NewDefaultRegistry()

	tests := []struct {
		name    string
		input   string
		wantIdx int // index in PathPatterns that should match, -1 for no match
	}{
		{"SPEC path with ID", ".moai/specs/SPEC-AUTH/spec.md", 0},
		{"SPEC path with braces", ".moai/specs/SPEC-{ID}/spec.md", 0},
		{"SPEC path bare", ".moai/specs/SPEC-LOGIN/", 0},
		{"quality.yaml", ".moai/config/sections/quality.yaml", 1},
		{"docs path", ".moai/docs/", 2},
		{"docs subpath", ".moai/docs/api.md", 2},
		{"no match", "src/main.go", -1},
	}

	compiled := make([]*regexp.Regexp, len(reg.PathPatterns))
	for i, pp := range reg.PathPatterns {
		compiled[i] = regexp.MustCompile(pp.Pattern)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchIdx := -1
			for i, re := range compiled {
				if re.MatchString(tt.input) {
					matchIdx = i
					break
				}
			}
			if matchIdx != tt.wantIdx {
				t.Errorf("input %q: matched pattern index = %d, want %d", tt.input, matchIdx, tt.wantIdx)
			}
		})
	}
}

func TestIsWholeFilePersonaAgent(t *testing.T) {
	reg := NewDefaultRegistry()

	positives := []string{
		"manager-spec", "manager-ddd", "manager-tdd",
		"manager-project", "manager-quality", "team-quality",
	}
	for _, name := range positives {
		if !reg.IsWholeFilePersonaAgent(name) {
			t.Errorf("IsWholeFilePersonaAgent(%q) = false, want true", name)
		}
	}

	negatives := []string{
		"expert-backend", "expert-frontend", "manager-git",
		"team-backend-dev", "builder-agent", "",
	}
	for _, name := range negatives {
		if reg.IsWholeFilePersonaAgent(name) {
			t.Errorf("IsWholeFilePersonaAgent(%q) = true, want false", name)
		}
	}
}

func TestIsWholeFilePersonaSkill(t *testing.T) {
	reg := NewDefaultRegistry()

	positives := []string{
		"moai-foundation-core", "moai-foundation-quality",
		"moai-workflow-ddd", "moai-workflow-tdd",
		"moai-workflow-spec", "moai-workflow-project",
	}
	for _, name := range positives {
		if !reg.IsWholeFilePersonaSkill(name) {
			t.Errorf("IsWholeFilePersonaSkill(%q) = false, want true", name)
		}
	}

	negatives := []string{
		"do-foundation-claude", "custom-my-skill", "moai-unknown", "",
	}
	for _, name := range negatives {
		if reg.IsWholeFilePersonaSkill(name) {
			t.Errorf("IsWholeFilePersonaSkill(%q) = true, want false", name)
		}
	}
}

func TestIsWholeFilePersonaRule(t *testing.T) {
	reg := NewDefaultRegistry()

	positives := []string{"spec-workflow.md", "workflow-modes.md"}
	for _, name := range positives {
		if !reg.IsWholeFilePersonaRule(name) {
			t.Errorf("IsWholeFilePersonaRule(%q) = false, want true", name)
		}
	}

	negatives := []string{
		"moai-constitution.md", "coding-standards.md",
		"dev-environment.md", "",
	}
	for _, name := range negatives {
		if reg.IsWholeFilePersonaRule(name) {
			t.Errorf("IsWholeFilePersonaRule(%q) = true, want false", name)
		}
	}
}

func TestNewDefaultRegistry_ContentPatternCounts(t *testing.T) {
	reg := NewDefaultRegistry()

	if got := len(reg.ContentPatterns); got != 2 {
		t.Errorf("ContentPatterns count = %d, want 2", got)
	}
}

func TestNewDefaultRegistry_ContentPatternSlotIDs(t *testing.T) {
	reg := NewDefaultRegistry()

	expectedSlots := map[string]bool{
		"QUALITY_GATE_TEXT":       false,
		"QUALITY_PRINCIPLES_TEXT": false,
	}

	for _, cp := range reg.ContentPatterns {
		if _, ok := expectedSlots[cp.SlotID]; ok {
			expectedSlots[cp.SlotID] = true
		}
	}

	for slotID, found := range expectedSlots {
		if !found {
			t.Errorf("expected slot ID %q not found in ContentPatterns", slotID)
		}
	}
}

func TestCompileContentPatterns_AllCompile(t *testing.T) {
	reg := NewDefaultRegistry()

	compiled, err := reg.CompileContentPatterns()
	if err != nil {
		t.Fatalf("CompileContentPatterns() returned error: %v", err)
	}

	if len(compiled) != len(reg.ContentPatterns) {
		t.Errorf("compiled count = %d, want %d", len(compiled), len(reg.ContentPatterns))
	}

	for i, re := range compiled {
		if re == nil {
			t.Errorf("compiled[%d] is nil for pattern %q", i, reg.ContentPatterns[i].Pattern)
		}
	}
}

func TestCompileContentPatterns_MatchExamples(t *testing.T) {
	reg := NewDefaultRegistry()
	compiled, err := reg.CompileContentPatterns()
	if err != nil {
		t.Fatalf("CompileContentPatterns() error: %v", err)
	}

	tests := []struct {
		name    string
		input   string
		wantIdx int // index in compiled that should match, -1 for no match
	}{
		{"TRUST 5 quality gates", "TRUST 5 quality gates passed", 0},
		{"TRUST5 quality gates", "TRUST5 quality gates passed", 0},
		{"TRUST 5 principles", "Follow TRUST 5 principles", 1},
		{"no match on plain TRUST", "Using TRUST methodology", -1},
		{"no match on TRUST 5 alone", "TRUST 5 is great", -1},
		{"no match on quality alone", "quality gates are important", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchIdx := -1
			for i, re := range compiled {
				if re.MatchString(tt.input) {
					matchIdx = i
					break
				}
			}
			if matchIdx != tt.wantIdx {
				t.Errorf("input %q: matched pattern index = %d, want %d", tt.input, matchIdx, tt.wantIdx)
			}
		})
	}
}

func TestNewDefaultRegistry_PartialSkillPatternCounts(t *testing.T) {
	reg := NewDefaultRegistry()

	if got := len(reg.PartialSkillPatterns); got != 1 {
		t.Errorf("PartialSkillPatterns count = %d, want 1", got)
	}
}

func TestIsPartialSkill(t *testing.T) {
	reg := NewDefaultRegistry()

	if !reg.IsPartialSkill("moai-workflow-testing") {
		t.Error("IsPartialSkill(moai-workflow-testing) = false, want true")
	}

	negatives := []string{
		"moai-foundation-core",
		"moai-workflow-ddd",
		"do-foundation-claude",
		"unknown-skill",
		"",
	}
	for _, name := range negatives {
		if reg.IsPartialSkill(name) {
			t.Errorf("IsPartialSkill(%q) = true, want false", name)
		}
	}
}

func TestIsPartialPersonaModule(t *testing.T) {
	reg := NewDefaultRegistry()

	tests := []struct {
		name         string
		skillName    string
		moduleRel    string
		wantPersona  bool
	}{
		{
			name:        "ddd-context7 root module",
			skillName:   "moai-workflow-testing",
			moduleRel:   "modules/ddd-context7.md",
			wantPersona: true,
		},
		{
			name:        "ddd-context7 submodule",
			skillName:   "moai-workflow-testing",
			moduleRel:   "modules/ddd-context7/advanced-features.md",
			wantPersona: true,
		},
		{
			name:        "ddd core-classes submodule",
			skillName:   "moai-workflow-testing",
			moduleRel:   "modules/ddd/core-classes.md",
			wantPersona: true,
		},
		{
			name:        "core module - ai-debugging",
			skillName:   "moai-workflow-testing",
			moduleRel:   "modules/ai-debugging.md",
			wantPersona: false,
		},
		{
			name:        "core module - debugging subdir",
			skillName:   "moai-workflow-testing",
			moduleRel:   "modules/debugging/debugging-workflows.md",
			wantPersona: false,
		},
		{
			name:        "SKILL.md itself",
			skillName:   "moai-workflow-testing",
			moduleRel:   "SKILL.md",
			wantPersona: false,
		},
		{
			name:        "core module - performance",
			skillName:   "moai-workflow-testing",
			moduleRel:   "modules/performance/optimization-patterns.md",
			wantPersona: false,
		},
		{
			name:        "non-matching skill name",
			skillName:   "some-other-skill",
			moduleRel:   "modules/ddd/core-classes.md",
			wantPersona: false,
		},
		{
			name:        "empty module path",
			skillName:   "moai-workflow-testing",
			moduleRel:   "",
			wantPersona: false,
		},
		{
			name:        "scripts dir (non-module)",
			skillName:   "moai-workflow-testing",
			moduleRel:   "scripts/with_server.py",
			wantPersona: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reg.IsPartialPersonaModule(tt.skillName, tt.moduleRel)
			if got != tt.wantPersona {
				t.Errorf("IsPartialPersonaModule(%q, %q) = %v, want %v",
					tt.skillName, tt.moduleRel, got, tt.wantPersona)
			}
		})
	}
}
