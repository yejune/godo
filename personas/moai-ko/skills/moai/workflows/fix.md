---
name: moai-workflow-fix
description: >
  병렬 스캔 및 분류를 통한 원샷 자율 수정 워크플로우입니다. LSP 오류, 린팅 이슈,
  타입 오류를 찾아 심각도별로 분류하고, 에이전트 위임을 통해 안전한 수정을 적용한 후
  결과를 보고합니다. 오류 수정, 린팅 이슈 해결, 진단 실행 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "2.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-07"
  tags: "fix, auto-fix, lsp, linting, diagnostics, errors, type-check"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["fix", "auto-fix", "error", "lint", "diagnostic", "lsp", "type error"]
  agents: ["expert-debug", "expert-backend", "expert-frontend", "expert-refactoring"]
  phases: ["fix"]
---

# 워크플로우: Fix - 원샷 자동 수정

목적: 병렬 스캔 및 분류를 통한 원샷 자율 수정입니다. AI가 이슈를 찾아 심각도별로 분류하고, 안전한 수정을 적용한 후 결과를 보고합니다.

흐름: 병렬 스캔 -> 분류 -> 수정 -> 검증 -> 보고

## 지원 플래그

- --dry (별칭 --dry-run): 미리보기만 표시, 변경 사항 미적용
- --sequential (별칭 --seq): 병렬 대신 순차 스캔
- --level N: 적용할 최대 수정 레벨 (기본값 3)
- --errors (별칭 --errors-only): 오류만 수정, 경고 건너뜀
- --security (별칭 --include-security): 스캔에 보안 이슈 포함
- --no-fmt (별칭 --no-format): 포맷팅 수정 건너뜀
- --resume [ID] (별칭 --resume-from): 스냅샷에서 재개 (ID 없으면 최신 사용)
- --team: 팀 기반 디버깅 활성화 (경쟁 가설 조사에 대해서는 team-debug.md 참조)

## Phase 1: 병렬 스캔

run_in_background와 함께 Bash를 사용하여 세 가지 진단 도구를 동시에 실행합니다 (8초 vs 30초로 3-4배 속도 향상).

스캐너 1 - LSP 진단:
- 언어별 타입 검사 및 오류 감지
- Python: mypy --output json
- TypeScript: tsc --noEmit
- Go: go vet ./...

스캐너 2 - AST-grep 스캔:
- sgconfig.yml 규칙을 이용한 구조적 패턴 매칭
- 보안 패턴 및 코드 품질 규칙

스캐너 3 - 린터:
- 언어별 린팅
- Python: ruff check --output-format json
- TypeScript: eslint --format json
- Go: golangci-lint run --out-format json
- Rust: cargo clippy --message-format json

모든 스캐너 완료 후:
- 각 도구의 출력을 구조화된 이슈 목록으로 파싱
- 여러 스캐너에 나타난 중복 이슈 제거
- 심각도별 정렬: Critical, High, Medium, Low
- 효율적인 수정을 위해 파일 경로별 그룹화

언어 자동 감지는 지시자 파일을 사용합니다: pyproject.toml (Python), package.json (TypeScript/JavaScript), go.mod (Go), Cargo.toml (Rust). 16개 언어를 지원합니다.

오류 처리: 스캐너 중 하나라도 실패하면 성공한 스캐너의 결과로 계속 진행합니다. 보고서에 실패한 스캐너를 기록합니다.

--sequential 플래그 사용 시: LSP, AST-grep, 린터를 순서대로 실행합니다.

## Phase 2: 분류

이슈는 네 가지 레벨로 분류됩니다:

- 레벨 1 (즉시): 승인 불필요. 예: import 정렬, 공백, 포맷팅
- 레벨 2 (안전): 로그만 기록, 승인 불필요. 예: 변수 이름 변경, 타입 어노테이션 추가
- 레벨 3 (검토): 사용자 승인 필요. 예: 로직 변경, API 수정
- 레벨 4 (수동): 자동 수정 불허. 예: 보안 취약점, 아키텍처 변경

## Phase 3: 자동 수정

[HARD] 에이전트 위임 의무: 모든 수정 작업은 반드시 전문 에이전트에게 위임해야 합니다. 수정을 직접 실행하지 마세요.

수정 레벨별 에이전트 선택:
- 레벨 1 (import, 포맷팅): expert-backend 또는 expert-frontend 서브에이전트
- 레벨 2 (이름 변경, 타입): expert-refactoring 서브에이전트
- 레벨 3 (로직, API): expert-debug 또는 expert-backend 서브에이전트 (사용자 승인 후)

실행 순서:
- 레벨 1 수정은 에이전트 위임을 통해 자동 적용
- 레벨 2 수정은 로깅과 함께 자동 적용
- 레벨 3 수정은 AskUserQuestion 승인 요청 후 에이전트에 위임
- 레벨 4 수정은 수동 조치 항목으로 보고서에 나열

--dry 플래그 사용 시: 분류된 모든 이슈의 미리보기를 표시하고 변경 없이 종료합니다.

## Phase 4: 검증

- 수정된 파일에서 영향받은 진단을 재실행
- 수정이 대상 이슈를 해결했는지 확인
- 수정으로 인한 회귀 감지

## 작업 추적

[HARD] 작업 관리 도구 필수:
- 발견된 모든 이슈를 TaskCreate로 pending 상태로 추가
- 각 수정 전: TaskUpdate로 in_progress로 변경
- 각 수정 후: TaskUpdate로 completed로 변경

## 안전 개발 프로토콜

모든 수정은 CLAUDE.md 섹션 7 안전 개발 프로토콜을 따릅니다:
- 재현 우선: 수정 전에 버그를 재현하는 실패 테스트 작성
- 접근 방식 우선: 레벨 3+ 수정의 경우, 적용 전 접근 방식 설명
- 수정 후 검토: 각 수정 후 잠재적 부작용 나열

## 스냅샷 저장/재개

스냅샷 위치: $CLAUDE_PROJECT_DIR/.moai/cache/fix-snapshots/

스냅샷 내용:
- 타임스탬프
- 대상 경로
- 발견, 수정, 대기 중인 이슈 수
- 현재 수정 레벨
- TODO 상태
- 스캔 결과

재개 명령어:
- /moai:fix --resume (최신 스냅샷 사용)
- /moai:fix --resume fix-20260119-143052 (특정 스냅샷 사용)

## 팀 모드

--team 플래그가 제공되면 경쟁 가설을 사용하는 팀 기반 디버깅 워크플로우에 위임합니다.

팀 구성: 3개의 가설 에이전트 (haiku)가 병렬로 다양한 근본 원인을 탐구합니다.

상세한 팀 오케스트레이션 단계는 workflows/team-debug.md를 참조하세요.

폴백: 팀 모드를 사용할 수 없는 경우 표준 단일 에이전트 수정 워크플로우로 계속합니다.

## 실행 요약

1. 인수 파싱 (플래그 추출: --dry, --sequential, --level, --errors, --security, --resume)
2. --resume 사용 시: 스냅샷 로드 후 저장된 상태에서 계속
3. 지시자 파일에서 프로젝트 언어 감지
4. 병렬 스캔 실행 (LSP + AST-grep + 린터)
5. 결과 집계 및 중복 제거
6. 레벨 1-4로 분류
7. 발견된 모든 이슈에 대해 TaskCreate
8. --dry 사용 시: 미리보기 표시 후 종료
9. 에이전트 위임을 통해 레벨 1-2 수정 적용
10. AskUserQuestion을 통해 레벨 3 수정 승인 요청
11. 진단 재실행으로 수정 검증
12. $CLAUDE_PROJECT_DIR/.moai/cache/fix-snapshots/에 스냅샷 저장
13. 증거와 함께 보고 (file:line 변경 사항)

---

Version: 2.0.0
Source: fix.md command v2.2.0
