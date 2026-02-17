# CLI 통합: cmd/godo 리네임 + 명령 확장 + 모델 확장
상태: [ ] | 담당: expert-backend (에이전트 H) | 작성 언어: ko

## Problem Summary
- cmd/convert/ → cmd/godo/ 리네임하여 godo가 최종 바이너리명
- Phase 3에서 생성된 패키지들을 CLI 명령으로 노출 (hook, mode, lint, create 등)
- model/ 패키지에 새 에러 타입 + PersonaManifest에 Workflows 필드 추가

## Acceptance Criteria
> 완료 시 [ ] → [o]로 변경. [x]는 "실패"를 의미하므로 절대 사용 금지.
- [ ] cmd/convert/ → cmd/godo/ 리네임 완료
- [ ] internal/cli/hook.go 생성 — `godo hook <event>` 명령 라우팅
- [ ] internal/cli/mode.go 생성 — `godo mode [get|set]` 명령
- [ ] internal/cli/lint.go 생성 — `godo lint [--all|setup]` 명령
- [ ] internal/cli/create.go 생성 — `godo create <agent|skill>` 명령
- [ ] internal/cli/claude.go 생성 — `godo claude [--profile]` 명령
- [ ] internal/cli/spinner.go, statusline.go, rank.go, glm.go 생성
- [ ] model/errors.go 확장 — ErrHook, ErrMode, ErrLint 추가
- [ ] model/persona_manifest.go 확장 — Workflows []string 필드 추가
- [ ] `go build ./cmd/godo/` 성공
- [ ] `./godo --help` 출력에 모든 명령 표시 확인
- [ ] 커밋 완료

## Solution Approach
- cmd/convert/ → cmd/godo/ 디렉토리 리네임 (git mv)
- main.go에서 cli.Execute() 호출은 유지, 바이너리명만 변경
- internal/cli/root.go에 새 서브커맨드 등록:
  - hookCmd → internal/hook/Dispatcher.Dispatch() 호출
  - modeCmd → internal/mode/Manager 호출
  - lintCmd → internal/lint/Runner 호출
  - createCmd → internal/scaffold/Create() 호출
  - claudeCmd → internal/profile/Launch() 호출
  - spinnerCmd, statuslineCmd, rankCmd, glmCmd
- model/errors.go에 ErrHook, ErrMode, ErrLint 에러 타입 추가 (기존 패턴 따름)
- model/persona_manifest.go에 Workflows 필드 추가 (assembler가 workflows/ 인식하도록)
- 대안: 별도 cmd/godo/main.go를 새로 작성 → 기각 (기존 CLI 인프라 재사용이 효율적)

## Critical Files
- **리네임**: `cmd/convert/` → `cmd/godo/`
- **수정 대상**:
  - `internal/cli/root.go` — 새 서브커맨드 등록
  - `internal/model/errors.go` — 새 에러 타입
  - `internal/model/persona_manifest.go` — Workflows 필드
- **생성 대상**:
  - `internal/cli/hook.go`, `internal/cli/mode.go`, `internal/cli/lint.go`
  - `internal/cli/create.go`, `internal/cli/claude.go`
  - `internal/cli/spinner.go`, `internal/cli/statusline.go`
  - `internal/cli/rank.go`, `internal/cli/glm.go`
- **참조**: `architecture.md` 섹션 5.1 CLI 라우팅

## Risks
- cmd/ 리네임 시 go.mod, Makefile, CI 설정 등에 경로 참조 변경 필요
- 기존 extract/assemble 명령이 리네임 후에도 정상 동작 확인 필수
- assembler가 Workflows 필드를 올바르게 처리하는지 검증 필요

## Progress Log
(작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (cmd/godo/, cli/*.go, model/*.go)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 개선 조치:
