package extractor

import (
	"path/filepath"
	"strings"

	"github.com/do-focus/convert/internal/detector"
	"github.com/do-focus/convert/internal/model"
)

// SkillExtractor classifies and extracts persona-specific content from skill
// definition files (.claude/skills/**/*.md). Skills are classified at three levels:
//   - Whole-file persona: entire skill is persona (listed in WholeFileSkills)
//   - Partial persona: specific modules within a skill are persona (PartialSkillPatterns)
//   - Core: everything else passes through unchanged
type SkillExtractor struct {
	registry *detector.PatternRegistry
}

// NewSkillExtractor creates a SkillExtractor with the given pattern registry.
func NewSkillExtractor(reg *detector.PatternRegistry) *SkillExtractor {
	return &SkillExtractor{
		registry: reg,
	}
}

// Extract separates a parsed skill Document into core and persona parts.
//
// Classification hierarchy:
//  1. Whole-file persona skills (WholeFileSkills): entire skill -> persona
//  2. Partial persona skills (PartialSkillPatterns): module-level classification
//  3. Everything else: core passthrough
func (e *SkillExtractor) Extract(doc *model.Document) (*model.Document, *model.PersonaManifest, error) {
	manifest := &model.PersonaManifest{
		SlotContent: make(map[string]string),
	}

	// Check whole-file persona first (by frontmatter name).
	if doc.Frontmatter != nil && e.registry.IsWholeFilePersonaSkill(doc.Frontmatter.Name) {
		manifest.Skills = append(manifest.Skills, doc.Path)
		return nil, manifest, nil
	}

	// Check partial skill: classify individual modules as persona or core.
	skillName := extractSkillName(doc.Path)
	if skillName != "" && e.registry.IsPartialSkill(skillName) {
		moduleRelPath := extractModuleRelPath(doc.Path, skillName)
		if moduleRelPath != "" && e.registry.IsPartialPersonaModule(skillName, moduleRelPath) {
			// This module is persona-specific.
			manifest.Skills = append(manifest.Skills, doc.Path)
			return nil, manifest, nil
		}
	}

	// Core skill -- passthrough
	return doc, manifest, nil
}

// extractSkillName extracts the skill directory name from a skill file path.
// For "skills/moai-workflow-testing/modules/ddd/core-classes.md", returns "moai-workflow-testing".
// For "skills/moai-workflow-testing/SKILL.md", returns "moai-workflow-testing".
func extractSkillName(relPath string) string {
	normalized := filepath.ToSlash(relPath)
	parts := strings.Split(normalized, "/")
	// Need at least: skills/<name>/<file>
	if len(parts) < 3 || parts[0] != "skills" {
		return ""
	}
	return parts[1]
}

// extractModuleRelPath extracts the path relative to the skill directory.
// For "skills/moai-workflow-testing/modules/ddd/core-classes.md" with
// skillName "moai-workflow-testing", returns "modules/ddd/core-classes.md".
// Returns empty string if the path doesn't have content beyond the skill dir.
func extractModuleRelPath(relPath, skillName string) string {
	normalized := filepath.ToSlash(relPath)
	prefix := "skills/" + skillName + "/"
	if !strings.HasPrefix(normalized, prefix) {
		return ""
	}
	return normalized[len(prefix):]
}
