---
name: moai-workflow-plan
description: >
  Plan-Run-Sync 워크플로우의 첫 번째 단계로 EARS 형식을 사용하여 포괄적인 SPEC 문서를
  생성합니다. 프로젝트 탐색, SPEC 파일 생성, 검증, 그리고 선택적 Git 환경 설정(worktree
  또는 브랜치 생성)을 처리합니다. 기능 계획이나 명세 작성 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-03"
  tags: "plan, spec, ears, requirements, specification, design"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["plan", "spec", "design", "architect", "requirements", "feature request"]
  agents: ["manager-spec", "Explore", "manager-git"]
  phases: ["plan"]
---

# Plan 워크플로우 오케스트레이션

## 목적

Plan-Run-Sync 워크플로우의 첫 번째 단계로 EARS 형식을 사용하여 포괄적인 SPEC 문서를 생성합니다. 이 워크플로우는 프로젝트 탐색부터 SPEC 파일 생성 및 선택적 Git 환경 설정까지 모든 것을 처리합니다.

## 범위

- MoAI 4단계 워크플로우의 1-2단계 구현 (의도 파악, 계획 생성)
- 3-4단계는 각각 /moai run과 /moai sync가 처리

## 입력

- $ARGUMENTS: 세 가지 패턴 중 하나
  - 기능 설명: "사용자 인증 시스템"
  - 재개 명령어: resume SPEC-XXX
  - 플래그와 함께하는 기능 설명: "사용자 인증" --worktree 또는 --branch

## 지원 플래그

- --worktree: 격리된 Git worktree 환경 생성 (최우선순위)
- --branch: 전통적인 기능 브랜치 생성 (두 번째 우선순위)
- 플래그 없음: 기본적으로 SPEC만 생성; 설정에 따라 사용자에게 확인할 수 있음
- --team: 팀 기반 탐색 활성화 (병렬 리서치 팀은 team-plan.md 참조)
- resume SPEC-XXX: 마지막 저장된 초안 상태에서 계속

플래그 우선순위: --worktree가 --branch보다 우선하며, --branch가 기본값보다 우선합니다.

## 컨텍스트 로딩

실행 전 다음 필수 파일을 로드하세요:

- .moai/config/config.yaml (git 전략, 언어 설정)
- .moai/config/sections/git-strategy.yaml (auto_branch, 브랜치 생성 정책)
- .moai/config/sections/language.yaml (git_commit_messages 설정)
- .moai/project/product.md (제품 컨텍스트)
- .moai/project/structure.md (아키텍처 컨텍스트)
- .moai/project/tech.md (기술 컨텍스트)
- .moai/specs/ 디렉토리 목록 (중복 제거를 위한 기존 SPEC)

실행 전 명령어: git status, git branch, git log, git diff, find .moai/specs.

---

## 단계 순서

### Phase 1A: 프로젝트 탐색 (선택 사항)

에이전트: Explore 서브에이전트 (읽기 전용 코드베이스 분석)

실행 시기:

- 사용자가 모호하거나 비구조적인 요청을 제공한 경우
- 기존 파일 및 패턴을 발견해야 할 때
- 현재 프로젝트 상태가 불명확한 경우

건너뛰는 경우:

- 사용자가 명확한 SPEC 제목을 제공한 경우 (예: "인증 모듈 추가")
- 기존 SPEC 컨텍스트가 있는 재개 시나리오

Explore 서브에이전트의 작업:

- 사용자 요청의 키워드로 관련 파일 검색
- .moai/specs/에서 기존 SPEC 문서 찾기
- 구현 패턴 및 의존성 식별
- 프로젝트 설정 파일 탐색
- Phase 1B 컨텍스트를 위한 포괄적인 결과 보고

### Phase 1B: SPEC 계획 (필수)

에이전트: manager-spec 서브에이전트

입력: 사용자 요청 + Phase 1A 결과 (실행된 경우)

manager-spec의 작업:

- 프로젝트 문서 분석 (product.md, structure.md, tech.md)
- 적절한 이름을 가진 1-3개의 SPEC 후보 제안
- .moai/specs/에서 중복 SPEC 확인
- 각 후보에 대한 EARS 구조 설계
- 기술적 제약 조건을 포함한 구현 계획 수립
- 라이브러리 버전 식별 (프로덕션 안정 버전만, 베타/알파 제외)

출력: SPEC 후보, EARS 구조, 기술적 제약 조건을 포함한 구현 계획.

### 결정 지점 1: SPEC 생성 승인

도구: AskUserQuestion (오케스트레이터 레벨에서만)

옵션:

- SPEC 생성 진행
- 계획 수정 요청
- 초안으로 저장
- 취소

"진행": Phase 1.5 후 Phase 2 계속.
"수정": 피드백 수집 후 피드백 컨텍스트와 함께 Phase 1B 재실행.
"초안": plan.md를 초안 상태로 저장, 커밋 생성, 재개 명령어 출력, 종료.
"취소": 계획 삭제, 파일 생성 없이 종료.

### Phase 1.5: 생성 전 검증 게이트

목적: 파일 생성 전 흔한 SPEC 생성 오류 방지.

Step 1 - 문서 유형 분류:

- SPEC, 리포트, 문서로 분류하기 위한 키워드 감지
- 리포트는 .moai/reports/로, 문서는 .moai/docs/로 라우팅
- SPEC 유형 컨텐츠만 Phase 2로 진행

Step 2 - SPEC ID 검증 (모든 확인 통과 필요):

