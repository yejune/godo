---
name: do-workflow-test
description: >
  TDD RED-GREEN-REFACTOR workflow for quality assurance. Do-unique workflow
  that enforces test-first development discipline with real DB testing and
  AI anti-pattern prevention.
license: Apache-2.0
compatibility: Designed for Claude Code
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-15"
  tags: "test, tdd, red-green-refactor, quality, coverage"

# Do Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 4000

# Do Extension: Triggers
triggers:
  keywords: ["test", "tdd", "coverage", "테스트", "red-green", "refactor"]
  agents: ["expert-testing", "manager-tdd"]
  phases: ["test"]
---

# Test Workflow Orchestration

## Purpose

Execute TDD RED-GREEN-REFACTOR cycles for test-driven quality assurance. This workflow can be invoked standalone or as part of the do pipeline. It enforces test-first discipline and prevents common AI testing anti-patterns.

## Scope

- Standalone TDD workflow or integrated into do pipeline
- Enforces dev-testing.md rules throughout
- Supports both new code (TDD) and existing code (DDD characterization tests)

## Input

- $ARGUMENTS: Feature or component to test, or checklist reference
- Context: Existing code to test, or specification for new code

---

## Phase Sequence

### Phase 1: Test Strategy Assessment

Determine testing approach based on code state:

**New Code (no existing implementation)**:
- Apply TDD: RED -> GREEN -> REFACTOR
- Agent: manager-tdd subagent

**Existing Code (modification or refactoring)**:
- Apply DDD: ANALYZE -> PRESERVE (characterization tests) -> IMPROVE
- Agent: manager-ddd subagent

**Mixed (new functions in existing files)**:
- TDD for new functions, characterization tests for existing behavior
- Agent: manager-tdd with DDD context

### Phase 2: RED - Write Failing Tests

Agent: Task(expert-testing) or Task(manager-tdd)

[HARD] Tests MUST be written BEFORE implementation code.

Tasks:
1. Write test that describes desired behavior
2. Verify the test FAILS (confirms it tests something new)
3. One test at a time, focused and specific
4. Test names describe scenarios: `test_login_fails_with_expired_token`
5. Use Arrange-Act-Assert (Given-When-Then) pattern

Quality rules (from dev-testing.md):
- [HARD] Test behavior, not implementation details
- [HARD] FIRST: Fast, Independent, Repeatable, Self-validating, Timely
- [HARD] Real DB only -- no mock DB, no in-memory DB, no SQLite substitution
- [HARD] Each test sets up and cleans up its own data (transaction rollback or truncate)
- [HARD] Use factory/builder patterns for test data creation
- [HARD] Unique identifiers per test (UUID/timestamp suffix) for parallel safety

### Phase 3: GREEN - Minimal Implementation

Agent: Task(expert-backend) or Task(expert-frontend) or direct (Focus mode)

[HARD] Write the SIMPLEST code that makes the test pass.

Rules:
- No premature optimization
- No premature abstraction
- Focus on correctness, not elegance
- Run tests after each change to verify GREEN state

### Phase 4: REFACTOR - Improve Quality

Agent: Same as Phase 3

[HARD] Improve code while keeping tests GREEN.

Tasks:
1. Extract patterns, remove duplication
2. Apply SOLID principles where appropriate
3. Run ALL tests after each refactoring step
4. No behavior changes -- tests must stay green throughout

### Phase 5: Verification

Run comprehensive verification:

1. Execute individual test: verify new test passes
2. Execute full test suite: verify no regressions
3. Check coverage: target 85%+ for modified code
4. Cross-validate: view <-> DB <-> worker <-> model <-> business logic <-> controller

---

## AI Anti-Pattern Prevention [HARD]

The following are FORBIDDEN during test writing:

- [HARD] Weakening assertions: changing `assertEqual` to `assertContains`, exact values to `any()`
- [HARD] Swallowing errors with try/catch to make tests green
- [HARD] Adjusting expected values to match wrong output (fix code, not test)
- [HARD] Using `time.sleep()` or arbitrary delays (find real timing cause)
- [HARD] Deleting or commenting out failing tests
- [HARD] Using wildcard matchers (`any()`, `mock.ANY`) when exact values are known
- [HARD] Testing only happy path -- error paths, edge cases, boundary values are REQUIRED

### Mutation Testing Mindset

Apply this thinking exercise to every test: "If I change one line of this code, does a test fail?"
If no test fails, the test coverage is insufficient. This is not a tool requirement -- it is a
discipline that every agent must apply when writing tests.

---

## Test Execution Rules [HARD]

- [HARD] All tests MUST pass -- no skip, no timeout bypass
- [HARD] "Timed out" is not an acceptable report -- must be resolved
- [HARD] When tests fail, fix the CODE, not the test
- [HARD] Individual tests first, then full suite when confident
- [HARD] Individual test < 3 minutes, full suite < 10 minutes
- [HARD] Tests run inside Docker containers: `docker compose exec <service> <test-command>`

---

## Characterization Tests (Existing Code)

For existing code that needs modification:

1. ANALYZE: Read existing code, map behavior and dependencies
2. PRESERVE: Write characterization tests capturing CURRENT behavior
   - These tests document what the code does NOW, not what it should do
   - They serve as a safety net for refactoring
3. IMPROVE: Make changes with characterization tests as guardrails
   - Run tests after each change
   - Any test failure means unintended behavior change

---

## Completion Criteria

- RED: Failing test written and verified to fail
- GREEN: Minimal implementation makes test pass
- REFACTOR: Code improved, all tests still pass
- Coverage: 85%+ for modified/new code
- Full suite: All tests pass (no regressions)
- Anti-patterns: None of the forbidden patterns used
- Commit: Test and implementation committed together with descriptive message

---

Version: 1.0.0
Updated: 2026-02-16
