package detector

import (
	"regexp"
	"strings"

	"github.com/do-focus/convert/internal/model"
)

// compiledHeaderPattern pairs a HeaderPattern with its pre-compiled regexp.
type compiledHeaderPattern struct {
	pattern HeaderPattern
	re      *regexp.Regexp
}

// PersonaDetector identifies persona-specific content within parsed documents.
// It uses a PatternRegistry to match headers, paths, and skills that indicate
// methodology-coupled content (TRUST5, SPEC workflow, etc.).
type PersonaDetector struct {
	registry *PatternRegistry
	compiled []*compiledHeaderPattern
}

// NewPersonaDetector creates a PersonaDetector with pre-compiled header regexps.
// Returns an error if any header pattern regex fails to compile.
func NewPersonaDetector(registry *PatternRegistry) (*PersonaDetector, error) {
	compiled := make([]*compiledHeaderPattern, len(registry.HeaderPatterns))
	for i, hp := range registry.HeaderPatterns {
		re, err := regexp.Compile(hp.Pattern)
		if err != nil {
			return nil, err
		}
		compiled[i] = &compiledHeaderPattern{
			pattern: hp,
			re:      re,
		}
	}
	return &PersonaDetector{
		registry: registry,
		compiled: compiled,
	}, nil
}

// Classify determines whether each section in a document is Core or Persona.
// It walks all sections (including children recursively), matches headers against
// compiled patterns, detects skill refs in frontmatter, and finds path patterns
// in section content.
func (d *PersonaDetector) Classify(doc *model.Document) *model.ClassificationResult {
	result := &model.ClassificationResult{
		DocPath: doc.Path,
	}

	// Walk all sections recursively and classify each
	d.walkSections(doc.Sections, result)

	// Detect persona skills in frontmatter
	result.SkillRefs = d.DetectSkillRefs(doc.Frontmatter)

	// Detect path patterns across all section content
	for _, sec := range doc.Sections {
		d.collectPathMatches(sec, result)
	}

	return result
}

// walkSections recursively walks sections and classifies each one.
func (d *PersonaDetector) walkSections(sections []*model.Section, result *model.ClassificationResult) {
	for _, sec := range sections {
		sc := d.classifySection(sec)
		result.Sections = append(result.Sections, sc)

		// Recurse into children
		if len(sec.Children) > 0 {
			d.walkSections(sec.Children, result)
		}
	}
}

// classifySection matches a single section's title against header patterns.
func (d *PersonaDetector) classifySection(sec *model.Section) model.SectionClassification {
	for _, chp := range d.compiled {
		if chp.re.MatchString(sec.Title) {
			return model.SectionClassification{
				Section:    sec,
				IsPersona:  true,
				Reason:     "header matches pattern: " + chp.pattern.Description,
				SlotID:     chp.pattern.SlotID,
				Confidence: 1.0,
			}
		}
	}

	// No match -- classify as core
	return model.SectionClassification{
		Section:    sec,
		IsPersona:  false,
		Reason:     "no persona pattern matched",
		Confidence: 0.0,
	}
}

// collectPathMatches finds path patterns in a section's content and its children.
func (d *PersonaDetector) collectPathMatches(sec *model.Section, result *model.ClassificationResult) {
	matches := d.DetectPathPatterns(sec.Content)
	result.PathRefs = append(result.PathRefs, matches...)

	for _, child := range sec.Children {
		d.collectPathMatches(child, result)
	}
}

// DetectSkillRefs identifies persona-specific skills in frontmatter.
// Returns a list of SkillClassification for each persona skill found.
func (d *PersonaDetector) DetectSkillRefs(fm *model.Frontmatter) []model.SkillClassification {
	if fm == nil {
		return nil
	}

	var refs []model.SkillClassification
	for _, skill := range fm.Skills {
		for _, sp := range d.registry.SkillPatterns {
			if skill == sp.SkillName {
				refs = append(refs, model.SkillClassification{
					SkillName: sp.SkillName,
					Category:  sp.Category,
				})
				break
			}
		}
	}
	return refs
}

// DetectPathPatterns finds hardcoded persona paths in content text.
// Returns matches with their line and column numbers (1-based).
func (d *PersonaDetector) DetectPathPatterns(content string) []model.PathMatch {
	var matches []model.PathMatch

	lines := strings.Split(content, "\n")
	for _, pp := range d.registry.PathPatterns {
		re, err := regexp.Compile(pp.Pattern)
		if err != nil {
			continue
		}

		for lineIdx, line := range lines {
			locs := re.FindAllStringIndex(line, -1)
			for _, loc := range locs {
				matches = append(matches, model.PathMatch{
					Original: line[loc[0]:loc[1]],
					SlotID:   pp.SlotID,
					Line:     lineIdx + 1, // 1-based
					Column:   loc[0] + 1,  // 1-based
				})
			}
		}
	}

	// Sort by line then column for deterministic output
	sortPathMatches(matches)

	return matches
}

// sortPathMatches sorts path matches by line number, then column.
func sortPathMatches(matches []model.PathMatch) {
	for i := 1; i < len(matches); i++ {
		for j := i; j > 0; j-- {
			if matches[j].Line < matches[j-1].Line ||
				(matches[j].Line == matches[j-1].Line && matches[j].Column < matches[j-1].Column) {
				matches[j], matches[j-1] = matches[j-1], matches[j]
			} else {
				break
			}
		}
	}
}
