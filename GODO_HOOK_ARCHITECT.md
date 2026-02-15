# godo Hook 아키텍처 + DB 연동 설계

**Version**: 2.0.0
**Date**: 2026-02-16
**Purpose**: godo hook 시스템의 현재 아키텍처를 서술하고, 체크리스트 상태를 DB로 동기화하는 미래 설계를 정의한다

---

## 1. 개요

godo는 Go 언어로 작성된 단일 바이너리다. Claude Code의 hook 이벤트를 수신하여, Do 프레임워크가 요구하는 페르소나 주입, 실행 모드 관리, 보안 정책 시행, 체크리스트 자동화를 수행한다. Claude Code가 도구를 호출하거나, 세션을 시작하거나, 종료를 시도할 때마다 godo가 실행되어 system-reminder 또는 additionalContext를 통해 AI의 행동을 제어한다.

godo의 이름은 "Go + Do"의 합성어다. Go 언어로 구현된 Do의 실행 엔진이라는 뜻이다.

현재 godo의 역할은 세 가지 축으로 나뉜다.

첫째, **상태 주입**이다. hook 이벤트가 발생할 때마다 godo는 현재 페르소나(호칭, 말투), 실행 모드(Do/Focus/Team), 체크리스트 상태를 system-reminder 또는 additionalContext로 Claude Code에 주입한다. 이 주입 덕분에 AI는 세션 전체에 걸쳐 일관된 캐릭터를 유지하고, 올바른 실행 모드에서 동작하며, 미완료 체크리스트를 인지한다.

둘째, **보안 정책 시행**이다. PreToolUse hook에서 Write, Edit, Bash 도구 호출을 가로채어 위험 명령 차단, 파일 경로 검증, 시크릿 감지, AST 보안 스캔을 수행한다. 이 계층이 AI 에이전트가 실수로(또는 주입 공격에 의해) 위험한 작업을 수행하는 것을 원천 차단한다.

셋째, **체크리스트 파이프라인**이다. AI가 plan.md를 작성하면 PostToolUse hook이 이를 감지하여 올바른 jobs 디렉토리로 이동시키고, checklist.md stub을 자동 생성하며, 이후 매 도구 호출마다 미작성/미분해 체크리스트를 감지하여 AI에게 작성을 촉구한다.

이 세 가지 역할은 모두 **파일 기반**으로 동작한다. godo는 `.do/` 디렉토리의 파일을 읽고, system-reminder/additionalContext를 통해 AI에 정보를 주입하며, 디렉토리 구조와 stub 파일을 생성할 뿐이다. 체크리스트 내용 자체를 작성하지는 않는다.

미래에 godo는 네 번째 축을 갖게 된다: **DB 동기화**. 체크리스트 파일의 상태 변경을 감지하여 SQLite에 append-only로 기록함으로써, 인간이 대시보드에서 작업 진행 상황을 조회할 수 있게 한다. 이때 파일은 여전히 Source of Truth이고, DB는 파일의 읽기 최적화 뷰(view)로만 존재한다.

---

## 2. 현재 Hook 아키텍처

### 2.1 셸 래퍼 제거: 직접 바이너리 호출

Do의 hook 아키텍처를 이해하려면 MoAI와의 차이에서 시작하는 것이 가장 명확하다.

MoAI는 7개의 셸 스크립트 래퍼(`.claude/hooks/moai/*.sh`)를 거쳐 `moai` 바이너리를 호출한다. 셸 래퍼가 stdin JSON을 파이프로 전달하는 구조다.

```
MoAI 실행 경로:
  settings.json → .claude/hooks/moai/handle-post-tool.sh → moai hook post-tool
```

Do는 이 간접 계층 전체를 제거했다. settings.json에서 godo 바이너리를 직접 호출한다.

```
Do 실행 경로:
  settings.json → godo hook post-tool-use
```

이 결정의 배경에는 역사적 경험이 있다. 셸 래퍼 패턴은 28가지 별개의 이슈를 발생시켰다. PATH 해석 문제, stdin 인코딩 이슈, SIGALRM 타이밍 문제, 크로스 플랫폼 호환성 문제 등이다. 직접 바이너리 호출은 셸 계층에서 발생하는 이 모든 문제를 원천적으로 제거한다. 동시에 프로세스 하나가 줄어들어 성능도 개선된다.

### 2.2 Hook 이벤트 매트릭스

settings.json에 등록된 8개 hook 이벤트와 각각의 역할을 정리하면 다음과 같다.

| 이벤트 | Matcher | godo 서브커맨드 | 핵심 역할 |
|--------|---------|----------------|----------|
| **SessionStart** | (전체) | `godo hook session-start` | 페르소나 스피너 적용, 모드 상태 읽기, 프로젝트 감지, 버전 확인 |
| **UserPromptSubmit** | (전체) | `godo hook user-prompt-submit` | 페르소나 호칭/말투 리마인더, 현재 모드 리마인더 |
| **PreToolUse** | `Write\|Edit\|Bash` | `godo hook pre-tool` | 보안 정책 시행 (위험 명령 차단, 파일 경로 검증, 시크릿 감지, AST 스캔) |
| **PostToolUse** | `.*` | `godo hook post-tool-use` | 페르소나 유지, 체크리스트 파이프라인, AI 푸터 제거, 미커밋 감지, 린트 |
| **PostToolUse** | (전체) | `godo hook compact` | 컨텍스트 압축 시 스냅샷 보존 |
| **SubagentStop** | (전체) | `godo hook subagent-stop` | 에이전트 태스크 진행률 추적, 완료 보고 |
| **Stop** | (전체) | `godo hook stop` | 활성 체크리스트 감지 시 종료 차단 |
| **SessionEnd** | (전체) | `godo hook session-end` | git 상태 스냅샷, 세션 요약, Rank API 제출 |

### 2.3 `.*` Matcher의 설계 근거

PostToolUse hook의 matcher가 `.*`(모든 도구)인 것은 Do의 의도적인 철학적 결정이다. 이 결정은 MoAI와 가장 큰 차이를 만드는 지점 중 하나이므로 상세히 서술한다.

MoAI의 PostToolUse matcher는 `Write|Edit`이다. 파일이 작성되거나 편집될 때만 hook이 실행된다. 토큰 효율을 우선시하는 선택이다. 처리할 파일 변경이 있을 때만 hook이 실행되므로, 읽기 전용 작업(Read, Grep, Glob 등) 중에는 hook 오버헤드가 전혀 없다.

Do의 PostToolUse matcher는 `.*`이다. Read, Grep, Glob, Bash, WebSearch 등 모든 도구 호출 후에 hook이 실행된다. 이것은 토큰 효율 측면에서 비용을 지불하는 선택이다. 하지만 Do는 이 비용을 감수할 만한 이유가 있다고 판단했다.

그 이유는 **페르소나 일관성**이다. godo의 PostToolUse hook은 additionalContext로 페르소나 호칭과 말투 리마인더를 AI에 주입한다. hook이 Write/Edit에서만 실행된다고 가정해보자. AI가 Read -> Grep -> Glob -> Read -> Grep -> ... 같은 긴 읽기-검색 시퀀스를 수행하는 동안, 페르소나 리마인더가 한 번도 주입되지 않는다. 이 시퀀스가 길어지면 AI는 점차 페르소나 캐릭터(호칭, 말투, 태도)를 잊어버린다. 컨텍스트 윈도우에서 페르소나 관련 토큰이 멀어지기 때문이다.

Do에서 페르소나는 장식이 아니다. 구조적 기능이다. 페르소나가 사용자에게 "승민선배"라고 부르는 것은 캐릭터 유지의 핵심이며, 이 호칭이 깨지면 사용자 경험 전체가 깨진다. 따라서 Do는 모든 도구 호출 후에 페르소나 리마인더를 주입하는 것을 선택했다.

이 트레이드오프의 비용은 다음과 같다. 매 도구 호출마다 godo 프로세스가 실행되고, 반환된 additionalContext가 컨텍스트 토큰을 소비한다. 현재 구현에서는 페르소나 리마인더를 짧은 한 줄(호칭 지시 + 말투 지시)로 최소화하여 토큰 오버헤드를 억제하고 있다.

향후 최적화 방안으로는 세 가지가 검토되고 있다. 첫째, 동일한 리마인더가 연속으로 주입될 때 생략하는 dedup 로직. 둘째, 읽기 전용 도구에서는 페르소나만, 쓰기 도구에서는 페르소나 + 체크리스트 상태를 주입하는 선택적 실행. 셋째, 특정 횟수(예: 5회)마다 한 번만 주입하는 주기적 실행이다.

### 2.4 실행 흐름 다이어그램

각 hook 이벤트의 실행 흐름을 단계별로 서술한다.

#### SessionStart

SessionStart는 Claude Code 세션이 시작될 때 단 한 번 실행된다. 세션 전체의 초기 상태를 설정하는 역할이다.

```
Claude Code 세션 시작
  |
  v
godo hook session-start 실행
  |
  v
stdin에서 JSON 파싱 (cwd, session_id 등)
  |
  +---> .do/.current-mode 파일에서 현재 모드 읽기 (do/focus/team)
  |       (없으면 DO_MODE 환경변수, 그것도 없으면 기본값 "do")
  |
  +---> 페르소나 타입에 따른 스피너 동사 적용
  |       (settings.json의 spinner 갱신)
  |
  +---> GitHub API로 최신 godo 버전 확인 (3초 타임아웃)
  |       (.do/.latest-version에 캐시)
  |
  +---> 프로젝트 정보 감지
  |       (go.mod, package.json 등에서 언어/프로젝트명 추출)
  |
  v
systemMessage 생성:
  "current_mode: {mode}\nproject: {name}, lang: {lang}"
  |
  v
stdout JSON 출력: {"continue": true, "systemMessage": "..."}
```

핵심 제약이 있다. SessionStart는 모드 상태를 **읽기만** 한다. 절대 덮어쓰지 않는다. 이 단방향 읽기 원칙은 사용자가 `godo mode set` 명령으로 설정한 모드를 hook이 임의로 변경하지 않도록 보장한다.

#### UserPromptSubmit

UserPromptSubmit은 사용자가 프롬프트를 제출할 때마다 실행된다. 컨텍스트 압축(`/clear`)이나 /compact 후에도 살아남는 리마인더를 주입하는 유일한 경로이므로, 페르소나와 모드 정보의 "마지막 방어선" 역할을 한다.

```
사용자가 프롬프트 제출
  |
  v
godo hook user-prompt-submit 실행
  |
  v
stdin JSON 파싱
  |
  +---> 페르소나 호칭/말투 리마인더 생성
  |       (예: "반드시 '승민선배'로 호칭할 것. 말투: 반말+존댓말 혼합")
  |
  +---> 현재 모드 리마인더 생성
  |       (예: "현재 실행 모드: do (응답 접두사: [Do])")
  |
  v
additionalContext 출력
```

#### PreToolUse (Write|Edit|Bash)

PreToolUse는 AI가 파일 수정이나 명령 실행을 시도하기 전에 실행되는 보안 게이트다. "allow", "deny", "ask" 세 가지 판정을 내린다.

```
Claude Code가 Write/Edit/Bash 도구 호출 시도
  |
  v
godo hook pre-tool 실행
  |
  v
stdin JSON에서 tool_name, tool_input 파싱
  |
  +---> [Bash 명령 검사]
  |       denyBashPatterns 매칭 --> "deny" (위험 명령 차단)
  |       askBashPatterns 매칭  --> "ask"  (사용자 확인 요청)
  |
  +---> [Write/Edit 파일 접근 검사]
  |       파일 경로를 절대 경로로 해석
  |       프로젝트 디렉토리 외부 접근 --> path traversal 차단
  |       denyFilePatterns 매칭 --> 보호 파일 차단 (.env, credentials)
  |       askFilePatterns 매칭  --> 중요 설정 파일 확인 요청
  |       콘텐츠 내 시크릿 패턴 감지 (API 키, 인증서 등)
  |
  +---> [AST 보안 스캔] (Write만, 선택적)
  |       ast-grep(sg) 설치 여부 확인
  |       임시 파일에 콘텐츠 작성 --> sg scan --json 실행 (10초 타임아웃)
  |       error severity 발견 --> "deny"
  |
  +---> [체크리스트 계층 검증] (Write만)
  |       .do/jobs/**/checklists/*.md 작성 시도 시:
  |         plan.md 존재 확인 --> 없으면 VIOLATION 차단
  |         checklist.md 존재 확인 --> 없으면 VIOLATION 차단
  |
  v
모든 검사 통과 --> "allow"
```

#### PostToolUse (.* matcher)

PostToolUse는 Do hook 시스템의 핵심이다. 모든 도구 호출 후에 실행되며, 8가지 하위 처리를 순차적으로 수행한다.

```
Claude Code가 임의의 도구를 호출한 후
  |
  v
godo hook post-tool-use 실행
  |
  v
stdin JSON에서 session_id, tool_name, tool_input, tool_output 파싱
  |
  +---> [1. 플랜 파일 파이프라인]
  |       Write/Edit이고 plan.md인 경우:
  |         .do/plans/에 쓴 경우 --> .do/jobs/{YY}/{MM}/{DD}/{slug}/plan.md로 이동
  |         checklist.md stub 자동 생성
  |         checklists/ 서브디렉토리 생성
  |         리마인더: "체크리스트를 작성하라"
  |
  +---> [2. AI 푸터 제거]
  |       Bash이고 "git commit" 포함 시:
  |         최근 커밋에서 "Co-Authored-By: Claude Code" 라인 제거
  |         (DO_AI_FOOTER=true이면 유지)
  |
  +---> [3. 페르소나 주입]
  |       DO_PERSONA + DO_USER_NAME에서 Persona 빌드
  |       additionalContext에 호칭 + 말투 리마인더 추가
  |
  +---> [4. 미작성 체크리스트 감지]
  |       plan.md 있고 checklist.md가 stub인 job 발견 시
  |       "[HARD] 미작성 체크리스트" 리마인더
  |
  +---> [5. 미분해 체크리스트 감지]
  |       checklist.md 작성됐지만 checklists/*.md 없는 job 발견 시
  |       "[HARD] 미분해 체크리스트" 리마인더
  |
  +---> [6. 에이전트 미커밋 감지]
  |       tool_name이 "Task"인 경우:
  |         git status로 .do/ 외부 미커밋 변경 확인
  |         발견 시 "[HARD] 커밋하라" 리마인더
  |
  +---> [7. 플랜 수정 시 문서 동기화 촉구]
  |       Edit이고 plan.md 수정인 경우:
  |         "checklist도 업데이트하라" 리마인더
  |
  +---> [8. 자동 린트]
  |       Write/Edit이고 코드 파일(.go, .ts, .py 등)인 경우:
  |         린트 실행 결과를 additionalContext에 추가
  |
  v
hookSpecificOutput로 additionalContext 조립 --> stdout JSON 출력
```

#### Stop

Stop hook은 Claude Code가 세션 종료를 시도할 때 실행된다. "시작되었지만 미완료"인 체크리스트가 있으면 종료를 차단하여, 에이전트가 작업을 중도 포기하지 않도록 강제한다.

```
Claude Code가 종료 시도
  |
  v
godo hook stop 실행
  |
  v
stdin JSON에서 stop_hook_active 확인
  |
  +---> stop_hook_active == true?
  |       예 --> 즉시 {} 반환 (무한 루프 방지: 이전 블로킹 후 재시도)
  |
  +---> .do/jobs/ 탐색: 최근 날짜의 checklist.md 파싱
  |       6종 상태 기호 카운트: [ ] [~] [*] [!] [o] [x]
  |
  +---> 판정:
  |       전체 완료([o]) --> 종료 허용
  |       작업 진행 중([~] 또는 [!] 존재), 미완료 --> 종료 차단
  |         {"decision": "block", "reason": "활성 체크리스트 있음 (3/5 완료, 1 진행중...)"}
  |       모든 항목 미시작([ ]) --> 종료 허용 (미시작 작업은 차단 안 함)
  |
  v
활성 체크리스트 없음 --> {} 반환 (종료 허용)
```

