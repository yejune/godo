package template

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/do-focus/convert/internal/model"
)

var (
	// inlineSlotRe matches {{slot:SLOT_ID}} markers in content.
	inlineSlotRe = regexp.MustCompile(model.InlineSlotPattern)

	// sectionBeginRe matches <!-- BEGIN_SLOT:ID --> markers.
	sectionBeginRe = regexp.MustCompile(`<!-- BEGIN_SLOT:([A-Z][A-Z0-9_]*) -->`)

	// sectionEndRe matches <!-- END_SLOT:ID --> markers.
	sectionEndRe = regexp.MustCompile(`<!-- END_SLOT:([A-Z][A-Z0-9_]*) -->`)
)

// InsertSectionSlot wraps content with section slot markers.
// Returns: "<!-- BEGIN_SLOT:ID -->\ncontent\n<!-- END_SLOT:ID -->"
func InsertSectionSlot(slotID string, content string) string {
	begin := fmt.Sprintf(model.SectionSlotBegin, slotID)
	end := fmt.Sprintf(model.SectionSlotEnd, slotID)
	return begin + "\n" + content + "\n" + end
}

// InsertInlineSlot returns the inline slot marker: "{{slot:SLOT_ID}}"
func InsertInlineSlot(slotID string) string {
	return "{{slot:" + slotID + "}}"
}

// ExtractSectionSlot finds content between BEGIN_SLOT and END_SLOT markers
// for the given slotID. Returns the content between markers and whether it
// was found.
func ExtractSectionSlot(content string, slotID string) (string, bool) {
	begin := fmt.Sprintf(model.SectionSlotBegin, slotID)
	end := fmt.Sprintf(model.SectionSlotEnd, slotID)

	beginIdx := strings.Index(content, begin)
	if beginIdx == -1 {
		return "", false
	}
	afterBegin := beginIdx + len(begin)

	endIdx := strings.Index(content[afterBegin:], end)
	if endIdx == -1 {
		return "", false
	}

	extracted := content[afterBegin : afterBegin+endIdx]
	// Trim the leading and trailing newline that InsertSectionSlot adds.
	extracted = strings.TrimPrefix(extracted, "\n")
	extracted = strings.TrimSuffix(extracted, "\n")
	return extracted, true
}

// ReplaceInlineSlots replaces all {{slot:SLOT_ID}} markers with values from the map.
// Returns the modified content and list of slot IDs that were replaced.
func ReplaceInlineSlots(content string, values map[string]string) (string, []string) {
	replacedSet := map[string]bool{}
	result := inlineSlotRe.ReplaceAllStringFunc(content, func(match string) string {
		// Extract slot ID from {{slot:SLOT_ID}}
		sub := inlineSlotRe.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		slotID := sub[1]
		if val, ok := values[slotID]; ok {
			replacedSet[slotID] = true
			return val
		}
		return match
	})

	replaced := make([]string, 0, len(replacedSet))
	for id := range replacedSet {
		replaced = append(replaced, id)
	}
	sort.Strings(replaced)
	return result, replaced
}

// FindAllSlotMarkers finds all slot markers (both section and inline) in content.
// Returns a deduplicated, sorted list of slot IDs found.
func FindAllSlotMarkers(content string) []string {
	seen := map[string]bool{}

	// Find inline slots: {{slot:SLOT_ID}}
	for _, match := range inlineSlotRe.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			seen[match[1]] = true
		}
	}

	// Find section begin slots: <!-- BEGIN_SLOT:ID -->
	for _, match := range sectionBeginRe.FindAllStringSubmatch(content, -1) {
		if len(match) >= 2 {
			seen[match[1]] = true
		}
	}

	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
