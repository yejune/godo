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
  updated: "2026-02-15"
  tags: "run, implementation, checklist, agent, dispatch"

# Do Extension: Triggers
triggers:
  keywords: ["run", "implement", "build", "create", "develop", "구현", "만들어"]
  agents: ["manager-ddd", "manager-tdd", "expert-backend", "expert-frontend"]
  phases: ["run"]
---

# Run Workflow Orchestration

## Purpose

Execute checklist items by dispatching agents with sub-checklist files. Each agent reads its sub-checklist, implements the work, tests it, commits, and updates the checklist status. The orchestrator monitors progress and handles interruptions.

## Scope

- Second step of Do's checklist-based workflow
- Receives plan.md + checklist.md from plan workflow
- Hands off to test/report workflows after completion

## Input

- Existing plan.md at `.do/jobs/{YY}/{MM}/{DD}/{title}/plan.md`
- Existing checklist.md at `.do/jobs/{YY}/{MM}/{DD}/{title}/checklist.md`
- Sub-checklists at `.do/jobs/{YY}/{MM}/{DD}/{title}/checklists/{NN}_{agent}.md`

## Prerequisites

- [HARD] plan.md MUST exist
- [HARD] checklist.md MUST exist (if missing, guide user to create via `/do:checklist` or plan workflow)
- [HARD] Checklist must NOT be in stub state (unwritten) -- if stub, write checklist first

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

Dispatch agents based on current execution mode.

#### Do Mode Dispatch (DO_MODE=do)

For each sub-checklist with incomplete items:

1. Identify independent sub-checklists (no cross-dependencies)
2. Launch independent agents in parallel via Task()
3. Launch dependent agents sequentially after prerequisites complete

Agent invocation MUST include these 4 items [HARD]:

1. **Task instruction**: What to implement (from sub-checklist Problem Summary)
2. **Sub-checklist path**: `.do/jobs/{YY}/{MM}/{DD}/{title}/checklists/{NN}_{agent}.md`
3. **Docker environment info**: Container names, service names, domains
4. **Commit instruction**: "After completing work, you MUST `git add` (specific files only) + `git commit`. Do NOT terminate without committing."

Additional agent instructions:
- "Read the sub-checklist file first to understand scope and acceptance criteria"
- "Update checklist status as you progress: `[ ]` -> `[~]` -> `[*]` -> `[o]`"
- "Record commit hash in Progress Log"
- "Write Lessons Learned before marking final `[o]`"

#### Focus Mode Dispatch (DO_MODE=focus)

Orchestrator implements directly (no Task() delegation):

1. Read sub-checklist items sequentially
2. Implement each item using Read/Write/Edit
3. Run tests after each item
4. Update checklist status directly
5. Commit after each logical unit

#### Team Mode Dispatch (DO_MODE=team)

Route to team-do.md workflow for TeamCreate-based execution.

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
2. Verify all tests pass (dev-testing.md rules)
3. Run syntax check: language-appropriate build/lint command
4. If failures: Create fix tasks, re-dispatch to agents

### Phase 5: Completion

1. Read final state of all checklists
2. Display checklist summary to user (each item with status symbol + one-line summary)
3. If incomplete items remain: Suggest next action via AskUserQuestion
4. If all complete: Guide to report workflow

---

## Agent Delegation Checklist [HARD]

Every agent invocation in Do mode MUST verify:

- [ ] Task instruction includes what to implement
- [ ] Sub-checklist file path is provided
- [ ] Docker environment info is included
- [ ] Commit instruction is explicit
- [ ] Agent is told to read sub-checklist first
- [ ] Agent is told to update checklist status
- [ ] Agent is told to record commit hash

---

## Agent Interruption & Resume [HARD]

When an agent is interrupted (token exhaustion, error, or manual stop):

1. Orchestrator reads the sub-checklist file
2. Items marked `[o]` with commit hashes are DONE -- skip these
3. Items marked `[~]` are IN PROGRESS -- new agent resumes from here
4. Items marked `[ ]` are NOT STARTED -- new agent picks these up
5. New agent receives: "Continue from incomplete items. Items marked [o] are already done."

This pattern ensures work continuity regardless of which agent instance runs.

---

## Completion Criteria

- Phase 1: All checklist items parsed, execution order determined
- Phase 2: All agents dispatched per mode (Do/Focus/Team)
- Phase 3: All agents completed, no uncommitted changes, all `[o]` items have commits
- Phase 4: Full test suite passes, syntax checks clean
- Phase 5: Checklist summary displayed to user, next steps presented

---

Version: 1.0.0
Updated: 2026-02-15
