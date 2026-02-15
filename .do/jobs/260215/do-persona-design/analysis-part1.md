# Do Persona 설계를 위한 MoAI vs Do-Focus 비교 분석 (Part 1)

**분석일시**: 2026-02-15
**대상**: moai persona (`/tmp/e2e4-extract/personas/moai/`) vs do-focus (`~/Work/do-focus.workspace/do-focus/`)
**목적**: do 페르소나의 오케스트레이터 스킬 설계를 위한 구조 분석

---

## 1. MoAI CLAUDE.md 섹션별 분석

**파일**: `/tmp/e2e4-extract/personas/moai/CLAUDE.md`
**총 줄 수**: 358줄 (Version 13.0.0)
**언어**: 영어

| # | 섹션명 | 줄 범위 | 핵심 내용 요약 |
|---|--------|---------|---------------|
| 1 | Core Identity | L3-25 | HARD 규칙 8개 선언. 오케스트레이터 원칙, 개발 안전장치 참조. "Agent delegation recommended" (SOFT) |
| 2 | Request Processing Pipeline | L28-67 | 4단계 파이프라인: Analyze → Route → Execute → Report. 핵심 스킬 3개 로드 참조 |
| 3 | Command Reference | L70-81 | `/moai` 단일 진입점. 서브커맨드: plan, run, sync, project, fix, loop, feedback. Allowed Tools 전체 나열 |
| 4 | Agent Catalog | L84-109 | Selection Decision Tree + 4카테고리 에이전트 목록 (Manager 8, Expert 8, Builder 3, Team 8). 상세는 rules 파일로 위임 |
| 5 | SPEC-Based Workflow | L113-133 | `/moai plan` → `/moai run` → `/moai sync` 3단계 파이프라인. Agent Chain 6단계 명시 |
| 6 | Quality Gates | L136-151 | TRUST 5 참조 + LSP Quality Gates (plan/run/sync 단계별 임계값) |
| 7 | Safe Development Protocol | L153-198 | HARD 규칙 4개: Approach-First, Multi-File Decomposition, Post-Implementation Review, Reproduction-First Bug Fix. Go 전용 가이드 포함 |
| 8 | User Interaction Architecture | L201-220 | AskUserQuestion 제약: MoAI만 사용 가능, 서브에이전트 불가. 올바른 패턴 5단계. 최대 4옵션 |
| 9 | Configuration Reference | L223-246 | `.moai/config/sections/` YAML 참조. Language Rules (응답/에이전트/코드 각각 다른 언어 정책) |
| 10 | Web Search Protocol | L248-263 | URL 검증 필수, 금지 패턴 3개, Sources 섹션 필수 |
| 11 | Error Handling | L265-282 | 에러 유형별 복구 전략 5가지. Resumable Agents (agentId 기반 재개) |
| 12 | MCP Servers & UltraThink | L284-293 | Sequential Thinking, Context7, Pencil, claude-in-chrome 4개 MCP 통합 |
| 13 | Progressive Disclosure System | L298-311 | 3단계 토큰 최적화: L1 Metadata(100tok), L2 Body(5K), L3 Bundled(on-demand). 67% 절감 |
| 14 | Parallel Execution Safeguards | L313-320 | 파일 충돌 방지, 에이전트 도구 요건, 루프 방지(3회 재시도), 플랫폼 호환성 |
| 15 | Agent Teams (Experimental) | L324-348 | 활성화 조건, 모드 선택(--team/--solo/auto), Team APIs, Hook Events |

### 핵심 설계 패턴

1. **참조 기반 경량화**: CLAUDE.md는 "요약 + @참조"로 구성. 상세 규칙은 `.claude/rules/moai/` 하위에 분산
2. **영어 단일 언어**: 모든 지시문이 영어. 사용자 응답만 conversation_language
3. **HARD/Recommendation 구분**: 필수 규칙은 [HARD] 태그, 권장은 일반 텍스트
4. **버전 관리**: 문서 말미에 Version, Last Updated, Language, Core Rule 메타데이터

---

## 2. MoAI 오케스트레이터 스킬 분석

### 2.1 SKILL.md 구조

**파일**: `/tmp/e2e4-extract/personas/moai/skills/moai/SKILL.md`
**총 줄 수**: 377줄 (Version 2.0.0)

#### Frontmatter (L1-14)
```yaml
name: moai
description: MoAI super agent - unified orchestrator for autonomous development
allowed-tools: Task AskUserQuestion TaskCreate TaskUpdate TaskList TaskGet Bash Read Write Edit Glob Grep
user-invocable: true
metadata:
  argument-hint: "[subcommand] [args] | \"natural language task\""
```

#### 본문 구조

| 섹션 | 줄 범위 | 역할 |
|------|---------|------|
| Pre-execution Context | L16-19 | `git status`, `git branch` 실행하여 컨텍스트 수집 |
| Essential Files | L21-23 | `@.moai/config/config.yaml` 로드 |
| Core Identity | L27-41 | 기본 원칙 6개 (위임, 직접구현 금지, AskUserQuestion MoAI only, 병렬실행, 언어감지, 태스크 추적) |
| Intent Router | L44-88 | **4단계 우선순위 라우팅**: (1) Explicit Subcommand → (2) SPEC-ID Detection → (3) NL Classification → (4) Default(AskUserQuestion) |
| Workflow Quick Reference | L92-158 | 8개 워크플로우 요약 (plan, run, sync, fix, loop, moai, project, feedback) + Team mode 참조 |
| Core Rules | L162-240 | Agent Delegation Mandate, User Interaction Architecture, Task Tracking, Completion Markers, Output Rules, Error Handling |
| Agent Catalog | L243-294 | Manager 7 + Expert 8 + Builder 3 + Team 8 + Decision Tree |
| Common Patterns | L297-314 | Parallel/Sequential/Resume/Context Propagation 패턴 |
| Additional Resources | L318-337 | 13개 workflow 파일 참조 목록 |
| Execution Directive | L340-373 | **10단계 실행 절차**: Parse Args → Route → Load Workflow → Read Config → Init Tasks → Execute → Track → Present → Completion Marker → Guide Next Steps |

### 핵심 발견

1. **SKILL.md = 오케스트레이터의 두뇌**: CLAUDE.md가 "무엇을"이라면, SKILL.md는 "어떻게" 실행하는지를 정의
2. **Intent Router가 핵심**: 4단계 우선순위로 사용자 의도를 워크플로우에 매핑
3. **Workflow는 외부 파일**: SKILL.md는 요약만, 상세는 `workflows/*.md`에 위임
4. **Execution Directive**: 스킬 활성화 시 실행할 10단계 절차를 명시적으로 기술

### 2.2 Workflow 파일별 분석

**디렉토리**: `/tmp/e2e4-extract/personas/moai/skills/moai/workflows/`
**총 13개 파일**

