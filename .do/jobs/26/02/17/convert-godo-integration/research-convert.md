# convert Project Analysis

## Source Location
/Users/max/Work/new/convert/ (15,366 lines Go, 8 internal packages)

## Module
github.com/do-focus/convert (go 1.25.0)

## CLI Commands
```
convert
├── extract --src <path> --out <path> [--persona <name>]
└── assemble --core <path> --persona <path> --out <path>
```

## Internal Packages

### model/ (360 lines)
Data structures: Document, Frontmatter, Section, ClassificationResult, PersonaManifest, Slot, DependsOn

### parser/ (1,207 lines)
Markdown parsing: frontmatter extraction, section tree building, code block awareness, round-trip fidelity

### detector/ (678 lines)
Persona pattern detection: header patterns, path patterns, skill patterns, content patterns. Classifies sections as core vs persona.

### extractor/ (4,564 lines)
Main pipeline: walks source .claude/, classifies files, extracts core vs persona.
Sub-extractors: Agent, Skill, Rule, Style, Command, Hooks, Settings, Character, Spinner, ClaudeMD
BrandSlotifier: replaces brand refs with {{slot:BRAND}}, {{slot:BRAND_DIR}}, {{slot:BRAND_CMD}}

### template/ (508 lines)
Slot system: section markers (<!-- BEGIN_SLOT -->), inline markers ({{slot:ID}})
Registry: version-tracked slot definitions in registry.yaml

### assembler/
Merge pipeline: copy core + fill slots + apply agent patches + copy persona files + skill mappings + merge settings + copy CLAUDE.md
BrandDeslotifier: replaces {{slot:BRAND}} back to brand values

### validator/ (611 lines)
Dependency validation: phases, artifacts, agents, env, services, checklist_items
Graph: dependency DAG + cycle detection

### cli/ (501 lines)
CLI commands: extract, assemble with flag parsing

## Data Flow
```
Extract: source .claude/ → walk → classify → slotify → core/ + personas/manifest.yaml
Assemble: core/ + manifest.yaml → fill slots → deslotify → agent patches → output .claude/
```

## Personas (current)
- moai (61 files): 6 agents, 20+ skills, 7 hooks, 3 styles
- do (82 files): 5 agents, 8 skills, 4 characters, 4 spinners, 6 commands, 3 styles
- moai-ko (564 files): Korean moai
- do-ko (583 files): Korean do

## Current Gaps
1. No godo integration (standalone binary, no hooks/mode/lint)
2. No CLI validation (pre-flight checks)
3. No incremental extraction (always full)
4. No assembly verification (unresolved slots silent)
5. No persona diff tool
6. Slot content: no conditional/parameterized
7. Agent patches: can't modify frontmatter fields
