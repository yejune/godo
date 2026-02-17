# 페르소나 워크플로우 생성: DO 스타일 Plan/Run/Report/Team 워크플로우
상태: [ ] | 담당: general-purpose (에이전트 D) | 작성 언어: ko

## Problem Summary
- DO 페르소나의 워크플로우를 5개의 씬(thin) 파일로 생성
- 각 워크플로우는 코어 패턴을 참조하고 DO만의 차이점만 선언 (SPEC 대신 문서 체인, .moai/ 대신 .do/jobs/)
- bootapp.md도 dev-environment.md의 bootapp 섹션을 반영하여 갱신

## Acceptance Criteria
> 완료 시 [ ] → [o]로 변경. [x]는 "실패"를 의미하므로 절대 사용 금지.
- [ ] personas/do/workflows/ 디렉토리 생성 완료
- [ ] personas/do/workflows/plan.md 생성 — 트리거 매핑, 문서 체인, 산출물 위치
- [ ] personas/do/workflows/run.md 생성 — 에이전트 실행 사이클 참조, 체크리스트 기반, 멱등 재개
- [ ] personas/do/workflows/report.md 생성 — 자동 트리거, Lessons Learned 종합
- [ ] personas/do/workflows/team-plan.md 생성 — 코어 plan-research 패턴 참조, 팀 이름 브랜드
- [ ] personas/do/workflows/team-run.md 생성 — 코어 implementation 패턴 참조, 체크리스트 기반
- [ ] personas/do/rules/bootapp.md 갱신 — dev-environment.md bootapp 네트워크 섹션 반영
- [ ] 커밋 완료

## Solution Approach
- architecture.md 섹션 4.2 내용 정의를 기준으로 5개 워크플로우 파일 작성
- 핵심 원칙: 코어 참조 + DO 차이점만 선언 (thin files)
  - plan.md: 한국어 트리거 매핑 + .do/jobs/ 경로 + 문서 체인 (analysis→architecture→plan→checklist)
  - run.md: READ-CLAIM-WORK-VERIFY-RECORD-COMMIT 사이클 참조 + DDD/TDD/Hybrid 모드 선택
  - report.md: 모든 체크리스트 [o] 시 자동 트리거 + report.md 템플릿 참조
  - team-plan.md: do-plan-{slug} 팀 이름 + 문서 체인 산출물
  - team-run.md: do-run-{slug} 팀 이름 + 체크리스트 기반 작업 분배
- bootapp.md: dev-environment.md의 네트워크 섹션(bootapp 도메인, .test TLD, SSL_DOMAINS 등) 반영
- 대안: MoAI 워크플로우 파일을 복사 후 수정 → 기각 (DO는 SPEC 없음, 구조가 근본적으로 다름)

## Critical Files
- **소스**:
  - `architecture.md` 섹션 4.2 (워크플로우 내용 정의)
  - `handoff.md` 섹션 2 (DO Methodology Decisions)
  - `~/Work/do-focus.workspace/do-focus/.claude/rules/do/development/environment.md` (bootapp 섹션)
- **생성 대상**:
  - `personas/do/workflows/plan.md`
  - `personas/do/workflows/run.md`
  - `personas/do/workflows/report.md`
  - `personas/do/workflows/team-plan.md`
  - `personas/do/workflows/team-run.md`
- **수정 대상**:
  - `personas/do/rules/bootapp.md`

## Risks
- 워크플로우 파일이 너무 얇으면 에이전트가 참조할 코어 패턴을 찾지 못할 수 있음 — 코어 스킬 이름을 명시적으로 참조
- bootapp.md 갱신 시 기존 내용 유실 주의 — 기존 내용 읽은 후 병합

## Progress Log
(작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (워크플로우 5개 + bootapp.md)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 개선 조치:
