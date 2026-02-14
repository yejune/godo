package assembler

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/do-focus/convert/internal/detector"
	"github.com/do-focus/convert/internal/extractor"
	"github.com/do-focus/convert/internal/model"
	"github.com/do-focus/convert/internal/template"
)

// setupSourceDir creates a temporary directory structure simulating a .claude/ dir.
func setupSourceDir(t *testing.T, files map[string]string) string {
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

// rekeyAgentPatches transforms AgentPatches keys from extractor format (agent
// name, e.g. "expert-backend") to assembler format (relative file path, e.g.
// "agents/expert-backend.md"). The extractor keys by frontmatter name; the
// assembler expects file paths relative to the core directory.
func rekeyAgentPatches(manifest *model.PersonaManifest, sourceFiles map[string]string) {
	if len(manifest.AgentPatches) == 0 {
		return
	}

	// Build name-to-path mapping from source agent files.
	nameToPath := make(map[string]string)
	for relPath := range sourceFiles {
		if strings.HasPrefix(relPath, "agents/") && strings.HasSuffix(relPath, ".md") {
			base := filepath.Base(relPath)
			name := strings.TrimSuffix(base, ".md")
			nameToPath[name] = relPath
		}
	}

	rekeyed := make(map[string]*model.AgentPatch, len(manifest.AgentPatches))
	for key, patch := range manifest.AgentPatches {
		if path, ok := nameToPath[key]; ok {
			rekeyed[path] = patch
		} else {
			// Key is already a path or unknown; keep as-is.
			rekeyed[key] = patch
		}
	}
	manifest.AgentPatches = rekeyed
}

// buildCoreAndPersonaDirs takes the extraction results and source directory,
// then builds separate core and persona directories that the assembler expects.
// Core dir: files NOT classified as persona-only, plus registry.yaml.
// Persona dir: files that ARE persona-only (listed in manifest).
func buildCoreAndPersonaDirs(
	t *testing.T,
	srcDir string,
	registry *template.Registry,
	manifest *model.PersonaManifest,
	sourceFiles map[string]string,
) (coreDir, personaDir string) {
	t.Helper()

	coreDir = t.TempDir()
	personaDir = t.TempDir()

	// Build set of persona-only file paths for quick lookup.
	personaSet := make(map[string]bool)
	for _, p := range manifest.Agents {
		personaSet[p] = true
	}
	for _, p := range manifest.Skills {
		personaSet[p] = true
	}
	for _, p := range manifest.Rules {
		personaSet[p] = true
	}
	for _, p := range manifest.Styles {
		personaSet[p] = true
	}
	for _, p := range manifest.Commands {
		personaSet[p] = true
	}
	for _, p := range manifest.HookScripts {
		personaSet[p] = true
	}
	if manifest.ClaudeMD != "" {
		personaSet[manifest.ClaudeMD] = true
	}

	// Copy files to core or persona based on classification.
	for relPath, content := range sourceFiles {
		if relPath == "settings.json" {
			// Settings are split: core gets the non-hooks part, persona gets hooks.
			coreSettings, _, err := extractor.ExtractSettings([]byte(content))
			if err != nil {
				t.Fatalf("extract settings: %v", err)
			}
			buf, err := json.MarshalIndent(coreSettings, "", "  ")
			if err != nil {
				t.Fatalf("marshal core settings: %v", err)
			}
			buf = append(buf, '\n')
			corePath := filepath.Join(coreDir, "settings.json")
			if err := os.MkdirAll(filepath.Dir(corePath), 0o755); err != nil {
				t.Fatalf("mkdir core settings: %v", err)
			}
			if err := os.WriteFile(corePath, buf, 0o644); err != nil {
				t.Fatalf("write core settings: %v", err)
			}
			continue
		}

		if personaSet[relPath] {
			// Copy to persona dir.
			dstPath := filepath.Join(personaDir, relPath)
			if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
				t.Fatalf("mkdir persona %s: %v", relPath, err)
			}
			if err := os.WriteFile(dstPath, []byte(content), 0o644); err != nil {
				t.Fatalf("write persona %s: %v", relPath, err)
			}
		} else {
			// Copy to core dir.
			dstPath := filepath.Join(coreDir, relPath)
			if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
				t.Fatalf("mkdir core %s: %v", relPath, err)
			}
			if err := os.WriteFile(dstPath, []byte(content), 0o644); err != nil {
				t.Fatalf("write core %s: %v", relPath, err)
			}
		}
	}

	// Save registry.yaml to core dir.
	if err := registry.Save(coreDir); err != nil {
		t.Fatalf("save registry: %v", err)
	}

	// Re-key AgentPatches from extractor format (agent name) to assembler
	// format (relative file path). This bridges the key format mismatch
	// between the two pipeline stages.
	rekeyAgentPatches(manifest, sourceFiles)

	return coreDir, personaDir
}

// newExtractOrchestrator creates an ExtractorOrchestrator wired with default patterns.
func newExtractOrchestrator(t *testing.T) *extractor.ExtractorOrchestrator {
	t.Helper()
	reg := detector.NewDefaultRegistry()
	det, err := detector.NewPersonaDetector(reg)
	if err != nil {
		t.Fatalf("NewPersonaDetector: %v", err)
	}
	return extractor.NewExtractorOrchestrator(det, reg)
}

// fullSourceFiles returns a complete .claude/ directory structure for E2E testing.
// It includes all file types: agents (core + persona), skills (core + persona),
// rules (core + persona), styles, settings, commands, hooks, and CLAUDE.md.
func fullSourceFiles() map[string]string {
	return map[string]string{
		// Core agent with persona skills in frontmatter.
		"agents/expert-backend.md": "---\nname: expert-backend\ndescription: >\n  Backend architecture specialist.\ntools: Read Write Edit Grep Glob Bash\nmodel: inherit\npermissionMode: acceptEdits\nskills:\n  - do-domain-backend\n  - do-lang-go\n  - moai-foundation-core\n  - moai-foundation-quality\n---\n\n## Role\n\nBackend development specialist for API design and server-side implementation.\n\n## Implementation Guidelines\n\n- Follow project conventions for error handling and logging\n- Use dependency injection for testability\n",

		// Whole-file persona agent (name in WholeFileAgents).
		"agents/moai/manager-spec.md": "---\nname: manager-spec\ndescription: >\n  SPEC document lifecycle manager.\ntools: Read Write Edit Grep Glob Bash\nskills:\n  - moai-workflow-spec\n---\n\n## SPEC Workflow\n\nManages the Plan phase of the SPEC workflow.\n",

		// Core skill (name NOT in WholeFileSkills).
		"skills/do-domain-backend.md": "---\nname: do-domain-backend\ndescription: Backend domain knowledge skill.\n---\n\n## Backend Patterns\n\nService layer patterns and conventions.\n",

		// Persona skill (name in WholeFileSkills).
		"skills/moai-foundation-core.md": "---\nname: moai-foundation-core\ndescription: Core foundation skill for methodology.\n---\n\n## Foundation\n\nCore foundation content for persona.\n",

		// Core rule (not in WholeFileRules).
		"rules/dev-testing.md": "# Testing Rules\n\n## Real DB Only\n\nAll database tests use the real database.\n\n## Test Quality\n\nWrite meaningful test names.\n",

		// Persona rule (filename matches WholeFileRules).
		"rules/do/workflow/spec-workflow.md": "# SPEC Workflow\n\nSPEC workflow rules for managing specifications.\n",

		// Style file (ALL styles are persona).
		"styles/pair.md": "---\nname: pair\ndescription: Friendly pair programming style.\n---\n\n## Style\n\nCollaborative and friendly pair programmer tone.\n",

		// Settings with core permissions + persona hooks.
		"settings.json": "{\"outputStyle\":\"pair\",\"permissions\":{\"allow\":[\"Read\",\"Write\",\"Edit\"]},\"hooks\":{\"PreToolUse\":[{\"command\":\"godo hook pre-tool\"}]}}",

		// Command file (persona).
		"commands/commit.md": "# Commit Command\n\nRuns the commit workflow with conventional commit messages.\n",

		// Hook script (persona).
		"hooks/pre-commit.sh": "#!/bin/bash\necho \"Running pre-commit checks\"\n",

		// CLAUDE.md (always persona).
		"CLAUDE.md": "# Do Execution Directive\n\n## Do/Focus/Team\n\nMain persona directive content.\n",
	}
}

// TestE2E_ExtractAssembleRoundTrip runs the full extract -> assemble round-trip.
// Starting from a .claude/ directory, it extracts into core + persona, then
// assembles back and verifies the output matches the original (or is
// semantically equivalent for files that undergo transformation).
func TestE2E_ExtractAssembleRoundTrip(t *testing.T) {
	sourceFiles := fullSourceFiles()
	srcDir := setupSourceDir(t, sourceFiles)

	// Step 1: Extract.
	orch := newExtractOrchestrator(t)
	registry, manifest, err := orch.Extract(srcDir)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	// Step 2: Build core and persona directories.
	coreDir, personaDir := buildCoreAndPersonaDirs(t, srcDir, registry, manifest, sourceFiles)

	// Step 3: Assemble.
	outputDir := t.TempDir()
	asm := NewAssembler(coreDir, personaDir, outputDir, manifest, registry)
	result, err := asm.Assemble()
	if err != nil {
		t.Fatalf("Assemble() error: %v", err)
	}

	t.Logf("Assembly result: %d files written, %d slots resolved, %d agents patched",
		result.FilesWritten, result.SlotsResolved, result.AgentsPatched)

	// Step 4: Verify output files exist and content is preserved.

	// 4a. Persona-only files should be byte-for-byte identical.
	personaOnlyFiles := []string{
		"styles/pair.md",
		"commands/commit.md",
		"hooks/pre-commit.sh",
		"CLAUDE.md",
		"agents/moai/manager-spec.md",
		"skills/moai-foundation-core.md",
		"rules/do/workflow/spec-workflow.md",
	}
	for _, relPath := range personaOnlyFiles {
		outputData, err := os.ReadFile(filepath.Join(outputDir, relPath))
		if err != nil {
			t.Errorf("persona file %q not found in output: %v", relPath, err)
			continue
		}
		original := sourceFiles[relPath]
		if string(outputData) != original {
			t.Errorf("persona file %q content mismatch.\noriginal:\n%s\noutput:\n%s",
				relPath, original, string(outputData))
		}
	}

	// 4b. Core-only files should be preserved (content identical, no slots affected).
	coreOnlyFiles := []string{
		"rules/dev-testing.md",
	}
	for _, relPath := range coreOnlyFiles {
		outputData, err := os.ReadFile(filepath.Join(outputDir, relPath))
		if err != nil {
			t.Errorf("core file %q not found in output: %v", relPath, err)
			continue
		}
		original := sourceFiles[relPath]
		if string(outputData) != original {
			t.Errorf("core file %q content mismatch.\noriginal:\n%s\noutput:\n%s",
				relPath, original, string(outputData))
		}
	}

	// 4c. Core skill should be preserved.
	skillData, err := os.ReadFile(filepath.Join(outputDir, "skills", "do-domain-backend.md"))
	if err != nil {
		t.Fatalf("core skill not found in output: %v", err)
	}
	if string(skillData) != sourceFiles["skills/do-domain-backend.md"] {
		t.Errorf("core skill content mismatch.\noriginal:\n%s\noutput:\n%s",
			sourceFiles["skills/do-domain-backend.md"], string(skillData))
	}

	// 4d. Mixed agent (expert-backend) should have persona skills restored.
	agentData, err := os.ReadFile(filepath.Join(outputDir, "agents", "expert-backend.md"))
	if err != nil {
		t.Fatalf("mixed agent not found in output: %v", err)
	}
	agentOutput := string(agentData)

	// The patched agent should contain the persona skills that were extracted.
	if !strings.Contains(agentOutput, "moai-foundation-core") {
		t.Error("assembled agent missing persona skill 'moai-foundation-core'")
	}
	if !strings.Contains(agentOutput, "moai-foundation-quality") {
		t.Error("assembled agent missing persona skill 'moai-foundation-quality'")
	}
	// Core skills should be preserved.
	if !strings.Contains(agentOutput, "do-domain-backend") {
		t.Error("assembled agent missing core skill 'do-domain-backend'")
	}
	if !strings.Contains(agentOutput, "do-lang-go") {
		t.Error("assembled agent missing core skill 'do-lang-go'")
	}
	// Core body sections should be preserved.
	if !strings.Contains(agentOutput, "## Role") {
		t.Error("assembled agent missing '## Role' section")
	}
	if !strings.Contains(agentOutput, "Backend development specialist") {
		t.Error("assembled agent missing body content")
	}

	// 4e. Settings.json in output should have core fields.
	settingsData, err := os.ReadFile(filepath.Join(outputDir, "settings.json"))
	if err != nil {
		t.Fatalf("settings.json not found in output: %v", err)
	}
	var settingsMap map[string]interface{}
	if err := json.Unmarshal(settingsData, &settingsMap); err != nil {
		t.Fatalf("parse output settings.json: %v", err)
	}
	if _, ok := settingsMap["outputStyle"]; !ok {
		t.Error("output settings.json missing 'outputStyle'")
	}
	if _, ok := settingsMap["permissions"]; !ok {
		t.Error("output settings.json missing 'permissions'")
	}

	// 4f. Result files list should contain expected entries.
	if len(result.Files) == 0 {
		t.Error("result.Files should not be empty")
	}
	// Files should be sorted (verified by Assembler).
	sorted := make([]string, len(result.Files))
	copy(sorted, result.Files)
	sort.Strings(sorted)
	for i := range sorted {
		if result.Files[i] != sorted[i] {
			t.Errorf("result.Files not sorted at index %d: got %q, expected %q",
				i, result.Files[i], sorted[i])
			break
		}
	}
}

// TestE2E_ExtractAssemblePreservesContent focuses on content preservation
// guarantees across the extract-assemble pipeline.
func TestE2E_ExtractAssemblePreservesContent(t *testing.T) {
	sourceFiles := fullSourceFiles()
	srcDir := setupSourceDir(t, sourceFiles)

	// Extract.
	orch := newExtractOrchestrator(t)
	registry, manifest, err := orch.Extract(srcDir)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	// Build dirs and assemble.
	coreDir, personaDir := buildCoreAndPersonaDirs(t, srcDir, registry, manifest, sourceFiles)
	outputDir := t.TempDir()
	asm := NewAssembler(coreDir, personaDir, outputDir, manifest, registry)
	if _, err := asm.Assemble(); err != nil {
		t.Fatalf("Assemble() error: %v", err)
	}

	// Test 1: YAML frontmatter preserved in core skill.
	t.Run("CoreSkillFrontmatterPreserved", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join(outputDir, "skills", "do-domain-backend.md"))
		if err != nil {
			t.Fatalf("read output skill: %v", err)
		}
		output := string(data)
		if !strings.Contains(output, "name: do-domain-backend") {
			t.Error("output skill missing 'name: do-domain-backend' in frontmatter")
		}
		if !strings.Contains(output, "description:") {
			t.Error("output skill missing 'description' in frontmatter")
		}
		if !strings.Contains(output, "## Backend Patterns") {
			t.Error("output skill missing '## Backend Patterns' section header")
		}
		if !strings.Contains(output, "Service layer patterns and conventions.") {
			t.Error("output skill missing body content")
		}
	})

	// Test 2: Persona-only files are byte-for-byte identical.
	t.Run("PersonaFilesBytePerfect", func(t *testing.T) {
		personaFiles := map[string]string{
			"styles/pair.md":                    sourceFiles["styles/pair.md"],
			"CLAUDE.md":                         sourceFiles["CLAUDE.md"],
			"commands/commit.md":                sourceFiles["commands/commit.md"],
			"hooks/pre-commit.sh":               sourceFiles["hooks/pre-commit.sh"],
			"agents/moai/manager-spec.md":       sourceFiles["agents/moai/manager-spec.md"],
			"skills/moai-foundation-core.md":    sourceFiles["skills/moai-foundation-core.md"],
			"rules/do/workflow/spec-workflow.md": sourceFiles["rules/do/workflow/spec-workflow.md"],
		}

		for relPath, expected := range personaFiles {
			data, err := os.ReadFile(filepath.Join(outputDir, relPath))
			if err != nil {
				t.Errorf("persona file %q missing in output: %v", relPath, err)
				continue
			}
			if string(data) != expected {
				t.Errorf("persona file %q byte mismatch.\nexpected len=%d, got len=%d",
					relPath, len(expected), len(data))
			}
		}
	})

	// Test 3: Markdown section content preserved in core rule.
	t.Run("CoreRuleContentPreserved", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join(outputDir, "rules", "dev-testing.md"))
		if err != nil {
			t.Fatalf("read output rule: %v", err)
		}
		output := string(data)
		if !strings.Contains(output, "# Testing Rules") {
			t.Error("missing '# Testing Rules' header")
		}
		if !strings.Contains(output, "## Real DB Only") {
			t.Error("missing '## Real DB Only' section")
		}
		if !strings.Contains(output, "## Test Quality") {
			t.Error("missing '## Test Quality' section")
		}
		if !strings.Contains(output, "Write meaningful test names.") {
			t.Error("missing body content under '## Test Quality'")
		}
	})

	// Test 4: Mixed agent has all skills after round-trip.
	t.Run("MixedAgentSkillsRestored", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join(outputDir, "agents", "expert-backend.md"))
		if err != nil {
			t.Fatalf("read output agent: %v", err)
		}
		output := string(data)

		wantSkills := []string{
			"do-domain-backend",
			"do-lang-go",
			"moai-foundation-core",
			"moai-foundation-quality",
		}
		for _, skill := range wantSkills {
			if !strings.Contains(output, skill) {
				t.Errorf("assembled agent missing skill %q", skill)
			}
		}
	})

	// Test 5: No extra whitespace or content loss in persona files.
	t.Run("NoExtraWhitespace", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join(outputDir, "styles", "pair.md"))
		if err != nil {
			t.Fatalf("read output style: %v", err)
		}
		original := sourceFiles["styles/pair.md"]
		if string(data) != original {
			origLines := strings.Split(original, "\n")
			outLines := strings.Split(string(data), "\n")
			maxLines := len(origLines)
			if len(outLines) > maxLines {
				maxLines = len(outLines)
			}
			for i := 0; i < maxLines; i++ {
				var origLine, outLine string
				if i < len(origLines) {
					origLine = origLines[i]
				}
				if i < len(outLines) {
					outLine = outLines[i]
				}
				if origLine != outLine {
					t.Errorf("line %d diff:\n  orig: %q\n  out:  %q", i+1, origLine, outLine)
				}
			}
		}
	})
}

