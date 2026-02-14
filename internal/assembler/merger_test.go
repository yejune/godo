package assembler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/do-focus/convert/internal/model"
	"github.com/do-focus/convert/internal/template"
)

func TestMergeFile_SlotFilling(t *testing.T) {
	reg := newTestRegistry(map[string]*template.SlotEntry{
		"QUALITY_FRAMEWORK": {
			Category:   "section",
			MarkerType: "section",
			Default:    "default quality content",
		},
	})

	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{
			"QUALITY_FRAMEWORK": "TRUST 5 quality gates enforced",
		},
	}

	// Set up core dir with a template file.
	coreDir := t.TempDir()
	coreFile := filepath.Join(coreDir, "rules", "quality.md")
	if err := os.MkdirAll(filepath.Dir(coreFile), 0o755); err != nil {
		t.Fatal(err)
	}
	coreContent := "# Quality\n<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->\nplaceholder\n<!-- END_SLOT:QUALITY_FRAMEWORK -->\n# End"
	if err := os.WriteFile(coreFile, []byte(coreContent), 0o644); err != nil {
		t.Fatal(err)
	}

	outputDir := t.TempDir()
	personaDir := t.TempDir()

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	result, err := m.MergeFile("rules/quality.md")
	if err != nil {
		t.Fatalf("MergeFile error: %v", err)
	}

	if result.SlotsResolved != 1 {
		t.Errorf("expected 1 slot resolved, got %d", result.SlotsResolved)
	}

	// Verify output file exists and has filled content.
	outPath := filepath.Join(outputDir, "rules", "quality.md")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}

	output := string(data)
	if !strings.Contains(output, "TRUST 5 quality gates enforced") {
		t.Errorf("expected filled slot content in output, got:\n%s", output)
	}
	if strings.Contains(output, "placeholder") {
		t.Errorf("expected placeholder to be replaced, got:\n%s", output)
	}
}

func TestMergeFile_MultipleSlots(t *testing.T) {
	reg := newTestRegistry(map[string]*template.SlotEntry{
		"QUALITY_FRAMEWORK": {
			Category:   "section",
			MarkerType: "section",
			Default:    "default quality",
		},
		"TOOL_NAME": {
			Category:   "path_pattern",
			MarkerType: "inline",
			Default:    "moai",
		},
		"SPEC_PATH": {
			Category:   "path_pattern",
			MarkerType: "inline",
			Default:    ".moai/specs/",
		},
	})

	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{
			"QUALITY_FRAMEWORK": "TRUST 5 gates",
			"TOOL_NAME":         "godo",
			"SPEC_PATH":         ".do/specs/",
		},
	}

	coreDir := t.TempDir()
	coreFile := filepath.Join(coreDir, "agents", "expert.md")
	if err := os.MkdirAll(filepath.Dir(coreFile), 0o755); err != nil {
		t.Fatal(err)
	}
	coreContent := `# Expert Agent
Use {{TOOL_NAME}} to read specs at {{SPEC_PATH}} path.

<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->
old quality
<!-- END_SLOT:QUALITY_FRAMEWORK -->

Run {{TOOL_NAME}} now.`
	if err := os.WriteFile(coreFile, []byte(coreContent), 0o644); err != nil {
		t.Fatal(err)
	}

	outputDir := t.TempDir()
	personaDir := t.TempDir()

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	result, err := m.MergeFile("agents/expert.md")
	if err != nil {
		t.Fatalf("MergeFile error: %v", err)
	}

	if result.SlotsResolved != 3 {
		t.Errorf("expected 3 slots resolved, got %d", result.SlotsResolved)
	}

	outPath := filepath.Join(outputDir, "agents", "expert.md")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}

	output := string(data)
	if !strings.Contains(output, "godo") {
		t.Errorf("expected TOOL_NAME filled with godo, got:\n%s", output)
	}
	if !strings.Contains(output, ".do/specs/") {
		t.Errorf("expected SPEC_PATH filled with .do/specs/, got:\n%s", output)
	}
	if !strings.Contains(output, "TRUST 5 gates") {
		t.Errorf("expected QUALITY_FRAMEWORK filled, got:\n%s", output)
	}
}

