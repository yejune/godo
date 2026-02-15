# hook-integration: 훅에 PersonaLoader 연결 (Phase 3)
상태: [ ] | 담당: expert-backend

## Problem Summary
- 현재 3개 훅 파일(session_start, post_tool_use, user_prompt_submit)과 spinner.go가 `buildPersona()` 하드코딩 함수를 직접 호출함
- Phase 2에서 구현한 PersonaLoader를 이 4개 파일에 연결해야 함
- Architecture의 Approach C에 따라 파일 기반 로딩을 시도하고, 실패 시 기존 `buildPersona()`로 폴백하는 구조가 필요
- 기존 동작과 정확히 동일한 출력을 보장해야 함 (행동 변경 없이 내부 구현만 교체)

## Acceptance Criteria
- [ ] hook_session_start.go: LoadCharacter + LoadSpinner 사용, 폴백 시 buildPersona() + 하드코딩 spinners
- [ ] hook_post_tool_use.go: LoadCharacter → BuildReminder 사용, 폴백 시 buildPersona()
- [ ] hook_user_prompt_submit.go: LoadCharacter → BuildReminder 사용, 폴백 시 buildPersona()
- [ ] spinner.go: LoadSpinner → BuildSpinnerVerbs 사용, 폴백 시 하드코딩 배열
- [ ] 기존 테스트 전체 통과 (회귀 없음)
- [ ] 파일 기반 로딩 시 기존과 동일한 hook 출력 생성
- [ ] 파일 미존재 시 stderr 경고 출력: `[godo] persona file not found, using built-in defaults`
- [ ] 커밋 완료

## Solution Approach
- 각 훅 파일에서 `buildPersona()` 호출부를 찾아 LoadCharacter() 호출로 교체
- 에러 처리 패턴: `pd, err := LoadCharacter(dir, typ); if err != nil { /* fallback to buildPersona() */ }`
- spinner.go에서 `getPersonaSpinnerVerbs()` 내부를 LoadSpinner() 호출로 교체
- 기존 buildPersona() 함수와 하드코딩 spinner 배열은 삭제하지 않음 (폴백용으로 유지)
- 대안 고려: buildPersona()를 제거하고 파일만 사용 → 기각 (하위 호환성 보장을 위해 폴백 유지)

## Test Strategy
- unit test: 기존 hook 테스트가 모두 통과하는지 확인 (회귀 테스트)
- integration: 파일이 있을 때 파일 기반 출력, 파일이 없을 때 폴백 출력이 동일한지 비교
- pass (수동 확인): DO_PERSONA 변경 후 각 훅의 출력 메시지가 올바른 호칭/말투인지 확인

## Critical Files

### 항목 #7: hook_session_start.go에 PersonaLoader 연결
- **수정 대상**: `cmd/godo/hook_session_start.go` — buildPersona() → LoadCharacter(), getPersonaSpinnerVerbs() → LoadSpinner()
- **참조 파일**: `cmd/godo/persona_loader.go` — LoadCharacter, LoadSpinner API
- **참조 파일**: `architecture-persona-system.md` Section 6.2 — 변경 명세

### 항목 #8: hook_post_tool_use.go에 PersonaLoader 연결
- **수정 대상**: `cmd/godo/hook_post_tool_use.go` — buildPersona() → LoadCharacter() + BuildReminder()
- **참조 파일**: `cmd/godo/persona_loader.go` — BuildReminder API
- **참조 파일**: `architecture-persona-system.md` Section 6.3 — 변경 명세

### 항목 #9: hook_user_prompt_submit.go에 PersonaLoader 연결
- **수정 대상**: `cmd/godo/hook_user_prompt_submit.go` — buildPersona() → LoadCharacter() + BuildReminder()
- **참조 파일**: `cmd/godo/persona_loader.go` — BuildReminder API
- **참조 파일**: `architecture-persona-system.md` Section 6.4 — 변경 명세

### 항목 #10: spinner.go에 LoadSpinner 연결
- **수정 대상**: `cmd/godo/spinner.go` — getPersonaSpinnerVerbs() 내부를 LoadSpinner() + BuildSpinnerVerbs()로 교체
- **참조 파일**: `cmd/godo/persona_loader.go` — LoadSpinner, BuildSpinnerVerbs API
- **참조 파일**: `architecture-persona-system.md` Section 6.5 — 변경 명세

## Risks
- hook 출력 형식이 미묘하게 달라질 수 있음: BuildReminder() 출력과 기존 하드코딩 문자열의 공백/줄바꿈 차이 주의
- session_start는 systemMessage를, post_tool_use/user_prompt_submit는 additionalContext를 출력 — 혼동 금지
- spinner.go의 suffix 조합 로직이 기존과 다르면 UI 스피너 메시지가 변경됨 — 기존 출력과 1:1 비교 필수
- 4개 파일을 한 에이전트가 수정하므로 3파일 제한에 근접 — 필요 시 #7/#8을 한 커밋, #9/#10을 다른 커밋으로 분리

## Progress Log
- (작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — hook_session_start.go, hook_post_tool_use.go, hook_user_prompt_submit.go, spinner.go만 스테이징
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 다음에 다르게 할 점:
