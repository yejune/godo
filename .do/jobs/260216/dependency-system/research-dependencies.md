# Dependency System Research Report

**Date**: 2026-02-16  
**Scope**: MoAI-ADK / Do Framework dependency system investigation  
**Analyzed**: Agent definitions, workflow documents, hook systems, checklist rules

---

## 1. CURRENT DEPENDENCY METADATA

### 1.1 Agent Orchestration Metadata (YAML frontmatter)

All agents in `.claude/agents/do/*.md` contain orchestration metadata:

```yaml
orchestration:
  can_resume: bool
  typical_chain_position: "initial|middle|terminal"
  depends_on: [list of agent names]
  spawns_subagents: bool
  token_budget: "low|medium|high"
  context_retention: "low|medium|high"
  output_format: string
```

### 1.2 Discovered Dependencies

#### Phase 1: Plan Phase
- **Agent**: `manager-spec`
  - Position: `initial` (workflow starter)
  - Dependencies: `[]` (no dependencies)
  - Output: SPEC document (`.moai/specs/SPEC-XXX/spec.md`)

#### Phase 2: Run Phase (DDD - Legacy Refactoring)
- **Agent**: `manager-ddd`
  - Position: `middle`
  - Dependencies: `["manager-spec"]` (must complete Plan first)
  - Input: SPEC document from Phase 1
  - Output: Refactored code, characterization tests

#### Phase 2: Run Phase (TDD - New Features)
- **Agent**: `manager-tdd`
  - Position: `middle`
  - Dependencies: `["manager-spec"]`
  - Input: SPEC document from Phase 1
  - Output: Implementation code with tests

#### Phase 3: Sync Phase
- **Agent**: `manager-docs`
  - Position: `terminal`
  - Dependencies: `["manager-ddd", "manager-quality"]`
  - Input: Implementation from Phase 2 + Quality validation
  - Output: Documentation (API docs, README, CHANGELOG)

#### Quality Gate (Terminal)
- **Agent**: `manager-quality`
  - Position: `terminal`
  - Dependencies: `["manager-ddd"]` (or `["manager-tdd"]`)
  - Input: Implemented code
  - Output: Quality verification report

#### Domain Expert Agents (Parallel execution in Run phase)
- **expert-backend**
  - Position: `middle`
  - Dependencies: `["manager-spec"]`
  - Used for: API design, database, server logic consultation

- **expert-frontend**
  - Position: `middle`
  - Dependencies: `["manager-spec"]`
  - Used for: UI/component consultation

- **expert-testing**
  - Position: `middle`
  - Dependencies: `["expert-backend", "expert-frontend", "manager-ddd"]`
  - Used for: Test strategy validation

- **expert-performance**
  - Position: `middle`
  - Dependencies: `["expert-backend", "expert-frontend", "expert-database"]`
  - Used for: Performance optimization validation

- **expert-security**
  - Position: `middle`
  - Dependencies: `["expert-backend", "expert-frontend"]`
  - Used for: Security review

- **expert-devops**
  - Position: `middle`
  - Dependencies: `["expert-backend", "expert-frontend"]`
  - Used for: Infrastructure/deployment consultation

#### Git Operations
- **Agent**: `manager-git`
  - Position: `terminal`
  - Dependencies: `["manager-quality", "manager-ddd"]` (or `["manager-tdd"]`)
  - Handles: Branch creation, commits, PR operations

#### Strategy & Architecture
- **Agent**: `manager-strategy`
  - Position: `initiator`
  - Dependencies: `["manager-spec"]`
  - Handles: System design, architecture decisions

#### Project Setup
- **Agent**: `manager-project`
  - Position: `initiator`
  - Dependencies: `[]` (no dependencies)
  - Handles: Project configuration, initialization

---

## 2. WORKFLOW PHASE REFERENCES

### 2.1 SPEC Workflow (3 Phases)

From `.claude/rules/do/workflow/spec-workflow.md`:

| Phase | Agent | Upstream Dependencies | Downstream Dependencies |
|-------|-------|----------------------|------------------------|
| **Plan** | manager-spec | None | Run phase blocks until SPEC approved |
| **Run** | manager-ddd/tdd | SPEC document (Plan) | Sync phase blocks until tests pass |
| **Sync** | manager-docs | Implementation (Run) + Quality (Quality gate) | None (terminal) |

**Phase Transitions**:
- Plan → Run: `/clear` mandatory between phases
- Run → Sync: `<moai>COMPLETE</moai>` marker + tests passing
- Sync → (none): Terminal phase

### 2.2 Workflow Modes (Affects Run Phase)

From `.claude/rules/do/workflow/workflow-modes.md`:

Development methodology selection depends on project state:
- **DDD Mode**: Use `manager-ddd` for legacy code refactoring
  - ANALYZE → PRESERVE → IMPROVE cycle
  - Depends on existing tests
  
- **TDD Mode**: Use `manager-tdd` for new features
  - RED → GREEN → REFACTOR cycle
  - Creates tests first
  
- **Hybrid Mode**: Both DDD and TDD
  - New code: TDD
  - Existing code: DDD

### 2.3 Team Mode Variant (Agent Teams API)

From `spec-workflow.md` Team Mode Phase Overview:

**Plan Phase (Team)**:
- Parallel agents: `team-researcher`, `team-analyst`, `team-architect`
- All three run simultaneously
- MoAI synthesizes findings into SPEC
- **Dependency**: All three must complete before SPEC synthesis

**Run Phase (Team)**:
- Parallel agents: `team-backend-dev`, `team-frontend-dev`, `team-tester`
- Shared task list with file ownership
- **Dependency**: File ownership prevents conflicts
- Quality gate validates after all complete

---

## 3. CHECKLIST SYSTEM DEPENDENCIES

### 3.1 Checklist Structure (`.do/jobs/` hierarchy)

From `.claude/rules/dev-checklist.md`:

```
.do/jobs/{YYMMDD}/{title-kebab-case}/
├── analysis.md                    # Complex tasks only
├── architecture.md                # Complex tasks only
├── plan.md                        # Depends on analysis/architecture
├── checklist.md                   # Depends on plan
├── report.md                      # Depends on checklists completion
└── checklists/
    ├── 01_expert-backend.md       # Item-level checklist
    ├── 02_expert-frontend.md
    └── 03_expert-testing.md
```

### 3.2 Item-Level Dependencies

**Dependency Declaration**:
```markdown
- [~] #2 로그인 API 구현 (depends on: #1)
- [ ] #3 프론트엔드 로그인 폼 (depends on: #2)
- [!] #4 소셜 로그인 연동 (depends on: #2, 블로커: OAuth 키 미발급)
```

**Rules**:
- Use `depends on: #N` keyword for item-level dependencies
- If dependency not complete → item marked as `[!]` blocker
- Status symbols: `[ ]` pending, `[~]` in-progress, `[*]` testing, `[o]` done, `[x]` failed, `[!]` blocked

### 3.3 Checklist as Agent State File

**Critical Rule**: Checklists are NOT just documentation - they are **agent persistent state**:
- Agent reads checklist at start
- Agent updates status as items complete
- On token exhaustion, new agent continues from last `[o]` item
- Guarantees task continuity across agent interruptions

---

## 4. WORKFLOW COMPLEXITY TRIGGERS (Analysis/Architecture)

From `.claude/rules/dev-workflow.md`:

**Complex Task (triggers Analysis → Architecture)**:
- ✓ 5+ files change expected
- ✓ New library/package/module creation
- ✓ System migration/tech stack change
- ✓ 3+ domain integration (backend + frontend + DB)
- ✓ Abstraction layer design needed
- ✓ Architecture changes (monolith → microservices)

**Simple Task** (Plan → Develop directly):
- ≤ 4 files change
- Existing pattern implementation
- Single domain
- No architecture changes

**Workflow if Complex**:
1. Analysis (expert-analyst) → analysis.md
2. Architecture (expert-architect) → architecture.md
3. Plan (user/manager) → plan.md
4. Checklist → checklist.md
5. Development
6. Testing
7. Report

