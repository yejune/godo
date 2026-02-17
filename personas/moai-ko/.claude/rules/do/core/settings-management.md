---
paths:
  - "**/.moai/config/**"
  - "**/.mcp.json"
  - "**/.claude/settings.json"
  - "**/.claude/settings.local.json"
---

# 설정 관리

Claude Code 및 MoAI 구성 관리 규칙입니다.

## 구성 파일

### Claude Code 설정

`.claude/settings.json` - 프로젝트 수준 설정:

- allowedTools: 허용된 도구 목록
- hooks: 훅 스크립트 정의
- permissions: 액세스 제어
- statusLine: Statusline 구성

### MCP 구성

`.mcp.json` - MCP 서버 정의:

- mcpServers: 서버 명령 및 인자
- 서버용 환경 변수

### MoAI 구성

`.moai/config/` - MoAI 전용 설정:

- config.yaml: 메인 구성
- sections/quality.yaml: 품질 게이트, 커버리지 목표
- sections/language.yaml: 언어 선호도
- sections/user.yaml: 사용자 정보

## 훅 구성

훅은 환경 변수를 지원하며 공백 처리를 위해 인용부호로 묶어야 합니다:

```json
{
  "hooks": {
    "SessionStart": [{
      "type": "command",
      "command": "\"$CLAUDE_PROJECT_DIR/.claude/hooks/moai/handle-session-start.sh\"",
      "timeout": 5
    }],
    "PreToolUse": [{
      "matcher": "Write|Edit|Bash",
      "command": "\"$CLAUDE_PROJECT_DIR/.claude/hooks/moai/handle-pre-tool.sh\"",
      "timeout": 5
    }]
  }
}
```

**중요**: 전체 경로를 인용부호로 묶으세요: `"\"$CLAUDE_PROJECT_DIR/path\""`가 아니라 `"$CLAUDE_PROJECT_DIR/path"`

## StatusLine 구성

StatusLine은 환경 변수를 지원하지 않습니다. 프로젝트 루트의 상대 경로를 사용하세요:

```json
{
  "statusLine": {
    "type": "command",
    "command": ".moai/status_line.sh"
  }
}
```

참조: GitHub Issue #7925 - statusline은 환경 변수를 확장하지 않습니다.

## 권한 관리

settings.json의 도구 권한:

- Read, Write, Edit: 파일 작업
- Bash: 셸 명령 실행
- Task: 에이전트 위임
- AskUserQuestion: 사용자 상호작용

## 품질 구성

quality.yaml의 품질 게이트:

- development_mode: ddd, tdd, 또는 hybrid
- test_coverage_target: 최소 커버리지 비율
- lsp_quality_gates: LSP 기반 검증

## 언어 설정

language.yaml의 언어 선호도:

- conversation_language: 사용자 응답 언어
- agent_prompt_language: 내부 통신
- code_comments: Code comment language

## Agent Teams Settings

Agent Teams require both an environment variable and workflow configuration.

### Environment Variable

Enable in `.claude/settings.json`:

```json
{
  "env": {
    "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1"
  }
}
```

This env var must be set for Claude Code to expose the Teams API.

### Workflow Configuration

Team behavior is controlled by the `workflow.team` section in `.moai/config/sections/workflow.yaml`:

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| team.enabled | boolean | false | Master switch for team mode |
| team.max_teammates | integer | 10 | Maximum teammates per team (2-10 recommended) |
| team.default_model | string | inherit | Default model for teammates (inherit/haiku/sonnet/opus) |
| team.require_plan_approval | boolean | true | Require plan approval before implementing |
| team.delegate_mode | boolean | true | Team lead coordination-only mode (no direct implementation) |
| team.teammate_display | string | auto | Display mode: auto, in-process, or tmux |

### Auto-Selection Thresholds

When `workflow.execution_mode` is `auto`, these thresholds determine when team mode activates:

| Setting | Default | Description |
|---------|---------|-------------|
| team.auto_selection.min_domains_for_team | 3 | Minimum distinct domains to trigger team mode |
| team.auto_selection.min_files_for_team | 10 | Minimum affected files to trigger team mode |
| team.auto_selection.min_complexity_score | 7 | Minimum complexity score (1-10) to trigger team mode |

## Rules

- Never commit secrets to settings files
- Use environment variables for sensitive data
- Keep settings minimal and focused
- Hook paths must be quoted when using environment variables
- StatusLine uses relative paths only (no env var expansion)
- Template sources (.tmpl files) belong in `internal/template/templates/` only
- Local projects should contain rendered results, not template sources

## MoAI Integration

- Skill("do-workflow-project") for project setup
- Skill("do-foundation-core") for quality framework
- See hooks-system.md for detailed hook configuration patterns
