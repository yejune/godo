package assembler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/do-focus/convert/internal/model"
	"github.com/do-focus/convert/internal/template"
)

// AssembleResult contains the summary of a full assembly run.
type AssembleResult struct {
	FilesWritten    int
	SlotsResolved   int
	SlotsUnfilled   int
	AgentsPatched   int
	SkillsMapped    int
	Warnings        []string
	Files           []string
}

// Assembler orchestrates the full assembly pipeline: core templates + persona
// manifest -> deployable .claude/ directory.
type Assembler struct {
	coreDir    string
	personaDir string
	outputDir  string
	manifest   *model.PersonaManifest
	registry   *template.Registry
}

// NewAssembler creates an Assembler with the given directories, manifest, and registry.
func NewAssembler(coreDir, personaDir, outputDir string, manifest *model.PersonaManifest, registry *template.Registry) *Assembler {
	return &Assembler{
		coreDir:    coreDir,
		personaDir: personaDir,
		outputDir:  outputDir,
		manifest:   manifest,
		registry:   registry,
	}
}

// Assemble runs the full assembly pipeline:
//  1. Copy core files to output, filling slots with persona content
//  2. Apply agent patches (append/remove skills, append content)
//  3. Copy persona-only files (agents, skills, rules, styles)
//  4. Apply skill mappings to all agent files in output
//  5. Merge settings.json (core + persona settings/hooks + manifest.Settings)
//  6. Copy persona CLAUDE.md
func (a *Assembler) Assemble() (*AssembleResult, error) {
	result := &AssembleResult{}
	merger := NewMerger(a.coreDir, a.personaDir, a.outputDir, a.manifest, a.registry)

	// Step 1: Walk core directory and merge each file to output.
	if err := a.copyCoreFiles(merger, result); err != nil {
		return nil, err
	}

	// Step 2: Apply agent patches.
	if err := a.applyAgentPatches(merger, result); err != nil {
		return nil, err
	}

	// Step 3: Copy persona-only files.
	if err := a.copyPersonaFiles(merger, result); err != nil {
		return nil, err
	}

	// Step 4: Apply skill mappings to all agent files.
	if err := a.applySkillMappings(merger, result); err != nil {
		return nil, err
	}

	// Step 5: Merge settings.json.
	if err := a.mergeSettings(merger, result); err != nil {
		return nil, err
	}

	// Step 6: Copy persona CLAUDE.md.
	if err := a.copyClaudeMD(merger, result); err != nil {
		return nil, err
	}

	sort.Strings(result.Files)
	return result, nil
}

// copyCoreFiles walks the core directory and merges each file to output,
// filling slot markers with persona content.
func (a *Assembler) copyCoreFiles(merger *Merger, result *AssembleResult) error {
	if _, err := os.Stat(a.coreDir); os.IsNotExist(err) {
		return nil
	}

	return filepath.Walk(a.coreDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(a.coreDir, path)
		if err != nil {
			return err
		}

		mergeResult, err := merger.MergeFile(relPath)
		if err != nil {
			return err
		}

		result.FilesWritten += mergeResult.FilesWritten
		result.SlotsResolved += mergeResult.SlotsResolved
		result.Warnings = append(result.Warnings, mergeResult.Warnings...)
		// Use remapped output path (brand prefix restored for skill dirs).
		outPath := mergeResult.OutputPath
		if outPath == "" {
			outPath = relPath
		}
		result.Files = append(result.Files, outPath)

		return nil
	})
}

// applyAgentPatches applies persona patches to core agent files that have
// already been copied to the output directory.
func (a *Assembler) applyAgentPatches(merger *Merger, result *AssembleResult) error {
	if a.manifest == nil || len(a.manifest.AgentPatches) == 0 {
		return nil
	}

	for relPath := range a.manifest.AgentPatches {
		if err := merger.PatchAgent(relPath); err != nil {
			return err
		}
		result.AgentsPatched++
	}

	return nil
}

// copyPersonaFiles copies persona-only files (agents, skills, rules, styles,
// commands, hook scripts) and any additional persona assets (scripts, templates,
// schemas, etc.) to the output directory.
func (a *Assembler) copyPersonaFiles(merger *Merger, result *AssembleResult) error {
	if a.manifest == nil {
		return nil
	}

	// Collect all persona-only file lists.
	var files []string
	files = append(files, a.manifest.Agents...)
	files = append(files, a.manifest.Skills...)
	files = append(files, a.manifest.Rules...)
	files = append(files, a.manifest.Styles...)
	files = append(files, a.manifest.Commands...)
	files = append(files, a.manifest.HookScripts...)

	// Track which files are already handled via the named slices.
	handled := make(map[string]bool, len(files))
	for _, relPath := range files {
		handled[relPath] = true
	}

	for _, relPath := range files {
		if _, err := merger.CopyPersonaFile(relPath); err != nil {
			// Check if file was already written by core copy (skip duplicate).
			if strings.Contains(err.Error(), "read persona file") {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("persona file %q not found, skipping", relPath))
				continue
			}
			return err
		}
		result.FilesWritten++
		result.Files = append(result.Files, relPath)
	}

	// Copy additional persona assets from PersonaFiles that aren't in named
	// slices (e.g., scripts/, templates/, schemas/ inside skill directories).
	// Also skip special files handled by other assembly steps.
	for relPath := range a.manifest.PersonaFiles {
		if handled[relPath] {
			continue
		}
		// Skip files handled by other steps (settings.json, CLAUDE.md).
		if relPath == "settings.json" || relPath == a.manifest.ClaudeMD {
			continue
		}
		if _, err := merger.CopyPersonaFile(relPath); err != nil {
			if strings.Contains(err.Error(), "read persona file") {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("persona asset %q not found, skipping", relPath))
				continue
			}
			return err
		}
		result.FilesWritten++
		result.Files = append(result.Files, relPath)
	}

	return nil
}

