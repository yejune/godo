---
description: "GitHub 워크플로우 - Agent Teams로 이슈 관리 및 PR 검토"
argument-hint: "issues [--all | --label LABEL | NUMBER] | pr [--all | NUMBER]"
type: local
allowed-tools: Read, Write, Edit, Grep, Glob, Bash, AskUserQuestion, Task, TeamCreate, SendMessage, TaskCreate, TaskUpdate, TaskList, TaskGet, TeamDelete
model: sonnet
version: 1.0.0
---

## GitHub 워크플로우 구성

- **저장소**: `gh repo view --json nameWithOwner`에서 자동 감지
- **기본 모드**: Agent Teams (AGENT_TEAMS 사용 불가 시 서브에이전트로 대체)
- **브랜치 접두사**: 버그는 `fix/issue-{number}`, 기능은 `feat/issue-{number}`
- **Git 전략**: `.moai/config/sections/system.yaml`의 `github.git_workflow` 읽기

---

## 실행 지침 - 즉시 시작

이것은 GitHub 워크플로우 명령입니다. $ARGUMENTS를 파싱하고 즉시 실행하세요.

### 인자 파싱

첫 번째 단어가 하위 명령을 결정:

- **issues** (별칭: issue, fix-issues): 이슈 수정 워크플로우
- **pr** (별칭: review, pull-request): PR 코드 검토 워크플로우
- 하위 명령 없음: AskUserQuestion을 사용하여 사용자가 선택

나머지 인자는 하위 명령 인자가 됨:

- `--all`: 모든 열린 항목 처리
- `--label LABEL`: 라벨로 필터링
- `--solo`: 서브에이전트 모드 강제 (Agent Teams 건너뛰기)
- `--merge`: CI 통과 후 PR 자동 병합 (issues만)
- `NUMBER`: 특정 이슈 또는 PR 번호

---

## 실행 전 컨텍스트

!gh repo view --json nameWithOwner --jq '.nameWithOwner'
!git branch --show-current
!git status --porcelain

@.moai/config/sections/system.yaml
@.moai/config/sections/language.yaml

---

## 팀 모드 (기본값)

Agent Teams 모드는 이 워크플로우의 기본값입니다. `--team` 플래그가 필요 없습니다.

전제 조건 확인 (시작 시 실행):
1. `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` 설정 확인
2. `.moai/config/sections/workflow.yaml`의 `workflow.team.enabled: true` 확인

두 전제 조건이 모두 충족: Agent Teams 모드 사용
전제 조건 중 하나 누락 또는 `--solo` 플래그: 서브에이전트 모드로 대체

---

# 하위 명령: issues

목적: GitHub 이슈를 가져오고, 근본 원인을 분석하며, 수정을 구현하고 PR을 생성합니다.

## Issues 1단계: 이슈 발견

### 1.1단계: 열린 이슈 가져오기

GitHub에서 모든 열린 이슈 가져오기:
`gh issue list --state open --limit 50 --json number,title,labels,assignees,body,createdAt`

### 1.2단계: 이슈 선택

NUMBER 인자 제공 시:
- 특정 이슈 가져오기: `gh issue view {number} --json number,title,labels,body,comments`
- 2단계로 바로 진행

--all 또는 인자 없음:
- 형식화된 테이블로 이슈 목록 표시
- AskUserQuestion을 사용하여 수정할 이슈를 사용자가 선택
- 옵션: 개별 이슈 번호 또는 일괄 모드용 "All"

--label LABEL인 경우:
- 필터링: `gh issue list --state open --label "{LABEL}" --json number,title,labels,body`
- 필터링된 목록 표시 및 사용자가 선택

### 1.3단계: 이슈 분류

선택된 각 이슈를 유형별로 분류:
- **bug**: 기존 동작 수정 (브랜치 접두사: `fix/issue-{number}`)
- **feature**: 새로운 기능 (브랜치 접두사: `feat/issue-{number}`)
- **enhancement**: 기존 기능 개선 (브랜치 접두사: `improve/issue-{number}`)
- **docs**: 문서만 (브랜치 접두사: `docs/issue-{number}`

분류 기준: 라벨, 제목 키워드, 본문 내용 분석

## Issues 2단계: 분석

### 팀 모드 (기본값)

병렬 이슈 분석을 위한 팀 생성:

```
TeamCreate(team_name: "github-issues-{repo-slug}")
```

선택된 각 이슈에 대해 작업 생성:
```
TaskCreate: "이슈 #{number} 분석: {title}"
TaskCreate: "이슈 #{number} 수정 구현" (분석 작업에 의해 차단됨)
TaskCreate: "이슈 #{number} 수정 검증" (구현 작업에 의해 차단됨)
```

병렬 분석 팀원 생성 (이슈당 하나, 최대 3개 동시):

```
Task(
  subagent_type: "team-researcher",
  team_name: "github-issues-{repo-slug}",
  name: "analyst-{number}",
  mode: "plan",
  prompt: "GitHub 이슈 #{number} 분석.
    Title: {title}
    Body: {body}
    Comments: {comments}
    코드베이스를 탐색하여 근본 원인, 영향받는 파일, 수정 접근 방식 식별.
    TaskUpdate를 통해 작업 완료로 표시하고 SendMessage로 결과 전송."
)
```

분석 완료 후 구현 팀원 생성:

```
Task(
  subagent_type: "team-backend-dev",  // 또는 영향받는 파일에 따라 team-frontend-dev
  team_name: "github-issues-{repo-slug}",
  name: "fixer-{number}",
  mode: "acceptEdits",
  prompt: "분석 결과를 바탕으로 GitHub 이슈 #{number} 수정.
    Analysis: {analyst_findings}
    Affected files: {file_list}
    기능 브랜치 생성: {prefix}/issue-{number}
    테스트 작성, 수정 구현, 테스트 통과 확인.
    TaskUpdate를 통해 작업 완료 표시 및 SendMessage로 결과 전송."
)
```

### 서브에이전트 모드 (--solo 또는 대체)

분류에 따라 적절한 전문가 에이전트에 위임:
- 버그 수정: expert-debug 서브에이전트
- 기능: expert-backend 또는 expert-frontend 서브에이전트
- 개선: expert-refactoring 서브에이전트
- 문서: manager-docs 서브에이전트

## Issues 3단계: 브랜치 및 수정

### 3.1단계: 기능 브랜치 생성

system.yaml에서 `github.git_workflow` 읽기:

**github_flow 또는 gitflow**:
1. main에 있는지 확인 (gitflow의 경우 develop): `git checkout main && git pull origin main`
2. 브랜치 생성: `git checkout -b {prefix}/issue-{number}`

**main_direct**:
- main에 유지, 브랜치 생성 없음

### 3.2단계: 수정 검증

구현 후:
1. 테스트 실행: 언어별 테스트 명령
2. 린터 실행: 언어별 린트 명령
3. 테스트 실패 시: 에러 컨텍스트와 함께 재시도 (최대 3회)
4. 여전히 실패 시: AskUserQuestion (재시도, 건너뛰기, 중단)

### 3.3단계: 변경 커밋

manager-git 서브에이전트에 위임.

커밋 메시지 형식:
```
fix(scope): description

Fixes #{issue_number}

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
```

## Issues 4단계: PR 생성

system.yaml에서 `github.git_workflow` 읽기:

**github_flow**:
1. 푸시: `git push -u origin {prefix}/issue-{number}`
2. PR 생성: `gh pr create --title "fix: {issue title}" --body "$(body)"`