**Workflow if Simple**:
1. Plan → plan.md
2. Checklist → checklist.md
3. Development
4. Testing
5. Report

---

## 5. SKILL SYSTEM PHASE TRIGGERS

From `.claude/rules/do/development/skill-authoring.md`:

Skills have conditional triggers:

```yaml
triggers:
  keywords: ["api", "database", "authentication"]
  agents: ["manager-spec", "expert-backend"]
  phases: ["plan", "run"]
  languages: ["python", "typescript"]
```

**Phase References in Skills**:
- `phases: ["plan"]`: Load during Plan phase only
- `phases: ["run"]`: Load during Run phase only
- `phases: ["plan", "run"]`: Load in both phases
- `phases: []`: Load always (fallback)

**Example**: Skill `do-workflow-spec` loads during Plan phase for SPEC creation.

---

## 6. EXISTING ENFORCEMENT MECHANISMS

### 6.1 In Codebase

**File**: `cmd/godo/statusline.go`
- Reads mode state from `.do/.current-mode` file
- Displays `[Do]`, `[Focus]`, `[Team]` prefix based on current mode
- **Enforcement**: No explicit dependency checking in statusline

**File**: `cmd/godo/moai_hook_contract.go`
- Hook system for lifecycle events (PreToolUse, PostToolUse, SubagentStop)
- Hooks can trigger shell scripts for validation
- **Enforcement**: Hook scripts could validate dependencies (but currently empty)

### 6.2 Hook Locations

**Expected**: `.claude/hooks/moai/handle-agent-hook.sh`
**Actual**: Hooks directory is **empty** (no validation scripts deployed)

**Hook Actions** (from `.claude/rules/do/core/agent-hooks.md`):
```
spec-completion     → After manager-spec finishes
ddd-pre-transformation, ddd-post-transformation → Around manager-ddd changes
ddd-completion → After manager-ddd finishes
tdd-pre-implementation, tdd-post-implementation → Around manager-tdd changes
tdd-completion → After manager-tdd finishes
quality-completion → After manager-quality finishes
docs-verification → After manager-docs writes
docs-completion → After manager-docs finishes
```

**Current Status**: Hooks defined in agent YAML but no implementation scripts.

---

## 7. FILE & ARTIFACT DEPENDENCIES

### 7.1 Artifact Dependencies (Must exist before downstream work)

| Artifact | Created by | Required by | Condition |
|----------|-----------|-----------|-----------|
| `plan.md` | Plan phase | Checklist creation | Always |
| `analysis.md` | expert-analyst | architecture.md creation | Complex tasks only |
| `architecture.md` | expert-architect | plan.md refinement | Complex tasks only |
| `checklist.md` | Orchestrator | Agent execution | Always (blocks development) |
| `checklists/*.md` | Agent | Report creation | Always |
| `.do/jobs/YYMMDD/` | Auto-created | All docs | Directory structure |

### 7.2 Plan File Dependencies

From CLAUDE.md (dev-workflow.md):

- **Complexity Assessment**: Required before choosing analysis/architecture path
- **TDD Decision**: Required before implementing ("TDD로 개발할까요?")
- **Docker Environment Info**: Required before agent delegation
- **Commit Instructions**: Required with agent delegation

---

## 8. IMPLICIT DEPENDENCIES (Not Currently Enforced)

### 8.1 Phase Sequencing

**Implicit Rule**: Phases must execute sequentially
- Plan completes → `/clear` called → Run begins
- Run completes → tests must pass → Sync begins
- Sync completes → ready for deployment

**Current Enforcement**: None (user responsible for `/clear`, phase transitions)

### 8.2 Agent Availability Dependencies

Agents require skills pre-loaded:

**Example - manager-spec**:
```yaml
skills: do-foundation-claude, do-foundation-core, do-foundation-philosopher, 
        do-workflow-spec, do-workflow-project, do-workflow-thinking, 
        do-lang-python, do-lang-typescript
```

If skill unavailable → agent degraded functionality

