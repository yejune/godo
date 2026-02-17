package assembler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/do-focus/convert/internal/model"
	"github.com/do-focus/convert/internal/template"
)

func TestAssemble_TemplatesAndPersonaContent(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()
	outputDir := t.TempDir()

	// Create a core template file with a slot.
	rulesDir := filepath.Join(coreDir, "rules")
	if err := os.MkdirAll(rulesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	coreContent := "# Quality\n<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->\nplaceholder\n<!-- END_SLOT:QUALITY_FRAMEWORK -->\n# End"
	if err := os.WriteFile(filepath.Join(rulesDir, "quality.md"), []byte(coreContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a second core file without slots.
	if err := os.WriteFile(filepath.Join(rulesDir, "plain.md"), []byte("# Plain\nNo slots here."), 0o644); err != nil {
		t.Fatal(err)
	}

	// Registry with the slot definition.
	reg := newTestRegistry(map[string]*template.SlotEntry{
		"QUALITY_FRAMEWORK": {
			Category:   "section",
			MarkerType: "section",
			Default:    "default quality",
		},
	})

	manifest := &model.PersonaManifest{
		Name: "test-persona",
		SlotContent: map[string]string{
			"QUALITY_FRAMEWORK": "TRUST 5 quality gates",
		},
	}

	orch := NewAssembler(coreDir, personaDir, outputDir, manifest, reg)
	result, err := orch.Assemble()
	if err != nil {
		t.Fatalf("Assemble error: %v", err)
	}

	if result.FilesWritten < 2 {
		t.Errorf("expected at least 2 files written, got %d", result.FilesWritten)
	}
	if result.SlotsResolved != 1 {
		t.Errorf("expected 1 slot resolved, got %d", result.SlotsResolved)
	}

	// Verify the slot-filled file.
	data, err := os.ReadFile(filepath.Join(outputDir, "rules", "quality.md"))
	if err != nil {
		t.Fatalf("read output quality.md: %v", err)
	}
	if !strings.Contains(string(data), "TRUST 5 quality gates") {
		t.Errorf("expected filled slot content, got:\n%s", string(data))
	}

	// Verify the plain file was copied.
	data, err = os.ReadFile(filepath.Join(outputDir, "rules", "plain.md"))
	if err != nil {
		t.Fatalf("read output plain.md: %v", err)
	}
	if !strings.Contains(string(data), "No slots here.") {
		t.Errorf("expected plain content copied, got:\n%s", string(data))
	}
}

func TestAssemble_PersonaOnlyFiles(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()
	outputDir := t.TempDir()

	// Create persona-only style files.
	stylesDir := filepath.Join(personaDir, "styles")
	if err := os.MkdirAll(stylesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	styleContent := "# Sprint Style\nCustom sprint style."
	if err := os.WriteFile(filepath.Join(stylesDir, "sprint.md"), []byte(styleContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create persona-only agent.
	agentsDir := filepath.Join(personaDir, "agents")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	agentContent := "---\nname: custom-agent\n---\n# Custom Agent"
	if err := os.WriteFile(filepath.Join(agentsDir, "custom.md"), []byte(agentContent), 0o644); err != nil {
		t.Fatal(err)
	}

	reg := newTestRegistry(nil)
	manifest := &model.PersonaManifest{
		Name:   "test-persona",
		Styles: []string{"styles/sprint.md"},
		Agents: []string{"agents/custom.md"},
	}

	orch := NewAssembler(coreDir, personaDir, outputDir, manifest, reg)
	result, err := orch.Assemble()
	if err != nil {
		t.Fatalf("Assemble error: %v", err)
	}

	if result.FilesWritten != 2 {
		t.Errorf("expected 2 files written, got %d", result.FilesWritten)
	}

	// Verify style file copied.
	data, err := os.ReadFile(filepath.Join(outputDir, "styles", "sprint.md"))
	if err != nil {
		t.Fatalf("read output style: %v", err)
	}
	if string(data) != styleContent {
		t.Errorf("expected style content:\n%s\ngot:\n%s", styleContent, string(data))
	}

	// Verify agent file copied.
	data, err = os.ReadFile(filepath.Join(outputDir, "agents", "test-persona", "custom.md"))
	if err != nil {
		t.Fatalf("read output agent: %v", err)
	}
	if string(data) != agentContent {
		t.Errorf("expected agent content:\n%s\ngot:\n%s", agentContent, string(data))
	}
}

func TestAssemble_EmptyRegistryAndManifest(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()
	outputDir := t.TempDir()

	reg := newTestRegistry(nil)
	manifest := &model.PersonaManifest{
		Name: "empty-persona",
	}

	orch := NewAssembler(coreDir, personaDir, outputDir, manifest, reg)
	result, err := orch.Assemble()
	if err != nil {
		t.Fatalf("Assemble error: %v", err)
	}

	if result.FilesWritten != 0 {
		t.Errorf("expected 0 files written, got %d", result.FilesWritten)
	}
	if result.SlotsResolved != 0 {
		t.Errorf("expected 0 slots resolved, got %d", result.SlotsResolved)
	}
	if result.AgentsPatched != 0 {
		t.Errorf("expected 0 agents patched, got %d", result.AgentsPatched)
	}
}

func TestAssemble_AgentPatchingApplied(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()
	outputDir := t.TempDir()

	// Create a core agent file that will be copied then patched.
	agentsDir := filepath.Join(coreDir, "agents", "test-persona")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
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

Core implementation.`
	if err := os.WriteFile(filepath.Join(agentsDir, "expert-backend.md"), []byte(agentContent), 0o644); err != nil {
		t.Fatal(err)
	}

	reg := newTestRegistry(nil)
	manifest := &model.PersonaManifest{
		Name: "test-persona",
		AgentPatches: map[string]*model.AgentPatch{
			"agents/expert-backend.md": {
				AppendSkills: []string{"do-quality"},
				RemoveSkills: []string{"do-legacy"},
			},
		},
	}

	orch := NewAssembler(coreDir, personaDir, outputDir, manifest, reg)
	result, err := orch.Assemble()
	if err != nil {
		t.Fatalf("Assemble error: %v", err)
	}

	if result.AgentsPatched != 1 {
		t.Errorf("expected 1 agent patched, got %d", result.AgentsPatched)
	}

	// Verify the patched file.
	data, err := os.ReadFile(filepath.Join(outputDir, "agents", "test-persona", "expert-backend.md"))
	if err != nil {
		t.Fatalf("read patched agent: %v", err)
	}
	output := string(data)

	if !strings.Contains(output, "do-quality") {
		t.Errorf("expected appended skill 'do-quality' in output:\n%s", output)
	}
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

func TestAssemble_ClaudeMDCopied(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()
	outputDir := t.TempDir()

	// Create persona CLAUDE.md.
	claudeContent := "# Do Persona\nThis is the persona CLAUDE.md."
	if err := os.WriteFile(filepath.Join(personaDir, "CLAUDE.md"), []byte(claudeContent), 0o644); err != nil {
		t.Fatal(err)
	}

	reg := newTestRegistry(nil)
	manifest := &model.PersonaManifest{
		Name:     "test-persona",
		ClaudeMD: "CLAUDE.md",
	}

	orch := NewAssembler(coreDir, personaDir, outputDir, manifest, reg)
	result, err := orch.Assemble()
	if err != nil {
		t.Fatalf("Assemble error: %v", err)
	}

	if result.FilesWritten < 1 {
		t.Errorf("expected at least 1 file written, got %d", result.FilesWritten)
	}

	// CLAUDE.md goes to the parent of .claude/ -- output dir root.
	data, err := os.ReadFile(filepath.Join(outputDir, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("read output CLAUDE.md: %v", err)
	}
	if string(data) != claudeContent {
		t.Errorf("expected CLAUDE.md content:\n%s\ngot:\n%s", claudeContent, string(data))
	}
}

func TestAssemble_CommandsAndHookScriptsCopied(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()
	outputDir := t.TempDir()

	// Create persona commands.
	cmdsDir := filepath.Join(personaDir, "commands")
	if err := os.MkdirAll(cmdsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cmdContent := "# /do:plan\nPlan command content."
	if err := os.WriteFile(filepath.Join(cmdsDir, "do-plan.md"), []byte(cmdContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create persona hook scripts.
	hooksDir := filepath.Join(personaDir, "hooks")
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		t.Fatal(err)
	}
	hookContent := "#!/bin/bash\necho 'hook executed'"
	if err := os.WriteFile(filepath.Join(hooksDir, "pre-tool.sh"), []byte(hookContent), 0o644); err != nil {
		t.Fatal(err)
	}

	reg := newTestRegistry(nil)
	manifest := &model.PersonaManifest{
		Name:        "test-persona",
		Commands:    []string{"commands/do-plan.md"},
		HookScripts: []string{"hooks/pre-tool.sh"},
	}

	orch := NewAssembler(coreDir, personaDir, outputDir, manifest, reg)
	result, err := orch.Assemble()
	if err != nil {
		t.Fatalf("Assemble error: %v", err)
	}

	if result.FilesWritten != 2 {
		t.Errorf("expected 2 files written, got %d", result.FilesWritten)
	}

	// Verify command copied.
	data, err := os.ReadFile(filepath.Join(outputDir, "commands", "test-persona", "do-plan.md"))
	if err != nil {
		t.Fatalf("read output command: %v", err)
	}
	if string(data) != cmdContent {
		t.Errorf("expected command content:\n%s\ngot:\n%s", cmdContent, string(data))
	}

	// Verify hook script copied.
	data, err = os.ReadFile(filepath.Join(outputDir, "hooks", "test-persona", "pre-tool.sh"))
	if err != nil {
		t.Fatalf("read output hook: %v", err)
	}
	if string(data) != hookContent {
		t.Errorf("expected hook content:\n%s\ngot:\n%s", hookContent, string(data))
	}
}

func TestAssemble_FileListReturnedSorted(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()
	outputDir := t.TempDir()

	// Create core files in multiple subdirs.
	for _, relPath := range []string{"rules/b.md", "agents/a.md", "skills/c.md"} {
		fullPath := filepath.Join(coreDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte("# "+relPath), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	reg := newTestRegistry(nil)
	manifest := &model.PersonaManifest{Name: "test-persona"}

	orch := NewAssembler(coreDir, personaDir, outputDir, manifest, reg)
	result, err := orch.Assemble()
	if err != nil {
		t.Fatalf("Assemble error: %v", err)
	}

	if len(result.Files) != 3 {
		t.Fatalf("expected 3 files, got %d: %v", len(result.Files), result.Files)
	}

	// Files should be sorted.
	sorted := make([]string, len(result.Files))
	copy(sorted, result.Files)
	sort.Strings(sorted)
	for i := range sorted {
		if result.Files[i] != sorted[i] {
			t.Errorf("files not sorted at index %d: got %q, expected %q", i, result.Files[i], sorted[i])
		}
	}
}

func TestAssemble_FullPipeline(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()
	outputDir := t.TempDir()

	// Core: agent with slot + plain rule.
	agentsDir := filepath.Join(coreDir, "agents", "godo")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	coreAgent := `---
name: expert-backend
description: Backend expert
skills:
    - do-foundation
---
# Expert Backend

Use {{slot:TOOL_NAME}} for development.`
	if err := os.WriteFile(filepath.Join(agentsDir, "expert-backend.md"), []byte(coreAgent), 0o644); err != nil {
		t.Fatal(err)
	}

	rulesDir := filepath.Join(coreDir, "rules")
	if err := os.MkdirAll(rulesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(rulesDir, "core-rule.md"), []byte("# Core Rule"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Persona: style file + CLAUDE.md + command.
	pStylesDir := filepath.Join(personaDir, "styles")
	if err := os.MkdirAll(pStylesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pStylesDir, "pair.md"), []byte("# Pair Style"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(personaDir, "CLAUDE.md"), []byte("# My Persona"), 0o644); err != nil {
		t.Fatal(err)
	}

	pCmdsDir := filepath.Join(personaDir, "commands")
	if err := os.MkdirAll(pCmdsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pCmdsDir, "plan.md"), []byte("# Plan Cmd"), 0o644); err != nil {
		t.Fatal(err)
	}

	reg := newTestRegistry(map[string]*template.SlotEntry{
		"TOOL_NAME": {
			Category:   "path_pattern",
			MarkerType: "inline",
			Default:    "moai",
		},
	})

	manifest := &model.PersonaManifest{
		Name:     "godo",
		ClaudeMD: "CLAUDE.md",
		Styles:   []string{"styles/pair.md"},
		Commands: []string{"commands/plan.md"},
		SlotContent: map[string]string{
			"TOOL_NAME": "godo",
		},
		AgentPatches: map[string]*model.AgentPatch{
			"agents/expert-backend.md": {
				AppendSkills: []string{"do-godo-extra"},
			},
		},
	}

	orch := NewAssembler(coreDir, personaDir, outputDir, manifest, reg)
	result, err := orch.Assemble()
	if err != nil {
		t.Fatalf("Assemble error: %v", err)
	}

	// 2 core files + 1 style + 1 CLAUDE.md + 1 command = 5
	if result.FilesWritten != 5 {
		t.Errorf("expected 5 files written, got %d", result.FilesWritten)
	}
	if result.SlotsResolved != 1 {
		t.Errorf("expected 1 slot resolved (TOOL_NAME), got %d", result.SlotsResolved)
	}
	if result.AgentsPatched != 1 {
		t.Errorf("expected 1 agent patched, got %d", result.AgentsPatched)
	}

	// Verify inline slot filled.
	data, err := os.ReadFile(filepath.Join(outputDir, "agents", "godo", "expert-backend.md"))
	if err != nil {
		t.Fatalf("read output agent: %v", err)
	}
	output := string(data)
	if !strings.Contains(output, "godo") {
		t.Errorf("expected TOOL_NAME filled with 'godo', got:\n%s", output)
	}
	if strings.Contains(output, "{{slot:TOOL_NAME}}") {
		t.Errorf("expected TOOL_NAME slot replaced, but marker still present:\n%s", output)
	}
	// Verify patch applied.
	if !strings.Contains(output, "do-godo-extra") {
		t.Errorf("expected appended skill 'do-godo-extra':\n%s", output)
	}

	// Verify persona style copied.
	data, err = os.ReadFile(filepath.Join(outputDir, "styles", "pair.md"))
	if err != nil {
		t.Fatalf("read output style: %v", err)
	}
	if string(data) != "# Pair Style" {
		t.Errorf("expected style content, got:\n%s", string(data))
	}

	// Verify CLAUDE.md at root.
	data, err = os.ReadFile(filepath.Join(outputDir, "CLAUDE.md"))
	if err != nil {
		t.Fatalf("read CLAUDE.md: %v", err)
	}
	if string(data) != "# My Persona" {
		t.Errorf("expected CLAUDE.md content, got:\n%s", string(data))
	}
}

func TestAssemble_SkillMappingsApplied(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()
	outputDir := t.TempDir()

	// Create core agent files.
	agentsDir := filepath.Join(coreDir, "agents")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	agent1 := `---
name: expert-backend
description: Backend expert
skills:
    - moai-foundation-quality
    - do-domain-backend
---
# Expert Backend`
	if err := os.WriteFile(filepath.Join(agentsDir, "expert-backend.md"), []byte(agent1), 0o644); err != nil {
		t.Fatal(err)
	}

	agent2 := `---
name: expert-frontend
description: Frontend expert
skills:
    - moai-foundation-quality
    - do-domain-frontend
---
# Expert Frontend`
	if err := os.WriteFile(filepath.Join(agentsDir, "expert-frontend.md"), []byte(agent2), 0o644); err != nil {
		t.Fatal(err)
	}

	reg := newTestRegistry(nil)
	manifest := &model.PersonaManifest{
		Name: "test-persona",
		SkillMappings: map[string]string{
			"moai-foundation-quality": "do-foundation-checklist",
		},
	}

	orch := NewAssembler(coreDir, personaDir, outputDir, manifest, reg)
	result, err := orch.Assemble()
	if err != nil {
		t.Fatalf("Assemble error: %v", err)
	}

	if result.SkillsMapped != 2 {
		t.Errorf("expected 2 agents with skill mappings applied, got %d", result.SkillsMapped)
	}

	// Verify both agents have the replaced skill.
	data1, err := os.ReadFile(filepath.Join(outputDir, "agents", "expert-backend.md"))
	if err != nil {
		t.Fatalf("read backend agent: %v", err)
	}
	if strings.Contains(string(data1), "moai-foundation-quality") {
		t.Error("backend agent should have moai-foundation-quality replaced")
	}
	if !strings.Contains(string(data1), "do-foundation-checklist") {
		t.Error("backend agent should have do-foundation-checklist after mapping")
	}

	data2, err := os.ReadFile(filepath.Join(outputDir, "agents", "expert-frontend.md"))
	if err != nil {
		t.Fatalf("read frontend agent: %v", err)
	}
	if strings.Contains(string(data2), "moai-foundation-quality") {
		t.Error("frontend agent should have moai-foundation-quality replaced")
	}
	if !strings.Contains(string(data2), "do-foundation-checklist") {
		t.Error("frontend agent should have do-foundation-checklist after mapping")
	}
}

func TestAssemble_SettingsMergeWithManifestSettings(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()
	outputDir := t.TempDir()

	// Core settings.json.
	coreSettings := map[string]interface{}{
		"outputStyle": "pair",
		"permissions": map[string]interface{}{
			"allow": []interface{}{"Read"},
		},
		"env": map[string]interface{}{
			"DO_MODE":     "do",
			"DO_LANGUAGE": "en",
		},
	}
	coreData, _ := json.MarshalIndent(coreSettings, "", "  ")
	if err := os.WriteFile(filepath.Join(coreDir, "settings.json"), coreData, 0o644); err != nil {
		t.Fatal(err)
	}

	reg := newTestRegistry(nil)
	manifest := &model.PersonaManifest{
		Name: "test-persona",
		Settings: map[string]interface{}{
			"outputStyle": "sprint",
			"env": map[string]interface{}{
				"DO_LANGUAGE": "ko",
				"DO_PERSONA":  "young-f",
			},
		},
	}

	orch := NewAssembler(coreDir, personaDir, outputDir, manifest, reg)
	result, err := orch.Assemble()
	if err != nil {
		t.Fatalf("Assemble error: %v", err)
	}

	if result.FilesWritten < 1 {
		t.Error("expected at least 1 file written (settings.json)")
	}

	// Read merged settings.
	data, err := os.ReadFile(filepath.Join(outputDir, "settings.json"))
	if err != nil {
		t.Fatalf("read settings.json: %v", err)
	}
	var merged map[string]interface{}
	if err := json.Unmarshal(data, &merged); err != nil {
		t.Fatalf("parse settings.json: %v", err)
	}

	// Top-level override.
	if merged["outputStyle"] != "sprint" {
		t.Errorf("expected outputStyle 'sprint' (persona override), got %v", merged["outputStyle"])
	}

	// Core permissions preserved.
	if _, ok := merged["permissions"]; !ok {
		t.Error("expected permissions preserved from core")
	}

	// Env deep merged.
	envRaw, ok := merged["env"]
	if !ok {
		t.Fatal("expected env field")
	}
	envMap := envRaw.(map[string]interface{})
	if envMap["DO_MODE"] != "do" {
		t.Errorf("expected DO_MODE preserved, got %v", envMap["DO_MODE"])
	}
	if envMap["DO_LANGUAGE"] != "ko" {
		t.Errorf("expected DO_LANGUAGE overridden to 'ko', got %v", envMap["DO_LANGUAGE"])
	}
	if envMap["DO_PERSONA"] != "young-f" {
		t.Errorf("expected DO_PERSONA added, got %v", envMap["DO_PERSONA"])
	}
}

func TestAssemble_PersonaSettingsJsonMergedWithCore(t *testing.T) {
	coreDir := t.TempDir()
	personaDir := t.TempDir()
	outputDir := t.TempDir()

	// Core settings.
	coreSettings := map[string]interface{}{
		"outputStyle": "pair",
		"permissions": map[string]interface{}{
			"allow": []interface{}{"Read", "Write"},
		},
	}
	coreData, _ := json.MarshalIndent(coreSettings, "", "  ")
	if err := os.WriteFile(filepath.Join(coreDir, "settings.json"), coreData, 0o644); err != nil {
		t.Fatal(err)
	}

	// Persona settings.json with override.
	personaSettings := map[string]interface{}{
		"outputStyle":    "direct",
		"plansDirectory": ".do/plans/",
	}
	personaData, _ := json.MarshalIndent(personaSettings, "", "  ")
	if err := os.WriteFile(filepath.Join(personaDir, "settings.json"), personaData, 0o644); err != nil {
		t.Fatal(err)
	}

	reg := newTestRegistry(nil)
	manifest := &model.PersonaManifest{
		Name: "test-persona",
	}

	orch := NewAssembler(coreDir, personaDir, outputDir, manifest, reg)
	_, err := orch.Assemble()
	if err != nil {
		t.Fatalf("Assemble error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outputDir, "settings.json"))
	if err != nil {
		t.Fatalf("read settings.json: %v", err)
	}
	var merged map[string]interface{}
	if err := json.Unmarshal(data, &merged); err != nil {
		t.Fatalf("parse settings.json: %v", err)
	}

	// Persona override.
	if merged["outputStyle"] != "direct" {
		t.Errorf("expected outputStyle 'direct' from persona settings.json, got %v", merged["outputStyle"])
	}
	// Persona addition.
	if merged["plansDirectory"] != ".do/plans/" {
		t.Errorf("expected plansDirectory from persona, got %v", merged["plansDirectory"])
	}
	// Core preserved.
	if _, ok := merged["permissions"]; !ok {
		t.Error("expected permissions preserved from core")
	}
}
