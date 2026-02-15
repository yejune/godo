# agents-philosophy: 에이전트 정의 Do 철학 반영 (Phase 6)
상태: [ ] | 담당: expert-backend

## Problem Summary
- personas/do/agents/do/ 의 5개 에이전트 정의 파일에 MoAI 잔재가 남아있다
- TRUST 5, SPEC, EARS, MoAI 설정 경로 등 MoAI 고유 개념이 혼재되어 있다
- Do의 체크리스트 기반 워크플로우, 커밋 기반 증명 철학이 반영되지 않았다
- plan-mode artifact writing rule (플랜 산출물 위치: .do/jobs/)이 반영되지 않았다

## Acceptance Criteria
- [ ] 5개 에이전트 파일에서 MoAI 잔재 완전 제거 (TRUST 5, SPEC, EARS, .moai/, XML 마커 등)
- [ ] Do 체크리스트 기반 워크플로우 반영 (살아있는 체크리스트, 상태 기호, 커밋 기반 증명)
- [ ] plan-mode artifact writing rule 반영 (.do/jobs/{YYMMDD}/{title}/ 경로)
- [ ] 에이전트 frontmatter 구문 유효 (YAML 파싱 가능)
- [ ] `grep -ri 'moai\|TRUST 5\|SPEC\|EARS\|\.moai/' personas/do/agents/do/` 결과 0건
- [ ] 커밋 완료

## Solution Approach
- 각 에이전트 파일을 읽고 MoAI 특유 용어/경로를 Do 대응 개념으로 치환
- TRUST 5 → Do 품질 기준 (테스트 통과 + 구문 검사 + 커밋 기반 증명)
- SPEC → Plan/Checklist 기반 워크플로우
- .moai/config/ → .do/jobs/ 또는 .claude/ 경로
- 대안: 에이전트 파일을 처음부터 재작성 → 기각 (기존 도메인 지식 보존이 중요하므로 점진적 치환 선택)

## Critical Files

### 항목 #18: manager-ddd.md, manager-tdd.md Do 철학 반영
- **수정 대상**: `personas/do/agents/do/manager-ddd.md` — DDD 워크플로우 에이전트, MoAI 잔재 제거
- **수정 대상**: `personas/do/agents/do/manager-tdd.md` — TDD 워크플로우 에이전트, MoAI 잔재 제거

### 항목 #19: manager-quality.md, manager-project.md, team-quality.md Do 철학 반영
- **수정 대상**: `personas/do/agents/do/manager-quality.md` — 품질 검증 에이전트, TRUST 5 → Do 품질 기준
- **수정 대상**: `personas/do/agents/do/manager-project.md` — 프로젝트 설정 에이전트, MoAI config → Do 설정
- **수정 대상**: `personas/do/agents/do/team-quality.md` — 팀 품질 에이전트, TRUST 5 → Do 품질 기준

## Risks
- 에이전트 정의의 도메인 전문성(DDD/TDD 워크플로우 지식)을 MoAI 제거 중 실수로 삭제할 수 있음: 변경 전 원본 백업, diff로 검증
- frontmatter 형식이 깨질 수 있음: 수정 후 YAML 파싱 검증

## Progress Log
- (작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (agents/do/ 5개 파일)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 다음에 다르게 할 점:
