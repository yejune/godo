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

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 150
  level2_tokens: 5000

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

### Commit-as-Proof

Completion evidence is a git commit hash -- not an in-memory marker or conversation signal.

- `[o]` transition requires commit hash: `[o] Done (commit: a1b2c3d)`
- Append-only: commit messages are never rewritten (no `--amend`, no `--force`)
- Atomic: one logical change = one commit
- Agent verification layer: Read(original) -> modify -> git diff(verify) -> confirm intended changes only
- Uncommitted work = incomplete work, regardless of code quality

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

```json
{
  "env": {
    "DO_USER_NAME": "name",
    "DO_LANGUAGE": "ko",
    "DO_COMMIT_LANGUAGE": "en",
    "DO_JOBS_LANGUAGE": "en",
    "DO_PERSONA": "young-f",
    "DO_STYLE": "pair",
    "DO_MODE": "do",
    "DO_AI_FOOTER": "false"
  }
}
```

Hooks access these as environment variables: `$DO_USER_NAME`, `$DO_LANGUAGE`, etc.

---

## Environment Variables

| Variable | Description | Default | Values |
|----------|-------------|---------|--------|
| `DO_MODE` | Execution mode | "do" | do, focus, team |
| `DO_USER_NAME` | User display name | "" | any string |
| `DO_LANGUAGE` | Conversation language | "en" | en, ko, ja, etc. |
| `DO_COMMIT_LANGUAGE` | Commit message language | "en" | en, ko |
| `DO_JOBS_LANGUAGE` | Jobs message language | "en" | en, ko |
| `DO_AI_FOOTER` | AI footer in commits | "false" | true, false |
| `DO_PERSONA` | Persona character type | "young-f" | young-f, young-m, senior-f, senior-m |

---

## Artifact Locations

All Do workflow artifacts are stored in `.do/jobs/`:

| Artifact | Path |
|----------|------|
| Research | `.do/jobs/{YY}/{MM}/{DD}/{title}/research_{topic}.md` |
| Analysis | `.do/jobs/{YY}/{MM}/{DD}/{title}/analysis.md` |
| Architecture | `.do/jobs/{YY}/{MM}/{DD}/{title}/architecture.md` |
| Plan | `.do/jobs/{YY}/{MM}/{DD}/{title}/plan.md` |
| Checklist | `.do/jobs/{YY}/{MM}/{DD}/{title}/checklist.md` |
| Sub-checklists | `.do/jobs/{YY}/{MM}/{DD}/{title}/checklists/{NN}_{agent}.md` |
| Report | `.do/jobs/{YY}/{MM}/{DD}/{title}/report.md` |

Date folder format: YY/MM/DD (e.g., 26/02/15 for 2026-02-15).
Title format: kebab-case (e.g., `user-authentication-api`).
Versioning: optionally add version postfix (e.g., analysis_v2.md)

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

- Plans save to `.do/jobs/{YY}/{MM}/{DD}/{title}/plan.md`
- NEVER use `~/.claude/plans/` regardless of system suggestion

---

## Quality Gates (Built-in Rules)

Do enforces 5 quality dimensions as always-active built-in rules (not a branded framework):

| Dimension | Enforcement | Source |
|-----------|------------|--------|
| Tested | FIRST principles, 85%+ coverage, Real DB only, AI anti-pattern 7 | dev-testing.md |
| Readable | Read Before Write, clear naming, match existing conventions | dev-workflow.md |
| Unified | Language-specific syntax checks (go vet, tsc --noEmit, ruff) | dev-environment.md |
| Secured | No secrets in commits, input validation, OWASP guidelines | core rules |
| Trackable | Atomic commits, WHY in messages, commit hash as proof | dev-workflow.md, dev-checklist.md |

### AI Anti-Pattern Prevention (7 Rules)

These are FORBIDDEN during all development work:
1. Weakening assertions (fake pass)
2. Swallowing errors with try/catch (fake success)
3. Adjusting expected values to wrong output (fake correctness)
4. Using time.sleep() for timing issues (fake stability)
5. Deleting/commenting failing tests (fake integrity)
6. Using wildcard matchers when exact values known (fake verification)
7. Testing only happy path (fake coverage)

### Test Strategy Pre-Declaration

Each checklist item declares test strategy BEFORE implementation:
- Testable code: specify test type + file path (e.g., `unit: handler_test.go`)
- Non-testable changes: declare `pass` with alternative verification (e.g., `pass (build check: go build ./...)`)

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
Last Updated: 2026-02-16
