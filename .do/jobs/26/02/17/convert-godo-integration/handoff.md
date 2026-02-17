# Convert-Godo Integration: Comprehensive Handoff

**Date**: 2026-02-17
**Job**: .do/jobs/26/02/17/convert-godo-integration/
**Status**: Research complete → Architecture needed

---

## 0. Working Directories (Critical)

Three directories are involved. Each is a separate git repo.

```
~/Work/do-focus.workspace/do-focus/     ← do-focus (LEGACY, will be deprecated)
│  ├── cmd/godo/                        ← godo source code (41 files, 9,796 lines)
│  ├── .claude/rules/do/development/    ← current DO rules (dev-*.md)
│  ├── .claude/rules/do/workflow/       ← current workflow rules
│  ├── .claude/skills/                  ← current DO skills
│  └── CLAUDE.md                        ← current DO identity
│
~/Work/new/convert/                     ← convert project (FUTURE, absorbs godo)
│  ├── cmd/convert/main.go             ← CLI entry (to become cmd/godo/)
│  ├── internal/                        ← Go packages (assembler, extractor, etc.)
│  ├── core/                            ← extracted core templates (396 files)
│  ├── personas/                        ← extracted personas (do, do-ko, moai, moai-ko)
│  │   ├── do/                          ← DO persona (82 files, TO BE RESTRUCTURED)
│  │   ├── do-ko/                       ← Korean DO persona
│  │   ├── moai/                        ← MoAI persona (61 files)
│  │   └── moai-ko/                     ← Korean MoAI persona
│  └── .do/jobs/26/02/17/convert-godo-integration/  ← THIS JOB's documents
│
~/Work/moai-adk/                        ← MoAI ADK (ORIGINAL SOURCE for extraction)
   └── (423 core + 20 persona files)
```

### Key Relationships
- **godo source** lives in do-focus but will MOVE to convert
- **convert** reads from moai-adk (extract) and outputs assembled .claude/
- **do-focus rules** (dev-*.md) contain content that must be DECOMPOSED into convert's core/ skills
- **Job documents** (handoff, research, architecture, plan) live in convert's .do/jobs/
- **Commits go to BOTH repos**: rule changes → do-focus, restructuring → convert

### Which Repo for What
| Change Type | Target Repo | Example |
|-------------|-------------|---------|
| Rule additions/fixes | do-focus | workflow.md agent delegation rules |
| Core skill restructuring | convert (core/) | decomposed dev-*.md content |
| Persona workflow creation | convert (personas/do/) | 5 thin workflow files |
| godo feature migration | convert (internal/) | hook, mode, lint packages |
| Job documents | convert (.do/jobs/) | handoff.md, architecture.md |

### WARNING: do-focus ↔ convert Sync
Currently both repos have copies of the same files. When do-focus rules are modified, convert personas must also be updated. This sync problem GOES AWAY after restructuring — convert becomes the single source of truth.

---

## 1. Project Goal

