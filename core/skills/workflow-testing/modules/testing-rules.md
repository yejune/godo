# Testing Mandatory Rules [HARD]

## Scope
- [HARD] These rules apply to **testable code** (business logic, APIs, data layer)
- [HARD] Non-testable changes (CSS, config files, docs, hook scripts, etc.) use alternative verification:
  - Build check: `go build`, `npm run build`, `docker compose config`
  - Manual check: browser verification, CLI execution check
  - Syntax check: lint, type check
- [HARD] Checklist items must specify verification method — e.g., `(verify: build check)`, `(test: unit test)`

## Testing Philosophy
- [HARD] Tests verify **behavior** — test what the code does, not how it does it
- [HARD] Follow FIRST principles: Fast, Independent, Repeatable, Self-validating, Timely
- [HARD] Follow Test Pyramid: Unit > Integration > E2E ratio — but use Real DB at all layers
- [HARD] Flaky tests are bugs — investigate root cause immediately, never ignore by re-running
- [HARD] Mutation testing mindset: "If I change this one line, will a test fail?" — if not, tests are insufficient

## Tests Must Pass
- [HARD] All tests must pass — no skip, no timeout bypass
- [HARD] Never report "did not complete due to timeout" — must resolve
- [HARD] No mocking or skipping — skipped tests must be implemented
- [HARD] On test failure, fix the code — never modify the test to pass

## AI Anti-Patterns Forbidden [HARD]
- [HARD] No assertion weakening — don't change `assertEqual` to `assertContains`, don't replace exact values with `any()`
- [HARD] No try/catch error swallowing — don't catch errors to make tests green
- [HARD] No adjusting expected values to match wrong output — fix the code instead
- [HARD] No `time.sleep()` / arbitrary delays — find the real cause of timing issues
- [HARD] No deleting/commenting out failing tests
- [HARD] No wildcard matchers (`any()`, `mock.ANY`) when exact values are known
- [HARD] No happy-path-only testing — error paths, edge cases, boundary values are mandatory

## Test Data Management [HARD]
- [HARD] Each test sets up and cleans its own data (Arrange-Act-Assert / Given-When-Then)
- [HARD] Cleanup via transaction rollback or truncate — no shared fixture dependencies
- [HARD] Test data is minimal — only what the specific test needs, nothing more
- [HARD] Use factory/builder patterns for test data creation — no raw SQL every time
- [HARD] Use unique identifiers per test (UUID/timestamp suffix) — prevent conflicts in parallel execution

## Test Quality
- [HARD] Test names describe the scenario: `test_login_fails_with_expired_token` — not `test_login_2`
- [HARD] One logical assertion per test (multiple asserts for same verification are allowed)
- [HARD] Failure messages clearly explain **what** went wrong and why
- [HARD] Test code follows same quality standards as production code — remove duplication, maintain readability

## Real DB Only
- [HARD] Database uses real queries only — no mock DB, in-memory DB, SQLite substitution
- [HARD] Tests connect to real DB from Docker Compose services
- [HARD] Only external APIs may be mocked — data layer is never mocked

## Parallelism & Time Limits
- [HARD] Tests must be designed for parallel execution safety (concurrency-safe)
- [HARD] Concurrency/Lock/Race condition tests are separated into methods within the same file for sequential execution
- [HARD] Sequential tests must not interfere with each other's transactions
- [HARD] For parallel tests, Docker Compose local environment/container info must be provided to agents
- [HARD] Individual test < 3 minutes, full suite < 10 minutes — split files/improve process if exceeding
- [HARD] Resource isolation between tests: DB schemas, ports, file paths — prevent shared resource conflicts

## Execution Order
- [HARD] Run individual tests first, then full suite when confident
- [HARD] Cross-verification required: view ↔ DB ↔ worker ↔ model ↔ business logic ↔ controller
- [HARD] Understand business logic and verify correct tests are written
- [HARD] Write tests that find problems and contribute to improvement — not tests for the sake of testing

## DB Transactions
- [HARD] Confirm with user (AskUserQuestion) whether logically separate DBs are on the same server
- [HARD] If same server: service uses X_DB → transaction on X_DB, reference via Z_DB.table
