---
description: "GitHub 워크플로우 - Agent Teams로 이슈 관리 및 PR 리뷰"
argument-hint: "issues [--all | --label LABEL | NUMBER] | pr [--all | NUMBER]"
type: local
allowed-tools: Read, Write, Edit, Grep, Glob, Bash, AskUserQuestion, Task, TeamCreate, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet, TeamDelete
model: sonnet
version: 1.0.0
---

## GitHub 워크플로우 설정

- **저장소**: `gh repo view --json nameWithOwner`에서 자동 감지
- **기본 모드**: Agent Teams (AGENT_TEAMS 불가 시 서브에이전트로 폴백)
- **브랜치 접두사**: 버그는 `fix/issue-{number}`, 기능은 `feat/issue-{number}`
- **Git 전략**: `.do/config/sections/system.yaml`의 `github.git_workflow` 읽기

---

## 실행 지시 - 즉시 시작

이것은 GitHub 워크플로우 명령입니다. $ARGUMENTS를 파싱하고 즉시 실행하세요.

### 인자 파싱

첫 번째 단어가 서브커맨드를 결정합니다:

- **issues** (별칭: issue, fix-issues): 이슈 수정 워크플로우
- **pr** (별칭: review, pull-request): PR 코드 리뷰 워크플로우
- 서브커맨드 없음: AskUserQuestion으로 선택

나머지 인자는 서브커맨드 인자가 됩니다:

- `--all`: 열린 항목 모두 처리
- `--label LABEL`: 레이블로 필터
- `--solo`: 서브에이전트 모드 강제 (Agent Teams 건너뜀)
- `--merge`: CI 통과 후 PR 자동 머지 (issues만)
- `NUMBER`: 특정 이슈 또는 PR 번호 대상

---

## 사전 실행 컨텍스트

!gh repo view --json nameWithOwner --jq '.nameWithOwner'
!git branch --show-current
!git status --porcelain

@.do/config/sections/system.yaml
@.do/config/sections/language.yaml

---

## 팀 모드 (기본값)

Agent Teams 모드가 이 워크플로우의 기본값입니다. `--team` 플래그 불필요.

시작 시 사전 요건 확인:
1. `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`이 설정되어 있는지 확인
2. `.do/config/sections/workflow.yaml`에 `workflow.team.enabled: true`인지 확인

두 사전 요건 충족 시: Agent Teams 모드 사용
사전 요건 미충족 또는 `--solo` 플래그: 서브에이전트 모드로 폴백

---

# 서브커맨드: issues

목적: GitHub 이슈를 가져오고, 근본 원인을 분석하고, 수정을 구현하고, PR을 생성합니다.

## Issues Phase 1: 이슈 탐색

### 1.1단계: 열린 이슈 가져오기

GitHub에서 모든 열린 이슈 가져오기:
`gh issue list --state open --limit 50 --json number,title,labels,assignees,body,createdAt`

### 1.2단계: 이슈 선택

NUMBER 인자가 있으면:
- 특정 이슈 가져오기: `gh issue view {number} --json number,title,labels,body,comments`
- Phase 2로 바로 진행

--all 또는 인자 없으면:
- 이슈 목록을 형식화된 테이블로 표시
- AskUserQuestion으로 수정할 이슈 선택
- 옵션: 개별 이슈 번호, 또는 일괄 모드를 위한 "All"

--label LABEL이면:
- 필터: `gh issue list --state open --label "{LABEL}" --json number,title,labels,body`
- 필터된 목록 표시 후 선택

### 1.3단계: 이슈 분류

선택된 각 이슈에 대해 유형 분류:
- **bug**: 기존 동작 수정 (브랜치 접두사: `fix/issue-{number}`)
- **feature**: 새 기능 (브랜치 접두사: `feat/issue-{number}`)
- **enhancement**: 기존 기능 개선 (브랜치 접두사: `improve/issue-{number}`)
- **docs**: 문서만 (브랜치 접두사: `docs/issue-{number}`)

## Issues Phase 2: 분석

### 팀 모드 (기본값)

병렬 이슈 분석을 위한 팀 생성:
분석 팀원 병렬 생성 (이슈당 하나, 최대 3개 동시)

### 서브에이전트 모드 (--solo 또는 폴백)

분류에 따라 적절한 전문 에이전트에 위임:
- 버그 수정: expert-debug 서브에이전트
- 기능: expert-backend 또는 expert-frontend 서브에이전트
- 개선: expert-refactoring 서브에이전트
- 문서: manager-docs 서브에이전트

## Issues Phase 3: 브랜치 및 수정

