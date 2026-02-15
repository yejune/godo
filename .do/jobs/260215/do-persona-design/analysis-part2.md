# Do Persona Design: Part 2 분석

**분석 일시**: 2026-02-15
**분석 대상**: moai persona override skills, output styles, hooks, settings, commands
**비교 대상**: do-focus 동등 파일들

---

## 1. Persona Override Skills 분석

moai persona는 6개의 override skill을 가지고 있다. core 스킬 디렉토리에는 이 6개와 이름이 일치하는 스킬이 없다 -- core에는 `foundation-claude`, `foundation-context`, `foundation-philosopher` 등 다른 foundation 스킬과 `workflow-testing`, `workflow-loop`, `workflow-thinking` 등 다른 workflow 스킬이 있다. 즉 이 6개는 **moai persona 전용 스킬**이다 (core에서 fork한 것이 아니라 persona 레이어에서 새로 추가된 것).

### 1.1 moai-foundation-core (v2.5.0)

**핵심 내용**: MoAI-ADK의 6가지 기반 원칙 정의
1. TRUST 5 Framework -- Tested, Readable, Unified, Secured, Trackable 품질 게이트
2. SPEC-First DDD -- 3단계 워크플로우 (SPEC 30K → DDD 180K → Docs 40K)
3. Delegation Patterns -- Task()를 통한 전문 에이전트 위임 (직접 실행 금지)
4. Token Optimization -- 200K 토큰 예산 관리, /clear 전략
5. Progressive Disclosure -- 3단계 지식 전달 (Quick/Implementation/Advanced)
6. Modular System -- SKILL.md 500줄 제한, modules/ 참조 아키텍처

**delta (core와 비교)**: core에 동일 이름 스킬 없음. 이 스킬은 moai의 "헌법" 역할.

**do 변환 시 필요 사항**:
- TRUST 5 → do-focus의 `.claude/rules/dev-testing.md`와 `.claude/rules/dev-workflow.md`에 이미 품질 규칙 존재. TRUST 5라는 브랜딩은 불필요.
- SPEC-First DDD → do-focus는 Plan → Checklist → Develop 워크플로우 사용. SPEC/EARS 개념 불필요.
- Delegation Patterns → do-focus의 CLAUDE.md에 Do 모드 위임 규칙으로 이미 존재.
- Token Optimization → do-focus에는 `.claude/rules/file-reading.md`로 축소 존재. 200K 예산 계획은 없음.
- Progressive Disclosure → do-focus는 skill 시스템 자체를 사용하지 않으므로 불필요.
- Modular System → 동일, 불필요.
- **결론**: 이 스킬의 내용은 do-focus의 rules/와 CLAUDE.md에 이미 분산 흡수됨. 별도 스킬 파일로 변환 불필요.

### 1.2 moai-foundation-quality (v2.2.0)

**핵심 내용**: 엔터프라이즈급 코드 품질 오케스트레이터
- TRUST 5 Validation 자동화
- Proactive Analysis (자동 이슈 감지)
- Best Practices Engine (Context7 연동)
- Multi-Language Support (25+ 언어)
- CI/CD 파이프라인 통합

**delta**: core에 동일 이름 없음. Python 클래스 기반 설계 (QualityOrchestrator, TRUST5Validator 등).

**do 변환 시 필요 사항**:
- do-focus는 코드 품질을 `dev-testing.md` + `dev-workflow.md`의 [HARD] 규칙으로 관리.
- Python 클래스 기반 설계는 do-focus에 맞지 않음 (Go 기반, godo binary).
- Context7 연동은 do-focus에서도 사용 가능하지만 스킬 형태가 아님.
- **결론**: 변환 불필요. do의 rules 파일이 이미 커버.

### 1.3 moai-workflow-ddd (v1.0.0)

**핵심 내용**: ANALYZE-PRESERVE-IMPROVE 리팩토링 워크플로우
- development_mode 설정 기반 선택 (quality.yaml)
- DDD vs TDD 비교표
- AST-grep 기반 구조 분석
- Characterization test 패턴
- Behavior preservation 원칙

**delta**: core에는 `workflow-loop`, `workflow-thinking` 등 다른 workflow가 있지만 DDD 전용은 없음.

**do 변환 시 필요 사항**:
- do-focus CLAUDE.md에 TDD 워크플로우가 기술되어 있음 (RED-GREEN-REFACTOR).
- DDD(ANALYZE-PRESERVE-IMPROVE)는 `.claude/rules/do/workflow/workflow-modes.md`에 이미 문서화됨.
- do-focus는 `.moai/config/sections/quality.yaml` 대신 사용자에게 직접 TDD 여부를 물어봄.
- **결론**: workflow-modes.md 규칙이 이미 커버. 별도 스킬 불필요.

### 1.4 moai-workflow-tdd (v1.0.0)

**핵심 내용**: RED-GREEN-REFACTOR 워크플로우
- quality.yaml의 development_mode 기반 선택
- TDD vs DDD 비교
- Specification tests (행위 정의 테스트)
- 80% 커밋별 최소 커버리지

**delta**: core에 없음. moai-workflow-ddd와 대칭 구조.

**do 변환 시 필요 사항**:
- do-focus CLAUDE.md의 TDD 섹션과 `dev-workflow.md`에 TDD 규칙 있음.
- `dev-testing.md`에 상세 테스트 규칙 있음.
- **결론**: 이미 do rules에 흡수됨. 별도 스킬 불필요.

### 1.5 moai-workflow-spec (v1.2.0)

**핵심 내용**: SPEC 문서 관리 + EARS 형식 요구사항
- EARS 5패턴: Ubiquitous, Event-Driven, State-Driven, Unwanted, Optional
- 4단계 요구사항 명확화 프로세스
- .moai/specs/SPEC-XXX/ 3파일 구조 (spec.md, plan.md, acceptance.md)
- SPEC Lifecycle (spec-first → spec-anchored → spec-as-source)
- Git Worktree 기반 병렬 개발
- Plan-Run-Sync 워크플로우 통합

**delta**: core에 없음. moai의 핵심 차별점.

**do 변환 시 필요 사항**:
- do-focus는 SPEC/EARS 개념을 사용하지 않음.
- 대신 `.do/jobs/{YYMMDD}/{title}/plan.md` 구조 사용.
- Analysis → Architecture → Plan → Checklist 워크플로우가 SPEC을 대체.
- **결론**: do에서는 완전히 다른 접근. 변환 대상 아님.

### 1.6 moai-workflow-project (v2.0.0)

**핵심 내용**: 프로젝트 초기화 + 문서 생성 + 다국어 지원
- 자동 프로젝트 타입 감지 (web, mobile, CLI, library, ML)
- 다국어 문서 생성 (en, ko, ja, zh)
- 템플릿 최적화
- SPEC 기반 문서 업데이트

**delta**: core에 없음. moai init/project 커맨드의 백엔드.

**do 변환 시 필요 사항**:
- do-focus의 `/do:setup` 커맨드가 초기화를 담당.
- godo binary가 프로젝트 설정 관리.
- 다국어 문서 생성은 do-focus 범위 밖.
- **결론**: godo의 `init`/`setup` 커맨드로 대체됨. 별도 스킬 불필요.

### Override Skills 종합 요약

