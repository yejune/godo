package detector

import (
	"testing"

	"github.com/do-focus/convert/internal/model"
)

func TestNewPersonaDetector_Success(t *testing.T) {
	reg := NewDefaultRegistry()
	det, err := NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}
	if det == nil {
		t.Fatal("NewPersonaDetector() returned nil")
	}
	if len(det.compiled) != len(reg.HeaderPatterns) {
		t.Errorf("compiled patterns count = %d, want %d", len(det.compiled), len(reg.HeaderPatterns))
	}
}

func TestNewPersonaDetector_InvalidPattern(t *testing.T) {
	reg := &PatternRegistry{
		HeaderPatterns: []HeaderPattern{
			{Pattern: "[invalid", SlotID: "TEST", Category: "test", Description: "bad regex"},
		},
	}
	_, err := NewPersonaDetector(reg)
	if err == nil {
		t.Fatal("NewPersonaDetector() with invalid regex should return error")
	}
}

func TestClassify_TRUST5Section_DetectedAsPersona(t *testing.T) {
	reg := NewDefaultRegistry()
	det, err := NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	doc := &model.Document{
		Path: "agents/moai/expert-backend.md",
		Sections: []*model.Section{
			{
				Level:     2,
				Title:     "Implementation Guidelines",
				Content:   "## Implementation Guidelines\n\nFollow clean code principles.",
				StartLine: 1,
				EndLine:   3,
			},
			{
				Level:     3,
				Title:     "TRUST 5 Compliance",
				Content:   "### TRUST 5 Compliance\n\n- Tested: 85%+ coverage\n- Readable: Clear naming\n- Unified: Consistent style\n- Secured: OWASP compliance\n- Trackable: Conventional commits",
				StartLine: 4,
				EndLine:   10,
			},
		},
	}

	result := det.Classify(doc)

	if result.DocPath != doc.Path {
		t.Errorf("DocPath = %q, want %q", result.DocPath, doc.Path)
	}

	if len(result.Sections) != 2 {
		t.Fatalf("Sections count = %d, want 2", len(result.Sections))
	}

	// First section should be core
	if result.Sections[0].IsPersona {
		t.Error("Section 0 (Implementation Guidelines) should be core, got persona")
	}

	// Second section should be persona with QUALITY_FRAMEWORK slot
	sec1 := result.Sections[1]
	if !sec1.IsPersona {
		t.Error("Section 1 (TRUST 5 Compliance) should be persona, got core")
	}
	if sec1.SlotID != "QUALITY_FRAMEWORK" {
		t.Errorf("Section 1 SlotID = %q, want %q", sec1.SlotID, "QUALITY_FRAMEWORK")
	}
	if sec1.Confidence != 1.0 {
		t.Errorf("Section 1 Confidence = %f, want 1.0", sec1.Confidence)
	}
}

func TestClassify_NoPersonaContent_AllCore(t *testing.T) {
	reg := NewDefaultRegistry()
	det, err := NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	doc := &model.Document{
		Path: "agents/expert-generic.md",
		Sections: []*model.Section{
			{
				Level:     2,
				Title:     "Overview",
				Content:   "## Overview\n\nGeneric agent overview.",
				StartLine: 1,
				EndLine:   3,
			},
			{
				Level:     2,
				Title:     "Guidelines",
				Content:   "## Guidelines\n\nFollow best practices.",
				StartLine: 4,
				EndLine:   6,
			},
		},
	}

	result := det.Classify(doc)

	for i, sc := range result.Sections {
		if sc.IsPersona {
			t.Errorf("Section %d (%q) should be core, got persona", i, sc.Section.Title)
		}
	}

	if len(result.SkillRefs) != 0 {
		t.Errorf("SkillRefs count = %d, want 0", len(result.SkillRefs))
	}
	if len(result.PathRefs) != 0 {
		t.Errorf("PathRefs count = %d, want 0", len(result.PathRefs))
	}
}

