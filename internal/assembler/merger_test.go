package assembler

import (
	"encoding/json"
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

func TestMergeSettings_CorePlusPersonaHooks(t *testing.T) {
	reg := newTestRegistry(nil)
	personaDir := t.TempDir()
	coreDir := t.TempDir()
	outputDir := t.TempDir()

	// Create core settings file with permissions and outputStyle.
	coreSettings := map[string]interface{}{
		"permissions": map[string]interface{}{
			"allow": []interface{}{"Read", "Write"},
		},
		"outputStyle": "pair",
	}
	coreData, _ := json.MarshalIndent(coreSettings, "", "  ")
	coreSettingsPath := filepath.Join(coreDir, "settings.json")
	if err := os.WriteFile(coreSettingsPath, coreData, 0o644); err != nil {
		t.Fatal(err)
	}

	// Manifest with persona hooks.
	manifest := &model.PersonaManifest{
		Hooks: map[string][]model.HookEntry{
			"PreToolUse": {
				{Command: "godo hook pre-tool", Timeout: 5000},
			},
			"PostToolUse": {
				{Command: "godo hook post-tool"},
			},
		},
	}

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	err := m.MergeSettings(coreSettingsPath)
	if err != nil {
		t.Fatalf("MergeSettings error: %v", err)
	}

	// Read output settings.json.
	outPath := filepath.Join(outputDir, "settings.json")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output settings: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal output settings: %v", err)
	}

	// Core fields preserved.
	if result["outputStyle"] != "pair" {
		t.Errorf("expected outputStyle 'pair', got %v", result["outputStyle"])
	}
	perms, ok := result["permissions"]
	if !ok {
		t.Fatal("expected permissions field preserved")
	}
	permMap, ok := perms.(map[string]interface{})
	if !ok {
		t.Fatalf("expected permissions as map, got %T", perms)
	}
	if _, ok := permMap["allow"]; !ok {
		t.Error("expected permissions.allow preserved")
	}

	// Persona hooks injected.
	hooks, ok := result["hooks"]
	if !ok {
		t.Fatal("expected hooks field in merged settings")
	}
	hooksMap, ok := hooks.(map[string]interface{})
	if !ok {
		t.Fatalf("expected hooks as map, got %T", hooks)
	}
	if _, ok := hooksMap["PreToolUse"]; !ok {
		t.Error("expected PreToolUse hook in merged settings")
	}
	if _, ok := hooksMap["PostToolUse"]; !ok {
		t.Error("expected PostToolUse hook in merged settings")
	}
}

func TestMergeSettings_NoPersonaHooks(t *testing.T) {
	reg := newTestRegistry(nil)
	personaDir := t.TempDir()
	coreDir := t.TempDir()
	outputDir := t.TempDir()

	coreSettings := map[string]interface{}{
		"outputStyle": "direct",
	}
	coreData, _ := json.MarshalIndent(coreSettings, "", "  ")
	coreSettingsPath := filepath.Join(coreDir, "settings.json")
	if err := os.WriteFile(coreSettingsPath, coreData, 0o644); err != nil {
		t.Fatal(err)
	}

	// No hooks in manifest.
	manifest := &model.PersonaManifest{}

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	err := m.MergeSettings(coreSettingsPath)
	if err != nil {
		t.Fatalf("MergeSettings error: %v", err)
	}

	outPath := filepath.Join(outputDir, "settings.json")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output settings: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal output settings: %v", err)
	}

	if result["outputStyle"] != "direct" {
		t.Errorf("expected outputStyle 'direct', got %v", result["outputStyle"])
	}
	// No hooks field should be present when no persona hooks.
	if _, ok := result["hooks"]; ok {
		t.Error("expected no hooks field when persona has no hooks")
	}
}

func TestCopyCommands(t *testing.T) {
	reg := newTestRegistry(nil)
	personaDir := t.TempDir()
	coreDir := t.TempDir()
	outputDir := t.TempDir()

	// Create persona command files.
	cmdDir := filepath.Join(personaDir, "commands")
	if err := os.MkdirAll(cmdDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cmd1Content := "# Plan command\nCreate a plan."
	if err := os.WriteFile(filepath.Join(cmdDir, "do-plan.md"), []byte(cmd1Content), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd2Content := "# Style command\nSwitch style."
	subDir := filepath.Join(cmdDir, "sub")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "do-style.md"), []byte(cmd2Content), 0o644); err != nil {
		t.Fatal(err)
	}

	manifest := &model.PersonaManifest{
		Commands: []string{"commands/do-plan.md", "commands/sub/do-style.md"},
	}

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	result, err := m.CopyCommands()
	if err != nil {
		t.Fatalf("CopyCommands error: %v", err)
	}

	if result.FilesWritten != 2 {
		t.Errorf("expected 2 files written, got %d", result.FilesWritten)
	}

	// Verify first command.
	out1 := filepath.Join(outputDir, "commands", "do-plan.md")
	data1, err := os.ReadFile(out1)
	if err != nil {
		t.Fatalf("read command 1: %v", err)
	}
	if string(data1) != cmd1Content {
		t.Errorf("command 1 content mismatch.\nexpected:\n%s\ngot:\n%s", cmd1Content, string(data1))
	}

	// Verify second command (nested).
	out2 := filepath.Join(outputDir, "commands", "sub", "do-style.md")
	data2, err := os.ReadFile(out2)
	if err != nil {
		t.Fatalf("read command 2: %v", err)
	}
	if string(data2) != cmd2Content {
		t.Errorf("command 2 content mismatch.\nexpected:\n%s\ngot:\n%s", cmd2Content, string(data2))
	}
}

func TestCopyCommands_EmptyList(t *testing.T) {
	reg := newTestRegistry(nil)
	personaDir := t.TempDir()
	coreDir := t.TempDir()
	outputDir := t.TempDir()

	manifest := &model.PersonaManifest{
		Commands: nil,
	}

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	result, err := m.CopyCommands()
	if err != nil {
		t.Fatalf("CopyCommands error: %v", err)
	}

	if result.FilesWritten != 0 {
		t.Errorf("expected 0 files written for empty commands, got %d", result.FilesWritten)
	}
}

