# Knowledge Management

Rules for information lookup priority and documentation placement.

## Reference Order [HARD]

- [HARD] Priority 1: `memory/MEMORY.md` and project documentation -- use existing knowledge first
- [HARD] Priority 2: Codebase search (Grep/Glob) -- check existing patterns and implementations
- [HARD] Priority 3: External search (WebSearch/Context7) -- only when the above two are insufficient

## Documentation Placement [HARD]

- [HARD] Discoveries during work -> Checklist Lessons Learned section (scoped to current task)
- [HARD] Cross-session patterns/debugging notes -> `memory/MEMORY.md` (auto-loaded into system prompt)
- [HARD] Project architecture decisions/rules -> Project documentation (README, CLAUDE.md, docs/)
- [HARD] When unsure where to write -> Write to memory first, organize later

## Rationale

This priority order ensures:
1. No redundant external searches when information is already available locally
2. Documentation lives close to where it is most useful
3. Knowledge compounds across sessions through persistent memory
4. Project-level decisions are discoverable by all contributors

## Anti-Patterns

- Searching the web for information already documented in the project
- Re-reading files that are already in the system prompt
- Writing discoveries only in chat output (lost after session ends)
- Duplicating information across multiple locations

## Related

- [File Reading Optimization](file-reading-optimization.md)
