# Agent Research Delegation

Purpose: Defines the rules for delegating research tasks to agents, preventing context explosion and ensuring research artifacts are written progressively to files rather than accumulated in memory.

Version: 1.0.0
Last Updated: 2026-02-17

---

## Quick Reference (30 seconds)

Research Delegation Rules:
- **File output needed** → use general-purpose agent (NOT Explore agent)
- **Code search only** → Explore agent is fine
- **Already loaded in system prompt** → summarize and pass, do NOT re-read
- **Large files (10K+ chars)** → never read multiple at once in one agent
- **Progressive Write** → agent writes findings to file as discovered, NOT all at end

---

## Implementation Guide (5 minutes)

### Agent Type Selection [HARD]

| Need | Agent Type | Reason |
|------|-----------|--------|
| Research with file output | general-purpose | Explore agent cannot Write — results pile up in context and explode |
| Code search / navigation | Explore | Read-only, lightweight, fast |
| Research with analysis doc | general-purpose | Must write analysis.md, architecture.md, etc. |

### Context Management [HARD]

- **Do NOT re-read files already in system prompt** — summarize the relevant parts and pass to the agent
- **Do NOT have an agent read multiple large files (10K+ chars) simultaneously** — context overflow risk
- **Progressive Write is mandatory**: Instruct the agent to write findings to the output file as they are discovered, not accumulate everything and write once at the end
- **Always pass the output file path** and the instruction: "Write findings to this file as you discover them"

### Research Agent Instruction Template

When delegating research to an agent, include:
1. Output file path (where to write results)
2. "Write findings progressively as you discover them"
3. Summarized context from system prompt (not raw file paths to re-read)
4. Scope constraints (what to investigate, what to skip)

### Anti-Patterns

WRONG — Using Explore agent for research that needs file output:
```
Task(subagent_type="Explore", prompt="Research X and write analysis to file")
# Explore cannot write files → results stay in context → context explosion
```

CORRECT — Using general-purpose agent:
```
Task(subagent_type="general-purpose", prompt="Research X, write findings progressively to /path/to/output.md")
```

WRONG — Having agent re-read files already in system prompt:
```
Task(prompt="Read CLAUDE.md and analyze...")
# CLAUDE.md is already loaded — just summarize the relevant part
```

CORRECT — Passing summarized context:
```
Task(prompt="Based on this project structure: [summary]. Analyze...")
```

---

## Works Well With

Foundation Modules:
- [Agent Execution Cycle](agent-execution-cycle.md) - Execution cycle that may trigger research
- [Agent Delegation](agent-delegation.md) - Interruption handling for long research tasks
- [Delegation Patterns](delegation-patterns.md) - How to structure research delegation

Skills:
- {{slot:BRAND}}-foundation-context - Token budget and progressive loading strategies

---

Version: 1.0.0
Last Updated: 2026-02-17
Status: Production Ready
