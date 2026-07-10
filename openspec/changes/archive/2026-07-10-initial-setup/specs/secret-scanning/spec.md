## ADDED Requirements

### Requirement: Incremental secret scan blocks CI on findings
The system SHALL run `gitleaks` over the incoming diff/commit range on every
push and pull request, as an independent job that gates the `build` job, and
SHALL fail the run if any finding is detected that is not covered by the
allowlist.

#### Scenario: Clean push passes the scan
- **WHEN** a push or pull request introduces no secret-like patterns
- **THEN** the `secret-scan` job succeeds and does not block `build`

#### Scenario: A leaked credential blocks the pipeline
- **WHEN** a push or pull request introduces a string matching a known
  secret pattern that is not allowlisted
- **THEN** the `secret-scan` job fails, and the `build` job does not run

### Requirement: One-time full-history baseline audit
The system SHALL include a one-time task that runs `gitleaks` across the
entire git history to confirm a clean baseline before ongoing incremental
scanning begins.

#### Scenario: Baseline audit finds no pre-existing leaks
- **WHEN** the full-history audit is run against the repository as it exists
  at the time of this change
- **THEN** it reports zero findings, establishing a confirmed-clean starting
  point for future incremental scans

### Requirement: Local pre-commit hook
The system SHALL provide a `gitleaks` pre-commit hook, documented in
`AGENTS.md`, so secrets are caught before a commit is created locally or in
any agent worktree.

#### Scenario: Pre-commit hook blocks a local commit
- **WHEN** a developer or agent attempts to commit a change containing a
  string matching a known secret pattern that is not allowlisted
- **THEN** the commit is rejected locally before it is created, with the
  finding reported to the terminal

### Requirement: Allowlist for known false positives
The system SHALL maintain a `.gitleaks.toml` allowlist for reviewed
non-secret matches (e.g. test fixtures, example configuration values), kept
minimal and version-controlled.

#### Scenario: Allowlisted pattern does not block the pipeline or a commit
- **WHEN** a string matches a secret pattern but is covered by an entry in
  `.gitleaks.toml`
- **THEN** neither the CI `secret-scan` job nor the local pre-commit hook
  reports it as a finding
