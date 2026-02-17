package parser

import (
	"strings"

	"github.com/yejune/godo/internal/model"
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

	// Sync DependsOn structured field to raw map
	if rawDep := dependsOnToRaw(fm.DependsOn); rawDep != nil {
		out["depends_on"] = rawDep
	} else {
		delete(out, "depends_on")
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


// dependsOnToRaw converts a DependsOn struct to a map[string]interface{} for YAML serialization.
// Returns nil if DependsOn is nil or entirely empty (all fields nil/empty).
func dependsOnToRaw(dep *model.DependsOn) map[string]interface{} {
	if dep == nil {
		return nil
	}

	out := make(map[string]interface{})

	if len(dep.Phases) > 0 {
		items := make([]interface{}, len(dep.Phases))
		for i, p := range dep.Phases {
			items[i] = p
		}
		out["phases"] = items
	}

	if len(dep.Artifacts) > 0 {
		items := make([]interface{}, len(dep.Artifacts))
		for i, a := range dep.Artifacts {
			items[i] = artifactDepToRaw(a)
		}
		out["artifacts"] = items
	}

	if len(dep.Agents) > 0 {
		items := make([]interface{}, len(dep.Agents))
		for i, a := range dep.Agents {
			items[i] = agentDepToRaw(a)
		}
		out["agents"] = items
	}

	if len(dep.Env) > 0 {
		items := make([]interface{}, len(dep.Env))
		for i, e := range dep.Env {
			items[i] = e
		}
		out["env"] = items
	}

	if len(dep.Services) > 0 {
		items := make([]interface{}, len(dep.Services))
		for i, s := range dep.Services {
			items[i] = serviceDepToRaw(s)
		}
		out["services"] = items
	}

	if len(dep.ChecklistItems) > 0 {
		items := make([]interface{}, len(dep.ChecklistItems))
		for i, c := range dep.ChecklistItems {
			items[i] = c
		}
		out["checklist_items"] = items
	}

	if len(out) == 0 {
		return nil
	}
	return out
}

// artifactDepToRaw converts an ArtifactDep to its YAML representation.
// Scalar shorthand: Required=true -> just the path string.
// Object form: Required=false -> map with path and required fields.
func artifactDepToRaw(a model.ArtifactDep) interface{} {
	if a.Required {
		return a.Path
	}
	return map[string]interface{}{
		"path":     a.Path,
		"required": a.Required,
	}
}

// agentDepToRaw converts an AgentDep to its YAML representation.
// Scalar shorthand: no Items -> just the name string.
// Object form: has Items -> map with name and items fields.
func agentDepToRaw(a model.AgentDep) interface{} {
	if len(a.Items) == 0 {
		return a.Name
	}
	items := make([]interface{}, len(a.Items))
	for i, item := range a.Items {
		items[i] = item
	}
	return map[string]interface{}{
		"name":  a.Name,
		"items": items,
	}
}

// serviceDepToRaw converts a ServiceDep to its YAML representation.
// Scalar shorthand: Healthcheck=false -> just the name string.
// Object form: Healthcheck=true -> map with name and healthcheck fields.
func serviceDepToRaw(s model.ServiceDep) interface{} {
	if !s.Healthcheck {
		return s.Name
	}
	return map[string]interface{}{
		"name":        s.Name,
		"healthcheck": s.Healthcheck,
	}
}

// PatchFrontmatterSkills modifies the skills field in raw YAML frontmatter text
// without re-serializing the entire document. This preserves original key order,
// formatting, and value styles (e.g., comma-separated inline vs YAML list).
//
// It detects the original skills format (inline comma-separated or YAML list)
// and writes the new skills in the same format. If skills were not present in the
// original YAML and newSkills is non-empty, a comma-separated skills line is appended.
// If newSkills is empty and skills existed, the skills field is removed.
func PatchFrontmatterSkills(rawYaml string, newSkills []string) string {
	lines := strings.Split(rawYaml, "\n")

	// Find the skills field
	skillsLineIdx := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "skills:") {
			skillsLineIdx = i
			break
		}
	}

	if skillsLineIdx < 0 {
		// No existing skills field
		if len(newSkills) == 0 {
			return rawYaml
		}
		// Append a comma-separated skills line at the end
		trimmed := strings.TrimRight(rawYaml, "\n")
		return trimmed + "\nskills: " + strings.Join(newSkills, ", ") + "\n"
	}

	// Determine format: inline (skills: a, b) or list (skills:\n  - a\n  - b)
	skillsLine := lines[skillsLineIdx]
	trimmedLine := strings.TrimSpace(skillsLine)

	// Check if it's inline format: "skills: value1, value2" (value after colon on same line)
	afterColon := strings.TrimPrefix(trimmedLine, "skills:")
	afterColon = strings.TrimSpace(afterColon)
	isInline := afterColon != ""

	if isInline {
		// Inline format: replace just this one line
		if len(newSkills) == 0 {
			// Remove the skills line
			lines = append(lines[:skillsLineIdx], lines[skillsLineIdx+1:]...)
		} else {
			// Preserve leading whitespace from original line
			leading := skillsLine[:len(skillsLine)-len(strings.TrimLeft(skillsLine, " \t"))]
			lines[skillsLineIdx] = leading + "skills: " + strings.Join(newSkills, ", ")
		}
		return strings.Join(lines, "\n")
	}

	// List format: "skills:\n    - a\n    - b"
	// Find the range of list items following the skills: line
	listStart := skillsLineIdx + 1
	listEnd := listStart
	for listEnd < len(lines) {
		trimmed := strings.TrimSpace(lines[listEnd])
		if strings.HasPrefix(trimmed, "- ") {
			listEnd++
		} else if trimmed == "" {
			// Skip blank lines within the list only if there's another list item after
			if listEnd+1 < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[listEnd+1]), "- ") {
				listEnd++
			} else {
				break
			}
		} else {
			break
		}
	}

	if len(newSkills) == 0 {
		// Remove skills: line and all list items
		lines = append(lines[:skillsLineIdx], lines[listEnd:]...)
		return strings.Join(lines, "\n")
	}

	// Detect indentation from first list item (if any)
	indent := "    " // default 4-space indent
	if listStart < listEnd {
		firstItem := lines[listStart]
		indent = firstItem[:len(firstItem)-len(strings.TrimLeft(firstItem, " \t"))]
	}

	// Build new list items
	var newListLines []string
	for _, skill := range newSkills {
		newListLines = append(newListLines, indent+"- "+skill)
	}

	// Replace: keep skills: line, replace list items
	result := make([]string, 0, len(lines)-listEnd+skillsLineIdx+1+len(newListLines))
	result = append(result, lines[:skillsLineIdx+1]...)
	result = append(result, newListLines...)
	result = append(result, lines[listEnd:]...)

	return strings.Join(result, "\n")
}
