---
name: moai-workflow-moai
description: >
  완전 자율 plan-run-sync 파이프라인입니다. 서브커맨드가 지정되지 않을 때의 기본
  워크플로우입니다. 병렬 탐색, SPEC 생성, 선택적 자동 수정 루프를 포함한 DDD/TDD
  구현, 문서화 동기화를 처리합니다.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "2.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-07"
  tags: "moai, autonomous, pipeline, plan-run-sync, default"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["moai", "autonomous", "pipeline", "build", "implement", "create"]
  agents: ["moai"]
  phases: ["plan", "run", "sync"]
---

# 워크플로우: MoAI - 자율 개발 오케스트레이션

목적: 완전 자율 워크플로우입니다. 사용자가 목표를 제공하면 MoAI가 plan -> run -> sync 파이프라인을 자율적으로 실행합니다. 서브커맨드가 지정되지 않을 때의 기본 워크플로우입니다.

흐름: 탐색 -> 계획 -> 실행 -> 동기화 -> 완료

## 지원 플래그

- --loop: 실행 단계 중 자동 반복 수정 활성화
- --max N: 루프의 최대 반복 횟수 (기본값 100)
- --branch: 기능 브랜치 자동 생성
- --pr: 완료 후 풀 리퀘스트 자동 생성
- --resume SPEC-XXX: 기존 SPEC에서 이전 작업 재개
- --team: plan 및 run 단계에서 Agent Teams 모드 강제 적용
- --solo: 서브에이전트 모드 강제 적용 (단계별 단일 에이전트)

**기본 동작 (플래그 없음)**: 복잡도에 따라 시스템이 자동 선택:
- 팀 모드: 멀티 도메인 작업 (도메인 3개 이상), 많은 파일 (10개 이상), 또는 높은 복잡도 (7점 이상)
- 서브에이전트 모드: 집중적인 단일 도메인 작업

## 설정 파일

- quality.yaml: TRUST 5 품질 임계값 AND development_mode 라우팅
- workflow.yaml: 실행 모드, 팀 설정, 루프 방지, 완료 마커

## 개발 모드 라우팅 (CRITICAL)

[HARD] Phase 2 구현 전, 반드시 `.moai/config/sections/quality.yaml`을 확인하세요:

```yaml
constitution:
  development_mode: hybrid    # ddd, tdd, 또는 hybrid
  hybrid_settings:
    new_features: tdd        # 새 코드는 TDD 사용
    legacy_refactoring: ddd  # 기존 코드는 DDD 사용
```

**라우팅 로직**:

| 기능 유형 | 모드: ddd | 모드: tdd | 모드: hybrid |
|--------------|-----------|-----------|--------------|
| **새 패키지/모듈** (기존 파일 없음) | DDD* | TDD | TDD |
| **기존 파일에 새 기능 추가** | DDD | TDD | TDD |
| **기존 코드 리팩토링** | DDD | 이 부분은 DDD 사용 | DDD |
| **기존 코드 버그 수정** | DDD | TDD | DDD |

*DDD는 그린필드에 맞게 조정됩니다 (요구사항 분석 → 스펙 테스트로 보존 → 개선)

**에이전트 선택**:
- **TDD 사이클**: `manager-tdd` 서브에이전트 (RED-GREEN-REFACTOR)
- **DDD 사이클**: `manager-ddd` 서브에이전트 (ANALYZE-PRESERVE-IMPROVE)

## Phase 0: 병렬 탐색

단일 응답에서 세 에이전트를 동시에 실행하여 2-3배 속도 향상 (15-30초 vs 45-90초).

에이전트 1 - 탐색 (subagent_type Explore):
- 작업 컨텍스트를 위한 코드베이스 분석
- 관련 파일, 아키텍처 패턴, 기존 구현

에이전트 2 - 리서치 (WebSearch/WebFetch에 초점을 맞춘 subagent_type Explore):
- 외부 문서 및 모범 사례
- API 문서, 라이브러리 문서, 유사 구현

에이전트 3 - 품질 (subagent_type manager-quality):
- 현재 프로젝트 품질 평가
- 테스트 커버리지 상태, 린트 상태, 기술 부채

모든 에이전트 완료 후:
- 각 에이전트 응답에서 출력 수집
- 탐색 (파일, 패턴), 리서치 (외부 지식), 품질 (커버리지 기준선)에서 핵심 발견사항 추출
- 통합 탐색 보고서로 종합
- 생성/수정할 파일 및 테스트 전략을 포함한 실행 계획 생성

오류 처리: 에이전트 중 하나라도 실패하면 성공한 에이전트의 결과로 계속 진행합니다. 계획에 누락된 정보를 기록합니다.

--sequential 플래그 사용 시: 탐색, 리서치, 품질을 순서대로 실행합니다.

## Phase 0 완료: 라우팅 결정

단일 도메인 라우팅:
- 작업이 단일 도메인인 경우 (예: "SQL 최적화"): 전문 에이전트에 직접 위임, SPEC 생성 건너뜀
- 작업이 멀티 도메인인 경우: SPEC 생성을 포함한 전체 워크플로우 진행

AskUserQuestion을 통한 사용자 승인 체크포인트:
- 옵션: SPEC 생성 진행, 접근 방식 수정, 취소

## Phase 1: SPEC 생성

- manager-spec 서브에이전트에 위임
- 출력: .moai/specs/SPEC-XXX/spec.md에 EARS 형식 SPEC 문서
- 요구사항, 인수 기준, 기술 접근 방식 포함

## Phase 2: 구현 (development_mode에 따라 TDD 또는 DDD)

[HARD] 에이전트 위임 의무: 모든 구현 작업은 반드시 전문 에이전트에게 위임해야 합니다. 자동 컴팩트 이후에도 직접 구현을 실행하지 마세요.

[HARD] `.moai/config/sections/quality.yaml`에 따른 방법론 선택:

- **새 기능** (hybrid_settings.new_features에 따라): `manager-tdd` 사용 (RED-GREEN-REFACTOR)
- **레거시 리팩토링** (hybrid_settings.legacy_refactoring에 따라): `manager-ddd` 사용 (ANALYZE-PRESERVE-IMPROVE)

전문 에이전트 선택 (도메인별 작업):
- 백엔드 로직: expert-backend 서브에이전트
- 프론트엔드 컴포넌트: expert-frontend 서브에이전트
- 테스트 생성: expert-testing 서브에이전트
- 버그 수정: expert-debug 서브에이전트
- 리팩토링: expert-refactoring 서브에이전트
- 보안 수정: expert-security 서브에이전트

루프 동작 (--loop 플래그 또는 workflow.yaml loop_prevention 설정 활성화 시):
- 이슈가 존재하고 반복 횟수가 최대치 미만인 동안:
  - 진단 실행 (기본적으로 병렬)
  - 적절한 전문 에이전트에 수정 위임
  - 수정 결과 검증
  - 완료 마커 확인
  - 마커 발견 시: 루프 종료

## Phase 3: 문서화 동기화

- manager-docs 서브에이전트에 위임
- 구현과 문서화 동기화
- SPEC-구현 불일치 감지 및 SPEC 문서 업데이트
- 구조적 변경 감지 시 프로젝트 문서 (.moai/project/) 조건부 업데이트
- 업데이트 전략을 위한 SPEC 라이프사이클 레벨 준수 (spec-first, spec-anchored, spec-as-source)
- 성공 시 완료 마커 추가

## 팀 모드

--team 플래그가 제공되거나 자동 선택된 경우 (workflow.yaml의 복잡도 임계값 기준):

- Phase 0 탐색: 병렬 리서치 팀 (researcher + analyst + architect)
- Phase 2 구현: 병렬 구현 팀 (backend-dev + frontend-dev + tester)
- Phase 3 동기화: 항상 서브에이전트 모드 (manager-docs)

팀 오케스트레이션 세부 사항:
- Plan 단계: workflows/team-plan.md 참조
- Run 단계: workflows/team-run.md 참조
- Sync 근거: workflows/team-sync.md 참조

모드 선택:
- --team: 모든 해당 단계에 팀 모드 강제 적용
- --solo: 서브에이전트 모드 강제 적용
- --auto (기본값): workflow.yaml 임계값에 따른 복잡도 기반 선택

## 작업 추적

[HARD] 모든 작업 추적에 작업 관리 도구 필수:
- 이슈 발견 시: TaskCreate로 pending 상태로 생성
- 작업 시작 전: TaskUpdate로 in_progress 상태로 변경
- 작업 완료 후: TaskUpdate로 completed 상태로 변경
- TODO 목록을 텍스트로 출력하지 않음

## 안전 개발 프로토콜

모든 구현은 CLAUDE.md 섹션 7 안전 개발 프로토콜을 따릅니다:
- 접근 방식 우선: 코드 작성 전 접근 방식 설명
- 멀티 파일 분해: 3개 이상 파일 수정 시 작업 분리
- 구현 후 검토: 잠재적 이슈 나열 및 테스트 제안
- 재현 우선: 버그 수정 전 실패 테스트 작성

## 완료 마커

작업 완료 시 AI는 마커를 추가해야 합니다:
- `<moai>DONE</moai>` - 작업 완료
- `<moai>COMPLETE</moai>` - 전체 완료

## 실행 요약

1. 인수 파싱 (플래그 추출: --loop, --max, --sequential, --branch, --pr, --resume, --team, --solo, --auto)
2. SPEC ID와 함께 --resume 사용 시: 기존 SPEC 로드 후 마지막 상태에서 계속
3. quality.yaml에서 development_mode 감지 (hybrid/ddd/tdd)
4. **팀 모드 결정**: workflow.yaml 팀 설정을 읽고 실행 모드 결정
   - `--team` 플래그: 팀 모드 강제 적용 (workflow.team.enabled: true AND CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 환경변수 필요)
   - `--solo` 플래그: 서브에이전트 모드 강제 적용 (팀 모드 완전 건너뜀)
   - `--auto` 또는 플래그 없음 (기본값): workflow.yaml auto_selection의 복잡도 임계값 확인 (도메인 3개 이상, 파일 10개 이상, 또는 점수 7 이상)
   - 팀 모드 선택됐으나 전제 조건 미충족 시: 사용자에게 경고 후 서브에이전트 모드로 폴백
5. Phase 0 실행 (병렬 또는 순차 탐색)
6. 라우팅 결정 (단일 도메인 직접 위임 vs 전체 워크플로우)
7. 발견된 작업에 대해 TaskCreate
8. AskUserQuestion을 통한 사용자 확인
9. **Phase 1 (Plan)**: 팀 모드 → workflows/team-plan.md 읽고 팀 오케스트레이션 따름. 아니면 → manager-spec 서브에이전트
10. **Phase 2 (Run)**: 팀 모드 → workflows/team-run.md 읽고 팀 오케스트레이션 따름. 아니면 → manager-tdd (새 기능) 또는 manager-ddd (레거시 리팩토링) 서브에이전트
11. **Phase 3 (Sync)**: 항상 manager-docs 서브에이전트 (동기화 단계는 팀 모드 사용 안 함)
12. 완료 마커와 함께 종료

---

Version: 2.0.0
Source: alfred.md에서 이름 변경. plan->run->sync 파이프라인 통합. 동기화 단계에 SPEC/프로젝트 문서 업데이트 추가.
