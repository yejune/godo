---
name: do-workflow-report
description: >
  Generates completion reports after all checklist items are done. Aggregates
  sub-checklist results, lessons learned, and test outcomes into a final
  report.md. Adapted from Do's checklist-based workflow.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-15"
  tags: "report, completion, summary, lessons-learned"

# Do Extension: Triggers
triggers:
  keywords: ["report", "completion", "summary", "완료", "보고"]
  agents: ["manager-docs", "manager-quality"]
  phases: ["report"]
---

# Report Workflow Orchestration

## Purpose

Generate a completion report after all checklist items are done. Aggregates results from sub-checklists, collects lessons learned, summarizes test outcomes, and produces a final report.md in the job directory.

## Scope

- Final step of Do's checklist-based workflow
- Consumes checklist.md + checklists/*.md + git history
- Produces report.md as the definitive record of work done

## Input

- Completed checklist at `.do/jobs/{YY}/{MM}/{DD}/{title}/checklist.md`
- Completed sub-checklists at `.do/jobs/{YY}/{MM}/{DD}/{title}/checklists/*.md`
- Git log for commit history during the job

## Prerequisites

- [HARD] All checklist items must be `[o]` (done) or `[x]` (failed) -- no `[ ]` or `[~]` remaining
- [HARD] If incomplete items exist, report workflow refuses to run and directs to run workflow
- [HARD] Test failures in sub-checklists block report generation -- fix tests first

---

## Phase Sequence

### Phase 1: Data Collection

Gather all information needed for the report:

1. **Checklist Summary**: Read checklist.md, count items by status
   - Total tasks, completed (`[o]`), failed (`[x]`), blocked (`[!]`)
2. **Sub-checklist Details**: Read each checklists/{NN}_{agent}.md
   - Extract Lessons Learned sections
   - Extract commit hashes from Progress Logs
   - Extract Critical Files (modified files)
3. **Git History**: Run `git log` and `git diff --stat` for the job period
   - Changed files with line counts
   - Commit messages and authors
4. **Test Results**: Collect latest test run output
   - Pass/fail counts, coverage percentage

### Phase 2: Report Generation

Agent: Task(report-agent) or direct generation

Generate `.do/jobs/{YY}/{MM}/{DD}/{title}/report.md` using the template from dev-checklist.md:

```markdown
## Completion Report

### Execution Summary
- Completed: {N}/{M} tasks
- Duration: {start} ~ {end}

### Plan vs Actual Changes
- (Differences from original plan, with reasons)
- (If no changes: "Executed as planned")

### Test Results
- Total: {pass}/{total} passed
- Coverage: {N}% (if measurable)
- Failures/Skips: none or details

### Changed Files Summary
- `path/to/file.go` -- one-line change description
- `path/to/test.go` -- added tests

### Unresolved Items
- (Follow-up work needed)
- (Known constraints)
- (If none: "None")

### Key Lessons
- (Synthesized from sub-checklist Lessons Learned)
- (Insights to share with team/project)
```

### Phase 3: Verification

Verify report accuracy:

1. **File list match**: Changed files in report match `git diff --stat`
2. **Task count match**: Execution summary matches checklist.md item counts
3. **No test failures**: Test results section shows no failures
4. **Lessons populated**: Key Lessons section is not empty
5. **Unresolved items tracked**: If any exist, suggest follow-up plan or issue creation

### Phase 4: Presentation

Display report summary to user:

1. Show execution summary (tasks completed, duration)
2. Show key lessons (top 3)
3. Show any unresolved items

AskUserQuestion with options:
- "Looks good, we're done" -> Finalize
- "Create follow-up plan for unresolved items" -> Trigger plan workflow
- "Review full report" -> Display complete report.md content

---

## Completion Criteria

- Phase 1: All data collected from checklists, git, and tests
- Phase 2: report.md generated at correct `.do/jobs/` path
- Phase 3: Report verified against git diff and checklist counts
- Phase 4: Summary presented to user with next step options
- [HARD] report.md contains no test failures (fix first, then report)
- [HARD] Changed files summary matches `git diff --stat`

---

Version: 1.0.0
Updated: 2026-02-15