핵심 설계 원칙이 있다. "시작된 작업"만 종료를 차단한다. 플랜만 세우고 아직 개발을 시작하지 않은 체크리스트(모든 항목이 `[ ]`)는 종료를 차단하지 않는다. 사용자가 계획을 세운 뒤 다음 세션에서 실행하는 패턴을 허용하기 위함이다.

#### SubagentStop

SubagentStop은 에이전트(Task tool)가 실행을 완료했을 때 발생한다. 백그라운드 태스크 추적과 전체 진행률 보고를 담당한다.

```
에이전트(Task tool) 실행 완료
  |
  v
godo hook subagent-stop 실행
  |
  v
stdin JSON에서 agent_id, status, description 파싱
  |
  +---> .do/cache/background-tasks.json에서 상태 로드
  +---> 완료 시각, 소요 시간 기록
  +---> 전체 진행률 재계산
  +---> stderr에 프로그레스 바 출력: [########--------] 4/8 (50%)
  |
  +---> 전체 완료 감지 (completed >= total):
  |       성공/실패 요약, 전체 소요 시간, 실패 태스크 목록 출력
  |       상태 리셋 (active = false)
  |
  v
{"continue": true} 출력
```

#### SessionEnd

SessionEnd는 세션이 종료된 후 실행된다. 2초 전체 타임아웃 내에서 병렬로 세션 데이터를 저장한다.

```
Claude Code 세션 종료
  |
  v
godo hook session-end 실행 (2초 전체 타임아웃)
  |
  +---> [goroutine 1] Rank API 세션 데이터 제출 (최선 노력, 논블로킹)
  |       인증 확인 --> 트랜스크립트 파싱 --> 토큰 사용량 집계 --> API 제출
  |
  +---> [goroutine 2] git 정보 수집
  |       현재 브랜치, 미커밋 파일 수, 최근 1시간 커밋 수
  |
  v
.do/memory/last-session-state.json에 스냅샷 저장
  |
  v
systemMessage 출력 (미커밋 WARNING 포함, 최근 커밋 수 포함)
```

### 2.5 페르소나 이중 주입 메커니즘

페르소나 리마인더는 두 개의 독립된 경로에서 AI에 주입된다. 이 이중 주입이 페르소나가 세션 전체에 걸쳐 깨지지 않는 메커니즘이다.

첫 번째 경로는 **UserPromptSubmit**이다. 사용자가 프롬프트를 제출할 때마다 호칭과 말투 리마인더가 주입된다. 이 경로의 핵심 특성은 `/clear`(컨텍스트 초기화) 이후에도 살아남는다는 것이다. `/clear`는 이전 대화 컨텍스트를 제거하지만, 다음 사용자 프롬프트 제출 시 UserPromptSubmit hook이 다시 실행되어 페르소나를 복원한다.

두 번째 경로는 **PostToolUse**이다. 모든 도구 호출 후에 페르소나 리마인더가 주입된다. 이 경로는 AI가 도구를 연속으로 사용하는 중간(사용자 프롬프트 없이 Read -> Grep -> Edit -> ... 시퀀스를 진행하는 동안)에도 페르소나를 유지하는 역할이다.

두 경로가 상호보완하는 이유는 다음과 같다. UserPromptSubmit만 있으면, AI가 도구를 10번 연속 호출하는 긴 시퀀스 동안 페르소나가 표류한다. PostToolUse만 있으면, `/clear` 직후 첫 번째 도구 호출 전까지 페르소나 정보가 없다. 두 경로가 함께 작동해야 "어떤 상황에서든" 페르소나가 유지된다.

4종 페르소나의 호칭과 말투 규칙은 `buildPersona(personaType, userName)` 함수에 하드코딩되어 있다.

| 페르소나 타입 | 호칭 형식 | 말투 | 캐릭터 |
|-------------|----------|------|--------|
| `young-f` (기본) | `{name}선배` | 반말+존댓말 혼합 (~할게요, ~했어요) | 밝고 에너지 넘치는 20대 여성 천재 개발자 |
| `young-m` | `{name}선배님` | 존댓말 (~하겠습니다, ~해보겠습니다) | 자신감 넘치는 20대 남성 천재 개발자 |
| `senior-f` | `{name}님` | 다정한 존댓말 (~해드릴게요, ~살펴볼까요) | 30년 경력의 레전드 50대 여성 개발자 |
| `senior-m` | `{name}씨` | 든든한 존댓말 (~해봅시다, ~확인해보죠) | 업계 전설의 50대 남성 시니어 아키텍트 |

### 2.6 모드 상태 관리

실행 모드(do/focus/team)는 `.do/.current-mode` 파일에 단일 문자열로 영속화된다. 이 파일은 `godo mode set <mode>` 명령으로만 변경되며, 세션 간에도 유지된다.