// applySkillMappings applies manifest.SkillMappings to all agent .md files
// in the output directory. This is a global replacement that runs after all
// files are assembled and agent patches are applied.
func (a *Assembler) applySkillMappings(merger *Merger, result *AssembleResult) error {
	if a.manifest == nil || len(a.manifest.SkillMappings) == 0 {
		return nil
	}

	mapped, err := merger.ApplySkillMappings()
	if err != nil {
		return err
	}
	result.SkillsMapped = mapped
	return nil
}

// mergeSettings handles settings.json assembly. Core settings.json is used as
// the base, persona settings.json (if present) is merged on top, then
// manifest.Settings is applied with deep merge for "env" and override for
// other top-level keys.
func (a *Assembler) mergeSettings(merger *Merger, result *AssembleResult) error {
	if a.manifest == nil {
		return nil
	}

	coreSettingsPath := filepath.Join(a.coreDir, "settings.json")
	personaSettingsPath := filepath.Join(a.personaDir, "settings.json")

	coreExists := true
	if _, err := os.Stat(coreSettingsPath); os.IsNotExist(err) {
		coreExists = false
	}
	personaExists := true
	if _, err := os.Stat(personaSettingsPath); os.IsNotExist(err) {
		personaExists = false
	}

	if !coreExists && !personaExists && len(a.manifest.Settings) == 0 && len(a.manifest.Hooks) == 0 {
		return nil
	}

	// Build base settings from core.
	var settings map[string]interface{}
	if coreExists {
		data, err := os.ReadFile(coreSettingsPath)
		if err != nil {
			return &model.ErrAssembly{
				Phase:   "merge_settings",
				File:    "settings.json",
				Message: fmt.Sprintf("read core settings: %v", err),
			}
		}
		if err := json.Unmarshal(data, &settings); err != nil {
			return &model.ErrAssembly{
				Phase:   "merge_settings",
				File:    "settings.json",
				Message: fmt.Sprintf("parse core settings: %v", err),
			}
		}
	} else {
		settings = make(map[string]interface{})
	}

	// Merge persona settings.json on top if present.
	if personaExists {
		data, err := os.ReadFile(personaSettingsPath)
		if err != nil {
			return &model.ErrAssembly{
				Phase:   "merge_settings",
				File:    "settings.json",
				Message: fmt.Sprintf("read persona settings: %v", err),
			}
		}
		var personaSettings map[string]interface{}
		if err := json.Unmarshal(data, &personaSettings); err != nil {
			return &model.ErrAssembly{
				Phase:   "merge_settings",
				File:    "settings.json",
				Message: fmt.Sprintf("parse persona settings: %v", err),
			}
		}
		mergeSettingsMap(settings, personaSettings)
	}

	// Inject persona hooks if present.
	if len(a.manifest.Hooks) > 0 {
		settings["hooks"] = a.manifest.Hooks
	}

	// Apply manifest.Settings overrides.
	if len(a.manifest.Settings) > 0 {
		mergeSettingsMap(settings, a.manifest.Settings)
	}

	buf, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return &model.ErrAssembly{
			Phase:   "merge_settings",
			File:    "settings.json",
			Message: fmt.Sprintf("marshal merged settings: %v", err),
		}
	}
	buf = append(buf, '\n')

	dstPath := filepath.Join(a.outputDir, "settings.json")
	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return &model.ErrAssembly{
			Phase:   "merge_settings",
			File:    "settings.json",
			Message: fmt.Sprintf("create output dir: %v", err),
		}
	}

	if err := os.WriteFile(dstPath, buf, 0o644); err != nil {
		return &model.ErrAssembly{
			Phase:   "merge_settings",
			File:    "settings.json",
			Message: fmt.Sprintf("write merged settings: %v", err),
		}
	}

	result.FilesWritten++
	result.Files = append(result.Files, "settings.json")
	return nil
}

// copyClaudeMD copies the persona's CLAUDE.md to the output directory root.
func (a *Assembler) copyClaudeMD(merger *Merger, result *AssembleResult) error {
	if a.manifest == nil || a.manifest.ClaudeMD == "" {
		return nil
	}

	if _, err := merger.CopyPersonaFile(a.manifest.ClaudeMD); err != nil {
		return &model.ErrAssembly{
			Phase:   "copy_claude_md",
			File:    a.manifest.ClaudeMD,
			Message: fmt.Sprintf("copy persona CLAUDE.md: %v", err),
		}
	}

	result.FilesWritten++
	result.Files = append(result.Files, a.manifest.ClaudeMD)
	return nil
}
