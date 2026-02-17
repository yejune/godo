# Syntax Check

Mandatory post-modification syntax validation rules.

## Rules [HARD]

- [HARD] After writing/modifying code, always run language-specific syntax checks:
  - Go: `go build ./...` or `go vet ./...`
  - TypeScript/JS: `npm run lint` or `npx tsc --noEmit`
  - Rust: `cargo check`
  - Python: `ruff check` or `flake8`
- [HARD] Prefer running syntax checks in containers -- if tools are not installed in the production image, running on the host is allowed (assuming source is shared via volume mount)

## Dependency Management [HARD]

- [HARD] Before adding new dependencies, check if existing dependencies can solve the problem
- [HARD] Leverage existing knowledge and experience first -- search and reference, document new discoveries

## Rationale

Syntax checks catch errors early before they reach testing or production. Running checks immediately after modification ensures:
- Parse errors are caught at write time, not commit time
- Type errors surface before tests run
- Import/dependency issues are identified immediately

## Related

- [Coding Discipline](coding-discipline.md)
- [Read Before Write](read-before-write.md)
