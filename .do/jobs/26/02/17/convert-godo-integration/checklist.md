# Convert-Godo 통합 체크리스트

**작성 일시**: 2026-02-17
**플랜 기반**: plan.md
**상태 범례**: [ ] 미시작 | [~] 진행중 | [*] 테스트중 | [!] 블로커 | [o] 완료 | [x] 실패

---

## Phase 1: 코어 스킬 콘텐츠 주입 (병렬 3개 에이전트)

### 에이전트 A: foundation-core 모듈 (general-purpose)
→ 서브 체크리스트: `checklists/01_core-foundation-core.md`

- [ ] #1 agent-execution-cycle.md 생성 (검증: 원본 줄 수 비교)
- [ ] #2 agent-delegation.md 생성 (검증: 원본 줄 수 비교)
- [ ] #3 agent-research.md 생성 (검증: 원본 줄 수 비교)

### 에이전트 B: foundation-quality + foundation-context 모듈 (general-purpose)
→ 서브 체크리스트: `checklists/02_core-foundation-quality-context.md`

- [ ] #4 read-before-write.md 생성 (검증: 원본 줄 수 비교)
- [ ] #5 coding-discipline.md 생성 (검증: 원본 줄 수 비교)
- [ ] #6 commit-discipline.md 생성 (검증: 원본 줄 수 비교)
- [ ] #7 parallel-agent-isolation.md 생성 (검증: 원본 줄 수 비교)
- [ ] #8 syntax-check.md 생성 (검증: 원본 줄 수 비교)
- [ ] #9 file-reading-optimization.md 생성 (검증: 원본 줄 수 비교)
- [ ] #10 knowledge-management.md 생성 (검증: 원본 줄 수 비교)

### 에이전트 C: workflow-* 모듈 (general-purpose)
→ 서브 체크리스트: `checklists/03_core-workflow-modules.md`

- [ ] #11 checklist-system.md 생성 (검증: 원본 줄 수 비교)
- [ ] #12 checklist-templates.md 생성 (검증: 원본 줄 수 비교)
- [ ] #13 analysis-template.md 생성 (검증: 원본 줄 수 비교)
- [ ] #14 architecture-template.md 생성 (검증: 원본 줄 수 비교)
- [ ] #15 report-template.md 생성 (검증: 원본 줄 수 비교)
- [ ] #16 complexity-check.md 생성 (검증: 원본 줄 수 비교)
- [ ] #17 testing-rules.md 생성 (검증: 원본 줄 수 비교)
- [ ] #18 bug-fix-workflow.md 생성 (검증: 원본 줄 수 비교)
- [ ] #19 docker-rules.md 생성 (검증: 원본 줄 수 비교)
- [ ] #20 ai-forbidden-patterns.md 생성 (검증: 원본 줄 수 비교)
- [ ] #21 tdd-cycle.md 생성 (검증: 원본 줄 수 비교)

---

## Phase 2: 페르소나 워크플로우 생성 + 오버라이드 제거 (병렬 2개 에이전트)

> depends on: Phase 1 전체 완료 (#1~#21)

### 에이전트 D: 워크플로우 생성 (general-purpose)
→ 서브 체크리스트: `checklists/04_persona-workflows.md`

- [ ] #22 personas/do/workflows/ 디렉토리 생성 (검증: ls 확인)
- [ ] #23 workflows/plan.md 생성 (검증: 내용 확인)
- [ ] #24 workflows/run.md 생성 (검증: 내용 확인)
- [ ] #25 workflows/report.md 생성 (검증: 내용 확인)
- [ ] #26 workflows/team-plan.md 생성 (검증: 내용 확인)
- [ ] #27 workflows/team-run.md 생성 (검증: 내용 확인)
- [ ] #28 rules/bootapp.md 갱신 — dev-environment.md bootapp 섹션 반영 (검증: diff 확인)

### 에이전트 E: 오버라이드 정리 (general-purpose)
→ 서브 체크리스트: `checklists/05_persona-cleanup.md`

- [ ] #29 오버라이드 스킬 9개 디렉토리 삭제 (검증: ls 확인)
- [ ] #30 rules/workflow/ 디렉토리 삭제 (검증: ls 확인)
- [ ] #31 manifest.yaml 갱신 — workflows 추가, 오버라이드 스킬 제거 (검증: yaml 파싱 확인)

---

## Phase 3: godo Go 코드 패키지 분리 (순차 2개 에이전트)

> Phase 1, 2와 독립 실행 가능 (Go 코드 ↔ 마크다운 무관)

### 에이전트 F: 핵심 패키지 — hook/mode/persona (expert-backend)
→ 서브 체크리스트: `checklists/06_go-core-packages.md`

- [ ] #32 internal/hook/types.go 생성 — Input/Output/HookSpecificOutput 타입 (검증: go build)
- [ ] #33 internal/hook/dispatcher.go 생성 — Dispatcher + Register + Dispatch (검증: go build)
- [ ] #34 internal/hook/session_start.go 생성 — SessionStart 핸들러 (검증: go build)
- [ ] #35 internal/hook/pre_tool.go 생성 — PreToolUse 핸들러 (검증: go build)
- [ ] #36 internal/hook/post_tool_use.go 생성 — PostToolUse 핸들러 (검증: go build)
- [ ] #37 internal/hook/security.go 생성 — 보안 패턴 deny 규칙 (검증: go build)
- [ ] #38 internal/hook/ 나머지 핸들러 생성 — user_prompt, subagent_stop, stop, compact, session_end, job_state (검증: go build)
- [ ] #39 internal/mode/mode.go 생성 — ExecutionMode/PermissionMode + Manager (검증: go build)
- [ ] #40 internal/persona/loader.go + spinner.go 생성 — Character/Loader (검증: go build)

### 에이전트 G: 보조 패키지 — lint/scaffold/profile/statusline/rank/glm (expert-backend)
→ 서브 체크리스트: `checklists/07_go-aux-packages.md`

- [ ] #41 internal/lint/ 패키지 생성 — runner.go, gate.go, setup.go (검증: go build)
- [ ] #42 internal/scaffold/create.go 생성 — 에이전트/스킬 스캐폴딩 (검증: go build)
- [ ] #43 internal/profile/profile.go 생성 — Claude 프로파일 관리 (검증: go build)
- [ ] #44 internal/statusline/statusline.go 생성 — 상태 줄 렌더링 (검증: go build)
- [ ] #45 internal/rank/ 패키지 생성 — auth.go, client.go, config.go, transcript.go (검증: go build)
- [ ] #46 internal/glm/glm.go 생성 — GLM 백엔드 (검증: go build)

---

## Phase 4: CLI 통합 (순차 1개 에이전트)

> depends on: Phase 3 완료 (#32~#46)

### 에이전트 H: CLI 통합 + 리네임 (expert-backend)
→ 서브 체크리스트: `checklists/08_cli-integration.md`

- [ ] #47 cmd/convert/ → cmd/godo/ 리네임 (검증: go build ./cmd/godo/)
- [ ] #48 internal/cli/ 확장 — hook, mode 명령 추가 (검증: go build)
- [ ] #49 internal/cli/ 확장 — lint, create, claude 명령 추가 (검증: go build)
- [ ] #50 internal/cli/ 확장 — spinner, statusline, rank, glm 명령 추가 (검증: go build)
- [ ] #51 model/errors.go 확장 — ErrHook, ErrMode, ErrLint 추가 (검증: go build)
- [ ] #52 model/persona_manifest.go 확장 — Workflows 필드 추가 (검증: go build)
- [ ] #53 전체 빌드 + godo --help 출력 확인 (검증: go build ./cmd/godo/ && ./godo --help)

---

## Phase 5: 테스트 + 검증 (순차 1개 에이전트)

> depends on: Phase 2 완료 (#22~#31) + Phase 4 완료 (#47~#53)

### 에이전트 I: 전체 테스트 (expert-testing)
→ 서브 체크리스트: `checklists/09_testing-validation.md`

- [ ] #54 internal/hook/ 단위 테스트 — dispatcher_test.go, types_test.go (검증: go test ./internal/hook/)
- [ ] #55 internal/hook/security_test.go — 보안 패턴 매칭 테스트 (검증: go test ./internal/hook/)
- [ ] #56 internal/mode/mode_test.go — 모드 get/set 테스트 (검증: go test ./internal/mode/)
- [ ] #57 internal/persona/loader_test.go — 캐릭터 로드 테스트 (검증: go test ./internal/persona/)
- [ ] #58 assembler E2E 테스트 확장 — 코어 모듈 주입 후 assemble 정상 동작 (검증: go test ./internal/assembler/)
- [ ] #59 CLI 통합 테스트 — cli/integration_test.go (검증: go test ./internal/cli/)
- [ ] #60 전체 테스트 스위트 실행 + 커버리지 확인 (검증: go test ./... -cover)
- [ ] #61 assemble 출력 검증 — 미해결 슬롯, 오버라이드 잔재 없는지 (검증: assemble 실행 후 출력 검사)

---

## 진행 현황 요약

| Phase | 항목 수 | 완료 | 진행중 | 미시작 | 블로커 |
|-------|---------|------|--------|--------|--------|
| 1. 코어 주입 | 21 | 0 | 0 | 21 | 0 |
| 2. 페르소나 | 10 | 0 | 0 | 10 | 0 |
| 3. Go 패키지 | 15 | 0 | 0 | 15 | 0 |
| 4. CLI 통합 | 7 | 0 | 0 | 7 | 0 |
| 5. 테스트 | 8 | 0 | 0 | 8 | 0 |
| **합계** | **61** | **0** | **0** | **61** | **0** |
