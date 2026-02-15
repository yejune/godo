# Do Persona Identity Document

## Executive Summary

> **관련 문서**: moai 업데이트 시 Do 정체성 검증 절차는 [RUNBOOK.md](./RUNBOOK.md) Section 9 참조.

**Do는 "말하면 한다"의 실행 철학을 가진, 한국어 문화에 뿌리를 둔 3모드 적응형 오케스트레이터이다.**

MoAI가 "전략적 오케스트레이터"(Strategic Orchestrator)로서 SPEC 기반 자율 실행을 추구한다면, Do는 "실행하는 자"(The Doer)로서 사용자의 의도를 즉시 현실로 변환한다. Do의 이름 자체가 동사이며, 선언문이 명령형이고, 모든 설계 결정이 "실행"이라는 하나의 원칙으로 수렴한다.

---

## 1. Core Philosophy -- Do만의 고유 철학

### 1.1 Action-First Identity (실행 우선 정체성)

Do의 정체성은 한국어 동사 "하다"(to do)에서 온다. 시스템이 **무엇인지**(what it is)가 아니라 **무엇을 하는지**(what it does)를 정의한다.

세 가지 선언문이 이를 보여준다:
- **Do**: "나는 Do다. 말하면 한다." -- 전략적 오케스트레이터
- **Focus**: "나는 Focus다. 집중해서 한다." -- 집중하는 실행자
- **Team**: "나는 Team이다. 팀을 이끈다." -- 병렬 팀 리더

MoAI의 선언문이 "MoAI is the Strategic Orchestrator for Claude Code"라는 영어 서술문인 것과 대조적이다. Do의 선언문은 한국어 명령형이며, 주어가 "나"(I)로 시작하여 1인칭 능동태를 취한다.

### 1.2 Appropriate Force Application (적정 실행력 원칙)

Do의 가장 핵심적인 철학적 차이는 **삼원 실행 구조**(三源 實行 構造)이다. 모든 작업에 전체 오케스트레이션 기계를 가동하지 않는다.

- **Focus**: 간단한 작업(1-3 파일)은 직접 코드를 작성한다. 위임 오버헤드 불필요.
- **Do**: 복잡한 작업(5-10 파일)은 전문 에이전트에게 위임한다.
- **Team**: 대규모 작업(10+ 파일)은 Agent Teams API로 병렬 팀을 구성한다.

MoAI는 항상 위임한다(always delegate). Do는 작업 규모에 맞게 적응한다(adapt to scale).

### 1.3 Checklist as Architecture (체크리스트가 곧 아키텍처)

Do에서 체크리스트는 단순한 문서가 아니다. **에이전트의 영속 상태 저장소**이다.

MoAI는 SPEC 문서 + TaskCreate/TaskUpdate로 상태를 관리한다. Do는 파일 기반 체크리스트(`[ ] [~] [*] [!] [o] [x]`)로 에이전트 간 작업 연속성을 보장한다. 이 6개 상태 기호는 표준 마크다운 체크박스와 의도적으로 다르며, 엄격한 상태 전이 규칙을 가진다.

### 1.4 Commit as Proof (커밋이 곧 증거)

Do에서 작업 완료의 증거는 git 커밋이다. 코드를 작성하고, 테스트를 통과하고, 커밋해야만 `[o]` 완료 상태로 전환할 수 있다. 체크리스트 Progress Log에 커밋 해시가 기록되지 않으면 작업은 미완료이다.

### 1.5 Korean Cultural Roots (한국 문화적 뿌리)

Do의 정체성은 한국어와 한국 직장 문화에 깊이 뿌리를 두고 있다:
- 선언문이 한국어 명령형
- 페르소나가 한국 직장 호칭 체계(선배/님/씨) 사용
- 말투가 반말+존댓말 혼합 (한국어 특유의 격식 조절)
- 기본 대화 언어가 한국어 (`DO_LANGUAGE=ko`)
- 페르소나 성격이 한국 직장 관계 역학 반영

