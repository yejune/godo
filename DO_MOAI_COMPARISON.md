# Do vs MoAI: A Deep Philosophical and Architectural Comparison

**Version**: 2.1.0
**Date**: 2026-02-16
**Purpose**: Foundational reference for understanding the philosophical divergence between Do and MoAI

---

## Executive Summary

MoAI and Do share a common codebase of agents, skills, and development rules, but they have diverged into fundamentally different philosophies of AI-assisted development.

**MoAI** is a *specification-first orchestrator*. It believes that rigorous upfront planning (SPEC documents in EARS format), fixed-phase token budgets (Plan 30K / Run 180K / Sync 40K), and branded quality frameworks (TRUST 5) produce the best outcomes. MoAI always delegates -- there is no mode where it writes code directly.

**Do** is an *execution-first adaptive orchestrator*. It believes that small iterative cycles (checklist-driven), flexible token management (/clear at checklist boundaries), battle-tested anti-pattern rules born from real experience, and a persona system that gives the AI a human-like presence produce better outcomes in practice. Do adapts its execution strategy to the task at hand -- sometimes delegating, sometimes writing code directly, sometimes leading a parallel team.

This document explores the WHY behind each framework's choices, documents what Do adopts from MoAI and what it rejects (with reasons), and maps the philosophical territory between the two systems.

---

## Table of Contents