### 3.1단계: 기능 브랜치 생성

system.yaml에서 `github.git_workflow` 읽기:

**github_flow 또는 gitflow**:
1. main 확인: `git checkout main && git pull origin main`
2. 브랜치 생성: `git checkout -b {prefix}/issue-{number}`

**main_direct**:
- main에 머물기, 브랜치 생성 없음

### 3.2단계: 수정 검증

구현 후:
1. 테스트 실행: 언어별 테스트 명령
2. 린터 실행: 언어별 린트 명령
3. 테스트 실패 시: 오류 컨텍스트로 재시도 (최대 3회)
4. 계속 실패 시: AskUserQuestion (재시도, 건너뛰기, 중단)

### 3.3단계: 변경사항 커밋

manager-git 서브에이전트에 위임.

커밋 메시지 형식:
```
fix(scope): description

Fixes #{issue_number}
```

## Issues Phase 4: PR 생성

system.yaml에서 `github.git_workflow` 읽기:

**github_flow**:
1. 푸시: `git push -u origin {prefix}/issue-{number}`
2. PR 생성: `gh pr create --title "fix: {issue title}" --body "$(body)"`

**main_direct**:
1. main에 직접 푸시

## Issues Phase 5: 정리 및 보고

팀 모드 사용 시:
1. SendMessage로 모든 팀원 종료 요청
2. 리소스 정리를 위해 TeamDelete

일괄 요약 표시

---

# 서브커맨드: pr

목적: PR을 가져오고, 다각도 코드 리뷰를 수행하고, 리뷰 댓글을 제출합니다.

## PR Phase 1: PR 탐색

### 1.1단계: 열린 PR 가져오기

`gh pr list --state open --limit 30 --json number,title,author,labels,additions,deletions,changedFiles,headRefName`

### 1.2단계: PR 선택

NUMBER 인자가 있으면:
- Phase 2로 진행

--all 또는 인자 없으면:
- PR 목록을 형식화된 테이블로 표시
- AskUserQuestion으로 리뷰할 PR 선택

## PR Phase 2: 코드 리뷰

### 팀 모드 (기본값)

병렬 다각도 분석을 위한 리뷰 팀 생성.
3명의 리뷰어 병렬 생성:
- **security-reviewer**: SQL 인젝션, XSS, 인증/인가 이슈, OWASP Top 10 준수 검토
- **perf-reviewer**: 알고리즘 복잡도, DB 쿼리 패턴, 메모리 누수, 동시성 이슈 검토
- **quality-reviewer**: 코드 정확성, 테스트 커버리지, 네이밍 컨벤션, 에러 처리 검토

### 서브에이전트 모드 (--solo 또는 폴백)

순차적으로 위임:
1. expert-security 서브에이전트: PR diff 보안 분석
2. expert-performance 서브에이전트: 성능 분석
3. manager-quality 서브에이전트: 코드 품질 리뷰

## PR Phase 3: 종합 및 리뷰 제출

모든 리뷰어 완료 후:

1. 모든 관점의 결과 수집
2. 심각도별 이슈 분류:
   - **Critical**: 머지 전 필수 수정 (보안 취약점, 데이터 손실 위험)
   - **Important**: 수정 권장 (성능 이슈, 누락된 에러 처리)
   - **Suggestion**: 있으면 좋음 (네이밍, 스타일, 소소한 개선)
3. GitHub 리뷰 형식으로 포맷

### 리뷰 제출

AskUserQuestion으로 리뷰 액션 확인:
- 승인: 요약과 함께 승인 제출
- 변경 요청: 필수 변경사항과 함께 제출
- 댓글만: 승인 결정 없이 댓글로 제출
- 건너뜀: 리뷰 제출 안 함

---

## 공통 규칙

- **[HARD] 에이전트 위임**: 모든 분석 및 수정은 에이전트에 위임해야 함
- **[HARD] 사용자 승인**: 이슈 수정 및 리뷰 제출은 사용자 확인 필요
- **팀 모드 기본값**: Agent Teams가 기본으로 사용, `--solo`로 재정의
- **Git 전략 인식**: system.yaml에서 `github.git_workflow` 읽기
- **이슈 연결**: 커밋/PR에 항상 `Fixes #{number}` 포함
- **이슈당 브랜치**: 각 이슈는 자체 브랜치 (main_direct 제외)
- **테스트 검증**: 모든 수정은 PR 생성 전 테스트 통과 필수

---

## 실행 시작

$ARGUMENTS를 파싱하여 서브커맨드 (issues 또는 pr)를 결정한 후 해당 워크플로우 단계를 즉시 실행하세요.
