package extractor

import (
	"github.com/do-focus/convert/internal/detector"
	"github.com/do-focus/convert/internal/model"
	"github.com/do-focus/convert/internal/template"
)

// AgentExtractor classifies and extracts persona-specific content from agent
// definition files (.claude/agents/*.md). It uses a PersonaDetector to classify
// sections and a PatternRegistry to identify whole-file persona agents.
type AgentExtractor struct {
	detector *detector.PersonaDetector
	registry *detector.PatternRegistry
}

// NewAgentExtractor creates an AgentExtractor with the given detector and registry.
func NewAgentExtractor(det *detector.PersonaDetector, reg *detector.PatternRegistry) *AgentExtractor {
	return &AgentExtractor{
		detector: det,
		registry: reg,
	}
}

// Extract separates a parsed agent Document into core and persona parts.
//
// For whole-file persona agents (listed in PatternRegistry.WholeFileAgents):
//   - Returns (nil, manifest, nil) where manifest.Agents contains the file path.
//
// For mixed agents with persona sections:
//   - Returns (coreDoc, manifest, nil) where persona sections are replaced with
//     slot markers and their content is stored in manifest.SlotContent.
//   - Persona skills are removed from core frontmatter and recorded in
//     manifest.AgentPatches.
//
// For pure core agents (no persona content):
//   - Returns (coreDoc, emptyManifest, nil) as a passthrough.
func (e *AgentExtractor) Extract(doc *model.Document) (*model.Document, *model.PersonaManifest, error) {
	manifest := &model.PersonaManifest{
		SlotContent:  make(map[string]string),
		AgentPatches: make(map[string]*model.AgentPatch),
	}

	// Check if this is a whole-file persona agent
	if doc.Frontmatter != nil && e.registry.IsWholeFilePersonaAgent(doc.Frontmatter.Name) {
		manifest.Agents = append(manifest.Agents, doc.Path)
		return nil, manifest, nil
	}

	// Classify all sections using the detector
	classification := e.detector.Classify(doc)

	// Build a lookup: section pointer -> classification
	sectionClass := make(map[*model.Section]*model.SectionClassification, len(classification.Sections))
	for i := range classification.Sections {
		sc := &classification.Sections[i]
		sectionClass[sc.Section] = sc
	}

	// Clone sections, replacing persona sections with slot markers
	coreSections := make([]*model.Section, len(doc.Sections))
	for i, sec := range doc.Sections {
		coreSections[i] = e.processSection(sec, sectionClass, manifest)
	}

	// Clone the document for core output
	coreDoc := &model.Document{
		Path:        doc.Path,
		Frontmatter: e.cloneFrontmatter(doc.Frontmatter),
		Sections:    coreSections,
		RawContent:  doc.RawContent,
	}

	// Extract persona skills from frontmatter
	if doc.Frontmatter != nil && len(classification.SkillRefs) > 0 {
		e.extractPersonaSkills(coreDoc, classification.SkillRefs, manifest)
	}

	return coreDoc, manifest, nil
}

// processSection handles a single section: if persona, replaces content with
// slot markers and stores original in manifest. Core sections are returned as-is
// (shallow copy).
func (e *AgentExtractor) processSection(
	sec *model.Section,
	classMap map[*model.Section]*model.SectionClassification,
	manifest *model.PersonaManifest,
) *model.Section {
	sc, ok := classMap[sec]

	// Clone the section (shallow)
	clone := &model.Section{
		Level:     sec.Level,
		Title:     sec.Title,
		Content:   sec.Content,
		StartLine: sec.StartLine,
		EndLine:   sec.EndLine,
	}

	if ok && sc.IsPersona && sc.SlotID != "" {
		// Store original persona content in manifest
		manifest.SlotContent[sc.SlotID] = sec.Content

		// Replace section content with slot marker
		clone.Content = template.InsertSectionSlot(sc.SlotID, "")
	}

	// Process children recursively
	if len(sec.Children) > 0 {
		clone.Children = make([]*model.Section, len(sec.Children))
		for i, child := range sec.Children {
			clone.Children[i] = e.processSection(child, classMap, manifest)
		}
	}

	return clone
}

// cloneFrontmatter creates a shallow copy of the frontmatter.
// Returns nil if input is nil.
func (e *AgentExtractor) cloneFrontmatter(fm *model.Frontmatter) *model.Frontmatter {
	if fm == nil {
		return nil
	}

	clone := &model.Frontmatter{
		Name:           fm.Name,
		Description:    fm.Description,
		Tools:          fm.Tools,
		Model:          fm.Model,
		PermissionMode: fm.PermissionMode,
		Memory:         fm.Memory,
	}

	// Deep copy skills slice
	if len(fm.Skills) > 0 {
		clone.Skills = make([]string, len(fm.Skills))
		copy(clone.Skills, fm.Skills)
	}

	// Copy raw map
	if fm.Raw != nil {
		clone.Raw = make(map[string]interface{}, len(fm.Raw))
		for k, v := range fm.Raw {
			clone.Raw[k] = v
		}
	}

	return clone
}

// extractPersonaSkills removes persona-specific skills from core frontmatter
// and records them in the manifest as an AgentPatch.
func (e *AgentExtractor) extractPersonaSkills(
	coreDoc *model.Document,
	skillRefs []model.SkillClassification,
	manifest *model.PersonaManifest,
) {
	if coreDoc.Frontmatter == nil {
		return
	}

	// Build a set of persona skill names
	personaSkills := make(map[string]bool, len(skillRefs))
	for _, sr := range skillRefs {
		personaSkills[sr.SkillName] = true
	}

	// Split skills into core and persona
	var coreSkills []string
	var extractedSkills []string
	for _, s := range coreDoc.Frontmatter.Skills {
		if personaSkills[s] {
			extractedSkills = append(extractedSkills, s)
		} else {
			coreSkills = append(coreSkills, s)
		}
	}

	// Update core frontmatter
	coreDoc.Frontmatter.Skills = coreSkills

	// Record persona skills in agent patch.
	// Use the document's relative file path as key (e.g., "agents/moai/expert-backend.md")
	// so it matches the assembler's PatchAgent expectation of file-path keys.
	if len(extractedSkills) > 0 {
		key := coreDoc.Path
		patch, ok := manifest.AgentPatches[key]
		if !ok {
			patch = &model.AgentPatch{}
			manifest.AgentPatches[key] = patch
		}
		patch.AppendSkills = append(patch.AppendSkills, extractedSkills...)
	}
}

