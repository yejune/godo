package parser

import (
	"os"
	"regexp"
	"strings"

	"github.com/do-focus/convert/internal/model"
)

var headerRegex = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

// ParseSections splits markdown body text into a tree of Section structs.
// Headers inside fenced code blocks (``` or ~~~) are ignored.
// Content before the first header becomes a level-0 preamble section.
func ParseSections(body string) []*model.Section {
	if body == "" {
		return nil
	}

	lines := strings.Split(body, "\n")
	var flat []*model.Section
	inCodeBlock := false

	currentLevel := 0
	currentTitle := ""
	var currentLines []string
	currentStart := 1

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Track fenced code block state (``` or ~~~)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inCodeBlock = !inCodeBlock
			currentLines = append(currentLines, line)
			continue
		}

		if inCodeBlock {
			currentLines = append(currentLines, line)
			continue
		}

		// Check for ATX header
		match := headerRegex.FindStringSubmatch(line)
		if match != nil {
			// Flush the current section
			if len(currentLines) > 0 || currentLevel > 0 {
				flat = append(flat, &model.Section{
					Level:     currentLevel,
					Title:     currentTitle,
					Content:   joinLines(currentLines),
					StartLine: currentStart,
					EndLine:   lineNum,
				})
			}

			// Start new section
			currentLevel = len(match[1])
			currentTitle = strings.TrimSpace(match[2])
			currentLines = []string{line}
			currentStart = lineNum
			continue
		}

		currentLines = append(currentLines, line)
	}

	// Flush final section
	if len(currentLines) > 0 || currentLevel > 0 {
		flat = append(flat, &model.Section{
			Level:     currentLevel,
			Title:     currentTitle,
			Content:   joinLines(currentLines),
			StartLine: currentStart,
			EndLine:   len(lines) + 1,
		})
	}

	return nestSections(flat)
}

// nestSections builds a parent-child tree from a flat list of sections.
// Uses a stack-based approach: when a section with a higher level number (deeper nesting)
// is encountered, it becomes a child of the previous lower-level section.
// Level-0 sections (preamble) are always roots and never act as parents.
func nestSections(flat []*model.Section) []*model.Section {
	if len(flat) == 0 {
		return nil
	}

	var roots []*model.Section
	// Stack tracks the current nesting path. Each entry is a section that
	// could be a parent of subsequent deeper-level sections.
	// Level-0 (preamble) sections are excluded from the stack.
	var stack []*model.Section

	for _, sec := range flat {
		if sec.Level == 0 {
			// Preamble is always a root and never a parent
			roots = append(roots, sec)
			continue
		}

		// Pop stack until we find a parent (strictly lower level) or stack is empty
		for len(stack) > 0 && stack[len(stack)-1].Level >= sec.Level {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			roots = append(roots, sec)
		} else {
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, sec)
		}

		stack = append(stack, sec)
	}

	return roots
}

// ParseDocument reads a markdown file from disk, parses frontmatter and sections,
// and returns a Document.
func ParseDocument(path string) (*model.Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseDocumentFromString(string(data), path)
}

// ParseDocumentFromString parses markdown content into a Document.
// The path parameter is stored in the Document for reference but is not read from disk.
func ParseDocumentFromString(content string, path string) (*model.Document, error) {
	doc := &model.Document{
		Path:       path,
		RawContent: content,
	}

	yamlContent, body, hasFM := SplitFrontmatter(content)
	if hasFM && yamlContent != "" {
		fm, err := ParseFrontmatter(yamlContent)
		if err != nil {
			return nil, err
		}
		doc.Frontmatter = fm
	} else if hasFM {
		// Empty frontmatter -- still valid, just no fields
		doc.Frontmatter = &model.Frontmatter{
			Raw: make(map[string]interface{}),
		}
	}

	doc.Sections = ParseSections(body)

	return doc, nil
}

// joinLines joins lines with newline separator.
func joinLines(lines []string) string {
	return strings.Join(lines, "\n")
}
