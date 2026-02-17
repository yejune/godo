---
name: do-workflow-plan
description: >
  EARS+MoSCoW 요구사항, 복잡도 평가, 분석-아키텍처-계획 파이프라인,
  Do의 체크리스트 기반 개발 방법론을 위한 체크리스트 생성을 포함하는
  계획 워크플로우 오케스트레이터입니다.
  계획 작성, 요구사항 정의, 인수 기준 설정, 복잡도 평가,
  또는 계획 단계 오케스트레이션 시 사용하세요.
  구현(do-workflow-ddd 또는 do-workflow-tdd 사용)이나
  문서 생성(do-workflow-project 사용)에는 사용하지 마세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Write Edit Bash Grep Glob mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-16"
  modularized: "false"
  tags: "workflow, plan, ears, moscow, requirements, checklist, analysis, architecture"
  agent: "manager-plan"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# Do Extension: Triggers
triggers:
  keywords:
    - "plan"
    - "requirement"
    - "EARS"
    - "MoSCoW"
    - "acceptance criteria"
    - "planning"
    - "analysis"
    - "architecture"
    - "checklist"
    - "complexity"
  phases: ["plan"]
  agents: ["manager-plan", "manager-strategy", "expert-analyst", "expert-architect"]
---

# Do 계획 워크플로우

## 빠른 참조

계획 워크플로우 오케스트레이션 -- 복잡도 평가, 선택적 분석/아키텍처 단계, EARS+MoSCoW 요구사항, 체크리스트 생성을 통한 포괄적인 계획 수립.

핵심 기능:

- 복잡도 평가: 단순 vs 복잡 워크플로우 경로 결정
- 분석 단계: 현재 시스템 분석, 요구사항 수집, 기술 비교 (복잡한 작업만)
- 아키텍처 단계: 솔루션 설계, 인터페이스 명세, 구현 순서 (복잡한 작업만)
- EARS + MoSCoW: 우선순위가 있는 구조화된 요구사항
- 계획 문서: 단계별 구현 로드맵
- 체크리스트 생성: 3단계 템플릿이 있는 에이전트별 서브 체크리스트

산출물 위치: `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/`

---

## 구현 가이드

### 복잡도 평가 [HARD]

모든 작업은 워크플로우 경로를 결정하기 위한 복잡도 평가로 시작합니다.

단순한 작업 (모두 해당되어야 함):
- 4개 이하 파일 변경
- 기존 패턴 내에서의 작업 (새 모듈 없음)
- 단일 도메인 작업
- 아키텍처 변경 없음
- 워크플로우: Plan -> Checklist -> Develop -> Test -> Report

복잡한 작업 (하나라도 해당되면 전체 파이프라인):
- 5개 이상 파일 변경 예상
- 새 라이브러리/패키지/모듈 생성
- 시스템 마이그레이션/전환
- 3개 이상 도메인 통합 (backend + frontend + DB)
- 추상화 계층 설계 필요
- 기존 시스템 아키텍처 변경
- 워크플로우: Analysis -> Architecture -> Plan -> Checklist -> Develop -> Test -> Report

불확실한 경우: AskUserQuestion으로 사용자에게 질문: "Analysis/Architecture가 필요한가요?"

### 단순 워크플로우: Plan -> Checklist

1단계 - 요구사항 수집:
- 범위와 제약 조건을 위해 사용자 요청 파싱
- 요구사항이 모호하면 명확한 질문 [HARD]
- AskUserQuestion으로 TDD 여부 확인: "TDD로 할까요?" [HARD]

2단계 - 계획 문서:
- `.do/jobs/{YY}/{MM}/{DD}/{title}/plan.md` 생성
- 포함 내용: 목표, 범위, 접근 방식, 단계, 위험
- 적절한 경우 요구사항에 EARS 형식 사용
- 다항목 계획에 MoSCoW 우선순위 적용

3단계 - 체크리스트 생성:
- 번호가 매겨진 항목으로 `checklist.md` 생성
- 각 항목: 1~3개 파일 변경 + 검증 방법 [HARD]
- 3개 파일 초과 항목: 반드시 분할 [HARD]
- 서브 체크리스트 생성: `checklists/{order}_{agent-topic}.md`
- 서브 체크리스트는 3단계 템플릿(사전/실행/사후) 따름

### 복잡 워크플로우: Analysis -> Architecture -> Plan

1단계 - 분석 단계:
- expert-analyst 에이전트에 위임
- 산출물: `.do/jobs/{YY}/{MM}/{DD}/{title}/analysis.md`
- 내용: 현재 시스템 분석, 요구사항(MoSCoW), 기술 비교, 변경 범위, 위험
- [HARD] Architecture 시작 전 Analysis 완료 필수

2단계 - 아키텍처 단계:
- expert-architect 에이전트에 위임 (analysis.md를 입력으로 받음)
- 산출물: `.do/jobs/{YY}/{MM}/{DD}/{title}/architecture.md`
- 내용: 시스템 구조 (ASCII 다이어그램), 디렉토리 레이아웃, 핵심 인터페이스 (코드 수준), 오류 처리, 컴포넌트 구현, 접근 방식 비교 (최소 2개), 테스트 전략, 구현 순서, 위험 완화
- [HARD] Plan 시작 전 Architecture 완료 필수

3단계 - 계획 단계:
- analysis + architecture 출력을 기반으로 plan.md 생성
- analysis의 모든 MUST 요구사항이 계획에 반영됨
- architecture의 구현 순서가 체크리스트 순서를 결정함

4단계 - 체크리스트 생성:
- 단순 워크플로우와 동일하지만 아키텍처를 반영
- 항목 간 의존성에 `depends on:` 표기 사용
- 파일 소유권 경계가 아키텍처 컴포넌트 분할과 일치

### EARS + MoSCoW 통합

요구사항 작성을 위한 EARS 형식:

| 유형 | 패턴 | 예시 |
|------|---------|---------|
| Ubiquitous | 시스템은 항상 X를 한다 | 시스템은 모든 API 요청을 로그에 기록한다 |
| Event-driven | X 이벤트 발생 시 Y 수행 | 로그인 실패 시 재시도 전 3초 대기 |
| State-driven | X 조건이면 Y 수행 | 오프라인이면 로컬 캐시 사용 |
| Unwanted | 시스템은 X를 하지 않는다 | 시스템은 평문 비밀번호를 저장하지 않는다 |
| Optional | 가능하면 X를 한다 | 가능하면 이메일 알림을 전송한다 |

우선순위를 위한 MoSCoW:

- MUST: 체크리스트 항목 우선, 차단적
- SHOULD: 중요, 두 번째 우선순위
- COULD: 시간이 허락하면 좋은 것
- WON'T: 명시적으로 제외, 범위 외로 문서화

변환 흐름: EARS 요구사항 -> MoSCoW 우선순위 -> 구현 분해 -> 테스트 전략 -> 체크리스트 항목 -> 서브 체크리스트

### 테스트 전략 사전 선언 [HARD]

계획 단계에서 각 체크리스트 항목은 테스트 전략을 반드시 선언해야 합니다:

| 코드 유형 | 테스트 전략 | 예시 |
|-----------|--------------|---------|
| 비즈니스 로직, API, 데이터 계층 | 테스트 유형 + 대상 파일 | `unit: handler_test.go` |
| 복잡한 기능 | 여러 유형 | `unit: validator_test.go + E2E: flow_test.go` |
| CSS, 설정, 문서, 훅 | `pass` + 대안 | `pass (빌드 확인: go build ./...)` |

"pass"는 판단이지 건너뛰기가 아닙니다. 테스트가 불필요한 이유를 기록합니다.

---

## 함께 사용하면 좋은 것들

- do-foundation-core: 핵심 원칙 및 체크리스트 시스템
- do-workflow-ddd: 계획 이후 DDD 구현
- do-workflow-tdd: 계획 이후 TDD 구현
- do-workflow-team: 병렬 조사를 통한 팀 모드 계획
- expert-analyst: 분석 단계 실행
- expert-architect: 아키텍처 단계 실행
- manager-plan: 계획 생성 및 체크리스트 분해

---

Version: 1.0.0
Last Updated: 2026-02-16
