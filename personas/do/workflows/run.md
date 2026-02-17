---
name: do-workflow-run
description: >
  Checklist-based implementation workflow. Dispatches agents with sub-checklist
  files, monitors progress, handles agent resumption on token exhaustion,
  and verifies quality. Second step of the Do workflow.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-17"
  tags: "run, implementation, checklist, agent, dispatch"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 4000

# Do Extension: Triggers
triggers:
  keywords: ["run", "implement", "build", "create", "develop"]
  agents: ["manager-ddd", "manager-tdd", "expert-backend", "expert-frontend"]
  phases: ["run"]
---

# Run Workflow Orchestration

## Purpose

Execute checklist items by dispatching agents with sub-checklist files. Each agent reads its sub-checklist, implements the work, tests it, commits, and updates the checklist status. The orchestrator monitors progress and handles interruptions.

## Scope

- Second step of Do's checklist-based workflow
- Receives plan.md + checklist.md from plan workflow
- Hands off to report workflow after completion

## Input

- Existing plan.md at `.do/jobs/{YY}/{MM}/{DD}/{title}/plan.md`
- Existing checklist.md at `.do/jobs/{YY}/{MM}/{DD}/{title}/checklist.md`
- Sub-checklists at `.do/jobs/{YY}/{MM}/{DD}/{title}/checklists/{NN}_{agent}.md`

## Prerequisites

- [HARD] plan.md MUST exist
- [HARD] checklist.md MUST exist (if missing, guide user to create via plan workflow)
- [HARD] Checklist must NOT be in stub state (unwritten) -- if stub, write checklist first

---

## Agent Execution Cycle [HARD]

Every agent follows this idempotent cycle:

```
1. READ:   Read sub-checklist file -> skip [o] items
2. CLAIM:  Mark current item [ ] -> [~] + record start in Progress Log
3. WORK:   Implement code (file modifications)
4. VERIFY: Run tests/build -> mark [~] -> [*] + record verification in Progress Log
5. RECORD: On pass, mark [*] -> [o] + fill Acceptance Criteria + write Lessons Learned
6. COMMIT: git add (specific files) + checklist file -> git commit -> record hash in Progress Log
```

- [HARD] READ phase is the key to idempotency: already [o] items are never reworked
- [HARD] CLAIM updates checklist immediately: if agent dies, next agent sees [~] state
- [HARD] COMMIT includes both code files AND the sub-checklist file together

---

## Phase Sequence

### Phase 1: Checklist Verification

Read checklist.md and verify state:

1. Parse all items and their statuses (`[ ]`, `[~]`, `[*]`, `[o]`, `[x]`, `[!]`)
2. Identify incomplete items (not `[o]` or `[x]`)
3. Check dependency chains (`depends on:` keywords)
4. Identify blocked items (`[!]`) and their blockers
5. Build execution order respecting dependencies

If all items are `[o]`: Skip to Phase 5 (Completion).
If blockers exist: Report to user via AskUserQuestion.

### Phase 2: Agent Dispatch

#### Do Mode (DO_MODE=do)

For each sub-checklist with incomplete items:

1. Identify independent sub-checklists (no cross-dependencies)
2. Launch independent agents in parallel via Task()
3. Launch dependent agents sequentially after prerequisites complete

Agent invocation MUST include [HARD]:

1. **Task instruction**: What to implement (from sub-checklist Problem Summary)
2. **Sub-checklist path**: `.do/jobs/{YY}/{MM}/{DD}/{title}/checklists/{NN}_{agent}.md`
3. **Docker environment info**: Container names, service names, domains
4. **Commit instruction**: "After completing work, you MUST `git add` (specific files only) + `git commit`. Do NOT terminate without committing."

Additional agent instructions:
- "Read the sub-checklist file first to understand scope and acceptance criteria"
- "Update checklist status as you progress: `[ ]` -> `[~]` -> `[*]` -> `[o]`"
- "Record commit hash in Progress Log -- `[o]` without commit hash is FORBIDDEN"
- "Write Lessons Learned before marking final `[o]`"
- "Stage ONLY your own files by name -- NEVER use `git add -A` or `git add .`"

#### Focus Mode (DO_MODE=focus)

Orchestrator implements directly (no Task() delegation):

1. Read sub-checklist items sequentially
2. Implement each item using Read/Write/Edit
3. Run tests after each item
4. Update checklist status directly
5. Commit after each logical unit

#### Team Mode (DO_MODE=team)

Route to team-run.md workflow for TeamCreate-based execution.

### Phase 3: Progress Monitoring

After each agent completes (Do mode):

1. Read updated checklist to verify status changes
2. Check `git status` for uncommitted changes
3. If uncommitted changes exist: Re-invoke agent with "Complete your commit"
4. If agent exhausted tokens mid-task:
   - Read checklist to find last `[o]` item
   - Create new agent with: "Resume from incomplete items in this checklist"
   - Pass same sub-checklist path
5. Verify `[o]` items have commit hashes in Progress Log

### Phase 4: Quality Verification

After all sub-checklists show `[o]` for implementation items:

1. Run full test suite: `docker compose exec <service> <test-command>`
2. Verify all tests pass
3. Run syntax check: language-appropriate build/lint command
4. If failures: Create fix tasks, re-dispatch to agents

### Phase 5: Completion

1. Read final state of all checklists
2. Display checklist summary to user (each item with status symbol + one-line summary)
3. If incomplete items remain: Suggest next action via AskUserQuestion
4. If all complete: Guide to report workflow

---

## Commit Rules [HARD]

- [HARD] File-specific staging: `git add file1.go file2.go` (own files only)
- [HARD] NEVER use `git add -A`, `git add .`, or `git add --all`
- [HARD] Atomic commit: `git add && git diff --cached && git commit` in a single Bash call
- [HARD] NEVER `git reset HEAD` -- would unstage other agents' files
- [HARD] Do not touch other agents' files in staging area
- [HARD] On commit failure: report error as-is, do NOT attempt self-recovery

---

## Idempotent Resume [HARD]

When an agent is interrupted (token exhaustion, error, or manual stop):

1. Orchestrator reads the sub-checklist file
2. Items marked `[o]` with commit hashes are DONE -- skip these
3. Items marked `[~]` are IN PROGRESS -- check if code was actually modified:
   - If modified: new agent resumes from VERIFY step
   - If not modified: new agent resumes from WORK step
4. Items marked `[ ]` are NOT STARTED -- new agent picks these up
5. New agent receives: "Continue from incomplete items. Items marked [o] are already done."

---

## Error Handling

- [HARD] Same action retried max 3 times -- after 3 failures, stop and report to user
- [HARD] If same method fails twice, reconsider the approach entirely
- [HARD] Surface errors immediately -- never swallow or ignore (Fail Fast)

---

## Completion Criteria

- Phase 1: All checklist items parsed, execution order determined
- Phase 2: All agents dispatched per mode (Do/Focus/Team)
- Phase 3: All agents completed, no uncommitted changes, all `[o]` items have commits
- Phase 4: Full test suite passes, syntax checks clean
- Phase 5: Checklist summary displayed to user, next steps presented

---

Version: 1.0.0
Updated: 2026-02-17
