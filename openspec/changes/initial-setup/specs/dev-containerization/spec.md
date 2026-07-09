## ADDED Requirements

### Requirement: Containerized service image
The system SHALL provide a `Dockerfile` that builds a runnable container
image of the server via a multi-stage build (compile stage + minimal runtime
stage).

#### Scenario: Image builds and runs
- **WHEN** the `Dockerfile` is built into an image and that image is run
- **THEN** the container starts the server process and it serves `GET /health`
  as specified in the `health-check` capability

### Requirement: Local dev stack via Docker Compose
The system SHALL provide a `docker-compose.yml` that starts the application
container and a PostgreSQL container together, with the application waiting
for the database to be ready (via healthcheck) before accepting traffic.

#### Scenario: Full stack starts with one command
- **WHEN** a developer runs `docker compose up` in the repository
- **THEN** both the `app` and `db` services start, `db` becomes healthy
  before `app` is considered ready, and `GET /health` succeeds against `app`

### Requirement: Isolated per-worktree sandbox
The system SHALL allow multiple Docker Compose stacks — one per git worktree
— to run concurrently on the same machine without container name, network,
volume, or host-port collisions.

#### Scenario: Two worktrees run sandboxes simultaneously
- **WHEN** a developer starts the compose stack from two different git
  worktree directories at the same time
- **THEN** both stacks start successfully, each with its own containers,
  network, and volumes, and neither stack's host ports conflict with the
  other's

#### Scenario: Resolving a sandbox's assigned port
- **WHEN** a developer runs the sandbox-ports helper (`make sandbox-ports`)
  from within a worktree whose stack is running
- **THEN** it prints the host port currently mapped to that worktree's `app`
  container, allowing the developer to reach it directly
