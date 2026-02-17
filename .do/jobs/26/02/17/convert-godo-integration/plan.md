# Convert-Godo 통합 플랜

**작성 일시**: 2026-02-17
**대상 프로젝트**: convert (~/Work/new/convert/)
**기반 문서**: architecture.md, handoff.md

---

## 작업 개요

godo CLI(do-focus/cmd/godo/, 41파일)를 convert 프로젝트에 흡수 병합하고, do-focus의 dev-*.md 규칙 파일들을 토픽별 분해하여 코어 스킬 모듈에 주입하며, DO 페르소나 워크플로우를 5개의 씬(thin) 파일로 재구성한다.

**최종 목표**: convert가 단일 바이너리로 extract + assemble + 모든 godo 기능을 제공. do-focus는 레거시가 된다.

---

## Phase별 작업 계획

### Phase 1: 코어 스킬 콘텐츠 주입 (dev-*.md 분해)

**목적**: do-focus의 dev-*.md 규칙 파일들을 토픽별로 분해하여 코어 스킬의 modules/ 또는 references/에 새 .md 파일로 주입

**작업 범위**: 마크다운 파일 21개 생성 (Go 코드 변경 없음)

**에이전트 배정**: 3개 general-purpose 에이전트 병렬 실행
- **에이전트 A** (foundation-core 담당): 3개 모듈 생성
  - agent-execution-cycle.md, agent-delegation.md, agent-research.md
- **에이전트 B** (foundation-quality + foundation-context 담당): 7개 모듈 생성
  - read-before-write.md, coding-discipline.md, commit-discipline.md, parallel-agent-isolation.md, syntax-check.md
  - file-reading-optimization.md, knowledge-management.md
- **에이전트 C** (workflow-spec + workflow-testing + workflow-project + workflow-tdd 담당): 11개 모듈 생성
  - checklist-system.md, checklist-templates.md, analysis-template.md, architecture-template.md, report-template.md, complexity-check.md
  - testing-rules.md, bug-fix-workflow.md
  - docker-rules.md, ai-forbidden-patterns.md
  - tdd-cycle.md

**예상 변경 파일 수**: 21개 (신규 생성)

**소스 파일** (do-focus에서 읽을 것):
- `~/Work/do-focus.workspace/do-focus/.claude/rules/do/development/workflow.md`
- `~/Work/do-focus.workspace/do-focus/.claude/rules/do/development/testing.md`
- `~/Work/do-focus.workspace/do-focus/.claude/rules/do/development/checklist.md`
- `~/Work/do-focus.workspace/do-focus/.claude/rules/do/development/environment.md`
- `~/Work/do-focus.workspace/do-focus/.claude/rules/do/workflow/file-reading-optimization.md`

**검증**: 각 모듈 파일 생성 후 원본 대비 내용 누락 없는지 줄 수 비교

---

### Phase 2: 페르소나 워크플로우 생성 + 오버라이드 스킬 제거

**목적**: DO 페르소나의 5개 씬 워크플로우 파일 생성, 오버라이드 스킬 삭제, manifest.yaml 갱신

**작업 범위**: 파일 5개 생성 + 9개 디렉토리 삭제 + manifest.yaml/bootapp.md 갱신

**에이전트 배정**: 2개 general-purpose 에이전트
- **에이전트 D** (워크플로우 생성): 5개 워크플로우 파일 + bootapp.md 갱신
  - workflows/plan.md, run.md, report.md, team-plan.md, team-run.md
  - rules/bootapp.md 갱신
- **에이전트 E** (오버라이드 정리): 오버라이드 스킬 삭제 + manifest.yaml 갱신 + rules/workflow/ 삭제
  - 9개 오버라이드 스킬 디렉토리 삭제 (do-foundation-core, do-foundation-quality, do-workflow-ddd, do-workflow-plan, do-workflow-project, do-workflow-spec, do-workflow-tdd, do-workflow-team, do-workflow-testing)
  - manifest.yaml 갱신 (workflows 추가, 오버라이드 스킬 제거, rules 축소)
  - rules/workflow/ 삭제

**의존성**: Phase 1 완료 후 실행 (코어 주입 완료 확인 후 오버라이드 제거)

**예상 변경 파일 수**: ~15개 (생성 5 + 삭제 9디렉토리 + 수정 2)

**검증**: assemble 빌드 검증 (미해결 슬롯 없는지)

---

### Phase 3: godo Go 코드 패키지 분리

**목적**: godo의 `package main` 41파일을 독립 패키지로 분리하여 convert/internal/에 배치

**작업 범위**: 10개 새 패키지 생성 (hook, mode, persona, lint, scaffold, profile, statusline, rank, glm, cli 확장)

**에이전트 배정**: 2개 expert-backend 에이전트 순차 실행
- **에이전트 F** (핵심 패키지): hook/ + mode/ + persona/ 패키지 (가장 복잡, 의존성 많음)
  - internal/hook/ (12파일): types.go, dispatcher.go, session_start.go ~ session_end.go, security.go, job_state.go
  - internal/mode/ (1파일): mode.go
  - internal/persona/ (2파일): loader.go, spinner.go
- **에이전트 G** (보조 패키지): lint/ + scaffold/ + profile/ + statusline/ + rank/ + glm/
  - internal/lint/ (3파일): runner.go, gate.go, setup.go
  - internal/scaffold/ (1파일): create.go
  - internal/profile/ (1파일): profile.go
  - internal/statusline/ (1파일): statusline.go
  - internal/rank/ (4파일): auth.go, client.go, config.go, transcript.go
  - internal/glm/ (1파일): glm.go

**의존성**: Phase 1, 2 완료 불필요 (독립 실행 가능). 다만 Phase 4와 순차.

**예상 변경 파일 수**: ~26개 (신규 Go 파일)

**검증**: `go build ./...` 통과

---

### Phase 4: CLI 통합

**목적**: cmd/convert/ → cmd/godo/ 리네임, CLI 명령 확장, moai_sync 삭제

**에이전트 배정**: 1개 expert-backend 에이전트
- **에이전트 H** (CLI 통합):
  - cmd/convert/ → cmd/godo/ 리네임
  - internal/cli/ 확장 (hook, mode, lint, create, claude, spinner, statusline, rank, glm 명령)
  - moai_sync*.go 관련 코드 삭제 (해당 파일이 convert에 있을 경우)
  - model/errors.go 확장 (ErrHook, ErrMode, ErrLint 추가)
  - model/persona_manifest.go 확장 (Workflows 필드 추가)

**의존성**: Phase 3 완료 후 실행 (패키지가 존재해야 CLI에서 import 가능)

**예상 변경 파일 수**: ~15개 (리네임 + 신규 CLI 파일 + 모델 수정)

**검증**: `go build ./cmd/godo/` 성공 + `godo --help` 출력 확인

---

### Phase 5: 테스트 + 검증

**목적**: 모든 새 패키지 단위 테스트 + 통합 테스트 + 전체 파이프라인 검증

**에이전트 배정**: 1개 expert-testing 에이전트
- **에이전트 I** (테스트):
  - internal/hook/ 단위 테스트 (dispatcher_test.go, types_test.go, security_test.go)
  - internal/mode/ 단위 테스트 (mode_test.go)
  - internal/persona/ 단위 테스트 (loader_test.go)
  - assembler E2E 테스트 확장 (코어 모듈 주입 후 assemble 정상 동작)
  - CLI 통합 테스트 (cli/integration_test.go)
  - 전체 테스트 스위트 실행

**의존성**: Phase 4 완료 후 실행

**예상 변경 파일 수**: ~10개 (테스트 파일)

**검증**: `go test ./...` 전체 통과, 커버리지 85%+

---

## 의존성 관계

```
Phase 1 (코어 스킬 주입)
    │
    ▼
Phase 2 (페르소나 워크플로우 + 오버라이드 제거)  ←── Phase 1 완료 필수
    │
    │   Phase 3 (Go 패키지 분리)  ←── Phase 1, 2와 독립 (병렬 가능)
    │       │
    │       ▼
    │   Phase 4 (CLI 통합)  ←── Phase 3 완료 필수
    │       │
    ▼       ▼
Phase 5 (테스트 + 검증)  ←── Phase 2, 4 모두 완료 필수
```

**병렬 실행 가능 구간**:
- Phase 1 내부: 에이전트 A, B, C 병렬
- Phase 2: 에이전트 D, E 병렬 (Phase 1 완료 후)
- Phase 3: Phase 2와 병렬 가능 (Go 코드는 마크다운과 독립)

---

## 에이전트 배정 요약

| 에이전트 | 타입 | Phase | 담당 | 서브 체크리스트 |
|---------|------|-------|------|--------------|
| A | general-purpose | 1 | foundation-core 모듈 3개 | 01_core-foundation-core.md |
| B | general-purpose | 1 | foundation-quality + context 모듈 7개 | 02_core-foundation-quality-context.md |
| C | general-purpose | 1 | workflow-* 모듈 11개 | 03_core-workflow-modules.md |
| D | general-purpose | 2 | 워크플로우 5개 + bootapp | 04_persona-workflows.md |
| E | general-purpose | 2 | 오버라이드 삭제 + manifest | 05_persona-cleanup.md |
| F | expert-backend | 3 | hook/mode/persona 패키지 | 06_go-core-packages.md |
| G | expert-backend | 3 | lint/scaffold/profile/etc 패키지 | 07_go-aux-packages.md |
| H | expert-backend | 4 | CLI 통합 + 리네임 | 08_cli-integration.md |
| I | expert-testing | 5 | 전체 테스트 | 09_testing-validation.md |

---

## 예상 총 변경 파일 수

| Phase | 신규 | 수정 | 삭제 | 합계 |
|-------|------|------|------|------|
| 1 | 21 | 0 | 0 | 21 |
| 2 | 5 | 2 | ~40 (9 디렉토리) | ~47 |
| 3 | 26 | 0 | 0 | 26 |
| 4 | 10 | 5 | ~5 | ~20 |
| 5 | 10 | 2 | 0 | 12 |
| **합계** | **72** | **9** | **~45** | **~126** |

---

## 리스크 + 완화책

| 리스크 | 영향도 | 완화책 |
|--------|--------|--------|
| dev-*.md 분해 시 규칙 누락 | HIGH | 분해 전/후 줄 수 비교, 원본 대비 diff 검증 |
| 오버라이드 스킬 삭제 후 기능 누락 | HIGH | 삭제 전 오버라이드 내용이 코어에 100% 반영 확인 |
| godo 패키지 분리 시 내부 참조 깨짐 | MEDIUM | package main → package hook 전환 후 go build 즉시 검증 |
| assembler가 workflows/ 인식 실패 | MEDIUM | PersonaManifest에 Workflows 필드 추가 필요 |
| do-ko, moai-ko 동기화 누락 | MEDIUM | Phase 2 후 -ko 변형도 동일 구조 갱신 (별도 후속 작업) |
| moai_sync 삭제로 기존 사용자 영향 | MEDIUM | convert 동일 기능 확인 후 삭제, CHANGELOG 안내 |
| 에이전트 토큰 소진으로 중단 | LOW | 서브 체크리스트 기반 멱등 재개 — [o] 건너뛰고 미완료부터 |

---

**작성자**: Focus 오케스트레이터
**다음 단계**: checklist.md + checklists/*.md 작성