func TestMergeFile_WholePersonaAsset(t *testing.T) {
	reg := newTestRegistry(nil)
	personaDir := t.TempDir()

	// Create a persona-only file (styles are 100% persona).
	stylesDir := filepath.Join(personaDir, "styles")
	if err := os.MkdirAll(stylesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	personaContent := "# Persona Style\nCustom style content here."
	if err := os.WriteFile(filepath.Join(stylesDir, "sprint.md"), []byte(personaContent), 0o644); err != nil {
		t.Fatal(err)
	}

	manifest := &model.PersonaManifest{
		Styles: []string{"styles/sprint.md"},
	}

	coreDir := t.TempDir() // empty -- no core template for styles
	outputDir := t.TempDir()

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	result, err := m.CopyPersonaFile("styles/sprint.md")
	if err != nil {
		t.Fatalf("CopyPersonaFile error: %v", err)
	}

	if result.FilesWritten != 1 {
		t.Errorf("expected 1 file written, got %d", result.FilesWritten)
	}

	outPath := filepath.Join(outputDir, "styles", "sprint.md")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}

	if string(data) != personaContent {
		t.Errorf("expected persona content copied as-is.\nexpected:\n%s\ngot:\n%s", personaContent, string(data))
	}
}

func TestMergeFile_MissingPersonaContent_UsesDefault(t *testing.T) {
	reg := newTestRegistry(map[string]*template.SlotEntry{
		"QUALITY_FRAMEWORK": {
			Category:   "section",
			MarkerType: "section",
			Default:    "default quality content",
		},
	})

	// Manifest with empty slot content -- no persona value for the slot.
	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{},
	}

	coreDir := t.TempDir()
	coreFile := filepath.Join(coreDir, "rules", "quality.md")
	if err := os.MkdirAll(filepath.Dir(coreFile), 0o755); err != nil {
		t.Fatal(err)
	}
	coreContent := "# Quality\n<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->\nplaceholder\n<!-- END_SLOT:QUALITY_FRAMEWORK -->"
	if err := os.WriteFile(coreFile, []byte(coreContent), 0o644); err != nil {
		t.Fatal(err)
	}

	outputDir := t.TempDir()
	personaDir := t.TempDir()

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	result, err := m.MergeFile("rules/quality.md")
	if err != nil {
		t.Fatalf("MergeFile error: %v", err)
	}

	if result.SlotsResolved != 1 {
		t.Errorf("expected 1 slot resolved (via default), got %d", result.SlotsResolved)
	}
	if len(result.Warnings) != 1 {
		t.Errorf("expected 1 warning for default fallback, got %d: %v", len(result.Warnings), result.Warnings)
	}

	outPath := filepath.Join(outputDir, "rules", "quality.md")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}

	output := string(data)
	if !strings.Contains(output, "default quality content") {
		t.Errorf("expected default content used, got:\n%s", output)
	}
}

func TestMergeFile_SectionOrderPreserved(t *testing.T) {
	reg := newTestRegistry(map[string]*template.SlotEntry{
		"SLOT_A": {
			Category:   "section",
			MarkerType: "section",
			Default:    "default A",
		},
		"SLOT_B": {
			Category:   "section",
			MarkerType: "section",
			Default:    "default B",
		},
	})

	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{
			"SLOT_A": "filled A content",
			"SLOT_B": "filled B content",
		},
	}

	coreDir := t.TempDir()
	coreFile := filepath.Join(coreDir, "doc.md")
	// Section order: Header1, SlotA, Header2, SlotB, Header3
	coreContent := `# Header 1
intro text

<!-- BEGIN_SLOT:SLOT_A -->
placeholder A
<!-- END_SLOT:SLOT_A -->

## Header 2
middle text

<!-- BEGIN_SLOT:SLOT_B -->
placeholder B
<!-- END_SLOT:SLOT_B -->

## Header 3
end text`
	if err := os.WriteFile(coreFile, []byte(coreContent), 0o644); err != nil {
		t.Fatal(err)
	}

	outputDir := t.TempDir()
	personaDir := t.TempDir()

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	result, err := m.MergeFile("doc.md")
	if err != nil {
		t.Fatalf("MergeFile error: %v", err)
	}

	if result.SlotsResolved != 2 {
		t.Errorf("expected 2 slots resolved, got %d", result.SlotsResolved)
	}

	outPath := filepath.Join(outputDir, "doc.md")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}

	output := string(data)
	// Verify section order is preserved: Header1 < SlotA < Header2 < SlotB < Header3
	idxH1 := strings.Index(output, "# Header 1")
	idxA := strings.Index(output, "filled A content")
	idxH2 := strings.Index(output, "## Header 2")
	idxB := strings.Index(output, "filled B content")
	idxH3 := strings.Index(output, "## Header 3")

	if idxH1 >= idxA || idxA >= idxH2 || idxH2 >= idxB || idxB >= idxH3 {
		t.Errorf("section order not preserved.\nH1=%d, A=%d, H2=%d, B=%d, H3=%d\noutput:\n%s",
			idxH1, idxA, idxH2, idxB, idxH3, output)
	}
}

func TestPatchAgent_AppendSkills(t *testing.T) {
	reg := newTestRegistry(nil)

	manifest := &model.PersonaManifest{
		AgentPatches: map[string]*model.AgentPatch{
			"agents/expert-backend.md": {
				AppendSkills: []string{"do-quality", "do-spec"},
			},
		},
	}

	coreDir := t.TempDir()
	outputDir := t.TempDir()
	personaDir := t.TempDir()

	// Create core agent file in output (simulating post-copy state).
	agentDir := filepath.Join(outputDir, "agents")
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	agentContent := `---
name: expert-backend
description: Backend development expert
skills:
    - do-foundation
    - do-backend
---
# Expert Backend

Implementation agent.`
	agentPath := filepath.Join(agentDir, "expert-backend.md")
	if err := os.WriteFile(agentPath, []byte(agentContent), 0o644); err != nil {
		t.Fatal(err)
	}

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	err := m.PatchAgent("agents/expert-backend.md")
	if err != nil {
		t.Fatalf("PatchAgent error: %v", err)
	}

	data, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("read patched agent: %v", err)
	}

	output := string(data)
	if !strings.Contains(output, "do-quality") {
		t.Errorf("expected appended skill 'do-quality' in output:\n%s", output)
	}
	if !strings.Contains(output, "do-spec") {
		t.Errorf("expected appended skill 'do-spec' in output:\n%s", output)
	}
	// Original skills should still be there.
	if !strings.Contains(output, "do-foundation") {
		t.Errorf("expected original skill 'do-foundation' preserved:\n%s", output)
	}
}

func TestPatchAgent_RemoveSkills(t *testing.T) {
	reg := newTestRegistry(nil)

	manifest := &model.PersonaManifest{
		AgentPatches: map[string]*model.AgentPatch{
			"agents/expert-backend.md": {
				RemoveSkills: []string{"do-legacy"},
			},
		},
	}

	coreDir := t.TempDir()
	outputDir := t.TempDir()
	personaDir := t.TempDir()

	agentDir := filepath.Join(outputDir, "agents")
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	agentContent := `---
name: expert-backend
description: Backend expert
skills:
    - do-foundation
    - do-legacy
    - do-backend
---
# Expert Backend

Content here.`
	agentPath := filepath.Join(agentDir, "expert-backend.md")
	if err := os.WriteFile(agentPath, []byte(agentContent), 0o644); err != nil {
		t.Fatal(err)
	}

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	err := m.PatchAgent("agents/expert-backend.md")
	if err != nil {
		t.Fatalf("PatchAgent error: %v", err)
	}

	data, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("read patched agent: %v", err)
	}

	output := string(data)
	if strings.Contains(output, "do-legacy") {
		t.Errorf("expected skill 'do-legacy' removed, but still present:\n%s", output)
	}
	if !strings.Contains(output, "do-foundation") {
		t.Errorf("expected 'do-foundation' preserved:\n%s", output)
	}
	if !strings.Contains(output, "do-backend") {
		t.Errorf("expected 'do-backend' preserved:\n%s", output)
	}
}

func TestPatchAgent_AppendContent(t *testing.T) {
	reg := newTestRegistry(nil)
	personaDir := t.TempDir()

	// Create persona content file to append.
	contentDir := filepath.Join(personaDir, "content")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	appendContent := "\n## Persona Rules\n\nCustom rules for this persona."
	if err := os.WriteFile(filepath.Join(contentDir, "backend-extra.md"), []byte(appendContent), 0o644); err != nil {
		t.Fatal(err)
	}

	manifest := &model.PersonaManifest{
		AgentPatches: map[string]*model.AgentPatch{
			"agents/expert-backend.md": {
				AppendContent: "content/backend-extra.md",
			},
		},
	}

	coreDir := t.TempDir()
	outputDir := t.TempDir()

	agentDir := filepath.Join(outputDir, "agents")
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	agentContent := `---
name: expert-backend
description: Backend expert
---
# Expert Backend

Core content.`
	agentPath := filepath.Join(agentDir, "expert-backend.md")
	if err := os.WriteFile(agentPath, []byte(agentContent), 0o644); err != nil {
		t.Fatal(err)
	}

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	err := m.PatchAgent("agents/expert-backend.md")
	if err != nil {
		t.Fatalf("PatchAgent error: %v", err)
	}

	data, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("read patched agent: %v", err)
	}

	output := string(data)
	if !strings.Contains(output, "Core content.") {
		t.Errorf("expected original content preserved:\n%s", output)
	}
	if !strings.Contains(output, "Persona Rules") {
		t.Errorf("expected appended persona content:\n%s", output)
	}
}

func TestCopyPersonaFile_CreatesDirectories(t *testing.T) {
	reg := newTestRegistry(nil)
	personaDir := t.TempDir()

	// Create deeply nested persona file.
	nestedDir := filepath.Join(personaDir, "agents", "moai")
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	personaContent := "# Persona Agent\nPersona-only agent."
	if err := os.WriteFile(filepath.Join(nestedDir, "manager-spec.md"), []byte(personaContent), 0o644); err != nil {
		t.Fatal(err)
	}

	manifest := &model.PersonaManifest{
		Agents: []string{"agents/moai/manager-spec.md"},
	}

	coreDir := t.TempDir()
	outputDir := t.TempDir()

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	result, err := m.CopyPersonaFile("agents/moai/manager-spec.md")
	if err != nil {
		t.Fatalf("CopyPersonaFile error: %v", err)
	}

	if result.FilesWritten != 1 {
		t.Errorf("expected 1 file written, got %d", result.FilesWritten)
	}

	outPath := filepath.Join(outputDir, "agents", "moai", "manager-spec.md")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	if string(data) != personaContent {
		t.Errorf("content mismatch.\nexpected:\n%s\ngot:\n%s", personaContent, string(data))
	}
}

func TestMergeFile_NoSlots_CopiesAsIs(t *testing.T) {
	reg := newTestRegistry(nil)
	manifest := &model.PersonaManifest{}

	coreDir := t.TempDir()
	coreContent := "# Simple Document\nNo slots here, just plain content."
	if err := os.WriteFile(filepath.Join(coreDir, "plain.md"), []byte(coreContent), 0o644); err != nil {
		t.Fatal(err)
	}

	outputDir := t.TempDir()
	personaDir := t.TempDir()

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	result, err := m.MergeFile("plain.md")
	if err != nil {
		t.Fatalf("MergeFile error: %v", err)
	}

	if result.SlotsResolved != 0 {
		t.Errorf("expected 0 slots resolved, got %d", result.SlotsResolved)
	}

	outPath := filepath.Join(outputDir, "plain.md")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	if string(data) != coreContent {
		t.Errorf("expected content copied as-is.\nexpected:\n%s\ngot:\n%s", coreContent, string(data))
	}
}
