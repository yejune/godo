package model

// Document represents a parsed markdown file with optional YAML frontmatter.
// It preserves the original structure while allowing section-level manipulation.
type Document struct {
	Path        string       `yaml:"-"`           // Original file path relative to .claude/
	Frontmatter *Frontmatter `yaml:"frontmatter"` // Parsed YAML frontmatter (nil if absent)
	Sections    []*Section   `yaml:"-"`           // Ordered list of header-bounded sections
	RawContent  string       `yaml:"-"`           // Original raw content for fallback
}

// Frontmatter represents the YAML frontmatter of an agent/skill file.
// Fields are stored both as structured data and raw map for round-trip fidelity.
type Frontmatter struct {
	Name           string                 `yaml:"name"`
	Description    string                 `yaml:"description"`
	Tools          string                 `yaml:"tools,omitempty"`
	Model          string                 `yaml:"model,omitempty"`
	PermissionMode string                 `yaml:"permissionMode,omitempty"`
	Skills         []string               `yaml:"skills,omitempty,flow"`
	Memory         string                 `yaml:"memory,omitempty"`
	Raw            map[string]interface{} `yaml:"-"` // Full raw data for round-trip
}
