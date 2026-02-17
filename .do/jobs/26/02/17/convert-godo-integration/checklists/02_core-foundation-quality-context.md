# foundation-quality + foundation-context 모듈 주입: 코딩 규율/파일 읽기/지식 관리
상태: [ ] | 담당: general-purpose (에이전트 B) | 작성 언어: ko

## Problem Summary
- do-focus의 dev-workflow.md + dev-environment.md + file-reading-optimization.md에서 코딩 규율, 커밋 규율, 병렬 격리, 구문 검사, 파일 읽기 최적화, 지식 관리 규칙을 추출
- foundation-quality/modules/에 5개, foundation-context/modules/에 2개 — 총 7개 모듈 생성

## Acceptance Criteria
> 완료 시 [ ] → [o]로 변경. [x]는 "실패"를 의미하므로 절대 사용 금지.
- [ ] core/skills/foundation-quality/modules/read-before-write.md 생성 완료
- [ ] core/skills/foundation-quality/modules/coding-discipline.md 생성 완료
- [ ] core/skills/foundation-quality/modules/commit-discipline.md 생성 완료
- [ ] core/skills/foundation-quality/modules/parallel-agent-isolation.md 생성 완료
- [ ] core/skills/foundation-quality/modules/syntax-check.md 생성 완료
- [ ] core/skills/foundation-context/modules/file-reading-optimization.md 생성 완료
- [ ] core/skills/foundation-context/modules/knowledge-management.md 생성 완료
- [ ] 각 모듈 내용이 원본 해당 섹션과 일치
- [ ] 커밋 완료

## Solution Approach
- architecture.md 4.1 매핑 테이블 기준:
  - read-before-write.md ← dev-workflow.md "코딩 전 필수 행동" 섹션
  - coding-discipline.md ← dev-workflow.md "코딩 규율" + "에러 대응" 섹션
  - commit-discipline.md ← dev-workflow.md "커밋 규율" 섹션
  - parallel-agent-isolation.md ← dev-workflow.md "병렬 에이전트 커밋 격리" 섹션
  - syntax-check.md ← dev-environment.md "구문 검사 필수" 섹션
  - file-reading-optimization.md ← file-reading-optimization.md 전체
  - knowledge-management.md ← dev-workflow.md "지식 관리" 섹션
- 대안: foundation-quality에 모두 넣기 → 기각 (context 관련은 foundation-context에 배치해야 Progressive Disclosure 효율적)

## Critical Files
- **소스**:
  - `~/Work/do-focus.workspace/do-focus/.claude/rules/do/development/workflow.md`
  - `~/Work/do-focus.workspace/do-focus/.claude/rules/do/development/environment.md`
  - `~/Work/do-focus.workspace/do-focus/.claude/rules/do/workflow/file-reading-optimization.md`
- **생성 대상**:
  - `core/skills/foundation-quality/modules/read-before-write.md`
  - `core/skills/foundation-quality/modules/coding-discipline.md`
  - `core/skills/foundation-quality/modules/commit-discipline.md`
  - `core/skills/foundation-quality/modules/parallel-agent-isolation.md`
  - `core/skills/foundation-quality/modules/syntax-check.md`
  - `core/skills/foundation-context/modules/file-reading-optimization.md`
  - `core/skills/foundation-context/modules/knowledge-management.md`

## Risks
- foundation-quality/modules/에 기존 모듈(trust-5-*.md 등)과 내용 중복 가능 — 기존 모듈 확인 후 추가
- coding-discipline에 에러 대응 합류 시 범위 초과 주의 — 핵심 규칙만 포함

## Progress Log
(작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (7개 모듈 파일)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 개선 조치:
