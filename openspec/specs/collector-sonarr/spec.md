## Purpose

Prometheus collectors exposing Sonarr metrics including system status, health, series library, episode tracking, download queue, disk usage, and upcoming calendar.
## Requirements
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
`sonarr_episodes_total` SHALL be a gauge with the total episode count across all series, read from `Series[].Statistics.TotalEpisodeCount`. If `Statistics.TotalEpisodeCount` is 0, the collector SHALL fall back to the top-level `TotalEpisodeCount`.

#### Scenario: Episode count
- **WHEN** all series have a combined 2000 total episodes
- **THEN** `sonarr_episodes_total` is 2000

#### Scenario: Episode count from statistics
- **WHEN** all series have a combined `statistics.totalEpisodeCount` of 2000
- **THEN** `sonarr_episodes_total` is 2000

#### Scenario: Fallback to top-level
- **WHEN** a series has `statistics.totalEpisodeCount=0` and top-level `totalEpisodeCount=50`
- **THEN** the top-level value 50 is used for that series

### Requirement: Episodes downloaded metric
`sonarr_episodes_downloaded_total` SHALL be a gauge with the number of episodes that have files, read from `Series[].Statistics.EpisodeFileCount`. If `Statistics.EpisodeFileCount` is 0, the collector SHALL fall back to the top-level `EpisodeFileCount`.

#### Scenario: Downloaded episodes
- **WHEN** 1800 of 2000 episodes have files
- **THEN** `sonarr_episodes_downloaded_total` is 1800

#### Scenario: Downloaded episodes from statistics
- **WHEN** 1800 of 2000 episodes have files per `statistics.episodeFileCount`
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
`sonarr_series_size_bytes` SHALL be a gauge with the total size on disk across all series, read from `Series[].Statistics.SizeOnDisk`. If `Statistics.SizeOnDisk` is 0, the collector SHALL fall back to the top-level `SizeOnDisk`.

#### Scenario: Total library size
- **WHEN** all series consume 1.5TB
- **THEN** `sonarr_series_size_bytes` is 1649267441664

#### Scenario: Total library size from statistics
- **WHEN** all series consume 1.5TB per `statistics.sizeOnDisk`
- **THEN** `sonarr_series_size_bytes` is 1649267441664

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
`sonarr_seasons_total` SHALL be a gauge with the total season count aggregated across all series from `Series[].Statistics.SeasonCount`. If `Statistics.SeasonCount` is 0, the collector SHALL fall back to the top-level `SeasonCount`.

#### Scenario: Season count
- **WHEN** Sonarr has 3 series with 5, 3, and 8 seasons respectively
- **THEN** `sonarr_seasons_total` is 16

#### Scenario: Season count from statistics
- **WHEN** Sonarr has 3 series with `statistics.seasonCount` of 5, 3, and 8 respectively
- **THEN** `sonarr_seasons_total` is 16

### Requirement: Queue by state metric
`sonarr_queue_by_state_total{state}` SHALL be a gauge counting queue items grouped by `TrackedDownloadState` (e.g., downloading, importPending, importing, failedPending, warning).

#### Scenario: Mixed queue states
- **WHEN** the queue has 2 downloading, 1 importing, and 1 failedPending items
- **THEN** `sonarr_queue_by_state_total{state="downloading"}` is 2, `{state="importing"}` is 1, `{state="failedPending"}` is 1

#### Scenario: Empty queue
- **WHEN** the queue has no items
- **THEN** no `sonarr_queue_by_state_total` metrics are emitted

### Requirement: Backup metrics
`sonarr_backup_total` SHALL be a gauge with the number of backups. `sonarr_backup_latest_timestamp_seconds` SHALL be a gauge with the Unix timestamp of the most recent backup.

#### Scenario: Backups available
- **WHEN** Sonarr has 3 backups, latest at 2024-06-15T08:00:00Z
- **THEN** `sonarr_backup_total` is 3 and `sonarr_backup_latest_timestamp_seconds` is the Unix timestamp

#### Scenario: No backups
- **WHEN** Sonarr has no backups
- **THEN** `sonarr_backup_total` is 0 and `sonarr_backup_latest_timestamp_seconds` is not emitted

### Requirement: Update available metric
`sonarr_update_available{version}` SHALL be a gauge: 1 if an update is available (latest entry has Installed=false), 0 if up to date. The `version` label SHALL contain the latest available version.

#### Scenario: Update available
- **WHEN** the latest update entry has Version "4.1.0" and Installed=false
- **THEN** `sonarr_update_available{version="4.1.0"}` is 1

#### Scenario: Up to date
- **WHEN** the latest entry has Installed=true, Version "4.0.0"
- **THEN** `sonarr_update_available{version="4.0.0"}` is 0

#### Scenario: No update data
- **WHEN** `/api/v3/update` returns an empty array
- **THEN** no `sonarr_update_available` metric is emitted

### Requirement: Blocklist total metric
`sonarr_blocklist_total` SHALL be a gauge with the total number of blocked releases.

#### Scenario: Blocked releases
- **WHEN** the blocklist has 15 entries
- **THEN** `sonarr_blocklist_total` is 15

### Requirement: Download client info metric
`sonarr_download_client_info{name, protocol, priority}` SHALL be a gauge with value 1 for each configured download client where Enable=true.

#### Scenario: Enabled download client
- **WHEN** client "SABnzbd" with protocol "usenet" and priority 1 is enabled
- **THEN** `sonarr_download_client_info{name="SABnzbd", protocol="usenet", priority="1"}` is 1

#### Scenario: Disabled download client
- **WHEN** a client has Enable=false
- **THEN** no info metric is emitted for that client

