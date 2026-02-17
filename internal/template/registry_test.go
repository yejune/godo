package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yejune/godo/internal/model"
)

func TestRegistrySaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()

	original := &Registry{
		Version:       "1.0.0",
		Source:        "moai-adk",
		SourceVersion: "13.0.0",
		ExtractedAt:   "2026-02-14T12:00:00Z",
		Slots: map[string]*SlotEntry{
			"QUALITY_FRAMEWORK": {
				Category:    "section",
				Scope:       "agent_body",
				Description: "Quality compliance section",
				MarkerType:  "section",
				FoundIn: []SlotLocation{
					{Path: "agents/moai/expert-backend.md", Line: 640, OriginalHeader: "### TRUST 5 Compliance"},
				},
				Default: "### Quality\nFollow standards.",
			},
			"SPEC_PATH_PATTERN": {
				Category:    "path_pattern",
				Scope:       "agent_body",
				Description: "Specification directory path",
				MarkerType:  "inline",
				FoundIn: []SlotLocation{
					{Path: "agents/moai/expert-backend.md", Occurrences: 3},
				},
				Default: "the project's specification directory",
			},
		},
	}

	if err := original.Save(dir); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists.
	if _, err := os.Stat(filepath.Join(dir, "registry.yaml")); err != nil {
		t.Fatalf("registry.yaml not created: %v", err)
	}

	loaded, err := LoadRegistry(dir)
	if err != nil {
		t.Fatalf("LoadRegistry failed: %v", err)
	}

	// Verify top-level fields.
	if loaded.Version != original.Version {
		t.Errorf("Version: got %q, want %q", loaded.Version, original.Version)
	}
	if loaded.Source != original.Source {
		t.Errorf("Source: got %q, want %q", loaded.Source, original.Source)
	}
	if loaded.SourceVersion != original.SourceVersion {
		t.Errorf("SourceVersion: got %q, want %q", loaded.SourceVersion, original.SourceVersion)
	}
	if loaded.ExtractedAt != original.ExtractedAt {
		t.Errorf("ExtractedAt: got %q, want %q", loaded.ExtractedAt, original.ExtractedAt)
	}

	// Verify slots count.
	if len(loaded.Slots) != 2 {
		t.Fatalf("Slots count: got %d, want 2", len(loaded.Slots))
	}

	// Verify QUALITY_FRAMEWORK slot.
	qf := loaded.Slots["QUALITY_FRAMEWORK"]
	if qf == nil {
		t.Fatal("QUALITY_FRAMEWORK slot missing")
	}
	if qf.Category != "section" {
		t.Errorf("QUALITY_FRAMEWORK.Category: got %q, want %q", qf.Category, "section")
	}
	if qf.MarkerType != "section" {
		t.Errorf("QUALITY_FRAMEWORK.MarkerType: got %q, want %q", qf.MarkerType, "section")
	}
	if len(qf.FoundIn) != 1 {
		t.Fatalf("QUALITY_FRAMEWORK.FoundIn count: got %d, want 1", len(qf.FoundIn))
	}
	if qf.FoundIn[0].Line != 640 {
		t.Errorf("QUALITY_FRAMEWORK.FoundIn[0].Line: got %d, want 640", qf.FoundIn[0].Line)
	}
	if qf.Default != "### Quality\nFollow standards." {
		t.Errorf("QUALITY_FRAMEWORK.Default: got %q, want %q", qf.Default, "### Quality\nFollow standards.")
	}

	// Verify SPEC_PATH_PATTERN slot.
	sp := loaded.Slots["SPEC_PATH_PATTERN"]
	if sp == nil {
		t.Fatal("SPEC_PATH_PATTERN slot missing")
	}
	if sp.MarkerType != "inline" {
		t.Errorf("SPEC_PATH_PATTERN.MarkerType: got %q, want %q", sp.MarkerType, "inline")
	}
	if sp.FoundIn[0].Occurrences != 3 {
		t.Errorf("SPEC_PATH_PATTERN.FoundIn[0].Occurrences: got %d, want 3", sp.FoundIn[0].Occurrences)
	}
}

func TestResolveSlotWithPersonaContent(t *testing.T) {
	dir := t.TempDir()
	personaDir := t.TempDir()

	// Create a content file in persona dir.
	contentDir := filepath.Join(personaDir, "content")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatalf("mkdir content: %v", err)
	}
	contentFile := filepath.Join(contentDir, "quality.md")
	if err := os.WriteFile(contentFile, []byte("### TRUST 5\n- Tested\n- Readable\n"), 0o644); err != nil {
		t.Fatalf("write content file: %v", err)
	}

	reg := &Registry{
		Version: "1.0.0",
		Slots: map[string]*SlotEntry{
			"QUALITY_FRAMEWORK": {
				Category: "section",
				Default:  "### Quality\nGeneric compliance.",
			},
		},
	}
	if err := reg.Save(dir); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{
			"QUALITY_FRAMEWORK": "content/quality.md",
		},
	}

	content, err := reg.ResolveSlot("QUALITY_FRAMEWORK", manifest, personaDir)
	if err != nil {
		t.Fatalf("ResolveSlot failed: %v", err)
	}

	expected := "### TRUST 5\n- Tested\n- Readable"
	if content != expected {
		t.Errorf("ResolveSlot with persona content\ngot:  %q\nwant: %q", content, expected)
	}
}

func TestResolveSlotWithLiteralContent(t *testing.T) {
	reg := &Registry{
		Version: "1.0.0",
		Slots: map[string]*SlotEntry{
			"SPEC_PATH_PATTERN": {
				Category: "path_pattern",
				Default:  "the project's specification directory",
			},
		},
	}

	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{
			"SPEC_PATH_PATTERN": ".moai/specs/SPEC-{ID}",
		},
	}

	content, err := reg.ResolveSlot("SPEC_PATH_PATTERN", manifest, "")
	if err != nil {
		t.Fatalf("ResolveSlot failed: %v", err)
	}

	if content != ".moai/specs/SPEC-{ID}" {
		t.Errorf("ResolveSlot literal content: got %q, want %q", content, ".moai/specs/SPEC-{ID}")
	}
}

func TestResolveSlotWithDefaultFallback(t *testing.T) {
	reg := &Registry{
		Version: "1.0.0",
		Slots: map[string]*SlotEntry{
			"QUALITY_FRAMEWORK": {
				Category: "section",
				Default:  "### Quality\nGeneric compliance.",
			},
		},
	}

	// Manifest does not have QUALITY_FRAMEWORK in SlotContent.
	manifest := &model.PersonaManifest{
		SlotContent: map[string]string{},
	}

	content, err := reg.ResolveSlot("QUALITY_FRAMEWORK", manifest, "")
	if err != nil {
		t.Fatalf("ResolveSlot failed: %v", err)
	}

	expected := "### Quality\nGeneric compliance."
	if content != expected {
		t.Errorf("ResolveSlot default fallback\ngot:  %q\nwant: %q", content, expected)
	}
}

func TestResolveSlotNilManifest(t *testing.T) {
	reg := &Registry{
		Version: "1.0.0",
		Slots: map[string]*SlotEntry{
			"MY_SLOT": {
				Default: "default value",
			},
		},
	}

	content, err := reg.ResolveSlot("MY_SLOT", nil, "")
	if err != nil {
		t.Fatalf("ResolveSlot with nil manifest failed: %v", err)
	}
	if content != "default value" {
		t.Errorf("ResolveSlot nil manifest: got %q, want %q", content, "default value")
	}
}

func TestResolveSlotNotFound(t *testing.T) {
	reg := &Registry{
		Version: "1.0.0",
		Slots:   map[string]*SlotEntry{},
	}

	_, err := reg.ResolveSlot("NONEXISTENT", nil, "")
	if err == nil {
		t.Error("ResolveSlot should return error for nonexistent slot")
	}
}

func TestAddSlot(t *testing.T) {
	reg := NewRegistry()

	entry := &SlotEntry{
		Category:    "section",
		Scope:       "agent_body",
		Description: "Test slot",
		MarkerType:  "section",
		Default:     "default content",
	}
	reg.AddSlot("NEW_SLOT", entry)

	if len(reg.Slots) != 1 {
		t.Fatalf("Slots count after AddSlot: got %d, want 1", len(reg.Slots))
	}

	got := reg.Slots["NEW_SLOT"]
	if got == nil {
		t.Fatal("NEW_SLOT not found after AddSlot")
	}
	if got.Category != "section" {
		t.Errorf("NEW_SLOT.Category: got %q, want %q", got.Category, "section")
	}
	if got.Description != "Test slot" {
		t.Errorf("NEW_SLOT.Description: got %q, want %q", got.Description, "Test slot")
	}

	// Update existing slot.
	updated := &SlotEntry{
		Category:    "path_pattern",
		Description: "Updated slot",
	}
	reg.AddSlot("NEW_SLOT", updated)

	if reg.Slots["NEW_SLOT"].Category != "path_pattern" {
		t.Errorf("AddSlot update: Category got %q, want %q", reg.Slots["NEW_SLOT"].Category, "path_pattern")
	}
}

func TestLoadRegistryNotFound(t *testing.T) {
	_, err := LoadRegistry("/nonexistent/path")
	if err == nil {
		t.Error("LoadRegistry should return error for nonexistent directory")
	}
}
