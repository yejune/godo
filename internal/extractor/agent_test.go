package extractor

import (
	"testing"

	"github.com/do-focus/convert/internal/detector"
	"github.com/do-focus/convert/internal/model"
)

// newTestDetector creates a PersonaDetector with the default registry for tests.
func newTestDetector(t *testing.T) *detector.PersonaDetector {
	t.Helper()
	reg := detector.NewDefaultRegistry()
	det, err := detector.NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector() error: %v", err)
	}
	return det
}

// newTestRegistry returns the default pattern registry for tests.
func newTestRegistry(t *testing.T) *detector.PatternRegistry {
	t.Helper()
	return detector.NewDefaultRegistry()
}

func TestExtractAgent_PureCoreAgent_Passthrough(t *testing.T) {
	det := newTestDetector(t)
	reg := newTestRegistry(t)

	doc := &model.Document{
		Path: "agents/expert-generic.md",
		Frontmatter: &model.Frontmatter{
			Name:        "expert-generic",
			Description: "A generic expert agent",
			Tools:       "Read Write Edit Grep Glob Bash",
			Skills:      []string{"do-foundation-claude"},
		},
		Sections: []*model.Section{
			{
				Level:     2,
				Title:     "## Overview",
				Content:   "## Overview\n\nGeneric agent overview.",
				StartLine: 1,
				EndLine:   3,
			},
			{
				Level:     2,
				Title:     "## Guidelines",
				Content:   "## Guidelines\n\nFollow best practices.",
				StartLine: 4,
				EndLine:   6,
			},
		},
	}

	ext := NewAgentExtractor(det, reg)
	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	// Core document should have all sections intact
	if len(coreDoc.Sections) != 2 {
		t.Errorf("core sections count = %d, want 2", len(coreDoc.Sections))
	}
	if coreDoc.Path != doc.Path {
		t.Errorf("core doc path = %q, want %q", coreDoc.Path, doc.Path)
	}

	// Frontmatter should be preserved
	if coreDoc.Frontmatter == nil {
		t.Fatal("core doc frontmatter is nil")
	}
	if coreDoc.Frontmatter.Name != "expert-generic" {
		t.Errorf("core doc frontmatter name = %q, want %q", coreDoc.Frontmatter.Name, "expert-generic")
	}

	// Manifest should have no slot content and no agent patches
	if len(manifest.SlotContent) != 0 {
		t.Errorf("manifest slot content count = %d, want 0", len(manifest.SlotContent))
	}
	if len(manifest.AgentPatches) != 0 {
		t.Errorf("manifest agent patches count = %d, want 0", len(manifest.AgentPatches))
	}
}

func TestExtractAgent_WithPersonaSections_Extraction(t *testing.T) {
	det := newTestDetector(t)
	reg := newTestRegistry(t)

	doc := &model.Document{
		Path: "agents/moai/expert-backend.md",
		Frontmatter: &model.Frontmatter{
			Name:        "expert-backend",
			Description: "Backend development expert",
			Skills:      []string{"do-foundation-claude"},
		},
		Sections: []*model.Section{
			{
				Level:     2,
				Title:     "## Implementation Guidelines",
				Content:   "## Implementation Guidelines\n\nFollow clean code principles.",
				StartLine: 1,
				EndLine:   3,
			},
			{
				Level:     3,
				Title:     "### TRUST 5 Compliance",
				Content:   "### TRUST 5 Compliance\n\n- Tested: 85%+ coverage\n- Readable: Clear naming\n- Unified: Consistent style\n- Secured: OWASP compliance\n- Trackable: Conventional commits",
				StartLine: 4,
				EndLine:   10,
			},
			{
				Level:     2,
				Title:     "## Error Handling",
				Content:   "## Error Handling\n\nHandle errors properly.",
				StartLine: 11,
				EndLine:   13,
			},
		},
	}

	ext := NewAgentExtractor(det, reg)
	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	// Core document should have 3 sections but the persona one replaced with slot marker
	if len(coreDoc.Sections) != 3 {
		t.Fatalf("core sections count = %d, want 3", len(coreDoc.Sections))
	}

	// Section 0: core (unchanged)
	if coreDoc.Sections[0].Title != "## Implementation Guidelines" {
		t.Errorf("section 0 title = %q, want '## Implementation Guidelines'", coreDoc.Sections[0].Title)
	}

	// Section 1: persona section should have slot marker in content
	sec1Content := coreDoc.Sections[1].Content
	if sec1Content == doc.Sections[1].Content {
		t.Error("section 1 content should be replaced with slot marker, but is unchanged")
	}
	// Should contain BEGIN_SLOT and END_SLOT markers
	if !containsSlotMarker(sec1Content, "QUALITY_FRAMEWORK") {
		t.Errorf("section 1 content should contain QUALITY_FRAMEWORK slot markers, got:\n%s", sec1Content)
	}

	// Section 2: core (unchanged)
	if coreDoc.Sections[2].Title != "## Error Handling" {
		t.Errorf("section 2 title = %q, want '## Error Handling'", coreDoc.Sections[2].Title)
	}

	// Manifest should have the persona content stored
	if len(manifest.SlotContent) == 0 {
		t.Fatal("manifest slot content is empty, expected QUALITY_FRAMEWORK entry")
	}
	if _, ok := manifest.SlotContent["QUALITY_FRAMEWORK"]; !ok {
		t.Error("manifest slot content missing QUALITY_FRAMEWORK key")
	}
}

func TestExtractAgent_MixedCorePesonaSections_SlotMarkersInserted(t *testing.T) {
	det := newTestDetector(t)
	reg := newTestRegistry(t)

	doc := &model.Document{
		Path: "agents/moai/expert-backend.md",
		Frontmatter: &model.Frontmatter{
			Name:        "expert-backend",
			Description: "Backend expert",
			Skills:      []string{"do-foundation-claude", "moai-foundation-core"},
		},
		Sections: []*model.Section{
			{
				Level:     2,
				Title:     "## Core Guidelines",
				Content:   "## Core Guidelines\n\nAlways follow these.",
				StartLine: 1,
				EndLine:   3,
			},
			{
				Level:     3,
				Title:     "### TRUST 5 Validation",
				Content:   "### TRUST 5 Validation\n\nValidation rules here.",
				StartLine: 4,
				EndLine:   6,
			},
			{
				Level:     3,
				Title:     "### TAG Chain",
				Content:   "### TAG Chain\n\nTraceability entries.",
				StartLine: 7,
				EndLine:   9,
			},
			{
				Level:     2,
				Title:     "## Deployment",
				Content:   "## Deployment\n\nDeploy instructions.",
				StartLine: 10,
				EndLine:   12,
			},
		},
	}

	ext := NewAgentExtractor(det, reg)
	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	// All 4 sections should remain, persona ones with slot markers
	if len(coreDoc.Sections) != 4 {
		t.Fatalf("core sections count = %d, want 4", len(coreDoc.Sections))
	}

	// Section 0: core
	if containsSlotMarker(coreDoc.Sections[0].Content, "") {
		t.Error("section 0 (Core Guidelines) should not have slot markers")
	}

	// Section 1: persona with QUALITY_FRAMEWORK slot
	if !containsSlotMarker(coreDoc.Sections[1].Content, "QUALITY_FRAMEWORK") {
		t.Error("section 1 (TRUST 5 Validation) should have QUALITY_FRAMEWORK slot marker")
	}

	// Section 2: persona with TRACEABILITY_SYSTEM slot
	if !containsSlotMarker(coreDoc.Sections[2].Content, "TRACEABILITY_SYSTEM") {
		t.Error("section 2 (TAG Chain) should have TRACEABILITY_SYSTEM slot marker")
	}

	// Section 3: core
	if containsSlotMarker(coreDoc.Sections[3].Content, "") {
		t.Error("section 3 (Deployment) should not have slot markers")
	}

	// Manifest should have both slots
	if len(manifest.SlotContent) < 2 {
		t.Errorf("manifest slot content count = %d, want >= 2", len(manifest.SlotContent))
	}
	if _, ok := manifest.SlotContent["QUALITY_FRAMEWORK"]; !ok {
		t.Error("manifest slot content missing QUALITY_FRAMEWORK")
	}
	if _, ok := manifest.SlotContent["TRACEABILITY_SYSTEM"]; !ok {
		t.Error("manifest slot content missing TRACEABILITY_SYSTEM")
	}

	// Persona skills should be extracted into agent patches
	if manifest.AgentPatches == nil {
		t.Fatal("manifest agent patches is nil")
	}
	patch, ok := manifest.AgentPatches["expert-backend"]
	if !ok {
		t.Fatal("manifest agent patches missing 'expert-backend'")
	}
	foundMoaiSkill := false
	for _, s := range patch.AppendSkills {
		if s == "moai-foundation-core" {
			foundMoaiSkill = true
		}
	}
	if !foundMoaiSkill {
		t.Errorf("agent patch AppendSkills should contain 'moai-foundation-core', got %v", patch.AppendSkills)
	}

	// Core frontmatter should NOT contain persona skills
	for _, s := range coreDoc.Frontmatter.Skills {
		if s == "moai-foundation-core" {
			t.Error("core frontmatter skills should not contain persona skill 'moai-foundation-core'")
		}
	}
}

