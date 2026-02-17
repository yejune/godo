---
name: do-workflow-project
description: >
  코드베이스 분석 또는 사용자 입력으로부터 프로젝트 문서를 생성한다.
  .do/project/ 디렉토리에 product.md, structure.md, tech.md를 생성한다.
  신규 및 기존 프로젝트 유형과 LSP 서버 감지를 지원한다.
  프로젝트 초기화 또는 프로젝트 문서 생성 시 사용.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "2.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-07"
  tags: "project, documentation, initialization, codebase-analysis, setup"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# Do Extension: Triggers
triggers:
  keywords: ["project", "init", "documentation", "setup", "initialize"]
  agents: ["manager-project", "manager-docs", "Explore", "expert-devops"]
  phases: ["project"]
---

# 워크플로우: project - 프로젝트 문서 생성

목적: 코드베이스 분석으로부터 프로젝트 문서를 생성한다. .do/project/ 디렉토리에 product.md, structure.md, tech.md를 생성한다.

---

## Phase 0: 프로젝트 유형 감지

[HARD] 분석 전 AskUserQuestion으로 먼저 프로젝트 유형을 확인한다.

질문: 어떤 유형의 프로젝트를 작업 중이신가요?

옵션 (사용자의 conversation_language로):

- 신규 프로젝트: 처음부터 시작, 프로젝트 정보를 대화형으로 수집
- 기존 프로젝트: 기존 코드베이스 문서화, 코드를 자동으로 분석

라우팅:

- 신규 프로젝트 선택 시: Phase 0.5로 진행
- 기존 프로젝트 선택 시: Phase 1로 진행

---

## Phase 0.5: 신규 프로젝트 정보 수집 (신규 프로젝트만)

목표: 분석할 기존 코드가 없을 때 프로젝트 상세 정보 수집.

질문 1 - 프로젝트 목적 (AskUserQuestion):

- 웹 애플리케이션: 프론트엔드, 백엔드, 또는 풀스택 웹앱
- API 서비스: REST API, GraphQL, 또는 마이크로서비스
- CLI 도구: 커맨드라인 유틸리티 또는 자동화 도구
- 라이브러리/패키지: 재사용 가능한 코드 라이브러리 또는 SDK

질문 2 - 주 언어 (AskUserQuestion):

- Python: 백엔드, 데이터 과학, 자동화
- TypeScript/JavaScript: 웹, Node.js, 프론트엔드
- Go: 고성능 서비스, CLI 도구
- 기타: Rust, Java, Ruby 등 (상세 입력 요청)

질문 3 - 프로젝트 설명 (자유 텍스트 입력):

- 프로젝트명
- 주요 기능 또는 목표
- 대상 사용자

수집 후 사용자 입력으로 초기 문서를 생성하고 Phase 4로 진행.

---

## Phase 1: 코드베이스 분석 (기존 프로젝트만)

[HARD] 코드베이스 분석을 Explore 서브에이전트에 위임한다.

[SOFT] 포괄적인 분석을 위해 --ultrathink 적용.

Explore 에이전트에 전달하는 분석 목표:

- 프로젝트 구조: 주요 디렉토리, 진입점, 아키텍처 패턴
- 기술 스택: 언어, 프레임워크, 주요 의존성
- 핵심 기능: 주요 기능 및 비즈니스 로직 위치
- 빌드 시스템: 빌드 도구, 패키지 매니저, 스크립트

Explore 에이전트의 예상 출력:

- 감지된 주 언어
- 식별된 프레임워크
- 아키텍처 패턴 (MVC, Clean Architecture, Microservices 등)
- 매핑된 주요 디렉토리 (소스, 테스트, 설정, 문서)
- 목적별로 분류된 의존성
- 식별된 진입점

실행 모드:

- 신규 문서화: .do/project/가 비어있을 때 세 파일 모두 생성
- 문서 업데이트: 문서가 있을 때 기존 문서 읽기, 변경 사항 분석, 재생성할 파일 사용자 확인

---

## Phase 2: 사용자 확인

AskUserQuestion으로 분석 요약 제시.

사용자의 conversation_language로 표시:

- 감지된 언어
- 프레임워크
- 아키텍처
- 주요 기능 목록

옵션:

- 문서 생성으로 진행
- 먼저 특정 분석 세부사항 검토
- 취소 후 프로젝트 설정 조정

"세부사항 검토" 시: 상세 내역 제공, 수정 허용.
"진행" 시: Phase 3으로 계속.
"취소" 시: 안내와 함께 종료.

---

## Phase 3: 문서 생성

[HARD] 문서 생성을 manager-docs 서브에이전트에 위임한다.

manager-docs에 전달:

- Phase 1 분석 결과 (또는 Phase 0.5의 사용자 입력)
- Phase 2 사용자 확인
- 출력 디렉토리: .do/project/
- 언어: 설정의 conversation_language

출력 파일:

- product.md: 프로젝트명, 설명, 대상 독자, 핵심 기능, 사용 사례
- structure.md: 디렉토리 트리, 각 디렉토리 목적, 주요 파일 위치, 모듈 구성
- tech.md: 기술 스택 개요, 프레임워크 선택 근거, 개발 환경 요구사항, 빌드 및 배포 설정

---

## Phase 3.5: 개발 환경 확인

목표: 감지된 기술 스택에 LSP 서버가 설치되어 있는지 확인.

언어-LSP 매핑 (16개 언어):

- Python: pyright 또는 pylsp (확인: which pyright)
- TypeScript/JavaScript: typescript-language-server (확인: which typescript-language-server)
- Go: gopls (확인: which gopls)
- Rust: rust-analyzer (확인: which rust-analyzer)
- Java: jdtls (Eclipse JDT Language Server)
- Ruby: solargraph (확인: which solargraph)
- PHP: intelephense (npm으로 확인)
- C/C++: clangd (확인: which clangd)
- Kotlin: kotlin-language-server
- Scala: metals
- Swift: sourcekit-lsp
- Elixir: elixir-ls
- Dart/Flutter: dart language-server (Dart SDK에 포함)
- C#: OmniSharp 또는 csharp-ls
- R: languageserver (R 패키지)
- Lua: lua-language-server

LSP 서버가 설치되지 않은 경우 AskUserQuestion 제시:

- LSP 없이 계속: 완료로 진행
- 설치 지침 표시: 감지된 언어의 설정 가이드 표시
- 지금 자동 설치: expert-devops 서브에이전트 사용 (확인 필요)

---

## Phase 4: 완료

사용자의 conversation_language로 완료 메시지 표시:

- 생성된 파일: 생성된 파일 목록
- 위치: .do/project/
- 상태: 성공 또는 부분 완료

다음 단계 (AskUserQuestion):

- SPEC 작성: /do plan 실행하여 기능 명세 정의
- 문서 검토: 생성된 파일 열어서 확인
- 새 세션 시작: 컨텍스트 초기화 후 새로 시작

---

## 에이전트 체인 요약

- Phase 0-2: Do 오케스트레이터 (모든 사용자 상호작용을 위한 AskUserQuestion)
- Phase 1: Explore 서브에이전트 (코드베이스 분석)
- Phase 3: manager-docs 서브에이전트 (문서 생성)
- Phase 3.5: expert-devops 서브에이전트 (선택적 LSP 설치)

---

Version: 2.0.0
Last Updated: 2026-02-07
