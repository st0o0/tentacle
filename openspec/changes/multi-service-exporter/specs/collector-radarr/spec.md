## Purpose

Prometheus collectors exposing Radarr metrics including system status, health, movie library, download queue, disk usage, and upcoming calendar.

## ADDED Requirements

### Requirement: System up metric
`radarr_up` SHALL be a gauge: 1 if the Radarr API responds to `/api/v3/system/status`, 0 otherwise.

#### Scenario: Radarr reachable
- **WHEN** `/api/v3/system/status` returns HTTP 200
- **THEN** `radarr_up` is 1

#### Scenario: Radarr unreachable
- **WHEN** `/api/v3/system/status` fails
- **THEN** `radarr_up` is 0 and no other system metrics are emitted

### Requirement: System info metric
`radarr_system_info{version, branch, runtime}` SHALL be a gauge with constant value 1 exposing Radarr metadata.

#### Scenario: System info
- **WHEN** Radarr returns version "5.2.0", branch "main", runtime "dotnet"
- **THEN** `radarr_system_info{version="5.2.0", branch="main", runtime="dotnet"}` is 1

### Requirement: Health issues metric
`radarr_health_issues_total{type, source}` SHALL be a gauge counting health check issues grouped by type and source.

#### Scenario: Health warnings
- **WHEN** Radarr reports 1 error from "DownloadClientCheck"
- **THEN** `radarr_health_issues_total{type="error", source="DownloadClientCheck"}` is 1

### Requirement: Movies total metric
`radarr_movies_total` SHALL be a gauge with the total number of movies.

#### Scenario: Movie count
- **WHEN** Radarr has 500 movies
- **THEN** `radarr_movies_total` is 500

### Requirement: Movies monitored metric
`radarr_movies_monitored_total` SHALL be a gauge with the number of monitored movies.

#### Scenario: Monitored movies
- **WHEN** 450 of 500 movies are monitored
- **THEN** `radarr_movies_monitored_total` is 450

### Requirement: Movies downloaded metric
`radarr_movies_downloaded_total` SHALL be a gauge with the number of movies that have files.

#### Scenario: Downloaded movies
- **WHEN** 400 of 500 movies have files
- **THEN** `radarr_movies_downloaded_total` is 400

### Requirement: Movies missing metric
`radarr_movies_missing_total` SHALL be a gauge with the number of missing monitored movies from the wanted/missing endpoint.

#### Scenario: Missing movies
- **WHEN** `/api/v3/wanted/missing` returns TotalRecords=50
- **THEN** `radarr_movies_missing_total` is 50

### Requirement: Queue total metric
`radarr_queue_total` SHALL be a gauge with the total number of items in the download queue.

#### Scenario: Queue items
- **WHEN** the download queue has 2 items
- **THEN** `radarr_queue_total` is 2

### Requirement: Disk space metrics
`radarr_disk_total_bytes{path}` and `radarr_disk_free_bytes{path}` SHALL be gauges with disk space per root folder.

#### Scenario: Disk usage
- **WHEN** root folder `/movies` has 4TB total and 1TB free
- **THEN** `radarr_disk_total_bytes{path="/movies"}` is 4398046511104 and `radarr_disk_free_bytes{path="/movies"}` is 1099511627776

### Requirement: Calendar upcoming metric
`radarr_calendar_upcoming_total` SHALL be a gauge with the number of movies releasing in the next 30 days.

#### Scenario: Upcoming movies
- **WHEN** 3 movies release in the next 30 days
- **THEN** `radarr_calendar_upcoming_total` is 3

### Requirement: Movies size metric
`radarr_movies_size_bytes` SHALL be a gauge with the total size on disk across all movies.

#### Scenario: Total library size
- **WHEN** all movies consume 3TB
- **THEN** `radarr_movies_size_bytes` is 3298534883328
