---
description: Do 프레임워크 실행 모드 전환 (do/focus/auto)
allowed-tools: Bash, Read, Write, Edit, AskUserQuestion
---

# /do:mode - 실행 모드 전환

## 실행

### Step 1: 현재 모드 확인

`godo mode get` 실행하여 현재 모드 확인.

### Step 2: 모드 표시 또는 전환

**인자가 없으면**: 현재 모드와 사용 가능한 모드 표시

**인자가 있으면** ($ARGUMENTS):
1. 실행 모드: "do", "focus", "auto" → Step 3
2. 권한 모드: "bypass", "accept", "default", "plan" → Step 3-B
3. 유효하지 않으면 오류 메시지 표시하고 종료

### Step 3: 실행 모드 전환

두 가지 동시 업데이트:

1. **즉시 반영** (statusline): `godo mode set {모드}` 실행
2. **영구 저장** (다음 세션): `.claude/settings.local.json`의 `env.DO_MODE` 업데이트

### Step 3-B: 권한 모드 전환

`godo mode {bypass|accept|default|plan}` 실행.
settings.local.json의 defaultMode가 업데이트됨.
"claude --continue 로 재시작하면 적용됩니다" 안내.

### Step 4: 결과 표시

**현재 모드 표시 형식**:

```
현재 모드: {모드명} ({모드 설명})

실행 모드: do, focus, auto
권한 모드: bypass, accept, default, plan

전환: /do:mode <모드명>
```

**모드 전환 후 표시 형식**:

```
모드 전환: {이전 모드} -> {새 모드}
{모드 설명}
```

## 실행 모드

### do (Strategic Orchestrator)
- 모든 작업을 전문 에이전트에게 위임
- 병렬 실행 우선
- 복잡한 멀티 도메인 작업에 최적

### focus (Direct Executor)
- 간단한 1-3 파일 변경을 직접 처리
- 에이전트 위임 없이 빠른 실행
- 단순 수정, 작은 버그 수정에 적합

### auto (Automatic Selection)
- 작업 복잡도에 따라 자동으로 모드 선택
- 간단한 작업은 focus, 복잡한 작업은 do
- 권장 모드

## 권한 모드

### bypass (bypassPermissions)
- 모든 권한 프롬프트 건너뛰기
- 신뢰할 수 있는 환경에서만 사용
- 재시작 필요

### accept (acceptEdits)
- 파일 편집 자동 수락, Bash만 확인
- 빠른 작업에 적합
- 재시작 필요

### default
- 표준 동작 - 첫 사용 시 권한 요청
- 재시작 필요

### plan
- 읽기 전용 - 파일 수정/명령 실행 불가
- 분석/계획 전용
- 재시작 필요
