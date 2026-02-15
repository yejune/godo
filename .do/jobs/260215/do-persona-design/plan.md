# Do Persona 구현 플랜

## 목표

do-focus의 fat CLAUDE.md(395줄)를 moai와 동일한 converter 인터페이스(manifest.yaml)를 따르는 lean orchestrator + skill workflows 구조로 재설계하여, `personas/do/` 디렉토리에 25개 파일로 구성된 완전한 do persona를 생성한다.

## 워크플로우 변경사항

architecture에서 조정된 점:
- `team-plan.md`와 `team-run.md`를 **하나의 `team-plan.md`로 통합** (Team Plan + Run + Report combined)
- 따라서 워크플로우 파일은 5개가 아닌 **4개**: `plan.md`, `run.md`, `do.md`, `team-plan.md`
- manifest.yaml의 skills 목록도 7개에서 **6개**로 축소 (team-run.md 제거)

---

## Phase 1: Foundation (manifest.yaml, CLAUDE.md, settings.json)

### Task 1.1: manifest.yaml 작성

- **파일**: `personas/do/manifest.yaml`
- **내용**: architecture Section 6 기반. converter가 assemble 시 사용하는 전체 구조 정의
- **참조**:
  - 기존 extracted manifest: `~/Work/new/convert/extracted-do/personas/do/manifest.yaml`
  - moai manifest 패턴: `/tmp/e2e4-extract/personas/moai/` (구조 참고)
- **주요 변경점 (기존 extracted manifest 대비)**:
  - `version`: `""` → `"3.0.0"`
  - `description`: `""` → `"Do execution framework with 3 modes (Do/Focus/Team)"`
  - `skills`: `[]` → 6개 (SKILL.md + 4 workflows + reference.md)
  - `styles`: `moai.md/r2d2.md/yoda.md` → `output-styles/do/sprint.md/pair.md/direct.md`
  - `agents`: `manager-spec.md` 제거 (5개로 축소)
  - `slot_content`: TRUST 5/TAG Chain 텍스트를 do 버전으로 간소화
  - `agent_patches`: `{}` 유지 (override skills 없으므로)
  - `hook_scripts`: `[]` (godo 직접 호출, shell wrapper 불필요)
- **검증**: `python -c "import yaml; yaml.safe_load(open('manifest.yaml'))"`
- **예상 규모**: ~80줄

### Task 1.2: CLAUDE.md 재작성 (lean)

- **파일**: `personas/do/CLAUDE.md`
- **내용**: 현재 395줄 → ~200줄로 축약. architecture Section 2 기반
- **핵심 원칙**: 상세 규칙은 `skills/do/` + `rules/`로 이동, CLAUDE.md는 "무엇을(what)" 요약 + 참조만

#### CLAUDE.md에서 skills/do/로 이동할 섹션 (상세)

| 현재 CLAUDE.md 위치 | 현재 내용 | 이동 대상 | 이유 |
|---|---|---|---|
| L67-95 (Do 모드 상세) | Full Delegation [HARD], 에이전트 검증 레이어, Parallel Execution, Response Format | `skills/do/SKILL.md` Core Rules 섹션 | "어떻게 실행하는가"는 SKILL.md 담당 |
| L98-134 (Team 모드 상세) | 전제조건, 실행방식, Plan/Run 팀 구성 테이블, Do vs Team 비교 | `skills/do/references/reference.md` Mode Reference 섹션 | Team 상세는 참조 문서로 |
| L146-193 (Intent-to-Agent Mapping) | 12개 도메인 키워드→에이전트 매핑 테이블 (Backend~Architecture 도메인) | `skills/do/SKILL.md` Intent Router Priority 3 (NL Classification) | 오케스트레이터 스킬의 핵심 라우팅 로직 |
| L194-207 (설계/계획 요청 3단계 순차실행) | Analysis→Architecture→Plan 순차 실행 지시, 산출물 위치, 승인 요청 패턴 | `skills/do/workflows/plan.md` | workflow 상세는 workflow 파일에 |
| L208-216 (모드 전환 godo mode) | "포커스 호출해"→`godo mode focus` 매핑, statusline 동기화 [HARD] | `skills/do/SKILL.md` Mode Router + Intent Router Priority 2 | 모드 전환은 Intent Router의 일부 |
| L220-234 (Parallel Execution Pattern) | 병렬 실행 예시 (backend + security 동시 호출), Task 동시 호출 패턴 | `skills/do/references/reference.md` Execution Patterns 섹션 | 공통 패턴 참조 |
| L238-255 (Plan Mode 워크플로우 상세) | `.do/jobs/` 경로 강제, `~/.claude/plans/` 금지, 날짜 폴더 YYMMDD | `skills/do/workflows/plan.md` Plan Mode Integration 섹션 | workflow 상세 |

