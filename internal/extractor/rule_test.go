package extractor

import (
	"strings"
	"testing"

	"github.com/do-focus/convert/internal/detector"
	"github.com/do-focus/convert/internal/model"
)

func TestRuleExtractor_WholeFilePersona(t *testing.T) {
	reg := detector.NewDefaultRegistry()
	det, err := detector.NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	ext := NewRuleExtractor(reg, det)

	doc := &model.Document{
		Path: "rules/moai/workflow/spec-workflow.md",
		Sections: []*model.Section{
			{Level: 1, Title: "SPEC Workflow", Content: "# SPEC Workflow\n\nDetails..."},
		},
	}

	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	if coreDoc != nil {
		t.Error("Expected nil coreDoc for whole-file persona rule")
	}
	if len(manifest.Rules) != 1 || manifest.Rules[0] != doc.Path {
		t.Errorf("manifest.Rules = %v, want [%q]", manifest.Rules, doc.Path)
	}
}

func TestRuleExtractor_CoreRulePassthrough(t *testing.T) {
	reg := detector.NewDefaultRegistry()
	det, err := detector.NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	ext := NewRuleExtractor(reg, det)

	doc := &model.Document{
		Path: "rules/dev-testing.md",
		Sections: []*model.Section{
			{Level: 1, Title: "Testing Rules", Content: "# Testing Rules\n\nFollow testing guidelines."},
		},
	}

	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	if coreDoc == nil {
		t.Error("Expected non-nil coreDoc for core rule")
	}
	if len(manifest.Rules) != 0 {
		t.Errorf("manifest.Rules should be empty for core rule, got %v", manifest.Rules)
	}
}

func TestRuleExtractor_InlineContentPatternSlotting(t *testing.T) {
	reg := detector.NewDefaultRegistry()
	det, err := detector.NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	ext := NewRuleExtractor(reg, det)

	doc := &model.Document{
		Path: "rules/moai/workflow/workflow-modes.md-like-core.md",
		Sections: []*model.Section{
			{
				Level:   2,
				Title:   "DDD Mode",
				Content: "## DDD Mode\n\nSuccess Criteria:\n- All SPEC requirements implemented\n- TRUST 5 quality gates passed\n- 85%+ code coverage achieved",
			},
		},
	}

	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	if coreDoc == nil {
		t.Fatal("Expected non-nil coreDoc")
	}

	// Check that the content now has the inline slot marker
	if !strings.Contains(coreDoc.Sections[0].Content, "{{slot:QUALITY_GATE_TEXT}}") {
		t.Errorf("Expected {{slot:QUALITY_GATE_TEXT}} in content, got:\n%s", coreDoc.Sections[0].Content)
	}

	// Check that original text is stored in manifest
	val, ok := manifest.SlotContent["QUALITY_GATE_TEXT"]
	if !ok {
		t.Fatal("manifest.SlotContent missing QUALITY_GATE_TEXT")
	}
	if val != "TRUST 5 quality gates" {
		t.Errorf("SlotContent[QUALITY_GATE_TEXT] = %q, want %q", val, "TRUST 5 quality gates")
	}

	// Check that "passed" is still there (only the pattern match was replaced)
	if !strings.Contains(coreDoc.Sections[0].Content, "passed") {
		t.Error("Expected 'passed' to remain in content after slotting")
	}
}

