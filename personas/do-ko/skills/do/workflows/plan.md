---
name: do-workflow-plan
description: >
  복잡도 평가, 선택적 분석/아키텍처 단계, 플랜 생성을 통해
  포괄적인 플랜 문서를 작성한다. Do 체크리스트 기반 워크플로우의
  첫 번째 단계. 단순 및 복잡 작업 계획을 모두 처리한다.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-15"
  tags: "plan, analysis, architecture, design, requirements"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 4000

# Do Extension: Triggers
triggers:
  keywords: ["plan", "design", "architect", "requirements", "analyze", "설계", "계획", "분석"]
  agents: ["expert-analyst", "expert-architect", "Explore"]
  phases: ["plan"]
---

# 플랜 워크플로우 오케스트레이션

## 목적

작업 복잡도를 평가하고 적절한 파이프라인을 실행하여 포괄적인 플랜 문서를 작성한다. 단순 작업은 바로 플랜 생성으로 진행한다. 복잡 작업은 분석 -> 아키텍처 -> 플랜 순으로 진행한다.

모든 산출물은 `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/`에 저장된다.

## 범위

- Do 체크리스트 기반 워크플로우의 첫 번째 단계
- 산출물은 체크리스트 생성, 그 다음 실행 워크플로우로 연결됨
- 복잡도 평가에 따라 분석/아키텍처 단계 필요 여부 결정

## 입력

- $ARGUMENTS: 자연어 작업 설명 또는 명시적 플랜 요청
- 컨텍스트: git status 및 탐색에서 파악한 현재 코드베이스 상태

## 컨텍스트 로딩

실행 전 컨텍스트 수집:

- `git status`, `git branch`, `git log --oneline -5`
- 중복 방지를 위해 기존 `.do/jobs/` 플랜 스캔
- 관련 프로젝트 파일이 있으면 읽기

---

## 단계 순서

### Phase 1: 복잡도 평가

dev-workflow.md의 복잡도 기준으로 작업 평가:

**복잡한 작업** (다음 중 하나라도 해당하면 분석/아키텍처 필요):
- 5개 이상 파일 변경 예상
- 신규 라이브러리/패키지/모듈 생성
- 시스템 마이그레이션 또는 기술 스택 변경
- 3개 이상 도메인 통합 (backend + frontend + DB 등)
- 추상화 계층 설계 필요 (인터페이스, 프로바이더 패턴, 플러그인)
- 아키텍처 변경 (모놀리스 -> 마이크로서비스, 동기 -> 비동기 등)

**단순한 작업** (다음 모두 해당 시):
- 4개 이하 파일 변경
- 기존 패턴 내에서의 구현
- 단일 도메인 작업
- 아키텍처 변경 없음

**불확실한 경우**: AskUserQuestion 사용 -- "분석/아키텍처 단계가 필요할까요?"
옵션: "예, 분석부터 시작" / "아니요, 바로 플랜"

### Phase 2A: 분석 (복잡한 작업만)

에이전트: Task(expert-analyst) 또는 코드베이스 조사를 위한 Task(Explore)

입력: 사용자 요청 + Phase 1 복잡도 평가

