## MODIFIED Requirements

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
