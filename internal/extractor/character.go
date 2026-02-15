package extractor

import (
	"github.com/do-focus/convert/internal/model"
)

// CharacterExtractor handles character definition files (characters/*.md).
// ALL characters are persona-specific -- each persona defines its own
// character identities. There is no core character content.
type CharacterExtractor struct{}

// NewCharacterExtractor creates a CharacterExtractor.
func NewCharacterExtractor() *CharacterExtractor {
	return &CharacterExtractor{}
}

// Extract classifies a character Document as fully persona-specific.
//
// Returns (nil, manifest, nil) where manifest.Characters contains the file path.
// Character files are never core -- every character belongs to a persona.
func (e *CharacterExtractor) Extract(doc *model.Document) (*model.Document, *model.PersonaManifest, error) {
	manifest := &model.PersonaManifest{
		SlotContent: make(map[string]string),
	}

	manifest.Characters = append(manifest.Characters, doc.Path)
	return nil, manifest, nil
}
