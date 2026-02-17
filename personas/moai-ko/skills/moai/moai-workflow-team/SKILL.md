---
name: moai-workflow-team
description: >
  MoAI-ADK를 위한 Agent Teams 워크플로우 관리. 팀 생성, 팀원 스폰,
  작업 분해, 에이전트 간 메시지 전달, 정상적인 종료를 처리합니다.
  팀 기반 Plan 및 Run 단계를 위해 SPEC 워크플로우와 통합됩니다.
  팀을 사용할 수 없을 때 서브에이전트 모드로 자동 폴백하는
  이중 모드 실행을 지원합니다.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Task TeamCreate TeamDelete SendMessage TaskCreate TaskUpdate TaskList TaskGet Read Grep Glob AskUserQuestion
user-invocable: false
metadata:
  version: "1.1.0"
  category: "workflow"
  status: "experimental"
  updated: "2026-02-07"
  modularized: "false"
  tags: "team, agent-teams, collaboration, parallel, dual-mode"
  related-skills: "moai-workflow-spec, moai-workflow-ddd, moai-workflow-tdd"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 8000

# MoAI Extension: Triggers
triggers:
  keywords: ["team", "agent-team", "parallel", "collaborate", "teammates", "--team"]
  agents: ["moai"]
  phases: ["plan", "run"]
---

# MoAI Agent Teams 워크플로우

## 개요

이 스킬은 MoAI 워크플로우를 위한 Agent Teams 실행을 관리합니다. 팀 모드가 선택되면 (--team 플래그, 자동 감지, 또는 설정을 통해), MoAI는 지속적인 팀원을 조율하는 팀 리더로 작동합니다.

## 전제 조건

Agent Teams에는 다음이 필요합니다:
- settings.json env에 `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`
- workflow.yaml에 `workflow.team.enabled: true`
- Claude Code v2.1.32 이상

## 모드 선택

모드 선택기가 실행 전략을 결정합니다:

1. --team/--solo 플래그 확인 (사용자 재정의)
2. workflow.yaml execution_mode 설정 확인
3. "auto"인 경우: 복잡도 점수 분석
   - 도메인 수 >= 3: 팀 모드
   - 영향받는 파일 >= 10: 팀 모드
   - 복잡도 점수 >= 7: 팀 모드
   - 그 외: 서브에이전트 모드
4. AGENT_TEAMS 활성화 여부 확인
5. 활성화되지 않은 경우: 사용자에게 경고 후 서브에이전트로 폴백

## 팀 라이프사이클

### 1단계: 팀 생성

```
TeamCreate(team_name: "moai-{workflow}-{timestamp}")
```

팀 명명 규칙:
- Plan 단계: `moai-plan-SPEC-XXX`
- Run 단계: `moai-run-SPEC-XXX`
- 디버그: `moai-debug-{issue}`
- 리뷰: `moai-review-{target}`

### 2단계: 작업 분해

팀원을 스폰하기 전에 완전한 공유 작업 목록을 생성합니다:

```
TaskCreate(subject: "작업 설명", description: "상세 요구사항")
```

작업 분해 규칙:
- 각 작업은 자체 완결적이어야 합니다 (명확한 결과물 하나)
- 작업 간 의존성 정의 (addBlockedBy)
- 팀원 역할에 따라 파일 소유권 경계 할당
- 최적 흐름을 위해 팀원당 5-6개 작업 목표
- 해당하는 경우 작업이 SPEC 요구사항에 매핑되어야 함

### 3단계: 팀원 스폰

team_name 파라미터와 함께 Task 도구를 사용하여 팀원을 스폰합니다:

```
Task(
  subagent_type: "team-backend-dev",
  team_name: "moai-run-SPEC-XXX",
  name: "backend-dev",
  prompt: "당신은 이 팀의 백엔드 개발자입니다. 파일 소유권: {detected_ownership}. SPEC 컨텍스트: {spec_summary}",
  mode: "plan"
)
```

스폰 규칙:
- 스폰 프롬프트에 SPEC 컨텍스트 포함
- 프로젝트 구조에서 감지된 파일 소유권 경계 할당 (파일 소유권 감지 참조)
- 역할별 적절한 모델 사용 (조사는 haiku, 구현은 sonnet)

계획 승인 (workflow.yaml `team.require_plan_approval: true`인 경우):
- `mode: "plan"`으로 구현 팀원 스폰
- 팀원은 코드를 작성하기 전에 계획을 제출해야 함
- 팀 리더는 팀원으로부터 plan_approval_request 메시지 수신
- 팀 리더는 계획 범위, 파일 소유권 준수, 접근 방식을 검토
- 승인: `SendMessage(type: "plan_approval_response", request_id: "{id}", recipient: "{name}", approve: true)`
- 피드백과 함께 거부: `SendMessage(type: "plan_approval_response", request_id: "{id}", recipient: "{name}", approve: false, content: "피드백 내용")`
- 승인 후 팀원은 plan 모드를 종료하고 구현 시작
- `require_plan_approval`이 false인 경우 `mode: "acceptEdits"`로 스폰

### 4단계: 조율

MoAI는 팀 리더로서 모니터링하고 조율합니다:

1. 팀원으로부터 자동 메시지 수신 (진행 상황, 완료, 이슈)
2. 직접 조율을 위해 SendMessage 사용
3. 모든 팀원에게 중요 업데이트 브로드캐스트
4. 파일 소유권 충돌 해결
5. 팀원이 차단된 경우 작업 재할당

조율 패턴:
- 백엔드가 API 완료 시: 사용 가능한 엔드포인트를 frontend-dev에게 알림
- 구현 완료 시: 품질 검증 작업 할당
- 품질 검사에서 이슈 발견 시: 담당 팀원에게 수정 메시지 전달
- 모든 작업 완료 시: 종료 시퀀스 시작

위임 모드 (workflow.yaml `team.delegate_mode: true`인 경우):
- MoAI는 조율 전용 모드로 작동
- 작업 할당, 메시지 라우팅, 진행 상황 모니터링, 충돌 해결에 집중
- 코드를 직접 구현하거나 파일 수정 금지 (구현을 위한 Write, Edit, Bash 사용 금지)
- 모든 구현 작업을 작업 할당과 SendMessage를 통해 팀원에게 위임
- 컨텍스트 이해와 팀원 출력 검토를 위해 Read와 Grep은 허용
- 적합한 팀원이 없는 작업은 직접 구현 대신 새 팀원 스폰
- delegate_mode가 false인 경우 팀 리더가 팀원들과 함께 소규모 작업을 직접 구현 가능

### 5단계: 종료

정상적인 종료 시퀀스:

1. TaskList를 통해 모든 작업이 완료되었는지 확인
2. 각 팀원에게 shutdown_request 전송:
   ```
   SendMessage(type: "shutdown_request", recipient: "backend-dev")
   ```
3. 각 팀원으로부터 종료 승인 대기
4. 팀 리소스 정리:
   ```
   TeamDelete()
   ```

## 파일 소유권 전략

독점 파일 소유권을 할당하여 쓰기 충돌을 방지합니다.

[HARD] 팀 리더는 소유권을 할당하기 전에 반드시 프로젝트 구조를 분석해야 합니다. Explore 에이전트 또는 Glob/Grep을 사용하여 디렉토리 구조를 매핑하고 실제 프로젝트 레이아웃에 맞는 소유권 경계를 할당합니다. 다른 프로젝트 유형의 하드코딩된 패턴을 절대 사용하지 마세요.

### 파일 소유권 감지

소유권 패턴은 프로젝트 유형에 따라 다릅니다. 먼저 프로젝트 구조를 감지한 후 그에 맞게 할당합니다:

**Go 프로젝트:**

| 역할 | 소유권 |
|------|-----------|
| backend-dev | internal/**, pkg/**, cmd/** |
| tester | *_test.go, testdata/**, test/** |
| quality | (읽기 전용, 파일 소유권 없음) |

**웹 프로젝트 (React, Vue, Angular):**

| 역할 | 소유권 |
|------|-----------|
| backend-dev | src/api/**, src/models/**, src/services/** |
| frontend-dev | src/ui/**, src/components/**, src/pages/** |
| tester | tests/**, __tests__/**, *.test.*, *.spec.* |
| quality | (읽기 전용, 파일 소유권 없음) |

**풀스택 프로젝트 (별도 client/server):**

| 역할 | 소유권 |
|------|-----------|
| backend-dev | server/**, api/**, src/server/** |
| frontend-dev | client/**, app/**, src/client/** |
| data-layer | db/**, migrations/**, schema/** |
| tester | tests/**, __tests__/**, *_test.go, *.test.*, *.spec.* |
| quality | (읽기 전용, 파일 소유권 없음) |

**모노레포 프로젝트:**

| 역할 | 소유권 |
|------|-----------|
| 도메인별 팀원 | packages/<domain-name>/**, apps/<domain-name>/** |
| tester | **/tests/**, **/__tests__/**, **/*_test.go, **/*.test.*, **/*.spec.* |
| quality | (읽기 전용, 파일 소유권 없음) |

**Python 프로젝트:**

| 역할 | 소유권 |
|------|-----------|
| backend-dev | src/<package>/**, <package>/** |
| tester | tests/**, **/test_*.py, **/*_test.py |
| quality | (읽기 전용, 파일 소유권 없음) |

### 소유권 규칙

- 두 팀원이 동일한 파일을 소유할 수 없음
- 공유 타입/인터페이스: 생성한 팀원이 소유, 메시지를 통해 공유
- 설정 파일: 팀 리더가 소유하거나 명시적으로 할당
- 소유권 충돌 시: 팀 리더가 SendMessage로 해결
- 테스트 파일은 위치에 관계없이 항상 tester 역할에 속함

## 팀 패턴 참조

### 계획 조사 팀
- 역할: researcher (haiku), analyst (sonnet), architect (sonnet)
- 용도: 다각도 탐색이 필요한 복잡한 SPEC 생성
- 지속 시간: 단기 (탐색 단계만)

### 구현 팀
- 역할: backend-dev (sonnet), frontend-dev (sonnet), tester (sonnet)
- 용도: 교차 계층 기능 구현
- 지속 시간: 중기 (전체 run 단계)

### 풀스택 팀
- 역할: api-layer, ui-layer, data-layer, quality (모두 sonnet)
- 용도: 대규모 풀스택 기능
- 지속 시간: 중-장기

### 조사 팀
- 역할: hypothesis-1, hypothesis-2, hypothesis-3 (모두 haiku)
- 용도: 경쟁 이론을 가진 복잡한 디버깅
- 지속 시간: 단기

### 리뷰 팀
- 역할: security-reviewer, perf-reviewer, quality-reviewer (모두 sonnet)
- 용도: 다각도 코드 리뷰
- 지속 시간: 단기

## 에러 복구

- 팀원 충돌: 동일 역할과 재개 컨텍스트로 교체 스폰
- 작업 지연: 팀 리더가 다른 팀원에게 재할당
- 파일 충돌: 팀 리더가 SendMessage로 중재, 소유권 조정
- 모든 팀원 idle: 남은 작업 확인 후 할당 또는 종료
- 토큰 한도: 팀을 정상적으로 종료하고 나머지 작업은 서브에이전트로 폴백

---

Version: 1.1.0
Last Updated: 2026-02-07
