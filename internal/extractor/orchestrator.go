package extractor

import (
	"fmt"
	"os"
	"regexp"
	"path/filepath"
	"strings"

	"github.com/yejune/godo/internal/detector"
	"github.com/yejune/godo/internal/model"
	"github.com/yejune/godo/internal/parser"
	"github.com/yejune/godo/internal/template"
)

// skipDirs contains directory names to skip during directory walking.
var skipDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	".DS_Store":    true,
	"__pycache__":  true,
}

// fileType classifies a file path into a known extractor category.
type fileType int

const (
	fileTypeUnknown fileType = iota
	fileTypeAgent
	fileTypeSkill
	fileTypeRule
	fileTypeStyle
	fileTypeClaudeMD
	fileTypeSettings
	fileTypeCommand
	fileTypeHook
	fileTypeCharacter
	fileTypeSpinner
	fileTypeAsset // non-markdown files in known directories (skills, etc.)
)

// ExtractorOrchestrator walks a .claude/ source directory, routes each file
// to the appropriate sub-extractor, and aggregates all core templates into
// a TemplateRegistry and all persona content into a merged PersonaManifest.
type ExtractorOrchestrator struct {
	registry  *detector.PatternRegistry
	agent     *AgentExtractor
	skill     *SkillExtractor
	rule      *RuleExtractor
	style     *StyleExtractor
	character *CharacterExtractor
	claudeM   *ClaudeMDExtractor
}

// NewExtractorOrchestrator creates an orchestrator wired with all sub-extractors.
func NewExtractorOrchestrator(det *detector.PersonaDetector, reg *detector.PatternRegistry) *ExtractorOrchestrator {
	return &ExtractorOrchestrator{
		registry:  reg,
		agent:     NewAgentExtractor(det, reg),
		skill:     NewSkillExtractor(reg),
		rule:      NewRuleExtractor(reg, det),
		style:     NewStyleExtractor(),
		character: NewCharacterExtractor(),
		claudeM:   NewClaudeMDExtractor(),
	}
}

// Extract walks srcDir recursively, parses each relevant file, routes it to
// the correct sub-extractor, and returns:
//   - A TemplateRegistry containing all slot entries discovered during extraction.
//   - A merged PersonaManifest combining persona content from all files.
//   - An error if a critical failure occurs (parse errors are critical).
func (o *ExtractorOrchestrator) Extract(srcDir string) (*template.Registry, *model.PersonaManifest, error) {
	registry := template.NewRegistry()
	merged := &model.PersonaManifest{
		SlotContent:  make(map[string]string),
		AgentPatches: make(map[string]*model.AgentPatch),
		PersonaFiles: make(map[string]string),
		SourceDir:    srcDir,
	}

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Skip irrelevant directories
		if info.IsDir() {
			if skipDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		// Determine relative path from srcDir
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return fmt.Errorf("relative path for %s: %w", path, err)
		}

		ft := classifyFile(relPath)
		if ft == fileTypeUnknown {
			return nil // skip non-relevant files
		}

		// Settings, commands, hooks, and spinners use different extraction APIs (not Document-based)
		switch ft {
		case fileTypeSettings:
			// Settings are persona files
			merged.PersonaFiles[relPath] = path
			return o.extractSettings(path, merged)
		case fileTypeCommand:
			// Commands are tracked as persona paths; the full extraction
			// with copy is done by the assembler. Here we just record the path.
			merged.Commands = append(merged.Commands, relPath)
			merged.PersonaFiles[relPath] = path
			return nil
		case fileTypeHook:
			merged.HookScripts = append(merged.HookScripts, relPath)
			merged.PersonaFiles[relPath] = path
			return nil
		case fileTypeSpinner:
			// Spinners are YAML files tracked as persona paths.
			merged.Spinners = append(merged.Spinners, relPath)
			merged.PersonaFiles[relPath] = path
			return nil
		case fileTypeAsset:
			// Non-markdown files in known directories (e.g., .yml, .json, .py, .sh in skills/).
			// Classify as core or persona based on parent skill directory and module patterns.
			skillName := extractSkillDirName(relPath)
			if skillName != "" && o.registry.IsWholeFilePersonaSkillDir(skillName) {
				// Whole-directory persona skill: all assets are persona.
				merged.PersonaFiles[relPath] = path
			} else if skillName != "" && o.registry.IsWholeFilePersonaSkill(skillName) {
				// Whole-file persona skill: all assets are persona.
				merged.PersonaFiles[relPath] = path
			} else if skillName != "" && o.registry.IsPartialSkill(skillName) {
				// Partial skill: check module-level classification for assets.
				moduleRelPath := extractModuleRelPath(relPath, skillName)
				if moduleRelPath != "" && o.registry.IsPartialPersonaModule(skillName, moduleRelPath) {
					merged.PersonaFiles[relPath] = path
				} else {
					merged.CoreFiles = append(merged.CoreFiles, relPath)
				}
			} else {
				merged.CoreFiles = append(merged.CoreFiles, relPath)
			}
			return nil
		}

		// Parse the markdown document
		doc, err := parser.ParseDocument(path)
		if err != nil {
			return fmt.Errorf("parse %s: %w", relPath, err)
		}
		doc.Path = relPath // use relative path for portability

		// Route to the appropriate extractor
		coreDoc, manifest, err := o.route(ft, doc)
		if err != nil {
			return fmt.Errorf("extract %s: %w", relPath, err)
		}

		// Track file: core files go to CoreFiles, persona files go to PersonaFiles
		if coreDoc != nil {
			merged.CoreFiles = append(merged.CoreFiles, relPath)
		}
		if coreDoc == nil {
			// Whole-file persona: no core doc returned
			merged.PersonaFiles[relPath] = path
		}

		// Merge manifest into the combined result
		mergeManifest(merged, manifest)

		// Register slots discovered from the core document
		if coreDoc != nil {
			registerSlots(registry, coreDoc)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	// Check for CLAUDE.md at project root (parent of srcDir).
	// In real moai-adk, CLAUDE.md lives at the project root, not inside .claude/.
	if merged.ClaudeMD == "" {
		projectRoot := filepath.Dir(srcDir)
		claudePath := filepath.Join(projectRoot, "CLAUDE.md")
		if info, statErr := os.Stat(claudePath); statErr == nil && !info.IsDir() {
			doc, parseErr := parser.ParseDocument(claudePath)
			if parseErr != nil {
				return nil, nil, fmt.Errorf("parse project-root CLAUDE.md: %w", parseErr)
			}
			doc.Path = "CLAUDE.md" // canonical relative path

			_, manifest, extErr := o.route(fileTypeClaudeMD, doc)
			if extErr != nil {
				return nil, nil, fmt.Errorf("extract project-root CLAUDE.md: %w", extErr)
			}

			// Track project-root CLAUDE.md as a persona file
			merged.PersonaFiles["CLAUDE.md"] = claudePath

			mergeManifest(merged, manifest)
		}
	}

	// Auto-detect persona name from directory structure if not already set.
	if merged.Name == "" {
		merged.Name = detectPersonaName(merged.PersonaFiles)
	}

	return registry, merged, nil
}

// route dispatches a parsed Document to the correct sub-extractor based on file type.
func (o *ExtractorOrchestrator) route(ft fileType, doc *model.Document) (*model.Document, *model.PersonaManifest, error) {
	switch ft {
	case fileTypeAgent:
		return o.agent.Extract(doc)
	case fileTypeSkill:
		return o.skill.Extract(doc)
	case fileTypeRule:
		return o.rule.Extract(doc)
	case fileTypeStyle:
		return o.style.Extract(doc)
	case fileTypeCharacter:
		return o.character.Extract(doc)
	case fileTypeClaudeMD:
		return o.claudeM.Extract(doc)
	default:
		return doc, &model.PersonaManifest{SlotContent: make(map[string]string)}, nil
	}
}

// extractSettings reads and splits a settings.json file into core and persona parts,
// storing persona settings in the manifest.
func (o *ExtractorOrchestrator) extractSettings(path string, merged *model.PersonaManifest) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read settings %s: %w", path, err)
	}

	_, persona, err := ExtractSettings(data)
	if err != nil {
		return fmt.Errorf("extract settings %s: %w", path, err)
	}

	if merged.Settings == nil {
		merged.Settings = make(map[string]interface{})
	}
	for k, v := range persona {
		merged.Settings[k] = v
	}

	return nil
}

// classifyFile determines the extractor category for a relative file path.
func classifyFile(relPath string) fileType {
	// Normalize separators
	normalized := filepath.ToSlash(relPath)

	// CLAUDE.md at root
	base := filepath.Base(normalized)
	if strings.EqualFold(base, "CLAUDE.md") && !strings.Contains(normalized, "/") {
		return fileTypeClaudeMD
	}

	// settings.json at root
	if base == "settings.json" && !strings.Contains(normalized, "/") {
		return fileTypeSettings
	}

	// Directory-based classification
	parts := strings.SplitN(normalized, "/", 2)
	if len(parts) == 0 {
		return fileTypeUnknown
	}

	topDir := parts[0]
	switch topDir {
	case "agents":
		if strings.HasSuffix(normalized, ".md") {
			return fileTypeAgent
		}
	case "skills":
		if strings.HasSuffix(normalized, ".md") {
			return fileTypeSkill
		}
		return fileTypeAsset

	case "rules":
		if strings.HasSuffix(normalized, ".md") {
			return fileTypeRule
		}
	case "styles", "output-styles":
		if strings.HasSuffix(normalized, ".md") {
			return fileTypeStyle
		}
	case "characters":
		if strings.HasSuffix(normalized, ".md") {
			return fileTypeCharacter
		}
	case "spinners":
		if strings.HasSuffix(normalized, ".yaml") || strings.HasSuffix(normalized, ".yml") {
			return fileTypeSpinner
		}
	case "commands":
		return fileTypeCommand
	case "hooks":
		return fileTypeHook
	}

	return fileTypeUnknown
}

// mergeManifest merges a single extraction manifest into the combined result.
func mergeManifest(dst, src *model.PersonaManifest) {
	if src == nil {
		return
	}

	if src.ClaudeMD != "" {
		dst.ClaudeMD = src.ClaudeMD
	}

	dst.Agents = append(dst.Agents, src.Agents...)
	dst.Skills = append(dst.Skills, src.Skills...)
	dst.Rules = append(dst.Rules, src.Rules...)
	dst.Styles = append(dst.Styles, src.Styles...)
	dst.Characters = append(dst.Characters, src.Characters...)
	dst.Spinners = append(dst.Spinners, src.Spinners...)
	dst.Commands = append(dst.Commands, src.Commands...)
	dst.HookScripts = append(dst.HookScripts, src.HookScripts...)

	for k, v := range src.SlotContent {
		if dst.SlotContent == nil {
			dst.SlotContent = make(map[string]string)
		}
		dst.SlotContent[k] = v
	}

	for k, v := range src.AgentPatches {
		if dst.AgentPatches == nil {
			dst.AgentPatches = make(map[string]*model.AgentPatch)
		}
		existing, ok := dst.AgentPatches[k]
		if ok {
			existing.AppendSkills = append(existing.AppendSkills, v.AppendSkills...)
			existing.RemoveSkills = append(existing.RemoveSkills, v.RemoveSkills...)
		} else {
			dst.AgentPatches[k] = v
		}
	}

	for k, v := range src.Settings {
		if dst.Settings == nil {
			dst.Settings = make(map[string]interface{})
		}
		dst.Settings[k] = v
	}

	// Merge SourceDir (prefer non-empty)
	if src.SourceDir != "" && dst.SourceDir == "" {
		dst.SourceDir = src.SourceDir
	}

	// Merge CoreFiles
	dst.CoreFiles = append(dst.CoreFiles, src.CoreFiles...)

	// Merge PersonaFiles
	for k, v := range src.PersonaFiles {
		if dst.PersonaFiles == nil {
			dst.PersonaFiles = make(map[string]string)
		}
		dst.PersonaFiles[k] = v
	}

	// Merge SkillMappings
	for k, v := range src.SkillMappings {
		if dst.SkillMappings == nil {
			dst.SkillMappings = make(map[string]string)
		}
		dst.SkillMappings[k] = v
	}
}

// registerSlots scans a core document's sections for slot markers and adds
// corresponding entries to the registry.
func registerSlots(reg *template.Registry, doc *model.Document) {
	walkAndRegister(reg, doc.Sections, doc.Path)
}

// walkAndRegister recursively checks sections for slot markers and registers them.
func walkAndRegister(reg *template.Registry, sections []*model.Section, docPath string) {
	for _, sec := range sections {
		if strings.Contains(sec.Content, "<!-- BEGIN_SLOT:") {
			slotID := extractSlotID(sec.Content)
			if slotID != "" {
				reg.AddSlot(slotID, &template.SlotEntry{
					Category:    "section",
					Scope:       "agent",
					Description: fmt.Sprintf("Slot from section '%s'", sec.Title),
					MarkerType:  "section",
					FoundIn: []template.SlotLocation{
						{
							Path:           docPath,
							Line:           sec.StartLine,
							OriginalHeader: sec.Title,
						},
					},
				})
			}
		}
		// Register inline {{slot:SLOT_ID}} markers from content pattern extraction
		registerInlineSlots(reg, sec, docPath)

		if len(sec.Children) > 0 {
			walkAndRegister(reg, sec.Children, docPath)
		}
	}
}

// inlineSlotPattern matches {{slot:SLOT_ID}} in section content.
var inlineSlotExtractRe = regexp.MustCompile(`\{\{slot:([A-Z][A-Z0-9_]*)\}\}`)

// registerInlineSlots finds {{slot:SLOT_ID}} inline markers in section content
// and adds them to the template registry.
func registerInlineSlots(reg *template.Registry, sec *model.Section, docPath string) {
	matches := inlineSlotExtractRe.FindAllStringSubmatch(sec.Content, -1)
	for _, m := range matches {
		slotID := m[1]
		// Only add if not already registered (avoid duplicates from section slots)
		if _, exists := reg.Slots[slotID]; !exists {
			reg.AddSlot(slotID, &template.SlotEntry{
				Category:    "content_pattern",
				Scope:       "rule",
				Description: fmt.Sprintf("Inline content pattern slot in '%s'", sec.Title),
				MarkerType:  "inline",
				FoundIn: []template.SlotLocation{
					{
						Path: docPath,
						Line: sec.StartLine,
					},
				},
			})
		}
	}
}

// extractSlotID parses "<!-- BEGIN_SLOT:SLOT_ID -->" to get the slot ID.
func extractSlotID(content string) string {
	const prefix = "<!-- BEGIN_SLOT:"
	const suffix = " -->"

	idx := strings.Index(content, prefix)
	if idx < 0 {
		return ""
	}

	start := idx + len(prefix)
	rest := content[start:]
	endIdx := strings.Index(rest, suffix)
	if endIdx < 0 {
		return ""
	}

	return strings.TrimSpace(rest[:endIdx])
}

// personaDirs lists the top-level directories that may contain persona subdirectories.
var personaDirs = map[string]bool{
	"agents":       true,
	"hooks":        true,
	"commands":     true,
	"skills":       true,
	"rules":        true,
	"styles":       true,
	"output-styles": true,
}

// detectPersonaName scans persona file relative paths to find the most common
// first-level subdirectory name under known top-level directories (agents/, hooks/,
// commands/, etc.). For example, given paths like "agents/moai/manager-spec.md" and
// "hooks/moai/pre-tool.sh", the common subdirectory "moai" is returned.
// Returns empty string if no common subdirectory is found.
func detectPersonaName(personaFiles map[string]string) string {
	// Count how many top-level dirs each subdirectory name appears under.
	// e.g., "moai" -> {"agents": true, "hooks": true, "commands": true}
	subDirTopDirs := make(map[string]map[string]bool)

	for relPath := range personaFiles {
		normalized := filepath.ToSlash(relPath)
		parts := strings.Split(normalized, "/")
		// Need at least 3 parts: topDir/subDir/file (e.g., agents/moai/spec.md)
		if len(parts) < 3 {
			continue
		}

		topDir := parts[0]
		if !personaDirs[topDir] {
			continue
		}

		subDir := parts[1]
		if subDirTopDirs[subDir] == nil {
			subDirTopDirs[subDir] = make(map[string]bool)
		}
		subDirTopDirs[subDir][topDir] = true
	}

	if len(subDirTopDirs) == 0 {
		return ""
	}

	// Find the subdirectory name that appears under the most top-level dirs.
	var bestName string
	bestCount := 0
	for name, topDirs := range subDirTopDirs {
		if len(topDirs) > bestCount {
			bestCount = len(topDirs)
			bestName = name
		}
	}

	// Require at least 2 top-level dirs to be confident this is a persona name.
	if bestCount < 2 {
		return ""
	}

	return bestName
}

// extractSkillDirName extracts the skill directory name from a relative path.
// For example, "skills/moai-tool-ast-grep/rules/go.yml" returns "moai-tool-ast-grep".
// Returns empty string if the path doesn't have a skill subdirectory component.
func extractSkillDirName(relPath string) string {
	normalized := filepath.ToSlash(relPath)
	parts := strings.Split(normalized, "/")
	// Need at least: skills/<name>/<file> (3 parts)
	if len(parts) < 3 {
		return ""
	}
	if parts[0] != "skills" {
		return ""
	}
	return parts[1]
}
