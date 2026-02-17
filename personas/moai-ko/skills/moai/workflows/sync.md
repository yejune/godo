---
name: moai-workflow-sync
description: >
  코드 변경 사항과 문서화를 동기화하고, 프로젝트 품질을 검증하며, 풀 리퀘스트를
  완성합니다. Plan-Run-Sync 워크플로우의 세 번째 단계입니다. SPEC 편차 분석 및
  프로젝트 문서 업데이트를 포함합니다. 문서 동기화, PR 생성, 품질 검증 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.1.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-03"
  tags: "sync, documentation, pull-request, quality, verification, pr"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["sync", "docs", "pr", "documentation", "pull request", "changelog", "readme"]
  agents: ["manager-docs", "manager-quality", "manager-git"]
  phases: ["sync"]
---

# Sync 워크플로우 오케스트레이션

## 목적

코드 변경 사항과 문서화를 동기화하고, 프로젝트 품질을 검증하며, 풀 리퀘스트를 완성합니다. 이것은 Plan-Run-Sync 워크플로우의 세 번째 단계입니다.

## 범위

- MoAI 4단계 워크플로우의 4단계 구현 (보고 및 커밋)
- /moai run의 구현 산출물을 입력으로 받음
- 동기화된 문서화, 커밋, PR 준비 상태를 산출

## 입력

- $ARGUMENTS: 모드 및 선택적 경로
  - 모드: auto (기본값), force, status, project
  - 경로: 선택적 동기화 대상 경로 (예: src/auth/)
  - 플래그: --merge

## 지원 모드

- auto (기본값): 변경된 파일만 스마트 선택적 동기화. PR Ready 전환. 일상 개발 워크플로우.
- force: 모든 문서화 완전 재생성. 오류 복구 및 대규모 리팩토링 사용 사례.
- status: 읽기 전용 상태 확인. 변경 없는 빠른 프로젝트 상태 보고.
- project: 프로젝트 전체 문서화 업데이트. 마일스톤 완료 및 주기적 동기화 사용 사례.

## 지원 플래그

- --merge: 동기화 후 PR 자동 병합 및 브랜치 정리. Worktree/브랜치 환경은 git 컨텍스트에서 자동 감지.

## 컨텍스트 로딩

실행 전 다음 필수 파일을 로드하세요:

- .moai/config/config.yaml (git 전략, 언어 설정)
- .moai/config/sections/git-strategy.yaml (auto_branch, 브랜치 생성 정책)
- .moai/config/sections/language.yaml (git_commit_messages 설정)
- .moai/specs/ 디렉토리 목록 (동기화할 SPEC 문서)
- .moai/project/ 디렉토리 목록 (조건부 업데이트할 프로젝트 문서)
- README.md (현재 프로젝트 문서)

실행 전 명령어: git status, git diff, git branch, git log, find .moai/specs.

---

## 단계 순서

### Phase 0.5: 품질 검증 (병렬 진단)

목적: 동기화 시작 전 프로젝트 품질 검증. 이슈를 조기에 발견하기 위해 Phase 1 전에 실행됩니다.

#### Step 1: 프로젝트 언어 감지

우선순위 순서로 지시자 파일 확인 (첫 번째 일치 적용):

- Python: pyproject.toml, setup.py, requirements.txt, .python-version, Pipfile
- TypeScript: tsconfig.json, typescript 의존성이 있는 package.json
- JavaScript: tsconfig 없는 package.json
- Go: go.mod, go.sum
- Rust: Cargo.toml, Cargo.lock
- Ruby: Gemfile, .ruby-version, Rakefile
- Java: pom.xml, build.gradle, build.gradle.kts
- PHP: composer.json, composer.lock
- Kotlin: kotlin 플러그인이 있는 build.gradle.kts
- Swift: Package.swift, .xcodeproj, .xcworkspace
- C#/.NET: .csproj, .sln, .fsproj
- C++: CMakeLists.txt, C++ 내용이 있는 Makefile
- Elixir: mix.exs
- R: DESCRIPTION (R 패키지), .Rproj, renv.lock
- Flutter/Dart: pubspec.yaml
- Scala: build.sbt, build.sc
- 폴백: unknown (언어별 도구 건너뜀, 코드 검토로 진행)

