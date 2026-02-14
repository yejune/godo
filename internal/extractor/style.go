package extractor

import (
	"github.com/do-focus/convert/internal/model"
)

// StyleExtractor handles style definition files (.claude/styles/*.md).
// ALL styles are persona-specific -- each persona defines its own
// communication style. There is no core style content.
type StyleExtractor struct{}

// NewStyleExtractor creates a StyleExtractor.
func NewStyleExtractor() *StyleExtractor {
	return &StyleExtractor{}
}

// Extract classifies a style Document as fully persona-specific.
//
// Returns (nil, manifest, nil) where manifest.Styles contains the file path.
// Style files are never core -- every style belongs to a persona.
func (e *StyleExtractor) Extract(doc *model.Document) (*model.Document, *model.PersonaManifest, error) {
	manifest := &model.PersonaManifest{
		SlotContent: make(map[string]string),
	}

	manifest.Styles = append(manifest.Styles, doc.Path)
	return nil, manifest, nil
}