#### CLAUDE.md에 유지할 섹션 (축약)

| 섹션 | 현재 위치 | 유지 이유 | 축약 방식 |
|---|---|---|---|
| Core Identity (3모드 정의) | L3-63 | CLAUDE.md = "무엇을" 정의하는 곳 | 모드 선택 가이드 테이블 제거, 3모드 1-2줄 정의 + 에스컬레이션 조건만 유지 (~30줄) |
| Violation Detection | L137-143 | 빠른 참조용 HARD 규칙 | 그대로 유지 (~7줄) |
| 기본 규칙 (Git, Multirepo, 커밋, 릴리즈) | L259-332 | Do 고유 기능 + dev-workflow.md 참조 | 규칙 이름 + [HARD] 태그 + `@rules` 참조로 축약 (~15줄). 상세는 rules에 위임 |
| 설정 파일 구조 | L336-371 | 환경변수 테이블은 자주 참조됨 | 환경변수 테이블만 유지, 설명 축약 (~20줄) |
| 페르소나 시스템 | L373-381 | Do 고유 기능 | 그대로 유지 (~10줄) |
| 스타일 전환 | L385-393 | Do 고유 기능 | 그대로 유지 (~5줄) |

#### CLAUDE.md에 신규 추가할 섹션

| 섹션 | moai 대응 | 내용 | 줄 수 |
|---|---|---|---|
| Request Processing Pipeline | moai CLAUDE.md Section 2 | Analyze→Route→Execute→Report 4단계 (Do 버전) | ~15줄 |
| Command Reference | moai CLAUDE.md Section 3 | `/do` 통합 진입점 + 서브커맨드 목록 + Allowed Tools | ~10줄 |
| Agent Catalog (요약) | moai CLAUDE.md Section 4 | Selection Decision Tree + 4카테고리 1줄 요약 | ~25줄 |
| Checklist-Based Workflow (요약) | moai CLAUDE.md Section 5 | Plan→Checklist→Develop→Test→Report 파이프라인 요약 | ~15줄 |
| Quality Gates | moai CLAUDE.md Section 6 | dev-testing.md + dev-workflow.md 통합 참조 | ~10줄 |
| User Interaction Architecture | moai CLAUDE.md Section 8 | AskUserQuestion 패턴 (오케스트레이터 only, 서브에이전트 불가) | ~15줄 |
| Error Handling | moai CLAUDE.md Section 11 | 에러 유형별 복구 + 3회 재시도 | ~10줄 |
| Parallel Execution Safeguards | moai CLAUDE.md Section 14 | 파일 충돌 방지, 루프 방지 | ~10줄 |
| Agent Teams (Experimental) | moai CLAUDE.md Section 15 | 활성화 조건 + Team APIs 참조 | ~15줄 |

#### 최종 CLAUDE.md 섹션 구성 (목표 ~200줄)

```
1. Core Identity (~30줄)
2. Request Processing Pipeline (~15줄) [신규]
3. Command Reference (~10줄) [신규]
4. Agent Catalog (~25줄) [신규 구조화]
5. Checklist-Based Workflow (~15줄) [신규]
6. Quality Gates (~10줄) [신규 통합]
7. Safe Development Protocol (~15줄) [축약]
8. User Interaction Architecture (~15줄) [신규]
9. Configuration Reference (~20줄) [축약]
10. Persona System (~10줄) [유지]
11. Style Switching (~5줄) [유지]
12. Error Handling (~10줄) [독립 섹션화]
13. Parallel Execution Safeguards (~10줄) [안전장치 추가]
14. Agent Teams (Experimental) (~15줄) [축약]
합계: ~205줄
```

- **검증**: 현재 CLAUDE.md의 모든 [HARD] 규칙이 새 CLAUDE.md 또는 skills/rules에 존재하는지 체크리스트 대조
- **예상 규모**: ~200줄

### Task 1.3: settings.json 작성

- **파일**: `personas/do/settings.json`
- **내용**: godo hooks 직접 호출 설정. 현재 do-focus settings.json 거의 그대로 사용
- **참조**: `~/Work/do-focus.workspace/do-focus/.claude/settings.json`
- **변경점**: 없음 (현재 settings.json이 이미 godo 직접 호출 패턴). 그대로 복사.
- **검증**: `python -c "import json; json.load(open('settings.json'))"`
- **예상 규모**: ~135줄 (현재와 동일)

---

## Phase 2: Orchestrator Skill (핵심)

### Task 2.1: SKILL.md 작성

- **파일**: `personas/do/skills/do/SKILL.md`
- **내용**: Mode Router + Intent Router + Core Rules + Execution Directive
- **참조**:
  - moai SKILL.md 구조 (377줄): `/tmp/e2e4-extract/personas/moai/skills/moai/SKILL.md`
  - do CLAUDE.md의 기존 로직: `~/Work/do-focus.workspace/do-focus/CLAUDE.md`
  - architecture Section 3.1 설계

#### Frontmatter

```yaml
name: do
description: >
  Do super agent - unified orchestrator with 3 execution modes (Do/Focus/Team).
  Routes natural language or explicit subcommands (plan, run, checklist, mode,
  style, setup, check) to specialized agents or direct execution.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Task AskUserQuestion TaskCreate TaskUpdate TaskList TaskGet Bash Read Write Edit Glob Grep
user-invocable: true
metadata:
  argument-hint: "[subcommand] [args] | \"natural language task\""
```

#### CLAUDE.md에서 이동해 올 내용 (구체적)

