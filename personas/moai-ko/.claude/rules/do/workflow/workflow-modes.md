# 워크플로우 모드

MoAI-ADK SPEC 워크플로우를 위한 개발 방법론 참조입니다.

단계 개요, 토큰 전략, 전환은 @spec-workflow.md를 참조하세요

## 방법론 선택

Run 단계는 `.moai/config/sections/quality.yaml`의 `quality.development_mode`를 기반으로 워크플로우를 조정합니다:

| 모드 | 워크플로우 주기 | 최적 용도 | 에이전트 전략 |
|------|---------------|----------|----------------|
| DDD | ANALYZE-PRESERVE-IMPROVE | 기존 프로젝트, < 10% 커버리지 | 특성 테스트 우선 |
| TDD | RED-GREEN-REFACTOR | 신규 프로젝트, 50%+ 커버리지 | 구현 전 테스트 |
| Hybrid | 변경 유형별 혼합 | 부분적 커버리지 (10-49%) | 신규 코드: TDD, 레거시: DDD |

## DDD 모드 (기본값)

개발 방법론: Domain-Driven Development (ANALYZE-PRESERVE-IMPROVE)

**ANALYZE**: 기존 동작과 코드 구조 이해
- 기존 코드를 읽고 의존성 식별
- 도메인 경계 and 상호작용 패턴 매핑
- 부작용 and 암묵적 계약 식별

**PRESERVE**: 기존 동작에 대한 특성 테스트 생성
- 현재 동작을 포착하는 특성 테스트 작성
- 회규 감지를 위한 동작 스냅샷 생성
- 중요 경로의 테스트 커버리지 검증

**IMPROVE**: 동작 보존으로 변경 구현
- 작은 점진적 변경 수행
- 각 변경 후 특성 테스트 실행
- 테스트 검증으로 리팩토링

성공 기준:
- 모든 SPEC 요구사항 구현됨
- 특성 테스트 통과
- 동작 스냅샷 안정적 (회귀 없음)
- 85%+ 코드 커버리지 달성
- TRUST 5 품질 게이트 통과

## TDD 모드

개발 방법론: Test-Driven Development (RED-GREEN-REFACTOR)

**RED**: 실패하는 테스트 작성
- 원하는 동작을 설명하는 테스트 작성
- 테스트가 실패하는지 확인 (새로운 것을 테스트하는지 확인)
- 한 번에 하나의 테스트, 집중적이고 구체적

**GREEN**: 통과하기 위한 최소 코드 작성
- 테스트를 통과하게 하는 가장 간단한 구현 작성
- 조기 최적화나 추상화 없음
- 정확성에 초점, 우아함 아님

**REFACTOR**: 코드 품질 개선
- 테스트를 녹색 상태로 유지하면서 구현 정리
- 패턴 추출, 중복 제거
- 적절한 곳에 SOLID 원칙 적용

성공 기준:
- 모든 SPEC 요구사항 구현됨
- 모든 테스트 통과 (RED-GREEN-REFACTOR 완료)
- 커밋당 최소 커버리지: 80% (구성 가능)
- 구현 코드 후 작성된 테스트 없음
- TRUST 5 품질 게이트 통과

## Hybrid 모드

개발 방법론: Hybrid (신규는 TDD + 레거시는 DDD)

**NEW 코드용** (새 파일, 새 함수):
- TDD 워크플로우 적용 (RED-GREEN-REFACTOR)
- 엄격한 테스트 우선 요구사항
- 커버리지 목표: 신규 코드 85%

**EXISTING 코드용** (수정, 리팩토링):
- DDD 워크플로우 적용 (ANALYZE-PRESERVE-IMPROVE)
- 변경 전 특성 테스트
- Coverage target: 85% for modified code

**Classification Logic**:
- New files - TDD rules
- Modified existing files - DDD rules
- New functions in existing files - TDD rules for those functions
- Deleted code - Verify characterization tests still pass

Success Criteria:
- All SPEC requirements implemented
- New code has TDD-level coverage (85%+)
- Modified code has characterization tests
- Overall coverage improvement trend
- TRUST 5 quality gates passed

## Methodology Selection Guide

### Auto-Detection (via /moai project or /moai init)

The system automatically recommends a methodology based on project analysis:

| Project State | Test Coverage | Recommendation | Rationale |
|--------------|---------------|----------------|-----------|
| Greenfield (new) | N/A | Hybrid | Clean slate, TDD for features + DDD structure |
| Brownfield | >= 50% | TDD | Sufficient test base for test-first development |
| Brownfield | 10-49% | Hybrid | Partial tests, expand with DDD then TDD for new |
| Brownfield | < 10% | DDD | No tests, gradual characterization test creation |

### Manual Override

Users can override the auto-detected methodology:
- During init: Select in the wizard or use `--development-mode` flag
- After init: Edit `quality.development_mode` in `.moai/config/sections/quality.yaml`
- Per session: Set `MOAI_DEVELOPMENT_MODE` environment variable

### Methodology Comparison

| Aspect | DDD | TDD | Hybrid |
|--------|-----|-----|--------|
| Test timing | After analysis (PRESERVE) | Before code (RED) | Mixed |
| Coverage approach | Gradual improvement | Strict per-commit | Unified 85% target |
| Best for | Legacy refactoring only | Isolated modules (rare) | All development work |
| Risk level | Low (preserves behavior) | Medium (requires discipline) | Medium |
| Coverage exemptions | Allowed | Not allowed | Allowed for legacy only |
| Run Phase cycle | ANALYZE-PRESERVE-IMPROVE | RED-GREEN-REFACTOR | Both (per change type) |
