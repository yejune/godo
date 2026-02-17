---
name: moai-workflow-run
description: >
  SPEC 요구사항에 대한 DDD/TDD/Hybrid 구현 워크플로우입니다. Plan-Run-Sync 워크플로우의
  두 번째 단계입니다. quality.yaml의 development_mode 설정에 따라 manager-ddd 또는
  manager-tdd로 라우팅합니다.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "2.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-07"
  tags: "run, implementation, ddd, tdd, hybrid, spec"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["run", "implement", "build", "create", "develop", "code"]
  agents: ["manager-ddd", "manager-tdd", "manager-strategy", "manager-quality", "manager-git"]
  phases: ["run"]
---

# Run 워크플로우 오케스트레이션

## 목적

설정된 개발 방법론을 사용하여 SPEC 요구사항을 구현합니다. 방법론은 quality.yaml의 `development_mode`에 의해 결정됩니다:

- **ddd**: ANALYZE-PRESERVE-IMPROVE 사이클을 사용하는 도메인 주도 개발 (레거시 리팩토링 전용)
- **tdd**: RED-GREEN-REFACTOR 사이클을 사용하는 테스트 주도 개발 (격리된 신규 모듈용)
- **hybrid**: 새 코드에는 TDD, 기존 코드 수정에는 DDD를 사용하는 통합 접근 방식 (모든 개발에 권장)

이것은 Plan-Run-Sync 워크플로우의 두 번째 단계입니다.

## 범위

- MoAI 4단계 워크플로우의 3단계 구현 (작업 실행)
- /moai plan으로 생성된 SPEC 문서를 입력으로 받음
- 문서화 및 PR을 위해 /moai sync로 인계

## 입력

- $ARGUMENTS: 구현할 SPEC-ID (예: SPEC-AUTH-001)
- 재개: /moai run SPEC-XXX 재실행 시 마지막 성공한 단계 체크포인트에서 재개
- --team: 팀 기반 구현 활성화 (병렬 구현 팀은 team-run.md 참조)

## 컨텍스트 로딩

실행 전 다음 필수 파일을 로드하세요:

- .moai/config/config.yaml (git 전략, 자동화 설정)
- .moai/config/sections/quality.yaml (커버리지 목표, TRUST 5 설정)
- .moai/config/sections/git-strategy.yaml (auto_branch, 브랜치 생성 정책)
- .moai/config/sections/language.yaml (git_commit_messages 설정)
- .moai/specs/SPEC-{ID}/ 디렉토리 (spec.md, plan.md, acceptance.md)

실행 전 명령어: git status, git branch, git log, git diff.

---

## 단계 순서

모든 단계는 순차적으로 실행됩니다. 각 단계는 이전 모든 단계의 출력을 컨텍스트로 받습니다. DDD 방법론이 특정 순서를 요구하므로 병렬 실행은 허용되지 않습니다.

### Phase 1: 분석 및 계획

에이전트: manager-strategy 서브에이전트

입력: 제공된 SPEC-ID의 SPEC 문서 내용.

manager-strategy의 작업:

- SPEC 문서 전체 읽기 및 분석
- 요구사항 및 성공 기준 추출
- 구현 단계 및 개별 작업 식별
- 필요한 기술 스택 및 의존성 결정
- 복잡도 및 작업량 추정
- 단계별 접근 방식을 포함한 상세 실행 전략 수립

출력: plan_summary, 요구사항 목록, success_criteria, effort_estimate를 포함한 실행 계획.

### 결정 지점 1: 계획 승인

도구: AskUserQuestion (오케스트레이터 레벨)

옵션:

- 계획대로 진행 (Phase 1.5로 계속)
- 계획 수정 (피드백 수집, Phase 1 재실행)
- 연기 (종료, 나중에 계속)

사용자가 "진행"을 선택하지 않으면: 실행 종료.

### Phase 1.5: 작업 분해

에이전트: manager-strategy 서브에이전트 (계속)

목적: 승인된 실행 계획을 SDD 2025 표준에 따른 원자적이고 검토 가능한 작업으로 분해합니다.

manager-strategy의 작업:

- 계획을 원자적 구현 작업으로 분해
- 각 작업은 단일 DDD/TDD 사이클에서 완료 가능해야 함
- 각 작업에 우선순위 및 의존성 할당
- 진행 상황 가시성을 위한 작업 추적 항목 생성
- 작업 범위가 모든 SPEC 요구사항을 충족하는지 확인

