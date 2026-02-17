package extractor

import (
	"path/filepath"
	"regexp"
	"strings"
)

// BrandSlotifier replaces brand-specific references in core file content
// with slot variables (e.g., {{slot:BRAND}}, {{slot:BRAND_DIR}}, {{slot:BRAND_CMD}}).
// This makes core files brand-neutral so they can be reassembled with any persona.
type BrandSlotifier struct {
	brand string

	// brandCatchAllRe matches any remaining case-insensitive occurrence of the brand name.
	// Runs after specific pattern replacements to catch all remaining references.
	brandCatchAllRe *regexp.Regexp
}

// NewBrandSlotifier creates a BrandSlotifier for the given brand name (e.g., "moai").
// Returns nil if brand is empty.
func NewBrandSlotifier(brand string) *BrandSlotifier {
	if brand == "" {
		return nil
	}
	// Case-insensitive catch-all for any remaining brand references.
	// Uses (?i) flag for case-insensitive matching.
	catchAll := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(brand))

	return &BrandSlotifier{
		brand:           brand,
		brandCatchAllRe: catchAll,
	}
}

// SlotifyContent replaces brand-specific references in content with slot variables.
//
// Replacement order (most specific first to prevent double-replacement):
//  1. /brand: → /{{slot:BRAND_CMD}}:  (slash command with colon, e.g., /moai:1-plan)
//  2. /brand<space> → /{{slot:BRAND_CMD}}<space>  (slash command, e.g., /moai plan)
//  3. .brand/ → .{{slot:BRAND_DIR}}/  (directory path, e.g., .moai/specs/)
//  4. brand- → {{slot:BRAND}}-  (skill name prefix, e.g., moai-lang-python)
func (s *BrandSlotifier) SlotifyContent(content string) string {
	if s == nil {
		return content
	}

	// Phase 1: Specific patterns (most-specific first, semantic slot types).

	// 1a. Slash commands: /moai:1-plan, /moai plan → BRAND_CMD
	content = strings.ReplaceAll(content, "/"+s.brand+":", "/{{slot:BRAND_CMD}}:")
	content = strings.ReplaceAll(content, "/"+s.brand+" ", "/{{slot:BRAND_CMD}} ")

	// 1b. Config directory: .moai/ → BRAND_DIR
	content = strings.ReplaceAll(content, "."+s.brand+"/", ".{{slot:BRAND_DIR}}/")

	// 1c. Path segment: /moai/ → BRAND_DIR (hooks/moai/, agents/moai/ etc.)
	content = strings.ReplaceAll(content, "/"+s.brand+"/", "/{{slot:BRAND_DIR}}/")

	// Phase 2: Case-insensitive catch-all for any remaining brand references → BRAND.
	// This catches: moai, MoAI, Moai, MOAI, moai-, moai_, --moai, @moai, etc.
	content = s.brandCatchAllRe.ReplaceAllString(content, "{{slot:BRAND}}")

	return content
}

// StripBrandPrefix removes the brand prefix from a skill directory name.
// For example, with brand "moai": "moai-lang-python" → "lang-python".
// Returns the original name unchanged if it doesn't start with the brand prefix.
func (s *BrandSlotifier) StripBrandPrefix(dirName string) string {
	if s == nil {
		return dirName
	}
	prefix := s.brand + "-"
	if strings.HasPrefix(dirName, prefix) {
		return dirName[len(prefix):]
	}
	return dirName
}

// RemapCorePath transforms a core file's relative path by stripping the brand
// prefix from skill directory names.
//
// Examples (brand="moai"):
//
//	skills/moai-lang-python/SKILL.md → skills/lang-python/SKILL.md
//	skills/moai-domain-backend/modules/api.md → skills/domain-backend/modules/api.md
//	agents/expert-backend.md → agents/expert-backend.md (unchanged, not a skill)
//	rules/dev-testing.md → rules/dev-testing.md (unchanged, not a skill)
func (s *BrandSlotifier) RemapCorePath(relPath string) string {
	if s == nil {
		return relPath
	}

	normalized := filepath.ToSlash(relPath)
	parts := strings.Split(normalized, "/")

	// Only remap skill paths: skills/<dirName>/...
	if len(parts) < 3 || parts[0] != "skills" {
		return relPath
	}

	stripped := s.StripBrandPrefix(parts[1])
	if stripped == parts[1] {
		return relPath // no change
	}

	parts[1] = stripped
	return strings.Join(parts, "/")
}

// brandSubdirCategories lists directory categories that use {category}/{brand}/ pattern.
// Override skill packages like skills/moai-foundation-core/ are NOT included here
// because the brand is in the package NAME, not as a subdirectory.
var brandSubdirCategories = map[string]bool{
	"agents":        true,
	"rules":         true,
	"commands":      true,
	"hooks":         true,
	"output-styles": true,
	"skills":        true,
}

// StripBrandSubdir removes the brand subdirectory from a persona file path.
// This is used during extraction to store persona files without brand nesting.
//
// Only applies to paths matching {category}/{brand}/... where category is one of
// agents, rules, commands, hooks, output-styles, or skills.
//
// For skills, only strips the brand if it's a plain subdirectory (skills/moai/foo.md),
// NOT if it's a brand-prefixed package name (skills/moai-foundation-core/SKILL.md).
//
// Examples (brand="moai"):
//
//	agents/moai/manager-ddd.md → agents/manager-ddd.md
//	rules/moai/workflow/spec.md → rules/workflow/spec.md
//	hooks/moai/pre-tool.sh → hooks/pre-tool.sh
//	commands/moai/plan.md → commands/plan.md
//	agents/expert-backend.md → agents/expert-backend.md (unchanged, no brand subdir)
//	skills/moai-foundation-core/SKILL.md → skills/moai-foundation-core/SKILL.md (unchanged, package name)
func (s *BrandSlotifier) StripBrandSubdir(relPath string) string {
	if s == nil {
		return relPath
	}

	normalized := filepath.ToSlash(relPath)
	parts := strings.Split(normalized, "/")

	// Need at least 3 parts: category/brand/file
	if len(parts) < 3 {
		return relPath
	}

	category := parts[0]
	if !brandSubdirCategories[category] {
		return relPath
	}

	// Check if the second segment is exactly the brand name (not brand-prefixed).
	if parts[1] != s.brand {
		return relPath
	}

	// Remove the brand segment: [category, brand, rest...] → [category, rest...]
	newParts := make([]string, 0, len(parts)-1)
	newParts = append(newParts, parts[0])
	newParts = append(newParts, parts[2:]...)
	return strings.Join(newParts, "/")
}