MoAI도 "Korean-First, English-Always"를 표방하지만, 지시문(CLAUDE.md, agents, skills)은 모두 영어다. Do는 CLAUDE.md 자체가 한국어+영어 혼합이며, 한국어가 단순한 번역이 아니라 설계 언어이다.

---

## 2. Shared DNA -- MoAI에서 상속받은 것

Do는 MoAI-ADK의 core 계층을 그대로 상속받는다. 이들은 converter의 `core/` 디렉토리에 있으며, 기계적으로 공유된다.

### 2.1 에이전트 카탈로그 (22개 core agents)

| 카테고리 | 에이전트 수 | 에이전트 목록 |
|---------|-----------|-------------|
| Builder | 3 | builder-agent, builder-plugin, builder-skill |
| Expert | 9 | expert-backend, expert-chrome-extension, expert-debug, expert-devops, expert-frontend, expert-performance, expert-refactoring, expert-security, expert-testing |
| Manager | 3 | manager-docs, manager-git, manager-strategy |
| Team | 7 | team-analyst, team-architect, team-backend-dev, team-designer, team-frontend-dev, team-researcher, team-tester |

### 2.2 스킬 시스템 (40+ core skills)

모든 `do-*` 접두사 스킬이 core에서 온다:
- `do-foundation-*`: claude, context, philosopher
- `do-domain-*`: backend, frontend, database, uiux
- `do-lang-*`: 16개 프로그래밍 언어
- `do-library-*`, `do-platform-*`, `do-tool-*`, `do-framework-*` 등

### 2.3 공유 Rules

- `dev-environment.md`: Docker 필수, bootapp 도메인, .env 금지
- `dev-testing.md`: Real DB Only, AI 안티패턴 금지, FIRST 원칙
- `dev-workflow.md`: 복잡도 판단, Read Before Write, 에러 대응
- `dev-checklist.md`: 체크리스트 시스템, 상태 기호, 서브 체크리스트 템플릿
- `file-reading.md`: 파일 읽기 최적화 (4-tier 시스템)

### 2.4 TRUST 5 Quality Framework

- **T**ested: 85%+ coverage, characterization tests
- **R**eadable: Clear naming, English comments
- **U**nified: Consistent style
- **S**ecured: OWASP compliance, input validation
- **T**rackable: Conventional commits, issue references

Do는 TRUST 5라는 "브랜드명"을 사용하지 않지만, 동일한 품질 기준을 `dev-testing.md` + `dev-workflow.md`의 [HARD] 규칙으로 적용한다.

### 2.5 Development Methodology (DDD/TDD/Hybrid)

- DDD (ANALYZE-PRESERVE-IMPROVE): 레거시 리팩토링
- TDD (RED-GREEN-REFACTOR): 신규 기능 개발
- Hybrid: 혼합 (신규 코드는 TDD, 기존 코드는 DDD)

### 2.6 Progressive Disclosure (3-level)

스킬의 토큰 최적화:
- Level 1 (~100 tokens): 메타데이터만
- Level 2 (~5000 tokens): 스킬 본문
- Level 3 (variable): 번들 파일 (on-demand)

### 2.7 Agent Authoring / Skill Authoring 규격

`.claude/rules/do/development/` 하위의 agent-authoring.md, skill-authoring.md, coding-standards.md는 core에서 온다.

---

## 3. Unique Features -- Do만의 고유 기능

### 3.1 삼원 실행 구조 (Do/Focus/Team)

MoAI에는 없는 Do 고유 기능. 세 가지 실행 모드를 상황에 따라 전환한다.

| 모드 | 접두사 | 실행 방식 | 병렬성 | 적합 시나리오 |
|------|--------|----------|--------|-------------|
| Do | `[Do]` | 전문 에이전트 위임 | 항상 병렬 | 5+ 파일, 멀티 도메인 |
| Focus | `[Focus]` | 직접 코드 작성 | 순차적 | 1-3 파일, 단일 도메인 |
| Team | `[Team]` | Agent Teams API | 팀 병렬 | 10+ 파일, 3+ 도메인 |

