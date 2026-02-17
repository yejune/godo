---
description: 체크리스트 생성 및 관리
allowed-tools: Read, Write, Edit, Glob, Grep, Bash
---

# /do:checklist

체크리스트 시스템을 관리합니다.

## 동작

**인자가 없으면 ($ARGUMENTS가 비어있으면)**:
1. `.do/jobs/` 디렉토리에서 가장 최근 체크리스트 파일 찾기
2. 없으면 "활성 체크리스트가 없습니다" 표시
3. 있으면 현재 상태 요약 표시 (완료/전체 비율)

**인자가 `create {title}`이면**:
1. 현재 날짜로 폴더명 생성 (YY/MM/DD 형식)
2. 메인 파일 생성: `.do/jobs/{YY}/{MM}/{DD}/{title}/checklist.md`
3. 서브 디렉토리 생성: `.do/jobs/{YY}/{MM}/{DD}/{title}/checklists/`
4. 메인 파일에 템플릿 작성:

```
# Checklist: {title}
생성일: {yyyy-mm-dd HH:MM}

## 작업 목록

- [ ] (작업 항목을 추가하세요)

## 상태 범례
- `[ ]` 미시작 (pending)
- `[~]` 진행중 (in progress)
- `[*]` 테스트중 (testing)
- `[!]` 블로커 (blocked)
- `[o]` 완료 (done) -- 커밋 해시 필수
- `[x]` 실패 (failed)

> 금지된 전이: [ ]->[o] (테스트 없이 완료 불가), [ ]->[x], [ ]->[*]
> 상태 변경 시 변경일시 기록: `[o] 제목 (2026-02-11 17:30, commit: a1b2c3d)`
```

**인자가 `status`이면**:
1. `.do/jobs/` 내 모든 체크리스트 파일 스캔
2. 각 파일의 [ ], [~], [*], [o], [x] 개수 집계
3. 전체 진행률 표시

## 서브 체크리스트

에이전트 실행단위별 서브 체크리스트 생성 시:
- 파일: `.do/jobs/{YY}/{MM}/{DD}/{title}/checklists/{order}_{agent-topic}.md`
- 템플릿:

```markdown
# {agent-topic}: {작업 제목}
상태: [ ] | 담당: {에이전트}

## Problem Summary
- 무엇을 해결하는가
- 왜 이 작업이 필요한가

## Acceptance Criteria
- [ ] 측정 가능한 완료 조건
- [ ] 검증 완료 (테스트 또는 대안 검증)
- [ ] 커밋 완료

## Test Strategy
- {unit: file_test.go | pass (빌드 확인: go build ./...)}

## Solution Approach
- 선택한 접근법 (대안 최소 1개 언급)

## Critical Files
- **수정 대상**: `path/to/file` -- 변경 이유
- **참조 파일**: `path/to/ref` -- 참조 이유
- **테스트 파일**: `path/to/test`

## Risks
- 깨질 수 있는 것 / 주의할 점

## Progress Log
- {timestamp} [~] 작업 시작

## FINAL STEP: Commit (절대 생략 금지)
- [ ] `git add` -- 변경된 파일만 스테이징 (본인 파일만)
- [ ] `git diff --cached` -- 의도한 변경만 포함 확인
- [ ] `git commit` -- WHY 포함
- [ ] 커밋 해시를 Progress Log에 기록

## Lessons Learned (완료 시 작성)
- 잘된 점:
- 어려웠던 점:
- 다음에 다르게 할 점:
```
