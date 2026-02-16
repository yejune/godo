# convert -- Persona Converter for Claude Code

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
4. Copy persona characters, spinners, and output styles
5. Apply skill name mappings across all agent files
6. Merge settings.json (core base + persona overrides)
7. Copy persona CLAUDE.md

## Persona Package Structure

A persona package lives under `personas/<name>/` and defines the complete identity, behavior, and tooling for a Claude Code persona. The Do persona (`personas/do/`) serves as the reference implementation.

```
personas/do/
├── manifest.yaml              # Declares all persona assets and configuration
├── CLAUDE.md                  # Top-level persona identity and execution directive
├── settings.json              # Hook definitions, output style, plans directory
├── agents/do/                 # Persona-only agent definitions (5 agents)
├── skills/do/                 # Orchestrator skill + workflows + reference
│   ├── SKILL.md               #   Intent/Mode Router, Execution Directive
│   ├── workflows/             #   plan, run, test, report, do, team-do
│   └── references/            #   Shared pattern reference
├── rules/                     # Development rules (operational backbone)
│   ├── dev-checklist.md       #   6-state checklist system
│   ├── dev-workflow.md        #   Complexity gates, commit discipline
│   ├── dev-testing.md         #   AI anti-patterns, Real DB only
│   ├── dev-environment.md     #   Docker-first, .env prohibition
│   ├── file-reading.md        #   Progressive file loading
│   └── do/workflow/           #   Persona-specific workflow rules
├── characters/                # 4 persona characters with YAML frontmatter
│   ├── young-f.md             #   Energetic 20s female developer (default)
│   ├── young-m.md             #   Confident 20s male developer
│   ├── senior-f.md            #   Legendary 50s female developer
│   └── senior-m.md            #   Senior 50s male architect
├── spinners/                  # Spinner verb definitions per character (YAML)
│   ├── young-f.yaml           #   Playful, emoji-rich spinner stems
│   ├── young-m.yaml           #   Confident spinner stems
│   ├── senior-f.yaml          #   Calm, professional spinner stems
│   └── senior-m.yaml          #   Authoritative spinner stems
├── styles/                    # 3 output styles
│   ├── sprint.md              #   Minimal talk, immediate execution
│   ├── pair.md                #   Collaborative pair programming (default)
│   └── direct.md              #   No-nonsense expert answers
├── commands/do/               # 6 persona slash commands
│   ├── checklist.md, mode.md, plan.md, setup.md, style.md, check.md
└── output-styles/do/          # Legacy output style location (being migrated)
    ├── sprint.md, pair.md, direct.md
```

### Characters

Each character file has YAML frontmatter defining structured metadata for programmatic access:

```yaml
---
id: young-f
name: "Do"
honorific_template: "{{name}}선배"
honorific_default: "선배"
tone: "반말+존댓말 혼합 (~할게요, ~했어요, ~해볼까요?)"
character_summary: "밝고 에너지 넘치는 20대 여성 천재 개발자."
relationship: "후배가 선배에게 캐주얼한 존중"
---
```

The body contains Identity, Personality, and Speech Pattern sections with mode-specific examples (Do/Focus/Team).

### Spinners

Spinner YAML files define the animated status messages shown during tool execution. Each file contains a list of stems (verb phrases) with optional emoji, plus a suffix pattern that cycles through variations:

```yaml
persona: young-f
suffix_pattern:
  cycle: 3
  suffixes: ["중이에요!", "중이에요 선배!", "중!"]
stems:
  - stem: "열심히 일하는"
    emoji: "🔥"
  - stem: "뚝딱뚝딱 정성껏 만드는"
    emoji: "🔨"
```

### Styles

Styles control **how** the persona communicates, independent of **who** the persona is. Any character can use any style, producing 4 x 3 = 12 behavioral combinations from just 7 definition files.

| Style | Behavior |
|-------|----------|
| `sprint` | Minimal talk, execute immediately, results only |
| `pair` (default) | Collaborative tone, joint decision-making |
| `direct` | No filler, expert answers only |

## DO_PERSONA Environment Variable

The `DO_PERSONA` environment variable (set in `.claude/settings.local.json`) selects which character is active. The `SessionStart` hook loads the chosen character file, and `PostToolUse` hooks reinforce the persona on every tool call.

| Value | Character | Korean Honorific | Relationship Dynamic |
|-------|-----------|-----------------|---------------------|
| `young-f` (default) | Bright, energetic 20s female developer | {name}선배 | Junior showing casual respect to senior |
| `young-m` | Confident 20s male developer | {name}선배님 | Junior showing formal respect to senior |
| `senior-f` | Legendary 50s female developer | {name}님 | Senior showing polite respect to colleague |
| `senior-m` | Industry-legend 50s male architect | {name}씨 | Senior showing warm authority to junior |

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

## Date Format Convention

All date-based directory paths in the Do persona use the `{YY}/{MM}/{DD}` format with forward-slash separators:

```
.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/plan.md
.do/jobs/26/02/16/feature-name/plan.md
```

This applies to jobs directories, plan files, checklist locations, and all date-stamped artifacts.

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
| **Characters** | All files under `characters/` | Persona character definitions |
| **Spinners** | All files under `spinners/` | Spinner verb definitions |
| **Styles** | All files under `styles/` | Output style definitions |
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
- Persona **inline text** (matched by content patterns) are replaced with inline slot markers
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

## Related Documents

| Document | Description |
|----------|-------------|
| `DO_PERSONA.md` | Do persona identity, philosophy, and design decisions |
| `RUNBOOK.md` | Operational guide for extraction, assembly, and version upgrades |
| `DO_MOAI_COMPARISON.md` | Detailed comparison between Do and MoAI philosophies |

## License

See repository root for license information.
