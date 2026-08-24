## ADDED Requirements

### Requirement: Start time metric
`radarr_start_time_seconds` SHALL be a gauge with the Unix timestamp of when the Radarr process started, parsed from `SystemStatus.StartTime`.

#### Scenario: Radarr start time
- **WHEN** Radarr returns StartTime "2024-06-15T08:00:00Z"
- **THEN** `radarr_start_time_seconds` is the Unix timestamp of that date

#### Scenario: Uptime calculation
- **WHEN** `radarr_start_time_seconds` is emitted
- **THEN** uptime can be calculated in PromQL as `time() - radarr_start_time_seconds`

### Requirement: Movies by status metric
`radarr_movies_by_status_total{status}` SHALL be a gauge counting movies grouped by their Status field. Status values include: released, announced, inCinemas, deleted.

#### Scenario: Mixed movie statuses
- **WHEN** Radarr has 400 released, 50 announced, and 50 inCinemas movies
- **THEN** `radarr_movies_by_status_total{status="released"}` is 400, `{status="announced"}` is 50, `{status="inCinemas"}` is 50

#### Scenario: No movies of a status
- **WHEN** no movies have status "deleted"
- **THEN** `radarr_movies_by_status_total{status="deleted"}` is not emitted
