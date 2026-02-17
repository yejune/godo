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

### Build & Restart
- [HARD] On Dockerfile change: `docker compose build` then request user to run startup command
- [HARD] On docker-compose.yml change: request user to run startup command (auto-reflects)
- [HARD] After startup, verify service status: `docker compose ps` for container running + healthcheck passing
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
