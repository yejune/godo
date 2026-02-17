# 워크플로우: Team Plan - Agent Teams SPEC 작성

목적: 병렬 팀 기반 조사 및 분석을 통해 포괄적인 SPEC 문서를 작성한다. 플랜 단계가 다각도 탐색으로 효과가 있을 때 사용한다.

흐름: TeamCreate -> 병렬 조사 -> 종합 -> SPEC 문서 -> 종료

## 전제 조건

- workflow.team.enabled: true
- CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1
- 트리거: /do plan --team 또는 복잡도 >= 임계값 자동 감지

## Phase 0: 팀 설정

1. 설정 읽기:
   - 팀 설정을 위한 .do/config/sections/workflow.yaml
   - 개발 모드를 위한 .do/config/sections/quality.yaml

2. 팀 생성:
   ```
   TeamCreate(team_name: "do-plan-{feature-slug}")
   ```

3. 공유 태스크 목록 생성:
   ```
   TaskCreate: "코드베이스 아키텍처 및 의존성 탐색"
   TaskCreate: "요구사항, 사용자 스토리, 엣지 케이스 분석"
   TaskCreate: "기술 접근법 설계 및 대안 평가"
   TaskCreate: "결과를 SPEC 문서로 종합" (위 3개에 의해 차단됨)
   ```

## Phase 1: 조사 팀 소환

plan_research 패턴에서 3명의 팀원 소환:

팀원 1 - researcher (team-researcher 에이전트, haiku 모델):
- 프롬프트: "{feature_description}를 위한 코드베이스 탐색. 아키텍처 파악, 관련 파일 찾기, 의존성 및 패턴 식별. 팀 리더에게 결과 보고."

팀원 2 - analyst (team-analyst 에이전트, inherit 모델):
- 프롬프트: "{feature_description}에 대한 요구사항 분석. 사용자 스토리, 인수 기준, 엣지 케이스, 위험, 제약사항 식별. 팀 리더에게 결과 보고."

팀원 3 - architect (team-architect 에이전트, inherit 모델):
- 프롬프트: "{feature_description}의 기술 접근법 설계. 구현 대안 평가, 트레이드오프 검토, 아키텍처 제안. researcher가 찾은 기존 패턴을 고려할 것. 팀 리더에게 보고."

## Phase 2: 병렬 조사

팀원들이 독립적으로 작업:
- researcher가 코드베이스 탐색 (가장 빠름, haiku)
- analyst가 요구사항 정의 (중간)
- architect가 솔루션 설계 (researcher 결과 대기)

Do 모니터링:
- 자동으로 진행 메시지 수신
- researcher 결과가 도착하면 architect에게 전달
- 팀원의 질문 해결

## Phase 3: 종합

모든 조사 태스크 완료 후:
1. 세 팀원의 결과 수집
2. 모든 결과와 함께 manager-spec 서브에이전트 (팀원이 아님)에 SPEC 작성 위임
3. 코드베이스 분석, 요구사항, 기술 설계, 엣지 케이스 포함

SPEC 출력 위치: .do/specs/SPEC-XXX/spec.md

## Phase 4: 사용자 승인

AskUserQuestion 옵션:
- SPEC 승인 후 구현으로 진행
- 수정 요청 (어느 섹션인지 명시)
- 워크플로우 취소

## Phase 5: 정리

1. 모든 팀원 종료:
   ```
   SendMessage(type: "shutdown_request", recipient: "researcher")
   SendMessage(type: "shutdown_request", recipient: "analyst")
   SendMessage(type: "shutdown_request", recipient: "architect")
   ```
2. 리소스 정리를 위한 TeamDelete
3. 다음 단계를 위한 컨텍스트 확보를 위해 /clear 실행

## 폴백

팀 생성 실패 또는 AGENT_TEAMS 미활성화 시:
- 서브에이전트 플랜 워크플로우 (workflows/plan.md)로 폴백
- 팀 모드 사용 불가 경고 로그

---

Version: 1.1.0
