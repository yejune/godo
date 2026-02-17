package assembler

import (
	"fmt"
	"os"
	"strings"

	"github.com/do-focus/convert/internal/model"
)

// BrandDeslotifier reverses the slotification performed by extractor.BrandSlotifier.
// It replaces brand slot variables ({{slot:BRAND}}, {{slot:BRAND_DIR}}, {{slot:BRAND_CMD}})
// with actual persona values and prepends brand prefix to stripped skill directory names.
type BrandDeslotifier struct {
	brand    string
	brandDir string
	brandCmd string
}

// NewBrandDeslotifier creates a BrandDeslotifier from the manifest's brand fields.
// If Brand is empty but Name is set, brand is inferred from the manifest name
// and default conventions are applied (BrandDir = name, BrandCmd = name).
// The leading "." and "/" are literal in slot patterns, not part of the value.
// Returns nil if both Brand and Name are empty.
func NewBrandDeslotifier(manifest *model.PersonaManifest) *BrandDeslotifier {
	if manifest == nil {
		return nil
	}

	brand := manifest.Brand
	if brand == "" {
		brand = manifest.Name
	}
	if brand == "" {
		return nil
	}

	brandDir := manifest.BrandDir
	if brandDir == "" {
		brandDir = brand
	}

	brandCmd := manifest.BrandCmd
	if brandCmd == "" {
		brandCmd = brand
	}

	return &BrandDeslotifier{
		brand:    brand,
		brandDir: brandDir,
		brandCmd: brandCmd,
	}
}

// DeslotifyContent replaces brand slot variables in content with actual values.
// Only replaces slots whose corresponding brand field is non-empty.
// Other slot variables (e.g., {{slot:TOOL_NAME}}) are preserved.
func (d *BrandDeslotifier) DeslotifyContent(content string) string {
	if d == nil {
		return content
	}

	if d.brand != "" {
		content = strings.ReplaceAll(content, "{{slot:BRAND}}", d.brand)
	}
	if d.brandDir != "" {
		content = strings.ReplaceAll(content, "{{slot:BRAND_DIR}}", d.brandDir)
	}
	if d.brandCmd != "" {
		content = strings.ReplaceAll(content, "{{slot:BRAND_CMD}}", d.brandCmd)
	}

	return content
}

// DeslotifyFile reads a file, replaces brand slot variables, and writes back.
// Skips writing if no replacements were made.
func (d *BrandDeslotifier) DeslotifyFile(path string) error {
	if d == nil {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %s: %w", path, err)
	}

	original := string(data)
	filled := d.DeslotifyContent(original)

	if filled == original {
		return nil
	}

	if err := os.WriteFile(path, []byte(filled), 0o644); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}

	return nil
}

// RemapSkillPath prepends the brand prefix to skill directory names,
// reversing the stripping done by extractor.BrandSlotifier.StripBrandPrefix.
//
// Examples (brand="moai"):
//
//	skills/lang-python/SKILL.md → skills/moai-lang-python/SKILL.md
//	skills/domain-backend/modules/api.md → skills/moai-domain-backend/modules/api.md
//	agents/expert-backend.md → agents/expert-backend.md (unchanged)
func (d *BrandDeslotifier) RemapSkillPath(relPath string) string {
	if d == nil {
		return relPath
	}

	parts := strings.Split(relPath, "/")

	// Only remap skill paths: skills/<dirName>/...
	if len(parts) < 3 || parts[0] != "skills" {
		return relPath
	}

	dirName := parts[1]
	if dirName == "" {
		return relPath
	}

	// Don't double-prefix if already prefixed.
	prefix := d.brand + "-"
	if strings.HasPrefix(dirName, prefix) {
		return relPath
	}

	parts[1] = prefix + dirName
	return strings.Join(parts, "/")
}

// RemapBrandDirInPath replaces the source brand directory segment with the target brand
// directory in file paths. This handles core files under agents/moai/, rules/moai/ etc.
// being remapped to agents/do/, rules/do/ when assembling for a different persona.
//
// Example (sourceBrand="moai", target brandDir="do"):
//   agents/moai/expert-backend.md → agents/do/expert-backend.md
//   rules/moai/workflow/spec.md → rules/do/workflow/spec.md
func (d *BrandDeslotifier) RemapBrandDirInPath(relPath, sourceBrandDir string) string {
	if d == nil || sourceBrandDir == "" || sourceBrandDir == d.brandDir {
		return relPath
	}
	return strings.Replace(relPath, "/"+sourceBrandDir+"/", "/"+d.brandDir+"/", 1)
}

// brandSubdirCategories lists directory categories that use {category}/{brand}/ pattern.
var brandSubdirCategories = map[string]bool{
	"agents":        true,
	"rules":         true,
	"commands":      true,
	"hooks":         true,
	"output-styles": true,
	"skills":        true,
}

// AddBrandSubdir inserts the brand subdirectory into a persona file path.
// This is used during assembly to restore brand nesting that was stripped during extraction.
//
// Only applies to paths whose top-level directory is one of the brand subdir categories.
// Does NOT add brand subdir if one is already present (prevents double-nesting).
//
// Examples (brand="moai"):
//
//	agents/manager-ddd.md → agents/moai/manager-ddd.md
//	rules/workflow/spec.md → rules/moai/workflow/spec.md
//	hooks/pre-tool.sh → hooks/moai/pre-tool.sh
//	commands/plan.md → commands/moai/plan.md
//	agents/moai/manager-ddd.md → agents/moai/manager-ddd.md (unchanged, already has brand)
//	styles/pair.md → styles/pair.md (unchanged, not in brand categories)
//	CLAUDE.md → CLAUDE.md (unchanged, no category dir)
func (d *BrandDeslotifier) AddBrandSubdir(relPath string) string {
	if d == nil {
		return relPath
	}

	parts := strings.Split(relPath, "/")

	// Need at least 2 parts: category/file
	if len(parts) < 2 {
		return relPath
	}

	category := parts[0]
	if !brandSubdirCategories[category] {
		return relPath
	}

	// Don't double-add if brand subdir is already present.
	if len(parts) >= 2 && parts[1] == d.brand {
		return relPath
	}

	// Insert brand after category: [category, rest...] → [category, brand, rest...]
	newParts := make([]string, 0, len(parts)+1)
	newParts = append(newParts, parts[0])
	newParts = append(newParts, d.brand)
	newParts = append(newParts, parts[1:]...)
	return strings.Join(newParts, "/")
}
