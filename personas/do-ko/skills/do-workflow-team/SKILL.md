---
name: do-workflow-team
description: >
  Agent Teams workflow management for Do. Handles team creation,
  teammate spawning, task decomposition, inter-agent messaging,
  file ownership enforcement, and graceful shutdown.
  Integrates with checklist workflow for team-based Plan and Run phases.
  Supports dual-mode execution with automatic fallback to Do mode (sub-agent).
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Task TeamCreate TeamDelete SendMessage TaskCreate TaskUpdate TaskList TaskGet Read Grep Glob AskUserQuestion
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "experimental"
  updated: "2026-02-16"
  modularized: "false"
  tags: "team, agent-teams, collaboration, parallel, file-ownership"
  related-skills: "do-workflow-plan, do-workflow-ddd, do-workflow-tdd"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 8000

# Do Extension: Triggers
triggers:
  keywords: ["team", "agent-team", "parallel", "collaborate", "teammates", "team mode"]
  agents: ["do"]
  phases: ["plan", "run"]
---

# Do Agent Teams Workflow

## Overview

This skill manages Agent Teams execution for Do workflows. When Team mode is active (via `godo mode team` or auto-escalation), Do operates as Team Lead coordinating persistent teammates with file ownership enforcement and checklist integration.

## Prerequisites

Agent Teams requires:
- `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1` in settings.json env
- Claude Code v2.1.32 or later
- Team mode activated: `godo mode team`

## Mode Selection

Team mode is selected when:
1. User explicitly switches: `godo mode team` or "team mode" keywords
2. Auto-escalation from Do mode: 10+ files, 3+ domains
3. Plan phase determines parallel research is beneficial

Fallback: If Agent Teams unavailable, fall back to Do mode (sub-agent sequential/parallel).

## Team Lifecycle

### Phase 1: Team Creation

```
TeamCreate(team_name: "do-{workflow}-{timestamp}")
```

Team naming convention:
- Plan phase: `do-plan-{title}`
- Run phase: `do-run-{title}`
- Debug: `do-debug-{issue}`

### Phase 2: Task Decomposition

Before spawning teammates, create complete shared task list from checklist:

```
TaskCreate(subject: "Task description", description: "Detailed requirements")
```

Rules for task decomposition:
- Each task is self-contained (one clear deliverable)
- Define dependencies between tasks (addBlockedBy)
- Assign file ownership boundaries per teammate role
- Tasks map to checklist items and sub-checklists
- 1 item = 1-3 file changes (Do's granularity rule)

### Phase 3: Teammate Spawning

Spawn teammates using Task tool with team_name parameter:

```
Task(
  subagent_type: "team-backend-dev",
  team_name: "do-run-{title}",
  name: "backend-dev",
  prompt: "You are the backend developer. File ownership: {files}. Sub-checklist: {path}."
)
```

Spawning rules:
- Include sub-checklist path in spawn prompt
- Assign file ownership from Critical Files section
- Include Docker environment info (container name, service, domain)
- Include commit instructions: "git add specific files && git diff --cached && git commit in single Bash call"
- Include checklist update instructions: "Update sub-checklist state with each progress step"
- Include DO_JOBS_LANGUAGE for document language

### Phase 4: Coordination

Team Lead monitors and coordinates:

1. Receive automatic messages from teammates
2. Use SendMessage for direct coordination
3. Broadcast critical updates (use sparingly -- expensive)
4. Resolve file ownership conflicts
5. Reassign tasks if teammate is blocked

Coordination patterns:
- Backend completes API: notify frontend-dev of available endpoints
- Implementation completes: assign quality validation
- Quality finds issues: direct fix to responsible teammate
- All tasks complete: begin shutdown

### Phase 5: Shutdown

Graceful shutdown sequence:

1. Verify all tasks completed via TaskList
2. Verify all checklist items are `[o]`
3. Send shutdown_request to each teammate
4. Wait for shutdown approval
5. Clean up: TeamDelete()

## File Ownership Strategy [HARD]

Prevent write conflicts by assigning exclusive file ownership.

### One File = One Owner

- No two teammates own the same file
- Sub-checklist Critical Files section defines ownership boundary
- `git add -A` / `git add .` ABSOLUTELY PROHIBITED
- Each agent stages only their own files explicitly: `git add path/to/my-file.go`

### Git Staging Rules for Team Agents [CRITICAL]

1. File-specific staging: `git add file1.go file2.go` (own files only)
2. Broad staging prohibited: `git add -A`, `git add .`, `git add --all` NEVER allowed
3. Atomic commit required: `git add && git diff --cached && git commit` in SINGLE Bash call
4. `git reset HEAD` prohibited: Would unstage other agent's files, causing data loss
5. Non-interference: Never touch other agent's files in staging area
6. On commit failure: Report error as-is, do NOT self-resolve (may destroy other agent's work)

### Ownership Conflict Resolution

| Situation | Resolution |
|-----------|-----------|
| File A owned by agent 1 only | Normal -- no conflict |
| File A needed by agent 1 and 2 | PROHIBITED -- sequential via `depends on:` |
| Agent 2 references agent 1's file (read-only) | Allowed -- reading is not ownership |
| Unexpected file modification needed | Check owner in checklist; if owner is `[~]`, declare `[!]` blocker |

## Team Patterns

### Plan Research Team
- Roles: researcher (haiku), analyst (inherit), architect (inherit)
- Use: Complex plan creation requiring multi-angle exploration
- Duration: Short-lived (plan phase only)
- All run in permissionMode: plan (read-only, Bash heredoc for .do/jobs/ output)

### Implementation Team
- Roles: backend-dev, frontend-dev, tester (all inherit)
- Use: Cross-layer feature implementation
- Duration: Medium (full run phase)
- File ownership strictly enforced

### Full-Stack Team
- Roles: backend-dev, frontend-dev, data-layer, quality (all inherit)
- Use: Large-scale full-stack features
- Duration: Medium-long

### Investigation Team
- Roles: hypothesis-1, hypothesis-2, hypothesis-3 (all haiku)
- Use: Complex debugging with competing theories
- Duration: Short

## Checklist Integration

Team mode integrates with Do's checklist system:

1. Main checklist (`checklist.md`) defines all items with agent assignment
2. Sub-checklists (`checklists/{order}_{agent-topic}.md`) distributed to teammates
3. Each teammate receives their sub-checklist path in spawn prompt
4. Teammate follows idempotent cycle: READ -> CLAIM -> WORK -> VERIFY -> RECORD -> COMMIT
5. Teammate commits code + checklist state together atomically
6. Team Lead monitors checklist states for progress
7. `[o]` items with commit hash = completed; `[~]` items = in progress; `[ ]` = not started

## Error Recovery

- Teammate crash: Spawn replacement with same role; new agent reads checklist, resumes from first non-`[o]` item
- Task stuck: Team Lead reassigns to different teammate
- File conflict: Team Lead mediates via SendMessage, adjusts ownership
- All teammates idle: Check if tasks remain; assign or shutdown
- Token limit: Shutdown team gracefully, fall back to Do mode for remaining work

---

Version: 1.0.0
Last Updated: 2026-02-16
