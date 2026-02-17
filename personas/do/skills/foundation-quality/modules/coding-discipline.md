# Coding Discipline

Rules governing code modification behavior, verification, and error response.

## Agent Verification Layer [HARD]

- [HARD] Read(original) -> Modify -> git diff(verify changes) -> Confirm only intended changes
- [HARD] If unintended deletions/changes discovered -> rollback and retry
- [HARD] Single responsibility: each change does one thing well -- no mixing purposes
- [HARD] Do not add error handling for impossible scenarios -- trust internal code, validate only at boundaries
- [HARD] Do not add comments/docstrings/type annotations to code you did not change -- minimize diff noise
- [HARD] Do not "improve" surrounding code while fixing a bug -- stay focused on the current task

## Post-Coding Requirements [HARD]

- [HARD] After writing code: list what could break and suggest tests to cover it
- [HARD] Make meaningful small commits -- diff and message alone must convey intent
- [HARD] Verify tests pass before committing -- self-check for guideline violations

## Error Response [HARD]

- [HARD] Maximum 3 retries for the same action -- after 3 failures, stop and ask the user for alternatives
- [HARD] No blind repetition -- if the same approach fails twice, reconsider the approach itself
- [HARD] Surface errors immediately -- never silently swallow or ignore them (Fail Fast)

## Rationale

These rules prevent common AI agent anti-patterns:
- Modifying code without understanding context leads to unintended side effects
- Unbounded retries waste tokens and don't solve root causes
- Silent error swallowing hides bugs that compound over time

## Related

- [Read Before Write](read-before-write.md)
- [Commit Discipline](commit-discipline.md)
- [Parallel Agent Isolation](parallel-agent-isolation.md)
