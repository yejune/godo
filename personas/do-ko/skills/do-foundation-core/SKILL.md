---
name: do-foundation-core
description: >
  TRUST 5 품질 프레임워크, SPEC-First DDD 방법론, 위임 패턴, 점진적 공개,
  에이전트 카탈로그 참조 등 Do-ADK의 핵심 원칙을 제공합니다.
  TRUST 5 게이트, SPEC 워크플로우, EARS 형식, DDD 방법론,
  에이전트 위임 패턴, 또는 Do 오케스트레이션 규칙을 참조할 때 사용하세요.
  컨텍스트 및 토큰 관리(do-foundation-context 사용)나
  전략적 분석(do-foundation-philosopher 사용)에는 사용하지 마세요.
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

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# Do Extension: Triggers
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

# Do Foundation Core

Do-ADK의 AI 기반 개발 워크플로우를 구동하는 핵심 원칙 및 아키텍처 패턴입니다.

핵심 철학: 검증된 패턴과 자동화된 워크플로우를 통해 품질 우선, 도메인 주도, 모듈화되고 효율적인 AI 개발을 추구합니다.

## 빠른 참조

Do Foundation Core란?

AI 기반 개발에서 품질, 효율성, 확장성을 보장하는 여섯 가지 핵심 원칙:

1. TRUST 5 프레임워크 - 품질 게이트 시스템 (Tested, Readable, Unified, Secured, Trackable)
2. SPEC-First DDD - 명세 주도 도메인 기반 개발 워크플로우
3. 위임 패턴 - 전문 에이전트를 통한 작업 오케스트레이션 (직접 실행 금지)
4. 토큰 최적화 - 200K 예산 관리 및 컨텍스트 효율성
5. 점진적 공개 - 3단계 지식 전달 (빠른 참조, 구현, 고급)
6. 모듈 시스템 - 확장성을 위한 파일 분할 및 참조 아키텍처

빠른 접근:

- modules/trust-5-framework.md: 품질 기준
- modules/spec-first-ddd.md: 개발 워크플로우
- modules/delegation-patterns.md: 에이전트 협조
- modules/token-optimization.md: 예산 관리
- modules/progressive-disclosure.md: 콘텐츠 구조
- modules/modular-system.md: 파일 구성
- modules/agents-reference.md: 에이전트 카탈로그
- modules/commands-reference.md: 커맨드 참조
- modules/execution-rules.md: 보안 및 제약 사항

사용 사례:

- 품질 기준을 적용한 새 에이전트 생성
- 구조적 가이드라인을 따른 새 스킬 개발
- 복잡한 워크플로우 오케스트레이션
- 토큰 예산 계획 및 최적화
- 문서 아키텍처 설계
- 품질 게이트 구성

---

## 구현 가이드

### 1. TRUST 5 프레임워크 - 품질 보증 시스템

목적: 코드 품질, 보안, 유지 보수성을 보장하는 자동화된 품질 게이트.

다섯 가지 기둥:

Tested(테스트됨) 기둥: 행동 보존을 보장하는 특성화 테스트를 포함한 포괄적인 테스트 커버리지 유지. 커버리지 보고와 함께 pytest 실행. 실패 시 머지를 차단하고 누락된 테스트를 생성. 특성화 테스트는 레거시 코드의 현재 동작을 캡처하고, 명세 테스트는 새 코드의 도메인 요구사항을 검증합니다. 높은 커버리지는 코드 안정성을 보장하고 프로덕션 결함을 줄입니다. 리팩토링 중 동작을 보존하고 디버깅 시간을 60~70% 줄입니다.

Readable(읽기 쉬운) 기둥: 명확하고 설명적인 네이밍 규칙 사용. ruff 린터 검사 실행. 실패 시 경고를 내고 리팩토링 개선 사항을 제안. 명확한 네이밍은 코드 이해도와 팀 협업을 향상시킵니다. 온보딩 시간을 40% 줄이고 유지 보수 속도를 높입니다.

Unified(통일된) 기둥: 일관된 포맷 및 임포트 패턴 적용. black 포매터와 isort 검사 실행. 실패 시 코드를 자동으로 포맷하거나 경고 발생. 일관성은 스타일 논쟁과 머지 충돌을 제거합니다. 코드 리뷰 시간을 30% 줄이고 가독성을 향상시킵니다.

Secured(보안) 기둥: OWASP 보안 표준 준수. security-expert 에이전트 분석 실행. 실패 시 머지를 차단하고 보안 검토를 요구. 보안 취약점은 심각한 비즈니스 및 법적 위험을 야기합니다. 일반적인 보안 취약점의 95% 이상을 방지합니다.

Trackable(추적 가능한) 기둥: 명확하고 구조화된 커밋 메시지 작성. Git 커밋 메시지 정규식 패턴 매칭. 실패 시 올바른 커밋 메시지 형식 제안. 명확한 히스토리는 디버깅, 감사, 협업을 가능하게 합니다. 이슈 조사 시간을 50% 줄입니다.

통합 포인트: 자동 검증을 위한 pre-commit 훅, 품질 게이트 적용을 위한 CI/CD 파이프라인, core-quality 검증을 위한 에이전트 워크플로우, 품질 지표 문서화.

상세 참조: modules/trust-5-framework.md

---

### 2. SPEC-First DDD - 개발 워크플로우

목적: 구현 전 명확한 요구사항을 보장하는 명세 주도 개발.

3단계 워크플로우:

1단계 SPEC (/do:1-plan): workflow-spec이 EARS 형식 생성. 출력물은 .do/specs/SPEC-XXX/spec.md. /clear 실행으로 45~50K 토큰 절약.

2단계 DDD (/do:2-run): ANALYZE로 요구사항 파악, PRESERVE로 기존 동작 유지, IMPROVE로 개선. 최소 85% 커버리지로 검증.

3단계 Docs (/do:3-sync): API 문서화, 아키텍처 다이어그램, 프로젝트 보고서.

EARS 형식: 시스템 전체에 항상 활성화된 요구사항은 Ubiquitous. 트리거 기반 "X 발생 시 Y 수행" 요구사항은 Event-driven. 조건부 "X 상태에서 Y 수행" 요구사항은 State-driven. 금지 "X 수행 금지" 요구사항은 Unwanted. "가능하면 X 수행" 같은 선택적 요구사항은 Optional.

토큰 예산: SPEC은 30K, DDD는 180K, Docs는 40K, 총 250K.

핵심 실천: 1단계 완료 후 /clear를 실행하여 컨텍스트를 초기화합니다.

상세 참조: modules/spec-first-ddd.md

---

### 3. 위임 패턴 - 에이전트 오케스트레이션

목적: 직접 실행을 피하고 전문 에이전트에게 작업을 위임.

핵심 원칙: Do는 모든 작업을 Task()를 통해 전문 에이전트에게 위임해야 합니다. 직접 실행은 전문화, 품질 게이트, 토큰 최적화를 우회합니다. 올바른 위임은 작업 성공률을 40% 향상시키고 병렬 실행을 가능하게 합니다.

위임 구문: 전문 에이전트를 위한 subagent_type 매개변수, 명확하고 구체적인 작업을 위한 prompt 매개변수, 관련 데이터 딕셔너리를 위한 context 매개변수로 Task를 호출합니다.

세 가지 패턴:

의존성이 있는 순차 실행: 설계를 위해 api-designer에 Task를 호출하고, 그 다음 설계 컨텍스트를 포함한 구현을 위해 backend-expert에 Task를 호출합니다.

독립적 작업의 병렬 실행: backend-expert와 frontend-expert에 동시에 Promise.all로 Task를 호출합니다.

분석 기반 조건부 실행: 분석을 위해 debug-helper에 Task를 호출하고, 분석 결과 유형에 따라 security-expert 또는 다른 적절한 에이전트에 Task를 호출합니다.

에이전트 선택: 1개 파일의 간단한 작업은 1~2개 에이전트 순차 실행. 3~5개 파일의 중간 작업은 2~3개 에이전트 순차 실행. 10개 이상 파일의 복잡한 작업은 5개 이상 에이전트 혼합 실행.

상세 참조: modules/delegation-patterns.md

---

### 4. 토큰 최적화 - 예산 관리

목적: 전략적 컨텍스트 관리를 통한 효율적인 200K 토큰 예산 운영.

예산 배분:

SPEC 단계는 30K 토큰. 전략은 요구사항만 로드하고 완료 후 /clear 실행. 명세 단계는 요구사항 분석을 위한 최소 컨텍스트만 필요. 구현 단계를 위해 45~50K 토큰을 절약합니다.

DDD 단계는 180K 토큰. 전략은 선택적 파일 로드, 구현과 관련된 파일만 로드. 구현은 심층 컨텍스트가 필요하지만 전체 코드베이스는 불필요. 예산 내에서 70% 더 큰 구현을 가능하게 합니다.

Docs 단계는 40K 토큰. 전략은 결과 캐싱 및 템플릿 재사용. 문서화는 완성된 작업 아티팩트를 기반으로 합니다. 중복 파일 읽기를 60% 줄입니다.

총 예산은 모든 단계에 걸쳐 250K 토큰. 단계 간 컨텍스트 초기화를 통한 단계 분리는 깨끗한 컨텍스트 경계를 제공하고 토큰 팽창을 방지합니다. 동일한 예산 내에서 2~3배 더 큰 프로젝트를 가능하게 합니다.

토큰 절약 전략:

단계 분리: 단계 사이에 /clear 실행, /do:1-plan 이후 45~50K 절약, 컨텍스트가 150K 초과 시, 50개 이상 메시지 후.

선택적 로드: 필요한 파일만 로드.

컨텍스트 최적화: 목표 20~30K 토큰.

모델 선택: 품질을 위해 Sonnet, 속도와 비용을 위해 Haiku(70% 저렴, 총 60~70% 절감).

상세 참조: modules/token-optimization.md

---

### 5. 점진적 공개 - 콘텐츠 아키텍처

목적: 가치와 깊이의 균형을 맞추는 3단계 지식 전달.

세 가지 레벨:

빠른 참조 레벨: 30초 투자, 핵심 원칙 및 필수 개념, 약 1,000 토큰. 시간이 부족한 사용자에게 빠른 가치 전달. 사용자는 전체 시간의 5%만으로 80%의 이해를 얻습니다.

구현 레벨: 5분 투자, 워크플로우, 실용적 예제, 통합 패턴, 약 3,000 토큰. 개념에서 실행으로의 다리 역할, 실용적 지침 제공. 깊은 전문 지식 없이도 즉시 생산적인 작업이 가능합니다.

고급 레벨: 10분 이상 투자, 깊은 기술적 심화, 엣지 케이스, 최적화 기법, 약 5,000 토큰. 복잡한 시나리오를 위한 숙련도 수준의 지식 제공. 포괄적인 커버리지를 통해 에스컬레이션을 70% 줄입니다.

SKILL.md 구조 (최대 500줄): 빠른 참조 섹션, 구현 가이드 섹션, 고급 패턴 섹션, 함께 사용하면 좋은 것들 섹션.

모듈 아키텍처: 교차 참조가 있는 진입점으로서의 SKILL.md, 무제한 크기의 심층 탐구를 위한 modules 디렉토리, 작동 예제를 위한 examples.md, 외부 링크를 위한 reference.md.

500줄 초과 시 파일 분할: SKILL.md는 빠른 참조 80~120줄, 구현 180~250줄, 고급 80~140줄, 참조 10~20줄. 초과 콘텐츠는 modules/topic.md로 이동.

상세 참조: modules/progressive-disclosure.md

---

### 6. 모듈 시스템 - 파일 구성

목적: 무제한 콘텐츠를 가능하게 하는 확장 가능한 파일 구조.

표준 구조: .claude/skills/skill-name/ 디렉토리 생성. 500줄 미만의 핵심 파일인 SKILL.md, 무제한 크기의 확장 콘텐츠를 위한 modules 디렉토리(patterns.md 포함), 작동 예제를 위한 examples.md, 외부 링크를 위한 reference.md, 유틸리티를 위한 scripts 디렉토리(선택적), templates 디렉토리(선택적) 포함.

파일 원칙: SKILL.md는 점진적 공개와 교차 참조로 500줄 미만 유지. modules 디렉토리는 주제 집중, 제한 없음, 자체 완결적 콘텐츠. examples.md는 주석이 있는 복사-붙여넣기 가능 형식. reference.md는 API 문서 및 리소스 포함.

교차 참조 구문: modules/patterns.md의 세부 내용 참조, examples.md#auth의 예제 참조, reference.md#api의 외부 문서 참조.

발견 흐름: SKILL.md → 주제 → modules/topic.md → 심층 탐구.

상세 참조: modules/modular-system.md

---

## 고급 구현

모듈 간 통합, 품질 검증, 오류 처리를 포함한 고급 패턴은 상세 모듈 참조에서 확인할 수 있습니다.

주요 고급 주제:

- 모듈 간 통합: TRUST 5 + SPEC-First DDD 결합
- 토큰 최적화 위임: 컨텍스트 초기화를 통한 병렬 실행
- 점진적 에이전트 워크플로우: 에스컬레이션 패턴
- 품질 검증: 실행 전/후 검증
- 오류 처리: 위임 실패 복구

상세 참조: examples.md(작동 코드 샘플)

---

## 함께 사용하면 좋은 것들

에이전트: 핵심 원칙으로 에이전트를 생성하는 agent-factory, 모듈 아키텍처로 스킬을 생성하는 skill-factory, 자동화된 TRUST 5 검증을 위한 core-quality, EARS 형식 명세를 위한 workflow-spec, ANALYZE-PRESERVE-IMPROVE 실행을 위한 workflow-ddd, 점진적 공개로 문서화하는 workflow-docs.

스킬: 핵심 패턴이 포함된 CLAUDE.md를 위한 do-cc-claude-md, TRUST 5가 포함된 설정을 위한 do-cc-configuration, 토큰 최적화를 위한 do-cc-memory, MCP 통합을 위한 do-context7-integration.

도구: 직접 사용자 상호작용 및 명확화를 위한 AskUserQuestion.

커맨드: SPEC-First 1단계를 위한 /do:1-plan, DDD 2단계를 위한 /do:2-run, 문서화 3단계를 위한 /do:3-sync, 지속적 개선을 위한 /do:9-feedback, 토큰 관리를 위한 /clear.

Foundation 모듈 (확장 문서): 7계층 계층의 26개 에이전트 카탈로그를 위한 modules/agents-reference.md, 6개 핵심 커맨드 워크플로우를 위한 modules/commands-reference.md, 보안, Git 전략, 규정 준수를 위한 modules/execution-rules.md.

---

## 빠른 결정 가이드

새 에이전트: 주요 원칙은 TRUST 5와 위임. 보조 원칙은 토큰 최적화와 모듈.

새 스킬: 주요 원칙은 점진적 공개와 모듈. 보조 원칙은 TRUST 5와 토큰 최적화.

워크플로우: 주요 원칙은 위임 패턴. 보조 원칙은 SPEC-First와 토큰 최적화.

품질: 주요 원칙은 TRUST 5 프레임워크. 보조 원칙은 SPEC-First DDD.

예산: 주요 원칙은 토큰 최적화. 보조 원칙은 점진적 공개와 모듈.

문서: 주요 원칙은 점진적 공개와 모듈. 보조 원칙은 토큰 최적화.

모듈 심층 탐구: modules/trust-5-framework.md, modules/spec-first-ddd.md, modules/delegation-patterns.md, modules/token-optimization.md, modules/progressive-disclosure.md, modules/modular-system.md, modules/agents-reference.md, modules/commands-reference.md, modules/execution-rules.md.

전체 예제: examples.md
외부 리소스: reference.md
