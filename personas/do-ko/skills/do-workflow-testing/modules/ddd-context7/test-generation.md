# AI 기반 테스트 생성

> 모듈: Context7 강화 테스트 케이스 생성 및 명세
> 복잡도: 고급
> 소요 시간: 15분 이상
> 의존성: Python 3.8+, pytest, Context7 MCP, AST 분석

## 테스트 생성기 클래스

```python
class TestGenerator:
    """명세 기반 AI 테스트 케이스 생성기."""

    def __init__(self, context7_client=None):
        self.context7 = context7_client
        self.templates = self._load_test_templates()

    def _load_test_templates(self) -> Dict[str, str]:
        """다양한 시나리오에 대한 테스트 템플릿 로드."""
        return {
            'unit_function': '''
def test_{function_name}_{scenario}():
    """
    Test {description}

    Given: {preconditions}
    When: {action}
    Then: {expected_outcome}
    """
    # Arrange
    {setup_code}

    # Act
    result = {function_call}

    # Assert
    assert result == {expected_value}, f"Expected {expected_value}, got {result}"
''',
            'exception_test': '''
def test_{function_name}_raises_{exception}_{scenario}():
    """
    Test that {function_name} raises {exception} when {condition}
    """
    # Arrange
    {setup_code}

    # Act & Assert
    with pytest.raises({exception}) as exc_info:
        {function_call}

    assert "{expected_message}" in str(exc_info.value)
''',
            'parameterized_test': '''
@pytest.mark.parametrize("{param_names}", {test_values})
def test_{function_name}_{scenario}({param_names}):
    """
    Test {function_name} with different inputs: {description}
    """
    # Arrange
    {setup_code}

    # Act
    result = {function_call}

    # Assert
    assert result == {expected_value}, f"For {param_names}={{param_names}}, expected {expected_value}, got {{result}}"
'''
        }

    async def generate_test_case(
        self, specification: TestSpecification,
        context7_patterns: Dict[str, Any] = None
    ) -> str:
        """명세를 기반으로 테스트 코드 생성."""

        if self.context7 and context7_patterns:
            try:
                # Context7으로 테스트 생성 강화
                enhanced_spec = await self._enhance_specification_with_context7(
                    specification, context7_patterns
                )
                return self._generate_test_from_enhanced_spec(enhanced_spec)
            except Exception as e:
                print(f"Context7 test generation failed: {e}")

        return self._generate_test_from_specification(specification)

    async def _enhance_specification_with_context7(
        self, specification: TestSpecification,
        context7_patterns: Dict[str, Any]
    ) -> TestSpecification:
        """Context7 패턴으로 테스트 명세 강화."""

        # Context7 패턴을 기반으로 추가 엣지 케이스 추가
        additional_edge_cases = []

        testing_patterns = context7_patterns.get('python_testing', {})
        if testing_patterns:
            # 다양한 데이터 유형에 대한 공통 엣지 케이스 추가
            if any('number' in str(req).lower() for req in specification.requirements):
                additional_edge_cases.extend([
                    "Test with zero value",
                    "Test with negative value",
                    "Test with maximum/minimum values",
                    "Test with floating point edge cases"
                ])

            if any('string' in str(req).lower() for req in specification.requirements):
                additional_edge_cases.extend([
                    "Test with empty string",
                    "Test with very long string",
                    "Test with special characters",
                    "Test with unicode characters"
                ])

            if any('list' in str(req).lower() or 'array' in str(req).lower()
                   for req in specification.requirements):
                additional_edge_cases.extend([
                    "Test with empty list",
                    "Test with single element",
                    "Test with large list",
                    "Test with duplicate elements"
                ])

        # 원본과 추가 엣지 케이스 합치기
        combined_edge_cases = list(set(specification.edge_cases + additional_edge_cases))

        return TestSpecification(
            name=specification.name,
            description=specification.description,
            test_type=specification.test_type,
            requirements=specification.requirements,
            acceptance_criteria=specification.acceptance_criteria,
            edge_cases=combined_edge_cases,
            preconditions=specification.preconditions,
            postconditions=specification.postconditions,
            dependencies=specification.dependencies,
            mock_requirements=specification.mock_requirements
        )

    def _generate_test_from_enhanced_spec(self, spec: TestSpecification) -> str:
        """강화된 명세에서 테스트 코드 생성."""
        return self._generate_test_from_specification(spec)

    def _generate_test_from_specification(self, spec: TestSpecification) -> str:
        """명세에서 테스트 코드 생성."""

        # 테스트 유형에 따라 적절한 템플릿 선택
        if spec.test_type == TestType.UNIT:
            return self._generate_unit_test(spec)
        elif spec.test_type == TestType.INTEGRATION:
            return self._generate_integration_test(spec)
        else:
            return self._generate_generic_test(spec)

    def _generate_unit_test(self, spec: TestSpecification) -> str:
        """단위 테스트 코드 생성."""

        function_name = spec.name.lower().replace('test_', '').replace('_test', '')

        # 예외 테스트인지 확인
        if any('error' in criterion.lower() or 'exception' in criterion.lower()
               for criterion in spec.acceptance_criteria):
            return self._generate_exception_test(spec, function_name)

        # 매개변수화가 필요한지 확인
        if len(spec.acceptance_criteria) > 1 or len(spec.edge_cases) > 2:
            return self._generate_parameterized_test(spec, function_name)

        # 표준 단위 테스트 생성
        return self._generate_standard_unit_test(spec, function_name)

    def _generate_standard_unit_test(self, spec: TestSpecification, function_name: str) -> str:
        """표준 단위 테스트 생성."""

        template = self.templates['unit_function']

        setup_code = self._generate_setup_code(spec)
        function_call = self._generate_function_call(function_name, spec)
        assertions = self._generate_assertions(spec)

        return template.format(
            function_name=function_name,
            scenario=self._extract_scenario(spec),
            description=spec.description,
            preconditions=', '.join(spec.preconditions),
            action=self._describe_action(spec),
            expected_outcome=spec.acceptance_criteria[0] if spec.acceptance_criteria else "expected behavior",
            setup_code=setup_code,
            function_call=function_call,
            expected_value=self._extract_expected_value(spec),
            assertions=assertions
        )

    def _generate_exception_test(self, spec: TestSpecification, function_name: str) -> str:
        """예외 테스트 생성."""

        template = self.templates['exception_test']

        # 예상 예외와 메시지 추출
        exception_type = "Exception" # 기본값
        expected_message = "Error occurred"

        for criterion in spec.acceptance_criteria:
            if 'raise' in criterion.lower() or 'exception' in criterion.lower():
                # 예외 유형 추출 시도
                if 'valueerror' in criterion.lower():
                    exception_type = "ValueError"
                elif 'typeerror' in criterion.lower():
                    exception_type = "TypeError"
                elif 'attributeerror' in criterion.lower():
                    exception_type = "AttributeError"
                elif 'keyerror' in criterion.lower():
                    exception_type = "KeyError"

                # 예상 메시지 추출 시도
                if 'message:' in criterion.lower():
                    parts = criterion.split('message:')
                    if len(parts) > 1:
                        expected_message = parts[1].strip().strip('"\'')
                        break

        return template.format(
            function_name=function_name,
            exception=exception_type,
            scenario=self._extract_scenario(spec),
            condition=self._describe_condition(spec),
            setup_code=self._generate_setup_code(spec),
            function_call=self._generate_function_call(function_name, spec),
            expected_message=expected_message
        )

    def _generate_parameterized_test(self, spec: TestSpecification, function_name: str) -> str:
        """매개변수화 테스트 생성."""

        template = self.templates['parameterized_test']

        # 테스트 매개변수와 값 생성
        param_names, test_values = self._generate_test_parameters(spec)

        return template.format(
            function_name=function_name,
            scenario=self._extract_scenario(spec),
            description=spec.description,
            param_names=', '.join(param_names),
            test_values=test_values,
            setup_code=self._generate_setup_code(spec),
            function_call=self._generate_function_call(function_name, spec),
            expected_value=self._extract_expected_value(spec)
        )

    def _generate_integration_test(self, spec: TestSpecification) -> str:
        """통합 테스트 생성."""
        return f'''
def test_{spec.name.replace(' ', '_').lower()}():
    """Test: {spec.description}"""
    # Arrange: {', '.join(spec.preconditions[:2])}
    # Act: Call integration function
    # Assert: Verify integration behavior
'''

    def _generate_generic_test(self, spec: TestSpecification) -> str:
        """범용 테스트 코드 생성."""
        return f'''
def test_{spec.name.replace(' ', '_').lower()}():
    """Test: {spec.description}"""
    # TODO: Implement based on specification
    # Requirements: {len(spec.requirements)} items
    # Acceptance Criteria: {len(spec.acceptance_criteria)} items
'''

    def _extract_scenario(self, spec: TestSpecification) -> str:
        """명세에서 시나리오 이름 추출."""
        if '_' in spec.name:
            parts = spec.name.split('_')
            if len(parts) > 1:
                return '_'.join(parts[1:])
        return 'default'

    def _describe_action(self, spec: TestSpecification) -> str:
        """테스트 중인 액션 설명."""
        return f"Call {spec.name}"

    def _describe_condition(self, spec: TestSpecification) -> str:
        """예외 테스트의 조건 설명."""
        return spec.requirements[0] if spec.requirements else "invalid input"

    def _generate_setup_code(self, spec: TestSpecification) -> str:
        """명세를 기반으로 설정 코드 생성."""
        setup_lines = []

        # mock 요구사항 추가
        for mock_name, mock_config in spec.mock_requirements.items():
            if isinstance(mock_config, dict) and 'return_value' in mock_config:
                setup_lines.append(f"{mock_name} = Mock(return_value={mock_config['return_value']})")
            else:
                setup_lines.append(f"{mock_name} = Mock()")

        # 사전 조건을 설정 코드로 추가
        for condition in spec.preconditions:
            setup_lines.append(f"# {condition}")

        return '\n '.join(setup_lines) if setup_lines else "pass"

    def _generate_function_call(self, function_name: str, spec: TestSpecification) -> str:
        """인수를 포함한 함수 호출 생성."""

        # mock 요구사항 또는 요구사항에서 인수 추출
        args = []

        if spec.mock_requirements:
            args.extend(spec.mock_requirements.keys())

        if not args:
            # 요구사항을 기반으로 플레이스홀더 인수 추가
            for req in spec.requirements[:3]: # 첫 3개 요구사항으로 제한
                if 'input' in req.lower() or 'parameter' in req.lower():
                    args.append("test_input")
                    break

        return f"{function_name}({', '.join(args)})" if args else f"{function_name}()"

    def _generate_assertions(self, spec: TestSpecification) -> str:
        """인수 기준을 기반으로 단언문 생성."""
        assertions = []

        for criterion in spec.acceptance_criteria[:3]: # 첫 3개 기준으로 제한
            if 'returns' in criterion.lower() or 'result' in criterion.lower():
                assertions.append("assert result is not None")
            elif 'equals' in criterion.lower() or 'equal' in criterion.lower():
                assertions.append("assert result == expected_value")
            elif 'length' in criterion.lower():
                assertions.append("assert len(result) > 0")
            else:
                assertions.append(f"# {criterion}")

        return '\n '.join(assertions) if assertions else "assert True # Add specific assertions"

    def _extract_expected_value(self, spec: TestSpecification) -> str:
        """인수 기준에서 예상 값 추출."""
        for criterion in spec.acceptance_criteria:
            if 'returns' in criterion.lower():
                # 예상 값 추출 시도
                if 'true' in criterion.lower():
                    return "True"
                elif 'false' in criterion.lower():
                    return "False"
                elif 'none' in criterion.lower():
                    return "None"
                elif 'empty' in criterion.lower():
                    return "[]"
                else:
                    return "expected_result"
        return "expected_result"

    def _generate_test_parameters(self, spec: TestSpecification) -> tuple:
        """매개변수화 테스트를 위한 매개변수와 값 생성."""

        # 인수 기준과 엣지 케이스에서 테스트 케이스 생성
        test_cases = []

        # 인수 기준을 테스트 케이스로 추가
        for criterion in spec.acceptance_criteria:
            if 'input' in criterion.lower():
                # 입력 값 추출
                if 'valid' in criterion.lower():
                    test_cases.append(('valid_input', 'expected_output'))
                elif 'invalid' in criterion.lower():
                    test_cases.append(('invalid_input', 'exception'))

        # 엣지 케이스 추가
        for edge_case in spec.edge_cases:
            if 'zero' in edge_case.lower():
                test_cases.append((0, 'zero_result'))
            elif 'empty' in edge_case.lower():
                test_cases.append(('', 'empty_result'))
            elif 'null' in edge_case.lower() or 'none' in edge_case.lower():
                test_cases.append((None, 'none_result'))

        # pytest 형식으로 변환
        if test_cases:
            param_names = ['test_input', 'expected_output']
            test_values = str(test_cases).replace("'", '"')
            return param_names, test_values

        # 폴백
        return ['test_input', 'expected_output'], '[("test", "expected")]'
```

## Context7 통합

```python
class Context7TestIntegration:
    """테스트 생성 패턴을 위한 Context7 통합."""

    def __init__(self, context7_client=None):
        self.context7 = context7_client
        self.pattern_cache = {}

    async def load_test_generation_patterns(
        self, language: str = "python"
    ) -> Dict[str, Any]:
        """Context7에서 테스트 생성 패턴 로드."""

        cache_key = f"test_gen_patterns_{language}"
        if cache_key in self.pattern_cache:
            return self.pattern_cache[cache_key]

        patterns = {}

        if self.context7:
            try:
                # 테스트 생성 패턴 로드
                gen_patterns = await self.context7.get_library_docs(
                    context7_library_id="/testing/pytest",
                    topic="test generation patterns automation 2025",
                    tokens=3000
                )
                patterns['generation'] = gen_patterns

                # 엣지 케이스 패턴 로드
                edge_patterns = await self.context7.get_library_docs(
                    context7_library_id="/testing/edge-cases",
                    topic="edge case generation boundary testing 2025",
                    tokens=2000
                )
                patterns['edge_cases'] = edge_patterns

            except Exception as e:
                print(f"Failed to load Context7 patterns: {e}")
                patterns = self._get_default_patterns()
        else:
            patterns = self._get_default_patterns()

        self.pattern_cache[cache_key] = patterns
        return patterns

    def _get_default_patterns(self) -> Dict[str, Any]:
        """기본 테스트 생성 패턴 반환."""
        return {
            'generation': {
                'strategies': [
                    "명세에서 테스트 생성",
                    "코드 분석으로 누락된 테스트 케이스 식별",
                    "여러 시나리오에 대한 매개변수화 테스트 생성",
                    "에러 조건에 대한 예외 테스트 생성"
                ]
            },
            'edge_cases': {
                'categories': [
                    "경계 값 (최소, 최대, 바로 위/아래)",
                    "빈/null 입력",
                    "유효하지 않은 데이터 유형",
                    "특수 문자와 유니코드",
                    "대용량 입력 (성능 테스트)"
                ]
            }
        }
```

## 모범 사례

1. 명세 주도: 항상 명확한 명세에서 테스트 생성
2. 엣지 케이스 커버리지: Context7 패턴으로 포괄적인 엣지 케이스 테스트 확보
3. 읽기 쉬운 테스트: 의도를 명확히 표현하는 테스트 생성
4. 유지 보수 가능: 생성된 테스트를 단순하고 집중적으로 유지
5. 컨텍스트 인식: 언어별, 프레임워크별 패턴을 위해 Context7 활용

---

관련: [ANALYZE-PRESERVE-IMPROVE](./analyze-preserve-improve.md) | [테스트 패턴](./test-patterns.md)
