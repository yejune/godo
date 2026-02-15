# rules-update: 개발 규칙 갱신 (Phase 7)
상태: [ ] | 담당: expert-backend

## Problem Summary
- personas/do/rules/do/workflow/ 의 2개 규칙 파일에 MoAI 잔재가 남아있다
- SPEC 워크플로우, EARS 형식, MoAI config 경로 등 MoAI 고유 개념이 혼재되어 있다
- Do의 커밋 기반 증명(commit-as-proof), 살아있는 체크리스트, 멱등성 규칙이 반영되지 않았다
- Do 워크플로우 (Plan → Checklist → Develop → Test → Report) 정의가 필요하다

## Acceptance Criteria
- [ ] 2개 규칙 파일에서 MoAI 잔재 완전 제거 (SPEC, EARS, .moai/, moai 명령어 등)
- [ ] Do 워크플로우 반영 (Plan → Checklist → Develop → Test → Report)
- [ ] 커밋 기반 증명, 살아있는 체크리스트 개념 반영
- [ ] `grep -ri 'moai\|SPEC\|EARS\|\.moai/' personas/do/rules/do/` 결과 0건
- [ ] 커밋 완료

## Solution Approach
- spec-workflow.md: SPEC 3-Phase (Plan/Run/Sync)를 Do 워크플로우 (Plan → Checklist → Develop → Test → Report)로 재작성
- workflow-modes.md: DDD/TDD/Hybrid 모드 설명은 보존하되 SPEC 참조, MoAI 설정 경로 제거
- /moai 명령어 → /do 명령어로 치환
- 대안: 규칙 파일을 완전 삭제하고 CLAUDE.md에 통합 → 기각 (규칙 분리가 유지보수에 유리)

## Critical Files

### 항목 #20: rules/ 개발 규칙 갱신
- **수정 대상**: `personas/do/rules/do/workflow/spec-workflow.md` — SPEC 워크플로우 → Do 워크플로우로 재작성
- **수정 대상**: `personas/do/rules/do/workflow/workflow-modes.md` — DDD/TDD/Hybrid 모드, MoAI 참조 제거

## Risks
- spec-workflow.md를 재작성할 때 토큰 예산 관리 등 유용한 정보를 삭제할 수 있음: Do에서도 유효한 개념은 보존
- workflow-modes.md의 DDD/TDD 워크사이클이 Do에서도 동일하게 적용되는지 확인 필요

## Progress Log
- (작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (rules/do/workflow/ 2개 파일)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 다음에 다르게 할 점:
