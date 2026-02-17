# Go 보조 패키지 분리: lint/scaffold/profile/statusline/rank/glm
상태: [ ] | 담당: expert-backend (에이전트 G) | 작성 언어: ko

## Problem Summary
- godo의 `package main`에서 lint, scaffold, profile, statusline, rank, glm 관련 코드를 독립 패키지로 분리
- 핵심 패키지(hook/mode/persona)보다 단순하고 상호 의존성 낮음

## Acceptance Criteria
> 완료 시 [ ] → [o]로 변경. [x]는 "실패"를 의미하므로 절대 사용 금지.
- [ ] internal/lint/runner.go, gate.go, setup.go 생성
- [ ] internal/scaffold/create.go 생성
- [ ] internal/profile/profile.go 생성
- [ ] internal/statusline/statusline.go 생성
- [ ] internal/rank/auth.go, client.go, config.go, transcript.go 생성
- [ ] internal/glm/glm.go 생성
- [ ] `go build ./...` 통과
- [ ] 커밋 완료

## Solution Approach
- godo 소스에서 각 파일을 읽고 `package main` → 독립 패키지로 변환
- lint 패키지:
  - lint.go + lint_*.go → runner.go (실행기), gate.go (게이트 로직), setup.go (설정)
- scaffold 패키지:
  - create.go → create.go (에이전트/스킬 템플릿 생성)
- profile 패키지:
  - claude_profile.go → profile.go (--profile 플래그 처리)
- statusline 패키지:
  - statusline.go → statusline.go (상태 줄 렌더링)
- rank 패키지:
  - rank*.go (5파일) → auth.go, client.go, config.go, transcript.go (4파일로 재구성)
- glm 패키지:
  - glm.go → glm.go (GLM 백엔드)
- 대안: 작은 패키지들을 utils/로 합치기 → 기각 (각 기능이 독립적이므로 패키지 분리 유지)

## Critical Files
- **소스** (do-focus/cmd/godo/):
  - `lint.go`, `lint_setup.go`, `lint_gate.go`, `lint_runner.go`
  - `create.go`
  - `claude_profile.go`
  - `statusline.go`
  - `rank.go`, `rank_auth.go`, `rank_client.go`, `rank_config.go`, `rank_transcript.go`
  - `glm.go`
- **생성 대상** (convert/internal/):
  - `lint/runner.go`, `lint/gate.go`, `lint/setup.go`
  - `scaffold/create.go`
  - `profile/profile.go`
  - `statusline/statusline.go`
  - `rank/auth.go`, `rank/client.go`, `rank/config.go`, `rank/transcript.go`
  - `glm/glm.go`

## Risks
- rank 패키지가 외부 API에 의존할 수 있음 — HTTP 클라이언트 인터페이스로 추상화
- lint 패키지가 외부 도구(ruff, eslint 등)를 exec하므로 경로 처리 확인 필요

## Progress Log
(작업 시작 시 기록)

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` — 변경된 파일만 스테이징 (lint/*.go, scaffold/*.go, profile/*.go, statusline/*.go, rank/*.go, glm/*.go)
- [ ] `git diff --cached` — 의도한 변경만 포함되었는지 확인
- [ ] `git commit` — 커밋 메시지에 WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록
⚠️ 이 섹션을 완료하지 않으면 작업은 미완료(incomplete) 상태임

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 개선 조치:
