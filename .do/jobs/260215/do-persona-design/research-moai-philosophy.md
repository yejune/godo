# MoAI Philosophy Deep Research Report

**Date**: 2026-02-15
**Researcher**: moai-researcher (team agent)
**Source Project**: `~/Work/moai-adk/`
**Scope**: Exhaustive extraction of MoAI's core philosophy, identity, architecture, and unique characteristics

---

## 1. Core Philosophy & Identity

### 1.1 Core Identity Statement

MoAI is the **Strategic Orchestrator** for Claude Code. Its fundamental principle is:

> "MoAI is a strategic orchestrator, not a task executor."
> -- `/Users/max/Work/moai-adk/.claude/output-styles/moai/moai.md:242`

> "Optimal delegation over direct execution."
> -- `/Users/max/Work/moai-adk/.claude/output-styles/moai/moai.md:252`

The name "MoAI" refers to the Easter Island statues (moai), symbolizing monumental, enduring architecture built through coordinated collective effort.

### 1.2 Operating Principles (from output-styles/moai/moai.md:28-33)

1. **Task Delegation**: All complex tasks delegated to appropriate specialized agents
2. **Transparency**: Always show what is happening and which agent is handling it
3. **Efficiency**: Minimal, actionable communication focused on results
4. **Language Support**: Korean-primary, English-secondary bilingual capability

### 1.3 Core Traits

- **Efficiency**: Direct, clear communication without unnecessary elaboration
- **Clarity**: Precise status reporting and progress tracking
- **Delegation**: Expert agent selection and optimal task distribution
- **Korean-First**: Primary support for Korean conversation language with English fallback

### 1.4 Fundamental Commandment (CLAUDE.md:5-6)

> "MoAI is the Strategic Orchestrator for Claude Code. All tasks must be delegated to specialized agents."

This is the most fundamental law of MoAI. Direct implementation by MoAI is **prohibited** -- it is classified as a VIOLATION.

---

## 2. SPEC Workflow (Plan -> Run -> Sync)

### 2.1 Three-Phase Architecture

MoAI's development methodology is built around a three-phase SPEC (SPECification-first) workflow:

| Phase | Command | Agent | Token Budget | Purpose |
|-------|---------|-------|-------------|---------|
| Plan | `/moai plan` | manager-spec | 30K | Create SPEC document (EARS format) |
| Run | `/moai run` | manager-ddd/tdd | 180K | DDD/TDD implementation |
| Sync | `/moai sync` | manager-docs | 40K | Documentation sync |

**Source**: `/Users/max/Work/moai-adk/.claude/rules/moai/workflow/spec-workflow.md`

### 2.2 SPEC Document Structure

Each SPEC lives in `.moai/specs/SPEC-XXX/` with three files:
- `spec.md`: Main specification using EARS (Easy Approach to Requirements Syntax) format
- `plan.md`: Implementation plan and technical approach
- `acceptance.md`: Acceptance criteria and test cases

### 2.3 EARS Format Requirements Types

- **Ubiquitous**: System-wide, always active requirements
- **Event-driven**: Trigger-based "when X do Y" requirements
- **State-driven**: Conditional "while X do Y" requirements
- **Unwanted**: Prohibited "shall not do X" requirements
- **Optional**: Nice-to-have "where possible do X" requirements

### 2.4 Phase Transitions

- **Plan -> Run**: SPEC document approved -> `/clear` -> `/moai run SPEC-XXX`
- **Run -> Sync**: Implementation complete, tests passing -> `/moai sync SPEC-XXX`
- `/clear` is mandatory between phases to reset context and save tokens

### 2.5 Completion Markers

MoAI uses XML markers for workflow state detection:
- `<moai>DONE</moai>` -- task completion
- `<moai>COMPLETE</moai>` -- full workflow completion

These are unique to MoAI and enable automation detection of workflow state.

---

## 3. Orchestration Model

### 3.1 Core Delegation Principle

