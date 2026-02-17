---
name: moai-workflow-tdd
description: >
  테스트 우선 소프트웨어 개발을 위한 RED-GREEN-REFACTOR 사이클을
  사용하는 테스트 주도 개발 워크플로우 전문가입니다.
  처음부터 새 기능을 개발하거나, 독립된 모듈을 생성하거나,
  동작 명세가 구현을 주도할 때 사용합니다.
  기존 코드 리팩토링에는 사용하지 마세요 (moai-workflow-ddd 사용).
  동작 보존이 주요 목표인 경우에도 사용하지 마세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Write Edit Bash(pytest:*) Bash(ruff:*) Bash(npm:*) Bash(npx:*) Bash(node:*) Bash(jest:*) Bash(vitest:*) Bash(go:*) Bash(cargo:*) Bash(mix:*) Bash(uv:*) Bash(bundle:*) Bash(php:*) Bash(phpunit:*) Grep Glob mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-03"
  modularized: "true"
  tags: "workflow, tdd, test-driven, red-green-refactor, test-first"
  author: "MoAI-ADK Team"
  context: "fork"
  agent: "manager-tdd"
  related-skills: "moai-workflow-ddd, moai-workflow-testing, moai-foundation-quality"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["TDD", "test-driven development", "red-green-refactor", "test-first", "new feature", "greenfield"]
  phases: ["run"]
  agents: ["manager-tdd", "expert-backend", "expert-frontend", "expert-testing"]
---

# 테스트 주도 개발 (TDD) 워크플로우

## 개발 모드 설정 (중요)

[참고] 이 워크플로우는 `.moai/config/sections/quality.yaml`을 기반으로 선택됩니다:

```yaml
constitution:
  development_mode: hybrid    # 또는 ddd, tdd
  hybrid_settings:
    new_features: tdd        # 새 코드 → TDD 사용 (이 워크플로우)
    legacy_refactoring: ddd  # 기존 코드 → DDD 사용
```

**이 워크플로우를 사용하는 경우**:
- `development_mode: tdd` → 항상 TDD 사용
- `development_mode: hybrid` + 새 패키지/모듈 → TDD 사용
- `development_mode: hybrid` + 기존 코드 리팩토링 → 대신 DDD 사용 (moai-workflow-ddd)

**핵심 구분**:
- **새 파일/패키지** (아직 존재하지 않음) → TDD (이 워크플로우)
- **기존 코드** (파일이 이미 존재함) → DDD (ANALYZE-PRESERVE-IMPROVE)

## 빠른 참조

테스트 주도 개발은 구현 전에 테스트가 기대 동작을 정의하는, 새로운 기능을 만들기 위한 체계적인 접근 방식을 제공합니다.

핵심 사이클 - RED-GREEN-REFACTOR:

- RED: 원하는 동작을 정의하는 실패하는 테스트 작성
- GREEN: 테스트를 통과시키는 최소한의 코드 작성
- REFACTOR: 테스트를 통과시키면서 코드 구조 개선

TDD 사용 시기:

- 처음부터 새 기능을 만들 때
- 기존 의존성이 없는 독립 모듈을 구축할 때
- 동작 명세가 개발을 주도할 때
- 명확한 계약이 있는 새 API 엔드포인트
- 정의된 동작을 가진 새 UI 컴포넌트
- Greenfield 프로젝트 (드물게 - 보통 Hybrid가 더 나음)

TDD를 사용하지 않는 경우:

- 기존 코드 리팩토링 (대신 DDD 사용)
- 동작 보존이 주요 목표인 경우
- 테스트 커버리지가 없는 레거시 코드베이스 (먼저 DDD 사용)
- 기존 파일을 수정하는 경우 (Hybrid 모드 고려)

---

## 핵심 철학

### TDD vs DDD 비교

TDD 접근 방식:

- 사이클: RED-GREEN-REFACTOR
- 목표: 테스트를 통해 새 기능 생성
- 시작점: 코드가 존재하지 않음
- 테스트 유형: 기대 동작을 정의하는 명세서 테스트
- 결과: 테스트 커버리지를 갖춘 새로운 작동 코드

DDD 접근 방식:

- 사이클: ANALYZE-PRESERVE-IMPROVE
- 목표: 동작 변경 없이 구조 개선
- 시작점: 정의된 동작을 가진 기존 코드
- 테스트 유형: 현재 동작을 포착하는 특성화 테스트
- 결과: 동일한 동작을 가진 더 잘 구조화된 코드

### 테스트 우선 원칙

TDD의 황금 규칙은 구현 코드 전에 테스트를 먼저 작성해야 한다는 것입니다:

- 테스트가 계약을 정의합니다
- 테스트가 기대 동작을 문서화합니다
- 테스트가 즉시 회귀를 잡아냅니다
- 구현은 테스트 요구사항에 의해 주도됩니다

---

## 구현 가이드

### 1단계: RED - 실패하는 테스트 작성

RED 단계는 실패하는 테스트를 통해 원하는 동작을 정의하는 데 집중합니다.

#### 효과적인 테스트 작성

구현 코드를 작성하기 전에:

- 요구사항을 명확히 이해합니다
- 테스트 형식으로 기대 동작을 정의합니다
- 한 번에 하나의 테스트를 작성합니다
- 테스트를 집중적이고 구체적으로 유지합니다
- 동작을 문서화하는 설명적인 테스트 이름을 사용합니다

#### 테스트 구조

Arrange-Act-Assert 패턴을 따릅니다:

- Arrange: 테스트 데이터와 의존성 설정
- Act: 테스트 대상 코드 실행
- Assert: 기대 결과 검증

#### 검증

테스트는 처음에 실패해야 합니다:

- 테스트가 실제로 무언가를 테스트한다는 것을 확인합니다
- 테스트가 우연히 통과하지 않음을 보장합니다
- 현재 상태와 원하는 상태 사이의 간격을 문서화합니다

### 2단계: GREEN - 테스트 통과시키기

GREEN 단계는 테스트를 만족시키는 최소한의 코드를 작성하는 데 집중합니다.

#### 최소한의 구현

테스트를 통과시키기에 충분한 코드만 작성합니다:

- 과도하게 설계하지 않습니다
- 테스트에서 요구하지 않는 기능을 추가하지 않습니다
- 완벽함이 아닌 정확성에 집중합니다
- 필요한 경우 값을 하드코딩합니다 (나중에 리팩토링)

#### 검증

테스트를 실행하여 통과하는지 확인합니다:

- 모든 단언이 성공해야 합니다
- 다른 테스트가 깨지지 않아야 합니다
- 구현이 테스트 요구사항을 만족시킵니다

### 3단계: REFACTOR - 코드 개선

REFACTOR 단계는 동작을 유지하면서 코드 품질을 개선하는 데 집중합니다.

#### 안전한 리팩토링

통과하는 테스트를 안전망으로 활용하여:

- 중복을 제거합니다
- 명명과 가독성을 개선합니다
- 메서드와 클래스를 추출합니다
- 적절한 경우 디자인 패턴을 적용합니다

#### 지속적인 검증

각 리팩토링 단계 후:

- 모든 테스트를 실행합니다
- 테스트가 실패하면 즉시 되돌립니다
- 테스트가 통과하면 커밋합니다

---

## TDD 워크플로우 실행

### 표준 TDD 세션

manager-tdd를 통해 TDD를 실행할 때:

1단계 - 요구사항 이해:

- 기능 범위를 위한 SPEC 문서 읽기
- 인수 기준에서 테스트 케이스 식별
- 테스트 구현 순서 계획

2단계 - RED 단계:

- 첫 번째 실패하는 테스트 작성
- 올바른 이유로 테스트가 실패하는지 확인
- 기대 동작 문서화

3단계 - GREEN 단계:

- 최소한의 구현 작성
- 테스트를 실행하여 통과하는지 확인
- 다음 테스트로 이동

4단계 - REFACTOR 단계:

- 개선을 위한 코드 검토
- 테스트를 안전망으로 활용하여 리팩토링 적용
- 깔끔한 코드 커밋

5단계 - 반복:

- RED-GREEN-REFACTOR 사이클 계속
- 모든 요구사항이 구현될 때까지
- 모든 인수 기준이 통과될 때까지

### TDD 루프 패턴

여러 테스트 케이스가 필요한 기능의 경우:

- 처음에 모든 테스트 케이스를 식별합니다
- 의존성과 복잡성에 따라 우선순위를 정합니다
- 각각에 대해 RED-GREEN-REFACTOR를 실행합니다
- 누적 테스트 스위트를 유지합니다

---

## 품질 지표

### TDD 성공 기준

테스트 커버리지 (필수):

- 커밋당 최소 80% 커버리지
- 새 코드에 대해 90% 권장
- 모든 공개 인터페이스 테스트

코드 품질 (목표):

- 모든 테스트 통과
- 구현 이후에 테스트 작성 없음
- 동작을 문서화하는 명확한 테스트 이름
- 테스트를 만족시키는 최소한의 구현

### TDD 특화 TRUST 검증

TDD 포커스로 TRUST 5 프레임워크 적용:

- Testability: 테스트 우선 접근 방식이 테스트 가능성 보장
- Readability: 테스트가 기대 동작 문서화
- Understandability: 테스트가 살아있는 문서 역할
- Security: 구현 전에 보안 테스트 작성
- Transparency: 테스트 실패가 즉각적인 피드백 제공

---

## 통합 지점

### DDD 워크플로우와의 통합

TDD와 DDD는 상호 보완적입니다:

- 새 코드에는 TDD
- 기존 코드 리팩토링에는 DDD
- Hybrid 모드에서는 두 접근 방식을 결합

### 테스팅 워크플로우와의 통합

TDD는 테스팅 워크플로우와 통합됩니다:

- 명세서 테스트 사용
- 커버리지 도구와 통합
- 테스트 품질을 위한 뮤테이션 테스팅 지원

### 품질 프레임워크와의 통합

TDD 출력은 품질 평가에 반영됩니다:

- 커버리지 지표 추적
- 변경사항에 대한 TRUST 5 검증
- 품질 게이트가 기준 적용

---

## 문제 해결

### 일반적인 이슈

테스트가 너무 복잡한 경우:

- 더 작고 집중적인 테스트로 분해합니다
- 한 번에 하나의 동작을 테스트합니다
- 복잡한 설정에는 테스트 픽스처를 사용합니다

구현이 너무 빠르게 성장하는 경우:

- 테스트되지 않은 기능을 구현하려는 충동을 억제합니다
- 새 기능에는 RED 단계로 돌아갑니다
- GREEN 단계를 최소한으로 유지합니다

리팩토링이 테스트를 깨뜨리는 경우:

- 즉시 되돌립니다
- 더 작은 단계로 리팩토링합니다
- 테스트가 구현이 아닌 동작을 검증하는지 확인합니다

### 복구 절차

TDD 규율이 무너지면:

- 멈추고 현재 상태를 평가합니다
- 기존 코드에 대한 특성화 테스트를 작성합니다
- 나머지 기능에 대해 TDD를 재개합니다
- Hybrid 모드로 전환하는 것을 고려합니다

---

Version: 1.0.0
Status: Active
Last Updated: 2026-02-03
