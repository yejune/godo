# Do Development Workflow

Do's development workflow with flexible token management and commit-as-proof tracking.

## Workflow Overview

### Simple Tasks (4 or fewer files, single domain, no architecture change)

```
Plan -> Checklist -> Develop -> Test -> Report
```

### Complex Tasks (5+ files, new modules, migrations, 3+ domains, abstraction design)

```
Analysis -> Architecture -> Plan -> Checklist -> Develop -> Test -> Report
```

| Phase | Agent | Purpose |
|-------|-------|---------|
| Analysis | expert-analyst | Current system analysis, requirements (EARS+MoSCoW), tech comparison |
| Architecture | expert-architect | Solution design, interfaces, implementation order |
| Plan | manager-plan | Create plan document at `.do/jobs/{YYMMDD}/{title}/plan.md` |
| Checklist | (orchestrator delegates) | Decompose plan into agent-sized items (1-3 files each) |
| Develop | manager-ddd/tdd or experts | Implementation with living checklist updates |
| Test | (within develop) | Verification per Test Strategy (unit/integration/E2E/pass) |
| Report | (orchestrator delegates) | Final report at `.do/jobs/{YYMMDD}/{title}/report.md` |

## Plan Phase

Create plan using EARS+MoSCoW requirements format.

Output:
- Plan document at `.do/jobs/{YYMMDD}/{title-kebab-case}/plan.md`
- EARS format requirements with MoSCoW priority (MUST/SHOULD/COULD/WON'T)
- Acceptance criteria
- Technical approach
- Complexity assessment (simple vs complex)

Plan-mode artifact writing rule:
- Plan-mode agents (researcher, analyst, architect) run with `permissionMode: plan` (read-only)
- Write/Edit tools are blocked in plan mode to protect source code
- Artifacts are written via Bash heredoc to `.do/jobs/` directory only
- Source code modification via Bash is also forbidden — only new file creation in `.do/jobs/`
- Files are the source of truth, not messages — ensures idempotency across sessions

## Checklist Phase

Decompose plan into executable checklist items.

Rules:
- One item = 1-3 file changes + verification
- Items exceeding 3 files MUST be split further
- Each item gets a Test Strategy declaration (unit/integration/E2E/pass)
- Items assigned to agents via sub-checklists in `checklists/` directory
- Checklist is a living document — evolves during implementation

Output:
- Main checklist: `.do/jobs/{YYMMDD}/{title}/checklist.md`
- Sub-checklists: `.do/jobs/{YYMMDD}/{title}/checklists/{order}_{agent-topic}.md`

## Develop Phase

Implement checklist items using configured development methodology.

Development Methodology:
- DDD (ANALYZE-PRESERVE-IMPROVE): For legacy refactoring
- TDD (RED-GREEN-REFACTOR): For new features
- Hybrid: Mixed per change type (recommended for most projects)
- See @workflow-modes.md for detailed methodology cycles

Living Checklist Updates:
- Agent reads sub-checklist before starting
- Updates status as work progresses: [ ] -> [~] -> [*] -> [o]
- If agent stops (token exhaustion), next agent reads checklist and resumes from last state
- Checklist file IS the handoff mechanism

Commit-as-proof:
- Every [o] completion requires a recorded commit hash
- `git add <specific files>` only (NEVER `git add -A` or `git add .`)
- `git diff --cached --name-only` to verify before commit
- Commit message explains WHY (diff shows WHAT)
- Record hash in Progress Log: `[o] completed (commit: <hash>)`

Success Criteria:
- All checklist requirements implemented
- Methodology-specific tests passing
- Quality dimensions verified (Tested/Readable/Unified/Secured/Trackable)

## Test Phase

Verification is integrated into the Develop phase per Test Strategy declarations.

For testable code (business logic, API, data layer):
- dev-testing.md rules apply: FIRST principles, Real DB only, AI anti-pattern 7
- Behavior-based testing, not implementation testing
- Mutation testing mindset: "if I change this line, does a test fail?"

For non-testable changes (CSS, config, docs, hooks):
- Alternative verification: build check, manual check, syntax check
- Declared as `pass` in Test Strategy with justification

## Report Phase

Generate completion report after all checklist items are done.

Output:
- Report at `.do/jobs/{YYMMDD}/{title}/report.md`
- Execution summary, plan deviations, test results, changed files, lessons learned

## Context Management

/clear Strategy (flexible, not fixed phases):
- At checklist item boundaries (natural work units)
- When context exceeds threshold (not a fixed number — depends on task)
- Between major workflow transitions
- NOT at rigid pre-allocated phase boundaries

Progressive Disclosure:
- Level 1: Metadata only (~100 tokens)
- Level 2: Skill body when triggered (~5000 tokens)
- Level 3: Bundled files on-demand

## Agent Teams Variant

When team mode is enabled (CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1), phases can execute with Agent Teams.

### Team Mode Phase Overview

| Phase | Sub-agent Mode | Team Mode | Condition |
|-------|---------------|-----------|-----------|
| Plan | manager-plan (single) | researcher + analyst + architect (parallel) | Complexity >= threshold |
| Develop | manager-ddd/tdd (sequential) | backend-dev + frontend-dev + tester (parallel) | Domains >= 3 or files >= 10 |
| Report | (orchestrator) | (orchestrator) | Always single |

### Team Mode Plan Phase
- TeamCreate for parallel research team
- Teammates explore codebase, analyze requirements, design approach
- Each teammate writes artifacts to `.do/jobs/` via Bash heredoc (plan-mode idempotency)
- Do synthesizes into plan and checklist
- Shutdown team, /clear before Develop phase

### Team Mode Develop Phase
- TeamCreate for implementation team
- Task decomposition with file ownership boundaries (one file = one owner)
- Teammates self-claim tasks from shared checklist
- Each teammate commits own files only (NEVER `git add -A`)
- Quality validation after all implementation completes
- Shutdown team

### Mode Selection
- --team flag: Force team mode
- --solo flag: Force sub-agent mode
- No flag (default): Complexity-based selection

### Fallback
If team mode fails or is unavailable:
- Graceful fallback to sub-agent mode
- Continue from last completed checklist item (living checklist ensures continuity)
- No data loss or state corruption