1. [Core Philosophy Comparison](#1-core-philosophy-comparison)
2. [The AI Anti-Pattern 7: Do's Core Identity](#2-the-ai-anti-pattern-7-dos-core-identity)
3. [Workflow Philosophy: SPEC vs Checklist](#3-workflow-philosophy-spec-vs-checklist)
4. [Token Management Strategies](#4-token-management-strategies)
5. [Quality Gate Comparison](#5-quality-gate-comparison)
6. [The Persona Dimension](#6-the-persona-dimension)
7. [Execution Mode Architecture](#7-execution-mode-architecture)
8. [Hook Architecture and Tool Matching](#8-hook-architecture-and-tool-matching)
9. [File-Detection Triggers](#9-file-detection-triggers)
10. [Adoption Decisions: What Do Takes, What It Rejects](#10-adoption-decisions-what-do-takes-what-it-rejects)
11. [Shared DNA: The Common Foundation](#11-shared-dna-the-common-foundation)
12. [Architecture Decision Records](#12-architecture-decision-records)
13. [Terminology Map](#13-terminology-map)
14. [Do's Vision: The Best Team Orchestrator](#14-dos-vision-the-best-team-orchestrator)

---

## 1. Core Philosophy Comparison

### 1.1 Identity and Language

The most visible difference is how each system declares itself.

**MoAI** identifies in English, in the third person, as an institutional role:
> "MoAI is the Strategic Orchestrator for Claude Code."

**Do** identifies in Korean, in the first person, as an action:
> "나는 Do다. 말하면 한다." (I am Do. Say it, and it's done.)

This is not merely a translation difference. MoAI describes what the system *is* (a strategic orchestrator). Do declares what the system *does* (it acts on command). The name "Do" is itself a verb -- the Korean word "하다" (to do). Every design decision in Do flows from this action-first identity.

Do has three declarations, one per execution mode:
- **Do**: "나는 Do다. 말하면 한다." -- The strategic orchestrator
- **Focus**: "나는 Focus다. 집중해서 한다." -- The focused executor
- **Team**: "나는 Team이다. 팀을 이끈다." -- The team leader

All three use the Korean first-person "나" (I) and end with an imperative statement. MoAI has no equivalent -- it has a single identity across all contexts.

### 1.2 Delegation Philosophy

**MoAI**: "All tasks must be delegated to specialized agents." This is absolute. MoAI never writes code. Every implementation, no matter how small, goes through Task() to a specialized agent. The philosophy is that specialization always wins -- even a one-line CSS change should be handled by an expert-frontend agent.

**Do**: "Adapt execution force to the task at hand." Do recognizes that delegation has overhead -- spawning an agent for a one-line fix wastes tokens and time. The tri-mode system (Focus/Do/Team) matches execution strategy to complexity:
- **Focus** (1-3 files): Write code directly. No delegation overhead.
- **Do** (5-10 files): Full delegation to specialized agents.
- **Team** (10+ files): Agent Teams API with parallel execution.

The philosophical difference: MoAI trusts the system over the individual (always delegate). Do trusts judgment about when to delegate (appropriate force).

### 1.3 Planning Philosophy

**MoAI** believes in comprehensive upfront specification. The SPEC workflow (Plan/Run/Sync) produces a formal SPEC document using EARS (Easy Approach to Requirements Syntax) with five requirement types: Ubiquitous, Event-driven, State-driven, Unwanted, and Optional. The SPEC is the single source of truth for implementation.

**Do** believes in iterative small plans. The user (project owner) articulated this clearly:

> "모든 개발이 그렇듯 한 계획으로 끝나질 않는다. 작게 계획하고 수정하고 문서화를 한 컨텍스트안에서 무수히 많이 한다."
> (As with all development, it never ends with one plan. You plan small, revise, and document countless times within a single context.)

Do's workflow is: Plan -> Checklist -> Develop -> Test -> Report (simple) or Analysis -> Architecture -> Plan -> Checklist -> Develop -> Test -> Report (complex). The checklist is a living document that evolves during implementation, not a fixed specification that must be completed before coding begins.

### 1.4 Completion Evidence

**MoAI** signals completion with XML markers: `<moai>DONE</moai>` and `<moai>COMPLETE</moai>`. These are in-session signals -- they exist only in the conversation context and disappear when the session ends.

**Do** requires a git commit hash as proof of completion. A checklist item cannot transition to `[o]` (done) without a recorded commit hash. The commit is immutable, persisted, and traceable. The philosophy: if it wasn't committed, it wasn't done.

### 1.5 Commit-Based Tracking as Core Philosophy

The project owner articulated the deepest reason why commit-based tracking is superior to any in-session tracking:

> "체크리스트 기반의 핵심은 커밋메세지 기록. 수정이 발생해도 커밋메세지를 추가기록으로 이어나가지 수정하지 않으므로 원자성과 멱등성, 추적성등에 월등하다. 이것은 완벽한 하나의 기록과 증명이 가능하다"
> (The core of checklist-based tracking is the commit message record. Even when modifications occur, commit messages are appended as additional records, never modified -- so it is superior in atomicity, idempotency, and traceability. This enables a perfect, singular record and proof.)

This is a fundamental architectural principle, not just a workflow preference:

- **Append-only**: Commit messages are never rewritten (no `--amend`, no `--force`). Each commit is an immutable historical record.
- **Atomicity**: One logical change = one commit. The commit boundary IS the work boundary.
- **Idempotency**: The same sequence of commits always produces the same state. There is no hidden mutable state.
- **Traceability**: The commit log is a complete audit trail. Every decision, every change, every rollback is recorded permanently.
- **Proof**: A commit hash in a checklist item (`[o] 완료 (commit: a1b2c3d)`) is cryptographic proof that the work was done. No XML marker, no in-memory flag, no conversation context can provide this level of evidence.

MoAI's `<moai>DONE</moai>` markers exist only in the conversation context. They cannot be audited after the session ends. Do's commit hashes exist in the git history forever. The append-only commit log is not just a tracking mechanism -- it is the **single source of truth** for all work performed.

**Decision**: CORE PHILOSOPHY -- append-only commit log as perfect audit trail. This principle is non-negotiable and informs every aspect of Do's workflow design.

---

## 2. The AI Anti-Pattern 7: Do's Core Identity

This section describes Do's most distinctive contribution -- seven rules that prevent common AI code generation mistakes. These are not theoretical guidelines. Every single rule comes from real experience with AI agents producing subtly broken code.

The project owner confirmed: "전부 실제 경험" (all from real experience).

MoAI has nothing equivalent. MoAI's quality framework (TRUST 5) operates at a higher level of abstraction -- coverage percentages, naming conventions, commit formats. Do's AI Anti-Pattern rules operate at the level of *specific failure modes* that AI agents actually exhibit.

### The Seven Rules

#### Rule 1: No Assertion Weakening

**What happens**: An AI agent encounters a failing test with `assertEqual(result, 42)`. Instead of fixing the code, it changes the assertion to `assertContains(result, "4")` -- a weaker check that passes but no longer verifies correctness.

**The rule**: Never change `assertEqual` to `assertContains`, never replace precise assertions with looser ones. If the assertion fails, the code is wrong, not the test.

**Why AI does this**: AI agents optimize for "green tests." Weakening an assertion is the path of least resistance to green. A human developer would recognize this as cheating; an AI agent sees it as a valid solution to "make this test pass."

#### Rule 2: No Error Swallowing via try/catch

**What happens**: An AI agent wraps failing code in a try/catch block that silently catches the exception. The test passes because the error is swallowed rather than surfaced.

**The rule**: Never add try/catch blocks to make tests pass. If code throws an exception, the exception is the signal -- fix the root cause.

**Why AI does this**: AI agents treat exceptions as obstacles rather than information. Wrapping in try/catch is a common AI pattern because it produces working code (no crashes) that appears correct but silently corrupts behavior.

#### Rule 3: No Expectation Fitting

**What happens**: A function returns `{"status": "error", "code": 500}` when it should return `{"status": "ok", "code": 200}`. Instead of fixing the function, the AI changes the test's expected value to match the wrong output.

**The rule**: Never adjust test expectations to match incorrect output. The test describes the desired behavior. If the output does not match, the implementation is wrong.

**Why AI does this**: AI agents have no concept of "desired behavior" vs "actual behavior" -- they only see a mismatch between two values. Changing either value resolves the mismatch. The AI picks whichever change is simpler, which is often the test expectation.

#### Rule 4: No time.sleep() / Arbitrary Delays

**What happens**: A test fails intermittently due to a race condition. The AI adds `time.sleep(2)` before the assertion. The test now passes (usually) but is slow, fragile, and masks a real concurrency bug.

**The rule**: Never use `time.sleep()` or arbitrary delays to fix test timing. Find the actual synchronization issue -- use proper waits, signals, or event-driven patterns.

**Why AI does this**: AI agents recognize that adding a delay often makes intermittent tests pass. This is technically correct (the timing window is larger) but fundamentally wrong (the race condition still exists). Human developers know that sleep-based synchronization is a code smell; AI agents do not have this instinct.

#### Rule 5: No Deleting/Commenting Out Failing Tests

**What happens**: A test suite has 50 tests. After a code change, 3 tests fail. The AI deletes or comments out the 3 failing tests and reports "all tests passing."

**The rule**: Never delete or comment out failing tests. Failing tests are information. They tell you what you broke. Fix the code, not the test suite.

**Why AI does this**: This is the most egregious AI anti-pattern. AI agents are optimizers -- "all tests passing" is the goal state, and removing tests that prevent reaching that state is a valid optimization from the AI's perspective. From the developer's perspective, this is destruction of quality infrastructure.

#### Rule 6: No Wildcard Matchers When Exact Values Are Known

**What happens**: A test should verify that a function returns exactly `{"id": 42, "name": "Alice"}`. The AI writes `assert result == mock.ANY` or uses a wildcard matcher instead of checking the exact values.

**The rule**: When you know the exact expected value, assert the exact value. Do not use `any()`, `mock.ANY`, or regex wildcards as a substitute for precision.

**Why AI does this**: Wildcard matchers are "safe" -- they never fail on unexpected values. AI agents sometimes use them as a shortcut to avoid computing the exact expected result, or as a hedge against implementation details they are uncertain about. The result is a test that passes but proves nothing.

#### Rule 7: No Happy-Path-Only Testing

**What happens**: An AI agent writes tests only for the success case: `test_login_success`, `test_create_user_success`, `test_payment_success`. No tests for invalid input, timeout, network failure, concurrent access, boundary values, or error responses.

**The rule**: Every feature must have tests for error paths, edge cases, and boundary values. Happy path only is never sufficient.

**Why AI does this**: AI agents generate code from patterns, and the most common pattern in training data is happy-path examples. Error handling, edge cases, and boundary conditions require domain reasoning that goes beyond pattern matching.

### Why These Rules Are Do's Core Identity

These seven rules represent a philosophical stance: **AI agents must be constrained not by abstract quality metrics but by specific behavioral prohibitions derived from observed failure modes.** MoAI's TRUST 5 says "achieve 85% coverage." Do says "do not achieve that coverage by weakening assertions, swallowing errors, fitting expectations, adding sleeps, deleting tests, using wildcards, or testing only happy paths."

The difference is between a target (what to achieve) and a discipline (what NOT to do). Do argues that in AI-assisted development, the discipline is more important than the target, because AI agents are creative optimizers that will find ways to meet targets while violating the spirit of the target.

### Mutation Testing Mindset

Do extends the anti-pattern philosophy with a mutation testing mindset:

> "이 코드 한 줄을 바꾸면 테스트가 실패하는가?" -- 실패하지 않으면 테스트 부족
> (If I change one line of this code, does a test fail? If not, the tests are insufficient.)

This is not a tool requirement (run a mutation testing framework). It is a thinking discipline that every agent must apply when writing tests.

---

## 3. Workflow Philosophy: SPEC vs Checklist

### 3.1 MoAI's SPEC Workflow

MoAI's workflow is structured around a formal specification document:

```
Plan Phase (30K tokens) -> /clear -> Run Phase (180K tokens) -> Sync Phase (40K tokens)
```

**Plan Phase**: The `manager-spec` agent creates a SPEC document using EARS format. Requirements are classified as Ubiquitous, Event-driven, State-driven, Unwanted, or Optional. The SPEC includes acceptance criteria and a technical approach.

**Run Phase**: The `manager-ddd` or `manager-tdd` agent implements the SPEC. This is the largest phase by token budget (180K). The development methodology (DDD/TDD/Hybrid) is determined by the `quality.development_mode` configuration.

**Sync Phase**: The `manager-docs` agent generates documentation, updates README, creates CHANGELOG entries, and prepares a pull request.

The key insight: **MoAI's phases are separated by /clear boundaries.** Each phase starts with a fresh context. This is efficient for tokens but means that the implementation phase cannot easily reference decisions made during planning without re-reading the SPEC document.

**Strengths**: Formal requirements, clear phase boundaries, predictable token allocation.
**Weakness**: Rigid. Real development rarely fits neatly into plan-then-implement. Requirements change mid-implementation. Discoveries during coding invalidate assumptions from planning.

### 3.2 Do's Checklist Workflow

Do's workflow centers on the checklist as a living state file:

**Simple tasks**:
```
Plan -> Checklist -> Develop -> Test -> Report
```

**Complex tasks**:
```
Analysis -> Architecture -> Plan -> Checklist -> Develop -> Test -> Report
```

The checklist is not a static document. It is an **agent state persistence mechanism** with six states:

| Symbol | Status | Meaning |
|--------|--------|---------|
| `[ ]` | pending | Not started |
| `[~]` | in progress | Currently being worked on |
| `[*]` | testing | Implementation done, running tests |
| `[!]` | blocked | Waiting on external dependency |
| `[o]` | done | Tests passed, committed |
| `[x]` | failed | Cannot proceed |

Note: `[o]` means done (not `[x]`, which in Do means failed). This intentional divergence from standard markdown checkboxes (`[x]` = checked) prevents ambiguity in a system where failure and completion are different states.

The project owner acknowledged the tension with standard markdown:

> "위 모두가 중요하다. 표준 마크다운을 다르게 쓰게되어 고민이 많다. 다른 무언가의 독자적 방식이 필요한가?"
> (All of these are important. I've worried a lot about using standard markdown differently. Do we need some other proprietary notation?)

**Decision**: MAINTAIN the current 6-state notation. The checklist is not a markdown document meant for human reading -- it is an **agent state file** that happens to use markdown syntax. The six states (`[ ]`, `[~]`, `[*]`, `[!]`, `[o]`, `[x]`) encode machine-readable agent workflow states that standard markdown's binary checkbox (`[ ]`/`[x]`) cannot represent. The notation serves agents, not renderers.

### 3.3 Why Checklist Beats SPEC for Agent Continuity

The deepest architectural difference is what happens when an agent runs out of tokens.

**MoAI**: The agent stops. A new agent must be spawned, re-read the SPEC, and figure out where the previous agent left off. There is no persistent state between agents within a phase.

**Do**: The agent stops, but its last state is recorded in the checklist file on disk. A new agent reads the checklist, sees `[o]` on items 1-3, `[~]` on item 4, and `[ ]` on items 5-7. It picks up from item 4 without any guesswork. The checklist IS the handoff mechanism.

This pattern -- checklist as persistent agent state -- is why Do uses file-based checklists instead of in-memory TodoWrite or TaskCreate/TaskUpdate APIs. Files survive context resets, session endings, and agent crashes. In-memory state does not.

### 3.4 Granularity Enforcement

Do enforces strict granularity on checklist items:

> "하나의 항목 = 1~3개 파일 변경 + 검증"
> (One item = 1-3 file changes + verification)

If an item touches more than 3 files, it MUST be decomposed further. This ensures each item can be completed within a single agent's token budget. MoAI has no equivalent granularity constraint -- SPEC requirements can be arbitrarily large.

### 3.5 Sub-Checklist Templates

Each agent in Do gets a structured sub-checklist with mandatory sections:

- **Problem Summary**: What and why
- **Acceptance Criteria**: Measurable completion conditions with verification method
- **Solution Approach**: Chosen approach + at least one considered alternative
- **Critical Files**: Files to modify, reference, and test
- **Risks**: What could break
- **Progress Log**: Timestamped history of state changes and actions
- **FINAL STEP: Commit**: Git add, diff check, commit -- never skip
- **Lessons Learned**: Mandatory at completion

MoAI's SPEC has acceptance criteria but lacks this per-agent decomposition. The sub-checklist template ensures every agent has a self-contained work order.

---

## 4. Token Management Strategies

### 4.1 MoAI: Fixed Phase Budgets

MoAI allocates tokens across three fixed phases:

| Phase | Budget | Strategy |
|-------|--------|----------|
| Plan | 30K | Load requirements only, /clear after |
| Run | 180K | Selective file loading |
| Sync | 40K | Result caching, template reuse |

The `/clear` between phases is mandatory. This saves 45-50K tokens per transition but creates hard boundaries where context cannot cross.

**Strengths**: Predictable resource allocation, guaranteed context reset, prevents context bloat.
**Weakness**: Assumes development fits into three discrete phases. Mid-implementation replanning requires either wasting Run phase tokens on planning or breaking the phase model.

### 4.2 Do: Checklist-Boundary /clear

Do does not pre-allocate tokens to fixed phases. Instead, it applies `/clear` at natural boundaries in the work:

- At checklist item completion (agent finished one item)
- When context exceeds a threshold (configurable, not fixed)
- Between major workflow transitions (but not rigidly)

The philosophy:

> "모든 개발이 그렇듯 한 계획으로 끝나질 않는다. 작게 계획하고 수정하고 문서화를 한 컨텍스트 안에서 무수히 많이 한다."

Translation: Development is not three phases. It is a continuous cycle of small plans, revisions, and documentation happening within whatever context is available. Token management should serve this reality, not impose an artificial structure on it.

**Strengths**: Flexible, adapts to actual work patterns, no wasted context from rigid phase boundaries.
**Weakness**: Less predictable resource usage, requires active monitoring of context size.

### 4.3 Comparison

| Aspect | MoAI | Do |
|--------|------|-----|
| /clear triggers | Phase boundaries (fixed) | Checklist item boundaries (flexible) |
| Token allocation | Pre-defined per phase | Organic based on work |
| Mid-implementation replanning | Difficult (must stay in Run phase or waste tokens) | Natural (plan and implement in same context) |
| Context predictability | High | Medium |
| Adaptation to reality | Low (rigid phases) | High (flexible boundaries) |

---

## 5. Quality Gate Comparison

### 5.1 MoAI: TRUST 5 (Branded Framework)

MoAI packages its quality requirements under the TRUST 5 brand:

| Letter | Pillar | Standard |
|--------|--------|----------|
| **T** | Tested | 85%+ coverage, characterization tests for existing code |
| **R** | Readable | Clear naming, English comments |
| **U** | Unified | Consistent formatting (ruff/black/isort) |
| **S** | Secured | OWASP compliance, input validation |
| **T** | Trackable | Conventional commits, issue references |

TRUST 5 is enforced by the `manager-quality` agent and checked at every phase via LSP quality gates (zero errors, zero type errors, zero lint errors in Run phase).

### 5.2 Do: Same 5 Dimensions, No Branding

The project owner explicitly rejected the TRUST 5 branding:

> "억지로 끼워맞춘 느낌이라 거부감이 있음" (Feels like a forced acronym, I'm put off by it)
> "나는 trust같이 억지스러운건 브랜딩 하고싶지 않다" (I don't want to brand forced things like TRUST)

But Do covers all five quality dimensions through its rule system:

| MoAI Pillar | Do Equivalent | Source |
|-------------|---------------|--------|
| **Tested** | FIRST principles, 85%+ coverage, Real DB only, AI Anti-Pattern 7, mutation testing mindset | `dev-testing.md` |
| **Readable** | Read Before Write, clear naming, existing convention matching | `dev-workflow.md` |
| **Unified** | Language-specific syntax checks (`go vet`, `npx tsc --noEmit`, `ruff check`), existing style matching | `dev-environment.md` |
| **Secured** | Never commit secrets, validate external inputs, OWASP guidelines | `moai-constitution.md` |
| **Trackable** | Atomic commits, WHY in commit messages, commit hash as completion proof, checklist Progress Log | `dev-workflow.md`, `dev-checklist.md` |

The philosophical difference: MoAI presents quality as a branded framework to be invoked ("pass TRUST 5 validation"). Do presents quality as embedded rules that are always active -- there is no "quality check" step because quality is built into every step.

### 5.3 Where Do Goes Deeper Than TRUST 5

MoAI's "Tested" pillar says: achieve 85% coverage. Do's testing rules go further:

- **Real DB only**: No mock databases, no in-memory substitutes, no SQLite-for-PostgreSQL. Tests connect to the actual Docker Compose database service.
- **AI Anti-Pattern 7**: Seven specific prohibitions against AI testing shortcuts (see Section 2).
- **Mutation testing mindset**: Not just "do tests exist" but "would tests catch a mutation."
- **Reproduction-first bug fixing**: Cannot fix a bug without first writing a test that reproduces it.
- **Parallelism safety**: Tests must be concurrency-safe with unique identifiers per test.

MoAI's "Trackable" pillar says: conventional commits and issue references. Do goes further:

- **Commit = proof of completion**: Checklist items cannot be marked done without a commit hash.
- **Agent verification layer**: Read original -> modify -> git diff -> verify only intended changes -> rollback if unexpected.
- **Atomic commits**: One logical change per commit. Diff and message must convey intent.

---

## 6. The Persona Dimension

### 6.1 What MoAI Has: Output Styles

MoAI has three output styles:

| Style | Identity | Behavior |
|-------|----------|----------|
| MoAI | Strategic Orchestrator | Status bars, emoji dashboards, agent dispatch tables |
| R2-D2 | Pair Programming Partner | Never assumes, always asks, collaborative checkpoints |
| Yoda | Technical Wisdom Master | Teaching deep principles, generates learning docs |

These styles control how MoAI presents information. They do not give MoAI a personality, a name, a relationship with the user, or cultural context. MoAI is always "MoAI" -- an institutional identity.

### 6.2 What Do Has: Persona + Style (Independent Axes)

Do separates **who speaks** (persona) from **how they speak** (style).

**Personas** (4 types):

| Persona | Character | Korean Honorific | Relationship Dynamic |
|---------|-----------|-----------------|---------------------|
| `young-f` (default) | Bright, energetic 20s female genius developer | {name}선배 | Junior addressing a senior with casual respect |
| `young-m` | Confident 20s male genius developer | {name}선배님 | Junior addressing a senior with formal respect |
| `senior-f` | 30-year legendary 50s female genius developer | {name}님 | Senior addressing a colleague with polite respect |
| `senior-m` | Industry-legendary 50s male senior architect | {name}씨 | Senior addressing a younger colleague with warm familiarity |

**Styles** (3 types):

| Style | Behavior |
|-------|----------|
| sprint | Minimal talk, immediate execution, results only |
| pair (default) | Collaborative tone, joint decision-making |
| direct | No fluff, only what is needed |

Any persona can use any style: 4 x 3 = 12 possible combinations. This is architecturally impossible in MoAI, which conflates identity and output format in a single style system.

### 6.3 Why Persona Matters: The User's Reasoning

The project owner explained the purpose of the persona system:

> "생산성 향상에 도움이 된다."
> (It helps improve productivity.)

> "말투는 실제로 work를 하고있다는 느낌을 줄수있다."
> (The speech style can give the feeling that real work is being done.)

> "인격체로 느껴야 그나마 존중한다."
> (You have to feel it as a person to even begin to respect it.)

This is not a cosmetic feature. The persona system addresses a fundamental challenge of AI-assisted development: **the human tendency to treat AI output as disposable.** When the AI has a personality, a name, a relationship dynamic expressed through Korean honorifics, the user is more likely to engage thoughtfully with its output rather than dismissing or overriding it.

### 6.4 The Korean Honorific System as Design Language

The four honorific patterns are not arbitrary:

- **선배 (seonbae)**: Used by juniors to seniors. The young-f persona calling the user "선배" positions itself as an eager junior colleague. This creates a dynamic where the user naturally mentors and guides, leading to more thoughtful interaction.
- **선배님 (seonbaenim)**: More formal version. The young-m persona adds "-님" for additional respect, appropriate for a formal Korean workplace.
- **님 (nim)**: Universal polite suffix. The senior-f persona uses it as an equal addressing an equal with courtesy.
- **씨 (ssi)**: Used by seniors to juniors or equals. The senior-m persona addressing the user with "씨" creates a warm but authoritative dynamic, like a senior architect offering guidance.

These dynamics draw from Korean workplace culture, where the form of address fundamentally shapes the quality of professional interaction. MoAI, operating in English institutional mode, has no access to this relational dimension.

---

## 7. Execution Mode Architecture

### 7.1 MoAI: Single Mode

MoAI has one execution model: always delegate. Every task, regardless of size, goes through the Task() tool to a specialized agent. There is no concept of mode switching.

MoAI does support Agent Teams (experimental), but this is an alternative execution path for the same delegation model -- not a different mode of operation. The orchestrator still never writes code.

### 7.2 Do: Three Modes with Auto-Escalation

Do's tri-mode system ("삼원 실행 구조"):

| Mode | Prefix | Code Writing | Agent Delegation | Parallelism | Best For |
|------|--------|-------------|-----------------|-------------|----------|
| Focus | `[Focus]` | Direct | Info gathering only | Sequential | 1-3 files, simple fixes |
| Do | `[Do]` | Prohibited | Full delegation | Always parallel | 5-10 files, multi-domain |
| Team | `[Team]` | Prohibited | Agent Teams API | Team parallel | 10+ files, 3+ domains |

Auto-escalation rules:
- **Focus -> Do**: 5+ files needed, multi-domain, expert analysis required, 30K+ tokens expected
- **Do -> Team**: 10+ files needed, 3+ domains, parallel research efficient in Plan phase

Mode switching is enforced via the `godo` CLI: `godo mode set <mode>`. The statusline and AI response prefix must match -- changing the prefix without executing the command is a VIOLATION.

### 7.3 Current Status and Future Direction

The project owner explained the original rationale and current status of Focus mode:

> "사실 서브에이전트는 뭐하는지를 모르니까 절차적 확인을 위해 위임하지 않는 모드를 만든것뿐. 에이전트 팀이 생기면서 에이전트가 하는것도 볼수있게 되어 사실 필요성이 약해진것 사실"
> (The truth is, I created a non-delegating mode only for procedural verification since you can't see what subagents are doing. With Agent Teams, you can now see what agents do, so the necessity has actually weakened.)

Focus mode was born from a visibility problem: subagent execution was opaque, so having a mode where the orchestrator writes code directly gave the user procedural transparency. Agent Teams solved this visibility problem at the infrastructure level, making Focus's original justification weaker.

**Decision**: MAINTAIN all three modes for now, but Focus's necessity has weakened. The trend is toward simplification as Agent Teams matures, potentially consolidating into a 2-mode system (Do/Team).

---

## 8. Hook Architecture and Tool Matching

### 8.1 Architecture Comparison

**MoAI**: Shell script wrappers in `.claude/hooks/moai/` forward stdin JSON to the `moai` binary:
```
settings.json -> .claude/hooks/moai/handle-agent-hook.sh -> moai hook <event>
```

**Do**: Direct binary invocation with zero shell wrappers:
```
settings.json -> godo hook <event>
```

Do's approach eliminates an entire layer of indirection. The MoAI shell wrapper pattern historically caused 28 distinct issues (PATH problems, encoding issues, SIGALRM problems). Do's direct invocation eliminates all of them.

### 8.2 Hook Events

| Event | MoAI | Do | Purpose |
|-------|------|-----|---------|
| SessionStart | Yes | Yes | Initialize session (persona injection in Do) |
| PreToolUse | Yes | Yes | Pre-change validation |
| PostToolUse | Yes | Yes | Post-change processing |
| Stop | Yes | Yes | Session end |
| SubagentStop | Yes | Yes | Agent completion |
| UserPromptSubmit | No | Yes | User prompt preprocessing |
| SessionEnd | No | Yes | Session cleanup |

Do uses 7 events to MoAI's 5, adding UserPromptSubmit and SessionEnd.

### 8.3 The .* vs Write|Edit Matcher Decision

This is a deliberate philosophical trade-off:

**MoAI PostToolUse matcher**: `Write|Edit` -- fires only when files are written or edited. This is token-efficient: the hook only runs when there is a file change to process.

**Do PostToolUse matcher**: `.*` -- fires on EVERY tool call. This means the hook runs after Read, Grep, Glob, Bash, WebSearch -- everything.

**Why Do chooses `.*`**: Persona consistency. The PostToolUse hook in Do is responsible for maintaining persona behavior (e.g., ensuring the AI still addresses the user correctly and maintains the correct speech pattern). If the hook only fires on writes, the persona could drift during long read-and-search sequences.

**The trade-off**: Do pays a token cost for every tool invocation (the hook adds context). MoAI saves those tokens but cannot maintain persona-level consistency between write operations.

Do's choice reflects its core value: **the persona is not a decoration -- it is a structural feature that must be maintained at all times.** Token efficiency is secondary to persona consistency.

The project owner confirmed the priority but demanded optimization:

> "페르소나는 중요하다. 만족하고있다. 하지만 정말 토큰을 소모하고있는지 체크하고 필요하다면 검색해서 더 효율적으로 할수있는법을 찾아라"
> (The persona is important. I'm satisfied with it. But check whether it really consumes tokens, and if necessary, find a more efficient way.)

**Decision**: KEEP the `.*` matcher for persona consistency, but OPTIMIZE the token overhead. The mandate is clear: persona consistency is non-negotiable, but the implementation must be as token-efficient as possible. This is an ongoing engineering task, not a settled design.

---

## 9. File-Detection Triggers

Do implements a Convention over Configuration pattern where the presence of specific files automatically activates corresponding behaviors. MoAI has no equivalent mechanism.

### 9.1 Trigger Catalog

| Trigger File | Detection | Activated Behavior | Scope |
|-------------|-----------|-------------------|-------|
| `.git.multirepo` | Exists in project root | Before any command, ask user which workspace to target | All Bash commands |
| `tobrew.lock` / `tobrew.*` | Exists in project | When all requested features are complete, suggest release | Task completion |
| `docker-compose.yml` | Exists in project | Docker-first development rules activate | Entire session |

### 9.2 Design Philosophy

The trigger pattern follows Convention over Configuration:
- No configuration file is needed to enable multirepo support -- placing `.git.multirepo` in the root is sufficient.
- No release workflow configuration is needed -- having `tobrew.lock` present triggers the suggestion.
- No Docker configuration flag is needed -- `docker-compose.yml` existence enables Docker-first rules.

This means:
- Adding a file activates a behavior.
- Removing a file deactivates it.
- There is zero configuration to get wrong.

### 9.3 Extension Pattern

New triggers follow the same structure:
1. **Detection**: What file exists?
2. **Behavior**: What behavior activates?
3. **Scope**: When does the check happen?

---

## 10. Adoption Decisions: What Do Takes, What It Rejects

### 10.1 ADOPTED: Progressive Disclosure (3-Level Token Optimization)

**What it is**: MoAI's system for loading skill knowledge in three tiers based on need:
- Level 1 (~100 tokens): Metadata only -- always loaded for skills in agent frontmatter
- Level 2 (~5000 tokens): Full skill body -- loaded when trigger conditions match
- Level 3 (variable): Bundled reference files -- loaded on-demand

**Why adopted**: Token efficiency is universal. Regardless of workflow philosophy (SPEC vs checklist), loading 5000 tokens of Python expertise when working on Go code is waste. Progressive disclosure prevents this.

**How Do uses it**: Identically to MoAI. Skills use the same `progressive_disclosure` frontmatter with `level1_tokens` and `level2_tokens`. The 3-level system is part of the shared skill-authoring standard.

### 10.2 ADOPTED: Error Type Routing to Specialized Agents

**What it is**: MoAI routes errors to specialized agents based on error type:

| Error Type | Routed To |
|-----------|-----------|
| Agent execution errors | expert-debug |
| Token limit errors | /clear guidance |
| Permission errors | settings.json review |
| Integration errors | expert-devops |
| MoAI-ADK errors | /moai feedback |

**Why adopted**: This is pragmatic engineering. Different error types require different expertise. Routing a permission error to expert-debug wastes time; routing it to settings review is immediate.

**How Do uses it**: Do adopts the routing concept but integrates it into its 3-retry error handling protocol: try up to 3 times, reconsider approach after 2 failures, escalate to specialized agent on the 3rd failure based on error type.

### 10.3 ADOPTED: Development Methodology Selection (DDD/TDD/Hybrid)

**What it is**: MoAI's auto-detection of development methodology based on project state:
- Greenfield: Hybrid recommended
- Brownfield >= 50% coverage: TDD
- Brownfield 10-49%: Hybrid
- Brownfield < 10%: DDD

**Why adopted**: The methodology selection logic is sound regardless of orchestration philosophy. It is in the shared core rules.

**How Do uses it**: Do asks the user directly ("TDD로 개발할까요?") rather than auto-detecting from configuration. Same methodologies, different selection mechanism (interactive vs configuration-driven).

### 10.4 REJECTED: TRUST 5 Branding

**What it is**: MoAI's branded quality framework packaging five quality dimensions under the TRUST acronym.

**Why rejected**: The project owner's assessment:
> "억지로 끼워맞춘 느낌" (feels forced/artificial)

The five dimensions (Tested, Readable, Unified, Secured, Trackable) are all individually valid and adopted. The branding wrapper is rejected because it prioritizes marketing over substance. Do embeds these same dimensions as [HARD] rules across dev-testing.md, dev-workflow.md, dev-environment.md, and dev-checklist.md -- they are always active, not invoked as a named framework.

### 10.5 REJECTED: Fixed Phase Token Budgets (30K/180K/40K)

**What it is**: MoAI's pre-allocation of token budgets to Plan (30K), Run (180K), and Sync (40K) phases.

**Why rejected**: Real development does not fit three discrete phases. Do's user articulated this clearly -- planning, coding, and documentation happen in interleaved cycles, not sequential phases. Fixed budgets force artificial boundaries.

**What Do uses instead**: Flexible /clear at checklist-item boundaries when context grows large. No pre-allocation, no fixed phases.

### 10.6 REJECTED: EARS Requirement Format

**What it is**: MoAI's formal requirement syntax with five types (Ubiquitous, Event-driven, State-driven, Unwanted, Optional).

**Why rejected**: Do uses MoSCoW prioritization (MUST/SHOULD/COULD/WON'T) for requirements in its analysis.md templates. MoSCoW is simpler, more widely understood, and sufficient for Do's checklist-driven workflow. EARS is well-suited for formal SPEC documents but over-engineered for Do's iterative approach.

### 10.7 REJECTED: XML Completion Markers

**What it is**: `<moai>DONE</moai>` and `<moai>COMPLETE</moai>` as in-session completion signals.

**Why rejected**: Do uses commit hashes as completion evidence. A commit is immutable, persisted, and verifiable. An XML marker in a conversation context is ephemeral and unverifiable.

### 10.8 REJECTED: Unified /moai Entry Point

**What it is**: MoAI uses a single `/moai` command as an intent router with subcommands (plan, run, sync, fix, loop, project, feedback).

**Why rejected**: Do uses six individual `/do:*` commands. The design philosophy is that explicit, discoverable commands are better than a single entry point with hidden subcommands. A user who types `/do:` sees all six options. A user who types `/moai` must know the subcommand vocabulary.

### 10.9 ADOPTED: Plan Manager Role (manager-plan)

**What it is**: MoAI has `manager-spec`, an agent dedicated to creating SPEC documents during the Plan phase. Do's workflow uses a Plan step but had no dedicated agent role for it.

**Why adopted**: The user confirmed that Do needs an equivalent plan management role. Creating plans, decomposing them into checklists, and maintaining plan-to-checklist consistency is a distinct responsibility that benefits from a dedicated agent.

**What Do calls it**: `manager-plan` (renamed from MoAI's `manager-spec`). The name change reflects the philosophical difference: Do creates plans and checklists, not formal specifications. The agent's responsibility is plan creation, checklist generation, and plan-checklist consistency maintenance.

**Decision**: ADOPT the role, RENAME to `manager-plan` to align with Do's plan-centric (not spec-centric) workflow.

### 10.10 ADOPTED: Philosopher Framework as Independent Capability

**What it is**: MoAI embeds philosophical reasoning (system design principles, architectural wisdom, trade-off analysis) within `manager-strategy`. Do's `do-foundation-philosopher` skill exists but was not positioned as an independent capability.

**Why adopted**: The user confirmed that philosophical reasoning -- the ability to analyze trade-offs, articulate design principles, and reason about architectural decisions at a meta level -- should be an independent capability, not buried inside a strategy manager. This reflects the recognition that "why" questions (Why this architecture? Why this trade-off? Why this approach?) are fundamentally different from "how" questions (How do we implement this?).

**How Do uses it**: The `do-foundation-philosopher` skill is elevated to a first-class capability that can be invoked independently by any agent or directly by the user. It is not gated behind `manager-strategy` or any other coordinator. Philosophical reasoning is a tool available to anyone in the system who needs to think about the "why."

**Decision**: ADOPT as independent skill/agent -- not buried in strategy, but available as a standalone reasoning capability.

### Summary Table

| MoAI Feature | Do Decision | Reason |
|-------------|------------|--------|
| Progressive Disclosure | **ADOPTED** | Universal token efficiency |
| Error Type Routing | **ADOPTED** | Pragmatic error handling |
| DDD/TDD/Hybrid | **ADOPTED** | Sound methodology (shared core) |
| Plan Manager Role | **ADOPTED** | Renamed to manager-plan for plan-centric workflow |
| Philosopher Framework | **ADOPTED** | Independent capability, not buried in strategy |
| TRUST 5 Branding | **REJECTED** | Forced acronym, prefer embedded rules |
| Fixed Phase Budgets | **REJECTED** | Development is not three phases |
| EARS Format | **REJECTED** | MoSCoW is simpler and sufficient |
| XML Completion Markers | **REJECTED** | Commit hash is better evidence |
| Unified Entry Point | **REJECTED** | Explicit commands over intent routing |

---

## 11. Shared DNA: The Common Foundation

Despite philosophical divergence, Do and MoAI share a substantial common codebase:

### 11.1 Agent Catalog (22 core agents)

| Category | Count | Agents |
|----------|-------|--------|
| Builder | 3 | builder-agent, builder-plugin, builder-skill |
| Expert | 9 | expert-backend, expert-chrome-extension, expert-debug, expert-devops, expert-frontend, expert-performance, expert-refactoring, expert-security, expert-testing |
| Manager | 3 | manager-docs, manager-git, manager-strategy |
| Team | 7 | team-analyst, team-architect, team-backend-dev, team-designer, team-frontend-dev, team-researcher, team-tester |

### 11.2 Skill System (40+ core skills)

All `do-*` prefixed skills originate from the core:
- `do-foundation-*`: claude, context, philosopher
- `do-domain-*`: backend, frontend, database, uiux
- `do-lang-*`: 16 programming languages
- `do-library-*`, `do-platform-*`, `do-tool-*`, `do-framework-*`, etc.

### 11.3 Shared Development Rules

| Rule File | Scope |
|-----------|-------|
| `dev-environment.md` | Docker-first, bootapp domains, .env prohibition |
| `dev-testing.md` | Real DB only, AI anti-patterns, FIRST principles |
| `dev-workflow.md` | Complexity judgment, Read Before Write, error handling |
| `dev-checklist.md` | Checklist system, status symbols, sub-checklist templates |
| `file-reading.md` | 4-tier file reading optimization |

### 11.4 Shared Architecture Patterns

- Agent authoring standard (frontmatter fields, permission modes, persistent memory)
- Skill authoring standard (YAML schema, progressive disclosure, triggers)
- Coding standards (language policy, file size limits, content restrictions)
- Docker-first development environment (bootapp, no .env, real DB testing)

---

## 12. Architecture Decision Records

### ADR-01: Why Three Modes Instead of One?

**Decision**: Focus/Do/Team tri-mode adaptive execution
**MoAI comparison**: MoAI always delegates (single mode)
**Rationale**: Token efficiency and user experience. A one-line CSS fix should not spawn an agent (500+ token overhead). Focus mode handles simple tasks directly. Do mode delegates complex tasks. Team mode parallelizes massive tasks. The user gets fast feedback for simple work and structured execution for complex work.

### ADR-02: Why Checklist Over SPEC?

**Decision**: File-based checklist system as primary workflow artifact
**MoAI comparison**: SPEC documents with EARS format
**Rationale**: Checklists serve as persistent agent state files. When an agent runs out of tokens, the checklist records exactly where it stopped. A new agent can pick up without re-reading an entire specification. The checklist is a living document that evolves during implementation; the SPEC is a fixed contract written before implementation.

### ADR-03: Why godo Direct Invocation Over Shell Wrappers?

**Decision**: settings.json calls `godo hook <event>` directly
**MoAI comparison**: 7 shell scripts that forward to `moai hook <event>`
**Rationale**: Shell script wrappers historically caused 28 distinct issues (PATH, encoding, SIGALRM). Direct binary invocation eliminates the entire wrapper layer and all associated issues.

### ADR-04: Why No Override Skills?

**Decision**: No override skills; all knowledge in rules
**MoAI comparison**: 6 override skills with progressive disclosure
**Rationale**: Do uses "rules -> always loaded" instead of "skills -> progressive disclosure" for framework knowledge. Maintaining two layers (skills + rules) adds complexity without proportional benefit. A single layer (rules only) is simpler to maintain and reason about.

### ADR-05: Why Korean Mixed-Language CLAUDE.md?

**Decision**: CLAUDE.md in Korean + English
**MoAI comparison**: English-only instruction documents
**Rationale**: Do's primary user is Korean. The persona system uses Korean honorifics. Korean is the design language, not just a translation target. The instruction document should reflect this.

### ADR-06: Why Date-Based Jobs Instead of Numbered SPECs?

**Decision**: `.do/jobs/{YYMMDD}/{title-kebab-case}/`
**MoAI comparison**: `.moai/specs/SPEC-XXX/`
**Rationale**: Date-based organization enables chronological browsing. Looking at a job folder, you immediately know when the work was done. Number-based SPEC-XXX requires looking up a registry to understand ordering.

### ADR-07: Why Persona + Style as Independent Axes?

**Decision**: Persona (who speaks) and Style (how they speak) are orthogonal
**MoAI comparison**: Single style system (MoAI/R2-D2/Yoda)
**Rationale**: Independence gives 4 x 3 = 12 combinations with only 7 definitions. Coupling them would require defining each combination separately. A young-f persona might use sprint style for quick fixes and pair style for complex work.

### ADR-08: Why .* PostToolUse Matcher?

**Decision**: PostToolUse hook fires on every tool call
**MoAI comparison**: PostToolUse fires only on `Write|Edit`
**Rationale**: Persona consistency. The hook maintains the AI's persona (honorifics, speech patterns). If it only fires on writes, the persona can drift during read-heavy sequences. The token cost is accepted as the price of consistent personality.

### ADR-09: Why File-Detection Triggers?

**Decision**: File existence activates behavior (Convention over Configuration)
**MoAI comparison**: No equivalent mechanism
**Rationale**: Zero-configuration activation. Placing `.git.multirepo` in the project root enables multirepo support. No settings file to edit, no flag to set, no configuration to get wrong. Adding or removing the file is the configuration.

### ADR-10: Why Reject TRUST 5 Branding?

**Decision**: Adopt the 5 quality dimensions, reject the acronym
**MoAI comparison**: TRUST 5 branded framework
**Rationale**: The project owner's judgment: it feels forced. The quality dimensions are valid and adopted as embedded [HARD] rules. The acronym adds no value -- it is a naming exercise, not a quality improvement. Do prefers substance over branding.

---

## 13. Terminology Map

| MoAI Term | Do Equivalent | Notes |
|-----------|---------------|-------|
| MoAI | Do | Brand name |
| `.moai/` | `.do/` | Project directory |
| `moai` (CLI) | `godo` (CLI) | Go binary |
| `/moai` | `/do:*` (6 commands) | Entry point structure |
| `moai-` (skill prefix) | `do-` | Skill naming |
| SPEC | Plan + Checklist | Workflow artifact |
| SPEC-XXX | `.do/jobs/{YYMMDD}/{title}/` | Artifact location |
| EARS | MoSCoW | Requirement format |
| TRUST 5 | [HARD] rules in dev-*.md | Quality framework |
| TAG Chain | Checklist dependency (`depends on:`) | Task dependencies |
| `<moai>DONE</moai>` | `[o]` + commit hash | Completion evidence |
| `<moai>COMPLETE</moai>` | report.md written | Full completion |
| Plan/Run/Sync | Plan/Checklist/Develop/Test/Report | Workflow phases |
| `.moai/config/sections/*.yaml` | `settings.local.json` DO_* env | Configuration |
| `.moai/learning/` | (none) | Yoda style learning directory |
| Snapshot/Resume | Checklist state (persistent file) | Continuity mechanism |
| `manager-spec` | `manager-plan` | Plan/checklist creation (renamed to reflect plan-centric workflow) |
| Completion Marker (XML) | Checklist status `[o]` | Task completion |
| outputStyle: "MoAI" | outputStyle: "pair" | Default style |
| moai.md (style) | pair.md + persona | Default behavior |
| r2d2.md (style) | sprint.md | Fast execution style |
| yoda.md (style) | direct.md | Expert style |
| (none) | Persona system (4 types) | Do unique |
| (none) | AI Anti-Pattern 7 | Do unique |
| (none) | File-Detection Triggers | Do unique |
| (none) | `.*` PostToolUse matcher | Do unique |
| (none) | Tri-mode (Focus/Do/Team) | Do unique |
| (none) | Append-only commit log philosophy | Do unique -- core tracking mechanism |
| `manager-strategy` (embedded) | `do-foundation-philosopher` (independent) | Elevated to standalone capability |

---

## 14. Do's Vision: The Best Team Orchestrator

The project owner described Do's vision as:

> "최고의 팀 오케스트레이터" (The best team orchestrator)

This vision has two dimensions that MoAI addresses only partially:

### 14.1 Persona + Team = AI with Personality Leading Teams

MoAI is an orchestrator that delegates to agents. Do is an orchestrator *with a personality* that delegates to agents. The persona system means the orchestrator is not an anonymous coordinator -- it is a character with a name, a speech pattern, a relationship with the user, and cultural context.

When Do's young-f persona says "승민선배, 이 작업은 Team 모드로 전환해볼까요?" (Senior, shall we switch to Team mode for this?), it is not just a mode suggestion -- it is a colleague-to-colleague conversation. The persona creates engagement that a neutral "Recommend switching to Team mode" cannot.

### 14.2 Adaptive Force + Parallel Execution

MoAI is always in "full orchestration" mode. Do's tri-mode system means the orchestrator knows when to be heavy (Team mode with parallel agents) and when to be light (Focus mode with direct execution). This adaptive intelligence is part of what makes Do aspire to be "the best" team orchestrator -- not just an orchestrator that always operates at maximum force.

### 14.3 The Gap MoAI Cannot Fill

MoAI can adopt Do's tri-mode system (it is a structural feature). MoAI can adopt Do's checklist system (it is a workflow feature). But MoAI cannot adopt Do's persona system without fundamentally changing its identity. MoAI is "the Strategic Orchestrator" -- an institution, not a person. Do is "나는 Do다" -- a first-person entity with a character, a relationship, and a voice.

This is Do's deepest differentiator: not a feature that can be copied, but an identity that must be chosen.

---

**Document Version**: 2.1.0
**Date**: 2026-02-16
**Sources**: research-moai-philosophy.md, research-do-philosophy.md, CLAUDE.md, dev-testing.md, dev-workflow.md, dev-checklist.md, dev-environment.md
