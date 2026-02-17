---
name: do-workflow-do
description: >
  완전 자율 plan -> checklist -> run -> test -> report 파이프라인.
  서브커맨드가 지정되지 않았을 때의 기본 워크플로우. 탐색,
  플랜 생성, 체크리스트 생성, 구현, 완료를 처리한다.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-15"
  tags: "do, autonomous, pipeline, default"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 3000

# Do Extension: Triggers
triggers:
  keywords: ["do", "autonomous", "pipeline", "build", "implement", "create"]
  agents: ["do"]
  phases: ["plan", "run", "test", "report"]
---

# 워크플로우: Do - 자율 개발 파이프라인

목적: 완전 자율 워크플로우. 사용자가 목표를 제시하면 Do가 plan -> checklist -> run -> test -> report 파이프라인을 자율로 실행한다. 서브커맨드가 지정되지 않았을 때의 기본 워크플로우.

흐름: 탐색 -> 플랜 -> 체크리스트 -> 실행 -> 테스트 -> 보고 -> 완료

## 지원 플래그

- --team: 적용 가능한 모든 단계에 Team 모드 강제 적용 (team-do.md로 라우팅)
- --solo: Do 모드 강제 적용 (단계당 에이전트 단일)

**기본 동작 (플래그 없음)**: 현재 DO_MODE 설정을 사용한다.

## Phase 0: 탐색

탐색 에이전트 실행 (Do 모드에서는 병렬, Focus 모드에서는 직접 실행):

에이전트 1 - 탐색 (Explore 서브에이전트):
- 작업 컨텍스트를 위한 코드베이스 분석
- 관련 파일, 아키텍처 패턴, 기존 구현체 파악

에이전트 2 - 조사 (넓은 범위의 Explore 서브에이전트):
- 필요 시 외부 문서 및 모범 사례 수집
- 코드베이스 내 유사 패턴 파악

탐색 완료 후:
- 결과를 통합 컨텍스트로 종합
- 복잡도 판단 (단순 vs 복잡 작업)

--team 플래그 지정 시: workflows/team-do.md로 라우팅.

## Phase 0 완료: 라우팅 결정

- 단일 도메인 단순 작업: 전문가 에이전트에 직접 위임 고려 (전체 파이프라인 생략)
- 멀티 도메인 또는 복잡 작업: 전체 파이프라인 진행

AskUserQuestion으로 사용자 승인:
- "플랜 생성으로 진행"
- "전문가 에이전트에 직접 위임" (단순 작업의 경우)
- "취소"

## Phase 1: 플랜

workflows/plan.md 실행:
- 복잡도 평가 -> 분석 (복잡한 경우) -> 아키텍처 (복잡한 경우) -> 플랜 생성
- 산출물: `.do/jobs/{YY}/{MM}/{DD}/{title}/plan.md` (복잡한 경우 + analysis.md, architecture.md)
- 진행 전 사용자 승인 체크포인트

## Phase 2: 체크리스트

플랜 승인 후:

1. plan.md에서 checklist.md 생성:
   - 각 플랜 태스크가 체크리스트 항목이 됨
   - 항목은 각각 1-3개 파일 변경 단위로 분해
   - 항목별 검증 방법 명시
   - `depends on:`으로 의존성 연결

2. 에이전트별 서브 체크리스트 생성:
   - `.do/jobs/{YY}/{MM}/{DD}/{title}/checklists/{NN}_{agent-topic}.md`
   - 각각 dev-checklist.md의 서브 체크리스트 템플릿 준수
   - 문제 요약, 인수 기준, 핵심 파일, FINAL STEP: Commit 포함

3. TDD 여부 사용자 확인 (플랜 단계에서 미결정 시):
   - AskUserQuestion: "TDD로 개발할까요?"
   - "예, TDD" -> 실행에 테스트 워크플로우 통합
   - "아니요, 먼저 구현" -> 표준 실행 후 검증

## Phase 3: 실행

workflows/run.md 실행:
- 체크리스트 읽기, 모드별 에이전트 디스패치 (Do/Focus/Team)
- 진행 상황 모니터링, 중단 처리
- 모든 에이전트는 체크리스트 갱신과 함께 작업 커밋

## Phase 4: 테스트

workflows/test.md 실행 (실행 단계에서 TDD로 완료되지 않은 경우):
- 전체 테스트 스위트 실행
- 커버리지 목표 검증 (85%+)
- 진행 전 실패 수정

## Phase 5: 보고

workflows/report.md 실행:
- 모든 서브 체크리스트 결과 집계
- 완료 보고서 생성
- 사용자에게 요약 및 다음 단계 제시

## 실행 요약

1. 인수 파싱 (플래그 추출: --team, --solo)
2. DO_MODE 확인 및 실행 전략 결정
3. --team 지정 시: workflows/team-do.md로 라우팅
4. Phase 0 실행 (탐색)
5. 라우팅 결정 (단순 직접 위임 vs 전체 파이프라인)
6. AskUserQuestion으로 사용자 확인
7. Phase 1 (플랜): workflows/plan.md 읽기
8. Phase 2 (체크리스트): checklist.md + 서브 체크리스트 생성
9. Phase 3 (실행): workflows/run.md 읽기
10. Phase 4 (테스트): workflows/test.md 읽기 (TDD가 아닌 경우)
11. Phase 5 (보고): workflows/report.md 읽기
12. 사용자에게 최종 요약 표시

---

Version: 1.0.0
Updated: 2026-02-16
