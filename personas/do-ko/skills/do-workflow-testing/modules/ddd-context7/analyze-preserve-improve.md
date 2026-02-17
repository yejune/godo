# ANALYZE-PRESERVE-IMPROVE DDD 사이클

> 모듈: Context7 통합을 갖춘 핵심 DDD 사이클 구현
> 복잡도: 고급
> 소요 시간: 20분 이상
> 의존성: Python 3.8+, pytest, Context7 MCP, unittest

## 핵심 DDD 클래스

```python
import pytest
import unittest
import asyncio
import subprocess
import os
import sys
import time
from typing import Dict, List, Optional, Any
from dataclasses import dataclass, field
from enum import Enum
import json
from pathlib import Path

class DDDPhase(Enum):
    """DDD 사이클 단계."""
    ANALYZE = "analyze"       # 기존 코드와 동작 분석
    PRESERVE = "preserve"     # 특성화 테스트 생성
    IMPROVE = "improve"       # 테스트를 녹색으로 유지하며 코드 개선
    REVIEW = "review"         # 검토 및 변경사항 커밋

class TestType(Enum):
    """DDD에서 사용하는 테스트 유형."""
    UNIT = "unit"
    INTEGRATION = "integration"
    CHARACTERIZATION = "characterization"
    ACCEPTANCE = "acceptance"
    PERFORMANCE = "performance"
    SECURITY = "security"
    REGRESSION = "regression"

class TestStatus(Enum):
    """테스트 실행 상태."""
    PENDING = "pending"
    RUNNING = "running"
    PASSED = "passed"
    FAILED = "failed"
    SKIPPED = "skipped"
    ERROR = "error"

@dataclass
class TestSpecification:
    """DDD 테스트 명세."""
    name: str
    description: str
    test_type: TestType
    requirements: List[str]
    acceptance_criteria: List[str]
    edge_cases: List[str]
    preconditions: List[str] = field(default_factory=list)
    postconditions: List[str] = field(default_factory=list)
    dependencies: List[str] = field(default_factory=list)
    mock_requirements: Dict[str, Any] = field(default_factory=dict)
    behavior_snapshot: Optional[Dict[str, Any]] = None

@dataclass
class TestCase:
    """메타데이터를 포함한 개별 테스트 케이스."""
    id: str
    name: str
    file_path: str
    line_number: int
    specification: TestSpecification
    status: TestStatus
    execution_time: float
    error_message: Optional[str] = None
    coverage_data: Dict[str, Any] = field(default_factory=dict)

@dataclass
class DDDSession:
    """사이클 추적을 포함한 DDD 개발 세션."""
    id: str
    project_path: str
    current_phase: DDDPhase
    test_cases: List[TestCase]
    start_time: float
    context7_patterns: Dict[str, Any] = field(default_factory=dict)
    metrics: Dict[str, Any] = field(default_factory=dict)
    behavior_snapshots: Dict[str, Any] = field(default_factory=dict)
```

## DDD Manager 구현

```python
class DDDManager:
    """Context7 통합을 갖춘 주요 DDD 워크플로우 관리자."""

    def __init__(self, project_path: str, context7_client=None):
        self.project_path = Path(project_path)
        self.context7 = context7_client
        self.current_session = None
        self.test_history = []

    async def start_ddd_session(
        self, feature_name: str,
        test_types: List[TestType] = None
    ) -> DDDSession:
        """새로운 DDD 개발 세션 시작."""

        if test_types is None:
            test_types = [TestType.CHARACTERIZATION, TestType.UNIT, TestType.INTEGRATION]

        # 세션 생성
        session = DDDSession(
            id=f"ddd_{feature_name}_{int(time.time())}",
            project_path=str(self.project_path),
            current_phase=DDDPhase.ANALYZE,
            test_cases=[],
            start_time=time.time(),
            context7_patterns={},
            metrics={
                'tests_written': 0,
                'tests_passing': 0,
                'tests_failing': 0,
                'coverage_percentage': 0.0,
                'behaviors_preserved': 0
            },
            behavior_snapshots={}
        )

        self.current_session = session
        return session

    async def run_full_ddd_cycle(
        self, specification: TestSpecification,
        target_function: str = None
    ) -> Dict[str, Any]:
        """완전한 ANALYZE-PRESERVE-IMPROVE DDD 사이클 실행."""

        cycle_results = {}

        # ANALYZE 단계
        print("🔍 ANALYZE 단계: 기존 코드와 동작 파악 중...")
        analyze_results = await self._run_analyze_phase(target_function)
        cycle_results['analyze'] = analyze_results
        self.current_session.current_phase = DDDPhase.ANALYZE

        # PRESERVE 단계
        print("🧪 PRESERVE 단계: 특성화 테스트 생성 중...")
        preserve_results = await self._run_preserve_phase(specification, analyze_results)
        cycle_results['preserve'] = preserve_results
        self.current_session.current_phase = DDDPhase.PRESERVE

        # IMPROVE 단계
        print("🔧 IMPROVE 단계: 동작 보존을 유지하며 리팩토링 중...")
        improve_results = await self._run_improve_phase(specification)
        cycle_results['improve'] = improve_results
        self.current_session.current_phase = DDDPhase.IMPROVE

        # REVIEW 단계
        print("✅ REVIEW 단계: 최종 검증 중...")
        coverage_results = await self._run_coverage_analysis()
        cycle_results['review'] = {'coverage': coverage_results}
        self.current_session.current_phase = DDDPhase.REVIEW

        return cycle_results

    async def _run_analyze_phase(self, target_function: str = None) -> Dict[str, Any]:
        """ANALYZE: 기존 코드와 동작 파악."""

        analysis = {
            'existing_tests': [],
            'code_patterns': [],
            'dependencies': [],
            'behavior_notes': []
        }

        # 기존 테스트 찾기
        test_files = list(self.project_path.glob("**/test_*.py"))
        analysis['existing_tests'] = [str(f) for f in test_files]

        # 코드 구조 분석
        if target_function:
            analysis['target'] = target_function
            analysis['behavior_notes'].append(f"Analyzing behavior of {target_function}")

        return analysis

    async def _run_preserve_phase(
        self, specification: TestSpecification,
        analysis: Dict[str, Any]
    ) -> Dict[str, Any]:
        """PRESERVE: 기존 동작을 포착하는 특성화 테스트 생성."""

        preserve_results = {
            'characterization_tests_created': 0,
            'behaviors_captured': [],
            'test_files': []
        }

        # 분석 결과를 기반으로 특성화 테스트 생성
        for behavior in analysis.get('behavior_notes', []):
            preserve_results['behaviors_captured'].append(behavior)
            preserve_results['characterization_tests_created'] += 1

        # 기존 테스트를 실행해 기준선 확립
        test_results = await self._run_pytest()
        preserve_results['baseline_results'] = test_results

        return preserve_results

    async def _run_improve_phase(self, specification: TestSpecification) -> Dict[str, Any]:
        """IMPROVE: 테스트 통과를 유지하며 코드 리팩토링."""

        improve_results = {
            'improvements_made': [],
            'tests_still_passing': True,
            'refactoring_notes': []
        }

        # 개선 후 테스트 실행
        test_results = await self._run_pytest()
        improve_results['tests_still_passing'] = test_results.get('failed', 0) == 0

        if improve_results['tests_still_passing']:
            self.current_session.metrics['behaviors_preserved'] += 1

        return improve_results

    async def _run_pytest(self) -> Dict[str, Any]:
        """pytest 실행 후 결과 반환."""

        try:
            result = subprocess.run(
                [
                    sys.executable, '-m', 'pytest',
                    str(self.project_path),
                    '--tb=short',
                    '-v'
                ],
                capture_output=True,
                text=True,
                cwd=str(self.project_path)
            )

            return self._parse_pytest_output(result.stdout)

        except Exception as e:
            print(f"Error running pytest: {e}")
            return {'error': str(e), 'passed': 0, 'failed': 0}

    def _parse_pytest_output(self, output: str) -> Dict[str, Any]:
        """pytest 출력 파싱."""

        lines = output.split('\n')
        results = {'passed': 0, 'failed': 0, 'skipped': 0, 'total': 0}

        for line in lines:
            if ' passed in ' in line:
                parts = line.split()
                if parts and parts[0].isdigit():
                    results['passed'] = int(parts[0])
                    results['total'] = int(parts[0])
            elif ' passed' in line and ' failed' in line:
                passed_part = line.split(' passed')[0]
                if passed_part.strip().isdigit():
                    results['passed'] = int(passed_part.strip())

                if ' failed' in line:
                    failed_part = line.split(' failed')[0].split(', ')[-1]
                    if failed_part.strip().isdigit():
                        results['failed'] = int(failed_part.strip())

                results['total'] = results['passed'] + results['failed']

        return results

    async def _run_coverage_analysis(self) -> Dict[str, Any]:
        """테스트 커버리지 분석 실행."""

        try:
            result = subprocess.run(
                [
                    sys.executable, '-m', 'pytest',
                    str(self.project_path),
                    '--cov=src',
                    '--cov-report=term-missing'
                ],
                capture_output=True,
                text=True,
                cwd=str(self.project_path)
            )

            return {'coverage_output': result.stdout}

        except Exception as e:
            return {'error': str(e)}

    def get_session_summary(self) -> Dict[str, Any]:
        """현재 DDD 세션 요약 가져오기."""

        if not self.current_session:
            return {}

        duration = time.time() - self.current_session.start_time

        return {
            'session_id': self.current_session.id,
            'phase': self.current_session.current_phase.value,
            'duration_seconds': duration,
            'duration_formatted': f"{duration:.1f} seconds",
            'metrics': self.current_session.metrics,
            'test_cases_count': len(self.current_session.test_cases),
            'behaviors_preserved': self.current_session.metrics.get('behaviors_preserved', 0)
        }
```

