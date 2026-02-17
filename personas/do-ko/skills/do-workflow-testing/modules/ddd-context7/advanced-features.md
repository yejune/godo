# Context7를 활용한 고급 DDD 기능

> 모듈: AI 기반 포괄적 테스트 스위트 생성 및 분석
> 복잡도: 전문가
> 소요 시간: 25분 이상
> 의존성: Python 3.8+, pytest, Context7 MCP, AST 분석, asyncio

## 강화된 테스트 생성기

```python
import ast
import inspect

class EnhancedTestGenerator(TestGenerator):
    """Context7 통합이 강화된 테스트 생성기."""

    async def generate_comprehensive_test_suite(
        self, function_code: str,
        context7_patterns: Dict[str, Any]
    ) -> List[str]:
        """함수 코드에서 포괄적인 테스트 스위트 생성."""

        # AST를 사용해 함수 분석
        function_analysis = self._analyze_function_code(function_code)

        # 다양한 시나리오에 대한 테스트 생성
        test_cases = []

        # 정상 경로 테스트
        happy_path_tests = await self._generate_happy_path_tests(
            function_analysis, context7_patterns
        )
        test_cases.extend(happy_path_tests)

        # 엣지 케이스 테스트
        edge_case_tests = await self._generate_edge_case_tests(
            function_analysis, context7_patterns
        )
        test_cases.extend(edge_case_tests)

        # 에러 처리 테스트
        error_tests = await self._generate_error_handling_tests(
            function_analysis, context7_patterns
        )
        test_cases.extend(error_tests)

        # 성능 임계 함수에 대한 성능 테스트
        if self._is_performance_critical(function_analysis):
            perf_tests = await self._generate_performance_tests(function_analysis)
            test_cases.extend(perf_tests)

        return test_cases

    def _analyze_function_code(self, code: str) -> Dict[str, Any]:
        """테스트 요구사항 추출을 위한 함수 코드 분석."""

        try:
            tree = ast.parse(code)

            analysis = {
                'functions': [],
                'parameters': [],
                'return_statements': [],
                'exceptions': [],
                'external_calls': []
            }

            for node in ast.walk(tree):
                if isinstance(node, ast.FunctionDef):
                    analysis['functions'].append({
                        'name': node.name,
                        'args': [arg.arg for arg in node.args.args],
                        'decorators': [d.id if isinstance(d, ast.Name) else str(d)
                                      for d in node.decorator_list]
                    })

                elif isinstance(node, ast.Raise):
                    analysis['exceptions'].append({
                        'type': node.exc.func.id if node.exc and hasattr(node.exc, 'func') else 'Exception',
                        'message': node.exc.msg if node.exc and hasattr(node.exc, 'msg') else None
                    })

                elif isinstance(node, ast.Call):
                    if isinstance(node.func, ast.Attribute):
                        analysis['external_calls'].append(f"{node.func.value.id}.{node.func.attr}")
                    elif isinstance(node.func, ast.Name):
                        analysis['external_calls'].append(node.func.id)

            return analysis

        except Exception as e:
            print(f"Error analyzing function code: {e}")
            return {}

    async def _generate_happy_path_tests(
        self, analysis: Dict[str, Any],
        context7_patterns: Dict[str, Any]
    ) -> List[str]:
        """정상 경로 테스트 케이스 생성."""

        tests = []

        for func in analysis.get('functions', []):
            # 정상 동작에 대한 테스트 생성
            test_code = f"""
def test_{func['name']}_happy_path():
    '''
    Test {func['name']} with valid inputs.

    Given: Valid input parameters
    When: {func['name']} is called
    Then: Expected result is returned
    '''
    # Arrange
    # Add setup code based on parameters: {', '.join(func['args'])}

    # Act
    # result = {func['name']}(*args)

    # Assert
    # assert result is not None
"""
            tests.append(test_code)

        return tests

    async def _generate_edge_case_tests(
        self, analysis: Dict[str, Any],
        context7_patterns: Dict[str, Any]
    ) -> List[str]:
        """엣지 케이스 테스트 시나리오 생성."""

        tests = []

        for func in analysis.get('functions', []):
            # 엣지 케이스 테스트 생성
            edge_cases = [
                ("empty_input", "Test with empty input"),
                ("null_input", "Test with None/null input"),
                ("boundary_value", "Test with boundary values"),
                ("max_input", "Test with maximum allowed input"),
                ("min_input", "Test with minimum allowed input")
            ]

            for case_name, description in edge_cases:
                test_code = f"""
def test_{func['name']}_{case_name}():
    '''
    {description}

    Given: Edge case input ({case_name})
    When: {func['name']} is called
    Then: Function handles edge case appropriately
    '''
    # Arrange
    # Setup edge case input

    # Act
    # result = {func['name']}(*edge_case_args)

    # Assert
    # Verify function handles edge case
"""
                tests.append(test_code)

        return tests

    async def _generate_error_handling_tests(
        self, analysis: Dict[str, Any],
        context7_patterns: Dict[str, Any]
    ) -> List[str]:
        """에러 처리 테스트 케이스 생성."""

        tests = []

        for exc in analysis.get('exceptions', []):
            exc_type = exc.get('type', 'Exception')
            test_code = f"""
def test_error_handling_{exc_type.lower()}():
    '''
    Test {exc_type} error handling.

    Given: Invalid input or error condition
    When: Function is called with invalid input
    Then: Appropriate exception is raised
    '''
    # Arrange
    # Setup invalid input

    # Act & Assert
    with pytest.raises({exc_type}):
        # function_call()
        pass
"""
            tests.append(test_code)

        return tests

    async def _generate_performance_tests(
        self, analysis: Dict[str, Any]
    ) -> List[str]:
        """성능 테스트 케이스 생성."""

        tests = []

        for func in analysis.get('functions', []):
            test_code = f"""
def test_{func['name']}_performance():
    '''
    Test {func['name']} performance characteristics.

    Given: Large input dataset
    When: {func['name']} is called
    Then: Function completes within acceptable time
    '''
    # Arrange
    import time
    large_input = list(range(10000))

    # Act
    start_time = time.time()
    # result = {func['name']}(large_input)
    execution_time = time.time() - start_time

    # Assert
    assert execution_time < 1.0, f"Function too slow: {{execution_time}}s"
"""
            tests.append(test_code)

        return tests

    def _is_performance_critical(self, analysis: Dict[str, Any]) -> bool:
        """함수가 성능 임계인지 판단."""

        # 성능 관련 지표 확인
        func_names = [f['name'] for f in analysis.get('functions', [])]

        performance_keywords = ['process', 'calculate', 'compute', 'parse', 'transform']

        return any(
            any(keyword in name.lower() for keyword in performance_keywords)
            for name in func_names
        )
```

## Context7 강화 테스트

```python
class Context7EnhancedTesting:
    """Context7 통합을 활용한 고급 테스트 기능."""

    def __init__(self, context7_client=None):
        self.context7 = context7_client

    async def get_intelligent_test_suggestions(
        self, codebase_context: Dict[str, Any]
    ) -> Dict[str, Any]:
        """Context7을 사용해 AI 기반 테스트 제안 가져오기."""

        if not self.context7:
            return self._get_rule_based_suggestions()

        try:
            # 고급 테스트 패턴 가져오기
            advanced_patterns = await self.context7.get_library_docs(
                context7_library_id="/testing/advanced",
                topic="property-based testing mutation testing 2025",
                tokens=5000
            )

            # 테스트 품질 지표 가져오기
            quality_patterns = await self.context7.get_library_docs(
                context7_library_id="/testing/quality",
                topic="test quality analysis coverage gaps 2025",
                tokens=3000
            )

            return {
                'advanced_patterns': advanced_patterns,
                'quality_metrics': quality_patterns,
                'suggestions': self._generate_intelligent_suggestions(
                    advanced_patterns, quality_patterns, codebase_context
                )
            }

        except Exception as e:
            print(f"Context7 test suggestions failed: {e}")
            return self._get_rule_based_suggestions()

    def _generate_intelligent_suggestions(
        self, advanced_patterns: Dict, quality_patterns: Dict, context: Dict
    ) -> List[str]:
        """지능적인 테스트 제안 생성."""

        suggestions = []

        # 코드베이스 컨텍스트 분석
        coverage = context.get('coverage_percentage', 0)

        if coverage < 80:
            suggestions.append("테스트 커버리지를 최소 80%까지 높이세요")

        # 누락된 테스트 유형 확인
        test_types = context.get('test_types', [])
        if 'integration' not in test_types:
            suggestions.append("컴포넌트 상호작용을 위한 통합 테스트 추가")

        if 'performance' not in test_types:
            suggestions.append("주요 경로에 대한 성능 테스트 추가")

        if 'security' not in test_types:
            suggestions.append("인증 및 권한 부여를 위한 보안 테스트 추가")

        return suggestions

    def _get_rule_based_suggestions(self) -> Dict[str, Any]:
        """규칙 기반 테스트 제안 가져오기."""

        return {
            'suggestions': [
                "리팩토링 전 기존 동작 분석 (DDD)",
                "높은 테스트 커버리지 목표 (80% 이상)",
                "긍정적 경우와 부정적 경우 모두 테스트",
                "외부 의존성에 mocking 사용",
                "여러 시나리오에 대한 테스트 매개변수화",
                "임계 함수에 성능 테스트 추가",
                "데이터 검증을 위한 속성 기반 테스트 구현",
                "변이 테스트로 테스트 품질 검증"
            ]
        }
```

## 속성 기반 테스트

```python
from hypothesis import given, strategies as st

class PropertyBasedTests:
    """Hypothesis를 사용한 속성 기반 테스트."""

    @given(st.integers(), st.integers())
    def test_addition_commutative(self, a, b):
        """덧셈이 교환 법칙을 만족하는지 테스트."""
        assert add(a, b) == add(b, a)

    @given(st.lists(st.integers()))
    def test_sort_idempotent(self, lst):
        """정렬이 멱등성을 만족하는지 테스트."""
        result = sort(lst)
        assert sort(result) == result

    @given(st.text())
    def test_reverse_inverse(self, text):
        """역순이 자기 자신의 역함수임을 테스트."""
        assert reverse(reverse(text)) == text

    @given(st.integers(min_value=0, max_value=1000))
    def test_square_non_negative(self, x):
        """어떤 수의 제곱도 음수가 아님을 테스트."""
        assert square(x) >= 0
```

## 변이 테스트

```python
class MutationTesting:
    """테스트 품질 검증을 위한 변이 테스트."""

    def __init__(self, project_path: str):
        self.project_path = Path(project_path)

    async def run_mutation_tests(self) -> Dict[str, Any]:
        """테스트 스위트 품질 확인을 위한 변이 테스트 실행."""

        try:
            # Python 변이 테스트 도구 mutmut 사용
            result = subprocess.run(
                [
                    sys.executable, '-m', 'mutmut',
                    'run',
                    '--paths-to-mutate', 'src'
                ],
                capture_output=True,
                text=True,
                cwd=str(self.project_path)
            )

            return self._parse_mutation_results(result.stdout)

        except Exception as e:
            return {'error': str(e), 'killed_mutants': 0, 'survived_mutants': 0}

    def _parse_mutation_results(self, output: str) -> Dict[str, Any]:
        """변이 테스트 결과 파싱."""

        # 출력 파싱으로 변이 통계 추출
        lines = output.split('\n')

        results = {
            'total_mutations': 0,
            'killed_mutants': 0,
            'survived_mutants': 0,
            'mutation_score': 0.0
        }

        for line in lines:
            if 'killed' in line.lower():
                parts = line.split()
                if len(parts) >= 2 and parts[0].isdigit():
                    results['killed_mutants'] = int(parts[0])

            elif 'survived' in line.lower():
                parts = line.split()
                if len(parts) >= 2 and parts[0].isdigit():
                    results['survived_mutants'] = int(parts[0])

        # 변이 점수 계산
        results['total_mutations'] = results['killed_mutants'] + results['survived_mutants']

        if results['total_mutations'] > 0:
            results['mutation_score'] = (
                results['killed_mutants'] / results['total_mutations']
            ) * 100

        return results
```

## 지속적 테스트 통합

```python
class ContinuousTesting:
    """DDD 워크플로우를 위한 지속적 테스트 통합."""

    def __init__(self, project_path: str):
        self.project_path = Path(project_path)
        self.test_watcher = None

    async def start_watch_mode(self):
        """파일 변경 감지 및 자동 테스트 실행 시작."""

        try:
            # pytest-watch 또는 pytest-xdist로 지속적 테스트
            result = subprocess.run(
                [
                    sys.executable, '-m', 'pytest_watch',
                    '--', str(self.project_path)
                ],
                capture_output=False,
                cwd=str(self.project_path)
            )

        except Exception as e:
            print(f"Watch mode error: {e}")

    async def run_parallel_tests(self, num_workers: int = 4) -> Dict[str, Any]:
        """빠른 피드백을 위한 병렬 테스트 실행."""

        try:
            result = subprocess.run(
                [
                    sys.executable, '-m', 'pytest',
                    '-n', str(num_workers),
                    str(self.project_path)
                ],
                capture_output=True,
                text=True,
                cwd=str(self.project_path)
            )

            return {
                'output': result.stdout,
                'success': result.returncode == 0
            }

        except Exception as e:
            return {'error': str(e), 'success': False}
```

## 모범 사례

1. 포괄적 테스트: Context7 패턴을 사용해 완전한 테스트 커버리지 확보
2. 속성 기반 테스트: 데이터 검증 함수에 속성 기반 테스트 추가
3. 변이 테스트: 변이 테스트로 테스트 스위트 품질 검증
4. 지속적 테스트: 즉각적인 피드백을 위한 감시 모드 구현
5. 성능 테스트: 주요 경로에 성능 테스트 추가
6. 보안 테스트: 인증 및 권한 부여를 위한 보안 테스트 포함
7. 통합 테스트: 컴포넌트 상호작용을 철저히 테스트
8. 테스트 문서화: 테스트 의도와 예상 동작 문서화
9. Context7 통합: 최신 테스트 패턴과 관행을 위해 Context7 활용
10. 자동화된 분석: AI 기반 테스트 제안으로 공백 식별

---

관련: [ANALYZE-PRESERVE-IMPROVE](./analyze-preserve-improve.md) | [테스트 생성](./test-generation.md) | [테스트 패턴](./test-patterns.md)
