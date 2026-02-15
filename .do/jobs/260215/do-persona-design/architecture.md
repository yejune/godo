# Do Persona Architecture Design

## Overview

do persona는 moai의 "SPEC 기반 자율 실행" 아키텍처를 "체크리스트 기반 3모드 실행"으로 재설계한다. 핵심은 fat CLAUDE.md를 lean orchestrator + skill workflows로 분리하고, godo binary 직접 호출 패턴을 유지하면서 moai와 동일한 converter 인터페이스(manifest.yaml)를 제공하는 것이다.

```
personas/do/
  CLAUDE.md                          <- lean orchestrator (~200줄, 요약+참조)
  manifest.yaml                      <- converter 인터페이스
  settings.json                      <- godo hooks 직접 호출
  skills/do/                         <- 오케스트레이터 스킬
  │ SKILL.md                         <- Intent Router + Execution Directive
  │ workflows/                       <- 워크플로우 상세
  │   plan.md                        <- Analysis->Architecture->Plan
  │   run.md                         <- Checklist 기반 구현
  │   do.md                          <- plan->run 자동 파이프라인
  │   team-plan.md                   <- Team Plan
  │   team-run.md                    <- Team Run
  │ references/
  │   reference.md                   <- 공통 패턴/플래그
  output-styles/do/                  <- 응답 스타일 3종
  │ sprint.md / pair.md / direct.md
  commands/do/                       <- 6개 커맨드 (기존 유지)
  rules/do/workflow/                 <- persona 전용 규칙 2개
  agents/do/                         <- persona 전용 에이전트 5개
```

---

## 1. Directory Structure

```
personas/do/
  ├── CLAUDE.md                          # lean orchestrator (~200줄)
  ├── manifest.yaml                      # converter manifest
  ├── settings.json                      # hooks + permissions + env
  │
  ├── skills/
  │   └── do/                            # 오케스트레이터 스킬 (핵심)
  │       ├── SKILL.md                   # Intent Router + Core Rules + Execution Directive
  │       ├── workflows/
  │       │   ├── plan.md                # Analysis → Architecture → Plan 3단계
  │       │   ├── run.md                 # Checklist 기반 DDD/TDD 구현
  │       │   ├── do.md                  # plan → run 자동 파이프라인
  │       │   ├── team-plan.md           # Team 모드 Plan (병렬 조사)
  │       │   └── team-run.md            # Team 모드 Run (병렬 구현)
  │       └── references/
  │           └── reference.md           # 공통 패턴, 설정 참조, 환경변수
  │
  ├── output-styles/
  │   └── do/                            # 3종 스타일
  │       ├── sprint.md                  # 민첩한 실행자
  │       ├── pair.md                    # 친절한 동료 (기본)
  │       └── direct.md                  # 직설적 전문가
  │
  ├── commands/
  │   └── do/                            # 6개 슬래시 커맨드
  │       ├── check.md                   # 체크리스트 상태 확인
  │       ├── checklist.md               # 체크리스트 생성/관리
  │       ├── mode.md                    # Do/Focus/Team 모드 전환
  │       ├── plan.md                    # 플랜 생성
  │       ├── setup.md                   # 초기 설정
  │       └── style.md                   # 스타일 전환
  │
  ├── rules/
  │   └── do/
  │       └── workflow/
  │           ├── spec-workflow.md        # Do 3단계 워크플로우 (Plan→Run→Report)
  │           └── workflow-modes.md       # DDD/TDD/Hybrid 방법론
  │
  └── agents/
      └── do/                            # persona 전용 에이전트 5개
          ├── manager-ddd.md             # DDD 구현 매니저
          ├── manager-project.md         # 프로젝트 설정 매니저
          ├── manager-quality.md         # 품질 검증 매니저
          ├── manager-tdd.md             # TDD 구현 매니저
          └── team-quality.md            # Team 모드 품질 검증
```

### moai와의 구조 대응

