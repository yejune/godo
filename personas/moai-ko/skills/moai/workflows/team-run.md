# 워크플로우: Team Run - Agent Teams 구현

목적: 병렬 팀 기반 개발을 통해 SPEC 요구사항을 구현합니다. 각 팀원은 충돌 방지를 위해 특정 파일/도메인을 소유합니다.

흐름: TeamCreate -> 작업 분해 -> 병렬 구현 -> 품질 검증 -> 종료

## 전제 조건

- .moai/specs/SPEC-XXX/에 승인된 SPEC 문서
- workflow.team.enabled: true
- CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1
- 트리거: /moai run SPEC-XXX --team 또는 복잡도 >= 임계값 자동 감지

## Phase 0: SPEC 분석 및 작업 분해

1. SPEC 문서 읽기 및 범위 분석
2. 개발 모드를 위해 quality.yaml 읽기:
   - hybrid (신규 프로젝트 기본값): 새 코드는 TDD, 기존 코드는 DDD
   - ddd (기존 프로젝트): 모든 코드에 ANALYZE-PRESERVE-IMPROVE 적용
   - tdd (명시적 선택): 모든 코드에 RED-GREEN-REFACTOR 적용

3. SPEC을 구현 작업으로 분해:
   - 도메인 경계 식별 (backend, frontend, data, tests)
   - 도메인별 파일 소유권 할당
   - 명확한 의존성이 있는 작업 생성
   - 팀원당 5-6개 작업 목표

4. 팀 생성:
   ```
   TeamCreate(team_name: "moai-run-SPEC-XXX")
   ```

5. 의존성이 있는 공유 작업 목록 생성:
   ```
   TaskCreate: "데이터 모델 및 스키마 구현" (의존성 없음)
   TaskCreate: "API 엔드포인트 구현" (데이터 모델에 의해 차단됨)
   TaskCreate: "UI 컴포넌트 구현" (API 엔드포인트에 의해 차단됨)
   TaskCreate: "단위 및 통합 테스트 작성" (API + UI에 의해 차단됨)
   TaskCreate: "품질 검증 - TRUST 5" (위 모든 항목에 의해 차단됨)
   ```

## Phase 1: 구현 팀 소집

SPEC 범위에 따라 팀 패턴 선택:

크로스 레이어 기능 (implementation 패턴):
- backend-dev (team-backend-dev, inherit): 서버 측 구현
- frontend-dev (team-frontend-dev, inherit): 클라이언트 측 구현
- tester (team-tester, inherit): 테스트 생성 및 커버리지

디자인이 포함된 크로스 레이어 기능 (design_implementation 패턴):
- designer (team-designer, inherit): Pencil/Figma MCP를 사용한 UI/UX 설계
- backend-dev (team-backend-dev, inherit): 서버 측 구현
- frontend-dev (team-frontend-dev, inherit): 클라이언트 측 구현
- tester (team-tester, inherit): 테스트 생성 및 커버리지

풀스택 기능 (full_stack 패턴):
- api-layer (team-backend-dev, inherit): API 및 비즈니스 로직
- ui-layer (team-frontend-dev, inherit): UI 및 컴포넌트
- data-layer (team-backend-dev, inherit): 데이터베이스 및 스키마
- quality (team-quality, inherit): 품질 검증

소집 프롬프트에 포함 필수:
- SPEC 요약 및 해당 팀원의 특정 요구사항
- 파일 소유권 경계 (프로젝트 구조에서 감지, SKILL.md 파일 소유권 감지 참조)
- 개발 방법론 (새 코드는 TDD, 기존 코드는 DDD)
- 품질 목표 (커버리지, 린트, 타입 검사)

### 계획 승인 모드

workflow.yaml `team.require_plan_approval: true`인 경우:
- `mode: "acceptEdits"` 대신 `mode: "plan"`으로 구현 팀원 소집
- 각 팀원은 코드 구현 전에 계획을 제출해야 함
- 팀 리더는 제안된 접근 방식과 함께 `plan_approval_request` 메시지를 받음
- 팀 리더 검토: 파일 소유권 준수, SPEC과의 접근 방식 일치, 범위 정확성
- 승인: `SendMessage(type: "plan_approval_response", request_id: "{id}", recipient: "{name}", approve: true)`
- 피드백과 함께 거절: `SendMessage(type: "plan_approval_response", request_id: "{id}", recipient: "{name}", approve: false, content: "X를 수정하세요")`
- 승인 후 팀원은 자동으로 plan 모드를 종료하고 구현 시작
- `require_plan_approval`이 false이면 `mode: "acceptEdits"`로 직접 소집

## Phase 2: 병렬 구현

팀원들은 공유 목록에서 작업을 스스로 가져가 독립적으로 작업합니다:

디자인 (team-designer 포함 시):
- Pencil MCP 또는 Figma MCP를 사용하여 UI/UX 설계 (작업 0)
- 디자인 토큰, 스타일 가이드, 컴포넌트 명세 제작
- SendMessage를 통해 frontend-dev와 디자인 명세 공유
- 디자인 파일 소유 (.pen, 디자인 토큰, 스타일 설정)

백엔드 개발:
- 데이터 모델 및 스키마 생성 (작업 1)
- API 엔드포인트 및 비즈니스 로직 구현 (작업 2)
- 새 코드에는 TDD 적용: 테스트 작성 -> 구현 -> 리팩토링
- 기존 코드에는 DDD 적용: 분석 -> 테스트로 보존 -> 개선
- API 계약이 준비되면 frontend-dev에게 알림

프론트엔드 개발:
- backend-dev의 API 계약 및 designer의 디자인 명세 대기
- UI 컴포넌트 및 페이지 구현 (작업 3)
- 새 컴포넌트에는 TDD 적용
- SendMessage를 통해 데이터 형태는 backend와, 시각적 명세는 designer와 조율

테스트:
- 구현 작업 완료 대기
- API와 UI에 걸친 통합 테스트 작성 (작업 4)
- 커버리지 목표 검증
- 담당 팀원에게 테스트 실패 보고

MoAI 조율:
- backend에서 frontend로 API 계약 정보 전달
- 차단 이슈 해결
- TaskList를 통해 작업 진행 상황 모니터링
- 접근 방식이 효과적이지 않을 경우 팀원 재지시

### 위임 모드

workflow.yaml `team.delegate_mode: true`인 경우:
- MoAI는 전체 실행 단계 동안 조율 전용 모드로 운영
- 집중 사항: 작업 할당, 메시지 라우팅, 진행 모니터링, 충돌 해결
- 구현 코드 직접 작성 또는 파일 수정 금지 (구현을 위한 Write, Edit, Bash 사용 금지)
- 모든 구현을 작업 할당 및 SendMessage를 통해 팀원에게 위임
- 컨텍스트 이해 및 팀원 출력 검토를 위해 Read와 Grep은 허용
- 적합한 팀원이 없는 작업은 직접 구현하지 말고 새 팀원 소집
- delegate_mode가 false이면 팀 리더가 팀원들과 함께 소규모 작업을 직접 구현 가능

## Phase 3: 품질 검증

구현 및 테스트 작업 완료 후:

옵션 A (quality 팀원 포함 시):
- 작업 5를 quality 팀원에게 할당
- quality가 TRUST 5 검증 실행
- 팀 리더에게 발견사항 보고
- 팀 리더가 담당 팀원에게 수정 지시

옵션 B (더 작은 팀의 경우 서브에이전트 사용):
- manager-quality 서브에이전트에게 품질 검증 위임
- 발견사항 검토 및 필요 시 수정 작업 생성
- 기존 팀원에게 수정 할당

품질 게이트 (모두 통과 필요):
- 린트 오류 없음
- 타입 오류 없음
- 커버리지 목표 충족 (전체 85%+, 새 코드 90%+)
- Critical 보안 이슈 없음
- 모든 SPEC 인수 기준 검증

## Phase 4: Git 작업

품질 검증 통과 후:
- manager-git 서브에이전트에 위임 (팀원이 아님)
- 컨벤셔널 커밋 형식으로 의미 있는 커밋 생성
- 커밋 메시지에 SPEC ID 참조

## Phase 5: 정리

1. 모든 팀원 우아하게 종료
2. 리소스 정리를 위해 TeamDelete
3. 사용자에게 구현 요약 보고

## 폴백

어떤 시점에서든 팀 모드 실패 시:
- 남은 팀원 우아하게 종료
- 서브에이전트 실행 워크플로우 (workflows/run.md)로 폴백
- 마지막으로 완료된 작업에서 계속
- 팀 모드 실패에 대한 경고 기록

## 작업 추적

[HARD] TaskUpdate를 통한 모든 작업 상태 변경:
- pending -> in_progress: 팀원이 작업을 가져갈 때
- in_progress -> completed: 작업이 검증될 때
- 일반 텍스트 TODO 목록 사용 금지

---

Version: 1.2.0