- ID 형식: SPEC-{DOMAIN}-{NUMBER} 패턴과 일치해야 함 (예: SPEC-AUTH-001)
- 도메인 이름: 승인된 도메인 목록에서 선택 (AUTH, API, UI, DB, REFACTOR, FIX, UPDATE, PERF, TEST, DOCS, INFRA, DEVOPS, SECURITY 등)
- ID 고유성: .moai/specs/를 검색하여 중복 없음 확인
- 디렉토리 구조: 디렉토리를 생성해야 함, 플랫 파일 절대 금지

복합 도메인 규칙: 최대 2개 도메인 권장 (예: UPDATE-REFACTOR-001), 최대 3개 허용.

### Phase 2: SPEC 문서 생성

에이전트: manager-spec 서브에이전트

입력: Phase 1B에서 승인된 계획, Phase 1.5에서 검증된 SPEC ID.

파일 생성 (세 파일 동시 생성):

- .moai/specs/SPEC-{ID}/spec.md
  - 7개 필수 필드를 포함한 YAML 프론트매터 (id, version, status, created, updated, author, priority)
  - 프론트매터 바로 뒤에 HISTORY 섹션
  - 5가지 요구사항 유형을 모두 포함한 완전한 EARS 구조
  - conversation_language로 작성된 내용

- .moai/specs/SPEC-{ID}/plan.md
  - 작업 분해를 포함한 구현 계획
  - 기술 스택 명세 및 의존성
  - 리스크 분석 및 완화 전략

- .moai/specs/SPEC-{ID}/acceptance.md
  - 최소 2개의 Given/When/Then 테스트 시나리오
  - 엣지 케이스 테스트 시나리오
  - 성능 및 품질 게이트 기준

품질 제약 조건:

- SPEC당 요구사항 모듈 5개 이하
- 인수 기준 최소 2개의 Given/When/Then 시나리오
- 기술 용어 및 함수 이름은 영어 유지

### Phase 3: Git 환경 설정 (조건부)

실행 조건: Phase 2 성공적으로 완료 AND 다음 중 하나:

- --worktree 플래그 제공
- --branch 플래그 제공 또는 사용자가 브랜치 생성 선택
- 설정이 브랜치 생성 허용 (git_strategy 설정)

건너뛰는 경우: develop_direct 워크플로우, 플래그 없고 사용자가 "현재 브랜치 사용" 선택.

#### Worktree 경로 (--worktree 플래그)

전제 조건: worktree 생성 전 SPEC 파일을 반드시 커밋해야 합니다.

- SPEC 파일 스테이징: git add .moai/specs/SPEC-{ID}/
- 커밋 생성: feat(spec): Add SPEC-{ID} - {title}
- feature/SPEC-{ID} 브랜치로 WorktreeManager를 통해 worktree 생성
- worktree 경로 및 탐색 방법 표시

#### 브랜치 경로 (--branch 플래그 또는 사용자 선택)

에이전트: manager-git 서브에이전트

- 브랜치 생성: feature/SPEC-{ID}-{description}
- 원격 저장소 존재 시 업스트림 추적 설정
- 새 브랜치로 전환
- 팀 모드: manager-git 서브에이전트를 통해 초안 PR 생성

#### 현재 브랜치 경로 (플래그 없음 또는 사용자 선택)

- 브랜치 생성 없음, manager-git 호출 없음
- SPEC 파일은 현재 브랜치에 유지

### 결정 지점 2: 개발 환경 선택

도구: AskUserQuestion (prompt_always 설정이 true이고 auto_branch가 true일 때)

옵션:

- Worktree 생성 (병렬 SPEC 개발에 권장)
- 브랜치 생성 (전통적인 워크플로우)
- 현재 브랜치 사용

### 결정 지점 3: 다음 액션 선택

도구: AskUserQuestion (SPEC 생성 완료 후)

옵션:

- 구현 시작 (/moai run SPEC-{ID} 실행)
- 계획 수정
- 새 기능 추가 (추가 SPEC 생성)

---

## 팀 모드 라우팅

--team 플래그가 제공되거나 자동 선택된 경우, plan 단계는 반드시 팀 오케스트레이션으로 전환해야 합니다:

1. 전제 조건 확인: workflow.team.enabled == true AND CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 환경변수 설정
2. 전제 조건 충족 시: workflows/team-plan.md를 읽고 팀 워크플로우 실행 (researcher + analyst + architect와 함께 TeamCreate)
3. 전제 조건 미충족 시: "팀 모드는 workflow.yaml에서 workflow.team.enabled: true 설정과 CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 환경변수가 필요합니다"라고 사용자에게 경고 후 표준 서브에이전트 모드 (manager-spec)로 폴백

팀 구성: researcher (haiku) + analyst (inherit) + architect (inherit)

상세한 팀 오케스트레이션 단계는 workflows/team-plan.md를 참조하세요.

---

## 완료 기준

다음 사항을 모두 확인해야 합니다:

- Phase 1: manager-spec이 프로젝트를 분석하고 SPEC 후보 제안
- SPEC 생성 전 AskUserQuestion을 통한 사용자 승인 획득
- Phase 2: 3개의 SPEC 파일 모두 생성 (spec.md, plan.md, acceptance.md)
- 디렉토리 이름이 .moai/specs/SPEC-{ID}/ 형식을 따름
- YAML 프론트매터에 7개 필수 필드 모두 포함
- EARS 구조 완성
- Phase 3: 플래그와 사용자 선택에 따른 적절한 git 작업 수행
- --worktree: worktree 생성 전 SPEC 커밋
- 사용자에게 다음 단계 제시

---

Version: 2.0.0
Updated: 2026-02-07
Source: .claude/commands/moai/1-plan.md v5.1.0에서 추출. 팀 모드 지원 및 --team 플래그 추가.