상태 읽기에는 우선순위가 있다. 첫째 `.do/.current-mode` 파일(godo mode set으로 설정), 둘째 `DO_MODE` 환경변수(settings.local.json에서 설정), 셋째 기본값 "do"이다.

모드 정보는 세 곳에서 AI에 주입된다. SessionStart에서는 systemMessage로, UserPromptSubmit에서는 additionalContext로, StatusLine에서는 프롬프트 접두사(`[Do]` / `[Focus]` / `[Team]`)로 주입된다.

설계 원칙은 단방향 흐름이다. Hook은 모드를 읽기만 한다. 변경은 오직 `godo mode set` CLI를 통해서만 가능하다. AI가 "모드를 전환하겠습니다"라고 선언하고 접두사만 바꾸는 것은 VIOLATION이다. 반드시 `godo mode set <mode>` 명령이 실행되어야 한다.

### 2.7 체크리스트 자동화 파이프라인

체크리스트 자동화는 PostToolUse hook의 핵심 기능이다. AI가 플랜 파일을 작성하는 순간부터 체크리스트 작성, 분해, 실행까지 전체 파이프라인을 자동으로 관리한다.

```
[플랜 작성]
사용자가 "플랜 짜줘" 요청
  --> AI가 Write 도구로 plan.md 작성
  --> PostToolUse hook 트리거 --> handlePlanFile() 감지

  Case A: .do/plans/에 작성 (Claude Plan Mode)
    --> 파일 제목에서 slug 추출
    --> .do/jobs/{YY}/{MM}/{DD}/{slug}/ 디렉토리 생성
    --> plan.md를 올바른 위치로 이동
    --> checklist.md stub 자동 생성 (6종 상태 범례 포함)
    --> checklists/ 서브디렉토리 생성
    --> 리마인더: "체크리스트를 작성하라"

  Case B: .do/jobs/**/ 에 직접 작성 (/do:plan 커맨드)
    --> 이미 올바른 위치
    --> checklist.md stub 자동 생성 (없으면)
    --> checklists/ 서브디렉토리 생성 (없으면)
    --> 리마인더: "체크리스트를 작성하라"

[지속 감시]
이후 매 도구 호출마다:
  --> checkUnfilledChecklists(): stub 상태 체크리스트 발견 시 작성 촉구
  --> checkUndecomposedChecklists(): 메인만 있고 서브 없는 상태 발견 시 분해 촉구

[플랜 수정 감지]
기존 plan.md가 Edit으로 수정될 때:
  --> checklist.md가 존재하면 "체크리스트도 업데이트하라" 동기화 리마인더
```

설계 원칙이 있다. Hook은 체크리스트를 **직접 작성하지 않는다**. Hook은 디렉토리 구조와 stub 파일만 생성하고, 실제 내용 작성은 AI 에이전트에게 맡긴다. Hook의 역할은 "이것을 해야 한다"는 리마인더를 주입하는 것이지, "이것을 대신 해주는 것"이 아니다.

### 2.8 StatusLine

StatusLine은 hook이 아니라 별도의 godo 서브커맨드(`godo statusline`)로 구현된다. Claude Code가 프롬프트 라인에 표시하는 상태 정보를 한 줄로 조립한다.

```
[Do] | opus | used:23% 12m | ~/Work/project | main +3 | $0.42 | 3.0.0
 |     |      |              |                |         |       +-- godo 버전
 |     |      |              |                |         +-- 세션 비용
 |     |      |              |                +-- git 브랜치 + 변경 파일 수
 |     |      |              +-- 프로젝트 경로 (~/ 축약)
 |     |      +-- 컨텍스트 사용률 (녹/노/빨 색상) + 세션 시간
 |     +-- 모델 약칭 (opus/sonnet/haiku)
 +-- 실행 모드 접두사
```

---

## 3. DB 연동 설계

### 3.1 설계 원칙

DB 연동의 설계 원칙은 네 가지다. 이 원칙들은 Do의 핵심 철학("파일이 Source of Truth", "append-only 기록")에서 직접 도출된다.

**파일 = Source of Truth, DB = 읽기 뷰.** 에이전트는 체크리스트 파일을 직접 읽고 쓴다. DB는 인간이 대시보드에서 조회하기 위한 읽기 최적화 뷰(materialized view)일 뿐이다. 파일이 git으로 추적되고, 에이전트 간 인수인계 메커니즘으로 기능하며, 세션 종료 후에도 영속하는 Source of Truth이다. DB가 이 역할을 대체하지 않는다.

**단방향 동기화: 파일에서 DB로만.** 데이터 흐름은 항상 파일 --> DB 한 방향이다. DB에서 파일을 수정하는 역방향은 존재하지 않는다. 역방향을 허용하는 순간 두 개의 Source of Truth가 생기고, 동기화 충돌이 불가피해진다. 이 원칙은 타협 불가능하다.

**Append-only: INSERT만, UPDATE/DELETE 없음.** `state_transitions` 테이블은 절대 UPDATE나 DELETE하지 않는다. 새로운 상태 전이만 INSERT한다. 이것은 Do의 핵심 철학 -- "기존 기록을 고치지 않고 새로운 커밋을 추가한다" -- 의 DB 표현이다. 커밋 로그가 append-only이듯, 상태 전이 로그도 append-only다.