자동 에스컬레이션:
- Focus -> Do: 5+ 파일, 멀티 도메인, 전문가 분석 필요
- Do -> Team: 10+ 파일, 3+ 도메인

### 3.2 페르소나 시스템 (4종 캐릭터)

MoAI에는 없는 Do 고유 기능. `DO_PERSONA` 환경변수로 선택.

| 페르소나 | 설명 | 호칭 |
|---------|------|------|
| `young-f` (기본) | 밝고 에너지 넘치는 20대 여성 천재 개발자 | {name}선배 |
| `young-m` | 자신감 넘치는 20대 남성 천재 개발자 | {name}선배님 |
| `senior-f` | 30년 경력의 레전드 50대 여성 천재 개발자 | {name}님 |
| `senior-m` | 업계 전설의 50대 남성 시니어 아키텍트 | {name}씨 |

SessionStart hook에서 `godo hook session-start`가 페르소나를 시스템 메시지로 주입한다.

### 3.3 스타일 시스템 (3종)

Do 고유. `DO_STYLE` 환경변수 또는 `/do:style` 커맨드로 선택.

| 스타일 | 설명 |
|--------|------|
| sprint | 민첩한 실행자 -- 말 최소화, 바로 실행 |
| pair (기본) | 친절한 동료 -- 협업적 톤 |
| direct | 직설적 전문가 -- 군더더기 없음 |

**페르소나 vs 스타일**: 독립적 축. 페르소나는 "누가" 말하는지, 스타일은 "어떻게" 말하는지. 어떤 페르소나든 어떤 스타일이든 조합 가능.

### 3.4 godo CLI Binary (직접 호출 패턴)

Do는 Go 바이너리(`godo`)를 hook에서 직접 호출한다. MoAI는 7개 shell wrapper 스크립트를 통해 `moai` 바이너리를 간접 호출한다.

```
MoAI: settings.json -> .claude/hooks/moai/*.sh -> moai binary
Do:   settings.json -> godo binary (직접)
```

Do가 MoAI보다 2개 더 많은 hook 이벤트를 처리한다:
- MoAI에 없는 것: SubagentStop, UserPromptSubmit
- Do의 PostToolUse matcher: `.*` (모든 도구) vs MoAI의 `Write|Edit` (쓰기만)

### 3.5 모드 전환 (godo mode)

`godo mode set <mode>` 명령으로 즉시 전환. statusline과 AI 응답 접두사가 동기화되어야 한다.
- [HARD] 모드 전환 시 반드시 `godo mode <mode>` 실행
- [HARD] 실행 없이 접두사만 바꾸는 것은 VIOLATION

### 3.6 Jobs 디렉토리 (날짜 기반 산출물)

`.do/jobs/{YYMMDD}/{title-kebab-case}/` 구조로 모든 작업 산출물 관리.
- MoAI의 `.moai/specs/SPEC-XXX/` (번호 기반)과 대조적.
- 날짜 기반이므로 시간순 탐색이 자연스럽다.
- `~/.claude/plans/` 전역 디렉토리 사용 금지.

### 3.7 Characters 디렉토리

`personas/do/characters/`에 4종 캐릭터 정의 파일. MoAI에는 캐릭터 시스템 자체가 없다.

### 3.8 6개 Do 전용 커맨드

| 커맨드 | 목적 |
|--------|------|
| `/do:check` | 설치/환경 확인 |
| `/do:checklist` | 체크리스트 생성/관리 |
| `/do:mode` | Do/Focus/Team 모드 전환 |
| `/do:plan` | 플랜 생성 |
| `/do:setup` | 사용자 설정 (이름, 언어, 페르소나) |
| `/do:style` | 출력 스타일 전환 |

MoAI는 `/moai` 단일 진입점 + 서브커맨드, Do는 개별 커맨드 방식.

### 3.9 Multirepo 지원

`.git.multirepo` 파일 감지 시 작업 위치를 사용자에게 확인한다. MoAI에는 없는 기능.

### 3.10 릴리즈 워크플로우

