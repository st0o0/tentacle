## MODIFIED Requirements

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

## ADDED Requirements

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

### Requirement: Disable playback flag
`TENTACLE_DISABLE_PLAYBACK` SHALL remain as a Jellyfin-specific flag. It SHALL only be relevant when Jellyfin is configured.

#### Scenario: Disable playback without Jellyfin
- **WHEN** `TENTACLE_DISABLE_PLAYBACK=true` is set but Jellyfin is not configured
- **THEN** the flag is ignored and config loads successfully
