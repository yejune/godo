# Checklist System Rules [HARD]

## Creation Timing
- [HARD] Create checklist after plan is confirmed (one per agent task)
- [HARD] Never create checklist without a plan — the plan is the checklist's basis
- PostToolUse hook automatically creates checklist stub when plan file is generated
- [HARD] If plan file exists and checklist is in stub state (not yet written), prioritize checklist creation over all other work
- [HARD] Checklist creation procedure: Read plan file → Decompose into agent tasks → Invoke agent (Task tool) to write checklist
- [HARD] Items must be granular enough for an agent to complete without token exhaustion
- [HARD] One item = 1-3 file changes + verification — split if exceeding this scope
- [HARD] Specify verification method per item: testable → unit/integration test, not testable → build check/manual check/`docker compose config` etc.

### Decomposition Procedure [HARD]
- [HARD] Step 1: Estimate how many files each item touches
- [HARD] Step 2: More than 3 files → must decompose
- [HARD] Step 3: Each decomposed item must be independently completable/verifiable
- [HARD] Step 4: If dependencies exist between items, link with `depends on:`
- [HARD] Decomposition example:
  - "Implement API" ← too large (5+ files)
  - → "Define router" (1 file), "Implement handler" (1 file), "Add validation logic" (1 file), "Error response handling" (1 file), "Unit tests" (1 file)
- [HARD] Proceeding with development while checklist is unwritten is a VIOLATION

## Authoring Rules [HARD]
- [HARD] All documents in jobs folder (checklist.md, report.md, checklists/*.md) must be delegated to agents (Task tool)
- [HARD] Same for both Do/Focus modes — orchestrator never directly writes/edits jobs folder files
- [HARD] Reason: Prevent orchestrator context token waste — document authoring is the agent's responsibility
- [HARD] Only exception: plan.md — Plan Mode hook auto-generates/moves it (not written by orchestrator)

## Checklist = Agent State File [HARD]
- [HARD] Checklist is not just a document — it is the **agent's persistent state store**
- [HARD] Agent reads checklist at task start → determines work scope
- [HARD] Updates checklist state on each item completion → progress recorded in file
- [HARD] On agent token exhaustion/interruption → last state remains in checklist
- [HARD] When new agent receives same checklist → skips `[o]` items, resumes from incomplete
- [HARD] This pattern guarantees **work continuity** — any agent can pick up where another left off

## File Structure (jobs directory integration) [HARD]
- [HARD] One task = one folder — all artifacts in same directory:
  - Analysis: `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/analysis.md` (complex tasks only)
  - Architecture: `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/architecture.md` (complex tasks only)
  - Plan: `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/plan.md`
  - Checklist: `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/checklist.md`
  - Report: `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/report.md`
  - Sub files: `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/checklists/{order}_{agent-topic}.md`
- [HARD] Sub file `{order}` is two-digit number: `01`, `02`, ... `99`
- Directories auto-created if missing

### Example (complex task)
```
.do/jobs/26/02/11/queue-library-migration/
  ├── analysis.md
  ├── architecture.md
  ├── plan.md
  ├── checklist.md
  ├── report.md
  └── checklists/
      ├── 01_expert-backend.md
      ├── 02_expert-security.md
      └── 03_expert-testing.md
```

### Example (simple task)
```
.do/jobs/26/02/11/login-api-security/
  ├── plan.md
  ├── checklist.md
  ├── report.md
  └── checklists/
      ├── 01_expert-backend.md
      └── 02_expert-testing.md
```

## State Management

### State Symbols
| Symbol | State | Meaning |
|--------|-------|---------|
| `[ ]` | pending | Not yet started |
| `[~]` | in progress | Currently being worked on |
| `[*]` | testing | Implementation done, test verification in progress |
| `[!]` | blocked | Waiting on external dependency/decision |
| `[o]` | done | Tests passed, work complete |
| `[x]` | failed | Final failure, no further progress possible |

> **[HARD] `[x]` = FAILED, `[o]` = DONE. Never confuse them.**
> - Completed items in `.do/jobs/` checklists must use `[o]`
> - `[x]` means "final failure" — different from standard markdown `[x]` ("checked")
> - Same applies to Acceptance Criteria, FINAL STEP, and all other items
> - Marking completion with `[x]` is a VIOLATION

### State Transition Rules [HARD]
- [HARD] Allowed transitions:
  ```
  [ ] → [~]        start
  [~] → [*]        implementation done → testing
  [~] → [!]        blocker occurred
  [*] → [o]        tests passed → done
  [*] → [~]        test failed → rework (regression)
  [*] → [x]        test final failure → failed
  [~] → [x]        implementation final failure → failed
  [!] → [~]        blocker resolved → resume
  ```
- [HARD] Forbidden transitions: `[ ] → [o]` (cannot complete without testing), `[ ] → [x]` (cannot fail without working), `[ ] → [*]` (cannot test without working)
- [HARD] Record history on state changes (never overwrite)

### Blocker Recording Rules [HARD]
- [HARD] On `[!]` transition, must record 3 things:
  1. **What** is blocking (specific reason)
  2. **Who** can resolve it (owner/external system)
  3. **When** it was blocked (timestamp)

### State History Example
```
[~] Login API implementation
    - [ ] 2026-02-11 14:00:00 created
    - [~] 2026-02-11 14:05:12 work started
    - [!] 2026-02-11 15:00:33 blocker: Redis config incomplete (owner: infra team, awaiting resolution)
    - [~] 2026-02-11 16:00:05 blocker resolved, resumed
    - [*] 2026-02-11 17:00:41 testing
    - [~] 2026-02-11 17:30:18 test failed (JWT expiry logic error) → rework
    - [*] 2026-02-11 18:00:22 re-testing
    - [o] 2026-02-11 18:30:55 done
```

## Dependency Management [HARD]
- [HARD] Declare dependencies between items with `depends on:` keyword
- [HARD] If dependency target is incomplete, the item is automatically treated as `[!]` blocked
- [HARD] Dependencies managed in main checklist (cross-referencing sub files)

### Dependency Notation
```
## Task List
- [o] #1 DB schema migration
- [~] #2 Login API implementation (depends on: #1)
- [ ] #3 Frontend login form (depends on: #2)
- [!] #4 Social login integration (depends on: #2, blocker: OAuth key not issued)
```

## [~] Start Enforcement [HARD]
- [HARD] Main checklist Phase/item must be set to `[~]` before sub-checklist work can begin
- [HARD] Sub-checklist item must be set to `[~]` before actual coding work can begin
- [HARD] No completion without start — `[ ] → [o]` direct transition is VIOLATION
- [HARD] Agent must update main checklist item to `[~]` first when starting work
- [HARD] This rule guarantees real-time tracking of progress

## Individual Item Tracking [HARD]
- [HARD] Items must transition state individually — batch completion is VIOLATION
- [HARD] Complete one item → update checklist state → start next item
- [HARD] Marking multiple items as `[o]` simultaneously is VIOLATION

## Progress Summary Table Update [HARD]
- [HARD] Main checklist.md progress summary table must be updated on every state change
- [HARD] When agent completes sub-checklist item, also update main checklist summary table
- [HARD] User must be able to see full progress status by reading main checklist.md alone

## Sub-Checklist Split Threshold [HARD]
- [HARD] Review splitting when sub-checklist has more than 5 items
- [HARD] Mandatory split when sub-checklist has more than 10 items — create additional agents
- [HARD] Each split sub-checklist must be independently completable

## Project Git Protection [HARD]

Destructive git commands that overwrite/delete files are forbidden project-wide.

- [HARD] Banned commands:
  - `git checkout .` / `git checkout -- <path>` — overwrites working content
  - `git reset --hard` — destroys committed changes
  - `git reset HEAD <path>` — unstaging causes commit omission
  - `git clean -f` / `git clean -fd` — deletes untracked files
  - `git stash` — risk of loss during temporary storage
  - `git restore .` / `git restore <path>` — same as checkout overwrite
  - `git rebase` (with squash/drop) — alters commit history
  - `git push --force` / `git push -f` — destroys remote history
  - `rm -rf` — direct deletion
- [HARD] Fix mistakes by `git add` → `git commit` — preserve history
- [HARD] Exception only when user explicitly instructs with full risk acknowledgment

## Checklist Path Tracing [HARD]
- [HARD] Include related sub-checklist path in commit message when committing code changes
- [HARD] Commit message format: `ref: .do/jobs/{path}/checklists/{NN}_{agent}.md`
- [HARD] Also include checklist path as comments in modified source code where possible
  - Example: `// ref: .do/jobs/26/02/17/login-api/checklists/01_expert-backend.md`
  - Use language-appropriate comment syntax (Go: `//`, Python: `#`, JS/TS: `//`)
  - Place near file top or changed function
- [HARD] Same format for subsequent modifications via add/commit — enables traceability

## Agent Token Exhaustion Protocol [HARD]
- [HARD] When agent reaches ~10% remaining tokens, self-assess:
  - If completable: finish the work
  - If not: record current state in checklist → report to super agent
- [HARD] Super agent creates new agent to continue work on report
- [HARD] Checklist is the sole handoff mechanism — never depend on message passing

## Agent Prompt Token Optimization [HARD]
- [HARD] Never instruct agent to read entire large files (500+ lines)
- [HARD] Extract only relevant sections and inject directly into prompt
- [HARD] For large artifacts (architecture.md, analysis.md), provide only the relevant Phase excerpt

## Jobs Continuation Strategy [HARD]
- [HARD] Job folders are immutable records -- once created, content must not be modified
- [HARD] Only allowed modification to existing job: adding "Continued in: {path}" reference
- [HARD] When execution continues or modifies a previous job's work, create a NEW job folder with current date
- [HARD] New job format: `.do/jobs/{YY}/{MM}/{DD}/{title}/`
- [HARD] New job's plan.md must start with: `Continues from: {path to previous job}`
- [HARD] New job's checklist.md must specify which previous sub-checklists are being continued
- [HARD] Bidirectional linking required:
  - Previous job checklist: add `Continued in: {new job path}` at bottom
  - New job checklist: add `Continues from: {previous job path}` at top
- [HARD] Sub-checklists also linked: new sub-checklist header references the previous one it continues
- [HARD] Research/analysis artifacts stay in original job -- only execution moves to new job
