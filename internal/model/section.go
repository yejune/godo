package model

// Section represents a header-bounded block within a markdown document.
// A section spans from its header line to the next header of equal or higher level.
type Section struct {
	Level     int        `yaml:"level"`      // Header level (1-6, 0 for preamble before first header)
	Title     string     `yaml:"title"`      // Header text without the '#' prefix
	Content   string     `yaml:"content"`    // Full content including header line
	StartLine int        `yaml:"start_line"` // 1-based line number in original file
	EndLine   int        `yaml:"end_line"`   // 1-based line number (exclusive)
	Children  []*Section `yaml:"children"`   // Nested subsections
}