Merge godo CLI (do-focus/cmd/godo/) into the convert project (~/Work/new/convert/).
Convert becomes standalone tool that does everything:
- Extract MoAI monolith → core + persona split
- Inject DO persona (with DO methodology, not MoAI's)
- Assemble final .claude/ output
- All godo features (hooks, mode, lint, create, profiles) embedded

do-focus becomes LEGACY after this. convert is the future.

---

## 2. DO Methodology Decisions (Critical)

### 2.1 NO SPEC Documents
MoAI uses SPEC documents (.moai/specs/SPEC-XXX/). DO does NOT.
DO uses a document chain instead:

```
analysis.md → architecture.md → plan.md → checklist.md → Run → report.md
```

All stored in: `.do/jobs/{YY}/{MM}/{DD}/{title}/`

### 2.2 TWO Phases Only: Plan + Run
- **Plan**: Research → Analysis → Architecture → Planning → Checklist decomposition
- **Run**: Checklist-based agent execution with VERIFY built into each agent's cycle

NO sync phase. MoAI's sync was a monolith (quality + docs + git delivery in one).
- Quality: already in Run (each agent's VERIFY step)
- Docs: checklist item if needed
- Git: manager-git agent call when needed

report.md is auto-generated when all checklist items are [o]. Not a separate phase.

### 2.3 Trigger Mapping (Natural Language → Execution)

| User says | Executes | Output |
|-----------|----------|--------|
| "파악해", "분석해", "조사해" | Analysis only | analysis.md |
| "설계해", "구조 짜줘" | Architecture only | architecture.md |
| "계획해", "플랜 짜줘" | Plan only | plan.md + checklist.md |
| "고민하고 계획해", "제대로 계획해" | Full chain | analysis → architecture → plan → checklist |
| "만들어줘", "구현해줘" | Auto complexity check | auto-select scope |
| Complexity >= threshold (5+ files, new module, 3+ domains) | Full chain auto | all documents |

### 2.4 Plan Phase Detail
```
Plan phase: Research → Analysis → Architecture → Plan → Checklist

Solo mode:  [research+analysis] → [architecture] → [plan+checklist]
Team mode:  [researcher]  ──┐
            [analyst]    ──┼→ synthesis → [plan.md + checklist.md]
            [architect]  ──┘
```

Research is a precursor to analysis (codebase exploration). In solo mode, one agent does both. In team mode, researcher runs in parallel with analyst and architect.

### 2.5 Run Phase Detail
Each agent follows the execution cycle:
```
READ → CLAIM [~] → WORK → VERIFY [*] → RECORD [o] → COMMIT
                     ↑         |
                     └─ fail ──┘  (retry)
```
Quality verification is INSIDE Run, not after it. This enables iterative improvement.

---

## 3. Brand Naming Convention

| Usage | Rule | Example |
|-------|------|---------|
| Folders/commands | lowercase `do` | `.do/`, `/do:plan`, `do-plan-{slug}` |
| Proper noun in docs | uppercase `DO` | "DO monitors:", "DO's Plan phase" |
| Environment variables | uppercase `DO_` | `DO_MODE`, `DO_LANGUAGE` |
| Slot template | `{{BRAND}}` = `DO` | "{{BRAND}} monitors:" → "DO monitors:" |

Reason: lowercase "do" reads as English verb. "DO" is unambiguous proper noun.
Compare: moai (folder), MoAI/Moai (proper noun).

---

## 4. Core vs Persona Classification

### 4.1 Core = Everything methodology-agnostic
- ALL expert-* agents (backend, frontend, security, debug, performance, testing, refactoring, devops)
- ALL team-* agents (researcher, analyst, architect, backend-dev, frontend-dev, tester, quality, designer)
- ALL manager-* agents (git, docs, strategy, ddd, tdd, quality, spec, project)
- ALL builder-* agents (agent, skill, plugin)
- ALL lang-*, domain-*, platform-*, library-*, tool-* skills
- ALL foundation-* skills (core, quality, context, claude, philosopher)
- ALL workflow-* skills (ddd, tdd, testing, spec, project, team, templates, etc.)
- Team orchestration MECHANICS (how to spawn, monitor, cleanup)

### 4.2 Persona = Brand-specific only
- CLAUDE.md (brand identity + DO/Focus/Team execution model)
- 5 workflow files (see section 5)
- Commands: do-check, do-checklist, do-mode, do-plan, do-setup, do-style
- Output styles: sprint, pair, direct
- Characters: young-f, young-m, senior-f, senior-m
- Spinners: Korean verb animations
- bootapp.md (DO-specific infrastructure rules)

### 4.3 Override Skills → ELIMINATED
Previously: persona had "override skills" (copies of core skills with modifications).
New approach: decompose rule content by topic into core skill sockets. No override layer.

---

## 5. Persona Workflows (Only 5)

```
persona/do/workflows/
├── plan.md          ← DO-style plan (research → analysis → architecture → plan → checklist)
├── run.md           ← DO-style run (checklist execution with VERIFY)
├── report.md        ← DO-style report (checklist aggregation → report.md)
├── team-plan.md     ← DO-style team plan (refs core plan-research pattern, outputs document chain)
└── team-run.md      ← DO-style team run (refs core implementation pattern, checklist-based)
```

These are THIN files. They reference core patterns and only declare what's different:
- Output format (document chain vs SPEC)
- Output location (.do/jobs/ vs .moai/specs/)
- Brand-specific team names (do-plan-{slug})

### Approach: Core Reference + Persona Override (Type A)
NOT full copies of MoAI workflows (Type B).
Core has the team mechanics. Persona declares "what's different".

---

## 6. dev-*.md Decomposition Mapping (KEY TASK - NOT YET DONE)

Current do-focus rules files contain MIXED concerns. They must be decomposed by TOPIC and placed into appropriate core skill sockets:

### dev-workflow.md (184 lines) → 7 destinations
| Content | Lines (approx) | Target Core Skill |
|---------|----------------|-------------------|
| Complexity check, Analysis/Architecture stages | ~40 | workflow-spec |
| TDD RED-GREEN-REFACTOR | ~10 | workflow-tdd |
| Agent execution cycle, delegation, resume, idempotent resume | ~50 | foundation-core |
| Agent research delegation rules (NEW, added a0f0f20) | ~7 | foundation-core or foundation-quality |
| Read before Write, coding/commit discipline, parallel agent isolation | ~40 | foundation-quality |
| Bug fix workflow | ~10 | workflow-testing |
| Knowledge management (reference order, documentation location) | ~15 | foundation-context |

### dev-testing.md (67 lines) → 1 destination
| Content | Lines | Target Core Skill |
|---------|-------|-------------------|
| Entire file (Real DB, FIRST, AI anti-patterns, data management, test quality) | 67 | workflow-testing |

### dev-checklist.md (523 lines) → 1 destination
| Content | Lines | Target Core Skill |
|---------|-------|-------------------|
| Checklist system, state management, templates, sub-checklist template, Analysis/Architecture templates, report template | 523 | workflow-spec |

### dev-environment.md (92 lines) → 3 destinations
| Content | Lines (approx) | Target Core Skill |
|---------|----------------|-------------------|
| Docker mandatory, 12-Factor, AI forbidden patterns | ~60 | workflow-project |
| bootapp domain rules | ~20 | **persona-only** (DO-specific) |
| Syntax check requirements | ~10 | foundation-quality |

### file-reading-optimization.md → 1 destination
| Content | Target |
|---------|--------|
| Progressive Loading, token budget awareness | foundation-context |

### coding-standards.md → stays as core rule
Already language-agnostic. Goes to core/rules/ as-is.

### agent-authoring.md → stays as core rule
Already generic. Goes to core/rules/ as-is.

### skill-authoring.md → stays as core rule
Already generic. Goes to core/rules/ as-is.

---

## 7. Core Team Patterns (to be created)

Extracted from MoAI's team-*.md workflows. Brand references removed.

### plan-research pattern
- Spawn: researcher + analyst + architect (parallel)
- TeamCreate, TaskCreate with dependencies
- Monitoring: forward researcher findings to architect
- Synthesis: collect all → delegate to persona workflow for output
- Cleanup: shutdown → TeamDelete → /clear

### implementation pattern
- Spawn: backend-dev + frontend-dev + tester (parallel)
- Variants: +designer (design_implementation), full-stack (api+ui+data+quality)
- File ownership boundaries per teammate
- Quality validation after implementation
- Cleanup

### Current MoAI team-*.md files
- team-plan.md (93 lines): 100% structurally identical between moai/do, only brand names differ
- team-run.md (180 lines): same
- team-sync.md (35 lines): MoAI-only concept, DO doesn't have sync
- team-review.md, team-debug.md: need verification

ALL reference SPEC documents which DO doesn't use. Must be rewritten for document chain.

---

## 8. godo → convert Integration

### godo source: do-focus/cmd/godo/ (41 files, 9,796 lines)

### Features to absorb:
| Feature | Files | Integration Plan |
|---------|-------|-----------------|
| Hook system (7 events) | hook*.go, moai_hook_types.go | Core: hook dispatcher + I/O contract |
| Security policy | hook_pre_tool.go, security_patterns.go | Core: reusable patterns |
| Mode system | mode.go | Core: execution + permission modes |
| moai sync (AST transform) | moai_sync*.go | Core: convert's extract/assemble replaces this |
| Persona loader | persona_loader.go | Persona: character/spinner loading |
| Claude profiles | claude_profile.go | Core: profile management |
| Lint | lint*.go | Core: language-specific linting |
| Create scaffolding | create.go | Core: agent/skill templates |
| Spinner | spinner.go | Persona: Korean verb animations |
| Statusline | statusline.go | Core: status rendering |
| Rank | rank*.go | Core: ranking system |
| GLM | glm.go | Core: GLM backend |

### convert's current structure:
```
convert/
├── cmd/convert/main.go        ← Entry point (to become godo)
├── internal/
│   ├── assembler/             ← Core + persona → .claude/
│   ├── cli/                   ← CLI commands
│   ├── detector/              ← Core vs persona detection
│   ├── extractor/             ← .claude/ → core + persona
│   ├── model/                 ← Data structures
│   ├── parser/                ← Markdown parsing
│   ├── template/              ← Slot system
│   └── validator/             ← Dependency validation
├── core/                      ← Extracted core (396 files)
├── personas/                  ← Extracted personas
│   ├── do/ (82 files)
│   ├── do-ko/ (583 files)
│   ├── moai/ (61 files)
│   └── moai-ko/ (564 files)
└── testdata/
```

### Target structure (after integration):
```
convert/ (renamed to godo?)
├── cmd/godo/main.go           ← Unified CLI entry
├── internal/
│   ├── assembler/             ← (existing)
│   ├── cli/                   ← Extended with godo commands
│   ├── detector/              ← (existing)
│   ├── extractor/             ← (existing)
│   ├── hook/                  ← NEW: hook system from godo
│   ├── lint/                  ← NEW: linting from godo
│   ├── mode/                  ← NEW: mode system from godo
│   ├── model/                 ← (existing, extended)
│   ├── parser/                ← (existing)
│   ├── persona/               ← NEW: persona loader from godo
│   ├── profile/               ← NEW: claude profiles from godo
│   ├── template/              ← (existing)
│   └── validator/             ← (existing)
├── core/                      ← Core templates (restructured)
│   ├── agents/
│   ├── rules/
│   ├── skills/
│   │   ├── workflow-testing/  ← + dev-testing.md content
│   │   ├── workflow-spec/     ← + dev-checklist.md content
│   │   ├── workflow-tdd/      ← + TDD rules
│   │   ├── workflow-project/  ← + Docker/12-Factor rules
│   │   ├── workflow-team/
│   │   │   └── patterns/     ← NEW: plan-research, implementation
│   │   ├── foundation-core/   ← + agent execution cycle
│   │   ├── foundation-quality/ ← + coding discipline
│   │   └── foundation-context/ ← + file reading, knowledge mgmt
│   └── team-patterns/         ← NEW: core team orchestration
├── personas/
│   ├── do/
│   │   ├── manifest.yaml
│   │   ├── CLAUDE.md
│   │   ├── workflows/        ← 5 files only (plan, run, report, team-plan, team-run)
│   │   ├── commands/
│   │   ├── characters/
│   │   ├── spinners/
│   │   ├── output-styles/
│   │   └── rules/bootapp.md  ← DO-specific only
│   ├── do-ko/
│   ├── moai/
│   └── moai-ko/
└── testdata/
```

---

## 9. Current State (What's Done)

### Research Phase [COMPLETE]
- [o] research-godo.md: godo CLI analysis (commit 812c885)
- [o] research-convert.md: convert project analysis (commit 812c885)
- [o] research-decisions.md: methodology decisions (commit 812c885)

### Rules Update
- [o] Agent research delegation rules added to do-focus workflow.md (commit a0f0f20)
- [!] MISTAKE: Synced old dev-workflow.md to convert personas/do/ (commit be68084) — this is WRONG because dev-*.md should be DECOMPOSED, not copied as-is. This commit should be reverted or the file should be rewritten later.

### Architecture Phase [NOT STARTED]
- [ ] architecture.md — the next step

### Plan Phase [NOT STARTED]
- [ ] plan.md
- [ ] checklist.md + checklists/

### Run Phase [NOT STARTED]
- [ ] Actual implementation

---

## 10. Mistakes Made & Lessons

### Mistake 1: Used Explore agents for file-producing research
- Explore agents have no Write tool → results accumulate in context → overflow
- Fix: Use general-purpose agents for research that needs file output
- Rule added to workflow.md (commit a0f0f20)

### Mistake 2: Asked agent to re-read files already in system prompt
- CLAUDE.md and rules files are already loaded → agent re-reads → double context
- Fix: Pass summaries to agents, don't ask them to read system prompt files

### Mistake 3: Synced old dev-workflow.md instead of rewriting
- We decided dev-*.md should be decomposed into core skill sockets
- Instead, I copied the old monolithic file to convert persona
- Fix: Revert or rewrite during architecture/implementation phase

---

## 11. Next Steps (For Next Context)

### Immediate: Architecture
1. Create architecture.md in .do/jobs/26/02/17/convert-godo-integration/
2. Design the full restructured project layout
3. Define interfaces between godo features and convert infrastructure
4. Map each dev-*.md section to its core skill destination with exact content

### Then: Plan + Checklist
1. Create plan.md with implementation order
2. Create checklist.md with agent assignments
3. Each checklist item = 1-3 files max

### Then: Run
1. Phase 1: Core skill content injection (decompose dev-*.md)
2. Phase 2: Persona workflow creation (5 thin files)
3. Phase 3: godo feature packages (hook, mode, lint, etc.)
4. Phase 4: CLI unification (cmd/godo/main.go)
5. Phase 5: Tests + validation

### Key Files to Read
- This handoff: .do/jobs/26/02/17/convert-godo-integration/handoff.md
- Research: research-godo.md, research-convert.md, research-decisions.md (same directory)
- Current convert code: ~/Work/new/convert/internal/
- Current godo code: ~/Work/do-focus.workspace/do-focus/cmd/godo/
- Current dev rules: ~/Work/do-focus.workspace/do-focus/.claude/rules/do/development/

---

## 12. Key Principles to Remember

1. **Core = HOW** (mechanics, patterns). **Persona = WHAT** (output format, brand identity)
2. **No SPEC documents** for DO. Document chain: analysis → architecture → plan → checklist
3. **Two phases** only: Plan + Run. No sync.
4. **VERIFY inside Run**, not after
5. **Thin persona workflows** that reference core patterns
6. **dev-*.md must be DECOMPOSED** by topic, not copied as-is
7. **DO (uppercase)** in docs, `do` (lowercase) in paths
8. **jobs directory** maintained: .do/jobs/{YY}/{MM}/{DD}/{title}/
9. **Override skills eliminated** — content goes to core sockets
10. **godo absorbed into convert** — single binary
