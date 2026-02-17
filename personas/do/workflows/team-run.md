---
name: do-workflow-team-run
description: >
  Team-based implementation using Agent Teams API. Spawns parallel development
  team with file ownership boundaries, manages shared task list via
  TaskCreate/TaskUpdate, enforces commit isolation between teammates.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-17"
  tags: "team, run, implementation, parallel, agent-teams"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 4000

# Do Extension: Triggers
triggers:
  keywords: ["team run", "team implement", "parallel build"]
  agents: ["team-backend-dev", "team-frontend-dev", "team-tester"]
  phases: ["run"]
---

# Team Run Workflow Orchestration

## Purpose

Implement checklist items through parallel team-based development. Each teammate owns specific files/domains to prevent conflicts. Shared task list coordinates work via TaskCreate/TaskUpdate. Quality validation runs after all implementation completes.

## Scope

- Team variant of the run workflow
- Uses Agent Teams API for parallel implementation
- Enforces file ownership to prevent merge conflicts
- Produces the same outcome as solo run workflow but with parallelism

## Prerequisites

- Approved plan with checklist at `.do/jobs/{YY}/{MM}/{DD}/{title}/`
- Agent Teams feature enabled (CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1)
- Triggered by: Team mode active OR auto-detected complexity >= threshold
- Fallback: If Agent Teams unavailable, falls back to solo run workflow (run.md)

---

## Phase Sequence

### Phase 0: Task Decomposition

1. Read plan.md and checklist.md
2. Analyze scope to determine team composition:
   - Backend-only: team-backend-dev + team-tester
   - Frontend-only: team-frontend-dev + team-tester
   - Full-stack: team-backend-dev + team-frontend-dev + team-tester
   - With design: + team-designer
3. Assign file ownership per teammate (exclusive -- no overlapping files)
4. Create team:
   ```
   TeamCreate(team_name: "do-run-{feature-slug}")
   ```
5. Create shared task list with dependencies:
   ```
   TaskCreate: "Implement backend logic" (no deps)
   TaskCreate: "Implement frontend components" (blocked by backend)
   TaskCreate: "Write integration tests" (blocked by backend + frontend)
   TaskCreate: "Quality validation" (blocked by all above)
   ```

### Phase 1: Spawn Implementation Team

Select teammates based on Phase 0 analysis.

**team-backend-dev** (inherit model, acceptEdits mode):
- Spawn prompt includes: task summary, assigned sub-checklist path, file ownership list, Docker environment info
- Owns: server-side code, API routes, data models, business logic

**team-frontend-dev** (inherit model, acceptEdits mode):
- Spawn prompt includes: task summary, assigned sub-checklist path, file ownership list, Docker environment info
- Owns: UI components, pages, client-side state, styles

**team-tester** (inherit model, acceptEdits mode):
- Spawn prompt includes: test strategy, assigned sub-checklist path, coverage targets
- Owns: test files exclusively (no production code)

**team-designer** (inherit model, acceptEdits mode, when needed):
- Spawn prompt includes: design requirements, UI/UX specifications
- Owns: design files (.pen, design tokens, style configs)

### Phase 2: Parallel Implementation

Teammates self-claim tasks from the shared list and work independently:

- Each teammate follows the agent execution cycle: READ -> CLAIM -> WORK -> VERIFY -> RECORD -> COMMIT
- Teammates communicate via SendMessage for cross-cutting concerns (API contracts, data shapes)
- Orchestrator monitors via TaskList and resolves blocking issues

Coordination patterns:
- Backend notifies frontend when API contracts are ready
- Designer shares design specs with frontend via SendMessage
- Tester waits for implementation completion before integration tests
- Orchestrator forwards cross-team information as needed

### Phase 3: Quality Validation

After all implementation tasks complete:

**Option A** (with team-quality teammate):
- Assign quality validation task to team-quality
- team-quality runs TRUST 5 checks (read-only mode)
- Reports findings to orchestrator
- Orchestrator directs fixes to responsible teammates

**Option B** (with subagent):
- Delegate to manager-quality subagent (NOT a teammate)
- Review findings and create fix tasks
- Assign fixes to existing teammates

Quality gates (must all pass):
- Zero lint errors
- Zero type errors
- Coverage targets met (85%+ overall)
- All acceptance criteria verified

### Phase 4: Cleanup

1. Shutdown all teammates gracefully:
   ```
   SendMessage(type: "shutdown_request", recipient: "backend-dev")
   SendMessage(type: "shutdown_request", recipient: "frontend-dev")
   SendMessage(type: "shutdown_request", recipient: "tester")
   ```
2. Report implementation summary to user
3. Guide to report workflow

---

## File Ownership [HARD]

- [HARD] Each teammate owns specific files -- no two teammates modify the same file
- [HARD] Ownership is declared at spawn time and enforced throughout
- [HARD] If a change requires cross-ownership files, orchestrator reassigns or splits the task
- [HARD] Test files are owned exclusively by team-tester

---

## Parallel Commit Isolation [HARD]

- [HARD] `git add file1.go file2.go && git diff --cached && git commit` in a single Bash call
- [HARD] Stage ONLY own files by name -- NEVER `git add -A` or `git add .`
- [HARD] NEVER `git reset HEAD` -- would unstage other teammates' files
- [HARD] Do not touch other teammates' files in staging area
- [HARD] On commit failure: report error as-is, do NOT attempt self-recovery
- [HARD] Orchestrator resolves commit conflicts -- teammates never self-fix

---

## Task Tracking [HARD]

All task status changes via TaskUpdate:
- pending -> in_progress: When teammate claims task
- in_progress -> completed: When task work is verified and committed
- Never use plain text TODO lists -- TaskCreate/TaskUpdate is the single source of truth

---

## Fallback

If team mode fails at any point:
- Shutdown remaining teammates gracefully
- Fall back to solo run workflow (workflows/run.md)
- Continue from the last completed task (checklist preserves state)
- Log warning about team mode failure

---

## Completion Criteria

- Phase 0: Team created, tasks decomposed, file ownership assigned
- Phase 1: All teammates spawned with correct sub-checklists
- Phase 2: All implementation tasks completed, committed
- Phase 3: Quality validation passed, all gates green
- Phase 4: Teammates shut down, summary presented to user

---

Version: 1.0.0
Updated: 2026-02-17
