package extractor

import (
	"path/filepath"

	"github.com/do-focus/convert/internal/detector"
	"github.com/do-focus/convert/internal/model"
)

// RuleExtractor classifies and extracts persona-specific content from rule
// files (.claude/rules/**/*.md). Classification is based on the
// PatternRegistry's WholeFileRules list -- specific filenames that are
// persona-specific regardless of their directory location.
//
// Rules in .claude/rules/do/ are generally core framework rules.
// Only files explicitly listed in WholeFileRules are persona-specific.
type RuleExtractor struct {
	registry *detector.PatternRegistry
}

// NewRuleExtractor creates a RuleExtractor with the given pattern registry.
func NewRuleExtractor(reg *detector.PatternRegistry) *RuleExtractor {
	return &RuleExtractor{
		registry: reg,
	}
}

// Extract separates a parsed rule Document into core and persona parts.
//
// Classification logic:
//   - Whole-file persona rules (filename in WholeFileRules):
//     Returns (nil, manifest, nil) where manifest.Rules contains the file path.
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

	// Core rule -- passthrough
	return doc, manifest, nil
}
