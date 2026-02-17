# 에이전트 작성

MoAI-ADK에서 사용자 정의 에이전트를 생성하기 위한 지침입니다.

## 에이전트 정의 위치

사용자 정의 에이전트는 `.claude/agents/*.md` 또는 `.claude/agents/**/*.md` (하위 디렉토리 지원)에 정의됩니다.

디렉토리 규칙:
- 사용자 정의 에이전트: `.claude/agents/<agent-name>.md` (루트 레벨)
- MoAI-ADK 시스템 에이전트: `.claude/agents/moai/<agent-name>.md` (moai 하위 디렉토리)

## 지원되는 Frontmatter 필드

모든 에이전트 정의는 YAML frontmatter를 사용합니다. 다음 필드를 사용할 수 있습니다:

| 필드 | 필수 | 기본값 | 설명 |
|-------|-------|----------|------|
| name | 예 | - | 고유 식별자, 소문자와 하이픈 |
| description | 예 | - | Claude가 이 에이전트에게 위임할 때 |
| tools | 아니오 | 전체 상속 | 에이전트가 사용할 수 있는 도구 (허용목록 접근) |
| disallowedTools | 아니오 | 없음 | 거부할 도구 (거부목록 접근, tools의 대안) |
| model | 아니오 | inherit | 모델 선택: sonnet, opus, haiku, 또는 inherit |
| permissionMode | 아니오 | default | 에이전트의 권한 동작 |
| maxTurns | 아니오 | 무제한 | 중지 전 최대 에이전트 턴 수 |
| skills | 아니오 | 없음 | 시작 시 에이전트 컨텍스트에 주입되는 스킬 |
| mcpServers | 아니오 | 없음 | 이 에이전트가 사용할 수 있는 MCP 서버 |
| hooks | 아니오 | 없음 | 이 에이전트로 범위 지정된 수명 주기 훅 |
| memory | 아니오 | 없음 | 교차 세션 학습을 위한 영구 메모리 범위 |

### 필드 상세

**tools**: 지정하면 에이전트는 나열된 도구만 사용할 수 있습니다. 생략하면 에이전트는 부모로부터 모든 도구를 상속받습니다. disallowedTools와 상호 배타적입니다.

**disallowedTools**: 거부목록 접근. 에이전트는 나열된 도구를 제외한 모든 도구를 상속받습니다. tools와 상호 배타적입니다.

**skills**: 전체 스킬 내용이 에이전트 컨텍스트에 주입되며, 호출을 위해 제공되는 것이 아닙니다. 에이전트는 부모로부터 스킬을 상속받지 않습니다. 나열된 각 스킬은 `.claude/skills/`에 존재해야 합니다.

**mcpServers**: `.mcp.json`의 키와 일치하는 서버 이름 참조 또는 명령과 인자가 있는 인라인 서버 정의입니다.

**hooks**: 이 에이전트로 범위 지정된 PreToolUse, PostToolUse, SubagentStop 이벤트를 지원합니다. 구성 형식은 @hooks-system.md를 참조하세요.

## Task(agent_type) 제한

`tools` 필드는 에이전트가 생성할 수 있는 서브에이전트를 제한하는 `Task(worker, researcher)` 구문을 지원합니다.

- `claude --agent`를 통해 메인 스레드로 실행되는 에이전트에만 적용됩니다
- 서브에이전트 정의에는 영향이 없습니다 (서브에이전트는 다른 서브에이전트를 생성할 수 없습니다)
- MoAI 에이전트는 서브에이전트로 실행되므로 이 제한은 현재 적용되지 않습니다
- 메인 스레드로 실행되는 코디네이터 에이전트 생성에 유용합니다

## 권한 모드

`permissionMode` 필드는 에이전트가 권한 확인을 처리하는 방식을 제어합니다:

| 모드 | 동작 | 사용 사례 |
|------|----------|----------|
| default | 사용자 프롬프트와 함께 표준 권한 확인 | 일반용 에이전트 |
| acceptEdits | 파일 편집 작업 자동 수락 | 신뢰할 수 있는 구현 에이전트 |
| delegate | 조정 전용 모드, 팀 관리 도구로 제한 | 팀 리드 에이전트 |
| dontAsk | 모든 권한 프롬프트 자동 거부 | 엄격한 샌드박스 에이전트 |
| bypassPermissions | 모든 권한 확인 건너뜀 (주의 사용) | 완전 신뢰할 수 있는 자동화 |
| plan | 읽기 전용 탐색 모드, 쓰기 작업 없음 | 연구 및 분석 에이전트 |

## 영구 메모리

`memory` 필드는 에이전트의 교차 세션 학습을 가능하게 합니다. 세 가지 범위 수준:

| 범위 | 저장 위치 | VCS로 공유 | 사용 사례 |
|-------|-----------------|----------------|----------|
| user | ~/.claude/agent-memory/\<name\>/ | 아니오 | 프로젝트 간 학습, 개인 선호도 |
| project | .claude/agent-memory/\<name\>/ | 예 | 프로젝트별 지식, 팀 공유 컨텍스트 |
| local | .claude/agent-memory-local/\<name\>/ | 아니오 | 프로젝트별 지식, 공유 안 됨 |

## 에이전트 카테고리

### 매니저 에이전트 (7)

워크플로우 및 다단계 프로세스를 조정합니다:

- manager-spec: SPEC document creation
- manager-ddd: DDD implementation cycle
- manager-tdd: TDD implementation cycle
- manager-docs: Documentation generation
- manager-quality: Quality gates validation
- manager-project: Project configuration
- manager-strategy: System design, architecture decisions
- manager-git: Git operations, branching strategy

### Expert Agents (8)

Domain-specific implementation:

- expert-backend: API and server development
- expert-frontend: UI and client development
- expert-security: Security analysis
- expert-devops: CI/CD and infrastructure
- expert-performance: Performance optimization
- expert-debug: Debugging and troubleshooting
- expert-testing: Test creation and strategy
- expert-refactoring: Code refactoring

### Builder Agents (3)

Create new MoAI components:

- builder-agent: New agent definitions
- builder-skill: New skill creation
- builder-plugin: Plugin creation

### Team Agents (8) - Experimental

Agents for Claude Code Agent Teams (v2.1.32+, requires CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1):

| Agent | Model | Phase | Mode | Purpose |
|-------|-------|-------|------|---------|
| team-researcher | haiku | plan | plan (read-only) | Codebase exploration and research |
| team-analyst | inherit | plan | plan (read-only) | Requirements analysis |
| team-architect | inherit | plan | plan (read-only) | Technical design |
| team-backend-dev | inherit | run | acceptEdits | Server-side implementation |
| team-designer | inherit | run | acceptEdits | UI/UX design with Pencil/Figma MCP |
| team-frontend-dev | inherit | run | acceptEdits | Client-side implementation |
| team-tester | inherit | run | acceptEdits | Test creation with exclusive test file ownership |
| team-quality | inherit | run | plan (read-only) | TRUST 5 quality validation |

## Rules

- Write agent definitions in English
- Define expertise domain clearly in description
- Minimize tool permissions (least privilege)
- Include relevant trigger keywords
- Use permissionMode: plan for read-only agents
- Preload skills for domain expertise instead of relying on runtime loading

## Tool Permissions

Recommended tool sets by category:

Manager agents: Read, Write, Edit, Grep, Glob, Bash, Task, TaskCreate, TaskUpdate

Expert agents: Read, Write, Edit, Grep, Glob, Bash

Builder agents: Read, Write, Edit, Grep, Glob

Team implementation agents: Read, Write, Edit, Grep, Glob, Bash (+ skills preloading for domain expertise)

Team research agents: Read, Grep, Glob, Bash (read-only via permissionMode: plan)

Notes:
- Use `skills` field to preload domain-specific knowledge into team agents
- Team agents with permissionMode: plan cannot write files regardless of tools listed
- Prefer skills preloading over large tool lists for domain expertise

## Agent Invocation

Invoke agents via Task tool:

- "Use the expert-backend subagent to implement the API"
- Task tool with subagent_type parameter

## MoAI Integration

- Use builder-agent subagent for creation
- Skill("do-foundation-claude") for patterns
- Follow skill-authoring.md for YAML schema
- See @hooks-system.md for agent hook configuration
