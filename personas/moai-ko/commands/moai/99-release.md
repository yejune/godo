---
description: "MoAI-ADK v2.x 프로덕션 릴리스 - 에이전트 위임 및 품질 검증 포함. 타겟 브랜치는 항상 main. 태그 형식 vX.Y.Z는 GoReleaser 트리거. 모든 git 작업은 manager-git에 위임. 품질 실패는 expert-debug로 에스컬레이션."
argument-hint: "[VERSION] - 선택적 타겟 버전 (예: 2.1.0). 생략 시 patch/minor/major 선택을 묻습니다."
type: local
allowed-tools: Read, Write, Edit, Grep, Glob, Bash, TodoWrite, AskUserQuestion, Task
model: sonnet
version: 3.0.0
metadata:
  release_target: "production"
  branch: "main"
  tag_format: "vX.Y.Z"
  changelog_format: "korean_first"
  release_notes_format: "bilingual"
  git_delegation: "required"
  quality_escalation: "expert-debug"
---

## 릴리스 구성

- **브랜치**: `main`
- **태그 형식**: `vX.Y.Z` (표준 semver, `.github/workflows/release.yml`을 통해 GoReleaser 트리거)
- **릴리스 URL**: https://github.com/modu-ai/moai-adk/releases/tag/vX.Y.Z
- **바이너리**: darwin-arm64, darwin-amd64, linux-arm64, linux-amd64, windows-amd64

---

## 실행 지침 - 즉시 시작

이것은 릴리스 명령입니다. 아래 워크플로우를 순서대로 실행하세요. 단계를 설명만 하지 말고 실제로 명령을 실행하세요.

제공된 인자: $ARGUMENTS

- VERSION 인자 제공 시: 타겟 버전으로 사용, 버전 선택 건너뜀기
- 인자 없음: 사용자에게 버전 유형(patch/minor/major) 선택 요청

---

## 실행 전 컨텍스트

!git status --porcelain
!git branch --show-current
!git tag --list --sort=-v:refname | head -5
!git log --oneline -10

@go.mod
@pkg/version/version.go

---

## 0단계: 사전 점검

릴리스 프로세스 시작 전, 작업 디렉토리가 깨끗한지 확인:

1. **커밋되지 않은 변경 확인**:
   ```bash
   git status --porcelain
   ```

2. **커밋되지 않은 파일 처리**:
   - `.claude/`에 추적되지 않은 파일이 있으면: 커밋할지 삭제할지 확인
   - 정리 명령:
   ```bash
   git checkout -- .claude/
   git clean -fd .claude/ internal/template/templates/.claude/
   ```

3. **브랜치 확인**:
   - `main` 브랜치에 있어야 함
   - 아니면 main 체크아웃: `git checkout main`

4. **최신 변경 풀**:
   ```bash
   git pull origin main
   ```

---

## 1단계: 품질 게이트

이 항목으로 TodoWrite를 만들고, 가능한 곳에서 병렬로 각 검사 실행:

1. 모든 테스트 실행: `go test -race ./... -count=1 2>&1 | tail -30`
2. go vet 실행: `go vet ./... 2>&1 | tail -10`
3. go fmt 확인: `gofumpt -l . 2>/dev/null | head -10`

포맷팅 이슈 발견 시 `make fmt`로 수정하고 커밋:
`git add -A && git commit -m "style: auto-fix formatting issues"`

품질 요약 표시:

- tests: PASS 또는 FAIL (FAIL이면 중단하고 보고)
- go vet: PASS 또는 WARNING
- gofmt: PASS 또는 FIXED

### 에러 처리

품질 게이트 실패 시:

- **expert-debug 서브에이전트 사용**하여 이슈 진단 및 해결
- 모든 게이트 통과 후에만 릴리스 워크플로우 재개

---

## 2단계: 코드 검토

마지막 태그 이후 커밋 가져오기:
`git log $(git describe --tags --abbrev=0 2>/dev/null || echo HEAD~20)..HEAD --oneline`

diff 통계 가져오기:
`git diff $(git describe --tags --abbrev=0 2>/dev/null || echo HEAD~20)..HEAD --stat`

다음에 대해 변경 사항 분석:

- 버그 잠재성
- 보안 이슈
- 호환성 깨지는 변경
- 테스트 커버리지 격차

권장사항과 함께 검토 보고서 표시: PROCEED 또는 REVIEW_NEEDED

---

## 3단계: 버전 선택

VERSION 인자가 제공된 경우 (예: "2.1.0"):

- 해당 버전을 직접 사용
- AskUserQuestion 건너뛰기

VERSION 인자 없는 경우:

- `pkg/version/version.go`에서 현재 버전 읽기
- AskUserQuestion으로 질문: patch/minor/major

새 버전 계산 및 모든 버전 파일 업데이트:

1. `pkg/version/version.go`의 `Version` 변수 편집
2. `.moai/config/sections/system.yaml`의 `moai.version` 편집
3. `internal/template/templates/.moai/config/sections/system.yaml`의 `moai.version` 편집
4. 커밋: `git add pkg/version/version.go .moai/config/sections/system.yaml internal/template/templates/.moai/config/sections/system.yaml && git commit -m "chore: bump version to vX.Y.Z"`

버전 파일 체크리스트:
- [ ] pkg/version/version.go: Version = "vX.Y.Z"
- [ ] .moai/config/sections/system.yaml: moai.version: "X.Y.Z"
- [ ] internal/template/templates/.moai/config/sections/system.yaml: moai.version: "X.Y.Z"

---

## 4단계: CHANGELOG 생성 (이국어: 영어 우선)

### [HARD] 영어 우선 이국어 형식

CHANGELOG.md와 GitHub 릴리스 노트는 영어 우선 이국어 구조를 따라야 합니다. 이는 국제 사용자가 영어를 먼저 보면서 한국어 문서를 유지합니다.

changelog용 커밋 가져오기: `git log $(git describe --tags --abbrev=0)..HEAD --pretty=format:"- %s (%h)"`

### CHANGELOG.md 구조

이 구조로 CHANGELOG.md 맨 앞에 새 버전 항목 추가:

```
## [X.Y.Z] - YYYY-MM-DD

### 요약
[영어: 핵심 기능과 개선 사항을 2-3줄로 요약]

### 주요 변경 사항 (Breaking Changes)
[영어: 호환성을 깨는 변경 사항 목록, 없으면 "없음"]

### 추가된 기능 (Added)
- [영어 추가 1]
- [영어 추가 2]

### 변경된 기능 (Changed)
- [영어 변경 1]
- [영어 변경 2]

### 수정된 버그 (Fixed)
- [영어 수정 1]
- [영어 수정 2]

### 설치 및 업데이트 (Installation & Update)

\`\`\`bash
# 최신 버전으로 업데이트
moai update

# 버전 확인
moai version
\`\`\`

---

## [X.Y.Z] - YYYY-MM-DD (한국어)

### Summary
[한국어: 핵심 기능과 개선 사항을 2-3줄로 요약]

### 주요 변경 사항 (Breaking Changes)
[한국어: 호환성을 깨는 변경 사항 목록, 없으면 "없음"]

### 추가된 기능 (Added)
- [한국어 추가 1]
- [한국어 추가 2]

### 변경된 기능 (Changed)
- [한국어 변경 1]
- [한국어 변경 2]

### 수정된 버그 (Fixed)
- [한국어 수정 1]
- [한국어 수정 2]

### 설치 및 업데이트 (Installation & Update)

\`\`\`bash
# 최신 버전으로 업데이트
moai update

# 버전 확인
moai version
\`\`\`
```

### [HARD] 영어 우선 규칙

각 섹션에서 **영어를 먼저 작성**하고 그 다음 한국어를 제공:
- 영어 섹션: 1-3줄 간결한 요약
- 한국어 섹션: 영어와 동일한 내용을 한국어로 번역
- 한국어 섹션의 제목만 한국어로 (예: "### 요약", "### 추가된 기능")

이유: 국제 사용자에게 영어 우선 접근성 제공하면서 한국어 문서 유지

---

## 5단계: Git 태그 및 푸시

manager-git 서브에이전트에 위임:

```bash
git tag -a vX.Y.Z -m "Release vX.Y.Z"
git push origin main
git push origin vX.Y.Z
```

---

## 6단계: GoReleaser 트리거

태그 푸시는 `.github/workflows/release.yml` 워크플로우를 트리거합니다:
- 다중 플랫폼 바이너리 빌드
- GitHub 릴리스에 자동 업로드
- Homebrew tap 업데이트

워크플로우 상태 확인:
`gh run list`

릴리스 확인:
`gh release view vX.Y.Z`

---

## 7단계: 완료 보고

성공적인 릴리스 후:

- 새 버전 번호: vX.Y.Z
- 릴리스 URL: https://github.com/modu-ai/moai-adk/releases/tag/vX.Y.Z
- 바이너리: darwin-arm64, darwin-amd64, linux-arm64, linux-amd64, windows-amd64
- 업데이트 방법: `moai update`

사용자에게 새 버전 사용 안내 및 변경 사항 요약 제공
