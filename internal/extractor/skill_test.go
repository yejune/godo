package extractor

import (
	"testing"

	"github.com/do-focus/convert/internal/detector"
	"github.com/do-focus/convert/internal/model"
)

func TestSkillExtractor_WholeFilePersona(t *testing.T) {
	reg := detector.NewDefaultRegistry()
	ext := NewSkillExtractor(reg)

	doc := &model.Document{
		Path: "skills/moai-workflow-ddd/SKILL.md",
		Frontmatter: &model.Frontmatter{
			Name: "moai-workflow-ddd",
		},
	}

	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}
	if coreDoc != nil {
		t.Error("expected nil coreDoc for whole-file persona skill")
	}
	if len(manifest.Skills) != 1 || manifest.Skills[0] != doc.Path {
		t.Errorf("manifest.Skills = %v, want [%q]", manifest.Skills, doc.Path)
	}
}

func TestSkillExtractor_CoreSkill(t *testing.T) {
	reg := detector.NewDefaultRegistry()
	ext := NewSkillExtractor(reg)

	doc := &model.Document{
		Path: "skills/do-foundation-claude/SKILL.md",
		Frontmatter: &model.Frontmatter{
			Name: "do-foundation-claude",
		},
	}

	coreDoc, manifest, err := ext.Extract(doc)
	if err != nil {
		t.Fatalf("Extract() error: %v", err)
	}
	if coreDoc == nil {
		t.Error("expected non-nil coreDoc for core skill")
	}
	if len(manifest.Skills) != 0 {
		t.Errorf("manifest.Skills = %v, want empty", manifest.Skills)
	}
}

func TestSkillExtractor_PartialSkill_PersonaModule(t *testing.T) {
	reg := detector.NewDefaultRegistry()
	ext := NewSkillExtractor(reg)

	tests := []struct {
		name string
		path string
	}{
		{
			name: "ddd-context7 root module",
			path: "skills/moai-workflow-testing/modules/ddd-context7.md",
		},
		{
			name: "ddd-context7 submodule",
			path: "skills/moai-workflow-testing/modules/ddd-context7/advanced-features.md",
		},
		{
			name: "ddd submodule",
			path: "skills/moai-workflow-testing/modules/ddd/core-classes.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &model.Document{
				Path: tt.path,
				// Module files typically don't have name frontmatter
				Frontmatter: &model.Frontmatter{},
			}

			coreDoc, manifest, err := ext.Extract(doc)
			if err != nil {
				t.Fatalf("Extract() error: %v", err)
			}
			if coreDoc != nil {
				t.Error("expected nil coreDoc for persona module")
			}
			if len(manifest.Skills) != 1 || manifest.Skills[0] != tt.path {
				t.Errorf("manifest.Skills = %v, want [%q]", manifest.Skills, tt.path)
			}
		})
	}
}

func TestSkillExtractor_PartialSkill_CoreModule(t *testing.T) {
	reg := detector.NewDefaultRegistry()
	ext := NewSkillExtractor(reg)

	tests := []struct {
		name string
		path string
	}{
		{
			name: "SKILL.md stays core",
			path: "skills/moai-workflow-testing/SKILL.md",
		},
		{
			name: "ai-debugging module stays core",
			path: "skills/moai-workflow-testing/modules/ai-debugging.md",
		},
		{
			name: "debugging subdir stays core",
			path: "skills/moai-workflow-testing/modules/debugging/debugging-workflows.md",
		},
		{
			name: "performance module stays core",
			path: "skills/moai-workflow-testing/modules/performance/optimization-patterns.md",
		},
		{
			name: "smart-refactoring stays core",
			path: "skills/moai-workflow-testing/modules/smart-refactoring.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &model.Document{
				Path: tt.path,
				Frontmatter: &model.Frontmatter{},
			}

			coreDoc, manifest, err := ext.Extract(doc)
			if err != nil {
				t.Fatalf("Extract() error: %v", err)
			}
			if coreDoc == nil {
				t.Error("expected non-nil coreDoc for core module")
			}
			if len(manifest.Skills) != 0 {
				t.Errorf("manifest.Skills = %v, want empty", manifest.Skills)
			}
		})
	}
}

func TestExtractSkillName(t *testing.T) {
	tests := []struct {
		relPath string
		want    string
	}{
		{"skills/moai-workflow-testing/SKILL.md", "moai-workflow-testing"},
		{"skills/moai-workflow-testing/modules/ddd/core.md", "moai-workflow-testing"},
		{"skills/do-foundation/SKILL.md", "do-foundation"},
		{"agents/expert-backend.md", ""},
		{"skills/short", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.relPath, func(t *testing.T) {
			got := extractSkillName(tt.relPath)
			if got != tt.want {
				t.Errorf("extractSkillName(%q) = %q, want %q", tt.relPath, got, tt.want)
			}
		})
	}
}

func TestExtractModuleRelPath(t *testing.T) {
	tests := []struct {
		relPath   string
		skillName string
		want      string
	}{
		{
			"skills/moai-workflow-testing/modules/ddd/core.md",
			"moai-workflow-testing",
			"modules/ddd/core.md",
		},
		{
			"skills/moai-workflow-testing/SKILL.md",
			"moai-workflow-testing",
			"SKILL.md",
		},
		{
			"skills/other-skill/modules/foo.md",
			"moai-workflow-testing",
			"",
		},
		{
			"agents/expert-backend.md",
			"moai-workflow-testing",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.relPath, func(t *testing.T) {
			got := extractModuleRelPath(tt.relPath, tt.skillName)
			if got != tt.want {
				t.Errorf("extractModuleRelPath(%q, %q) = %q, want %q",
					tt.relPath, tt.skillName, got, tt.want)
			}
		})
	}
}
