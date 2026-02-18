# godo Hook System Reference

Complete reference for `internal/hook/`. Use this to verify correct implementation
without comparing source code each time.

Last updated: 2026-02-18 | Version: v0.3.14

---

## Architecture Overview

```
Claude Code
    │  (stdin: JSON)
    ▼
godo hook <event-type>         ← internal/cli/hook.go
    │
    ├─ ReadInput()              ← reads stdin → *hook.Input
    ├─ Contract.Validate()      ← checks ctx + workdir
    ├─ handler(*Input)          ← dispatches to specific handler
    └─ WriteOutput(*Output)     ← writes JSON to stdout
```

Exit codes: **0** = allow/success, **2** = block, other = non-blocking error.

---

## settings.json Hook Configuration

File: `.claude/settings.json`

| Event | CLI command | tool matcher |
|---|---|---|
| SessionStart | `godo hook session-start` | (all sessions) |
| PreToolUse | `godo hook pre-tool` | `Write\|Edit\|Bash` |
| PostToolUse | `godo hook post-tool-use` | `.*` |
| PostToolUse (compact) | `godo hook compact` | (none = all tools) |
| SessionEnd | `godo hook session-end` | (all sessions) |
| Stop | `godo hook stop` | (always) |
| SubagentStop | `godo hook subagent-stop` | (always) |
| UserPromptSubmit | `godo hook user-prompt-submit` | (every prompt) |

---

## Protocol: Input Fields (`types.go`)

```go
type Input struct {
    // Common (all events)
    SessionID      string
    TranscriptPath string
    CWD            string          // project root = $CLAUDE_PROJECT_DIR
    PermissionMode string          // "default"|"acceptEdits"|"bypassPermissions"|etc.
    HookEventName  string

    // Tool events (PreToolUse, PostToolUse)
    ToolName     string            // "Write"|"Edit"|"Bash"|"Task"|...
    ToolInput    json.RawMessage   // tool call arguments as JSON object
    ToolOutput   json.RawMessage
    ToolUseID    string

    // SessionStart only
    Source    string               // "human"|"api"
    Model     string               // active model ID
    AgentType string               // "main"|subagent type

    // SessionEnd only
    Reason string

    // Stop / SubagentStop
    StopHookActive bool            // true = stop hook already running (loop guard)

    // SubagentStop only
    AgentID             string
    AgentTranscriptPath string

    // PreCompact only
    Trigger            string      // "auto"|"manual"
    CustomInstructions string

    // PostToolUse failure
    Error       string
    IsInterrupt bool

    // UserPromptSubmit only
    Prompt string
}
```

### ToolInput JSON shapes

```json
// Write / Edit
{ "file_path": "/absolute/path/to/file", "content": "..." }

// Bash
{ "command": "go build ./..." }

// Task
{ "prompt": "use expert-backend subagent, read .do/jobs/26/02/18/title/checklists/01_backend.md" }
```

---

## Protocol: Output Fields (`types.go`)

```go
type Output struct {
    // Universal
    Continue       bool    // continue session (SessionStart/SessionEnd/SubagentStop)
    SystemMessage  string  // injected into Claude system context (SessionStart)
    SuppressOutput bool

    // Top-level block (Stop / PostToolUse)
    Decision string         // "block"
    Reason   string

    // Nested (PreToolUse / PostToolUse)
    HookSpecificOutput *SpecificOutput
}

type SpecificOutput struct {
    HookEventName            string  // "PreToolUse"|"PostToolUse"
    PermissionDecision       string  // "allow"|"deny"|"ask"  (PreToolUse only)
    PermissionDecisionReason string
    AdditionalContext        string  // injected into Claude context
}
```

### Output constructors

| Constructor | When to use |
|---|---|
| `NewAllowOutput()` | PreToolUse: allow |
| `NewDenyOutput(reason)` | PreToolUse: deny with reason shown to user |
| `NewAskOutput(reason)` | PreToolUse: ask user for confirmation |
| `NewAllowOutputWithWarning(msg)` | PreToolUse: allow + inject warning into context |
| `NewSessionOutput(continue, msg)` | SessionStart / SessionEnd |
| `NewPostToolOutput(ctx)` | PostToolUse: inject additional context |
| `NewStopBlockOutput(reason)` | Stop: prevent Claude from stopping |
| `NewPostToolBlockOutput(reason, ctx)` | PostToolUse: block after tool execution |

---

## Contract (`contract.go`)

Every hook call validates the execution environment:

```go
contract := hook.NewContract(workDir)   // workDir = os.Getwd()
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
err := contract.Validate(ctx)
// Checks: context not done | WorkDir non-empty | WorkDir is accessible directory
```

**Claude Code guarantees:**
- stdin: valid JSON (hook protocol)
- CWD = project root (`$CLAUDE_PROJECT_DIR`)

**NOT guaranteed:**
- User PATH (binary must be in system PATH)
- Shell env vars (`.bashrc`/`.zshrc` not loaded)
- Python/Node.js runtime availability

---

## Hook Handlers (all 8)

### 1. SessionStart (`session_start.go`)

**CLI**: `godo hook session-start`
**Trigger**: Once per session start

**Input used**: `input.CWD`

**Behavior (in order)**:
1. Read `DO_PERSONA` env → default `"young-f"`
2. Read current mode: `mode.ReadState()`
3. Apply spinner verbs:
   ```
   persona.GetSpinnerVerbs(personaType) → persona.ApplySpinnerToSettings(verbs)
   ```
   - Writes `spinnerVerbs` to `CLAUDE_CONFIG_DIR/settings.json`
   - **BUG**: `GetClaudeSettingsPath()` ignores `CLAUDE_CONFIG_DIR` → always writes to `~/.claude/settings.json`
   - YAML source: `.claude/personas/do/spinners/<type>.yaml` (hardcoded fallback if not found)
4. Version check: `checkLatestVersion()`
   - Cache file: `.do/.latest-version` (24h TTL)
   - HTTP GET GitHub releases API (3s timeout)
   - **BUG**: uses `!=` not `IsNewer()` → stale cache causes false "update available"
5. Detect project from `cwd`:
   - `go.mod` → name (module last path segment), type=go, lang=go
   - `package.json` + `tsconfig.json` → lang=typescript, type=node
   - `Cargo.toml` → name from `name =` field, lang=rust, type=rust
   - `pyproject.toml` or `setup.py` → lang=python, type=python
   - Env overrides: `DO_PROJECT_NAME`, `DO_PROJECT_TYPE`, `DO_PROJECT_LANG`
   - Fallback name: `filepath.Base(cwd)`
6. Build systemMessage: `"current_mode: do\nproject: godo, type: go, lang: go"`

**Output**: `NewSessionOutput(true, systemMsg)` → `{continue: true, systemMessage: "..."}`

---

### 2. SessionEnd (`session_end.go`)

**CLI**: `godo hook session-end`
**Behavior**: Stub. Allow session to end normally.
**Output**: `&Output{Continue: true}`

---

### 3. PreToolUse (`pre_tool.go`)

**CLI**: `godo hook pre-tool` | **Timeout**: 5s
**Trigger**: Before `Write`, `Edit`, `Bash` tool calls

**Decision flow**:

```
Phase A — Checklist enforcement (only if .do/jobs/ exists):
  Tool == Task only:
    → extractPrompt(input.ToolInput)  ← looks for "prompt" JSON field
    → if no prompt → DENY "체크리스트 없음: ..."
    → if prompt lacks ".do/jobs/" AND "/checklists/" → DENY
    → extract checklist path (must end in .md, contain /checklists/)
    → if file not on disk → DENY

Phase B — Security policy:
  Write|Edit|Read|Glob → checkFileAccess(policy, input)
    file_path from ToolInput ("file_path" or "path" key)
    → DenyFilePatterns match → DENY
    → AskFilePatterns match (Write/Edit only) → ASK
    → else → ALLOW

  Bash → checkBashCommand(policy, input)
    command from ToolInput "command" key
    → DenyBashPatterns match → DENY
    → AskBashPatterns match → ASK
    → SensitiveContentPatterns match → DENY
    → else → ALLOW

  other tools → ALLOW
```

**Security patterns** (defined in `security.go` — see source for full regex list):

DenyFilePatterns — files that must NEVER be modified:
- Secrets/credentials files
- SSH keys, TLS certs (`.pem`, `.key`, `.crt`)
- `.git/` internals
- Cloud credential dirs (`.aws/`, `.gcloud/`, `.azure/`, `.kube/`)
- Token files

AskFilePatterns — require user confirmation on Write/Edit:
- All lock files (`package-lock.json`, `yarn.lock`, `Cargo.lock`, `uv.lock`, etc.)
- Critical project configs (`tsconfig.json`, `pyproject.toml`, `Cargo.toml`, `package.json`)
- Docker files (`Dockerfile`, `docker-compose.yaml`, `.dockerignore`)
- CI/CD configs (`.github/workflows/`, `.gitlab-ci.yaml`, `Jenkinsfile`)
- Infrastructure (`terraform/`, `kubernetes/`, `k8s/`)

DenyBashPatterns — commands that must NEVER execute (see `DenyBashPatternStrings` in `security.go`):
- Cloud DB destruction (supabase/neon/railway/vercel project delete)
- SQL DDL destruction (drop/truncate operations)
- Destructive file system operations on root, home, or git dirs
- Windows CMD/PowerShell destructive equivalents
- Dangerous git operations (force push to main/master, delete main branch)
- Cloud infrastructure destruction (terraform/pulumi/gcloud/aws/az delete)
- Docker bulk pruning operations
- Classic exploits (fork bomb, disk format/overwrite)

AskBashPatterns — require confirmation (see `AskBashPatternStrings` in `security.go`):
- ORM schema push/reset (prisma, drizzle-kit)
- Destructive git ops (force push, hard reset, clean)
- Package manager cache purges

SensitiveContentPatterns — deny bash containing (see `SensitiveContentPatternStrings` in `security.go`):
- PEM private key / certificate headers
- Known API key formats: OpenAI, GitHub, GitLab, Slack, AWS, Google OAuth

---

### 4. PostToolUse (`post_tool_use.go`)

**CLI**: `godo hook post-tool-use`
**Behavior**: Stub. Non-intrusive baseline.
**Output**: `&Output{}`

---

### 5. Compact / PreCompact (`compact.go`)

**CLI**: `godo hook compact`
**Registered under**: PostToolUse (no matcher) → runs after EVERY tool call
**Behavior**: Stub. Non-intrusive.
**Output**: `&Output{Continue: true}`

> **Issue**: Listed under PostToolUse in settings.json, not PreCompact. Runs on every
> tool, not just on compact events. Add a separate PreCompact entry if needed.

---

### 6. Stop (`stop.go`)

**CLI**: `godo hook stop`
**Trigger**: When Claude is about to stop responding

**Behavior**:
1. `input.StopHookActive == true` → allow immediately (loop guard)
2. `checkActiveChecklist()`:
   - Walk `.do/jobs/{YY}/{MM}/{DD}/{task}/checklist.md` (most recent job first — sorted desc)
   - Parse item states: `[o]`/`[O]`=done, `[~]`=inProgress, `[!]`=blocked, `[ ]`=pending, `[*]`=testing
   - **No block** if: total=0, all done, or zero in-progress AND zero blocked
   - **Block** if in-progress or blocked items exist:
     ```
     "활성 체크리스트가 있습니다 (N/M 완료, N 진행중, N 대기, N 블로커).
     체크리스트 파일(path)을 읽고 현재 상태를 사용자에게 표시한 뒤 종료하세요."
     ```

**Output**:
- Allow: `&Output{}` (empty)
- Block: `{decision: "block", reason: "활성 체크리스트..."}`

---

### 7. SubagentStop (`subagent_stop.go`)

**CLI**: `godo hook subagent-stop`
**Behavior**: Stub. Allow subagent to stop.
**Output**: `&Output{Continue: true}`

---

### 8. UserPromptSubmit (`user_prompt_submit.go`)

**CLI**: `godo hook user-prompt-submit`
**Trigger**: Each user prompt turn (reinforces mode + persona every turn)

**Behavior**:
1. `mode.ReadState()` → current mode (do/focus/team)
2. `DO_USER_NAME`, `DO_PERSONA` env vars
3. Mode prefix: capitalize first letter → `"현재 실행 모드: do (응답 접두사: [Do])"`
4. Persona reminder: `persona.LoadCharacter(personaDir, personaType).BuildReminder(userName)`
   - `personaDir` = `persona.ResolveDir()`:
     - Uses `CLAUDE_PROJECT_DIR` env or `os.Getwd()`
     - Checks `{dir}/.claude/personas/do` → then `{dir}/personas/do`
5. Join with `"\n"` → `additionalContext`

**Output**:
```go
&Output{
    HookSpecificOutput: &SpecificOutput{
        HookEventName:    "UserPromptSubmit",
        AdditionalContext: "현재 실행 모드: do (응답 접두사: [Do])\n<persona-specific reminder>",
    },
}
```
This injects mode + persona context into EVERY user turn without appearing in conversation.

---

## Shared Utilities

### `checklist.go` — Checklist parsing

```go
type ChecklistStats struct {
    Total, Pending, InProgress, Testing, Blocked, Done, Failed int
}
stats.HasIncomplete() bool    // InProgress > 0 || Blocked > 0
stats.Summary() string        // "[o]3 [~]1 [ ]2"

ParseChecklistFile(path) (*ChecklistStats, error)
ParseChecklistContent(content string) *ChecklistStats
FindLatestChecklist() string  // most recent checklist.md under .do/jobs/
```

State symbols:
| Symbol | Meaning |
|---|---|
| `[ ]` | Pending |
| `[~]` | InProgress |
| `[*]` | Testing |
| `[!]` | Blocked |
| `[o]`/`[O]` | Done |
| `[x]` | Failed |

### `git.go` — Git status

```go
var GitStatus = func() (bool, string)
// hasChanges: false if not a git repo, clean, or only untracked (??) files
// summary: tracked changes only (excludes ?? lines)
```

### `job_state.go` — Workflow state

```go
type JobState struct {
    JobID, WorkflowType string              // "simple"|"complex"
    Phases map[string]PhaseState            // status: "pending"|"in_progress"|"complete"
    Agents map[string]AgentState            // status: "pending"|"in_progress"|"complete"|"failed"|"blocked"
    AutoResolveAttempts map[string]bool
}
LoadJobState(path) (*JobState, error)
SaveJobState(path, state) error
```

### `dispatcher.go` — I/O

```go
ReadInput() *Input                      // parse stdin → structured Input
WriteOutput(*Output)                    // marshal Output → stdout
GetStringField(data, key, fallback)     // safe string extraction from map
```

---

## Known Bugs

### Bug 1: Version check false positive (`session_start.go:45`)

```go
// Current: stale cache "0.2.53" != current "0.3.14" → false "downgrade available"
if latestVer != "" && latestVer != Version && Version != "dev" && Version != "" {

// Fix: replace with IsNewer(latestVer, Version) semver comparison
```

### Bug 2: Spinner written to wrong file when using profiles (`persona/settings.go:11`)

```go
// Current: always writes to ~/.claude/settings.json (ignores CLAUDE_CONFIG_DIR)
func GetClaudeSettingsPath() string {
    homeDir, _ := os.UserHomeDir()
    return filepath.Join(homeDir, ".claude", "settings.json")
}

// Fix: respect CLAUDE_CONFIG_DIR
func GetClaudeSettingsPath() string {
    if configDir := os.Getenv("CLAUDE_CONFIG_DIR"); configDir != "" {
        return filepath.Join(configDir, "settings.json")
    }
    homeDir, _ := os.UserHomeDir()
    return filepath.Join(homeDir, ".claude", "settings.json")
}
```

Effect: `godo cc -p work` sets `CLAUDE_CONFIG_DIR=~/.claude-profiles/work/` but
spinner verbs are written to `~/.claude/settings.json` instead of
`~/.claude-profiles/work/settings.json`. Each profile should have its own spinner.

### Bug 3: Session sharing between profiles

`CLAUDE_CONFIG_DIR` isolates Claude's config (settings, conversation history).
However, `settings.local.json` is read from the PROJECT's `.claude/` directory,
not from the profile dir. So `DO_CLAUDE_*` settings are always project-local,
not profile-local. Intentional? Worth documenting clearly.

### Bug 4: Checklist check incorrectly applied to Write/Edit (FIXED in current code)

Original buggy code:
```go
case "Task", "Write", "Edit":
    if output := checkChecklistRequirement(input); output != nil {
        return output
    }
```

`checkChecklistRequirement` extracts `prompt` from ToolInput. Write/Edit have
`file_path`, not `prompt`. This caused ALL Write/Edit to be denied when `.do/jobs/` existed.

Fixed: checklist check applies to Task only.

---

## Adding a New Hook Handler

1. `internal/hook/<event>.go`:
   ```go
   func HandleMyEvent(input *Input) *Output { ... }
   ```

2. `internal/cli/hook.go` — add to `hookHandlers` map:
   ```go
   "my-event": hook.HandleMyEvent,
   ```

3. `.claude/settings.json` — add hook entry:
   ```json
   "MyEvent": [{"hooks": [{"command": "godo hook my-event", "type": "command"}]}]
   ```

4. `internal/hook/types.go` — add constant:
   ```go
   EventMyEvent EventType = "MyEvent"
   ```

---

## Reference: Original Source Locations

- `~/Work/do-focus.workspace/do-focus/` — do-focus workspace (primary original)
- `~/Abyss/Workspace/do-focus/` — secondary copy
- `~/Work/moai-adk/` — moai-adk (hook system origin)
- `~/Abyss/Workspace/moai-adk/` — secondary copy
