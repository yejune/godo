package parser

import (
	"testing"
)

func TestParseSections_BasicSplitting(t *testing.T) {
	body := "# Heading 1\nContent under h1\n## Heading 2\nContent under h2\n# Heading 3\nContent under h3\n"
	sections := ParseSections(body)

	if len(sections) != 2 {
		t.Fatalf("expected 2 root sections, got %d", len(sections))
	}
	if sections[0].Title != "Heading 1" {
		t.Errorf("expected title 'Heading 1', got %q", sections[0].Title)
	}
	if sections[0].Level != 1 {
		t.Errorf("expected level 1, got %d", sections[0].Level)
	}
	if sections[1].Title != "Heading 3" {
		t.Errorf("expected title 'Heading 3', got %q", sections[1].Title)
	}
}

func TestParseSections_Nesting(t *testing.T) {
	body := "# Top\nTop content\n## Child A\nChild A content\n### Grandchild\nGrandchild content\n## Child B\nChild B content\n"
	sections := ParseSections(body)

	if len(sections) != 1 {
		t.Fatalf("expected 1 root section, got %d", len(sections))
	}
	root := sections[0]
	if root.Title != "Top" {
		t.Errorf("expected root title 'Top', got %q", root.Title)
	}
	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(root.Children))
	}
	if root.Children[0].Title != "Child A" {
		t.Errorf("expected child 'Child A', got %q", root.Children[0].Title)
	}
	if len(root.Children[0].Children) != 1 {
		t.Fatalf("expected 1 grandchild, got %d", len(root.Children[0].Children))
	}
	if root.Children[0].Children[0].Title != "Grandchild" {
		t.Errorf("expected grandchild 'Grandchild', got %q", root.Children[0].Children[0].Title)
	}
	if root.Children[1].Title != "Child B" {
		t.Errorf("expected child 'Child B', got %q", root.Children[1].Title)
	}
}

func TestParseSections_Preamble(t *testing.T) {
	body := "This is preamble text.\nBefore any header.\n# First Header\nHeader content\n"
	sections := ParseSections(body)

	if len(sections) != 2 {
		t.Fatalf("expected 2 sections (preamble + header), got %d", len(sections))
	}

	preamble := sections[0]
	if preamble.Level != 0 {
		t.Errorf("expected preamble level 0, got %d", preamble.Level)
	}
	if preamble.Title != "" {
		t.Errorf("expected empty preamble title, got %q", preamble.Title)
	}
	if preamble.StartLine != 1 {
		t.Errorf("expected preamble start line 1, got %d", preamble.StartLine)
	}

	header := sections[1]
	if header.Title != "First Header" {
		t.Errorf("expected title 'First Header', got %q", header.Title)
	}
	if header.Level != 1 {
		t.Errorf("expected level 1, got %d", header.Level)
	}
}

func TestParseSections_CodeBlockSkipsHeaders(t *testing.T) {
	body := "# Real Header\nSome content\n```\n# Not a header\n## Also not a header\n```\n## Real Subheader\nMore content\n"
	sections := ParseSections(body)

	if len(sections) != 1 {
		t.Fatalf("expected 1 root section, got %d", len(sections))
	}
	root := sections[0]
	if root.Title != "Real Header" {
		t.Errorf("expected 'Real Header', got %q", root.Title)
	}
	// The code block headers should not create separate sections
	if len(root.Children) != 1 {
		t.Fatalf("expected 1 child (Real Subheader), got %d", len(root.Children))
	}
	if root.Children[0].Title != "Real Subheader" {
		t.Errorf("expected 'Real Subheader', got %q", root.Children[0].Title)
	}
}

func TestParseSections_TildeCodeBlock(t *testing.T) {
	body := "# Header\n~~~\n# Fake header inside tilde block\n~~~\n## Sub\nContent\n"
	sections := ParseSections(body)

	if len(sections) != 1 {
		t.Fatalf("expected 1 root, got %d", len(sections))
	}
	if len(sections[0].Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(sections[0].Children))
	}
	if sections[0].Children[0].Title != "Sub" {
		t.Errorf("expected child title 'Sub', got %q", sections[0].Children[0].Title)
	}
}