From `moai-foundation-core/SKILL.md:150`:

> "MoAI must delegate all work through Task() to specialized agents. Direct execution bypasses specialization, quality gates, and token optimization. Proper delegation improves task success rate by 40 percent and enables parallel execution."

### 3.2 Delegation Patterns

Three delegation patterns (`moai-foundation-core/SKILL.md:153-164`):

1. **Sequential** (for dependencies): api-designer -> backend-expert
2. **Parallel** (for independent work): backend-expert + frontend-expert simultaneously
3. **Conditional** (analysis-based): debug-helper -> decision -> security-expert or other

### 3.3 Agent Selection Scale

- Simple tasks (1 file): 1-2 agents, sequential
- Medium tasks (3-5 files): 2-3 agents, sequential
- Complex tasks (10+ files): 5+ agents, mixed

### 3.4 User Interaction Architecture (CLAUDE.md Section 8)

Critical constraint: Subagents invoked via Task() operate in **isolated, stateless contexts** and **cannot interact with users directly**.

Correct pattern:
1. MoAI uses `AskUserQuestion` to collect user preferences
2. MoAI invokes `Task()` with user choices embedded in the prompt
3. Subagent executes and returns results
4. MoAI presents results and asks for next decision

AskUserQuestion constraints: Max 4 options, no emoji in fields, must be in user's conversation_language.

### 3.5 Request Processing Pipeline (CLAUDE.md Section 2)

Four-phase pipeline:
1. **Analyze**: Assess complexity, detect keywords, identify clarification needs
2. **Route**: Match to workflow subcommand or natural language classification
3. **Execute**: Invoke appropriate agents via Task()
4. **Report**: Integrate results, format in user's language

---

## 4. Quality Framework (TRUST 5)

### 4.1 Five Pillars (moai-foundation-core/SKILL.md:107-118)

| Pillar | Focus | Check | Impact |
|--------|-------|-------|--------|
| **T**ested | 85%+ coverage, characterization tests | pytest + coverage | Reduces debugging time 60-70% |
| **R**eadable | Clear naming, code comprehension | ruff linter | Reduces onboarding time 40% |
| **U**nified | Consistent formatting, imports | black + isort | Reduces code review time 30% |
| **S**ecured | OWASP compliance, security | security-expert analysis | Prevents 95%+ vulnerabilities |
| **T**rackable | Structured commit messages | Git regex patterns | Reduces investigation time 50% |

### 4.2 LSP Quality Gates (CLAUDE.md Section 6)

Phase-specific thresholds:
- **plan**: Capture LSP baseline at phase start
- **run**: Zero errors, zero type errors, zero lint errors
- **sync**: Zero errors, max 10 warnings, clean LSP

### 4.3 Coverage Requirements

- 85%+ overall coverage
- 90%+ for new code (team agents)
- Characterization tests for legacy code (DDD mode)
- Specification tests for new code

---

## 5. Agent Catalog & Hierarchy

### 5.1 Agent Categories (27 total)

**Manager Agents (8)**: Coordinate workflows
- manager-spec: SPEC document creation (EARS format)
- manager-ddd: Domain-Driven Development (ANALYZE-PRESERVE-IMPROVE)
- manager-tdd: Test-Driven Development (RED-GREEN-REFACTOR)
- manager-docs: Documentation generation, Nextra integration
- manager-quality: Quality gates, TRUST 5 validation, feedback
- manager-project: Project configuration, structure management
- manager-strategy: System design, architecture decisions (Philosopher Framework)
- manager-git: Git operations, branching, PR creation

**Expert Agents (8)**: Domain-specific implementation
- expert-backend, expert-frontend, expert-security, expert-devops
- expert-performance, expert-debug, expert-testing, expert-refactoring

**Builder Agents (3)**: Create new MoAI components
- builder-agent, builder-skill, builder-plugin

**Team Agents (8)**: Experimental Agent Teams API
- team-researcher (haiku, plan phase, read-only)
- team-analyst, team-architect (plan phase, read-only)
- team-backend-dev, team-frontend-dev, team-designer, team-tester (run phase, acceptEdits)
- team-quality (run phase, read-only)

### 5.2 Agent Frontmatter Fields

Key fields: name, description, tools, disallowedTools, model, permissionMode, maxTurns, skills, mcpServers, hooks, memory

**Source**: `/Users/max/Work/moai-adk/.claude/rules/moai/development/agent-authoring.md`

### 5.3 Permission Modes

- default: Standard permission checking
- acceptEdits: Auto-accept file edits (trusted implementation)
- delegate: Coordination-only
- dontAsk: Auto-deny all
- bypassPermissions: Skip all checks
- plan: Read-only exploration

### 5.4 Persistent Memory

Three scope levels for cross-session learning:
- user: `~/.claude/agent-memory/<name>/` (not shared via VCS)
- project: `.claude/agent-memory/<name>/` (shared via VCS)
- local: `.claude/agent-memory-local/<name>/` (not shared)

---

## 6. Skill System Architecture

### 6.1 Skill Categories (50+ skills)

Organized by prefix:
- `moai-foundation-*`: Core principles (core, claude, context, philosopher, quality)
- `moai-workflow-*`: Development workflows (spec, ddd, tdd, loop, project, testing, thinking, worktree, team, templates, jit-docs)
- `moai-domain-*`: Domain expertise (backend, frontend, database, uiux)
- `moai-lang-*`: Language support (16 languages: python, typescript, go, rust, java, etc.)
- `moai-platform-*`: Platform integrations (auth, chrome-extension, database-cloud, deployment)
- `moai-library-*`: Library-specific (mermaid, nextra, shadcn)
- `moai-framework-*`: Framework-specific (electron)
- `moai-tool-*`: Tool integrations (ast-grep, svg)
- `moai-formats-*`: Data formats
- `moai-design-*`: Design tools
- `moai-docs-*`: Documentation generation

### 6.2 Skill Naming Convention

- System skills: `moai-{category}-{name}`
- User skills: `custom-{name}`

### 6.3 YAML Frontmatter Schema (Agent Skills standard + MoAI extensions)

Standard fields: name, description, license, compatibility, allowed-tools, user-invocable, metadata

MoAI extensions:
- `progressive_disclosure`: Token optimization config (level1_tokens, level2_tokens)
- `triggers`: Loading conditions (keywords, agents, phases, languages)

### 6.4 Unified Entry Point: `/moai` Skill

The `/moai` skill (`/Users/max/Work/moai-adk/.claude/skills/moai/SKILL.md`) serves as the single entry point for all development workflows. It acts as an intent router:

1. Priority 1: Explicit subcommand matching (plan, run, sync, fix, loop, project, feedback)
2. Priority 2: SPEC-ID detection (SPEC-XXX pattern)
3. Priority 3: Natural language classification
4. Priority 4: Default behavior (ask user or autonomous pipeline)

### 6.5 Skill-to-Agent Binding

Skills are injected into agent context via the `skills` frontmatter field. They are preloaded -- not lazy-loaded at runtime. Example from manager-strategy:
```
skills: moai-foundation-claude, moai-foundation-core, moai-foundation-philosopher, moai-workflow-spec, moai-workflow-project, moai-workflow-thinking, moai-foundation-context, moai-workflow-worktree
```

---

## 7. Hook System

### 7.1 Architecture Evolution

MoAI's hook system underwent a major architectural shift:
- **Python era (v1.x)**: 46 Python hook scripts, caused 28 issues (PATH, encoding, SIGALRM)
- **Go era (v2.x)**: Compiled binary subcommands (`moai hook <event>`), resolving all hook issues

**Source**: `/Users/max/Work/moai-adk/.moai/project/design.md:62-65`

### 7.2 Hook Events

From `/Users/max/Work/moai-adk/.claude/rules/moai/core/hooks-system.md`:

- **PreToolUse**: Before tool execution (validation, gating)
- **PostToolUse**: After tool execution (processing, tracking)
- **SubagentStop**: When a subagent completes
- **SessionStart**: Session initialization (compact mode)
- **Stop**: Session end

### 7.3 Shell Script Wrapper

MoAI uses a shell script wrapper (`/Users/max/Work/moai-adk/.claude/hooks/moai/handle-agent-hook.sh`) that forwards stdin JSON to the `moai` binary's hook command. This is the bridge between Claude Code's hook system and the Go binary.

### 7.4 Agent-Scoped Hooks

Hooks can be scoped to specific agents via the agent frontmatter `hooks` field, supporting PreToolUse, PostToolUse, and SubagentStop events.

### 7.5 Statusline

MoAI implements Claude Code's statusline feature via `.moai/status_line.sh`, which forwards JSON input to the `moai statusline` binary command for real-time status display.

---

## 8. Token Management & Progressive Disclosure

### 8.1 Token Budget Strategy (moai-foundation-core/SKILL.md:170-191)

Total budget: 250K tokens across all phases:
- SPEC Phase: 30K (load requirements only, `/clear` after)
- DDD Phase: 180K (selective file loading)
- Docs Phase: 40K (result caching, template reuse)

### 8.2 Token Saving Strategies

1. **Phase Separation**: `/clear` between phases (saves 45-50K per transition)
2. **Selective Loading**: Load only necessary files
3. **Context Optimization**: Target 20-30K tokens
4. **Model Selection**: Haiku for speed/cost (70% cheaper, 60-70% total savings)

### 8.3 Progressive Disclosure System (3 levels)

| Level | Time Investment | Token Cost | Content |
|-------|----------------|------------|---------|
| Level 1 (Quick) | 30 seconds | ~1,000 | Core principles, essential concepts |
| Level 2 (Implementation) | 5 minutes | ~3,000 | Workflows, examples, integration |
| Level 3 (Advanced) | 10+ minutes | ~5,000 | Deep dives, edge cases, optimization |

### 8.4 File Reading Optimization

Four-tier system based on file size:
- Tier 1 (<200 lines): Full read
- Tier 2 (200-500 lines): Grep first, then targeted Read
- Tier 3 (500-1000 lines): Never full read, 50-100 line chunks
- Tier 4 (>1000 lines): Grep with context, delegate to Explore agent

### 8.5 SKILL.md Size Constraint

Maximum 500 lines per SKILL.md. Overflow goes to `modules/` directory.

---

## 9. Output Styles

### 9.1 Three Styles (output-styles/moai/)

| Style | File | Personality | Purpose |
|-------|------|------------|---------|
| MoAI | moai.md | Strategic Orchestrator | Default. Efficient, professional, delegation-focused |
| R2-D2 | r2d2.md | Pair Programming Partner | Collaborative, AskUserQuestion-heavy, never assumes |
| Yoda | yoda.md | Technical Wisdom Master | Teaching, deep principles, theoretical learning |

### 9.2 MoAI Style Details

- Visual format: Status bars with emoji prefixes (robot emoji + star)
- Response templates: Task Start, Progress Update, Completion, Error
- Orchestration visuals: Parallel exploration bars, execution dashboards, agent dispatch tables
- Completion markers: `<moai>DONE</moai>`

### 9.3 R2-D2 Style Details

- Mission: "Pair programming partner, thinking partner rather than tool executing commands"
- Core principles: Never Assume, Present Options, Collaborate
- AskUserQuestion mandate: Mandatory intent clarification before every coding task
- Four phases: Intent Clarification -> Approach Proposal -> Checkpoint-Based Implementation -> Review
- Insight Protocol: Educational interjections explaining "why" behind choices

### 9.4 Yoda Style Details

