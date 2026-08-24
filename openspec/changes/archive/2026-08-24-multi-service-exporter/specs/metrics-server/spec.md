## MODIFIED Requirements

### Requirement: Scrape instrumentation
Every collector SHALL emit `<namespace>_scrape_duration_seconds{collector}` and `<namespace>_scrape_success{collector}` on each scrape, where `<namespace>` matches the service namespace (e.g., `jellyfin`, `sonarr`, `audiobookshelf`).

#### Scenario: Jellyfin collector scrape metrics
- **WHEN** the system collector for Jellyfin completes a scrape
- **THEN** `jellyfin_scrape_success{collector="system"}` and `jellyfin_scrape_duration_seconds{collector="system"}` are emitted

#### Scenario: Sonarr collector scrape metrics
- **WHEN** the system collector for Sonarr completes a scrape
- **THEN** `sonarr_scrape_success{collector="system"}` and `sonarr_scrape_duration_seconds{collector="system"}` are emitted

#### Scenario: One service failure does not affect others
- **WHEN** the Sonarr API is unreachable but Jellyfin responds normally
- **THEN** Sonarr collectors emit `scrape_success=0` while Jellyfin collectors emit their metrics normally
