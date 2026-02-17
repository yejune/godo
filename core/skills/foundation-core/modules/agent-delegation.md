# Agent Delegation — Interruption and Resumption

Purpose: Defines the rules for handling agent interruptions (token exhaustion, crashes) and idempotent resumption, ensuring no work is lost and no duplicate work is performed.

Version: 1.0.0
Last Updated: 2026-02-17

---

## Quick Reference (30 seconds)

When an agent stops mid-work:
1. Orchestrator checks the checklist file
2. `[o]` items → skip (already done)
3. `[ ]` or `[~]` items → delegate to a new agent starting from those
4. Pass the same sub-checklist file path → continuity guaranteed

Key Principle: The checklist file IS the state. A new agent reads it, sees what's done, and picks up where the previous one left off.

---

## Implementation Guide (5 minutes)

### Agent Interruption & Resumption (Idempotent Resume) [HARD]

- When an agent is interrupted (token exhaustion, crash), the orchestrator checks the checklist
- Items marked `[o]` are skipped; items marked `[ ]` or `[~]` are delegated to a new agent
- The new agent receives the instruction: "Continue from incomplete items in this checklist"
- The same sub-checklist file path is always passed — ensuring continuity
- **Idempotent resumption guaranteed**: The new agent's READ step reads existing state, so no duplicate work occurs

### Handling `[~]` (In-Progress) Items [HARD]

When a new agent encounters a `[~]` item left by a previous agent:
1. Check whether the code was actually modified (inspect files listed in Critical Files)
2. If modified → resume from VERIFY step (run tests/build to validate)
3. If not modified → resume from WORK step (implement from scratch)

### Orchestrator Responsibilities

- After agent completion, always run `git status` to detect uncommitted changes
- If uncommitted changes exist, re-invoke agent to perform the commit
- When delegating to a new agent after interruption:
  - Pass the sub-checklist file path (same file, not a copy)
  - Include instruction: "Continue from incomplete items"
  - Provide Docker environment info (same as original delegation)

---

## Works Well With

Foundation Modules:
- [Agent Execution Cycle](agent-execution-cycle.md) - The execution cycle this module extends
- [Agent Research](agent-research.md) - Research delegation constraints
- [Delegation Patterns](delegation-patterns.md) - Task delegation strategies

Skills:
- {{slot:BRAND}}-workflow-spec - Checklist state management

---

Version: 1.0.0
Last Updated: 2026-02-17
Status: Production Ready
