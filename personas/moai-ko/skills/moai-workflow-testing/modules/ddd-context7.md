# Context7 통합 도메인 주도 개발

> 모듈: Context7 패턴과 AI 기반 테스팅을 포함한 ANALYZE-PRESERVE-IMPROVE DDD 사이클
> 복잡도: 고급
> 소요 시간: 25분+
> 의존성: Python 3.8+, pytest, Context7 MCP, unittest, asyncio

## 개요

DDD Context7 통합은 AI 기반 테스트 생성, Context7 강화 테스팅 패턴, 자동화된 모범 사례 적용을 포함한 포괄적인 도메인 주도 개발 워크플로우를 제공합니다.

### 주요 기능

- AI 기반 테스트 생성: 명세에서 포괄적인 테스트 스위트 생성
- Context7 통합: 최신 테스팅 패턴 및 모범 사례 접근
- ANALYZE-PRESERVE-IMPROVE 사이클: 완전한 DDD 워크플로우 구현
- 고급 테스팅: 속성 기반 테스팅, 변이 테스팅, 지속적 테스팅
- 테스트 패턴: 포괄적인 테스팅 패턴 및 픽스처 라이브러리

## 빠른 시작

### 기본 DDD 사이클

```python
from moai_workflow_testing import DDDManager, TestSpecification, TestType

# DDD Manager 초기화
ddd_manager = DDDManager(
    project_path="/path/to/project",
    context7_client=context7
)

# DDD 세션 시작
session = await ddd_manager.start_ddd_session("user_authentication_refactor")

# 테스트 명세 생성
test_spec = TestSpecification(
    name="test_user_login_behavior_preservation",
    description="리팩토링 중 기존 로그인 동작 보존",
    test_type=TestType.CHARACTERIZATION,
    requirements=[
        "기존 로그인 흐름이 계속 작동해야 함",
        "오류 메시지가 일관성 있게 유지되어야 함"
    ],
    acceptance_criteria=[
        "유효한 자격증명이 사용자 토큰을 반환함 (기존 동작)",
        "잘못된 자격증명이 동일한 오류 메시지를 발생시킴"
    ],
    edge_cases=[
        "빈 이메일로 테스트 (기존 동작)",
        "빈 비밀번호로 테스트 (기존 동작)"
    ]
)

# 전체 DDD 사이클 실행
cycle_results = await ddd_manager.run_full_ddd_cycle(
    specification=test_spec,
    target_function="authenticate_user"
)
```

## 핵심 컴포넌트

### DDD 사이클 단계

1. ANALYZE 단계: 기존 코드 이해
   - 기존 코드 구조 및 패턴 분석
   - 코드 읽기를 통한 현재 동작 식별
   - 의존성 및 부작용 문서화
   - 테스트 커버리지 공백 매핑

2. PRESERVE 단계: 특성화 테스트 생성
   - 기존 동작에 대한 특성화 테스트 작성
   - 현재 동작을 "황금 표준"으로 캡처
   - 현재 구현으로 테스트가 통과하는지 확인
   - 복잡한 출력에 대한 동작 스냅샷 생성

3. IMPROVE 단계: 동작 보존과 함께 리팩토링
   - 테스트를 통과하는 상태를 유지하며 코드 리팩토링
   - 작고 점진적인 변경 수행
   - 각 변경 후 테스트 실행
   - 동작 보존 유지

4. REVIEW 단계: 검증 및 커밋
   - 모든 특성화 테스트가 여전히 통과하는지 확인
   - 코드 품질 및 문서화 검토
   - 동작 변경 사항 확인
   - 명확한 메시지로 변경 사항 커밋

### Context7 통합

DDD Context7 통합은 다음을 제공합니다:

- 패턴 로딩: Context7에서 최신 테스팅 패턴 접근
- AI 테스트 생성: Context7 패턴을 사용한 강화된 테스트 생성
- 모범 사례: 업계 표준 테스팅 관행
- 엣지 케이스 감지: 자동 엣지 케이스 식별
- 테스트 제안: AI 기반 테스트 개선 제안

## 모듈 구조

### 핵심 모듈

**ANALYZE-PRESERVE-IMPROVE 구현** (`ddd-context7/analyze-preserve-improve.md`)
- DDD 사이클 구현
- 테스트 실행 및 검증
- 커버리지 분석
- 세션 관리

**테스트 생성** (`ddd-context7/test-generation.md`)
- AI 기반 테스트 생성
- 명세 기반 생성
- Context7 강화 생성
- 템플릿 기반 생성

**테스트 패턴** (`ddd-context7/test-patterns.md`)
- 테스팅 패턴 및 모범 사례
- Pytest 픽스처 및 구성
- 테스트 탐색 구조
- 커버리지 설정

**고급 기능** (`ddd-context7/advanced-features.md`)
- 포괄적인 테스트 스위트 생성
- 속성 기반 테스팅
- 변이 테스팅
- 지속적 테스팅

## 일반적인 사용 사례

### 동작 보존

```python
# 특성화 테스트 명세
char_spec = TestSpecification(
    name="test_calculate_sum_existing_behavior",
    description="기존 합계 계산 동작 보존",
    test_type=TestType.CHARACTERIZATION,
    requirements=["함수가 두 수를 더해야 함 (기존 동작)"],
    acceptance_criteria=["현재 구현대로 올바른 합계 반환"],
    edge_cases=["0 값", "음수", "큰 수"]
)

test_code = await test_generator.generate_test_case(char_spec)
```

### 테스트와 함께 리팩토링

```python
# 리팩토링을 위한 통합 테스트 명세
refactor_spec = TestSpecification(
    name="test_database_integration_refactor",
    description="리팩토링 중 데이터베이스 동작 보존",
    test_type=TestType.INTEGRATION,
    requirements=["데이터베이스 연결", "쿼리 실행"],
    acceptance_criteria=["이전과 동일하게 연결 성공", "쿼리가 동일한 데이터 반환"],
    edge_cases=["연결 실패 처리", "빈 결과", "대용량 데이터셋"]
)
```

### 예외 동작 보존

```python
# 예외 테스트 명세
exception_spec = TestSpecification(
    name="test_divide_by_zero_existing_behavior",
    description="0으로 나누기 예외 처리 보존",
    test_type=TestType.CHARACTERIZATION,
    requirements=["나누기 함수", "에러 처리"],
    acceptance_criteria=["이전과 동일한 ZeroDivisionError 발생"],
    edge_cases=["제수가 0", "피제수가 0"]
)
```

## 모범 사례

### 테스트 설계

1. 특성화 우선: 코드 변경 전에 기존 동작을 캡처하는 테스트 작성
2. 서술적 이름: 테스트 이름은 어떤 동작이 보존되는지 명확하게 설명해야 함
3. Arrange-Act-Assert: 명확성을 위해 이 패턴으로 테스트 구성
4. 독립적인 테스트: 테스트가 서로에게 의존하지 않아야 함
5. 빠른 실행: 빠른 피드백을 위해 테스트를 빠르게 유지

### Context7 통합

1. 패턴 로딩: 최신 모범 사례를 위해 Context7 패턴 로드
2. 엣지 케이스 감지: 누락된 엣지 케이스를 식별하기 위해 Context7 사용
3. 테스트 제안: 테스트 개선을 위해 AI 제안 활용
4. 품질 분석: 테스트 품질 분석을 위해 Context7 사용

### DDD 워크플로우

1. 먼저 분석: 코드 변경 전에 항상 기존 동작 이해
2. 테스트로 보존: 리팩토링 전에 특성화 테스트 생성
3. 테스트를 통과 상태로 유지: 실패하는 테스트는 커밋하지 않음
4. 작은 증분: 작고 점진적인 변경 수행
5. 지속적 테스팅: 매 변경 후 테스트 실행

## 고급 기능

### 속성 기반 테스팅

많은 무작위 입력에 걸쳐 코드 속성을 검증하기 위해 Hypothesis를 사용한 속성 기반 테스팅을 활용하세요.

### 변이 테스팅

코드 변이를 도입하고 테스트가 이를 감지하는지 확인하여 테스트 스위트 품질을 검증하기 위해 변이 테스팅을 사용하세요.

### 지속적 테스팅

파일 변경 시 자동 테스트 실행을 위한 감시 모드를 구현하세요.

### AI 기반 생성

지능적인 테스트 생성 및 제안을 위해 Context7을 활용하세요.

## 성능 고려사항

- 테스트 실행: 더 빠른 피드백을 위해 병렬 테스트 실행 사용
- 테스트 격리: 간섭 방지를 위해 테스트가 격리되어 있는지 확인
- 외부 의존성 모킹: 빠르고 신뢰할 수 있는 테스트를 위해 외부 서비스 모킹
- 설정 최적화: 효율적인 테스트 설정을 위해 픽스처 및 테스트 팩토리 사용

## 문제 해결

### 일반적인 이슈

1. 간헐적으로 실패하는 테스트
   - 테스트 간 공유 상태 확인
   - 테스트 격리 검증
   - 픽스처에 적절한 정리 추가

2. 느린 테스트 실행
   - 병렬 테스트 실행 사용
   - 외부 의존성 모킹
   - 테스트 설정 최적화

3. Context7 통합 이슈
   - Context7 클라이언트 설정 확인
   - 네트워크 연결 확인
   - 폴백으로 기본 패턴 사용

## 리소스

### 상세 모듈

- [ANALYZE-PRESERVE-IMPROVE 구현](./ddd-context7/analyze-preserve-improve.md) - 핵심 DDD 사이클
- [테스트 생성](./ddd-context7/test-generation.md) - AI 기반 생성
- [테스트 패턴](./ddd-context7/test-patterns.md) - 패턴 및 모범 사례
- [고급 기능](./ddd-context7/advanced-features.md) - 고급 테스팅 기법

### 관련 모듈

- [AI 디버깅](./ai-debugging.md) - 디버깅 기법
- [성능 최적화](./performance-optimization.md) - 성능 테스팅
- [스마트 리팩토링](./smart-refactoring.md) - 테스트와 함께 리팩토링

### 외부 리소스

- [Pytest 문서](https://docs.pytest.org/)
- [Python 테스팅 모범 사례](https://docs.python-guide.org/writing/tests/)
- [Hypothesis 속성 기반 테스팅](https://hypothesis.works/)
- [Context7 MCP 문서](https://context7.io/docs)

---

모듈: `modules/ddd-context7.md`
Version: 2.0.0 (DDD 마이그레이션)
Last Updated: 2026-01-17
