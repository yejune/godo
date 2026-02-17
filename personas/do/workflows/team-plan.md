---
name: do-workflow-team-plan
description: >
  Team-based plan creation using Agent Teams API. Spawns parallel research
  team (researcher + analyst + architect), synthesizes findings into
  analysis.md + architecture.md, then generates plan.md + checklist.md.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-17"
  tags: "team, plan, research, parallel, agent-teams"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 3500

# Do Extension: Triggers
triggers:
  keywords: ["team plan", "parallel plan", "team research"]
  agents: ["team-researcher", "team-analyst", "team-architect"]
  phases: ["plan"]
---

# Team Plan Workflow Orchestration

## Purpose

Create comprehensive plan documents through parallel team-based research and analysis. Three specialized teammates explore the codebase, analyze requirements, and design the technical approach simultaneously. Results are synthesized into unified plan artifacts.

## Scope

- Team variant of the plan workflow
- Uses Agent Teams API for parallel execution
- Produces the same artifacts as solo plan workflow but faster

## Prerequisites

- Agent Teams feature enabled (CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1)
- Triggered by: Team mode active OR auto-detected complexity >= threshold
- Fallback: If Agent Teams unavailable, falls back to solo plan workflow (plan.md)

---

## Phase Sequence

### Phase 0: Team Setup

1. Create team:
   ```
   TeamCreate(team_name: "do-plan-{feature-slug}")
   ```

2. Create shared task list:
   ```
   TaskCreate: "Explore codebase architecture and dependencies"
   TaskCreate: "Analyze requirements and identify edge cases"
   TaskCreate: "Design technical approach and evaluate alternatives"
   TaskCreate: "Synthesize findings into plan documents" (blocked by above 3)
   ```

### Phase 1: Spawn Research Team

Spawn 3 teammates in parallel:

**Teammate 1 - researcher** (team-researcher agent, haiku model):
- Prompt: "Explore the codebase for {feature_description}. Map architecture, find relevant files, identify dependencies and patterns. Report findings to the team lead."
- Role: Codebase exploration (fastest, uses haiku)
- Output: Architecture map, dependency graph, relevant file list

**Teammate 2 - analyst** (team-analyst agent, inherit model):
- Prompt: "Analyze requirements for {feature_description}. Identify user stories, acceptance criteria, edge cases, risks, and constraints. Report findings to the team lead."
- Role: Requirements analysis
- Output: MoSCoW requirements, risk assessment, edge cases

**Teammate 3 - architect** (team-architect agent, inherit model):
- Prompt: "Design the technical approach for {feature_description}. Evaluate implementation alternatives, assess trade-offs, propose architecture. Consider existing patterns found by the researcher. Report to the team lead."
- Role: Technical design
- Output: Architecture proposal, approach comparison, implementation order

### Phase 2: Parallel Research

Teammates work independently:
- researcher explores codebase (fastest, haiku)
- analyst defines requirements (medium)
- architect designs solution (receives researcher findings when available)

Orchestrator coordination:
- Receive progress messages automatically
- Forward researcher findings to architect when available via SendMessage
- Resolve any questions from teammates

### Phase 3: Synthesis

After all research tasks complete:

1. Collect findings from all three teammates
2. Generate artifacts (delegate to subagent or direct):
   - `analysis.md` from analyst + researcher findings
   - `architecture.md` from architect findings
   - `plan.md` synthesizing all three into actionable plan
   - `checklist.md` + `checklists/*.md` from plan decomposition

All artifacts at: `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/`

### Phase 4: User Approval

AskUserQuestion with options:
- "Approve and proceed to implementation" -> Guide to team-run or run workflow
- "Request modifications" -> Collect feedback, revise
- "Cancel" -> Exit

### Phase 5: Cleanup

1. Shutdown all teammates:
   ```
   SendMessage(type: "shutdown_request", recipient: "researcher")
   SendMessage(type: "shutdown_request", recipient: "analyst")
   SendMessage(type: "shutdown_request", recipient: "architect")
   ```
2. Execute /clear to free context for next phase

---

## Fallback

If team creation fails or AGENT_TEAMS not enabled:
- Fall back to solo plan workflow (workflows/plan.md)
- Log warning about team mode unavailability
- Continue with sequential Analysis -> Architecture -> Plan

---

## Completion Criteria

- Phase 0: Team created, task list established
- Phase 1: All 3 teammates spawned and working
- Phase 2: All research tasks complete, findings collected
- Phase 3: analysis.md + architecture.md + plan.md + checklist.md generated
- Phase 4: User approval obtained
- Phase 5: All teammates shut down, context cleared

---

Version: 1.0.0
Updated: 2026-02-17
