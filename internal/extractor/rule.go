package extractor

import (
	"path/filepath"

	"github.com/yejune/godo/internal/detector"
	"github.com/yejune/godo/internal/model"
)

// RuleExtractor classifies and extracts persona-specific content from rule
// files (.claude/rules/**/*.md). Classification is based on the
// PatternRegistry's WholeFileRules list -- specific filenames that are
// persona-specific regardless of their directory location.
//
// For non-whole-file rules, inline content patterns are detected and replaced
// with {{slot:SLOT_ID}} markers so persona-specific text (like quality framework
// references) can be swapped during assembly.
type RuleExtractor struct {
	registry *detector.PatternRegistry
	detector *detector.PersonaDetector
}

// NewRuleExtractor creates a RuleExtractor with the given pattern registry
// and persona detector. The detector may be nil if content pattern detection
// is not needed.
func NewRuleExtractor(reg *detector.PatternRegistry, det *detector.PersonaDetector) *RuleExtractor {
	return &RuleExtractor{
		registry: reg,
		detector: det,
	}
}

// Extract separates a parsed rule Document into core and persona parts.
//
// Classification logic:
//   - Whole-file persona rules (filename in WholeFileRules):
//     Returns (nil, manifest, nil) where manifest.Rules contains the file path.
//   - Core rules with inline content patterns:
//     Returns (doc, manifest, nil) where doc has inline text replaced with
//     {{slot:SLOT_ID}} markers and manifest.SlotContent has original text.
//   - Core rules (all others):
//     Returns (doc, emptyManifest, nil) as a passthrough.
func (e *RuleExtractor) Extract(doc *model.Document) (*model.Document, *model.PersonaManifest, error) {
	manifest := &model.PersonaManifest{
		SlotContent: make(map[string]string),
	}

	// Check if this specific file is a whole-file persona rule
	filename := filepath.Base(doc.Path)
	if e.registry.IsWholeFilePersonaRule(filename) {
		manifest.Rules = append(manifest.Rules, doc.Path)
		return nil, manifest, nil
	}

	// Check for inline content patterns in rule body
	if e.detector != nil {
		for _, sec := range doc.Sections {
			e.slotContentPatterns(sec, manifest)
		}
	}

	return doc, manifest, nil
}

// slotContentPatterns replaces inline content pattern matches in a section's
// content with {{slot:SLOT_ID}} markers and stores original text in manifest.
func (e *RuleExtractor) slotContentPatterns(sec *model.Section, manifest *model.PersonaManifest) {
	matches := e.detector.DetectContentPatterns(sec.Content)
	if len(matches) > 0 {
		// Replace matches in reverse order to preserve offsets
		content := sec.Content
		for i := len(matches) - 1; i >= 0; i-- {
			m := matches[i]
			manifest.SlotContent[m.SlotID] = m.Original
			replacement := "{{slot:" + m.SlotID + "}}"
			content = content[:m.Start] + replacement + content[m.End:]
		}
		sec.Content = content
	}

	// Recurse into children
	for _, child := range sec.Children {
		e.slotContentPatterns(child, manifest)
	}
}
