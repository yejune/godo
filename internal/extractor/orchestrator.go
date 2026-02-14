package extractor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/do-focus/convert/internal/detector"
	"github.com/do-focus/convert/internal/model"
	"github.com/do-focus/convert/internal/parser"
	"github.com/do-focus/convert/internal/template"
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
)

// ExtractorOrchestrator walks a .claude/ source directory, routes each file
// to the appropriate sub-extractor, and aggregates all core templates into
// a TemplateRegistry and all persona content into a merged PersonaManifest.
type ExtractorOrchestrator struct {
	agent   *AgentExtractor
	skill   *SkillExtractor
	rule    *RuleExtractor
	style   *StyleExtractor
	claudeM *ClaudeMDExtractor
}

// NewExtractorOrchestrator creates an orchestrator wired with all sub-extractors.
func NewExtractorOrchestrator(det *detector.PersonaDetector, reg *detector.PatternRegistry) *ExtractorOrchestrator {
	return &ExtractorOrchestrator{
		agent:   NewAgentExtractor(det, reg),
		skill:   NewSkillExtractor(reg),
		rule:    NewRuleExtractor(reg),
		style:   NewStyleExtractor(),
		claudeM: NewClaudeMDExtractor(),
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

		// Settings and commands/hooks use different extraction APIs (not Document-based)
		switch ft {
		case fileTypeSettings:
			return o.extractSettings(path, merged)
		case fileTypeCommand:
			// Commands are tracked as persona paths; the full extraction
			// with copy is done by the assembler. Here we just record the path.
			merged.Commands = append(merged.Commands, relPath)
			return nil
		case fileTypeHook:
			merged.HookScripts = append(merged.HookScripts, relPath)
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
			mergeManifest(merged, manifest)
		}
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
	case "rules":
		if strings.HasSuffix(normalized, ".md") {
			return fileTypeRule
		}
	case "styles", "output-styles":
		if strings.HasSuffix(normalized, ".md") {
			return fileTypeStyle
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
		if len(sec.Children) > 0 {
			walkAndRegister(reg, sec.Children, docPath)
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
