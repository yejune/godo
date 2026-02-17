package persona

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Data holds parsed character data from a character .md file.
type Data struct {
	ID                string `yaml:"id"`
	Name              string `yaml:"name"`
	HonorificTemplate string `yaml:"honorific_template"`
	HonorificDefault  string `yaml:"honorific_default"`
	Tone              string `yaml:"tone"`
	CharacterSummary  string `yaml:"character_summary"`
	Relationship      string `yaml:"relationship"`
	FullContent       string `yaml:"-"` // markdown body (not from YAML)
}

// SpinnerStem holds a single spinner stem entry with optional emoji.
type SpinnerStem struct {
	Stem  string `yaml:"stem"`
	Emoji string `yaml:"emoji"`
}

// SpinnerSuffixPattern holds the suffix cycling configuration.
type SpinnerSuffixPattern struct {
	Cycle    int      `yaml:"cycle"`
	Suffixes []string `yaml:"suffixes"`
}

// SpinnerData holds parsed spinner configuration from a .yaml file.
type SpinnerData struct {
	Persona       string               `yaml:"persona"`
	SuffixPattern SpinnerSuffixPattern `yaml:"suffix_pattern"`
	Stems         []SpinnerStem        `yaml:"stems"`
}

// ResolveDir returns the first existing persona directory from the
// resolution order:
//  1. {projectDir}/.claude/personas/do/ (assembled output)
//  2. {projectDir}/personas/do/ (source, for development)
//
// Returns empty string if neither exists.
func ResolveDir() string {
	projectDir := os.Getenv("CLAUDE_PROJECT_DIR")
	if projectDir == "" {
		projectDir, _ = os.Getwd()
	}
	if projectDir == "" {
		return ""
	}

	candidates := []string{
		filepath.Join(projectDir, ".claude", "personas", "do"),
		filepath.Join(projectDir, "personas", "do"),
	}

	for _, dir := range candidates {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir
		}
	}
	return ""
}

// ParseFrontmatter extracts YAML frontmatter from markdown content.
// Frontmatter is delimited by --- on the first line and a closing ---.
// Returns the YAML bytes and the remaining markdown body.
func ParseFrontmatter(content []byte) (yamlBytes []byte, body string, err error) {
	s := string(content)
	// Must start with ---
	if !strings.HasPrefix(s, "---") {
		return nil, s, fmt.Errorf("no frontmatter found: file does not start with ---")
	}

	// Find the closing ---
	rest := s[3:]
	// Skip the newline after opening ---
	if idx := strings.IndexByte(rest, '\n'); idx >= 0 {
		rest = rest[idx+1:]
	}

	closeIdx := strings.Index(rest, "\n---")
	if closeIdx < 0 {
		return nil, s, fmt.Errorf("no closing --- found for frontmatter")
	}

	yamlPart := rest[:closeIdx]
	bodyPart := rest[closeIdx+4:] // skip \n---
	// Skip newline after closing ---
	if len(bodyPart) > 0 && bodyPart[0] == '\n' {
		bodyPart = bodyPart[1:]
	}

	return []byte(yamlPart), bodyPart, nil
}

// LoadCharacter reads and parses a character .md file from the given persona directory.
// personaDir is the base directory (e.g., "personas/do").
// personaType is the character ID (e.g., "young-f").
func LoadCharacter(personaDir, personaType string) (*Data, error) {
	path := filepath.Join(personaDir, "characters", personaType+".md")

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read character file %s: %w", path, err)
	}

	yamlBytes, body, err := ParseFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter in %s: %w", path, err)
	}

	var pd Data
	if err := yaml.Unmarshal(yamlBytes, &pd); err != nil {
		return nil, fmt.Errorf("failed to parse YAML in %s: %w", path, err)
	}

	pd.FullContent = strings.TrimSpace(body)

	return &pd, nil
}

// LoadSpinner reads and parses a spinner .yaml file from the given persona directory.
func LoadSpinner(personaDir, personaType string) (*SpinnerData, error) {
	path := filepath.Join(personaDir, "spinners", personaType+".yaml")

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read spinner file %s: %w", path, err)
	}

	var sd SpinnerData
	if err := yaml.Unmarshal(content, &sd); err != nil {
		return nil, fmt.Errorf("failed to parse spinner YAML %s: %w", path, err)
	}

	return &sd, nil
}

// BuildHonorific applies the user name to the honorific template.
// If userName is empty, returns HonorificDefault.
func (p *Data) BuildHonorific(userName string) string {
	if userName == "" {
		return p.HonorificDefault
	}
	return strings.ReplaceAll(p.HonorificTemplate, "{{name}}", userName)
}

// BuildReminder produces the short one-line reminder for PostToolUse/UserPromptSubmit.
// Format: 반드시 "{honorific}"로 호칭할 것. 말투: {tone}
func (p *Data) BuildReminder(userName string) string {
	honorific := p.BuildHonorific(userName)
	if honorific == "" {
		return ""
	}
	return "반드시 \"" + honorific + "\"로 호칭할 것. 말투: " + p.Tone
}

// BuildSpinnerVerbs generates the full spinner verb list from SpinnerData.
// Each stem is combined with a suffix from the cycling pattern, and emoji
// is appended after the suffix if present.
func (s *SpinnerData) BuildSpinnerVerbs() []string {
	cycle := s.SuffixPattern.Cycle
	if cycle < 1 {
		cycle = 1
	}
	suffixes := s.SuffixPattern.Suffixes
	if len(suffixes) == 0 {
		return nil
	}

	verbs := make([]string, len(s.Stems))
	for i, stem := range s.Stems {
		suffixIdx := i % cycle
		if suffixIdx >= len(suffixes) {
			suffixIdx = suffixIdx % len(suffixes)
		}
		suffix := suffixes[suffixIdx]

		if stem.Emoji != "" {
			verbs[i] = stem.Stem + suffix + " " + stem.Emoji
		} else {
			verbs[i] = stem.Stem + suffix
		}
	}

	return verbs
}
