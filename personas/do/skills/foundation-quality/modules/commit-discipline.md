# Commit Discipline

Rules for atomic, meaningful commits that maintain codebase integrity.

## Rules [HARD]

- [HARD] Atomic commits: one logical change = one commit
- [HARD] Commit message explains WHY -- WHAT is shown by the diff
- [HARD] Never commit: .env files, credentials, large binaries, generated files
- [HARD] Pre-commit self-check: "Does this diff contain only intended changes?"
- [HARD] AI footer controlled by DO_AI_FOOTER environment variable:
  - DO_AI_FOOTER=true: append AI-generated footer to commit message
  - DO_AI_FOOTER=false (default): no footer
- [HARD] These rules apply identically to development agents and manager-git -- Single Source of Truth

## Commit Message Format

```
type: what was done (50 chars max)

Why it was done and how (optional body)
```

Types: feat, fix, refactor, docs, test, chore

The commit message and diff together must fully convey the modification intent.

## Anti-Patterns

- Committing multiple unrelated changes in one commit
- Writing "fix bug" without explaining which bug or why
- Including debug code, temporary files, or commented-out code
- Committing without running tests first

## Related

- [Coding Discipline](coding-discipline.md)
- [Parallel Agent Isolation](parallel-agent-isolation.md)
