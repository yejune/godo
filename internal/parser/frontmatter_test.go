package parser

import (
	"strings"
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

func TestPatchFrontmatterSkills_InlineFormat_Append(t *testing.T) {
	rawYaml := "name: builder-agent\ndescription: Agent creation specialist\nskills: moai-foundation-claude, moai-workflow-project\nmemory: user\n"
	newSkills := []string{"moai-foundation-claude", "moai-workflow-project", "moai-persona-custom"}

	result := PatchFrontmatterSkills(rawYaml, newSkills)

	expected := "name: builder-agent\ndescription: Agent creation specialist\nskills: moai-foundation-claude, moai-workflow-project, moai-persona-custom\nmemory: user\n"
	if result != expected {
		t.Errorf("inline append mismatch.\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestPatchFrontmatterSkills_InlineFormat_Remove(t *testing.T) {
	rawYaml := "name: test-agent\nskills: skill-a, skill-b, skill-c\nmemory: user\n"
	newSkills := []string{"skill-a", "skill-c"}

	result := PatchFrontmatterSkills(rawYaml, newSkills)

	expected := "name: test-agent\nskills: skill-a, skill-c\nmemory: user\n"
	if result != expected {
		t.Errorf("inline remove mismatch.\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestPatchFrontmatterSkills_ListFormat_Append(t *testing.T) {
	rawYaml := "name: expert-backend\nskills:\n    - do-foundation\n    - do-backend\nmemory: project\n"
	newSkills := []string{"do-foundation", "do-backend", "do-quality"}

	result := PatchFrontmatterSkills(rawYaml, newSkills)

	expected := "name: expert-backend\nskills:\n    - do-foundation\n    - do-backend\n    - do-quality\nmemory: project\n"
	if result != expected {
		t.Errorf("list append mismatch.\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestPatchFrontmatterSkills_ListFormat_Remove(t *testing.T) {
	rawYaml := "name: expert-backend\nskills:\n    - do-foundation\n    - do-legacy\n    - do-backend\nmemory: project\n"
	newSkills := []string{"do-foundation", "do-backend"}

	result := PatchFrontmatterSkills(rawYaml, newSkills)

	expected := "name: expert-backend\nskills:\n    - do-foundation\n    - do-backend\nmemory: project\n"
	if result != expected {
		t.Errorf("list remove mismatch.\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestPatchFrontmatterSkills_NoExistingSkills_AddNew(t *testing.T) {
	rawYaml := "name: test-agent\ndescription: Test agent\nmemory: user\n"
	newSkills := []string{"skill-a", "skill-b"}

	result := PatchFrontmatterSkills(rawYaml, newSkills)

	expected := "name: test-agent\ndescription: Test agent\nmemory: user\nskills: skill-a, skill-b\n"
	if result != expected {
		t.Errorf("no existing skills mismatch.\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestPatchFrontmatterSkills_RemoveAllSkills(t *testing.T) {
	rawYaml := "name: test-agent\nskills: skill-a, skill-b\nmemory: user\n"
	newSkills := []string{}

	result := PatchFrontmatterSkills(rawYaml, newSkills)

	if strings.Contains(result, "skills") {
		t.Errorf("expected skills field removed, got:\n%s", result)
	}
	if !strings.Contains(result, "name: test-agent") {
		t.Errorf("expected other fields preserved, got:\n%s", result)
	}
}

func TestPatchFrontmatterSkills_RemoveAllSkills_ListFormat(t *testing.T) {
	rawYaml := "name: test-agent\nskills:\n    - skill-a\n    - skill-b\nmemory: user\n"
	newSkills := []string{}

	result := PatchFrontmatterSkills(rawYaml, newSkills)

	if strings.Contains(result, "skills") {
		t.Errorf("expected skills field removed, got:\n%s", result)
	}
	if !strings.Contains(result, "name: test-agent") || !strings.Contains(result, "memory: user") {
		t.Errorf("expected other fields preserved, got:\n%s", result)
	}
}

func TestPatchFrontmatterSkills_PreservesKeyOrder(t *testing.T) {
	rawYaml := "description: |\n  Agent creation specialist\nmemory: user\nmodel: inherit\npermissionMode: bypassPermissions\nskills: moai-foundation-claude, moai-workflow-project\n"
	newSkills := []string{"moai-foundation-claude", "moai-workflow-project", "moai-extra"}

	result := PatchFrontmatterSkills(rawYaml, newSkills)

	// Verify key order is preserved: description, memory, model, permissionMode, skills
	descIdx := strings.Index(result, "description:")
	memIdx := strings.Index(result, "memory:")
	modelIdx := strings.Index(result, "model:")
	permIdx := strings.Index(result, "permissionMode:")
	skillsIdx := strings.Index(result, "skills:")

	if descIdx >= memIdx || memIdx >= modelIdx || modelIdx >= permIdx || permIdx >= skillsIdx {
		t.Errorf("key order not preserved.\ndesc=%d, mem=%d, model=%d, perm=%d, skills=%d\nresult:\n%s",
			descIdx, memIdx, modelIdx, permIdx, skillsIdx, result)
	}
}

func TestPatchFrontmatterSkills_NoSkills_NoChange(t *testing.T) {
	rawYaml := "name: test-agent\ndescription: Test\n"
	result := PatchFrontmatterSkills(rawYaml, nil)
	if result != rawYaml {
		t.Errorf("expected no change when no skills and nil newSkills.\nexpected:\n%s\ngot:\n%s", rawYaml, result)
	}
}

func TestPatchFrontmatterSkills_ListFormat_2SpaceIndent(t *testing.T) {
	rawYaml := "name: test\nskills:\n  - skill-a\n  - skill-b\n"
	newSkills := []string{"skill-a", "skill-b", "skill-c"}

	result := PatchFrontmatterSkills(rawYaml, newSkills)

	expected := "name: test\nskills:\n  - skill-a\n  - skill-b\n  - skill-c\n"
	if result != expected {
		t.Errorf("2-space indent mismatch.\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}
