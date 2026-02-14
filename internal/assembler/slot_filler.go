package assembler

import (
	"fmt"
	"os"
	"regexp"
	"sort"

	"github.com/do-focus/convert/internal/model"
	"github.com/do-focus/convert/internal/template"
)

// sectionSlotRe matches the full section slot block including content between markers.
// Group 1: begin slot ID, Group 2: content between markers, Group 3: end slot ID.
var sectionSlotRe = regexp.MustCompile(
	`(?s)<!-- BEGIN_SLOT:([A-Z][A-Z0-9_]*) -->\n?(.*?)\n?<!-- END_SLOT:([A-Z][A-Z0-9_]*) -->`,
)

// inlineSlotRe matches {{SLOT_ID}} inline markers.
var inlineSlotRe = regexp.MustCompile(model.InlineSlotPattern)


// SlotFiller replaces slot markers in template content with persona-specific values.
type SlotFiller struct {
	registry   *template.Registry
	manifest   *model.PersonaManifest
	personaDir string
}

// NewSlotFiller creates a SlotFiller with the given registry, manifest, and persona directory.
func NewSlotFiller(registry *template.Registry, manifest *model.PersonaManifest, personaDir string) *SlotFiller {
	return &SlotFiller{
		registry:   registry,
		manifest:   manifest,
		personaDir: personaDir,
	}
}

// FillContent replaces all slot markers in content with persona values.
//
// Section slots (<!-- BEGIN_SLOT:ID -->...<!-- END_SLOT:ID -->) have their
// content replaced with persona content looked up via the registry.
//
// Inline slots ({{ID}}) are replaced with persona values.
//
// Returns the filled content, a sorted list of resolved slot IDs, and a list
// of warnings for slots that fell back to registry defaults.
func (f *SlotFiller) FillContent(content string) (string, []string, []string) {
	resolvedSet := map[string]bool{}
	var warnings []string

	// 1. Fill section slots.
	result := sectionSlotRe.ReplaceAllStringFunc(content, func(match string) string {
		sub := sectionSlotRe.FindStringSubmatch(match)
		if len(sub) < 4 {
			return match
		}
		slotID := sub[1]
		// Verify begin/end IDs match.
		if sub[1] != sub[3] {
			warnings = append(warnings, fmt.Sprintf("slot %q: mismatched BEGIN/END markers", slotID))
			return match
		}

		resolved, err := f.registry.ResolveSlot(slotID, f.manifest, f.personaDir)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("slot %q: %v", slotID, err))
			return match
		}

		resolvedSet[slotID] = true

		// Check if this fell back to registry default (no persona content).
		if f.manifest == nil || f.manifest.SlotContent == nil {
			warnings = append(warnings, fmt.Sprintf("slot %q: using registry default", slotID))
		} else if _, hasPersona := f.manifest.SlotContent[slotID]; !hasPersona {
			warnings = append(warnings, fmt.Sprintf("slot %q: using registry default", slotID))
		}

		begin := fmt.Sprintf(model.SectionSlotBegin, slotID)
		end := fmt.Sprintf(model.SectionSlotEnd, slotID)
		return begin + "\n" + resolved + "\n" + end
	})

	// 2. Fill inline slots.
	result = inlineSlotRe.ReplaceAllStringFunc(result, func(match string) string {
		sub := inlineSlotRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		slotID := sub[1]

		resolved, err := f.registry.ResolveSlot(slotID, f.manifest, f.personaDir)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("slot %q: %v", slotID, err))
			return match
		}

		resolvedSet[slotID] = true

		// Check if this fell back to registry default.
		if f.manifest == nil || f.manifest.SlotContent == nil {
			warnings = append(warnings, fmt.Sprintf("slot %q: using registry default", slotID))
		} else if _, hasPersona := f.manifest.SlotContent[slotID]; !hasPersona {
			warnings = append(warnings, fmt.Sprintf("slot %q: using registry default", slotID))
		}

		return resolved
	})

	resolved := make([]string, 0, len(resolvedSet))
	for id := range resolvedSet {
		resolved = append(resolved, id)
	}
	sort.Strings(resolved)

	return result, resolved, warnings
}

// FillFile reads a file, fills all slot markers, and writes the result back.
// Returns the number of resolved slots, number of warnings, and any I/O error.
func (f *SlotFiller) FillFile(path string) (int, int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, fmt.Errorf("read file %s: %w", path, err)
	}

	content := string(data)
	filled, resolved, warnings := f.FillContent(content)

	if len(resolved) == 0 {
		// No slots found; nothing to write.
		return 0, len(warnings), nil
	}

	if err := os.WriteFile(path, []byte(filled), 0o644); err != nil {
		return 0, 0, fmt.Errorf("write file %s: %w", path, err)
	}

	return len(resolved), len(warnings), nil
}
