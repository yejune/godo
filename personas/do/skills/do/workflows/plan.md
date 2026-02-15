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
  updated: "2026-02-15"
  tags: "plan, analysis, architecture, design, requirements"

# Do Extension: Triggers
triggers:
  keywords: ["plan", "design", "architect", "requirements", "analyze", "설계", "계획", "분석"]
  agents: ["expert-analyst", "expert-architect", "Explore"]
  phases: ["plan"]
---

# Plan Workflow Orchestration

## Purpose

Create comprehensive plan documents by assessing task complexity and running the appropriate pipeline. Simple tasks go directly to plan generation. Complex tasks go through Analysis -> Architecture -> Plan.

All artifacts are stored at `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/`.

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

Evaluate the task against complexity criteria from dev-workflow.md:

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
- Classify requirements using MoSCoW (MUST/SHOULD/COULD/WON'T)
- Compare at least 2 technical candidates with pros/cons
- Identify risks with impact levels (HIGH/MEDIUM/LOW)
- Map change scope and affected files/modules

Output: `.do/jobs/{YY}/{MM}/{DD}/{title}/analysis.md`
Template: dev-checklist.md Analysis template (Sections 1-6)

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
Template: dev-checklist.md Architecture template (Sections 1-10)

[HARD] Architecture must complete before Plan proceeds.

### Phase 3: Plan Generation

Agent: Task(plan-agent) or direct generation

Input: User request (simple) OR analysis.md + architecture.md (complex)

Tasks:
- Create detailed implementation plan with task decomposition
- Each task: 1-3 files max, independently completable and verifiable
- Include dependency graph between tasks
- Specify verification method per task (test file path or build check)
- Ask user about TDD preference via AskUserQuestion: "TDD로 개발할까요?"

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
- [HARD] Create `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/` directory if it doesn't exist
- [HARD] Date folder uses YY/MM/DD format (e.g., 26/02/15)

---

## Design/Plan Request Detection [HARD]

When user's natural language contains these patterns, trigger the plan workflow with the full Analysis -> Architecture -> Plan pipeline:

- Design: "설계해줘", "설계해", "design", "아키텍처 설계", "구조 설계"
- Plan: "플랜 짜줘", "계획 세워줘", "계획해줘", "plan", "플랜", "로드맵"
- Implementation questions: "어떻게 구현해야해?", "어떻게 만들어야해?", "구현 방법"
- Analysis: "분석해줘", "조사해줘", "파악해줘", "현황 분석"
- Composite: "~하고 싶어", "~만들고 싶어", "~개발하려면"

Execute 3 phases sequentially (each phase's output feeds the next):
1. **Analysis**: expert-analyst -> `analysis.md`
2. **Architecture**: expert-architect -> `architecture.md`
3. **Plan**: plan agent -> `plan.md`

After completion: AskUserQuestion "설계 완료! 구현 진행할까요?"
If approved -> Checklist generation -> Develop

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
Updated: 2026-02-15
