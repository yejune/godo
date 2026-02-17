# Checklist System Rules [HARD]

## 생성 시점
- [HARD] 플랜 확정 후 반드시 체크리스트 생성 (에이전트 태스크당 1개)
- [HARD] 플랜 없이 체크리스트 생성 금지 — 플랜이 체크리스트의 근거
- PostToolUse hook이 플랜 파일 생성 시 체크리스트 stub을 자동 생성함
- [HARD] 플랜 파일이 존재하고 체크리스트가 stub 상태(미작성)이면, 다른 모든 작업보다 체크리스트 작성을 우선 수행
- [HARD] 체크리스트 작성 절차: 플랜 파일 읽기 → 에이전트별 작업 분해 → 에이전트(Task tool) 호출하여 체크리스트 작성
- [HARD] 항목은 에이전트가 토큰 소진 없이 완료할 수 있는 **아주 작은 단위**로 세분화
- [HARD] 하나의 항목 = 1~3개 파일 변경 + 검증 — 이 범위를 초과하면 쪼갤 것
- [HARD] 검증 방법은 항목별로 명시: 테스트 가능 → 단위/통합 테스트, 테스트 불가 → 빌드 확인/수동 확인/`docker compose config` 등

### 분해 절차 [HARD]
- [HARD] 1단계: 항목 하나가 몇 개 파일을 건드리는지 추정
- [HARD] 2단계: 3파일 초과 → 반드시 분해
- [HARD] 3단계: 분해된 각 항목은 독립적으로 완료/검증 가능해야 함
- [HARD] 4단계: 항목 간 의존성 있으면 `depends on:` 으로 연결
- [HARD] 분해 예시:
  - "API 구현" ← 너무 큼 (5+ 파일)
  - → "라우터 정의" (1파일), "핸들러 구현" (1파일), "검증 로직 추가" (1파일), "에러 응답 처리" (1파일), "단위 테스트" (1파일)
- [HARD] 체크리스트 미작성 상태에서 개발 작업 진행은 VIOLATION

## 체크리스트 항목 최소 요구사항 [HARD]

### 플랜보다 상세해야 하는 이유
- [HARD] 플랜 = 전략 (무엇을, 왜) / 체크리스트 = 실행 명세 (어떻게, 어디서, 뭘 검증)
- [HARD] 체크리스트 항목은 에이전트가 **플랜을 다시 읽지 않아도** 바로 코딩 가능한 수준이어야 함
- [HARD] 플랜의 요약이 아닌, 플랜의 **구체화**여야 함

### 각 항목 필수 포함 요소 [HARD]
- [HARD] **수정 파일**: 절대 경로 또는 프로젝트 루트 기준 상대 경로 (1~3개)
- [HARD] **수정 내용**: 어떤 함수/구조체/로직을 추가/변경하는지 구체적으로
  - 함수명, 파라미터, 반환값 수준까지 명시
  - "API 구현" ← 금지. "loginHandler(w, r) 함수 추가, POST /api/login, JWT 토큰 반환" ← 필수
- [HARD] **입력/출력**: 해당 코드가 받는 입력과 기대 출력 (해당되는 경우)
- [HARD] **검증 방법**: 구체적인 명령어 또는 테스트 파일 경로
  - "테스트 통과" ← 금지. "go test ./cmd/godo/ -run TestConvertAgent" ← 필수
  - 테스트 불가 시: "go build ./cmd/godo/ 성공" 또는 "grep -c 'moai-' result.md == 0"
- [HARD] **의존성**: 다른 항목에 의존하면 `depends on: #N` 명시

### 항목 상세도 검증 기준 [HARD]
- [HARD] "이 항목만 읽고 코딩할 수 있는가?" — Yes면 통과, No면 분해 부족
- [HARD] 파일 경로 없는 항목 = VIOLATION
- [HARD] 검증 방법 없는 항목 = VIOLATION
- [HARD] 함수명/로직 설명 없이 "구현" 만 적은 항목 = VIOLATION