`tobrew.lock` 파일 감지 시 릴리즈 프로세스를 제안한다.

### 3.11 설정 아키텍처 차이

Do는 `settings.local.json`의 `DO_*` 환경변수로 설정을 관리한다.

| 변수 | 설명 | 기본값 |
|-----|------|-------|
| `DO_MODE` | 실행 모드 | do |
| `DO_USER_NAME` | 사용자 이름 | "" |
| `DO_LANGUAGE` | 대화 언어 | ko |
| `DO_COMMIT_LANGUAGE` | 커밋 메시지 언어 | en |
| `DO_AI_FOOTER` | AI 푸터 | false |
| `DO_PERSONA` | 페르소나 타입 | young-f |
| `DO_STYLE` | 출력 스타일 | pair |

MoAI는 `.moai/config/sections/*.yaml` (YAML 분리)로 설정을 관리한다.

---

## 4. Philosophy Diff Table

| # | 항목 | MoAI | Do | 비고 |
|---|------|------|-----|------|
| 1 | **정체성 선언** | "MoAI is the Strategic Orchestrator" (영어 서술문) | "나는 Do다. 말하면 한다." (한국어 명령형) | Do는 1인칭 능동태 |
| 2 | **실행 모드** | 단일 모드 (항상 위임) | 3모드 (Do/Focus/Team) 적응형 | Do의 핵심 차별점 |
| 3 | **위임 원칙** | "All tasks must be delegated" (무조건 위임) | "작업 규모에 맞게 위임" (적정 실행력) | Focus 모드에서 직접 실행 |
| 4 | **워크플로우 체계** | SPEC 기반 (Plan/Run/Sync) | 체크리스트 기반 (Plan/Checklist/Develop/Test/Report) | SPEC vs Checklist |
| 5 | **상태 관리** | TaskCreate/TaskUpdate (세션 스코프) | 파일 기반 체크리스트 (영속) | 체크리스트가 에이전트 상태 파일 |
| 6 | **완료 증거** | `<moai>DONE</moai>` XML 마커 | 커밋 해시 (`[o]` + commit hash) | 커밋이 곧 증거 |
| 7 | **요구사항 형식** | EARS (Easy Approach to Requirements Syntax) | MoSCoW (MUST/SHOULD/COULD/WON'T) | 형식론 차이 |
| 8 | **산출물 위치** | `.moai/specs/SPEC-XXX/` (번호 기반) | `.do/jobs/{YYMMDD}/{title}/` (날짜 기반) | 시간순 vs 번호순 |
| 9 | **지시문 언어** | 영어 단일 언어 | 한국어+영어 혼합 | Do의 CLAUDE.md는 이중 언어 |
| 10 | **설정 저장소** | `.moai/config/sections/*.yaml` | `settings.local.json` `DO_*` env | YAML 분리 vs JSON 환경변수 |
| 11 | **Hook 아키텍처** | Shell wrapper (7개 .sh) -> binary | godo binary 직접 호출 (0개 .sh) | Do가 더 간결 |
| 12 | **Hook 이벤트 수** | 5개 | 7개 (+SubagentStop, UserPromptSubmit) | Do가 2개 더 많음 |
| 13 | **PostToolUse 범위** | `Write\|Edit` (쓰기만) | `.*` (모든 도구) | Do가 더 광범위 감시 |
| 14 | **페르소나 시스템** | 없음 | 4종 캐릭터 (young-f/m, senior-f/m) | Do 고유 |
| 15 | **스타일 시스템** | 3종 (moai, r2d2, yoda) | 3종 (sprint, pair, direct) | 이름과 성격이 다름 |
| 16 | **스타일 성격** | moai=오케스트레이터, r2d2=페어, yoda=교육 | sprint=빠름, pair=협업, direct=간결 | 다른 설계 철학 |
| 17 | **품질 브랜딩** | TRUST 5 (명시적 브랜드) | rules의 [HARD] 규칙 (암묵적) | 동일 기준, 다른 표현 |
| 18 | **토큰 예산 계획** | 200K 명시 분배 (30K+180K+40K) | 없음 (file-reading 최적화만) | MoAI가 더 체계적 |
| 19 | **진입점 구조** | `/moai` 단일 + Intent Router | `/do:*` 6개 개별 커맨드 | 통합 vs 분리 |
| 20 | **CLI 도구** | `moai` binary (Go) | `godo` binary (Go) | 이름만 다름 |
| 21 | **StatusLine** | `.moai/status_line.sh` (shell wrapper) | `godo statusline` (직접 명령) | Do가 직접 호출 |
| 22 | **모드 전환** | 없음 (단일 모드) | `godo mode set <mode>` + statusline 동기화 | Do 고유 |
| 23 | **복잡도 판단** | SPEC 기반 (자동/수동) | 파일 수 + 도메인 수 기반 명시 규칙 | Do가 더 명시적 |
| 24 | **에러 핸들링** | 에러 유형별 에이전트 매핑 (5종) | 3회 재시도 + 접근법 재검토 | MoAI가 더 세분화 |
| 25 | **MCP 통합 문서** | 독립 섹션 (4개 MCP) | 없음 (미기술) | MoAI가 더 상세 |
| 26 | **Override Skills** | 6개 | 0개 (rules에 흡수) | 아키텍처 차이 |
| 27 | **지식 주입 방식** | skill -> progressive disclosure | rules -> 항상 로드 | 근본적 아키텍처 차이 |
| 28 | **agent_patches** | 20+ 에이전트에 override skill 주입 | 비어있음 | Do는 rules 참조로 충분 |
| 29 | **Multirepo 지원** | 없음 | `.git.multirepo` 파일 기반 | Do 고유 |
| 30 | **릴리즈 워크플로우** | sync 단계에서 처리 | tobrew 전용 독립 워크플로우 | Do 고유 |

---

## 5. Terminology Map

| MoAI 용어 | Do 용어 | 설명 |
|-----------|---------|------|
| MoAI | Do | 브랜드명 |
| `.moai/` | `.do/` | 프로젝트 설정/상태 디렉토리 |
| `moai` (CLI) | `godo` (CLI) | Go 바이너리 CLI 도구 |
| `/moai` | `/do:*` (개별 커맨드) | 슬래시 커맨드 진입점 |
| `moai-` (skill prefix) | `do-` (skill prefix) | 스킬 네이밍 접두사 |
| `agents/moai/` | `agents/do/` | 페르소나 에이전트 디렉토리 |
| `output-styles/moai/` | `output-styles/do/` | 출력 스타일 디렉토리 |
| `rules/moai/` | `rules/do/` | 페르소나 규칙 디렉토리 |
| `commands/moai/` | `commands/do/` | 커맨드 디렉토리 |
| `hooks/moai/*.sh` | (없음 -- godo 직접 호출) | Hook 스크립트 |
| `moai-constitution.md` | `do-constitution.md` | 핵심 원칙 파일 |
| SPEC | Plan (+ Checklist) | 워크플로우 산출물 |
| SPEC-XXX | `.do/jobs/{YYMMDD}/{title}/` | 산출물 식별자/경로 |
| EARS | MoSCoW | 요구사항 형식 |
| TAG Chain | Checklist dependency (`depends on:`) | 작업 의존성 |
| `<moai>DONE</moai>` | `[o]` + commit hash | 완료 마커 |
| `<moai>COMPLETE</moai>` | report.md 작성 완료 | 전체 완료 마커 |
| Plan/Run/Sync | Plan/Checklist/Develop/Test/Report | 워크플로우 단계 |
| `.moai/specs/` | `.do/jobs/` | 산출물 루트 디렉토리 |
| `.moai/config/sections/*.yaml` | `settings.local.json` `DO_*` env | 설정 저장소 |
| `.moai/learning/` | (없음) | Yoda 스타일 학습 디렉토리 |
| `MOAI_DEVELOPMENT_MODE` | 사용자에게 직접 질문 (AskUserQuestion) | DDD/TDD 선택 방식 |
| `MOAI_CONFIG_SOURCE` | (없음) | MoAI 전용 설정 |
| `moai hook <event>` | `godo hook <event>` | Hook 명령 |
| `moai statusline` | `godo statusline` | StatusLine 명령 |
| `moai mode` | `godo mode set <mode>` | 모드 전환 (Do 고유) |
| Snapshot/Resume | 체크리스트 상태 (영속) | 중단/재개 메커니즘 |
| `manager-spec` | (없음 -- plan workflow로 대체) | SPEC 문서 생성 에이전트 |
| Completion Marker | Checklist status `[o]` | 작업 완료 표시 |
| Philosopher Framework | (core manager-strategy에 포함) | 전략적 사고 프레임워크 |
| UltraThink | (core MCP로 사용 가능) | Sequential Thinking MCP |
| outputStyle: "MoAI" | outputStyle: "pair" | 기본 출력 스타일 |
| moai.md (style) | pair.md | 기본 스타일 파일 |
| r2d2.md (style) | sprint.md | 빠른 실행 스타일 |
| yoda.md (style) | direct.md | 간결한 전문가 스타일 |

---

## 6. Identity Boundaries -- 절대 바꾸면 안 되는 것

moai가 업데이트되더라도 다음 Do 정체성 요소는 **절대 변경하면 안 된다**.

### 6.1 삼원 실행 구조 (CRITICAL)

- **파일**: `personas/do/CLAUDE.md` Section 1 (Do/Focus/Team 정의)
- **파일**: `personas/do/skills/do/SKILL.md` Mode Router 섹션
- **패턴**: `[Do]`, `[Focus]`, `[Team]` 접두사
- **패턴**: "나는 Do다. 말하면 한다." / "나는 Focus다. 집중해서 한다." / "나는 Team이다. 팀을 이끈다."
- **로직**: 자동 에스컬레이션 (Focus->Do->Team)
- **이유**: Do의 가장 근본적인 차별점. MoAI의 단일 모드와 완전히 다른 설계.

### 6.2 체크리스트 시스템 (CRITICAL)

- **파일**: `rules/dev-checklist.md` (core rule이지만 Do의 핵심 워크플로우)
- **파일**: `personas/do/skills/do/workflows/run.md` (Checklist 기반 실행)
- **패턴**: `[ ] [~] [*] [!] [o] [x]` 상태 기호
- **패턴**: `.do/jobs/{YYMMDD}/{title}/checklist.md` 경로
- **패턴**: `checklists/{NN}_{agent}.md` 서브 체크리스트
- **로직**: 상태 전이 규칙 (`[ ]`->`[~]`->`[*]`->`[o]`)
- **로직**: 커밋 해시 = 완료 증거
- **이유**: SPEC 기반이 아닌 체크리스트 기반이 Do의 핵심 워크플로우.

### 6.3 페르소나 시스템 (HIGH)

- **파일**: `personas/do/characters/*.md` (4개 파일)
- **파일**: `personas/do/CLAUDE.md` 페르소나 시스템 섹션
- **설정**: `DO_PERSONA` 환경변수
- **패턴**: `{name}선배`, `{name}선배님`, `{name}님`, `{name}씨`
- **로직**: SessionStart hook에서 페르소나 주입
- **이유**: MoAI에 없는 Do 고유 기능.

### 6.4 스타일 시스템 (HIGH)

- **파일**: `personas/do/output-styles/do/sprint.md`, `pair.md`, `direct.md`
- **설정**: `DO_STYLE` 환경변수
- **패턴**: sprint/pair/direct (moai의 moai/r2d2/yoda와 다름)
- **이유**: Do 고유 스타일.

### 6.5 godo 직접 호출 패턴 (HIGH)

- **파일**: `personas/do/settings.json` hooks 섹션
- **패턴**: `"command": "godo hook <event>"` (shell wrapper 없음)
- **패턴**: `hook_scripts: []` in manifest.yaml
- **이유**: MoAI의 shell wrapper 패턴과 아키텍처적으로 다름.

### 6.6 한국어 선언문 (HIGH)

- **패턴**: "나는 Do다. 말하면 한다."
- **패턴**: "나는 Focus다. 집중해서 한다."
- **패턴**: "나는 Team이다. 팀을 이끈다."
- **이유**: Do의 정체성 그 자체.

### 6.7 6개 개별 커맨드 (MEDIUM)

- **파일**: `personas/do/commands/do/*.md` (6개)
- **패턴**: `/do:check`, `/do:checklist`, `/do:mode`, `/do:plan`, `/do:setup`, `/do:style`
- **이유**: MoAI의 `/moai` 단일 진입점과 다른 설계.

### 6.8 Jobs 디렉토리 경로 (MEDIUM)

- **패턴**: `.do/jobs/{YYMMDD}/{title-kebab-case}/`
- **설정**: `plansDirectory: ".do/jobs"`
- **이유**: MoAI의 `.moai/specs/SPEC-XXX/`와 다른 조직 방식.

### 6.9 DO_* 환경변수 체계 (MEDIUM)

- **패턴**: `DO_MODE`, `DO_USER_NAME`, `DO_LANGUAGE`, `DO_COMMIT_LANGUAGE`, `DO_AI_FOOTER`, `DO_PERSONA`, `DO_STYLE`
- **이유**: MoAI의 `.moai/config/` YAML 체계와 다른 설정 방식.

---

## 7. Conversion Invariants -- 변환 시 보존 규칙

### 7.1 기계적 치환 (Safe to automate)

| 패턴 | 치환 | 범위 |
|------|------|------|
| `agents/moai/` | `agents/do/` | 경로 참조 |
| `.moai/` | `.do/` | 디렉토리 참조 |
| `moai-` (skill prefix) | `do-` | 스킬명 참조 |
| `skills/moai/` | `skills/do/` | 스킬 경로 |
| `commands/moai/` | `commands/do/` | 커맨드 경로 |
| `output-styles/moai/` | `output-styles/do/` | 스타일 경로 |
| `rules/moai/` | `rules/do/` | 규칙 경로 |
| `hooks/moai/` | (삭제) | Hook 경로 |
| `moai hook` | `godo hook` | CLI 명령 |
| `moai statusline` | `godo statusline` | CLI 명령 |

### 7.2 구조적 변환 (Manual review required)

| 영역 | moai 패턴 | do 변환 | 이유 |
|------|----------|---------|------|
| 워크플로우 | SPEC-XXX -> plan/run/sync | plan/checklist/develop/test/report | 완전히 다른 워크플로우 |
| 요구사항 | EARS 형식 | MoSCoW 분류 | 다른 형식론 |
| 상태 추적 | TaskCreate/TaskUpdate | 체크리스트 `[ ][~][*][!][o][x]` | 다른 상태 관리 |
| 완료 마커 | `<moai>DONE</moai>` | `[o]` + commit hash | 다른 완료 증거 |
| 스타일 매핑 | moai/r2d2/yoda | pair/sprint/direct | 이름과 내용이 다름 |
| Shell hooks | 7개 .sh wrapper | godo 직접 호출 | 아키텍처 차이 |
| 설정 | `.moai/config/*.yaml` | `settings.local.json` env | 저장 방식 차이 |
| Override skills | 6개 SKILL.md | (없음 -- rules에 흡수) | 지식 주입 방식 차이 |

### 7.3 절대 건드리면 안 되는 패턴

1. **삼원 선언문**: "나는 Do다. 말하면 한다." 등 3개 선언문
2. **상태 기호**: `[ ] [~] [*] [!] [o] [x]` 6개 기호와 전이 규칙
3. **페르소나 호칭**: `{name}선배`, `{name}선배님`, `{name}님`, `{name}씨`
4. **Jobs 경로 패턴**: `.do/jobs/{YYMMDD}/{title-kebab-case}/`
5. **DO_* 환경변수**: 7개 환경변수명과 의미
6. **godo 직접 호출**: settings.json의 `"command": "godo hook ..."` 패턴
7. **VIOLATION 정의**: Do가 직접 코드 작성, 에이전트 위임 없이 파일 수정 등

---

## 8. Architecture Decisions -- 설계 결정과 근거

### 8.1 왜 3모드인가? (단일 모드 대신)

**결정**: Focus/Do/Team 3모드 적응형
**MoAI와의 차이**: MoAI는 항상 위임
**근거**: 토큰 효율성. CSS 1줄 수정에 에이전트를 호출하면 500+ 토큰 오버헤드. Focus 모드로 직접 수행하면 즉시 완료. Do의 사용자는 혼자 작업하는 개발자이므로 간단한 작업에 대한 빠른 피드백이 중요하다.

### 8.2 왜 체크리스트인가? (SPEC 대신)

**결정**: 파일 기반 체크리스트 시스템
**MoAI와의 차이**: MoAI는 SPEC 문서 + EARS 형식
**근거**: 체크리스트는 에이전트 간 상태 전달에 최적화되어 있다. 에이전트가 토큰을 소진해도 체크리스트 파일에 마지막 상태가 남아있으므로, 새 에이전트가 즉시 이어받을 수 있다. SPEC은 한 번 작성하고 참조하는 문서지만, 체크리스트는 실시간으로 갱신되는 상태 파일이다.

### 8.3 왜 godo 직접 호출인가? (Shell wrapper 대신)

**결정**: settings.json에서 godo binary를 직접 호출
**MoAI와의 차이**: MoAI는 7개 shell wrapper 스크립트
**근거**: 중간 레이어 제거로 디버깅 단순화. Shell 스크립트는 PATH 문제, 인코딩 문제, SIGALRM 문제 등 28가지 이슈를 일으켰다. 직접 호출은 이 모든 문제를 원천 제거한다.

### 8.4 왜 Override Skills가 없는가? (MoAI는 6개)

**결정**: Override skills 생성하지 않음
**MoAI와의 차이**: MoAI는 6개 override skill로 progressive disclosure
**근거**: Do는 "rules -> 항상 로드" 방식으로 지식을 주입한다. 두 계층(skill + rules)을 유지하는 복잡성보다 단일 계층(rules only)의 단순함이 Do에 더 적합하다.

### 8.5 왜 한국어 혼합 CLAUDE.md인가? (영어 단일 대신)

**결정**: CLAUDE.md를 한국어+영어 혼합으로 작성
**MoAI와의 차이**: MoAI는 영어 단일 언어
**근거**: Do의 주 사용자가 한국어 화자이고, 페르소나 시스템이 한국어 문화에 기반하므로, 지시문도 한국어가 자연스럽다.

### 8.6 왜 날짜 기반 Jobs인가? (번호 기반 SPEC 대신)

**결정**: `.do/jobs/{YYMMDD}/{title}/` 날짜 기반
**MoAI와의 차이**: MoAI는 `.moai/specs/SPEC-XXX/` 번호 기반
**근거**: 날짜 기반은 시간순 탐색이 자연스럽고, 이름만 보고 언제 한 작업인지 알 수 있다.

### 8.7 왜 페르소나+스타일을 독립 축으로?

**결정**: 페르소나(캐릭터)와 스타일(응답 형식)을 분리
**MoAI와의 차이**: MoAI에는 페르소나가 없고 스타일만 있음
**근거**: 독립 축이면 4+3=7개만 정의하면 되고 어떤 조합이든 가능하다.

### 8.8 왜 `[o]`가 완료이고 `[x]`가 실패인가?

**결정**: 커스텀 상태 기호 체계
**근거**: 표준 마크다운 `[x]`는 "체크됨"과 "실패"를 구분할 수 없다. Do의 체크리스트는 6단계 상태 머신이므로 각 상태가 고유한 기호를 가져야 한다.

---

**작성자**: synthesizer agent
**작성일**: 2026-02-15
**소스**: research-do-philosophy.md, research-moai-philosophy.md, analysis-part1.md, analysis-part2.md, architecture.md
