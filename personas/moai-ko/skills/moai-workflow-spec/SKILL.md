---
name: moai-workflow-spec
description: >
  EARS 형식 요구사항, 인수 기준, MoAI-ADK 개발 방법론을 위한
  Plan-Run-Sync 통합을 갖춘 SPEC 워크플로우 오케스트레이션입니다.
  SPEC 문서 생성, EARS 요구사항 작성, 인수 기준 정의,
  기능 계획, 또는 /moai plan 단계 오케스트레이션 시 사용합니다.
  구현에는 사용하지 마세요 (moai-workflow-ddd 사용).
  문서 생성에는 사용하지 마세요 (moai-workflow-project 사용).
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Write Edit Bash(git:*) Bash(ls:*) Bash(wc:*) Bash(mkdir:*) Grep Glob mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "1.2.0"
  category: "workflow"
  status: "active"
  updated: "2026-01-08"
  modularized: "true"
  tags: "workflow, spec, ears, requirements, moai-adk, planning"
  author: "MoAI-ADK Team"
  context: "fork"
  agent: "Plan"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords:
    [
      "SPEC",
      "requirement",
      "EARS",
      "acceptance criteria",
      "user story",
      "planning",
      "specification",
      "requirements gathering",
    ]
  phases: ["plan"]
  agents: ["manager-spec", "manager-strategy", "Plan"]
---

# SPEC 워크플로우 관리

## 빠른 참조 (30초)

SPEC 워크플로우 오케스트레이션 - 체계적인 요구사항 정의 및 Plan-Run-Sync 워크플로우 통합을 위한 EARS 형식을 사용한 포괄적인 명세서 관리.

핵심 기능:

- EARS 형식 명세서: 명확성을 위한 다섯 가지 요구사항 패턴
- 요구사항 명확화: 4단계 체계적 프로세스
- SPEC 문서 템플릿: 일관성을 위한 표준화된 구조
- Plan-Run-Sync 통합: 원활한 워크플로우 연결
- 병렬 개발: Git Worktree 기반 SPEC 격리
- 품질 게이트: TRUST 5 프레임워크 검증

EARS 다섯 가지 패턴:

- Ubiquitous: 시스템은 항상 조치를 수행해야 함 - 항상 활성화
- Event-Driven: 이벤트 발생 시 조치 실행 - 트리거-응답
- State-Driven: 조건이 참이면 조치 실행 - 조건부 동작
- Unwanted: 시스템은 조치를 수행하지 않아야 함 - 금지
- Optional: 가능한 경우 기능을 제공함 - 있으면 좋음

사용 시기:

- 기능 계획 및 요구사항 정의
- SPEC 문서 생성 및 유지 관리
- 병렬 기능 개발 조율
- 품질 보증 및 검증 계획

빠른 명령어:

- 새 SPEC 생성: /moai:1-plan "user authentication system"
- Worktrees로 병렬 SPEC 생성: /moai:1-plan "login feature" "signup feature" --worktree
- 새 브랜치로 SPEC 생성: /moai:1-plan "payment processing" --branch
- 기존 SPEC 업데이트: /moai:1-plan SPEC-001 "add OAuth support"

---

## 구현 가이드 (5분)

### 핵심 개념

SPEC 우선 개발 철학:

- EARS 형식은 명확한 요구사항을 보장
- 요구사항 명확화는 범위 확장을 방지
- 테스트 시나리오를 통한 체계적 검증
- 구현을 위한 DDD 워크플로우와의 통합
- 품질 게이트는 완료 기준을 적용
- Constitution 참조는 프로젝트 전반의 일관성을 보장

### Constitution 참조 (SDD 2025 표준)

Constitution은 모든 SPEC이 준수해야 하는 프로젝트 DNA를 정의합니다. SPEC을 생성하기 전에 `.moai/project/tech.md`에 정의된 프로젝트 constitution과의 일치를 확인합니다.

Constitution 구성 요소:

- 기술 스택: 필수 버전 및 프레임워크
- 명명 규칙: 변수, 함수, 파일 명명 기준
- 금지 라이브러리: 대안이 있는 명시적으로 금지된 라이브러리
- 아키텍처 패턴: 레이어링 규칙 및 의존성 방향
- 보안 기준: 인증 패턴 및 암호화 요구사항
- 로깅 기준: 로그 형식 및 구조화된 로깅 요구사항

Constitution 검증:

- 모든 SPEC 기술 선택은 Constitution 스택 버전과 일치해야 함
- SPEC은 금지된 라이브러리나 패턴을 도입할 수 없음
- SPEC은 Constitution에 정의된 명명 규칙을 따라야 함
- SPEC은 아키텍처 경계 및 레이어링을 존중해야 함

이유: Constitution은 아키텍처 드리프트를 방지하고 유지 관리성을 보장
영향: Constitution과 일치하는 SPEC은 통합 충돌을 상당히 줄임

### SPEC 워크플로우 단계

1단계 - 사용자 입력 분석: 자연어 기능 설명 파싱
2단계 - 요구사항 명확화: 4단계 체계적 프로세스
3단계 - EARS 패턴 적용: 다섯 가지 패턴으로 요구사항 구조화
4단계 - 성공 기준 정의: 완료 지표 수립
5단계 - 테스트 시나리오 생성: 검증 테스트 케이스 생성
6단계 - SPEC 문서 생성: 표준화된 마크다운 출력 생성

### EARS 형식 심층 분석

Ubiquitous 요구사항 - 항상 활성화:

- 사용 사례: 시스템 전반의 품질 속성
- 예시: 로깅, 입력 검증, 에러 처리
- 테스트 전략: 모든 기능 테스트 스위트에 공통 검증으로 포함

Event-Driven 요구사항 - 트리거-응답:

- 사용 사례: 사용자 상호작용 및 시스템 간 통신
- 예시: 버튼 클릭, 파일 업로드, 결제 완료
- 테스트 전략: 예상 응답 검증을 통한 이벤트 시뮬레이션

State-Driven 요구사항 - 조건부 동작:

- 사용 사례: 접근 제어, 상태 머신, 조건부 비즈니스 로직
- 예시: 계정 상태 확인, 재고 검증, 권한 확인
- 테스트 전략: 조건부 동작 검증을 통한 상태 설정

Unwanted 요구사항 - 금지된 조치:

- 사용 사례: 보안 취약점, 데이터 무결성 보호
- 예시: 평문 비밀번호 없음, 무단 접근 없음, 로그에 개인 정보 없음
- 테스트 전략: 금지된 동작 검증을 갖춘 부정 테스트 케이스

Optional 요구사항 - 향상 기능:

- 사용 사례: MVP 범위 정의, 기능 우선순위
- 예시: OAuth 로그인, 다크 모드, 오프라인 모드
- 테스트 전략: 구현 상태에 따른 조건부 테스트 실행

### 요구사항 명확화 프로세스

0단계 - 가정 분석 (Philosopher 프레임워크):

범위를 정의하기 전에 AskUserQuestion을 사용하여 기본 가정을 표면화하고 검증합니다.

가정 카테고리:

- 기술적 가정: 기술 역량, API 가용성, 성능 특성
- 비즈니스 가정: 사용자 동작, 시장 요구사항, 타임라인 실현 가능성
- 팀 가정: 기술 가용성, 리소스 할당, 지식 격차
- 통합 가정: 서드파티 서비스 신뢰성, 호환성 기대치

가정 문서화:

- 가정 진술: 가정하는 내용의 명확한 설명
- 신뢰 수준: 증거에 기반한 높음, 중간, 낮음
- 증거 근거: 이 가정을 지지하는 것
- 틀렸을 경우 위험: 가정이 거짓으로 판명될 경우의 결과
- 검증 방법: 상당한 노력을 투입하기 전에 확인하는 방법

0.5단계 - 근본 원인 분석:

기능 요청 또는 문제 중심 SPEC의 경우, 5가지 이유를 적용합니다:

- 표면 문제: 사용자가 무엇을 관찰하거나 요청하는가?
- 첫 번째 이유: 이 요청을 유발하는 즉각적인 필요는?
- 두 번째 이유: 그 필요를 만드는 근본적인 문제는?
- 세 번째 이유: 어떤 시스템적 요인이 기여하는가?
- 근본 원인: 솔루션이 해결해야 하는 근본적인 이슈는?

1단계 - 범위 정의:

- 지원되는 인증 방법 식별
- 검증 규칙 및 제약 정의
- 실패 처리 전략 결정
- 세션 관리 접근 방식 수립

2단계 - 제약 추출:

- 성능 요구사항: 응답 시간 목표
- 보안 요구사항: OWASP 준수, 암호화 기준
- 호환성 요구사항: 지원되는 브라우저 및 기기
- 확장성 요구사항: 동시 사용자 목표

3단계 - 성공 기준 정의:

- 테스트 커버리지: 최소 비율 목표
- 응답 시간: 백분위수 목표 (P50, P95, P99)
- 기능 완료: 모든 정상 시나리오 검증 통과
- 품질 게이트: 린터 경고 없음, 보안 취약점 없음

4단계 - 테스트 시나리오 생성:

- 정상 케이스: 예상 출력이 있는 유효한 입력
- 에러 케이스: 에러 처리가 있는 유효하지 않은 입력
- 엣지 케이스: 경계 조건 및 코너 케이스
- 보안 케이스: 인젝션 공격, 권한 상승 시도

### Plan-Run-Sync 워크플로우 통합

PLAN 단계 (/moai:1-plan):

- manager-spec 에이전트가 사용자 입력 분석
- EARS 형식 요구사항 생성
- 사용자 상호작용을 통한 요구사항 명확화
- .moai/specs/ 디렉토리에 SPEC 문서 생성
- Git 브랜치 생성 (선택적 --branch 플래그)
- Git Worktree 설정 (선택적 --worktree 플래그)

RUN 단계 (/moai:2-run):

- manager-ddd 에이전트가 SPEC 문서 로드
- ANALYZE-PRESERVE-IMPROVE DDD 사이클 실행
- 테스트 패턴에 대한 moai-workflow-testing 스킬 참조
- 도메인 전문가 에이전트 위임 (expert-backend, expert-frontend 등)
- manager-quality 에이전트를 통한 품질 검증

SYNC 단계 (/moai:3-sync):

- manager-docs 에이전트가 문서 동기화
- SPEC에서 API 문서 생성
- README 및 아키텍처 문서 업데이트
- CHANGELOG 항목 생성
- SPEC 참조가 있는 버전 제어 커밋

### Git Worktree를 통한 병렬 개발

Worktree 개념:

- 여러 브랜치를 위한 독립적인 작업 디렉토리
- 각 SPEC이 격리된 개발 환경을 가짐
- 병렬 작업을 위한 브랜치 전환 불필요
- 기능 격리를 통한 병합 충돌 감소

Worktree 생성:

- /moai:1-plan "login feature" "signup feature" --worktree 명령어로 여러 SPEC 생성
- 결과로 SPEC별 하위 디렉토리가 있는 project-worktrees 디렉토리 생성

Worktree 장점:

- 병렬 개발: 여러 기능을 동시에 개발
- 팀 협업: SPEC별 명확한 소유권 경계
- 의존성 격리: 기능별 다른 라이브러리 버전
- 위험 감소: 불안정한 코드가 다른 기능에 영향을 미치지 않음

---

## 고급 구현 (10분 이상)

SPEC 템플릿, 검증 자동화, 워크플로우 최적화를 포함한 고급 패턴은 다음을 참조합니다:

- [고급 패턴](modules/advanced-patterns.md): 커스텀 SPEC 템플릿, 검증 자동화
- [참조 가이드](reference.md): SPEC 메타데이터 스키마, 통합 예제
- [예제](examples.md): 실제 SPEC 문서, 워크플로우 시나리오

## 리소스

### SPEC 파일 구성

디렉토리 구조 (표준 3파일 형식):

- .moai/specs/SPEC-{ID}/: 3개의 필수 파일을 포함하는 SPEC 문서 디렉토리
  - spec.md: EARS 형식 명세서 (Environment, Assumptions, Requirements, Specifications)
  - plan.md: 구현 계획, 마일스톤, 기술적 접근 방식
  - acceptance.md: 상세 인수 기준, 테스트 시나리오 (Given-When-Then 형식)
