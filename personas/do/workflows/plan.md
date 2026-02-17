---
name: do-workflow-plan
description: >
  Creates comprehensive plan documents through complexity assessment, optional
  analysis/architecture phases, and plan generation. First step of the
  Do checklist-based workflow. Handles simple and complex task planning.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-17"
  tags: "plan, analysis, architecture, design, requirements"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 4000

# Do Extension: Triggers
triggers:
  keywords: ["plan", "design", "architect", "requirements", "analyze"]
  agents: ["expert-analyst", "expert-architect", "Explore"]
  phases: ["plan"]
---

# Plan Workflow Orchestration

## Purpose

Create comprehensive plan documents by assessing task complexity and running the appropriate pipeline. Simple tasks go directly to plan generation. Complex tasks go through Analysis -> Architecture -> Plan.

All artifacts are stored at `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/`.
The `{title-kebab-case}` language is determined by `DO_JOBS_LANGUAGE` env var ("en" default).

## Scope

- First step of Do's checklist-based workflow
- Outputs feed into checklist generation, then run workflow
- Complexity assessment determines whether Analysis/Architecture phases are needed

## Input

- $ARGUMENTS: Natural language task description or explicit plan request
- Context: Current codebase state from git status and exploration

## Context Loading

Before execution, gather context:

- `git status`, `git branch`, `git log --oneline -5`
- Scan for existing `.do/jobs/` plans to avoid duplication
- Read relevant project files if they exist

---

## Phase Sequence

### Phase 1: Complexity Assessment

Evaluate the task against complexity criteria:

**Complex task** (ANY one condition triggers Analysis/Architecture):
- 5+ files expected to change
- New library/package/module creation
- System migration or technology stack change
- 3+ domain integration (backend + frontend + DB, etc.)
- Abstraction layer design needed (interfaces, provider patterns, plugins)
- Architecture change (monolith -> microservice, sync -> async, etc.)

**Simple task** (ALL conditions met):
- 4 or fewer files to change
- Implementation within existing patterns
- Single domain work
- No architecture changes

**Uncertain**: Use AskUserQuestion -- "Analysis/Architecture phases needed?"
Options: "Yes, start from analysis" / "No, go straight to plan"

### Phase 2A: Analysis (Complex Tasks Only)

Agent: Task(expert-analyst) or Task(Explore) for codebase research

Input: User request + Phase 1 complexity assessment

Tasks for analyst:
- Reverse-engineer existing code and system (code-based, not guesswork)
- Classify requirements using MoSCoW priority (MUST/SHOULD/COULD/WON'T)
- Compare at least 2 technical candidates with pros/cons
- Identify risks with impact levels (HIGH/MEDIUM/LOW)
- Map change scope and affected files/modules

Output: `.do/jobs/{YY}/{MM}/{DD}/{title}/analysis.md`

[HARD] Analysis must complete before Architecture proceeds.

### Phase 2B: Architecture (Complex Tasks Only)

Agent: Task(expert-architect) or domain-appropriate expert

Input: analysis.md from Phase 2A

Tasks for architect:
- Design system structure with ASCII diagram
- Define core interfaces at code level (not pseudocode)
- Compare at least 2 approaches with selection rationale
- Plan implementation order by file with phase numbering
- Define testing strategy (Unit/Integration with file paths)
- Cross-verify all MUST/SHOULD requirements from analysis.md are addressed

Output: `.do/jobs/{YY}/{MM}/{DD}/{title}/architecture.md`

[HARD] Architecture must complete before Plan proceeds.

### Phase 3: Plan Generation

Agent: Task(plan-agent) or direct generation

Input: User request (simple) OR analysis.md + architecture.md (complex)

Tasks:
- Create detailed implementation plan with task decomposition
- Each task: 1-3 files max, independently completable and verifiable
- Include dependency graph between tasks
- Specify verification method per task (test file path or build check)
- Ask user about TDD preference via AskUserQuestion

Output: `.do/jobs/{YY}/{MM}/{DD}/{title}/plan.md`

### Phase 4: User Approval

Tool: AskUserQuestion

Display plan summary and present options:

- "Proceed to implementation" -> Guide to checklist generation, then run workflow
- "Modify plan" -> Collect feedback, re-run Phase 3
- "Cancel" -> Exit with no further action

### Plan Mode (Shift+Tab) Integration [HARD]

When Claude Code Plan Mode is entered (Shift+Tab):

- [HARD] Save location: `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/plan.md`
- [HARD] NEVER use `~/.claude/plans/` -- even if system suggests it
- [HARD] Create directory if it doesn't exist
- [HARD] Date folder uses YY/MM/DD nested format (e.g., 26/02/17)

---

## Solo vs Team Execution

### Solo Mode (Do/Focus)

A single agent performs both exploration and analysis:
- One expert-analyst agent handles codebase research + requirements analysis
- One expert-architect agent handles design
- Sequential execution: Analysis -> Architecture -> Plan

### Team Mode

Parallel research with specialized teammates:
- See `team-plan.md` workflow for team-based plan execution
- researcher + analyst + architect work in parallel
- Results synthesized into unified plan

---

## Completion Criteria

- Phase 1: Complexity assessed (simple or complex determined)
- Phase 2A (complex only): analysis.md created with MoSCoW requirements and 2+ candidates
- Phase 2B (complex only): architecture.md created with ASCII diagram and 2+ approaches
- Phase 3: plan.md created with task decomposition and verification methods
- Phase 4: User approval obtained
- Plan file saved at correct `.do/jobs/` path (never `~/.claude/plans/`)

---

Version: 1.0.0
Updated: 2026-02-17