func TestClassify_NestedChildren_PersonaDetected(t *testing.T) {
	reg := NewDefaultRegistry()
	det, err := NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	child := &model.Section{
		Level:     3,
		Title:     "TAG Chain",
		Content:   "### TAG Chain\n\nTraceability system.",
		StartLine: 4,
		EndLine:   6,
	}

	doc := &model.Document{
		Path: "agents/moai/expert-backend.md",
		Sections: []*model.Section{
			{
				Level:     2,
				Title:     "Methodology",
				Content:   "## Methodology\n\nMethodology section.",
				StartLine: 1,
				EndLine:   6,
				Children:  []*model.Section{child},
			},
		},
	}

	result := det.Classify(doc)

	// Should find the TAG Chain child section as persona
	foundTagChain := false
	for _, sc := range result.Sections {
		if sc.Section == child && sc.IsPersona {
			foundTagChain = true
			if sc.SlotID != "TRACEABILITY_SYSTEM" {
				t.Errorf("TAG Chain SlotID = %q, want %q", sc.SlotID, "TRACEABILITY_SYSTEM")
			}
			if sc.Confidence != 1.0 {
				t.Errorf("TAG Chain Confidence = %f, want 1.0", sc.Confidence)
			}
		}
	}
	if !foundTagChain {
		t.Error("TAG Chain child section was not detected as persona")
	}
}

func TestDetectSkillRefs_PersonaSkillsDetected(t *testing.T) {
	reg := NewDefaultRegistry()
	det, err := NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	fm := &model.Frontmatter{
		Skills: []string{
			"moai-foundation-core",
			"do-foundation-claude",
			"moai-workflow-ddd",
		},
	}

	refs := det.DetectSkillRefs(fm)

	if len(refs) != 2 {
		t.Fatalf("DetectSkillRefs count = %d, want 2", len(refs))
	}

	if refs[0].SkillName != "moai-foundation-core" {
		t.Errorf("refs[0].SkillName = %q, want %q", refs[0].SkillName, "moai-foundation-core")
	}
	if refs[0].Category != "TRUST5 + SPEC-First DDD" {
		t.Errorf("refs[0].Category = %q, want %q", refs[0].Category, "TRUST5 + SPEC-First DDD")
	}

	if refs[1].SkillName != "moai-workflow-ddd" {
		t.Errorf("refs[1].SkillName = %q, want %q", refs[1].SkillName, "moai-workflow-ddd")
	}
}

func TestDetectSkillRefs_NoPersonaSkills_Empty(t *testing.T) {
	reg := NewDefaultRegistry()
	det, err := NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	fm := &model.Frontmatter{
		Skills: []string{
			"do-foundation-claude",
			"custom-my-skill",
		},
	}

	refs := det.DetectSkillRefs(fm)

	if len(refs) != 0 {
		t.Errorf("DetectSkillRefs count = %d, want 0", len(refs))
	}
}

func TestDetectSkillRefs_NilFrontmatter(t *testing.T) {
	reg := NewDefaultRegistry()
	det, err := NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	refs := det.DetectSkillRefs(nil)

	if len(refs) != 0 {
		t.Errorf("DetectSkillRefs(nil) count = %d, want 0", len(refs))
	}
}

func TestDetectPathPatterns_MoaiPathDetected(t *testing.T) {
	reg := NewDefaultRegistry()
	det, err := NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	content := "Some intro text here.\n\nRead specification files from .moai/specs/SPEC-001/spec.md for details.\n\nCheck quality config at .moai/config/sections/quality.yaml before proceeding.\n"

	matches := det.DetectPathPatterns(content)

	if len(matches) != 2 {
		t.Fatalf("DetectPathPatterns count = %d, want 2", len(matches))
	}

	if matches[0].SlotID != "SPEC_PATH_PATTERN" {
		t.Errorf("matches[0].SlotID = %q, want %q", matches[0].SlotID, "SPEC_PATH_PATTERN")
	}
	if matches[0].Line != 3 {
		t.Errorf("matches[0].Line = %d, want 3", matches[0].Line)
	}
	if matches[0].Column < 1 {
		t.Errorf("matches[0].Column = %d, want >= 1", matches[0].Column)
	}

	if matches[1].SlotID != "QUALITY_CONFIG_PATH" {
		t.Errorf("matches[1].SlotID = %q, want %q", matches[1].SlotID, "QUALITY_CONFIG_PATH")
	}
	if matches[1].Line != 5 {
		t.Errorf("matches[1].Line = %d, want 5", matches[1].Line)
	}
}

func TestDetectPathPatterns_NoMoaiPaths_Empty(t *testing.T) {
	reg := NewDefaultRegistry()
	det, err := NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	content := "This is a normal document.\n\nIt references src/main.go and package.json but no .moai paths.\n"

	matches := det.DetectPathPatterns(content)

	if len(matches) != 0 {
		t.Errorf("DetectPathPatterns count = %d, want 0", len(matches))
	}
}

