# 워크플로우: Team Plan - Agent Teams SPEC 생성

목적: 병렬 팀 기반 리서치 및 분석을 통해 포괄적인 SPEC 문서를 생성합니다. plan 단계가 다각적 탐색으로 이점을 얻을 때 사용합니다.

흐름: TeamCreate -> 병렬 리서치 -> 종합 -> SPEC 문서 -> 종료

## 전제 조건

- workflow.team.enabled: true
- CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1
- 트리거: /moai plan --team 또는 복잡도 >= 임계값 자동 감지

## Phase 0: 팀 설정

1. 설정 읽기:
   - 팀 설정을 위한 .moai/config/sections/workflow.yaml
   - 개발 모드를 위한 .moai/config/sections/quality.yaml

2. 팀 생성:
   ```
   TeamCreate(team_name: "moai-plan-{feature-slug}")
   ```

3. 공유 작업 목록 생성:
   ```
   TaskCreate: "코드베이스 아키텍처 및 의존성 탐색"
   TaskCreate: "요구사항, 사용자 스토리, 엣지 케이스 분석"
   TaskCreate: "기술 접근 방식 설계 및 대안 평가"
   TaskCreate: "발견사항을 SPEC 문서로 종합" (위 3개에 의해 차단됨)
   ```

## Phase 1: 리서치 팀 소집

plan_research 패턴에서 3명의 팀원 소집:

팀원 1 - researcher (team-researcher 에이전트, haiku 모델):
- 프롬프트: "{feature_description}을 위해 코드베이스를 탐색하세요. 아키텍처를 매핑하고, 관련 파일을 찾고, 의존성과 패턴을 식별하세요. 팀 리더에게 발견사항을 보고하세요."

팀원 2 - analyst (team-analyst 에이전트, inherit 모델):
- 프롬프트: "{feature_description}에 대한 요구사항을 분석하세요. 사용자 스토리, 인수 기준, 엣지 케이스, 리스크, 제약 조건을 식별하세요. 팀 리더에게 발견사항을 보고하세요."

팀원 3 - architect (team-architect 에이전트, inherit 모델):
- 프롬프트: "{feature_description}에 대한 기술 접근 방식을 설계하세요. 구현 대안을 평가하고, 트레이드오프를 평가하고, 아키텍처를 제안하세요. researcher가 찾은 기존 패턴을 고려하세요. 팀 리더에게 보고하세요."

## Phase 2: 병렬 리서치

팀원들은 독립적으로 작업합니다:
- researcher가 코드베이스 탐색 (가장 빠름, haiku)
- analyst가 요구사항 정의 (중간)
- architect가 솔루션 설계 (researcher 발견사항 대기)

MoAI 모니터링:
- 자동으로 진행 메시지 수신
- 가능할 때 researcher 발견사항을 architect에게 전달
- 팀원들의 질문 해결

## Phase 3: 종합

모든 리서치 작업 완료 후:
1. 세 팀원의 발견사항 수집
2. 모든 발견사항과 함께 manager-spec 서브에이전트에게 SPEC 생성 위임 (팀원이 아님)
3. 포함 내용: 코드베이스 분석, 요구사항, 기술 설계, 엣지 케이스

SPEC 출력 위치: .moai/specs/SPEC-XXX/spec.md

## Phase 4: 사용자 승인

AskUserQuestion 옵션:
- SPEC 승인 및 구현 진행
- 수정 요청 (어느 섹션인지 명시)
- 워크플로우 취소

## Phase 5: 정리

1. 모든 팀원 종료:
   ```
   SendMessage(type: "shutdown_request", recipient: "researcher")
   SendMessage(type: "shutdown_request", recipient: "analyst")
   SendMessage(type: "shutdown_request", recipient: "architect")
   ```
2. 리소스 정리를 위해 TeamDelete
3. 다음 단계를 위한 컨텍스트 확보를 위해 /clear 실행

## 폴백

팀 생성 실패 또는 AGENT_TEAMS 미활성화 시:
- 서브에이전트 plan 워크플로우 (workflows/plan.md)로 폴백
- 팀 모드 사용 불가에 대한 경고 기록

---

Version: 1.1.0