#### Step 2: 병렬 진단 실행

세 가지 백그라운드 작업을 동시에 실행:

- 테스트 실행기: 언어별 테스트 명령어 (pytest, npm test, go test, cargo test 등)
- 린터: 언어별 린트 명령어 (ruff, eslint, golangci-lint, clippy 등)
- 타입 검사기: 언어별 타입 검사 (mypy, tsc --noEmit, go vet 등)

타임아웃을 두고 모든 결과 수집 (테스트 180초, 기타 120초). 부분 실패를 우아하게 처리.

#### Step 3: 테스트 실패 처리

테스트가 실패하면 AskUserQuestion 사용:

- 계속: 실패에도 불구하고 동기화 진행
- 중단: 동기화 중지, 먼저 테스트 수정 (Phase 4 우아한 종료로 이동)

#### Step 4: 코드 검토

에이전트: manager-quality 서브에이전트

프로젝트 언어와 무관하게 호출됩니다. TRUST 5 품질 검증 실행 및 포괄적인 품질 보고서 생성.

#### LSP 품질 게이트

동기화 단계는 quality.yaml에 설정된 LSP 기반 품질 게이트를 적용합니다:
- 오류 없음 필수 (lsp_quality_gates.sync.max_errors: 0)
- 최대 10개 경고 허용 (lsp_quality_gates.sync.max_warnings: 10)
- 클린 LSP 상태 필수 (lsp_quality_gates.sync.require_clean_lsp: true)

#### Step 5: 품질 보고서 생성

테스트 실행기, 린터, 타입 검사기, 코드 검토의 상태를 보여주는 품질 보고서로 모든 결과 집계. 전반적인 상태 결정 (PASS 또는 WARN).

status 모드 조기 종료: 모드가 "status"인 경우 품질 보고서를 표시하고 종료합니다. 이후 단계는 실행되지 않습니다.

### Phase 1: 분석 및 계획

#### Step 1.1: 전제 조건 확인

- .moai/ 디렉토리가 존재해야 함
- .claude/ 디렉토리가 존재해야 함
- 프로젝트가 Git 저장소 내에 있어야 함

#### Step 1.2: 프로젝트 상태 분석

- Git 변경 사항 분석: git status, git diff, 변경된 파일 분류
- 프로젝트 설정 읽기: git_strategy.mode, conversation_language, spec_git_workflow
- $ARGUMENTS에서 동기화 모드 결정
- Worktree 컨텍스트 감지: git 디렉토리에 worktrees/ 컴포넌트가 있는지 확인
- 브랜치 컨텍스트 감지: 현재 브랜치 이름 확인

#### Step 1.3: 프로젝트 상태 검증

(변경된 파일만이 아닌) 모든 소스 파일을 스캔:

- 끊어진 참조 및 불일치
- 정확한 위치를 포함한 이슈
- 심각도 분류 (Critical, High, Medium, Low)

#### Step 1.4: 동기화 계획

에이전트: manager-docs 서브에이전트

Git 변경 사항, 모드, 프로젝트 검증 결과를 기반으로 동기화 전략 수립. 출력: 업데이트할 문서, 동기화가 필요한 SPEC, 필요한 프로젝트 개선, 예상 범위.

#### Step 1.5: SPEC-구현 편차 분석

목적: 문서화 정확성을 위해 원래 SPEC 계획과 실제 구현 간의 차이를 감지합니다.

현재 동기화와 관련된 각 SPEC에 대해:

- Step 1.5.1: SPEC 문서 로드
  - spec.md (요구사항), plan.md (구현 계획), acceptance.md (기준) 읽기
  - 계획된 파일, 계획된 기능, 계획된 범위 추출