| 파일 | 줄 수 | 핵심 역할 | Phase 구조 |
|------|-------|----------|-----------|
| **moai.md** | 219줄 | 기본 자율 파이프라인 (plan→run→sync) | Phase 0(탐색) → Routing → Phase 1(SPEC) → Phase 2(구현) → Phase 3(문서) |
| **plan.md** | 269줄 | SPEC 문서 생성 (EARS 형식) | Phase 1A(탐색) → 1B(SPEC 계획) → 1.5(검증) → Phase 2(SPEC 생성) → Phase 3(Git) |
| **run.md** | 361줄 | DDD/TDD/Hybrid 구현 | Phase 1(분석) → 1.5(태스크 분해) → Phase 2(구현, 모드별 분기) → 2.5(품질) → Phase 3(Git) → Phase 4(안내) |
| **sync.md** | 503줄 | 문서 동기화 + PR 생성 | Phase 0.5(품질검증) → Phase 1(분석) → Phase 2(문서동기화) → Phase 3(Git+PR) → Phase 4(완료) |
| **fix.md** | 172줄 | 1회성 자동 수정 | Phase 1(병렬스캔) → Phase 2(분류 L1-L4) → Phase 3(수정) → Phase 4(검증) |
| **loop.md** | 169줄 | 반복 자동 수정 | 반복: 완료체크 → 메모리체크 → 진단 → 수정 → 검증 → 반복 |
| **project.md** | 211줄 | 프로젝트 문서 생성 | Phase 0(타입감지) → Phase 1(분석) → Phase 2(확인) → Phase 3(생성) → Phase 3.5(LSP) → Phase 4(완료) |
| **feedback.md** | 142줄 | GitHub 이슈 생성 | Phase 1(수집) → Phase 2(이슈 생성) |
| **team-plan.md** | 92줄 | 팀 기반 SPEC 생성 | TeamCreate → Spawn(researcher+analyst+architect) → 병렬조사 → 종합 → 승인 → Shutdown |
| **team-run.md** | 179줄 | 팀 기반 구현 | TeamCreate → 태스크분해 → Spawn(backend+frontend+tester) → 병렬구현 → 품질검증 → Git → Shutdown |
| **team-sync.md** | 34줄 | Sync는 항상 sub-agent | 이유 설명만 (순차적 특성, 파일 적음, 토큰 예산 적음, 일관성) |
| **team-debug.md** | 79줄 | 경쟁 가설 디버깅 | TeamCreate → 가설배정(3개) → 병렬조사 → 증거종합 → 수정 → Shutdown |
| **team-review.md** | 75줄 | 다관점 코드 리뷰 | TeamCreate → 관점배정(security+perf+quality) → 병렬리뷰 → 보고서통합 → Shutdown |

### 2.3 Team Workflow Skill

**파일**: `/tmp/e2e4-extract/personas/moai/skills/moai/moai-workflow-team/SKILL.md`
**총 줄 수**: 259줄 (Version 1.1.0)

핵심 구조:
- **Mode Selection**: `--team`/`--solo`/auto (복잡도 기반: 도메인>=3, 파일>=10, 점수>=7)
- **Team Lifecycle**: 5단계 (TeamCreate → Task Decomposition → Teammate Spawning → Coordination → Shutdown)
- **File Ownership Strategy**: 프로젝트 유형별(Go/Web/Full-stack/Monorepo/Python) 소유권 테이블
- **Team Patterns**: 5종 (Plan Research, Implementation, Full-Stack, Investigation, Review)
- **Error Recovery**: 5가지 시나리오 (충돌, 정체, 파일충돌, idle, 토큰한도)

### 2.4 Reference Skill

**파일**: `/tmp/e2e4-extract/personas/moai/skills/moai/references/reference.md`
**총 줄 수**: 250줄 (Version 1.1.0)

공통 패턴과 참조 정보:
- Execution Patterns (Parallel/Sequential/Hybrid)
- Resume Pattern (--resume 플래그 동작)
- Context Propagation (페이즈 간 컨텍스트 전달)
- Flag Reference (Global/Plan/Run/Sync/Fix/Loop/MoAI 전체 플래그)
- Legacy Command Mapping (`/moai:X-Y` → `/moai subcommand`)
- Configuration Files Reference (Core/Project/SPEC/Release/Version)
- Completion Markers
- Error Handling Delegation

---

## 3. Do-Focus CLAUDE.md 대응 매핑

**파일**: `~/Work/do-focus.workspace/do-focus/CLAUDE.md`
**총 줄 수**: 395줄 (Version 3.0.0)
**언어**: 한국어 + 영어 혼합

### 3.1 Do CLAUDE.md 섹션 구조

| # | 섹션명 | 줄 범위 | 핵심 내용 |
|---|--------|---------|----------|
| 1 | Do/Focus/Team: 삼원 실행 구조 | L3-63 | 3모드 정의 (Do/Focus/Team), 선택 가이드 테이블, 자동 에스컬레이션 |
| 2 | Do 모드 상세 | L67-95 | Full Delegation [HARD], 에이전트 검증 레이어, Parallel Execution, Response Format |
| 3 | Team 모드 상세 | L98-134 | 전제조건, 실행방식, Plan/Run 팀 구성 테이블, Do vs Team 비교 |
| 4 | Violation Detection | L137-143 | 3가지 위반 사례 정의 |
| 5 | Intent-to-Agent Mapping | L146-216 | 12개 도메인 키워드→에이전트 매핑. 설계/계획 요청 3단계 순차실행. 모드 전환 `godo mode` |
| 6 | Parallel Execution Pattern | L220-234 | 병렬 실행 예시 (backend + security 동시 호출) |
| 7 | Plan Mode 지침 | L238-255 | `.do/jobs/` 경로 강제, `~/.claude/plans/` 금지 |
| 8 | 기본 규칙 | L259-332 | Git, Multirepo, 커밋메시지, 릴리즈, 코드스타일, 테스트, 안전규칙, 필수개발규칙 |
| 9 | 설정 파일 구조 | L336-371 | settings.json(공유) + settings.local.json(개인) + 환경변수 테이블 |
| 10 | 페르소나 시스템 | L373-381 | 4종 캐릭터 (young-f/m, senior-f/m), SessionStart hook 주입 |
| 11 | 스타일 전환 | L385-393 | sprint/pair/direct 3종 |

### 3.2 비교 매핑 테이블