**최소 침습: 기존 동작을 변경하지 않는다.** DB 로직은 기존 hook 흐름에 추가되는 것이지, 기존 동작을 변경하는 것이 아니다. DB가 다운되어도 페르소나 주입, 보안 정책 시행, 체크리스트 리마인더 등 기존 hook 기능은 모두 정상 작동해야 한다. DB 쓰기 실패는 경고를 출력하되, hook 실행을 차단하지 않는다(graceful degradation).

### 3.2 SQLite 스키마

```sql
-- ================================================================
-- 작업 단위 (job)
-- ================================================================
-- .do/jobs/{YY}/{MM}/{DD}/{slug}/ 디렉토리 하나가 jobs 테이블의 한 행에 대응한다.
CREATE TABLE jobs (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT NOT NULL,                    -- kebab-case 제목 (slug)
    path        TEXT NOT NULL UNIQUE,             -- .do/jobs/26/02/16/login-api/
    status      TEXT NOT NULL DEFAULT 'active',   -- active, completed, archived
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_path ON jobs(path);

-- ================================================================
-- 메인 체크리스트 항목
-- ================================================================
-- checklist.md의 각 "- [{상태}] #N ..." 줄이 한 행에 대응한다.
CREATE TABLE checklist_items (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    job_id          INTEGER NOT NULL REFERENCES jobs(id),
    item_order      INTEGER NOT NULL,             -- #1, #2, ...
    subject         TEXT NOT NULL,                 -- 항목 제목
    status          TEXT NOT NULL DEFAULT '[ ]',   -- [ ] [~] [*] [!] [o] [x]
    owner           TEXT,                          -- 에이전트명 (expert-backend 등)
    depends_on      TEXT,                          -- 의존 항목 ID (JSON array, 예: [1,3])
    sub_checklist_path TEXT,                       -- checklists/01_expert-backend.md
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_checklist_items_job ON checklist_items(job_id);
CREATE INDEX idx_checklist_items_status ON checklist_items(status);
CREATE INDEX idx_checklist_items_owner ON checklist_items(owner);

-- ================================================================
-- 서브 체크리스트 메타데이터
-- ================================================================
-- checklists/*.md 파일 하나가 한 행에 대응한다.
-- 서브 체크리스트의 구조화된 섹션(Problem Summary, Test Strategy 등)을
-- 파싱하여 컬럼에 저장한다.
CREATE TABLE sub_checklists (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    item_id         INTEGER NOT NULL REFERENCES checklist_items(id),
    agent           TEXT NOT NULL,                 -- 담당 에이전트
    file_path       TEXT NOT NULL,                 -- checklists/01_expert-backend.md
    problem_summary TEXT,                          -- Problem Summary 섹션 내용
    test_strategy   TEXT,                          -- unit/integration/E2E/pass
    triggered_by    TEXT,                          -- 피드백 루프 출처 (NULL이면 원본)
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sub_checklists_item ON sub_checklists(item_id);
CREATE INDEX idx_sub_checklists_agent ON sub_checklists(agent);

-- ================================================================
-- 상태 변경 이력 (APPEND-ONLY -- 핵심 테이블)
-- ================================================================
-- 이 테이블은 절대 UPDATE/DELETE하지 않는다.
-- 체크리스트 항목의 모든 상태 전이가 시간순으로 기록된다.
-- Do의 "append-only commit 로그" 철학의 DB 표현이다.
--
-- 체크리스트의 "현재 상태"(checklist_items.status)는 사실 파생 값이다.
-- 진짜 데이터는 이 테이블의 상태 전이 이력이다.
-- [ ] -> [~] -> [*] -> [o] 시퀀스를 보면 작업의 전체 생애주기를 알 수 있다.
-- [ ] -> [~] -> [!] -> [~] -> [*] -> [~] -> [*] -> [o]를 보면
-- 블로커와 테스트 실패 후 재작업이 있었음을 알 수 있다.
-- 이 이력은 현재 상태 값만으로는 복원 불가능하다.
CREATE TABLE state_transitions (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    target_type TEXT NOT NULL,                     -- 'checklist_item' 또는 'sub_checklist'
    target_id   INTEGER NOT NULL,                  -- checklist_items.id 또는 sub_checklists.id
    from_state  TEXT NOT NULL,                     -- [ ] [~] [*] [!] [o] [x]
    to_state    TEXT NOT NULL,                     -- [ ] [~] [*] [!] [o] [x]
    detail      TEXT,                              -- 변경 사유, 블로커 이유 등
    commit_hash TEXT,                              -- [o] 전환 시 필수 기록
    timestamp   DATETIME DEFAULT CURRENT_TIMESTAMP
    -- NO UPDATE, NO DELETE -- append-only
);

CREATE INDEX idx_state_transitions_target ON state_transitions(target_type, target_id);
CREATE INDEX idx_state_transitions_to ON state_transitions(to_state);
CREATE INDEX idx_state_transitions_time ON state_transitions(timestamp);

-- ================================================================
-- 파일 소유권 레지스트리
-- ================================================================
-- Team 모드에서 에이전트 간 파일 충돌을 원천 차단한다.
-- 한 파일은 한 job 내에서 한 에이전트만 소유할 수 있다.
-- UNIQUE(job_id, file_path) 제약이 DB 레벨에서 충돌을 차단한다.
-- INSERT가 실패하면 "이 파일은 이미 다른 에이전트가 소유"라는 의미다.
CREATE TABLE file_ownership (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    job_id      INTEGER NOT NULL REFERENCES jobs(id),
    file_path   TEXT NOT NULL,                     -- 소스 파일 경로
    owner_agent TEXT NOT NULL,                     -- 소유 에이전트명
    item_id     INTEGER REFERENCES checklist_items(id),
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(job_id, file_path)                      -- 한 파일 = 한 소유자
);

CREATE INDEX idx_file_ownership_job ON file_ownership(job_id);
CREATE INDEX idx_file_ownership_agent ON file_ownership(owner_agent);
```