| SKILL.md 섹션 | 원본 위치 (CLAUDE.md) | 변환 방식 |
|---|---|---|
| **Pre-execution Context** | 없음 (신규) | moai 패턴 차용: `git status`, `git branch`, `godo statusline` |
| **Core Identity** (~20줄) | L67-79 (Do 모드 HARD 규칙 3개) | 6개 원칙으로 재구성: 위임, 직접구현금지, AskUserQuestion 오케스트레이터 only, 병렬실행, 언어감지, 모드 인식 |
| **Mode Router** (~25줄) | L208-216 (모드 전환) + L41-63 (에스컬레이션) | 신규 섹션: DO_MODE 환경변수 확인 → 모드별 행동 분기 → 에스컬레이션 규칙 |
| **Intent Router Priority 1** (~15줄) | 없음 (개별 commands에 분산) | 신규: plan/run/checklist/mode/style/setup/check 서브커맨드 매핑 |
| **Intent Router Priority 2** (~10줄) | L208-216 (모드 전환 키워드) | "포커스", "Focus", "Do 모드", "팀 모드" → `godo mode <mode>` 실행 |
| **Intent Router Priority 3** (~30줄) | L146-193 (Intent-to-Agent Mapping 12개 도메인) | 도메인 키워드→에이전트 매핑을 NL Classification으로 재구성 |
| **Intent Router Priority 4** (~5줄) | 없음 (신규) | Ambiguous → AskUserQuestion with top 2-3 options |
| **Workflow Quick Reference** (~50줄) | L194-207 (설계/계획 3단계) + L238-255 (Plan Mode) | 4개 워크플로우 + 6개 유틸리티 요약 (상세는 workflows/*.md) |
| **Core Rules - Agent Delegation** (~15줄) | L67-79 (Full Delegation [HARD]) | [HARD] ALL implementation via Task(). NEVER implement directly. |
| **Core Rules - User Interaction** (~15줄) | 없음 (CLAUDE.md에 누락이었음) | moai SKILL.md L189-205에서 차용: AskUserQuestion 패턴 5단계 |
| **Core Rules - Checklist Tracking** (~10줄) | L326-332 (필수 개발 규칙 참조) | moai의 Task Tracking을 체크리스트 시스템으로 대체 |
| **Core Rules - Output Rules** (~10줄) | L385-393 (스타일) + 분산된 응답 규칙 | 언어 규칙 + 마크다운 + XML 금지 + 페르소나 스타일 |
| **Core Rules - Error Handling** (~10줄) | L321-322 (안전 규칙 에러 핸들링) | 3회 재시도, expert-debug 위임, /clear 전략 |
| **Agent Catalog** (~50줄) | L146-193 (도메인→에이전트 매핑) | Manager 6 + Expert 8 + Builder 3 + Team 8 + Decision Tree |
| **Common Patterns** (~20줄) | L220-234 (Parallel Execution Pattern) | Parallel/Sequential/Resume(체크리스트 기반)/Context Propagation |
| **Execution Directive** (~40줄) | 없음 (신규, moai 10단계→8단계) | 8단계: Parse→CheckMode→Route→ReadConfig→Execute→UpdateChecklist→Present→GuideNext |

- **검증**: Intent Router의 12개 도메인 키워드가 현재 CLAUDE.md L146-193과 1:1 매핑되는지 확인
- **예상 규모**: ~350줄

### Task 2.2: plan.md workflow

- **파일**: `personas/do/skills/do/workflows/plan.md`
- **내용**: Analysis→Architecture→Plan 3단계 상세
- **참조**:
  - do CLAUDE.md L194-207 (설계/계획 요청 3단계)
  - do CLAUDE.md L238-255 (Plan Mode 워크플로우)
  - `dev-workflow.md` 복잡도 판단 기준
  - `dev-checklist.md` Analysis/Architecture 템플릿
  - moai workflows/plan.md (269줄, 구조 참고)
- **CLAUDE.md에서 이동해 올 내용**:
  - L194-207: Analysis→Architecture→Plan 3단계 순차 실행 규칙, 산출물 위치 `.do/jobs/`, 승인 요청 패턴
  - L238-255: Plan Mode(Shift+Tab) 저장 위치 강제 `.do/jobs/{YYMMDD}/{title}/plan.md`, `~/.claude/plans/` 금지, 날짜 폴더 YYMMDD 형식
- **섹션 구조**:
  - 목적, 복잡도 판단 (dev-workflow.md 참조)
  - Phase 1: Complexity Assessment
  - Phase 2A: Analysis (복잡한 작업만, Task(expert-analyst))
  - Phase 2B: Architecture (복잡한 작업만, Task(expert-architect))
  - Phase 3: Plan Generation
  - Phase 4: User Approval (AskUserQuestion)
  - Plan Mode (Shift+Tab) Integration [HARD]
- **검증**: 복잡도 판단 기준이 dev-workflow.md와 일치하는지 확인
- **예상 규모**: ~200줄

### Task 2.3: run.md workflow

- **파일**: `personas/do/skills/do/workflows/run.md`
- **내용**: Checklist 기반 DDD/TDD 구현
- **참조**:
  - `dev-workflow.md` 에이전트 위임 필수 전달사항
  - `dev-checklist.md` 상태 관리/의존성/서브 체크리스트 템플릿
  - moai workflows/run.md (361줄, 구조 참고)
- **CLAUDE.md에서 이동해 올 내용**: 직접 이동 없음. dev-workflow.md와 dev-checklist.md의 규칙을 워크플로우 형태로 재구성.
- **섹션 구조**:
  - 목적, 전제조건 (plan.md + checklist.md 존재)
  - Phase 1: Checklist Verification (미완료 항목 확인, 의존성 확인)
  - Phase 2: Agent Dispatch (Do Mode → 병렬 Task, Focus Mode → 직접 실행, Team Mode → TeamCreate)
  - Phase 3: Progress Monitoring (에이전트 완료 후 체크리스트 확인, 미커밋 재호출, 토큰 소진 재개)
  - Phase 4: Quality Verification (전체 테스트 스위트)
  - Phase 5: Completion (체크리스트 표시, report.md 작성)
- **검증**: 에이전트 위임 시 4가지 필수 전달사항이 포함되어 있는지 확인
- **예상 규모**: ~250줄

### Task 2.4: do.md workflow

- **파일**: `personas/do/skills/do/workflows/do.md`
- **내용**: plan→run 자동 파이프라인
- **참조**:
  - moai workflows/moai.md (219줄, 구조 참고)
  - 위의 plan.md, run.md workflow
- **CLAUDE.md에서 이동해 올 내용**: 직접 이동 없음. plan + checklist + run을 파이프라인으로 연결하는 신규 워크플로우.
- **섹션 구조**:
  - 목적 (plan → checklist → run 자동 파이프라인)
  - Phase 0: Exploration (병렬: 코드베이스 분석 + 요구사항 파악)
  - Phase 1: Plan (workflows/plan.md 실행 → 사용자 승인 대기)
  - Phase 2: Checklist (plan.md 기반 checklist.md + 서브 체크리스트 자동 생성)
  - Phase 3: Run (workflows/run.md 실행 → 모드별 분기)
  - Phase 4: Report (완료 보고서 + 다음 단계 안내)
- **검증**: Phase 간 전환 시 사용자 승인 포인트가 명확한지 확인
- **예상 규모**: ~100줄

---

## Phase 3: Team Workflow + Reference

### Task 3.1: team-plan.md workflow (Plan + Run + Report 통합)

- **파일**: `personas/do/skills/do/workflows/team-plan.md`
- **내용**: Team 모드에서의 전체 워크플로우 (Plan + Run + Report를 하나의 파일로 통합)
- **참조**:
  - moai workflows/team-plan.md (92줄) + team-run.md (179줄) 구조 참고
  - architecture Section 3.2의 team-plan.md + team-run.md 설계 통합
- **섹션 구조**:
  - 목적: Team 모드 전체 워크플로우
  - Part 1: Team Plan (병렬 조사)
    - Team Composition: team-researcher(haiku, plan) + team-analyst(inherit, plan) + team-architect(inherit, plan)
    - Execution: TeamCreate → Spawn 3 teammates → 병렬 조사 → 결과 종합 → plan.md 생성 → Shutdown team
  - Part 2: Team Run (병렬 구현)
    - Team Composition: team-backend-dev(acceptEdits) + team-frontend-dev(acceptEdits) + team-tester(acceptEdits) + team-quality(plan, read-only)
    - File Ownership: 체크리스트 서브 파일 기반 소유권 배정
    - Git Staging Rules: 개별 파일 git add, broad staging 금지 (MEMORY.md 참조)
    - Execution: TeamCreate → 태스크 분해 → Spawn → 각 팀원 서브 체크리스트 실행 → 품질 검증 → Shutdown
  - Part 3: Report
    - 완료 보고서 생성 (dev-checklist.md 보고서 템플릿)
    - 체크리스트 최종 상태 표시
- **검증**: Git staging 안전 규칙이 MEMORY.md의 사고 사례를 반영하는지 확인
- **예상 규모**: ~200줄

### Task 3.2: reference.md

- **파일**: `personas/do/skills/do/references/reference.md`
- **내용**: 공통 패턴, 환경변수, 설정 참조
- **참조**:
  - moai references/reference.md (250줄, 구조 참고)
  - architecture Section 3.3 설계
- **CLAUDE.md에서 이동해 올 내용**:
  - L98-134 (Team 모드 상세): 팀 구성 테이블, Do vs Team 비교표
  - L220-234 (Parallel Execution Pattern): 병렬 실행 예시
- **섹션 구조**:
  - Execution Patterns (Parallel/Sequential/Resume)
  - Mode Reference (Do/Focus/Team 상세 비교)
  - Auto-Escalation Thresholds
  - Configuration Reference (settings.json, settings.local.json)
  - Environment Variables (DO_MODE, DO_USER_NAME 등 테이블)
  - Artifact Locations (.do/jobs/ 경로 패턴)
  - Persona System (4종 캐릭터 상세)
  - Do-Specific Features (Multirepo, Release, Mode Switching)
- **검증**: 환경변수 테이블이 CLAUDE.md L336-371과 일치하는지 확인
- **예상 규모**: ~150줄

---

## Phase 4: Styles + Commands + Rules (복사 + 조정)

> Phase 1 완료 후 독립 실행 가능. Phase 2/3과 병렬 가능.

### Task 4.1: Output Styles 3개

- **파일**:
  - `personas/do/output-styles/do/pair.md` (기본 스타일)
  - `personas/do/output-styles/do/sprint.md`
  - `personas/do/output-styles/do/direct.md`
- **출처**: `~/Work/do-focus.workspace/do-focus/.claude/styles/` 기반
- **변환 규칙**:
  - `moai.md` → `pair.md`: "Strategic Orchestrator" 톤을 "친절한 동료" 톤으로 재작성. `moai-constitution` 참조를 `do-constitution`으로 변경. `.moai/` 경로를 `.do/` 경로로 변경
  - `r2d2.md` → `sprint.md`: "Pair Programming" 톤을 "민첩한 실행자" 톤으로 재작성. 최소 대화, 빠른 실행 강조
  - `yoda.md` → `direct.md`: "Technical Wisdom" 톤을 "직설적 전문가" 톤으로 재작성. 군더더기 없는 답변
- **주의**: 스타일 파일 내부의 `moai` 참조, `.moai/` 경로, SPEC 관련 용어를 모두 do 체계로 치환
- **검증**: 각 파일에 `moai` 문자열이 남아있지 않은지 grep 확인
- **예상 규모**: pair ~270줄, sprint ~580줄, direct ~360줄

### Task 4.2: Commands 6개

- **파일**:
  - `personas/do/commands/do/check.md`
  - `personas/do/commands/do/checklist.md`
  - `personas/do/commands/do/mode.md`
  - `personas/do/commands/do/plan.md`
  - `personas/do/commands/do/setup.md`
  - `personas/do/commands/do/style.md`
- **출처**: `~/Work/do-focus.workspace/do-focus/.claude/commands/do/` 그대로 복사
- **변경점**: 없음. 내용 동일.
- **검증**: 파일 존재 + 내용 diff 없음 확인
- **예상 규모**: 각 파일 원본 크기 그대로

### Task 4.3: Rules 2개

- **파일**:
  - `personas/do/rules/do/workflow/spec-workflow.md`
  - `personas/do/rules/do/workflow/workflow-modes.md`
- **출처**: `~/Work/do-focus.workspace/do-focus/.claude/rules/do/workflow/` 기반
- **변경점**:
  - `spec-workflow.md`: 현재 do-focus에서 이미 Do 체계로 작성되어 있으므로 그대로 복사
  - `workflow-modes.md`: DDD/TDD/Hybrid 방법론. 변경 없음. 그대로 복사
- **검증**: 파일 존재 확인
- **예상 규모**: 원본 크기 그대로

---

## Phase 5: Agents + E2E 검증

### Task 5.1: Agents 복사 (5개)

- **파일**:
  - `personas/do/agents/do/manager-ddd.md`
  - `personas/do/agents/do/manager-project.md`
  - `personas/do/agents/do/manager-quality.md`
  - `personas/do/agents/do/manager-tdd.md`
  - `personas/do/agents/do/team-quality.md`
- **출처**: `~/Work/do-focus.workspace/do-focus/.claude/agents/do/` 중 위 5개만
- **제외**: `manager-spec.md` (SPEC 체계 미사용. plan workflow 내 에이전트 위임으로 대체)
- **변경점**: 없음. 그대로 복사.
- **검증**: 5개 파일 존재 확인. `manager-spec.md`가 포함되지 않았는지 확인.
- **참고**: 나머지 에이전트 (expert-*, builder-*, team-*, manager-docs, manager-git, manager-strategy)는 **core_files**로 자동 포함됨 (persona 전용이 아닌 공유 에이전트)

### Task 5.2: E2E 검증

- **검증 항목**:
  1. **Assemble 실행**: converter의 assemble 명령으로 `core + do persona → .claude/` 디렉토리 생성
  2. **구조 비교**: 생성된 `.claude/`와 현재 `~/Work/do-focus.workspace/do-focus/.claude/`를 diff
     - 예상 차이: styles 이름 변경, CLAUDE.md 축약, skills/do/ 신규 추가
     - 예상 동일: agents/do/ (core + persona), rules/, commands/
  3. **슬롯 잔존 확인**: `grep -r "{{slot:" .claude/` → 0개여야 함. 미치환 슬롯 발견 시 manifest의 slot_content에 do 버전 추가
  4. **HARD 규칙 무결성**: 현재 CLAUDE.md의 모든 [HARD] 태그가 새 구조(CLAUDE.md + SKILL.md + rules)에 존재하는지 대조
  5. **Intent Router 검증**: 현재 CLAUDE.md L146-216의 12개 도메인 키워드가 SKILL.md Intent Router에 빠짐없이 매핑되어 있는지 확인
  6. **moai 잔존 확인**: `grep -ri "moai" personas/do/` → moai-constitution.md 참조(do-constitution으로 rename 예정)를 제외하고 0개여야 함

---

## 의존성 그래프

```
Phase 1 (Foundation)
  ├── Task 1.1 manifest.yaml ──┐
  ├── Task 1.2 CLAUDE.md ──────┤ (순차: manifest가 전체 구조 정의)
  └── Task 1.3 settings.json ──┘ (독립: 복사 수준)
         │
         ▼
Phase 2 (Orchestrator Skill) ─── Phase 1 완료 후 시작
  ├── Task 2.1 SKILL.md ────────┐ (CLAUDE.md 축약 내용이 SKILL.md로 이동하므로 의존)
  ├── Task 2.2 plan.md ─────────┤ (SKILL.md와 병렬 가능하나 Workflow Quick Ref가 참조)
  ├── Task 2.3 run.md ──────────┤ (plan.md와 병렬 가능)
  └── Task 2.4 do.md ───────────┘ (plan.md + run.md 참조하므로 마지막)
         │
         ▼
Phase 3 (Team + Reference) ─── Phase 2 완료 후 시작
  ├── Task 3.1 team-plan.md ────┐ (run.md의 Agent Dispatch 패턴 참조)
  └── Task 3.2 reference.md ───┘ (SKILL.md + CLAUDE.md 이동 내용 기반)
         │
         ▼
Phase 4 (Styles + Commands + Rules) ─── Phase 1 완료 후 독립 실행 가능
  ├── Task 4.1 Styles 3개 ──────┐ (병렬)
  ├── Task 4.2 Commands 6개 ────┤ (병렬, 단순 복사)
  └── Task 4.3 Rules 2개 ───────┘ (병렬, 단순 복사)
         │
         ▼
Phase 5 (Agents + E2E) ─── 모든 Phase 완료 후
  ├── Task 5.1 Agents 5개 (단순 복사)
  └── Task 5.2 E2E 검증
```

**병렬 실행 가능한 그룹**:
- Phase 2 내: Task 2.1 + 2.2 + 2.3 병렬 → Task 2.4 순차
- Phase 3 내: Task 3.1 + 3.2 병렬
- Phase 4 전체: Phase 1 완료 후 Phase 2/3과 병렬 가능
- Phase 4 내: Task 4.1 + 4.2 + 4.3 모두 병렬

---

## 에이전트 배정

### 순차 실행 (의존성 있음)
- **Phase 1**: 1명 에이전트가 manifest → CLAUDE.md → settings 순서로 작성
- **Phase 2**: SKILL.md 작성 에이전트 1명 + workflows 작성 에이전트 1명 (plan/run/do 순차)
- **Phase 3**: 1명 에이전트가 team-plan.md + reference.md 작성

### 병렬 실행 (독립 작업)
- **Phase 4**: 3명 병렬 (styles 담당 / commands 담당 / rules+agents 담당)
  - 또는 단순 복사이므로 1명이 순차 처리해도 무방

### 검증 전담
- **Phase 5**: 1명 에이전트가 assemble + diff + grep 검증

---

## 예상 파일 수

| 유형 | 파일 목록 | 파일 수 | 작업 난이도 |
|------|----------|---------|-----------|
| **신규 작성 (핵심)** | manifest.yaml, CLAUDE.md, SKILL.md, plan.md, run.md, do.md, team-plan.md, reference.md | 8개 | 높음 |
| **기존 기반 재작성** | pair.md, sprint.md, direct.md | 3개 | 중간 |
| **기존 복사** | 6 commands + 2 rules + 5 agents + 1 settings.json | 14개 | 낮음 |
| **총계** | | **25개** | |

---

## 변경 이력

- 2026-02-15: 초안 작성
  - architecture.md 기반 plan 생성
  - 사용자 조정 반영: team-plan.md + team-run.md → team-plan.md 통합 (4 workflows)

---

## 변경 이력 (추가)

- 2026-02-15 (사용자 피드백 반영):
  - **워크플로우 4개 → 6개**: test.md(TDD), report.md(완료보고서) 추가
  - **team-plan.md → team-do.md**: 팀 오토파일럿으로 변경
  - **characters/ 폴더 추가**: young-f.md, young-m.md, senior-f.md, senior-m.md (4개)
  - **모든 workflow는 moai workflow 양식을 참고**: Purpose, Trigger, Prerequisites, Execution Steps, Output, Error Handling 섹션 구조
  - 총 파일 수: 25개 → **29개** (workflow +2, characters +4, team-plan→team-do 변경)

### 최종 워크플로우 파일 구성 (6개)
| 파일 | moai 대응 | 내용 |
|------|----------|------|
| plan.md | workflows/plan.md (269줄) | Analysis→Architecture→Plan |
| run.md | workflows/run.md (361줄) | Checklist 기반 구현 |
| test.md | 없음 (do 고유) | TDD RED-GREEN-REFACTOR |
| report.md | workflows/sync.md 대응 | 완료 보고서 |
| do.md | workflows/moai.md (219줄) | plan→run→test→report 자동 파이프라인 |
| team-do.md | team-plan+team-run 통합 | 팀 오토파일럿 (plan+run+test+report) |

### 최종 디렉토리 구조
```
personas/do/
  CLAUDE.md
  manifest.yaml
  settings.json
  characters/
    young-f.md / young-m.md / senior-f.md / senior-m.md
  skills/do/
    SKILL.md
    workflows/
      plan.md / run.md / test.md / report.md / do.md / team-do.md
    references/reference.md
  output-styles/do/
    sprint.md / pair.md / direct.md
  commands/do/
    check.md / checklist.md / mode.md / plan.md / setup.md / style.md
  rules/do/workflow/
    spec-workflow.md / workflow-modes.md
  agents/do/
    manager-ddd.md / manager-project.md / manager-quality.md / manager-tdd.md / team-quality.md
```
