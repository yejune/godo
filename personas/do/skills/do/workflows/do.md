---
name: do-workflow-do
description: >
  Full autonomous plan -> checklist -> run -> test -> report pipeline.
  Default workflow when no subcommand is specified. Handles exploration,
  plan generation, checklist creation, implementation, and completion.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-15"
  tags: "do, autonomous, pipeline, default"

# Do Extension: Triggers
triggers:
  keywords: ["do", "autonomous", "pipeline", "build", "implement", "create"]
  agents: ["do"]
  phases: ["plan", "run", "test", "report"]
---

# Workflow: Do - Autonomous Development Pipeline

Purpose: Full autonomous workflow. User provides a goal, Do autonomously executes plan -> checklist -> run -> test -> report pipeline. This is the default workflow when no subcommand is specified.

Flow: Explore -> Plan -> Checklist -> Run -> Test -> Report -> Done

## Supported Flags

- --team: Force Team mode for all applicable phases (route to team-do.md)
- --solo: Force Do mode (single agent per phase)

**Default Behavior (no flag)**: Uses current DO_MODE setting.

## Phase 0: Exploration

Launch exploration agents (parallel in Do mode, direct in Focus mode):

Agent 1 - Explore (Explore subagent):
- Codebase analysis for task context
- Relevant files, architecture patterns, existing implementations

Agent 2 - Research (Explore subagent with broader scope):
- External documentation and best practices if needed
- Existing similar patterns in the codebase

After exploration completes:
- Synthesize findings into unified context
- Determine complexity (simple vs complex task)

If --team flag: Route to workflows/team-do.md instead.

## Phase 0 Completion: Routing Decision

- Single-domain simple task: Consider delegating directly to expert agent (skip full pipeline)
- Multi-domain or complex task: Proceed to full pipeline

User approval via AskUserQuestion:
- "Proceed to plan creation"
- "Delegate directly to expert agent" (for simple tasks)
- "Cancel"

## Phase 1: Plan

Execute workflows/plan.md:
- Complexity assessment -> Analysis (if complex) -> Architecture (if complex) -> Plan generation
- Output: `.do/jobs/{YY}/{MM}/{DD}/{title}/plan.md` (+ analysis.md, architecture.md if complex)
- User approval checkpoint before proceeding

## Phase 2: Checklist

After plan is approved:

1. Generate checklist.md from plan.md:
   - Each plan task becomes a checklist item
   - Items decomposed to 1-3 file changes each
   - Verification method specified per item
   - Dependencies linked with `depends on:`

2. Generate sub-checklists per agent:
   - `.do/jobs/{YY}/{MM}/{DD}/{title}/checklists/{NN}_{agent-topic}.md`
   - Each follows the sub-checklist template from dev-checklist.md
   - Problem Summary, Acceptance Criteria, Critical Files, FINAL STEP: Commit

3. Ask user about TDD preference (if not already decided in plan phase):
   - AskUserQuestion: "TDD로 개발할까요?"
   - "Yes, TDD" -> test workflow integrated into run
   - "No, implement first" -> standard run then verify

## Phase 3: Run

Execute workflows/run.md:
- Read checklist, dispatch agents per mode (Do/Focus/Team)
- Monitor progress, handle interruptions
- All agents commit their work with checklist updates

## Phase 4: Test

Execute workflows/test.md (if not already done via TDD in run phase):
- Run full test suite
- Verify coverage targets (85%+)
- Fix any failures before proceeding

## Phase 5: Report

Execute workflows/report.md:
- Aggregate results from all sub-checklists
- Generate completion report
- Present summary and next steps to user

## Execution Summary

1. Parse arguments (extract flags: --team, --solo)
2. Check DO_MODE and determine execution strategy
3. If --team: Route to workflows/team-do.md
4. Execute Phase 0 (exploration)
5. Routing decision (simple direct delegation vs full pipeline)
6. User confirmation via AskUserQuestion
7. Phase 1 (Plan): Read workflows/plan.md
8. Phase 2 (Checklist): Generate checklist.md + sub-checklists
9. Phase 3 (Run): Read workflows/run.md
10. Phase 4 (Test): Read workflows/test.md (if not TDD)
11. Phase 5 (Report): Read workflows/report.md
12. Display final summary to user

---

Version: 1.0.0
Updated: 2026-02-15