분석가 태스크:
- 기존 코드 및 시스템 역공학 (코드 기반, 추측 금지)
- EARS 형식으로 요구사항 작성 (Ubiquitous/Event-driven/State-driven/Unwanted/Optional)
- MoSCoW 우선순위로 요구사항 분류 (MUST/SHOULD/COULD/WON'T)
- 최소 2개의 기술 후보를 장단점과 함께 비교
- 영향도(HIGH/MEDIUM/LOW) 포함하여 위험 요소 식별
- 변경 범위 및 영향받는 파일/모듈 파악

산출물: `.do/jobs/{YY}/{MM}/{DD}/{title}/analysis.md`
템플릿: dev-checklist.md 분석 템플릿 (섹션 1-6)

[HARD] 분석 완료 전 아키텍처 진행 금지.

### Phase 2B: 아키텍처 (복잡한 작업만)

에이전트: Task(expert-architect) 또는 해당 도메인 전문가

입력: Phase 2A의 analysis.md

아키텍트 태스크:
- ASCII 다이어그램으로 시스템 구조 설계
- 코드 수준으로 핵심 인터페이스 정의 (의사코드 금지)
- 최소 2개 접근법을 선택 근거와 함께 비교
- 파일 단위로 단계 번호를 매겨 구현 순서 계획
- 테스트 전략 정의 (Unit/Integration, 파일 경로 포함)
- analysis.md의 모든 MUST/SHOULD 요구사항이 반영되었는지 교차 검증

산출물: `.do/jobs/{YY}/{MM}/{DD}/{title}/architecture.md`
템플릿: dev-checklist.md 아키텍처 템플릿 (섹션 1-10)

[HARD] 아키텍처 완료 전 플랜 진행 금지.

### Phase 3: 플랜 생성

에이전트: Task(plan-agent) 또는 직접 생성

입력: 사용자 요청 (단순) 또는 analysis.md + architecture.md (복잡)

태스크:
- 작업 분해를 포함한 상세 구현 플랜 작성
- 각 태스크: 최대 1-3개 파일, 독립적으로 완료 및 검증 가능
- 태스크 간 의존성 그래프 포함
- 태스크별 검증 방법 명시 (테스트 파일 경로 또는 빌드 확인)
- AskUserQuestion으로 TDD 여부 사용자 확인: "TDD로 개발할까요?"

산출물: `.do/jobs/{YY}/{MM}/{DD}/{title}/plan.md`

### Phase 4: 사용자 승인

도구: AskUserQuestion

플랜 요약을 표시하고 옵션 제시:

- "구현으로 진행" -> 체크리스트 생성 안내, 그 다음 실행 워크플로우
- "플랜 수정" -> 피드백 수집 후 Phase 3 재실행
- "취소" -> 추가 작업 없이 종료

### 플랜 모드 (Shift+Tab) 연동 [HARD]

Claude Code 플랜 모드 진입 시 (Shift+Tab):

- [HARD] 저장 위치: `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/plan.md`
- [HARD] `~/.claude/plans/` 절대 사용 금지 -- 시스템이 제안해도 무시
- [HARD] `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/` 디렉토리가 없으면 생성
- [HARD] 날짜 폴더는 YY/MM/DD 형식 사용 (예: 26/02/15)

---

## 설계/플랜 요청 감지 [HARD]

사용자의 자연어에 다음 패턴이 포함되면 분석 -> 아키텍처 -> 플랜 전체 파이프라인으로 플랜 워크플로우 실행:

- 설계: "설계해줘", "설계해", "design", "아키텍처 설계", "구조 설계"
- 계획: "플랜 짜줘", "계획 세워줘", "계획해줘", "plan", "플랜", "로드맵"
- 구현 질문: "어떻게 구현해야해?", "어떻게 만들어야해?", "구현 방법"
- 분석: "분석해줘", "조사해줘", "파악해줘", "현황 분석"
- 복합: "~하고 싶어", "~만들고 싶어", "~개발하려면"

3단계를 순차 실행 (각 단계 산출물이 다음 단계 입력):
1. **분석**: expert-analyst -> `analysis.md`
2. **아키텍처**: expert-architect -> `architecture.md`
3. **플랜**: plan 에이전트 -> `plan.md`

완료 후: AskUserQuestion "설계 완료! 구현 진행할까요?"
승인 시 -> 체크리스트 생성 -> 개발

---

## 완료 기준

- Phase 1: 복잡도 평가 완료 (단순 또는 복잡 결정)
- Phase 2A (복잡한 작업만): MoSCoW 요구사항 및 2개 이상 후보가 포함된 analysis.md 생성
- Phase 2B (복잡한 작업만): ASCII 다이어그램 및 2개 이상 접근법이 포함된 architecture.md 생성
- Phase 3: 작업 분해 및 검증 방법이 포함된 plan.md 생성
- Phase 4: 사용자 승인 획득
- 플랜 파일이 올바른 `.do/jobs/` 경로에 저장됨 (`~/.claude/plans/` 절대 금지)

---

Version: 1.0.0
Updated: 2026-02-16