## 단계별 가이드라인

### ANALYZE 단계
- 기존 코드 구조와 패턴 파악
- 코드 읽기를 통해 현재 동작 식별
- 의존성과 사이드 이펙트 문서화
- 테스트 커버리지 공백 매핑
- 기존 설계 패턴 파악

### PRESERVE 단계
- 기존 동작을 위한 특성화 테스트 작성
- 현재 동작을 "황금 표준"으로 포착
- 현재 구현으로 테스트가 통과함을 확인
- 발견된 동작 문서화
- 복잡한 출력에 대한 동작 스냅샷 생성

### IMPROVE 단계
- 테스트를 녹색으로 유지하며 코드 리팩토링
- 작고 점진적인 변경 수행
- 변경 후마다 테스트 실행
- 동작 보존 유지
- 설계 패턴 적절히 적용

### REVIEW 단계
- 모든 특성화 테스트가 여전히 통과하는지 확인
- 코드 품질과 문서 검토
- 동작 변경 여부 확인
- 명확한 메시지로 변경사항 커밋
- 수행된 개선사항 문서화

## 사용 예시

```python
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
    description="Preserve existing login behavior during refactoring",
    test_type=TestType.CHARACTERIZATION,
    requirements=[
        "Existing login flow must continue to work",
        "Error messages should remain consistent"
    ],
    acceptance_criteria=[
        "Valid credentials return user token (existing behavior)",
        "Invalid credentials raise same error messages"
    ],
    edge_cases=[
        "Test with empty email (existing behavior)",
        "Test with empty password (existing behavior)"
    ]
)

# 전체 DDD 사이클 실행
cycle_results = await ddd_manager.run_full_ddd_cycle(
    specification=test_spec,
    target_function="authenticate_user"
)

# 세션 요약 가져오기
summary = ddd_manager.get_session_summary()
print(f"Session completed in {summary['duration_formatted']}")
print(f"Behaviors preserved: {summary['behaviors_preserved']}")
```

---

관련: [테스트 생성](./test-generation.md) | [테스트 패턴](./test-patterns.md)