### 예시

#### BAD (VIOLATION)
```
- [ ] 에이전트 컨버전 로직 구현
```

#### GOOD
```
- [ ] sync_agents.go: convertAgent() 함수 구현
      수정: cmd/godo/sync_agents.go
      내용: YAML skills 필드의 moai-* → do-* 접두사 치환
            본문의 "moai" → "do" 치환 (코드블록 내부 제외)
      입력: moai-adk/.claude/agents/moai/expert-backend.md
      출력: target/.claude/agents/do/expert-backend.md
      검증: go build ./cmd/godo/ && grep -c "moai-" output == 0
      depends on: #1
```

## 작성 방식 [HARD]
- [HARD] jobs 폴더 내 **모든 문서**(checklist.md, report.md, checklists/*.md)는 반드시 에이전트(Task tool)에게 위임
- [HARD] Do/Focus 모드 모두 동일 — 오케스트레이터가 직접 jobs 폴더 파일을 Write/Edit 하지 않음
- [HARD] 이유: 오케스트레이터 컨텍스트 토큰 낭비 방지 — 문서 작성은 에이전트 책임
- [HARD] plan.md만 예외: Plan Mode 훅이 자동 생성/이동 (오케스트레이터가 쓰는 게 아님)

## 체크리스트 = 에이전트 상태 파일 [HARD]
- [HARD] 체크리스트는 단순 문서가 아닌 **에이전트의 영속 상태 저장소**
- [HARD] 에이전트는 작업 시작 시 체크리스트를 읽고 → 작업 범위 파악
- [HARD] 항목 완료할 때마다 체크리스트 상태 갱신 → 진행 상황 파일에 기록
- [HARD] 에이전트 토큰 소진/중단 시 → 체크리스트에 마지막 상태가 남아있음
- [HARD] 새 에이전트가 동일 체크리스트를 받으면 → `[o]` 건너뛰고 미완료 항목부터 재개
- [HARD] 이 패턴으로 **작업 연속성 보장** — 어떤 에이전트든 이어받기 가능

## 파일 구조 (jobs 디렉토리 통합) [HARD]
- [HARD] 하나의 작업 = 하나의 폴더 — 모든 산출물이 같은 디렉토리에 위치:
  - 플랜: `.do/jobs/{YYMMDD}/{title-kebab-case}/plan.md`
  - 체크리스트: `.do/jobs/{YYMMDD}/{title-kebab-case}/checklist.md`
  - 완료 보고서: `.do/jobs/{YYMMDD}/{title-kebab-case}/report.md`
  - 서브 파일: `.do/jobs/{YYMMDD}/{title-kebab-case}/checklists/{order}_{agent-topic}.md`
- [HARD] 서브 파일의 `{order}`는 두 자리 숫자: `01`, `02`, ... `99`
- 디렉토리 없으면 자동 생성

### 예시
```
.do/jobs/260211/login-api-security/
  ├── plan.md                        ← 플랜
  ├── checklist.md                   ← 메인 체크리스트
  ├── report.md                      ← 완료 보고서
  └── checklists/                    ← 에이전트별 서브
      ├── 01_expert-backend.md
      ├── 02_expert-security.md
      └── 03_expert-testing.md
```

## Phase Gate 통과 조건 [HARD]

- [HARD] Phase 1 (PLAN_ONLY): checklist.md를 stub에서 실제 내용으로 채우면 Phase 2로 전이
- [HARD] Phase 2 (CHECKLIST_DRAFT): checklists/*.md 서브 체크리스트 생성하면 Phase 3으로 전이
- [HARD] Phase 3 이상에서만 구현 작업 가능
- [HARD] PreToolUse Hook이 Phase 1~2에서 Write/Edit/MultiEdit/Task 차단

## Plan Coverage 필수 [HARD]

- [HARD] checklist.md 상단에 `## Plan Coverage` 테이블 필수
- [HARD] plan.md의 모든 요구사항이 체크리스트 항목에 매핑됨을 보장
- [HARD] 매핑 안 된 요구사항이 있으면 PostToolUse Hook이 경고

## 상태 관리

### 상태 기호
| 기호 | 상태 | 의미 |
|------|------|------|
| `[ ]` | 미시작 (pending) | 아직 작업 시작 안 됨 |
| `[~]` | 진행중 (in progress) | 현재 작업 중 |
| `[*]` | 테스트중 (testing) | 구현 완료, 테스트 검증 중 |
| `[!]` | 블로커 (blocked) | 외부 의존성/결정 대기 중 |
| `[o]` | 완료 (done) | 테스트 통과, 커밋 해시 기록 |
| `[x]` | 실패 (fail) | 3회 재시도 후 해결 불가 |

### 상태 전이 규칙 [HARD]
- [HARD] 허용된 전이:
  ```
  [ ] → [~]        시작
  [~] → [*]        구현 완료 → 테스트
  [~] → [!]        블로커 발생
  [*] → [o]        테스트 통과 → 완료 + 커밋 해시 기록
  [*] → [~]        테스트 실패 → 재작업 (회귀)
  [*] → [x]        3회 회귀 후에도 실패 → fail
  [!] → [~]        블로커 해소 → 재개
  ```
- [HARD] 금지된 전이: `[ ] → [o]` (테스트 없이 완료 불가), `[ ] → [x]` (작업 없이 실패 불가), `[ ] → [*]` (작업 없이 테스트 불가)
- [HARD] 상태 변경 시 히스토리로 기록 (덮어쓰지 않음)

### 블로커 기록 규칙 [HARD]
- [HARD] `[!]` 전환 시 반드시 3가지 기록:
  1. **무엇이** 블로킹하는지 (구체적 이유)
  2. **누가** 해소할 수 있는지 (담당자/외부 시스템)
  3. **언제** 블로킹되었는지 (타임스탬프)

### 상태 히스토리 예시
```
[o] 로그인 API 구현
    - [ ] 2026-02-11 14:00 생성
    - [~] 2026-02-11 14:05 진행 시작
    - [!] 2026-02-11 15:00 블로커: Redis 설정 미완 (담당: infra팀, 해소 대기)
    - [~] 2026-02-11 16:00 블로커 해소, 재개
    - [*] 2026-02-11 17:00 테스트중
    - [~] 2026-02-11 17:30 테스트 실패 (JWT 만료 로직 오류) → 재작업
    - [*] 2026-02-11 18:00 재테스트
    - [o] 2026-02-11 18:30 완료 (commit: a1b2c3d)
```

## 의존성 관리 [HARD]
- [HARD] 항목 간 의존성은 `depends on:` 키워드로 선언
- [HARD] 의존 대상이 미완료면 해당 항목은 자동으로 `[!]` 블로커 취급
- [HARD] 의존성은 메인 체크리스트에서 관리 (서브 파일 간 참조)

### 의존성 표기법
```
## 작업 목록
- [o] #1 DB 스키마 마이그레이션
- [~] #2 로그인 API 구현 (depends on: #1)
- [ ] #3 프론트엔드 로그인 폼 (depends on: #2)
- [!] #4 소셜 로그인 연동 (depends on: #2, 블로커: OAuth 키 미발급)
```

## 서브 체크리스트 템플릿 [HARD]

각 서브 파일(`{order}_{agent-topic}.md`)은 다음 섹션을 포함:

```markdown
# {agent-topic}: {작업 제목}
상태: [ ] | 담당: {에이전트}

## Problem Summary
- 무엇을 해결하는가
- 왜 이 작업이 필요한가

## Acceptance Criteria
- [ ] 측정 가능한 완료 조건 1
- [ ] 측정 가능한 완료 조건 2
- [ ] 검증 완료 (아래 중 해당하는 방식):
  - 테스트 필요: `path/to/file_test.go` 작성 및 통과
  - 테스트 불필요: 검증 방법 명시 (빌드 확인, 수동 확인 등)
- [ ] 커밋 완료

## Solution Approach
- 선택한 접근 방식
- 왜 이 방식인가 (고려한 대안과 기각 이유)

## Critical Files
- **수정 대상**: `path/to/file.go` -- 변경 이유
- **참조 파일**: `path/to/ref.go` -- 참조 이유
- **테스트 파일**: `path/to/file_test.go`

## Risks
- 깨질 수 있는 것: (구체적으로)
- 주의할 점: (사이드이펙트, 성능, 호환성)

## Progress Log
- 2026-02-11 14:00 [~] 작업 시작: 초기 구조 설계
- 2026-02-11 15:30 [~] JWT 토큰 발급 로직 구현 완료
- 2026-02-11 16:00 [*] 단위 테스트 작성 및 실행
- 2026-02-11 16:30 [o] 모든 테스트 통과, 커밋 완료 (commit: a1b2c3d)

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 다음에 다르게 할 점:
```

### 템플릿 필수 규칙
- [HARD] Problem Summary, Acceptance Criteria, Critical Files는 작업 시작 전 반드시 작성
- [HARD] Acceptance Criteria에 검증 방법 필수 명시 — 테스트 파일 경로 또는 대안 검증 방법
- [HARD] 에이전트는 코드 작성 → 검증(테스트/빌드) → 통과 → 커밋까지가 한 세트 — 코드만 쓰고 끝내기 금지
- [HARD] 커밋 후 Progress Log에 커밋 해시 기록 필수 — 예: `[o] 완료 (commit: a1b2c3d)`
- [HARD] 커밋 해시 없는 `[o]` 완료 전환 금지 — 커밋이 곧 완료의 증거
- [HARD] Solution Approach는 구현 시작 시 작성 (대안 최소 1개 언급)
- [HARD] Progress Log는 상태 변경뿐 아니라 **무엇을 했는지** 기록 (작업 내용 추적)
- [HARD] Lessons Learned는 `[o]` 완료 전환 시 반드시 작성 -- 빈 칸 금지
- Risks는 식별된 것이 없으면 "식별된 리스크 없음" 기재

## 완료 보고 [HARD]

- [HARD] 모든 체크리스트 완료 시 해당 job 폴더의 `report.md`에 최종 완료 보고서 작성

### 완료 보고서 템플릿
```markdown
## 완료 보고서

### 실행 요약
- 완료: {N}/{M} 태스크 (예: 3/3)
- 기간: {시작일시} ~ {종료일시}

### 플랜 대비 변경사항
- (원래 계획과 달라진 점, 왜 변경했는지)
- (변경 없으면 "플랜대로 진행")

### 테스트 결과
- 전체: {pass}/{total} 통과
- 커버리지: {N}% (측정 가능한 경우)
- 실패/스킵: 없음 또는 상세 내역

### 변경 파일 요약
- `path/to/file.go` -- 변경 내용 한 줄 요약
- `path/to/test.go` -- 추가된 테스트

### 미해결 사항
- (후속 작업이 필요한 항목)
- (알려진 제약사항)
- (없으면 "없음")

### 핵심 교훈
- (서브 태스크 Lessons Learned 종합)
- (팀/프로젝트에 공유할 인사이트)
```

### 완료 보고 규칙
- [HARD] report.md 작성 시 checklist.md와 checklists/*.md를 참조하여 종합
  - 실행 요약 ← checklist.md 상태 집계
  - 핵심 교훈 ← 각 서브 체크리스트의 Lessons Learned 종합
  - 변경 파일 ← 각 서브 체크리스트의 Critical Files + `git diff --stat`
- [HARD] `미해결 사항`이 있으면 후속 플랜 또는 이슈로 등록
- [HARD] 테스트 결과에 실패가 있으면 완료 보고 금지 — 먼저 해결
- [HARD] 변경 파일 요약은 `git diff --stat`과 일치해야 함
