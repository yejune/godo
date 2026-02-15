---
name: do-workflow-team-do
description: >
  Team mode full pipeline combining plan, run, test, and report phases
  with Agent Teams API. Parallel research team for planning, parallel
  implementation team for execution, quality validation, and report.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-15"
  tags: "team, autopilot, parallel, agent-teams, plan, run, report"

# MoAI Extension: Triggers
triggers:
  keywords: ["team", "autopilot", "parallel", "팀"]
  agents: ["team-researcher", "team-analyst", "team-architect", "team-backend-dev", "team-frontend-dev", "team-tester", "team-quality"]
  phases: ["plan", "run", "test", "report"]
---

# Workflow: Team Do - Team Autopilot Pipeline

Purpose: Full autonomous pipeline using Agent Teams API for parallel execution. Combines plan (parallel research), run (parallel implementation), test, and report into a single team-orchestrated workflow.

Flow: Team Plan -> Checklist -> Team Run -> Quality -> Report -> Shutdown

## Prerequisites

- CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1 environment variable set
- Team mode active (DO_MODE=team or --team flag)
- If prerequisites not met: Warn user and fallback to Do mode (workflows/do.md)

---

## Part 1: Team Plan (Parallel Research)

### Team Composition

| Teammate | Agent | Model | Mode | Purpose |
|----------|-------|-------|------|---------|
| researcher | team-researcher | haiku | plan (read-only) | Codebase exploration |
| analyst | team-analyst | inherit | plan (read-only) | Requirements analysis |
| architect | team-architect | inherit | plan (read-only) | Technical design |

### Execution

1. Create team:
   ```
   TeamCreate(team_name: "do-plan-{feature-slug}")
   ```

2. Create shared task list:
   ```
   TaskCreate: "Explore codebase architecture and dependencies"
   TaskCreate: "Analyze requirements, user stories, and edge cases"
   TaskCreate: "Design technical approach and evaluate alternatives"
   ```

3. Spawn 3 teammates with investigation prompts:

   - **researcher**: "Explore the codebase for {feature_description}. Map architecture, find relevant files, identify dependencies and patterns. Report findings to the team lead."
   - **analyst**: "Analyze requirements for {feature_description}. Identify user stories, acceptance criteria, edge cases, risks, and constraints. Report findings to the team lead."
   - **architect**: "Design the technical approach for {feature_description}. Evaluate implementation alternatives, assess trade-offs, propose architecture. Report to the team lead."

4. Monitor parallel research:
   - Receive progress messages automatically
   - Forward researcher findings to architect when available
   - Resolve any questions from teammates

5. After all research tasks complete:
   - Collect findings from all three teammates
   - Synthesize into plan: Delegate to plan agent (sub-agent, NOT teammate) with all findings
   - Generate: analysis.md, architecture.md, plan.md at `.do/jobs/{YYMMDD}/{title}/`

6. User approval:
   - AskUserQuestion: "Plan complete. Proceed to implementation?"
   - "Proceed" -> Continue to checklist + Team Run
   - "Modify" -> Collect feedback, re-run synthesis
   - "Cancel" -> Shutdown team, exit

7. Shutdown plan team:
   ```
   SendMessage(type: "shutdown_request", recipient: "researcher")
   SendMessage(type: "shutdown_request", recipient: "analyst")
   SendMessage(type: "shutdown_request", recipient: "architect")
   ```

### Checklist Generation

After plan approval, generate checklist.md + sub-checklists (same as do.md Phase 2).

---

## Part 2: Team Run (Parallel Implementation)

### Team Composition

| Teammate | Agent | Model | Mode | Purpose |
|----------|-------|-------|------|---------|
| backend-dev | team-backend-dev | inherit | acceptEdits | Server-side implementation |
| frontend-dev | team-frontend-dev | inherit | acceptEdits | Client-side implementation |
| tester | team-tester | inherit | acceptEdits | Test creation and coverage |
| quality | team-quality | inherit | plan (read-only) | Quality validation |

Adjust team based on project needs (e.g., skip frontend-dev for backend-only tasks).

### File Ownership [HARD]

Assign file ownership based on sub-checklist Critical Files sections:

- Each teammate owns specific files from their sub-checklist
- No two teammates may modify the same file
- If overlap detected: Resolve before spawning (split the file's changes or assign to one owner)

### Git Staging Safety [HARD]

Every teammate MUST follow these rules:

- [HARD] Stage only own files: `git add file1.go file2.go` (individual files only)
- [HARD] NEVER use broad staging: `git add -A`, `git add .`, `git add --all` are FORBIDDEN
- [HARD] Verify before commit: `git diff --cached --name-only` must show only owned files
- [HARD] Check for foreign files: Other teammates may have staged but not committed. If foreign files are staged, unstage them: `git reset HEAD <file>`
- [HARD] Unstage before commit: If anything unexpected is staged, `git reset HEAD <file>` first

### Execution

1. Create team:
   ```
   TeamCreate(team_name: "do-run-{feature-slug}")
   ```

2. Create shared task list from checklist items with dependencies

3. Spawn teammates with sub-checklist paths:
   - Each teammate prompt includes: SPEC summary, sub-checklist path, file ownership list, git staging rules, quality targets

4. Parallel implementation:
   - backend-dev: Implements server-side code (Task 1-2)
   - frontend-dev: Implements client-side code, waits for API contracts (Task 3)
   - tester: Writes integration tests after implementation completes (Task 4)
   - Teammates coordinate via SendMessage for API contracts and dependencies

5. Orchestrator coordination:
   - Forward API contract info from backend to frontend
   - Monitor task progress via TaskList
   - Resolve blocking issues
   - Redirect teammates if approach isn't working

6. Quality validation (after implementation + tests complete):
   - Assign quality validation task to quality teammate
   - Quality runs full test suite and checks coverage
   - Reports findings to team lead
   - If issues found: Direct fixes to responsible teammates

### Quality Gates

All must pass before proceeding:
- Zero lint errors
- Zero type errors
- Coverage targets met (85%+ overall)
- All checklist acceptance criteria verified
- All tests pass

---

## Part 3: Report

After quality validation passes:

1. Generate completion report (same as report.md workflow):
   - Aggregate sub-checklist results
   - Collect lessons learned from all teammates
   - Summarize test outcomes

2. Display checklist final state to user

3. AskUserQuestion with next steps

---

## Part 4: Cleanup

1. Shutdown all teammates gracefully:
   ```
   SendMessage(type: "shutdown_request", recipient: "backend-dev")
   SendMessage(type: "shutdown_request", recipient: "frontend-dev")
   SendMessage(type: "shutdown_request", recipient: "tester")
   SendMessage(type: "shutdown_request", recipient: "quality")
   ```

2. Display final summary to user

---

## Fallback

If team creation fails or Agent Teams not enabled at any point:
- Shutdown remaining teammates gracefully
- Fall back to Do mode workflow (workflows/do.md)
- Continue from last completed phase
- Log warning about team mode unavailability

---

Version: 1.0.0
Updated: 2026-02-15
