---
name: do-workflow-tdd
description: >
  Test-Driven Development workflow specialist using RED-GREEN-REFACTOR
  cycle for test-first software development.
  Use when developing new features from scratch, creating isolated modules,
  or when behavior specification drives implementation.
  Do NOT use for refactoring existing code (use do-workflow-ddd instead)
  or when behavior preservation is the primary goal.
license: Apache-2.0
compatibility: Designed for Claude Code
allowed-tools: Read Write Edit Bash(pytest:*) Bash(ruff:*) Bash(npm:*) Bash(npx:*) Bash(node:*) Bash(jest:*) Bash(vitest:*) Bash(go:*) Bash(cargo:*) Bash(mix:*) Bash(uv:*) Bash(bundle:*) Bash(php:*) Bash(phpunit:*) Grep Glob mcp__context7__resolve-library-id mcp__context7__get-library-docs
user-invocable: false
metadata:
  version: "1.0.0"
  category: "workflow"
  status: "active"
  updated: "2026-02-03"
  modularized: "true"
  tags: "workflow, tdd, test-driven, red-green-refactor, test-first"
  author: "MoAI-ADK Team"
  context: "fork"
  agent: "manager-tdd"
  related-skills: "do-workflow-ddd, do-workflow-testing, do-foundation-quality"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 5000

# MoAI Extension: Triggers
triggers:
  keywords: ["TDD", "test-driven development", "red-green-refactor", "test-first", "new feature", "greenfield"]
  phases: ["run"]
  agents: ["manager-tdd", "expert-backend", "expert-frontend", "expert-testing"]
---

# 테스트 주도 개발 (TDD) 워크플로우

## 개발 모드 구성 (중요)

[참고] 이 워크플로우는 `.moai/config/sections/quality.yaml`을 기반으로 선택됩니다:

```yaml
constitution:
  development_mode: hybrid    # 또는 ddd, tdd
  hybrid_settings:
    new_features: tdd        # 신규 코드 → TDD 사용 (이 워크플로우)
    legacy_refactoring: ddd  # 기존 코드 → DDD 사용
```

**이 워크플로우를 사용할 때:**
- `development_mode: tdd` → 항상 TDD 사용
- `development_mode: hybrid` + 새 패키지/모듈 → TDD 사용
- `development_mode: hybrid` + 기존 코드 리팩토링 → 대신 DDD 사용 (do-workflow-ddd)

**핵심 구분:**
- **새 파일/패키지** (아직 존재하지 않음) → TDD (이 워크플로우)
- **기존 코드** (파일이 이미 존재함) → DDD (ANALYZE-PRESERVE-IMPROVE)

## 빠른 참조

테스트 주도 개발은 구현 전에 테스트가 예상 동작을 정의하는 새로운 기능 생성을 위한 접근 방식을 제공합니다.

핵심 주기 - RED-GREEN-REFACTOR:

- RED: 원하는 동작을 정의하는 실패하는 테스트 작성
- GREEN: 테스트를 통과하게 하는 최소 코드 작성
- REFACTOR: 테스트를 녹색 상태로 유지하면서 코드 구조 개선

TDD를 사용할 때:

- 처음부터 새로운 기능 생성
- 기존 의존성이 없는 격리된 모듈 구축
- 동작 명세서가 개발을 주도할 때
- 명확한 계약이 있는 새 API 엔드포인트
- 정의된 동작이 있는 새 UI 컴포넌트
- Greenfield 프로젝트 (드뭄 - 일반적으로 Hybrid가 더 좋음)

When NOT to Use TDD:

- Refactoring existing code (use DDD instead)
- When behavior preservation is the primary goal
- Legacy codebase without test coverage (use DDD first)
- When modifying existing files (consider Hybrid mode)

---

## Core Philosophy

### TDD vs DDD Comparison

TDD Approach:

- Cycle: RED-GREEN-REFACTOR
- Goal: Create new functionality through tests
- Starting Point: No code exists
- Test Type: Specification tests that define expected behavior
- Outcome: New working code with test coverage

DDD Approach:

- Cycle: ANALYZE-PRESERVE-IMPROVE
- Goal: Improve structure without behavior change
- Starting Point: Existing code with defined behavior
- Test Type: Characterization tests that capture current behavior
- Outcome: Better structured code with identical behavior

### Test-First Principle

The golden rule of TDD is that tests must be written before implementation code:

- Tests define the contract
- Tests document expected behavior
- Tests catch regressions immediately
- Implementation is driven by test requirements

---

## Implementation Guide

### Phase 1: RED - Write a Failing Test

The RED phase focuses on defining the desired behavior through a failing test.

#### Writing Effective Tests

Before writing any implementation code:

- Understand the requirement clearly
- Define the expected behavior in test form
- Write one test at a time
- Keep tests focused and specific
- Use descriptive test names that document behavior

#### Test Structure

Follow the Arrange-Act-Assert pattern:

- Arrange: Set up test data and dependencies
- Act: Execute the code under test
- Assert: Verify the expected outcome

#### Verification

The test must fail initially:

- Confirms the test actually tests something
- Ensures the test is not passing by accident
- Documents the gap between current and desired state

### Phase 2: GREEN - Make the Test Pass

The GREEN phase focuses on writing minimal code to satisfy the test.

#### Minimal Implementation

Write only enough code to make the test pass:

- Do not over-engineer
- Do not add features not required by tests
- Focus on correctness, not perfection
- Hardcode values if necessary (refactor later)

#### Verification

Run the test to confirm it passes:

- All assertions must succeed
- No other tests should break
- Implementation satisfies the test requirements

### Phase 3: REFACTOR - Improve the Code

The REFACTOR phase focuses on improving code quality while maintaining behavior.

#### Safe Refactoring

With passing tests as a safety net:

- Remove duplication
- Improve naming and readability
- Extract methods and classes
- Apply design patterns where appropriate

#### Continuous Verification

After each refactoring step:

- Run all tests
- If any test fails, revert immediately
- Commit when tests pass

---

## TDD Workflow Execution

### Standard TDD Session

When executing TDD through manager-tdd:

Step 1 - Understand Requirements:

- Read SPEC document for feature scope
- Identify test cases from acceptance criteria
- Plan test implementation order

Step 2 - RED Phase:

- Write first failing test
- Verify test fails for the right reason
- Document expected behavior

Step 3 - GREEN Phase:

- Write minimal implementation
- Run test to verify it passes
- Move to next test

Step 4 - REFACTOR Phase:

- Review code for improvements
- Apply refactoring with tests as safety net
- Commit clean code

Step 5 - Repeat:

- Continue RED-GREEN-REFACTOR cycle
- Until all requirements are implemented
- Until all acceptance criteria pass

### TDD Loop Pattern

For features requiring multiple test cases:

- Identify all test cases upfront
- Prioritize by dependency and complexity
- Execute RED-GREEN-REFACTOR for each
- Maintain cumulative test suite

---

## Quality Metrics

### TDD Success Criteria

Test Coverage (Required):

- Minimum 80% coverage per commit
- 90% recommended for new code
- All public interfaces tested

Code Quality (Goals):

- All tests pass
- No test written after implementation
- Clear test names documenting behavior
- Minimal implementation satisfying tests

### TDD-Specific TRUST Validation

Apply TRUST 5 framework with TDD focus:

- Testability: Test-first approach ensures testability
- Readability: Tests document expected behavior
- Understandability: Tests serve as living documentation
- Security: Security tests written before implementation
- Transparency: Test failures provide immediate feedback

---

## Integration Points

### With DDD Workflow

TDD and DDD are complementary:

- TDD for new code
- DDD for existing code refactoring
- Hybrid mode combines both approaches

### With Testing Workflow

TDD integrates with testing workflow:

- Uses specification tests
- Integrates with coverage tools
- Supports mutation testing for test quality

### With Quality Framework

TDD outputs feed into quality assessment:

- Coverage metrics tracked
- TRUST 5 validation for changes
- Quality gates enforce standards

---

## Troubleshooting

### Common Issues

Test is Too Complex:

- Break into smaller, focused tests
- Test one behavior at a time
- Use test fixtures for complex setup

Implementation Grows Too Fast:

- Resist urge to implement untested features
- Return to RED phase for new functionality
- Keep GREEN phase minimal

Refactoring Breaks Tests:

- Revert immediately
- Refactor in smaller steps
- Ensure tests verify behavior, not implementation

### Recovery Procedures

When TDD discipline breaks down:

- Stop and assess current state
- Write characterization tests for existing code
- Resume TDD for remaining features
- Consider switching to Hybrid mode

---

Version: 1.0.0
Status: Active
Last Updated: 2026-02-03