**Current Enforcement**: None (skills auto-loaded from YAML, fallback undefined)

### 8.3 Environment Dependencies

From `dev-environment.md`:

- Docker Compose must be running: `docker bootapp up`
- Services must be healthy before tests
- Port conflicts prevent parallel execution
- Database migration must complete before tests

**Current Enforcement**: Manual (user responsibility)

---

## 9. DEPENDENCY CANDIDATES FOR FORMALIZATION

### 9.1 Phase Dependencies (HIGH PRIORITY)

```
Plan Phase
├── Input: User request + project docs
├── Agent: manager-spec
├── Output: SPEC document
└── Blocking: Run phase cannot start without SPEC approval

Run Phase
├── Input: SPEC document (from Plan)
├── Agent: manager-ddd OR manager-tdd (per quality.yaml)
├── Output: Implementation code + tests
└── Blocking: Sync phase cannot start without tests passing

Sync Phase
├── Input: Implementation (from Run) + Quality validation
├── Agent: manager-docs
├── Output: Documentation
└── Blocking: Terminal phase (no downstream)
```

**Validation Points**:
- Before Run: Verify `.moai/specs/SPEC-XXX/spec.md` exists
- Before Sync: Verify all tests passing (exit code 0)
- Before Sync: Verify quality gate completed

### 9.2 Checklist Item Dependencies (MEDIUM PRIORITY)

```
Item Dependency Graph:
- Parse `depends on: #N` from checklist items
- Build directed graph of item IDs
- Validate no circular dependencies
- Block item execution until dependencies complete
- Auto-mark items `[!]` if dependency blocked
```

**Validation Points**:
- Parse checklist.md for `depends on:` annotations
- Build dependency graph at checklist creation
- Validate no cycles (would deadlock)
- Block agent from claiming item if dependencies incomplete

### 9.3 File Artifact Dependencies (MEDIUM PRIORITY)

```
Complexity Path Triggers:
- IF 5+ files OR new library OR multiple domains OR architecture change:
  analysis.md MUST exist before architecture.md
  architecture.md MUST exist before plan.md
  
Simple Path:
- plan.md created directly (no analysis/architecture)
```

**Validation Points**:
- Check complexity criteria before workflow selection
- If complex: verify analysis.md exists before creating architecture.md
- If complex: verify architecture.md exists before creating plan.md

### 9.4 Team Mode Task Dependencies (LOW PRIORITY - Already handled by TaskCreate/TaskUpdate)

```
Plan Phase Team:
- researcher task → analyst task (analyze requires research)
- analyst task → architect task (design requires analysis)
- All complete → synthesis by MoAI

Run Phase Team:
- File ownership prevents conflicts
- TaskCreate/TaskUpdate manages sequencing
- Quality validation after all tasks done
```

**Validation Points**:
- TaskCreate sets explicit ordering
- TaskUpdate marks status transitions
- Team lead validates before synthesis

---

## 10. DEPENDENCY SYSTEM IMPLEMENTATION GAPS

### 10.1 Missing Enforcement

1. **Phase Sequencing**: No validation that Plan → Run → Sync order is respected
   - Users could try to run implementation without SPEC
   - `/clear` between phases is advisory, not enforced

2. **Checklist Existence**: No blocking of agent execution without checklist
   - Agents could theoretically work without checklist
   - Current: Trust-based (HARD rule in docs)

3. **Artifact Existence**: No validation before downstream phases
   - Run phase could start without SPEC document
   - Sync phase could start with failing tests

4. **Complexity Assessment**: No automation of analysis/architecture decision
   - User must manually decide (AskUserQuestion)
   - Could be automated based on file count, domain detection

5. **Skill Availability**: No fallback if skill unavailable
   - Agent degrades but continues (documented in YAML)
   - Silent failure possible

### 10.2 Hook System (Defined but Not Implemented)

Agent hooks exist in YAML frontmatter:

```yaml
hooks:
  PreToolUse:
    - matcher: "Write|Edit|MultiEdit"
      hooks:
        - type: command
          command: "\"$CLAUDE_PROJECT_DIR/.claude/hooks/moai/handle-agent-hook.sh\" ddd-pre-transformation"
