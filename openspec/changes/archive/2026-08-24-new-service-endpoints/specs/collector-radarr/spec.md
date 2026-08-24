## ADDED Requirements

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
