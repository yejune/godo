# SPEC 워크플로우

토큰 예산 관리와 함께하는 MoAI의 3단계 개발 워크플로우입니다.

## 단계 개요

| 단계 | 명령어 | 에이전트 | 토큰 예산 | 목적 |
|-------|---------|-------|--------------|---------|
| Plan | /moai plan | manager-spec | 30K | SPEC 문서 생성 |
| Run | /moai run | manager-ddd/tdd (quality.yaml 기준) | 180K | DDD/TDD 구현 |
| Sync | /moai sync | manager-docs | 40K | 문서화 동기화 |

## Plan 단계

EARS 형식을 사용하여 포괄적인 명세서를 생성합니다.

토큰 전략:
- 할당: 30,000 토큰
- 요구사항만 로드
- 완료 후 /clear 실행
- 구현을 위해 45-50K 토큰 절약

출력:
- `.moai/specs/SPEC-XXX/spec.md`의 SPEC 문서
- EARS 형식 요구사항
- 인수 조건
- 기술적 접근 방식

## Run 단계

구성된 개발 방법론을 사용하여 명세서를 구현합니다.

토큰 전략:
- 할당: 180,000 토큰
- 선택적 파일 로딩
- 70% 더 큰 구현 가능

개발 방법론:
- quality.yaml에 구성됨 (development_mode: ddd, tdd, 또는 hybrid)
- 자세한 방법론 주기는 @workflow-modes.md 참조

성공 기준:
- 모든 SPEC 요구사항 구현됨
- 방법론별 테스트 통과
- 85%+ 코드 커버리지
- TRUST 5 품질 게이트 통과

## Sync 단계

문서화를 생성하고 배포를 준비합니다.

토큰 전략:
- 할당: 40,000 토큰
- 결과 캐싱
- 60% fewer 중복 파일 읽기

출력:
- API 문서
- 업데이트된 README
- CHANGELOG 항목
- Pull request

## 완료 마커

AI는 작업 완료를 신호하기 위해 마커를 사용합니다:
- `<moai>DONE</moai>` - 작업 완료
- `<moai>COMPLETE</moai>` - 전체 완료

## 컨텍스트 관리

/clear 전략:
- /moai plan 완료 후 (필수)
- 컨텍스트가 150K 토큰을 초과할 때
- 주요 단계 전환 전

점진적 공개:
- Level 1: 메타데이터만 (~100 토큰)
- Level 2: 트리거 시 스킬 본문 (~5000 토큰)
- Level 3: 번들 파일 요청 시

## Phase Transitions

Plan to Run:
- Trigger: SPEC document approved
- Action: Execute /clear, then /moai run SPEC-XXX

Run to Sync:
- Trigger: Implementation complete, tests passing
- Action: Execute /moai sync SPEC-XXX

## Agent Teams Variant

When team mode is enabled (workflow.team.enabled and AGENT_TEAMS env), phases can execute with Agent Teams instead of sub-agents.

### Team Mode Phase Overview

| Phase | Sub-agent Mode | Team Mode | Condition |
|-------|---------------|-----------|-----------|
| Plan | manager-spec (single) | researcher + analyst + architect (parallel) | Complexity >= threshold |
| Run | manager-ddd/tdd (sequential) | backend-dev + frontend-dev + tester (parallel) | Domains >= 3 or files >= 10 |
| Sync | manager-docs (single) | manager-docs (always sub-agent) | N/A |

### Team Mode Plan Phase
- TeamCreate for parallel research team
- Teammates explore codebase, analyze requirements, design approach
- MoAI synthesizes into SPEC document
- Shutdown team, /clear before Run phase

### Team Mode Run Phase
- TeamCreate for implementation team
- Task decomposition with file ownership boundaries
- Teammates self-claim tasks from shared list
- Quality validation after all implementation completes
- Shutdown team

### Mode Selection
- --team flag: Force team mode
- --solo flag: Force sub-agent mode
- No flag (default): Complexity-based selection
- See workflow.yaml team.auto_selection for thresholds

### Fallback
If team mode fails or is unavailable:
- Graceful fallback to sub-agent mode
- Continue from last completed task
- No data loss or state corruption
