package assembler

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/do-focus/convert/internal/model"
	"github.com/do-focus/convert/internal/template"
)

// newTestRegistry creates a Registry with the given slot entries.
func newTestRegistry(slots map[string]*template.SlotEntry) *template.Registry {
	r := template.NewRegistry()
	for id, entry := range slots {
		r.AddSlot(id, entry)
	}
	return r
}

func TestFillContent_SectionSlot(t *testing.T) {
	reg := newTestRegistry(map[string]*template.SlotEntry{
		"QUALITY_FRAMEWORK": {
			Category:   "section",
			MarkerType: "section",
			Default:    "default quality content",
		},
	})

	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{
			"QUALITY_FRAMEWORK": "TRUST 5 quality gates enforced",
		},
	}

	filler := NewSlotFiller(reg, manifest, "")

	input := "# Rules\n<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->\nold content here\n<!-- END_SLOT:QUALITY_FRAMEWORK -->\n# End"
	filled, resolved, warnings := filler.FillContent(input)

	if len(resolved) != 1 || resolved[0] != "QUALITY_FRAMEWORK" {
		t.Errorf("expected resolved=[QUALITY_FRAMEWORK], got %v", resolved)
	}
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}

	expected := "# Rules\n<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->\nTRUST 5 quality gates enforced\n<!-- END_SLOT:QUALITY_FRAMEWORK -->\n# End"
	if filled != expected {
		t.Errorf("filled content mismatch.\nexpected:\n%s\ngot:\n%s", expected, filled)
	}
}

func TestFillContent_InlineSlot(t *testing.T) {
	reg := newTestRegistry(map[string]*template.SlotEntry{
		"SPEC_PATH": {
			Category:   "path_pattern",
			MarkerType: "inline",
			Default:    ".moai/specs/",
		},
	})

	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{
			"SPEC_PATH": ".do/specs/",
		},
	}

	filler := NewSlotFiller(reg, manifest, "")

	input := "Read specs from {{SPEC_PATH}} directory."
	filled, resolved, warnings := filler.FillContent(input)

	if len(resolved) != 1 || resolved[0] != "SPEC_PATH" {
		t.Errorf("expected resolved=[SPEC_PATH], got %v", resolved)
	}
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}

	expected := "Read specs from .do/specs/ directory."
	if filled != expected {
		t.Errorf("filled content mismatch.\nexpected: %s\ngot: %s", expected, filled)
	}
}

func TestFillContent_UndefinedSlotPreservedAsIs(t *testing.T) {
	reg := newTestRegistry(map[string]*template.SlotEntry{
		"QUALITY_FRAMEWORK": {
			Category:   "section",
			MarkerType: "section",
			Default:    "default quality content",
		},
		"SPEC_PATH": {
			Category:   "path_pattern",
			MarkerType: "inline",
			Default:    ".moai/specs/",
		},
	})

	// Manifest with no slot content — undefined slots must be preserved.
	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{},
	}

	filler := NewSlotFiller(reg, manifest, "")

	// Section slot: not in slot_content → falls back to registry default.
	sectionInput := "<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->\noriginal\n<!-- END_SLOT:QUALITY_FRAMEWORK -->"
	filled, resolved, warnings := filler.FillContent(sectionInput)

	if len(resolved) != 1 {
		t.Errorf("section: expected 1 resolved slot, got %d: %v", len(resolved), resolved)
	}
	if len(warnings) != 1 {
		t.Errorf("section: expected 1 warning about registry default, got %d: %v", len(warnings), warnings)
	}
	expectedSection := "<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->\ndefault quality content\n<!-- END_SLOT:QUALITY_FRAMEWORK -->"
	if filled != expectedSection {
		t.Errorf("section: expected registry default content.\nexpected:\n%s\ngot:\n%s", expectedSection, filled)
	}

	// Inline slot: not in slot_content, must be preserved as-is.
	inlineInput := "Read from {{SPEC_PATH}} now."
	filled2, resolved2, warnings2 := filler.FillContent(inlineInput)

	if len(resolved2) != 0 {
		t.Errorf("inline: expected no resolved slots, got %v", resolved2)
	}
	if len(warnings2) != 0 {
		t.Errorf("inline: expected no warnings, got %v", warnings2)
	}
	if filled2 != inlineInput {
		t.Errorf("inline: expected content preserved as-is.\nexpected: %s\ngot: %s", inlineInput, filled2)
	}
}

