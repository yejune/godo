# skills-philosophy: 스킬 Do 철학 반영 (Phase 8)
상태: [ ] | 담당: expert-backend

## Problem Summary
- personas/do/skills/do/ 의 8개 스킬 파일에 MoAI 잔재가 남아있다
- TRUST 5, SPEC, XML 마커 (<moai>DONE</moai> 등), MoAI 설정 경로 등이 혼재되어 있다
- Progressive Disclosure 패턴이 Do 철학에 맞게 정비되지 않았다
- 스킬 YAML frontmatter에 MoAI 전용 메타데이터가 포함되어 있을 수 있다

## Acceptance Criteria
- [ ] 8개 스킬 파일에서 MoAI 잔재 완전 제거 (TRUST 5, SPEC, XML 마커, .moai/ 등)
- [ ] Progressive Disclosure 패턴 Do 철학에 맞게 반영
- [ ] 스킬 YAML frontmatter 파싱 가능 (구문 유효)
- [ ] `grep -ri 'moai\|TRUST 5\|SPEC\|EARS\|<moai>\|\.moai/' personas/do/skills/do/` 결과 0건
- [ ] 커밋 완료

## Solution Approach
- SKILL.md: 스킬 시스템 개요에서 MoAI 참조 제거, Do 스킬 구조로 재정의
- reference.md: 참조 문서에서 MoAI 개념 치환
- workflows/*.md: 각 워크플로우 스킬에서 SPEC → Do Plan/Checklist 워크플로우로 치환
- XML 마커 (<moai>DONE</moai>) → Do 완료 신호 방식으로 변경
- 대안: 스킬 파일을 모두 삭제하고 새로 작성 → 기각 (Progressive Disclosure 구조와 워크플로우 지식 보존)

## Critical Files

### 항목 #21: SKILL.md + reference.md 스킬 철학 반영
- **수정 대상**: `personas/do/skills/do/SKILL.md` — 스킬 시스템 개요, MoAI 잔재 제거
- **수정 대상**: `personas/do/skills/do/references/reference.md` — 참조 문서, MoAI 개념 치환

### 항목 #22: workflows/ 스킬 파일 Do 철학 반영
- **수정 대상**: `personas/do/skills/do/workflows/do.md` — Do 메인 워크플로우 스킬
- **수정 대상**: `personas/do/skills/do/workflows/team-do.md` — Team 워크플로우 스킬
- **수정 대상**: `personas/do/skills/do/workflows/report.md` — 리포트 워크플로우 스킬
- **수정 대상**: `personas/do/skills/do/workflows/run.md` — Run 워크플로우 스킬
- **수정 대상**: `personas/do/skills/do/workflows/plan.md` — Plan 워크플로우 스킬
- **수정 대상**: `personas/do/skills/do/workflows/test.md` — Test 워크플로우 스킬

## Risks
- Progressive Disclosure의 token 계산이 내용 변경 후 부정확해질 수 있음: frontmatter의 토큰 수치 재계산 필요
- 워크플로우 스킬의 핵심 로직(DDD/TDD 사이클, Plan/Run 전환 등)을 MoAI 제거 중 실수로 삭제할 수 있음: 변경 전/후 diff 검증

## Progress Log
- (작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (skills/do/ 8개 파일)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 다음에 다르게 할 점:
