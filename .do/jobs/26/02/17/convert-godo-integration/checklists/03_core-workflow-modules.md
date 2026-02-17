# workflow-* 모듈 주입: 체크리스트/테스팅/프로젝트/TDD 규칙
상태: [ ] | 담당: general-purpose (에이전트 C) | 작성 언어: ko

## Problem Summary
- do-focus의 dev-checklist.md, dev-testing.md, dev-environment.md, dev-workflow.md에서 워크플로우 관련 규칙을 추출
- workflow-spec/modules/ 6개, workflow-testing/modules/ 2개, workflow-project/references/ 2개, workflow-tdd/modules/ 1개 — 총 11개 모듈 생성

## Acceptance Criteria
> 완료 시 [ ] → [o]로 변경. [x]는 "실패"를 의미하므로 절대 사용 금지.
- [ ] core/skills/workflow-spec/modules/checklist-system.md 생성 완료
- [ ] core/skills/workflow-spec/modules/checklist-templates.md 생성 완료
- [ ] core/skills/workflow-spec/modules/analysis-template.md 생성 완료
- [ ] core/skills/workflow-spec/modules/architecture-template.md 생성 완료
- [ ] core/skills/workflow-spec/modules/report-template.md 생성 완료
- [ ] core/skills/workflow-spec/modules/complexity-check.md 생성 완료
- [ ] core/skills/workflow-testing/modules/testing-rules.md 생성 완료
- [ ] core/skills/workflow-testing/modules/bug-fix-workflow.md 생성 완료
- [ ] core/skills/workflow-project/references/docker-rules.md 생성 완료
- [ ] core/skills/workflow-project/references/ai-forbidden-patterns.md 생성 완료
- [ ] core/skills/workflow-tdd/modules/tdd-cycle.md 생성 완료
- [ ] 각 모듈 내용이 원본 해당 섹션과 일치
- [ ] 커밋 완료

## Solution Approach
- architecture.md 4.1 매핑 테이블 기준:
  - **workflow-spec** (6개):
    - checklist-system.md ← dev-checklist.md 생성 시점/작성 방식/상태 관리/의존성/블로커 규칙
    - checklist-templates.md ← dev-checklist.md 서브 체크리스트 템플릿 전체
    - analysis-template.md ← dev-checklist.md Analysis 문서 템플릿
    - architecture-template.md ← dev-checklist.md Architecture 문서 템플릿
    - report-template.md ← dev-checklist.md 완료 보고서 템플릿 + 체크리스트 표시 의무
    - complexity-check.md ← dev-workflow.md 복잡도 판단 + 워크플로우 선택 + Analysis/Architecture 단계
  - **workflow-testing** (2개):
    - testing-rules.md ← dev-testing.md 전체
    - bug-fix-workflow.md ← dev-workflow.md 버그 수정 워크플로우
  - **workflow-project** (2개):
    - docker-rules.md ← dev-environment.md Docker 필수 + 12-Factor + 환경변수 관리 (bootapp 제외)
    - ai-forbidden-patterns.md ← dev-environment.md AI 에이전트 금지 패턴
  - **workflow-tdd** (1개):
    - tdd-cycle.md ← dev-workflow.md TDD RED-GREEN-REFACTOR 섹션
- 대안: dev-checklist.md를 하나의 큰 모듈로 유지 → 기각 (523줄은 너무 큼, 6개로 분할)

## Critical Files
- **소스**:
  - `~/Work/do-focus.workspace/do-focus/.claude/rules/do/development/checklist.md` (523줄)
  - `~/Work/do-focus.workspace/do-focus/.claude/rules/do/development/testing.md` (67줄)
  - `~/Work/do-focus.workspace/do-focus/.claude/rules/do/development/environment.md` (92줄)
  - `~/Work/do-focus.workspace/do-focus/.claude/rules/do/development/workflow.md` (184줄)
- **생성 대상**: 위 Acceptance Criteria 11개 파일

## Risks
- dev-checklist.md 523줄 분해가 가장 큰 작업 — 섹션 경계를 정확히 잡아야 함
- docker-rules.md에 bootapp 내용이 섞이지 않도록 주의 (bootapp은 페르소나 전용)
- workflow-spec/modules/에 기존 모듈과 겹치는 내용 확인 필요

## Progress Log
(작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (11개 모듈 파일)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 개선 조치:
