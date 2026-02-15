# commands-update: 커맨드 확인 및 갱신 (Phase 9)
상태: [ ] | 담당: expert-backend

## Problem Summary
- personas/do/commands/do/ 의 5개 커맨드 파일이 Do 워크플로우에 맞게 갱신되지 않았다
- MoAI 잔재 (moai 명령어, .moai/ 경로, SPEC 참조 등)가 포함되어 있을 수 있다
- Do 워크플로우 (Plan → Checklist → Develop → Test → Report)와 커맨드가 정합하지 않을 수 있다

## Acceptance Criteria
- [ ] 5개 커맨드 파일에서 MoAI 잔재 완전 제거 (moai 명령어, .moai/ 경로, SPEC 등)
- [ ] Do 워크플로우에 맞는 커맨드 구조 확인 (/do:plan, /do:checklist, /do:mode, /do:style, /do:setup)
- [ ] `grep -ri 'moai\|SPEC\|EARS\|\.moai/' personas/do/commands/do/` 결과 0건
- [ ] 커밋 완료

## Solution Approach
- 각 커맨드 파일을 읽고 MoAI 참조를 Do 대응 개념으로 치환
- /moai 명령어 → /do 명령어 (또는 godo 명령어)
- .moai/config/ → .claude/ 또는 .do/ 경로
- SPEC → Plan/Checklist 워크플로우
- 대안: 커맨드 파일을 처음부터 재작성 → 기각 (기존 커맨드 구조 보존이 중요)

## Critical Files

### 항목 #23: commands/ 확인 및 갱신
- **수정 대상**: `personas/do/commands/do/setup.md` — 프로젝트 초기화 커맨드, MoAI 잔재 제거
- **수정 대상**: `personas/do/commands/do/style.md` — 스타일 전환 커맨드, MoAI 잔재 제거
- **수정 대상**: `personas/do/commands/do/mode.md` — 모드 전환 커맨드, MoAI 잔재 제거
- **수정 대상**: `personas/do/commands/do/checklist.md` — 체크리스트 커맨드, MoAI 잔재 제거
- **수정 대상**: `personas/do/commands/do/plan.md` — 플랜 커맨드, MoAI 잔재 제거

### 항목 #24: Phase 6-9 전체 검증
- **검증 대상**: `personas/do/agents/do/` — 전체 에이전트 파일
- **검증 대상**: `personas/do/rules/do/` — 전체 규칙 파일
- **검증 대상**: `personas/do/skills/do/` — 전체 스킬 파일
- **검증 대상**: `personas/do/commands/do/` — 전체 커맨드 파일

## Risks
- 커맨드 파일의 $ARGUMENTS 플레이스홀더나 slash command 구문이 깨질 수 있음: 수정 후 구문 검증
- mode.md의 모드 전환 로직이 godo와 정합하는지 확인 필요

## Progress Log
- (작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (commands/do/ 5개 파일)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 다음에 다르게 할 점:
