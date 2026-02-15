# convert -- MoAI-ADK Persona Converter

A CLI tool that splits a `.claude/` directory into reusable **core templates** and a **persona layer**, then reassembles them with any persona to produce a complete, deployable `.claude/` directory.

## How It Works

The converter operates in two steps:

```
Extract (split)
  .claude/  ──>  core/              (methodology-agnostic templates)
                 personas/<name>/   (persona-specific content + manifest)

Assemble (merge)
  core/  +  persona manifest  ──>  .claude/   (complete, ready to use)
```

During **extraction**, each file is classified as core or persona based on pattern matching (headers, paths, skills, content patterns). Persona-specific sections within core files are replaced with slot markers, and the original content is recorded in the persona manifest.

During **assembly**, core templates are merged with a chosen persona manifest. Slot markers are filled with persona content, agent patches are applied, persona-only files are copied, settings are merged, and skill name mappings are resolved.

## Installation

```bash
go build -o convert ./cmd/convert
```

## Commands

### extract

Split a `.claude/` directory into core templates and a persona manifest.

**From a local directory:**

```bash
convert extract --src ~/Work/moai-adk/.claude --out ./output
```

**From a GitHub repository:**

```bash
convert extract --repo org/moai-adk --out ./output
convert extract --repo https://github.com/org/moai-adk.git --branch main --out ./output
```

**Flags:**

| Flag | Required | Description |
|------|----------|-------------|
| `--src` | One of src/repo | Path to a local `.claude/` directory |
| `--repo` | One of src/repo | GitHub repo URL or shorthand (e.g., `org/repo`) |
| `--branch` | No | Branch or tag to clone (default: repo default) |
| `--out` | Yes | Output directory for extracted layers |
| `--persona` | No | Persona name (default: auto-detected from source) |

`--src` and `--repo` are mutually exclusive.

**Output structure:**

```
output/
├── core/                          # Shared, methodology-agnostic files
│   ├── registry.yaml              # Slot registry (all discovered slots)
│   ├── agents/                    # Core agent definitions (with slot markers)
│   ├── rules/                     # Core rules
│   ├── skills/                    # Core skills
│   └── styles/                    # Core output styles
└── personas/moai/                 # Persona-specific content
    ├── manifest.yaml              # Persona manifest (assets, slots, patches)
    ├── agents/                    # Persona-only agents
    ├── skills/                    # Persona-only skills
    ├── rules/                     # Persona-only rules + development rules
    │   ├── dev-checklist.md       #   Checklist system, 6-state transitions, commit-as-proof
    │   ├── dev-workflow.md        #   Workflow orchestration, complexity gates, commit discipline
    │   ├── dev-testing.md         #   AI anti-patterns, Real DB only, FIRST principles
    │   ├── dev-environment.md     #   Docker-first, bootapp domains, .env prohibition
    │   └── file-reading.md        #   Progressive file loading, token budget awareness
    ├── commands/                  # Persona slash commands
    ├── hooks/                     # Persona hook scripts
    └── CLAUDE.md                  # Persona CLAUDE.md
```

### assemble

Merge core templates with a persona manifest to produce a complete `.claude/` directory.

```bash
convert assemble --core ./output/core --persona ./output/personas/moai/manifest.yaml --out ./result
```

**Flags:**

| Flag | Required | Description |
|------|----------|-------------|
| `--core` | Yes | Path to the core templates directory |
| `--persona` | Yes | Path to the persona manifest file (`manifest.yaml`) |
| `--out` | Yes | Output directory for the assembled `.claude/` |

**Assembly pipeline:**

1. Copy core files to output, filling slot markers with persona content
2. Apply agent patches (append/remove skills in frontmatter)
3. Copy persona-only files (agents, skills, rules, commands, hooks)
4. Apply skill name mappings across all agent files
5. Merge settings.json (core base + persona overrides)
6. Copy persona CLAUDE.md

## Development Rules

The persona package includes 5 development rules (`dev-*.md` and `file-reading.md`) that enforce Do's core development philosophy. These are not optional extras -- they are the operational backbone that makes the persona's workflow system function.

| Rule file | What it enforces |
|-----------|-----------------|
| `dev-checklist.md` | 6-state checklist transitions, sub-checklist templates, commit hash as completion proof |
| `dev-workflow.md` | Complexity-based workflow selection, Analysis/Architecture gates, Read Before Write, retry discipline |
| `dev-testing.md` | 7 AI anti-patterns (assertion weakening, test deletion, etc.), Real DB only, mutation testing mindset |
| `dev-environment.md` | Docker-first with bootapp domains, `.env` prohibition, container-only execution, 12-Factor principles |
| `file-reading.md` | 4-tier progressive file loading by size, Grep-first strategy, token budget awareness |

Without these rules, the persona's checklist system loses its state transition enforcement, the testing philosophy has no concrete prohibitions, and commit discipline becomes an unenforced suggestion. A persona package missing these files will produce a `.claude/` directory that declares a workflow but cannot enforce it.

## Classification Rules

Files are classified as **core** (shared template) or **persona** (methodology-specific) based on these rules:

### Persona markers

| Category | Rule | Examples |
|----------|------|---------|
| **Agents** | Names in the whole-file agent list | `manager-spec`, `manager-ddd`, `manager-tdd`, `manager-quality` |
| **Skills** | Names in the whole-file skill list | `moai-foundation-core`, `moai-workflow-ddd`, `moai-workflow-spec` |
| **Skill dirs** | Directory prefix match | `skills/moai/`, `skills/moai-workflow-project/` |
| **Rules** | Filenames in the whole-file rule list | `spec-workflow.md`, `workflow-modes.md` |
| **Commands** | All files under `commands/` | Persona slash commands |
| **Hooks** | All files under `hooks/` | Persona hook scripts |
| **Settings** | `settings.json` at `.claude/` root | Persona-specific config fields |
| **CLAUDE.md** | Project root `CLAUDE.md` | Top-level persona identity |
| **Headers** | Regex match on section titles | Sections titled "TRUST 5 Compliance", "TAG Chain" |
| **Content** | Inline text patterns in body | References to "TRUST 5 quality gates" |

### Core (shared templates)

- `rules/` files not in the whole-file rule list
- `agents/` files not in the whole-file agent list
- `skills/` files not matching persona skill patterns
- Shared utilities and generic configurations

### Partial classification

Some files contain both core and persona content. During extraction:
- Persona **sections** (matched by header patterns) are replaced with section slot markers
- Persona **inline text** (matched by content patterns) is replaced with inline slot markers
- Persona **skill references** in agent frontmatter are recorded as agent patches
- The rest of the file remains as a core template

### Slot syntax

Slots use a namespaced syntax to avoid collision with project template variables:

| Type | Syntax | Purpose |
|------|--------|---------|
| **Persona slot** | `{{slot:SLOT_ID}}` | Inline placeholder for persona content |
| **Section slot** | `<!-- BEGIN_SLOT:ID -->` ... `<!-- END_SLOT:ID -->` | Block placeholder for entire sections |
| **Project variable** | `{{VAR_NAME}}` | Plain `{{VAR}}` syntax, not processed by convert |

The `slot:` namespace prefix ensures persona slots are distinguishable from project-level template variables like `{{PRIMARY_USERS}}`.

## Full Roundtrip Example

```bash
# 1. Build the tool
go build -o convert ./cmd/convert

# 2. Extract: split .claude/ into core + persona
./convert extract --src ~/Work/moai-adk/.claude --out ./layers

# 3. Inspect the output
ls ./layers/core/           # shared templates + registry.yaml
ls ./layers/personas/moai/  # persona manifest + persona files
cat ./layers/personas/moai/manifest.yaml

# 4. Assemble: merge core + persona into a complete .claude/
./convert assemble \
  --core ./layers/core \
  --persona ./layers/personas/moai/manifest.yaml \
  --out ./assembled

# 5. Verify: the assembled output should match the original structure
ls ./assembled/             # complete .claude/ directory
diff -rq ~/Work/moai-adk/.claude ./assembled  # compare with original
```

## Project Structure

```
cmd/convert/           CLI entry point
internal/
  cli/                 Cobra command definitions (extract, assemble)
  detector/            Persona detection (pattern matching, classification)
  extractor/           Extract pipeline (agent, skill, rule, style, CLAUDE.md, settings, hooks, commands)
  assembler/           Assemble pipeline (merger, slot filler, orchestrator)
  parser/              Markdown parsing (frontmatter, sections)
  template/            Slot registry and slot operations
  model/               Shared types (Document, Section, PersonaManifest, Slot)
```

## License

See repository root for license information.
