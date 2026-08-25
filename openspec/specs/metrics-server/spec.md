## Purpose

HTTP server exposing Prometheus metrics and a health endpoint, with graceful shutdown support.

## Requirements

### Requirement: Prometheus metrics endpoint
The metrics server SHALL serve Prometheus metrics on the configured listen address at `/metrics`.

#### Scenario: Scrape endpoint
- **WHEN** Prometheus scrapes `GET /metrics`
- **THEN** all registered collector metrics are returned in Prometheus exposition format

### Requirement: Health endpoint
The server SHALL serve a JSON health response at `GET /healthz` containing status, version, and uptime.

#### Scenario: Health check
- **WHEN** `GET /healthz` is called
- **THEN** the response is `{"status":"ok","version":"<version>","uptime":"<duration>"}`

### Requirement: Healthcheck CLI subcommand
The binary SHALL support a `healthcheck` subcommand that sends `GET http://<addr>/healthz` with a 3-second timeout and exits 0 on HTTP 200, 1 otherwise. The default address is `:9594`, overridable via an optional CLI argument.

#### Scenario: Healthy server
- **WHEN** `tentacle healthcheck` is run and the server responds with 200
- **THEN** exit code is 0

#### Scenario: Unreachable server
- **WHEN** `tentacle healthcheck` is run and the server is unreachable
- **THEN** an error is written to stderr and exit code is 1

#### Scenario: Custom address
- **WHEN** `tentacle healthcheck :8080` is run
- **THEN** the healthcheck targets `http://:8080/healthz`

### Requirement: Docker healthcheck
The Dockerfile SHALL define a `HEALTHCHECK` using `["/tentacle", "healthcheck"]` with interval 30s, timeout 10s, start-period 15s, retries 3.

#### Scenario: Container health
- **WHEN** Docker runs the healthcheck
- **THEN** it executes `/tentacle healthcheck` at the configured interval

### Requirement: Version subcommand
The binary SHALL support a `version` subcommand that prints the build version and exits 0.

#### Scenario: Version output
- **WHEN** `tentacle version` is run
- **THEN** the injected version string is printed to stdout

### Requirement: Graceful shutdown
The server SHALL shut down gracefully on SIGTERM or SIGINT with a 5-second shutdown timeout.

#### Scenario: SIGTERM received
- **WHEN** the process receives SIGTERM
- **THEN** the HTTP server stops accepting new connections and existing requests complete within 5 seconds

### Requirement: Scrape instrumentation
Every service SHALL have a single `ServiceCollector` that emits `<namespace>_scrape_duration_seconds{collector}` and `<namespace>_scrape_success{collector}` for each of its sub-collectors on every scrape. Individual sub-collectors SHALL NOT emit these metrics directly.

#### Scenario: Jellyfin collector scrape metrics
- **WHEN** the Jellyfin `ServiceCollector` completes a scrape
- **THEN** `jellyfin_scrape_success{collector="system"}`, `jellyfin_scrape_success{collector="users"}`, and all other sub-collector scrape metrics are emitted by the wrapper

#### Scenario: Sonarr collector scrape metrics
- **WHEN** the Sonarr `ServiceCollector` completes a scrape
- **THEN** `sonarr_scrape_success{collector="system"}` and all other sub-collector scrape metrics are emitted by the wrapper

#### Scenario: One service failure does not affect others
- **WHEN** the Sonarr API is unreachable but Jellyfin responds normally
- **THEN** Sonarr sub-collectors emit `scrape_success=0` while Jellyfin sub-collectors emit their metrics normally