각 분해된 작업의 구조:

- 작업 ID: SPEC 내 순차적 (TASK-001, TASK-002 등)
- 설명: 명확한 액션 문장
- 요구사항 매핑: 충족하는 SPEC 요구사항
- 의존성: 선행 작업 목록
- 인수 기준: 완료 검증 방법

제약 조건: 각 작업이 단일 DDD/TDD 사이클에서 완료되는 원자적 작업으로 분해합니다. 작업 수에 인위적인 제한 없음. SPEC 자체가 너무 복잡한 경우 SPEC 분할을 고려하세요.

출력: coverage_verified 플래그가 true로 설정된 작업 목록.

### 개발 모드 라우팅

Phase 2 전에 `.moai/config/sections/quality.yaml`을 읽어 개발 방법론을 결정합니다:

**development_mode가 "ddd"인 경우:**
- 모든 작업을 manager-ddd 서브에이전트로 라우팅
- ANALYZE-PRESERVE-IMPROVE 사이클 사용
- 특성화 테스트를 통한 동작 보존에 집중

**development_mode가 "tdd"인 경우:**
- 모든 작업을 manager-tdd 서브에이전트로 라우팅
- RED-GREEN-REFACTOR 사이클 사용
- 명세 테스트를 통한 테스트 우선 개발에 집중

**development_mode가 "hybrid" (권장)인 경우:**
- 변경 유형별로 각 작업 분류:
  - 새 파일 → manager-tdd로 라우팅 (TDD 워크플로우)
  - 기존 파일의 새 함수 → manager-tdd로 라우팅 (TDD 워크플로우)
  - 기존 코드 수정 → manager-ddd로 라우팅 (DDD 워크플로우)
  - 기존 코드 리팩토링 → manager-ddd로 라우팅 (DDD 워크플로우)
- 의존성 순서대로 작업 실행, 작업별로 적절한 에이전트로 라우팅

### Phase 2: 구현 (모드 의존적)

#### Phase 2A: DDD 구현 (ddd 모드 또는 hybrid 모드의 레거시 코드)

에이전트: manager-ddd 서브에이전트

입력: Phase 1의 승인된 실행 계획 + Phase 1.5의 작업 분해.

DDD 사이클은 세 단계로 실행됩니다:

- ANALYZE: 도메인 경계, 결합 지표, 리팩토링 대상 식별. 기존 코드 읽기 및 의존성 매핑.
- PRESERVE: 기존 테스트 확인. 변경 전 안전망 구축을 위해 커버되지 않은 코드 경로에 대한 특성화 테스트 생성.
- IMPROVE: 지속적인 검증을 통한 점진적 변환 적용. 각 변환 후 모든 테스트 실행.

요구사항:

- 리팩토링 단계 전반의 진행 상황 추적을 위한 작업 추적 초기화
- 완전한 ANALYZE-PRESERVE-IMPROVE 사이클 실행
- 각 변환 후 모든 기존 테스트 통과 확인
- 커버되지 않은 코드 경로에 대한 특성화 테스트 생성
- 테스트 커버리지 85% 이상 달성

출력: files_modified 목록, characterization_tests_created 목록, test_results (모두 통과), behavior_preserved 플래그, structural_metrics 비교, implementation_divergence 보고서.

구현 편차 추적:

manager-ddd 서브에이전트는 구현 중 원래 SPEC 계획과의 편차를 추적해야 합니다:

- planned_files: plan.md에 나열된 생성 또는 수정 예정 파일
- actual_files: DDD 사이클 중 실제 생성 또는 수정된 파일
- additional_features: 원래 SPEC 범위를 넘어 구현된 기능 또는 기능 (근거 포함)
- scope_changes: 구현 중 수행된 범위 조정 설명 (확장, 연기 또는 대체)
- new_dependencies: 도입된 새 라이브러리, 패키지 또는 외부 의존성
- new_directories: 생성된 새 디렉토리 구조

이 편차 데이터는 SPEC 문서 업데이트 및 프로젝트 문서 동기화를 위해 /moai sync에서 사용됩니다.

#### Phase 2B: TDD 구현 (tdd 모드 또는 hybrid 모드의 새 코드)

에이전트: manager-tdd 서브에이전트

입력: Phase 1의 승인된 실행 계획 + Phase 1.5의 작업 분해.

TDD 사이클은 세 단계로 실행됩니다:

- RED: 예상 동작을 정의하는 명세 테스트 작성. 테스트는 초기에 실패해야 합니다 (새로운 것을 테스트한다는 것을 확인).
- GREEN: 테스트를 통과하는 최소한의 구현 코드 작성. 정확성에 집중, 우아함은 나중에.
- REFACTOR: 테스트를 통과하는 상태를 유지하면서 코드 구조 개선. 클린 코드 원칙 적용.

요구사항:

- TDD 사이클 전반의 진행 상황 추적을 위한 작업 추적 초기화
- 각 기능에 대한 완전한 RED-GREEN-REFACTOR 사이클 실행
- 구현 전 테스트 작성 (테스트 우선 원칙)
- 커밋당 최소 80% 커버리지 (새 코드에는 85% 권장)

출력: files_created 목록, specification_tests_created 목록, test_results (모두 통과), coverage 퍼센트, refactoring_improvements 목록, implementation_divergence 보고서.

구현 편차 추적:

manager-tdd 서브에이전트는 구현 중 원래 SPEC 계획과의 편차를 추적해야 합니다:

- planned_files: plan.md에 나열된 생성 예정 파일
- actual_files: TDD 사이클 중 실제 생성된 파일
- additional_features: 원래 SPEC 범위를 넘어 구현된 기능 또는 기능 (근거 포함)
- scope_changes: 구현 중 수행된 범위 조정 설명
- new_dependencies: 도입된 새 라이브러리, 패키지 또는 외부 의존성
- new_directories: 생성된 새 디렉토리 구조

이 편차 데이터는 SPEC 문서 업데이트 및 프로젝트 문서 동기화를 위해 /moai sync에서 사용됩니다.

### Phase 2.5: 품질 검증

에이전트: manager-quality 서브에이전트

입력: Phase 1 계획 컨텍스트와 Phase 2 구현 결과 모두.

TRUST 5 검증 항목:

- Tested: 변경 전 테스트가 존재하고 통과. 테스트 주도 설계 원칙 유지.
- Readable: 코드가 프로젝트 관례를 따르고 문서 포함.
- Unified: 구현이 기존 프로젝트 패턴을 따름.
- Secured: 보안 취약점 없음. OWASP 준수 확인.
- Trackable: 모든 변경 사항이 명확한 커밋 메시지로 기록됨. 히스토리 분석 지원.

추가 검증 (모드 의존적):

DDD 모드의 경우:
- 테스트 커버리지 85% 이상
- 동작 보존: 모든 기존 테스트가 변경 없이 통과
- 특성화 테스트 통과: 동작 스냅샷 일치
- 구조적 개선: 결합 및 응집 지표 개선

TDD 모드의 경우:
- 커밋당 테스트 커버리지 80% 이상 (새 코드에는 85% 권장)
- 테스트 우선 원칙: 실패 테스트 없이 코드 작성 금지
- 모든 명세 테스트 통과
- REFACTOR 단계에서 클린 코드 원칙 적용

Hybrid 모드의 경우:
- 새 코드: TDD 커버리지 목표 (새 파일 85%)
- 수정된 코드: DDD 커버리지 목표 (동작 보존 포함 85%)
- 전반적인 커버리지 개선 추세 유지

출력: 항목별 trust_5_validation 결과, coverage 퍼센트, 전반적인 상태 (PASS, WARNING, 또는 CRITICAL), issues_found 목록.

### 품질 게이트 결정

상태가 CRITICAL인 경우:

- AskUserQuestion을 통해 사용자에게 품질 이슈 제시
- 수정을 위해 구현 단계로 돌아가는 옵션
- 현재 실행 흐름 종료

상태가 PASS 또는 WARNING인 경우: Phase 3으로 계속.

### LSP 품질 게이트

실행 단계는 quality.yaml에 설정된 LSP 기반 품질 게이트를 적용합니다:
- LSP 오류 없음 필수 (lsp_quality_gates.run.max_errors: 0)
- 타입 오류 없음 필수 (lsp_quality_gates.run.max_type_errors: 0)
- 린트 오류 없음 필수 (lsp_quality_gates.run.max_lint_errors: 0)
- 기준선 대비 회귀 불허 (lsp_quality_gates.run.allow_regression: false)

### Phase 3: Git 작업 (조건부)

에이전트: manager-git 서브에이전트

입력: Phase 1, 2, 2.5의 전체 컨텍스트.

실행 조건:

- quality_status가 PASS 또는 WARNING
- 설정 git_strategy.automation.auto_branch가 true인 경우: feature/SPEC-{ID} 기능 브랜치 생성
- auto_branch가 false인 경우: 현재 브랜치에 직접 커밋

manager-git의 작업:

- 기능 브랜치 생성 (auto_branch 활성화 시)
- 모든 관련 구현 및 테스트 파일 스테이징
- 컨벤셔널 커밋 메시지로 커밋 생성
- 각 커밋이 성공적으로 생성됐는지 확인

출력: branch_name, commits 배열 (sha와 message), files_staged 수, status.

### Phase 4: 완료 및 안내

도구: AskUserQuestion (오케스트레이터 레벨)

구현 요약 표시:

- 생성된 파일 수
- 통과된 테스트 수
- 커버리지 퍼센트
- 커밋 수

옵션:

- 문서 동기화 (권장): /moai sync 실행하여 문서 동기화 및 PR 생성
- 다른 기능 구현: 추가 SPEC을 위해 /moai plan으로 돌아가기
- 결과 검토: 구현 및 테스트 커버리지 로컬 검토
- 완료: 세션 종료

---

## 팀 모드 라우팅

--team 플래그가 제공되거나 자동 선택된 경우, 실행 단계는 반드시 팀 오케스트레이션으로 전환해야 합니다:

1. 전제 조건 확인: workflow.team.enabled == true AND CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 환경변수 설정
2. 전제 조건 충족 시: workflows/team-run.md를 읽고 팀 워크플로우 실행 (backend-dev + frontend-dev + tester + quality와 함께 TeamCreate)
3. 전제 조건 미충족 시: "팀 모드는 workflow.yaml에서 workflow.team.enabled: true 설정과 CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 환경변수가 필요합니다"라고 사용자에게 경고 후 표준 서브에이전트 모드 (development_mode 기반 manager-ddd/tdd)로 폴백

팀 구성: backend-dev (inherit) + frontend-dev (inherit) + tester (inherit) + quality (inherit, 읽기 전용)

상세한 팀 오케스트레이션 단계는 workflows/team-run.md를 참조하세요.

---

## 컨텍스트 전파

컨텍스트는 모든 단계를 통해 앞으로 전달됩니다:

- Phase 1 → Phase 2: 아키텍처 결정을 포함한 실행 계획이 구현을 안내
- Phase 2 → Phase 2.5: 구현 코드와 계획 컨텍스트가 컨텍스트 인식 검증을 가능하게 함
- Phase 2.5 → Phase 3: 품질 발견사항이 의미 있는 커밋 메시지를 가능하게 함
- Phase 2 → /moai sync: 구현 편차 보고서가 SPEC 및 프로젝트 문서의 정확한 업데이트를 가능하게 함

이점: 단계 간 재분석 없음. 아키텍처 결정이 자연스럽게 전파됨. 커밋이 변경된 내용과 이유 모두 설명. 편차 추적으로 동기화 단계가 SPEC 및 프로젝트 문서를 정확하게 업데이트할 수 있음.

---

## 완료 기준

다음 사항을 모두 확인해야 합니다:

- Phase 1: manager-strategy가 요구사항 및 성공 기준을 포함한 실행 계획 반환
- 사용자 승인 체크포인트가 사용자 확인 전까지 Phase 2를 차단
- Phase 1.5: 요구사항 추적성을 가진 작업 분해 완료
- Phase 2: development_mode에 따른 구현 완료:
  - DDD 모드: manager-ddd가 85%+ 커버리지로 ANALYZE-PRESERVE-IMPROVE 실행
  - TDD 모드: manager-tdd가 85%+ 커버리지로 RED-GREEN-REFACTOR 실행
  - Hybrid 모드: 작업 유형별 적절한 에이전트로 85%+ 통합 커버리지 목표 달성
- Phase 2.5: manager-quality가 PASS 또는 WARNING 상태로 TRUST 5 검증 완료
- 품질 게이트가 상태 CRITICAL인 경우 Phase 3 차단
- Phase 3: 품질이 허용된 경우에만 manager-git이 커밋 생성 (브랜치 또는 직접)
- Phase 4: 사용자에게 다음 단계 옵션 제시

---

Version: 2.0.0
Updated: 2026-02-07
Source: .claude/commands/moai/2-run.md v5.0.0에서 추출. 구현 편차 추적, development_mode 라우팅 (ddd/tdd/hybrid), 팀 모드 지원, LSP 품질 게이트 추가.
