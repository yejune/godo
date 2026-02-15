# Do Persona 재작업 지시서

## 원칙
moai persona 파일이 템플릿이다. moai persona를 복제한 뒤 do 페르소나로 수정한다.
do-focus/.claude/에서 복사하는 것이 아니다.

## 소스
- moai persona: `/tmp/e2e4-extract/personas/moai/` (또는 새로 extract)
- moai-adk 원본: `~/Work/moai-adk/.claude/`

## 현재 상태
- personas/do/ 폴더에 파일이 있지만, 잘못된 소스(do-focus)에서 복사한 것이 섞여있음
- 올바른 소스(moai persona)에서 새로 시작해야 함

## 재작업 대상 파일

### A. agents/do/ (5개) — 전부 폐기 후 재작업
소스: `/tmp/e2e4-extract/personas/moai/agents/moai/`
1. moai/manager-ddd.md → do/manager-ddd.md
2. moai/manager-tdd.md → do/manager-tdd.md  
3. moai/manager-quality.md → do/manager-quality.md
4. moai/manager-project.md → do/manager-project.md (또는 moai/manager-spec.md 기반)
5. moai/team-quality.md → do/team-quality.md

변환 규칙:
- 파일명: moai/ → do/
- hooks 섹션: moai shell script 경로 → godo hook 직접 호출로 변환
  - `.claude/hooks/moai/handle-agent-hook.sh` → `godo hook agent-{name}`
- memory 경로: `.moai/` → `.do/`
- skills 참조: `moai-*` → `do-*`
- commands: `/moai` → `/do`
- 브랜드: `MoAI` → `Do`

### B. rules/do/workflow/ (2개) — 전부 폐기 후 재작업
소스: `/tmp/e2e4-extract/personas/moai/rules/moai/workflow/`
1. spec-workflow.md → do 버전
2. workflow-modes.md → do 버전

변환: moai SPEC 워크플로우 → do 체크리스트 워크플로우

### C. output-styles/do/ (3개) — 확인 필요
소스: `/tmp/e2e4-extract/personas/moai/output-styles/moai/`
현재 do-focus styles/에서 복사했는데, moai 원본에서 복사+변환이 맞음
- moai.md → pair.md (또는 do 스타일에 맞게)
- r2d2.md → sprint.md
- yoda.md → direct.md

### D. commands/do/ (6개) — 유지 가능
이건 do 고유 커맨드라서 moai에 대응하는 것이 다름.
moai: 2개 (99-release, github)
do: 6개 (check, checklist, mode, plan, setup, style)
→ do 커맨드는 do-focus 원본에서 복사한 것이 맞을 수 있음. 단 moai 참조 있으면 수정.

### E. characters/ (4개) — 유지
do 고유 기능. moai에 없음. 새로 작성한 것이 맞음.

### F. skills/do/ (SKILL.md + workflows + reference) — 확인 필요
orchestrator-agent가 새로 작성한 파일. moai SKILL.md를 참고해서 만들었으므로 맞을 수 있음.
단 moai 참조 잔존 확인 필요.

### G. CLAUDE.md, manifest.yaml, settings.json — 확인 필요
새로 작성한 파일. moai 참조 잔존 확인.

## 실행 순서
1. moai persona를 새로 extract (또는 기존 /tmp/e2e4-extract/personas/moai/ 사용)
2. agents/do/ 5개 폐기 → moai agents에서 복사+변환
3. rules/do/ 2개 폐기 → moai rules에서 복사+변환
4. output-styles/do/ 3개 확인 → 필요시 moai에서 재복사+변환
5. 나머지 파일 moai 참조 grep → 수정
6. YYMMDD → YY/MM/DD 전환 확인
7. E2E: assemble 테스트

## 변환 치환표
| moai 패턴 | do 치환 | 비고 |
|-----------|--------|------|
| `agents/moai/` | `agents/do/` | 경로 |
| `.claude/hooks/moai/*.sh` | godo hook 직접 호출 | 구조 변경 |
| `.moai/` | `.do/` | 디렉토리 |
| `/moai` | `/do` | 명령어 |
| `moai-` (skill prefix) | `do-` | 스킬명 |
| `MoAI` | `Do` | 브랜드명 |
| `SPEC-XXX` | checklist 기반 | 워크플로우 |
| `YYMMDD` | `YY/MM/DD` | 날짜 형식 |
