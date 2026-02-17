package persona

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_ParseFrontmatter_valid(t *testing.T) {
	content := []byte("---\nid: young-f\nname: Test\n---\nBody content here")

	yamlBytes, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("ParseFrontmatter: %v", err)
	}

	if string(yamlBytes) != "id: young-f\nname: Test" {
		t.Errorf("yaml: got %q", string(yamlBytes))
	}
	if body != "Body content here" {
		t.Errorf("body: got %q", body)
	}
}

func Test_ParseFrontmatter_no_frontmatter(t *testing.T) {
	content := []byte("No frontmatter here")
	_, _, err := ParseFrontmatter(content)
	if err == nil {
		t.Error("expected error for missing frontmatter")
	}
}

func Test_ParseFrontmatter_no_closing(t *testing.T) {
	content := []byte("---\nid: test\nNo closing delimiter")
	_, _, err := ParseFrontmatter(content)
	if err == nil {
		t.Error("expected error for missing closing ---")
	}
}

func Test_LoadCharacter_valid(t *testing.T) {
	tmp := t.TempDir()
	charDir := filepath.Join(tmp, "characters")
	os.MkdirAll(charDir, 0755)

	content := `---
id: young-f
name: "Test Character"
honorific_template: "{{name}}선배"
honorific_default: "선배님"
tone: "반말+존댓말 혼합"
character_summary: "bright developer"
relationship: "colleague"
---
Character body markdown here.
`
	os.WriteFile(filepath.Join(charDir, "young-f.md"), []byte(content), 0644)

	data, err := LoadCharacter(tmp, "young-f")
	if err != nil {
		t.Fatalf("LoadCharacter: %v", err)
	}

	if data.ID != "young-f" {
		t.Errorf("ID: got %q, want %q", data.ID, "young-f")
	}
	if data.Name != "Test Character" {
		t.Errorf("Name: got %q, want %q", data.Name, "Test Character")
	}
	if data.HonorificTemplate != "{{name}}선배" {
		t.Errorf("HonorificTemplate: got %q", data.HonorificTemplate)
	}
	if data.HonorificDefault != "선배님" {
		t.Errorf("HonorificDefault: got %q", data.HonorificDefault)
	}
	if data.FullContent != "Character body markdown here." {
		t.Errorf("FullContent: got %q", data.FullContent)
	}
}

func Test_LoadCharacter_missing_file(t *testing.T) {
	tmp := t.TempDir()
	_, err := LoadCharacter(tmp, "nonexistent")
	if err == nil {
		t.Error("expected error for missing character file")
	}
}

func Test_BuildHonorific_with_name(t *testing.T) {
	d := &Data{
		HonorificTemplate: "{{name}}선배",
		HonorificDefault:  "선배님",
	}
	got := d.BuildHonorific("승민")
	if got != "승민선배" {
		t.Errorf("got %q, want %q", got, "승민선배")
	}
}

func Test_BuildHonorific_empty_name(t *testing.T) {
	d := &Data{
		HonorificTemplate: "{{name}}선배",
		HonorificDefault:  "선배님",
	}
	got := d.BuildHonorific("")
	if got != "선배님" {
		t.Errorf("got %q, want %q", got, "선배님")
	}
}

func Test_BuildReminder(t *testing.T) {
	d := &Data{
		HonorificTemplate: "{{name}}선배",
		HonorificDefault:  "선배님",
		Tone:              "반말+존댓말 혼합",
	}

	got := d.BuildReminder("승민")
	expected := `반드시 "승민선배"로 호칭할 것. 말투: 반말+존댓말 혼합`
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func Test_BuildReminder_empty_honorific(t *testing.T) {
	d := &Data{
		HonorificTemplate: "",
		HonorificDefault:  "",
		Tone:              "test",
	}
	got := d.BuildReminder("")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func Test_LoadSpinner_valid(t *testing.T) {
	tmp := t.TempDir()
	spinnerDir := filepath.Join(tmp, "spinners")
	os.MkdirAll(spinnerDir, 0755)

	content := `persona: young-f
suffix_pattern:
  cycle: 2
  suffixes: ["는 중", "고 있어요"]
stems:
  - stem: "코딩하"
    emoji: ""
  - stem: "분석하"
    emoji: ""
`
	os.WriteFile(filepath.Join(spinnerDir, "young-f.yaml"), []byte(content), 0644)

	sd, err := LoadSpinner(tmp, "young-f")
	if err != nil {
		t.Fatalf("LoadSpinner: %v", err)
	}

	if sd.Persona != "young-f" {
		t.Errorf("Persona: got %q", sd.Persona)
	}
	if len(sd.Stems) != 2 {
		t.Fatalf("expected 2 stems, got %d", len(sd.Stems))
	}
	if sd.Stems[0].Stem != "코딩하" {
		t.Errorf("first stem: got %q", sd.Stems[0].Stem)
	}
}

func Test_BuildSpinnerVerbs(t *testing.T) {
	sd := &SpinnerData{
		SuffixPattern: SpinnerSuffixPattern{
			Cycle:    2,
			Suffixes: []string{"는 중", "고 있어요"},
		},
		Stems: []SpinnerStem{
			{Stem: "코딩하", Emoji: ""},
			{Stem: "분석하", Emoji: ""},
			{Stem: "테스트하", Emoji: ""},
		},
	}

	verbs := sd.BuildSpinnerVerbs()
	if len(verbs) != 3 {
		t.Fatalf("expected 3 verbs, got %d", len(verbs))
	}
	if verbs[0] != "코딩하는 중" {
		t.Errorf("verb[0]: got %q, want %q", verbs[0], "코딩하는 중")
	}
	if verbs[1] != "분석하고 있어요" {
		t.Errorf("verb[1]: got %q, want %q", verbs[1], "분석하고 있어요")
	}
	if verbs[2] != "테스트하는 중" {
		t.Errorf("verb[2]: got %q, want %q", verbs[2], "테스트하는 중")
	}
}

func Test_BuildSpinnerVerbs_with_emoji(t *testing.T) {
	sd := &SpinnerData{
		SuffixPattern: SpinnerSuffixPattern{
			Cycle:    1,
			Suffixes: []string{"는 중"},
		},
		Stems: []SpinnerStem{
			{Stem: "코딩하", Emoji: "💻"},
		},
	}

	verbs := sd.BuildSpinnerVerbs()
	if len(verbs) != 1 {
		t.Fatalf("expected 1 verb, got %d", len(verbs))
	}
	if verbs[0] != "코딩하는 중 💻" {
		t.Errorf("got %q, want %q", verbs[0], "코딩하는 중 💻")
	}
}

func Test_BuildSpinnerVerbs_empty_suffixes(t *testing.T) {
	sd := &SpinnerData{
		SuffixPattern: SpinnerSuffixPattern{
			Cycle:    1,
			Suffixes: nil,
		},
		Stems: []SpinnerStem{
			{Stem: "test"},
		},
	}

	verbs := sd.BuildSpinnerVerbs()
	if verbs != nil {
		t.Errorf("expected nil, got %v", verbs)
	}
}
