# Sub-Checklist Template [HARD]

Each sub file (`{order}_{agent-topic}.md`) must include the following sections:

```markdown
# {agent-topic}: {task title}
Status: [ ] | Owner: {agent} | Language: per DO_JOBS_LANGUAGE env var (default: en)

## Problem Summary
- What is being solved
- Why this work is needed

## Acceptance Criteria
> Mark `[ ]` → `[o]` on completion. `[x]` means "failure" — never use for completion.
- [ ] Measurable completion condition 1
- [ ] Measurable completion condition 2
- [ ] Verification complete (one of the following):
  - Test required: `path/to/file_test.go` written and passing
  - Test not required: verification method specified (build check, manual check, etc.)
- [ ] Committed

## Solution Approach
- Chosen approach
- Why this approach (alternatives considered and rejection reasons)

## Critical Files
- **Modify**: `path/to/file.go` -- reason for change
- **Reference**: `path/to/ref.go` -- reason for reference
- **Test**: `path/to/file_test.go`

## Risks
- What could break: (specifics)
- Cautions: (side effects, performance, compatibility)

## Progress Log
- 2026-02-11 14:00:00 [~] Work started: initial structure design
- 2026-02-11 15:30:22 [~] JWT token issuance logic implemented
- 2026-02-11 16:00:45 [*] Unit tests written and executed
- 2026-02-11 16:30:10 [o] All tests passed, committed (commit: a1b2c3d)

## FINAL STEP: Commit (never skip)
- [ ] `git add` — stage only changed files
- [ ] `git diff --cached` — verify only intended changes included
- [ ] `git commit` — include WHY in commit message
- [ ] Record commit hash in Progress Log
⚠️ Work is incomplete if this section is not completed

## Lessons Learned (write on completion)
- What went well:
- What was difficult:
- Improvement actions: (concrete actions applied to rules/code/process + commit hash)
```

## Template Mandatory Rules
- [HARD] Problem Summary, Acceptance Criteria, Critical Files must be written before starting work
- [HARD] Acceptance Criteria must specify verification method — test file path or alternative verification
- [HARD] Agent workflow is: write code → verify (test/build) → pass → commit — never just write code and stop
- [HARD] Must record commit hash in Progress Log after commit — e.g., `[o] Done (commit: a1b2c3d)`
- [HARD] No `[o]` completion without commit hash — the commit is proof of completion
- [HARD] Solution Approach written at implementation start (mention at least 1 alternative)
- [HARD] Progress Log records **what was done**, not just state changes (track work content)
- [HARD] Progress Log timestamps must include seconds: `YYYY-MM-DD HH:MM:SS` — distinguishable when multiple state changes in same minute
- [HARD] No multiple state transitions at same timestamp — each step reflects actual work timing
- [HARD] Lessons Learned must be written on `[o]` completion — empty is forbidden
- [HARD] Improvement actions are **actions**, not observations: rule changes, code improvements, memory records — with commit hash
- [HARD] "Will do next time" is forbidden — improve now, or specify why improvement is not possible
- Risks: write "No risks identified" if none found