func TestRuleExtractor_MultipleContentPatterns(t *testing.T) {
	reg := detector.NewDefaultRegistry()
	det, err := detector.NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	ext := NewRuleExtractor(reg, det)

	doc := &model.Document{
		Path: "rules/language/typescript.md",
		Sections: []*model.Section{
			{
				Level:   2,
				Title:   "Quality",
				Content: "## Quality\n\n- TRUST 5 quality gates passed\n- Follow TRUST 5 principles in depth",
			},
		},
	}

	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	if coreDoc == nil {
		t.Fatal("Expected non-nil coreDoc")
	}

	if !strings.Contains(coreDoc.Sections[0].Content, "{{slot:QUALITY_GATE_TEXT}}") {
		t.Error("Expected {{slot:QUALITY_GATE_TEXT}} in content")
	}
	if !strings.Contains(coreDoc.Sections[0].Content, "{{slot:QUALITY_PRINCIPLES_TEXT}}") {
		t.Error("Expected {{slot:QUALITY_PRINCIPLES_TEXT}} in content")
	}

	if _, ok := manifest.SlotContent["QUALITY_GATE_TEXT"]; !ok {
		t.Error("manifest.SlotContent missing QUALITY_GATE_TEXT")
	}
	if _, ok := manifest.SlotContent["QUALITY_PRINCIPLES_TEXT"]; !ok {
		t.Error("manifest.SlotContent missing QUALITY_PRINCIPLES_TEXT")
	}
}

func TestRuleExtractor_NestedChildrenSlotting(t *testing.T) {
	reg := detector.NewDefaultRegistry()
	det, err := detector.NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	ext := NewRuleExtractor(reg, det)

	child := &model.Section{
		Level:   3,
		Title:   "Success Criteria",
		Content: "### Success Criteria\n\n- TRUST 5 quality gates passed",
	}

	doc := &model.Document{
		Path: "rules/dev-workflow.md",
		Sections: []*model.Section{
			{
				Level:    2,
				Title:    "DDD Mode",
				Content:  "## DDD Mode\n\nMethodology details.",
				Children: []*model.Section{child},
			},
		},
	}

	_, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	// The child section should have been slotted
	if !strings.Contains(child.Content, "{{slot:QUALITY_GATE_TEXT}}") {
		t.Errorf("Expected {{slot:QUALITY_GATE_TEXT}} in child content, got:\n%s", child.Content)
	}

	if _, ok := manifest.SlotContent["QUALITY_GATE_TEXT"]; !ok {
		t.Error("manifest.SlotContent missing QUALITY_GATE_TEXT for child section")
	}
}

func TestRuleExtractor_NilDetector_NoSlotting(t *testing.T) {
	reg := detector.NewDefaultRegistry()

	ext := NewRuleExtractor(reg, nil)

	doc := &model.Document{
		Path: "rules/dev-testing.md",
		Sections: []*model.Section{
			{
				Level:   2,
				Title:   "Testing",
				Content: "## Testing\n\n- TRUST 5 quality gates passed",
			},
		},
	}

	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	if coreDoc == nil {
		t.Fatal("Expected non-nil coreDoc")
	}

	// With nil detector, no slotting should happen
	if strings.Contains(coreDoc.Sections[0].Content, "{{slot:QUALITY_GATE_TEXT}}") {
		t.Error("Expected no slotting with nil detector")
	}

	if len(manifest.SlotContent) != 0 {
		t.Errorf("Expected empty SlotContent with nil detector, got %v", manifest.SlotContent)
	}
}

func TestRuleExtractor_NoContentMatch_Passthrough(t *testing.T) {
	reg := detector.NewDefaultRegistry()
	det, err := detector.NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}

	ext := NewRuleExtractor(reg, det)

	originalContent := "## Coding Rules\n\nFollow clean code principles.\nWrite readable code."
	doc := &model.Document{
		Path: "rules/coding-standards.md",
		Sections: []*model.Section{
			{Level: 2, Title: "Coding Rules", Content: originalContent},
		},
	}

	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	if coreDoc == nil {
		t.Fatal("Expected non-nil coreDoc")
	}

	// Content should be unchanged
	if coreDoc.Sections[0].Content != originalContent {
		t.Errorf("Content was modified when no patterns match:\ngot:  %q\nwant: %q",
			coreDoc.Sections[0].Content, originalContent)
	}

	if len(manifest.SlotContent) != 0 {
		t.Errorf("Expected empty SlotContent, got %v", manifest.SlotContent)
	}
}
