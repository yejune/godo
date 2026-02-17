# Do 개발 워크플로우

Do의 개발 워크플로우 — 유연한 토큰 관리와 커밋-증거 추적 방식.

## 워크플로우 개요

### 단순한 작업 (파일 4개 이하, 단일 도메인, 아키텍처 변경 없음)

```
Plan -> Checklist -> Develop -> Test -> Report
```

### 복잡한 작업 (파일 5개 이상, 신규 모듈, 마이그레이션, 3개 이상 도메인, 추상화 설계)

```
Analysis -> Architecture -> Plan -> Checklist -> Develop -> Test -> Report
```

| 단계 | 에이전트 | 목적 |
|-------|-------|---------|
| Analysis | expert-analyst | 현재 시스템 분석, 요구사항(EARS+MoSCoW), 기술 비교 |
| Architecture | expert-architect | 솔루션 설계, 인터페이스, 구현 순서 결정 |
| Plan | manager-plan | `.do/jobs/{YY}/{MM}/{DD}/{title}/plan.md`에 플랜 문서 생성 |
| Checklist | (오케스트레이터 위임) | 플랜을 에이전트 단위 항목으로 분해 (각 항목 1-3파일) |
| Develop | manager-ddd/tdd 또는 experts | 체크리스트 실시간 갱신과 함께 구현 |
| Test | (개발 단계 내) | Test Strategy별 검증 (unit/integration/E2E/pass) |
| Report | (오케스트레이터 위임) | `.do/jobs/{YY}/{MM}/{DD}/{title}/report.md`에 최종 보고서 |

## Plan 단계

EARS+MoSCoW 요구사항 형식으로 플랜 작성.

