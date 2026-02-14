package template

import (
	"testing"
)

func TestInsertSectionSlot(t *testing.T) {
	result := InsertSectionSlot("QUALITY_FRAMEWORK", "### Quality\nFollow standards.")
	expected := "<!-- BEGIN_SLOT:QUALITY_FRAMEWORK -->\n### Quality\nFollow standards.\n<!-- END_SLOT:QUALITY_FRAMEWORK -->"
	if result != expected {
		t.Errorf("InsertSectionSlot mismatch\ngot:  %q\nwant: %q", result, expected)
	}
}

func TestInsertSectionSlotEmpty(t *testing.T) {
	result := InsertSectionSlot("EMPTY_SLOT", "")
	expected := "<!-- BEGIN_SLOT:EMPTY_SLOT -->\n\n<!-- END_SLOT:EMPTY_SLOT -->"
	if result != expected {
		t.Errorf("InsertSectionSlot empty content mismatch\ngot:  %q\nwant: %q", result, expected)
	}
}

func TestInsertInlineSlot(t *testing.T) {
	result := InsertInlineSlot("SPEC_PATH_PATTERN")
	expected := "{{SPEC_PATH_PATTERN}}"
	if result != expected {
		t.Errorf("InsertInlineSlot mismatch\ngot:  %q\nwant: %q", result, expected)
	}
}

func TestExtractSectionSlot(t *testing.T) {
	content := "some preamble\n<!-- BEGIN_SLOT:MY_SLOT -->\nhello world\n<!-- END_SLOT:MY_SLOT -->\nsome epilogue"
	extracted, found := ExtractSectionSlot(content, "MY_SLOT")
	if !found {
		t.Fatal("ExtractSectionSlot: expected to find MY_SLOT but did not")
	}
	if extracted != "hello world" {
		t.Errorf("ExtractSectionSlot content mismatch\ngot:  %q\nwant: %q", extracted, "hello world")
	}
}

func TestExtractSectionSlotMultiLine(t *testing.T) {
	inner := "line one\nline two\nline three"
	wrapped := InsertSectionSlot("MULTI", inner)
	content := "before\n" + wrapped + "\nafter"

	extracted, found := ExtractSectionSlot(content, "MULTI")
	if !found {
		t.Fatal("ExtractSectionSlot: expected to find MULTI but did not")
	}
	if extracted != inner {
		t.Errorf("ExtractSectionSlot multiline mismatch\ngot:  %q\nwant: %q", extracted, inner)
	}
}

func TestExtractSectionSlotNotFound(t *testing.T) {
	_, found := ExtractSectionSlot("no markers here", "MISSING")
	if found {
		t.Error("ExtractSectionSlot: expected not found for MISSING slot")
	}
}

func TestReplaceInlineSlots(t *testing.T) {
	content := "Read specs from {{SPEC_PATH_PATTERN}} and config from {{QUALITY_CONFIG_PATH}}."
	values := map[string]string{
		"SPEC_PATH_PATTERN":  ".moai/specs/SPEC-{ID}",
		"QUALITY_CONFIG_PATH": ".moai/config/sections/quality.yaml",
	}
	result, replaced := ReplaceInlineSlots(content, values)

	expected := "Read specs from .moai/specs/SPEC-{ID} and config from .moai/config/sections/quality.yaml."
	if result != expected {
		t.Errorf("ReplaceInlineSlots content mismatch\ngot:  %q\nwant: %q", result, expected)
	}
	if len(replaced) != 2 {
		t.Errorf("ReplaceInlineSlots replaced count: got %d, want 2", len(replaced))
	}
}

func TestReplaceInlineSlotsPartial(t *testing.T) {
	content := "{{KNOWN}} stays but {{UNKNOWN}} is left"
	values := map[string]string{
		"KNOWN": "resolved",
	}
	result, replaced := ReplaceInlineSlots(content, values)

	expected := "resolved stays but {{UNKNOWN}} is left"
	if result != expected {
		t.Errorf("ReplaceInlineSlots partial mismatch\ngot:  %q\nwant: %q", result, expected)
	}
	if len(replaced) != 1 || replaced[0] != "KNOWN" {
		t.Errorf("ReplaceInlineSlots replaced: got %v, want [KNOWN]", replaced)
	}
}

func TestReplaceInlineSlotsNoMatch(t *testing.T) {
	content := "no slots here"
	result, replaced := ReplaceInlineSlots(content, map[string]string{"X": "y"})
	if result != content {
		t.Errorf("ReplaceInlineSlots should not modify content without markers")
	}
	if len(replaced) != 0 {
		t.Errorf("ReplaceInlineSlots replaced should be empty, got %v", replaced)
	}
}

func TestFindAllSlotMarkers(t *testing.T) {
	content := `Some text with {{INLINE_A}} and {{INLINE_B}}.
<!-- BEGIN_SLOT:SECTION_C -->
section content
<!-- END_SLOT:SECTION_C -->
Also {{INLINE_A}} appears again.`

	ids := FindAllSlotMarkers(content)
	expected := []string{"INLINE_A", "INLINE_B", "SECTION_C"}

	if len(ids) != len(expected) {
		t.Fatalf("FindAllSlotMarkers count: got %d, want %d\ngot: %v", len(ids), len(expected), ids)
	}
	for i, id := range ids {
		if id != expected[i] {
			t.Errorf("FindAllSlotMarkers[%d]: got %q, want %q", i, id, expected[i])
		}
	}
}

func TestFindAllSlotMarkersEmpty(t *testing.T) {
	ids := FindAllSlotMarkers("no markers at all")
	if len(ids) != 0 {
		t.Errorf("FindAllSlotMarkers should return empty for no markers, got %v", ids)
	}
}
