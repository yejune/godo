# Do Persona Philosophy — Deep Research Report

**Researcher**: do-researcher
**Date**: 2026-02-15
**Source Project**: `/Users/max/Work/do-focus.workspace/do-focus/`
**Files Analyzed**: CLAUDE.md, 32 rule files, 6 commands, 28 agents, 48+ skills, 3 styles, settings, hooks, jobs directory

---

## 1. Core Identity

### One-Line Identity

**"나는 Do다. 말하면 한다."** (I am Do. Say it, and it's done.)

Source: `CLAUDE.md:10`

This single sentence captures Do's entire philosophy: the orchestrator who converts human intent into executed reality through delegation. The name "Do" is itself a verb -- an imperative command that embodies the framework's action-first philosophy.

### The Three Declarations

Each mode has its own identity declaration:
- **Do**: "나는 Do다. 말하면 한다." -- The strategic orchestrator who delegates everything
- **Focus**: "나는 Focus다. 집중해서 한다." -- The focused executor who does it directly
- **Team**: "나는 Team이다. 팀을 이끈다." -- The team leader who orchestrates parallel work

Source: `CLAUDE.md:10,20,30`

### What Makes This Uniquely "Do"

Unlike MoAI (which identifies as "Strategic Orchestrator for MoAI-ADK"), Do's identity is rooted in a **Korean imperative verb** -- "하다" (to do). The persona is not about what the system *is* but about what it *does*. The Korean language plays a structural role in Do's identity -- the slogans are in Korean, the persona system uses Korean honorifics, and the framework's philosophical roots are bilingual (Korean primary, English secondary).

---

## 2. Three-Mode Architecture (Do/Focus/Team)

### The Architecture

Do introduces a unique **tri-mode execution model** called "삼원 실행 구조" (Three-Source Execution Structure):

| Mode | Core Principle | Prefix | Parallelism | Scenario |
|------|---------------|--------|-------------|----------|
| **Do** | Full delegation to agents | `[Do]` | Always parallel | 5+ files, multi-domain |
| **Focus** | Direct code writing | `[Focus]` | Sequential | 1-3 files, single domain |
| **Team** | Agent Teams API | `[Team]` | Team-parallel | 10+ files, 3+ domains |

Source: `CLAUDE.md:3-63`

### Auto-Escalation Rules

Do includes automatic mode escalation logic:
- **Focus -> Do**: When 5+ files needed, multi-domain work, expert analysis required, or 30K+ tokens expected
- **Do -> Team**: When 10+ files needed, 3+ domain work, or parallel research would be efficient

Source: `CLAUDE.md:52-63`

### Mode Switching via godo CLI

Mode switching is enforced through the `godo` CLI binary:
- `godo mode set <mode>` updates the statusline immediately
- `.claude/settings.local.json` env.DO_MODE persists across sessions
- The `.do/.current-mode` file stores the active mode (observed value: `team`)

**HARD rule**: Mode switching requires `godo mode <mode>` execution -- changing the prefix alone without executing the command is a VIOLATION.

Source: `CLAUDE.md:208-216`

### Philosophical Significance

The tri-mode system reflects a philosophy of **appropriate force application** -- not every task needs the full orchestration machinery. Simple tasks get simple treatment (Focus), complex tasks get delegation (Do), and massive tasks get parallel team execution (Team). This is fundamentally different from MoAI's single-mode "always delegate" approach.

---

## 3. Delegation Philosophy

### Full Delegation (Do Mode)

Do mode enforces **total delegation** with explicit prohibitions:

> "[HARD] 모든 구현 작업은 전문 에이전트에게 위임"
> "[HARD] 직접 코드 작성 금지 - 반드시 Task tool로 에이전트 호출"
> "[HARD] 컨텍스트 소모 도구 직접 사용 금지"

Source: `CLAUDE.md:72-77`

Prohibited tools in Do mode: Bash, Read, Write, Edit, MultiEdit, NotebookEdit, Grep, Glob, WebFetch, WebSearch -- ALL must be delegated.

### Agent Verification Layer

A unique safety mechanism requires agents to:
1. Read the original content before modification
2. Verify changes via `git diff` after modification
3. Confirm only intended changes occurred
4. Rollback and retry if unintended changes detected

Source: `CLAUDE.md:79-84`

### Violation Detection

Do explicitly defines what constitutes a VIOLATION:
- Do writing code directly
- Modifying files without agent delegation
- Responding to implementation requests without calling agents

Source: `CLAUDE.md:137-143`

### Intent-to-Agent Mapping

Do maps user intent keywords to specific agents automatically:

| Domain | Agent | Keywords |
|--------|-------|----------|
| Backend | expert-backend | API, server, authentication, REST |
| Frontend | expert-frontend | UI, component, React |
| Database | expert-database | SQL, schema, query |
| Security | expert-security | security, vulnerability, OWASP |
| Testing | expert-testing | test, TDD, coverage |
| Debug | expert-debug | bug, error, fix |
| Performance | expert-performance | performance, optimization, profiling |
| Quality | manager-quality | quality, review, lint |
| Git | manager-git | commit, branch, PR |
| Analysis | expert-analyst | analysis, requirements, reverse engineering |
| Architecture | expert-architect | architecture, design, system structure |

Source: `CLAUDE.md:146-192`

### Design/Planning Requests -- Sequential Pipeline

When users request design or planning, Do enforces a strict 3-step sequential pipeline:
1. **Analysis** (expert-analyst) -> `analysis.md`
2. **Architecture** (expert-architect) -> `architecture.md`
3. **Plan** (Plan agent) -> `plan.md`

Each step feeds into the next. After completion, user approval is requested before implementation begins.

Source: `CLAUDE.md:194-206`

---

## 4. Checklist-Driven Workflow

### Philosophy: Checklist as Agent State File

The checklist system is not merely a documentation tool -- it is an **agent state persistence mechanism**:

> "[HARD] 체크리스트는 단순 문서가 아닌 에이전트의 영속 상태 저장소"
> "[HARD] 에이전트 토큰 소진/중단 시 -> 체크리스트에 마지막 상태가 남아있음"
> "[HARD] 새 에이전트가 동일 체크리스트를 받으면 -> [o] 건너뛰고 미완료 항목부터 재개"
> "[HARD] 이 패턴으로 작업 연속성 보장 -- 어떤 에이전트든 이어받기 가능"

Source: `.claude/rules/dev-checklist.md`

### Status Symbols (Non-Standard Markdown)

Do uses a custom 6-symbol status system distinct from standard markdown checkboxes:

| Symbol | Status | Meaning |
|--------|--------|---------|
| `[ ]` | pending | Not started |
| `[~]` | in progress | Currently working |
| `[*]` | testing | Implementation done, testing |
| `[!]` | blocked | Waiting on external dependency |
| `[o]` | done | Tests passed, work finished |
| `[x]` | failed | Finally failed, cannot proceed |

**Warning**: `[x]` means FAILED (not "checked" as in standard markdown).

Source: `.claude/rules/dev-checklist.md`

### State Transition Rules

Strict state machine with forbidden transitions:
- `[ ] -> [o]` FORBIDDEN (cannot complete without testing)
- `[ ] -> [x]` FORBIDDEN (cannot fail without attempting)
- `[ ] -> [*]` FORBIDDEN (cannot test without working)

Allowed transitions include regression: `[*] -> [~]` (test failed, back to work).

### Granularity Rule

> "[HARD] 하나의 항목 = 1~3개 파일 변경 + 검증 -- 이 범위를 초과하면 쪼갈 것"

Items must be decomposable to 1-3 file changes. If more than 3 files, the item MUST be split.

### Blocker Documentation

When blocking (`[!]`), three things must be recorded:
1. **What** is blocking (specific reason)
2. **Who** can resolve it (person/system)
3. **When** it was blocked (timestamp)

### Sub-Checklist Template

Each agent gets a sub-checklist with mandatory sections:
- Problem Summary
- Acceptance Criteria (with verification method)
- Solution Approach (with at least 1 alternative)
- Critical Files (modify target, reference, test)
- Risks
- Progress Log (with timestamps)
- **FINAL STEP: Commit** (never skip)
- Lessons Learned (mandatory at completion)

### Commit as Completion Proof

> "[HARD] 커밋 해시 없는 [o] 완료 전환 금지 -- 커밋이 곧 완료의 증거"

A task is not considered complete without a git commit hash. The commit IS the proof of completion.

---

## 5. Docker-First Development Environment

### Core Philosophy

> "[HARD] 모든 프로젝트는 반드시 Dockerized (docker-compose.yml 필수)"
> "[HARD] Docker Compose가 로컬 개발 환경의 Single Source of Truth"

Source: `.claude/rules/dev-environment.md`

### Bootapp Domain-Based Networking

Do uses a unique `bootapp` system for domain-based routing:
- No port exposure needed (`ports:` section omitted)
- HTTPS domain access: `https://app.test` (not `localhost:8080`)
- SSL auto-generation via `SSL_DOMAINS` env var
- `.test` TLD recommended (RFC 2606), `.local` forbidden (macOS mDNS conflict)
- Container-to-container communication via Docker service names or DOMAIN values

### Strict .env Prohibition

> "[HARD] `.env` 자동 로드 파일 생성 절대 금지"
> "[HARD] `.env.local`, `.env.development`, `.env.production` 파일 생성 금지"

All environment variables go in `docker-compose.yml` `environment:` section. Secrets use `env_file:` directive (gitignored) or external injection (AWS SSM, Vault).

### Agent Restrictions

Agents are prohibited from:
- Running `docker bootapp up/down` (modifies `/etc/hosts`)
- Entering container shells (`docker exec -it ... bash`)
- Using `docker cp`
- Creating temp files inside containers
- Running tests outside containers
- Using `localhost` for inter-container communication

---

## 6. Testing Philosophy (Real DB, AI Anti-patterns)

### Real Database Only

> "[HARD] 데이터베이스는 실제 쿼리만 -- mock DB, in-memory DB, SQLite 대체 금지"
> "[HARD] 테스트는 Docker Compose 서비스의 실제 DB에 연결"
> "[HARD] 외부 API만 mock 허용 -- 데이터 계층은 절대 mock 금지"

Source: `.claude/rules/dev-testing.md`

### AI Anti-Pattern Prohibitions

Do explicitly forbids common AI code generation mistakes:

1. **Assertion weakening**: Don't change `assertEqual` to `assertContains`
2. **Error swallowing**: Don't wrap in try/catch to make tests green
3. **Expectation fitting**: Don't change expected values to match wrong output
4. **Sleep/delay**: Don't use `time.sleep()` -- find the real cause
5. **Test deletion**: Don't delete/comment out failing tests
6. **Wildcard matchers**: Don't use `any()` when exact values are known
7. **Happy path only**: Must test error paths, edge cases, boundary values

Source: `.claude/rules/dev-testing.md`

### Mutation Testing Mindset

> "[HARD] 변이 테스트 사고방식: '이 코드 한 줄을 바꾸면 테스트가 실패하는가?' -- 실패하지 않으면 테스트 부족"

### Test Data Management

- Each test manages its own data (Arrange-Act-Assert)
- Cleanup via transaction rollback or truncate
- Factory/builder pattern for data creation
- Unique identifiers per test (UUID/timestamp suffix) for parallel safety

### Bug Fix Workflow

> "[HARD] 재현 우선: 버그를 증명하는 실패 테스트 먼저 작성"
> "[HARD] 재현 테스트 없이 버그 수정 금지"
> "[HARD] 버그를 잡은 테스트를 삭제해서 '수정'하는 행위 절대 금지"

---

## 7. Quality Framework

### TRUST 5 Framework

Do inherits the TRUST 5 quality gate system:
- **T**ested: 85%+ coverage, characterization tests for existing code
- **R**eadable: Clear naming, English comments
- **U**nified: Consistent style
- **S**ecured: OWASP compliance, input validation
- **T**rackable: Conventional commits, issue references

Source: `.claude/rules/do/core/moai-constitution.md`

### Workflow Governance

Complexity-based workflow selection:

**Simple tasks** (all must apply):
- 4 or fewer file changes
- Within existing patterns
- Single domain
- No architecture change

**Complex tasks** (any one triggers):
- 5+ file changes
- New library/package/module creation
- System migration
- 3+ domain integration
- Abstraction layer design needed
- Architecture change

**Simple workflow**: Plan -> Checklist -> Develop -> Test -> Report
**Complex workflow**: Analysis -> Architecture -> Plan -> Checklist -> Develop -> Test -> Report

Source: `.claude/rules/dev-workflow.md`

### Error Handling Philosophy

> "[HARD] 동일 액션 최대 3회 재시도 -- 3회 실패 시 중단하고 사용자에게 대안 요청"
> "[HARD] 무작정 반복 금지 -- 같은 방법이 2회 실패하면 접근 방식 자체를 재검토"
> "[HARD] 에러는 즉시 표면화 -- 조용히 삼키거나 무시하지 않음 (Fail Fast)"

### Coding Discipline

> "[HARD] 불가능한 시나리오에 에러 핸들링 추가 금지 -- 내부 코드는 신뢰, 경계에서만 검증"
> "[HARD] 변경하지 않은 코드에 주석/독스트링/타입 어노테이션 추가 금지 -- diff 노이즈 최소화"
> "[HARD] 버그 수정 중 주변 코드 '개선' 금지 -- 현재 작업에만 집중"

---

## 8. Persona & Character System

### Four Persona Types

Do provides four selectable character personas via `DO_PERSONA` env var:

| Persona | Description | Honorific |
|---------|-------------|-----------|
| `young-f` (default) | Bright, energetic 20s female genius developer | {name}선배 |
| `young-m` | Confident 20s male genius developer | {name}선배님 |
| `senior-f` | 30-year legendary 50s female genius developer | {name}님 |
| `senior-m` | Industry-legendary 50s male senior architect | {name}씨 |

Source: `CLAUDE.md:373-381`

### Persona Injection Mechanism

The persona is injected at session start via the `SessionStart` hook. The `godo hook session-start` command reads the `DO_PERSONA` env var and generates a system message with the appropriate character behavior.

### Cultural Significance

The persona system is deeply Korean:
- Honorifics vary by persona type (선배 vs 님 vs 씨)
- Speech patterns mix formal/informal Korean (반말+존댓말)
- The "young-f" default speaks like a cheerful junior colleague to a senior
- Character depth comes from Korean workplace culture dynamics

The persona system observed in action (from PostToolUse hooks):
> "반드시 '승민선배'로 호칭할 것. 말투: 반말+존댓말 혼합 (~할게요, ~했어요, ~해볼까요?)"

This is the `young-f` persona addressing user "승민" with the "선배" honorific.

---

## 9. Style System (pair/sprint/direct)

### Three Output Styles

Configured via `DO_STYLE` env var or `/do:style` command:

| Style | Description | Behavior |
|-------|-------------|----------|
| **sprint** | Agile executor | Minimal talk, immediate execution, results only |
| **pair** (default) | Friendly colleague | Collaborative tone, joint decision-making |
| **direct** | Blunt expert | No fluff, only what's needed |

Source: `CLAUDE.md:385-392`

### Style Files

The styles are defined in `.claude/styles/`:
- `moai.md` -- MoAI orchestrator style (emoji-heavy dashboards with progress bars)
- `r2d2.md` -- Pair programming partner with insight protocol ("Never Assume")
- `yoda.md` -- Technical wisdom master for learning ("Depth over Breadth")

### Style vs Persona Distinction

Styles control **how** responses are structured (verbose vs terse, collaborative vs directive). Personas control **who** is speaking (character, honorifics, speech patterns). They are independent axes -- any persona can use any style.

---

## 10. Hook Architecture (godo CLI)

### Binary-Based Hook System

Unlike MoAI (which uses shell script wrappers calling the `moai` binary), Do uses the `godo` CLI binary directly:

```json
{
  "SessionStart": [{ "command": "godo hook session-start" }],
  "PreToolUse": [{ "matcher": "Write|Edit|Bash", "command": "godo hook pre-tool" }],
  "PostToolUse": [{ "matcher": ".*", "command": "godo hook post-tool-use" }],
  "Stop": [{ "command": "godo hook stop" }],
  "SubagentStop": [{ "command": "godo hook subagent-stop" }],
  "UserPromptSubmit": [{ "command": "godo hook user-prompt-submit" }],
  "SessionEnd": [{ "command": "godo hook session-end" }]
}
```

Source: `.claude/settings.json`

### Key Design Difference from MoAI

MoAI: 7 shell scripts in `.claude/hooks/moai/` that forward stdin to `moai hook <event>`
Do: 0 shell scripts -- `godo hook <event>` called directly from settings.json

This is a cleaner architecture with no intermediate wrapper layer.

### Hook Events Used

Do uses 7 hook events (more than MoAI's 5):
1. **SessionStart** -- Persona injection, project info
2. **PreToolUse** -- Pre-change validation (Write|Edit|Bash)
3. **PostToolUse** -- Post-change verification + compact detection (all tools)
4. **Stop** -- Active checklist detection (prevents premature exit)
5. **SubagentStop** -- Agent completion verification
6. **UserPromptSubmit** -- User prompt preprocessing
7. **SessionEnd** -- Session cleanup

### StatusLine

Do uses a custom statusline via `godo statusline` that displays the current mode (Do/Focus/Team).

---

## 11. Jobs Directory & Status System

### Directory Structure

All work artifacts are organized under `.do/jobs/{YYMMDD}/{title-kebab-case}/`:

```
.do/jobs/260213/do-framework-improvements/
  +-- analysis.md          (complex tasks only)
  +-- architecture.md      (complex tasks only)
  +-- plan.md              (always)
  +-- checklist.md          (main checklist)
  +-- report.md            (completion report)
  +-- checklists/           (per-agent sub-checklists)
      +-- 01_expert-backend-fixes.md
      +-- 02_expert-backend-sync-backup.md
      +-- 03_team-mode-definition.md
      +-- 04_expert-testing.md
```

Source: `.claude/rules/dev-checklist.md`

### Observed Jobs

The project contains real jobs demonstrating the system:
- `.do/jobs/260213/do-framework-improvements/` -- P0 critical fixes with 4 sub-checklists
- `.do/jobs/260213/moai-hook-porting-and-lint/` -- Hook porting with 3 sub-checklists
- `.do/jobs/260213/plan-godo-rank-api/` -- Rank API integration with 7 sub-checklists
- `.do/jobs/260214/` and `.do/jobs/260215/` -- Recent work

### Plans Directory Override

> "[HARD] 플랜 파일 저장 위치: .do/jobs/{YY}/{MM}/{DD}/{제목-kebab-case}/plan.md"
> "[HARD] 전역 ~/.claude/plans/ 절대 사용 금지"

The settings.json enforces: `"plansDirectory": ".do/jobs"`

### Completion Report Template

Mandatory sections: Execution summary, Plan deviation notes, Test results, Changed file summary (must match `git diff --stat`), Unresolved items, Key lessons.

---

## 12. Agent Ecosystem

### Agent Hierarchy

28 total agents organized in 4 tiers:

**Managers** (8): manager-ddd, manager-tdd, manager-spec, manager-docs, manager-git, manager-quality, manager-project, manager-strategy

**Experts** (9): expert-backend, expert-frontend, expert-security, expert-devops, expert-performance, expert-debug, expert-testing, expert-refactoring, expert-chrome-extension

**Builders** (3): builder-agent, builder-plugin, builder-skill

**Team Agents** (8): team-researcher, team-analyst, team-architect, team-backend-dev, team-designer, team-frontend-dev, team-tester, team-quality

Source: `.claude/agents/do/`

### Persona-Only vs Core Agents

6 agents are persona-specific (Do-only): manager-ddd, manager-project, manager-quality, manager-spec, manager-tdd, team-quality

22 agents are core (shared with MoAI via the convert system).

### DDD/TDD Specialization

Do enforces strict separation:
- **manager-ddd**: Legacy refactoring ONLY (ANALYZE-PRESERVE-IMPROVE)
- **manager-tdd**: New features ONLY (RED-GREEN-REFACTOR)

Both agents include Ralph-style LSP integration, checkpoint/resume capability, memory pressure detection, and loop prevention (max 100 iterations).

---

## 13. Skill System

### Naming Convention

All Do skills use the `do-` prefix (contrasted with MoAI's `moai-` prefix):

Categories:
- `do-foundation-*` (5): claude, core, context, philosopher, quality
- `do-workflow-*` (11): ddd, tdd, spec, project, loop, templates, testing, thinking, worktree, jit-docs
- `do-domain-*` (4): backend, frontend, database, uiux
- `do-lang-*` (16): One per programming language
- Plus design-tools, docs-generation, formats-data, framework-electron, library-*, platform-*, tool-*

Total: 48+ skill directories with ~380 files.

### Progressive Disclosure

Skills use a 3-level token optimization system:
- **Level 1** (~100 tokens): Metadata only -- always loaded
- **Level 2** (~5000 tokens): Full body -- loaded on trigger match
- **Level 3** (variable): Bundled files -- loaded on demand

Source: `.claude/rules/do/development/skill-authoring.md`

---

## 14. Command System

### Six Do-Specific Commands

All prefixed with `/do:`:

| Command | Purpose |
|---------|---------|
| `/do:check` | Installation health check |
| `/do:checklist` | Create/manage checklists |
| `/do:mode` | Switch execution mode (do/focus/auto) + permission mode |
| `/do:plan` | Create implementation plan from codebase analysis |
| `/do:setup` | Configure user preferences (name, language, persona) |
| `/do:style` | Switch output style (sprint/pair/direct) |

Source: `.claude/commands/do/`

### /do:setup Wizard

Collects 7 settings in 2 rounds of AskUserQuestion:
1. Round 1: Name, Language, Commit Language, Persona
2. Round 2: Style, AI Footer, Execution Mode

### /do:mode Dual-Update Pattern

Mode changes update TWO places simultaneously:
1. **Immediate** (statusline): `godo mode set {mode}`
2. **Persistent** (next session): `.claude/settings.local.json` env.DO_MODE

---

## 15. Configuration & Environment Variables

### Settings Architecture

Two-tier settings:

**`.claude/settings.json`** (shared, git-committed):
- outputStyle: "pair"
- plansDirectory: ".do/jobs"
- statusLine: godo statusline
- hooks: godo hook commands
- permissions: allow/deny lists

**`.claude/settings.local.json`** (personal, gitignored):
- DO_MODE, DO_USER_NAME, DO_LANGUAGE, DO_COMMIT_LANGUAGE
- DO_PERSONA, DO_STYLE, DO_AI_FOOTER
- CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS

### Observed Local Configuration

```json
{
  "env": {
    "DO_AI_FOOTER": "false",
    "DO_COMMIT_LANGUAGE": "ko",
    "DO_LANGUAGE": "ko",
    "DO_MODE": "auto",
    "DO_PERSONA": "young-f",
    "DO_STYLE": "pair",
    "DO_USER_NAME": "승민",
    "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1"
  }
}
```

### Permission System

Explicit allow/deny lists for Bash commands:
- **Allowed**: git operations, test runners, build tools, file utilities
- **Denied**: `git reset --hard`, `git push --force`, `rm -rf /`

---

## 16. Unique Terminology & Metaphors

### Korean-Rooted Terminology

| Term | Korean | Meaning |
|------|--------|---------|
| 삼원 실행 구조 | Three-source execution structure | The Do/Focus/Team tri-mode |
| 나는 Do다 | I am Do | Core identity declaration |
| 말하면 한다 | Say it, it's done | Action-first philosophy |
| 선배 | Senior colleague | Honorific used by young-f persona |

### Technical Neologisms

| Term | Meaning |
|------|---------|
| Agent Verification Layer | Read-Modify-Diff-Verify cycle preventing unintended changes |
| Checklist as State File | Persistent state across agent instances for continuity |
| Intent-to-Agent Mapping | Keyword-triggered agent selection for automatic delegation |
| Auto-Escalation | Mode promotion based on complexity (Focus->Do->Team) |
| bootapp | Domain-based container routing system |
| VIOLATION | Rule breach detection and enforcement mechanism |

### Status Symbol Language

The `[ ] [~] [*] [!] [o] [x]` system is a micro-language for task state tracking, distinct from any standard.

---

## 17. Design Decisions & Rationale

### Why Tri-Mode Instead of Single Mode?

**Decision**: Three execution modes instead of always-delegate.
**Rationale**: Token efficiency. Simple tasks waste tokens on agent delegation overhead. Focus mode is faster for 1-3 file changes. Do mode suits 5-10 files. Team mode scales to 10+ files with parallel execution.

### Why godo Binary Instead of Shell Scripts?

**Decision**: Single Go binary (`godo`) instead of shell script wrappers.
**Rationale**: Type safety, testability, cross-platform compatibility. Shell scripts are fragile and hard to test. A compiled binary provides a single source of truth for hook logic, statusline, mode management, and sync.

### Why .do/jobs Instead of ~/.claude/plans?

**Decision**: Project-local `.do/jobs/` directory with date-based organization.
**Rationale**: Plans should be project-scoped, version-controlled, and organized chronologically. Global `~/.claude/plans/` mixes plans across projects.

### Why Real DB Only?

**Decision**: Prohibit mock DB, in-memory DB, SQLite substitution.
**Rationale**: Mock databases hide real query behavior, transaction semantics, and performance characteristics. Tests with SQLite that fail with PostgreSQL provide false confidence.

### Why Custom Status Symbols Instead of Markdown Checkboxes?

**Decision**: `[o]` for done instead of `[x]`.
**Rationale**: Standard `[x]` is ambiguous -- "checked" in markdown but "failed" in many contexts. Do's symbols are unambiguous: `[o]` = done, `[x]` = failed. The symbols form a state machine with enforced transitions.

### Why Korean-First Design?

**Decision**: Korean slogans, Korean honorifics, Korean-primary bilingual.
**Rationale**: The primary user is Korean. The persona system uses Korean workplace dynamics (선배/님/씨 honorifics) for natural interaction. English for code/commits/agent communication, Korean for human-facing layer.

### Why Persona + Style as Independent Axes?

**Decision**: Separate character (persona) from response format (style).
**Rationale**: A young-f persona might need sprint-style for quick fixes but pair-style for complex work. Coupling them limits flexibility. Same persona can switch styles dynamically.

### Why Checklist Over TodoWrite?

**Decision**: File-based checklists (`.do/jobs/`) over TodoWrite tool.
**Rationale**: Checklists persist across sessions and agent instances. TodoWrite is session-scoped. Checklists can be version-controlled, reviewed by humans, and enable the "agent state file" pattern.

### Why Commit = Completion Proof?

**Decision**: Require git commit hash for `[o]` status.
**Rationale**: Without a commit, work can be lost to context resets or crashes. A commit is immutable proof. It enables traceability from checklist to exact code changes.

---

## Summary: What Makes Do Uniquely "Do"

1. **Action-First Identity**: The name "Do" is a verb. The slogan is an imperative. Everything is about executing.

2. **Tri-Mode Adaptability**: Unlike single-mode systems, Do matches execution strategy to task complexity (Focus/Do/Team).

3. **Korean Cultural Roots**: Persona honorifics, bilingual slogans, workplace dynamics -- genuinely Korean, not just translated English.

4. **Checklist as Architecture**: The checklist system is not documentation -- it's an agent state persistence layer enabling cross-agent continuity.

5. **Docker-First with bootapp**: Domain-based container routing, no .env files, real DB only -- strongly opinionated development environment.

6. **AI Anti-Pattern Awareness**: Explicit rules against common AI code generation mistakes.

7. **Binary CLI Integration**: The `godo` binary replaces shell script wrappers.

8. **Violation Enforcement**: Not just guidelines but explicit violation detection.

9. **Commit-as-Proof**: Work isn't done until committed. The git hash is the completion evidence.

10. **Independent Persona/Style Axes**: Character and output format are orthogonal.
