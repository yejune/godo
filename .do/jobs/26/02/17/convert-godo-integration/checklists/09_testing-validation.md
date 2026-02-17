# 테스트 + 검증: 단위/통합/E2E 테스트 + 전체 파이프라인 검증
상태: [ ] | 담당: expert-testing (에이전트 I) | 작성 언어: ko

## Problem Summary
- Phase 3-4에서 생성된 모든 새 Go 패키지에 대해 단위 테스트 작성
- assembler E2E 테스트를 확장하여 코어 모듈 주입 후 assemble 정상 동작 검증
- CLI 통합 테스트로 전체 명령 라우팅 검증
- 최종: 전체 테스트 스위트 통과 + 커버리지 85%+ 확인

## Acceptance Criteria
> 완료 시 [ ] → [o]로 변경. [x]는 "실패"를 의미하므로 절대 사용 금지.
- [ ] internal/hook/dispatcher_test.go + types_test.go 작성 — 이벤트 라우팅, JSON 직렬화 검증
- [ ] internal/hook/security_test.go 작성 — 보안 패턴 매칭, deny 규칙 검증
- [ ] internal/mode/mode_test.go 작성 — 실행/권한 모드 get/set, 파일 영속성
- [ ] internal/persona/loader_test.go 작성 — 캐릭터 YAML 로드, 리마인더 생성
- [ ] internal/assembler/ E2E 테스트 확장 — 코어 모듈 주입 후 assemble 정상 동작
- [ ] internal/cli/ 통합 테스트 — 전체 명령 라우팅 + 실행 검증
- [ ] `go test ./...` 전체 통과
- [ ] assemble 출력 검증 — 미해결 슬롯, 오버라이드 잔재 없는지 확인
- [ ] 커밋 완료

## Solution Approach
- 테스트 전략 (architecture.md 섹션 8 기준):
  - **Unit**: hook (JSON round-trip, 패턴 매칭, 라우팅), mode (파일 I/O), persona (YAML 파싱)
  - **Integration**: assembler E2E (testdata/ 사용), CLI (임시 디렉토리)
- hook 테스트:
  - dispatcher_test.go: Register + Dispatch, 알 수 없는 이벤트 → 에러
  - types_test.go: Input/Output JSON 직렬화/역직렬화 round-trip
  - security_test.go: deny 패턴 매칭 (git push --force 등)
- mode 테스트:
  - mode_test.go: Get/Set + 파일 영속성 (임시 디렉토리 사용)
- persona 테스트:
  - loader_test.go: testdata/ 캐릭터 YAML로 로드 + 리마인더 빌드
- assembler E2E:
  - 기존 e2e_test.go 확장: 코어에 새 모듈이 추가된 상태에서 assemble 실행 → 출력 검증
- CLI 통합:
  - integration_test.go: 각 서브커맨드 실행 → exit code + 출력 확인
- 대안: mock 기반 테스트 우선 → 기각 (Real DB 원칙은 아니지만 파일 I/O도 실제 경로로)

## Critical Files
- **생성 대상**:
  - `internal/hook/dispatcher_test.go`
  - `internal/hook/types_test.go`
  - `internal/hook/security_test.go`
  - `internal/mode/mode_test.go`
  - `internal/persona/loader_test.go`
  - `internal/cli/integration_test.go` (신규)
- **수정 대상**:
  - `internal/assembler/` E2E 테스트 파일 (기존 확장)
- **참조**: `testdata/` 디렉토리 (테스트 데이터)

## Risks
- assembler E2E 테스트가 코어 디렉토리 구조에 의존 — 코어 변경(Phase 1) 후 testdata도 갱신 필요
- CLI 통합 테스트에서 바이너리 빌드 필요 — go build 후 실행하는 패턴
- 커버리지 85% 미달 시 — 추가 테스트 작성 필요 (lint, scaffold 등 보조 패키지)

## Progress Log
(작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (*_test.go 파일들)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 개선 조치:
