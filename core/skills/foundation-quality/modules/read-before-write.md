# Read Before Write

Mandatory pre-coding behavior rules ensuring agents understand existing code before making changes.

## Rules [HARD]

- [HARD] Read existing code before writing any new code -- writing without reading is prohibited
- [HARD] Understand existing patterns first: naming conventions, error handling style, project structure
- [HARD] Check if similar functionality already exists -- prevent duplicate creation
- [HARD] Prefer modifying existing files over creating new ones -- prevent file bloat
- [HARD] New abstractions require justification -- three similar lines of code is better than premature abstraction (YAGNI)

## Practical Application

Before any code change:
1. Read the target file and surrounding files to understand context
2. Identify naming conventions, error handling patterns, and architectural style
3. Search for existing utilities/helpers that solve the same problem
4. Only then begin implementation, following discovered patterns

## Rationale

Agents that skip reading tend to:
- Introduce inconsistent naming and style
- Create duplicate utilities that already exist
- Break implicit contracts between modules
- Over-engineer solutions that don't fit the codebase

## Related

- [Coding Discipline](coding-discipline.md)
- [Commit Discipline](commit-discipline.md)
