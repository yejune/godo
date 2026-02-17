# Bootapp Infrastructure Rules [HARD]

## Network -- bootapp Domain-Based
- [HARD] No external port exposure needed -- bootapp provides domain-based routing
- [HARD] `ports:` section in docker-compose.yml can be omitted -- access via domain
- [HARD] Service access via DOMAIN environment variable (e.g., `DOMAIN: app.test`, `DOMAIN: db.test`)
- [HARD] SSL_DOMAINS environment variable required -- HTTPS domain communication without ports, auto certificate generation/trust
- [HARD] Host access: `https://app.test` (NOT `localhost:8080`)
- [HARD] Inter-container communication: use Docker service name or DOMAIN value (e.g., `db.test`, `redis.test`)
- [HARD] Never access other services via `localhost` from inside container code
- [HARD] TLD should be `.test` (RFC 2606 reserved) -- `.local` forbidden (macOS mDNS conflict)

## Build & Restart -- bootapp Commands
- [HARD] `docker bootapp up/down` modifies `/etc/hosts`, so agents must NEVER run it directly -- inform the user of the command to run and request execution (AskUserQuestion)
- [HARD] Project start: request user to run `docker bootapp up` -- subnet allocation, domain registration, SSL auto-handling
- [HARD] Project end: request user to run `docker bootapp down` -- domain/hosts cleanup
- [HARD] On Dockerfile change: `docker compose build` then request user to run `docker bootapp up`
- [HARD] On docker-compose.yml change: request user to run `docker bootapp up` (auto-reflects)
- [HARD] After bootapp up, verify service status: `docker compose ps` for container running + healthcheck passing
- [HARD] On dependency addition (package.json, go.mod, etc.): install inside container or rebuild image
