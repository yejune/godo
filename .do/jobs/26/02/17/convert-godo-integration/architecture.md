# Convert-Godo 통합 아키텍처 설계

**설계 일시**: 2026-02-17
**대상 프로젝트**: convert (~/Work/new/convert/)
**입력 문서**: handoff.md, research-godo.md, research-convert.md, research-decisions.md

---

## Overview

godo CLI(do-focus/cmd/godo/, 41파일 11,924줄)를 convert 프로젝트에 흡수 병합하고, do-focus의 dev-*.md 규칙 파일들을 토픽별로 분해하여 코어 스킬 소켓에 주입하며, DO 페르소나 워크플로우를 5개의 씬(thin) 파일로 재구성한다. 완료 후 convert가 단일 바이너리로 extract + assemble + 모든 godo 기능을 제공하고, do-focus는 레거시가 된다.

```
┌─────────────────────────────────────────────────────────────────┐
│                        convert (godo)                           │
│                                                                 │
│  cmd/godo/main.go ─── CLI Entry Point                          │
│       │                                                         │
│       ├── extract ──→ internal/extractor/                       │
│       ├── assemble ─→ internal/assembler/                       │
│       ├── hook ─────→ internal/hook/         ← NEW (from godo) │
│       ├── mode ─────→ internal/mode/         ← NEW (from godo) │
│       ├── lint ─────→ internal/lint/         ← NEW (from godo) │
│       ├── create ───→ internal/scaffold/     ← NEW (from godo) │
│       ├── claude ───→ internal/profile/      ← NEW (from godo) │
│       ├── spinner ──→ internal/persona/      ← NEW (from godo) │
│       ├── statusline→ internal/statusline/   ← NEW (from godo) │
│       ├── rank ─────→ internal/rank/         ← NEW (from godo) │
│       └── glm ──────→ internal/glm/         ← NEW (from godo) │
│                                                                 │
│  core/                        personas/do/                      │
│  ├── skills/                  ├── CLAUDE.md                     │
│  │   ├── foundation-core/     ├── workflows/ (5 thin files)    │
│  │   │   └── +에이전트 실행     │   ├── plan.md                  │
│  │   │      사이클 주입         │   ├── run.md                   │
│  │   ├── foundation-quality/  │   ├── report.md                │
│  │   │   └── +코딩 규율 주입    │   ├── team-plan.md             │
│  │   ├── foundation-context/  │   └── team-run.md              │
│  │   │   └── +파일 읽기 최적화  ├── commands/                    │
│  │   ├── workflow-spec/       ├── characters/                   │
│  │   │   └── +체크리스트 시스템  ├── spinners/                    │
│  │   ├── workflow-testing/    ├── output-styles/                │
│  │   │   └── +테스팅 규칙 주입  └── rules/bootapp.md             │
│  │   ├── workflow-project/                                      │
│  │   │   └── +Docker/12F 주입                                   │
│  │   └── workflow-tdd/                                          │
│  │       └── +TDD 규칙 주입                                      │
│  └── rules/                                                     │
│      ├── coding-standards.md  (as-is)                           │
│      ├── agent-authoring.md   (as-is)                           │
│      └── skill-authoring.md   (as-is)                           │
└─────────────────────────────────────────────────────────────────┘
```

---

## 1. Directory Structure

### 1.1 Go 코드 구조 (병합 후)

```
convert/
├── cmd/
│   └── godo/
│       └── main.go                    # 통합 CLI 진입점 (convert → godo 리네임)
├── internal/
│   ├── assembler/                     # [기존] core + persona → .claude/ 조립
│   │   ├── orchestrator.go            # 6단계 파이프라인
│   │   ├── merger.go                  # 파일 병합 + 슬롯 채우기
│   │   ├── slot_filler.go             # 슬롯 마커 치환
│   │   └── brand_deslotifier.go       # {{slot:BRAND}} → 실제 브랜드값
│   ├── cli/                           # [확장] CLI 명령 라우팅
│   │   ├── root.go                    # 기존 extract/assemble
│   │   ├── hook.go                    # NEW: godo hook <event> 라우팅
│   │   ├── mode.go                    # NEW: godo mode [get|set]
│   │   ├── lint.go                    # NEW: godo lint [--all|setup]
│   │   ├── create.go                  # NEW: godo create <agent|skill>
│   │   ├── claude.go                  # NEW: godo claude [--profile]
│   │   ├── spinner.go                 # NEW: godo spinner [apply|restore]
│   │   ├── statusline.go             # NEW: godo statusline
│   │   ├── rank.go                    # NEW: godo rank [login|status]
│   │   └── glm.go                     # NEW: godo glm [setup]
│   ├── detector/                      # [기존] 코어 vs 페르소나 분류
│   ├── extractor/                     # [기존] .claude/ → core + persona 추출
│   ├── hook/                          # NEW: 훅 디스패처 + I/O 계약
│   │   ├── dispatcher.go             # 이벤트 → 핸들러 라우팅
│   │   ├── types.go                   # HookInput, HookOutput, 결정 상수
│   │   ├── session_start.go           # SessionStart 핸들러
│   │   ├── pre_tool.go               # PreToolUse 핸들러 (보안 정책)
│   │   ├── post_tool_use.go          # PostToolUse 핸들러
│   │   ├── user_prompt.go            # UserPromptSubmit 핸들러
│   │   ├── subagent_stop.go          # SubagentStop 핸들러
│   │   ├── stop.go                    # Stop 핸들러
│   │   ├── compact.go                # Compact 핸들러
│   │   ├── session_end.go            # SessionEnd 핸들러
│   │   └── security.go               # 보안 패턴 (deny 규칙)
│   ├── lint/                          # NEW: 코드 린트 오케스트레이션
│   │   ├── runner.go                  # 린트 실행기
│   │   ├── gate.go                    # 게이트 로직
│   │   └── setup.go                   # 린트 설정
│   ├── mode/                          # NEW: 실행/권한 모드 관리
│   │   └── mode.go                    # do/focus/team + bypass/accept/default/plan
│   ├── model/                         # [기존, 확장] 데이터 구조
│   │   ├── document.go               # Document, Section
│   │   ├── persona_manifest.go       # PersonaManifest, AgentPatch
│   │   ├── classification.go         # ClassificationResult
│   │   ├── slot.go                    # Slot 정의
│   │   ├── depends_on.go             # DependsOn 검증
│   │   └── errors.go                  # 에러 타입 (+ hook/mode/lint 에러 추가)
│   ├── parser/                        # [기존] 마크다운 파싱
│   ├── persona/                       # NEW: 페르소나 로더
│   │   ├── loader.go                  # 캐릭터 YAML 로드, 리마인더 생성
│   │   └── spinner.go                 # 한국어 스피너 동사
│   ├── profile/                       # NEW: Claude 프로파일 관리
│   │   └── profile.go                # --profile 플래그 처리
│   ├── rank/                          # NEW: 랭크 시스템
│   │   ├── auth.go                    # 인증
│   │   ├── client.go                  # API 클라이언트
│   │   ├── config.go                  # 설정
│   │   └── transcript.go             # 트랜스크립트
│   ├── scaffold/                      # NEW: 에이전트/스킬 스캐폴딩
│   │   └── create.go                  # godo create <type> <name>
│   ├── statusline/                    # NEW: 상태 줄 렌더링
│   │   └── statusline.go
│   ├── glm/                           # NEW: GLM 백엔드
│   │   └── glm.go
│   ├── template/                      # [기존] 슬롯 시스템
│   └── validator/                     # [기존] 의존성 검증
├── core/                              # 코어 템플릿 (재구조화 대상)
├── personas/                          # 페르소나 (재구조화 대상)
└── testdata/                          # 테스트 데이터
```

### 1.2 코어 스킬 구조 (dev-*.md 분해 후)

```
core/skills/
├── foundation-core/
│   ├── SKILL.md                       # [기존 MoAI 원본 유지]
│   ├── modules/
│   │   ├── (기존 16개 modules 유지)
│   │   ├── agent-execution-cycle.md   # NEW ← dev-workflow.md 에이전트 실행 사이클
│   │   ├── agent-delegation.md        # NEW ← dev-workflow.md 위임+중단+재개 규칙
│   │   └── agent-research.md          # NEW ← dev-workflow.md 리서치 위임 규칙
│   └── references/
├── foundation-quality/
│   ├── SKILL.md                       # [기존 MoAI 원본 유지]
│   ├── modules/
│   │   ├── (기존 4개 modules 유지)
│   │   ├── read-before-write.md       # NEW ← dev-workflow.md 코딩 전 필수 행동
│   │   ├── coding-discipline.md       # NEW ← dev-workflow.md 코딩 규율 + 에러 대응
│   │   ├── commit-discipline.md       # NEW ← dev-workflow.md 커밋 규율
│   │   ├── parallel-agent-isolation.md # NEW ← dev-workflow.md 병렬 에이전트 격리
│   │   └── syntax-check.md            # NEW ← dev-environment.md 구문 검사
│   └── references/
├── foundation-context/
│   ├── SKILL.md                       # [기존 MoAI 원본 유지]
│   ├── modules/
│   │   ├── (기존 modules 유지)
│   │   ├── file-reading-optimization.md  # NEW ← file-reading-optimization.md 전체
│   │   └── knowledge-management.md       # NEW ← dev-workflow.md 지식 관리
│   └── references/
├── workflow-spec/
│   ├── SKILL.md                       # [기존 MoAI 원본 유지]
│   ├── modules/
│   │   ├── (기존 modules 유지)
│   │   ├── checklist-system.md        # NEW ← dev-checklist.md 시스템+상태관리+의존성
│   │   ├── checklist-templates.md     # NEW ← dev-checklist.md 서브 체크리스트 템플릿
│   │   ├── analysis-template.md       # NEW ← dev-checklist.md Analysis 문서 템플릿
│   │   ├── architecture-template.md   # NEW ← dev-checklist.md Architecture 문서 템플릿
│   │   ├── report-template.md         # NEW ← dev-checklist.md 완료 보고서 템플릿
│   │   └── complexity-check.md        # NEW ← dev-workflow.md 복잡도 판단+워크플로우 선택
│   └── references/
├── workflow-testing/
│   ├── SKILL.md                       # [기존 MoAI 원본 유지]
│   ├── modules/
│   │   ├── (기존 modules 유지)
│   │   ├── testing-rules.md           # NEW ← dev-testing.md 전체
│   │   └── bug-fix-workflow.md        # NEW ← dev-workflow.md 버그 수정 워크플로우
│   └── ...
├── workflow-project/
│   ├── SKILL.md                       # [기존 MoAI 원본 유지]
│   ├── references/
│   │   ├── (기존 references 유지)
│   │   ├── docker-rules.md            # NEW ← dev-environment.md Docker+12-Factor+env 관리
│   │   └── ai-forbidden-patterns.md   # NEW ← dev-environment.md AI 에이전트 금지 패턴
│   └── templates/
├── workflow-tdd/
│   ├── SKILL.md                       # [기존 MoAI 원본 유지]
│   ├── modules/
│   │   ├── (기존 modules 유지)
│   │   └── tdd-cycle.md               # NEW ← dev-workflow.md TDD RED-GREEN-REFACTOR
│   └── ...
└── (나머지 스킬들 변경 없이 유지)
```

### 1.3 DO 페르소나 구조 (재구조화 후)

```
personas/do/
├── manifest.yaml                      # 수정: workflows 추가, 오버라이드 스킬 제거
├── CLAUDE.md                          # 브랜드 아이덴티티 (유지)
├── settings.json                      # 페르소나별 settings (유지)
├── workflows/                         # NEW: 5개 씬 워크플로우
│   ├── plan.md                        # DO 스타일 plan (문서 체인)
│   ├── run.md                         # DO 스타일 run (체크리스트 기반)
│   ├── report.md                      # DO 스타일 report (집계)
│   ├── team-plan.md                   # DO 스타일 팀 plan
│   └── team-run.md                    # DO 스타일 팀 run
├── agents/                            # 페르소나 전용 에이전트 (유지)
├── characters/                        # 캐릭터 4종 (유지)
├── spinners/                          # 한국어 스피너 (유지)
├── commands/                          # DO 전용 슬래시 명령 (유지)
├── output-styles/                     # 출력 스타일 3종 (유지)
├── rules/
│   └── bootapp.md                     # DO 전용: bootapp 도메인 규칙
└── skills/
    └── do/                            # 유지: 오케스트레이터 스킬만 남김
        ├── SKILL.md
        └── references/
    # do-foundation-*, do-workflow-* 오버라이드 스킬들 → 전부 삭제
```

---

## 2. Core Interfaces

### 2.1 Hook 시스템 (internal/hook/)

godo의 훅 시스템을 `package main`에서 독립 패키지로 추출.

```go
package hook

import "encoding/json"

// EventType은 Claude Code 훅 이벤트 타입.
type EventType string

const (
    EventSessionStart EventType = "SessionStart"
    EventPreToolUse   EventType = "PreToolUse"
    EventPostToolUse  EventType = "PostToolUse"
    EventSessionEnd   EventType = "SessionEnd"
    EventStop         EventType = "Stop"
    EventSubagentStop EventType = "SubagentStop"
    EventPreCompact   EventType = "PreCompact"
    EventUserPrompt   EventType = "UserPromptSubmit"
)

// Decision 상수 (Claude Code 프로토콜).
const (
    DecisionAllow = "allow"
    DecisionDeny  = "deny"
    DecisionAsk   = "ask"
    DecisionBlock = "block"
)

// Input은 Claude Code가 stdin으로 보내는 JSON 페이로드.
type Input struct {
    SessionID      string          `json:"session_id,omitempty"`
    TranscriptPath string          `json:"transcript_path,omitempty"`
    CWD            string          `json:"cwd,omitempty"`
    PermissionMode string          `json:"permission_mode,omitempty"`
    HookEventName  string          `json:"hook_event_name,omitempty"`
    ToolName       string          `json:"tool_name,omitempty"`
    ToolInput      json.RawMessage `json:"tool_input,omitempty"`
    ToolOutput     json.RawMessage `json:"tool_output,omitempty"`
    Source         string          `json:"source,omitempty"`
    Model          string          `json:"model,omitempty"`
    AgentType      string          `json:"agent_type,omitempty"`
    Prompt         string          `json:"prompt,omitempty"`
    Reason         string          `json:"reason,omitempty"`
}

// Output은 Claude Code에 stdout으로 보내는 JSON 응답.
type Output struct {
    Continue           bool                `json:"continue,omitempty"`
    StopReason         string              `json:"stopReason,omitempty"`
    SystemMessage      string              `json:"systemMessage,omitempty"`
    SuppressOutput     bool                `json:"suppressOutput,omitempty"`
    Decision           string              `json:"decision,omitempty"`
    Reason             string              `json:"reason,omitempty"`
    HookSpecificOutput *HookSpecificOutput `json:"hookSpecificOutput,omitempty"`
}

// HookSpecificOutput은 PreToolUse/PostToolUse 전용 응답.
type HookSpecificOutput struct {
    HookEventName            string `json:"hookEventName,omitempty"`
    PermissionDecision       string `json:"permissionDecision,omitempty"`
    PermissionDecisionReason string `json:"permissionDecisionReason,omitempty"`
    AdditionalContext        string `json:"additionalContext,omitempty"`
}

// Handler는 개별 훅 이벤트 핸들러 인터페이스.
type Handler interface {
    Handle(input json.RawMessage) (*Output, error)
}

// Dispatcher는 이벤트를 핸들러로 라우팅.
type Dispatcher struct {
    handlers map[EventType]Handler
}

func NewDispatcher() *Dispatcher                                  { ... }
func (d *Dispatcher) Register(event EventType, handler Handler)   { ... }
func (d *Dispatcher) Dispatch(event string, input json.RawMessage) error { ... }
```

### 2.2 Mode 시스템 (internal/mode/)

```go
package mode

// ExecutionMode는 DO 실행 모드.
type ExecutionMode string

const (
    ModeDo    ExecutionMode = "do"
    ModeFocus ExecutionMode = "focus"
    ModeTeam  ExecutionMode = "team"
)

// PermissionMode는 Claude Code 권한 모드.
type PermissionMode string

const (
    PermBypass  PermissionMode = "bypassPermissions"
    PermAccept  PermissionMode = "acceptEdits"
    PermDefault PermissionMode = "default"
    PermPlan    PermissionMode = "plan"
)

// Manager는 실행/권한 모드 상태를 관리.
type Manager struct {
    stateFile    string // .do/.current-mode
    settingsFile string // .claude/settings.local.json
}

func NewManager() *Manager                                    { ... }
func (m *Manager) GetExecutionMode() ExecutionMode            { ... }
func (m *Manager) SetExecutionMode(mode ExecutionMode) error  { ... }
func (m *Manager) SetPermissionMode(mode PermissionMode) error { ... }
```

### 2.3 Persona Loader (internal/persona/)

```go
package persona

// Character는 캐릭터 YAML에서 로드된 페르소나 데이터.
type Character struct {
    Type         string `yaml:"type"`
    Honorific    string `yaml:"honorific"`
    Tone         string `yaml:"tone"`
    Relationship string `yaml:"relationship"`
}

// Loader는 페르소나 캐릭터를 로드.
type Loader struct {
    personaDir string
}

func NewLoader(personaDir string) *Loader                              { ... }
func (l *Loader) LoadCharacter(characterType string) (*Character, error) { ... }
func (l *Loader) BuildReminder(char *Character, userName string) string  { ... }
```

---

## 3. Error Handling

### 에러 계층 구조

```go
package model

// 기존 에러 타입 (유지)
type ErrExtraction struct { Phase, File, Message string }
type ErrAssembly struct { Phase, File, Message string }
type ErrValidation struct { Field, Message string }

// NEW: 훅 에러
type ErrHook struct {
    Event   string
    Message string
}

// NEW: 모드 에러
type ErrMode struct {
    Mode    string
    Message string
}

// NEW: 린트 에러
type ErrLint struct {
    Tool    string
    File    string
    Message string
}
```

### 에러 처리 전략

| 계층 | 에러 처리 | 예시 |
|------|----------|------|
| CLI (cmd/) | 에러 메시지 출력 + exit code | `godo hook unknown-event` → stderr + exit(1) |
| Hook I/O | JSON 파싱 실패 → 빈 input으로 진행 | stdin 비었을 때 graceful 처리 |
| Security (pre_tool) | deny 결정 반환 (프로세스 종료 안 함) | 금지 패턴 매칭 → deny + reason |
| Assembler | ErrAssembly 반환 + 파이프라인 중단 | 슬롯 미해결 → warning (에러 아님) |
| Extractor | ErrExtraction 반환 + 파이프라인 중단 | 파일 읽기 실패 → 즉시 에러 |

---

## 4. Component Implementations

### 4.1 dev-*.md 분해 상세 매핑

핵심 작업. 각 소스 파일의 내용이 어떤 코어 스킬 모듈로 가는지 정확한 매핑.

#### dev-workflow.md (184줄) → 7개 모듈

| 줄 범위 (대략) | 내용 | 대상 코어 스킬 | 대상 모듈 파일 |
|---------------|------|--------------|--------------|
| 1-34 | 복잡도 판단, 워크플로우 선택, Analysis/Architecture 단계 | workflow-spec | modules/complexity-check.md |
| 35-42 | TDD 선택 시 RED-GREEN-REFACTOR | workflow-tdd | modules/tdd-cycle.md |
| 43-68 | 에이전트 위임 필수 전달사항, 실행 사이클 | foundation-core | modules/agent-execution-cycle.md |
| 69-80 | 에이전트 중단 & 재개, 멱등 재개 | foundation-core | modules/agent-delegation.md |
| 81-87 | 에이전트 리서치 위임 규칙 | foundation-core | modules/agent-research.md |
| 88-120 | Read before Write, 코딩/커밋 규율, 병렬 격리 | foundation-quality | 4개 모듈로 분할 |
| 121-135 | 버그 수정 워크플로우 | workflow-testing | modules/bug-fix-workflow.md |
| 136-150 | 지식 관리 (참조 순서, 문서화 위치) | foundation-context | modules/knowledge-management.md |
| 151-160 | 에러 대응 (3회 재시도) | foundation-quality | modules/coding-discipline.md 합류 |

#### dev-testing.md (67줄) → 1개 모듈

| 내용 | 대상 코어 스킬 | 대상 모듈 파일 |
|------|--------------|--------------|
| 적용 범위, 테스트 철학(FIRST, Pyramid), AI 안티패턴, Real DB Only, 병렬성, 실행 순서, DB 트랜잭션 | workflow-testing | modules/testing-rules.md |

#### dev-checklist.md (523줄) → 6개 모듈

| 줄 범위 (대략) | 내용 | 대상 모듈 파일 |
|---------------|------|--------------|
| 1-100 | 생성 시점, 작성 방식, 상태 파일, 구조, 상태 관리, 블로커, 의존성 | checklist-system.md |
| 101-180 | 서브 체크리스트 템플릿 (전체 양식) | checklist-templates.md |
| 181-210 | 체크리스트 표시 의무, 완료 보고 | report-template.md |
| 211-340 | Analysis 문서 템플릿 | analysis-template.md |
| 341-523 | Architecture 문서 템플릿 | architecture-template.md |
| (all above) | 모두 workflow-spec 소속 | workflow-spec/modules/ |

#### dev-environment.md (92줄) → 3개 목적지

| 줄 범위 | 내용 | 대상 | 대상 파일 |
|---------|------|------|----------|
| 1-50 | Docker 필수, 명령 실행 구분, 빌드 재시작 | workflow-project | references/docker-rules.md |
| 51-65 | bootapp 네트워크 도메인 규칙 | 페르소나 전용 | personas/do/rules/bootapp.md |
| 66-80 | 환경변수 관리, 12-Factor, .env 금지 | workflow-project | references/docker-rules.md 합류 |
| 81-92 | AI 금지 패턴, 구문 검사 | workflow-project + foundation-quality | ai-forbidden-patterns.md + syntax-check.md |

#### file-reading-optimization.md → 1개 모듈

| 내용 | 대상 코어 스킬 | 대상 모듈 파일 |
|------|--------------|--------------|
| Progressive Loading, 토큰 예산 | foundation-context | modules/file-reading-optimization.md |

#### coding-standards.md, agent-authoring.md, skill-authoring.md → 그대로

이미 범용적. core/rules/로 복사 (변경 없음).

### 4.2 페르소나 워크플로우 5개 파일 내용 정의

각 파일은 씬(thin) — 코어 패턴을 참조하고 DO만의 차이점만 선언.

#### workflows/plan.md

```
역할: DO 스타일 Plan 페이즈 정의
내용:
- 트리거 매핑 (한국어 → 실행): "분석해"→Analysis, "설계해"→Architecture, "계획해"→Plan
- 복잡도 자동 판단 → 워크플로우 선택 (단순: Plan만, 복잡: 풀체인)
- 문서 체인: analysis.md → architecture.md → plan.md → checklist.md
- 산출물 위치: .do/jobs/{YY}/{MM}/{DD}/{title}/
- 참조: core workflow-spec (체크리스트 시스템, 복잡도 판단, 템플릿)
MoAI 대비 차이: SPEC 문서 대신 문서 체인, .moai/specs/ 대신 .do/jobs/
```

#### workflows/run.md

```
역할: DO 스타일 Run 페이즈 정의
내용:
- 에이전트 실행 사이클 참조 (READ-CLAIM-WORK-VERIFY-RECORD-COMMIT)
- 체크리스트 기반 진행: 서브 체크리스트 = 에이전트 상태 파일
- 멱등 재개: [o] 건너뛰기, [~] 이어받기
- VERIFY가 Run 안에 내장 (별도 sync 페이즈 없음)
- DDD/TDD/Hybrid 모드 선택 (core workflow-modes 참조)
MoAI 대비 차이: sync 페이즈 없음, 품질 검증이 Run 내장
```

#### workflows/report.md

```
역할: DO 스타일 완료 보고 정의
내용:
- 모든 체크리스트 [o] 시 자동 트리거
- 서브 체크리스트 Lessons Learned 종합
- report.md 템플릿 참조 (core workflow-spec report-template)
- 산출물 위치: .do/jobs/{YY}/{MM}/{DD}/{title}/report.md
MoAI 대비 차이: 독립 페이즈 아닌 Run의 마무리, docs sync 없음
```

#### workflows/team-plan.md

```
역할: DO 스타일 팀 Plan (Agent Teams API)
내용:
- 코어 plan-research 패턴 참조 (researcher + analyst + architect 병렬)
- 팀 이름: do-plan-{slug}
- 산출물: 문서 체인 (analysis → architecture → plan → checklist)
- 산출물 위치: .do/jobs/{YY}/{MM}/{DD}/{title}/
MoAI 대비 차이: SPEC 대신 문서 체인, .do/jobs/ 경로, 팀 이름 브랜드
```

#### workflows/team-run.md

```
역할: DO 스타일 팀 Run (Agent Teams API)
내용:
- 코어 implementation 패턴 참조 (backend-dev + frontend-dev + tester 병렬)
- 팀 이름: do-run-{slug}
- 체크리스트 기반 작업 분배 (파일 소유권 경계)
- VERIFY 내장 (team-quality가 Run 중 검증)
MoAI 대비 차이: 체크리스트 기반 실행, sync 없음
```

### 4.3 오버라이드 스킬 제거

| 오버라이드 스킬 | 처리 |
|---------------|------|
| do-foundation-core/ | 삭제 — 코어 foundation-core/modules/에 주입 |
| do-foundation-quality/ | 삭제 — 코어 foundation-quality/modules/에 주입 |
| do-workflow-spec/ | 삭제 — 코어 workflow-spec/modules/에 주입 |
| do-workflow-testing/ | 삭제 — 코어 workflow-testing/modules/에 주입 |
| do-workflow-project/ | 삭제 — 코어 workflow-project/references/에 주입 |
| do-workflow-tdd/ | 삭제 — 코어 workflow-tdd/modules/에 주입 |
| do-workflow-ddd/ | 삭제 — DDD 내용 이미 코어에 존재 |
| do-workflow-plan/ | 삭제 — personas/do/workflows/plan.md로 대체 |
| do-workflow-team/ | 삭제 — personas/do/workflows/team-*.md로 대체 |
| do/ (오케스트레이터) | 유지 — 페르소나 오케스트레이터 스킬 |

### 4.4 godo Go 코드 패키지 분리

현재 `package main` 41파일 → 독립 패키지:

| 현재 파일 | 대상 패키지 | 비고 |
|----------|-----------|------|
| hook.go, moai_hook_types.go, moai_hook_contract.go | internal/hook/ | 디스패처 + 타입 + 팩토리 |
| hook_session_start.go ~ hook_session_end.go (8파일) | internal/hook/ | 이벤트 핸들러 |
| security_patterns.go | internal/hook/ | 보안 패턴 |
| job_state.go | internal/hook/ | job 상태 추적 |
| mode.go | internal/mode/ | 실행/권한 모드 |
| lint.go, lint_*.go (4파일) | internal/lint/ | 린트 시스템 |
| create.go | internal/scaffold/ | 스캐폴딩 |
| claude_profile.go | internal/profile/ | 프로파일 |
| persona_loader.go, spinner.go | internal/persona/ | 페르소나 로더 |
| statusline.go | internal/statusline/ | 상태 줄 |
| rank*.go (5파일) | internal/rank/ | 랭크 시스템 |
| glm.go | internal/glm/ | GLM 백엔드 |
| agent.go | internal/cli/ | 에이전트 지원 |
| moai_sync*.go (5파일) | 삭제 | convert extract/assemble이 대체 |
| main.go | cmd/godo/ | 통합 CLI 진입점 |

---

## 5. Integration Layer

### 5.1 CLI 라우팅 통합

```
cmd/godo/main.go
  → cli.Execute()
    → cli.rootCmd
      ├── extractCmd     (기존)
      ├── assembleCmd    (기존)
      ├── hookCmd        → internal/hook/Dispatcher.Dispatch()
      ├── modeCmd        → internal/mode/Manager
      ├── lintCmd        → internal/lint/Runner
      ├── createCmd      → internal/scaffold/Create()
      ├── claudeCmd      → internal/profile/Launch()
      ├── spinnerCmd     → internal/persona/SpinnerApply()
      ├── statuslineCmd  → internal/statusline/Render()
      ├── rankCmd        → internal/rank/
      ├── glmCmd         → internal/glm/
      ├── selfupdateCmd  → (brew 연동)
      └── versionCmd     → (빌드 정보)
```

### 5.2 Hook ↔ Persona 연동

```
SessionStart 핸들러:
  1. internal/persona/Loader.LoadCharacter(DO_PERSONA 환경변수)
  2. Loader.BuildReminder(character, DO_USER_NAME)
  3. 리마인더 → Output.SystemMessage로 반환

PostToolUse 핸들러:
  1. 동일 리마인더 → Output.AdditionalContext로 반환

PreToolUse 핸들러:
  1. internal/hook/security.go 패턴 매칭
  2. internal/mode/Manager.GetExecutionMode() 현재 모드 확인
  3. 모드 + 보안 정책 조합으로 allow/deny/ask 결정
```

### 5.3 Assembler ↔ 새 구조 연동

assembler는 기존 코드 변경 없이 동작. manifest.yaml에 워크플로우 경로 추가:

```yaml
# manifest.yaml 추가
workflows:
  - workflows/plan.md
  - workflows/run.md
  - workflows/report.md
  - workflows/team-plan.md
  - workflows/team-run.md
```

assembler의 copyPersonaFiles()가 자동으로 workflows/ 파일들을 복사.
PersonaManifest 구조체에 Workflows []string 필드 추가 필요.

---

## 6. Configuration

### 6.1 Go 모듈

```
module: github.com/do-focus/convert (유지)
go: 1.25.0 (유지)
외부 의존성: gopkg.in/yaml.v3 (유지, 추가 없음)
```

godo 외부 의존성이 yaml.v3 하나뿐이라 충돌 없음.

### 6.2 빌드

```
바이너리명: godo
진입점: cmd/godo/main.go (cmd/convert/ → cmd/godo/ 리네임)
```

### 6.3 manifest.yaml 변경사항

```yaml
# 추가
workflows:
  - workflows/plan.md
  - workflows/run.md
  - workflows/report.md
  - workflows/team-plan.md
  - workflows/team-run.md

# rules 축소 (dev-*.md 코어 이동 후)
rules:
  - rules/bootapp.md

# skills에서 오버라이드 스킬 제거 (do/ 오케스트레이터만 유지)
```

