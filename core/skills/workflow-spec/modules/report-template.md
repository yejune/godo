# Completion Report Template and Checklist Display Rules [HARD]

## Checklist Display Obligation [HARD]

- [HARD] When all tasks are complete, orchestrator must display current checklist state to user
- [HARD] Display format: state symbol (`[o]`/`[x]`/`[~]`/`[ ]`/`[!]`) + one-line summary per item
- [HARD] If incomplete items exist, suggest next action (resume or request user decision)
- [HARD] Stop hook blocks termination when active checklist detected — can terminate after displaying checklist

## Completion Report [HARD]

- [HARD] When all checklists are complete, write final report to `report.md` in the job folder

### Report Template
```markdown
## Completion Report

### Execution Summary
- Completed: {N}/{M} tasks (e.g., 3/3)
- Duration: {start datetime} ~ {end datetime}

### Changes from Plan
- (What changed from original plan, why)
- (If no changes: "Proceeded as planned")

### Test Results
- Total: {pass}/{total} passed
- Coverage: {N}% (if measurable)
- Failures/Skips: none or detailed breakdown

### Changed Files Summary
- `path/to/file.go` -- one-line summary of change
- `path/to/test.go` -- added tests

### Unresolved Items
- (Items requiring follow-up work)
- (Known constraints)
- (If none: "None")

### Key Lessons
- (Synthesis of sub-task Lessons Learned)
- (Insights to share with team/project)
```

### Report Rules
- [HARD] When writing report.md, reference checklist.md and checklists/*.md for synthesis
  - Execution summary ← checklist.md state aggregation
  - Key lessons ← synthesis of each sub-checklist's Lessons Learned
  - Changed files ← each sub-checklist's Critical Files + `git diff --stat`
- [HARD] If `Unresolved Items` exist, register as follow-up plan or issue
- [HARD] Cannot write completion report if test results contain failures — resolve first
- [HARD] Changed files summary must match `git diff --stat`
