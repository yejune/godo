---
name: do-reference
description: >
  Common execution patterns, mode reference, configuration file paths,
  environment variables, artifact locations, persona system, and
  Do-specific features used across all Do workflows.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "foundation"
  status: "active"
  updated: "2026-02-15"
  tags: "reference, patterns, configuration, environment, persona"

# Do Extension: Triggers
triggers:
  keywords: ["reference", "pattern", "config", "environment", "persona", "mode"]
  agents: ["manager-ddd", "manager-tdd", "manager-quality", "manager-git"]
  phases: ["plan", "run", "test", "report"]
---

# Do Skill Reference

Common patterns, mode reference, configuration, and Do-specific features used across all workflows.

---

## Execution Patterns

### Parallel Execution

When multiple operations are independent, invoke them in a single response. Claude Code runs multiple Task() calls in parallel (up to 10 concurrent).

Use Cases:
- Exploration Phase: Launch codebase analysis + documentation research simultaneously
- Multi-agent dispatch: Send independent sub-checklist agents in parallel
- Diagnostic Scan: Run tests + lint + type check simultaneously

### Sequential Execution

When operations have dependencies, chain sequentially. Each Task() receives context from previous phases.

Use Cases:
- Plan pipeline: Analysis -> Architecture -> Plan (each feeds the next)
- Checklist dependencies: Task B depends on Task A completing first
- Quality gates: Implementation must pass before report generation

### Resume Pattern (Checklist-Based)

When a workflow or agent is interrupted:
- Read checklist file to find last `[o]` completed item
- Skip all `[o]` items (already done with committed proof)
- Resume from first `[ ]` or `[~]` item
- New agent receives same sub-checklist path for continuity
- This pattern ensures ANY agent can pick up where another left off

### Context Propagation

Each phase passes results forward:
- Include previous phase outputs in Task() prompts
- Reference specific file paths rather than inlining large content
- Checklist files serve as the canonical state across phases

---

## Mode Reference

### Do Mode (DO_MODE=do)

- Response prefix: `[Do]`
- Execution: ALL implementation via Task() delegation
- Parallelism: Independent tasks run in parallel
- File access: Only through delegated agents
- Best for: 5-10 files, complex multi-domain tasks

### Focus Mode (DO_MODE=focus)

- Response prefix: `[Focus]`
- Execution: Direct Read/Write/Edit by orchestrator
- Parallelism: Sequential (one task at a time)
- File access: Direct orchestrator access
- Best for: 1-3 files, simple bugs, CSS changes, small refactors

### Team Mode (DO_MODE=team)

- Response prefix: `[Team]`
- Execution: Agent Teams API (TeamCreate/SendMessage)
- Parallelism: Teammates work simultaneously on independent tasks
- File access: Each teammate owns specific files exclusively
- Best for: 10+ files, 3+ domains, large-scale parallel work

### Do vs Team Comparison

| Aspect | Do Mode | Team Mode |
|--------|---------|-----------|
| Execution | Task(subagent) sequential/parallel | Agent Teams simultaneous |
| File Ownership | None (conflict possible) | Per-teammate exclusive |
| State Management | Checklist | Shared task list + Checklist |
| Suitable Scale | 5-10 files | 10+ files |
| Fallback | Focus | Do |

---

## Auto-Escalation Thresholds

| Condition | Trigger | Action |
|-----------|---------|--------|
| Focus -> Do | 5+ files, multi-domain, expert needed, 30K+ tokens | Suggest mode switch |
| Do -> Team | 10+ files, 3+ domains, parallel research beneficial | Suggest mode switch |
| Team unavailable | AGENT_TEAMS env not set | Fallback to Do mode |

---

## Configuration Reference

### settings.json (Project-Shared, Git-Committed)

Standard Claude Code settings fields:
- `outputStyle`: Current output style (sprint/pair/direct)
- `plansDirectory`: Plan file location (`.do/jobs`)
- `statusLine`: Statusline command (`godo statusline`)
- `hooks`: Hook definitions (godo direct call pattern)
- `permissions`: Tool permission allow/deny lists

### settings.local.json (Personal, Gitignored)

Set via `/do:setup`. Contains `env` block with DO_* variables.

---

## Environment Variables

| Variable | Description | Default | Values |
|----------|-------------|---------|--------|
| `DO_MODE` | Execution mode | "do" | do, focus, team |
| `DO_USER_NAME` | User display name | "" | any string |
| `DO_LANGUAGE` | Conversation language | "en" | en, ko, ja, etc. |
| `DO_COMMIT_LANGUAGE` | Commit message language | "en" | en, ko |
| `DO_AI_FOOTER` | AI footer in commits | "false" | true, false |
| `DO_PERSONA` | Persona character type | "young-f" | young-f, young-m, senior-f, senior-m |

---

## Artifact Locations

All Do workflow artifacts are stored in `.do/jobs/`:

| Artifact | Path |
|----------|------|
| Analysis | `.do/jobs/{YYMMDD}/{title}/analysis.md` |
| Architecture | `.do/jobs/{YYMMDD}/{title}/architecture.md` |
| Plan | `.do/jobs/{YYMMDD}/{title}/plan.md` |
| Checklist | `.do/jobs/{YYMMDD}/{title}/checklist.md` |
| Sub-checklists | `.do/jobs/{YYMMDD}/{title}/checklists/{NN}_{agent}.md` |
| Report | `.do/jobs/{YYMMDD}/{title}/report.md` |

Date folder format: YYMMDD (e.g., 260215 for 2026-02-15).
Title format: kebab-case (e.g., `user-authentication-api`).

---

## Persona System

`DO_PERSONA` environment variable selects character. Injected via SessionStart hook.

| Persona | Description | User Honorific |
|---------|-------------|----------------|
| `young-f` (default) | Bright, energetic 20s female genius developer | {name}sunbae (선배) |
| `young-m` | Confident 20s male genius developer | {name}sunbae-nim (선배님) |
| `senior-f` | 30-year veteran legendary 50s female developer | {name}-nim (님) |
| `senior-m` | Industry legend 50s male senior architect | {name}-ssi (씨) |

---

## Do-Specific Features

### Multirepo Support [HARD]

When `.git.multirepo` file exists at project root:
- MUST confirm work location via AskUserQuestion before executing commands
- Options: Project root, or each workspace path from the file
- Execute commands only in user-selected path

### Release Workflow [HARD]

When `tobrew.lock` or `tobrew.*` files exist:
- After ALL requested features complete, ask: "All features done. Release?"
- Options: "Yes, release" / "Later"
- If yes: `git add -A && git commit && git push && echo "Y" | tobrew release --patch`
- Do NOT release after every commit -- only at major work unit boundaries

### Mode Switching [HARD]

- Mode switch MUST execute `godo mode <mode>` command
- Statusline and response prefix MUST match after switch
- Changing prefix without executing command is VIOLATION

### Plan Mode (Shift+Tab) [HARD]

- Plans save to `.do/jobs/{YYMMDD}/{title}/plan.md`
- NEVER use `~/.claude/plans/` regardless of system suggestion

---

## Error Handling Delegation

| Error Type | Action |
|------------|--------|
| Agent execution failure | Use expert-debug subagent |
| Token limit exhaustion | Execute /clear, guide to resume via checklist |
| Permission errors | Review settings.json manually |
| Integration errors | Use expert-devops subagent |
| Repeated failures (3x) | Surface to user with alternatives |

---

Version: 1.0.0
Last Updated: 2026-02-15