---

## 7. Approach Comparison

### Approach A: 코어 스킬 모듈 주입 (권장)

dev-*.md 내용을 토픽별 분해하여 기존 코어 스킬의 modules/ 또는 references/에 새 .md 파일로 추가.

| 항목 | 평가 |
|------|------|
| 복잡도 | 중간 — 분해 매핑은 복잡하나 실행은 단순 (파일 생성) |
| 확장성 | 높음 — 모든 페르소나가 동일 코어 혜택 |
| 테스트 용이성 | 높음 — assembler 기존 테스트로 검증 가능 |
| 호환성 | 높음 — 기존 구조 유지, 모듈만 추가 |

장점:
- Single Source of Truth: 규칙이 코어에 한 곳만 존재
- 모든 페르소나(do, moai, do-ko, moai-ko)가 동일 규칙 혜택
- Progressive Disclosure 자연 적용 (Level 3 on-demand)

단점:
- 분해 정확도가 핵심 — 잘못 분류하면 코어 오염
- 기존 코어 모듈과 내용 중복 가능성 확인 필요

### Approach B: 페르소나 오버라이드 스킬 유지 (기각)

현재 구조 유지. 오버라이드 스킬에 dev-*.md 내용 보강.

| 항목 | 평가 |
|------|------|
| 복잡도 | 낮음 — 기존 구조에 내용 추가만 |
| 확장성 | 낮음 — 다른 페르소나에서 재사용 불가 |
| 테스트 용이성 | 중간 — 오버라이드 우선순위 검증 필요 |
| 호환성 | 높음 — 기존 변경 없음 |

기각 이유:
- dev-*.md 규칙은 범용적 (Docker, TDD, 체크리스트) — 페르소나 전용 아님
- 오버라이드 스킬이 코어와 중복 → 유지보수 비용 증가
- handoff.md 결정사항: "Override skills → ELIMINATED"

### 결론: Approach A 선택

코어에 규칙 주입 + 오버라이드 스킬 제거. handoff 결정사항과 일치, 장기 유지보수 유리.

---

## 8. Testing Strategy

### Unit Tests

| 대상 패키지 | 테스트 파일 | 내용 |
|-----------|-----------|------|
| internal/hook/ | dispatcher_test.go | 이벤트 라우팅, 알 수 없는 이벤트 처리 |
| internal/hook/ | types_test.go | Input/Output JSON 직렬화/역직렬화 |
| internal/hook/ | security_test.go | 보안 패턴 매칭, deny 규칙 |
| internal/mode/ | mode_test.go | 실행/권한 모드 get/set, 파일 영속성 |
| internal/persona/ | loader_test.go | 캐릭터 YAML 로드, 리마인더 생성 |
| internal/lint/ | runner_test.go | 린트 실행, 게이트 로직 |
| internal/scaffold/ | create_test.go | 에이전트/스킬 템플릿 생성 |
| internal/profile/ | profile_test.go | 프로파일 읽기/쓰기 |

### Integration Tests

| 대상 | 테스트 파일 | 내용 |
|------|-----------|------|
| assembler + 재구조화 코어 | assembler/e2e_test.go 확장 | 코어 모듈 주입 후 assemble 정상 동작 |
| CLI 통합 | cli/integration_test.go (NEW) | 전체 명령 라우팅 + 실행 검증 |
| 훅 E2E | hook/e2e_test.go (NEW) | stdin → 핸들러 → stdout 전체 파이프라인 |

### Test Matrix

| Layer | Method | Infrastructure |
|-------|--------|---------------|
| hook/types | Unit (JSON round-trip) | 없음 |
| hook/security | Unit (패턴 매칭) | 없음 |
| hook/dispatcher | Unit (라우팅) | 없음 |
| mode/manager | Unit (파일 I/O) | 임시 디렉토리 |
| persona/loader | Unit (YAML 파싱) | testdata/ |
| assembler pipeline | Integration (E2E) | testdata/ core + persona |
| CLI commands | Integration | 임시 디렉토리 |

---

## 9. Implementation Order

### Phase 1: 코어 스킬 콘텐츠 주입 (dev-*.md 분해)

분해 작업은 파일 생성만. Go 코드 변경 없이 진행 가능. 코어 완전성 확보를 위해 최우선 실행.

```
1-1.  foundation-core/modules/agent-execution-cycle.md 생성
1-2.  foundation-core/modules/agent-delegation.md 생성
1-3.  foundation-core/modules/agent-research.md 생성
1-4.  foundation-quality/modules/read-before-write.md 생성
1-5.  foundation-quality/modules/coding-discipline.md 생성
1-6.  foundation-quality/modules/commit-discipline.md 생성
1-7.  foundation-quality/modules/parallel-agent-isolation.md 생성
1-8.  foundation-quality/modules/syntax-check.md 생성
1-9.  foundation-context/modules/file-reading-optimization.md 생성
1-10. foundation-context/modules/knowledge-management.md 생성
1-11. workflow-spec/modules/checklist-system.md 생성
1-12. workflow-spec/modules/checklist-templates.md 생성
1-13. workflow-spec/modules/analysis-template.md 생성
1-14. workflow-spec/modules/architecture-template.md 생성
1-15. workflow-spec/modules/report-template.md 생성
1-16. workflow-spec/modules/complexity-check.md 생성
1-17. workflow-testing/modules/testing-rules.md 생성
1-18. workflow-testing/modules/bug-fix-workflow.md 생성
1-19. workflow-project/references/docker-rules.md 생성
1-20. workflow-project/references/ai-forbidden-patterns.md 생성
1-21. workflow-tdd/modules/tdd-cycle.md 생성
```

