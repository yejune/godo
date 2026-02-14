package model

// PersonaManifest describes a persona's assets and how they integrate
// with core files. Loaded from persona_dir/persona.yaml.
type PersonaManifest struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`

	// Top-level persona definition
	ClaudeMD string `yaml:"claude_md"` // Path to persona's CLAUDE.md

	// Asset directories (relative to persona root)
	Agents   []string `yaml:"agents"`
	Skills   []string `yaml:"skills"`
	Rules    []string `yaml:"rules"`
	Styles   []string `yaml:"styles"`
	Commands []string `yaml:"commands"`

	// Slot content mappings: slot_id -> content file path
	SlotContent map[string]string `yaml:"slot_content"`

	// Frontmatter patches for core agents
	AgentPatches map[string]*AgentPatch `yaml:"agent_patches"`
}

// AgentPatch defines modifications to apply to a core agent file.
type AgentPatch struct {
	AppendSkills  []string `yaml:"append_skills"`
	RemoveSkills  []string `yaml:"remove_skills"`
	AppendContent string   `yaml:"append_content"` // Path to content to append
}
