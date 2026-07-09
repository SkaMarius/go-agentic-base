## Why

This repository is currently empty aside from OpenSpec scaffolding. There is no
runnable Go service, no database layer, no containerization, and no CI. Before
any feature work can start, the project needs a working, reusable base: a
service skeleton that boots, connects to a real database, runs the same way in
every environment (local, CI, agent sandbox), and is verified automatically on
every change.

## What Changes

- Add a Go HTTP server (chi router) exposing `GET /health` → `200` with body
  `service is running`, listening on port `8080`.
- Add PostgreSQL integration with `golang-migrate`-based schema migrations,
  including one example migration to demonstrate the pattern end-to-end.
- Add Docker containerization: a `Dockerfile` for the service and a
  `docker-compose.yml` for local dev (app + Postgres).
- Add per-git-worktree sandbox isolation for docker-compose: each worktree runs
  its own compose project (distinct project name, ports, volumes, network) so
  multiple agents/worktrees can run the full stack concurrently without
  colliding.
- Add a GitHub Actions CI pipeline with four parallel gates — typecheck, lint, test, and secret-scan (via gitleaks) — where build runs only once all four succeed. No deployment stage; the pipeline stops at a successful build artifact. Includes a one-time full-history baseline audit (run now, while the repo is empty) and a local pre-commit hook to catch secrets before they leave a machine or worktree.
- Add `AGENTS.md` documenting the git worktree workflow used for agent
  isolation (how to create a worktree, how it maps to an isolated
  docker-compose sandbox, cleanup).
- OpenSpec is already initialized in this repo (`openspec/config.yaml`,
  `changes/`, `specs/`); this change is the first one to actually populate
  `openspec/specs/` with real capability specs.

## Capabilities

### New Capabilities
- `health-check`: HTTP endpoint reporting service liveness (`GET /health`).
- `db-migrations`: Versioned PostgreSQL schema migrations applied via
  `golang-migrate`.
- `dev-containerization`: Docker image and docker-compose stack for local
  development, with isolated per-worktree sandbox instances.
- `ci-pipeline`: Automated typecheck/lint/test/build verification on every
  push and pull request.
- `secret-scanning`: Gitleaks-based detection of committed secrets — as a
  blocking incremental CI gate, a one-time full-history baseline audit, and
  a local pre-commit hook.

### Modified Capabilities
- None — this is the first change in the repo; no existing specs to modify.

## Impact

- **New code**: `cmd/server`, `internal/...` (Go module init), `migrations/`,
  `Dockerfile`, `docker-compose.yml`, `.github/workflows/ci.yml`, `AGENTS.md`,
  `Makefile` (dev/test/migrate/sandbox helper targets).
- **New dependencies**: `go-chi/chi`, `golang-migrate/migrate`, a Postgres
  driver (`lib/pq` or `jackc/pgx`), and CI-only tools (`golangci-lint`).
- **Infra**: Requires Docker + Docker Compose locally and in any agent
  sandbox; no external hosting/deployment target is introduced.
- **Dev workflow**: Establishes the git-worktree-per-agent pattern as the
  supported way to run multiple isolated instances of the stack side by side.
