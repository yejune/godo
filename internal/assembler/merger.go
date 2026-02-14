package assembler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/do-focus/convert/internal/model"
	"github.com/do-focus/convert/internal/parser"
	"github.com/do-focus/convert/internal/template"
)

// MergeResult contains the summary of a single merge operation.
type MergeResult struct {
	FilesWritten  int
	SlotsResolved int
	Warnings      []string
}

// Merger combines core template files with persona-specific content
// to produce a final assembled output directory.
type Merger struct {
	coreDir    string
	personaDir string
	outputDir  string
	manifest   *model.PersonaManifest
	registry   *template.Registry
	filler     *SlotFiller
}

// NewMerger creates a Merger with the given directories, manifest, and registry.
func NewMerger(coreDir, personaDir, outputDir string, manifest *model.PersonaManifest, registry *template.Registry) *Merger {
	return &Merger{
		coreDir:    coreDir,
		personaDir: personaDir,
		outputDir:  outputDir,
		manifest:   manifest,
		registry:   registry,
		filler:     NewSlotFiller(registry, manifest, personaDir),
	}
}

// MergeFile reads a core template file, fills slot markers with persona content,
// and writes the result to the output directory. The relPath is relative to coreDir.
//
// If the file contains no slot markers, it is copied as-is.
func (m *Merger) MergeFile(relPath string) (*MergeResult, error) {
	srcPath := filepath.Join(m.coreDir, relPath)
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, &model.ErrAssembly{
			Phase:   "merge",
			File:    relPath,
			Message: fmt.Sprintf("read core file: %v", err),
		}
	}

	content := string(data)
	filled, resolved, warnings := m.filler.FillContent(content)

	result := &MergeResult{
		FilesWritten:  1,
		SlotsResolved: len(resolved),
		Warnings:      warnings,
	}

	dstPath := filepath.Join(m.outputDir, relPath)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return nil, &model.ErrAssembly{
			Phase:   "merge",
			File:    relPath,
			Message: fmt.Sprintf("create output dir: %v", err),
		}
	}

	if err := os.WriteFile(dstPath, []byte(filled), 0o644); err != nil {
		return nil, &model.ErrAssembly{
			Phase:   "merge",
			File:    relPath,
			Message: fmt.Sprintf("write output file: %v", err),
		}
	}

	return result, nil
}

// CopyPersonaFile copies a persona-only file (no core template) to the output
// directory. The relPath is relative to personaDir.
func (m *Merger) CopyPersonaFile(relPath string) (*MergeResult, error) {
	srcPath := filepath.Join(m.personaDir, relPath)
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return nil, &model.ErrAssembly{
			Phase:   "copy",
			File:    relPath,
			Message: fmt.Sprintf("read persona file: %v", err),
		}
	}

	dstPath := filepath.Join(m.outputDir, relPath)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return nil, &model.ErrAssembly{
			Phase:   "copy",
			File:    relPath,
			Message: fmt.Sprintf("create output dir: %v", err),
		}
	}

	if err := os.WriteFile(dstPath, data, 0o644); err != nil {
		return nil, &model.ErrAssembly{
			Phase:   "copy",
			File:    relPath,
			Message: fmt.Sprintf("write output file: %v", err),
		}
	}

	return &MergeResult{FilesWritten: 1}, nil
}

// PatchAgent applies persona patches to a core agent file that has already been
// copied to the output directory. Patches can append/remove skills in frontmatter
// and append content sections to the body.
func (m *Merger) PatchAgent(relPath string) error {
	patch, ok := m.manifest.AgentPatches[relPath]
	if !ok {
		return nil // No patch defined for this agent.
	}

	agentPath := filepath.Join(m.outputDir, relPath)
	data, err := os.ReadFile(agentPath)
	if err != nil {
		return &model.ErrAssembly{
			Phase:   "patch_agent",
			File:    relPath,
			Message: fmt.Sprintf("read agent file: %v", err),
		}
	}

	content := string(data)
	doc, err := parser.ParseDocumentFromString(content, relPath)
	if err != nil {
		return &model.ErrAssembly{
			Phase:   "patch_agent",
			File:    relPath,
			Message: fmt.Sprintf("parse agent: %v", err),
		}
	}

	if doc.Frontmatter == nil {
		return &model.ErrAssembly{
			Phase:   "patch_agent",
			File:    relPath,
			Message: "agent has no frontmatter to patch",
		}
	}

	// Append skills.
	if len(patch.AppendSkills) > 0 {
		existing := make(map[string]bool, len(doc.Frontmatter.Skills))
		for _, s := range doc.Frontmatter.Skills {
			existing[s] = true
		}
		for _, s := range patch.AppendSkills {
			if !existing[s] {
				doc.Frontmatter.Skills = append(doc.Frontmatter.Skills, s)
			}
		}
	}

	// Remove skills.
	if len(patch.RemoveSkills) > 0 {
		removeSet := make(map[string]bool, len(patch.RemoveSkills))
		for _, s := range patch.RemoveSkills {
			removeSet[s] = true
		}
		filtered := doc.Frontmatter.Skills[:0]
		for _, s := range doc.Frontmatter.Skills {
			if !removeSet[s] {
				filtered = append(filtered, s)
			}
		}
		doc.Frontmatter.Skills = filtered
	}

	// Serialize frontmatter back.
	fmStr, err := parser.SerializeFrontmatter(doc.Frontmatter)
	if err != nil {
		return &model.ErrAssembly{
			Phase:   "patch_agent",
			File:    relPath,
			Message: fmt.Sprintf("serialize frontmatter: %v", err),
		}
	}

	// Extract body (everything after frontmatter).
	_, body, hasFM := parser.SplitFrontmatter(content)
	if !hasFM {
		body = content
	}

	// Append content from persona file if specified.
	if patch.AppendContent != "" {
		appendPath := filepath.Join(m.personaDir, patch.AppendContent)
		appendData, err := os.ReadFile(appendPath)
		if err != nil {
			return &model.ErrAssembly{
				Phase:   "patch_agent",
				File:    relPath,
				Message: fmt.Sprintf("read append content %s: %v", patch.AppendContent, err),
			}
		}
		appendStr := string(appendData)
		// Ensure there is a newline separator before appended content.
		if !strings.HasSuffix(body, "\n") {
			body += "\n"
		}
		body += appendStr
	}

	// Reconstruct the full file.
	var output string
	if hasFM {
		output = fmStr + body
	} else {
		output = body
	}

	if err := os.WriteFile(agentPath, []byte(output), 0o644); err != nil {
		return &model.ErrAssembly{
			Phase:   "patch_agent",
			File:    relPath,
			Message: fmt.Sprintf("write patched agent: %v", err),
		}
	}

	return nil
}