### 3.3 Hook에서 DB로의 파이프라인

PostToolUse hook에서 `.do/jobs/` 하위 파일의 Write 또는 Edit을 감지하면, 기존 로직(페르소나, 체크리스트 리마인더 등)을 실행한 뒤 DB INSERT를 추가로 수행한다.

```
PostToolUse (Write|Edit on .do/jobs/**)
  |
  v
기존 로직 실행 (페르소나, 보안, 체크리스트 리마인더)
  |
  v
[NEW] 파일 경로 패턴 매칭:
  |
  +---> */plan.md (생성/수정)
  |       jobs 테이블 UPSERT (path 기준)
  |       plan 내용에서 메타데이터 추출 (제목, 날짜 등)
  |
  +---> */checklist.md (생성/수정)
  |       parse_main_checklist() 호출
  |       정규식으로 상태 기호 추출: \- \[([ ~*!ox])\]
  |       이전 DB 상태와 비교 (diff)
  |       변경된 항목마다:
  |         INSERT INTO state_transitions (from, to, detail)
  |         UPDATE checklist_items SET status = {to_state}
  |
  +---> */checklists/*.md (생성/수정)
  |       parse_sub_checklist() 호출
  |       섹션별 파싱:
  |         Problem Summary --> sub_checklists.problem_summary
  |         Test Strategy   --> sub_checklists.test_strategy
  |         Progress Log    --> state_transitions 배치 INSERT
  |         FINAL STEP      --> commit_hash 추출
  |         Critical Files  --> file_ownership 배치 INSERT
  |
  +---> 기타 파일 --> DB 무시
  |
  v
DB 쓰기 실패 시:
  stderr에 경고 출력 (비차단)
  기존 hook 동작은 정상 계속 (graceful degradation)
```

파싱 전략의 핵심은 정규식 기반 상태 추출이다. checklist.md에서 `- [{상태}]` 패턴을 가진 줄을 항목으로 인식하고, 대괄호 안의 문자(공백, ~, *, !, o, x)를 상태로 추출한다. `#번호`에서 항목 순서를, `(depends on: #1, #3)`에서 의존성을, `(checklists/01_xxx.md)`에서 서브 체크리스트 경로를 추출한다.

### 3.4 쿼리 예시

DB가 구축되면 다음과 같은 쿼리로 작업 상황을 조회할 수 있다.

```sql
-- 1. 현재 활성 작업의 체크리스트 진행 상황
-- 대시보드의 메인 뷰에 해당한다.
SELECT ci.item_order, ci.subject, ci.status, ci.owner
FROM checklist_items ci
JOIN jobs j ON ci.job_id = j.id
WHERE j.status = 'active'
ORDER BY ci.item_order;

-- 2. 특정 항목의 전체 상태 이력 (append-only 추적)
-- "이 항목이 어떤 경로를 거쳐 완료되었는가"를 시간순으로 추적한다.
SELECT from_state, to_state, detail, commit_hash, timestamp
FROM state_transitions
WHERE target_type = 'checklist_item' AND target_id = ?
ORDER BY timestamp;

-- 3. 에이전트별 작업 통계
-- 어떤 에이전트가 몇 개를 완료했고, 평균 소요 시간은 얼마인지 집계한다.
SELECT
    ci.owner AS agent,
    COUNT(*) AS total_items,
    SUM(CASE WHEN ci.status = '[o]' THEN 1 ELSE 0 END) AS completed,
    ROUND(
        AVG(
            CASE WHEN ci.status = '[o]' THEN
                julianday(
                    (SELECT MAX(st.timestamp) FROM state_transitions st
                     WHERE st.target_id = ci.id AND st.to_state = '[o]')
                ) - julianday(
                    (SELECT MIN(st.timestamp) FROM state_transitions st
                     WHERE st.target_id = ci.id AND st.to_state = '[~]')
                )
            END
        ) * 24, 1
    ) AS avg_hours
FROM checklist_items ci
WHERE ci.owner IS NOT NULL
GROUP BY ci.owner;

-- 4. 파일 소유권 충돌 사전 감지 (Team 모드)
-- 특정 파일을 누가 소유하고 있는지 확인한다.
-- INSERT 전에 이 쿼리로 충돌 여부를 확인할 수 있다.
SELECT fo.file_path, fo.owner_agent
FROM file_ownership fo
WHERE fo.job_id = ? AND fo.file_path = ?;
```

---

## 4. 구현 로드맵

### Phase 1: SQLite 로컬 (MVP)

첫 번째 단계는 godo에 체크리스트 상태 추적 DB를 최소한으로 추가하는 것이다.

DB 파일 위치는 `.do/do.db`이다. SQLite WAL 모드, 단일 연결로 운영한다. Docker가 필요 없다. SQLite이므로 godo 바이너리에 내장되어 로컬에서 즉시 시작할 수 있다.

구현 범위는 다음과 같다. godo에 `db` 서브커맨드 그룹을 추가한다(`godo db init`, `godo db status`, `godo db sync`). PostToolUse hook의 checklist 감지 로직에 DB INSERT를 추가한다. DB 쓰기 실패 시 기존 동작을 유지하는 graceful degradation을 보장한다.

구현 순서는 여섯 단계다. 첫째 스키마 정의와 `godo db init` 명령(테이블 생성). 둘째 `parse_main_checklist()` 파서 구현(정규식 기반). 셋째 PostToolUse에서 checklist.md Write/Edit 감지 시 DB INSERT. 넷째 `godo db status` 명령(터미널에 현재 진행 상황 출력). 다섯째 state_transitions INSERT 로직. 여섯째 file_ownership INSERT 로직이다.

예상 변경 파일은 네 개다. `cmd/godo/db.go`(신규: db 서브커맨드), `cmd/godo/db_schema.go`(신규: 스키마와 마이그레이션), `cmd/godo/db_parser.go`(신규: checklist.md 파서), `cmd/godo/hook_post_tool_use.go`(수정: DB INSERT 추가).

### Phase 2: 대시보드

두 번째 단계는 체크리스트 상태를 시각적으로 확인할 수 있는 인터페이스를 제공하는 것이다. 두 가지 옵션이 있다.

옵션 A는 터미널 대시보드(`godo db status --watch`)다. 체크리스트 항목별 상태를 컬러로 표시하고, 에이전트별 진행률 바를 보여주며, 최근 state_transitions 타임라인과 블로커 하이라이트를 출력한다.

옵션 B는 웹 대시보드(`godo db serve`)다. SQLite를 읽기 전용 HTTP 서버로 노출하고, 실시간 상태 갱신(polling 또는 SSE), 상태 전이 타임라인 시각화, 파일 소유권 맵을 제공한다.

### Phase 3: PostgreSQL 확장 (선택)

세 번째 단계는 팀 공유 및 장기 이력 보관을 위한 PostgreSQL 지원이다. 이 단계는 여러 개발자가 같은 프로젝트에서 Do를 사용하고, 프로젝트 간 분석이 필요한 경우에만 진행한다.

변경 범위는 다음과 같다. `db.Adapter` 인터페이스를 추출하고, `db.SQLite`와 `db.PostgreSQL` 두 구현체를 만든다. `godo db` 서브커맨드에 `--driver` 옵션을 추가하고, 환경변수(`DO_DB_URL`)로 연결 문자열을 설정한다.

---

## 5. 설계 제약

다음 여섯 가지 제약은 위반 시 시스템의 정합성이 무너지는 하드 제약이다.

| 제약 | 규칙 | 위반 시 |
|------|------|---------|
| DB는 절대 Source of Truth가 아니다 | 파일이 원본, DB는 뷰 | DB를 기반으로 파일을 수정하면 데이터 불일치 발생 |
| DB에서 파일로의 역방향 수정 금지 | 파일 --> DB (단방향만) | 역방향을 추가하면 두 개의 Source of Truth가 생겨 충돌 불가피 |
| state_transitions는 DELETE/UPDATE 금지 | Append-only 강제 | 이력을 수정하면 감사 추적이 무너짐 |
| file_ownership의 UNIQUE 제약 유지 | DB 레벨에서 파일 소유권 충돌 차단 | 해제하면 동일 파일을 두 에이전트가 동시 수정 가능해짐 |
| DB 실패 시 hook 정상 동작 보장 | Graceful degradation | DB 의존성이 hook을 블로킹하면 AI 세션 전체가 영향받음 |
| godo 프로세스는 빠르게 종료 | Hook 실행 시간 최소화 | 느린 DB 쓰기가 매 도구 호출을 지연시키면 UX 저하 |

---

## 6. 참조 문서

| 문서 | 내용 | 위치 |
|------|------|------|
| DO_MOAI_COMPARISON.md | Do vs MoAI 비교 분석 -- hook 아키텍처, .* matcher 결정 근거, 체크리스트 시스템 철학 포함 | `./DO_MOAI_COMPARISON.md` |
| DO_PERSONA.md | Do의 독립 정체성과 개발 철학 -- 페르소나 시스템의 WHY와 4종 캐릭터 설계 | `./DO_PERSONA.md` |
| RUNBOOK.md | converter 운영 가이드 -- hook 파일 매핑, settings.json 구조 비교, core/persona 분해 구조 | `./RUNBOOK.md` |
| settings.json | 현재 hook 설정 (matcher, timeout 등) | `do-focus/.claude/settings.json` |
| dev-checklist.md | 체크리스트 시스템 규칙 -- 상태 기호, 전이 규칙, 서브 체크리스트 템플릿 | `do-focus/.claude/rules/dev-checklist.md` |
| dev-workflow.md | 개발 워크플로우 규칙 -- 에이전트 위임, 커밋 규율, 에러 대응 | `do-focus/.claude/rules/dev-workflow.md` |

---

**Document Version**: 2.0.0
**Date**: 2026-02-16
**Author**: godo hook system architecture
**Sources**: `cmd/godo/hook*.go`, `cmd/godo/mode.go`, `cmd/godo/statusline.go`, `DO_MOAI_COMPARISON.md`, `DO_PERSONA.md`, `RUNBOOK.md`
