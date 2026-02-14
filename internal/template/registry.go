package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/do-focus/convert/internal/model"
	"gopkg.in/yaml.v3"
)

const registryFilename = "registry.yaml"

// Registry manages the template slot definitions and their mappings.
// Persisted as registry.yaml in the core output directory.
type Registry struct {
	Version       string                `yaml:"version"`
	Source        string                `yaml:"source"`
	SourceVersion string                `yaml:"source_version"`
	ExtractedAt   string                `yaml:"extracted_at"`
	Slots         map[string]*SlotEntry `yaml:"slots"`
}

// SlotEntry is a single slot definition in the registry.
type SlotEntry struct {
	Category    string         `yaml:"category"`
	Scope       string         `yaml:"scope"`
	Description string         `yaml:"description"`
	MarkerType  string         `yaml:"marker_type"`
	FoundIn     []SlotLocation `yaml:"found_in"`
	Default     string         `yaml:"default"`
}

// SlotLocation records where a slot was found during extraction.
type SlotLocation struct {
	Path           string `yaml:"path"`
	Line           int    `yaml:"line,omitempty"`
	OriginalHeader string `yaml:"original_header,omitempty"`
	Occurrences    int    `yaml:"occurrences,omitempty"`
}

// NewRegistry creates a new empty Registry with initialized Slots map.
func NewRegistry() *Registry {
	return &Registry{
		Version: "1.0.0",
		Slots:   make(map[string]*SlotEntry),
	}
}

// LoadRegistry reads registry.yaml from a directory.
func LoadRegistry(dir string) (*Registry, error) {
	path := filepath.Join(dir, registryFilename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load registry %s: %w", path, err)
	}

	var reg Registry
	if err := yaml.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("parse registry %s: %w", path, err)
	}
	if reg.Slots == nil {
		reg.Slots = make(map[string]*SlotEntry)
	}
	return &reg, nil
}

// Save writes registry.yaml to a directory.
func (r *Registry) Save(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create registry dir %s: %w", dir, err)
	}

	data, err := yaml.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshal registry: %w", err)
	}

	path := filepath.Join(dir, registryFilename)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write registry %s: %w", path, err)
	}
	return nil
}

// ResolveSlot returns content for a slot given a persona manifest.
// Checks manifest.SlotContent[slotID] first: if the value looks like a
// persona content file (ends with .md and does not start with "."), it
// reads the file from personaDir. Otherwise the value is used as literal
// content. Falls back to the slot's Default if persona does not provide content.
func (r *Registry) ResolveSlot(slotID string, manifest *model.PersonaManifest, personaDir string) (string, error) {
	entry, ok := r.Slots[slotID]
	if !ok {
		return "", fmt.Errorf("slot %q not found in registry", slotID)
	}

	// Check persona manifest for slot content.
	if manifest != nil && manifest.SlotContent != nil {
		if val, ok := manifest.SlotContent[slotID]; ok {
			// If value looks like a file reference, read the file.
			if isFileRef(val) {
				path := filepath.Join(personaDir, val)
				data, err := os.ReadFile(path)
				if err != nil {
					return "", fmt.Errorf("read slot content file %s for slot %q: %w", path, slotID, err)
				}
				return strings.TrimRight(string(data), "\n"), nil
			}
			// Otherwise use the value as literal content.
			return val, nil
		}
	}

	// Fall back to registry default.
	return entry.Default, nil
}

// AddSlot adds or updates a slot entry in the registry.
func (r *Registry) AddSlot(slotID string, entry *SlotEntry) {
	if r.Slots == nil {
		r.Slots = make(map[string]*SlotEntry)
	}
	r.Slots[slotID] = entry
}

// isFileRef returns true if the value looks like a persona content file
// reference rather than literal content. Content files are relative paths
// within the persona directory (e.g., "content/quality.md") and always have
// a .md extension without starting with "." (which would indicate a dotpath
// like ".moai/specs/..." used as literal replacement text).
func isFileRef(val string) bool {
	if !strings.HasSuffix(val, ".md") {
		return false
	}
	// Dotpaths like ".moai/..." are literal values, not file references.
	if strings.HasPrefix(val, ".") {
		return false
	}
	return true
}
