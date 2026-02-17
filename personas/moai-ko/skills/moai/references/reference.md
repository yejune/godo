---
name: moai-reference
description: >
  모든 MoAI 워크플로우에서 사용되는 일반적인 실행 패턴, 플래그 참조,
  레거시 명령어 매핑, 설정 파일 경로, 에러 처리 위임.
  재개 패턴 및 컨텍스트 전파 안내를 제공합니다.
  실행 패턴, 플래그 세부 정보, 또는 설정 참조가 필요할 때 사용합니다.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.1.0"
  category: "foundation"
  status: "active"
  updated: "2026-02-03"
  tags: "reference, patterns, flags, configuration, legacy, resume, context"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["reference", "pattern", "flag", "config", "resume", "legacy", "mapping"]
  agents: ["manager-spec", "manager-ddd", "manager-docs", "manager-quality", "manager-git"]
  phases: ["plan", "run", "sync"]
---

# MoAI 스킬 참조

모든 MoAI 워크플로우에서 사용되는 일반적인 패턴, 플래그 참조, 레거시 명령어 매핑, 설정 파일.

---

## 실행 패턴

### 병렬 실행 패턴

여러 작업이 독립적인 경우, 단일 응답에서 호출합니다. Claude Code는 여러 Task() 호출을 자동으로 병렬로 실행합니다 (최대 10개 동시).

사용 사례:

- 탐색 단계: 코드베이스 분석, 문서 조사, 품질 평가를 별도의 Task() 호출로 동시에 실행
- 진단 스캔: LSP 진단, AST-grep 분석, 린터 검사를 병렬로 실행
- 다중 파일 생성: 분석이 완료되면 product.md, structure.md, tech.md를 동시에 생성

구현:

- 동일한 응답 메시지에 여러 Task() 호출 포함
- 각 Task()는 다른 서브에이전트 또는 동일 에이전트의 다른 범위를 대상으로 함
- 모든 병렬 작업이 완료되면 결과 수집
- 최적 처리량을 위해 최대 10개 동시 Task() 호출

### 순차 실행 패턴

작업에 의존성이 있는 경우, 순차적으로 연결합니다. 각 Task() 호출은 이전 단계 결과의 컨텍스트를 받습니다.

사용 사례:

- DDD 워크플로우: 1단계(계획)가 2단계(구현)로, 2단계가 2.5단계(품질 검증)로 이어짐
- SPEC 생성: Explore 에이전트 결과가 manager-spec 에이전트의 문서 생성으로 이어짐
- 릴리즈 파이프라인: 품질 게이트 통과 후 버전 선택, 버전 선택 완료 후 태깅

구현:

- 다음 Task()를 호출하기 전에 각 Task()의 반환 대기
- 다음 Task() 프롬프트에 이전 단계 출력을 컨텍스트로 포함
- 의미론적 연속성 보장: 각 에이전트가 독립적으로 작동하기에 충분한 컨텍스트 수신

### 하이브리드 실행 패턴

단일 워크플로우 내에서 병렬과 순차 패턴을 결합합니다.

사용 사례:

- Fix 워크플로우: 병렬 진단 스캔 (LSP + 린터 + AST-grep), 그 후 결합된 결과를 기반으로 순차적 수정 적용
- MoAI 워크플로우: 병렬 탐색 단계, 그 후 순차적 SPEC 생성과 DDD 구현
- Run 워크플로우: 병렬 품질 검사, 그 후 순차적 구현 작업

구현:

- 어떤 작업이 독립적인지 파악 (병렬화)
- 어떤 작업이 이전 결과에 의존하는지 파악 (순차화)
- 각 단계 시작 시 병렬 작업을 그룹화하고, 이어서 순차적 의존 작업 실행

---

## 재개 패턴

워크플로우가 중단되거나 이전 세션에서 계속해야 할 때 --resume 플래그를 사용합니다.

동작:

- .moai/specs/SPEC-XXX/에서 기존 SPEC 문서 읽기
- SPEC 상태 마커에서 마지막으로 완료된 단계 결정
- 완료된 단계를 건너뛰고 다음 보류 단계에서 재개
- 이전의 모든 분석, 결정, 생성된 산출물 보존

적용 가능한 워크플로우:

- plan --resume SPEC-XXX: 마지막 체크포인트에서 SPEC 생성 재개
- run --resume SPEC-XXX: 마지막으로 완료된 작업에서 DDD 구현 재개
- moai --resume SPEC-XXX: 마지막 단계에서 전체 자율 워크플로우 재개
- fix --resume: 마지막 진단 상태에서 수정 사이클 재개

---

## 단계 간 컨텍스트 전파

중복 분석을 피하기 위해 각 단계는 결과를 다음 단계로 전달해야 합니다.

필수 컨텍스트 요소:

- 탐색 결과: 파일 경로, 아키텍처 패턴, 기술 스택, 의존성 맵
- SPEC 데이터: 요구사항 목록, 인수 기준, 기술적 접근 방식, 범위 경계
- 구현 결과: 수정된 파일, 생성된 테스트, 커버리지 지표, 남은 작업
- 품질 결과: 테스트 통과/실패 수, 린트 오류, 타입 검사 결과, 보안 발견사항
- Git 상태: 현재 브랜치, 마지막 태그 이후 커밋 수, 태그 히스토리

전파 방법:

- Task() 프롬프트에 이전 단계 출력의 구조화된 요약 포함
- 대용량 콘텐츠 블록을 인라인으로 포함하는 대신 특정 파일 경로 참조
- 단계 전반에 걸쳐 SPEC 문서를 진실의 단일 출처로 사용

---

## 플래그 참조

### 전역 플래그 (모든 워크플로우에서 사용 가능)

- --resume [ID]: 마지막 체크포인트에서 워크플로우 재개 (SPEC-ID 또는 스냅샷 ID)
- --seq: 해당하는 경우 병렬 대신 순차 실행 강제
- --ultrathink: 실행 전 심층 분석을 위해 Sequential Thinking MCP 활성화

### Plan 플래그

- --worktree: SPEC 구현을 위한 격리된 git worktree 생성
- --branch: SPEC을 위한 기능 브랜치 생성 (기본 브랜치 명명: spec/SPEC-XXX)
- --resume SPEC-XXX: 중단된 plan 세션 재개

### Run 플래그

- --resume SPEC-XXX: 마지막으로 완료된 작업에서 DDD 구현 재개

### Sync 플래그

- 모드 (위치 인수): auto (기본값), force, status, project
- --merge: sync 후 PR 자동 병합 및 브랜치 정리

### Fix 플래그

- --dry: 수정 적용 없이 감지된 이슈 미리보기
- --level N: 수정 깊이 제어 (레벨 1: 자동 수정 가능, 레벨 2: 간단한 로직, 레벨 3: 복잡, 레벨 4: 아키텍처)
- --security: 스캔에 보안 이슈 포함

### Loop 플래그

- --max N: 최대 반복 횟수 (기본값: 100)
- --auto: 레벨 1-2에 대한 자동 수정 적용 활성화

### MoAI (기본값) 플래그

- --loop: run 단계에서 반복 수정 활성화
- --max N: --loop 활성화 시 최대 수정 반복 횟수
- --branch: 구현 전 기능 브랜치 생성
- --pr: 완료 후 풀 리퀘스트 생성

---

## 레거시 명령어 매핑

이전 /moai:X-Y 명령어 형식을 새 /moai 서브커맨드 형식으로 매핑:

- /moai:0-project → /moai project
- /moai:1-plan → /moai plan
- /moai:2-run → /moai run
- /moai:3-sync → /moai sync
- /moai:9-feedback → /moai feedback
- /moai:fix → /moai fix
- /moai:loop → /moai loop
- /moai:moai → /moai (기본 자율 워크플로우)

참고: /moai:99-release는 별도의 로컬 전용 명령어로, /moai 스킬의 일부가 아닙니다.

---

## 설정 파일 참조

### 핵심 설정

- .moai/config/config.yaml: 메인 설정 파일 (섹션 파일에서 병합됨)
- .moai/config/sections/language.yaml: 언어 설정 (conversation_language, agent_prompt_language, code_comments)
- .moai/config/sections/user.yaml: 사용자 식별 (name)
- .moai/config/sections/quality.yaml: TRUST 5 프레임워크 설정, LSP 품질 게이트, 테스트 커버리지 목표
- .moai/config/sections/system.yaml: 시스템 메타데이터 (moai.version)

### 프로젝트 문서

- .moai/project/product.md: 제품 개요, 기능, 사용자 가치
- .moai/project/structure.md: 프로젝트 아키텍처 및 디렉토리 구성
- .moai/project/tech.md: 기술 스택, 의존성, 기술적 결정

### SPEC 문서

- .moai/specs/SPEC-XXX/spec.md: EARS 형식 요구사항이 있는 명세서 문서
- .moai/specs/SPEC-XXX/plan.md: 작업 분해가 있는 실행 계획
- .moai/specs/SPEC-XXX/acceptance.md: 인수 기준 및 테스트 계획

### 릴리즈 산출물

- CHANGELOG.md: 이중 언어 변경 로그 (버전별 영어 + 한국어)
- .moai/cache/release-snapshots/latest.json: 복구를 위한 릴리즈 상태 스냅샷

### 버전 파일 (릴리즈 시 5개 파일 동기화)

- pyproject.toml: 권위 있는 버전 소스
- pkg/version/version.go: 빌드 타임 주입이 있는 런타임 버전
- .moai/config/config.yaml: 설정 표시 버전
- .moai/config/sections/system.yaml: 시스템 메타데이터 버전
- internal/template/templates/: 바이너리 번들링을 위한 임베드 템플릿 디렉토리

---

## 완료 마커

AI가 워크플로우 상태를 알리기 위해 마커를 추가합니다:

- `<moai>DONE</moai>`: 단일 작업 또는 단계 완료
- `<moai>COMPLETE</moai>`: 전체 워크플로우 완료 (모든 단계 완료)

이 마커들은 loop 워크플로우에서 자동화 감지 및 루프 종료를 가능하게 합니다.

---

## 에러 처리 위임

- 품질 게이트 실패: 진단 및 해결을 위해 expert-debug 서브에이전트 사용
- 에이전트 실행 실패: 조사를 위해 expert-debug 서브에이전트 사용
- 토큰 한도 오류: /clear 실행 후 사용자가 --resume 플래그로 재개하도록 안내
- 권한 오류: .claude/settings.json 수동 검토
- 통합 오류: expert-devops 서브에이전트 사용
- MoAI-ADK 오류: GitHub 이슈 생성을 위해 /moai feedback 제안

---

Version: 1.1.0
Last Updated: 2026-01-28
