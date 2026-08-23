## Purpose

Environment-variable-based configuration loading for the tentacle Jellyfin exporter, controlling connection settings, timeouts, logging, and listen address.

## Requirements

### Requirement: Environment variable configuration
The config SHALL be loaded exclusively from environment variables with the `TENTACLE_` prefix. The loader SHALL accept a `getenv func(string) string` parameter for testability.

#### Scenario: Testable without real env vars
- **WHEN** `Load(func(string) string)` is called with a custom getenv function
- **THEN** config reads values from the provided function, not `os.Getenv`

### Requirement: Required variables
`TENTACLE_JELLYFIN_ADDRESS` and `TENTACLE_JELLYFIN_TOKEN` SHALL be required. If either is missing, config loading SHALL return an error.

#### Scenario: Missing Jellyfin address
- **WHEN** `TENTACLE_JELLYFIN_ADDRESS` is not set
- **THEN** config loading returns an error "TENTACLE_JELLYFIN_ADDRESS is required"

#### Scenario: Missing Jellyfin token
- **WHEN** `TENTACLE_JELLYFIN_TOKEN` is not set
- **THEN** config loading returns an error "TENTACLE_JELLYFIN_TOKEN is required"

#### Scenario: Both required vars present
- **WHEN** `TENTACLE_JELLYFIN_ADDRESS=http://jellyfin:8096` and `TENTACLE_JELLYFIN_TOKEN=abc123` are set
- **THEN** config loads successfully

### Requirement: Listen address
`TENTACLE_LISTEN_ADDRESS` (default `:9594`) SHALL configure the HTTP server listen address.

#### Scenario: Default listen address
- **WHEN** `TENTACLE_LISTEN_ADDRESS` is not set
- **THEN** `config.ListenAddress` is `:9594`

#### Scenario: Custom listen address
- **WHEN** `TENTACLE_LISTEN_ADDRESS=:8080` is set
- **THEN** `config.ListenAddress` is `:8080`

### Requirement: Scrape timeout
`TENTACLE_SCRAPE_TIMEOUT` (default `10s`) SHALL set the per-collector scrape timeout. The value SHALL accept both plain integer seconds and Go duration format.

#### Scenario: Default timeout
- **WHEN** `TENTACLE_SCRAPE_TIMEOUT` is not set
- **THEN** `config.ScrapeTimeout` is 10 seconds

#### Scenario: Integer seconds
- **WHEN** `TENTACLE_SCRAPE_TIMEOUT=30` is set
- **THEN** `config.ScrapeTimeout` is 30 seconds

#### Scenario: Duration format
- **WHEN** `TENTACLE_SCRAPE_TIMEOUT=1m` is set
- **THEN** `config.ScrapeTimeout` is 60 seconds

#### Scenario: Invalid duration
- **WHEN** `TENTACLE_SCRAPE_TIMEOUT=banana` is set
- **THEN** config loading returns an error

### Requirement: Log level
`TENTACLE_LOG_LEVEL` (default `info`) SHALL support `debug`, `info`, `warn`/`warning`, `error`.

#### Scenario: Default log level
- **WHEN** `TENTACLE_LOG_LEVEL` is not set
- **THEN** slog level is `slog.LevelInfo`

#### Scenario: Debug logging
- **WHEN** `TENTACLE_LOG_LEVEL=debug` is set
- **THEN** slog level is `slog.LevelDebug`

#### Scenario: Warning alias
- **WHEN** `TENTACLE_LOG_LEVEL=warning` is set
- **THEN** slog level is `slog.LevelWarn`

#### Scenario: Invalid log level
- **WHEN** `TENTACLE_LOG_LEVEL=trace` is set
- **THEN** config loading returns an error

### Requirement: Log format
`TENTACLE_LOG_FORMAT` (default `json`) SHALL support `json` and `text`.

#### Scenario: Default log format
- **WHEN** `TENTACLE_LOG_FORMAT` is not set
- **THEN** `config.LogFormat` is `json`

#### Scenario: Text format
- **WHEN** `TENTACLE_LOG_FORMAT=text` is set
- **THEN** `config.LogFormat` is `text` and slog uses `TextHandler`

#### Scenario: Invalid log format
- **WHEN** `TENTACLE_LOG_FORMAT=yaml` is set
- **THEN** config loading returns an error
