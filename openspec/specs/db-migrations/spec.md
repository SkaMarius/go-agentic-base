# db-migrations Specification

## Purpose

TBD - created by archiving change initial-setup. Update Purpose after archive.

## Requirements

### Requirement: Versioned schema migrations
The system SHALL manage its PostgreSQL schema through versioned, ordered
`golang-migrate` SQL migration files stored in a `migrations/` directory,
each with a paired `up` and `down` script.

#### Scenario: Applying all migrations from an empty database
- **WHEN** the migration tool is run against a fresh, empty Postgres database
- **THEN** all migrations apply in order and the database schema matches the
  latest migration version, with no errors

#### Scenario: Rolling back the latest migration
- **WHEN** the migration tool is run with a "down" (rollback) command against
  a migrated database
- **THEN** the most recently applied migration's `down` script executes
  successfully and the schema reverts to the prior version

### Requirement: Example migration demonstrates the pattern
The system SHALL include one example migration that creates a table, so the
migration mechanism is demonstrated end-to-end rather than left empty.

#### Scenario: Example table exists after migrating
- **WHEN** all migrations are applied to a fresh database
- **THEN** the example table defined by the initial migration exists and is
  queryable

### Requirement: Database connectivity via configuration
The system SHALL connect to PostgreSQL using a connection string supplied via
the `DATABASE_URL` environment variable, with no hard-coded credentials.

#### Scenario: Server connects using DATABASE_URL
- **WHEN** the server starts with a valid `DATABASE_URL` environment variable set
- **THEN** it establishes a working connection pool to that Postgres instance

#### Scenario: Server fails fast on invalid database configuration
- **WHEN** the server starts with a missing or invalid `DATABASE_URL`
- **THEN** it fails to start with a clear error rather than starting in a
  degraded/unclear state
