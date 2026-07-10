## 1. Project skeleton

- [x] 1.1 Run `go mod init` for the module and pin a Go toolchain version (1.23.x)
- [x] 1.2 Create directory layout: `cmd/server/`, `internal/server/`, `internal/config/`, `internal/db/`, `migrations/`
- [x] 1.3 Add `.gitignore` (binaries, `.env`, local compose overrides)

## 2. HTTP server & health endpoint

- [x] 2.1 Add `go-chi/chi` dependency
- [x] 2.2 Implement `internal/config` — load `PORT` (default `8080`) and `DATABASE_URL` from env
- [x] 2.3 Implement `internal/server` router with `GET /health` returning `200` and body `service is running`
- [x] 2.4 Implement `cmd/server/main.go` — load config, build router, start HTTP server on configured port
- [x] 2.5 Write a handler test for `GET /health` asserting status `200` and exact body

## 3. Database & migrations

- [x] 3.1 Add `jackc/pgx/v5` (stdlib adapter) dependency
- [x] 3.2 Implement `internal/db` connection helper using `DATABASE_URL`, failing fast with a clear error on invalid/missing config
- [x] 3.3 Add `golang-migrate` as the migration tool (CLI usage documented in Makefile, no app dependency required unless self-migration is added later)
- [x] 3.4 Write initial migration pair (`0001_*.up.sql` / `0001_*.down.sql`) creating one example table
- [x] 3.5 Add Makefile targets: `migrate-up`, `migrate-down`, `migrate-create`
- [x] 3.6 Write a test that runs migrations against a test database and asserts the example table exists

## 4. Containerization & sandbox isolation

- [x] 4.1 Write multi-stage `Dockerfile` (build stage → minimal runtime stage) for the server
- [x] 4.2 Write `docker-compose.yml` with `app` + `db` (`postgres:16-alpine`) services
- [x] 4.3 Add a healthcheck to `db` and make `app` depend on `db` being healthy before starting
- [x] 4.4 Configure `app`'s port mapping as container-only (e.g. `"8080"`) so Docker assigns a random host port per compose project
- [x] 4.5 Add Makefile targets: `sandbox-up`, `sandbox-down`, `sandbox-ports` (wraps `docker compose port app 8080`)
- [x] 4.6 Manually verify: start the stack from two different git worktrees at once and confirm both come up with no container/network/volume/port collisions

## 5. CI pipeline

- [x] 5.1 Add `.github/workflows/ci.yml` triggered on `push` and `pull_request`
- [x] 5.2 Add `typecheck` job: `go build ./...` + `go vet ./...`
- [x] 5.3 Add `lint` job (independent of `typecheck`): install and run `golangci-lint`
- [x] 5.4 Add `test` job (independent of `typecheck`/`lint`): start a Postgres `services:` container, wait for health, run migrations, run `go test ./... -race -cover`
- [x] 5.5 Add `secret-scan` job (independent of the others): run `gitleaks detect` over the push/PR diff (see section 6)
- [x] 5.6 Add `build` job depending on `typecheck`, `lint`, `test`, and `secret-scan` all succeeding: `go build -o bin/server ./cmd/server`, upload as a workflow artifact
- [x] 5.7 Confirm no job publishes, deploys, or contacts an external host — pipeline ends at the build artifact

## 6. Secret scanning

- [x] 6.1 Add `.gitleaks.toml` with an initial (empty or minimal) allowlist
- [x] 6.2 Run a one-time full-history baseline audit (`gitleaks detect --source . --log-opts="--all"` or equivalent) and confirm zero findings before proceeding
- [x] 6.3 Wire the `secret-scan` CI job (`gitleaks/gitleaks-action` or equivalent) into `.github/workflows/ci.yml`, scoped to the incoming diff/commit range
- [x] 6.4 Add a local pre-commit hook (`.pre-commit-config.yaml` or a Makefile-wrapped git hook) running `gitleaks protect --staged`
- [x] 6.5 Pin the same `gitleaks` version in both the CI job and the pre-commit hook so behavior matches

## 7. Agent workflow documentation

- [x] 7.1 Write `AGENTS.md`: how to create a worktree (`git worktree add`), start an isolated sandbox from it (`make sandbox-up`), find its port (`make sandbox-ports`), and tear down (`make sandbox-down` + `git worktree remove`)
- [x] 7.2 Document the pre-commit hook install step in `AGENTS.md`
- [x] 7.3 Cross-link `AGENTS.md` from the repo root `README.md` (create a minimal `README.md` if one doesn't exist)

## 8. Verification

- [x] 8.1 Run `docker compose up` locally and confirm `GET /health` returns `200` / `service is running`
- [x] 8.2 Run the full CI workflow locally or via a pushed branch and confirm `typecheck`, `lint`, `test`, and `secret-scan` all pass and `build` runs only after
- [x] 8.3 Confirm the pre-commit hook blocks a deliberately-staged fake secret, and that the allowlist correctly excludes a known-safe test fixture
- [x] 8.4 Run `openspec validate initial-setup --strict` (or equivalent) and fix any structural issues before archiving
