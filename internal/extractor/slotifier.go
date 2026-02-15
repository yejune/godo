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

	// skillNameRe matches the brand prefix in skill name contexts:
	// "moai-" followed by a lowercase letter (e.g., moai-lang-python, moai-domain-backend).
	skillNameRe *regexp.Regexp
}

// NewBrandSlotifier creates a BrandSlotifier for the given brand name (e.g., "moai").
// Returns nil if brand is empty.
func NewBrandSlotifier(brand string) *BrandSlotifier {
	if brand == "" {
		return nil
	}
	// Match brand prefix at word boundary followed by hyphen and lowercase letter.
	// This targets skill name patterns like "moai-lang-python" while avoiding
	// false positives in unrelated text.
	re := regexp.MustCompile(`\b` + regexp.QuoteMeta(brand) + `-(?=[a-z])`)

	return &BrandSlotifier{
		brand:       brand,
		skillNameRe: re,
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

	// 1. Slash commands with colon: /moai:1-plan → /{{slot:BRAND_CMD}}:1-plan
	content = strings.ReplaceAll(content, "/"+s.brand+":", "/{{slot:BRAND_CMD}}:")

	// 2. Slash commands with space: /moai plan → /{{slot:BRAND_CMD}} plan
	content = strings.ReplaceAll(content, "/"+s.brand+" ", "/{{slot:BRAND_CMD}} ")

	// 3. Directory paths: .moai/ → .{{slot:BRAND_DIR}}/
	content = strings.ReplaceAll(content, "."+s.brand+"/", ".{{slot:BRAND_DIR}}/")

	// 4. Skill name prefix: moai-lang-python → {{slot:BRAND}}-lang-python
	// Uses word-boundary regex to avoid matching inside unrelated words.
	content = s.skillNameRe.ReplaceAllString(content, "{{slot:BRAND}}-")

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
