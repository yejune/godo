# AI Agent Forbidden Patterns [HARD]

## Prohibited Actions for AI Agents

- [HARD] Dependency installation runs on host — auto-reflected in container via volume mount (`npm install`, `pip install`, `go mod download`, etc.)
- [HARD] Syntax check/dev tools may also be installed on host per platform (lint, formatter, type checker, etc.)
- [HARD] Never use `localhost` for inter-container communication — use Docker service names or domains
- [HARD] Never expose ports with `ports:` mapping — use domain-based routing
- [HARD] Never ignore Docker Compose healthcheck status — wait for healthy before proceeding
- [HARD] Never create `.env` files — forbidden under any circumstances
- [HARD] Never run tests outside containers — use `docker compose exec`
- [HARD] Never enter container shell (`docker exec -it ... bash/sh`) — always use `docker compose exec <service> <command>` single commands
- [HARD] Never use `docker cp` — no copying files to containers, share via volume mount
- [HARD] Never create temporary files/scripts inside containers — all changes on host source, reflected via volume
- [HARD] Idempotency required — running the same command multiple times must produce identical results, never depend on temporary state
- [HARD] Never `COPY` source code in Dockerfile — share both source and dependencies via volume mount
