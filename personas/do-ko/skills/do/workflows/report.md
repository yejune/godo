---
name: do-workflow-report
description: >
  모든 체크리스트 항목 완료 후 완료 보고서를 생성한다. 서브 체크리스트
  결과, 교훈, 테스트 결과를 최종 report.md로 집계한다.
  Do의 체크리스트 기반 워크플로우에서 파생됨.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-15"
  tags: "report, completion, summary, lessons-learned"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 3000

# Do Extension: Triggers
triggers:
  keywords: ["report", "completion", "summary", "완료", "보고"]
  agents: ["manager-docs", "manager-quality"]
  phases: ["report"]
---

# 보고 워크플로우 오케스트레이션

## 목적

모든 체크리스트 항목이 완료된 후 완료 보고서를 생성한다. 서브 체크리스트 결과를 집계하고, 교훈을 수집하고, 테스트 결과를 요약하여 작업 디렉토리에 최종 report.md를 생성한다.

## 범위

- Do 체크리스트 기반 워크플로우의 마지막 단계
- checklist.md + checklists/*.md + git 이력을 소비
- report.md를 수행한 작업의 최종 기록으로 생성

## 입력

- `.do/jobs/{YY}/{MM}/{DD}/{title}/checklist.md`의 완료된 체크리스트
- `.do/jobs/{YY}/{MM}/{DD}/{title}/checklists/*.md`의 완료된 서브 체크리스트
- 작업 기간의 커밋 이력을 위한 git log

## 전제 조건

- [HARD] 모든 체크리스트 항목이 `[o]` (완료) 또는 `[x]` (실패) 상태여야 함 -- `[ ]` 또는 `[~]` 남아있으면 안 됨
- [HARD] 미완료 항목이 있으면 보고 워크플로우 실행 거부 후 실행 워크플로우로 안내
- [HARD] 서브 체크리스트의 테스트 실패는 보고서 생성을 차단 -- 먼저 테스트 수정

---

## 단계 순서

### Phase 1: 데이터 수집

보고서에 필요한 모든 정보 수집:

1. **체크리스트 요약**: checklist.md 읽기, 상태별 항목 집계
   - 전체 태스크, 완료 (`[o]`), 실패 (`[x]`), 차단 (`[!]`)
2. **서브 체크리스트 상세**: 각 checklists/{NN}_{agent}.md 읽기
   - 교훈 섹션 추출
   - 진행 로그에서 커밋 해시 추출
   - 핵심 파일 (수정된 파일) 추출
3. **git 이력**: 작업 기간의 `git log` 및 `git diff --stat` 실행
   - 라인 수가 포함된 변경 파일
   - 커밋 메시지 및 작성자
4. **테스트 결과**: 최신 테스트 실행 출력 수집
   - 통과/실패 수, 커버리지 비율

### Phase 2: 보고서 생성

에이전트: Task(report-agent) 또는 직접 생성

dev-checklist.md 템플릿을 사용하여 `.do/jobs/{YY}/{MM}/{DD}/{title}/report.md` 생성:

```markdown
## 완료 보고서

### 실행 요약
- 완료: {N}/{M} 태스크
- 기간: {시작} ~ {종료}

### 플랜 대비 변경사항
- (원래 플랜과의 차이점, 이유 포함)
- (변경 없으면 "플랜대로 진행")

### 테스트 결과
- 전체: {pass}/{total} 통과
- 커버리지: {N}% (측정 가능한 경우)
- 실패/스킵: 없음 또는 상세 내역

### 변경 파일 요약
- `path/to/file.go` -- 변경 내용 한 줄 요약
- `path/to/test.go` -- 추가된 테스트

### 미해결 사항
- (후속 작업 필요 항목)
- (알려진 제약사항)
- (없으면 "없음")

### 핵심 교훈
- (서브 체크리스트 교훈 종합)
- (팀/프로젝트에 공유할 인사이트)
```

### Phase 3: 검증

보고서 정확성 검증:

1. **파일 목록 일치**: 보고서의 변경 파일이 `git diff --stat`과 일치
2. **태스크 수 일치**: 실행 요약이 checklist.md 항목 수와 일치
3. **테스트 실패 없음**: 테스트 결과 섹션에 실패 없음
4. **교훈 작성됨**: 핵심 교훈 섹션이 비어있지 않음
5. **미해결 사항 추적**: 있으면 후속 플랜 또는 이슈 생성 제안

### Phase 4: 표시

사용자에게 보고서 요약 표시:

1. 실행 요약 표시 (완료된 태스크, 기간)
2. 핵심 교훈 표시 (상위 3개)
3. 미해결 사항 표시

AskUserQuestion 옵션:
- "좋아, 완료" -> 최종 확정
- "미해결 사항에 대한 후속 플랜 작성" -> 플랜 워크플로우 실행
- "전체 보고서 검토" -> 완전한 report.md 내용 표시

---

## 완료 기준

- Phase 1: 체크리스트, git, 테스트에서 모든 데이터 수집
- Phase 2: 올바른 `.do/jobs/` 경로에 report.md 생성
- Phase 3: git diff 및 체크리스트 수와 대조하여 보고서 검증
- Phase 4: 다음 단계 옵션과 함께 사용자에게 요약 표시
- [HARD] report.md에 테스트 실패 없음 (먼저 수정 후 보고)
- [HARD] 변경 파일 요약이 `git diff --stat`과 일치

---

Version: 1.0.0
Updated: 2026-02-16
