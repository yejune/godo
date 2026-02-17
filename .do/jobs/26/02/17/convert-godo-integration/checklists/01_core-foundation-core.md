# foundation-core 모듈 주입: 에이전트 실행/위임/리서치 규칙
상태: [ ] | 담당: general-purpose (에이전트 A) | 작성 언어: ko

## Problem Summary
- do-focus의 dev-workflow.md에서 에이전트 실행 사이클, 위임 규칙, 리서치 위임 규칙을 추출하여 convert의 코어 스킬 foundation-core/modules/에 3개 모듈로 주입
- 현재 이 내용은 do-focus 전용 규칙 파일에만 존재 — 코어에 없으면 다른 페르소나가 사용 불가

## Acceptance Criteria
> 완료 시 [ ] → [o]로 변경. [x]는 "실패"를 의미하므로 절대 사용 금지.
- [ ] core/skills/foundation-core/modules/agent-execution-cycle.md 생성 완료
- [ ] core/skills/foundation-core/modules/agent-delegation.md 생성 완료
- [ ] core/skills/foundation-core/modules/agent-research.md 생성 완료
- [ ] 각 모듈의 내용이 원본 dev-workflow.md 해당 섹션과 일치 (줄 수 비교)
- [ ] 커밋 완료

## Solution Approach
- do-focus의 dev-workflow.md를 읽어서 architecture.md 4.1 매핑 테이블 기준으로 내용 추출
- agent-execution-cycle.md ← 에이전트 위임 필수 전달사항 + 실행 사이클 (READ-CLAIM-WORK-VERIFY-RECORD-COMMIT)
- agent-delegation.md ← 에이전트 중단 & 재개 + 멱등 재개 규칙
- agent-research.md ← 에이전트 리서치 위임 규칙 (general-purpose 사용, Explore 제한 등)
- 대안: 3개를 하나의 큰 파일로 합치기 → 기각 (Progressive Disclosure 원칙: 모듈 단위 분리가 토큰 효율적)

## Critical Files
- **소스**: `~/Work/do-focus.workspace/do-focus/.claude/rules/do/development/workflow.md` — 에이전트 실행/위임/리서치 섹션
- **생성 대상**:
  - `core/skills/foundation-core/modules/agent-execution-cycle.md`
  - `core/skills/foundation-core/modules/agent-delegation.md`
  - `core/skills/foundation-core/modules/agent-research.md`
- **참조**: `architecture.md` 섹션 4.1 매핑 테이블

## Risks
- dev-workflow.md의 줄 범위가 정확하지 않을 수 있음 — 실제 내용 기준으로 분할할 것
- foundation-core/modules/에 이미 유사 모듈이 있을 수 있음 (delegation-patterns.md 등) — 중복 확인 필수

## Progress Log
(작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (3개 모듈 파일)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 개선 조치:
