---
description: "Do v2.x 프로덕션 릴리즈 - 에이전트 위임 및 품질 검증 포함. 대상 브랜치는 항상 main. 태그 형식 vX.Y.Z가 GoReleaser를 트리거합니다. 모든 git 작업은 manager-git에 위임. 품질 게이트 실패 시 expert-debug로 에스컬레이션."
argument-hint: "[VERSION] - 선택적 대상 버전 (예: 2.1.0). 생략 시 patch/minor/major 선택 프롬프트."
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

## 릴리즈 설정

- **브랜치**: `main`
- **태그 형식**: `vX.Y.Z` (표준 semver, `.github/workflows/release.yml`을 통해 GoReleaser 트리거)
- **릴리즈 URL**: https://github.com/do-focus/do-focus/releases/tag/vX.Y.Z
- **바이너리**: darwin-arm64, darwin-amd64, linux-arm64, linux-amd64, windows-amd64

---

## 실행 지시 - 즉시 시작

이것은 릴리즈 명령입니다. 아래 워크플로우를 순서대로 실행하세요. 단계를 설명하지 말고 실제로 명령을 실행하세요.

제공된 인자: $ARGUMENTS

- VERSION 인자가 있으면: 대상 버전으로 사용, 버전 선택 건너뜀
- 인자 없으면: 버전 유형 선택 (patch/minor/major) 질문

---

## 사전 실행 컨텍스트

!git status --porcelain
!git branch --show-current
!git tag --list --sort=-v:refname | head -5
!git log --oneline -10

@go.mod
@pkg/version/version.go

---

## PHASE 0: 사전 점검

릴리즈 프로세스 시작 전 작업 디렉토리가 깨끗한지 확인:

1. **미커밋 변경사항 확인**:
   ```bash
   git status --porcelain
   ```

2. **미커밋 파일 처리**:
   - `.claude/`의 추적되지 않은 파일이 있으면: 커밋 또는 폐기 여부 확인
   - 정리 명령:
   ```bash
   git checkout -- .claude/
   git clean -fd .claude/ internal/template/templates/.claude/
   ```

3. **브랜치 확인**:
   - `main` 브랜치에 있어야 함
   - 아니면 main으로 전환: `git checkout main`

4. **최신 변경사항 가져오기**:
   ```bash
   git pull origin main
   ```

---

## PHASE 1: 품질 게이트

다음 항목을 TodoWrite로 생성 후 가능한 경우 병렬로 실행:

1. 전체 테스트 실행: `go test -race ./... -count=1 2>&1 | tail -30`
2. go vet 실행: `go vet ./... 2>&1 | tail -10`
3. go fmt 검사: `gofumpt -l . 2>/dev/null | head -10`

포맷 문제 발견 시 `make fmt`로 수정 후 커밋:
`git add -A && git commit -m "style: auto-fix formatting issues"`

품질 요약 표시:

- tests: PASS 또는 FAIL (FAIL이면 중단 및 보고)
- go vet: PASS 또는 WARNING
- gofmt: PASS 또는 FIXED

### 오류 처리

품질 게이트 실패 시:

- **expert-debug 서브에이전트 사용**하여 문제 진단 및 해결
- 모든 게이트 통과 후에만 릴리즈 워크플로우 재개

---

## PHASE 2: 코드 리뷰

마지막 태그 이후 커밋 확인:
`git log $(git describe --tags --abbrev=0 2>/dev/null || echo HEAD~20)..HEAD --oneline`

diff 통계 확인:
`git diff $(git describe --tags --abbrev=0 2>/dev/null || echo HEAD~20)..HEAD --stat`

다음 항목 분석:

- 버그 가능성
- 보안 이슈
- 호환성 깨는 변경
- 테스트 커버리지 갭

검토 결과 표시: PROCEED 또는 REVIEW_NEEDED

---

## PHASE 3: 버전 선택

VERSION 인자가 있으면 (예: "2.1.0"):

- 해당 버전 직접 사용
- AskUserQuestion 건너뜀

VERSION 인자 없으면:

- `pkg/version/version.go`에서 현재 버전 읽기
- AskUserQuestion으로 질문: patch/minor/major

새 버전 계산 후 모든 버전 파일 업데이트:

1. `pkg/version/version.go`의 `Version` 변수 수정
2. `.do/config/sections/system.yaml`의 `do.version` 수정
3. `internal/template/templates/.do/config/sections/system.yaml`의 `do.version` 수정
4. 커밋: `git add pkg/version/version.go .do/config/sections/system.yaml internal/template/templates/.do/config/sections/system.yaml && git commit -m "chore: bump version to vX.Y.Z"`

버전 파일 체크리스트:
- [ ] pkg/version/version.go: Version = "vX.Y.Z"
- [ ] .do/config/sections/system.yaml: do.version: "X.Y.Z"
- [ ] internal/template/templates/.do/config/sections/system.yaml: do.version: "X.Y.Z"

---

## PHASE 4: CHANGELOG 생성 (이중 언어: 영어 우선)

### [HARD] 영어 우선 이중 언어 형식

CHANGELOG.md와 GitHub 릴리즈 노트는 영어 우선 이중 언어 구조를 따라야 합니다. 국제 사용자가 영어를 먼저 볼 수 있도록 하면서 한국어 문서도 유지합니다.

changelog용 커밋 가져오기: `git log $(git describe --tags --abbrev=0)..HEAD --pretty=format:"- %s (%h)"`

### CHANGELOG.md 구조

다음 구조로 CHANGELOG.md에 새 버전 항목 앞에 추가:

```
## [X.Y.Z] - YYYY-MM-DD

### Summary
[영어: 핵심 기능과 개선 사항 2-3줄 요약]

### Breaking Changes
[영어: 호환성 깨는 변경 목록, 없으면 "None"]

### Added
- [영어 추가 항목 1]

### Changed
- [영어 변경 항목 1]

### Fixed
- [영어 수정 항목 1]

---

## [X.Y.Z] - YYYY-MM-DD (한국어)

### 요약
[한국어: 핵심 기능과 개선 사항 2-3줄 요약]

### 주요 변경 사항 (Breaking Changes)
[한국어: 호환성 깨는 변경 목록, 없으면 "없음"]

### 추가됨 (Added)
- [한국어 추가 항목 1]

### 변경됨 (Changed)
- [한국어 변경 항목 1]

### 수정됨 (Fixed)
- [한국어 수정 항목 1]

---

[이전 버전 항목]
```

CHANGELOG.md 커밋:
`git add CHANGELOG.md && git commit -m "docs: update CHANGELOG for vX.Y.Z"`

---

## PHASE 5: 최종 승인

릴리즈 요약 표시:

- 버전 변경 (현재 -> 대상)
- 포함된 커밋 (수 및 주요 항목)
- 품질 게이트 결과
- 승인 후 발생할 일

AskUserQuestion 사용:

- 릴리즈: 태그 생성 및 main에 푸시
- 중단: 취소 (변경사항은 로컬 유지)

---

## PHASE 6: 태그 및 푸시 (에이전트 위임 필수)

**[HARD] 모든 git 작업은 manager-git 에이전트에 위임해야 합니다.**

승인되면 다음 컨텍스트로 manager-git 서브에이전트에 위임:

```
## 임무: 버전 X.Y.Z 릴리즈 Git 작업

### 컨텍스트
- 대상 버전: X.Y.Z
- 대상 브랜치: main
- 태그 형식: vX.Y.Z (표준 semver)
- 현재 상태: [현재 git 상태 설명]
- 품질 게이트: 모두 통과
- 포함된 커밋: [커밋 수 및 요약]

### 필요한 작업
1. 원격 상태 확인: 원격(origin)에 vX.Y.Z 태그가 있는지 확인
2. 태그 충돌 처리:
   - 원격에 vX.Y.Z가 없으면: 태그 생성 및 푸시
   - 원격에 이미 있으면: 상황 보고 및 옵션 제시
3. 푸시 실행: `git push origin main --tags`
4. GoReleaser 워크플로우 트리거 확인
```

---

## PHASE 7: GitHub 릴리즈 노트 (이중 언어: 영어 우선)

### 1단계: GoReleaser 대기

GoReleaser가 초기 릴리즈와 바이너리 에셋을 생성합니다.

**예상 에셋:**
- do-focus_X.Y.Z_darwin_arm64.tar.gz
- do-focus_X.Y.Z_darwin_amd64.tar.gz
- do-focus_X.Y.Z_linux_arm64.tar.gz
- do-focus_X.Y.Z_linux_amd64.tar.gz
- do-focus_X.Y.Z_windows_amd64.zip
- checksums.txt

### 2단계: 영어 우선 이중 언어 콘텐츠로 릴리즈 노트 교체

**[HARD] 영어 섹션 먼저, 한국어 섹션 두 번째. CHANGELOG와 동일한 이중 언어 형식 사용.**

### 3단계: 최종 검증

검증 체크리스트:
- [ ] 영어 섹션이 먼저 나타남
- [ ] 구분자 `---`가 섹션 사이에 있음
- [ ] 한국어 섹션이 두 번째로 나타남
- [ ] 양쪽 섹션에 설치 명령 포함
- [ ] Breaking changes 섹션 포함 (없으면 "없음")

---

## PHASE 8: 로컬 환경 업데이트

릴리즈 검증 후 로컬 개발 환경을 새 바이너리로 업데이트하고 템플릿을 동기화합니다.

**1. 로컬 바이너리를 릴리즈된 버전으로 업데이트:**

```bash
do update --binary
```

**2. 로컬 프로젝트 템플릿 동기화 (필요시):**

```bash
do update --templates-only
```

**3. 로컬 환경 확인:**

```bash
do version
```

---

## 핵심 규칙

- **대상 브랜치**: `main` (프로덕션 릴리즈)
- **태그 형식**: `vX.Y.Z` (release.yml을 통해 GoReleaser 트리거)
- 테스트는 반드시 통과해야 계속 진행 (패키지당 85%+ 커버리지)
- 3개 버전 파일이 일관되어야 함
- **[HARD] CHANGELOG 및 GitHub 릴리즈: 영어 먼저, 한국어 두 번째**
- **[HARD] 모든 git 작업은 manager-git 에이전트에 위임해야 함**
- **[HARD] 품질 게이트 실패는 expert-debug 에이전트에 위임해야 함**

---

## 실행 시작

지금 Phase 1을 시작하세요. TodoWrite를 생성하고 즉시 품질 게이트를 실행하세요.
