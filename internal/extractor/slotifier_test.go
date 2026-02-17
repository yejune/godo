package extractor

import (
	"testing"
)

func TestNewBrandSlotifier_EmptyBrand(t *testing.T) {
	s := NewBrandSlotifier("")
	if s != nil {
		t.Error("NewBrandSlotifier(\"\") should return nil")
	}
}

func TestSlotifyContent_NilSlotifier(t *testing.T) {
	var s *BrandSlotifier
	got := s.SlotifyContent("moai-lang-python")
	if got != "moai-lang-python" {
		t.Errorf("nil SlotifyContent should return input unchanged, got %q", got)
	}
}

func TestSlotifyContent_CommandWithColon(t *testing.T) {
	s := NewBrandSlotifier("moai")
	tests := []struct {
		input string
		want  string
	}{
		{"/moai:1-plan", "/{{slot:BRAND_CMD}}:1-plan"},
		{"/moai:plan", "/{{slot:BRAND_CMD}}:plan"},
		{"Use /moai:sync to sync", "Use /{{slot:BRAND_CMD}}:sync to sync"},
	}
	for _, tt := range tests {
		got := s.SlotifyContent(tt.input)
		if got != tt.want {
			t.Errorf("SlotifyContent(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSlotifyContent_CommandWithSpace(t *testing.T) {
	s := NewBrandSlotifier("moai")
	tests := []struct {
		input string
		want  string
	}{
		{"/moai plan", "/{{slot:BRAND_CMD}} plan"},
		{"/moai run", "/{{slot:BRAND_CMD}} run"},
		{"/moai sync", "/{{slot:BRAND_CMD}} sync"},
		{"Use /moai plan to create", "Use /{{slot:BRAND_CMD}} plan to create"},
	}
	for _, tt := range tests {
		got := s.SlotifyContent(tt.input)
		if got != tt.want {
			t.Errorf("SlotifyContent(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSlotifyContent_DirectoryPaths(t *testing.T) {
	s := NewBrandSlotifier("moai")
	tests := []struct {
		input string
		want  string
	}{
		{".moai/specs/", ".{{slot:BRAND_DIR}}/specs/"},
		{".moai/config/sections/quality.yaml", ".{{slot:BRAND_DIR}}/config/sections/quality.yaml"},
		{".moai/logs/", ".{{slot:BRAND_DIR}}/logs/"},
		{".moai/project/", ".{{slot:BRAND_DIR}}/project/"},
		{".moai/memory/", ".{{slot:BRAND_DIR}}/memory/"},
		{"Read from .moai/specs/SPEC-001/spec.md", "Read from .{{slot:BRAND_DIR}}/specs/SPEC-001/spec.md"},
	}
	for _, tt := range tests {
		got := s.SlotifyContent(tt.input)
		if got != tt.want {
			t.Errorf("SlotifyContent(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSlotifyContent_SkillNamePrefix(t *testing.T) {
	s := NewBrandSlotifier("moai")
	tests := []struct {
		input string
		want  string
	}{
		{"name: moai-lang-python", "name: {{slot:BRAND}}-lang-python"},
		{"name: moai-domain-backend", "name: {{slot:BRAND}}-domain-backend"},
		{"name: moai-foundation-core", "name: {{slot:BRAND}}-foundation-core"},
		{"name: moai-workflow-spec", "name: {{slot:BRAND}}-workflow-spec"},
	}
	for _, tt := range tests {
		got := s.SlotifyContent(tt.input)
		if got != tt.want {
			t.Errorf("SlotifyContent(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSlotifyContent_RelatedSkills(t *testing.T) {
	s := NewBrandSlotifier("moai")
	input := `related-skills: "moai-workflow-spec, moai-workflow-ddd"`
	want := `related-skills: "{{slot:BRAND}}-workflow-spec, {{slot:BRAND}}-workflow-ddd"`
	got := s.SlotifyContent(input)
	if got != want {
		t.Errorf("SlotifyContent related-skills:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestSlotifyContent_SkillFunctionCalls(t *testing.T) {
	s := NewBrandSlotifier("moai")
	tests := []struct {
		input string
		want  string
	}{
		{`Skill("moai-domain-backend")`, `Skill("{{slot:BRAND}}-domain-backend")`},
		{`Skill("moai-foundation-core")`, `Skill("{{slot:BRAND}}-foundation-core")`},
		{`Skill("moai-workflow-project")`, `Skill("{{slot:BRAND}}-workflow-project")`},
	}
	for _, tt := range tests {
		got := s.SlotifyContent(tt.input)
		if got != tt.want {
			t.Errorf("SlotifyContent(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSlotifyContent_CrossSkillRefsInBody(t *testing.T) {
	s := NewBrandSlotifier("moai")
	tests := []struct {
		input string
		want  string
	}{
		{"see moai-foundation-core for details", "see {{slot:BRAND}}-foundation-core for details"},
		{"  - moai-workflow-project", "  - {{slot:BRAND}}-workflow-project"},
		{"moai-lang-go and moai-lang-python", "{{slot:BRAND}}-lang-go and {{slot:BRAND}}-lang-python"},
	}
	for _, tt := range tests {
		got := s.SlotifyContent(tt.input)
		if got != tt.want {
			t.Errorf("SlotifyContent(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSlotifyContent_NoFalsePositives(t *testing.T) {
	s := NewBrandSlotifier("moai")
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "standalone brand word replaced",
			input: "The moai system is great",
			want:  "The {{slot:BRAND}} system is great",
		},
		{
			name:  "brand case-insensitive replaced",
			input: "MoAI Constitution",
			want:  "{{slot:BRAND}} Constitution",
		},
		{
			name:  "brand inside another word",
			input: "samoai is not a brand ref",
			want:  "samoai is not a brand ref",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.SlotifyContent(tt.input)
			if got != tt.want {
				t.Errorf("SlotifyContent(%q) = %q, want %q (false positive)", tt.input, got, tt.want)
			}
		})
	}
}

func TestSlotifyContent_MultipleReplacements(t *testing.T) {
	s := NewBrandSlotifier("moai")
	input := `---
name: moai-workflow-spec
description: >
  SPEC workflow skill using /moai plan and /moai:1-plan commands.
  Reads from .moai/specs/ directory.
related-skills: "moai-foundation-core, moai-workflow-ddd"
---

Use Skill("moai-domain-backend") for backend tasks.
Config is at .moai/config/sections/quality.yaml.
`
	want := `---
name: {{slot:BRAND}}-workflow-spec
description: >
  SPEC workflow skill using /{{slot:BRAND_CMD}} plan and /{{slot:BRAND_CMD}}:1-plan commands.
  Reads from .{{slot:BRAND_DIR}}/specs/ directory.
related-skills: "{{slot:BRAND}}-foundation-core, {{slot:BRAND}}-workflow-ddd"
---

Use Skill("{{slot:BRAND}}-domain-backend") for backend tasks.
Config is at .{{slot:BRAND_DIR}}/config/sections/quality.yaml.
`
	got := s.SlotifyContent(input)
	if got != want {
		t.Errorf("SlotifyContent (multi):\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestStripBrandPrefix(t *testing.T) {
	s := NewBrandSlotifier("moai")
	tests := []struct {
		input string
		want  string
	}{
		{"moai-lang-python", "lang-python"},
		{"moai-domain-backend", "domain-backend"},
		{"moai-foundation-core", "foundation-core"},
		{"moai-workflow-spec", "workflow-spec"},
		{"do-foundation-claude", "do-foundation-claude"},   // no moai prefix
		{"moai", "moai"},                                    // no hyphen after brand
		{"", ""},                                             // empty
		{"lang-python", "lang-python"},                      // no brand prefix
	}
	for _, tt := range tests {
		got := s.StripBrandPrefix(tt.input)
		if got != tt.want {
			t.Errorf("StripBrandPrefix(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestStripBrandPrefix_NilSlotifier(t *testing.T) {
	var s *BrandSlotifier
	got := s.StripBrandPrefix("moai-lang-python")
	if got != "moai-lang-python" {
		t.Errorf("nil StripBrandPrefix should return input unchanged, got %q", got)
	}
}

func TestRemapCorePath(t *testing.T) {
	s := NewBrandSlotifier("moai")
	tests := []struct {
		input string
		want  string
	}{
		{"skills/moai-lang-python/SKILL.md", "skills/lang-python/SKILL.md"},
		{"skills/moai-domain-backend/modules/api.md", "skills/domain-backend/modules/api.md"},
		{"skills/moai-foundation-core/SKILL.md", "skills/foundation-core/SKILL.md"},
		{"skills/moai-workflow-spec/reference.md", "skills/workflow-spec/reference.md"},
		{"skills/do-foundation-claude/SKILL.md", "skills/do-foundation-claude/SKILL.md"}, // no moai prefix
		{"agents/expert-backend.md", "agents/expert-backend.md"},                          // not a skill
		{"rules/dev-testing.md", "rules/dev-testing.md"},                                  // not a skill
		{"skills/short", "skills/short"},                                                   // too short path
	}
	for _, tt := range tests {
		got := s.RemapCorePath(tt.input)
		if got != tt.want {
			t.Errorf("RemapCorePath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRemapCorePath_NilSlotifier(t *testing.T) {
	var s *BrandSlotifier
	got := s.RemapCorePath("skills/moai-lang-python/SKILL.md")
	if got != "skills/moai-lang-python/SKILL.md" {
		t.Errorf("nil RemapCorePath should return input unchanged, got %q", got)
	}
}
