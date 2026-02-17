# Context7 테스트 패턴과 모범 사례

> 모듈: 테스트 패턴, Context7 통합, 업계 모범 사례
> 복잡도: 고급
> 소요 시간: 15분 이상
> 의존성: Python 3.8+, pytest, Context7 MCP, unittest.mock

## Context7 DDD 통합

```python
class Context7DDDIntegration:
    """DDD 패턴과 모범 사례를 위한 Context7 통합."""

    def __init__(self, context7_client=None):
        self.context7 = context7_client
        self.pattern_cache = {}

    async def load_ddd_patterns(self, language: str = "python") -> Dict[str, Any]:
        """Context7에서 DDD 패턴과 모범 사례 로드."""

        cache_key = f"ddd_patterns_{language}"
        if cache_key in self.pattern_cache:
            return self.pattern_cache[cache_key]

        patterns = {}

        if self.context7:
            try:
                # DDD 모범 사례 로드
                ddd_patterns = await self.context7.get_library_docs(
                    context7_library_id="/testing/pytest",
                    topic="DDD ANALYZE-PRESERVE-IMPROVE patterns best practices 2025",
                    tokens=4000
                )
                patterns['ddd_best_practices'] = ddd_patterns

                # 특정 언어 테스트 패턴 로드
                if language == "python":
                    python_patterns = await self.context7.get_library_docs(
                        context7_library_id="/python/pytest",
                        topic="advanced testing patterns mocking fixtures 2025",
                        tokens=3000
                    )
                    patterns['python_testing'] = python_patterns

                # 단언 패턴 로드
                assertion_patterns = await self.context7.get_library_docs(
                    context7_library_id="/testing/assertions",
                    topic="assertion patterns error messages test design 2025",
                    tokens=2000
                )
                patterns['assertions'] = assertion_patterns

                # mocking 패턴 로드
                mocking_patterns = await self.context7.get_library_docs(
                    context7_library_id="/python/unittest-mock",
                    topic="mocking strategies test doubles isolation patterns 2025",
                    tokens=3000
                )
                patterns['mocking'] = mocking_patterns

            except Exception as e:
                print(f"Failed to load Context7 patterns: {e}")
                patterns = self._get_default_patterns()
        else:
            patterns = self._get_default_patterns()

        self.pattern_cache[cache_key] = patterns
        return patterns

    def _get_default_patterns(self) -> Dict[str, Any]:
        """Context7 사용 불가 시 기본 DDD 패턴 반환."""
        return {
            'ddd_best_practices': {
                'analyze_phase': [
                    "기존 코드 구조와 패턴 파악",
                    "코드 읽기를 통해 현재 동작 식별",
                    "의존성과 사이드 이펙트 문서화",
                    "테스트 커버리지 공백 매핑"
                ],
                'preserve_phase': [
                    "기존 동작을 위한 특성화 테스트 작성",
                    "현재 동작을 황금 표준으로 포착",
                    "현재 구현으로 테스트가 통과함을 확인",
                    "복잡한 출력에 대한 동작 스냅샷 생성"
                ],
                'improve_phase': [
                    "테스트를 녹색으로 유지하며 코드 리팩토링",
                    "작고 점진적인 변경 수행",
                    "변경 후마다 테스트 실행",
                    "동작 보존 유지"
                ]
            },
            'python_testing': {
                'pytest_features': [
                    "여러 시나리오를 위한 매개변수화 테스트",
                    "테스트 설정과 정리를 위한 fixture",
                    "테스트 분류를 위한 마커",
                    "기능 강화를 위한 플러그인"
                ],
                'assertions': [
                    "pytest의 assert 문 사용",
                    "명확한 에러 메시지 제공",
                    "pytest.raises로 예상 예외 테스트",
                    "부동소수점 비교에 pytest.approx 사용"
                ]
            },
            'assertions': {
                'best_practices': [
                    "가능하면 테스트당 하나의 단언",
                    "명확하고 설명적인 단언 메시지",
                    "긍정적 경우와 부정적 경우 모두 테스트",
                    "적절한 단언 메서드 사용"
                ]
            },
            'mocking': {
                'strategies': [
                    "외부 의존성 mock 처리",
                    "테스트 가능성을 위한 의존성 주입 사용",
                    "복잡한 객체를 위한 테스트 더블 생성",
                    "mock과의 상호작용 검증"
                ]
            }
        }
```

## 테스트 패턴

### Given-When-Then 패턴

```python
def test_user_authentication_valid_credentials():
    """
    유효한 자격 증명으로 사용자 인증 테스트.

    Given: 유효한 자격 증명을 가진 등록된 사용자
    When: 사용자가 인증을 시도할 때
    Then: 시스템이 유효한 인증 토큰을 반환해야 함
    """
    # Given
    user = User(email="test@example.com", password="secure_password")
    auth_service = AuthenticationService()

    # When
    result = auth_service.authenticate(user.email, user.password)

    # Then
    assert result is not None
    assert result.token is not None
    assert result.expires_at > datetime.now()
```

### Arrange-Act-Assert 패턴

```python
def test_calculate_total_price_with_discount():
    """
    할인이 적용된 총 가격 계산 테스트.
    """
    # Arrange
    cart = ShoppingCart()
    cart.add_item("item1", price=100.0, quantity=2)
    cart.add_item("item2", price=50.0, quantity=1)
    discount_code = "SAVE10"

    # Act
    total = cart.calculate_total(discount_code)

    # Assert
    assert total == 225.0  # (200 + 50) * 0.9
```

### 매개변수화 테스트 패턴

```python
@pytest.mark.parametrize("input,expected", [
    (2, 4),    # 2^2 = 4
    (3, 9),    # 3^2 = 9
    (0, 0),    # 0^2 = 0
    (-1, 1),   # (-1)^2 = 1
    (10, 100)  # 10^2 = 100
])
def test_square_function(input, expected):
    """다양한 입력으로 제곱 함수 테스트."""
    result = square(input)
    assert result == expected
```

### 예외 테스트 패턴

```python
def test_divide_by_zero_raises_exception():
    """0으로 나누기가 ZeroDivisionError를 발생시키는지 테스트."""
    with pytest.raises(ZeroDivisionError) as exc_info:
        divide(10, 0)

    assert "division by zero" in str(exc_info.value)
```

### Mock 테스트 패턴

```python
def test_external_api_call_with_mock():
    """mock 응답으로 외부 API 호출 테스트."""
    # mock 생성
    mock_api = Mock()
    mock_api.get_data.return_value = {"status": "success", "data": [1, 2, 3]}

    # 테스트에서 mock 사용
    service = DataService(api_client=mock_api)
    result = service.fetch_data()

    # 상호작용 검증
    mock_api.get_data.assert_called_once()
    assert result == [1, 2, 3]
```

## Pytest Fixture

### 기본 Fixture

```python
@pytest.fixture
def sample_user():
    """테스트용 샘플 사용자 생성."""
    return User(
        email="test@example.com",
        username="testuser",
        password="secure_password"
    )

def test_user_email(sample_user):
    """사용자 이메일 속성 테스트."""
    assert sample_user.email == "test@example.com"
```

### 설정과 정리가 있는 Fixture

```python
@pytest.fixture
def database_connection():
    """정리 포함 데이터베이스 연결 생성."""
    # 설정
    conn = Database.connect(":memory:")
    conn.create_tables()

    yield conn  # 테스트에 연결 제공

    # 정리
    conn.close()

def test_database_query(database_connection):
    """fixture로 데이터베이스 쿼리 테스트."""
    result = database_connection.query("SELECT * FROM users")
    assert len(result) >= 0
```

### 매개변수화 Fixture

```python
@pytest.fixture(params=[
    ("valid_email@example.com", True),
    ("invalid_email", False),
    ("", False)
])
def email_validation_data(request):
    """이메일 검증 테스트 데이터 제공."""
    return request.param

def test_email_validation(email_validation_data):
    """다양한 입력으로 이메일 검증 테스트."""
    email, expected_valid = email_validation_data
    result = validate_email(email)
    assert result.is_valid == expected_valid
```

## 테스트 구성

### 테스트 탐색 구조

```
tests/
├── unit/              # 개별 컴포넌트의 단위 테스트
│   ├── test_models.py
│   ├── test_services.py
│   └── test_utils.py
├── integration/       # 컴포넌트 상호작용의 통합 테스트
│   ├── test_api_integration.py
│   └── test_database_integration.py
├── acceptance/        # 사용자 시나리오의 인수 테스트
│   ├── test_user_scenarios.py
│   └── test_business_workflows.py
└── conftest.py        # 공유 fixture와 설정
```

### 테스트 마커

```python
import pytest

# 사용자 정의 마커 정의
pytest.mark.slow = pytest.mark.slow
pytest.mark.integration = pytest.mark.integration
pytest.mark.unit = pytest.mark.unit

# 테스트에서 마커 사용
@pytest.mark.unit
def test_individual_function():
    """개별 함수에 대한 단위 테스트."""
    assert calculate(2, 2) == 4

@pytest.mark.integration
def test_database_integration():
    """데이터베이스 통합 테스트."""
    result = db.query("SELECT * FROM users")
    assert result is not None

@pytest.mark.slow
def test_performance_benchmark():
    """느린 성능 테스트."""
    result = expensive_operation()
    assert result is not None
```

## 테스트 커버리지

### 커버리지 분석 실행

```bash
# 커버리지와 함께 테스트 실행
pytest --cov=src --cov-report=html --cov-report=term

# 커버리지 보고서 생성
pytest --cov=src --cov-report=html
```

### 커버리지 설정

```ini
# .coveragerc
[run]
source = src
omit =
    */tests/*
    */__init__.py

[report]
exclude_lines =
    pragma: no cover
    def __repr__
    raise AssertionError
    raise NotImplementedError
    if __name__ == .__main__.:
```

## 모범 사례

1. 테스트 격리: 각 테스트는 독립적이며 다른 테스트에 의존하지 않아야 함
2. 설명적인 이름: 테스트 이름이 무엇을 테스트하는지 명확히 설명
3. 테스트당 하나의 단언: 단일 동작에 테스트를 집중
4. Arrange-Act-Assert: 이 패턴으로 테스트를 명확하게 구조화
5. 외부 의존성 Mock: 외부 서비스와 데이터베이스에 mock 사용
6. 엣지 케이스 테스트: 경계 조건과 에러 케이스에 대한 테스트 포함
7. 빠른 테스트: 빠른 피드백을 위해 단위 테스트를 빠르게 유지
8. 유지 보수 가능한 테스트: 테스트를 단순하고 이해하기 쉽게 유지
9. Context7 통합: 최신 테스트 패턴과 모범 사례를 위해 Context7 활용
10. 지속적 테스트: 모든 코드 변경에 자동으로 테스트 실행

---

관련: [ANALYZE-PRESERVE-IMPROVE](./analyze-preserve-improve.md) | [테스트 생성](./test-generation.md)
