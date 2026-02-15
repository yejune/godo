package assembler

import (
	"encoding/json"
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
//
// Frontmatter is patched at the raw YAML text level (not re-serialized) to
// preserve original key order, formatting, and value styles.
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

	// Split into raw YAML frontmatter and body.
	rawYaml, body, hasFM := parser.SplitFrontmatter(content)
	if !hasFM {
		// No frontmatter to patch; handle append-content only.
		if patch.AppendContent != "" {
			body, err = m.appendContentToBody(content, patch.AppendContent)
			if err != nil {
				return err
			}
			return os.WriteFile(agentPath, []byte(body), 0o644)
		}
		return &model.ErrAssembly{
			Phase:   "patch_agent",
			File:    relPath,
			Message: "agent has no frontmatter to patch",
		}
	}

	// Parse frontmatter to get structured skills (for dedup logic).
	fm, err := parser.ParseFrontmatter(rawYaml)
	if err != nil {
		return &model.ErrAssembly{
			Phase:   "patch_agent",
			File:    relPath,
			Message: fmt.Sprintf("parse frontmatter: %v", err),
		}
	}

	// Compute new skills list using structured data.
	skills := fm.Skills

	// Append skills (with dedup).
	if len(patch.AppendSkills) > 0 {
		existing := make(map[string]bool, len(skills))
		for _, s := range skills {
			existing[s] = true
		}
		for _, s := range patch.AppendSkills {
			if !existing[s] {
				skills = append(skills, s)
			}
		}
	}

	// Remove skills.
	if len(patch.RemoveSkills) > 0 {
		removeSet := make(map[string]bool, len(patch.RemoveSkills))
		for _, s := range patch.RemoveSkills {
			removeSet[s] = true
		}
		filtered := skills[:0]
		for _, s := range skills {
			if !removeSet[s] {
				filtered = append(filtered, s)
			}
		}
		skills = filtered
	}

	// Patch raw YAML text to update skills, preserving original format.
	patchedYaml := parser.PatchFrontmatterSkills(rawYaml, skills)

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
		if !strings.HasSuffix(body, "\n") {
			body += "\n"
		}
		body += appendStr
	}

	// Reconstruct the full file with patched frontmatter.
	output := "---\n" + patchedYaml + "---\n" + body

	if err := os.WriteFile(agentPath, []byte(output), 0o644); err != nil {
		return &model.ErrAssembly{
			Phase:   "patch_agent",
			File:    relPath,
			Message: fmt.Sprintf("write patched agent: %v", err),
		}
	}

	return nil
}

// appendContentToBody appends persona content from a file to the body text.
func (m *Merger) appendContentToBody(body, appendRelPath string) (string, error) {
	appendPath := filepath.Join(m.personaDir, appendRelPath)
	appendData, err := os.ReadFile(appendPath)
	if err != nil {
		return "", &model.ErrAssembly{
			Phase:   "patch_agent",
			File:    appendRelPath,
			Message: fmt.Sprintf("read append content: %v", err),
		}
	}
	if !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	return body + string(appendData), nil
}

// MergeSettings reads core settings.json, merges manifest Settings on top
// (with deep merge for the "env" sub-object), injects persona-specific hooks,
// and writes the merged result to the output directory.
func (m *Merger) MergeSettings(coreSettingsPath string) error {
	data, err := os.ReadFile(coreSettingsPath)
	if err != nil {
		return &model.ErrAssembly{
			Phase:   "merge_settings",
			File:    "settings.json",
			Message: fmt.Sprintf("read core settings: %v", err),
		}
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return &model.ErrAssembly{
			Phase:   "merge_settings",
			File:    "settings.json",
			Message: fmt.Sprintf("parse core settings: %v", err),
		}
	}

	// Inject persona hooks if present.
	if len(m.manifest.Hooks) > 0 {
		settings["hooks"] = m.manifest.Hooks
	}

	// Apply manifest.Settings overrides (deep merge for "env", override for top-level).
	if len(m.manifest.Settings) > 0 {
		mergeSettingsMap(settings, m.manifest.Settings)
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

	dstPath := filepath.Join(m.outputDir, "settings.json")
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

	return nil
}

// mergeSettingsMap merges persona settings into base settings.
// For the "env" key: deep merge (persona keys add/override, base keys preserved).
// For all other top-level keys: persona values override base values.
func mergeSettingsMap(base, persona map[string]interface{}) {
	for key, personaVal := range persona {
		if key == "env" {
			// Deep merge for env sub-object.
			baseEnv, baseOk := base["env"]
			if baseOk {
				baseEnvMap, baseIsMap := baseEnv.(map[string]interface{})
				personaEnvMap, personaIsMap := personaVal.(map[string]interface{})
				if baseIsMap && personaIsMap {
					for k, v := range personaEnvMap {
						baseEnvMap[k] = v
					}
					base["env"] = baseEnvMap
					continue
				}
			}
		}
		// Top-level override (including env when base has no env or types don't match).
		base[key] = personaVal
	}
}

// ApplySkillMappings walks all agent .md files in the output directory and
// replaces skill names in frontmatter according to the manifest's SkillMappings.
func (m *Merger) ApplySkillMappings() (int, error) {
	if len(m.manifest.SkillMappings) == 0 {
		return 0, nil
	}

	agentsDir := filepath.Join(m.outputDir, "agents")
	if _, err := os.Stat(agentsDir); os.IsNotExist(err) {
		return 0, nil
	}

	agentsModified := 0

	err := filepath.Walk(agentsDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		modified, err := m.applySkillMappingsToFile(path)
		if err != nil {
			relPath, _ := filepath.Rel(m.outputDir, path)
			return &model.ErrAssembly{
				Phase:   "skill_mappings",
				File:    relPath,
				Message: fmt.Sprintf("apply skill mappings: %v", err),
			}
		}
		if modified {
			agentsModified++
		}
		return nil
	})

	return agentsModified, err
}

// applySkillMappingsToFile reads a single agent file, replaces skill names
// in frontmatter per the manifest's SkillMappings, and writes back if changed.
func (m *Merger) applySkillMappingsToFile(agentPath string) (bool, error) {
	data, err := os.ReadFile(agentPath)
	if err != nil {
		return false, err
	}

	content := string(data)
	rawYaml, body, hasFM := parser.SplitFrontmatter(content)
	if !hasFM {
		return false, nil
	}

	fm, err := parser.ParseFrontmatter(rawYaml)
	if err != nil {
		return false, err
	}

	if len(fm.Skills) == 0 {
		return false, nil
	}

	// Apply mappings to the skills list.
	changed := false
	newSkills := make([]string, len(fm.Skills))
	for i, skill := range fm.Skills {
		if replacement, ok := m.manifest.SkillMappings[skill]; ok {
			newSkills[i] = replacement
			changed = true
		} else {
			newSkills[i] = skill
		}
	}

	if !changed {
		return false, nil
	}

	// Patch raw YAML to update skills, preserving format.
	patchedYaml := parser.PatchFrontmatterSkills(rawYaml, newSkills)
	output := "---\n" + patchedYaml + "---\n" + body

	if err := os.WriteFile(agentPath, []byte(output), 0o644); err != nil {
		return false, err
	}

	return true, nil
}

// CopyCommands copies persona-specific command files listed in the manifest
// from personaDir to outputDir, preserving relative paths.
func (m *Merger) CopyCommands() (*MergeResult, error) {
	result := &MergeResult{}

	for _, relPath := range m.manifest.Commands {
		srcPath := filepath.Join(m.personaDir, relPath)
		dstPath := filepath.Join(m.outputDir, relPath)

		data, err := os.ReadFile(srcPath)
		if err != nil {
			return nil, &model.ErrAssembly{
				Phase:   "copy_commands",
				File:    relPath,
				Message: fmt.Sprintf("read command file: %v", err),
			}
		}

		if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
			return nil, &model.ErrAssembly{
				Phase:   "copy_commands",
				File:    relPath,
				Message: fmt.Sprintf("create output dir: %v", err),
			}
		}

		if err := os.WriteFile(dstPath, data, 0o644); err != nil {
			return nil, &model.ErrAssembly{
				Phase:   "copy_commands",
				File:    relPath,
				Message: fmt.Sprintf("write command file: %v", err),
			}
		}

		result.FilesWritten++
	}

	return result, nil
}

// CopyHookScripts copies persona-specific hook scripts listed in the manifest
// from personaDir to outputDir, preserving relative paths and file permissions.
func (m *Merger) CopyHookScripts() (*MergeResult, error) {
	result := &MergeResult{}

	for _, relPath := range m.manifest.HookScripts {
		srcPath := filepath.Join(m.personaDir, relPath)

		info, err := os.Stat(srcPath)
		if err != nil {
			return nil, &model.ErrAssembly{
				Phase:   "copy_hook_scripts",
				File:    relPath,
				Message: fmt.Sprintf("stat hook script: %v", err),
			}
		}

		data, err := os.ReadFile(srcPath)
		if err != nil {
			return nil, &model.ErrAssembly{
				Phase:   "copy_hook_scripts",
				File:    relPath,
				Message: fmt.Sprintf("read hook script: %v", err),
			}
		}

		dstPath := filepath.Join(m.outputDir, relPath)
		if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
			return nil, &model.ErrAssembly{
				Phase:   "copy_hook_scripts",
				File:    relPath,
				Message: fmt.Sprintf("create output dir: %v", err),
			}
		}

		if err := os.WriteFile(dstPath, data, info.Mode().Perm()); err != nil {
			return nil, &model.ErrAssembly{
				Phase:   "copy_hook_scripts",
				File:    relPath,
				Message: fmt.Sprintf("write hook script: %v", err),
			}
		}

		result.FilesWritten++
	}

	return result, nil
}
