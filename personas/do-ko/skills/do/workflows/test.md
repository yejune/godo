---
name: do-workflow-test
description: >
  품질 보증을 위한 TDD RED-GREEN-REFACTOR 워크플로우. 실제 DB 테스트와
  AI 안티패턴 방지로 테스트 우선 개발 규율을 강제하는
  Do 전용 워크플로우.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-15"
  tags: "test, tdd, red-green-refactor, quality, coverage"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 4000

# Do Extension: Triggers
triggers:
  keywords: ["test", "tdd", "coverage", "테스트", "red-green", "refactor"]
  agents: ["expert-testing", "manager-tdd"]
  phases: ["test"]
---

# 테스트 워크플로우 오케스트레이션

## 목적

테스트 주도 품질 보증을 위한 TDD RED-GREEN-REFACTOR 사이클을 실행한다. 이 워크플로우는 단독으로 호출하거나 do 파이프라인의 일부로 실행할 수 있다. 테스트 우선 규율을 강제하고 일반적인 AI 테스트 안티패턴을 방지한다.

## 범위

- 독립 실행형 TDD 워크플로우 또는 do 파이프라인에 통합
- 전체적으로 dev-testing.md 규칙 강제
- 새 코드 (TDD)와 기존 코드 (DDD 특성화 테스트) 모두 지원

## 입력

- $ARGUMENTS: 테스트할 기능 또는 컴포넌트, 또는 체크리스트 참조
- 컨텍스트: 테스트할 기존 코드, 또는 새 코드에 대한 명세

---

## 단계 순서

### Phase 1: 테스트 전략 평가

코드 상태에 따라 테스트 접근 방식 결정:

**새 코드 (구현 없음)**:
- TDD 적용: RED -> GREEN -> REFACTOR
- 에이전트: manager-tdd 서브에이전트

**기존 코드 (수정 또는 리팩토링)**:
- DDD 적용: ANALYZE -> PRESERVE (특성화 테스트) -> IMPROVE
- 에이전트: manager-ddd 서브에이전트

**혼합 (기존 파일에 새 함수 추가)**:
- 새 함수는 TDD, 기존 동작은 특성화 테스트
- 에이전트: DDD 컨텍스트와 함께 manager-tdd

### Phase 2: RED - 실패하는 테스트 작성

에이전트: Task(expert-testing) 또는 Task(manager-tdd)

[HARD] 구현 코드 작성 전에 반드시 테스트를 먼저 작성해야 한다.

태스크:
1. 원하는 동작을 설명하는 테스트 작성
2. 테스트가 실패하는지 검증 (새로운 것을 테스트한다는 확인)
3. 한 번에 하나의 테스트, 집중적이고 구체적으로
4. 테스트 이름은 시나리오를 설명: `test_login_fails_with_expired_token`
5. Arrange-Act-Assert (Given-When-Then) 패턴 사용

품질 규칙 (dev-testing.md에서):
- [HARD] 구현 세부사항이 아닌 동작(behavior) 테스트
- [HARD] FIRST: Fast, Independent, Repeatable, Self-validating, Timely
- [HARD] Real DB만 -- mock DB, in-memory DB, SQLite 대체 금지
- [HARD] 각 테스트가 자체 데이터 설정 및 정리 (트랜잭션 롤백 또는 truncate)
- [HARD] 테스트 데이터 생성은 factory/builder 패턴 사용
- [HARD] 병렬 안전을 위해 테스트별 고유 식별자 (UUID/timestamp suffix)

### Phase 3: GREEN - 최소 구현

에이전트: Task(expert-backend) 또는 Task(expert-frontend) 또는 직접 (Focus 모드)

[HARD] 테스트를 통과하는 가장 단순한 코드를 작성한다.

규칙:
- 조기 최적화 금지
- 조기 추상화 금지
- 우아함이 아닌 정확성에 집중
- 변경 후 테스트를 실행하여 GREEN 상태 확인

### Phase 4: REFACTOR - 품질 개선

에이전트: Phase 3와 동일

[HARD] 테스트를 GREEN으로 유지하면서 코드를 개선한다.

태스크:
1. 패턴 추출, 중복 제거
2. 적절한 곳에 SOLID 원칙 적용
3. 각 리팩토링 단계 후 모든 테스트 실행
4. 동작 변경 없음 -- 테스트는 전체적으로 green을 유지해야 함

### Phase 5: 검증

포괄적인 검증 실행:

1. 개별 테스트 실행: 새 테스트 통과 확인
2. 전체 테스트 스위트 실행: 회귀 없음 확인
3. 커버리지 확인: 수정된 코드의 목표 85%+
4. 교차 검증: view <-> DB <-> worker <-> model <-> business logic <-> controller

---

## AI 안티패턴 방지 [HARD]

테스트 작성 시 다음은 금지:

- [HARD] 단언 약화: `assertEqual`을 `assertContains`로, 정확한 값을 `any()`로 변경 금지
- [HARD] try/catch로 오류를 삼켜 테스트를 녹색으로 만들기 금지
- [HARD] 잘못된 출력에 맞춰 기대값 조정 금지 (테스트가 아닌 코드를 수정할 것)
- [HARD] `time.sleep()` 또는 임의 지연 사용 금지 (타이밍 문제의 진짜 원인 찾기)
- [HARD] 실패하는 테스트 삭제 또는 주석 처리 금지
- [HARD] 정확한 값을 알 때 와일드카드 매처 (`any()`, `mock.ANY`) 사용 금지
- [HARD] happy path만 테스트 금지 -- 에러 경로, 엣지 케이스, 경계값 테스트 필수

### 변이 테스트 사고방식

모든 테스트에 다음 사고 연습 적용: "이 코드 한 줄을 바꾸면 테스트가 실패하는가?"
실패하지 않으면 테스트 커버리지가 부족한 것이다. 이것은 도구 요구사항이 아니라
테스트 작성 시 모든 에이전트가 적용해야 하는 규율이다.

---

## 테스트 실행 규칙 [HARD]

- [HARD] 모든 테스트는 반드시 통과 -- skip, timeout bypass 금지
- [HARD] "시간 초과됨"은 허용되지 않는 보고 -- 반드시 해결
- [HARD] 테스트 실패 시 테스트가 아닌 코드를 수정
- [HARD] 먼저 개별 테스트, 확신이 생기면 전체 스위트 실행
- [HARD] 개별 테스트 < 3분, 전체 스위트 < 10분
- [HARD] 테스트는 Docker 컨테이너 내부에서 실행: `docker compose exec <service> <test-command>`

---

## 특성화 테스트 (기존 코드)

수정이 필요한 기존 코드의 경우:

1. ANALYZE: 기존 코드 읽기, 동작 및 의존성 파악
2. PRESERVE: 현재 동작을 포착하는 특성화 테스트 작성
   - 이 테스트들은 코드가 지금 하는 일을 문서화 (해야 하는 일이 아님)
   - 리팩토링을 위한 안전망 역할
3. IMPROVE: 특성화 테스트를 가드레일로 사용하여 변경
   - 각 변경 후 테스트 실행
   - 테스트 실패는 의도치 않은 동작 변경을 의미

---

## 완료 기준

- RED: 실패하는 테스트 작성 및 실패 검증됨
- GREEN: 최소 구현으로 테스트 통과
- REFACTOR: 코드 개선됨, 모든 테스트 여전히 통과
- 커버리지: 수정/새 코드의 85%+
- 전체 스위트: 모든 테스트 통과 (회귀 없음)
- 안티패턴: 금지된 패턴 미사용
- 커밋: 테스트와 구현이 설명적인 메시지와 함께 함께 커밋됨

---

Version: 1.0.0
Updated: 2026-02-16
