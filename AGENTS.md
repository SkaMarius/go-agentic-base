# Agent workflow

This project is designed so that multiple agents (or humans) can work on
separate branches in parallel, each with its own isolated sandbox — no
shared containers, networks, volumes, or ports.

## One-time setup: install the pre-commit hook

Install [pre-commit](https://pre-commit.com/) and register the hook that
blocks commits containing secrets (see [Secret scanning](#secret-scanning)):

```sh
brew install pre-commit   # or: pip install pre-commit
make install-hooks
```

This only needs to be done once per clone of the repo (git hooks live in the
shared `.git` directory and apply across all worktrees).

## Working in a worktree

1. **Create a worktree** for the branch you're working on:

   ```sh
   git worktree add ../go-agentic-base-<name> -b feature/<name>
   cd ../go-agentic-base-<name>
   ```

2. **Start an isolated sandbox** from that worktree:

   ```sh
   make sandbox-up
   ```

   This runs `docker compose up --build -d`. Compose namespaces containers,
   networks, and volumes by project name, which defaults to the worktree's
   directory name — so sandboxes from different worktrees never collide.
   The app's port is published without a fixed host port, so Docker assigns
   a free one automatically.

3. **Find the sandbox's port**:

   ```sh
   make sandbox-ports
   ```

   Use the returned host port to reach the app, e.g.
   `curl http://localhost:<port>/health`.

4. **Tear down** the sandbox and remove the worktree when done:

   ```sh
   make sandbox-down
   cd ..
   git worktree remove go-agentic-base-<name>
   ```

## Secret scanning

Commits are scanned locally by the pre-commit hook (`gitleaks protect
--staged`, installed via `make install-hooks` above) and again in CI (the
`secret-scan` job in `.github/workflows/ci.yml`). Both are pinned to the same
`gitleaks` version so behavior matches. Known-safe strings that gitleaks
would otherwise flag can be excluded via the allowlist in `.gitleaks.toml`.
