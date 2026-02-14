package parser

import (
	"testing"

	"github.com/do-focus/convert/internal/model"
)

func TestSplitFrontmatter_WithValidFrontmatter(t *testing.T) {
	content := "---\nname: test-agent\ndescription: A test agent\n---\n# Body\nSome content\n"
	yamlContent, body, hasFM := SplitFrontmatter(content)

	if !hasFM {
		t.Fatal("expected frontmatter to be detected")
	}
	if yamlContent != "name: test-agent\ndescription: A test agent\n" {
		t.Errorf("unexpected yaml content: %q", yamlContent)
	}
	if body != "# Body\nSome content\n" {
		t.Errorf("unexpected body: %q", body)
	}
}

func TestSplitFrontmatter_NoFrontmatter(t *testing.T) {
	content := "# Just a heading\nSome content\n"
	yamlContent, body, hasFM := SplitFrontmatter(content)

	if hasFM {
		t.Fatal("expected no frontmatter")
	}
	if yamlContent != "" {
		t.Errorf("expected empty yaml content, got: %q", yamlContent)
	}
	if body != content {
		t.Errorf("expected body to be original content, got: %q", body)
	}
}

func TestSplitFrontmatter_EmptyFrontmatter(t *testing.T) {
	content := "---\n---\n# Body\n"
	yamlContent, body, hasFM := SplitFrontmatter(content)

	if !hasFM {
		t.Fatal("expected frontmatter to be detected")
	}
	if yamlContent != "" {
		t.Errorf("expected empty yaml content, got: %q", yamlContent)
	}
	if body != "# Body\n" {
		t.Errorf("unexpected body: %q", body)
	}
}

func TestSplitFrontmatter_NoClosingDelimiter(t *testing.T) {
	content := "---\nname: broken\n# No closing delimiter\n"
	_, body, hasFM := SplitFrontmatter(content)

	if hasFM {
		t.Fatal("expected no frontmatter when closing delimiter is missing")
	}
	if body != content {
		t.Errorf("expected body to be original content, got: %q", body)
	}
}

func TestSplitFrontmatter_DelimiterNotAtStart(t *testing.T) {
	content := "Some text before\n---\nname: test\n---\n"
	_, body, hasFM := SplitFrontmatter(content)

	if hasFM {
		t.Fatal("expected no frontmatter when delimiter is not at line 1")
	}
	if body != content {
		t.Errorf("expected body to be original content")
	}
}

func TestParseFrontmatter_FullAgent(t *testing.T) {
	yamlContent := "name: expert-backend\ndescription: Backend development specialist\ntools: Read Write Edit Grep Glob Bash\nmodel: inherit\npermissionMode: acceptEdits\nskills:\n  - moai-foundation-core\n  - moai-workflow-ddd\nmemory: project\n"

	fm, err := ParseFrontmatter(yamlContent)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if fm.Name != "expert-backend" {
		t.Errorf("expected name 'expert-backend', got %q", fm.Name)
	}
	if fm.Description != "Backend development specialist" {
		t.Errorf("unexpected description: %q", fm.Description)
	}
	if fm.Tools != "Read Write Edit Grep Glob Bash" {
		t.Errorf("unexpected tools: %q", fm.Tools)
	}
	if fm.Model != "inherit" {
		t.Errorf("unexpected model: %q", fm.Model)
	}
	if fm.PermissionMode != "acceptEdits" {
		t.Errorf("unexpected permissionMode: %q", fm.PermissionMode)
	}
	if len(fm.Skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(fm.Skills))
	}
	if fm.Skills[0] != "moai-foundation-core" || fm.Skills[1] != "moai-workflow-ddd" {
		t.Errorf("unexpected skills: %v", fm.Skills)
	}
	if fm.Memory != "project" {
		t.Errorf("unexpected memory: %q", fm.Memory)
	}
}

func TestParseFrontmatter_SkillsAsCommaSeparatedString(t *testing.T) {
	yamlContent := "name: test-agent\ndescription: Test\nskills: moai-foundation-core, moai-workflow-tdd\n"

	fm, err := ParseFrontmatter(yamlContent)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(fm.Skills) != 2 {
		t.Fatalf("expected 2 skills from comma-separated string, got %d: %v", len(fm.Skills), fm.Skills)
	}
	if fm.Skills[0] != "moai-foundation-core" || fm.Skills[1] != "moai-workflow-tdd" {
		t.Errorf("unexpected skills: %v", fm.Skills)
	}
}

func TestParseFrontmatter_UnknownFieldsPreserved(t *testing.T) {
	yamlContent := "name: test-agent\ndescription: Test\ncustom_field: custom_value\nmetadata:\n  version: \"1.0.0\"\n  tags: \"a, b, c\"\n"

	fm, err := ParseFrontmatter(yamlContent)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if fm.Raw["custom_field"] != "custom_value" {
		t.Errorf("expected custom_field in Raw, got: %v", fm.Raw["custom_field"])
	}
	metadata, ok := fm.Raw["metadata"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected metadata map in Raw, got: %T", fm.Raw["metadata"])
	}
	if metadata["version"] != "1.0.0" {
		t.Errorf("expected metadata.version '1.0.0', got: %v", metadata["version"])
	}
}

func TestSerializeFrontmatter_RoundTrip(t *testing.T) {
	original := &model.Frontmatter{
		Name:        "test-agent",
		Description: "A test agent",
		Tools:       "Read Write",
		Skills:      []string{"skill-a", "skill-b"},
		Raw: map[string]interface{}{
			"name":         "test-agent",
			"description":  "A test agent",
			"tools":        "Read Write",
			"skills":       []interface{}{"skill-a", "skill-b"},
			"custom_field": "preserved",
		},
	}

	serialized, err := SerializeFrontmatter(original)
	if err != nil {
		t.Fatalf("serialize error: %v", err)
	}

	// Verify it has delimiters
	if serialized[:4] != "---\n" {
		t.Errorf("expected opening delimiter, got: %q", serialized[:4])
	}
	if serialized[len(serialized)-4:] != "---\n" {
		t.Errorf("expected closing delimiter, got: %q", serialized[len(serialized)-4:])
	}

	// Parse back and verify
	yamlStr, _, hasFM := SplitFrontmatter(serialized)
	if !hasFM {
		t.Fatal("expected frontmatter in serialized output")
	}
	parsed, err := ParseFrontmatter(yamlStr)
	if err != nil {
		t.Fatalf("re-parse error: %v", err)
	}

	if parsed.Name != original.Name {
		t.Errorf("name mismatch: %q vs %q", parsed.Name, original.Name)
	}
	if parsed.Description != original.Description {
		t.Errorf("description mismatch: %q vs %q", parsed.Description, original.Description)
	}
	if parsed.Tools != original.Tools {
		t.Errorf("tools mismatch: %q vs %q", parsed.Tools, original.Tools)
	}
	if len(parsed.Skills) != len(original.Skills) {
		t.Errorf("skills count mismatch: %d vs %d", len(parsed.Skills), len(original.Skills))
	}
	if parsed.Raw["custom_field"] != "preserved" {
		t.Errorf("custom_field not preserved: %v", parsed.Raw["custom_field"])
	}
}

func TestSerializeFrontmatter_Nil(t *testing.T) {
	result, err := SerializeFrontmatter(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string for nil frontmatter, got: %q", result)
	}
}

func TestParseFrontmatter_EmptyYAML(t *testing.T) {
	fm, err := ParseFrontmatter("")
	if err != nil {
		t.Fatalf("parse error on empty yaml: %v", err)
	}
	if fm.Name != "" {
		t.Errorf("expected empty name, got: %q", fm.Name)
	}
	if fm.Raw == nil {
		t.Error("expected non-nil Raw map even for empty YAML")
	}
}