산출물:
- `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/plan.md`에 플랜 문서
- MoSCoW 우선순위가 부여된 EARS 형식 요구사항 (MUST/SHOULD/COULD/WON'T)
- 인수 기준
- 기술 접근 방식
- 복잡도 평가 (단순 vs 복잡)

plan-mode 산출물 작성 규칙:
- plan-mode 에이전트(researcher, analyst, architect)는 `permissionMode: plan`(읽기 전용)으로 실행
- plan 모드에서는 소스 코드 보호를 위해 Write/Edit 도구가 차단됨
- 산출물은 Bash heredoc을 통해 `.do/jobs/` 디렉토리에만 작성
- Bash를 통한 소스 코드 수정도 금지 — `.do/jobs/`에 새 파일 생성만 허용
- 파일이 진실의 원천(source of truth) — 메시지가 아닌 파일이 세션 간 멱등성을 보장

## Checklist 단계

플랜을 실행 가능한 체크리스트 항목으로 분해.

규칙:
- 항목 하나 = 파일 변경 1-3개 + 검증
- 3개 파일 초과 항목은 반드시 분해
- 각 항목에 Test Strategy 선언 (unit/integration/E2E/pass)
- 항목은 `checklists/` 디렉토리의 서브 체크리스트를 통해 에이전트에게 배정
- 체크리스트는 살아있는 문서 — 구현 중 계속 진화

산출물:
- 메인 체크리스트: `.do/jobs/{YY}/{MM}/{DD}/{title}/checklist.md`
- 서브 체크리스트: `.do/jobs/{YY}/{MM}/{DD}/{title}/checklists/{order}_{agent-topic}.md`

## Develop 단계

설정된 개발 방법론으로 체크리스트 항목 구현.

개발 방법론:
- DDD (ANALYZE-PRESERVE-IMPROVE): 레거시 리팩토링용
- TDD (RED-GREEN-REFACTOR): 신규 기능용
- Hybrid: 변경 유형별 혼합 (대부분의 프로젝트에 권장)
- 상세 방법론 사이클은 @workflow-modes.md 참조

체크리스트 실시간 갱신:
- 에이전트는 시작 전 서브 체크리스트를 읽음
- 작업 진행에 따라 상태 갱신: [ ] -> [~] -> [*] -> [o]
- 에이전트가 중단될 경우(토큰 소진), 다음 에이전트가 체크리스트를 읽고 마지막 상태부터 재개
- 체크리스트 파일이 인수인계 메커니즘 역할을 함

커밋-증거(Commit-as-proof):
- 모든 [o] 완료에는 커밋 해시 기록 필수
- `git add <특정 파일>`만 사용 (NEVER `git add -A` 또는 `git add .`)
- 커밋 전 `git diff --cached --name-only`로 검증
- 커밋 메시지는 WHY를 설명 (WHAT은 diff가 보여줌)
- Progress Log에 해시 기록: `[o] 완료 (commit: <hash>)`

성공 기준:
- 모든 체크리스트 요구사항 구현 완료
- 방법론별 테스트 통과
- 품질 차원 검증 (Tested/Readable/Unified/Secured/Trackable)

## Test 단계

검증은 Test Strategy 선언에 따라 Develop 단계 내에서 통합 실행.

테스트 가능한 코드(비즈니스 로직, API, 데이터 계층):
- dev-testing.md 규칙 적용: FIRST 원칙, Real DB only, AI 안티패턴 7가지
- 행위 기반 테스트, 구현 방식 테스트 금지
- 변이 테스트 사고방식: "이 줄을 바꾸면 테스트가 실패하는가?"

테스트 불가능한 변경 (CSS, 설정, 문서, 훅):
- 대안 검증: 빌드 확인, 수동 확인, 구문 검사
- Test Strategy에서 근거와 함께 `pass`로 선언

## Report 단계

모든 체크리스트 항목 완료 후 완료 보고서 생성.

산출물:
- `.do/jobs/{YY}/{MM}/{DD}/{title}/report.md`에 보고서
- 실행 요약, 플랜 대비 변경사항, 테스트 결과, 변경 파일, 핵심 교훈

## 컨텍스트 관리

/clear 전략 (유연하게, 고정 단계 아님):
- 체크리스트 항목 경계에서 (자연스러운 작업 단위)
- 컨텍스트가 임계값 초과 시 (고정 숫자 아님 — 작업에 따라 달라짐)
- 주요 워크플로우 전환 사이
- 엄격하게 사전 할당된 단계 경계에서는 실행하지 않음

Progressive Disclosure:
- Level 1: 메타데이터만 (~100 토큰)
- Level 2: 트리거 조건 충족 시 스킬 본문 (~5000 토큰)
- Level 3: 번들 파일 온디맨드

## Agent Teams 변형

팀 모드 활성화 시 (CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1), Agent Teams로 단계 실행 가능.

### 팀 모드 단계 개요

| 단계 | 서브에이전트 모드 | 팀 모드 | 조건 |
|-------|---------------|-----------|-----------|
| Plan | manager-plan (단일) | researcher + analyst + architect (병렬) | 복잡도 >= 임계값 |
| Develop | manager-ddd/tdd (순차) | backend-dev + frontend-dev + tester (병렬) | 도메인 >= 3 또는 파일 >= 10 |
| Report | (오케스트레이터) | (오케스트레이터) | 항상 단일 |

### 팀 모드 Plan 단계
- 병렬 리서치 팀을 위한 TeamCreate
- 팀원들이 코드베이스 탐색, 요구사항 분석, 접근 방식 설계
- 각 팀원은 Bash heredoc을 통해 `.do/jobs/`에 산출물 작성 (plan-mode 멱등성)
- Do가 플랜과 체크리스트로 종합
- 팀 셧다운 후 Develop 단계 전 /clear

### 팀 모드 Develop 단계
- 구현 팀을 위한 TeamCreate
- 파일 소유권 경계를 기반으로 작업 분해 (파일 하나 = 소유자 하나)
- 팀원들이 공유 체크리스트에서 자율적으로 작업 선택
- 각 팀원은 본인 파일만 커밋 (NEVER `git add -A`)
- 모든 구현 완료 후 품질 검증
- 팀 셧다운

### 모드 선택
- --team 플래그: 팀 모드 강제
- --solo 플래그: 서브에이전트 모드 강제
- 플래그 없음 (기본값): 복잡도 기반 자동 선택

### 폴백
팀 모드 실패 또는 사용 불가 시:
- 서브에이전트 모드로 자동 폴백
- 마지막으로 완료된 체크리스트 항목부터 재개 (체크리스트가 연속성 보장)
- 데이터 손실 또는 상태 오염 없음