| moai persona | do persona | 변경 사유 |
|---|---|---|
| `skills/moai/SKILL.md` (377줄) | `skills/do/SKILL.md` (~350줄) | Intent Router를 Do 3모드에 맞게 재설계 |
| `skills/moai/workflows/` (13개) | `skills/do/workflows/` (5개) | sync/fix/loop/feedback/project 제거 |
| `skills/moai/references/` | `skills/do/references/` | Do 설정체계에 맞게 재작성 |
| `output-styles/moai/` (moai,r2d2,yoda) | `output-styles/do/` (sprint,pair,direct) | Do 고유 스타일 |
| `hooks/moai/` (7개 .sh) | (없음) | godo 직접 호출, shell wrapper 불필요 |
| override skills 6개 | (없음) | rules/*.md에 이미 흡수됨 |
| `commands/moai/` (2개) | `commands/do/` (6개) | Do 고유 커맨드 |
| `agents/moai/` (6개) | `agents/do/` (5개) | manager-spec 제거 |

---

## 2. CLAUDE.md Redesign

### 설계 원칙

현재 do-focus CLAUDE.md는 395줄의 "fat" 문서로, 오케스트레이션 로직/워크플로우/규칙이 모두 인라인되어 있다. moai의 패턴(358줄 lean orchestrator + skills로 상세 위임)을 따라 재구조화한다.

### 섹션 구조 (목표: ~200줄)

| # | 섹션명 | 줄 수 | 내용 | 현재 위치 → 이동 |
|---|--------|-------|------|-----------------|
| 1 | Core Identity | ~30줄 | 3모드(Do/Focus/Team) 정의 + HARD 규칙 목록 | CLAUDE.md L3-63 → **유지하되 축약** |
| 2 | Request Processing Pipeline | ~15줄 | Analyze→Route→Execute→Report 4단계 | **신규 추가** (moai 대응) |
| 3 | Command Reference | ~10줄 | `/do` 통합 진입점 + 서브커맨드 목록 | CLAUDE.md commands 참조 → **신규 구조화** |
| 4 | Agent Catalog | ~25줄 | Selection Decision Tree + 4카테고리 요약 | CLAUDE.md L146-216 → **축약 (상세는 rules 참조)** |
| 5 | Checklist-Based Workflow | ~15줄 | Plan→Checklist→Develop→Test→Report 요약 | CLAUDE.md L326-332 → **요약만 유지** |
| 6 | Quality Gates | ~10줄 | dev-testing.md + dev-workflow.md 참조 | (분산) → **통합 참조 섹션** |
| 7 | Safe Development Protocol | ~15줄 | HARD 규칙 4개 요약 + dev-workflow.md 참조 | CLAUDE.md 기본규칙 → **축약** |
| 8 | User Interaction Architecture | ~15줄 | AskUserQuestion 패턴 (오케스트레이터 only) | **신규 추가** (moai 대응, 현재 누락) |
| 9 | Configuration Reference | ~20줄 | settings.local.json + 환경변수 테이블 | CLAUDE.md L336-371 → **축약** |
| 10 | Persona System | ~10줄 | 4종 캐릭터 + SessionStart hook | CLAUDE.md L373-381 → **유지** |
| 11 | Style Switching | ~5줄 | sprint/pair/direct + `/do:style` 참조 | CLAUDE.md L385-393 → **유지** |
| 12 | Error Handling | ~10줄 | 에러 유형별 복구 + 3회 재시도 | 안전규칙 → **독립 섹션화** |
| 13 | Parallel Execution Safeguards | ~10줄 | 파일 충돌 방지, 루프 방지 | CLAUDE.md L220-234 → **안전장치 추가** |
| 14 | Agent Teams (Experimental) | ~15줄 | 활성화 조건 + Team APIs 참조 | CLAUDE.md L98-134 → **축약** |

### 현재 CLAUDE.md에서 skills/do/로 이동할 내용

| 현재 위치 | 이동 대상 | 이유 |
|----------|----------|------|
| Intent-to-Agent Mapping (L146-216, 71줄) | `skills/do/SKILL.md` Intent Router | 오케스트레이터 스킬의 핵심 로직 |
| 설계/계획 요청 3단계 순차 실행 (L194-207) | `skills/do/workflows/plan.md` | workflow 상세는 skill에 |
| 모드 전환 godo mode (L208-216) | `skills/do/SKILL.md` Mode Switch | Intent Router의 일부 |
| Parallel Execution Pattern 예시 (L220-234) | `skills/do/references/reference.md` | 공통 패턴 참조 |
| Plan Mode 워크플로우 상세 (L238-255) | `skills/do/workflows/plan.md` | workflow 상세 |
| Do 모드 상세 + Violation (L67-143) | `skills/do/SKILL.md` Core Rules | SKILL.md가 "어떻게"를 담당 |
| Team 모드 상세 (L98-134) | `skills/do/references/reference.md` | Team 관련 참조 |

### 현재 CLAUDE.md에서 유지할 내용

| 섹션 | 이유 |
|------|------|
| Core Identity (3모드 정의) | CLAUDE.md = "무엇을" 정의 |
| HARD 규칙 목록 (요약) | 빠른 참조용 |
| Agent Catalog (요약) | Decision Tree만 |
| Configuration Reference | 환경변수 테이블 |
| Persona System | Do 고유 기능 |
| Style Switching | Do 고유 기능 |

### CLAUDE.md에서 참조할 외부 파일

```
@.claude/rules/dev-workflow.md          # 개발 워크플로우 상세
@.claude/rules/dev-testing.md           # 테스트 규칙 상세
@.claude/rules/dev-checklist.md         # 체크리스트 시스템
@.claude/rules/dev-environment.md       # Docker/환경 규칙
@.claude/rules/do/core/do-constitution.md  # 핵심 원칙 (moai-constitution에서 rename)
```

---

## 3. Orchestrator Skill Design (skills/do/)

### 3.1 SKILL.md 구조 (~350줄)

moai SKILL.md (377줄)을 참고하되 Do의 3모드 + 체크리스트 체계를 반영한다.

#### Frontmatter

```yaml
---
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
---
```

#### 본문 구조

| 섹션 | 줄 수 | moai 대응 | Do 특화 내용 |
|------|-------|----------|-------------|
| Pre-execution Context | ~5줄 | 동일 | `git status`, `git branch`, `godo statusline` (현재 모드 확인) |
| Essential Files | ~3줄 | `.moai/config/config.yaml` | 없음 (Do는 settings.local.json의 env로 관리) |
| Core Identity | ~20줄 | 동일 구조 | 3모드 원칙 6개: 위임, 직접구현금지, AskUserQuestion 오케스트레이터 only, 병렬실행, 언어감지, 모드 인식 |
| **Mode Router** | ~25줄 | 없음 (moai는 단일 모드) | **Do 고유**: 현재 모드(Do/Focus/Team) 확인 → 모드별 실행 전략 분기 |
| **Intent Router** | ~50줄 | Priority 1-4 | Priority 1: 명시적 서브커맨드, Priority 2: 모드 전환 감지, Priority 3: 자연어 분류, Priority 4: AskUserQuestion |
| Workflow Quick Reference | ~50줄 | 8개 워크플로우 | 5개 워크플로우 (plan, run, do, team-plan, team-run) + 6개 유틸리티 (mode, style, setup, check, checklist) |
| Core Rules | ~60줄 | 동일 구조 | Agent Delegation, User Interaction, Checklist Tracking (Task 대신 체크리스트), Output Rules, Error Handling |
| Agent Catalog | ~50줄 | 동일 구조 | Manager 6 + Expert 8 + Builder 3 + Team 8 + Decision Tree |
| Common Patterns | ~20줄 | 동일 | Parallel/Sequential/Resume(체크리스트 기반)/Context Propagation |
| Additional Resources | ~15줄 | 동일 | workflows/*.md + rules 참조 목록 |
| Execution Directive | ~40줄 | 10단계 | 8단계로 축소 (Task Tracking/Completion Marker 제거, 체크리스트로 대체) |

#### Mode Router (Do 고유 섹션)

```markdown
## Mode Router

Before routing to a workflow, determine current execution mode:

1. Check DO_MODE environment variable (set by godo mode command)
2. If not set, default to "do" mode

### Mode-Specific Behavior

- **Do Mode** (DO_MODE=do): Full delegation. ALL implementation via Task().
  Response prefix: [Do]
- **Focus Mode** (DO_MODE=focus): Direct execution. Read/Write/Edit directly.
  Response prefix: [Focus]
- **Team Mode** (DO_MODE=team): Agent Teams API. TeamCreate/SendMessage.
  Response prefix: [Team]

### Auto-Escalation Rules

- Focus → Do: 5+ files, multi-domain, expert analysis needed
- Do → Team: 10+ files, 3+ domains, parallel research beneficial
```

#### Intent Router (Do 버전)

```markdown
## Intent Router

Parse $ARGUMENTS to determine which workflow to execute.

### Priority 1: Explicit Subcommand Matching
- **plan**: Analysis → Architecture → Plan pipeline
- **run**: Checklist-based implementation
- **checklist**: Checklist management (create/view/update)
- **mode** [do|focus|team]: Execute godo mode <mode>
- **style** [sprint|pair|direct]: Style switching
- **setup**: Initial configuration wizard
- **check**: Installation diagnostics

### Priority 2: Mode Switch Detection
- Korean/English mode keywords → Execute godo mode <mode> [HARD]
- "포커스", "Focus", "Do 모드", "팀 모드" etc.

### Priority 3: Natural Language Classification
- Design/plan language → plan workflow
- Implementation language → run workflow (check checklist first)
- Bug/error language → expert-debug delegation
- Domain keywords → corresponding expert agent

### Priority 4: Default Behavior
- Ambiguous → AskUserQuestion with top 2-3 options
- Clear development task → do workflow (plan → run pipeline)
```

#### Execution Directive (Do 버전, 8단계)

```
Step 1 - Parse Arguments: Extract subcommand + flags from $ARGUMENTS
Step 2 - Check Mode: Read DO_MODE, verify statusline matches response prefix
Step 3 - Route to Workflow: Apply Intent Router (Priority 1-4)
Step 4 - Read Config: Load DO_* environment variables from settings.local.json
Step 5 - Execute Workflow: Read workflows/<name>.md, delegate to agents
Step 6 - Update Checklist: Mark items [ ]→[~]→[*]→[o] as work progresses
Step 7 - Present Results: Display in user's language with persona style
Step 8 - Guide Next Steps: AskUserQuestion for logical next actions
```

### 3.2 Workflow Files

#### workflows/plan.md (~200줄)

moai plan.md (269줄)을 체크리스트 체계로 재설계.

```markdown
# Do Plan Workflow

## 목적
사용자 요청을 분석하여 .do/jobs/{YYMMDD}/{title}/plan.md 생성

## 복잡도 판단 (dev-workflow.md 참조)
- 단순: 4개 이하 파일 → Plan만 생성
- 복잡: 5개+ 파일, 신규 모듈, 마이그레이션 → Analysis → Architecture → Plan

## Phase 1: Complexity Assessment
- 사용자 요청 분석
- 파일 변경 범위 추정
- 복잡/단순 판단 → 불확실하면 AskUserQuestion

## Phase 2A: Analysis (복잡한 작업만)
- Task(expert-analyst): 현황 조사 + 요구사항 + 기술 선택지
- 산출물: .do/jobs/{YYMMDD}/{title}/analysis.md
- 템플릿: dev-checklist.md Analysis 템플릿 준수

## Phase 2B: Architecture (복잡한 작업만)
- Task(expert-architect): 솔루션 설계 + 인터페이스 명세
- 입력: analysis.md
- 산출물: .do/jobs/{YYMMDD}/{title}/architecture.md
- 템플릿: dev-checklist.md Architecture 템플릿 준수

## Phase 3: Plan Generation
- Task(plan-agent): analysis + architecture 기반 작업 계획
- 산출물: .do/jobs/{YYMMDD}/{title}/plan.md
- TDD 여부 AskUserQuestion으로 확인

## Phase 4: User Approval
- AskUserQuestion: "설계 완료! 구현 진행할까요?"
- 승인 시 → checklist 생성 안내
- 거부 시 → 수정 사항 수집 후 Phase 3 반복

## Plan Mode (Shift+Tab) Integration
- [HARD] 저장 위치: .do/jobs/{YYMMDD}/{title}/plan.md
- [HARD] ~/.claude/plans/ 절대 금지
```

#### workflows/run.md (~250줄)

moai run.md (361줄)을 체크리스트 기반으로 재설계.

```markdown
# Do Run Workflow

## 목적
체크리스트 기반으로 구현 실행. 에이전트에게 서브 체크리스트를 전달하고 작업 추적.

## 전제조건
- plan.md 존재
- checklist.md 존재 (없으면 생성 안내)

## Phase 1: Checklist Verification
- checklist.md 읽기
- 미완료 항목 확인 ([~], [ ])
- 의존성 확인 (depends on: 해소 여부)

## Phase 2: Agent Dispatch
- 서브 체크리스트별 에이전트 할당
- 에이전트 호출 시 필수 전달사항 (dev-workflow.md):
  1. 작업 지시
  2. 서브 체크리스트 경로
  3. Docker 환경 정보
  4. 커밋 지시

### Do Mode Dispatch
- 독립 작업 → 병렬 Task() 호출
- 의존 작업 → 순차 실행

### Focus Mode Dispatch
- 직접 코드 작성 (Task 위임 안 함)
- 순차적 항목 처리

### Team Mode Dispatch
- TeamCreate → 팀원 Spawn
- 파일 소유권 기반 작업 분배
- SendMessage로 조율

## Phase 3: Progress Monitoring
- 에이전트 완료 후 체크리스트 상태 확인
- 미커밋 변경 있으면 에이전트 재호출
- 토큰 소진 시 새 에이전트로 재개

## Phase 4: Quality Verification
- 전체 테스트 스위트 실행
- dev-testing.md 규칙 준수 확인

## Phase 5: Completion
- checklist.md 전체 상태 표시
- 미완료 항목 있으면 다음 행동 제안
- 전체 완료 시 report.md 작성 안내
```

#### workflows/do.md (~100줄)

moai moai.md (219줄)의 자동 파이프라인 축소판.

```markdown
# Do Default Workflow (plan → run)

## 목적
plan → checklist → run을 자동 파이프라인으로 실행

## Phase 0: Exploration
- 병렬: 코드베이스 분석 + 요구사항 파악
- Explore subagent 또는 직접 Grep/Glob

## Phase 1: Plan
- Read workflows/plan.md 실행
- 사용자 승인 대기

## Phase 2: Checklist
- plan.md 기반 checklist.md 자동 생성
- 에이전트별 서브 체크리스트 생성

## Phase 3: Run
- Read workflows/run.md 실행
- 모드별(Do/Focus/Team) 분기

## Phase 4: Report
- 완료 보고서 생성
- 다음 단계 안내
```

#### workflows/team-plan.md (~80줄)

moai team-plan.md (92줄)과 유사 구조.

```markdown
# Do Team Plan Workflow

## 목적
Team 모드에서 병렬 조사 팀으로 Plan 생성

## Team Composition
- team-researcher: 코드베이스 탐색 (haiku, plan mode)
- team-analyst: 요구사항 분석 (inherit, plan mode)
- team-architect: 기술 설계 (inherit, plan mode)

## Execution
1. TeamCreate
2. Spawn 3 teammates with investigation prompts
3. 병렬 조사 결과 수집
4. 결과 종합 → plan.md 생성
5. Shutdown team
```

#### workflows/team-run.md (~120줄)

moai team-run.md (179줄)을 체크리스트 체계로 적용.

```markdown
# Do Team Run Workflow

## 목적
Team 모드에서 병렬 구현 팀으로 Checklist 실행

## Team Composition
- team-backend-dev: 서버 구현 (acceptEdits)
- team-frontend-dev: 클라이언트 구현 (acceptEdits)
- team-tester: 테스트 작성 (acceptEdits)
- team-quality: 품질 검증 (plan, read-only)

## File Ownership
- 체크리스트 서브 파일 기반 파일 소유권 배정
- Critical Files 섹션의 수정 대상을 소유권 근거로 사용

## Execution
1. TeamCreate
2. 태스크 분해 (checklist.md → 팀원별)
3. Spawn teammates with sub-checklist paths
4. 각 팀원: 서브 체크리스트 읽기 → 구현 → 커밋
5. team-quality: 전체 품질 검증
6. Git staging 규칙 준수 (개별 파일 git add, broad staging 금지)
7. Shutdown team
```

### 3.3 Reference File

#### references/reference.md (~150줄)

moai reference.md (250줄)에서 Do에 필요한 부분만 추출.

```markdown
# Do Reference

## Execution Patterns
- Parallel: 독립 Task() 동시 호출 (최대 10개)
- Sequential: 의존성 있는 작업 순차 실행
- Resume: 체크리스트 [o] 건너뛰고 미완료부터 재개

## Mode Reference
- Do: [Do] prefix, Task() delegation, parallel
- Focus: [Focus] prefix, direct execution, sequential
- Team: [Team] prefix, TeamCreate/SendMessage, parallel

## Auto-Escalation Thresholds
- Focus → Do: files >= 5, domains >= 2, expert needed
- Do → Team: files >= 10, domains >= 3

## Configuration Reference
- settings.json: outputStyle, plansDirectory, hooks, permissions
- settings.local.json: DO_* environment variables
- Environment variables: DO_MODE, DO_USER_NAME, DO_LANGUAGE,
  DO_COMMIT_LANGUAGE, DO_AI_FOOTER, DO_PERSONA

## Artifact Locations
- Plans: .do/jobs/{YYMMDD}/{title}/plan.md
- Checklists: .do/jobs/{YYMMDD}/{title}/checklist.md
- Sub-checklists: .do/jobs/{YYMMDD}/{title}/checklists/{NN}_{agent}.md
- Reports: .do/jobs/{YYMMDD}/{title}/report.md

## Persona System
- young-f: 밝은 20대 여성 천재 개발자, {name}선배
- young-m: 자신감 있는 20대 남성 천재 개발자, {name}선배님
- senior-f: 30년 경력 50대 여성 레전드, {name}님
- senior-m: 업계 전설 50대 남성 시니어, {name}씨

## Do-Specific Features (moai에 없는 것)
- Multirepo: .git.multirepo 파일 기반 작업 위치 확인
- Release: tobrew.lock 감지 시 릴리즈 워크플로우
- Mode Switching: godo mode <mode> (statusline 동기화 필수)
```

---

## 4. Persona Override Skills Decision

### 결정: Override skills 생성하지 않음

**근거:**

| moai override skill | do 대응 | 결정 |
|---|---|---|
| moai-foundation-core | rules/*.md + CLAUDE.md에 인라인 | 불필요 -- 이미 흡수 |
| moai-foundation-quality | dev-testing.md + dev-workflow.md | 불필요 -- rules가 더 상세 |
| moai-workflow-ddd | rules/do/workflow/workflow-modes.md | 불필요 -- 동일 내용 |
| moai-workflow-tdd | dev-workflow.md TDD 섹션 | 불필요 -- 동일 내용 |
| moai-workflow-spec | .do/jobs/ + dev-checklist.md | 불필요 -- 완전히 다른 체계 |
| moai-workflow-project | godo init/setup | 불필요 -- godo가 대체 |

**핵심 인사이트**: moai는 "skill → progressive disclosure"로 지식을 로드하지만, do는 "rules → 항상 로드"로 직접 주입한다. 아키텍처가 근본적으로 다르므로 skill 변환이 아니라 rules 유지가 올바른 전략이다.

**agent_patches도 비워둠**: moai는 agent_patches로 core 에이전트에 override skill을 주입했지만(예: expert-backend에 moai-foundation-core + moai-workflow-tdd + moai-workflow-ddd 추가), do는 이 override skills 자체가 없으므로 agent_patches가 필요 없다. core 에이전트의 rules 참조만으로 충분하다.

---

## 5. Styles / Hooks / Commands

### 5.1 Output Styles

**결정: do 고유 스타일 3종으로 교체**

| moai 스타일 | do 스타일 | 대응 |
|---|---|---|
| moai.md (Strategic Orchestrator) | pair.md (친절한 동료) | 기본값 |
| r2d2.md (Pair Programming) | sprint.md (민첩한 실행자) | 최소 대화 |
| yoda.md (Technical Wisdom) | direct.md (직설적 전문가) | 간결한 전문성 |

**폴더 규칙**: `output-styles/do/` (moai의 `output-styles/moai/`와 동일 패턴)

현재 do-focus에는 `styles/` 폴더에 moai.md, r2d2.md, yoda.md가 있지만 이는 moai에서 복사된 것이다. do persona로 전환 시 do 고유 스타일로 교체한다. 기존 3개 파일의 내용을 do 스타일 이름으로 재매핑:
- moai.md (오케스트레이터 스타일) → pair.md로 재작성 (협업적 톤 유지, moai 참조 제거)
- r2d2.md (페어 프로그래밍) → sprint.md로 재작성 (빠른 실행, 간결한 응답)
- yoda.md (지혜 마스터) → direct.md로 재작성 (직설적 전문가)

### 5.2 Hooks

**결정: shell wrapper 없음. godo binary 직접 호출 유지**

| Event | moai (shell wrapper) | do (직접 호출) |
|---|---|---|
| SessionStart | handle-session-start.sh → moai hook | `godo hook session-start` |
| SessionEnd | handle-session-end.sh → moai hook | `godo hook session-end` |
| PreToolUse | handle-pre-tool.sh → moai hook | `godo hook pre-tool` |
| PostToolUse | handle-post-tool.sh (Write\|Edit) | `godo hook post-tool-use` (.*) |
| PreCompact | handle-compact.sh → moai hook | `godo hook compact` (PostToolUse 두 번째) |
| Stop | handle-stop.sh → moai hook | `godo hook stop` |
| SubagentStop | (없음) | `godo hook subagent-stop` |
| UserPromptSubmit | (없음) | `godo hook user-prompt-submit` |

**manifest.yaml에서**: `hook_scripts: []` (빈 배열). hooks는 settings.json에서 직접 정의.

### 5.3 Commands

**결정: 기존 6개 커맨드 유지**

| 커맨드 | 역할 | 변경 사항 |
|---|---|---|
| `/do:check` | 설치/환경 확인 | 유지 |
| `/do:checklist` | 체크리스트 관리 | 유지 |
| `/do:mode` | 모드 전환 | 유지 |
| `/do:plan` | 플랜 생성 | 유지 |
| `/do:setup` | 초기 설정 | 유지 |
| `/do:style` | 스타일 전환 | 유지 |

moai의 `/moai:github`, `/moai:99-release`는 프로젝트 특화이므로 do에 포함하지 않음.

---

## 6. Manifest Design

```yaml
name: do
version: "3.0.0"
description: "Do execution framework with 3 modes (Do/Focus/Team)"
brand: do
brand_dir: do
brand_cmd: do
claude_md: CLAUDE.md

# Persona-specific agents (core agents are inherited automatically)
agents:
  - agents/do/manager-ddd.md
  - agents/do/manager-project.md
  - agents/do/manager-quality.md
  - agents/do/manager-tdd.md
  - agents/do/team-quality.md

# Orchestrator skill + workflows (핵심 추가)
skills:
  - skills/do/SKILL.md
  - skills/do/workflows/plan.md
  - skills/do/workflows/run.md
  - skills/do/workflows/do.md
  - skills/do/workflows/team-plan.md
  - skills/do/workflows/team-run.md
  - skills/do/references/reference.md

# Persona-specific rules
rules:
  - rules/do/workflow/spec-workflow.md
  - rules/do/workflow/workflow-modes.md

# Do-specific output styles
styles:
  - output-styles/do/sprint.md
  - output-styles/do/pair.md
  - output-styles/do/direct.md

# Do-specific commands
commands:
  - commands/do/check.md
  - commands/do/checklist.md
  - commands/do/mode.md
  - commands/do/plan.md
  - commands/do/setup.md
  - commands/do/style.md

# No shell wrapper scripts needed (godo direct call)
hook_scripts: []

# Settings template
settings:
  outputStyle: pair
  plansDirectory: .do/jobs
  statusLine:
    command: godo statusline
    type: command
  attribution:
    commit: ""
    pr: ""
  hooks:
    SessionStart:
      - hooks:
          - command: godo hook session-start
            type: command
    SessionEnd:
      - hooks:
          - command: godo hook session-end
            type: command
    PreToolUse:
      - hooks:
          - command: godo hook pre-tool
            timeout: 5
            type: command
        matcher: Write|Edit|Bash
    PostToolUse:
      - hooks:
          - command: godo hook post-tool-use
            type: command
        matcher: ".*"
      - hooks:
          - command: godo hook compact
            timeout: 5
            type: command
    Stop:
      - hooks:
          - command: godo hook stop
            type: command
    SubagentStop:
      - hooks:
          - command: godo hook subagent-stop
            type: command
    UserPromptSubmit:
      - hooks:
          - command: godo hook user-prompt-submit
            type: command

# No slot_content overrides (Do uses rules instead of slot injection)
slot_content: {}

# No agent patches (no override skills to inject)
agent_patches: {}
```

### moai manifest와의 차이

| 필드 | moai | do | 이유 |
|------|------|-----|------|
| skills | 16개 (orchestrator) + 6개 (override) | 7개 (orchestrator만) | override skills 불필요 |
| hook_scripts | 7개 .sh | 0개 | godo 직접 호출 |
| settings.hooks | shell wrapper 경로 | godo 직접 명령 | 아키텍처 차이 |
| slot_content | TRUST 5 + TAG Chain | 비어있음 | Do는 TRUST 5 브랜딩 미사용 |
| agent_patches | 20+ 에이전트 패치 | 비어있음 | override skills 없으므로 |
| settings.env | MOAI_CONFIG_SOURCE 등 | (없음, settings.local.json 사용) | Do는 DO_* env 패턴 |

---

## 7. Approach Comparison

### Approach A: Lean CLAUDE.md + Orchestrator Skill (권장)

moai의 패턴을 따라 CLAUDE.md를 ~200줄로 축약하고, 상세 로직을 `skills/do/SKILL.md` + `workflows/*.md`로 분리.

| 항목 | 평가 |
|---|---|
| 복잡도 | 중간 -- 파일 분리 작업 필요 |
| 확장성 | 높음 -- 워크플로우 추가 시 파일만 추가 |
| 토큰 효율 | 높음 -- SKILL.md는 /do 호출 시만 로드 |
| converter 호환 | 완전 -- moai와 동일 manifest 구조 |
| 유지보수 | 좋음 -- 관심사 분리 명확 |

**장점:**
- converter가 moai와 동일한 방식으로 처리 가능
- 워크플로우 추가/삭제가 파일 단위로 독립적
- CLAUDE.md가 짧아서 항상 로드되는 토큰 비용 감소
- skills progressive disclosure 활용 가능

**단점:**
- 초기 파일 분리 작업량이 많음 (12개 파일 신규 생성)
- 기존 do-focus 사용자가 구조 변화에 적응 필요

### Approach B: Fat CLAUDE.md 유지 + Minimal Persona

현재 do-focus CLAUDE.md를 거의 그대로 유지하고, persona에는 최소한의 파일(settings, styles, commands)만 포함.

| 항목 | 평가 |
|---|---|
| 복잡도 | 낮음 -- 현재 구조 유지 |
| 확장성 | 낮음 -- CLAUDE.md가 비대해짐 |
| 토큰 효율 | 낮음 -- 395줄 항상 로드 |
| converter 호환 | 부분적 -- skills 비어있음 |
| 유지보수 | 어려움 -- 하나의 파일에 모든 것 |

**장점:**
- 작업량 최소 (현재 파일 복사 수준)
- 기존 사용자에게 변화 없음

**단점:**
- converter의 핵심 기능(skills 관리)을 활용하지 못함
- moai와 구조적 불일치 (converter가 두 가지 패턴을 지원해야 함)
- CLAUDE.md 40K 제한에 근접할수록 유지보수 어려움

### 결론: Approach A 선택

converter 프로젝트의 목적이 "persona를 교체 가능한 단위로 관리"하는 것이므로, moai와 동일한 구조를 따르는 것이 올바르다. 초기 작업량이 많지만 한번 분리하면 장기적 유지보수가 훨씬 수월하다.

---

## 8. Implementation Order

### Phase 1: Foundation (3개 파일)

1. **`personas/do/manifest.yaml`** -- 전체 구조 정의
   - Section 6의 manifest 초안을 작성
   - 빈 파일 참조도 포함 (아직 생성 안 된 파일)

2. **`personas/do/CLAUDE.md`** -- lean orchestrator
   - 현재 do-focus CLAUDE.md (395줄) → ~200줄로 축약
   - 상세 내용을 skills/rules 참조로 교체
   - 섹션 구조는 Section 2 설계를 따름

3. **`personas/do/settings.json`** -- hooks + permissions
   - 현재 extracted-do settings.json 기반
   - godo 직접 호출 패턴 유지

### Phase 2: Orchestrator Skill (4개 파일)

4. **`personas/do/skills/do/SKILL.md`** -- 핵심 파일
   - Mode Router + Intent Router + Core Rules + Execution Directive
   - moai SKILL.md (377줄) 참조하되 Do 체계로 재작성
   - Section 3.1 설계를 따름

5. **`personas/do/skills/do/workflows/plan.md`**
   - Analysis → Architecture → Plan 3단계
   - Section 3.2 plan.md 설계를 따름

6. **`personas/do/skills/do/workflows/run.md`**
   - Checklist 기반 구현 실행
   - Section 3.2 run.md 설계를 따름

7. **`personas/do/skills/do/workflows/do.md`**
   - plan → run 자동 파이프라인
   - Section 3.2 do.md 설계를 따름

### Phase 3: Team Workflows + Reference (3개 파일)

8. **`personas/do/skills/do/workflows/team-plan.md`**
   - Team 모드 Plan 워크플로우

9. **`personas/do/skills/do/workflows/team-run.md`**
   - Team 모드 Run 워크플로우

10. **`personas/do/skills/do/references/reference.md`**
    - 공통 패턴, 설정 참조, 환경변수

### Phase 4: Styles + Commands + Rules (11개 파일)

11. **`personas/do/output-styles/do/pair.md`** -- 기본 스타일
12. **`personas/do/output-styles/do/sprint.md`**
13. **`personas/do/output-styles/do/direct.md`**
14. **`personas/do/commands/do/check.md`** -- 기존 복사
15. **`personas/do/commands/do/checklist.md`** -- 기존 복사
16. **`personas/do/commands/do/mode.md`** -- 기존 복사
17. **`personas/do/commands/do/plan.md`** -- 기존 복사
18. **`personas/do/commands/do/setup.md`** -- 기존 복사
19. **`personas/do/commands/do/style.md`** -- 기존 복사
20. **`personas/do/rules/do/workflow/spec-workflow.md`** -- 기존 복사 + Do 체계 수정
21. **`personas/do/rules/do/workflow/workflow-modes.md`** -- 기존 복사

### Phase 5: Agents + Verification (6개 파일)

22. **`personas/do/agents/do/manager-ddd.md`** -- 기존 복사
23. **`personas/do/agents/do/manager-project.md`** -- 기존 복사
24. **`personas/do/agents/do/manager-quality.md`** -- 기존 복사
25. **`personas/do/agents/do/manager-tdd.md`** -- 기존 복사
26. **`personas/do/agents/do/team-quality.md`** -- 기존 복사
27. **Assemble 검증** -- converter assemble 실행하여 .claude/ 디렉토리 생성 확인

### 총 파일 수: 26개 (신규 작성 10개 + 기존 기반 수정 16개)

| 유형 | 파일 수 | 작업 난이도 |
|------|---------|-----------|
| 신규 작성 (핵심) | 10개 (CLAUDE.md, manifest, SKILL.md, 5 workflows, reference, settings) | 높음 |
| 기존 기반 수정 | 3개 (styles 3종) | 중간 |
| 기존 복사 | 13개 (commands 6, rules 2, agents 5) | 낮음 |

---

## 9. Risk Mitigation

| Risk | Impact | Mitigation |
|---|---|---|
| CLAUDE.md 축약 시 핵심 규칙 누락 | HIGH | 축약 전/후 HARD 규칙 체크리스트 대조. rules 참조가 실제로 로드되는지 검증 |
| SKILL.md Intent Router가 기존 동작과 불일치 | HIGH | 현재 CLAUDE.md의 Intent-to-Agent Mapping 테이블을 1:1 매핑 검증 |
| converter assemble 시 skills 경로 오류 | MEDIUM | Phase 5에서 assemble 실행 후 .claude/skills/do/SKILL.md 존재 확인 |
| godo hook이 새 구조에서 경로 변경 영향 | MEDIUM | hooks는 settings.json에서 godo 직접 호출이므로 경로 무관. 변경 없음 |
| output-styles 경로 변경 (styles/ → output-styles/do/) | LOW | manifest.yaml의 styles 필드가 assembler에게 정확한 경로 제공 |
| slot_content 비워둠에 따른 core 파일 미치환 | MEDIUM | assemble 후 grep -r "{{slot:" .claude/ 실행하여 미치환 슬롯 확인 |
| manager-spec 제거에 따른 참조 깨짐 | LOW | SKILL.md에서 manager-spec 대신 plan workflow 내 에이전트 위임으로 대체 |

### slot_content 관련 추가 조사 필요

현재 moai manifest의 slot_content에는 `QUALITY_FRAMEWORK`, `QUALITY_GATE_TEXT`, `TRACEABILITY_SYSTEM`이 있다. core 파일 중 `{{slot:QUALITY_FRAMEWORK}}` 등을 참조하는 파일이 있을 수 있다. do persona에서 slot_content를 비우면 이 슬롯이 미치환 상태로 남을 수 있다.

**완화책**: assemble 후 `grep -r "{{slot:" .claude/` 실행하여 미치환 슬롯 확인. 발견 시 do 버전의 slot_content를 추가한다. 가능한 대체:
- `QUALITY_FRAMEWORK` → dev-testing.md + dev-workflow.md 참조 텍스트
- `QUALITY_GATE_TEXT` → "code quality rules defined in dev-testing.md"
- `TRACEABILITY_SYSTEM` → 제거 또는 간소화

---

**작성자**: architect agent
**검토 상태**: Architecture 설계 완료
**다음 단계**: Plan 생성 (implementation order 기반 체크리스트)
