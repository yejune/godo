package assembler

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yejune/godo/internal/model"
)

func TestBrandDeslotifier_DeslotifyContent(t *testing.T) {
	manifest := &model.PersonaManifest{
		Brand:    "moai",
		BrandDir: ".moai",
		BrandCmd: "/moai",
	}
	d := NewBrandDeslotifier(manifest)

	input := "Use {{slot:BRAND}} tool. Config in {{slot:BRAND_DIR}}/config. Run {{slot:BRAND_CMD}} plan."
	got := d.DeslotifyContent(input)
	expected := "Use moai tool. Config in .moai/config. Run /moai plan."
	if got != expected {
		t.Errorf("DeslotifyContent mismatch.\nexpected: %s\ngot:      %s", expected, got)
	}
}

func TestBrandDeslotifier_DeslotifyContent_BrandOnlyInfersDirAndCmd(t *testing.T) {
	// Only Brand set — BrandDir and BrandCmd are inferred as the brand name itself.
	// The leading "." and "/" are literal in slot patterns, not part of the value.
	manifest := &model.PersonaManifest{
		Brand: "do",
	}
	d := NewBrandDeslotifier(manifest)

	// Slot patterns: .{{slot:BRAND_DIR}}/ and /{{slot:BRAND_CMD}}
	input := "Use {{slot:BRAND}} with .{{slot:BRAND_DIR}}/ and /{{slot:BRAND_CMD}} plan."
	got := d.DeslotifyContent(input)
	expected := "Use do with .do/ and /do plan."
	if got != expected {
		t.Errorf("brand-only infer mismatch.\nexpected: %s\ngot:      %s", expected, got)
	}
}

func TestBrandDeslotifier_NilManifest(t *testing.T) {
	d := NewBrandDeslotifier(nil)
	if d != nil {
		t.Error("expected nil deslotifier for nil manifest")
	}

	// Nil-safe methods.
	got := d.DeslotifyContent("{{slot:BRAND}} test")
	if got != "{{slot:BRAND}} test" {
		t.Errorf("nil deslotifier should return content unchanged, got: %s", got)
	}

	remapped := d.RemapSkillPath("skills/lang-python/SKILL.md")
	if remapped != "skills/lang-python/SKILL.md" {
		t.Errorf("nil deslotifier should return path unchanged, got: %s", remapped)
	}
}

func TestBrandDeslotifier_EmptyBrandAndName(t *testing.T) {
	manifest := &model.PersonaManifest{Brand: "", Name: ""}
	d := NewBrandDeslotifier(manifest)
	if d != nil {
		t.Error("expected nil deslotifier for empty brand and name")
	}
}

func TestBrandDeslotifier_InferBrandFromName(t *testing.T) {
	// When Brand is empty but Name is set, brand should be inferred.
	// Inferred values have no prefix — the "." and "/" are literal in slot patterns.
	manifest := &model.PersonaManifest{Name: "moai"}
	d := NewBrandDeslotifier(manifest)
	if d == nil {
		t.Fatal("expected non-nil deslotifier when Name is set")
	}

	// Matches actual slot patterns from extractor: .{{slot:BRAND_DIR}}/ and /{{slot:BRAND_CMD}}
	input := "Use {{slot:BRAND}} tool. Config in .{{slot:BRAND_DIR}}/config. Run /{{slot:BRAND_CMD}} plan."
	got := d.DeslotifyContent(input)
	expected := "Use moai tool. Config in .moai/config. Run /moai plan."
	if got != expected {
		t.Errorf("infer brand mismatch.\nexpected: %s\ngot:      %s", expected, got)
	}
}

func TestBrandDeslotifier_ExplicitBrandOverridesName(t *testing.T) {
	// Explicit Brand takes precedence over Name.
	manifest := &model.PersonaManifest{
		Name:     "moai",
		Brand:    "do",
		BrandDir: ".do",
		BrandCmd: "/do",
	}
	d := NewBrandDeslotifier(manifest)
	if d == nil {
		t.Fatal("expected non-nil deslotifier")
	}

	input := "{{slot:BRAND}} at {{slot:BRAND_DIR}} run {{slot:BRAND_CMD}}"
	got := d.DeslotifyContent(input)
	expected := "do at .do run /do"
	if got != expected {
		t.Errorf("explicit brand mismatch.\nexpected: %s\ngot:      %s", expected, got)
	}
}

func TestBrandDeslotifier_RemapSkillPath(t *testing.T) {
	manifest := &model.PersonaManifest{Brand: "moai"}
	d := NewBrandDeslotifier(manifest)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "skill directory gets brand prefix",
			input:    "skills/lang-python/SKILL.md",
			expected: "skills/moai-lang-python/SKILL.md",
		},
		{
			name:     "nested skill file gets prefix",
			input:    "skills/domain-backend/modules/api.md",
			expected: "skills/moai-domain-backend/modules/api.md",
		},
		{
			name:     "non-skill path unchanged",
			input:    "agents/expert-backend.md",
			expected: "agents/expert-backend.md",
		},
		{
			name:     "rules path unchanged",
			input:    "rules/dev-testing.md",
			expected: "rules/dev-testing.md",
		},
		{
			name:     "skill with only one segment unchanged",
			input:    "skills/",
			expected: "skills/",
		},
		{
			name:     "already-prefixed skill not double-prefixed",
			input:    "skills/moai-lang-python/SKILL.md",
			expected: "skills/moai-lang-python/SKILL.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := d.RemapSkillPath(tt.input)
			if got != tt.expected {
				t.Errorf("RemapSkillPath(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestBrandDeslotifier_DeslotifyFile(t *testing.T) {
	manifest := &model.PersonaManifest{
		Brand:    "do",
		BrandDir: ".do",
		BrandCmd: "/do",
	}
	d := NewBrandDeslotifier(manifest)

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")
	content := "Use {{slot:BRAND}} tool at {{slot:BRAND_DIR}}/specs. Run {{slot:BRAND_CMD}} plan."
	if err := os.WriteFile(tmpFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := d.DeslotifyFile(tmpFile); err != nil {
		t.Fatalf("DeslotifyFile error: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	expected := "Use do tool at .do/specs. Run /do plan."
	if string(data) != expected {
		t.Errorf("DeslotifyFile content mismatch.\nexpected: %s\ngot:      %s", expected, string(data))
	}
}

func TestBrandDeslotifier_DeslotifyFile_NoSlots(t *testing.T) {
	manifest := &model.PersonaManifest{
		Brand:    "moai",
		BrandDir: ".moai",
		BrandCmd: "/moai",
	}
	d := NewBrandDeslotifier(manifest)

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "plain.md")
	content := "No brand slots here."
	if err := os.WriteFile(tmpFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := d.DeslotifyFile(tmpFile); err != nil {
		t.Fatalf("DeslotifyFile error: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != content {
		t.Errorf("expected unchanged content, got: %s", string(data))
	}
}

func TestBrandDeslotifier_PreservesOtherSlots(t *testing.T) {
	manifest := &model.PersonaManifest{
		Brand:    "moai",
		BrandDir: ".moai",
		BrandCmd: "/moai",
	}
	d := NewBrandDeslotifier(manifest)

	// BRAND slots should be replaced, but TOOL_NAME and SPEC_PATH should remain.
	input := "Use {{slot:BRAND}} with {{slot:TOOL_NAME}} at {{slot:BRAND_DIR}}/path."
	got := d.DeslotifyContent(input)
	expected := "Use moai with {{slot:TOOL_NAME}} at .moai/path."
	if got != expected {
		t.Errorf("preserve other slots mismatch.\nexpected: %s\ngot:      %s", expected, got)
	}
}
