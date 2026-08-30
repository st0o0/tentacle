## Purpose

Environment-variable-based configuration loading for the tentacle multi-service exporter, controlling per-service connection settings, timeouts, logging, and listen address.

## Requirements

### Requirement: Environment variable configuration
The config SHALL be loaded exclusively from environment variables with the `TENTACLE_` prefix. The loader SHALL accept a `getenv func(string) string` parameter for testability. The config SHALL support per-service configuration blocks, where each service has its own address and token variables.

#### Scenario: Testable without real env vars
- **WHEN** `Load(func(string) string)` is called with a custom getenv function
- **THEN** config reads values from the provided function, not `os.Getenv`

### Requirement: Required variables
No single service SHALL be globally required. At least one service MUST be configured (address + token pair set), otherwise config loading SHALL return an error.

#### Scenario: No services configured
- **WHEN** no service address/token pairs are set
- **THEN** config loading returns an error indicating at least one service must be configured

#### Scenario: Only Jellyfin configured
- **WHEN** only `TENTACLE_JELLYFIN_ADDRESS` and `TENTACLE_JELLYFIN_TOKEN` are set
- **THEN** config loads successfully with only the Jellyfin service enabled

#### Scenario: Multiple services configured
- **WHEN** Jellyfin and Sonarr address/token pairs are both set
- **THEN** config loads successfully with both services enabled

### Requirement: Per-service configuration
Each service SHALL have its own config block with address and token variables. A service is enabled when both its `TENTACLE_<SERVICE>_ADDRESS` and `TENTACLE_<SERVICE>_TOKEN` are set. If only one of the pair is set, config loading SHALL return an error.

#### Scenario: Jellyfin service config
- **WHEN** `TENTACLE_JELLYFIN_ADDRESS=http://jellyfin:8096` and `TENTACLE_JELLYFIN_TOKEN=abc` are set
- **THEN** the Jellyfin service config is populated and enabled

#### Scenario: Sonarr service config
- **WHEN** `TENTACLE_SONARR_ADDRESS=http://sonarr:8989` and `TENTACLE_SONARR_TOKEN=abc` are set
- **THEN** the Sonarr service config is populated and enabled

#### Scenario: Radarr service config
- **WHEN** `TENTACLE_RADARR_ADDRESS=http://radarr:7878` and `TENTACLE_RADARR_TOKEN=abc` are set
- **THEN** the Radarr service config is populated and enabled

#### Scenario: Prowlarr service config
- **WHEN** `TENTACLE_PROWLARR_ADDRESS=http://prowlarr:9696` and `TENTACLE_PROWLARR_TOKEN=abc` are set
- **THEN** the Prowlarr service config is populated and enabled

#### Scenario: Audiobookshelf service config
- **WHEN** `TENTACLE_AUDIOBOOKSHELF_ADDRESS=http://abs:13378` and `TENTACLE_AUDIOBOOKSHELF_TOKEN=abc` are set
- **THEN** the Audiobookshelf service config is populated and enabled

#### Scenario: Seerr service config
- **WHEN** `TENTACLE_SEERR_ADDRESS=http://seerr:5055` and `TENTACLE_SEERR_TOKEN=abc` are set
- **THEN** the Seerr service config is populated and enabled

#### Scenario: Partial service config
- **WHEN** `TENTACLE_SONARR_ADDRESS` is set but `TENTACLE_SONARR_TOKEN` is not
- **THEN** config loading returns an error indicating the token is required when address is set

### Requirement: Service-independent settings
`TENTACLE_LISTEN_ADDRESS`, `TENTACLE_SCRAPE_TIMEOUT`, `TENTACLE_LOG_LEVEL`, `TENTACLE_LOG_FORMAT` SHALL remain global settings shared across all services. Their defaults and behavior SHALL not change.

#### Scenario: Global settings with multiple services
- **WHEN** `TENTACLE_SCRAPE_TIMEOUT=30s` is set with Jellyfin and Sonarr configured
- **THEN** both services use the same 30-second scrape timeout

### Requirement: Cache intervals
`TENTACLE_CACHE_WARM_INTERVAL` (default `10m`) and `TENTACLE_CACHE_COLD_INTERVAL` (default `30m`) SHALL configure the background refresh intervals for cached collectors. Both SHALL accept plain integer seconds and Go duration format.

#### Scenario: Default warm interval
- **WHEN** `TENTACLE_CACHE_WARM_INTERVAL` is not set
- **THEN** `config.CacheWarmInterval` is 10 minutes

#### Scenario: Default cold interval
- **WHEN** `TENTACLE_CACHE_COLD_INTERVAL` is not set
- **THEN** `config.CacheColdInterval` is 30 minutes

#### Scenario: Custom warm interval
- **WHEN** `TENTACLE_CACHE_WARM_INTERVAL=5m` is set
- **THEN** `config.CacheWarmInterval` is 5 minutes

#### Scenario: Custom cold interval as integer seconds
- **WHEN** `TENTACLE_CACHE_COLD_INTERVAL=3600` is set
- **THEN** `config.CacheColdInterval` is 1 hour

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
