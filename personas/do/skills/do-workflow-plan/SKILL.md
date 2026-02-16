---
name: do-workflow-plan
description: >
  Plan workflow orchestration with EARS+MoSCoW requirements, complexity assessment,
  Analysis-Architecture-Plan pipeline, and checklist generation for Do's
  checklist-driven development methodology.
  Use when creating plans, writing requirements, defining acceptance criteria,
  assessing complexity, or orchestrating the plan phase.
  Do NOT use for implementation (use do-workflow-ddd or do-workflow-tdd instead)
  or documentation generation (use do-workflow-project instead).
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Write Edit Bash Grep Glob mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-16"
  modularized: "false"
  tags: "workflow, plan, ears, moscow, requirements, checklist, analysis, architecture"
  agent: "manager-plan"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# Do Extension: Triggers
triggers:
  keywords:
    - "plan"
    - "requirement"
    - "EARS"
    - "MoSCoW"
    - "acceptance criteria"
    - "planning"
    - "analysis"
    - "architecture"
    - "checklist"
    - "complexity"
  phases: ["plan"]
  agents: ["manager-plan", "manager-strategy", "expert-analyst", "expert-architect"]
---

# Do Plan Workflow

## Quick Reference

Plan Workflow Orchestration -- comprehensive planning through complexity assessment, optional analysis/architecture phases, EARS+MoSCoW requirements, and checklist generation.

Core Capabilities:

- Complexity Assessment: Determines simple vs complex workflow path
- Analysis Phase: Current system analysis, requirements gathering, tech comparison (complex only)
- Architecture Phase: Solution design, interface specs, implementation order (complex only)
- EARS + MoSCoW: Structured requirements with prioritization
- Plan Document: Implementation roadmap with phases
- Checklist Generation: Agent-specific sub-checklists with 3-phase template

Output Location: `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/`

---

## Implementation Guide

### Complexity Assessment [HARD]

Every task starts with complexity assessment to determine workflow path.

Simple Task (ALL must apply):
- 4 or fewer file changes
- Within existing patterns (no new modules)
- Single domain work
- No architecture changes
- Workflow: Plan -> Checklist -> Develop -> Test -> Report

Complex Task (ANY triggers full pipeline):
- 5+ file changes expected
- New library/package/module creation
- System migration/transition
- 3+ domain integration (backend + frontend + DB)
- Abstraction layer design needed
- Existing system architecture change
- Workflow: Analysis -> Architecture -> Plan -> Checklist -> Develop -> Test -> Report

When uncertain: Ask user via AskUserQuestion: "Analysis/Architecture needed?"

### Simple Workflow: Plan -> Checklist

Step 1 - Requirements Gathering:
- Parse user request for scope and constraints
- Ask clarifying questions if requirements are ambiguous [HARD]
- Determine TDD preference via AskUserQuestion: "TDD?" [HARD]

Step 2 - Plan Document:
- Create `.do/jobs/{YY}/{MM}/{DD}/{title}/plan.md`
- Include: objectives, scope, approach, phases, risks
- EARS format for requirements when appropriate
- MoSCoW prioritization for multi-item plans

Step 3 - Checklist Generation:
- Generate `checklist.md` with numbered items
- Each item: 1-3 file changes + verification method [HARD]
- Items over 3 files: MUST split [HARD]
- Generate sub-checklists: `checklists/{order}_{agent-topic}.md`
- Sub-checklists follow 3-phase template (pre/execute/post)

### Complex Workflow: Analysis -> Architecture -> Plan

Step 1 - Analysis Phase:
- Delegate to expert-analyst agent
- Output: `.do/jobs/{YY}/{MM}/{DD}/{title}/analysis.md`
- Contents: current system analysis, requirements (MoSCoW), tech comparison, change scope, risks
- [HARD] Analysis must complete before Architecture begins

Step 2 - Architecture Phase:
- Delegate to expert-architect agent (receives analysis.md as input)
- Output: `.do/jobs/{YY}/{MM}/{DD}/{title}/architecture.md`
- Contents: system structure (ASCII diagram), directory layout, core interfaces (code-level), error handling, component implementations, approach comparison (min 2), testing strategy, implementation order, risk mitigation
- [HARD] Architecture must complete before Plan begins

Step 3 - Plan Phase:
- Create plan.md based on analysis + architecture outputs
- All MUST requirements from analysis reflected in plan
- Implementation order from architecture drives checklist sequence

Step 4 - Checklist Generation:
- Same as simple workflow but informed by architecture
- Dependencies between items use `depends on:` notation
- File ownership boundaries align with architecture component split

### EARS + MoSCoW Integration

EARS format for writing requirements:

| Type | Pattern | Example |
|------|---------|---------|
| Ubiquitous | System always does X | System logs all API requests |
| Event-driven | WHEN event THEN action | WHEN login fails THEN wait 3s before retry |
| State-driven | IF condition THEN action | IF offline THEN use local cache |
| Unwanted | System shall NOT X | System shall NOT store plaintext passwords |
| Optional | Where possible, X | Where possible, send email notifications |

MoSCoW for prioritization:

- MUST: Checklist items first, blocking
- SHOULD: Important, second priority
- COULD: Nice-to-have if time permits
- WON'T: Explicitly excluded, documented as out-of-scope

Conversion flow: EARS requirements -> MoSCoW priority -> implementation decomposition -> Test Strategy -> checklist items -> sub-checklists

### Test Strategy Pre-Declaration [HARD]

During plan phase, each checklist item MUST declare its test strategy:

| Code Type | Test Strategy | Example |
|-----------|--------------|---------|
| Business logic, API, data layer | Test type + target file | `unit: handler_test.go` |
| Complex features | Multiple types | `unit: validator_test.go + E2E: flow_test.go` |
| CSS, config, docs, hooks | `pass` + alternative | `pass (build check: go build ./...)` |

"pass" is a judgment, not a skip. It records WHY testing is unnecessary.

---

## Works Well With

- do-foundation-core: Core principles and checklist system
- do-workflow-ddd: DDD implementation after plan
- do-workflow-tdd: TDD implementation after plan
- do-workflow-team: Team mode plan with parallel research
- expert-analyst: Analysis phase execution
- expert-architect: Architecture phase execution
- manager-plan: Plan creation and checklist decomposition

---

Version: 1.0.0
Last Updated: 2026-02-16