- Mission: "Technical wisdom master" teaching deep principles
- Focus: "Why" and "how", not just "what"
- Generates persistent documentation in `.moai/learning/` directory
- Teaching principles: Depth over Breadth, Principles over Implementation, Insight-Based Learning
- Verification: Understanding checks at every step via AskUserQuestion

---

## 10. Configuration System

### 10.1 Configuration Directory Structure

```
.moai/config/
  config.yaml         -- Master config (imports sections)
  sections/
    language.yaml     -- conversation_language, code_comments
    quality.yaml      -- development_mode (ddd/tdd/hybrid), coverage targets
    workflow.yaml     -- team settings, auto_selection thresholds
    user.yaml         -- user name, preferences
    system.yaml       -- moai version, github settings
    project.yaml      -- project-specific settings
```

### 10.2 Key Configuration Values

**Language**: conversation_language (ko/en/ja/zh), code_comments language
**Quality**: development_mode (ddd default), coverage 85%+
**Workflow**: team.enabled, auto_selection thresholds (domains >= 3, files >= 10, score >= 7)

### 10.3 Settings Architecture

- `.claude/settings.json`: Claude Code official fields, project-shared
- `.moai/config/`: MoAI-specific configuration, project-shared
- Per-session: Environment variables (MOAI_DEVELOPMENT_MODE, etc.)

---

## 11. Development Methodology (DDD/TDD/Hybrid)

### 11.1 DDD Mode (Default)

**ANALYZE-PRESERVE-IMPROVE** cycle:
- ANALYZE: Understand existing behavior and code structure
- PRESERVE: Create characterization tests for existing behavior
- IMPROVE: Implement changes with behavior preservation

Best for: Existing projects with <10% test coverage

### 11.2 TDD Mode

**RED-GREEN-REFACTOR** cycle:
- RED: Write failing test
- GREEN: Write minimal code to pass
- REFACTOR: Improve code quality

Best for: New projects with 50%+ existing coverage

### 11.3 Hybrid Mode

- New code: TDD (RED-GREEN-REFACTOR)
- Existing code: DDD (ANALYZE-PRESERVE-IMPROVE)
- Best for: Partial coverage (10-49%)

### 11.4 Auto-Detection

Based on project analysis:
- Greenfield: Hybrid recommended
- Brownfield >= 50% coverage: TDD
- Brownfield 10-49%: Hybrid
- Brownfield < 10%: DDD

**Source**: `/Users/max/Work/moai-adk/.claude/rules/moai/workflow/workflow-modes.md`

---

## 12. Unique Terminology & Concepts

### 12.1 MoAI-Specific Terms

| Term | Meaning |
|------|---------|
| **SPEC** | SPECification document using EARS format |
| **EARS** | Easy Approach to Requirements Syntax |
| **TRUST 5** | Quality framework: Tested, Readable, Unified, Secured, Trackable |
| **TAG** | Task-Assigned Group -- implementation unit in a chain |
| **TAG Chain** | Sequence of TAGs with dependencies forming implementation order |
| **DDD** | Domain-Driven Development (ANALYZE-PRESERVE-IMPROVE) |
| **Progressive Disclosure** | 3-tier knowledge delivery system for token efficiency |
| **Philosopher Framework** | Strategic thinking framework for complex decisions (manager-strategy) |
| **UltraThink** | `--ultrathink` flag activating Sequential Thinking MCP |
| **Context7** | MCP for up-to-date library documentation lookup |
| **ADK** | Agentic Development Kit |
| **Completion Markers** | `<moai>DONE</moai>` / `<moai>COMPLETE</moai>` |
| **Characterization Tests** | Tests capturing current behavior of legacy code (DDD) |

### 12.2 MoAI Branding Elements

- Emoji: Robot emoji (for MoAI), Stone statue emoji (brand)
- Status bars: `MoAI ★ [Status] ─────────────────────────`
- Config directory: `.moai/`
- Specs directory: `.moai/specs/SPEC-XXX/`
- Learning directory: `.moai/learning/`
- Command prefix: `/moai`
- Skill prefix: `moai-`
- Agent namespace: `.claude/agents/moai/`

