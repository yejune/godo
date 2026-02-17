# Do 설치 상태 확인

다음을 확인하고 결과를 보고해주세요:

## 1. 핵심 파일 존재 여부
- [ ] `.claude/settings.json` 존재
- [ ] `.claude/hooks/do/` 디렉토리 존재
- [ ] `.claude/skills/` 디렉토리 존재
- [ ] `.claude/agents/` 디렉토리 존재
- [ ] `CLAUDE.md` 존재

## 2. Hook 실행 테스트
```bash
CLAUDE_PROJECT_DIR=$(pwd) uv run .claude/hooks/do/session_start__show_project_info.py
```
- 정상 JSON 출력 여부 확인

## 3. 의존성 확인
- [ ] `uv` 설치 여부
- [ ] Python 3.9+ 사용 가능 여부

## 4. 결과 요약
모든 항목을 체크하고 다음 형식으로 보고:

```
Do 설치 상태: ✅ 정상 / ⚠️ 일부 문제 / ❌ 미설치

핵심 파일: X/5
Hook 실행: 성공/실패
의존성: 정상/누락

문제 발견 시 해결 방법 제시
```
