# Context7와 함께하는 DDD - 핵심 클래스

> 서브 모듈: DDD 워크플로우 관리를 위한 핵심 클래스 구현
> 상위 모듈: [Context7와 함께하는 DDD](../ddd-context7.md)

## 열거형

### DDDPhase 열거형

```python
class DDDPhase(Enum):
    """DDD 사이클의 단계."""
    ANALYZE = "analyze"        # 기존 코드와 동작 분석
    PRESERVE = "preserve"      # 특성화 테스트 생성
    IMPROVE = "improve"        # 동작 보존을 유지하며 리팩토링
    REVIEW = "review"          # 검증 및 문서화
```

### TestType 열거형

```python
class TestType(Enum):
    """DDD 워크플로우에서 사용하는 테스트 유형."""
    UNIT = "unit"
    INTEGRATION = "integration"
    ACCEPTANCE = "acceptance"
    PERFORMANCE = "performance"
    SECURITY = "security"
    REGRESSION = "regression"
    CHARACTERIZATION = "characterization"
```

### TestStatus 열거형

```python
class TestStatus(Enum):
    """테스트 케이스의 상태."""
    PENDING = "pending"
    PASSED = "passed"
    FAILED = "failed"
    SKIPPED = "skipped"
    ERROR = "error"
```

## 데이터 클래스

### TestSpecification

```python
@dataclass
class TestSpecification:
    """테스트 케이스 생성을 위한 명세."""
    name: str
    description: str
    test_type: TestType
    requirements: List[str]
    acceptance_criteria: List[str]
    edge_cases: List[str]
    mock_requirements: List[str] = field(default_factory=list)
    fixture_requirements: List[str] = field(default_factory=list)
    timeout: Optional[int] = None
    tags: List[str] = field(default_factory=list)
```

### TestCase

```python
@dataclass
class TestCase:
    """메타데이터를 포함한 개별 테스트 케이스."""
    id: str
    name: str
    file_path: str
    line_number: int
    test_type: TestType
    specification: TestSpecification
    status: TestStatus
    execution_time: Optional[float] = None
    error_message: Optional[str] = None
    code: str = ""
    coverage_impact: float = 0.0
```

### DDDSession

```python
@dataclass
class DDDSession:
    """모든 사이클 활동을 추적하는 DDD 세션."""
    id: str
    project_path: str
    current_phase: DDDPhase
    test_cases: List[TestCase]
    implementation_files: List[str]
    metrics: Dict[str, Any]
    context7_patterns: Dict[str, Any]
    started_at: float
    last_activity: float
```

### DDDCycleResult

```python
@dataclass
class DDDCycleResult:
    """완전한 DDD 사이클의 결과."""
    session_id: str
    test_specification: TestSpecification
    test_file_path: str
    implementation_file_path: str
    analyze_phase_result: Dict[str, Any]
    preserve_phase_result: Dict[str, Any]
    improve_phase_result: Dict[str, Any]
    final_coverage: float
    total_time: float
    context7_patterns_applied: List[str]
    behavior_preserved: bool
```

## DDDManager 클래스

