# Sub-Checklist Template [HARD]

Each sub file (`{order}_{agent-topic}.md`) must include the following sections:

```markdown
# {agent-topic}: {task title}
Status: [ ] | Owner: {agent} | Language: per DO_JOBS_LANGUAGE env var (default: en)

**State Legend**: [ ] pending | [~] in progress | [*] testing | [!] blocked | [o] done | [x] failed

> **Token Budget Warning**: Check file size before reading — never read 500+ line files in full (use Grep for relevant sections only). At ~10% tokens remaining, record current state in checklist and report to super agent.

## Prompt (RESTART INSTRUCTION)
> **이 섹션은 재실행 시 에이전트가 정확히 무엇을 해야 하는지 명시합니다.**
> 새 에이전트가 이 체크리스트를 받으면 이 프롬프트만 읽고도 작업을 이어갈 수 있어야 합니다.

```
Task: {간결한 작업 요약}

Context:
- Previous work: {이전에 완료된 작업}
- Current state: {현재 상태 - [~], [*], 등}
- Remaining: {남은 작업}

Exact commands to run:
1. {정확한 명령어 또는 작업 단계}
2. ...

Files to modify:
- {파일 경로}: {수정 내용}

Verification:
- Run: {테스트 명령어}
- Expected: {예상 결과}

Commit when done:
git add {files} && git commit -m "..."
```

## Agent Instructions
> This section contains the exact orchestrator prompt for this task.
> A new agent receiving this sub-checklist should execute based on these instructions alone.

{Orchestrator writes the exact task prompt here — including:
- What to do (task description)
- Exact content to write/modify (verbatim, not file references)
- File paths to modify
- Verification commands
- Git commit command}

## Problem Summary
- What is being solved
- Why this work is needed

## Acceptance Criteria
> Mark `[ ]` → `[o]` on completion. `[x]` means "failure" — never use for completion.
- [ ] Measurable completion condition 1
- [ ] Measurable completion condition 2
- [ ] **TEST REQUIRED**: `go test ./...` or `pytest` must pass (NO EXCEPTIONS for code changes)
- [ ] Verification complete (one of the following):
  - Test required: `path/to/file_test.go` written and passing
  - Test not required: verification method specified (build check, manual check, etc.)
- [ ] Committed (commit: {hash}, files: {modified file list})

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
- [ ] `go test ./...` or equivalent — ALL tests must pass before commit
- [ ] `git add` — stage only changed files
- [ ] `git diff --cached` — verify only intended changes included
- [ ] `git commit` — include WHY in commit message
- [ ] Record commit hash in Progress Log
⚠️ Work is incomplete if this section is not completed
⚠️ NEVER commit without running tests first

## Lessons Learned (write on completion)
- What went well:
- What was difficult:
- Improvement actions: (concrete actions applied to rules/code/process + commit hash)
```

## Template Mandatory Rules
- [HARD] **Prompt 섹션 필수**: 재실행 시 정확히 무엇을 해야 하는지 명시 — 새 에이전트가 이것만 읽고 작업 가능해야 함
- [HARD] Problem Summary, Acceptance Criteria, Critical Files must be written before starting work
- [HARD] Acceptance Criteria must specify verification method — test file path or alternative verification
- [HARD] **테스트 필수**: 코드 변경 시 반드시 `go test ./...` 또는 equivalent 실행 — 테스트 없이 커밋 금지
- [HARD] Agent workflow is: write code → **run tests** → pass → commit — never just write code and stop
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
