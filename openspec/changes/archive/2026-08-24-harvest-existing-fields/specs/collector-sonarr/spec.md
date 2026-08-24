## ADDED Requirements

### Requirement: Start time metric
`sonarr_start_time_seconds` SHALL be a gauge with the Unix timestamp of when the Sonarr process started, parsed from `SystemStatus.StartTime`.

#### Scenario: Sonarr start time
- **WHEN** Sonarr returns StartTime "2024-06-15T08:00:00Z"
- **THEN** `sonarr_start_time_seconds` is the Unix timestamp of that date

#### Scenario: Uptime calculation
- **WHEN** `sonarr_start_time_seconds` is emitted
- **THEN** uptime can be calculated in PromQL as `time() - sonarr_start_time_seconds`

### Requirement: Series by status metric
`sonarr_series_by_status_total{status}` SHALL be a gauge counting series grouped by their Status field. Status values include: continuing, ended, upcoming, deleted.

#### Scenario: Mixed series statuses
- **WHEN** Sonarr has 30 continuing, 15 ended, and 5 upcoming series
- **THEN** `sonarr_series_by_status_total{status="continuing"}` is 30, `{status="ended"}` is 15, `{status="upcoming"}` is 5

#### Scenario: No series of a status
- **WHEN** no series have status "deleted"
- **THEN** `sonarr_series_by_status_total{status="deleted"}` is not emitted

### Requirement: Seasons total metric
`sonarr_seasons_total` SHALL be a gauge with the total season count aggregated across all series from `Series[].SeasonCount`.

#### Scenario: Season count
- **WHEN** Sonarr has 3 series with 5, 3, and 8 seasons respectively
- **THEN** `sonarr_seasons_total` is 16
