# 워크플로우: Team Run - Agent Teams 구현

목적: 병렬 팀 기반 개발을 통해 SPEC 요구사항을 구현한다. 각 팀원이 충돌을 방지하기 위해 특정 파일/도메인을 소유한다.

흐름: TeamCreate -> 태스크 분해 -> 병렬 구현 -> 품질 검증 -> 종료

## 전제 조건

- .do/specs/SPEC-XXX/에 승인된 SPEC 문서
- workflow.team.enabled: true
- CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1
- 트리거: /do run SPEC-XXX --team 또는 복잡도 >= 임계값 자동 감지

## Phase 0: SPEC 분석 및 태스크 분해

1. SPEC 문서 읽기 및 범위 분석
2. 개발 모드를 위한 quality.yaml 읽기:
   - hybrid (신규 프로젝트 기본값): 새 코드는 TDD, 기존 코드는 DDD
   - ddd (기존 프로젝트용): 모든 코드에 ANALYZE-PRESERVE-IMPROVE 적용
   - tdd (명시적 선택): 모든 코드에 RED-GREEN-REFACTOR 적용

3. SPEC을 구현 태스크로 분해:
   - 도메인 경계 식별 (backend, frontend, data, tests)
   - 도메인별 파일 소유권 배정
   - 명확한 의존성이 있는 태스크 생성
   - 팀원당 5-6개 태스크 목표

4. 팀 생성:
   ```
   TeamCreate(team_name: "do-run-SPEC-XXX")
   ```

5. 의존성이 있는 공유 태스크 목록 생성:
   ```
   TaskCreate: "데이터 모델 및 스키마 구현" (의존성 없음)
   TaskCreate: "API 엔드포인트 구현" (데이터 모델에 의해 차단됨)
   TaskCreate: "UI 컴포넌트 구현" (API 엔드포인트에 의해 차단됨)
   TaskCreate: "단위 및 통합 테스트 작성" (API + UI에 의해 차단됨)
   TaskCreate: "품질 검증 - TRUST 5" (위 모두에 의해 차단됨)
   ```

## Phase 1: 구현 팀 소환

SPEC 범위에 따라 팀 패턴 선택:

크로스 레이어 기능 (implementation 패턴):
- backend-dev (team-backend-dev, inherit): 서버 사이드 구현
- frontend-dev (team-frontend-dev, inherit): 클라이언트 사이드 구현
- tester (team-tester, inherit): 테스트 작성 및 커버리지

설계가 포함된 크로스 레이어 기능 (design_implementation 패턴):
- designer (team-designer, inherit): Pencil/Figma MCP를 사용한 UI/UX 설계
- backend-dev (team-backend-dev, inherit): 서버 사이드 구현
- frontend-dev (team-frontend-dev, inherit): 클라이언트 사이드 구현
- tester (team-tester, inherit): 테스트 작성 및 커버리지

풀스택 기능 (full_stack 패턴):
- api-layer (team-backend-dev, inherit): API 및 비즈니스 로직
- ui-layer (team-frontend-dev, inherit): UI 및 컴포넌트
- data-layer (team-backend-dev, inherit): 데이터베이스 및 스키마
- quality (team-quality, inherit): 품질 검증

소환 프롬프트에 반드시 포함:
- SPEC 요약 및 해당 팀원의 구체적인 요구사항
- 파일 소유권 경계 (프로젝트 구조에서 감지, SKILL.md 파일 소유권 감지 참조)
- 개발 방법론 (새 코드는 TDD, 기존 코드는 DDD)
- 품질 목표 (커버리지, 린트, 타입 검사)

### 플랜 승인 모드

workflow.yaml의 `team.require_plan_approval: true` 시:
- 구현 팀원을 `mode: "acceptEdits"` 대신 `mode: "plan"`으로 소환
- 각 팀원은 코드 구현 전 플랜을 제출해야 함
- 팀 리더가 제안된 접근 방식과 함께 `plan_approval_request` 메시지 수신
- 팀 리더 검토 항목: 파일 소유권 준수, SPEC과의 접근 방식 정합성, 범위 적절성
- 승인: `SendMessage(type: "plan_approval_response", request_id: "{id}", recipient: "{name}", approve: true)`
- 피드백과 함께 거부: `SendMessage(type: "plan_approval_response", request_id: "{id}", recipient: "{name}", approve: false, content: "X를 수정하라")`
- 승인 후 팀원이 자동으로 플랜 모드를 종료하고 구현 시작
- `require_plan_approval`이 false이면 `mode: "acceptEdits"`로 바로 소환

## Phase 2: 병렬 구현

팀원들이 공유 목록에서 태스크를 자율 선택하여 독립적으로 작업:

설계 (team-designer 포함 시):
- Pencil MCP 또는 Figma MCP를 사용하여 UI/UX 설계 (태스크 0)
- 디자인 토큰, 스타일 가이드, 컴포넌트 명세 작성
- SendMessage로 frontend-dev에게 디자인 명세 공유
- 디자인 파일 소유 (.pen, 디자인 토큰, 스타일 설정)

백엔드 개발:
- 데이터 모델 및 스키마 생성 (태스크 1)
- API 엔드포인트 및 비즈니스 로직 구현 (태스크 2)
- 새 코드는 TDD 따름: 테스트 작성 -> 구현 -> 리팩토링
- 기존 코드는 DDD 따름: 분석 -> 테스트로 보존 -> 개선
- API 계약 준비되면 frontend-dev에게 알림

프론트엔드 개발:
- backend-dev의 API 계약 및 designer의 디자인 명세 대기
- UI 컴포넌트 및 페이지 구현 (태스크 3)
- 새 컴포넌트는 TDD 따름
- SendMessage로 데이터 형태는 백엔드와, 시각적 명세는 designer와 협력

테스트:
- 구현 태스크 완료 대기
- API와 UI를 아우르는 통합 테스트 작성 (태스크 4)
- 커버리지 목표 검증
- 담당 팀원에게 테스트 실패 보고

Do 조율:
- 백엔드에서 프론트엔드로 API 계약 정보 전달
- 차단 이슈 해결
- TaskList로 태스크 진행 상황 모니터링
- 접근 방식이 효과 없으면 팀원 방향 전환

### 위임 모드

workflow.yaml의 `team.delegate_mode: true` 시:
- Do가 실행 단계 전체에서 조율 전용 모드로 운영
- 집중 항목: 태스크 배정, 메시지 라우팅, 진행 모니터링, 충돌 해결
- Do가 직접 코드 구현 또는 파일 수정 금지 (구현을 위한 Write, Edit, Bash 사용 금지)
- 태스크 배정 및 SendMessage를 통해 모든 구현을 팀원에게 위임
- 컨텍스트 파악 및 팀원 출력 검토를 위한 Read와 Grep은 허용
- 적합한 팀원이 없는 태스크의 경우 직접 구현 대신 새 팀원 소환
- delegate_mode가 false이면 팀 리더가 팀원과 함께 소규모 태스크를 직접 구현 가능

## Phase 3: 품질 검증

구현 및 테스트 태스크 완료 후:

옵션 A (quality 팀원 포함):
- quality 팀원에게 태스크 5 배정
- quality가 TRUST 5 검증 실행
- 팀 리더에게 결과 보고
- 팀 리더가 담당 팀원에게 수정 지시

옵션 B (소규모 팀에서 서브에이전트 사용):
- manager-quality 서브에이전트에 품질 검증 위임
- 결과 검토 후 필요 시 수정 태스크 생성
- 기존 팀원에게 수정 배정

품질 게이트 (모두 통과해야 함):
- 린트 오류 없음
- 타입 오류 없음
- 커버리지 목표 충족 (85%+ 전체, 90%+ 새 코드)
- critical 보안 이슈 없음
- 모든 SPEC 인수 기준 검증됨

## Phase 4: Git 작업

품질 검증 통과 후:
- manager-git 서브에이전트에 위임 (팀원이 아님)
- conventional commit 형식으로 의미 있는 커밋 생성
- 커밋 메시지에 SPEC ID 참조

## Phase 5: 정리

1. 모든 팀원 정상 종료
2. 리소스 정리를 위한 TeamDelete
3. 사용자에게 구현 요약 보고

## 폴백

어느 시점에서든 팀 모드 실패 시:
- 남은 팀원 정상 종료
- 서브에이전트 실행 워크플로우 (workflows/run.md)로 폴백
- 마지막 완료된 태스크부터 계속
- 팀 모드 실패 경고 로그

## 태스크 추적

[HARD] 모든 태스크 상태 변경은 TaskUpdate 사용:
- pending -> in_progress: 팀원이 태스크 선택 시
- in_progress -> completed: 태스크 작업이 검증될 때
- 일반 텍스트 TODO 목록 절대 사용 금지

---

Version: 1.2.0