func TestParseSections_SkippedLevels(t *testing.T) {
	// Test ## followed by #### (skipping ###)
	body := "## Level 2\nContent\n#### Level 4\nDeep content\n## Another Level 2\nMore content\n"
	sections := ParseSections(body)

	if len(sections) != 2 {
		t.Fatalf("expected 2 root sections, got %d", len(sections))
	}
	if sections[0].Title != "Level 2" {
		t.Errorf("expected 'Level 2', got %q", sections[0].Title)
	}
	if len(sections[0].Children) != 1 {
		t.Fatalf("expected 1 child of Level 2, got %d", len(sections[0].Children))
	}
	if sections[0].Children[0].Title != "Level 4" {
		t.Errorf("expected 'Level 4', got %q", sections[0].Children[0].Title)
	}
	if sections[1].Title != "Another Level 2" {
		t.Errorf("expected 'Another Level 2', got %q", sections[1].Title)
	}
}

func TestParseSections_EmptyBody(t *testing.T) {
	sections := ParseSections("")
	if len(sections) != 0 {
		t.Errorf("expected 0 sections for empty body, got %d", len(sections))
	}
}

func TestParseSections_LineNumbers(t *testing.T) {
	body := "# H1\nLine 2\nLine 3\n## H2\nLine 5\n"
	sections := ParseSections(body)

	if len(sections) != 1 {
		t.Fatalf("expected 1 root, got %d", len(sections))
	}
	if sections[0].StartLine != 1 {
		t.Errorf("expected H1 start at line 1, got %d", sections[0].StartLine)
	}
	if len(sections[0].Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(sections[0].Children))
	}
	child := sections[0].Children[0]
	if child.StartLine != 4 {
		t.Errorf("expected H2 start at line 4, got %d", child.StartLine)
	}
}

func TestParseDocumentFromString_FullDocument(t *testing.T) {
	content := "---\nname: test-agent\ndescription: A test\n---\n# Overview\nSome overview\n## Details\nDetailed content\n"

	doc, err := ParseDocumentFromString(content, "agents/test-agent.md")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if doc.Path != "agents/test-agent.md" {
		t.Errorf("unexpected path: %q", doc.Path)
	}
	if doc.Frontmatter == nil {
		t.Fatal("expected frontmatter to be parsed")
	}
	if doc.Frontmatter.Name != "test-agent" {
		t.Errorf("expected name 'test-agent', got %q", doc.Frontmatter.Name)
	}
	if len(doc.Sections) != 1 {
		t.Fatalf("expected 1 root section, got %d", len(doc.Sections))
	}
	if doc.Sections[0].Title != "Overview" {
		t.Errorf("expected section title 'Overview', got %q", doc.Sections[0].Title)
	}
	if len(doc.Sections[0].Children) != 1 {
		t.Fatalf("expected 1 child section, got %d", len(doc.Sections[0].Children))
	}
	if doc.Sections[0].Children[0].Title != "Details" {
		t.Errorf("expected child title 'Details', got %q", doc.Sections[0].Children[0].Title)
	}
	if doc.RawContent != content {
		t.Error("expected RawContent to match original input")
	}
}

func TestParseDocumentFromString_NoFrontmatter(t *testing.T) {
	content := "# Just a doc\nWith some content\n"

	doc, err := ParseDocumentFromString(content, "rules/dev-workflow.md")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if doc.Frontmatter != nil {
		t.Error("expected nil frontmatter for document without ---")
	}
	if len(doc.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(doc.Sections))
	}
}

func TestParseDocumentFromString_EmptyFrontmatter(t *testing.T) {
	content := "---\n---\n# Body\n"

	doc, err := ParseDocumentFromString(content, "test.md")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if doc.Frontmatter == nil {
		t.Fatal("expected non-nil frontmatter for empty frontmatter block")
	}
	if doc.Frontmatter.Name != "" {
		t.Errorf("expected empty name, got %q", doc.Frontmatter.Name)
	}
}