| moai override skill | do-focus 동등물 | 변환 필요? |
|---|---|---|
| moai-foundation-core | CLAUDE.md + rules/*.md | 불필요 (이미 흡수) |
| moai-foundation-quality | dev-testing.md + dev-workflow.md | 불필요 (이미 흡수) |
| moai-workflow-ddd | rules/do/workflow/workflow-modes.md | 불필요 (이미 흡수) |
| moai-workflow-tdd | dev-workflow.md TDD 섹션 | 불필요 (이미 흡수) |
| moai-workflow-spec | .do/jobs/ 구조 + dev-checklist.md | 불필요 (다른 접근) |
| moai-workflow-project | godo init/setup | 불필요 (godo 대체) |

**핵심 인사이트**: moai는 "스킬 파일"에 지식을 담고 progressive disclosure로 로드하는 방식. do-focus는 "rules 파일"에 [HARD] 규칙으로 직접 주입하는 방식. 아키텍처가 근본적으로 다르므로 스킬 → 스킬 변환이 아니라 스킬 → rules 변환이 이미 완료된 상태.

---

## 2. Output Styles 분석

### 2.1 폴더 구조 비교

| 항목 | moai | do-focus |
|---|---|---|
| 경로 | `personas/moai/output-styles/moai/` | `.claude/styles/` |
| 파일 수 | 3개 (r2d2.md, moai.md, yoda.md) | 3개 (r2d2.md, moai.md, yoda.md) |
| 네이밍 | 동일 | 동일 |

### 2.2 파일별 비교

#### R2-D2 (pair programming partner)

**moai 버전 vs do 버전**: 내용 100% 동일 (바이트 단위 일치).
- 583줄, v2.2.0
- Pair Programming Protocol (4 Phase)
- AskUserQuestion 필수 사용
- Skills + Context7 Integration
- Insight Protocol

**do 변환 시 차이점**: 없음. 이미 동일 파일.

#### MoAI (Strategic Orchestrator)

**moai 버전 vs do 버전**: 거의 동일하지만 1줄 차이.
- moai 버전: `@.claude/rules/moai/core/moai-constitution.md` 참조
- do 버전: `@.claude/rules/moai/core/do-constitution.md` 참조
- 나머지 268줄은 동일 (v4.0.0)

**do 변환 시 차이점**: 참조 경로 1곳만 `moai-constitution` → `do-constitution`으로 변경됨.

#### Yoda (Technical Wisdom Master)

**moai 버전 vs do 버전**: 내용 100% 동일 (바이트 단위 일치).
- 360줄, v2.1.0
- Deep Understanding 프레임워크
- .moai/learning/ 디렉토리 문서 생성
- Insight Exercise 패턴

**do 변환 시 차이점**: 없음.

### 2.3 Output Styles 종합

**결론**: do-focus의 output styles는 moai에서 그대로 복사된 것. moai.md의 참조 경로 1줄만 다름. 컨버터가 output-styles를 처리할 때:
1. 파일을 `personas/{name}/output-styles/{name}/` → `.claude/styles/`로 복사
2. 내용 중 `moai-constitution` → persona 이름에 맞게 치환
3. `.moai/` 경로 참조를 `.do/` 또는 해당 persona 경로로 치환

---

## 3. Hooks 분석

### 3.1 moai의 7개 shell scripts

모든 hook은 동일한 패턴의 "shell wrapper"이다:

```
1. mktemp으로 임시 파일 생성
2. stdin을 임시 파일에 저장
3. moai binary를 찾아서 실행 (3단계 폴백)
   - command -v moai (PATH에서 검색)
   - /Users/goos/go/bin/moai (하드코딩된 Go bin 경로)
   - 동일 경로 재시도
4. 못 찾으면 exit 0 (graceful fail)
```

| Hook Script | moai 명령 | Claude Code Event |
|---|---|---|
| handle-session-start.sh | `moai hook session-start` | SessionStart |
| handle-session-end.sh | `moai hook session-end` | SessionEnd |
| handle-pre-tool.sh | `moai hook pre-tool` | PreToolUse (Write\|Edit\|Bash) |
| handle-post-tool.sh | `moai hook post-tool` | PostToolUse (Write\|Edit) |
| handle-compact.sh | `moai hook compact` | PreCompact |
| handle-stop.sh | `moai hook stop` | Stop |
| handle-agent-hook.sh | `moai hook agent <action>` | 에이전트별 커스텀 |

**handle-agent-hook.sh**의 특이점:
- 다른 6개와 달리 `$1` 인자를 받음 (action 이름)
- Manager/Expert 에이전트별 action 목록 정의 (ddd-pre-transformation, backend-validation 등)
- `moai hook agent "$action"` 형태로 전달

### 3.2 do-focus의 godo 대체 방식

do-focus는 shell wrapper 없이 **godo binary를 직접 호출**한다.

| do-focus Hook | godo 명령 | Claude Code Event |
|---|---|---|
| (직접) | `godo hook session-start` | SessionStart |
| (직접) | `godo hook session-end` | SessionEnd |
| (직접) | `godo hook pre-tool` | PreToolUse (Write\|Edit\|Bash) |
| (직접) | `godo hook post-tool-use` | PostToolUse (.*) |
| (직접) | `godo hook compact` | PostToolUse (두 번째 항목) |
| (직접) | `godo hook stop` | Stop |
| (직접) | `godo hook subagent-stop` | SubagentStop |
| (직접) | `godo hook user-prompt-submit` | UserPromptSubmit |

### 3.3 차이점 분석

| 항목 | moai | do-focus |
|---|---|---|
| 실행 방식 | shell wrapper → moai binary | godo binary 직접 호출 |
| 래퍼 파일 | 7개 .sh 파일 필요 | 0개 (.sh 파일 없음) |
| Binary 탐색 | 3단계 폴백 (PATH → 하드코딩 → 하드코딩) | PATH만 사용 (godo가 PATH에 있어야 함) |
| PostToolUse matcher | `Write\|Edit` (쓰기만) | `.*` (모든 도구) |
| PreCompact | 전용 이벤트 | PostToolUse 두 번째 항목으로 compact 처리 |
| SubagentStop | 없음 | 있음 |
| UserPromptSubmit | 없음 | 있음 |
| agent-hook | 있음 (에이전트별 커스텀) | 없음 (subagent-stop으로 통합) |

**핵심 차이**:
1. moai는 **shell wrapper 패턴** (`.sh` 파일이 stdin을 받아 binary로 전달). do-focus는 **직접 호출 패턴** (settings.json에서 godo를 직접 실행).
2. do-focus가 2개 더 많은 이벤트 처리 (SubagentStop, UserPromptSubmit).
3. do-focus의 PostToolUse는 모든 도구에 반응 (`.*`), moai는 Write/Edit만.
4. moai의 agent-hook 시스템은 do-focus에 없음 -- do는 SubagentStop 하나로 통합.

**컨버터 설계 시 고려**:
- moai의 shell wrapper 7개를 생성할 필요 없음 (do 방식은 직접 호출)
- 대신 settings.json의 hooks 섹션에 godo 명령을 직접 매핑
- agent-hook은 do에서 subagent-stop으로 단순화

---

## 4. Settings 분석

### 4.1 구조 비교

| 항목 | moai settings.json | do settings.json |
|---|---|---|
| **outputStyle** | `"MoAI"` | `"pair"` |
| **plansDirectory** | `".moai/plans"` | `".do/jobs"` |
| **statusLine** | `.moai/status_line.sh` | `godo statusline` |
| **attribution.commit** | `"MoAI <email@mo.ai.kr>"` (이모지 포함) | `""` (비어있음) |
| **attribution.pr** | `"MoAI <email@mo.ai.kr>"` (이모지 포함) | `""` (비어있음) |

### 4.2 Hooks 비교 (settings.json 내부)

| Event | moai | do-focus |
|---|---|---|
| SessionStart | `.claude/hooks/moai/handle-session-start.sh` | `godo hook session-start` |
| SessionEnd | `.claude/hooks/moai/handle-session-end.sh` | `godo hook session-end` |
| PreToolUse | `.claude/hooks/moai/handle-pre-tool.sh` (Write\|Edit\|Bash) | `godo hook pre-tool` (Write\|Edit\|Bash) |
| PostToolUse | `.claude/hooks/moai/handle-post-tool.sh` (Write\|Edit) | `godo hook post-tool-use` (.*) + `godo hook compact` |
| PreCompact | `.claude/hooks/moai/handle-compact.sh` | (PostToolUse에 통합) |
| Stop | `.claude/hooks/moai/handle-stop.sh` | `godo hook stop` |
| SubagentStop | 없음 | `godo hook subagent-stop` |
| UserPromptSubmit | 없음 | `godo hook user-prompt-submit` |

### 4.3 Permissions 비교

**moai**: 매우 상세한 허용 목록 (194개 allow, 4개 ask, 65개 deny)
- 모든 언어 도구 허용 (Python, Node, Go, Rust, Ruby, Elixir, PHP 등)
- MCP 도구 허용 (context7, sequential-thinking)
- Team 도구 허용 (TaskCreate, TaskUpdate, TaskList, TaskGet)
- 상세한 deny 목록 (rm -rf, git force push, DB 삭제 등)

**do-focus**: 최소한의 허용 목록 (21개 allow, 0개 ask, 5개 deny)
- Git + Go + npm 기본 명령만
- MCP 도구 미포함
- Team 도구 미포함
- deny도 최소한 (rm -rf, git force push만)

**핵심 차이**:
- moai는 "모든 것을 명시적으로 허용"하는 화이트리스트 방식
- do-focus는 "최소한만 허용"하는 간결한 방식
- do-focus에 Team 도구(TaskCreate 등)와 MCP 도구가 빠져있는 것은 추후 추가 필요할 수 있음

### 4.4 Environment Variables

**moai**: 5개 env 설정
- `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`
- `CLAUDE_CODE_FILE_READ_MAX_OUTPUT_TOKENS=64000`
- `ENABLE_TOOL_SEARCH=1`
- `MOAI_CONFIG_SOURCE=sections`
- 매우 긴 PATH (Go, Node, Python, Bun, Cargo 등 모든 바이너리 경로)

**do-focus**: env 없음 (settings.json에). 대신 `settings.local.json`에서 `DO_*` 환경변수 관리.

### 4.5 기타

| 항목 | moai | do-focus |
|---|---|---|
| cleanupPeriodDays | 30 | 없음 |
| enableAllProjectMcpServers | true | 없음 |
| respectGitignore | true | 없음 |
| spinnerTipsEnabled | true | 없음 |

**컨버터 설계 시 고려**:
- settings.json은 persona별로 자동 생성해야 함
- outputStyle, plansDirectory, statusLine, attribution은 persona 설정에서 결정
- hooks는 godo 직접 호출 방식으로 통일
- permissions는 do-focus 방식(간결)으로 유지하되, 프로젝트 필요에 따라 확장 가능하게

---

## 5. Commands 분석

### 5.1 moai의 2개 커맨드

#### /moai:github (v1.0.0)
- **기능**: GitHub 이슈 수정 + PR 코드 리뷰 워크플로우
- **서브커맨드**: `issues` (이슈 분석/수정/PR생성), `pr` (멀티 관점 코드 리뷰)
- **Agent Teams 기본 사용**: 팀 모드가 기본, `--solo`로 서브에이전트 폴백
- **상세 기능**:
  - issues: 이슈 발견 → 분석 → 브랜치/수정 → PR 생성 → 리포트
  - pr: PR 발견 → 3명 병렬 리뷰(security/perf/quality) → 종합 → 제출
- **model**: sonnet
- **allowed-tools**: 22개 (Read, Write, Edit, ..., TeamCreate, SendMessage, TaskCreate 등)

#### /moai:99-release (v3.0.0)
- **기능**: MoAI-ADK 프로덕션 릴리즈 워크플로우
- **8단계 Phase**: Pre-flight → Quality Gates → Code Review → Version Selection → CHANGELOG → Final Approval → Tag & Push → GitHub Release Notes
- **특이사항**: GoReleaser 연동, 바이링구얼(영/한) CHANGELOG/릴리즈 노트
- **model**: sonnet
- **moai-adk 전용**: Go 프로젝트 빌드, 바이너리 배포

### 5.2 do-focus의 6개 커맨드

| 커맨드 | 기능 |
|---|---|
| `/do:check` | 체크리스트 상태 확인 |
| `/do:checklist` | 체크리스트 생성/관리 |
| `/do:mode` | Do/Focus/Team 모드 전환 |
| `/do:plan` | 플랜 생성 |
| `/do:setup` | 프로젝트 초기화 설정 |
| `/do:style` | 출력 스타일 전환 |

### 5.3 매핑 비교

| moai 커맨드 | do-focus 동등물 | 비고 |
|---|---|---|
| /moai:github | 없음 | do에 GitHub 워크플로우 없음 |
| /moai:99-release | 없음 | moai-adk 전용, do에 해당 없음 |
| /moai:1-plan | /do:plan | SPEC vs .do/jobs 차이 |
| /moai:2-run | (없음, CLAUDE.md 워크플로우로 대체) | do는 체크리스트 기반 실행 |
| /moai:3-sync | (없음) | do는 문서 동기화 미지원 |
| /moai:9-feedback | (없음) | do는 피드백 루프 미지원 |
| (없음) | /do:check | moai에 체크리스트 확인 커맨드 없음 |
| (없음) | /do:checklist | moai에 없음 (SPEC 기반이라 불필요) |
| (없음) | /do:mode | moai에 없음 (단일 모드) |
| (없음) | /do:setup | /moai:init (별도 커맨드, persona에 없음) |
| (없음) | /do:style | moai는 outputStyle 설정으로 관리 |

**핵심 차이**:
1. moai는 "SPEC 워크플로우 중심" (plan → run → sync → feedback)
2. do-focus는 "체크리스트 워크플로우 중심" (plan → checklist → check → mode/style)
3. moai의 github/release는 프로젝트 특화 커맨드 -- 범용 컨버터로는 변환 불가
4. do-focus의 mode/style 커맨드는 moai에 없는 고유 기능 (삼원 실행 구조)

**컨버터 설계 시 고려**:
- commands는 persona별로 자유롭게 정의 가능해야 함
- moai → do 변환 시 github/release는 제외 (프로젝트 특화)
- SPEC 관련 커맨드(1-plan, 2-run, 3-sync)는 do의 plan/checklist로 매핑하되, 내용은 완전히 다름

---

## 6. 종합 결론

### moai persona의 구성 요소 vs do-focus 아키텍처

| moai persona 구성 | 역할 | do-focus 동등물 | 변환 방식 |
|---|---|---|---|
| 6 override skills | 지식/워크플로우 정의 | rules/*.md + CLAUDE.md | 이미 흡수 완료 |
| 3 output styles | 응답 스타일 | .claude/styles/*.md | 거의 동일 (1줄 차이) |
| 7 hook scripts | 이벤트 핸들링 | godo binary 직접 호출 | shell wrapper 불필요 |
| 1 settings.json | 설정 | settings.json | 구조 동일, 내용 상이 |
| 2 commands | 커맨드 | 6 commands | 목적이 다름 |

### 컨버터가 persona에서 처리해야 할 것

1. **output-styles**: 그대로 복사 + 경로 치환 (`.moai/` → persona 경로)
2. **hooks**: shell wrapper 생성 불필요. settings.json에 godo 직접 호출 패턴 생성
3. **settings.json**: persona 이름 기반으로 자동 생성 (outputStyle, plansDirectory, attribution, hooks)
4. **commands**: persona별 커스텀. 공통 변환 규칙 없음 -- 프로젝트마다 다름
5. **override skills**: do-focus 아키텍처에서는 rules 파일로 대체. 스킬 파일 생성 불필요

### 아키텍처 차이 요약

```
moai: CLAUDE.md ← skills(progressive disclosure) ← rules ← agents ← hooks(shell wrapper → binary)
do:   CLAUDE.md ← rules(직접 주입) ← hooks(binary 직접 호출)
```

moai는 "스킬 시스템 + progressive disclosure"라는 중간 레이어가 있고, do-focus는 그 레이어를 제거하고 rules에 직접 주입한다. 이것이 가장 근본적인 아키텍처 차이다.

---

**작성자**: analysis agent
**검토 상태**: 분석 완료
**다음 단계**: Architecture 설계에서 컨버터의 persona 처리 로직 확정
