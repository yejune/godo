package assembler

import (
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
	FilesWritten  int
	SlotsResolved int
	SlotsUnfilled int
	AgentsPatched int
	Warnings      []string
	Files         []string
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
//  4. Copy persona commands
//  5. Copy persona hook scripts
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

	// Step 4: Copy persona CLAUDE.md.
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
		result.Files = append(result.Files, relPath)

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
// commands, hook scripts) to the output directory.
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