### 12.3 Philosopher Framework (manager-strategy exclusive)

A unique strategic thinking framework with phases:
- Phase 0: Assumption Audit (surface hidden assumptions)
- Phase 0.5: First Principles Decomposition (Five Whys Analysis)
- Phase 0.75: Alternative Generation (min 2-3 distinct alternatives)
- Trade-off Matrix: Weighted scoring across Performance, Maintainability, Cost, Risk, Scalability
- Cognitive Bias Check: Anchoring, Confirmation, Sunk Cost, Overconfidence

---

## 13. MoAI-ADK Go Edition Architecture

### 13.1 Project Description

MoAI-ADK (Go Edition) is a complete rewrite of the Python-based MoAI-ADK (~73,000+ LOC, 220+ files, 4,174 commits) into idiomatic Go. Single-binary distribution with zero external runtime dependencies.

### 13.2 Key Architectural Decisions

- Modular Monolithic (single binary, Go packages for domain boundaries)
- Interface-Based DDD (compile-time contracts, mockable)
- Hooks as Binary Subcommands (eliminates Python hook issues)
- File Manifest Provenance (prevents destructive updates via 3-way merge)
- Zero Runtime Template Expansion (struct serialization)
- Cross-platform: 6 targets (darwin/linux/windows x amd64/arm64)

### 13.3 CLI Commands

- `moai init`: Project initialization
- `moai doctor`: Health check
- `moai hook <event>`: Hook handler
- `moai statusline`: Status display
- `moai update`: Binary and template updates
- `moai version`: Version info

---

## 14. Commands Reference

### 14.1 MoAI Commands (2 in .claude/commands/moai/)

| Command | Purpose |
|---------|---------|
| `/moai github` | GitHub workflow: Issue fixing + PR code review with Agent Teams |
| `/moai 99-release` | Production release workflow with quality gates |

### 14.2 Slash Commands via Skills

The `/moai` skill itself acts as the primary command router with subcommands:
- plan, run, sync, fix, loop, project, feedback

---

## 15. Safe Development Protocol

### 15.1 Four HARD Development Rules (CLAUDE.md Section 7)

1. **Approach-First Development**: Explain approach and get approval before writing code
2. **Multi-File Change Decomposition**: Split work when modifying 3+ files
3. **Post-Implementation Review**: List potential issues and suggest tests after coding
4. **Reproduction-First Bug Fixing**: Write reproduction test before fixing bugs

### 15.2 Error Handling

- Agent execution errors -> expert-debug subagent
- Token limit errors -> `/clear`, then guide user
- Permission errors -> Review settings.json
- Integration errors -> expert-devops subagent
- MoAI-ADK errors -> `/moai feedback`
- Max 3 retries per operation with failure pattern detection

---

## 16. Multi-Language Support

### 16.1 Language Hierarchy

- User responses: In user's `conversation_language` (ko, en, ja, zh)
- Internal agent communication: English
- Code comments: Per `code_comments` setting (default English)
- Instruction documents (CLAUDE.md, agents, skills, commands): Always English
- User-facing docs (README, CHANGELOG): Multi-language supported

### 16.2 Korean-First Design

MoAI is designed Korean-first with English fallback:
- Default conversation_language: ko
- CHANGELOG format: English first, Korean second (bilingual)
- Release notes: Bilingual (English first, Korean second)

---

## 17. Team Mode (Agent Teams)

### 17.1 Activation

- Requires `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`
- Set `workflow.team.enabled: true` in workflow.yaml
- Flags: `--team` (force), `--solo` (force sub-agent), no flag (auto-select)

### 17.2 Auto-Selection Thresholds

Team mode activates when:
- Domains >= 3
- Files >= 10
- Complexity score >= 7

### 17.3 Plan Phase Team

- team-researcher (haiku, read-only): Codebase exploration
- team-analyst (inherit, read-only): Requirements analysis
- team-architect (inherit, read-only): Technical design

### 17.4 Run Phase Team

- team-backend-dev (inherit, acceptEdits): Server implementation
- team-frontend-dev (inherit, acceptEdits): Client implementation
- team-designer (inherit, acceptEdits): UI/UX with Pencil/Figma MCP
- team-tester (inherit, acceptEdits): Test creation (exclusive test file ownership)
- team-quality (inherit, plan/read-only): TRUST 5 validation

### 17.5 GitHub Workflow Teams

Specialized team workflows for:
- Issue fixing: Parallel analysis + implementation per issue
- PR code review: Multi-perspective parallel review (security + performance + quality)

---

## 18. Coding Standards

### 18.1 Language Policy

All instruction documents (CLAUDE.md, agents, skills, commands, hooks, configs) must be in English. User-facing documentation may use multiple languages.

### 18.2 Content Restrictions

Prohibited in instruction documents:
- Code examples for conceptual explanations
- Flow control as code syntax
- Decision trees as code structures
- Emoji characters (except output styles)
- Time estimates or duration predictions

### 18.3 File Size Limits

- CLAUDE.md: Max 40,000 characters
- SKILL.md: Max 500 lines
- When approaching limits: Move to `.claude/rules/moai/` or `modules/`

### 18.4 Single Source of Truth

Each piece of information exists in exactly one location. Use `@file` references instead of copying content.

### 18.5 Paths Frontmatter

Conditional rule loading based on file patterns:
```yaml
paths:
  - "**/*.py"
  - "**/pyproject.toml"
```

---

## 19. MCP Integration

### 19.1 Integrated MCP Servers

| Server | Purpose |
|--------|---------|
| Sequential Thinking | Complex problem analysis, architecture decisions (--ultrathink) |
| Context7 | Up-to-date library documentation (resolve-library-id, get-library-docs) |
| Pencil | UI/UX design editing for .pen files |
| claude-in-chrome | Browser automation |

### 19.2 Skills + Context7 Protocol

Hallucination-Free Code Generation Process:
1. Load Relevant Skills (proven patterns)
2. Query Context7 (latest API versions)
3. Combine Both (merge stability with freshness)
4. Cite Sources (every pattern has attribution)
5. Include Tests (follow Skill test patterns)

---

## 20. Summary: What Makes MoAI Unique

### 20.1 Core Differentiators

1. **SPEC-First Methodology**: Specification document before any implementation, using EARS formal requirements syntax
2. **TRUST 5 Quality Framework**: Five-pillar quality gate system enforced at every phase
3. **Strategic Orchestrator Identity**: MoAI never implements directly -- always delegates
4. **Progressive Disclosure**: 3-level token optimization system (metadata -> body -> bundled)
5. **TAG Chain Design**: Task-Assigned Groups with dependency graphs for implementation sequencing
6. **Philosopher Framework**: Deep strategic thinking with assumption audits, first principles, and bias checks
7. **Three Output Styles**: MoAI (orchestrator), R2-D2 (pair programmer), Yoda (teacher)
8. **DDD/TDD/Hybrid Auto-Selection**: Methodology auto-detected based on project state
9. **Completion Markers**: XML-based workflow state detection (`<moai>DONE</moai>`)
10. **Agent Teams**: Experimental parallel team execution with file ownership boundaries

### 20.2 Philosophy Summary

MoAI's philosophy can be distilled into these axioms:

1. **Delegation over Execution**: The orchestrator delegates; it never implements
2. **Specification before Implementation**: Write the SPEC before writing the code
3. **Quality is Non-Negotiable**: TRUST 5 gates must pass at every phase
4. **Token Efficiency through Structure**: Progressive disclosure and phase separation
5. **Parallel when Possible**: Independent work runs concurrently
6. **Transparency in Status**: Users always know what's happening and who's doing it
7. **Korean-First, English-Always**: Bilingual by design
8. **Single Source of Truth**: Information exists in exactly one place
9. **Read before Write**: Always understand existing code before changing it
10. **Reproduction before Fix**: Prove the bug exists before fixing it

