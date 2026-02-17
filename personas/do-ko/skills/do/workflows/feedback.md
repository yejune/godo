---
name: do-workflow-feedback
description: >
  사용자 피드백, 버그 리포트, 기능 제안을 수집하고
  manager-quality 에이전트를 통해 GitHub 이슈를 자동 생성한다.
  버그 리포트, 기능 요청, 질문을 우선순위 분류와 함께 지원한다.
  피드백 제출, 버그 신고, 기능 요청 시 사용.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "2.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-07"
  tags: "feedback, bug-report, feature-request, github-issues, quality"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# Do Extension: Triggers
triggers:
  keywords: ["feedback", "bug", "issue", "suggestion", "report", "feature request"]
  agents: ["manager-quality"]
  phases: ["feedback"]
---

# 워크플로우: feedback - GitHub 이슈 생성

목적: 사용자 피드백, 버그 리포트, 기능 제안을 수집하고 manager-quality 에이전트를 통해 GitHub 이슈를 자동 생성한다.

전제 조건: `gh` CLI가 설치 및 인증되어 있어야 한다 (`gh auth status`). 미설치 시 https://cli.github.com/ 을 통한 설치 안내.

---

## Phase 1: 피드백 수집

### Step 1: 피드백 유형 결정

[HARD] $ARGUMENTS가 있는 경우 피드백 유형을 우선 해석한다 (issue, suggestion, question).

$ARGUMENTS가 비어있는 경우 AskUserQuestion 사용:

질문: 어떤 유형의 피드백을 제출하시겠습니까?

옵션:

- 버그 리포트: 기술적 문제나 오류 발생
- 기능 요청: 개선 사항이나 신규 기능 제안
- 질문: 명확화 또는 도움 요청

### Step 2: 상세 내용 수집

[HARD] AskUserQuestion으로 피드백 제목을 사용자에게 요청한다 (자유 텍스트 입력).

[HARD] AskUserQuestion으로 상세 설명을 사용자에게 요청한다 (자유 텍스트 입력).

[SOFT] 사용자에게 우선순위 수준 요청:

- 낮음: 사소한 문제, 우회 방법 존재
- 중간: 보통 수준의 영향, 긴급 우회 불필요
- 높음: 큰 영향, 워크플로우 차단

---

## Phase 2: GitHub 이슈 생성

[HARD] 수집한 피드백 상세 내용을 manager-quality 서브에이전트에 위임한다.

manager-quality에 전달:

- 피드백 유형: 버그 리포트, 기능 요청, 또는 질문
- 제목: 사용자 입력 제목
- 설명: 사용자 입력 설명
- 우선순위: 선택한 우선순위 수준
- 대화 언어: 설정에서 가져옴

### GitHub 이슈 레이블

- 버그 리포트: labels "bug"
- 기능 요청: labels "enhancement"
- 질문: labels "question"

### 이슈 생성 커맨드

manager-quality 에이전트가 실행: gh issue create --repo do-focus/do-focus

이슈 본문은 다음을 포함하는 일관된 템플릿 사용:

- 피드백 유형 헤더
- 설명 내용
- 우선순위 수준
- 환경 정보 (Do 버전, OS)

### 결과 보고

[HARD] 생성된 이슈 URL을 사용자에게 제공한다.
[HARD] 피드백 제출 성공 여부를 사용자에게 확인한다.

사용자의 conversation_language로 표시:

- 이슈 번호 및 제목
- 생성된 이슈 직접 URL
- 적용된 레이블

---

## 제출 후 옵션

성공적인 제출 후 AskUserQuestion 사용:

- 개발 계속: 현재 개발 워크플로우로 복귀
- 추가 피드백 제출: 다른 이슈나 제안 신고
- 이슈 보기: 생성된 GitHub 이슈를 브라우저에서 열기

---

## 실행 패턴

이 워크플로우는 단순 순차 실행을 사용한다 (병렬 처리 불필요):

- Phase 1에서 Do 오케스트레이터 수준에서 모든 사용자 입력 수집
- Phase 2에서 완전한 컨텍스트와 함께 manager-quality에 위임
- 단일 에이전트가 전체 제출 프로세스 처리
- 일반적으로 30초 이내에 실행 완료

재개 지원: 해당 없음 (원자적 작업).

---

## 에이전트 체인 요약

- Phase 1: Do 오케스트레이터 (피드백 수집을 위한 AskUserQuestion)
- Phase 2: manager-quality 서브에이전트 (gh CLI를 통한 GitHub 이슈 생성)

---

Version: 2.0.0
Last Updated: 2026-02-07