// TestE2E_MultiplePersonaAssembly verifies that extracting once and assembling
// with different persona manifests produces independent, correct outputs.
func TestE2E_MultiplePersonaAssembly(t *testing.T) {
	sourceFiles := fullSourceFiles()
	srcDir := setupSourceDir(t, sourceFiles)

	// Extract once.
	orch := newExtractOrchestrator(t)
	registry, manifest, err := orch.Extract(srcDir)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}

	// Build core dir (shared across assemblies).
	coreDir, personaDirA := buildCoreAndPersonaDirs(t, srcDir, registry, manifest, sourceFiles)

	// Assembly A: use the extracted manifest as-is.
	outputA := t.TempDir()
	asmA := NewAssembler(coreDir, personaDirA, outputA, manifest, registry)
	resultA, err := asmA.Assemble()
	if err != nil {
		t.Fatalf("Assemble A error: %v", err)
	}

	// Assembly B: create a different persona with different CLAUDE.md and style.
	personaDirB := t.TempDir()
	claudeB := "# Different Persona\n\nThis is persona B.\n"
	if err := os.WriteFile(filepath.Join(personaDirB, "CLAUDE.md"), []byte(claudeB), 0o644); err != nil {
		t.Fatal(err)
	}
	stylesDirB := filepath.Join(personaDirB, "styles")
	if err := os.MkdirAll(stylesDirB, 0o755); err != nil {
		t.Fatal(err)
	}
	styleBContent := "---\nname: direct\ndescription: Direct expert style.\n---\n\n## Style\n\nNo-nonsense, direct communication.\n"
	if err := os.WriteFile(filepath.Join(stylesDirB, "direct.md"), []byte(styleBContent), 0o644); err != nil {
		t.Fatal(err)
	}

	manifestB := &model.PersonaManifest{
		Name:     "persona-b",
		ClaudeMD: "CLAUDE.md",
		Styles:   []string{"styles/direct.md"},
	}

	outputB := t.TempDir()
	asmB := NewAssembler(coreDir, personaDirB, outputB, manifestB, registry)
	resultB, err := asmB.Assemble()
	if err != nil {
		t.Fatalf("Assemble B error: %v", err)
	}

	// Verify assemblies are independent.
	t.Run("AssemblyA_HasOriginalClaudeMD", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join(outputA, "CLAUDE.md"))
		if err != nil {
			t.Fatalf("read A CLAUDE.md: %v", err)
		}
		if !strings.Contains(string(data), "Do Execution Directive") {
			t.Error("Assembly A should have original CLAUDE.md content")
		}
	})

	t.Run("AssemblyB_HasDifferentClaudeMD", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join(outputB, "CLAUDE.md"))
		if err != nil {
			t.Fatalf("read B CLAUDE.md: %v", err)
		}
		if !strings.Contains(string(data), "Different Persona") {
			t.Error("Assembly B should have persona B CLAUDE.md content")
		}
		if strings.Contains(string(data), "Do Execution Directive") {
			t.Error("Assembly B should NOT have persona A CLAUDE.md content")
		}
	})

	t.Run("AssemblyA_HasPairStyle", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join(outputA, "styles", "pair.md"))
		if err != nil {
			t.Fatalf("read A pair.md: %v", err)
		}
		if !strings.Contains(string(data), "pair programming") {
			t.Error("Assembly A should have pair.md style")
		}
	})

	t.Run("AssemblyB_HasDirectStyle", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join(outputB, "styles", "direct.md"))
		if err != nil {
			t.Fatalf("read B direct.md: %v", err)
		}
		if !strings.Contains(string(data), "No-nonsense") {
			t.Error("Assembly B should have direct.md style")
		}
	})

	t.Run("SharedCoreFilesIdentical", func(t *testing.T) {
		dataA, err := os.ReadFile(filepath.Join(outputA, "rules", "dev-testing.md"))
		if err != nil {
			t.Fatalf("read A core rule: %v", err)
		}
		dataB, err := os.ReadFile(filepath.Join(outputB, "rules", "dev-testing.md"))
		if err != nil {
			t.Fatalf("read B core rule: %v", err)
		}
		if string(dataA) != string(dataB) {
			t.Error("core rule files should be identical across assemblies")
		}
	})

	// Verify both completed successfully.
	if resultA.FilesWritten == 0 {
		t.Error("Assembly A wrote 0 files")
	}
	if resultB.FilesWritten == 0 {
		t.Error("Assembly B wrote 0 files")
	}
	t.Logf("Assembly A: %d files, Assembly B: %d files", resultA.FilesWritten, resultB.FilesWritten)
}
