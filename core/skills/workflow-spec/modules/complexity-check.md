# Complexity Check and Workflow Selection [HARD]

## Workflow Selection

### Complexity Assessment [HARD]

- [HARD] Assess complexity at task start to select workflow
- [HARD] **Complex tasks**: Analysis → Architecture → Plan → Checklist → Develop → Test → Report
- [HARD] **Simple tasks**: Plan → Checklist → Develop → Test → Report
- [HARD] TDD workflow (applicable to both simple/complex): Replace Develop phase with Test(RED) → Develop(GREEN) → Refactor

### Complex Task Criteria (Analysis/Architecture required if any apply) [HARD]
- [HARD] 5+ files expected to change
- [HARD] Creating new library/package/module
- [HARD] System migration/transition (tech stack replacement, data migration, etc.)
- [HARD] Multi-domain integration (backend + frontend + DB — 3+ domains)
- [HARD] Abstraction layer design needed (interfaces, provider patterns, plugin architecture, etc.)
- [HARD] Architecture change of existing system (monolith → microservices, sync → async, etc.)

### Simple Task Criteria (Analysis/Architecture can be skipped if all apply)
- 4 or fewer files changing
- Implementation within existing patterns (new endpoint, bug fix, etc.)
- Single domain work
- No architecture changes

### When Assessment is Uncertain [HARD]
- [HARD] If complexity is ambiguous, ask user (AskUserQuestion): "Does this task need Analysis/Architecture phases?"
- [HARD] Options: "Yes, start from analysis" / "No, go straight to plan"

## Analysis Phase [HARD]

- [HARD] Purpose: Understand current state + organize requirements + compare technology options
- [HARD] Owner: analyst agent (expert-analyst or relevant domain expert)
- [HARD] Artifact: `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/analysis.md`
- [HARD] Cannot proceed to Architecture before Analysis is complete

Analysis agent responsibilities:
1. **Current state investigation**: Reverse-engineer existing code/system, analyze usage patterns, trace data flows
2. **Requirements analysis**: Classify functional requirements (MUST/SHOULD/COULD/WON'T), identify non-functional requirements
3. **Technology options comparison**: List candidate libraries/approaches, create pros/cons analysis table
4. **Change scope identification**: List affected files/modules, draft migration strategy
5. **Risk identification**: Technical risks, compatibility issues, performance impact

## Architecture Phase [HARD]

- [HARD] Purpose: Solution design + interface specification + implementation order
- [HARD] Owner: architect agent (expert-backend, expert-frontend, etc. by domain)
- [HARD] Input: analysis.md (Analysis phase artifact)
- [HARD] Artifact: `.do/jobs/{YY}/{MM}/{DD}/{title-kebab-case}/architecture.md`
- [HARD] Cannot proceed to Plan before Architecture is complete

Architecture agent responsibilities:
1. **System structure design**: Full architecture diagram (ASCII or text), layer/component relationships
2. **Directory structure**: File/folder tree design, role of each file specified
3. **Interface specification**: Core type/interface definitions (code level), contracts between components
4. **Implementation details**: Internal logic of each component, error handling strategy, configuration structure
5. **Alternative comparison**: Compare at least 2 approaches (selection reason + rejection reason)
6. **Implementation order**: Phase-by-phase implementation order, dependency graph
7. **Risk mitigation**: Specific mitigation strategies for risks identified in Analysis

## General Rules
- [HARD] At plan phase, confirm TDD with user (AskUserQuestion): "Develop with TDD?"
- [HARD] If requirements are ambiguous, ask clarifying questions (AskUserQuestion) before any planning
- [HARD] One agent changes maximum 3 files — decompose into smaller tasks if exceeding
- [HARD] Checklist items must be granular enough to complete within agent token budget — split large items
- [HARD] On requirement addition/change, must update all artifacts: plan.md → checklist.md → checklists/*.md → report.md
- [HARD] Mismatch between documentation and actual work is a VIOLATION — documents must always reflect current state

## TDD Selection [HARD]
- [HARD] RED: Write failing test first — must fail since implementation doesn't exist yet
- [HARD] GREEN: Write minimal code to pass test — no over-engineering
- [HARD] REFACTOR: Clean up code while keeping tests green — no behavior changes

## Non-TDD Selection
- Proceed in implement → verify order
- [HARD] Testable code → follow testing-rules.md (behavior-based, FIRST, Real DB, etc.)
- [HARD] Non-testable changes (CSS, config, docs, hooks, etc.) → specify alternative verification (build check, manual check, etc.)