```

But `.claude/hooks/moai/` directory is **empty**.

**Implementation Opportunity**:
- `ddd-pre-transformation`: Validate checklist exists
- `ddd-post-transformation`: Update checklist status
- `spec-completion`: Validate SPEC format (EARS compliance)
- `tdd-completion`: Verify test coverage

---

## 11. SUMMARY TABLE: Dependency Types

| Type | Current State | Enforcement | Examples |
|------|--------------|-------------|----------|
| **Phase** | Documented | Trust-based | Plan → Run → Sync |
| **Agent** | Metadata (depends_on) | Not enforced | manager-spec → manager-ddd |
| **Artifact** | Rule-based | Trust-based | analysis.md → architecture.md → plan.md |
| **Item-level** | Documented | Not enforced | checklist items with `depends on:` |
| **Skill** | YAML triggers | Auto-load | phases: ["plan", "run"] |
| **Environment** | Documentation | Manual | Docker, healthcheck |
| **File Ownership** | Team mode only | Task system | Per-agent file assignments |

---

## 12. RECOMMENDATIONS

### Immediate (High-Priority Gaps)

1. **Implement Hook Scripts** (`.claude/hooks/moai/handle-agent-hook.sh`)
   - Validate artifacts exist before phase transitions
   - Validate checklist format and item dependencies
   - Enforce EARS syntax in SPEC documents

2. **Add Phase Validation** (godo or hook system)
   - Block Run without SPEC document
   - Block Sync without test passing status
   - Require `/clear` between Plan and Run

3. **Automate Complexity Assessment**
   - File count detection
   - Domain keyword matching
   - Auto-trigger analysis/architecture if complex

### Medium-Priority Enhancements

4. **Dependency Graph Validation**
   - Parse checklist item dependencies
   - Detect circular dependencies at checklist creation
   - Visualize dependency graph for debugging

5. **Skill Fallback System**
   - Define graceful degradation if skill unavailable
   - Log warnings when skill triggers fail
   - Suggest alternative skills/agents

6. **Team Mode Dependency Tracking**
   - Extend TaskCreate with explicit dependency graph
   - Visualize task DAG in status reports
   - Prevent team members from claiming out-of-order tasks

### Long-Term (Nice-to-Have)

7. **Dependency Specification Language**
   - YAML syntax for declaring complex workflows
   - Reusable workflow templates
   - Conditional dependencies (if X then Y)

8. **Observability**
   - Dependency timeline visualization
   - Critical path analysis
   - Bottleneck detection (which dependencies block most items)

---

## 13. FILES ANALYZED

### Agent Definitions
- `.claude/agents/do/manager-spec.md` - Plan phase
- `.claude/agents/do/manager-ddd.md` - Run phase (legacy)
- `.claude/agents/do/manager-tdd.md` - Run phase (new)
- `.claude/agents/do/manager-docs.md` - Sync phase
- `.claude/agents/do/manager-quality.md` - Quality validation
- `.claude/agents/do/expert-backend.md`, `expert-frontend.md`, `expert-testing.md`, etc.
- `.claude/agents/do/team-researcher.md`, `team-analyst.md`, `team-architect.md` - Team mode

### Workflow Documents
- `.claude/rules/do/workflow/spec-workflow.md` - Phase overview & transitions
- `.claude/rules/do/workflow/workflow-modes.md` - DDD/TDD/Hybrid methodology
- `.claude/rules/dev-workflow.md` - Complexity-based workflow selection
- `.claude/rules/dev-checklist.md` - Checklist structure & dependencies
- `.claude/rules/dev-environment.md` - Docker & environment dependencies
- `.claude/rules/do/development/skill-authoring.md` - Skill phase triggers
- `.claude/rules/do/development/agent-authoring.md` - Agent definition format

### Code
- `cmd/godo/statusline.go` - Mode detection
- `cmd/godo/moai_hook_contract.go` - Hook system definition