- Step 1.5.2: 실제 구현 분석
  - git diff와 git log를 사용하여 실행 단계 중 생성, 수정 또는 삭제된 모든 파일 식별
  - 도메인별 변경 사항 분류 (backend, frontend, tests, config, docs)

- Step 1.5.3: 계획 vs 현실 비교
  - 원래 plan.md에 없었던 생성된 파일 식별
  - 원래 spec.md 범위를 넘어 구현된 기능 또는 엔드포인트 식별
  - 구현되지 않은 계획된 항목 식별 (연기 또는 삭제)
  - 계획에 없던 리팩토링 또는 의존성 변경 식별

- Step 1.5.4: 편차 보고서 생성
  - 편차 분류: scope_expansion, unplanned_additions, deferred_items, structural_changes
  - 포함 내용: new_directories_created, new_dependencies_added, new_features_implemented
  - 이 보고서는 Phase 2.2 (SPEC 업데이트) 및 Phase 2.2.5 (프로젝트 문서 업데이트)에 사용됩니다

- Step 1.5.5: SPEC 라이프사이클 레벨 확인
  - 라이프사이클 레벨에 대한 SPEC 메타데이터 읽기 (지정되지 않은 경우 기본값: spec-first)
  - 레벨 1 (spec-first): SPEC에 구현 요약을 추가하여 완료 표시
  - 레벨 2 (spec-anchored): 실제 구현을 반영하도록 SPEC 내용 업데이트
  - 레벨 3 (spec-as-source): 불일치를 경고로 표시 (구현이 SPEC과 정확히 일치해야 함)

#### Step 1.6: 사용자 승인

도구: AskUserQuestion

동기화 계획 보고서를 표시하고 옵션 제시:

- 동기화 진행
- 수정 요청 (Phase 1 재실행)
- 세부 사항 검토 (전체 프로젝트 결과 표시, 재질문)
- 중단 (변경 없이 종료)

### Phase 2: 문서 동기화 실행

#### Step 2.1: 안전 백업 생성

수정 전:

- 타임스탬프 식별자 생성
- 백업 디렉토리 생성: .moai/backups/sync-{timestamp}/
- 중요 파일 복사: README.md, docs/, .moai/specs/
- 백업 무결성 확인 (비어 있지 않은 디렉토리 확인)

#### Step 2.2: 문서 동기화

에이전트: manager-docs 서브에이전트

입력: 승인된 동기화 계획, 프로젝트 검증 결과, 변경된 파일 목록, Phase 1.5의 편차 보고서.

manager-docs의 작업:

- 변경된 코드를 Living Documents에 반영
- API 문서 자동 생성 및 업데이트
- 필요 시 README 업데이트
- 아키텍처 문서 동기화
- 프로젝트 이슈 수정 및 끊어진 참조 복원
- 편차 분석 및 라이프사이클 레벨에 따라 SPEC 문서 업데이트 (Step 2.2.1 참조)
- 변경된 도메인 감지 및 도메인별 업데이트 생성
- 동기화 보고서 생성: .moai/reports/sync-report-{timestamp}.md

모든 문서 업데이트는 conversation_language 설정을 사용합니다.

##### Step 2.2.1: SPEC 문서 업데이트 (편차 보고서 기반)

Phase 1.5.5에서 감지된 SPEC 라이프사이클 레벨에 따라 업데이트 적용:

레벨 1 (spec-first):
- spec.md에 실제 구현을 요약하는 "구현 노트" 섹션 추가
- 범위 변경 기록: 계획을 넘어 추가된 기능, 연기된 항목
- SPEC을 완료로 표시 (이후 유지보수 불필요)

레벨 2 (spec-anchored):
- 실제 구현을 반영하도록 spec.md 요구사항 업데이트
- 원래 범위를 넘어 구현된 기능에 대해 새 EARS 형식 요구사항 추가
- 실제 수행된 구현 단계로 plan.md 업데이트
- 추가된 기능에 대한 새 인수 기준으로 acceptance.md 업데이트
- 변경된 경우 "as-implemented" 주석과 함께 원래 요구사항 보존

레벨 3 (spec-as-source):
- SPEC 내용 수정 금지
- 구현의 SPEC 편차를 나열하는 불일치 보고서 생성
- 수동 검토를 위해 동기화 보고서에 경고로 표시
- SPEC 업데이트 또는 구현 조정 권장

#### Step 2.2.5: 프로젝트 문서 업데이트 (조건부)

목적: 중요한 구조적 변경이 감지될 때 .moai/project/ 문서 업데이트.

조건: Phase 1.5의 편차 보고서가 다음을 나타낼 때만 이 단계 실행:
- 프로젝트에 새 디렉토리 생성
- 새 의존성 또는 기술 추가
- 새로운 주요 기능 또는 기능 구현
- 중요한 아키텍처 변경 발생

건너뛰는 조건: .moai/project/ 디렉토리가 없거나 파일이 없으면 이 단계를 완전히 건너뜁니다.

에이전트: manager-docs 서브에이전트

manager-docs의 작업:

- 새 디렉토리 생성 시: 새 디렉토리 설명과 목적으로 structure.md 업데이트
- 새 의존성 추가 시: 새 기술 스택 항목과 근거로 tech.md 업데이트
- 새 기능 구현 시: 새 기능 설명과 사용 사례로 product.md 업데이트
- 아키텍처 변경 시: 수정된 아키텍처 패턴으로 structure.md 업데이트

제약 조건:
- 감지된 변경 사항과 관련된 섹션만 업데이트 (전체 파일 재생성 금지)
- 기존 내용 보존, 점진적으로 추가 또는 수정
- 모든 업데이트에 conversation_language 설정 사용

#### Step 2.3: 동기화 후 품질 검증

에이전트: manager-quality 서브에이전트

TRUST 5 기준으로 동기화 품질 검증:

- 모든 프로젝트 링크 완전
- 문서 잘 포맷됨
- 모든 문서 일관성
- 크리덴셜 노출 없음
- 모든 SPEC 적절하게 연결됨

#### Step 2.4: SPEC 상태 업데이트

라이프사이클 레벨 및 구현 완성도에 따라 SPEC 상태 업데이트:

- 레벨 1 (spec-first): 상태를 "completed"로 설정. 이후 유지보수 불필요.
- 레벨 2 (spec-anchored): 모든 요구사항 충족 시 "completed"로 설정, 부분적인 경우 "in-progress". 분기별 유지보수 정책에 따라 다음 검토 예약.
- 레벨 3 (spec-as-source): 구현-SPEC 일치도에 따라 상태 설정. 해결을 위해 불일치 표시.

버전 변경, 상태 전환, 편차 요약 기록. 동기화 보고서에 포함.

### Phase 3: Git 작업 및 전달

#### Step 3.0: Git 워크플로우 전략 감지

`.moai/config/sections/system.yaml`에서 `github.git_workflow`를 읽습니다. 이것이 변경 사항이 전달되는 방식을 결정합니다.

| 전략 | 브랜치 모델 | PR 동작 | 적합한 경우 |
|----------|-------------|-------------|----------|
| github_flow | main에서 기능 브랜치 | main으로 PR 자동 생성 | 팀/오픈소스 프로젝트 |
| main_direct | main에 직접 커밋 | PR 생성 안 함 | 개인 개발 |
| gitflow | develop/release/hotfix 브랜치 | 적절한 베이스로 PR | 엔터프라이즈 프로젝트 |

기본 전략 (설정 없을 시): `github_flow`

SPEC 브랜치 처리를 결정하기 위해 `github.spec_git_workflow`도 읽습니다:
- `feature_branch`: 각 SPEC이 자체 브랜치를 가짐 (github_flow/gitflow에 권장)
- `main_direct`: SPEC 변경 사항을 현재 브랜치에 커밋 (git_workflow가 main_direct일 때만)

#### Step 3.1: 변경 사항 커밋

