---
name: do-workflow-team-do
description: >
  Agent Teams API를 사용하여 plan, run, test, report 단계를 결합한
  Team 모드 전체 파이프라인. 계획을 위한 병렬 조사 팀, 실행을 위한
  병렬 구현 팀, 품질 검증, 보고를 포함한다.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-15"
  tags: "team, autopilot, parallel, agent-teams, plan, run, report"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# Do Extension: Triggers
triggers:
  keywords: ["team", "autopilot", "parallel", "팀"]
  agents: ["team-researcher", "team-analyst", "team-architect", "team-backend-dev", "team-frontend-dev", "team-tester", "team-quality"]
  phases: ["plan", "run", "test", "report"]
---

# 워크플로우: Team Do - 팀 자동 파이프라인

목적: Agent Teams API를 사용한 병렬 실행 완전 자율 파이프라인. 플랜 (병렬 조사), 실행 (병렬 구현), 테스트, 보고를 단일 팀 오케스트레이션 워크플로우로 결합한다.

흐름: 팀 플랜 -> 체크리스트 -> 팀 실행 -> 품질 검증 -> 보고 -> 종료

## 전제 조건

- CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 환경변수 설정
- Team 모드 활성화 (DO_MODE=team 또는 --team 플래그)
- 전제 조건 미충족 시: 사용자에게 경고 후 Do 모드 (workflows/do.md)로 폴백

---

## Part 1: 팀 플랜 (병렬 조사)

### 팀 구성

| 팀원 | 에이전트 | 모델 | 모드 | 목적 |
|------|---------|------|------|------|
| researcher | team-researcher | haiku | plan (읽기 전용) | 코드베이스 탐색 |
| analyst | team-analyst | inherit | plan (읽기 전용) | 요구사항 분석 |
| architect | team-architect | inherit | plan (읽기 전용) | 기술 설계 |

### 실행

1. 팀 생성:
   ```
   TeamCreate(team_name: "do-plan-{feature-slug}")
   ```

2. 공유 태스크 목록 생성:
   ```
   TaskCreate: "코드베이스 아키텍처 및 의존성 탐색"
   TaskCreate: "요구사항, 사용자 스토리, 엣지 케이스 분석"
   TaskCreate: "기술 접근법 설계 및 대안 평가"
   ```

3. 조사 프롬프트와 함께 3명의 팀원 소환:

   - **researcher**: "{feature_description}를 위한 코드베이스 탐색. 아키텍처 파악, 관련 파일 찾기, 의존성 및 패턴 식별. 팀 리더에게 결과 보고."
   - **analyst**: "{feature_description}에 대한 요구사항 분석. 사용자 스토리, 인수 기준, 엣지 케이스, 위험, 제약사항 식별. 팀 리더에게 결과 보고."
   - **architect**: "{feature_description}의 기술 접근법 설계. 구현 대안 평가, 트레이드오프 검토, 아키텍처 제안. 팀 리더에게 보고."

4. 병렬 조사 모니터링:
   - 자동으로 진행 메시지 수신
   - researcher 결과가 도착하면 architect에게 전달
   - 팀원의 질문 해결

5. 모든 조사 태스크 완료 후:
   - 세 팀원의 결과 수집
   - 플랜으로 종합: 모든 결과와 함께 plan 에이전트 (팀원이 아닌 서브에이전트)에 위임
   - `.do/jobs/{YY}/{MM}/{DD}/{title}/`에 analysis.md, architecture.md, plan.md 생성

6. 사용자 승인:
   - AskUserQuestion: "플랜 완료. 구현으로 진행할까요?"
   - "진행" -> 체크리스트 + 팀 실행으로 계속
   - "수정" -> 피드백 수집, 종합 재실행
   - "취소" -> 팀 종료, 종료

7. 플랜 팀 종료:
   ```
   SendMessage(type: "shutdown_request", recipient: "researcher")
   SendMessage(type: "shutdown_request", recipient: "analyst")
   SendMessage(type: "shutdown_request", recipient: "architect")
   ```

### 체크리스트 생성

플랜 승인 후 checklist.md + 서브 체크리스트 생성 (do.md Phase 2와 동일).

---

## Part 2: 팀 실행 (병렬 구현)

### 팀 구성

| 팀원 | 에이전트 | 모델 | 모드 | 목적 |
|------|---------|------|------|------|
| backend-dev | team-backend-dev | inherit | acceptEdits | 서버 사이드 구현 |
| frontend-dev | team-frontend-dev | inherit | acceptEdits | 클라이언트 사이드 구현 |
| tester | team-tester | inherit | acceptEdits | 테스트 작성 및 커버리지 |
| quality | team-quality | inherit | plan (읽기 전용) | 품질 검증 |

