## Purpose

Prometheus collectors exposing Sonarr metrics including system status, health, series library, episode tracking, download queue, disk usage, and upcoming calendar.

## ADDED Requirements

### Requirement: System up metric
`sonarr_up` SHALL be a gauge: 1 if the Sonarr API responds to `/api/v3/system/status`, 0 otherwise.

#### Scenario: Sonarr reachable
- **WHEN** `/api/v3/system/status` returns HTTP 200
- **THEN** `sonarr_up` is 1

#### Scenario: Sonarr unreachable
- **WHEN** `/api/v3/system/status` fails
- **THEN** `sonarr_up` is 0 and no other system metrics are emitted

### Requirement: System info metric
`sonarr_system_info{version, branch, runtime}` SHALL be a gauge with constant value 1 exposing Sonarr metadata.

#### Scenario: System info
- **WHEN** Sonarr returns version "4.0.0", branch "main", runtime "dotnet"
- **THEN** `sonarr_system_info{version="4.0.0", branch="main", runtime="dotnet"}` is 1

### Requirement: Health issues metric
`sonarr_health_issues_total{type, source}` SHALL be a gauge counting health check issues grouped by type (warning, error) and source.

#### Scenario: Health warnings
- **WHEN** Sonarr reports 2 warnings from "IndexerStatusCheck"
- **THEN** `sonarr_health_issues_total{type="warning", source="IndexerStatusCheck"}` is 2

#### Scenario: No health issues
- **WHEN** Sonarr health endpoint returns an empty array
- **THEN** no `sonarr_health_issues_total` metrics are emitted

### Requirement: Series total metric
`sonarr_series_total` SHALL be a gauge with the total number of series.

#### Scenario: Series count
- **WHEN** Sonarr has 50 series
- **THEN** `sonarr_series_total` is 50

### Requirement: Series monitored metric
`sonarr_series_monitored_total` SHALL be a gauge with the number of monitored series.

#### Scenario: Monitored series
- **WHEN** 40 of 50 series are monitored
- **THEN** `sonarr_series_monitored_total` is 40

### Requirement: Episodes total metric
`sonarr_episodes_total` SHALL be a gauge with the total episode count across all series.

#### Scenario: Episode count
- **WHEN** all series have a combined 2000 total episodes
- **THEN** `sonarr_episodes_total` is 2000

### Requirement: Episodes downloaded metric
`sonarr_episodes_downloaded_total` SHALL be a gauge with the number of episodes that have files.

#### Scenario: Downloaded episodes
- **WHEN** 1800 of 2000 episodes have files
- **THEN** `sonarr_episodes_downloaded_total` is 1800

### Requirement: Episodes missing metric
`sonarr_episodes_missing_total` SHALL be a gauge with the number of missing monitored episodes from the wanted/missing endpoint.

#### Scenario: Missing episodes
- **WHEN** `/api/v3/wanted/missing` returns TotalRecords=12
- **THEN** `sonarr_episodes_missing_total` is 12

### Requirement: Queue total metric
`sonarr_queue_total` SHALL be a gauge with the total number of items in the download queue.

#### Scenario: Queue items
- **WHEN** the download queue has 3 items
- **THEN** `sonarr_queue_total` is 3

### Requirement: Disk space metrics
`sonarr_disk_total_bytes{path}` and `sonarr_disk_free_bytes{path}` SHALL be gauges with disk space per root folder.

#### Scenario: Disk usage
- **WHEN** root folder `/tv` has 2TB total and 500GB free
- **THEN** `sonarr_disk_total_bytes{path="/tv"}` is 2199023255552 and `sonarr_disk_free_bytes{path="/tv"}` is 536870912000

### Requirement: Calendar upcoming metric
`sonarr_calendar_upcoming_total` SHALL be a gauge with the number of episodes airing in the next 7 days.

#### Scenario: Upcoming episodes
- **WHEN** 5 episodes air in the next 7 days
- **THEN** `sonarr_calendar_upcoming_total` is 5

### Requirement: Series size metric
`sonarr_series_size_bytes` SHALL be a gauge with the total size on disk across all series.

#### Scenario: Total library size
- **WHEN** all series consume 1.5TB
- **THEN** `sonarr_series_size_bytes` is 1649267441664