| MoAI 섹션 | MoAI 위치 | Do 대응 | Do 위치 | 차이점/누락 |
|-----------|----------|---------|---------|------------|
| **1. Core Identity** | CLAUDE.md L3-25 | Do/Focus/Team 삼원 구조 | CLAUDE.md L3-63 | MoAI: 단일 오케스트레이터. Do: 3모드 전환 구조. Do에 HARD 규칙 목록 없음 (분산됨) |
| **2. Request Processing Pipeline** | CLAUDE.md L28-67 | **없음** | - | **Do에 누락**. Analyze→Route→Execute→Report 파이프라인 없음. Intent-to-Agent Mapping이 부분 대체하지만 체계적이지 않음 |
| **3. Command Reference** | CLAUDE.md L70-81 | 개별 commands/do/*.md | .claude/commands/do/ | MoAI: `/moai` 단일 진입점 + 서브커맨드. Do: `/do:plan`, `/do:checklist` 등 개별 커맨드. **통합 스킬 없음** |
| **4. Agent Catalog** | CLAUDE.md L84-109 | Intent-to-Agent Mapping | CLAUDE.md L146-216 | MoAI: 구조화된 카탈로그(카테고리별). Do: 키워드 기반 매핑만. Decision Tree 없음 |
| **5. SPEC-Based Workflow** | CLAUDE.md L113-133 | 필수 개발 규칙 | CLAUDE.md L326-332 | MoAI: SPEC(plan→run→sync). Do: Plan→Checklist→Develop→Test→Report. **SPEC 문서 체계 없음**, 체크리스트가 대체 |
| **6. Quality Gates** | CLAUDE.md L136-151 | **없음** (rules에 분산) | dev-testing.md, dev-workflow.md | MoAI: TRUST 5 + LSP Quality Gates 명시. Do: 테스트/품질 규칙은 rules 파일에만 |
| **7. Safe Development Protocol** | CLAUDE.md L153-198 | 코딩 규율/에러 대응 | dev-workflow.md | MoAI: 4개 HARD 규칙 + Go 전용. Do: dev-workflow.md에 더 상세 (Read Before Write, 3회 재시도 등) |
| **8. User Interaction Architecture** | CLAUDE.md L201-220 | **없음** | - | **Do에 누락**. AskUserQuestion 사용 패턴 미정의. 서브에이전트 제약 미문서화 |
| **9. Configuration Reference** | CLAUDE.md L223-246 | 설정 파일 구조 | CLAUDE.md L336-371 | MoAI: `.moai/config/` YAML 체계. Do: `settings.local.json` env 기반. **Do가 더 단순** |
| **10. Web Search Protocol** | CLAUDE.md L248-263 | 안전 규칙 일부 | CLAUDE.md L321 | MoAI: 독립 섹션. Do: 한 줄로 축약 |
| **11. Error Handling** | CLAUDE.md L265-282 | 에러 대응 | dev-workflow.md | MoAI: 에러유형별 에이전트 매핑. Do: "3회 재시도 후 사용자에게" 일반 규칙만 |
| **12. MCP Servers** | CLAUDE.md L284-293 | **없음** | - | **Do에 누락**. MCP 통합 문서 없음 |
| **13. Progressive Disclosure** | CLAUDE.md L298-311 | **없음** | - | **Do에 누락**. 토큰 최적화 시스템 미구현 |
| **14. Parallel Execution Safeguards** | CLAUDE.md L313-320 | Parallel Execution Pattern | CLAUDE.md L220-234 | MoAI: 안전장치 중심(충돌방지, 루프방지). Do: 예시 중심. 안전장치 미흡 |
| **15. Agent Teams** | CLAUDE.md L324-348 | Team 모드 상세 | CLAUDE.md L98-134 | 구조 유사하나 Do에 Team APIs, Hook Events 상세 없음 |
| - | - | **페르소나 시스템** | CLAUDE.md L373-381 | **MoAI에 없음**. Do 고유 기능 |
| - | - | **스타일 전환** | CLAUDE.md L385-393 | **MoAI에 없음**. Do 고유 기능 |
| - | - | **Multirepo 환경** | CLAUDE.md L266-279 | **MoAI에 없음**. Do 고유 기능 |
| - | - | **릴리즈 워크플로우** | CLAUDE.md L298-304 | **MoAI에 없음** (MoAI는 sync에서 처리). Do: tobrew 전용 |
| - | - | **모드 전환(godo)** | CLAUDE.md L208-216 | **MoAI에 없음**. Do 고유 기능 (statusline 동기화) |

### 3.3 Rules 파일 분포 비교

| 영역 | MoAI | Do-Focus |
|------|------|----------|
| Core 원칙 | `.claude/rules/moai/core/moai-constitution.md` | CLAUDE.md에 인라인 |
| 개발 환경 | `.moai/config/` YAML 기반 | `.claude/rules/dev-environment.md` (Docker 필수) |
| 테스트 | quality.yaml (커버리지 타겟) | `.claude/rules/dev-testing.md` (매우 상세) |
| 워크플로우 | `.claude/rules/moai/workflow/` | `.claude/rules/dev-workflow.md` + `dev-checklist.md` |
| 코딩 표준 | `.claude/rules/moai/development/` | `.claude/rules/dev-workflow.md`에 통합 |

---

## 4. Do 오케스트레이터 스킬 설계 방향

### 4.1 현재 Do의 구조적 한계

1. **통합 진입점 부재**: MoAI는 `/moai` 단일 스킬이 Intent Router로 모든 서브커맨드를 분배. Do는 `/do:plan`, `/do:checklist`, `/do:mode`, `/do:style`, `/do:setup`, `/do:check` 6개가 분리되어 있음
2. **오케스트레이션 로직 부재**: CLAUDE.md가 직접 오케스트레이터 역할. MoAI처럼 Execution Directive(10단계)가 없음
3. **워크플로우 파일 미분리**: MoAI는 `workflows/` 디렉토리에 13개 워크플로우를 분리. Do는 CLAUDE.md + dev-workflow.md + dev-checklist.md에 혼재
4. **Resume/Snapshot 없음**: MoAI는 `--resume` 플래그와 스냅샷으로 중단/재개 지원. Do는 체크리스트 상태(`[o]`/`[~]`)만으로 수동 재개
5. **Progressive Disclosure 없음**: 토큰 최적화 없이 모든 규칙이 항상 로드

### 4.2 MoAI 워크플로우 → Do 커맨드 매핑

| MoAI 서브커맨드 | MoAI 워크플로우 파일 | Do 현재 대응 | 매핑 방안 |
|----------------|---------------------|-------------|----------|
| `/moai plan` | workflows/plan.md | `/do:plan` | `/do plan` 서브커맨드로 통합. Analysis→Architecture→Plan 3단계를 workflow 파일로 분리 |
| `/moai run` | workflows/run.md | CLAUDE.md "Checklist → Develop" | `/do run` 서브커맨드. Checklist 기반 구현 실행. DDD/TDD 라우팅 추가 |
| `/moai sync` | workflows/sync.md | **없음** | `/do sync` 추가 검토. Do는 PR/문서 자동화가 약함. 당장은 불필요할 수 있음 |
| `/moai fix` | workflows/fix.md | **없음** | `/do fix` 추가 검토. 병렬 스캔 + 분류 + 자동수정. 유용하지만 복잡도 높음 |
| `/moai loop` | workflows/loop.md | **없음** | `/do loop` 추가 검토. 반복 수정. fix 없이 단독 의미 약함 |
| `/moai project` | workflows/project.md | **없음** | `/do project` 추가 검토. 프로젝트 초기화 문서 생성 |
| `/moai feedback` | workflows/feedback.md | **없음** | 우선순위 낮음. Do는 자체 프로젝트용이 아님 |
| `/moai` (default) | workflows/moai.md | **없음** (수동 파이프라인) | `/do` 기본 동작으로 plan→run 자동 파이프라인 |
| `--team` 플래그 | workflows/team-*.md (5개) | Team 모드 (CLAUDE.md) | Team 모드 활성화 시 team workflow 파일 참조 |

### 4.3 Do 오케스트레이터 스킬 구조 제안

MoAI의 구조를 참고하되, Do의 고유 특성(3모드, 페르소나, 체크리스트, Docker)을 반영:

```
.claude/skills/do/
  SKILL.md                    <- 오케스트레이터 스킬 (MoAI의 SKILL.md 대응)
  workflows/
    plan.md                   <- Analysis → Architecture → Plan (MoAI plan.md 대응)
    run.md                    <- Checklist 기반 구현 (MoAI run.md 대응)
    do.md                     <- plan → run 자동 파이프라인 (MoAI moai.md 대응)
    team-plan.md              <- Team 모드 Plan (MoAI team-plan.md 대응)
    team-run.md               <- Team 모드 Run (MoAI team-run.md 대응)
  references/
    reference.md              <- 공통 패턴, 플래그, 설정 참조 (MoAI reference.md 대응)
```

### 4.4 SKILL.md 핵심 설계 요소

MoAI SKILL.md에서 가져와야 할 것:

1. **Intent Router**: Do도 4단계 우선순위 라우팅 필요
   - Priority 1: 명시적 서브커맨드 (`/do plan`, `/do run`, `/do checklist`)
   - Priority 2: 모드 전환 감지 ("포커스 모드", "Do 모드", "팀 모드")
   - Priority 3: 자연어 분류 ("설계해줘" → plan, "구현해줘" → run)
   - Priority 4: 기본 동작 (AskUserQuestion으로 확인)

2. **Execution Directive**: 스킬 활성화 시 실행 절차
   - Step 1: Parse Arguments (서브커맨드 + 플래그 추출)
   - Step 2: Check Mode (현재 Do/Focus/Team 모드 확인)
   - Step 3: Route to Workflow (워크플로우 파일 로드)
   - Step 4: Read Config (settings.local.json 환경변수 로드)
   - Step 5: Execute Workflow
   - Step 6: Present Results (페르소나 + 스타일에 맞게 보고)

3. **Pre-execution Context**: MoAI처럼 `git status`, `git branch` 자동 수집

4. **Core Rules**: 위임 규칙, 사용자 인터랙션, 완료 조건을 SKILL.md에 집중

### 4.5 Do가 가져가지 않아도 되는 것

1. **SPEC 문서 체계** (EARS 형식, spec.md/plan.md/acceptance.md 3파일 세트): Do는 체크리스트 시스템이 더 실용적
2. **LSP Quality Gates**: Do는 dev-testing.md의 Real DB + 구문검사가 이미 충분
3. **Completion Markers** (`<moai>DONE</moai>`): Do는 체크리스트 `[o]` 상태로 충분
4. **SPEC Lifecycle Levels**: Do는 단일 plan.md → checklist.md 흐름이 더 직관적
5. **Snapshot/Resume 시스템**: 체크리스트 상태가 이미 영속 저장소 역할 수행
6. **16개 언어 LSP 감지**: Do는 프로젝트별 Docker 환경에 의존

### 4.6 기존 /do:* 커맨드의 마이그레이션

| 현재 커맨드 | 역할 | 마이그레이션 방향 |
|------------|------|-----------------|
| `/do:plan` | 플랜 작성 | `/do plan` 서브커맨드로 통합. workflow/plan.md에 상세 절차 |
| `/do:checklist` | 체크리스트 관리 | `/do checklist` 서브커맨드. 상태 조회/생성 기능 유지 |
| `/do:mode` | 모드 전환 | `/do mode` 서브커맨드. godo mode 실행 로직 유지 |
| `/do:style` | 스타일 선택 | `/do style` 서브커맨드. 간단하므로 workflow 파일 불필요 |
| `/do:setup` | 초기 설정 | `/do setup` 서브커맨드. 설정 wizard 유지 |
| `/do:check` | 설치 확인 | `/do check` 서브커맨드. 진단 로직 유지 |

**통합 후 진입점**: `/do [subcommand] [args]` 또는 `/do "자연어 요청"`

---

## 5. 핵심 차이 요약

| 관점 | MoAI | Do-Focus |
|------|------|----------|
| **철학** | 단일 오케스트레이터, SPEC 기반, 자율 실행 | 3모드 전환, 체크리스트 기반, 사용자 주도 |
| **진입점** | `/moai` 하나 (Intent Router) | `/do:*` 6개 분리 |
| **워크플로우** | 13개 workflow 파일로 분리 | CLAUDE.md + rules에 혼재 |
| **상태 관리** | SPEC 문서 + TaskCreate/TaskUpdate | 체크리스트 (.do/jobs/) |
| **구현 방법론** | DDD/TDD/Hybrid (quality.yaml 설정) | TDD 또는 일반 (사용자 선택) |
| **품질 관리** | TRUST 5 + LSP Quality Gates | dev-testing.md 규칙 (Real DB, FIRST) |
| **토큰 최적화** | Progressive Disclosure 3단계 | file-reading.md 최적화 규칙만 |
| **설정 저장** | `.moai/config/` YAML 분리 | `.claude/settings.local.json` env |
| **산출물 위치** | `.moai/specs/SPEC-XXX/` | `.do/jobs/{YYMMDD}/{title}/` |
| **고유 기능** | SPEC lifecycle, Snapshot/Resume, fix/loop | 페르소나 4종, 스타일 3종, Docker 필수, Multirepo, 릴리즈 |
| **언어** | 영어 (지시문 전체) | 한국어+영어 혼합 |

---

**작성자**: analysis agent
**다음 단계**: Part 2에서 SKILL.md 상세 설계 + CLAUDE.md 리팩토링 방안 작성
