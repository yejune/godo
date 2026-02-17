---
name: moai-foundation-core
description: >
  TRUST 5 품질 프레임워크, SPEC 우선 DDD 방법론, 위임 패턴, 점진적 공개,
  에이전트 카탈로그 참조를 포함한 MoAI-ADK 기초 원칙을 제공합니다.
  TRUST 5 게이트, SPEC 워크플로우, EARS 형식, DDD 방법론,
  에이전트 위임 패턴 또는 MoAI 오케스트레이션 규칙을 참조할 때 사용합니다.
  컨텍스트 및 토큰 관리에는 사용하지 마세요 (moai-foundation-context 사용).
  전략적 분석에는 사용하지 마세요 (moai-foundation-philosopher 사용).
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Grep Glob mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "2.5.0"
  category: "foundation"
  status: "active"
  updated: "2026-01-21"
  modularized: "true"
  tags: "foundation, core, orchestration, agents, commands, trust-5, spec-first-ddd"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords:
    - "trust-5"
    - "spec-first"
    - "ddd"
    - "delegation"
    - "agent"
    - "token"
    - "progressive disclosure"
    - "modular"
    - "workflow"
    - "orchestration"
    - "quality gate"
    - "spec"
    - "ears format"
  agents:
    - "manager-spec"
    - "manager-ddd"
    - "manager-strategy"
    - "manager-quality"
    - "builder-agent"
    - "builder-skill"
  phases:
    - "plan"
    - "run"
    - "sync"
---

# MoAI Foundation Core

MoAI-ADK의 AI 기반 개발 워크플로우를 구동하는 기초 원칙 및 아키텍처 패턴.

핵심 철학: 검증된 패턴과 자동화된 워크플로우를 통한 품질 우선, 도메인 주도, 모듈식, 효율적인 AI 개발.

## 빠른 참조

MoAI Foundation Core란?

AI 기반 개발에서 품질, 효율성, 확장성을 보장하는 여섯 가지 핵심 원칙:

1. TRUST 5 프레임워크 - 품질 게이트 시스템 (Tested, Readable, Unified, Secured, Trackable)
2. SPEC 우선 DDD - 명세서 기반 도메인 주도 개발 워크플로우
3. 위임 패턴 - 전문 에이전트를 통한 작업 오케스트레이션 (직접 실행 금지)
4. 토큰 최적화 - 200K 예산 관리 및 컨텍스트 효율성
5. 점진적 공개 - 3단계 지식 전달 (빠른 참조, 구현, 고급)
6. 모듈식 시스템 - 확장성을 위한 파일 분할 및 참조 아키텍처

빠른 접근:

- 품질 기준: modules/trust-5-framework.md
- 개발 워크플로우: modules/spec-first-ddd.md
- 에이전트 조율: modules/delegation-patterns.md
- 예산 관리: modules/token-optimization.md
- 콘텐츠 구조: modules/progressive-disclosure.md
- 파일 구성: modules/modular-system.md
- 에이전트 카탈로그: modules/agents-reference.md
- 명령어 참조: modules/commands-reference.md
- 보안 및 제약: modules/execution-rules.md

사용 사례:

- 품질 기준을 갖춘 새 에이전트 생성
- 구조적 가이드라인을 갖춘 새 스킬 개발
- 복잡한 워크플로우 오케스트레이션
- 토큰 예산 계획 및 최적화
- 문서 아키텍처 설계
- 품질 게이트 설정

---

## 구현 가이드

### 1. TRUST 5 프레임워크 - 품질 보증 시스템

목적: 코드 품질, 보안, 유지 관리성을 보장하는 자동화된 품질 게이트.

다섯 가지 기둥:

Tested 기둥: 동작 보존을 보장하는 특성화 테스트로 포괄적인 테스트 커버리지 유지. coverage 보고와 함께 pytest 실행. 실패 시 병합 차단 및 누락 테스트 생성. 특성화 테스트는 레거시 코드의 현재 동작을 포착하고, 명세서 테스트는 새 코드의 도메인 요구사항을 검증. 높은 커버리지는 코드 신뢰성을 보장하고 프로덕션 결함을 60-70% 줄임.

Readable 기둥: 명확하고 설명적인 명명 규칙 사용. ruff 린터 검사 실행. 실패 시 경고 및 리팩토링 개선 제안. 명확한 명명은 코드 이해 및 팀 협업을 향상. 온보딩 시간을 40% 줄이고 유지 관리 속도 향상.

Unified 기둥: 일관된 포맷팅 및 import 패턴 적용. black 포맷터 및 isort 검사 실행. 실패 시 코드 자동 포맷 또는 경고. 일관성은 스타일 논쟁과 병합 충돌을 제거. 코드 리뷰 시간을 30% 줄이고 가독성 향상.

Secured 기둥: OWASP 보안 기준 준수. security-expert 에이전트 분석 실행. 실패 시 병합 차단 및 보안 검토 요구. 보안 취약점은 심각한 비즈니스 및 법적 위험을 초래. 일반적인 보안 취약점의 95% 이상을 예방.

Trackable 기둥: 명확하고 구조화된 커밋 메시지 작성. Git 커밋 메시지 정규식 패턴 매칭. 실패 시 적절한 커밋 메시지 형식 제안. 명확한 히스토리는 디버깅, 감사, 협업을 가능하게 함. 이슈 조사 시간을 50% 줄임.

통합 지점: 자동화된 검증을 위한 pre-commit 훅, 품질 게이트 적용을 위한 CI/CD 파이프라인, core-quality 검증을 위한 에이전트 워크플로우, 품질 지표를 위한 문서.

상세 참조: modules/trust-5-framework.md

---

### 2. SPEC 우선 DDD - 개발 워크플로우

목적: 구현 전 명확한 요구사항을 보장하는 명세서 기반 개발.

3단계 워크플로우:

단계 1 SPEC (/moai:1-plan): workflow-spec가 EARS 형식 생성. 출력은 .moai/specs/SPEC-XXX/spec.md. /clear로 45-50K 토큰 절약.

단계 2 DDD (/moai:2-run): 요구사항 분석(ANALYZE), 기존 동작 보존(PRESERVE), 개선(IMPROVE). 최소 85% 커버리지로 검증.

단계 3 Docs (/moai:3-sync): API 문서, 아키텍처 다이어그램, 프로젝트 보고서.

EARS 형식: 시스템 전반에 항상 활성화된 요구사항을 위한 Ubiquitous. X 발생 시 Y 수행하는 트리거 기반 요구사항을 위한 Event-driven. X 동안 Y 수행하는 조건부 요구사항을 위한 State-driven. X를 절대 수행하지 않는 금지 요구사항을 위한 Unwanted. 가능한 경우 X를 제공하는 선택적 요구사항을 위한 Optional.

토큰 예산: SPEC은 30K, DDD는 180K, Docs는 40K, 총계는 250K.

핵심 실천: /clear는 단계 1 완료 후 실행하여 컨텍스트 초기화.

상세 참조: modules/spec-first-ddd.md

---

### 3. 위임 패턴 - 에이전트 오케스트레이션

목적: 직접 실행을 피하고 전문 에이전트에게 작업 위임.

핵심 원칙: MoAI는 모든 작업을 Task()를 통해 전문 에이전트에게 위임해야 합니다. 직접 실행은 전문화, 품질 게이트, 토큰 최적화를 우회합니다. 적절한 위임은 작업 성공률을 40% 향상시키고 병렬 실행을 가능하게 합니다.

위임 구문: 전문 에이전트를 위한 subagent_type 파라미터, 명확하고 구체적인 작업을 위한 prompt 파라미터, 관련 데이터 딕셔너리를 위한 context 파라미터로 Task를 호출합니다.

세 가지 패턴:

의존성을 위한 순차: api-designer에게 설계를 위한 Task 호출 후, 설계 컨텍스트와 함께 구현을 위해 backend-expert에게 Task 호출.

독립 작업을 위한 병렬: backend-expert와 frontend-expert에게 동시에 Promise.all로 Task 호출.

분석 기반 조건부: 분석을 위해 debug-helper에게 Task 호출 후, analysis.type에 따라 security-expert 또는 다른 적절한 에이전트에게 Task 호출.

에이전트 선택: 1개 파일의 단순 작업은 1-2개 에이전트 순차. 3-5개 파일의 중간 작업은 2-3개 에이전트 순차. 10개 이상 파일의 복잡한 작업은 5개 이상 에이전트 혼합.

상세 참조: modules/delegation-patterns.md

---

### 4. 토큰 최적화 - 예산 관리

목적: 전략적 컨텍스트 관리를 통한 효율적인 200K 토큰 예산 사용.

예산 할당:

SPEC 단계는 30K 토큰. 요구사항만 로드하고 완료 후 /clear 실행 전략. 명세서 단계는 요구사항 분석을 위해 최소한의 컨텍스트만 필요. 구현 단계를 위해 45-50K 토큰 절약.

DDD 단계는 180K 토큰. 선택적 파일 로드, 구현 관련 파일만 로드 전략. 구현은 깊은 컨텍스트가 필요하지만 전체 코드베이스는 불필요. 예산 내에서 70% 더 큰 구현 가능.

Docs 단계는 40K 토큰. 결과 캐싱 및 템플릿 재사용 전략. 문서는 완료된 작업 산출물을 기반으로 구축. 중복 파일 읽기를 60% 줄임.

모든 단계의 총 예산은 250K 토큰. 단계 간 컨텍스트 초기화로 단계 분리는 깨끗한 컨텍스트 경계를 제공하고 토큰 비대화를 방지. 동일 예산 내에서 2-3배 더 큰 프로젝트 가능.

토큰 절약 전략:

단계 분리: 단계 사이에 /clear 실행, /moai:1-plan 완료 후 45-50K 절약, 컨텍스트가 150K 초과 시, 50개 이상 메시지 후.

선택적 로드: 필요한 파일만 로드.

컨텍스트 최적화: 20-30K 토큰 목표.

모델 선택: 품질을 위한 Sonnet, 속도와 비용을 위한 Haiku (70% 저렴, 총 60-70% 절약).

상세 참조: modules/token-optimization.md

---

### 5. 점진적 공개 - 콘텐츠 아키텍처

목적: 깊이와 가치의 균형을 맞추는 3단계 지식 전달.

세 가지 레벨:

빠른 참조 레벨: 30초 투자 시간, 핵심 원칙 및 필수 개념, 약 1,000 토큰. 시간이 제한된 사용자를 위한 신속한 가치 제공. 사용자는 5%의 시간으로 80% 이해 달성.

구현 레벨: 5분 투자 시간, 워크플로우, 실제 예제, 통합 패턴, 약 3,000 토큰. 실행 가능한 안내로 개념과 실행을 연결. 깊은 전문성 없이 즉시 생산적인 작업 가능.

고급 레벨: 10분 이상 투자 시간, 심층 기술 분석, 엣지 케이스, 최적화 기법, 약 5,000 토큰. 복잡한 시나리오를 위한 마스터 수준의 지식 제공. 포괄적인 커버리지로 에스컬레이션을 70% 줄임.

SKILL.md 구조 (최대 500줄): 빠른 참조 섹션, 구현 가이드 섹션, 고급 패턴 섹션, 잘 어울리는 것들 섹션.

모듈 아키텍처: 교차 참조가 있는 진입점으로서의 SKILL.md, 무제한 크기의 심층 분석을 위한 modules 디렉토리, 작동하는 샘플을 위한 examples.md, 외부 링크를 위한 reference.md.

500줄 초과 시 파일 분할: SKILL.md는 빠른 참조 80-120줄, 구현 180-250줄, 고급 80-140줄, 참조 10-20줄 포함. 넘치는 내용은 modules/topic.md로.

상세 참조: modules/progressive-disclosure.md

---

### 6. 모듈식 시스템 - 파일 구성

목적: 무제한 콘텐츠를 가능하게 하는 확장 가능한 파일 구조.

표준 구조: 500줄 미만의 핵심 파일인 SKILL.md, 패턴 포함 무제한 크기의 확장 콘텐츠를 위한 modules 디렉토리, 작업 샘플을 위한 examples.md, 외부 링크를 위한 reference.md, 유틸리티를 위한 scripts 디렉토리 (선택), templates 디렉토리 (선택)를 포함하는 .claude/skills/skill-name/ 디렉토리를 생성합니다.

파일 원칙: SKILL.md는 점진적 공개 및 교차 참조를 포함하여 500줄 미만 유지. modules 디렉토리는 주제 중심으로 제한 없고 자체 완결적 콘텐츠. examples.md는 주석이 있는 복사-붙여넣기 가능 형식. reference.md는 API 문서 및 리소스 포함.

교차 참조 구문: modules/patterns.md의 세부 정보 참조, examples.md#auth의 예제 참조, reference.md#api의 외부 문서 참조.

탐색 흐름: SKILL.md에서 주제로, modules/topic.md에서 심층 분석으로.

상세 참조: modules/modular-system.md

---

## 고급 구현

교차 모듈 통합, 품질 검증, 에러 처리를 포함한 고급 패턴은 상세 모듈 참조에서 확인할 수 있습니다.

주요 고급 주제:

- 교차 모듈 통합: TRUST 5 + SPEC 우선 DDD 결합
- 토큰 최적화 위임: 컨텍스트 초기화를 통한 병렬 실행
- 점진적 에이전트 워크플로우: 에스컬레이션 패턴
- 품질 검증: 실행 전후 검증
- 에러 처리: 위임 실패 복구

상세 참조: 작동하는 코드 샘플을 위한 examples.md

---

## 잘 어울리는 것들

에이전트: 기초 원칙으로 에이전트를 생성하기 위한 agent-factory, 모듈식 아키텍처로 스킬을 생성하기 위한 skill-factory, 자동화된 TRUST 5 검증을 위한 core-quality, EARS 형식 명세서를 위한 workflow-spec, ANALYZE-PRESERVE-IMPROVE 실행을 위한 workflow-ddd, 점진적 공개를 통한 문서화를 위한 workflow-docs.

스킬: 기초 패턴을 갖춘 CLAUDE.md를 위한 moai-cc-claude-md, TRUST 5를 갖춘 설정을 위한 moai-cc-configuration, 토큰 최적화를 위한 moai-cc-memory, MCP 통합을 위한 moai-context7-integration.

도구: 직접 사용자 상호작용 및 명확화 필요 시 AskUserQuestion.

명령어: SPEC 우선 단계 1을 위한 /moai:1-plan, DDD 단계 2를 위한 /moai:2-run, 문서 단계 3을 위한 /moai:3-sync, 지속적 개선을 위한 /moai:9-feedback, 토큰 관리를 위한 /clear.

Foundation 모듈 (확장 문서): 7단계 계층 구조를 가진 26개 에이전트 카탈로그를 위한 modules/agents-reference.md, 6개 핵심 명령어 워크플로우를 위한 modules/commands-reference.md, 보안, Git 전략, 규정 준수를 위한 modules/execution-rules.md.

---

## 빠른 결정 가이드

새 에이전트: 주요 원칙은 TRUST 5와 위임. 지원 원칙은 토큰 최적화와 모듈식.

새 스킬: 주요 원칙은 점진적 공개와 모듈식. 지원 원칙은 TRUST 5와 토큰 최적화.

워크플로우: 주요 원칙은 위임 패턴. 지원 원칙은 SPEC 우선과 토큰 최적화.

품질: 주요 원칙은 TRUST 5 프레임워크. 지원 원칙은 SPEC 우선 DDD.

예산: 주요 원칙은 토큰 최적화. 지원 원칙은 점진적 공개와 모듈식.

문서: 주요 원칙은 점진적 공개와 모듈식. 지원 원칙은 토큰 최적화.

모듈 심층 분석: modules/trust-5-framework.md, modules/spec-first-ddd.md, modules/delegation-patterns.md, modules/token-optimization.md, modules/progressive-disclosure.md, modules/modular-system.md, modules/agents-reference.md, modules/commands-reference.md, modules/execution-rules.md.

전체 예제: examples.md
외부 리소스: reference.md