func TestCopyHookScripts(t *testing.T) {
	reg := newTestRegistry(nil)
	personaDir := t.TempDir()
	coreDir := t.TempDir()
	outputDir := t.TempDir()

	// Create persona hook script files.
	hookDir := filepath.Join(personaDir, "hooks")
	if err := os.MkdirAll(hookDir, 0o755); err != nil {
		t.Fatal(err)
	}
	script1 := "#!/bin/bash\ngodo hook pre-tool \"$@\""
	if err := os.WriteFile(filepath.Join(hookDir, "pre-tool.sh"), []byte(script1), 0o755); err != nil {
		t.Fatal(err)
	}
	script2 := "#!/bin/bash\ngodo hook post-tool \"$@\""
	if err := os.WriteFile(filepath.Join(hookDir, "post-tool.sh"), []byte(script2), 0o755); err != nil {
		t.Fatal(err)
	}

	manifest := &model.PersonaManifest{
		HookScripts: []string{"hooks/pre-tool.sh", "hooks/post-tool.sh"},
	}

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	result, err := m.CopyHookScripts()
	if err != nil {
		t.Fatalf("CopyHookScripts error: %v", err)
	}

	if result.FilesWritten != 2 {
		t.Errorf("expected 2 files written, got %d", result.FilesWritten)
	}

	// Verify first script.
	out1 := filepath.Join(outputDir, "hooks", "pre-tool.sh")
	data1, err := os.ReadFile(out1)
	if err != nil {
		t.Fatalf("read hook script 1: %v", err)
	}
	if string(data1) != script1 {
		t.Errorf("hook script 1 content mismatch.\nexpected:\n%s\ngot:\n%s", script1, string(data1))
	}

	// Verify executable permission preserved.
	info, err := os.Stat(out1)
	if err != nil {
		t.Fatalf("stat hook script 1: %v", err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Errorf("expected hook script to be executable, got mode %v", info.Mode())
	}

	// Verify second script.
	out2 := filepath.Join(outputDir, "hooks", "post-tool.sh")
	data2, err := os.ReadFile(out2)
	if err != nil {
		t.Fatalf("read hook script 2: %v", err)
	}
	if string(data2) != script2 {
		t.Errorf("hook script 2 content mismatch.\nexpected:\n%s\ngot:\n%s", script2, string(data2))
	}
}

func TestCopyHookScripts_EmptyList(t *testing.T) {
	reg := newTestRegistry(nil)
	personaDir := t.TempDir()
	coreDir := t.TempDir()
	outputDir := t.TempDir()

	manifest := &model.PersonaManifest{
		HookScripts: nil,
	}

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	result, err := m.CopyHookScripts()
	if err != nil {
		t.Fatalf("CopyHookScripts error: %v", err)
	}

	if result.FilesWritten != 0 {
		t.Errorf("expected 0 files written for empty hook scripts, got %d", result.FilesWritten)
	}
}

func TestPatchAgent_PreservesInlineSkillsFormat(t *testing.T) {
	reg := newTestRegistry(nil)

	manifest := &model.PersonaManifest{
		AgentPatches: map[string]*model.AgentPatch{
			"agents/moai/builder-agent.md": {
				AppendSkills: []string{"moai-persona-custom"},
			},
		},
	}

	coreDir := t.TempDir()
	outputDir := t.TempDir()
	personaDir := t.TempDir()

	// Create agent file with inline comma-separated skills and specific key order.
	agentDir := filepath.Join(outputDir, "agents", "moai")
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	agentContent := `---
description: |
  Agent creation specialist for building new MoAI agents.
memory: user
model: inherit
permissionMode: bypassPermissions
skills: moai-foundation-claude, moai-workflow-project
---
# Builder Agent

Creates new agent definitions.`
	agentPath := filepath.Join(agentDir, "builder-agent.md")
	if err := os.WriteFile(agentPath, []byte(agentContent), 0o644); err != nil {
		t.Fatal(err)
	}

	m := NewMerger(coreDir, personaDir, outputDir, manifest, reg)
	err := m.PatchAgent("agents/moai/builder-agent.md")
	if err != nil {
		t.Fatalf("PatchAgent error: %v", err)
	}

	data, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("read patched agent: %v", err)
	}

	output := string(data)

	// Verify skills stayed inline comma-separated (not converted to YAML list).
	if strings.Contains(output, "    - moai-foundation-claude") {
		t.Errorf("skills should remain inline, but got YAML list format:\n%s", output)
	}
	if !strings.Contains(output, "skills: moai-foundation-claude, moai-workflow-project, moai-persona-custom") {
		t.Errorf("expected inline skills with appended skill, got:\n%s", output)
	}

	// Verify key order is preserved (description before memory before model before permissionMode before skills).
	descIdx := strings.Index(output, "description:")
	memIdx := strings.Index(output, "memory:")
	modelIdx := strings.Index(output, "model:")
	permIdx := strings.Index(output, "permissionMode:")
	skillsIdx := strings.Index(output, "skills:")

	if descIdx >= memIdx || memIdx >= modelIdx || modelIdx >= permIdx || permIdx >= skillsIdx {
		t.Errorf("key order not preserved.\ndesc=%d, mem=%d, model=%d, perm=%d, skills=%d\noutput:\n%s",
			descIdx, memIdx, modelIdx, permIdx, skillsIdx, output)
	}

	// Verify body is preserved.
	if !strings.Contains(output, "# Builder Agent") {
		t.Errorf("expected body content preserved:\n%s", output)
	}
}

func TestPatchAgent_PreservesListSkillsFormat(t *testing.T) {
	reg := newTestRegistry(nil)

	manifest := &model.PersonaManifest{
		AgentPatches: map[string]*model.AgentPatch{
			"agents/expert-backend.md": {
				AppendSkills: []string{"do-quality"},
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
	// Use YAML list format with 4-space indent.
	agentContent := `---
name: expert-backend
description: Backend development expert
skills:
    - do-foundation
    - do-backend
memory: project
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

	// Verify skills stayed in list format with same indentation.
	if !strings.Contains(output, "skills:\n    - do-foundation\n    - do-backend\n    - do-quality") {
		t.Errorf("expected list format with 4-space indent preserved:\n%s", output)
	}

	// Verify key order preserved.
	nameIdx := strings.Index(output, "name:")
	descIdx := strings.Index(output, "description:")
	skillsIdx := strings.Index(output, "skills:")
	memIdx := strings.Index(output, "memory:")

	if nameIdx >= descIdx || descIdx >= skillsIdx || skillsIdx >= memIdx {
		t.Errorf("key order not preserved:\n%s", output)
	}
}
