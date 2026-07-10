## ADDED Requirements

### Requirement: Ordered verification pipeline
The system SHALL run an automated pipeline on every push and pull request
that executes, in strict order: typecheck, then lint, then test, then build —
stopping at the first failing stage.

#### Scenario: All stages pass
- **WHEN** a commit that compiles cleanly, passes lint rules, passes all
  tests, and builds successfully is pushed
- **THEN** the pipeline reports all four stages (typecheck, lint, test,
  build) as successful

#### Scenario: A failing stage halts later stages
- **WHEN** a commit fails the lint stage (for example)
- **THEN** the pipeline reports lint as failed and does not run the test or
  build stages for that run

### Requirement: Test stage runs against a real database
The system SHALL run the test stage against a real, migrated PostgreSQL
instance provisioned within the CI run itself (not mocked), so migration and
data-access code is exercised as it would run in production.

#### Scenario: Tests run with a live migrated database
- **WHEN** the test stage runs
- **THEN** a Postgres instance is available, all migrations have been applied
  to it, and the test suite connects to it successfully

### Requirement: No deployment stage
The pipeline SHALL NOT include any deployment, publishing, or hosting step —
it verifies buildability only.

#### Scenario: Pipeline completes without deploying anywhere
- **WHEN** the build stage completes successfully
- **THEN** the pipeline run ends with a build artifact produced locally to
  the CI run, and no external service, registry, or host is contacted to
  deploy it
