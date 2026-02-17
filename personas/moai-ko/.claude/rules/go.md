---
paths:
  - "**/*.go"
  - "**/go.mod"
  - "**/go.sum"
---

# Go Development Rules

## Tooling
- Format: gofmt
- Lint: go vet
- Test: go test -race ./...
- Coverage target: 85%+

## Code Rules
- context.Context as first parameter
- Explicit error handling (no blank identifier _)
- defer for cleanup (file handles, locks)
- errgroup for concurrent operations
- No global mutable state
- No panic for normal error flow

## File Structure
- cmd/godo/ - CLI entry points
- internal/ - private packages
- Test files: *_test.go alongside source

## Testing
- Table-driven tests
- t.Parallel() for independent tests
- testify/assert for assertions
- Test edge cases and error paths

## Naming
- Exported: PascalCase
- Unexported: camelCase
- Acronyms: consistent case (URL, HTTP, ID)
- Interfaces: -er suffix (Reader, Writer)
