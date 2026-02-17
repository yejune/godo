# godo CLI Analysis

## Source Location
/Users/max/Work/do-focus.workspace/do-focus/cmd/godo/ (41 files, 9,796 lines)

## Dependencies
- go module: github.com/yejune/do-focus (go 1.21)
- External: only gopkg.in/yaml.v3

## CLI Command Tree
```
godo
├── sync [--dry-run] [--init-custom]    # Framework install/update
├── selfupdate                           # Self-update via brew
├── hook <event>                         # Claude Code hook dispatcher
│   ├── session-start                    # Persona injection, system message
│   ├── pre-tool                         # Security policy (1,181 lines)
│   ├── post-tool-use                    # Transcript compaction
│   ├── user-prompt-submit               # Persona reminder
│   ├── subagent-stop                    # Agent finalization
│   ├── stop                             # Session termination guard
│   ├── compact                          # Transcript compression
│   └── session-end                      # Cleanup
├── mode [get|set <mode>]               # do/focus/team + permission modes
├── moai sync --source <path>           # MoAI asset sync (AST transform)
├── claude [--profile <name>]           # Claude launcher + profiles
├── create <agent|skill> <name>         # Scaffolding
├── lint [--all|setup]                  # Code linting
├── glm [setup]                         # GLM backend
├── spinner [apply|restore]             # Korean spinner verbs
├── statusline                          # Status rendering
├── rank [login|status|logout]          # Rank system
└── version / help
```

## Key Components

### moai sync (AST Transform Pipeline)
1. Load sync metadata (.do/moai-sync.json)
2. Git pull source repo
3. Detect changed files (git log)
4. Block if go.mod/go.sum changed
5. AST transform: package hook → package main, prefix exports with "moai"
6. Build verify (go build, rollback on fail)
7. Copy assets with brand prefix conversion
8. Merge hooks into settings.json
9. Save metadata

### Hook System (I/O Contract)
- Input: JSON via stdin (session info, tool info, agent info)
- Output: JSON to stdout (continue, decision, system message)
- PreToolUse: security deny patterns, dependency validation
- PostToolUse: transcript compaction, persona reminder injection

### Persona Loader
- Resolves .claude/personas/do/ or personas/do/
- Loads character YAML (honorific, tone, relationship)
- 4 types: young-f (default), young-m, senior-f, senior-m
- BuildReminder() → honorific + tone string for hooks

## Core vs DO-Specific

### Reusable (→ convert)
- Hook dispatcher + I/O contract
- Security policy patterns
- Mode system (execution + permission)
- AST transform pipeline
- Claude profile management
- Lint orchestration
- Scaffolding (create agent/skill)

### DO-Specific (→ persona)
- Persona loader (character selection)
- Korean spinner verb system
- Hook persona injection logic
- Job state tracking (.do/ paths)
- Rank system (authentication + transcript)