### Phase 2: 페르소나 구조 정리 + 양방향 동기화 (REVISED)

Original plan had 10 items (2-1 to 2-10). After research, revised to 6 work items:

```
2-1.  [CANCELLED] personas/do/workflows/ 디렉토리 생성 — 15파일이 skills/do/workflows/에 이미 존재, 이동 불필요
2-2.  [CANCELLED] thin workflow 파일 생성 — SKILL.md 오케스트레이터가 라우팅 처리
2-3.  [DONE] manifest.yaml 갱신 — 9개 누락 워크플로우 파일을 skills: 섹션에 추가
2-4.  [DONE] rules/bootapp.md 생성 — 네트워크 8규칙 + 빌드/재시작 7규칙
2-5.  [DONE] docker-rules.md 갱신 — Network 섹션 + Build & Restart 확장 + AI Anti-Patterns 12규칙 추가
2-6.  [DONE] output-styles/ 삭제 — styles/가 정본 (manifest 참조), output-styles/는 레거시 중복
2-7.  [ALREADY DONE] 오버라이드 스킬 삭제 — Phase 1에서 이미 완료 (injection modules만 잔존, 의도된 설계)
2-8.  [NOT APPLICABLE] rules/workflow/ 삭제 — spec-workflow.md/workflow-modes.md는 원본 dev-*.md가 아닌 페르소나 규칙
2-9.  [DONE] do-focus checklist.md 백포팅 — convert에서 추가된 8개 신규 규칙 + 3개 템플릿 개선
2-10. [DONE] architecture.md Phase 2 섹션 갱신 — 실제 결정 반영
```

### Phase 3: godo Go 코드 패키지 분리

```
3-1.  internal/hook/ 패키지 생성 (types.go, dispatcher.go, contract.go)
3-2.  internal/hook/ 핸들러 이동 (session_start.go ~ session_end.go, security.go)
3-3.  internal/mode/ 패키지 생성 (mode.go)
3-4.  internal/persona/ 패키지 생성 (loader.go, spinner.go)
3-5.  internal/lint/ 패키지 생성 (runner.go, gate.go, setup.go)
3-6.  internal/scaffold/ 패키지 생성 (create.go)
3-7.  internal/profile/ 패키지 생성 (profile.go)
3-8.  internal/statusline/ 패키지 생성 (statusline.go)
3-9.  internal/rank/ 패키지 생성 (auth.go, client.go, config.go, transcript.go)
3-10. internal/glm/ 패키지 생성 (glm.go)
```

### Phase 4: CLI 통합

```
4-1. cmd/convert/ → cmd/godo/ 리네임
4-2. internal/cli/ 확장 (hook, mode, lint 등 명령 추가)
4-3. moai_sync*.go 삭제 (extract/assemble이 대체)
4-4. go build 검증
```

### Phase 5: 테스트 + 검증

```
5-1. internal/hook/ 단위 테스트
5-2. internal/mode/ 단위 테스트
5-3. internal/persona/ 단위 테스트
5-4. assembler E2E 테스트 확장
5-5. CLI 통합 테스트
5-6. 전체 테스트 스위트 실행
5-7. assemble 출력 검증 (미해결 슬롯, 오버라이드 잔재 없는지)
```

---

## 10. Risk Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| dev-*.md 분해 시 누락된 규칙 | HIGH | 분해 전/후 줄 수 비교, 원본 대비 diff 검증 |
| 오버라이드 스킬 삭제 후 기능 누락 | HIGH | 삭제 전 오버라이드 내용이 코어에 100% 반영 확인 |
| godo 패키지 분리 시 내부 참조 깨짐 | MEDIUM | package main → package hook 전환 후 go build 즉시 검증 |
| 코어 모듈 추가로 Progressive Disclosure 깨짐 | LOW | modules/는 Level 3 on-demand이므로 토큰 영향 최소 |
| assembler가 workflows/ 인식 실패 | MEDIUM | PersonaManifest에 Workflows 필드 추가 or 기존 Skills 활용 |
| do-ko, moai-ko 동기화 누락 | MEDIUM | Phase 2 후 -ko 변형도 동일 구조 갱신 (별도 작업) |
| bootapp.md 분리 시 환경 규칙 불완전 | LOW | 원본 vs (bootapp + docker-rules) 합산 줄 수 비교 |
| moai_sync 삭제로 기존 사용자 영향 | MEDIUM | convert 동일 기능 확인 후 삭제, CHANGELOG 마이그레이션 안내 |

---

**작성자**: architecture agent
**검토 상태**: 설계 완료
**다음 단계**: Plan + Checklist 작성
