# Agent Execution Cycle

Purpose: Defines the mandatory handoff protocol for delegating work to agents and the idempotent execution cycle that guarantees resumability and progress tracking through checklist state files.

Version: 1.0.0
Last Updated: 2026-02-17

---

## Quick Reference (30 seconds)

Agent Delegation Checklist — every agent call MUST include:
1. Task instruction (what to do)
2. Sub-checklist file path (the agent's work specification)
3. Docker environment info (container name, service name, domain)
4. Commit instruction: "After completion, run `git add` + `git commit`. Never exit without committing."
5. Checklist update instruction: "Update the sub-checklist status as you progress. A commit without checklist update is considered incomplete."
6. Jobs document language: Follow `DO_JOBS_LANGUAGE` env var (default: "en")

Language Rules:
- `DO_JOBS_LANGUAGE` controls jobs document language (analysis.md, plan.md, checklist.md, report.md)
- `DO_LANGUAGE` (conversation language) and `DO_JOBS_LANGUAGE` (document language) are independent

---

## Implementation Guide (5 minutes)

### Idempotent Execution Cycle [HARD]

Every agent repeats this cycle. Running it multiple times MUST produce the same result.

```
1. READ:   Read sub-checklist file → skip [o] items
2. CLAIM:  Mark item [ ] → [~] + log start in Progress Log
3. WORK:   Implement code (modify/create files)
4. VERIFY: Run tests/build → mark [~] → [*] + log verification in Progress Log
5. RECORD: On pass, mark [*] → [o] + check Acceptance Criteria + write Lessons Learned
6. COMMIT: Commit code + checklist together → log commit hash in Progress Log
```

### Key Invariants

- **READ is the core**: Items already marked [o] are NEVER reworked — this is the foundation of idempotency
- **CLAIM immediately**: Mark [~] the moment work starts — if interrupted, the next agent knows the state
- **COMMIT includes checklist**: `git add` MUST include both code files AND the sub-checklist file
- **Checklist = agent state file**: Progress is persisted to disk, not held in memory
- **Commit messages** follow the commit discipline rules (atomic, explain WHY)
- **Orchestrator verifies**: After agent completion, run `git status` to check for uncommitted changes — if found, re-invoke agent to commit

### Checklist Non-Update = VIOLATION [HARD]

- Code committed but sub-checklist still at `[ ]` → VIOLATION
- Main checklist at `[o]` but sub-checklist at `[ ]` → inconsistency VIOLATION
- Progress Log empty → VIOLATION (minimum 2 entries: start + completion)
- Lessons Learned empty on `[o]` completion → VIOLATION

---

## Works Well With

Foundation Modules:
- [Delegation Patterns](delegation-patterns.md) - Task delegation strategies
- [Agent Delegation](agent-delegation.md) - Interruption and resumption rules
- [Agent Research](agent-research.md) - Research delegation constraints

Skills:
- {{slot:BRAND}}-workflow-spec - Checklist system and templates

---

Version: 1.0.0
Last Updated: 2026-02-17
Status: Production Ready
