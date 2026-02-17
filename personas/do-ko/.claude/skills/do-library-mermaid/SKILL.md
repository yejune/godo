---
name: do-library-mermaid
description: >
  MCP Playwright와 함께 Claude Code용 Enterprise Mermaid 다이어그램밍 스킬입니다.
  아키텍처 다이어그램, 플우차트, 시퀀스 다이어그램,
  시각화 문서 생성 시 사용하세요.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Grep Glob Bash(npx:*) Bash(mmdc:*) mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "7.1.0"
  category: "library"
  modularized: "true"
  status: "active"
  updated: "2026-01-11"
  tags: "library, mermaid, diagrams, flowchart, sequence, visualization, documentation"

# MoAI Extension: Triggers
triggers:
  keywords: ["diagram", "flowchart", "sequence", "architecture", "mermaid", "visualization", "chart", "graph"]
---

## Quick Reference

Mermaid Diagram Expert v7.1.0 - MCP Playwright 통합이 있는 pure skill-based 렌더링을 위한 Claude Code입니다.

이 스킬은 21가지 다이어그램 타입에 대한 완전한 Mermaid 11.12.2 구문, SVG 및 PNG 출력을 위한 MCP Playwright 통합, 사용 준비 예제 및 참조 문서, 엔터프라이즈 다이어그램을 위한 best practices를 제공합니다.

이 스킬을 호출하려면 skill 이름 do-library-mermaid를 사용하여 표준 skill 호출 패턴을 사용하세요.

### 지원되는 다이어그램 타입

구조조 다이어그램은 Flowchart(process flows, decision trees), Sequence(interaction sequences, message flows), Class(object-oriented class relationships), ER(entity-relationship databases), Block(block diagram structures), State(state machines and stateful flows)를 포함합니다.

Timeline 및 추적 다이어그램은 Timeline(chronological events, milestones), Gantt(project scheduling, timelines), Gitgraph(Git workflow, branching visualization)를 포함합니다.

아키텍처 및 디자인 다이어그램은 C4(Context, Container, Component, Code architecture diagrams), Architecture(system architecture diagrams), Requirement(requirements and traceability documentation)를 포함합니다.

데이터 시각화 다이그램은 Pie Chart(pie & donut charts), XY Chart(scatter and line charts), Sankey(flow diagrams with proportional width), Radar(multi-variable comparison charts)를 포함합니다.

User 및 프로세스 다이어그램은 Mindmap(hierarchical mind mapping), User Journey(user experience flows), Kanban(board state visualization), Packet(network packet structures)를 포함합니다.

### MCP Playwright 통합

이 스킬은 다이어그램 렌더링을 위한 MCP Playwright와 통합됩니다. 프로젝트 mcp.json 파일의 MCP 구성에서 Playwright 서버를 활성화해야 합니다.

다이어그램을 렌더링하려면 프로젝트 MCP 설정에서 Playwright 서버가 구성되어야 하고, 시스템에 Node.js가 설치되어야 하며, npx를 통해 Playwright를 사용할 수 있어야 합니다.

---

## Implementation Guide

### Diagram Syntax Patterns

Flowchart 다이어그램:

flowchart 키워드와 함께 방향 지시자(TD for top-down, LR for left-right)를 사용하세요. 둥근 모양, 각형, 중간 다이아몬드, 스테이득 형태를 위한 괄호, 더블 대시 더블 각형, 스타디웃 형태를 사용하여 노드를 정의하세요. 화살표와 함께 이중 대시, 각�진 화살표와 함께 단일 화살표, 대시 더블 더블 각형을 사용하여 연결합니다. 접두사이 pipe로 연결되고 콜마(,)와 함께 경로 설명을 추가하여 edge description을 제공합니다. subgraph 키워드와 제목과 종료 구분자를 사용하여 관련 노드를 그룹화합니다.

Sequence 다이어그램:

participant 키워드로 참여자를 먼저 정의하세요. 화살표와 쌍을 포함한 arrow 표기법을 사용하여 상호작용을 표시합니다. 쌍 화살표는 동기 호출을, 대시 화살표는 응답 또는 비동기 메시지를 나타냅니다. activate 및 deactivate 키워드 또는 +/- 바로 가장 shorthand를 사용하여 활성 기간을 나타냅니다. 참고 right, reference left, over, note를 사용하여 참여자에 대해 메모, left over, over를 추가할 수 있습니다.

C4 Context 다이그램:

C4Context 키워드를 사용하세요. Enterprise_Boundary 또는 System_Boundary 함수를 사용하여 시스템 경계를 정의합니다. Person(id, name, optional description) 함수를 사용하여 사람을 정의합니다. System 내부 시스템을 위한 System을, 외부 시스템을 위한 System_Ext를 정의합니다. Rel 함수로 source, target, description, optional technology를 연결합니다.

Class 다이어그램:

classDiagram 키워드를 사용하세요. 클래스에서 attributes와 methods를 정의하세요. 상속을 위해 화살표와 파이프 문자, 집합을 위해 asterisk, 연관을 위해 circle, 연관을 위해 대시를 사용합니다. 공용 멤버는 plus, private 멤버는 minus, protected 멤버는 hash로 표시합니다.

State 다이어그램:

stateDiagram-v2 키워드를 사용하여 최신 구문을 사용하세요. state 키워드와 대괄형의 state를 정의합니다. 화살표를 사용하여 전이를 정의하고 arrow 표기와 함께 optional label을 추가합니다. 상태 블록 내에 중첩 상태를 포함합니다.

### 다이어그램 카테고리

Process and Flow 다이어그램: Flowchart, Sequence, State, Timeline, Gitgraph, User Journey 다이어그램 타입을 포함합니다. 이는 동적 프로세스와 시간적 시퀀스를 나타냅니다.

Structure and Design 다이어그램: Class, ER, Block, Architecture, C4 타입을 포함합니다. 이는 정적 구조와 시스템 구성을 나타냅니다.

Data and Analytics 다이그램: Pie Chart, XY Chart, Sankey, Radar 타입을 포함합니다. 이는 정량 데이터와 비교 메트릭스를 시각화합니다.

Planning and Organization 다이그램: Gantt, Mindmap, Kanban, Requirement 타입을 포함합니다. 이는 프로젝트 관리 및 요구사항 추적을 지원합니다.

Network and Technical 다이그램: 현재 Packet 타입이 포함되며 추가 확장이 예약되어 있습니다.

### Best Practices

명확성과 가독성: 모든 노드에 설명적 레이블을 사용하고 20-30개 노드로 복잡성을 유지하며 전체적으로 일관된 스타일과 color scheme를 사용하세요.

성능: 복잡한 다이어그램을 여러 작은 다이어그램으로 분할하고, 큰 flowchart를 subgraph로 구성하고, 렌더링 성능을 유지하기 위해 노드 내 텍스트 길이를 제한하세요.

접근성: 모든 다이어그램에 텍스트 대안을 제공하고, color만이 아닌 색상/명암 차별을 사용하고, 문맥에 제목과 legend를 포함하세요.

조직화: related 다이어그램을 디렉토리로 그룹화하고, 일관된 명명 규칙을 사용하며, 소스 파일 내 주석으로 다이어그램 목적을 문서화하세요.

---

## Advanced Patterns

### MoAI-ADK와의 통합

이 스킬은 Claude Code의 다양한 개발 단계에서 사용하기 위해 설계되었습니다:

/moai:1-plan 명령을 사용하는 아키텍처 단계에서 시스템 설계 다이어그램을 생성하여 제안된 솔루션과 컴포넌트 관계를 시각화합니다.

/moai:3-sync 명령을 사용하는 문서화 단계에서 flowcharts, sequence diagrams, architecture overviews를 포함한 시각 문서를 생성합니다.

코드 리뷰 단계에서 시스템 설계를 시각적으로 전달하고 우려 사항을 강조하기 위해 다이어그램을 사용합니다.

온보딩 프로세스에서 새 팀원이 아키텍처와 이해를 돕습니다.

### 일반 아키텍처 패턴

API 아키텍처: C4 다이어그램을 사용하여 API gateway, backend services, database layer, cache layer 관계를 보여줍니다.

Microservices Flow: sequence 다이어그램을 사용하여 클라이언트 요청이 API gateway를 통해 개별 서비스와 그 data store로 흐르는 것을 나타냅니다.

Data Pipeline: flowchart 다이그램을 사용하여 extract, transform, load, validate, report 단계를 통해 데이터 이동을 나타냅니다.

### Context7 통합

최신 Mermaid 문서의 경우 Context7 library resolution 및 documentation tools를 사용하세요.

2025년 12월 기준 현재 안정 버전은 Mermaid 11.12.2입니다.

공식 문서는 mermaid.js.org/intro에 일반 문서가 있고 mermaid.js.org/config/setup/modules/mermaidAPI.html에 API 참조가 있습니다.

릴리스 노트와 마이그레이션 가이드는 Mermaid GitHub repository releases 섹션에 있습니다.

### Learning Resources

공식 Mermaid 사이트는 mermaid.js.org에 있습니다. 테스트용 대화형 live editor는 mermaid.live에 있습니다. 전체 구문 가이드는 mermaid.js.org/syntax/에 있습니다.

모든 21가지 다이어그램 타입에 대한 작동 예제는 이 스킬 디렉토리의 examples.md 파일을 참조하세요. 확장 참조 문서는 reference.md를 참조하세요. 최적화 기법은 optimization.md를 참조하세요. 복잡 다이어그램 패턴은 advanced-patterns.md를 참조하세요.

---

## Works Well With

에이 스킬은 여러 agent 및 기타 스킬과 통합합니다:

workflow-docs documentation with diagrams, workflow-spec for SPEC diagrams and requirements visualization, design-uiux for architecture visualization and interface documentation와 잘 작동합니다.

Skills: do-docs-generation comprehensive documentation generation, do-workflow-docs for diagram validation and documentation workflows, do-library-nextra for architecture documentation sites를 포함합니다.

Commands: moai:3-sync for documentation with embedded diagrams, moai:1-plan for SPEC creation with visual architecture diagrams를 포함합니다.

Focus: MCP Playwright 통합이 있는 pure skill-based Mermaid 렌더링
