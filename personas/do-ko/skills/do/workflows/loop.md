---
name: do-workflow-loop
description: >
  모든 이슈가 해결되거나 최대 반복 횟수에 도달할 때까지
  스캔, 수정, 검증, 반복하는 자율 반복 수정 워크플로우.
  메모리 압박 감지 및 스냅샷 기반 재개를 포함한다.
  반복적인 오류 해결 또는 지속적인 수정이 필요할 때 사용.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "2.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-07"
  tags: "loop, iterative, auto-fix, diagnostics, testing, coverage"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# Do Extension: Triggers
triggers:
  keywords: ["loop", "iterate", "repeat", "until done", "keep fixing", "all errors"]
  agents: ["expert-debug", "expert-backend", "expert-frontend", "expert-testing"]
  phases: ["loop"]
---

# 워크플로우: Loop - 반복 자율 수정

목적: 모든 이슈가 해결될 때까지 반복 자율 수정. AI가 스캔, 수정, 검증을 반복하며 완료 조건이 충족되거나 최대 반복 횟수에 도달할 때까지 계속한다.

흐름: 완료 확인 -> 메모리 확인 -> 진단 -> 수정 -> 검증 -> 반복

## 지원 플래그

- --max N (별칭 --max-iterations): 최대 반복 횟수 (기본값 100)
- --auto-fix: 자동 수정 활성화 (기본값 수준 1)
- --sequential (별칭 --seq): 병렬 대신 순차 진단
- --errors (별칭 --errors-only): 오류만 수정, 경고 생략
- --coverage (별칭 --include-coverage): 커버리지 임계값 포함 (기본값 85%)
- --memory-check: 메모리 압박 감지 활성화
- --resume ID (별칭 --resume-from): 스냅샷에서 복원

## 반복당 사이클

각 반복은 다음 단계를 순서대로 실행한다:

Step 1 - 완료 확인:
- 이전 반복 응답에서 완료 마커 확인
- 마커 유형: `<do>DONE</do>`, `<do>COMPLETE</do>`
- 마커 발견 시: 성공으로 루프 종료

Step 2 - 메모리 압박 확인 (--memory-check 활성화 시):
- 시작 시간부터 세션 지속 시간 계산
- 반복 시간에서 GC 압박 징후 모니터링 (반복 시간이 두 배로 늘어남)
- 세션 지속 시간이 25분을 초과하거나 반복 시간이 두 배가 될 경우:
  - $CLAUDE_PROJECT_DIR/.do/cache/loop-snapshots/memory-pressure.json에 선제적 체크포인트 저장
  - 메모리 압박에 대해 사용자에게 경고
  - /do:loop --resume memory-pressure로 재개 제안
- 메모리 안전 한도 도달 시 (50회 반복): 체크포인트와 함께 종료

Step 3 - 병렬 진단:
- run_in_background와 함께 Bash를 사용해 네 가지 진단 도구 동시 실행
- 도구 1: 감지된 언어에 대한 LSP 진단
- 도구 2: sgconfig.yml 규칙을 사용한 AST-grep 스캔
- 도구 3: 감지된 언어의 테스트 러너 (pytest, jest, go test, cargo test)
- 도구 4: 커버리지 측정 (coverage.py, c8, go test -cover, cargo tarpaulin)
- 각 백그라운드 태스크에 TaskOutput으로 결과 수집
- 메트릭이 포함된 통합 진단 보고서로 집계: 오류 수, 경고 수, 테스트 통과율, 커버리지 비율

--sequential 플래그 지정 시: LSP, 그 다음 AST-grep, 그 다음 테스트, 그 다음 커버리지 순으로 순차 실행.

Step 4 - 완료 조건 확인:
- 조건: 오류 없음 AND 모든 테스트 통과 AND 커버리지 임계값 충족
- 모든 조건 충족 시: 완료 마커 추가 또는 계속 진행 여부 사용자에게 확인

Step 5 - 태스크 생성:
- [HARD] 새로 발견된 모든 이슈에 pending 상태로 TaskCreate

Step 6 - 수정 실행:
- [HARD] 각 수정 전: TaskUpdate로 항목을 in_progress로 변경
- [HARD] 에이전트 위임 강제: 모든 수정 작업은 반드시 전문 에이전트에게 위임해야 한다. 직접 수정 절대 금지.

이슈 유형별 에이전트 선택:
- 타입 오류, 로직 버그: expert-debug 서브에이전트
- import/모듈 이슈: expert-backend 또는 expert-frontend 서브에이전트
- 테스트 실패: expert-testing 서브에이전트
- 보안 이슈: expert-security 서브에이전트
- 성능 이슈: expert-performance 서브에이전트

--auto 설정별 적용 수정 수준:
- 수준 1 (즉시): 승인 불필요. import 정렬, 공백
- 수준 2 (안전): 로그만 기록. 변수 이름 변경, 타입 추가
- 수준 3 (승인): AskUserQuestion 필요. 로직 변경, API 수정
- 수준 4 (수동): 자동 수정 불가. 보안, 아키텍처

Step 7 - 검증:
- [HARD] 각 수정 후: TaskUpdate로 항목을 completed로 변경

Step 8 - 스냅샷 저장:
- $CLAUDE_PROJECT_DIR/.do/cache/loop-snapshots/에 반복 스냅샷 저장
- 반복 카운터 증가

Step 9 - 반복 또는 종료:
- 최대 반복 횟수 도달 시: 남은 이슈와 옵션 표시
- 그렇지 않으면: Step 1로 돌아감

## 완료 조건

다음 조건 중 하나가 충족될 때 루프가 종료된다:
- 응답에서 완료 마커 감지
- 모든 조건 충족: 오류 없음 + 테스트 통과 + 커버리지 임계값
- 최대 반복 횟수 도달 (남은 이슈 표시)
- 메모리 압박 임계값 초과 (체크포인트 저장)
- 사용자 중단 (상태 자동 저장)

## 스냅샷 관리

스냅샷 위치: $CLAUDE_PROJECT_DIR/.do/cache/loop-snapshots/

파일:
- iteration-001.json, iteration-002.json 등 (반복별 스냅샷)
- latest.json (최신 파일에 대한 심볼릭 링크)
- memory-pressure.json (메모리 압박 시 선제적 체크포인트)

루프 상태 파일: $CLAUDE_PROJECT_DIR/.do/cache/.do_loop_state.json

재개 커맨드:
- /do:loop --resume latest
- /do:loop --resume iteration-002
- /do:loop --resume memory-pressure

## 언어별 커맨드

Python: pytest --tb=short (테스트), coverage run -m pytest (커버리지)
TypeScript: npm test 또는 jest (테스트), npm run coverage (커버리지)
Go: go test ./... (테스트), go test -cover ./... (커버리지)
Rust: cargo test (테스트), cargo tarpaulin (커버리지)

언어 감지: pyproject.toml (Python), package.json (TypeScript/JavaScript), go.mod (Go), Cargo.toml (Rust)

## 취소

루프를 중단하려면 아무 메시지나 전송. 상태는 session_end 훅을 통해 자동 저장된다.

## 안전 개발 프로토콜

루프 내 모든 수정은 CLAUDE.md 섹션 7 안전 개발 프로토콜을 따른다:
- 재현 우선: 버그 수정 전 실패 테스트 작성
- 수정 후 검토: 각 수정 사이클 후 잠재적 부작용 나열
- 개별 작업당 최대 3회 재시도 (CLAUDE.md 헌법 기준)

## 실행 요약

1. 인수 파싱 (플래그 추출: --max, --auto-fix, --sequential, --errors, --coverage, --memory-check, --resume)
2. --resume 지정 시: 지정된 스냅샷에서 상태 로드 후 계속
3. 지시자 파일에서 프로젝트 언어 감지
4. 반복 카운터 및 메모리 추적 초기화 (시작 시간)
5. 루프: 반복당 사이클 실행 (Step 1-9)
6. 종료 시: 근거와 함께 최종 요약 보고
7. 메모리 체크포인트 생성 시: 재개 지침 표시

---

Version: 2.0.0
Source: loop.md command v2.2.0
