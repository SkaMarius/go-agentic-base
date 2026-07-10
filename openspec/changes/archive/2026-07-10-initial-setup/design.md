## Context

`go-agentic-base` is currently an empty repository (only OpenSpec scaffolding
exists). This change establishes the first real code: a minimal but complete
Go service base that future features build on top of. Key constraints:

- No deployment target exists or is planned as part of this change — CI must
  fully verify the code without shipping it anywhere.
- The repo is expected to be worked on by multiple agents concurrently, each
  in its own `git worktree`, and each needs to be able to run the full stack
  (app + Postgres) without colliding with another agent's running stack.
- Nothing here has a predecessor to be compatible with — this is a greenfield
  base, so decisions favor simplicity and convention over flexibility.

## Goals / Non-Goals

**Goals:**
- A Go HTTP service that boots, serves `GET /health`, and is easy to extend.
- Real Postgres connectivity with versioned, reversible schema migrations.
- Local dev via `docker compose up`, with zero manual port bookkeeping when
  multiple worktrees run the stack simultaneously.
- CI that proves the code compiles, lints clean, passes tests (against a real
  Postgres instance), and builds — on every push/PR.
- `AGENTS.md` that tells an agent exactly how to get an isolated sandbox
  running from a fresh worktree.

**Non-Goals:**
- No production deployment, hosting, or infra-as-code (this change explicitly
  stops before any "ship it" step).
- No authentication, business domain models, or API beyond `/health`.
- No orchestration of multiple agents themselves — only the environment
  isolation (worktree + compose project) they rely on.
- No Kubernetes/Helm — Docker Compose is sufficient for this stage.

## Decisions

**1. Router: `go-chi/chi` (v5)**
Idiomatic, `net/http`-compatible, minimal, well-maintained, good middleware
ecosystem for when routes grow beyond `/health`. Alternative considered:
plain `net/http` `ServeMux` — viable for one route today, but chi avoids a
rewrite the moment a second route or middleware (logging, recovery) is
needed, at negligible cost now.

**2. Migrations: `golang-migrate/migrate`**
Plain SQL `up`/`down` files, usable both as a CLI (for local/CI) and as a Go
library (if the app needs to self-migrate on boot later). Alternative
considered: `goose` — similar, but also supports Go-code migrations, which
this project doesn't need; plain SQL keeps migrations reviewable and
DB-engine-native.

**3. Postgres driver: `jackc/pgx` (v5) via its `database/sql` stdlib adapter**
`pgx` is the actively maintained, higher-performance option; `lib/pq` is in
maintenance mode. Using it through `database/sql` (rather than pgx's native
API) keeps the base project idiomatic and swappable — code depends on
`database/sql` interfaces, not a driver-specific API.

**4. Project layout**
```
cmd/server/main.go        - entrypoint: load config, connect DB, start HTTP server
internal/server/          - router + handlers (health.go, router.go)
internal/config/          - env-based config loading
internal/db/              - connection + migration-runner helpers
migrations/               - golang-migrate SQL files (*.up.sql / *.down.sql)
```
Standard, unsurprising Go layout; nothing under `internal/` is importable
outside this module, which is correct for a base project with no public API
yet.

**5. Configuration: environment variables only**
`PORT` (default `8080`), `DATABASE_URL`. No config files or flag parsing —
there's nothing yet that needs more than two settings. Add structure later
if/when config grows.

**6. Containerization**
Multi-stage `Dockerfile` (Go build stage → minimal runtime stage) for a small
final image. `docker-compose.yml` defines `app` + `db` (`postgres:16-alpine`),
with a healthcheck on `db` so `app` waits for real readiness rather than a
fixed sleep.

**7. Per-worktree sandbox isolation — ephemeral host ports, no manual offsets**
Docker Compose already scopes container/network/volume names by *project
name*, which defaults to the current directory's basename. Since every
`git worktree add` creates a distinctly named directory, running
`docker compose up` from inside a worktree is *already* isolated at the
container/network/volume level — no extra config needed.
The one remaining collision point is **host port bindings**. Resolution:
declare ports in compose as container-only (e.g. `"8080"` instead of
`"8080:8080"`), which tells Docker to bind to a random free host port per
project. `docker compose port app 8080` (wrapped in a `make sandbox-ports`
target) reveals the actual assigned port for that worktree's stack.
Alternative considered: computing a deterministic port offset from the
worktree path (e.g. hash → `8080 + offset`) — rejected as unnecessary
bookkeeping and a source of subtle collisions (hash collisions, running out
of range) when Docker already solves this natively.

**8. CI pipeline order and jobs**
Single GitHub Actions workflow. `typecheck`, `lint`, `test`, and `secret-scan`
run as independent parallel jobs (none of them depend on each other's
output); `build` depends on all four succeeding:
- **typecheck** — `go build ./...` (compilation is Go's typecheck) and
  `go vet ./...`.
- **lint** — `golangci-lint run`.
- **test** — `go test ./... -race -cover`, against a real Postgres brought
  up via the workflow's `services:` block, migrated with the `migrate` CLI
  before tests run.
- **secret-scan** — `gitleaks detect` over the incoming push/PR diff (see
  Decision 10 below); independent of the Go toolchain entirely, so it gains
  nothing from being sequenced after the others.
- **build** (depends on the above four) — `go build -o bin/server
  ./cmd/server`, uploaded as a build artifact (proof of a shippable binary)
  — nothing is deployed anywhere.

```
secret-scan ─┐
typecheck ───┼─▶ build
lint ─────────┤
test ─────────┘
```

No stage is more than "run once and check exit code" — no external
credentials or deploy targets are required, so this fully addresses the
"can this exist without deploying anywhere?" question from the proposal.

**9. `AGENTS.md` workflow**
Documents: creating a worktree (`git worktree add ../go-agentic-base-<id>
<branch>`), starting an isolated sandbox from inside it (`make sandbox-up`),
finding its assigned ports (`make sandbox-ports`), and tearing down both the
compose stack and the worktree when done (`make sandbox-down` +
`git worktree remove`).

**10. Secret scanning: `gitleaks`**
Chosen over `trufflehog` for speed and simplicity: `gitleaks` is
pattern/entropy-based with no outbound network calls, has a mature GitHub
Action (`gitleaks/gitleaks-action`) and a pre-commit hook, and needs no
per-provider verification logic. Three layers, each catching leaks at a
different point:
1. **Incremental CI gate** — runs on every push/PR, scoped to the
   diff/commit range, blocking merge on any finding.
2. **One-time full-history baseline audit** — run once, now, while the repo
   has no real code yet, so there's a confirmed-clean history before any
   accumulates. Not a recurring CI job — a task done as part of this change.
3. **Local pre-commit hook** — installed via the same workflow documented in
   `AGENTS.md`, catching a leak before it's committed at all, in either the
   main checkout or any agent's worktree.
Alternative considered: `trufflehog`, for its live-verification of found
secrets (confirming a credential actually still works, not just that it
matches a pattern) — rejected for this base project on the grounds that it
adds outbound network calls from CI/pre-commit to third-party APIs using
whatever string it finds, which is a heavier and more surprising default
than a new base project should ship with. Worth revisiting if false
positives from pattern-matching become a recurring problem.

## Risks / Trade-offs

- **[Risk]** Ephemeral host ports mean `localhost:8080` won't reliably work
  when more than one sandbox is running → **Mitigation**: `make
  sandbox-ports` Makefile target wraps `docker compose port` so the actual
  URL is always one command away; documented in `AGENTS.md`.
- **[Risk]** Running a real Postgres service container in CI adds startup
  time and a small flakiness surface → **Mitigation**: use Postgres's
  built-in healthcheck and a wait step before running migrations/tests.
- **[Risk]** `golang-migrate` CLI must be available wherever migrations run
  (dev machine, CI, sandbox) → **Mitigation**: install it via a pinned
  version in the CI workflow and document the local install step (or a
  `go run` wrapper) in `AGENTS.md`/Makefile so there's one source of truth
  for the version.
- **[Risk]** `gitleaks` is pattern/entropy-based, not verification-based, so
  it can flag non-secrets (test fixtures, example env values, high-entropy
  strings that aren't credentials) → **Mitigation**: maintain a
  `.gitleaks.toml` allowlist for known-safe patterns/paths; keep it small and
  reviewed rather than broad.
- **[Risk]** A pre-commit hook is one more tool contributors/agents must have
  installed locally → **Mitigation**: document the install step in
  `AGENTS.md` and pin the `gitleaks` version there and in CI so behavior
  matches between local and CI runs.
- **[Risk]** Concurrent `docker compose build` across many worktrees could
  contend for local Docker build cache/resources on one machine →
  **Mitigation**: acceptable at this scale; revisit with BuildKit cache
  scoping if it becomes a real bottleneck.

## Migration Plan

Greenfield — there is no existing system or data to migrate from or roll
back to. "Rollout" here just means implementing the tasks in order (see
`tasks.md`); if something is wrong after merge, revert the commit.

## Open Questions

- Target Go version: assumed to be the latest stable minor (1.23.x) unless
  there's an existing constraint — flag if a different version is required.
- Whether `staticcheck` should be added on top of `golangci-lint` (which can
  include it as one of its linters) — defaulting to enabling it as a
  `golangci-lint` linter rather than a separate tool/step, to keep the "lint"
  stage single-tool.
