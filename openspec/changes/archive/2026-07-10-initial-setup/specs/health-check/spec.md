## ADDED Requirements

### Requirement: Health endpoint
The system SHALL expose an HTTP `GET /health` endpoint that returns HTTP
status `200` with a plain-text (or text-compatible) body of exactly
`service is running`, with no authentication required.

#### Scenario: Health check succeeds
- **WHEN** a client sends `GET /health` to the running server
- **THEN** the server responds with status `200` and body `service is running`

### Requirement: Server listens on port 8080
The system SHALL listen for HTTP connections on port `8080` by default.

#### Scenario: Server starts and accepts connections on port 8080
- **WHEN** the server process starts with no port override configured
- **THEN** it accepts HTTP connections on `0.0.0.0:8080`
