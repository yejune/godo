# Bug Fix Workflow [HARD]

## Reproduction-First Approach
- [HARD] Reproduction first: Write a failing test that proves the bug
- [HARD] Fix code until the test passes
- [HARD] Regression verification: Run the entire related test suite to check for side effects
- [HARD] Never fix a bug without a reproduction test
- [HARD] Never delete the test that caught the bug to "fix" it
