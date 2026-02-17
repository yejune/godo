# 페르소나 오버라이드 정리: 스킬 삭제 + manifest 갱신
상태: [ ] | 담당: general-purpose (에이전트 E) | 작성 언어: ko

## Problem Summary
- Phase 1에서 코어에 규칙을 주입했으므로 DO 페르소나의 오버라이드 스킬은 불필요
- 9개 오버라이드 스킬 디렉토리 삭제 + rules/workflow/ 삭제 + manifest.yaml 갱신
- handoff.md 결정사항: "Override skills → ELIMINATED"

## Acceptance Criteria
> 완료 시 [ ] → [o]로 변경. [x]는 "실패"를 의미하므로 절대 사용 금지.
- [ ] personas/do/skills/ 오버라이드 스킬 9개 디렉토리 삭제 완료:
  - do-foundation-core, do-foundation-quality
  - do-workflow-ddd, do-workflow-plan, do-workflow-project
  - do-workflow-spec, do-workflow-tdd, do-workflow-team, do-workflow-testing
- [ ] personas/do/skills/do/ (오케스트레이터 스킬) 유지 확인
- [ ] personas/do/rules/workflow/ 디렉토리 삭제 완료
- [ ] personas/do/manifest.yaml 갱신 완료:
  - workflows 섹션 추가 (5개 워크플로우 경로)
  - 오버라이드 스킬 참조 제거
  - rules 섹션 축소 (bootapp.md만)
- [ ] 커밋 완료

## Solution Approach
- 삭제 대상 디렉토리 목록 (architecture.md 섹션 4.3):
  ```
  personas/do/skills/do-foundation-core/
  personas/do/skills/do-foundation-quality/
  personas/do/skills/do-workflow-ddd/
  personas/do/skills/do-workflow-plan/
  personas/do/skills/do-workflow-project/
  personas/do/skills/do-workflow-spec/
  personas/do/skills/do-workflow-tdd/
  personas/do/skills/do-workflow-team/
  personas/do/skills/do-workflow-testing/
  personas/do/rules/workflow/
  ```
- manifest.yaml 갱신 (architecture.md 섹션 6.3):
  - workflows 추가: plan.md, run.md, report.md, team-plan.md, team-run.md
  - rules 축소: bootapp.md만 유지
  - skills에서 오버라이드 스킬 제거 (do/ 오케스트레이터만 유지)
- 대안: 오버라이드 스킬을 deprecated 마크만 하기 → 기각 (코어에 이미 주입 완료, 중복 유지는 혼란)

## Critical Files
- **삭제 대상**: 위 10개 디렉토리
- **수정 대상**: `personas/do/manifest.yaml`
- **유지 확인**: `personas/do/skills/do/` (오케스트레이터 스킬)
- **참조**: `architecture.md` 섹션 4.3, 6.3

## Risks
- 오버라이드 스킬 중 코어에 미반영된 내용이 있을 수 있음 — 삭제 전 각 오버라이드 내용과 코어 모듈 대조
- manifest.yaml 구조가 예상과 다를 수 있음 — 현재 manifest 읽은 후 수정
- do/ 오케스트레이터 스킬을 실수로 삭제하면 치명적 — rm 대상 경로 정확히 확인

## Progress Log
(작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (삭제된 파일 + manifest.yaml)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 개선 조치:
