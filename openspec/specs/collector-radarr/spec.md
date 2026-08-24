## Purpose

Prometheus collectors exposing Radarr metrics including system status, health, movie library, download queue, disk usage, and upcoming calendar.

## Requirements

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

### Requirement: Queue by state metric
`radarr_queue_by_state_total{state}` SHALL be a gauge counting queue items grouped by `TrackedDownloadState` (e.g., downloading, importPending, importing, failedPending, warning).

#### Scenario: Mixed queue states
- **WHEN** the queue has 1 downloading and 1 failedPending items
- **THEN** `radarr_queue_by_state_total{state="downloading"}` is 1 and `{state="failedPending"}` is 1

#### Scenario: Empty queue
- **WHEN** the queue has no items
- **THEN** no `radarr_queue_by_state_total` metrics are emitted

### Requirement: Backup metrics
`radarr_backup_total` SHALL be a gauge with the number of backups. `radarr_backup_latest_timestamp_seconds` SHALL be a gauge with the Unix timestamp of the most recent backup.

#### Scenario: Backups available
- **WHEN** Radarr has 5 backups, latest at 2024-06-15T08:00:00Z
- **THEN** `radarr_backup_total` is 5 and `radarr_backup_latest_timestamp_seconds` is the Unix timestamp

#### Scenario: No backups
- **WHEN** Radarr has no backups
- **THEN** `radarr_backup_total` is 0 and `radarr_backup_latest_timestamp_seconds` is not emitted

### Requirement: Update available metric
`radarr_update_available{version}` SHALL be a gauge: 1 if an update is available, 0 if up to date.

#### Scenario: Update available
- **WHEN** the latest update entry has Version "5.3.0" and Installed=false
- **THEN** `radarr_update_available{version="5.3.0"}` is 1

#### Scenario: Up to date
- **WHEN** the latest entry has Installed=true
- **THEN** `radarr_update_available{version="5.2.0"}` is 0

#### Scenario: No update data
- **WHEN** `/api/v3/update` returns an empty array
- **THEN** no `radarr_update_available` metric is emitted

### Requirement: Blocklist total metric
`radarr_blocklist_total` SHALL be a gauge with the total number of blocked releases.

#### Scenario: Blocked releases
- **WHEN** the blocklist has 8 entries
- **THEN** `radarr_blocklist_total` is 8

### Requirement: Download client info metric
`radarr_download_client_info{name, protocol, priority}` SHALL be a gauge with value 1 for each configured download client where Enable=true.

#### Scenario: Enabled download client
- **WHEN** client "qBittorrent" with protocol "torrent" and priority 1 is enabled
- **THEN** `radarr_download_client_info{name="qBittorrent", protocol="torrent", priority="1"}` is 1
