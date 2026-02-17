# Go 핵심 패키지 분리: hook/mode/persona
상태: [ ] | 담당: expert-backend (에이전트 F) | 작성 언어: ko

## Problem Summary
- godo의 `package main` 41파일 중 hook(12파일), mode(1파일), persona(2파일) 관련 코드를 독립 패키지로 분리
- hook 시스템이 가장 복잡 (8개 이벤트 핸들러 + 디스패처 + 보안 정책 + job 상태 추적)
- convert/internal/ 하위에 hook/, mode/, persona/ 패키지 생성

## Acceptance Criteria
> 완료 시 [ ] → [o]로 변경. [x]는 "실패"를 의미하므로 절대 사용 금지.
- [ ] internal/hook/types.go 생성 — Input, Output, HookSpecificOutput, EventType, Decision 상수
- [ ] internal/hook/dispatcher.go 생성 — Dispatcher 구조체, Register, Dispatch 메서드
- [ ] internal/hook/session_start.go 생성 — SessionStart 핸들러 (페르소나 리마인더 주입)
- [ ] internal/hook/pre_tool.go 생성 — PreToolUse 핸들러 (보안 정책 + 모드 확인)
- [ ] internal/hook/post_tool_use.go 생성 — PostToolUse 핸들러
- [ ] internal/hook/security.go 생성 — 보안 패턴 deny 규칙
- [ ] internal/hook/ 나머지 핸들러 생성 — user_prompt.go, subagent_stop.go, stop.go, compact.go, session_end.go, job_state.go
- [ ] internal/mode/mode.go 생성 — ExecutionMode, PermissionMode, Manager
- [ ] internal/persona/loader.go 생성 — Character, Loader, BuildReminder
- [ ] internal/persona/spinner.go 생성 — 한국어 스피너 동사
- [ ] `go build ./...` 통과
- [ ] 커밋 완료

## Solution Approach
- godo 소스(do-focus/cmd/godo/)에서 해당 파일들을 읽고, `package main` → 독립 패키지로 변환
- hook 패키지:
  - moai_hook_types.go → types.go (Input/Output 구조체)
  - hook.go + moai_hook_contract.go → dispatcher.go (디스패처 + 팩토리)
  - hook_session_start.go → session_start.go (각 핸들러 1파일)
  - security_patterns.go → security.go
  - job_state.go → job_state.go
- mode 패키지:
  - mode.go → mode.go (ExecutionMode/PermissionMode + Manager)
- persona 패키지:
  - persona_loader.go → loader.go
  - spinner.go → spinner.go
- 내부 참조: hook → mode (PreToolUse에서 현재 모드 확인), hook → persona (SessionStart에서 리마인더)
- 대안: hook만 먼저 분리하고 mode/persona는 다음 단계 → 기각 (hook이 mode/persona를 import하므로 동시 생성 필요)

## Critical Files
- **소스** (do-focus/cmd/godo/):
  - `hook.go`, `moai_hook_types.go`, `moai_hook_contract.go`
  - `hook_session_start.go`, `hook_pre_tool.go`, `hook_post_tool_use.go`
  - `hook_user_prompt.go`, `hook_subagent_stop.go`, `hook_stop.go`
  - `hook_compact.go`, `hook_session_end.go`
  - `security_patterns.go`, `job_state.go`
  - `mode.go`, `persona_loader.go`, `spinner.go`
- **생성 대상** (convert/internal/):
  - `hook/types.go`, `hook/dispatcher.go`, `hook/session_start.go`, `hook/pre_tool.go`
  - `hook/post_tool_use.go`, `hook/security.go`, `hook/user_prompt.go`
  - `hook/subagent_stop.go`, `hook/stop.go`, `hook/compact.go`
  - `hook/session_end.go`, `hook/job_state.go`
  - `mode/mode.go`
  - `persona/loader.go`, `persona/spinner.go`

## Risks
- godo 코드가 전역 변수/함수에 의존할 수 있음 — 패키지 분리 시 의존성 주입으로 전환
- hook ↔ mode ↔ persona 순환 참조 방지 — 인터페이스로 의존성 역전
- go.mod에 새 의존성 불필요 확인 (yaml.v3만 사용)

## Progress Log
(작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (hook/*.go, mode/*.go, persona/*.go)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 개선 조치:
