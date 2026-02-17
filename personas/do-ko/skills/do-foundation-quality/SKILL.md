---
name: do-foundation-quality
description: >
  TRUST 5 검증, 선제적 코드 분석, 린팅 기준, 자동화된 모범 사례를 적용하는
  코드 품질 오케스트레이터입니다.
  코드 리뷰, 품질 게이트 검사, 린트 구성, TRUST 5 준수 검증,
  또는 코딩 기준 확립 시 사용하세요.
  테스트 작성(do-workflow-testing 사용)이나
  런타임 오류 디버깅(expert-debug 에이전트 사용)에는 사용하지 마세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Grep Glob mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "2.2.0"
  category: "foundation"
  status: "active"
  updated: "2026-01-11"
  modularized: "true"
  tags: "foundation, quality, testing, validation, trust-5, best-practices, code-review"
  aliases: "do-foundation-quality"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# Do Extension: Triggers
triggers:
  keywords:
    - "quality"
    - "testing"
    - "test"
    - "validation"
    - "trust-5"
    - "best practice"
    - "code review"
    - "linting"
    - "coverage"
    - "pytest"
    - "security"
    - "ci/cd"
    - "quality gate"
    - "proactive"
    - "code smell"
    - "technical debt"
    - "refactoring"
  agents:
    - "manager-quality"
    - "manager-ddd"
    - "expert-testing"
    - "expert-security"
    - "expert-refactoring"
  phases:
    - "run"
    - "sync"
  languages:
    - "python"
    - "javascript"
    - "typescript"
    - "java"
    - "go"
    - "rust"
    - "cpp"
    - "csharp"
---

# 엔터프라이즈 코드 품질 오케스트레이터

체계적인 코드 리뷰, 선제적 개선 제안, 자동화된 모범 사례 적용을 결합한 엔터프라이즈급 코드 품질 관리 시스템입니다. Context7 통합을 통해 실시간 모범 사례를 제공하며 TRUST 5 프레임워크 검증으로 포괄적인 품질 보증을 제공합니다.

## 빠른 참조 (30초)

핵심 기능:

- TRUST 5 검증: Testable, Readable, Unified, Secured, Trackable 품질 게이트
- 선제적 분석: 자동화된 이슈 감지 및 개선 제안
- 모범 사례 적용: Context7 기반 실시간 기준 검증
- 다중 언어 지원: 25개 이상 프로그래밍 언어의 전문화된 규칙
- 엔터프라이즈 통합: CI/CD 파이프라인, 품질 지표, 보고

주요 패턴:

- 품질 게이트 파이프라인: 구성 가능한 임계값으로 자동 검증
- 선제적 스캐너: 개선 권장 사항을 포함한 지속적 분석
- 모범 사례 엔진: Context7 기반 기준 적용
- 품질 지표 대시보드: 포괄적인 보고 및 트렌드 분석

사용 시점:

- 코드 리뷰 자동화 및 품질 게이트 적용
- 선제적 코드 품질 향상 및 기술 부채 감소
- 엔터프라이즈 코딩 기준 적용 및 준수 검증
- 자동화된 품질 검사를 통한 CI/CD 파이프라인 통합

빠른 접근:

- TRUST 5 프레임워크: [trust5-validation.md](modules/trust5-validation.md) 참조
- 선제적 분석: [proactive-analysis.md](modules/proactive-analysis.md) 참조
- 모범 사례: [best-practices.md](modules/best-practices.md) 참조
- 통합 패턴: [integration-patterns.md](modules/integration-patterns.md) 참조

## 구현 가이드

### 시작하기

기본 품질 검증: trust5_enabled, proactive_analysis, best_practices_enforcement, context7_integration을 모두 True로 설정하여 QualityOrchestrator를 초기화합니다. path 매개변수를 소스 디렉토리로, languages 목록에 python, javascript, typescript를 포함하고, quality_threshold를 0.85로 설정하여 analyze_codebase 메서드를 호출합니다. 이 메서드는 포괄적인 품질 결과를 반환합니다.

TRUST 5를 통한 품질 게이트 검증의 경우, QualityGate 인스턴스를 생성하고 codebase_path, test_coverage_threshold 0.90, complexity_threshold 10으로 validate_trust5를 호출합니다.

선제적 품질 분석: context7_client와 BestPracticesEngine rule_engine으로 ProactiveQualityScanner를 초기화합니다. path와 security, performance, maintainability, testing을 포함한 scan_types 목록으로 scan_codebase를 호출합니다. issues, priority를 high로, auto_fix를 활성화하여 generate_recommendations를 호출하여 권장 사항을 생성합니다.

### 핵심 컴포넌트

#### 품질 오케스트레이션 엔진

QualityOrchestrator 클래스는 TRUST 5 프레임워크를 통해 엔터프라이즈 품질 오케스트레이션을 제공합니다. QualityConfig로 초기화하고 TRUST5Validator, ProactiveScanner, BestPracticesEngine, Context7Client, QualityMetricsCollector 인스턴스를 생성합니다.

analyze_codebase 메서드는 네 단계로 포괄적인 분석을 수행합니다. 1단계에서 지정된 임계값으로 코드베이스에 TRUST 5 검증을 실행합니다. 2단계에서 포커스 영역을 스캔하는 선제적 분석을 수행합니다. 3단계에서 Context7 문서를 활성화하여 지정된 언어에 대한 모범 사례를 확인합니다. 4단계에서 모든 분석 결과로부터 포괄적인 지표를 수집합니다.

이 메서드는 trust5_validation, proactive_analysis, best_practices, metrics, 그리고 모든 결과로부터 계산된 overall_score를 포함하는 QualityResult를 반환합니다.

상세 구현은 모듈에서 확인 가능합니다:

- [trust5-validation.md](modules/trust5-validation.md): TRUST 5 검증기 구현
- [proactive-analysis.md](modules/proactive-analysis.md): 선제적 스캐너 구현
- [best-practices.md](modules/best-practices.md): 모범 사례 엔진 구현

### 구성 및 사용자 정의

품질 구성: quality_orchestration 섹션이 있는 quality-config.yaml을 생성합니다.

trust5_framework 하위에 enabled를 true로, overall(0.85), testable(0.90), readable(0.80), unified(0.85), secured(0.90), trackable(0.80) 임계값을 설정합니다.

proactive_analysis 하위에 enabled를 true로, scan_frequency를 daily로, focus_areas 목록에 performance, security, maintainability, technical_debt를 포함합니다.

auto_fix 하위에 enabled를 true로, severity_threshold를 medium으로, confirmation_required를 true로 설정합니다.

best_practices 하위에 enabled를 true로, context7_integration을 true로, auto_update_standards를 true로, compliance_target을 0.85로 설정합니다.

language_rules 하위에 python은 pep8 style_guide, black formatter, ruff linter, mypy type_checker로 구성합니다. javascript는 airbnb style_guide, prettier formatter, eslint linter로 구성합니다. typescript는 google style_guide, prettier formatter, eslint linter로 구성합니다.

reporting 하위에 enabled를 true로, metrics_retention_days를 90으로, trend_analysis를 true로, executive_dashboard를 true로 설정합니다.

notifications 하위에 quality_degradation, security_vulnerabilities, technical_debt_increase를 활성화합니다.

통합 예제: CI/CD 파이프라인 통합, GitHub Actions 통합, Quality-as-Service REST API, 프로젝트 간 벤치마킹은 [통합 패턴](modules/integration-patterns.md)을 참조하세요.

## 고급 패턴

### 사용자 정의 품질 규칙

name, validator 호출 가능 객체, severity(기본값 medium)를 가진 CustomQualityRule 클래스를 생성합니다. validate 비동기 메서드는 코드베이스에 validator를 실행하고 try-except로 감쌉니다. 성공 시 rule_name, passed 상태, severity, details, recommendations가 포함된 RuleResult를 반환합니다. 예외 발생 시 passed가 false이고, severity가 error이며, 오류 세부 정보와 수정 권장 사항이 포함된 RuleResult를 반환합니다.

전체 예제는 [모범 사례 - 사용자 정의 규칙](modules/best-practices.md#custom-quality-rules)을 참조하세요.

### 머신 러닝 품질 예측

코드 특성 추출 및 예측 모델을 사용한 ML 기반 품질 이슈 예측. 구현 세부 사항은 [선제적 분석 - ML 예측](modules/proactive-analysis.md#machine-learning-quality-prediction)을 참조하세요.

### 실시간 품질 모니터링

품질 저하 및 보안 취약점에 대한 자동화된 알림을 포함한 지속적 품질 모니터링. 구현 세부 사항은 [선제적 분석 - 실시간 모니터링](modules/proactive-analysis.md#real-time-quality-monitoring)을 참조하세요.

### 프로젝트 간 품질 벤치마킹

업계 유사 프로젝트와 프로젝트 품질 지표를 비교합니다. 구현 세부 사항은 [통합 패턴 - 벤치마킹](modules/integration-patterns.md#cross-project-quality-benchmarking)을 참조하세요.

## 모듈 참조

### 핵심 모듈

- [TRUST 5 검증](modules/trust5-validation.md) - 포괄적인 품질 프레임워크 검증
- [선제적 분석](modules/proactive-analysis.md) - 자동화된 이슈 감지 및 개선
- [모범 사례](modules/best-practices.md) - Context7 기반 기준 적용
- [통합 패턴](modules/integration-patterns.md) - CI/CD 및 엔터프라이즈 통합

### 모듈별 주요 컴포넌트

TRUST 5 검증: 5개 기둥 품질 검증을 위한 TRUST5Validator, 테스트 커버리지 및 품질을 위한 TestableValidator, 보안 및 OWASP 준수를 위한 SecuredValidator, 품질 게이트 파이프라인 통합.

선제적 분석: 자동화된 이슈 감지를 위한 ProactiveQualityScanner, ML 기반 예측을 위한 QualityPredictionEngine, 지속적 모니터링을 위한 RealTimeQualityMonitor, 성능 및 유지 보수성 분석.

모범 사례: 기준 검증을 위한 BestPracticesEngine, 최신 문서를 위한 Context7 통합, 사용자 정의 품질 규칙, 언어별 검증기.

통합 패턴: CI/CD 파이프라인 통합, GitHub Actions 워크플로우, Quality-as-Service REST API, 프로젝트 간 벤치마킹.

## Context7 라이브러리 매핑

품질 분석 도구 및 프레임워크를 위한 필수 라이브러리 매핑. 전체 목록은 [모범 사례 - 라이브러리 매핑](modules/best-practices.md#context7-library-mappings)을 참조하세요.

## 함께 사용하면 좋은 것들

에이전트:

- core-planner - 품질 요구사항 계획
- workflow-ddd - DDD 구현 검증
- security-expert - 보안 취약점 분석
- code-backend - 백엔드 코드 품질
- code-frontend - 프론트엔드 코드 품질

스킬:

- do-foundation-core - TRUST 5 프레임워크 참조
- do-workflow-ddd - DDD 워크플로우 검증
- do-security-owasp - 보안 준수
- do-context7-integration - Context7 모범 사례
- do-performance-optimization - 성능 분석

커맨드:

- /do:2-run - DDD 검증 통합
- /do:3-sync - 문서 품질 검사
- /do:9-feedback - 품질 개선 피드백

## 빠른 참조 요약

핵심 기능: TRUST 5 검증, 선제적 스캐닝, Context7 기반 모범 사례, 다중 언어 지원, 엔터프라이즈 통합

주요 클래스: QualityOrchestrator, TRUST5Validator, ProactiveQualityScanner, BestPracticesEngine, QualityMetricsCollector

필수 메서드: analyze_codebase(), validate_trust5(), scan_for_issues(), validate_best_practices(), generate_quality_report()

통합 준비: CI/CD 파이프라인, GitHub Actions, REST API, 실시간 모니터링, 프로젝트 간 벤치마킹

엔터프라이즈 기능: 사용자 정의 규칙, ML 예측, 실시간 모니터링, 벤치마킹, 포괄적인 보고

품질 기준: OWASP 준수, TRUST 5 프레임워크, Context7 통합, 자동화된 개선 권장 사항
