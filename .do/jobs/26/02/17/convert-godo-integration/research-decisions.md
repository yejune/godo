# DO Framework Restructuring - Methodology Decisions

## 1. Core vs Persona Classification

### Rules → Core Skills Decomposition
Current DO rules files contain mixed concerns. Content must be decomposed by topic into appropriate core skill sockets:

| Source File | Content | Target Core Skill |
|-------------|---------|-------------------|
| dev-testing.md | Real DB, FIRST, AI anti-patterns, data management | workflow-testing |
| dev-workflow.md | Complexity判定, Analysis/Architecture stages | workflow-spec |
| dev-workflow.md | TDD RED-GREEN-REFACTOR | workflow-tdd |
| dev-workflow.md | Agent execution cycle, delegation, resume | foundation-core |
| dev-workflow.md | Read before Write, coding/commit discipline, parallel isolation | foundation-quality |
| dev-workflow.md | Bug fix workflow | workflow-testing |
| dev-workflow.md | Knowledge management | foundation-context |
| dev-checklist.md | Checklist system, templates, state management | workflow-spec |
| dev-environment.md | Docker, 12-Factor, AI forbidden patterns | workflow-project |
| dev-environment.md | bootapp domain rules | persona-only (DO specific) |
| dev-environment.md | Syntax check requirements | foundation-quality |
| file-reading.md | Progressive Loading, token budget | foundation-context |

### What stays in Persona (brand-specific only)
- CLAUDE.md (brand identity + DO/Focus/Team execution model)
- orchestrator SKILL.md + workflows (brand orchestration)
- commands (brand slash commands)
- output-styles, characters, spinners (brand aesthetics)
- bootapp.md (DO-specific infrastructure)

## 2. DO Methodology (differs from MoAI)

### Document Chain (NO SPEC)
```
analysis.md → architecture.md → plan.md → checklist.md → Run → report.md
```

### Two Phases Only: Plan and Run
- **Plan**: Research → Analysis → Architecture → Planning → Checklist decomposition
- **Run**: Checklist-based agent execution with VERIFY built into each agent cycle
- report.md auto-generated when all checklist items are [o]

### NO sync phase
- Quality verification: built into Run (each agent's VERIFY step)
- Document updates: checklist item if needed
- Git delivery: manager-git agent call when needed

### Trigger Mapping
| User says | Executes | Output |
|-----------|----------|--------|
| "파악해", "분석해" | Analysis only | analysis.md |
| "설계해", "구조 짜줘" | Architecture only | architecture.md |
| "계획해", "플랜 짜줘" | Plan only | plan.md + checklist.md |
| "고민하고 계획해", "제대로 계획해" | Full chain | analysis → architecture → plan → checklist |
| "만들어줘" | Auto complexity判定 | auto-select |
| Complexity >= threshold | Full chain auto | analysis → architecture → plan → checklist |

## 3. Brand Naming Convention

| Usage | Rule | Example |
|-------|------|---------|
| Folders/commands | lowercase `do` | `.do/`, `/do:plan`, `do-plan-{slug}` |
| Proper noun in docs | uppercase `DO` | "DO monitors:", "DO's Plan phase" |
| Environment variables | uppercase `DO_` | `DO_MODE`, `DO_LANGUAGE` |
| Slot template | `{{BRAND}}` = `DO` | `{{BRAND}} monitors:` → `DO monitors:` |

Reason: lowercase "do" in English sentences reads as verb "do". Uppercase "DO" is unambiguous proper noun.

## 4. Persona Workflows (5 only)

```
persona/do/workflows/
├── plan.md          ← DO-style plan (research → analysis → architecture → plan → checklist)
├── run.md           ← DO-style run (checklist execution with VERIFY)
├── report.md        ← DO-style report (checklist aggregation → report.md)
├── team-plan.md     ← DO-style team plan (references core plan-research pattern)
└── team-run.md      ← DO-style team run (references core implementation pattern)
```

Persona workflows are THIN - they reference core patterns and only declare what's different (output format, output location).

## 5. Core Team Patterns (extracted from MoAI monolith)

### plan-research pattern
- Spawn: researcher + analyst + architect (parallel)
- Mechanics: TeamCreate, TaskCreate, monitoring, synthesis, cleanup
- Persona injects: what documents to produce

### implementation pattern
- Spawn: backend-dev + frontend-dev + tester (parallel)
- Variants: +designer, full-stack
- Persona injects: checklist-based vs SPEC-based execution

### Key principle
- Core = HOW to organize teams (mechanics)
- Persona = WHAT to produce (output format)

## 6. Jobs Directory (maintained)

```
.do/jobs/{YY}/{MM}/{DD}/{title}/
├── analysis.md
├── architecture.md
├── plan.md
├── checklist.md
├── checklists/
│   ├── 01_agent-topic.md
│   └── 02_agent-topic.md
└── report.md
```

## 7. Integration Goal

- godo app absorbed into convert project
- convert becomes standalone tool
- do-focus becomes legacy
- convert = "godo moai sync" functionality:
  - Extract MoAI monolith → core + persona split
  - Inject DO persona overrides
  - Assemble final .claude output
  - All in one binary

## 8. Override Skills → Eliminated

Previous approach had "override skills" (persona copies of core skills with modifications).
New approach: decompose rules by topic into core skill sockets. No override layer needed.
Persona only contains brand-specific content that has no core equivalent.
