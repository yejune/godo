# Docker and Environment Rules [HARD]

## Docker Required

### Core Principles
- [HARD] All projects must be Dockerized (docker-compose.yml required)
- [HARD] Avoid creating Dockerfiles — use official images (`node:20`, `python:3.13`, `golang:1.23`, etc.) + docker-compose `command:`
- [HARD] Dockerfile needed only when: system package installation, multi-stage builds, etc. where official images alone are insufficient
- [HARD] All code runs inside Docker containers — commands target containers
- [HARD] Docker Compose is the Single Source of Truth for local development environment
- [HARD] Host and container connected via volume mounts — code changes reflect immediately in container

### Command Execution
- [HARD] Commands on running container: `docker compose exec <service> <command>`
- [HARD] One-off commands (migrations, seeds, etc.): `docker compose run --rm <service> <command>`
- [HARD] Check container status with `docker compose ps` before running commands — if stopped, `docker compose up -d` first
- [HARD] For services with healthcheck defined, wait for healthy status before running tests/migrations

### Network -- bootapp Domain-Based
- [HARD] No external port exposure needed -- bootapp provides domain-based routing
- [HARD] `ports:` section in docker-compose.yml can be omitted -- access via domain
- [HARD] Service access via DOMAIN environment variable (e.g., `DOMAIN: app.test`, `DOMAIN: db.test`)
- [HARD] SSL_DOMAINS environment variable required -- HTTPS domain communication without ports, auto certificate generation/trust
- [HARD] Host access: `https://app.test` (NOT `localhost:8080`)
- [HARD] Inter-container communication: use Docker service name or DOMAIN value (e.g., `db.test`, `redis.test`)
- [HARD] Never access other services via `localhost` from inside container code
- [HARD] TLD should be `.test` (RFC 2606 reserved) -- `.local` forbidden (macOS mDNS conflict)

### Build & Restart
- [HARD] `docker bootapp up/down` modifies `/etc/hosts`, so agents must NEVER run it directly -- inform the user of the command to run and request execution (AskUserQuestion)
- [HARD] Project start: request user to run `docker bootapp up` -- subnet allocation, domain registration, SSL auto-handling
- [HARD] Project end: request user to run `docker bootapp down` -- domain/hosts cleanup
- [HARD] On Dockerfile change: `docker compose build` then request user to run `docker bootapp up`
- [HARD] On docker-compose.yml change: request user to run `docker bootapp up` (auto-reflects)
- [HARD] After bootapp up, verify service status: `docker compose ps` for container running + healthcheck passing
- [HARD] On dependency addition (package.json, go.mod, etc.): install inside container or rebuild image

## Environment Variable Management

### 12-Factor Principles
- [HARD] Configuration via environment variables — never hardcode in source
- [HARD] Connection strings, API URLs must be read from environment variables
- [HARD] Defaults target container environment (e.g., `DB_HOST=postgres`, NOT `DB_HOST=localhost`)

### docker-compose.yml Environment Variables
- [HARD] General env vars defined directly in docker-compose.yml `environment` section
- [HARD] Group by service concern — DB-related, cache-related, app settings, etc. separated by comments
- [HARD] Only secrets (AWS keys, DB passwords, etc.) separated — Docker secrets or external injection (AWS SSM, Vault, etc.)

### Forbidden
- [HARD] Never create `.env` auto-load file (the `.env` docker-compose reads implicitly)
- [HARD] Never create `.env.local`, `.env.development`, `.env.production` files

### Secret Injection (env_file Allowed)
- env_file: directive is allowed **for secrets only** — local dev secret substitution
- Local: `env_file:` → Production: AWS SSM / Vault external injection — app code is identical (idempotent)
- env_file target files must be in `.gitignore` — never commit
- General env vars still defined directly in `environment:` section

## Code Quality

### Read Before Write
- [HARD] Before writing code, understand existing patterns — file structure, naming, error handling style
- [HARD] New code follows project's existing conventions — no new patterns without justification
- [HARD] Reuse existing utilities/helpers if similar ones exist — no duplicate file creation

### Syntax Check Required
- [HARD] After writing/modifying code, run language-specific syntax checks:
  - Go: `go build ./...` or `go vet ./...`
  - TypeScript/JS: `npm run lint` or `npx tsc --noEmit`
  - Rust: `cargo check`
  - Python: `ruff check` or `flake8`
- [HARD] Syntax check prefers container — if tools not installed in production image, allow host execution (assuming volume mount shares source)

### Dependency Management
- [HARD] Before adding new dependency, check if existing dependencies can solve it
- [HARD] Leverage existing knowledge and experience first — search and reference, document new discoveries

## AI Agent Anti-Patterns

- [HARD] Dependency installation runs on host -- auto-reflected in container via volume mount (`npm install`, `pip install`, `go mod download`, etc.)
- [HARD] Syntax check/dev tools also allowed on host platform-specifically (lint, formatter, type checker, etc.)
- [HARD] Never use `localhost` for inter-container communication -- use Docker service name or domain
- [HARD] No port exposure via `ports:` mapping -- use bootapp domains
- [HARD] Never ignore Docker Compose healthcheck status -- wait for healthy before subsequent work
- [HARD] Never create `.env` family files -- forbidden for any reason
- [HARD] Never run tests outside containers -- use `docker compose exec`
- [HARD] Never enter container shell (`docker exec -it ... bash/sh`) -- always use `docker compose exec <service> <command>` single commands
- [HARD] `docker cp` forbidden -- never copy files to container, share via volume mount
- [HARD] Never create temp files/scripts inside container -- all changes on host source, reflected via volume
- [HARD] Idempotency required -- same command multiple times yields same result, never depend on temporary state
- [HARD] Never `COPY` source code in Dockerfile -- both source and dependencies shared via volume mount