func TestExtractAgent_WholeFilePersonaAgent_EntireFileInManifest(t *testing.T) {
	det := newTestDetector(t)
	reg := newTestRegistry(t)

	doc := &model.Document{
		Path: "agents/moai/manager-spec.md",
		Frontmatter: &model.Frontmatter{
			Name:        "manager-spec",
			Description: "SPEC document creation manager",
			Skills:      []string{"moai-workflow-spec"},
		},
		Sections: []*model.Section{
			{
				Level:     2,
				Title:     "## SPEC Management",
				Content:   "## SPEC Management\n\nCreate and manage SPEC documents.",
				StartLine: 1,
				EndLine:   3,
			},
		},
		RawContent: "---\nname: manager-spec\ndescription: SPEC document creation manager\nskills:\n  - moai-workflow-spec\n---\n\n## SPEC Management\n\nCreate and manage SPEC documents.\n",
	}

	ext := NewAgentExtractor(det, reg)
	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	// Core document should be nil for whole-file persona agents
	if coreDoc != nil {
		t.Error("core doc should be nil for whole-file persona agent, got non-nil")
	}

	// Manifest should list this agent in Agents
	foundAgent := false
	for _, a := range manifest.Agents {
		if a == "agents/moai/manager-spec.md" {
			foundAgent = true
		}
	}
	if !foundAgent {
		t.Errorf("manifest Agents should contain the agent path, got %v", manifest.Agents)
	}
}

func TestExtractAgent_NilFrontmatter_NoError(t *testing.T) {
	det := newTestDetector(t)
	reg := newTestRegistry(t)

	doc := &model.Document{
		Path: "agents/simple.md",
		Sections: []*model.Section{
			{
				Level:     2,
				Title:     "## Simple Agent",
				Content:   "## Simple Agent\n\nDoes simple things.",
				StartLine: 1,
				EndLine:   3,
			},
		},
	}

	ext := NewAgentExtractor(det, reg)
	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	if coreDoc == nil {
		t.Fatal("core doc should not be nil for non-persona agent")
	}
	if len(coreDoc.Sections) != 1 {
		t.Errorf("core sections count = %d, want 1", len(coreDoc.Sections))
	}
	if len(manifest.SlotContent) != 0 {
		t.Errorf("manifest slot content should be empty, got %d", len(manifest.SlotContent))
	}
}

func TestExtractAgent_PersonaSkillsExtractedFromFrontmatter(t *testing.T) {
	det := newTestDetector(t)
	reg := newTestRegistry(t)

	doc := &model.Document{
		Path: "agents/moai/expert-backend.md",
		Frontmatter: &model.Frontmatter{
			Name:   "expert-backend",
			Skills: []string{"do-foundation-claude", "moai-foundation-core", "moai-workflow-ddd", "custom-skill"},
		},
		Sections: []*model.Section{
			{
				Level:     2,
				Title:     "## Backend",
				Content:   "## Backend\n\nBackend stuff.",
				StartLine: 1,
				EndLine:   3,
			},
		},
	}

	ext := NewAgentExtractor(det, reg)
	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	// Core frontmatter should only contain non-persona skills
	coreSkills := coreDoc.Frontmatter.Skills
	for _, s := range coreSkills {
		if s == "moai-foundation-core" || s == "moai-workflow-ddd" {
			t.Errorf("core frontmatter skills should not contain persona skill %q", s)
		}
	}
	if len(coreSkills) != 2 {
		t.Errorf("core skills count = %d, want 2 (do-foundation-claude, custom-skill); got %v", len(coreSkills), coreSkills)
	}

	// Manifest agent patch should list extracted persona skills
	patch, ok := manifest.AgentPatches["expert-backend"]
	if !ok {
		t.Fatal("manifest agent patches missing 'expert-backend'")
	}
	if len(patch.AppendSkills) != 2 {
		t.Errorf("patch AppendSkills count = %d, want 2; got %v", len(patch.AppendSkills), patch.AppendSkills)
	}
}

func TestExtractAgent_EmptyDocument(t *testing.T) {
	det := newTestDetector(t)
	reg := newTestRegistry(t)

	doc := &model.Document{
		Path:     "agents/empty.md",
		Sections: []*model.Section{},
	}

	ext := NewAgentExtractor(det, reg)
	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	if coreDoc == nil {
		t.Fatal("core doc should not be nil for empty non-persona doc")
	}
	if len(coreDoc.Sections) != 0 {
		t.Errorf("core sections count = %d, want 0", len(coreDoc.Sections))
	}
	if len(manifest.SlotContent) != 0 {
		t.Errorf("manifest slot content should be empty, got %d", len(manifest.SlotContent))
	}
}

// containsSlotMarker checks if content contains a BEGIN_SLOT marker.
// If slotID is empty, checks for any slot marker.
func containsSlotMarker(content string, slotID string) bool {
	if slotID == "" {
		return len(content) > 0 &&
			(contains(content, "<!-- BEGIN_SLOT:") || contains(content, "{{"))
	}
	return contains(content, "<!-- BEGIN_SLOT:"+slotID+" -->")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
