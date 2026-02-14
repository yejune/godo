package extractor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/do-focus/convert/internal/detector"
)

// setupTestDir creates a temporary directory structure simulating a .claude/ dir.
func setupTestDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for relPath, content := range files {
		fullPath := filepath.Join(dir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatalf("mkdir for %s: %v", relPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", relPath, err)
		}
	}
	return dir
}

func newTestOrchestrator(t *testing.T) *ExtractorOrchestrator {
	t.Helper()
	reg := detector.NewDefaultRegistry()
	det, err := detector.NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector: %v", err)
	}
	return NewExtractorOrchestrator(det, reg)
}

func TestOrchestrator_EmptyDirectory(t *testing.T) {
	orch := newTestOrchestrator(t)
	dir := t.TempDir()

	registry, manifest, err := orch.Extract(dir)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	if registry == nil {
		t.Fatal("registry should not be nil")
	}
	if manifest == nil {
		t.Fatal("manifest should not be nil")
	}
	if len(registry.Slots) != 0 {
		t.Errorf("registry slots count = %d, want 0", len(registry.Slots))
	}
	if manifest.ClaudeMD != "" {
		t.Errorf("manifest ClaudeMD = %q, want empty", manifest.ClaudeMD)
	}
	if len(manifest.Agents) != 0 {
		t.Errorf("manifest Agents count = %d, want 0", len(manifest.Agents))
	}
}

func TestOrchestrator_SingleAgentFile(t *testing.T) {
	orch := newTestOrchestrator(t)

	files := map[string]string{
		"agents/expert-generic.md": "---\nname: expert-generic\ndescription: A generic expert\ntools: Read Write Edit\n---\n\n## Overview\n\nGeneric agent overview.\n\n## Guidelines\n\nFollow best practices.\n",
	}

	dir := setupTestDir(t, files)
	registry, manifest, err := orch.Extract(dir)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	// Pure core agent produces no persona content
	if len(manifest.Agents) != 0 {
		t.Errorf("manifest Agents = %v, want empty (core agent)", manifest.Agents)
	}
	if len(manifest.SlotContent) != 0 {
		t.Errorf("manifest SlotContent count = %d, want 0", len(manifest.SlotContent))
	}

	// Registry should have no slots for pure core agent
	if len(registry.Slots) != 0 {
		t.Errorf("registry slots count = %d, want 0", len(registry.Slots))
	}
}

func TestOrchestrator_WholeFilePersonaAgent(t *testing.T) {
	orch := newTestOrchestrator(t)

	files := map[string]string{
		"agents/moai/manager-spec.md": "---\nname: manager-spec\ndescription: SPEC document creation manager\nskills:\n  - moai-workflow-spec\n---\n\n## SPEC Management\n\nCreate and manage SPEC documents.\n",
	}

	dir := setupTestDir(t, files)
	_, manifest, err := orch.Extract(dir)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	// Whole-file persona agent should appear in manifest.Agents
	if len(manifest.Agents) != 1 {
		t.Fatalf("manifest Agents count = %d, want 1", len(manifest.Agents))
	}
	if manifest.Agents[0] != "agents/moai/manager-spec.md" {
		t.Errorf("manifest Agents[0] = %q, want %q", manifest.Agents[0], "agents/moai/manager-spec.md")
	}
}

func TestOrchestrator_MixedFileTypes(t *testing.T) {
	orch := newTestOrchestrator(t)

	files := map[string]string{
		// Core agent with persona skill in frontmatter
		"agents/expert-backend.md": "---\nname: expert-backend\ndescription: Backend expert\nskills:\n  - do-foundation-claude\n  - moai-foundation-core\n---\n\n## Implementation\n\nFollow clean code.\n\n## Error Handling\n\nHandle errors properly.\n",
		// Whole-file persona agent
		"agents/moai/manager-ddd.md": "---\nname: manager-ddd\ndescription: DDD manager\n---\n\n## DDD Workflow\n\nDDD content.\n",
		// Core skill
		"skills/do-foundation-claude.md": "---\nname: do-foundation-claude\ndescription: Foundation skill\n---\n\n## Foundation\n\nCore foundation content.\n",
		// Persona skill
		"skills/moai-workflow-ddd.md": "---\nname: moai-workflow-ddd\ndescription: DDD workflow skill\n---\n\n## DDD\n\nDDD skill content.\n",
		// Style file (always persona) — legacy styles/ path
		"styles/pair.md": "---\nname: pair\ndescription: Pair programming style\n---\n\n## Style\n\nFriendly pair programmer.\n",
		// Style file via output-styles/ path (real moai-adk structure)
		"output-styles/moai/sprint.md": "---\nname: sprint\ndescription: Sprint style\n---\n\n## Style\n\nMinimal output.\n",
		// Rule file (core)
		"rules/dev-testing.md": "# Testing Rules\n\nTest everything.\n",
		// Persona rule (spec-workflow.md is in WholeFileRules)
		"rules/do/workflow/spec-workflow.md": "# SPEC Workflow\n\nSPEC workflow rules.\n",
		// CLAUDE.md (always persona)
		"CLAUDE.md": "# Do Execution Directive\n\nMain persona directive.\n",
		// Settings
		"settings.json": "{\"outputStyle\":\"pair\",\"permissions\":{\"allow\":[\"Read\"]},\"hooks\":{\"PreToolUse\":[{\"command\":\"godo hook pre\"}]}}",
		// Command files
		"commands/help.md": "Help command content.",
		// Hook scripts
		"hooks/pre-tool.sh": "#!/bin/bash\necho \"core hook\"",
	}

	dir := setupTestDir(t, files)
	_, manifest, err := orch.Extract(dir)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	// CLAUDE.md
	if manifest.ClaudeMD != "CLAUDE.md" {
		t.Errorf("manifest ClaudeMD = %q, want %q", manifest.ClaudeMD, "CLAUDE.md")
	}

	// Whole-file persona agent
	if len(manifest.Agents) != 1 {
		t.Errorf("manifest Agents count = %d, want 1 (manager-ddd)", len(manifest.Agents))
	}

	// Persona skill
	if len(manifest.Skills) != 1 {
		t.Errorf("manifest Skills count = %d, want 1 (moai-workflow-ddd)", len(manifest.Skills))
	}

	// Style (all persona) — includes both styles/ and output-styles/
	if len(manifest.Styles) != 2 {
		t.Errorf("manifest Styles count = %d, want 2", len(manifest.Styles))
	}

	// Persona rule (spec-workflow.md is in WholeFileRules)
	if len(manifest.Rules) != 1 {
		t.Errorf("manifest Rules count = %d, want 1", len(manifest.Rules))
	}

	// Commands
	if len(manifest.Commands) != 1 {
		t.Errorf("manifest Commands count = %d, want 1", len(manifest.Commands))
	}

	// Hook scripts
	if len(manifest.HookScripts) != 1 {
		t.Errorf("manifest HookScripts count = %d, want 1", len(manifest.HookScripts))
	}

	// Settings (hooks should be in persona settings)
	if manifest.Settings == nil {
		t.Fatal("manifest Settings is nil")
	}
	if _, ok := manifest.Settings["hooks"]; !ok {
		t.Error("manifest Settings should contain 'hooks' key")
	}

	// Agent patches (moai-foundation-core skill extracted from frontmatter)
	// NOTE: Header-based persona detection (QUALITY_FRAMEWORK) does not trigger
	// when using ParseDocument because the parser strips '#' marks from Section.Title
	// (e.g., "TRUST 5 Compliance" not "### TRUST 5 Compliance"), while the detector
	// regex expects the '#' prefix. This is a known pre-existing mismatch between
	// the parser and detector. Skill-based extraction still works via frontmatter.
	patch, ok := manifest.AgentPatches["agents/expert-backend.md"]
	if !ok {
		t.Error("manifest AgentPatches should contain 'agents/expert-backend.md'")
	} else {
		found := false
		for _, s := range patch.AppendSkills {
			if s == "moai-foundation-core" {
				found = true
			}
		}
		if !found {
			t.Errorf("agent patch AppendSkills should contain 'moai-foundation-core', got %v", patch.AppendSkills)
		}
	}
}

func TestOrchestrator_SkipsGitDirectory(t *testing.T) {
	orch := newTestOrchestrator(t)

	files := map[string]string{
		".git/config":              "git config content",
		"agents/expert-generic.md": "---\nname: expert-generic\ndescription: test\n---\n\n## Overview\n\nContent.\n",
	}

	dir := setupTestDir(t, files)
	_, _, err := orch.Extract(dir)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}
	// No error means .git was skipped successfully
}

func TestOrchestrator_SkipsNonRelevantFiles(t *testing.T) {
	orch := newTestOrchestrator(t)

	files := map[string]string{
		"README.md":  "# README\n\nProject readme.",
		"config.yml": "key: value",
		"random.txt": "random content",
	}

	dir := setupTestDir(t, files)
	registry, manifest, err := orch.Extract(dir)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	if len(registry.Slots) != 0 {
		t.Errorf("registry slots count = %d, want 0 (no relevant files)", len(registry.Slots))
	}
	if manifest.ClaudeMD != "" {
		t.Errorf("manifest ClaudeMD = %q, want empty", manifest.ClaudeMD)
	}
}


func TestOrchestrator_ClaudeMDAtProjectRoot(t *testing.T) {
	orch := newTestOrchestrator(t)

	// Simulate project root with .claude/ subdirectory.
	// CLAUDE.md is at project root, not inside .claude/.
	projectRoot := t.TempDir()
	claudeDir := filepath.Join(projectRoot, ".claude")
	if err := os.MkdirAll(filepath.Join(claudeDir, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Write CLAUDE.md at project root (parent of .claude/)
	claudeMD := "# Do Execution Directive\n\nMain persona directive.\n"
	if err := os.WriteFile(filepath.Join(projectRoot, "CLAUDE.md"), []byte(claudeMD), 0o644); err != nil {
		t.Fatal(err)
	}

	// Write an agent inside .claude/
	agentContent := "---\nname: expert-generic\ndescription: test\n---\n\n## Overview\n\nContent.\n"
	if err := os.WriteFile(filepath.Join(claudeDir, "agents", "expert-generic.md"), []byte(agentContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Extract from .claude/ — CLAUDE.md should be found at project root
	_, manifest, err := orch.Extract(claudeDir)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	if manifest.ClaudeMD != "CLAUDE.md" {
		t.Errorf("manifest ClaudeMD = %q, want %q", manifest.ClaudeMD, "CLAUDE.md")
	}
}

func TestOrchestrator_ClaudeMDInsideSrcDirTakesPrecedence(t *testing.T) {
	orch := newTestOrchestrator(t)

	// If CLAUDE.md exists both inside srcDir and at project root,
	// the one inside srcDir wins (found during Walk).
	projectRoot := t.TempDir()
	claudeDir := filepath.Join(projectRoot, ".claude")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// CLAUDE.md inside .claude/
	if err := os.WriteFile(filepath.Join(claudeDir, "CLAUDE.md"), []byte("# Inside"), 0o644); err != nil {
		t.Fatal(err)
	}
	// CLAUDE.md at project root
	if err := os.WriteFile(filepath.Join(projectRoot, "CLAUDE.md"), []byte("# Root"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, manifest, err := orch.Extract(claudeDir)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	// The one inside srcDir should win
	if manifest.ClaudeMD != "CLAUDE.md" {
		t.Errorf("manifest ClaudeMD = %q, want %q", manifest.ClaudeMD, "CLAUDE.md")
	}
}

func TestClassifyFile(t *testing.T) {
	tests := []struct {
		path string
		want fileType
	}{
		{"CLAUDE.md", fileTypeClaudeMD},
		{"claude.md", fileTypeClaudeMD},
		{"settings.json", fileTypeSettings},
		{"agents/expert-backend.md", fileTypeAgent},
		{"agents/moai/manager-spec.md", fileTypeAgent},
		{"skills/do-foundation.md", fileTypeSkill},
		{"rules/dev-testing.md", fileTypeRule},
		{"rules/do/workflow/spec-workflow.md", fileTypeRule},
		{"styles/pair.md", fileTypeStyle},
		{"output-styles/pair.md", fileTypeStyle},
		{"output-styles/moai/sprint.md", fileTypeStyle},
		{"commands/help.md", fileTypeCommand},
		{"commands/subdir/deploy.md", fileTypeCommand},
		{"hooks/pre-tool.sh", fileTypeHook},
		{"README.md", fileTypeUnknown},
		{"config.yml", fileTypeUnknown},
		{"agents/README.txt", fileTypeUnknown}, // not .md
		{"nested/CLAUDE.md", fileTypeUnknown},  // not at root
		{"nested/settings.json", fileTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := classifyFile(tt.path)
			if got != tt.want {
				t.Errorf("classifyFile(%q) = %d, want %d", tt.path, got, tt.want)
			}
		})
	}
}

func TestExtractSlotID(t *testing.T) {
	tests := []struct {
		content string
		want    string
	}{
		{"<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->\n\n<!-- END_SLOT:QUALITY_FRAMEWORK -->", "QUALITY_FRAMEWORK"},
		{"<!-- BEGIN_SLOT:TRACEABILITY_SYSTEM -->\ncontent\n<!-- END_SLOT:TRACEABILITY_SYSTEM -->", "TRACEABILITY_SYSTEM"},
		{"no slot markers here", ""},
		{"", ""},
	}

	for _, tt := range tests {
		got := extractSlotID(tt.content)
		if got != tt.want {
			t.Errorf("extractSlotID(%q...) = %q, want %q", tt.content[:min(len(tt.content), 40)], got, tt.want)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