- .moai/memory/: 세션 상태 파일 (last-session-state.json)
- .moai/docs/: 생성된 문서 (api-documentation.md)

[HARD] 필수 파일 세트:
모든 SPEC 디렉토리는 3개의 파일(spec.md, plan.md, acceptance.md)을 모두 포함해야 함
이유: 완전한 SPEC 구조는 추적 가능성, 구현 안내, 품질 검증을 보장
영향: 누락된 파일은 불완전한 요구사항을 만들고 적절한 워크플로우 실행을 방해

### SPEC 메타데이터 스키마

필수 필드:

- SPEC ID: 순차 번호 (SPEC-001, SPEC-002 등)
- 제목: 영어로 된 기능 이름
- 생성일: ISO 8601 타임스탬프
- 상태: Planned, In Progress, Completed, Blocked
- 우선순위: High, Medium, Low
- 담당자: 구현 책임 에이전트

선택적 필드:

- 관련 SPEC: 의존성 및 관련 기능
- Epic: 부모 기능 그룹
- 예상 노력: 시간 또는 스토리 포인트 단위 시간 추정
- 레이블: 분류를 위한 태그

### SPEC 라이프사이클 관리 (SDD 2025 표준)

라이프사이클 레벨 필드:

레벨 1 - spec-first:

- 설명: 구현 전에 작성되고 완료 후 폐기되는 SPEC
- 사용 사례: 일회성 기능, 프로토타입, 실험
- 유지 관리 정책: 구현 후 유지 관리 불필요

레벨 2 - spec-anchored:

- 설명: 진화를 위해 구현과 함께 유지 관리되는 SPEC
- 사용 사례: 핵심 기능, API 계약, 통합 지점
- 유지 관리 정책: 분기별 검토, 구현 변경 시 업데이트

레벨 3 - spec-as-source:

- 설명: SPEC이 진실의 단일 출처; 인간은 SPEC만 편집
- 사용 사례: 중요 시스템, 규제 환경, 코드 생성 워크플로우
- 유지 관리 정책: SPEC 변경으로 구현 재생성 트리거

라이프사이클 전환 규칙:

- spec-first에서 spec-anchored로: 기능이 프로덕션 핵심이 될 때
- spec-anchored에서 spec-as-source로: 준수 또는 재생성 워크플로우 필요 시
- 다운그레이드 허용되지만 SPEC 히스토리에 명시적인 정당성 필요

### 품질 지표

SPEC 품질 지표:

- 요구사항 명확성: 모든 EARS 패턴이 적절하게 사용됨
- 테스트 커버리지: 모든 요구사항에 해당하는 테스트 시나리오가 있음
- 제약 완전성: 기술적 및 비즈니스 제약이 정의됨
- 성공 기준 측정 가능성: 정량화 가능한 완료 지표

검증 체크리스트:

- 모든 EARS 요구사항이 테스트 가능
- 모호한 언어 없음 (should, might, usually)
- 모든 에러 케이스가 문서화됨
- 성능 목표가 정량화됨
- 보안 요구사항이 OWASP 준수

### 잘 어울리는 것들

- moai-foundation-core: SPEC 우선 DDD 방법론 및 TRUST 5 프레임워크
- moai-workflow-testing: DDD 구현 및 테스트 자동화
- moai-workflow-project: 프로젝트 초기화 및 설정
- moai-workflow-worktree: 병렬 개발을 위한 Git Worktree 관리
- manager-spec: SPEC 생성 및 요구사항 분석 에이전트
- manager-ddd: SPEC 요구사항 기반 DDD 구현
- manager-quality: TRUST 5 품질 검증 및 게이트 적용

### 통합 예제

순차 워크플로우:

- 1단계 PLAN: /moai:1-plan "user authentication system"
- 2단계 RUN: /moai:2-run SPEC-001
- 3단계 SYNC: /moai:3-sync SPEC-001

병렬 워크플로우:

- 여러 SPEC 생성: /moai:1-plan "backend API" "frontend UI" "database schema" --worktree
- 세션 1: /moai:2-run SPEC-001 (backend API)
- 세션 2: /moai:2-run SPEC-002 (frontend UI)
- 세션 3: /moai:2-run SPEC-003 (database schema)

### 토큰 관리

세션 전략:

- PLAN 단계는 세션 토큰의 약 30% 사용
- RUN 단계는 세션 토큰의 약 60% 사용
- SYNC 단계는 세션 토큰의 약 10% 사용

컨텍스트 최적화:

- SPEC 문서는 .moai/specs/ 디렉토리에 지속
- 세션 간 컨텍스트를 위한 .moai/memory/의 세션 메모리
- SPEC ID 참조를 통한 최소 컨텍스트 전달
- 토큰 오버헤드를 줄이는 에이전트 위임

---

## SPEC 범위 및 분류

### .moai/specs/에 속하는 것

`.moai/specs/` 디렉토리는 구현할 기능을 정의하는 SPEC 문서만을 위한 것입니다.

유효한 SPEC 콘텐츠:

- EARS 형식의 기능 요구사항
- 마일스톤이 있는 구현 계획
- Given/When/Then 시나리오를 갖춘 인수 기준
- 새로운 기능을 위한 기술 명세서
- 명확한 결과물이 있는 사용자 스토리

SPEC 특성:

- 미래 지향적: 구축될 것을 설명
- 실행 가능: 구현 안내 포함
- 테스트 가능: 인수 기준 포함
- 구조화됨: EARS 형식 패턴 사용

### .moai/specs/에 속하지 않는 것

| 문서 유형 | SPEC이 아닌 이유 | 올바른 위치 |
| --------------------- | ----------------------------- | ----------------------------------------- |
| 보안 감사 | 기존 코드 분석 | `.moai/reports/security-audit-{DATE}/` |
| 성능 보고서 | 현재 지표 문서화 | `.moai/reports/performance-{DATE}/` |
| 의존성 분석 | 기존 의존성 검토 | `.moai/reports/dependency-review-{DATE}/` |
| 아키텍처 개요 | 현재 상태 문서화 | `.moai/docs/architecture.md` |
| API 참조 | 기존 API 문서화 | `.moai/docs/api-reference.md` |
| 회의 노트 | 내린 결정 기록 | `.moai/reports/meeting-{DATE}/` |
| 회고 | 과거 작업 분석 | `.moai/reports/retro-{DATE}/` |

### 제외 규칙

[HARD] 보고서 vs SPEC 구분:

보고서는 존재하는 것을 분석 → `.moai/reports/`
SPEC은 구축될 것을 정의 → `.moai/specs/`

[HARD] 문서 vs SPEC 구분:

문서는 사용 방법을 설명 → `.moai/docs/`
SPEC은 구축할 것을 정의 → `.moai/specs/`

---

## 레거시 파일 마이그레이션 가이드

### 시나리오 1: 플랫 SPEC 파일 → 디렉토리 변환

문제: `.moai/specs/SPEC-AUTH-001.md`가 단일 파일로 존재

해결 단계:

1. 디렉토리 생성: `mkdir -p .moai/specs/SPEC-AUTH-001/`
2. 콘텐츠 이동: `mv .moai/specs/SPEC-AUTH-001.md .moai/specs/SPEC-AUTH-001/spec.md`
3. 누락된 파일 생성:
   - 구현 계획 추출 → `plan.md`
   - 인수 기준 추출 → `acceptance.md`
4. 구조 확인: 3개의 파일 모두 존재
5. 커밋: `git add . && git commit -m "refactor(spec): Convert SPEC-AUTH-001 to directory structure"`

검증 명령어:

```bash
# 플랫 SPEC 파일 확인 (빈 결과여야 함)
find .moai/specs -maxdepth 1 -name "SPEC-*.md" -type f
```

### 시나리오 2: 번호 없는 SPEC ID → 번호 할당

문제: 번호 없는 `SPEC-REDESIGN` 또는 `SPEC-SDK-INTEGRATION`

해결 단계:

1. 다음 사용 가능한 번호 찾기:
   ```bash
   ls -d .moai/specs/SPEC-*-[0-9][0-9][0-9] 2>/dev/null | sort -t- -k3 -n | tail -1
   ```
2. 번호 할당: `SPEC-REDESIGN` → `SPEC-REDESIGN-001`
3. 디렉토리 이름 변경:
   ```bash
   mv .moai/specs/SPEC-REDESIGN .moai/specs/SPEC-REDESIGN-001
   ```
4. spec.md frontmatter의 내부 참조 업데이트
5. 커밋: `git commit -m "refactor(spec): Assign number to SPEC-REDESIGN → SPEC-REDESIGN-001"`

### 시나리오 3: SPEC 디렉토리의 보고서 → 분리

문제: `.moai/specs/`에 분석/감사 문서

해결 단계:

1. 콘텐츠에서 문서 유형 식별
2. 보고서 디렉토리 생성:
   ```bash
   mkdir -p .moai/reports/security-audit-2025-01/
   ```
3. 콘텐츠 이동:
   ```bash
   mv .moai/specs/SPEC-SECURITY-AUDIT/* .moai/reports/security-audit-2025-01/
   rmdir .moai/specs/SPEC-SECURITY-AUDIT
   ```
4. 필요한 경우 메인 파일을 report.md로 이름 변경
5. 커밋: `git commit -m "refactor: Move security audit from specs to reports"`

### 시나리오 4: 중복 SPEC ID → 해결

문제: 동일한 SPEC ID를 가진 두 개의 디렉토리

해결 단계:

1. 생성 날짜 비교:
   ```bash
   ls -la .moai/specs/ | grep SPEC-AUTH-001
   ```
2. 어떤 것이 정본인지 결정 (보통 오래된 것)
3. 새것을 다음 사용 가능한 번호로 다시 번호 매기기:
   ```bash
   mv .moai/specs/SPEC-AUTH-001-duplicate .moai/specs/SPEC-AUTH-002
   ```
4. 내부 참조 업데이트
5. 커밋: `git commit -m "fix(spec): Resolve duplicate SPEC-AUTH-001 → SPEC-AUTH-002"`

### 검증 스크립트

SPEC 구성 이슈를 식별하기 위해 이 스크립트를 실행합니다:

```bash
#!/bin/bash
# SPEC 구성 검증기

echo "=== SPEC 구성 확인 ==="

# 확인 1: specs 루트의 플랫 파일
echo -e "\n[확인 1] 플랫 SPEC 파일 (비어있어야 함):"
find .moai/specs -maxdepth 1 -name "SPEC-*.md" -type f

# 확인 2: 필수 파일이 없는 디렉토리
echo -e "\n[확인 2] 필수 파일이 없는 SPEC 디렉토리:"
for dir in .moai/specs/SPEC-*/; do
  if [ -d "$dir" ]; then
    missing=""
    [ ! -f "${dir}spec.md" ] && missing="${missing}spec.md "
    [ ! -f "${dir}plan.md" ] && missing="${missing}plan.md "
    [ ! -f "${dir}acceptance.md" ] && missing="${missing}acceptance.md "
    [ -n "$missing" ] && echo "$dir: 누락 $missing"
  fi
done

# 확인 3: 번호 없는 SPEC
echo -e "\n[확인 3] 적절한 번호 없는 SPEC:"
ls -d .moai/specs/SPEC-*/ 2>/dev/null | grep -v -E 'SPEC-[A-Z]+-[0-9]{3}'

# 확인 4: specs에 있을 수 있는 보고서
echo -e "\n[확인 4] specs에 있을 수 있는 보고서 (수동 확인):"
grep -l -r "findings\|recommendations\|audit\|analysis" .moai/specs/*/spec.md 2>/dev/null

echo -e "\n=== 확인 완료 ==="
```

---

Version: 1.3.0 (SDD 2025 표준 통합 + SPEC 범위 분류)
Last Updated: 2026-01-21
Integration Status: 완료 - SDD 2025 기능 및 마이그레이션 가이드를 갖춘 완전한 Plan-Run-Sync 워크플로우
