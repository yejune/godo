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
- `[ ]` 미시작
- `[~]` 진행중
- `[*]` 테스트중
- `[o]` 완료
- `[x]` 실패

> 상태 변경 시 변경일시 기록: `[o] 제목 (2026-02-11 17:30)`
```

**인자가 `status`이면**:
1. `.do/jobs/` 내 모든 체크리스트 파일 스캔
2. 각 파일의 [ ], [~], [*], [o], [x] 개수 집계
3. 전체 진행률 표시

## 서브 체크리스트

에이전트 실행단위별 서브 체크리스트 생성 시:
- 파일: `.do/jobs/{YY}/{MM}/{DD}/{title}/checklists/{order}_{agent-topic}.md`
- 템플릿:

```
# {agent-topic}

## 문제 요약

## 해결해야 하는 것

## 해결 방안

## Critical Files
1. **수정 대상:**
2. **참조 파일:**
3. **테스트 파일:**

## 해결 과정에서 겪은 문제

## 최종 소회 / 배운 점
```
