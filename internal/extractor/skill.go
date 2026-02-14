package extractor

import (
	"github.com/do-focus/convert/internal/detector"
	"github.com/do-focus/convert/internal/model"
)

// SkillExtractor classifies and extracts persona-specific content from skill
// definition files (.claude/skills/*.md). Skills are classified as whole-file:
// either entirely core or entirely persona, based on the skill name.
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
// Skills are classified at the whole-file level using PatternRegistry:
//   - Whole-file persona skills (listed in WholeFileSkills):
//     Returns (nil, manifest, nil) where manifest.Skills contains the file path.
//   - Core skills (everything else):
//     Returns (doc, emptyManifest, nil) as a passthrough.
func (e *SkillExtractor) Extract(doc *model.Document) (*model.Document, *model.PersonaManifest, error) {
	manifest := &model.PersonaManifest{
		SlotContent: make(map[string]string),
	}

	// Skills are whole-file: check name against registry
	if doc.Frontmatter != nil && e.registry.IsWholeFilePersonaSkill(doc.Frontmatter.Name) {
		manifest.Skills = append(manifest.Skills, doc.Path)
		return nil, manifest, nil
	}

	// Core skill -- passthrough
	return doc, manifest, nil
}