```python
class DDDManager:
    """Context7 통합을 갖춘 DDD 워크플로우 관리자."""

    def __init__(self, project_path: str, context7_client=None):
        self.project_path = project_path
        self.context7 = context7_client
        self.context7_integration = Context7DDDIntegration(context7_client)
        self.test_generator = TestGenerator(context7_client)
        self.current_session: Optional[DDDSession] = None

    async def start_ddd_session(self, feature_name: str) -> DDDSession:
        """새로운 DDD 세션 시작."""
        session_id = f"ddd_{feature_name}_{int(time.time())}"

        # Context7 패턴 로드
        patterns = await self.context7_integration.load_ddd_patterns()

        session = DDDSession(
            id=session_id,
            project_path=self.project_path,
            current_phase=DDDPhase.ANALYZE,
            test_cases=[],
            implementation_files=[],
            metrics={'analyze_phases': 0, 'preserve_phases': 0, 'improve_phases': 0},
            context7_patterns=patterns,
            started_at=time.time(),
            last_activity=time.time()
        )

        self.current_session = session
        return session

    async def run_full_ddd_cycle(
        self,
        specification: TestSpecification,
        target_function: str
    ) -> DDDCycleResult:
        """완전한 ANALYZE-PRESERVE-IMPROVE 사이클 실행."""
        if not self.current_session:
            self.current_session = await self.start_ddd_session("default")

        cycle_start = time.time()
        context7_patterns_applied = []

        # ANALYZE 단계 - 기존 코드와 동작 파악
        analyze_result = await self._execute_analyze_phase(specification)
        self.current_session.metrics['analyze_phases'] += 1

        # PRESERVE 단계 - 특성화 테스트 생성
        preserve_result = await self._execute_preserve_phase(
            specification, target_function, analyze_result
        )
        self.current_session.metrics['preserve_phases'] += 1

        # IMPROVE 단계 - 동작 보존을 유지하며 리팩토링
        improve_result = await self._execute_improve_phase(
            specification, preserve_result
        )
        self.current_session.metrics['improve_phases'] += 1
        context7_patterns_applied.extend(improve_result.get('patterns_applied', []))

        # 최종 커버리지 실행
        coverage = await self._run_coverage_analysis()

        return DDDCycleResult(
            session_id=self.current_session.id,
            test_specification=specification,
            test_file_path=preserve_result.get('test_file_path', ''),
            implementation_file_path=improve_result.get('implementation_file_path', ''),
            analyze_phase_result=analyze_result,
            preserve_phase_result=preserve_result,
            improve_phase_result=improve_result,
            final_coverage=coverage.get('total_coverage', 0.0),
            total_time=time.time() - cycle_start,
            context7_patterns_applied=context7_patterns_applied,
            behavior_preserved=improve_result.get('behavior_preserved', True)
        )
```

## 단계 실행 메서드

### ANALYZE 단계

```python
async def _execute_analyze_phase(
    self, specification: TestSpecification
) -> Dict[str, Any]:
    """ANALYZE 단계 실행 - 기존 코드와 동작 파악."""
    self.current_session.current_phase = DDDPhase.ANALYZE

    # 기존 코드 구조 분석
    code_analysis = await self._analyze_existing_code(specification)

    # 동작 패턴 식별
    behavior_patterns = await self._identify_behavior_patterns(code_analysis)

    # 리팩토링 대상 결정
    refactoring_targets = await self._identify_refactoring_targets(code_analysis)

    return {
        'code_analysis': code_analysis,
        'behavior_patterns': behavior_patterns,
        'refactoring_targets': refactoring_targets,
        'phase_success': True
    }
```

### PRESERVE 단계

```python
async def _execute_preserve_phase(
    self, specification: TestSpecification,
    target_function: str,
    analyze_result: Dict[str, Any]
) -> Dict[str, Any]:
    """PRESERVE 단계 실행 - 특성화 테스트 생성."""
    self.current_session.current_phase = DDDPhase.PRESERVE

    # 기존 동작에 대한 특성화 테스트 생성
    test_code = await self.test_generator.generate_characterization_test(
        specification, analyze_result['behavior_patterns']
    )

    # 테스트 파일 경로 결정
    test_file_path = self._get_test_file_path(specification)

    # 파일에 테스트 작성
    self._write_test_file(test_file_path, test_code)

    # 테스트 실행 - 통과해야 함 (기존 동작 테스트)
    test_result = await self._run_tests(test_file_path)

    # 테스트 케이스 레코드 생성
    test_case = TestCase(
        id=f"tc_{specification.name}",
        name=specification.name,
        file_path=test_file_path,
        line_number=1,
        test_type=TestType.CHARACTERIZATION,
        specification=specification,
        status=TestStatus.PASSED if test_result['failed'] == 0 else TestStatus.FAILED,
        execution_time=test_result.get('execution_time', 0),
        code=test_code
    )

    self.current_session.test_cases.append(test_case)

    return {
        'test_code': test_code,
        'test_file_path': test_file_path,
        'test_result': test_result,
        'test_case': test_case,
        'phase_success': test_result['failed'] == 0  # PRESERVE 단계에서는 통과해야 함
    }
```

### IMPROVE 단계

```python
async def _execute_improve_phase(
    self, specification: TestSpecification,
    preserve_result: Dict[str, Any]
) -> Dict[str, Any]:
    """IMPROVE 단계 실행 - 동작 보존을 유지하며 리팩토링."""
    self.current_session.current_phase = DDDPhase.IMPROVE

    # Context7에서 개선 패턴 가져오기
    improve_patterns = await self.context7_integration.get_improvement_patterns()

    # 개선사항 생성
    improvements = await self._generate_improvements(
        preserve_result.get('implementation', ''),
        improve_patterns
    )

    patterns_applied = []
    successful_improvements = []
    behavior_preserved = True

    for improvement in improvements:
        # 개선사항 적용
        improved = await self._apply_improvement(
            preserve_result.get('implementation_file_path', ''),
            improvement
        )

        if improved['success']:
            # 특성화 테스트를 실행해 동작 보존 검증
            test_result = await self._run_tests(preserve_result['test_file_path'])

            if test_result['failed'] == 0:
                successful_improvements.append(improvement)
                patterns_applied.append(improvement.get('pattern', 'custom'))
            else:
                # 실패한 개선사항 롤백 - 동작이 보존되지 않음
                await self._rollback_improvement(
                    preserve_result.get('implementation_file_path', '')
                )
                behavior_preserved = False

    return {
        'improvements_suggested': len(improvements),
        'improvements_applied': len(successful_improvements),
        'patterns_applied': patterns_applied,
        'behavior_preserved': behavior_preserved,
        'phase_success': behavior_preserved
    }
```

## 헬퍼 메서드

```python
def _get_test_file_path(self, specification: TestSpecification) -> str:
    """명세를 기반으로 테스트 파일 경로 결정."""
    test_dir = os.path.join(self.project_path, 'tests')
    os.makedirs(test_dir, exist_ok=True)

    test_type_dir = specification.test_type.value
    full_test_dir = os.path.join(test_dir, test_type_dir)
    os.makedirs(full_test_dir, exist_ok=True)

    return os.path.join(full_test_dir, f"test_{specification.name}.py")

def _get_implementation_file_path(self, target_function: str) -> str:
    """구현 파일 경로 결정."""
    src_dir = os.path.join(self.project_path, 'src')
    os.makedirs(src_dir, exist_ok=True)
    return os.path.join(src_dir, f"{target_function}.py")

async def _run_tests(self, test_path: str) -> Dict[str, Any]:
    """지정된 경로에서 pytest 실행."""
    result = subprocess.run(
        ['pytest', test_path, '-v', '--tb=short', '--json-report'],
        capture_output=True,
        text=True,
        cwd=self.project_path
    )

    return {
        'passed': result.stdout.count('PASSED'),
        'failed': result.stdout.count('FAILED'),
        'errors': result.stdout.count('ERROR'),
        'execution_time': 0.0,  # 출력에서 파싱
        'output': result.stdout
    }

async def _run_coverage_analysis(self) -> Dict[str, Any]:
    """커버리지 분석 실행."""
    result = subprocess.run(
        ['pytest', '--cov=src', '--cov-report=json'],
        capture_output=True,
        text=True,
        cwd=self.project_path
    )

    try:
        coverage_file = os.path.join(self.project_path, 'coverage.json')
        with open(coverage_file) as f:
            coverage_data = json.load(f)
            return {'total_coverage': coverage_data.get('totals', {}).get('percent_covered', 0)}
    except Exception:
        return {'total_coverage': 0.0}
```

## 관련 서브 모듈

- [테스트 생성](./test-generation.md) - AI 기반 테스트 생성
- [Context7 패턴](./context7-patterns.md) - 패턴 통합

---

서브 모듈: `modules/ddd/core-classes.md`