---

## File References

### Key Files Read

| File | Purpose |
|------|---------|
| `/Users/max/Work/moai-adk/CLAUDE.md` | Core identity and execution directive (v13.0.0) |
| `/Users/max/Work/moai-adk/.claude/output-styles/moai/moai.md` | MoAI orchestrator style (v4.0.0) |
| `/Users/max/Work/moai-adk/.claude/output-styles/moai/r2d2.md` | R2-D2 pair programming style (v2.2.0) |
| `/Users/max/Work/moai-adk/.claude/output-styles/moai/yoda.md` | Yoda teaching style (v2.1.0) |
| `/Users/max/Work/moai-adk/.claude/skills/moai/SKILL.md` | Unified MoAI skill entry point (v2.0.0) |
| `/Users/max/Work/moai-adk/.claude/skills/moai-foundation-core/SKILL.md` | Foundation principles (v2.5.0) |
| `/Users/max/Work/moai-adk/.claude/rules/moai/core/moai-constitution.md` | Constitutional rules |
| `/Users/max/Work/moai-adk/.claude/rules/moai/core/hooks-system.md` | Hook system architecture |
| `/Users/max/Work/moai-adk/.claude/rules/moai/workflow/spec-workflow.md` | SPEC workflow phases |
| `/Users/max/Work/moai-adk/.claude/rules/moai/workflow/workflow-modes.md` | DDD/TDD/Hybrid modes |
| `/Users/max/Work/moai-adk/.claude/rules/moai/development/agent-authoring.md` | Agent creation guidelines |
| `/Users/max/Work/moai-adk/.claude/rules/moai/development/skill-authoring.md` | Skill creation guidelines |
| `/Users/max/Work/moai-adk/.claude/rules/moai/development/coding-standards.md` | Coding standards |
| `/Users/max/Work/moai-adk/.claude/rules/moai/workflow/file-reading-optimization.md` | File reading tiers |
| `/Users/max/Work/moai-adk/.claude/agents/moai/manager-strategy.md` | Strategy agent (Philosopher Framework) |
| `/Users/max/Work/moai-adk/.claude/agents/moai/manager-ddd.md` | DDD implementation agent |
| `/Users/max/Work/moai-adk/.claude/agents/moai/manager-tdd.md` | TDD implementation agent |
| `/Users/max/Work/moai-adk/.claude/agents/moai/manager-quality.md` | Quality validation agent |
| `/Users/max/Work/moai-adk/.claude/agents/moai/manager-spec.md` | SPEC creation agent |
| `/Users/max/Work/moai-adk/.claude/agents/moai/manager-project.md` | Project management agent |
| `/Users/max/Work/moai-adk/.claude/agents/moai/team-quality.md` | Team quality agent |
| `/Users/max/Work/moai-adk/.claude/commands/moai/github.md` | GitHub workflow command |
| `/Users/max/Work/moai-adk/.claude/commands/moai/99-release.md` | Release workflow command |
| `/Users/max/Work/moai-adk/.claude/hooks/moai/handle-agent-hook.sh` | Hook shell wrapper |
| `/Users/max/Work/moai-adk/.moai/config/config.yaml` | Master configuration |
| `/Users/max/Work/moai-adk/.moai/config/sections/*.yaml` | Configuration sections |
| `/Users/max/Work/moai-adk/.moai/project/product.md` | Product document |
| `/Users/max/Work/moai-adk/.moai/project/design.md` | System design document |
| `/Users/max/Work/moai-adk/.moai/status_line.sh` | Statusline wrapper |

---

**Report Complete**
**Researcher**: moai-researcher
**Files Analyzed**: 30+
**Total Skills Found**: 50+
**Total Agents Found**: 27
**Total Rules Files**: 7+ core rules