func TestClassify_FullDocument_AllDetected(t *testing.T) {
	reg := NewDefaultRegistry()
	det, err := NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	doc := &model.Document{
		Path: "agents/moai/expert-backend.md",
		Frontmatter: &model.Frontmatter{
			Name:   "expert-backend",
			Skills: []string{"moai-foundation-core", "do-foundation-claude", "moai-workflow-tdd"},
		},
		Sections: []*model.Section{
			{
				Level:     2,
				Title:     "Backend Development",
				Content:   "## Backend Development\n\nImplement server-side features.\nRead specs from .moai/specs/SPEC-{ID}/spec.md",
				StartLine: 1,
				EndLine:   4,
			},
			{
				Level:     3,
				Title:     "TRUST 5 Compliance",
				Content:   "### TRUST 5 Compliance\n\n- Tested: 85%+\n- Readable: Clean code\n- Unified: Consistent\n- Secured: OWASP\n- Trackable: Commits",
				StartLine: 5,
				EndLine:   12,
			},
			{
				Level:     3,
				Title:     "Error Handling",
				Content:   "### Error Handling\n\nHandle errors properly.",
				StartLine: 13,
				EndLine:   15,
			},
		},
	}

	result := det.Classify(doc)

	if result.DocPath != "agents/moai/expert-backend.md" {
		t.Errorf("DocPath = %q, want %q", result.DocPath, "agents/moai/expert-backend.md")
	}

	if len(result.Sections) != 3 {
		t.Fatalf("Sections count = %d, want 3", len(result.Sections))
	}

	if result.Sections[0].IsPersona {
		t.Error("Section 0 (Backend Development) should be core")
	}
	if !result.Sections[1].IsPersona {
		t.Error("Section 1 (TRUST 5 Compliance) should be persona")
	}
	if result.Sections[1].SlotID != "QUALITY_FRAMEWORK" {
		t.Errorf("Section 1 SlotID = %q, want QUALITY_FRAMEWORK", result.Sections[1].SlotID)
	}
	if result.Sections[2].IsPersona {
		t.Error("Section 2 (Error Handling) should be core")
	}

	if len(result.SkillRefs) != 2 {
		t.Fatalf("SkillRefs count = %d, want 2", len(result.SkillRefs))
	}

	skillNames := map[string]bool{}
	for _, sr := range result.SkillRefs {
		skillNames[sr.SkillName] = true
	}
	if !skillNames["moai-foundation-core"] {
		t.Error("moai-foundation-core not detected in SkillRefs")
	}
	if !skillNames["moai-workflow-tdd"] {
		t.Error("moai-workflow-tdd not detected in SkillRefs")
	}

	if len(result.PathRefs) < 1 {
		t.Fatalf("PathRefs count = %d, want >= 1", len(result.PathRefs))
	}

	foundSpecPath := false
	for _, pr := range result.PathRefs {
		if pr.SlotID == "SPEC_PATH_PATTERN" {
			foundSpecPath = true
		}
	}
	if !foundSpecPath {
		t.Error("SPEC_PATH_PATTERN not detected in PathRefs")
	}
}

func TestClassify_NilFrontmatter_NoSkillRefs(t *testing.T) {
	reg := NewDefaultRegistry()
	det, err := NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	doc := &model.Document{
		Path: "rules/dev-testing.md",
		Sections: []*model.Section{
			{
				Level:     2,
				Title:     "Testing Rules",
				Content:   "## Testing Rules\n\nFollow testing guidelines.",
				StartLine: 1,
				EndLine:   3,
			},
		},
	}

	result := det.Classify(doc)

	if len(result.SkillRefs) != 0 {
		t.Errorf("SkillRefs count = %d, want 0 for nil frontmatter", len(result.SkillRefs))
	}
}

func TestDetectPathPatterns_LineColumnAccuracy(t *testing.T) {
	reg := NewDefaultRegistry()
	det, err := NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	content := "line 1\nline 2\npath: .moai/docs/api.md here\nline 4"

	matches := det.DetectPathPatterns(content)

	if len(matches) != 1 {
		t.Fatalf("DetectPathPatterns count = %d, want 1", len(matches))
	}

	if matches[0].Line != 3 {
		t.Errorf("Line = %d, want 3", matches[0].Line)
	}
	// ".moai/docs/" starts at column 7 (1-based: "path: " is 6 chars)
	if matches[0].Column != 7 {
		t.Errorf("Column = %d, want 7", matches[0].Column)
	}
	if matches[0].SlotID != "DOCS_PATH_PATTERN" {
		t.Errorf("SlotID = %q, want DOCS_PATH_PATTERN", matches[0].SlotID)
	}
}