에이전트: manager-git 서브에이전트

- 변경된 모든 문서 파일, 보고서, README, docs/ 스테이징
- 동기화된 문서, 프로젝트 수정, SPEC 업데이트를 나열하는 설명적인 메시지로 단일 커밋 생성
- 커밋 메시지 언어는 `language.git_commit_messages` 설정을 따름
- git log로 커밋 확인

#### Step 3.2: 푸시 및 전달 (전략 인식)

동작은 `github.git_workflow` 설정과 현재 브랜치 컨텍스트에 따라 다릅니다.

##### 전략: github_flow

현재 브랜치 감지:

**기능 브랜치** (main 외의 모든 브랜치):
1. 원격으로 브랜치 푸시: `git push -u origin <branch>`
2. PR이 이미 존재하는지 확인: `gh pr list --head <branch> --json number`
3. PR이 없으면: `gh pr create`로 PR 생성
   - 제목: SPEC 제목 또는 브랜치 이름에서 도출
   - 본문: 동기화 요약, 변경된 파일, 품질 보고서 포함
   - 베이스: main
   - 레이블: 변경된 파일에서 자동 감지
4. PR이 있으면: 동기화 변경 사항을 요약하는 댓글로 업데이트
5. 사용자에게 PR URL 표시

**main 브랜치** (직접 커밋):
- 직접 푸시: `git push origin main`
- 푸시 확인 표시
- 참고: main 직접 커밋은 허용되지만 기능 브랜치를 권장함

**Worktree 컨텍스트** (git 디렉토리 구조에서 감지):
- Worktree 브랜치를 원격으로 푸시
- PR이 없으면 생성 (기능 브랜치 흐름과 동일)
- PR URL 및 worktree 컨텍스트 표시

##### 전략: main_direct

모든 커밋은 직접 main으로, PR 없음:
1. main으로 푸시: `git push origin main`
2. 푸시 확인 표시
3. 브랜치 이름과 무관하게 PR 생성 안 함

##### 전략: gitflow

현재 브랜치 유형을 감지하여 라우팅:

**feature/* 브랜치** → `develop`으로 PR:
1. 브랜치 푸시: `git push -u origin <branch>`
2. `develop` 브랜치를 대상으로 PR 생성 또는 업데이트
3. PR URL 표시

**release/* 브랜치** → `main`으로 PR:
1. 브랜치 푸시: `git push -u origin <branch>`
2. `main` 브랜치를 대상으로 PR 생성 또는 업데이트
3. PR URL 표시

**hotfix/* 브랜치** → `main`으로 PR (develop으로 역병합):
1. 브랜치 푸시: `git push -u origin <branch>`
2. `main` 브랜치를 대상으로 PR 생성 또는 업데이트
3. 병합 후: 역병합을 위해 `develop`으로 후속 PR 생성
4. PR URL 표시

**develop 브랜치** → 직접 푸시:
1. develop으로 푸시: `git push origin develop`
2. 푸시 확인 표시

**main 브랜치** → 오류:
- gitflow에서 main에 직접 커밋 불허
- 대신 hotfix 또는 release 브랜치 생성 제안

#### Step 3.3: PR Ready 전환 (팀 모드)

Step 3.2에서 PR이 생성된 경우에만 적용:

- 팀 모드 활성화 및 PR이 초안인 경우: `gh pr ready`로 준비 상태로 전환
- 설정된 경우 검토자 및 레이블 할당
- 팀 모드 비활성화 시: 자동 전환 금지 (사용자가 준비 상태 제어)

#### Step 3.4: 자동 병합 (--merge 플래그 설정 시)

Step 3.2에서 PR이 생성된 경우에만 적용.

실행 조건 [HARD]:
- 플래그가 명시적으로 설정되어야 함: --merge
- 모든 CI/CD 검사 통과 필요
- PR에 병합 충돌이 없어야 함
- 최소 검토자 승인 획득 (팀 모드인 경우)

자동 병합 실행:
1. `gh pr checks --watch`로 CI/CD 상태 확인 (완료 대기)
2. `gh pr view --json mergeable`로 병합 충돌 확인
3. 통과하고 병합 가능한 경우: `gh pr merge --squash --delete-branch` 실행
4. 대상 브랜치 체크아웃, 최신 내용 가져오기
5. 로컬이 원격과 동기화됐는지 확인

자동 병합 실패:
- CI/CD 실패 시: 실패 보고, 오류 세부 사항 표시, 병합 금지
- 병합 충돌 시: 충돌 보고, 수동 해결 안내 제공, 병합 금지
- 승인 미달 (팀 모드) 시: 대기 중인 승인 보고, 병합 금지

### Phase 4: 완료 및 다음 단계

#### 완료 보고서

다음을 포함한 요약 표시:
- 사용된 Git 워크플로우 전략 (github_flow, main_direct, 또는 gitflow)
- 동기화 모드 및 범위
- 업데이트 및 생성된 파일
- 프로젝트 개선 사항
- 업데이트된 문서
- 생성된 보고서
- 백업 위치
- PR URL (생성된 경우) 또는 푸시 대상 (직접 푸시인 경우)

#### 컨텍스트 인식 다음 단계

도구: 전달 결과에 맞춰진 옵션과 함께 AskUserQuestion:

**PR이 생성된 경우 (github_flow 기능 브랜치, 또는 gitflow):**
- GitHub에서 PR 검토
- PR 자동 병합 (/moai sync --merge)
- 다음 SPEC 생성 (/moai plan)
- 새 세션 시작 (/clear)

**직접 푸시인 경우 (main_direct, 또는 github_flow main 브랜치):**
- 다음 SPEC 생성 (/moai plan)
- 개발 계속
- 새 세션 시작 (/clear)

**Worktree 컨텍스트인 경우:**
- 브라우저에서 PR 검토
- 메인 디렉토리로 돌아가기
- 이 Worktree 제거

---

## 팀 모드

동기화 단계는 다른 단계에서 --team이 활성화된 경우에도 항상 서브에이전트 모드 (manager-docs)를 사용합니다. 문서화 동기화는 순차적 일관성과 프로젝트 상태에 대한 단일 권위 있는 뷰가 필요합니다.

근거 및 세부 사항은 workflows/team-sync.md를 참조하세요.

---

## 우아한 종료

사용자가 어떤 결정 지점에서 중단하는 경우:

- 문서, Git 히스토리, 브랜치 상태에 변경 사항 없음
- 프로젝트는 현재 상태 유지
- 재시도 명령어 표시: /moai sync [mode]
- 코드 0으로 종료

---

## 완료 기준

다음 사항을 모두 확인해야 합니다:

- Phase 0.5: 품질 검증 완료 (테스트, 린터, 타입 검사기, 코드 검토)
- Phase 1: 전제 조건 확인, 프로젝트 분석, 편차 분석 완료, 사용자가 동기화 계획 승인
- Phase 2: 안전 백업 생성 및 확인, 문서 동기화, 라이프사이클 레벨에 따른 SPEC 문서 업데이트, 프로젝트 문서 업데이트 (해당 시), 품질 검증, SPEC 상태 업데이트
- Phase 3: 변경 사항 커밋, git_workflow 전략에 따른 전달 (github_flow/gitflow는 PR 생성, main_direct는 직접 푸시), 자동 병합 실행 (플래그 설정 및 PR 존재 시)
- Phase 4: 전달 결과를 포함한 완료 보고서 표시, 전략 및 컨텍스트에 따른 적절한 다음 단계 제시

---

Version: 3.0.0
Updated: 2026-02-07
Source: .claude/commands/moai/3-sync.md v3.4.0에서 추출. SPEC 편차 분석, 프로젝트 문서 업데이트, SPEC 라이프사이클 인식, 팀 모드 섹션, LSP 품질 게이트, github_flow/main_direct/gitflow를 지원하는 전략 인식 git 전달 추가.