프로젝트 필요에 따라 팀 구성 조정 (예: 백엔드 전용 태스크는 frontend-dev 생략).

### 파일 소유권 [HARD]

서브 체크리스트 핵심 파일 섹션에 따라 파일 소유권 배정:

- 각 팀원은 자신의 서브 체크리스트에서 특정 파일을 소유
- 두 팀원이 같은 파일을 수정하면 안 됨
- 겹침이 감지되면: 소환 전에 해결 (파일의 변경 사항을 분리하거나 한 소유자에게 배정)

### Git 스테이징 안전 규칙 [HARD]

모든 팀원은 반드시 다음 규칙을 따라야 한다:

- [HARD] 본인 파일만 스테이징: `git add file1.go file2.go` (개별 파일만)
- [HARD] 광범위 스테이징 절대 금지: `git add -A`, `git add .`, `git add --all` 금지
- [HARD] 커밋 전 검증: `git diff --cached --name-only`에 본인 파일만 표시되어야 함
- [HARD] 외부 파일 확인: 다른 팀원이 스테이징만 하고 커밋하지 않았을 수 있음. 외부 파일이 스테이징되어 있으면 unstage: `git reset HEAD <file>`
- [HARD] 커밋 전 unstage: 예상치 못한 것이 스테이징되어 있으면 먼저 `git reset HEAD <file>`

### 실행

1. 팀 생성:
   ```
   TeamCreate(team_name: "do-run-{feature-slug}")
   ```

2. 의존성이 있는 체크리스트 항목으로 공유 태스크 목록 생성

3. 서브 체크리스트 경로와 함께 팀원 소환:
   - 각 팀원 프롬프트 포함: 플랜 요약, 서브 체크리스트 경로, 파일 소유권 목록, git 스테이징 규칙, 품질 목표

4. 병렬 구현:
   - backend-dev: 서버 사이드 코드 구현 (태스크 1-2)
   - frontend-dev: 클라이언트 사이드 코드 구현, API 계약 대기 (태스크 3)
   - tester: 구현 완료 후 통합 테스트 작성 (태스크 4)
   - 팀원들은 API 계약 및 의존성을 위해 SendMessage로 협력

5. 오케스트레이터 조율:
   - 백엔드에서 프론트엔드로 API 계약 정보 전달
   - TaskList로 태스크 진행 상황 모니터링
   - 차단 이슈 해결
   - 접근 방식이 효과 없으면 팀원 방향 전환

6. 품질 검증 (구현 + 테스트 완료 후):
   - quality 팀원에게 품질 검증 태스크 배정
   - quality가 전체 테스트 스위트 실행 및 커버리지 확인
   - 팀 리더에게 결과 보고
   - 이슈 발견 시: 담당 팀원에게 수정 지시

### 품질 게이트 (내장 규칙)

진행 전 모두 통과해야 함 (항상 활성화된 규칙으로 적용, 별도 프레임워크 아님):
- 린트 오류 없음 (Unified)
- 타입 오류 없음 (Unified)
- 커버리지 목표 충족 -- 85%+ 전체 (Tested)
- 모든 체크리스트 인수 기준 검증됨 (Trackable)
- Real DB로 모든 테스트 통과 -- mock DB 대체 금지 (Tested)
- 테스트 코드에 AI 안티패턴 없음 (7가지 금지 패턴)
- 모든 `[o]` 항목에 증거로 커밋 해시 있음 (Commit-as-Proof)

---

## Part 3: 보고

품질 검증 통과 후:

1. 완료 보고서 생성 (report.md 워크플로우와 동일):
   - 서브 체크리스트 결과 집계
   - 모든 팀원의 교훈 수집
   - 테스트 결과 요약

2. 사용자에게 체크리스트 최종 상태 표시

3. 다음 단계와 함께 AskUserQuestion

---

## Part 4: 정리

1. 모든 팀원을 정상적으로 종료:
   ```
   SendMessage(type: "shutdown_request", recipient: "backend-dev")
   SendMessage(type: "shutdown_request", recipient: "frontend-dev")
   SendMessage(type: "shutdown_request", recipient: "tester")
   SendMessage(type: "shutdown_request", recipient: "quality")
   ```

2. 사용자에게 최종 요약 표시

---

## 폴백

어느 시점에서든 팀 생성 실패 또는 Agent Teams 미활성화 시:
- 남은 팀원 정상 종료
- Do 모드 워크플로우 (workflows/do.md)로 폴백
- 마지막 완료된 단계부터 계속
- 팀 모드 사용 불가 경고 로그

---

Version: 1.0.0
Updated: 2026-02-16
