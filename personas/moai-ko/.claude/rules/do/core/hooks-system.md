---
paths:
  - "**/.claude/hooks/**"
  - "**/.claude/settings.json"
  - "**/.claude/settings.local.json"
---

# 훅 시스템

사용자 정의 스크립트로 기능을 확장하기 위한 Claude Code 훅입니다.

## 훅 이벤트

사용 가능한 14개 훅 이벤트 유형:

| Event | Matcher | 차단 가능 | 설명 |
|-------|---------|-----------|-------------|
| UserPromptSubmit | No | Yes | 사용자가 프롬프트를 제출할 때 실행, 처리 전 |
| SessionStart | No | No | 새 세션이 시작될 때 실행 |
| PreCompact | No | No | 컨텍스트 압축 전에 실행 |
| PreToolUse | Tool name | Yes | 도구가 실행되기 전에 실행 |
| PostToolUse | Tool name | No | 도구가 성공적으로 완료된 후에 실행 |
| PostToolUseFailure | Tool name | No | 도구 실행 실패 후에 실행 |
| PermissionRequest | Tool name | Yes | 권한 대화상자가 나타날 때 실행 |
| Notification | Type | No | Claude Code가 알림을 보낼 때 실행 |
| SubagentStart | Agent type | No | 서브에이전트가 생성될 때 실행 |
| SubagentStop | No | No | 서브에이전트가 종료될 때 실행 |
| Stop | No | No | 대화가 중단될 때 실행 |
| TeammateIdle | No | Yes | 에이전트 팀 팀원이 유휴 상태가 되려 할 때 실행 |
| TaskCompleted | No | Yes | 작업이 완료로 표시되려 할 때 실행 |
| SessionEnd | Reason | No | 세션이 종료될 때 실행 |

### 이벤트 카테고리

**Lifecycle Events**: SessionStart, SessionEnd, Stop, PreCompact

**Prompt Events**: UserPromptSubmit, PermissionRequest, Notification

**Tool Events**: PreToolUse, PostToolUse, PostToolUseFailure

**Agent Events**: SubagentStart, SubagentStop, TeammateIdle, TaskCompleted

## 훅 이벤트 stdin/stdout 참조

| Event | stdin | stdout | 참고 |
|-------|-------|--------|-------|
| UserPromptSubmit | `prompt` | `additionalContext`, `reason` | Exit 2는 프롬프트 차단 |
| PermissionRequest | `toolName`, `toolInput` | `reason` | Exit 0 = 허용, exit 2 = 거부 |
| PostToolUseFailure | `toolName`, `toolInput`, `error`, `is_interrupt` | `systemMessage` | 차단 없음 |
| Notification | `type`, `message` | - | 유형: permission_prompt, idle_prompt, auth_success, elicitation_dialog |
| SubagentStart | `agentType`, `agentName` | `additionalContext` | 서브에이전트에 컨텍스트 주입 |
| TeammateIdle | `agentType`, `agentName`, `tasksSummary` | `systemMessage` | Exit 2 = 계속 작업. 팀 품질에 중요 |
| TaskCompleted | `taskId`, `taskSummary`, `agentName` | `reason` | Exit 2 = 완료 거부. 팀 품질에 중요 |
| SessionEnd | `reason`, `sessionId` | - | 이유: clear, logout, prompt_input_exit, bypass_permissions_disabled, other |

표준 이벤트 (SessionStart, PreCompact, PreToolUse, PostToolUse, Stop)는 일반적인 stdin/stdout 패턴을 사용합니다: stdin은 이벤트별 필드를 수신하고, stdout은 선택적 `systemMessage`를 수락합니다.

## Hook Execution Types

### Command Hooks (type: "command")

Default hook type. Executes a shell command, communicates via stdin/stdout JSON.

- Configuration: `type`, `command`, `timeout`
- stdin: JSON with event data
- stdout: JSON with response (optional `systemMessage`, `additionalContext`, `reason`)
- Exit codes: 0 = success, 1 = error (shown to user), 2 = block/reject (for blocking events)

### Prompt Hooks (type: "prompt")

Send hook input to an LLM for single-turn evaluation. The LLM receives the event data and returns a judgment.

- Configuration: `type`, `prompt`, `model`, `timeout`
- The `prompt` field contains instructions for the LLM evaluator
- Returns JSON: `ok` (boolean), `reason` (string explanation)
- When `ok` is false on a blocking event, the operation is blocked with the provided reason

### Agent Hooks (type: "agent")

Spawn a subagent with tool access to verify conditions. The agent can read files, search code, and make informed decisions.

- Configuration: `type`, `prompt`, `model`, `timeout`
- Agent has access to: Read, Grep, Glob
- Returns JSON: `ok` (boolean), `reason` (string explanation)
- Same blocking behavior as prompt hooks

### Async Command Hooks (async: true)

Run command hooks in the background without blocking the conversation.

- Only available for `type: "command"` hooks
- Configuration: Add `async: true` to any command hook definition
- Results are delivered on the next conversation turn via `systemMessage`
- Useful for long-running validations (linting, test execution, deployments)

## Agent-Specific Hooks

Agent hooks are defined in agent frontmatter and executed for agent lifecycle events. For detailed configuration, actions table, and handler architecture, see @agent-hooks.md.

## Hook Location

Hooks are defined in `.claude/hooks/` directory:

- Shell scripts: `*.sh`
- Python scripts: `*.py`

## Configuration

Define hooks in `.claude/settings.json`:

```json
{
  "hooks": {
    "SessionStart": [{
      "type": "command",
      "command": "\"$CLAUDE_PROJECT_DIR/.claude/hooks/moai/handle-session-start.sh\"",
      "timeout": 5
    }],
    "PreCompact": [{
      "command": "\"$CLAUDE_PROJECT_DIR/.claude/hooks/moai/handle-compact.sh\"",
      "timeout": 5
    }],
    "PreToolUse": [{
      "matcher": "Write|Edit|Bash",
      "command": "\"$CLAUDE_PROJECT_DIR/.claude/hooks/moai/handle-pre-tool.sh\"",
      "timeout": 5
    }],
    "PostToolUse": [{
      "matcher": "Write|Edit",
      "command": "\"$CLAUDE_PROJECT_DIR/.claude/hooks/moai/handle-post-tool.sh\"",
      "timeout": 60
    }],
    "Stop": [{
      "command": "\"$CLAUDE_PROJECT_DIR/.claude/hooks/moai/handle-stop.sh\"",
      "timeout": 5
    }]
  }
}
```

## Path Syntax Rules

Hooks support `$CLAUDE_PROJECT_DIR` and `$HOME` environment variables:

```json
{
  "command": "\"$CLAUDE_PROJECT_DIR/.claude/hooks/moai/hook.sh\""
}
```

**Important**: Quote the entire path to handle project folders with spaces:
- Correct: `"\"$CLAUDE_PROJECT_DIR/.claude/hooks/moai/hook.sh\""`
- Wrong: `"$CLAUDE_PROJECT_DIR/.claude/hooks/moai/hook.sh"`

For StatusLine path configuration, see @settings-management.md (StatusLine does NOT support environment variables).

## Hook Wrappers

MoAI-ADK generates hook wrapper scripts during `moai init` that:

1. Read stdin JSON from Claude Code
2. Forward it to the moai binary via `moai hook <event>` command
3. Support multiple moai binary locations:
   - `moai` command in PATH
   - Detected Go bin path from initialization
   - Default `~/go/bin/moai`

Wrapper scripts are located at:
- `.claude/hooks/moai/handle-session-start.sh`
- `.claude/hooks/moai/handle-compact.sh`
- `.claude/hooks/moai/handle-pre-tool.sh`
- `.claude/hooks/moai/handle-post-tool.sh`
- `.claude/hooks/moai/handle-stop.sh`

## Rules

- Hook feedback is treated as user input
- When blocked, suggest alternatives
- Avoid infinite loops (no recursive tool calls)
- Keep hooks lightweight for performance
- Use proper path quoting to handle spaces in project paths
- Prompt and agent hooks return JSON with `ok` and `reason` fields
- Async hooks deliver results via `systemMessage` on the next turn
- Exit code 2 is the universal "block/reject" signal for blocking events

## Error Handling

- Failed hooks should exit with non-zero code
- Error messages are displayed to user
- Hooks can block operations by returning error
- Missing hooks exit silently (Claude Code handles gracefully)
- Prompt/agent hooks that fail return `ok: false` with a reason

## Security

- Hooks run in sandbox by default
- Validate all hook inputs
- Do not store secrets in hook scripts
- Agent hooks (type: "agent") have read-only tool access (Read, Grep, Glob)

## MoAI Integration

- Skill("do-foundation-claude") for detailed patterns
- Hook scripts must follow coding-standards.md
- Hook wrappers are managed by `internal/hook/` package
- TeammateIdle and TaskCompleted hooks are critical for Agent Teams quality enforcement
