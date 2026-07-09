# go-agentic-base

A Go HTTP server with a Postgres-backed database layer, containerized for
local development and CI-checked on every push/PR.

## Quick start

```sh
make sandbox-up      # build and start app + db
make sandbox-ports    # find the app's assigned host port
make sandbox-down     # tear down
```

## Development

- Database migrations: `make migrate-up`, `make migrate-down`, `make migrate-create name=<name>`
- Run tests: `go test ./...`

See [AGENTS.md](AGENTS.md) for the full agent/worktree workflow, including
sandbox isolation and the pre-commit secret-scanning hook.
