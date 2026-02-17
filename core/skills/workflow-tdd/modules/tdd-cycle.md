# TDD RED-GREEN-REFACTOR Cycle [HARD]

## TDD Additional Steps [HARD]

### RED Phase
- [HARD] Write a failing test first — must fail because implementation does not yet exist
- Test describes the desired behavior
- Test name documents the requirement
- Only one test at a time, focused and specific

### GREEN Phase
- [HARD] Write minimal code to pass the test — no over-engineering
- Only satisfy the current test
- Do not add extra features or premature abstractions
- Focus on correctness, not elegance

### REFACTOR Phase
- [HARD] Clean up code while keeping tests green — no behavior changes
- Remove duplication
- Improve naming
- Extract methods/functions
- Apply SOLID principles where appropriate

## Non-TDD Path
- When TDD is not selected, proceed in implement → verify order
- [HARD] Testable code → follow testing-rules.md (behavior-based, FIRST, Real DB, etc.)
- [HARD] Non-testable changes (CSS, config, docs, hooks, etc.) → specify alternative verification (build check, manual check, etc.)

## TDD Quality Checks
- Test written BEFORE implementation code — never after
- Each RED-GREEN-REFACTOR cycle covers one behavior
- Minimum coverage per commit: 80% (configurable)
- TRUST 5 quality gates apply to TDD code as well
