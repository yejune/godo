package parser

import (
	"strings"

	"github.com/do-focus/convert/internal/model"
	"gopkg.in/yaml.v3"
)

const frontmatterDelimiter = "---"

// SplitFrontmatter splits raw markdown content into frontmatter YAML and body.
// Returns (yamlContent, body, hasFrontmatter).
// Frontmatter must start at the very first line with "---".
func SplitFrontmatter(content string) (string, string, bool) {
	lines := strings.SplitAfter(content, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != frontmatterDelimiter {
		return "", content, false
	}

	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == frontmatterDelimiter {
			yamlContent := strings.Join(lines[1:i], "")
			// Body starts after the closing delimiter line
			body := strings.Join(lines[i+1:], "")
			return yamlContent, body, true
		}
	}

	// No closing delimiter found -- treat entire content as body (no frontmatter)
	return "", content, false
}

// ParseFrontmatter parses YAML string into a Frontmatter struct.
// It populates both structured fields and the Raw map for round-trip fidelity.
// Handles skills as both []string and comma-separated string.
func ParseFrontmatter(yamlContent string) (*model.Frontmatter, error) {
	fm := &model.Frontmatter{}

	// Parse into Raw map first for round-trip preservation
	raw := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(yamlContent), &raw); err != nil {
		return nil, err
	}
	fm.Raw = raw

	// Check if skills is a plain string before struct unmarshal.
	// yaml.Unmarshal into []string fails when the YAML value is a scalar string,
	// so we detect this case and handle it manually.
	var skillsFromString []string
	if rawSkills, ok := raw["skills"]; ok {
		if skillStr, ok := rawSkills.(string); ok {
			skillsFromString = splitAndTrim(skillStr, ",")
			// Remove skills from raw before struct unmarshal to avoid type error
			delete(raw, "skills")
			// Re-marshal without the problematic field
			cleaned, err := yaml.Marshal(raw)
			if err != nil {
				return nil, err
			}
			if err := yaml.Unmarshal(cleaned, fm); err != nil {
				return nil, err
			}
			// Restore skills in both places
			fm.Skills = skillsFromString
			raw["skills"] = rawSkills // restore original raw value
			fm.Raw = raw
			return fm, nil
		}
	}

	// Normal case: skills is already a []string (or absent)
	if err := yaml.Unmarshal([]byte(yamlContent), fm); err != nil {
		return nil, err
	}

	return fm, nil
}

// SerializeFrontmatter converts a Frontmatter back to YAML string with delimiters.
// It uses the Raw map as the base and overlays structured field changes,
// preserving unknown fields that were in the original YAML.
func SerializeFrontmatter(fm *model.Frontmatter) (string, error) {
	if fm == nil {
		return "", nil
	}

	// Start from Raw map to preserve unknown fields
	out := make(map[string]interface{})
	for k, v := range fm.Raw {
		out[k] = v
	}

	// Overlay structured fields (these may have been modified)
	if fm.Name != "" {
		out["name"] = fm.Name
	}
	if fm.Description != "" {
		out["description"] = fm.Description
	}
	if fm.Tools != "" {
		out["tools"] = fm.Tools
	} else {
		delete(out, "tools")
	}
	if fm.Model != "" {
		out["model"] = fm.Model
	} else {
		delete(out, "model")
	}
	if fm.PermissionMode != "" {
		out["permissionMode"] = fm.PermissionMode
	} else {
		delete(out, "permissionMode")
	}
	if len(fm.Skills) > 0 {
		skills := make([]interface{}, len(fm.Skills))
		for i, s := range fm.Skills {
			skills[i] = s
		}
		out["skills"] = skills
	} else {
		delete(out, "skills")
	}
	if fm.Memory != "" {
		out["memory"] = fm.Memory
	} else {
		delete(out, "memory")
	}

	data, err := yaml.Marshal(out)
	if err != nil {
		return "", err
	}

	return frontmatterDelimiter + "\n" + string(data) + frontmatterDelimiter + "\n", nil
}

// splitAndTrim splits a string by separator and trims whitespace from each part.
// Empty parts after trimming are excluded.
func splitAndTrim(s string, sep string) []string {
	parts := strings.Split(s, sep)
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