func TestFillContent_MixedDefinedAndUndefinedSlots(t *testing.T) {
	reg := newTestRegistry(map[string]*template.SlotEntry{
		"TOOL_NAME": {
			Category:   "path_pattern",
			MarkerType: "inline",
			Default:    "moai",
		},
		"PRIMARY_USERS": {
			Category:   "path_pattern",
			MarkerType: "inline",
			Default:    "",
		},
	})

	// Only TOOL_NAME defined; PRIMARY_USERS is NOT in slot_content.
	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{
			"TOOL_NAME": "godo",
		},
	}

	filler := NewSlotFiller(reg, manifest, "")

	input := "Use {{TOOL_NAME}} for {{PRIMARY_USERS}} users."
	filled, resolved, warnings := filler.FillContent(input)

	if len(resolved) != 1 || resolved[0] != "TOOL_NAME" {
		t.Errorf("expected resolved=[TOOL_NAME], got %v", resolved)
	}
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}

	expected := "Use godo for {{PRIMARY_USERS}} users."
	if filled != expected {
		t.Errorf("mixed fill mismatch.\nexpected: %s\ngot: %s", expected, filled)
	}
}

func TestFillContent_MixedSlots(t *testing.T) {
	reg := newTestRegistry(map[string]*template.SlotEntry{
		"QUALITY_FRAMEWORK": {
			Category:   "section",
			MarkerType: "section",
			Default:    "default quality",
		},
		"SPEC_PATH": {
			Category:   "path_pattern",
			MarkerType: "inline",
			Default:    ".moai/specs/",
		},
		"TOOL_NAME": {
			Category:   "path_pattern",
			MarkerType: "inline",
			Default:    "moai",
		},
	})

	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{
			"QUALITY_FRAMEWORK": "TRUST 5 gates",
			"SPEC_PATH":         ".do/specs/",
			"TOOL_NAME":         "godo",
		},
	}

	filler := NewSlotFiller(reg, manifest, "")

	input := `# Config
Use {{TOOL_NAME}} to manage specs at {{SPEC_PATH}} path.

<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->
old quality content
<!-- END_SLOT:QUALITY_FRAMEWORK -->

Run {{TOOL_NAME}} now.`

	filled, resolved, warnings := filler.FillContent(input)

	if len(resolved) != 3 {
		t.Errorf("expected 3 resolved slots, got %d: %v", len(resolved), resolved)
	}
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}

	expected := `# Config
Use godo to manage specs at .do/specs/ path.

<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->
TRUST 5 gates
<!-- END_SLOT:QUALITY_FRAMEWORK -->

Run godo now.`

	if filled != expected {
		t.Errorf("mixed filled mismatch.\nexpected:\n%s\ngot:\n%s", expected, filled)
	}
}

func TestFillContent_SectionSlotFromFile(t *testing.T) {
	// Create a temp persona directory with a content file.
	personaDir := t.TempDir()
	contentDir := filepath.Join(personaDir, "content")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	contentFile := filepath.Join(contentDir, "quality.md")
	if err := os.WriteFile(contentFile, []byte("File-based quality content"), 0o644); err != nil {
		t.Fatal(err)
	}

	reg := newTestRegistry(map[string]*template.SlotEntry{
		"QUALITY_FRAMEWORK": {
			Category:   "section",
			MarkerType: "section",
			Default:    "default quality",
		},
	})

	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{
			"QUALITY_FRAMEWORK": "content/quality.md",
		},
	}

	filler := NewSlotFiller(reg, manifest, personaDir)

	input := "<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->\nplaceholder\n<!-- END_SLOT:QUALITY_FRAMEWORK -->"
	filled, resolved, warnings := filler.FillContent(input)

	if len(resolved) != 1 || resolved[0] != "QUALITY_FRAMEWORK" {
		t.Errorf("expected resolved=[QUALITY_FRAMEWORK], got %v", resolved)
	}
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %v", warnings)
	}

	expected := "<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->\nFile-based quality content\n<!-- END_SLOT:QUALITY_FRAMEWORK -->"
	if filled != expected {
		t.Errorf("file-based fill mismatch.\nexpected:\n%s\ngot:\n%s", expected, filled)
	}
}

func TestFillFile(t *testing.T) {
	reg := newTestRegistry(map[string]*template.SlotEntry{
		"TOOL_NAME": {
			Category:   "path_pattern",
			MarkerType: "inline",
			Default:    "moai",
		},
	})

	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{
			"TOOL_NAME": "godo",
		},
	}

	filler := NewSlotFiller(reg, manifest, "")

	// Create temp file with slot markers.
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(tmpFile, []byte("Run {{TOOL_NAME}} now."), 0o644); err != nil {
		t.Fatal(err)
	}

	resolvedCount, warnCount, err := filler.FillFile(tmpFile)
	if err != nil {
		t.Fatalf("FillFile error: %v", err)
	}
	if resolvedCount != 1 {
		t.Errorf("expected 1 resolved, got %d", resolvedCount)
	}
	if warnCount != 0 {
		t.Errorf("expected 0 warnings, got %d", warnCount)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "Run godo now." {
		t.Errorf("file content mismatch: got %q", string(data))
	}
}
