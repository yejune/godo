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
> **[HARD] 대용량 파일(architecture.md, handoff.md 등) 읽기 금지 — 핵심 내용을 이 섹션에 직접 기술**

```
Task: {간결한 작업 요약 - 한 줄}

Context:
- Current state: {[~] 또는 [ ] + 마지막 Progress Log 내용}
- What's done: {완료된 작업 요약}
- What's remaining: {남은 작업}

Implementation (핵심 내용 직접 기술, 파일 참조 금지):
{실제 구현해야 할 코드/로직을 여기에 직접 작성}
{architecture.md 등에서 필요한 부분만 발췌해서 여기에}

Files to modify (최대 1-2개):
1. {파일 경로}: {수정 내용}
2. {파일 경로}: {수정 내용}

Test command:
{go test ./path/to/package -v 또는 equivalent}

Commit when done:
git add {files} && git commit -m "{type}: {description}"
```

### Prompt 작성 규칙 [HARD]
- [HARD] **대용량 파일 읽기 금지**: architecture.md(500줄+), handoff.md 등은 참조만 하고 내용을 prompt에 직접 기술
- [HARD] **핵심만 발췌**: 전체 파일이 아닌, 해당 작업에 필요한 섹션만 추출해서 prompt에 포함
- [HARD] **소스 파일 최대 2개**: 한 에이전트가 수정할 파일은 1-2개로 제한
- [HARD] **커밋 명령어 포함**: 정확한 `git add` 대상과 커밋 메시지까지 명시
- [HARD] **테스트 명령어 포함**: 검증 방법을 구체적으로 명시

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
