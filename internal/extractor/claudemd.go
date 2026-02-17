package extractor

import (
	"github.com/yejune/godo/internal/model"
)

// ClaudeMDExtractor handles the top-level CLAUDE.md file.
// The entire CLAUDE.md is persona-specific -- each persona gets its own
// CLAUDE.md defining identity, workflow, and configuration.
type ClaudeMDExtractor struct{}

// NewClaudeMDExtractor creates a ClaudeMDExtractor.
func NewClaudeMDExtractor() *ClaudeMDExtractor {
	return &ClaudeMDExtractor{}
}

// Extract classifies a CLAUDE.md Document as fully persona-specific.
//
// Returns (nil, manifest, nil) where manifest.ClaudeMD contains the file path.
// CLAUDE.md files are never core -- every persona defines its own.
func (e *ClaudeMDExtractor) Extract(doc *model.Document) (*model.Document, *model.PersonaManifest, error) {
	manifest := &model.PersonaManifest{
		SlotContent: make(map[string]string),
	}

	manifest.ClaudeMD = doc.Path
	return nil, manifest, nil
}
