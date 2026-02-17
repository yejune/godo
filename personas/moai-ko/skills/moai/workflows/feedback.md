---
name: moai-workflow-feedback
description: >
  사용자 피드백, 버그 리포트, 기능 제안을 수집하고 manager-quality 에이전트를 통해
  GitHub 이슈를 자동으로 생성합니다. 버그 리포트, 기능 요청, 질문을 우선순위 분류와 함께
  지원합니다. 피드백 제출, 버그 리포트, 기능 요청 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "2.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-07"
  tags: "feedback, bug-report, feature-request, github-issues, quality"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["feedback", "bug", "issue", "suggestion", "report", "feature request"]
  agents: ["manager-quality"]
  phases: ["feedback"]
---

# 워크플로우: feedback - GitHub 이슈 생성

목적: 사용자 피드백, 버그 리포트, 기능 제안을 수집하고 manager-quality 에이전트를 통해 GitHub 이슈를 자동으로 생성합니다.

전제 조건: `gh` CLI가 설치되어 있고 인증되어 있어야 합니다 (`gh auth status`). 설치되지 않은 경우 https://cli.github.com/ 을 통해 설치를 안내하세요.

---

## Phase 1: 피드백 수집

### Step 1: 피드백 유형 결정

[HARD] $ARGUMENTS가 제공된 경우 피드백 유형을 확인합니다 (issue, suggestion, question).

$ARGUMENTS가 비어 있으면 AskUserQuestion을 사용하세요:

질문: 어떤 유형의 피드백을 제출하시겠습니까?

옵션:

- 버그 리포트: 발생한 기술적 문제나 오류
- 기능 요청: 개선 사항이나 새 기능에 대한 제안
- 질문: 설명이나 도움이 필요한 사항

### Step 2: 세부 내용 수집

[HARD] AskUserQuestion을 통해 피드백 제목을 사용자에게 입력받습니다 (자유 텍스트 입력).

[HARD] AskUserQuestion을 통해 상세 설명을 사용자에게 입력받습니다 (자유 텍스트 입력).

[SOFT] 사용자에게 우선순위를 입력받습니다:

- 낮음: 사소한 문제, 우회 방법 있음
- 중간: 보통 수준의 영향, 긴급한 우회 방법 불필요
- 높음: 중요한 영향, 워크플로우 차단

---

## Phase 2: GitHub 이슈 생성

[HARD] 수집된 피드백 세부 내용과 함께 manager-quality 서브에이전트에게 위임합니다.

manager-quality에 전달할 내용:

- 피드백 유형: 버그 리포트, 기능 요청 또는 질문
- 제목: 사용자가 입력한 제목
- 설명: 사용자가 입력한 설명
- 우선순위: 선택한 우선순위 수준
- 대화 언어: 설정에서 가져옴

### GitHub 이슈 레이블

- 버그 리포트: 레이블 "bug"
- 기능 요청: 레이블 "enhancement"
- 질문: 레이블 "question"

### 이슈 생성 명령어

manager-quality 에이전트가 실행합니다: gh issue create --repo modu-ai/moai-adk

이슈 본문은 다음을 포함하는 일관된 템플릿을 사용합니다:

- 피드백 유형 헤더
- 설명 내용
- 우선순위 수준
- 환경 정보 (MoAI 버전, OS)

### 결과 보고

[HARD] 생성된 이슈 URL을 사용자에게 제공합니다.
[HARD] 피드백 제출 성공을 사용자에게 확인합니다.

사용자의 conversation_language로 표시:

- 이슈 번호와 제목
- 생성된 이슈의 직접 URL
- 적용된 레이블

---

## 제출 후 옵션

성공적으로 제출된 후 AskUserQuestion을 사용하세요:

- 개발 계속: 현재 개발 워크플로우로 돌아가기
- 추가 피드백 제출: 다른 이슈나 제안 리포트
- 이슈 보기: 생성된 GitHub 이슈를 브라우저에서 열기

---

## 실행 패턴

이 워크플로우는 단순한 순차 실행을 사용합니다 (병렬 처리 불필요):

- Phase 1에서 MoAI 오케스트레이터 레벨에서 모든 사용자 입력을 수집합니다
- Phase 2에서 완전한 컨텍스트와 함께 manager-quality에 위임합니다
- 단일 에이전트가 전체 제출 프로세스를 처리합니다
- 일반적인 실행은 30초 이내에 완료됩니다

재개 지원: 해당 없음 (원자적 작업).

---

## 에이전트 체인 요약

- Phase 1: MoAI 오케스트레이터 (피드백 수집을 위한 AskUserQuestion)
- Phase 2: manager-quality 서브에이전트 (gh CLI를 통한 GitHub 이슈 생성)

---

Version: 2.0.0
Last Updated: 2026-02-07
